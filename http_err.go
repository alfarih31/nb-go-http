package noob

import (
	"errors"
	"github.com/alfarih31/nb-go-http/app_err"
	"net/http"
)

type httpError struct {
	BadRequest     *apperr.AppErr
	Unauthorized   *apperr.AppErr
	Forbidden      *apperr.AppErr
	Internal       *apperr.AppErr
	BadGateway     *apperr.AppErr
	NotFound       *apperr.AppErr
	NoMethod       *apperr.AppErr
	TooManyRequest *apperr.AppErr
}

var HTTPError = httpError{
	BadRequest: &apperr.AppErr{
		Code: http.StatusBadRequest,
		Err:  errors.New("bad request"),
	},
	Unauthorized: &apperr.AppErr{
		Code: http.StatusUnauthorized,
		Err:  errors.New("unauthorized"),
	},
	Forbidden: &apperr.AppErr{
		Code: http.StatusForbidden,
		Err:  errors.New("forbidden"),
	},
	Internal: &apperr.AppErr{
		Code: http.StatusInternalServerError,
		Err:  errors.New("internal error"),
	},
	BadGateway: &apperr.AppErr{
		Code: http.StatusBadGateway,
		Err:  errors.New("bad gateway"),
	},
	NotFound: &apperr.AppErr{
		Code: http.StatusNotFound,
		Err:  errors.New("not found"),
	},
	NoMethod: &apperr.AppErr{
		Code: http.StatusMethodNotAllowed,
		Err:  errors.New("no method"),
	},
	TooManyRequest: &apperr.AppErr{
		Code: http.StatusTooManyRequests,
		Err:  errors.New("too many request"),
	},
}
