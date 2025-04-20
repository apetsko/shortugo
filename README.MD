# Shortugo

Shortugo is a high-performance URL shortening service written in Go. It provides an API for generating short links, managing user-specific URLs, and interacting with a persistent database.

## 🚀 Features

- URL shortening (plain text, JSON, batch)
- Retrieve all user URLs
- Delete user URLs
- Expand shortened URLs to original
- Health check endpoint for database connectivity

## 📋 Endpoints

| Method   | Path                      | Description                             |
|----------|---------------------------|-----------------------------------------|
| `POST`   | `/`                       | Shorten URL (plain text)                |
| `POST`   | `/api/shorten`            | Shorten URL (JSON)                      |
| `POST`   | `/api/shorten/batch`      | Batch URL shortening                    |
| `GET`    | `/api/user/urls`          | Retrieve user's URLs                    |
| `DELETE` | `/api/user/urls`          | Delete user's URLs                      |
| `GET`    | `/{id}`                   | Expand shortened URL                    |
| `GET`    | `/ping`                   | Check database connectivity             |

## ⚙️ Middleware

- `RealIP` — extracts the real client IP
- `Recoverer` — handles panics and returns 500 errors
- `LogMiddleware` — logs requests and responses
- `GzipMiddleware` — compresses responses using gzip

> ❌ The `RequestID` middleware was removed as part of performance optimization.

---

## ⚡ Optimizations

The following changes were made to improve performance:

- 🔥 Removed `RequestID` middleware
- ⚙️ Replaced `json.Marshal` with `json.NewEncoder(&buf)` to reduce memory allocations
- ⚙ Replaced `fmt.Sprintf` with string concatenation using `+` to reduce CPU usage.
- 🔧 Modified `GzipMiddleware` to use a `sync.Pool` for `gzip.Reader`, improving performance and reducing GC pressure.

```go
var gzipReaderPool = sync.Pool{
	New: func() any {
		return new(gzip.Reader)
	},
}
```
---

## 📈 Profiling

### 📊 Heap Profiling

```bash
go tool pprof -top -diff_base=./profile/list_user_urls_heap_base.pprof ./profile/list_user_urls_heap_result.pprof
```

**Result Summary:**
```
File: wgo_20250406124054_1464
Build ID: 6137533df8ba87f1f2ecb5a5bfe7865bc8fd4cbf
Type: inuse_space
Time: 2025-04-06 12:30:00 +05
Duration: 90.04s, Total samples = 111.84MB 
Showing nodes accounting for 12.24MB, 10.94% of 111.84MB total
      flat  flat%   sum%        cum   cum%
  -22.47MB 20.09% 20.09%    ...
```

### 🧠 CPU Profiling

```bash
go tool pprof -top -diff_base=./profile/list_user_urls_profile_base.pprof ./profile/list_user_urls_profile_result.pprof
```

**Result Summary:**
```
File: wgo_20250406124054_1464
Build ID: 87d19bc98f8b5040a8724bf0462b6d91e5a17db5
Type: cpu
Time: 2025-04-06 12:21:18 +05
Duration: 60.35s, Total samples = 230.69s (382.26%)
Showing nodes accounting for 0.68s, 0.29% of 230.69s total
Dropped 13 nodes (cum <= 1.15s)
      flat  flat%   sum%        cum   cum%
     4.01s  1.74%  1.74%     12.95s  5.61%  runtime.concatstrings
    -3.89s  1.69%  0.052%   -16.02s  6.94%  fmt.(*pp).doPrintf
     3.80s  1.65%  1.70%      4.05s  1.76%  encoding/json.appendString[go.shape.string]
     3.54s  1.53%  3.23%      3.54s  1.53%  internal/runtime/syscall.Syscall6
     3.18s  1.38%  4.61%      3.18s  1.38%  runtime.memclrNoHeapPointers
    -2.37s  1.03%  3.58%    -10.48s  4.54%  fmt.(*pp).printArg
    -2.11s  0.91%  2.67%     -8.29s  3.59%  runtime.mallocgcSmallScanNoHeader
     2.01s  0.87%  3.54%      2.01s  0.87%  runtime.memmove
    -1.80s  0.78%  2.76%     -1.80s  0.78%  runtime.nextFreeFast (inline)
    -1.73s  0.75%  2.01%     -4....
```

These improvements led to significant reductions in memory usage and CPU time, especially around `fmt` and `encoding/json` operations.

---

## 🧪 Usage

Use `curl`, Postman, or integrate with your frontend/backend services to shorten and manage URLs using the API above.

## 📄 License

MIT

