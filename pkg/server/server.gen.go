// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
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

// PostRequest A post object to be created
type PostRequest struct {
	Content string `json:"content"`

	// Description A short description of the post for the index page
	Description string `json:"description"`

	// Keywords Keywords for the post for SEO and search purposes
	Keywords *[]string `json:"keywords,omitempty"`

	// Slug The URL slug of the post. Should be unique and URL-friendly.
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

// PostResponse A post object after it's been created or fetched
type PostResponse struct {
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Id          int       `json:"id"`
	Keywords    *[]string `json:"keywords,omitempty"`

	// ReadingTime Approximate post reading time in seconds
	ReadingTime int       `json:"reading_time"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PostsListItem defines model for PostsListItem.
type PostsListItem struct {
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	Keywords    *[]string `json:"keywords,omitempty"`

	// ReadingTime Approximate post reading time in seconds
	ReadingTime int    `json:"reading_time"`
	Slug        string `json:"slug"`
	Title       string `json:"title"`
}

// PostsListResponse A list of posts
type PostsListResponse struct {
	Posts []PostsListItem `json:"posts"`
	Total int             `json:"total"`
}

// PutPostRequest A post object to be updated
type PutPostRequest struct {
	Content string `json:"content"`

	// Description A short description of the post for the index page
	Description string `json:"description"`

	// Keywords Keywords for the post for SEO and search purposes
	Keywords *[]string `json:"keywords,omitempty"`
	Title    string    `json:"title"`
}

// RequestError defines model for RequestError.
type RequestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// GetPostsParams defines parameters for GetPosts.
type GetPostsParams struct {
	// Page Page number
	Page *int `form:"page,omitempty" json:"page,omitempty"`

	// Limit Number of items per page
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// PostLoginGithubAuthorizeJSONRequestBody defines body for PostLoginGithubAuthorize for application/json ContentType.
type PostLoginGithubAuthorizeJSONRequestBody = GitHubAuthRequestBody

// PostPostsJSONRequestBody defines body for PostPosts for application/json ContentType.
type PostPostsJSONRequestBody = PostRequest

// PutPostsSlugJSONRequestBody defines body for PutPostsSlug for application/json ContentType.
type PutPostsSlugJSONRequestBody = PutPostRequest

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
	// Get all posts
	// (GET /posts)
	GetPosts(ctx echo.Context, params GetPostsParams) error
	// Create a new post
	// (POST /posts)
	PostPosts(ctx echo.Context) error
	// Get a post by slug
	// (GET /posts/{slug})
	GetPostsSlug(ctx echo.Context, slug string) error
	// Update a post by slug
	// (PUT /posts/{slug})
	PutPostsSlug(ctx echo.Context, slug string) error
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

// GetPosts converts echo context to params.
func (w *ServerInterfaceWrapper) GetPosts(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetPostsParams
	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPosts(ctx, params)
	return err
}

// PostPosts converts echo context to params.
func (w *ServerInterfaceWrapper) PostPosts(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPosts(ctx)
	return err
}

// GetPostsSlug converts echo context to params.
func (w *ServerInterfaceWrapper) GetPostsSlug(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "slug" -------------
	var slug string

	err = runtime.BindStyledParameterWithOptions("simple", "slug", ctx.Param("slug"), &slug, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter slug: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPostsSlug(ctx, slug)
	return err
}

// PutPostsSlug converts echo context to params.
func (w *ServerInterfaceWrapper) PutPostsSlug(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "slug" -------------
	var slug string

	err = runtime.BindStyledParameterWithOptions("simple", "slug", ctx.Param("slug"), &slug, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter slug: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PutPostsSlug(ctx, slug)
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
	router.GET(baseURL+"/posts", wrapper.GetPosts)
	router.POST(baseURL+"/posts", wrapper.PostPosts)
	router.GET(baseURL+"/posts/:slug", wrapper.GetPostsSlug)
	router.PUT(baseURL+"/posts/:slug", wrapper.PutPostsSlug)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
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
