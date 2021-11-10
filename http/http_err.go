package http

import (
	"github.com/alfarih31/nb-go-http/app_error"
)

type httpError struct {
	BadRequest   *apperror.Err
	Unauthorized *apperror.Err
	Forbidden    *apperror.Err
	Internal     *apperror.Err
	BadGateway   *apperror.Err
}

var HTTPError = httpError{
	BadRequest: apperror.NewError(apperror.Err{
		Code:    "400",
		Message: "Bad Request",
	}),
	Unauthorized: apperror.NewError(apperror.Err{
		Code:    "401",
		Message: "Unauthorized",
	}),
	Forbidden: apperror.NewError(apperror.Err{
		Code:    "403",
		Message: "Forbidden",
	}),
	Internal: apperror.NewError(apperror.Err{
		Code:    "500",
		Message: "Internal Error",
	}),
	BadGateway: apperror.NewError(apperror.Err{
		Code:    "502",
		Message: "Bad Gateway",
	}),
}
