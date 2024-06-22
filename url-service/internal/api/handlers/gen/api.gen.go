// Package handlers provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.3.0 DO NOT EDIT.
package handlers

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
)

// DefaultErrorBody defines model for DefaultErrorBody.
type DefaultErrorBody struct {
	Code    int     `json:"code"`
	ErrorID *string `json:"errorID,omitempty"`
	Message string  `json:"message"`
}

// CreateURLJSONBody defines parameters for CreateURL.
type CreateURLJSONBody struct {
	OriginalUrl string `json:"originalUrl"`
}

// CreateURLJSONRequestBody defines body for CreateURL for application/json ContentType.
type CreateURLJSONRequestBody CreateURLJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Create short url representation
	// (POST /data/shorten)
	CreateURL(w http.ResponseWriter, r *http.Request)
	// Redirect to actual URL
	// (GET /{shortenedUrl})
	RedirectURL(w http.ResponseWriter, r *http.Request, shortenedUrl string)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Create short url representation
// (POST /data/shorten)
func (_ Unimplemented) CreateURL(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Redirect to actual URL
// (GET /{shortenedUrl})
func (_ Unimplemented) RedirectURL(w http.ResponseWriter, r *http.Request, shortenedUrl string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// CreateURL operation middleware
func (siw *ServerInterfaceWrapper) CreateURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateURL(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

// RedirectURL operation middleware
func (siw *ServerInterfaceWrapper) RedirectURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var err error

	// ------------- Path parameter "shortenedUrl" -------------
	var shortenedUrl string

	err = runtime.BindStyledParameterWithOptions("simple", "shortenedUrl", chi.URLParam(r, "shortenedUrl"), &shortenedUrl, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "shortenedUrl", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.RedirectURL(w, r, shortenedUrl)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/data/shorten", wrapper.CreateURL)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/{shortenedUrl}", wrapper.RedirectURL)
	})

	return r
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9RWTW/cNhD9KwO2hxbQSspuUiQ6NU1cwEAaB1lvCzQxYJqalRhLJDukvF4Y/u/FUNJ+",
	"t0F6MNCbRM7H4/DNGz4IZVtnDZrgRfEgvKqxlfHzLS5l14QzIku/2HLNa3gvW9dg/OT1N7ZEUex8J/33",
	"+VtRiK7T5bjwG3ovK7aN8aAd/h8T4cg6pKAxZlUx4oMIa8fW2gSskNhuE3iz6QNpU/HeGO547zERhH91",
	"mrAUxScREake6eh0lYxO9uYLqiAe2Uubpe0BmSBViEdupW5EISpLtsTgb/V69urnildTZVsGUqJXpF3Q",
	"1ohCXNbag/Ygwce6weLjO5jXlgIaJPBId1oh3EiPJVgDoUa4cGhefziHWZqDd6j0UivJ8VL4Q4caAscc",
	"PBNY2w6UNKAIZUDwQ+wS7pC8tsaDXUJjTcWpPQQLrbxFTtRCawmhlUZWKG8aBGlK8LWk+Jd+Np/NZX0I",
	"mZHJprErH1MHC45sRbJtZdBKNs0aKjbcgOnzLi2BNOsNkhQu7YhZbg0T8GhKkPDhYn4JfG/oAyfhwlxn",
	"pQwyG454DWhKZ7UJsOrLgmBJV9rIJmKWPq5p47qQAp+Esa900wBh6MiMiWO52EP3F0DonTW+r8Dctgid",
	"x2XXQKPNrS8+mwl8InTW62BpzfVlL0eWuXP1Qx2C80WWVTrU3Q3zIqtyystZ8LfrrKNmMmalH2MshuZt",
	"R6pHWOJSG80X/m2x4PLi7YVIRKMVGh9bwciWWT1sdMTc5ZBFlq1Wq1Q6qWpMLVXZ4OSzd+dvzt7PzybT",
	"NE/r0DbM6YDU+ovlvOfcNoZfyapCSrXNokkmEhF0YHkQi11sIhEDHUUh8jRPn3FY69BIp0UhZmmezkQi",
	"nAx1FIG9m+YFZ33swP32esME0qYaCMQHjGEp9st5OZrg4uM70esA+jBqGTc2mhhWOtcMXZZ98dZslfBA",
	"9EaGLTa19EMxB6N4R8zyCaM5Urc9/69p1a7xaY06UQ7cFiN2Xey4vjLb2IE6jMl6pkdo0zz/r1XBe6cJ",
	"/esgCjHNp88n+U+TWX45nRUvXhUvXv0pErHptMUuCxurZFNbH4qX+cs8kzfquGQ7wR/E0lIrOU0pA06C",
	"blnHj6aBLk8OiX8v/iHGU7fz1SuYd0qh96wWGx5y6OffWNzvCZeiEN9l2+mcDaM5O5rLJ3CcmzvZ6HKj",
	"iCxvvuOUWEY80+mT4vmd0cTggPcK3ViXF09el4DEE4LFDAn6x8gEFgbvHaqAZRytSBBfGz3E2ZNCHGQW",
	"FkbeSd3EyTyJ82tApj2ojghNaNbQmWgQLNTSlA0OEywKXaSs79pW0vqUPBA6Qo8m9DRNRJCVZ+Fhubhi",
	"5+xhtyce+XgVnhDij1hqQrUZ1ZtBTNhPtiNdHj16ZXaSZIsBidMfBp/vDelg4QZBWeN1iYRlFLmlbgLG",
	"LuV3myjiJBHJOAD3GvtQB5Oduzvs+KsDjZzlz/758H2fP39Ssry3AX61nSlhAvt1MjbAknf+F222R9Rd",
	"NkkVuv5Bd8xPdopRetb075tMOp3dPRN8d4P94YWd3SGtQ80PByZPR82135Ilzu2rx78DAAD//zPx+bwZ",
	"DQAA",
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
