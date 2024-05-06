package clienterrors

import (
	"encoding/json"
	"io"
)

type error struct {
	Error       string `json:"error"`
	Description string `json:"error_description"`
}

type Error string

const (
	// The request is missing a required parameter, includes an
	// unsupported parameter value (other than grant type),
	// repeats a parameter, includes multiple credentials,
	// utilizes more than one mechanism for authenticating the
	// client, or is otherwise malformed.
	ErrInvalidRequest Error = "invalid_request"

	// Client authentication failed (e.g., unknown client, no
	// client authentication included, or unsupported
	// authentication method).  The authorization server MAY
	// return an HTTP 401 (Unauthorized) status code to indicate
	// which HTTP authentication schemes are supported.  If the
	// client attempted to authenticate via the "Authorization"
	// request header field, the authorization server MUST
	// respond with an HTTP 401 (Unauthorized) status code and
	// include the "WWW-Authenticate" response header field
	// matching the authentication scheme used by the client.
	ErrInvalidClient Error = "invalid_client"

	// The provided authorization grant (e.g., authorization
	// code, resource owner credentials) or refresh token is
	// invalid, expired, revoked, does not match the redirection
	// URI used in the authorization request, or was issued to
	// another client.
	ErrInvalidGrant Error = "invalid_grant"

	// The authenticated client is not authorized to use this
	// authorization grant type.
	ErrUnauthorizedClient Error = "unauthorized_client"

	// The authorization grant type is not supported by the
	// authorization server.
	ErrUnsupportedGrantType Error = "unsupported_grant_type"

	// The requested scope is invalid, unknown, malformed, or
	// exceeds the scope granted by the resource owner.
	ErrUnvalidScope Error = "invalid_scope"
)

func Write(
	w io.Writer,
	err Error,
	description string,
) {
	json.NewEncoder(w).Encode(&error{
		Error:       string(err),
		Description: description,
	})
}
