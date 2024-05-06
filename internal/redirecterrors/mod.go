package redirecterrors

import (
	"fmt"
	"net/url"
)

type Error string

func (err Error) String() string {
	return string(err)
}

const (
	// The request is missing a required parameter, includes an
	// invalid parameter value, includes a parameter more than
	// once, or is otherwise malformed.
	ErrInvalidRequest Error = "invalid_request"
	// The client is not authorized to request an authorization
	// code using this method.
	ErrUnauthorizedClient Error = "unauthorized_client"
	// The resource owner or authorization server denied the
	// request.
	ErrAccessDenied Error = "access_denied"
	// The authorization server does not support obtaining an
	// authorization code using this method.
	ErrUnsupportedResponseType Error = "unsupported_response_type"
	// The requested scope is invalid, unknown, or malformed.
	ErrInvalidScope Error = "invalid_scope"
	// The authorization server encountered an unexpected
	// condition that prevented it from fulfilling the request.
	// (This error code is needed because a 500 Internal Server
	// Error HTTP status code cannot be returned to the client
	// via an HTTP redirect.)
	ErrServerError Error = "server_error"
	// The authorization server is currently unable to handle
	// the request due to a temporary overloading or maintenance
	// of the server.  (This error code is needed because a 503
	// Service Unavailable HTTP status code cannot be returned
	// to the client via an HTTP redirect.)
	ErrTemporarilyUnavailable Error = "temporarily_unavailable"
)

func URI(
	uri string,
	err Error,
	description string,
	state string,
) string {
	return fmt.Sprintf("%s?%s", uri, url.Values{
		"error":             {err.String()},
		"error_description": {description},
		"state":             {state},
	}.Encode())
}
