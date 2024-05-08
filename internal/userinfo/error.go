package userinfo

import (
	"fmt"
	"net/http"
)

type Error string

func (err *Error) Error() string {
	return string(*err)
}

const (
	// The request is missing a required parameter, includes an
	// unsupported parameter or parameter value, repeats the same
	// parameter, uses more than one method for including an access
	// token, or is otherwise malformed. The resource server SHOULD
	// respond with the HTTP 400 (Bad Request) status code
	ErrInvalidRequest Error = "invalid_request"
	// The access token provided is expired, revoked, malformed, or
	// invalid for other reasons. The resource SHOULD respond with
	// the HTTP 401 (Unauthorized) status code. The client MAY
	// request a new access token and retry the protected resource
	// request.
	ErrInvalidToken Error = "invalid_token"
	// The request requires higher privileges than provided by the
	// access token. The resource server SHOULD respond with the HTTP
	// 403 (Forbidden) status code and MAY include the "scope"
	// attribute with the scope necessary to access the protected
	// resource.
	ErrInsufficientScope Error = "insufficient_scope"
)

func writeErr(w http.ResponseWriter, err Error, description string, httpStatus int) {
	key := "WWW-Authenticate"
	w.Header().Add(key, fmt.Sprintf(`Bearer error="%s"`, err.Error()))
	w.Header().Add(key, fmt.Sprintf(`error_description="%s"`, description))
	w.WriteHeader(httpStatus)
}
