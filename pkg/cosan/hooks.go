package cosan

import "net/http"

// hooks stores router-level hooks for lifecycle events
type hooks struct {
	beforeRequest []RequestHook
	afterResponse []ResponseHook
	errorHandler  ErrorHandler
}

// BeforeRequest registers a hook to run before each request
func (r *router) BeforeRequest(hook RequestHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.hooks == nil {
		r.hooks = &hooks{}
	}
	r.hooks.beforeRequest = append(r.hooks.beforeRequest, hook)
}

// AfterResponse registers a hook to run after each response
func (r *router) AfterResponse(hook ResponseHook) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.hooks == nil {
		r.hooks = &hooks{}
	}
	r.hooks.afterResponse = append(r.hooks.afterResponse, hook)
}

// SetErrorHandler sets a custom error handler for the router
func (r *router) SetErrorHandler(handler ErrorHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	if r.hooks == nil {
		r.hooks = &hooks{}
	}
	r.hooks.errorHandler = handler
}

// executeBeforeHooks runs all before-request hooks
func (r *router) executeBeforeHooks(req *http.Request) error {
	if r.hooks == nil {
		return nil
	}
	
	for _, hook := range r.hooks.beforeRequest {
		if err := hook(req); err != nil {
			return err
		}
	}
	
	return nil
}

// executeAfterHooks runs all after-response hooks
func (r *router) executeAfterHooks(req *http.Request, statusCode int) {
	if r.hooks == nil {
		return
	}
	
	for _, hook := range r.hooks.afterResponse {
		hook(req, statusCode)
	}
}

// handleError handles errors using custom handler if set
func (r *router) handleError(ctx Context, err error) {
	if r.hooks != nil && r.hooks.errorHandler != nil {
		r.hooks.errorHandler(ctx, err)
		return
	}
	
	// Default error handling
	ctx.String(500, "Internal Server Error: "+err.Error())
}
