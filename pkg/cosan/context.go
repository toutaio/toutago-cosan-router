package cosan

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// context is the default implementation of the Context interface.
type context struct {
	req    *http.Request
	res    http.ResponseWriter
	params map[string]string
	values map[string]interface{}
}

// newContext creates a new context for a request.
func newContext(w http.ResponseWriter, r *http.Request, params map[string]string) Context {
	return &context{
		req:    r,
		res:    w,
		params: params,
		values: make(map[string]interface{}),
	}
}

// Request returns the underlying *http.Request.
func (c *context) Request() *http.Request {
	return c.req
}

// Response returns the underlying http.ResponseWriter.
func (c *context) Response() http.ResponseWriter {
	return c.res
}

// Param returns the value of the named path parameter.
func (c *context) Param(key string) string {
	return c.params[key]
}

// Params returns all path parameters as a map.
func (c *context) Params() map[string]string {
	return c.params
}

// Query returns the first value of the named query parameter.
func (c *context) Query(key string) string {
	return c.req.URL.Query().Get(key)
}

// QueryAll returns all values of the named query parameter.
func (c *context) QueryAll(key string) []string {
	return c.req.URL.Query()[key]
}

// Bind parses the request body into the provided struct.
// For Phase 1, this only supports JSON.
func (c *context) Bind(v interface{}) error {
	contentType := c.req.Header.Get("Content-Type")
	
	// For Phase 1, only support JSON
	if contentType != "application/json" && contentType != "" {
		return fmt.Errorf("unsupported content type: %s", contentType)
	}

	decoder := json.NewDecoder(c.req.Body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	return nil
}

// BodyBytes returns the raw request body as bytes.
func (c *context) BodyBytes() ([]byte, error) {
	return io.ReadAll(c.req.Body)
}

// JSON writes a JSON response with the given status code.
func (c *context) JSON(code int, v interface{}) error {
	c.res.Header().Set("Content-Type", "application/json")
	c.res.WriteHeader(code)

	encoder := json.NewEncoder(c.res)
	if err := encoder.Encode(v); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// String writes a formatted string response with the given status code.
func (c *context) String(code int, format string, args ...interface{}) error {
	c.res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	c.res.WriteHeader(code)
	fmt.Fprintf(c.res, format, args...)
	return nil
}

// HTML writes an HTML response with the given status code.
func (c *context) HTML(code int, html string) error {
	c.res.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.res.WriteHeader(code)
	_, err := c.res.Write([]byte(html))
	return err
}

// Status sets the HTTP status code.
func (c *context) Status(code int) {
	c.res.WriteHeader(code)
}

// Header returns the response header map.
func (c *context) Header() http.Header {
	return c.res.Header()
}

// Write writes the response body bytes.
func (c *context) Write(b []byte) (int, error) {
	return c.res.Write(b)
}

// Set stores a value in the context for the request lifetime.
func (c *context) Set(key string, value interface{}) {
	c.values[key] = value
}

// Get retrieves a value from the context.
func (c *context) Get(key string) interface{} {
	return c.values[key]
}
