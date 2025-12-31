package cosan

import (
	"net/http"
	"sync"
)

// contextPool manages the recycling of Context instances to reduce allocations
var contextPool = sync.Pool{
	New: func() interface{} {
		return &context{
			params: make(map[string]string, 4),
			values: make(map[string]interface{}, 4),
		}
	},
}

// acquireContext gets a Context from the pool
func acquireContext(w http.ResponseWriter, r *http.Request) *context {
	ctx := contextPool.Get().(*context)
	ctx.req = r
	ctx.res = w
	return ctx
}

// releaseContext returns a Context to the pool after cleaning it
func releaseContext(ctx *context) {
	// Clear maps to prevent memory leaks
	for k := range ctx.params {
		delete(ctx.params, k)
	}
	for k := range ctx.values {
		delete(ctx.values, k)
	}

	// Reset fields
	ctx.req = nil
	ctx.res = nil

	// Return to pool
	contextPool.Put(ctx)
}
