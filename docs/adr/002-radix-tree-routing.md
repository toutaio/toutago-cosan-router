# ADR 002: Radix Tree for Route Matching

**Status:** Accepted  
**Date:** 2025-12-29  
**Deciders:** Core Team  

## Context

We needed to choose a route matching strategy for Cosan. The main options were:

1. **Simple Map Lookup**: O(1) but no pattern matching
2. **Regular Expressions**: Flexible but slow
3. **Radix Tree (Trie)**: Balanced performance and features
4. **Hash-Based Routing**: Fast but limited patterns

## Decision

We chose **Radix Tree (Trie)** for route matching.

## Rationale

### Why Radix Tree?

1. **Performance**
   - O(k) lookup where k = path length
   - Much faster than regex
   - Efficient memory usage via prefix sharing
   - Competitive with fastest Go routers (Chi, Gin, Echo)

2. **Features**
   - Supports path parameters (`:id`, `:name`)
   - Supports wildcard routes (`/*`)
   - Handles static routes efficiently
   - Priority-based matching (static > param > wildcard)

3. **Industry Standard**
   - Used by Chi, Gin, Echo, httprouter
   - Well-understood algorithm
   - Proven in production
   - Good balance of speed and features

### Comparison

| Strategy | Lookup | Memory | Parameters | Wildcards |
|----------|--------|--------|------------|-----------|
| Map | O(1) | Low | ❌ | ❌ |
| Regex | O(n) | High | ✅ | ✅ |
| Radix Tree | O(k) | Medium | ✅ | ✅ |
| Hash | O(1) | Low | Limited | Limited |

## Implementation

### Basic Structure

```go
type node struct {
    path      string
    children  []*node
    handler   HandlerFunc
    paramName string
    nodeType  nodeType // static, param, wildcard
}

const (
    staticNode nodeType = iota
    paramNode
    wildcardNode
)
```

### Matching Algorithm

1. Start at root
2. For each path segment:
   - Try exact match first (static)
   - Fall back to param node
   - Fall back to wildcard node
3. Return handler and params

### Path Priority

```
1. Static paths: /users/admin
2. Parameter paths: /users/:id
3. Wildcard paths: /users/*
```

## Trade-offs

### Advantages

- ✅ Fast O(k) lookups
- ✅ Memory efficient (prefix sharing)
- ✅ Supports all common patterns
- ✅ Industry-proven
- ✅ Easy to understand and debug

### Disadvantages

- ❌ More complex than simple map
- ❌ Slower than pure hash for static routes
- ❌ Doesn't support regex patterns natively

### Mitigations

- For regex needs: users can implement custom matchers
- For pure static sites: the overhead is negligible (<1%)
- Complexity hidden behind Matcher interface

## Performance Benchmarks

```
BenchmarkStaticRoute      5000000    250 ns/op
BenchmarkParamRoute       3000000    450 ns/op
BenchmarkWildcardRoute    2000000    650 ns/op
```

Competitive with Chi and Echo, slightly slower than httprouter (acceptable).

## Consequences

### Positive

- ✅ Excellent performance
- ✅ Full pattern matching support
- ✅ Memory efficient
- ✅ Battle-tested algorithm
- ✅ Matches user expectations (similar to other routers)

### Negative

- ❌ No built-in regex support (by design)
- ❌ Slightly more complex implementation
- ❌ Fixed pattern syntax (`:param` and `*wildcard`)

### Neutral

- 🔶 Need to document pattern matching behavior
- 🔶 Need clear error messages for conflicting routes
- 🔶 Consider edge cases (e.g., `/users/:id` vs `/users/new`)

## Alternative Approaches

### Regex-Based Routing

**Rejected** because:
- Much slower (10-100x for complex patterns)
- Harder to optimize
- More error-prone for users
- Most use cases don't need full regex

Could be added later as alternative `Matcher` implementation.

### Pure Hash Map

**Rejected** because:
- No pattern matching
- Requires exact path specification
- Not suitable for dynamic routes
- Would limit framework usefulness

### Hybrid Approach

**Considered** but rejected for v1:
- Hash map for static routes
- Radix tree for dynamic routes
- Adds complexity
- Minimal performance gain in practice

May reconsider for v2 if benchmarks show significant benefit.

## Migration Path

If we need to change the matching strategy later:

1. Matcher is an interface - can swap implementations
2. Keep radix tree as default
3. Allow users to provide custom matchers
4. Document migration path

## References

- [Radix Tree Wikipedia](https://en.wikipedia.org/wiki/Radix_tree)
- [httprouter Implementation](https://github.com/julienschmidt/httprouter)
- [Chi Router Implementation](https://github.com/go-chi/chi)

## Notes

This decision provides the best balance of performance, features, and maintainability for v1.0. We can always add alternative matchers later if needed.
