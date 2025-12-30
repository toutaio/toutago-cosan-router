package cosan

import (
	"fmt"
	"sync"
)

// simpleMatcher is a basic exact-match implementation of the Matcher interface.
// This is used in Phase 1 for simple routing. Phase 2 will add radix tree matcher.
type simpleMatcher struct {
	routes   map[string]*route // key: "METHOD:PATH"
	compiled bool
	mu       sync.RWMutex
}

// newSimpleMatcher creates a new simple matcher.
func newSimpleMatcher() Matcher {
	return &simpleMatcher{
		routes:   make(map[string]*route),
		compiled: false,
	}
}

// Register adds a route to the matcher.
func (m *simpleMatcher) Register(method, pattern string, handler HandlerFunc) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.compiled {
		return fmt.Errorf("matcher already compiled")
	}

	key := method + ":" + pattern
	if _, exists := m.routes[key]; exists {
		return fmt.Errorf("duplicate route: %s %s", method, pattern)
	}

	m.routes[key] = &route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}

	return nil
}

// Match finds a route matching the method and path.
// For Phase 1, this only does exact matching.
func (m *simpleMatcher) Match(method, path string) (*Route, map[string]string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := method + ":" + path
	if rt, found := m.routes[key]; found {
		// Return the route as Route interface pointer
		var routeInterface Route = rt
		// No params for exact match
		return &routeInterface, make(map[string]string), true
	}

	return nil, nil, false
}

// Compile optimizes the matcher (no-op for simple matcher).
func (m *simpleMatcher) Compile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.compiled = true
	return nil
}
