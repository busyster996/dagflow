package types

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"strconv"
)

const (
	Version                  = "1.0.0"
	HeaderUploadOffset       = "Upload-Offset"
	HeaderUploadLength       = "Upload-Length"
	HeaderUploadDeferLength  = "Upload-Defer-Length"
	HeaderUploadMetadata     = "Upload-Metadata"
	HeaderUploadConcat       = "Upload-Concat"
	HeaderUploadChecksum     = "Upload-Checksum"
	HeaderContent            = "Content-Type"
	HeaderContentDisposition = "Content-Disposition"
	HeaderCacheControl       = "Cache-Control"
	HeaderLocation           = "Location"
	HeaderVersion            = "Tus-Version"
	HeaderResumable          = "Tus-Resumable"
	HeaderMaxSize            = "Tus-Max-Size"
	HeaderExtension          = "Tus-Extension"
	HeaderChecksumAlgorithm  = "Tus-Checksum-Algorithm"
)

type FileInfoChanges struct {
	ID       string
	MetaData map[string]string
}

type FileInfo struct {
	ID             string            `json:"id"`
	Size           int64             `json:"size,omitempty"`
	SizeIsDeferred bool              `json:"sizeIsDeferred"`
	Offset         int64             `json:"offset"`
	MetaData       map[string]string `json:"metaData"`
	IsPartial      bool              `json:"isPartial"`
	IsFinal        bool              `json:"isFinal"`
	PartialIDs     []string          `json:"partialIDs,omitempty"`
}

type HookEvent struct {
	Context     context.Context
	Upload      FileInfo
	HTTPRequest *http.Request
}

type HTTPResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

func (resp HTTPResponse) WriteTo(w http.ResponseWriter) {
	headers := w.Header()
	for key, value := range resp.Headers {
		headers.Set(key, value)
	}

	if len(resp.Body) > 0 {
		headers.Set("Content-Length", strconv.Itoa(len(resp.Body)))
	}

	w.WriteHeader(resp.StatusCode)

	if len(resp.Body) > 0 {
		_, _ = w.Write([]byte(resp.Body))
	}
}

func (resp HTTPResponse) MergeWith(resp2 HTTPResponse) HTTPResponse {
	// Clone the response 1 and use it as a basis
	newResp := resp

	if resp2.StatusCode != 0 {
		newResp.StatusCode = resp2.StatusCode
	}

	if len(resp2.Body) > 0 {
		newResp.Body = resp2.Body
	}

	newResp.Headers = make(map[string]string, len(resp.Headers)+len(resp2.Headers))
	for key, value := range resp.Headers {
		newResp.Headers[key] = value
	}

	for key, value := range resp2.Headers {
		newResp.Headers[key] = value
	}
	return newResp
}

// Uid returns a unique id. These ids consist of 128 bits from a
// cryptographically strong pseudo-random generator and are like uuids, but
// without the dashes and significant bits.
//
// See: http://en.wikipedia.org/wiki/UUID#Random_UUID_probability_of_duplicates
func Uid() string {
	id := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		// This is probably an appropriate way to handle errors from our source
		// for random bits.
		panic(err)
	}
	return hex.EncodeToString(id)
}
