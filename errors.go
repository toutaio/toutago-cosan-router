package cosan

import "errors"

// Common errors returned by the router.
var (
	// ErrRouterAlreadyCompiled is returned when trying to register routes after compilation.
	ErrRouterAlreadyCompiled = errors.New("cosan: router already compiled, cannot register new routes")

	// ErrConflictingRoutes is returned when two routes conflict.
	ErrConflictingRoutes = errors.New("cosan: conflicting routes detected")

	// ErrInvalidPattern is returned for invalid route patterns.
	ErrInvalidPattern = errors.New("cosan: invalid route pattern")
)
