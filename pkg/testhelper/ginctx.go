package testhelper

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// Builder builds a *gin.Context for handler/unit tests with a fluent API.
// It supports:
//   - path params, headers, and query strings
//   - three body modes (JSON, urlencoded form, multipart form with file)
//   - context injection hooks (WithContext, WithContextValue, WithContextFunc)
type Builder struct {
	// Recorder receives the handler/engine response.
	Recorder *httptest.ResponseRecorder

	// HTTP request metadata
	method  string
	path    string
	headers http.Header
	queries url.Values
	params  []gin.Param

	// request bodies
	jsonBody      any
	formFields    map[string]string
	formFileField string
	formFilePath  string

	// base context for request context
	baseCtx context.Context
	ctxFns  []func(context.Context) context.Context
}

// NewGinCtx creates a new Builder in gin.TestMode.
// If resp is nil, an internal httptest.ResponseRecorder is created.
//
// Typical usage:
//
//	rec := httptest.NewRecorder()
//	ctx := testhelper.NewGinCtx(rec).
//	    Method(http.MethodPost).
//	    Path("/v1/orders/:id").
//	    Param("id", "42").
//	    Query("include", "items").
//	    Header("X-Trace-Id", "abc").
//	    JSONBody(map[string]any{"note":"hello"}).
//	    WithContextFunc(func(c context.Context) context.Context { return c }).
//	    MustBuild(t)
//
// Then pass ctx into your handler/router under test.
func NewGinCtx(resp *httptest.ResponseRecorder) *Builder {
	gin.SetMode(gin.TestMode)
	if resp == nil {
		resp = httptest.NewRecorder()
	}
	return &Builder{
		Recorder: resp,
		method:   http.MethodGet,
		path:     "/",
		headers:  make(http.Header),
		queries:  make(url.Values),
		params:   make([]gin.Param, 0),
	}
}

// Method sets the HTTP method, default is GET.
func (b *Builder) Method(m string) *Builder { b.method = m; return b }

// Path sets the request path (may include :params to pair with Param()).
func (b *Builder) Path(p string) *Builder { b.path = p; return b }

// Header sets a single header key to value (overwrites previous value).
func (b *Builder) Header(k, v string) *Builder { b.headers.Set(k, v); return b }

// Headers sets multiple headers (overwrites per key).
func (b *Builder) Headers(h map[string]string) *Builder {
	for k, v := range h {
		b.headers.Set(k, v)
	}
	return b
}

// Query appends a query string key=value. Value is stringified via fmt.Sprint.
func (b *Builder) Query(k string, v any) *Builder {
	b.queries.Add(k, fmt.Sprint(v))
	return b
}

// Queries appends multiple query string pairs.
func (b *Builder) Queries(q map[string]any) *Builder {
	for k, v := range q {
		b.Query(k, v)
	}
	return b
}

// Param adds a Gin path parameter.
func (b *Builder) Param(k, v string) *Builder {
	b.params = append(b.params, gin.Param{Key: k, Value: v})
	return b
}

// Params appends a slice of Gin path parameters.
func (b *Builder) Params(p []gin.Param) *Builder { b.params = append(b.params, p...); return b }

// JSONBody sets a JSON request body. Mutually exclusive with any form body.
func (b *Builder) JSONBody(v any) *Builder {
	b.jsonBody = v
	b.formFields = nil
	b.formFileField, b.formFilePath = "", ""
	return b
}

// FormFields sets an application/x-www-form-urlencoded body (if no file is set)
// or contributes additional fields to a multipart/form-data body (if a file is set).
// Mutually exclusive with JSONBody.
func (b *Builder) FormFields(m map[string]string) *Builder {
	b.formFields = m
	b.jsonBody = nil
	// DO NOT clear file fields here; allows fields+file multipart
	return b
}

// FormFile sets a single file upload. If FormFields were also set, the body will be multipart
// including both fields and file. Mutually exclusive with JSONBody.
func (b *Builder) FormFile(field, filePath string) *Builder {
	b.formFileField, b.formFilePath = field, filePath
	b.jsonBody = nil
	// Keep b.formFields as-is to support multipart with fields.
	return b
}

// WithContext sets the base context for the request.
func (b *Builder) WithContext(ctx context.Context) *Builder {
	b.baseCtx = ctx
	return b
}

// WithContextValue adds a single key/value to the request context.
func (b *Builder) WithContextValue(key, val any) *Builder {
	if b.baseCtx == nil {
		b.baseCtx = context.Background()
	}
	b.baseCtx = context.WithValue(b.baseCtx, key, val)
	return b
}

// WithContextFunc appends a transformer applied to the context just before Build() returns.
func (b *Builder) WithContextFunc(fn func(context.Context) context.Context) *Builder {
	b.ctxFns = append(b.ctxFns, fn)
	return b
}

// Build assembles and returns a *gin.Context.
// Rules:
//   - JSON cannot be combined with any form body
//   - FormFields + FormFile is allowed (multipart)
//   - FormFields alone results in application/x-www-form-urlencoded
func (b *Builder) Build() (*gin.Context, error) {
	if b.method == "" {
		return nil, errors.New("method is required")
	}
	if b.path == "" {
		return nil, errors.New("path is required")
	}

	// Exclusive body rule: JSON cannot mix with any form body
	if b.jsonBody != nil && (b.formFields != nil || b.formFileField != "") {
		return nil, errors.New("cannot combine JSON body with form fields/files")
	}

	// build body + content-type
	var reqBody io.Reader
	contentType := "" // set only when we actually have a body

	switch {
	case b.jsonBody != nil:
		buf, err := json.Marshal(b.jsonBody)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(buf)
		contentType = gin.MIMEJSON

	case b.formFileField == "" && b.formFields != nil:
		// pure application/x-www-form-urlencoded
		form := url.Values{}
		for k, v := range b.formFields {
			form.Set(k, v)
		}
		reqBody = bytes.NewBufferString(form.Encode())
		contentType = "application/x-www-form-urlencoded"

	case b.formFileField != "":
		// multipart (fields + optional file)
		var err error
		reqBody, contentType, err = buildMultipart(b.formFields, b.formFileField, b.formFilePath)
		if err != nil {
			return nil, err
		}
	}

	// url
	u := &url.URL{Path: b.path, RawQuery: b.queries.Encode()}

	// request
	req := httptest.NewRequest(b.method, u.String(), reqBody)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	for k, vals := range b.headers {
		if len(vals) > 0 {
			// single-value semantics for this builder (matches Header/Headers API)
			req.Header.Set(k, vals[len(vals)-1])
		}
	}

	// compose context
	ctxBase := b.baseCtx
	if ctxBase == nil {
		ctxBase = context.Background()
	}
	for _, fn := range b.ctxFns {
		ctxBase = fn(ctxBase)
	}
	if ctxBase != nil {
		req = req.WithContext(ctxBase)
	}

	// gin context
	ctx, _ := gin.CreateTestContext(b.Recorder)
	ctx.Request = req
	if len(b.params) > 0 {
		ctx.Params = b.params
	}
	return ctx, nil
}

// MustBuild is like Build but fails the test immediately on error.
func (b *Builder) MustBuild(t *testing.T) *gin.Context {
	t.Helper()
	ctx, err := b.Build()
	if err != nil {
		require.NoError(t, err)
	}
	return ctx
}

// buildMultipart constructs a multipart/form-data body with optional fields and a single file.
func buildMultipart(fields map[string]string, fileField, filePath string) (io.Reader, string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	// normal fields
	for k, v := range fields {
		if k == fileField {
			continue // skip reserve key for file
		}
		if err := w.WriteField(k, v); err != nil {
			return nil, "", err
		}
	}

	// file field (optional)
	if fileField != "" {
		if filePath == "" {
			return nil, "", errors.New("form file path is empty")
		}
		part, err := w.CreateFormFile(fileField, filepath.Base(filePath))
		if err != nil {
			return nil, "", err
		}
		f, err := os.Open(filePath)
		if err != nil {
			return nil, "", err
		}
		defer f.Close()
		if _, err := io.Copy(part, f); err != nil {
			return nil, "", err
		}
	}

	if err := w.Close(); err != nil {
		return nil, "", err
	}
	return &buf, w.FormDataContentType(), nil
}
