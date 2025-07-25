// Code generated by 'yaegi extract github.com/imroc/req/v3'. DO NOT EDIT.

package libs

import (
	"github.com/imroc/req/v3"
	"go/constant"
	"go/token"
	"reflect"
)

func init() {
	Symbols["github.com/imroc/req/v3/req"] = map[string]reflect.Value{
		// function, constant and variable definitions
		"AddCommonQueryParam":                      reflect.ValueOf(req.AddCommonQueryParam),
		"AddCommonQueryParams":                     reflect.ValueOf(req.AddCommonQueryParams),
		"AddCommonRetryCondition":                  reflect.ValueOf(req.AddCommonRetryCondition),
		"AddCommonRetryHook":                       reflect.ValueOf(req.AddCommonRetryHook),
		"AddQueryParam":                            reflect.ValueOf(req.AddQueryParam),
		"AddQueryParams":                           reflect.ValueOf(req.AddQueryParams),
		"AddRetryCondition":                        reflect.ValueOf(req.AddRetryCondition),
		"AddRetryHook":                             reflect.ValueOf(req.AddRetryHook),
		"AllowedDomainRedirectPolicy":              reflect.ValueOf(req.AllowedDomainRedirectPolicy),
		"AllowedHostRedirectPolicy":                reflect.ValueOf(req.AllowedHostRedirectPolicy),
		"AlwaysCopyHeaderRedirectPolicy":           reflect.ValueOf(req.AlwaysCopyHeaderRedirectPolicy),
		"C":                                        reflect.ValueOf(req.C),
		"ClearCookies":                             reflect.ValueOf(req.ClearCookies),
		"DefaultClient":                            reflect.ValueOf(req.DefaultClient),
		"DefaultRedirectPolicy":                    reflect.ValueOf(req.DefaultRedirectPolicy),
		"Delete":                                   reflect.ValueOf(req.Delete),
		"DevMode":                                  reflect.ValueOf(req.DevMode),
		"DisableAllowGetMethodPayload":             reflect.ValueOf(req.DisableAllowGetMethodPayload),
		"DisableAutoDecode":                        reflect.ValueOf(req.DisableAutoDecode),
		"DisableAutoReadResponse":                  reflect.ValueOf(req.DisableAutoReadResponse),
		"DisableCompression":                       reflect.ValueOf(req.DisableCompression),
		"DisableDebugLog":                          reflect.ValueOf(req.DisableDebugLog),
		"DisableDumpAll":                           reflect.ValueOf(req.DisableDumpAll),
		"DisableForceChunkedEncoding":              reflect.ValueOf(req.DisableForceChunkedEncoding),
		"DisableForceHttpVersion":                  reflect.ValueOf(req.DisableForceHttpVersion),
		"DisableForceMultipart":                    reflect.ValueOf(req.DisableForceMultipart),
		"DisableH2C":                               reflect.ValueOf(req.DisableH2C),
		"DisableInsecureSkipVerify":                reflect.ValueOf(req.DisableInsecureSkipVerify),
		"DisableKeepAlives":                        reflect.ValueOf(req.DisableKeepAlives),
		"DisableTrace":                             reflect.ValueOf(req.DisableTrace),
		"DisableTraceAll":                          reflect.ValueOf(req.DisableTraceAll),
		"EnableAllowGetMethodPayload":              reflect.ValueOf(req.EnableAllowGetMethodPayload),
		"EnableAutoDecode":                         reflect.ValueOf(req.EnableAutoDecode),
		"EnableAutoReadResponse":                   reflect.ValueOf(req.EnableAutoReadResponse),
		"EnableCloseConnection":                    reflect.ValueOf(req.EnableCloseConnection),
		"EnableCompression":                        reflect.ValueOf(req.EnableCompression),
		"EnableDebugLog":                           reflect.ValueOf(req.EnableDebugLog),
		"EnableDump":                               reflect.ValueOf(req.EnableDump),
		"EnableDumpAll":                            reflect.ValueOf(req.EnableDumpAll),
		"EnableDumpAllAsync":                       reflect.ValueOf(req.EnableDumpAllAsync),
		"EnableDumpAllTo":                          reflect.ValueOf(req.EnableDumpAllTo),
		"EnableDumpAllToFile":                      reflect.ValueOf(req.EnableDumpAllToFile),
		"EnableDumpAllWithoutBody":                 reflect.ValueOf(req.EnableDumpAllWithoutBody),
		"EnableDumpAllWithoutHeader":               reflect.ValueOf(req.EnableDumpAllWithoutHeader),
		"EnableDumpAllWithoutRequest":              reflect.ValueOf(req.EnableDumpAllWithoutRequest),
		"EnableDumpAllWithoutRequestBody":          reflect.ValueOf(req.EnableDumpAllWithoutRequestBody),
		"EnableDumpAllWithoutResponse":             reflect.ValueOf(req.EnableDumpAllWithoutResponse),
		"EnableDumpAllWithoutResponseBody":         reflect.ValueOf(req.EnableDumpAllWithoutResponseBody),
		"EnableDumpEachRequest":                    reflect.ValueOf(req.EnableDumpEachRequest),
		"EnableDumpEachRequestWithoutBody":         reflect.ValueOf(req.EnableDumpEachRequestWithoutBody),
		"EnableDumpEachRequestWithoutHeader":       reflect.ValueOf(req.EnableDumpEachRequestWithoutHeader),
		"EnableDumpEachRequestWithoutRequest":      reflect.ValueOf(req.EnableDumpEachRequestWithoutRequest),
		"EnableDumpEachRequestWithoutRequestBody":  reflect.ValueOf(req.EnableDumpEachRequestWithoutRequestBody),
		"EnableDumpEachRequestWithoutResponse":     reflect.ValueOf(req.EnableDumpEachRequestWithoutResponse),
		"EnableDumpEachRequestWithoutResponseBody": reflect.ValueOf(req.EnableDumpEachRequestWithoutResponseBody),
		"EnableDumpTo":                             reflect.ValueOf(req.EnableDumpTo),
		"EnableDumpToFile":                         reflect.ValueOf(req.EnableDumpToFile),
		"EnableDumpWithoutBody":                    reflect.ValueOf(req.EnableDumpWithoutBody),
		"EnableDumpWithoutHeader":                  reflect.ValueOf(req.EnableDumpWithoutHeader),
		"EnableDumpWithoutRequest":                 reflect.ValueOf(req.EnableDumpWithoutRequest),
		"EnableDumpWithoutRequestBody":             reflect.ValueOf(req.EnableDumpWithoutRequestBody),
		"EnableDumpWithoutResponse":                reflect.ValueOf(req.EnableDumpWithoutResponse),
		"EnableDumpWithoutResponseBody":            reflect.ValueOf(req.EnableDumpWithoutResponseBody),
		"EnableForceChunkedEncoding":               reflect.ValueOf(req.EnableForceChunkedEncoding),
		"EnableForceHTTP1":                         reflect.ValueOf(req.EnableForceHTTP1),
		"EnableForceHTTP2":                         reflect.ValueOf(req.EnableForceHTTP2),
		"EnableForceHTTP3":                         reflect.ValueOf(req.EnableForceHTTP3),
		"EnableForceMultipart":                     reflect.ValueOf(req.EnableForceMultipart),
		"EnableH2C":                                reflect.ValueOf(req.EnableH2C),
		"EnableHTTP3":                              reflect.ValueOf(req.EnableHTTP3),
		"EnableInsecureSkipVerify":                 reflect.ValueOf(req.EnableInsecureSkipVerify),
		"EnableKeepAlives":                         reflect.ValueOf(req.EnableKeepAlives),
		"EnableTrace":                              reflect.ValueOf(req.EnableTrace),
		"EnableTraceAll":                           reflect.ValueOf(req.EnableTraceAll),
		"ErrorState":                               reflect.ValueOf(req.ErrorState),
		"Get":                                      reflect.ValueOf(req.Get),
		"GetClient":                                reflect.ValueOf(req.GetClient),
		"GetCookies":                               reflect.ValueOf(req.GetCookies),
		"GetTLSClientConfig":                       reflect.ValueOf(req.GetTLSClientConfig),
		"Head":                                     reflect.ValueOf(req.Head),
		"HeaderOderKey":                            reflect.ValueOf(constant.MakeFromLiteral("\"__header_order__\"", token.STRING, 0)),
		"ImpersonateChrome":                        reflect.ValueOf(req.ImpersonateChrome),
		"ImpersonateFirefox":                       reflect.ValueOf(req.ImpersonateFirefox),
		"ImpersonateSafari":                        reflect.ValueOf(req.ImpersonateSafari),
		"MaxRedirectPolicy":                        reflect.ValueOf(req.MaxRedirectPolicy),
		"MustDelete":                               reflect.ValueOf(req.MustDelete),
		"MustGet":                                  reflect.ValueOf(req.MustGet),
		"MustHead":                                 reflect.ValueOf(req.MustHead),
		"MustOptions":                              reflect.ValueOf(req.MustOptions),
		"MustPatch":                                reflect.ValueOf(req.MustPatch),
		"MustPost":                                 reflect.ValueOf(req.MustPost),
		"MustPut":                                  reflect.ValueOf(req.MustPut),
		"NewClient":                                reflect.ValueOf(req.NewClient),
		"NewLogger":                                reflect.ValueOf(req.NewLogger),
		"NewLoggerFromStandardLogger":              reflect.ValueOf(req.NewLoggerFromStandardLogger),
		"NewRequest":                               reflect.ValueOf(req.NewRequest),
		"NewTransport":                             reflect.ValueOf(req.NewTransport),
		"NoBody":                                   reflect.ValueOf(&req.NoBody).Elem(),
		"NoRedirectPolicy":                         reflect.ValueOf(req.NoRedirectPolicy),
		"OnAfterResponse":                          reflect.ValueOf(req.OnAfterResponse),
		"OnBeforeRequest":                          reflect.ValueOf(req.OnBeforeRequest),
		"Options":                                  reflect.ValueOf(req.Options),
		"Patch":                                    reflect.ValueOf(req.Patch),
		"Post":                                     reflect.ValueOf(req.Post),
		"PseudoHeaderOderKey":                      reflect.ValueOf(constant.MakeFromLiteral("\"__pseudo_header_order__\"", token.STRING, 0)),
		"Put":                                      reflect.ValueOf(req.Put),
		"R":                                        reflect.ValueOf(req.R),
		"SameDomainRedirectPolicy":                 reflect.ValueOf(req.SameDomainRedirectPolicy),
		"SameHostRedirectPolicy":                   reflect.ValueOf(req.SameHostRedirectPolicy),
		"SetAutoDecodeAllContentType":              reflect.ValueOf(req.SetAutoDecodeAllContentType),
		"SetAutoDecodeContentType":                 reflect.ValueOf(req.SetAutoDecodeContentType),
		"SetAutoDecodeContentTypeFunc":             reflect.ValueOf(req.SetAutoDecodeContentTypeFunc),
		"SetBaseURL":                               reflect.ValueOf(req.SetBaseURL),
		"SetBasicAuth":                             reflect.ValueOf(req.SetBasicAuth),
		"SetBearerAuthToken":                       reflect.ValueOf(req.SetBearerAuthToken),
		"SetBody":                                  reflect.ValueOf(req.SetBody),
		"SetBodyBytes":                             reflect.ValueOf(req.SetBodyBytes),
		"SetBodyJsonBytes":                         reflect.ValueOf(req.SetBodyJsonBytes),
		"SetBodyJsonMarshal":                       reflect.ValueOf(req.SetBodyJsonMarshal),
		"SetBodyJsonString":                        reflect.ValueOf(req.SetBodyJsonString),
		"SetBodyString":                            reflect.ValueOf(req.SetBodyString),
		"SetBodyXmlBytes":                          reflect.ValueOf(req.SetBodyXmlBytes),
		"SetBodyXmlMarshal":                        reflect.ValueOf(req.SetBodyXmlMarshal),
		"SetBodyXmlString":                         reflect.ValueOf(req.SetBodyXmlString),
		"SetCertFromFile":                          reflect.ValueOf(req.SetCertFromFile),
		"SetCerts":                                 reflect.ValueOf(req.SetCerts),
		"SetCommonBasicAuth":                       reflect.ValueOf(req.SetCommonBasicAuth),
		"SetCommonBearerAuthToken":                 reflect.ValueOf(req.SetCommonBearerAuthToken),
		"SetCommonContentType":                     reflect.ValueOf(req.SetCommonContentType),
		"SetCommonCookies":                         reflect.ValueOf(req.SetCommonCookies),
		"SetCommonDigestAuth":                      reflect.ValueOf(req.SetCommonDigestAuth),
		"SetCommonDumpOptions":                     reflect.ValueOf(req.SetCommonDumpOptions),
		"SetCommonError":                           reflect.ValueOf(req.SetCommonError),
		"SetCommonErrorResult":                     reflect.ValueOf(req.SetCommonErrorResult),
		"SetCommonFormData":                        reflect.ValueOf(req.SetCommonFormData),
		"SetCommonFormDataFromValues":              reflect.ValueOf(req.SetCommonFormDataFromValues),
		"SetCommonHeader":                          reflect.ValueOf(req.SetCommonHeader),
		"SetCommonHeaderOrder":                     reflect.ValueOf(req.SetCommonHeaderOrder),
		"SetCommonHeaders":                         reflect.ValueOf(req.SetCommonHeaders),
		"SetCommonPathParam":                       reflect.ValueOf(req.SetCommonPathParam),
		"SetCommonPathParams":                      reflect.ValueOf(req.SetCommonPathParams),
		"SetCommonPseudoHeaderOder":                reflect.ValueOf(req.SetCommonPseudoHeaderOder),
		"SetCommonQueryParam":                      reflect.ValueOf(req.SetCommonQueryParam),
		"SetCommonQueryParams":                     reflect.ValueOf(req.SetCommonQueryParams),
		"SetCommonQueryString":                     reflect.ValueOf(req.SetCommonQueryString),
		"SetCommonRetryBackoffInterval":            reflect.ValueOf(req.SetCommonRetryBackoffInterval),
		"SetCommonRetryCondition":                  reflect.ValueOf(req.SetCommonRetryCondition),
		"SetCommonRetryCount":                      reflect.ValueOf(req.SetCommonRetryCount),
		"SetCommonRetryFixedInterval":              reflect.ValueOf(req.SetCommonRetryFixedInterval),
		"SetCommonRetryHook":                       reflect.ValueOf(req.SetCommonRetryHook),
		"SetCommonRetryInterval":                   reflect.ValueOf(req.SetCommonRetryInterval),
		"SetContentType":                           reflect.ValueOf(req.SetContentType),
		"SetContext":                               reflect.ValueOf(req.SetContext),
		"SetCookieJar":                             reflect.ValueOf(req.SetCookieJar),
		"SetCookies":                               reflect.ValueOf(req.SetCookies),
		"SetDefaultClient":                         reflect.ValueOf(req.SetDefaultClient),
		"SetDial":                                  reflect.ValueOf(req.SetDial),
		"SetDialTLS":                               reflect.ValueOf(req.SetDialTLS),
		"SetDigestAuth":                            reflect.ValueOf(req.SetDigestAuth),
		"SetDownloadCallback":                      reflect.ValueOf(req.SetDownloadCallback),
		"SetDownloadCallbackWithInterval":          reflect.ValueOf(req.SetDownloadCallbackWithInterval),
		"SetDumpOptions":                           reflect.ValueOf(req.SetDumpOptions),
		"SetError":                                 reflect.ValueOf(req.SetError),
		"SetErrorResult":                           reflect.ValueOf(req.SetErrorResult),
		"SetFile":                                  reflect.ValueOf(req.SetFile),
		"SetFileBytes":                             reflect.ValueOf(req.SetFileBytes),
		"SetFileReader":                            reflect.ValueOf(req.SetFileReader),
		"SetFileUpload":                            reflect.ValueOf(req.SetFileUpload),
		"SetFiles":                                 reflect.ValueOf(req.SetFiles),
		"SetFormData":                              reflect.ValueOf(req.SetFormData),
		"SetFormDataAnyType":                       reflect.ValueOf(req.SetFormDataAnyType),
		"SetFormDataFromValues":                    reflect.ValueOf(req.SetFormDataFromValues),
		"SetHTTP2ConnectionFlow":                   reflect.ValueOf(req.SetHTTP2ConnectionFlow),
		"SetHTTP2HeaderPriority":                   reflect.ValueOf(req.SetHTTP2HeaderPriority),
		"SetHTTP2MaxHeaderListSize":                reflect.ValueOf(req.SetHTTP2MaxHeaderListSize),
		"SetHTTP2PingTimeout":                      reflect.ValueOf(req.SetHTTP2PingTimeout),
		"SetHTTP2PriorityFrames":                   reflect.ValueOf(req.SetHTTP2PriorityFrames),
		"SetHTTP2ReadIdleTimeout":                  reflect.ValueOf(req.SetHTTP2ReadIdleTimeout),
		"SetHTTP2SettingsFrame":                    reflect.ValueOf(req.SetHTTP2SettingsFrame),
		"SetHTTP2StrictMaxConcurrentStreams":       reflect.ValueOf(req.SetHTTP2StrictMaxConcurrentStreams),
		"SetHTTP2WriteByteTimeout":                 reflect.ValueOf(req.SetHTTP2WriteByteTimeout),
		"SetHeader":                                reflect.ValueOf(req.SetHeader),
		"SetHeaderOrder":                           reflect.ValueOf(req.SetHeaderOrder),
		"SetHeaders":                               reflect.ValueOf(req.SetHeaders),
		"SetJsonMarshal":                           reflect.ValueOf(req.SetJsonMarshal),
		"SetJsonUnmarshal":                         reflect.ValueOf(req.SetJsonUnmarshal),
		"SetLogger":                                reflect.ValueOf(req.SetLogger),
		"SetMultipartBoundaryFunc":                 reflect.ValueOf(req.SetMultipartBoundaryFunc),
		"SetOrderedFormData":                       reflect.ValueOf(req.SetOrderedFormData),
		"SetOutput":                                reflect.ValueOf(req.SetOutput),
		"SetOutputDirectory":                       reflect.ValueOf(req.SetOutputDirectory),
		"SetOutputFile":                            reflect.ValueOf(req.SetOutputFile),
		"SetPathParam":                             reflect.ValueOf(req.SetPathParam),
		"SetPathParams":                            reflect.ValueOf(req.SetPathParams),
		"SetProxy":                                 reflect.ValueOf(req.SetProxy),
		"SetProxyURL":                              reflect.ValueOf(req.SetProxyURL),
		"SetPseudoHeaderOrder":                     reflect.ValueOf(req.SetPseudoHeaderOrder),
		"SetQueryParam":                            reflect.ValueOf(req.SetQueryParam),
		"SetQueryParams":                           reflect.ValueOf(req.SetQueryParams),
		"SetQueryParamsAnyType":                    reflect.ValueOf(req.SetQueryParamsAnyType),
		"SetQueryString":                           reflect.ValueOf(req.SetQueryString),
		"SetRedirectPolicy":                        reflect.ValueOf(req.SetRedirectPolicy),
		"SetResponseBodyTransformer":               reflect.ValueOf(req.SetResponseBodyTransformer),
		"SetResult":                                reflect.ValueOf(req.SetResult),
		"SetResultStateCheckFunc":                  reflect.ValueOf(req.SetResultStateCheckFunc),
		"SetRetryBackoffInterval":                  reflect.ValueOf(req.SetRetryBackoffInterval),
		"SetRetryCondition":                        reflect.ValueOf(req.SetRetryCondition),
		"SetRetryCount":                            reflect.ValueOf(req.SetRetryCount),
		"SetRetryFixedInterval":                    reflect.ValueOf(req.SetRetryFixedInterval),
		"SetRetryHook":                             reflect.ValueOf(req.SetRetryHook),
		"SetRetryInterval":                         reflect.ValueOf(req.SetRetryInterval),
		"SetRootCertFromString":                    reflect.ValueOf(req.SetRootCertFromString),
		"SetRootCertsFromFile":                     reflect.ValueOf(req.SetRootCertsFromFile),
		"SetScheme":                                reflect.ValueOf(req.SetScheme),
		"SetSuccessResult":                         reflect.ValueOf(req.SetSuccessResult),
		"SetTLSClientConfig":                       reflect.ValueOf(req.SetTLSClientConfig),
		"SetTLSFingerprint":                        reflect.ValueOf(req.SetTLSFingerprint),
		"SetTLSFingerprint360":                     reflect.ValueOf(req.SetTLSFingerprint360),
		"SetTLSFingerprintAndroid":                 reflect.ValueOf(req.SetTLSFingerprintAndroid),
		"SetTLSFingerprintChrome":                  reflect.ValueOf(req.SetTLSFingerprintChrome),
		"SetTLSFingerprintEdge":                    reflect.ValueOf(req.SetTLSFingerprintEdge),
		"SetTLSFingerprintFirefox":                 reflect.ValueOf(req.SetTLSFingerprintFirefox),
		"SetTLSFingerprintIOS":                     reflect.ValueOf(req.SetTLSFingerprintIOS),
		"SetTLSFingerprintQQ":                      reflect.ValueOf(req.SetTLSFingerprintQQ),
		"SetTLSFingerprintRandomized":              reflect.ValueOf(req.SetTLSFingerprintRandomized),
		"SetTLSFingerprintSafari":                  reflect.ValueOf(req.SetTLSFingerprintSafari),
		"SetTLSHandshakeTimeout":                   reflect.ValueOf(req.SetTLSHandshakeTimeout),
		"SetTimeout":                               reflect.ValueOf(req.SetTimeout),
		"SetURL":                                   reflect.ValueOf(req.SetURL),
		"SetUnixSocket":                            reflect.ValueOf(req.SetUnixSocket),
		"SetUploadCallback":                        reflect.ValueOf(req.SetUploadCallback),
		"SetUploadCallbackWithInterval":            reflect.ValueOf(req.SetUploadCallbackWithInterval),
		"SetUserAgent":                             reflect.ValueOf(req.SetUserAgent),
		"SetXmlMarshal":                            reflect.ValueOf(req.SetXmlMarshal),
		"SetXmlUnmarshal":                          reflect.ValueOf(req.SetXmlUnmarshal),
		"SuccessState":                             reflect.ValueOf(req.SuccessState),
		"T":                                        reflect.ValueOf(req.T),
		"UnknownState":                             reflect.ValueOf(req.UnknownState),
		"WrapRoundTrip":                            reflect.ValueOf(req.WrapRoundTrip),
		"WrapRoundTripFunc":                        reflect.ValueOf(req.WrapRoundTripFunc),

		// type definitions
		"Client":                   reflect.ValueOf((*req.Client)(nil)),
		"ContentDisposition":       reflect.ValueOf((*req.ContentDisposition)(nil)),
		"DownloadCallback":         reflect.ValueOf((*req.DownloadCallback)(nil)),
		"DownloadInfo":             reflect.ValueOf((*req.DownloadInfo)(nil)),
		"DumpOptions":              reflect.ValueOf((*req.DumpOptions)(nil)),
		"ErrorHook":                reflect.ValueOf((*req.ErrorHook)(nil)),
		"FileUpload":               reflect.ValueOf((*req.FileUpload)(nil)),
		"GetContentFunc":           reflect.ValueOf((*req.GetContentFunc)(nil)),
		"GetRetryIntervalFunc":     reflect.ValueOf((*req.GetRetryIntervalFunc)(nil)),
		"HttpRoundTripFunc":        reflect.ValueOf((*req.HttpRoundTripFunc)(nil)),
		"HttpRoundTripWrapper":     reflect.ValueOf((*req.HttpRoundTripWrapper)(nil)),
		"HttpRoundTripWrapperFunc": reflect.ValueOf((*req.HttpRoundTripWrapperFunc)(nil)),
		"Logger":                   reflect.ValueOf((*req.Logger)(nil)),
		"ParallelDownload":         reflect.ValueOf((*req.ParallelDownload)(nil)),
		"RedirectPolicy":           reflect.ValueOf((*req.RedirectPolicy)(nil)),
		"Request":                  reflect.ValueOf((*req.Request)(nil)),
		"RequestMiddleware":        reflect.ValueOf((*req.RequestMiddleware)(nil)),
		"Response":                 reflect.ValueOf((*req.Response)(nil)),
		"ResponseMiddleware":       reflect.ValueOf((*req.ResponseMiddleware)(nil)),
		"ResultState":              reflect.ValueOf((*req.ResultState)(nil)),
		"RetryConditionFunc":       reflect.ValueOf((*req.RetryConditionFunc)(nil)),
		"RetryHookFunc":            reflect.ValueOf((*req.RetryHookFunc)(nil)),
		"RoundTripFunc":            reflect.ValueOf((*req.RoundTripFunc)(nil)),
		"RoundTripWrapper":         reflect.ValueOf((*req.RoundTripWrapper)(nil)),
		"RoundTripWrapperFunc":     reflect.ValueOf((*req.RoundTripWrapperFunc)(nil)),
		"RoundTripper":             reflect.ValueOf((*req.RoundTripper)(nil)),
		"TraceInfo":                reflect.ValueOf((*req.TraceInfo)(nil)),
		"Transport":                reflect.ValueOf((*req.Transport)(nil)),
		"UploadCallback":           reflect.ValueOf((*req.UploadCallback)(nil)),
		"UploadInfo":               reflect.ValueOf((*req.UploadInfo)(nil)),

		// interface wrapper definitions
		"_Logger":       reflect.ValueOf((*_github_com_imroc_req_v3_Logger)(nil)),
		"_RoundTripper": reflect.ValueOf((*_github_com_imroc_req_v3_RoundTripper)(nil)),
	}
}

// _github_com_imroc_req_v3_Logger is an interface wrapper for Logger type
type _github_com_imroc_req_v3_Logger struct {
	IValue  interface{}
	WDebugf func(format string, v ...any)
	WErrorf func(format string, v ...any)
	WWarnf  func(format string, v ...any)
}

func (W _github_com_imroc_req_v3_Logger) Debugf(format string, v ...any) {
	W.WDebugf(format, v...)
}
func (W _github_com_imroc_req_v3_Logger) Errorf(format string, v ...any) {
	W.WErrorf(format, v...)
}
func (W _github_com_imroc_req_v3_Logger) Warnf(format string, v ...any) {
	W.WWarnf(format, v...)
}

// _github_com_imroc_req_v3_RoundTripper is an interface wrapper for RoundTripper type
type _github_com_imroc_req_v3_RoundTripper struct {
	IValue     interface{}
	WRoundTrip func(a0 *req.Request) (*req.Response, error)
}

func (W _github_com_imroc_req_v3_RoundTripper) RoundTrip(a0 *req.Request) (*req.Response, error) {
	return W.WRoundTrip(a0)
}
