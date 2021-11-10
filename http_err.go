package nbgohttp

type httpError struct {
	BadRequest   *Err
	Unauthorized *Err
	Forbidden    *Err
	Internal     *Err
	BadGateway   *Err
	NotFound     *Err
	NoMethod     *Err
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
	NotFound: NewError(Err{
		Code:    "404",
		Message: "Not Found",
	}),
	NoMethod: NewError(Err{
		Code:    "405",
		Message: "No Method",
	}),
}
