package noob

type HTTPStatusCode uint

func (c HTTPStatusCode) Copy() *HTTPStatusCode {
	nc := c
	return &nc
}

// List gathered from net/http/status
const (
	StatusContinue           HTTPStatusCode = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols HTTPStatusCode = 101 // RFC 7231, 6.2.2
	StatusProcessing         HTTPStatusCode = 102 // RFC 2518, 10.1
	StatusEarlyHints         HTTPStatusCode = 103 // RFC 8297

	StatusOK                   HTTPStatusCode = 200 // RFC 7231, 6.3.1
	StatusCreated              HTTPStatusCode = 201 // RFC 7231, 6.3.2
	StatusAccepted             HTTPStatusCode = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo HTTPStatusCode = 203 // RFC 7231, 6.3.4
	StatusNoContent            HTTPStatusCode = 204 // RFC 7231, 6.3.5
	StatusResetContent         HTTPStatusCode = 205 // RFC 7231, 6.3.6
	StatusPartialContent       HTTPStatusCode = 206 // RFC 7233, 4.1
	StatusMultiStatus          HTTPStatusCode = 207 // RFC 4918, 11.1
	StatusAlreadyReported      HTTPStatusCode = 208 // RFC 5842, 7.1
	StatusIMUsed               HTTPStatusCode = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   HTTPStatusCode = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  HTTPStatusCode = 301 // RFC 7231, 6.4.2
	StatusFound             HTTPStatusCode = 302 // RFC 7231, 6.4.3
	StatusSeeOther          HTTPStatusCode = 303 // RFC 7231, 6.4.4
	StatusNotModified       HTTPStatusCode = 304 // RFC 7232, 4.1
	StatusUseProxy          HTTPStatusCode = 305 // RFC 7231, 6.4.5
	_                       HTTPStatusCode = 306 // RFC 7231, 6.4.6 (Unused)
	StatusTemporaryRedirect HTTPStatusCode = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect HTTPStatusCode = 308 // RFC 7538, 3

	StatusBadRequest                   HTTPStatusCode = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 HTTPStatusCode = 401 // RFC 7235, 3.1
	StatusPaymentRequired              HTTPStatusCode = 402 // RFC 7231, 6.5.2
	StatusForbidden                    HTTPStatusCode = 403 // RFC 7231, 6.5.3
	StatusNotFound                     HTTPStatusCode = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             HTTPStatusCode = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                HTTPStatusCode = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            HTTPStatusCode = 407 // RFC 7235, 3.2
	StatusRequestTimeout               HTTPStatusCode = 408 // RFC 7231, 6.5.7
	StatusConflict                     HTTPStatusCode = 409 // RFC 7231, 6.5.8
	StatusGone                         HTTPStatusCode = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               HTTPStatusCode = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           HTTPStatusCode = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        HTTPStatusCode = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            HTTPStatusCode = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         HTTPStatusCode = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable HTTPStatusCode = 416 // RFC 7233, 4.4
	StatusExpectationFailed            HTTPStatusCode = 417 // RFC 7231, 6.5.14
	StatusTeapot                       HTTPStatusCode = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           HTTPStatusCode = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          HTTPStatusCode = 422 // RFC 4918, 11.2
	StatusLocked                       HTTPStatusCode = 423 // RFC 4918, 11.3
	StatusFailedDependency             HTTPStatusCode = 424 // RFC 4918, 11.4
	StatusTooEarly                     HTTPStatusCode = 425 // RFC 8470, 5.2.
	StatusUpgradeRequired              HTTPStatusCode = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         HTTPStatusCode = 428 // RFC 6585, 3
	StatusTooManyRequests              HTTPStatusCode = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  HTTPStatusCode = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   HTTPStatusCode = 451 // RFC 7725, 3

	StatusInternalServerError           HTTPStatusCode = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                HTTPStatusCode = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    HTTPStatusCode = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            HTTPStatusCode = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                HTTPStatusCode = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       HTTPStatusCode = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         HTTPStatusCode = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           HTTPStatusCode = 507 // RFC 4918, 11.5
	StatusLoopDetected                  HTTPStatusCode = 508 // RFC 5842, 7.2
	StatusNotExtended                   HTTPStatusCode = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired HTTPStatusCode = 511 // RFC 6585, 6
)
