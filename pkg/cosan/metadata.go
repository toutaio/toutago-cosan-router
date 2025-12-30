package cosan

// RouteMetadata contains metadata about a route for documentation and introspection
type RouteMetadata struct {
	Name        string
	Description string
	Tags        []string
	Deprecated  bool
	Version     string
}

// RouteInfo contains information about a registered route
type RouteInfo struct {
	Method      string
	Pattern     string
	Name        string
	Description string
	Tags        []string
	Deprecated  bool
	Version     string
}

// routeWithMetadata extends route with metadata capabilities
type routeWithMetadata struct {
	*route
	metadata *RouteMetadata
}

// WithName sets the name of the route for documentation
func WithName(name string) RouteOption {
	return func(r *route) {
		if r.metadata == nil {
			r.metadata = &RouteMetadata{}
		}
		r.metadata.Name = name
	}
}

// WithDescription sets the description of the route for documentation
func WithDescription(desc string) RouteOption {
	return func(r *route) {
		if r.metadata == nil {
			r.metadata = &RouteMetadata{}
		}
		r.metadata.Description = desc
	}
}

// WithTags adds tags to the route for categorization
func WithTags(tags ...string) RouteOption {
	return func(r *route) {
		if r.metadata == nil {
			r.metadata = &RouteMetadata{}
		}
		r.metadata.Tags = append(r.metadata.Tags, tags...)
	}
}

// Deprecated marks the route as deprecated
func Deprecated() RouteOption {
	return func(r *route) {
		if r.metadata == nil {
			r.metadata = &RouteMetadata{}
		}
		r.metadata.Deprecated = true
	}
}

// WithVersion sets the API version for the route
func WithVersion(version string) RouteOption {
	return func(r *route) {
		if r.metadata == nil {
			r.metadata = &RouteMetadata{}
		}
		r.metadata.Version = version
	}
}

// RouteOption is a functional option for configuring route metadata
type RouteOption func(*route)

// GetRoutes returns all registered routes with metadata for introspection
func (r *router) GetRoutes() []RouteInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	routes := make([]RouteInfo, 0, len(r.routes))
	for _, route := range r.routes {
		info := RouteInfo{
			Method:  route.method,
			Pattern: route.pattern,
		}
		
		if route.metadata != nil {
			info.Name = route.metadata.Name
			info.Description = route.metadata.Description
			info.Tags = route.metadata.Tags
			info.Deprecated = route.metadata.Deprecated
			info.Version = route.metadata.Version
		}
		
		routes = append(routes, info)
	}
	
	return routes
}

// FindRoute finds a route by name
func (r *router) FindRoute(name string) *RouteInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, route := range r.routes {
		if route.metadata != nil && route.metadata.Name == name {
			return &RouteInfo{
				Method:      route.method,
				Pattern:     route.pattern,
				Name:        route.metadata.Name,
				Description: route.metadata.Description,
				Tags:        route.metadata.Tags,
				Deprecated:  route.metadata.Deprecated,
				Version:     route.metadata.Version,
			}
		}
	}
	
	return nil
}
