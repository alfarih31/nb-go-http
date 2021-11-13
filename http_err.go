package noob

import "github.com/alfarih31/nb-go-http/app_err"

type httpError struct {
	BadRequest   apperr.AppErr
	Unauthorized apperr.AppErr
	Forbidden    apperr.AppErr
	Internal     apperr.AppErr
	BadGateway   apperr.AppErr
	NotFound     apperr.AppErr
	NoMethod     apperr.AppErr
}

var HTTPError = httpError{
	BadRequest: apperr.New(
		"Bad Request",
		"400",
	),
	Unauthorized: apperr.New(
		"Unauthorized",
		"401",
	),
	Forbidden: apperr.New(
		"Forbidden",
		"403",
	),
	Internal: apperr.New(
		"Internal Error",
		"500",
	),
	BadGateway: apperr.New(
		"Bad Gateway",
		"502",
	),
	NotFound: apperr.New(
		"Not Found",
		"404",
	),
	NoMethod: apperr.New(
		"No Method",
		"405",
	),
}
