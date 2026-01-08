package cosan

import (
	"net/http/httptest"
	"testing"
)

func TestContextPool_AcquireRelease(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	// Acquire a context
	ctx := acquireContext(w, r)
	if ctx == nil {
		t.Fatal("Failed to acquire context")
	}
	if ctx.req != r {
		t.Error("Request not set correctly")
	}
	if ctx.res != w {
		t.Error("ResponseWriter not set correctly")
	}

	// Add some data
	ctx.params["id"] = "123"
	ctx.values["key"] = "value"

	// Release it
	releaseContext(ctx)

	// Verify it was cleaned
	if len(ctx.params) != 0 {
		t.Error("Params not cleared after release")
	}
	if len(ctx.values) != 0 {
		t.Error("Values not cleared after release")
	}
	if ctx.req != nil {
		t.Error("Request not cleared after release")
	}
	if ctx.res != nil {
		t.Error("ResponseWriter not cleared after release")
	}
}

func TestContextPool_Reuse(t *testing.T) {
	w1 := httptest.NewRecorder()
	r1 := httptest.NewRequest("GET", "/test1", nil)

	// Acquire and release first context
	ctx1 := acquireContext(w1, r1)
	ctx1.params["id"] = "123"
	releaseContext(ctx1)

	// Acquire second context (should reuse from pool)
	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest("GET", "/test2", nil)
	ctx2 := acquireContext(w2, r2)

	// Verify it's clean and has new request/response
	if len(ctx2.params) != 0 {
		t.Error("Context from pool not clean")
	}
	if ctx2.req != r2 {
		t.Error("Request not updated")
	}
	if ctx2.res != w2 {
		t.Error("ResponseWriter not updated")
	}

	releaseContext(ctx2)
}

func TestContextPool_ConcurrentUsage(t *testing.T) {
	const concurrency = 100

	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(id int) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/test", nil)

			ctx := acquireContext(w, r)
			ctx.params["id"] = string(rune(id))
			ctx.values["count"] = id

			// Simulate some work
			_ = ctx.Request()
			_ = ctx.Response()

			releaseContext(ctx)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < concurrency; i++ {
		<-done
	}
}

func TestContextPool_MemoryLeaks(t *testing.T) {
	// Test that large maps don't cause memory leaks
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	ctx := acquireContext(w, r)

	// Add many params and values
	for i := 0; i < 1000; i++ {
		key := string(rune(i))
		ctx.params[key] = key
		ctx.values[key] = i
	}

	// Release should clear everything
	releaseContext(ctx)

	if len(ctx.params) != 0 {
		t.Errorf("Expected params to be cleared, got %d items", len(ctx.params))
	}
	if len(ctx.values) != 0 {
		t.Errorf("Expected values to be cleared, got %d items", len(ctx.values))
	}
}

func BenchmarkContextPool_AcquireRelease(b *testing.B) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := acquireContext(w, r)
		ctx.params["id"] = "123"
		releaseContext(ctx)
	}
}

func BenchmarkContextPool_NoPool(b *testing.B) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := &context{
			req:    r,
			res:    w,
			params: make(map[string]string, 4),
			values: make(map[string]interface{}, 4),
		}
		ctx.params["id"] = "123"
		_ = ctx
	}
}
