package auth

import (
	"net/http"
	"net/url"

	"github.com/tierklinik-dobersberg/cis-idm/internal/policy"
)

type ForwardAuthInput struct {
	// Subject holds the authenticated user that is performing the request, if any.
	Subject *policy.SubjectInput `json:"subject"`

	// Method is the HTTP method used.
	Method string `json:"method"`

	// Path is the path of the HTTP request.
	Path string `json:"path"`

	// Host is the requested hostname.
	Host string `json:"host"`

	// Headers holds all request headers.
	Headers http.Header `json:"headers"`

	// Query holds all query values.
	Query url.Values `json:"query"`

	// ClientIP holds the IP of the client in it's string form.
	ClientIP string `json:"client_ip"`
}

type ForwardAuthPolicyResult struct {
	// Allow should be set to true if the request should be allowed.
	Allow bool `json:"allow"`

	// StatusCode may be set to the status code to return to the client.
	// This is only used when the request is denied.
	StatusCode int `json:"status_code"`

	// Headers holds additional headers that are added to the response to
	// the reverse proxy.
	//
	// If the request is denied (Allow = false) then the headers are directly
	// returned to the client making the request.
	//
	// If the request is allowed (Allow = true) then those headers are sent
	// to the reverse proxy which might decide to forward those headers to
	// the upstream server. Note that forwarding of headers might require
	// configuration on the reverse proxy side.
	Headers map[string][]string `json:"headers"`
}
