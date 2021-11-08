package nbgohttp

type httpError struct {
	BadRequest   Err
	Unauthorized Err
	Forbidden    Err
	Internal     Err
	BadGateway   Err
}

func HttpError() httpError {
	return httpError{
		BadRequest: NewError(Err{
			Code:    StatusErrorBadRequest,
			Message: "Bad Request",
		}),
		Unauthorized: NewError(Err{
			Code:    StatusErrorUnauthorized,
			Message: "Unauthorized",
		}),
		Forbidden: NewError(Err{
			Code:    StatusErrorForbidden,
			Message: "Forbidden",
		}),
		Internal: NewError(Err{
			Code:    StatusErrorInternal,
			Message: "Internal Error",
		}),
		BadGateway: NewError(Err{
			Code:    StatusErrorBadGateway,
			Message: "Bad Gateway",
		}),
	}
}
