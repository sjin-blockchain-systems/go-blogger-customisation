// Package client provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package client

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// GitHubAuthRequestBody defines model for GitHubAuthRequestBody.
type GitHubAuthRequestBody struct {
	Code string `json:"code"`
}

// HealthCheckResponse defines model for HealthCheckResponse.
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// JWTToken defines model for JWTToken.
type JWTToken struct {
	Token string `json:"token"`
}

// RequestError defines model for RequestError.
type RequestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PostLoginGithubAuthorizeJSONRequestBody defines body for PostLoginGithubAuthorize for application/json ContentType.
type PostLoginGithubAuthorizeJSONRequestBody = GitHubAuthRequestBody

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client which conforms to the OpenAPI3 specification for this service.
type Client struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Client) error

// Creates a new Client, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Client, error) {
	// create a client with sane default values
	client := Client{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Client) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Client) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// GetHealth request
	GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostLoginGithubAuthorizeWithBody request with any body
	PostLoginGithubAuthorizeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error)

	PostLoginGithubAuthorize(ctx context.Context, body PostLoginGithubAuthorizeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error)

	// PostLoginRefresh request
	PostLoginRefresh(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Client) GetHealth(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetHealthRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostLoginGithubAuthorizeWithBody(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostLoginGithubAuthorizeRequestWithBody(c.Server, contentType, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostLoginGithubAuthorize(ctx context.Context, body PostLoginGithubAuthorizeJSONRequestBody, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostLoginGithubAuthorizeRequest(c.Server, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Client) PostLoginRefresh(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewPostLoginRefreshRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewGetHealthRequest generates requests for GetHealth
func NewGetHealthRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/health")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewPostLoginGithubAuthorizeRequest calls the generic PostLoginGithubAuthorize builder with application/json body
func NewPostLoginGithubAuthorizeRequest(server string, body PostLoginGithubAuthorizeJSONRequestBody) (*http.Request, error) {
	var bodyReader io.Reader
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	bodyReader = bytes.NewReader(buf)
	return NewPostLoginGithubAuthorizeRequestWithBody(server, "application/json", bodyReader)
}

// NewPostLoginGithubAuthorizeRequestWithBody generates requests for PostLoginGithubAuthorize with any type of body
func NewPostLoginGithubAuthorizeRequestWithBody(server string, contentType string, body io.Reader) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/login/github/authorize")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)

	return req, nil
}

// NewPostLoginRefreshRequest generates requests for PostLoginRefresh
func NewPostLoginRefreshRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/login/refresh")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Client) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// GetHealthWithResponse request
	GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error)

	// PostLoginGithubAuthorizeWithBodyWithResponse request with any body
	PostLoginGithubAuthorizeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostLoginGithubAuthorizeResponse, error)

	PostLoginGithubAuthorizeWithResponse(ctx context.Context, body PostLoginGithubAuthorizeJSONRequestBody, reqEditors ...RequestEditorFn) (*PostLoginGithubAuthorizeResponse, error)

	// PostLoginRefreshWithResponse request
	PostLoginRefreshWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PostLoginRefreshResponse, error)
}

type GetHealthResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *HealthCheckResponse
}

// Status returns HTTPResponse.Status
func (r GetHealthResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetHealthResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostLoginGithubAuthorizeResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *JWTToken
	JSON400      *RequestError
	JSON403      *RequestError
}

// Status returns HTTPResponse.Status
func (r PostLoginGithubAuthorizeResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostLoginGithubAuthorizeResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type PostLoginRefreshResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *JWTToken
	JSON400      *RequestError
	JSON401      *RequestError
}

// Status returns HTTPResponse.Status
func (r PostLoginRefreshResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r PostLoginRefreshResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// GetHealthWithResponse request returning *GetHealthResponse
func (c *ClientWithResponses) GetHealthWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetHealthResponse, error) {
	rsp, err := c.GetHealth(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetHealthResponse(rsp)
}

// PostLoginGithubAuthorizeWithBodyWithResponse request with arbitrary body returning *PostLoginGithubAuthorizeResponse
func (c *ClientWithResponses) PostLoginGithubAuthorizeWithBodyWithResponse(ctx context.Context, contentType string, body io.Reader, reqEditors ...RequestEditorFn) (*PostLoginGithubAuthorizeResponse, error) {
	rsp, err := c.PostLoginGithubAuthorizeWithBody(ctx, contentType, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostLoginGithubAuthorizeResponse(rsp)
}

func (c *ClientWithResponses) PostLoginGithubAuthorizeWithResponse(ctx context.Context, body PostLoginGithubAuthorizeJSONRequestBody, reqEditors ...RequestEditorFn) (*PostLoginGithubAuthorizeResponse, error) {
	rsp, err := c.PostLoginGithubAuthorize(ctx, body, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostLoginGithubAuthorizeResponse(rsp)
}

// PostLoginRefreshWithResponse request returning *PostLoginRefreshResponse
func (c *ClientWithResponses) PostLoginRefreshWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*PostLoginRefreshResponse, error) {
	rsp, err := c.PostLoginRefresh(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParsePostLoginRefreshResponse(rsp)
}

// ParseGetHealthResponse parses an HTTP response from a GetHealthWithResponse call
func ParseGetHealthResponse(rsp *http.Response) (*GetHealthResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetHealthResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest HealthCheckResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	}

	return response, nil
}

// ParsePostLoginGithubAuthorizeResponse parses an HTTP response from a PostLoginGithubAuthorizeWithResponse call
func ParsePostLoginGithubAuthorizeResponse(rsp *http.Response) (*PostLoginGithubAuthorizeResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostLoginGithubAuthorizeResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest JWTToken
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest RequestError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 403:
		var dest RequestError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON403 = &dest

	}

	return response, nil
}

// ParsePostLoginRefreshResponse parses an HTTP response from a PostLoginRefreshWithResponse call
func ParsePostLoginRefreshResponse(rsp *http.Response) (*PostLoginRefreshResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &PostLoginRefreshResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest JWTToken
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest RequestError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest RequestError
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	}

	return response, nil
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Health check
	// (GET /health)
	GetHealth(ctx echo.Context) error
	// Authorize with GitHub
	// (POST /login/github/authorize)
	PostLoginGithubAuthorize(ctx echo.Context) error
	// Refresh the JWT token
	// (POST /login/refresh)
	PostLoginRefresh(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetHealth converts echo context to params.
func (w *ServerInterfaceWrapper) GetHealth(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetHealth(ctx)
	return err
}

// PostLoginGithubAuthorize converts echo context to params.
func (w *ServerInterfaceWrapper) PostLoginGithubAuthorize(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostLoginGithubAuthorize(ctx)
	return err
}

// PostLoginRefresh converts echo context to params.
func (w *ServerInterfaceWrapper) PostLoginRefresh(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostLoginRefresh(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/health", wrapper.GetHealth)
	router.POST(baseURL+"/login/github/authorize", wrapper.PostLoginGithubAuthorize)
	router.POST(baseURL+"/login/refresh", wrapper.PostLoginRefresh)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RVUW/bNhD+K4fbgO5Bs5yk6zYFwZAUnpM2W4LEQx+WIaWos8REIlXylMQt/N8HUrId",
	"2wraAl1e9mIT5PG+4933ffqE0lS10aTZYfIJnSyoEmE5VnzcpIcNFxf0oSHHRyab+YPamposKwph0mTk",
	"/+lBVHVJmODO7t7Ln179/MuvQ5HKjKYYIc9qf+LYKp3jfB6hpQ+NspRh8neb4p9llElvSDLOIzwmUXLx",
	"uiB5e0GuNtrRdgGOBTduvYSzt58F7a71wb55N5mYW9LbWLzYXkHd3PM+nAsuDuJ9OGauz3Q524dLko2l",
	"z1bRJuwromv6yFpjv6Tr5AOvw/4WaIQVOSfyjSshNyyOvmhKq0zbNfsLSk9NW59mIdkvtah81KWoYGw+",
	"VkJjhI0tMcGCuXZJHOeKiyYdSFPFTlR5CIpz82Namjyf+fIzctKqmpXRmOAhOOWfAFfNcLj7CkqVF3xP",
	"/hdSIW9JZzA1FjK6o9J3zb0A/2u0KMEndRhhqSR1fOoq/ONkAqfd7teVGKelSeNKKB2fnrwe/Xk58kWz",
	"4tDmsYGjEAaH5ycY4R1Z1z5kOBgOdnyoqUmLWmGCe4OdwRAjrAUXYdJxEUTglzmFhq43I6gD1BS4IHBk",
	"78iCctDUIHQGttHaDzNAWOHvnGS+JuJWXOiH3EorwO0Oh4v5kQ5woq5LJcPV+MYZvbIJv/re0hQT/C5e",
	"+UjcmUjcJ9/AkvUXnL0NXHNNVQk7w6STPUh/MRzFpcmV7oYQi4YLY9XH1guM62nK6EEWQucEAlobA89e",
	"+GFqTeWncKXHowm875lui2Q8xgroN1kq0nytsoOWcpYyZUnydWPVwfsr7dkm4M27CQRBX+mthp8bx6c+",
	"9zhgHS7f0Irskb1+k9b3m/d8XdNsG5r/h/Nf+ugTQ4/w5TdEWzPMHsQjkUEX00LvPRv078amKstIQ3Dp",
	"hVob12pVGwZRluaeMmADQkpyLkSIrFIaaqGp3NDIkkBwr7joWP5YLJamllzxtEYu2oCAs2SuL03xCwf0",
	"UAeSRFiQyMgGcixARZtiM+PkcSbcZFr0qJmrb9ARCUs2OPmebMXjl14XU2MrwZhgGmJ6PlDzJ1XWPQ7/",
	"v+zeeTbov/TSKLOvIPgGoXvp6OHm/wYAAP//A8n3cKQKAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
