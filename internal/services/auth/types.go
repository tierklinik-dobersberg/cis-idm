package auth

import (
	"net/http"
	"net/url"

	"github.com/tierklinik-dobersberg/cis-idm/internal/policy"
)

type ForwardAuthInput struct {
	// Subject holds the authenticated user that is performing the request, if any.
	Subject *policy.SubjectInput `json:"subject,omitempty"`

	// Method is the HTTP method used.
	Method string `json:"method,omitempty"`

	// Path is the path of the HTTP request.
	Path string `json:"path,omitempty"`

	// Host is the requested hostname.
	Host string `json:"host,omitempty"`

	// Headers holds all request headers.
	Headers http.Header `json:"headers,omitempty"`

	// Query holds all query values.
	Query url.Values `json:"query,omitempty"`

	// ClientIP holds the IP of the client in it's string form.
	ClientIP string `json:"client_ip,omitempty"`
}

type ForwardAuthPolicyResult struct {
	// Allow should be set to true if the request should be allowed.
	Allow bool `mapstructure:"allow"`

	// StatusCode may be set to the status code to return to the client.
	// This is only used when the request is denied.
	StatusCode int `mapstructure:"status_code"`

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
	Headers map[string][]string `mapstructure:"headers"`

	// ResponseBody is the HTTP response body in case the initial request
	// is not allowed.
	// Note that response_body is only sent if status_code is a non-zero value.
	ResponseBody string `mapstructure:"response_body"`

	// AssignSubject may be used to assign a different user to the request.
	// This field is only evaluated if the request is allowed.
	AssignSubject string `mapstructure:"assign_subject"`
}
