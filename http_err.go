package noob

import (
	"errors"
	"github.com/alfarih31/nb-go-http/app_err"
)

const StatusBadRequest = "400"
const StatusUnauthorized = "401"
const StatusForbidden = "403"
const StatusInternalError = "500"
const StatusBadGateway = "502"
const StatusNotFound = "404"
const StatusNoMethod = "405"
const StatusTooManyRequest = "429"

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
		Code: StatusBadRequest,
		Err:  errors.New("bad request"),
	},
	Unauthorized: &apperr.AppErr{
		Code: StatusUnauthorized,
		Err:  errors.New("unauthorized"),
	},
	Forbidden: &apperr.AppErr{
		Code: StatusForbidden,
		Err:  errors.New("forbidden"),
	},
	Internal: &apperr.AppErr{
		Code: StatusInternalError,
		Err:  errors.New("internal error"),
	},
	BadGateway: &apperr.AppErr{
		Code: StatusBadGateway,
		Err:  errors.New("bad gateway"),
	},
	NotFound: &apperr.AppErr{
		Code: StatusNotFound,
		Err:  errors.New("not found"),
	},
	NoMethod: &apperr.AppErr{
		Code: StatusNoMethod,
		Err:  errors.New("no method"),
	},
	TooManyRequest: &apperr.AppErr{
		Code: StatusTooManyRequest,
		Err:  errors.New("too many request"),
	},
}
