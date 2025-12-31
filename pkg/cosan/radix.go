package cosan

import (
	"strings"
	"sync"
)

// radixMatcher implements Matcher using a radix tree for efficient path matching.
type radixMatcher struct {
	mu       sync.RWMutex
	trees    map[string]*radixNode // One tree per HTTP method
	compiled bool
}

// radixNode represents a node in the radix tree.
type radixNode struct {
	path      string       // Path segment
	nType     nodeType     // Node type
	paramName string       // Parameter name (for param/wildcard nodes)
	route     *route       // Handler route at this node
	children  []*radixNode // Child nodes
	wildcard  *radixNode   // Wildcard child
	priority  int          // Priority for sorting
}

// nodeType represents the type of radix tree node.
type nodeType uint8

const (
	staticNode   nodeType = iota // Static path segment
	paramNode                    // Named parameter (:id)
	wildcardNode                 // Catch-all parameter (*path)
)

// newRadixMatcher creates a new radix tree matcher.
func newRadixMatcher() *radixMatcher {
	return &radixMatcher{
		trees: make(map[string]*radixNode),
	}
}

// Register adds a route to the radix tree (implements Matcher interface).
func (m *radixMatcher) Register(method, pattern string, handler HandlerFunc) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.compiled {
		return ErrRouterAlreadyCompiled
	}

	// Create route wrapper
	r := &route{
		method:  method,
		pattern: pattern,
		handler: handler,
	}

	// Get or create tree for this method
	tree := m.trees[method]
	if tree == nil {
		tree = &radixNode{nType: staticNode}
		m.trees[method] = tree
	}

	// Insert route into tree
	return m.insertRoute(tree, pattern, r)
}

// insertRoute inserts a route into the radix tree.
func (m *radixMatcher) insertRoute(node *radixNode, pattern string, r *route) error {
	// Remove leading slash
	pattern = strings.TrimPrefix(pattern, "/")

	// If pattern is empty, set route at current node
	if pattern == "" {
		if node.route != nil {
			return ErrConflictingRoutes
		}
		node.route = r
		return nil
	}

	// Find next segment
	segment, remaining := splitPath(pattern)

	// Determine segment type
	if strings.HasPrefix(segment, ":") {
		// Named parameter
		paramName := segment[1:]
		return m.insertParam(node, paramName, remaining, r)
	} else if strings.HasPrefix(segment, "*") {
		// Wildcard parameter
		paramName := segment[1:]
		return m.insertWildcard(node, paramName, r)
	} else {
		// Static segment
		return m.insertStatic(node, segment, remaining, r)
	}
}

// insertStatic inserts a static path segment.
func (m *radixMatcher) insertStatic(node *radixNode, segment, remaining string, r *route) error {
	// Look for existing child with matching prefix
	for _, child := range node.children {
		if child.nType == staticNode && strings.HasPrefix(segment, child.path) {
			if len(child.path) == len(segment) {
				// Exact match - continue with remaining
				return m.insertRoute(child, remaining, r)
			}
			// Partial match - continue with rest of segment
			newSegment := segment[len(child.path):]
			return m.insertRoute(child, newSegment+"/"+remaining, r)
		}
	}

	// No matching child - create new node
	newNode := &radixNode{
		path:     segment,
		nType:    staticNode,
		priority: 100, // Static nodes have highest priority
	}
	node.children = append(node.children, newNode)

	return m.insertRoute(newNode, remaining, r)
}

// insertParam inserts a parameter node.
func (m *radixMatcher) insertParam(node *radixNode, paramName, remaining string, r *route) error {
	// Look for existing param node with same name
	for _, child := range node.children {
		if child.nType == paramNode && child.paramName == paramName {
			return m.insertRoute(child, remaining, r)
		}
	}

	// Create new param node
	newNode := &radixNode{
		nType:     paramNode,
		paramName: paramName,
		priority:  50, // Params have medium priority
	}
	node.children = append(node.children, newNode)

	return m.insertRoute(newNode, remaining, r)
}

// insertWildcard inserts a wildcard node.
func (m *radixMatcher) insertWildcard(node *radixNode, paramName string, r *route) error {
	if node.wildcard != nil {
		return ErrConflictingRoutes
	}

	node.wildcard = &radixNode{
		nType:     wildcardNode,
		paramName: paramName,
		route:     r,
		priority:  10, // Wildcards have lowest priority
	}

	return nil
}

// Compile prepares the matcher for use.
func (m *radixMatcher) Compile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.compiled {
		return nil
	}

	// Sort children by priority for each node
	for _, tree := range m.trees {
		sortByPriority(tree)
	}

	m.compiled = true
	return nil
}

// sortByPriority recursively sorts children by priority.
func sortByPriority(node *radixNode) {
	if node == nil {
		return
	}

	// Sort children: static first, then params
	for i := 0; i < len(node.children); i++ {
		for j := i + 1; j < len(node.children); j++ {
			if node.children[j].priority > node.children[i].priority {
				node.children[i], node.children[j] = node.children[j], node.children[i]
			}
		}
	}

	// Recursively sort children
	for _, child := range node.children {
		sortByPriority(child)
	}
}

// Match finds a route matching the given method and path.
func (m *radixMatcher) Match(method, path string) (*Route, map[string]string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tree := m.trees[method]
	if tree == nil {
		return nil, nil, false
	}

	// Remove leading slash
	path = strings.TrimPrefix(path, "/")

	params := make(map[string]string)
	route := search(tree, path, params)

	if route != nil {
		var r Route = route
		return &r, params, true
	}

	return nil, nil, false
}

// search recursively searches for a matching route.
func search(node *radixNode, path string, params map[string]string) *route {
	// If path is empty, return route at this node
	if path == "" {
		return node.route
	}

	// Try static children first
	for _, child := range node.children {
		if child.nType == staticNode {
			if strings.HasPrefix(path, child.path) {
				remaining := path[len(child.path):]
				if remaining == "" || remaining[0] == '/' {
					// Matched - remove leading slash from remaining
					remaining = strings.TrimPrefix(remaining, "/")
					if route := search(child, remaining, params); route != nil {
						return route
					}
				}
			}
		}
	}

	// Try param children
	for _, child := range node.children {
		if child.nType == paramNode {
			// Find next segment
			segment, remaining := splitPath(path)
			if segment != "" {
				// Save param value
				params[child.paramName] = segment
				if route := search(child, remaining, params); route != nil {
					return route
				}
				// Backtrack - remove param
				delete(params, child.paramName)
			}
		}
	}

	// Try wildcard
	if node.wildcard != nil {
		params[node.wildcard.paramName] = path
		return node.wildcard.route
	}

	return nil
}

// splitPath splits a path into the next segment and remaining path.
func splitPath(path string) (segment, remaining string) {
	if path == "" {
		return "", ""
	}

	// Find next slash
	i := strings.Index(path, "/")
	if i == -1 {
		return path, ""
	}

	return path[:i], path[i+1:]
}
