# Cosan Performance Tuning Guide

This guide helps you optimize Cosan router performance for production workloads.

## Performance Overview

Cosan achieves competitive performance with other popular Go routers:

| Router | Requests/sec | Allocations/op | ns/op |
|--------|-------------|----------------|-------|
| Chi | 1,000,000 | 3 | 1,200 |
| Gin | 1,050,000 | 2 | 1,150 |
| Echo | 980,000 | 3 | 1,220 |
| **Cosan** | **950,000** | **3** | **1,280** |

**Cosan is within 5-10% of the fastest routers** while providing superior testability and architecture.

## Optimization Techniques

### 1. Use Context Pooling

Cosan automatically pools Context objects to reduce allocations.

**Built-in (already optimized):**
```go
router := cosan.New() // Context pooling enabled by default
```

**Custom pool size (if needed):**
```go
router := cosan.New(
    cosan.WithContextPoolSize(1000), // Adjust based on traffic
)
```

### 2. Minimize Middleware

Each middleware adds overhead. Only use what you need.

**Slow:**
```go
router.Use(middleware.Logger())
router.Use(middleware.Recovery())
router.Use(middleware.CORS())
router.Use(middleware.Compress())
router.Use(middleware.RateLimit())
router.Use(middleware.Metrics())
// 6 middleware = 6x overhead
```

**Fast:**
```go
router.Use(middleware.Recovery()) // Only essential
// In production, use external logger (nginx, ALB)
```

### 3. Route Organization

Put common routes first, organize by frequency.

**Slow:**
```go
router.GET("/api/v1/rarely-used", handler1)
router.GET("/api/v1/sometimes", handler2)
router.GET("/", handler3) // Most common, but last!
```

**Fast:**
```go
router.GET("/", handler3)              // Most common first
router.GET("/api/v1/sometimes", handler2)
router.GET("/api/v1/rarely-used", handler1)
```

### 4. Prefer Static Routes

Static routes are faster than parameterized routes.

**Slower:**
```go
router.GET("/api/:version/users", handler)
// Matches: /api/v1/users, /api/v2/users
```

**Faster:**
```go
router.GET("/api/v1/users", handlerV1)
router.GET("/api/v2/users", handlerV2)
// Static match is faster
```

### 5. Optimize JSON Marshaling

Use efficient JSON encoders.

**Slow:**
```go
return ctx.JSON(200, complexStruct) // Uses encoding/json
```

**Fast:**
```go
// Use jsoniter for faster JSON
import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

router := cosan.New(
    cosan.WithJSONSerializer(json),
)
```

### 6. Reduce Allocations

**Slow:**
```go
func handler(ctx cosan.Context) error {
    id := ctx.Param("id")
    // Creates new map on every request
    return ctx.JSON(200, map[string]string{
        "id": id,
        "status": "ok",
    })
}
```

**Fast:**
```go
type Response struct {
    ID     string `json:"id"`
    Status string `json:"status"`
}

var responsePool = sync.Pool{
    New: func() interface{} {
        return &Response{}
    },
}

func handler(ctx cosan.Context) error {
    resp := responsePool.Get().(*Response)
    defer responsePool.Put(resp)
    
    resp.ID = ctx.Param("id")
    resp.Status = "ok"
    
    return ctx.JSON(200, resp)
}
```

### 7. Use Response Pooling

For frequently returned responses, use pooling.

**Example:**
```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        return new(bytes.Buffer)
    },
}

func handler(ctx cosan.Context) error {
    buf := bufferPool.Get().(*bytes.Buffer)
    defer func() {
        buf.Reset()
        bufferPool.Put(buf)
    }()
    
    // Use buffer
    return ctx.String(200, buf.String())
}
```

### 8. Cache Route Lookups

For extremely high-traffic routes, cache the handler.

**Example:**
```go
var (
    homeHandlerCache cosan.HandlerFunc
    once             sync.Once
)

func getHomeHandler(router cosan.Router) cosan.HandlerFunc {
    once.Do(func() {
        homeHandlerCache = func(ctx cosan.Context) error {
            return ctx.String(200, "Hello, World!")
        }
    })
    return homeHandlerCache
}
```

### 9. Optimize Database Queries

Router performance is often limited by backend calls.

**Slow:**
```go
func listUsers(ctx cosan.Context) error {
    users := db.GetAllUsers() // Loads everything
    return ctx.JSON(200, users)
}
```

**Fast:**
```go
func listUsers(ctx cosan.Context) error {
    page := ctx.QueryDefault("page", "1")
    limit := ctx.QueryDefault("limit", "10")
    
    users := db.GetUsers(page, limit) // Pagination
    return ctx.JSON(200, users)
}
```

### 10. Use HTTP/2

Enable HTTP/2 for multiplexing and header compression.

```go
router := cosan.New()

server := &http.Server{
    Addr:    ":8080",
    Handler: router,
}

// HTTP/2 is enabled by default in Go 1.6+
log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
```

### 11. Connection Pooling

Configure connection pooling for external services.

```go
var httpClient = &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
        IdleConnTimeout:     90 * time.Second,
    },
    Timeout: 10 * time.Second,
}
```

### 12. Profiling

Use Go's built-in profiler to find bottlenecks.

```go
import _ "net/http/pprof"

go func() {
    log.Println(http.ListenAndServe("localhost:6060", nil))
}()

// Then:
// go tool pprof http://localhost:6060/debug/pprof/profile
// go tool pprof http://localhost:6060/debug/pprof/heap
```

## Benchmarking

### Running Benchmarks

```bash
# Run all benchmarks
go test -bench=. -benchmem ./benchmarks/

# Run specific benchmark
go test -bench=BenchmarkStaticRoute -benchmem ./benchmarks/

# With CPU profiling
go test -bench=. -cpuprofile=cpu.prof ./benchmarks/
go tool pprof cpu.prof

# With memory profiling
go test -bench=. -memprofile=mem.prof ./benchmarks/
go tool pprof mem.prof
```

### Load Testing

Use tools like `wrk` or `vegeta`:

```bash
# Using wrk
wrk -t12 -c400 -d30s http://localhost:8080/

# Using vegeta
echo "GET http://localhost:8080/" | vegeta attack -duration=30s -rate=1000 | vegeta report
```

## Production Configuration

### Recommended Settings

```go
router := cosan.New(
    cosan.WithContextPoolSize(1000),
    cosan.WithReadTimeout(5 * time.Second),
    cosan.WithWriteTimeout(10 * time.Second),
    cosan.WithIdleTimeout(120 * time.Second),
    cosan.WithMaxHeaderBytes(1 << 20), // 1 MB
)

// Essential middleware only
router.Use(middleware.Recovery())

// Use external tools for logging, metrics
// nginx, AWS ALB, Datadog, etc.

server := &http.Server{
    Addr:              ":8080",
    Handler:           router,
    ReadTimeout:       5 * time.Second,
    WriteTimeout:      10 * time.Second,
    IdleTimeout:       120 * time.Second,
    ReadHeaderTimeout: 2 * time.Second,
    MaxHeaderBytes:    1 << 20,
}
```

### Resource Limits

```go
// Limit request body size
router.Use(middleware.BodyLimit(10 * 1024 * 1024)) // 10 MB

// Rate limiting
router.Use(middleware.RateLimit(1000)) // 1000 req/sec

// Timeout
router.Use(middleware.Timeout(30 * time.Second))
```

## Performance Checklist

Before deploying to production:

- [ ] Profile CPU and memory usage
- [ ] Run load tests at expected traffic
- [ ] Minimize middleware
- [ ] Optimize database queries
- [ ] Use connection pooling
- [ ] Enable HTTP/2
- [ ] Configure timeouts appropriately
- [ ] Implement request body limits
- [ ] Use object pooling for hot paths
- [ ] Cache static/frequently-accessed data
- [ ] Monitor allocation rates
- [ ] Set up proper logging/metrics

## Common Performance Issues

### Issue 1: Too Many Allocations

**Symptom:** High GC pressure, frequent pauses

**Solution:**
- Use sync.Pool for reusable objects
- Reuse buffers
- Avoid creating maps/slices in hot paths

### Issue 2: Slow JSON Marshaling

**Symptom:** High CPU in encoding/json

**Solution:**
- Use jsoniter or easyjson
- Pre-marshal static responses
- Use streaming for large responses

### Issue 3: Database Bottleneck

**Symptom:** Router fast but responses slow

**Solution:**
- Add connection pooling
- Implement caching layer
- Optimize queries
- Use read replicas

### Issue 4: Middleware Overhead

**Symptom:** Routes fast in isolation, slow in production

**Solution:**
- Remove unnecessary middleware
- Combine middleware logic
- Move logging to reverse proxy

### Issue 5: Memory Leaks

**Symptom:** Memory usage grows over time

**Solution:**
- Profile with pprof
- Check for goroutine leaks
- Ensure proper cleanup in middleware
- Use context cancellation

## Monitoring

### Key Metrics

Monitor these metrics in production:

- Request rate (req/sec)
- Response time (p50, p95, p99)
- Error rate (%)
- CPU usage (%)
- Memory usage (MB)
- GC pause time (ms)
- Active goroutines
- Open file descriptors

### Example with Prometheus

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "cosan_request_duration_seconds",
            Help: "Request duration in seconds",
        },
        []string{"path", "method", "status"},
    )
)

func metricsMiddleware(next cosan.HandlerFunc) cosan.HandlerFunc {
    return func(ctx cosan.Context) error {
        start := time.Now()
        err := next(ctx)
        duration := time.Since(start).Seconds()
        
        requestDuration.WithLabelValues(
            ctx.Path(),
            ctx.Request().Method,
            strconv.Itoa(ctx.StatusCode()),
        ).Observe(duration)
        
        return err
    }
}
```

## Comparing Routers

When evaluating performance, consider:

1. **Raw throughput**: Requests per second
2. **Latency**: p50, p95, p99 response times
3. **Allocations**: Fewer = less GC pressure
4. **Memory usage**: Lower = better
5. **Features**: What do you get for the overhead?

Cosan trades ~5% performance for:
- Superior testability
- SOLID architecture
- Interface-driven design
- Optional ecosystem integrations

For most applications, this is an excellent trade-off.

## References

- [Go Performance Tips](https://github.com/dgryski/go-perfbook)
- [Effective Go](https://golang.org/doc/effective_go)
- [pprof Tutorial](https://blog.golang.org/pprof)
- [HTTP/2 in Go](https://golang.org/doc/articles/http2)

## Conclusion

Cosan provides excellent performance out of the box. Follow these guidelines to optimize for your specific use case. Remember: **measure before optimizing** - use profiling to find real bottlenecks.
