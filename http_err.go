package nbgohttp

type httpError struct {
	BadRequest   *Err
	Unauthorized *Err
	Forbidden    *Err
	Internal     *Err
	BadGateway   *Err
}

var HTTPError = httpError{
	BadRequest: NewError(Err{
		Code:    "400",
		Message: "Bad Request",
	}),
	Unauthorized: NewError(Err{
		Code:    "401",
		Message: "Unauthorized",
	}),
	Forbidden: NewError(Err{
		Code:    "403",
		Message: "Forbidden",
	}),
	Internal: NewError(Err{
		Code:    "500",
		Message: "Internal Error",
	}),
	BadGateway: NewError(Err{
		Code:    "502",
		Message: "Bad Gateway",
	}),
}
