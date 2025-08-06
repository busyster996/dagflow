package tus

import (
	"context"
	"encoding/base64"
	"fmt"
	"mime"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/busyster996/dagflow/pkg/tus/storage"
	"github.com/busyster996/dagflow/pkg/tus/types"
)

var (
	reForwardedHost  = regexp.MustCompile(`host="?([^;"]+)`)
	reForwardedProto = regexp.MustCompile(`proto=(https?)`)
	reValidUploadId  = regexp.MustCompile(`^[A-Za-z0-9\-._~%!$'()*+,;=/:@]*$`)
)

type STusx struct {
	config        *SConfig
	basePath      string
	isBasePathAbs bool
	logger        ILogger
	storage       storage.IStorage
	events        *sMemoryBroker
	extensions    []string
	algorithms    []string
}

func New(config *SConfig) (*STusx, error) {
	if err := config.validate(); err != nil {
		return nil, err
	}
	return &STusx{
		config:        config,
		basePath:      config.BasePath,
		isBasePathAbs: config.isAbs,
		storage:       config.Store,
		logger:        config.Logger,
		events:        newMemoryBroker(config.Logger),
		extensions:    []string{"creation", "creation-with-upload", "checksum", "expiration", "termination", "concatenation"},
		algorithms:    []string{"sha1", "sha256", "sha512", "md5"},
	}, nil
}

func (s *STusx) Close(ctx context.Context) error {
	s.events.Shutdown(ctx)
	return nil
}

func (s *STusx) SubscribeCompleteUploads(ctx context.Context, callback func(hook types.HookEvent) error) {
	s.events.SubscribeEvent(ctx, "upload.finished", callback)
}

func (s *STusx) SubscribeTerminatedUploads(ctx context.Context, callback func(hook types.HookEvent) error) {
	s.events.SubscribeEvent(ctx, "upload.terminated", callback)
}

func (s *STusx) SubscribeCreatedUploads(ctx context.Context, callback func(hook types.HookEvent) error) {
	s.events.SubscribeEvent(ctx, "upload.created", callback)
}

func (s *STusx) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.setCommonHeaders(w, r)
	if r.Header.Get("X-HTTP-Method-Override") != "" {
		r.Method = r.Header.Get("X-HTTP-Method-Override")
	}

	if r.Method == http.MethodOptions {
		s.handleOptions(w, r)
		return
	}
	tusResumable := r.Header.Get(types.HeaderResumable)
	if tusResumable != types.Version && r.Method != http.MethodGet {
		w.Header().Set(types.HeaderVersion, types.Version)
		http.Error(w, "Unsupported version", http.StatusPreconditionFailed)
		return
	}

	if r.URL.Path == s.basePath {
		if r.Method == http.MethodPost {
			s.handlePost(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else if strings.HasPrefix(r.URL.Path, s.basePath) {
		uploadID := strings.TrimPrefix(r.URL.Path, s.basePath)
		switch r.Method {
		case http.MethodHead:
			s.handleHead(w, r, uploadID)
		case http.MethodPatch:
			s.handlePatch(w, r, uploadID)
		case http.MethodDelete:
			s.handleDelete(w, r, uploadID)
		case http.MethodGet:
			s.handleGet(w, r, uploadID)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	} else {
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (s *STusx) handlePost(w http.ResponseWriter, r *http.Request) {
	info, err := s.parseUploadInfo(r)
	if err != nil {
		s.logger.Errorf("Error parsing upload info: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if info.IsFinal && r.ContentLength != 0 {
		s.logger.Errorf("Final uploads cannot have a body")
		http.Error(w, "Final uploads cannot have a body", http.StatusBadRequest)
		return
	}

	if s.config.MaxSize > 0 && info.Size > s.config.MaxSize {
		s.logger.Errorf("Upload size exceeds maximum allowed: %v", s.config.MaxSize)
		http.Error(w, "Request Entity Too Large", http.StatusRequestEntityTooLarge)
		return
	}

	resp := types.HTTPResponse{
		StatusCode: http.StatusCreated,
		Headers:    make(map[string]string),
	}
	if s.config.PreUploadCreateCallback != nil {
		var resp2 types.HTTPResponse
		var changes types.FileInfoChanges
		resp2, changes, err = s.config.PreUploadCreateCallback(types.HookEvent{
			Context:     r.Context(),
			HTTPRequest: r,
			Upload:      info,
		})
		if err != nil {
			s.logger.Errorf("failed to run PreUploadCreateCallback: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp = resp.MergeWith(resp2)
		if changes.ID != "" {
			if err = s.validateUploadId(changes.ID); err != nil {
				s.logger.Errorf("failed to validate upload ID: %v", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			info.ID = changes.ID
		}

		info.MetaData = s.mergeMetadata(info.MetaData, changes.MetaData)
	}

	upload, err := s.storage.NewUpload(r.Context(), info)
	if err != nil {
		s.logger.Errorf("Error creating upload: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	info, err = upload.GetInfo(r.Context())
	if err != nil {
		s.logger.Errorf("Error getting upload inf: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(types.HeaderLocation, s.absFileURL(r, info.ID))
	s.events.PublishEvent("upload.created", types.HookEvent{
		Context:     r.Context(),
		HTTPRequest: r,
		Upload:      info,
	})

	// 处理Creation With Upload
	if r.ContentLength > 0 {
		contentType := r.Header.Get(types.HeaderContent)
		if contentType != "application/offset+octet-stream" {
			s.logger.Errorf("Unsupported Media Type: %v", contentType)
			http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
			return
		}
		var written int64
		written, err = s.wrapWithChecksum(r, upload, 0)
		if err != nil {
			s.logger.Errorf("Error parsing upload info: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set(types.HeaderUploadOffset, strconv.FormatInt(written, 10))
	}

	if info.IsFinal {
		var partialUploads []storage.IUpload
		for _, partialID := range info.PartialIDs {
			var partialUpload storage.IUpload
			partialUpload, err = s.storage.GetUpload(r.Context(), partialID)
			if err != nil {
				s.logger.Errorf("Error getting partial upload: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			partialUploads = append(partialUploads, partialUpload)
		}
		err = upload.ConcatUploads(r.Context(), partialUploads)
		if err != nil {
			s.logger.Errorf("Error concatenating uploads: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if s.config.PreFinishResponseCallback != nil {
			var resp2 types.HTTPResponse
			resp2, err = s.config.PreFinishResponseCallback(types.HookEvent{
				Context:     r.Context(),
				HTTPRequest: r,
				Upload:      info,
			})
			resp = resp.MergeWith(resp2)
		}
		s.events.PublishEvent("upload.finished", types.HookEvent{
			Context:     r.Context(),
			HTTPRequest: r,
			Upload:      info,
		})
		resp.WriteTo(w)
	}
}

func (s *STusx) handleHead(w http.ResponseWriter, r *http.Request, uploadID string) {
	upload, err := s.storage.GetUpload(r.Context(), uploadID)
	if err != nil {
		s.logger.Errorf("Error getting upload: %v", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	info, err := upload.GetInfo(r.Context())
	if err != nil {
		s.logger.Errorf("Error getting upload info: %v", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set(types.HeaderUploadOffset, strconv.FormatInt(info.Offset, 10))
	w.Header().Set(types.HeaderUploadLength, strconv.FormatInt(info.Size, 10))

	if len(info.MetaData) > 0 {
		metadata := s.encodeMetadata(info.MetaData)
		w.Header().Set(types.HeaderUploadMetadata, metadata)
	}

	if info.IsPartial {
		w.Header().Set(types.HeaderUploadConcat, "partial")
	} else if info.IsFinal {
		concat := "final"
		for _, partialID := range info.PartialIDs {
			concat += ";" + s.basePath + partialID
		}
		w.Header().Set(types.HeaderUploadConcat, concat)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *STusx) handlePatch(w http.ResponseWriter, r *http.Request, uploadID string) {
	contentType := r.Header.Get(types.HeaderContent)
	if contentType != "application/offset+octet-stream" {
		s.logger.Errorf("UnsupportedMedia Type: %v", contentType)
		http.Error(w, "Unsupported Media Type", http.StatusUnsupportedMediaType)
		return
	}

	upload, err := s.storage.GetUpload(r.Context(), uploadID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			s.logger.Errorf("Error getting upload: %v", err)
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			s.logger.Errorf("Error getting upload: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	info, err := upload.GetInfo(r.Context())
	if err != nil {
		s.logger.Errorf("Error getting upload info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.IsFinal {
		s.logger.Errorf("Cannot patch final upload: %v", uploadID)
		http.Error(w, "Cannot patch final upload", http.StatusForbidden)
		return
	}

	offsetHeader := r.Header.Get(types.HeaderUploadOffset)
	offset, err := strconv.ParseInt(offsetHeader, 10, 64)
	if err != nil || offset < 0 {
		s.logger.Errorf("Invalid Upload-Offset header: %v", offsetHeader)
		http.Error(w, "Invalid Upload-Offset header", http.StatusBadRequest)
		return
	}

	if offset != info.Offset {
		s.logger.Errorf(fmt.Sprintf("Offset mismatch: %d != %d", offset, info.Offset))
		http.Error(w, "Offset mismatch", http.StatusConflict)
		return
	}

	var written int64
	written, err = s.wrapWithChecksum(r, upload, offset)
	if err != nil {
		s.logger.Errorf("Error writing chunk: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newOffset := offset + written
	resp := types.HTTPResponse{
		StatusCode: http.StatusNoContent,
		Headers: map[string]string{
			types.HeaderUploadOffset: strconv.FormatInt(newOffset, 10),
		},
	}

	if s.config.PreFinishResponseCallback != nil {
		var resp2 types.HTTPResponse
		resp2, err = s.config.PreFinishResponseCallback(types.HookEvent{
			Context:     r.Context(),
			HTTPRequest: r,
			Upload:      info,
		})
		resp = resp.MergeWith(resp2)
	}
	s.events.PublishEvent("upload.progress", types.HookEvent{
		Context:     r.Context(),
		HTTPRequest: r,
		Upload:      info,
	})
	resp.WriteTo(w)
}

func (s *STusx) handleDelete(w http.ResponseWriter, r *http.Request, uploadID string) {
	upload, err := s.storage.GetUpload(r.Context(), uploadID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			s.logger.Errorf("Error getting upload: %v", err)
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			s.logger.Errorf("Error getting upload: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	info, err := upload.GetInfo(r.Context())
	if err != nil {
		s.logger.Errorf("Error getting upload info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := types.HTTPResponse{
		StatusCode: http.StatusNoContent,
	}

	if s.config.PreUploadTerminateCallback != nil {
		var resp2 types.HTTPResponse
		resp2, err = s.config.PreUploadTerminateCallback(types.HookEvent{
			Context:     r.Context(),
			HTTPRequest: r,
			Upload:      info,
		})
		if err != nil {
			s.logger.Errorf("failed to emit finish events: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		resp = resp.MergeWith(resp2)
	}

	err = upload.Terminate(r.Context())
	if err != nil {
		s.logger.Errorf("Error terminating upload: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.events.PublishEvent("upload.terminated", types.HookEvent{
		Context:     r.Context(),
		HTTPRequest: r,
		Upload:      info,
	})
	resp.WriteTo(w)
}

func (s *STusx) handleGet(w http.ResponseWriter, r *http.Request, uploadID string) {
	upload, err := s.storage.GetUpload(r.Context(), uploadID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			s.logger.Errorf("Error getting upload: %v", err)
			http.Error(w, "Not found", http.StatusNotFound)
		} else {
			s.logger.Errorf("Error getting upload: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	info, err := upload.GetInfo(r.Context())
	if err != nil {
		s.logger.Errorf("Error getting upload info: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	contentType, contentDisposition := s.filterContentType(info)
	w.Header().Set(types.HeaderContent, contentType)
	w.Header().Set(types.HeaderContentDisposition, contentDisposition)
	w.Header().Set(types.HeaderUploadLength, strconv.FormatInt(info.Size, 10))
	err = upload.ServeContent(r.Context(), w, r)
}

func (s *STusx) setCommonHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(types.HeaderResumable, types.Version)
	w.Header().Set(types.HeaderCacheControl, "no-store")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, HEAD, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Upload-Length, Upload-Offset, Tus-Resumable, Upload-Metadata, Upload-Defer-Length, Upload-Concat, Upload-Checksum")
	w.Header().Set("Access-Control-Expose-Headers", "Upload-Offset, Location, Upload-Length, Tus-Version, Tus-Resumable, Tus-Max-Size, Tus-Extension, Upload-Metadata, Upload-Defer-Length, Upload-Concat, Upload-Checksum, Tus-Checksum-Algorithm")
}

func (s *STusx) handleOptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(types.HeaderVersion, types.Version)
	w.Header().Set(types.HeaderMaxSize, strconv.FormatInt(s.config.MaxSize, 10))
	w.Header().Set(types.HeaderExtension, strings.Join(s.extensions, ","))
	w.Header().Set(types.HeaderChecksumAlgorithm, strings.Join(s.algorithms, ","))
	w.WriteHeader(http.StatusNoContent)
}

func (s *STusx) wrapWithChecksum(r *http.Request, upload storage.IUpload, offset int64) (written int64, err error) {
	checksumHeader := r.Header.Get(types.HeaderUploadChecksum)
	if checksumHeader == "" {
		return upload.WriteChunk(r.Context(), offset, r.Body)
	}

	parts := strings.SplitN(checksumHeader, " ", 2)
	if len(parts) != 2 {
		s.logger.Errorf("Invalid checksum header format: %v", checksumHeader)
		return 0, fmt.Errorf("invalid checksum header format")
	}

	algorithm := parts[0]
	expectedChecksum := parts[1]

	// 检查算法支持
	supported := false
	for _, alg := range s.algorithms {
		if alg == algorithm {
			supported = true
			break
		}
	}
	if !supported {
		s.logger.Errorf("Algorithm not supported: %v", algorithm)
		return 0, fmt.Errorf("algorithm not supported %s", algorithm)
	}
	sumReader, err := NewShaSumReader(algorithm, r.Body)
	if err != nil {
		return 0, err
	}

	defer func() {
		calculatedSum := sumReader.ChecksumBase64()
		if calculatedSum != expectedChecksum {
			s.logger.Errorf("checksum mismatch: %v", expectedChecksum)
			err = fmt.Errorf("checksum verification failed: expected %s, got %s",
				expectedChecksum, calculatedSum)
		}
	}()

	return upload.WriteChunk(r.Context(), 0, r.Body)
}

func (s *STusx) parseUploadInfo(r *http.Request) (info types.FileInfo, err error) {
	info.IsPartial, info.IsFinal, info.PartialIDs, err = s.parseConcat(r.Header.Get("Upload-Concat"))
	if err != nil {
		s.logger.Errorf("Error parsing upload info: %v", err)
		return info, err
	}

	lengthHeader := r.Header.Get(types.HeaderUploadLength)
	deferLengthHeader := r.Header.Get(types.HeaderUploadDeferLength)
	if !info.IsFinal {
		if deferLengthHeader == "1" {
			info.SizeIsDeferred = true
		} else if lengthHeader != "" {
			info.Size, err = strconv.ParseInt(lengthHeader, 10, 64)
			if err != nil {
				s.logger.Errorf("Invalid Upload-Length header: %v", lengthHeader)
				return info, fmt.Errorf("invalid Upload-Length header")
			}
		} else {
			s.logger.Errorf("Missing Upload-Length or Upload-Defer-Length header")
			return info, fmt.Errorf("missing Upload-Length or Upload-Defer-Length header")
		}
	}

	metadataHeader := r.Header.Get(types.HeaderUploadMetadata)
	if metadataHeader != "" {
		info.MetaData, err = s.parseMetadata(metadataHeader)
		if err != nil {
			s.logger.Errorf("Error parsing upload info: %v", err)
			return info, err
		}
	}

	return info, nil
}

func (s *STusx) parseMetadata(header string) (map[string]string, error) {
	metadata := make(map[string]string)

	if header == "" {
		return metadata, nil
	}

	pairs := strings.Split(header, ",")
	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), " ", 2)
		if len(parts) < 1 {
			continue
		}

		key := parts[0]
		var value string

		if len(parts) == 2 {
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid base64 in metadata: %v", err)
			}
			value = string(decoded)
		}

		metadata[key] = value
	}

	return metadata, nil
}

func (s *STusx) encodeMetadata(metadata map[string]string) string {
	var pairs []string
	for key, value := range metadata {
		if value == "" {
			pairs = append(pairs, key)
		} else {
			encoded := base64.StdEncoding.EncodeToString([]byte(value))
			pairs = append(pairs, key+" "+encoded)
		}
	}
	return strings.Join(pairs, ",")
}

func (s *STusx) mergeMetadata(old, new map[string]string) map[string]string {
	for key, value := range new {
		old[key] = value
	}
	return old
}

func (s *STusx) parseConcat(header string) (isPartial bool, isFinal bool, partialUploads []string, err error) {
	if len(header) == 0 {
		return
	}

	if header == "partial" {
		isPartial = true
		return
	}

	l := len("final;")
	if strings.HasPrefix(header, "final;") && len(header) > l {
		isFinal = true

		list := strings.Split(header[l:], " ")
		for _, value := range list {
			value = strings.TrimSpace(value)
			if value == "" {
				continue
			}

			id, extractErr := s.extractIDFromURL(value, s.basePath)
			if extractErr != nil {
				err = extractErr
				return
			}

			partialUploads = append(partialUploads, id)
		}
	}

	// If no valid partial upload ids are extracted this is not a final upload.
	if len(partialUploads) == 0 {
		isFinal = false
		err = fmt.Errorf("invalid Upload-Concat header")
	}

	return
}

func (s *STusx) extractIDFromURL(url string, basePath string) (string, error) {
	_, id, ok := strings.Cut(url, basePath)
	if !ok {
		return "", fmt.Errorf("upload not found %s", url)
	}

	return strings.Trim(id, "/"), nil
}

func (s *STusx) validateUploadId(newId string) error {
	if newId == "" {
		return nil
	}
	if strings.HasPrefix(newId, "/") || strings.HasSuffix(newId, "/") {
		return fmt.Errorf("validation error in FileInfoChanges: ID must not begin or end with a forward slash (got: %s)", newId)
	}
	if !reValidUploadId.MatchString(newId) {
		return fmt.Errorf("validation error in FileInfoChanges: ID must contain only URL-safe character: %s (got: %s)", reValidUploadId.String(), newId)
	}
	return nil
}

var mimeInlineBrowserWhitelist = map[string]struct{}{
	"text/plain": {},

	"image/png":  {},
	"image/jpeg": {},
	"image/gif":  {},
	"image/bmp":  {},
	"image/webp": {},

	"audio/wave":     {},
	"audio/wav":      {},
	"audio/x-wav":    {},
	"audio/x-pn-wav": {},
	"audio/webm":     {},
	"audio/ogg":      {},

	"video/mp4":  {},
	"video/webm": {},
	"video/ogg":  {},

	"application/ogg": {},
}

func (s *STusx) filterContentType(info types.FileInfo) (contentType string, contentDisposition string) {
	filetype := info.MetaData["filetype"]

	if ft, _, err := mime.ParseMediaType(filetype); err == nil {
		// If the filetype from metadata is well-formed, we forward use this for the Content-Type header.
		// However, only allowlisted mime types	will be allowed to be shown inline in the browser
		contentType = filetype
		if _, isWhitelisted := mimeInlineBrowserWhitelist[ft]; isWhitelisted {
			contentDisposition = "inline"
		} else {
			contentDisposition = "attachment"
		}
	} else {
		// If the filetype from the metadata is not well-formed, we use a
		// default type and force the browser to download the content.
		contentType = "application/octet-stream"
		contentDisposition = "attachment"
	}

	// Add a filename to Content-Disposition if one is available in the metadata
	if filename, ok := info.MetaData["filename"]; ok {
		contentDisposition += ";filename=" + strconv.Quote(filename)
	}

	return contentType, contentDisposition
}

func (s *STusx) absFileURL(r *http.Request, id string) string {
	if s.isBasePathAbs {
		return s.basePath + id
	}

	// Read origin and protocol from request
	host, proto := s.getHostAndProtocol(r)

	url := proto + "://" + host + s.basePath + id

	return url
}

func (s *STusx) getHostAndProtocol(r *http.Request) (host, proto string) {
	if r.TLS != nil {
		proto = "https"
	} else {
		proto = "http"
	}

	host = r.Host
	if h := r.Header.Get("X-Forwarded-Host"); h != "" {
		host = h
	}

	if h := r.Header.Get("X-Forwarded-Proto"); h == "http" || h == "https" {
		proto = h
	}

	if h := r.Header.Get("Forwarded"); h != "" {
		if _r := reForwardedHost.FindStringSubmatch(h); len(_r) == 2 {
			host = _r[1]
		}

		if _r := reForwardedProto.FindStringSubmatch(h); len(_r) == 2 {
			proto = _r[1]
		}
	}

	// Remove default ports
	if proto == "http" {
		host = strings.TrimSuffix(host, ":80")
	}
	if proto == "https" {
		host = strings.TrimSuffix(host, ":443")
	}

	return
}
