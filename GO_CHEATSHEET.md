# Go Syntax Cheat Sheet

## 1. BASICS

### Package Declaration & Imports
```go
package main

import (
    "fmt"
    "math"
    "strings"

    // Aliased import
    str "strings"

    // Blank import (side effects only)
    _ "net/http/pprof"
)
```

### Variables
```go
// Declaration with var
var name string = "Habeeb"
var age int                      // zero value: 0
var active bool                  // zero value: false

// Short declaration (inside functions only)
name := "Habeeb"
age := 25

// Multiple declaration
var x, y, z int = 1, 2, 3
a, b := "hello", 42

// Constants
const Pi = 3.14159
const (
    StatusOK    = 200
    StatusError = 500
)

// iota (auto-incrementing constants)
const (
    Sunday = iota  // 0
    Monday         // 1
    Tuesday        // 2
)
```

### Basic Types
```go
// Numeric
int, int8, int16, int32, int64
uint, uint8, uint16, uint32, uint64
float32, float64
complex64, complex128
byte   // alias for uint8
rune   // alias for int32 (Unicode code point)

// Other
string
bool

// Type conversion (Go has no implicit casting)
i := 42
f := float64(i)
s := string(rune(65))  // "A"
```

### Printing
```go
fmt.Println("hello", name)            // prints with newline
fmt.Printf("name: %s age: %d\n", name, age)  // formatted
fmt.Sprintf("hello %s", name)         // returns string

// Format verbs
// %s  string
// %d  integer
// %f  float         %.2f for 2 decimal places
// %v  any value     default format
// %+v struct        with field names
// %T  type
// %p  pointer
// %t  boolean
```

### Control Flow
```go
// If-else
if x > 10 {
    fmt.Println("big")
} else if x > 5 {
    fmt.Println("medium")
} else {
    fmt.Println("small")
}

// If with init statement
if err := doSomething(); err != nil {
    log.Fatal(err)
}

// For loop (the only loop in Go)
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While-style
for x < 100 {
    x *= 2
}

// Infinite loop
for {
    break // use break to exit
}

// Range loop
for index, value := range slice {
    fmt.Println(index, value)
}
for key, value := range myMap {
    fmt.Println(key, value)
}
for i, char := range "hello" {
    fmt.Println(i, string(char))
}

// Skip index or value with _
for _, value := range slice {
    fmt.Println(value)
}

// Switch
switch day {
case "Monday":
    fmt.Println("start of week")
case "Friday":
    fmt.Println("almost weekend")
default:
    fmt.Println("regular day")
}

// Switch with no condition (cleaner if-else)
switch {
case score >= 90:
    grade = "A"
case score >= 80:
    grade = "B"
default:
    grade = "C"
}

// Type switch
switch v := i.(type) {
case int:
    fmt.Println("int:", v)
case string:
    fmt.Println("string:", v)
}
```

### Functions
```go
// Basic function
func add(a int, b int) int {
    return a + b
}

// Multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Named return values
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return // naked return
}

// Variadic function
func sum(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
// Call: sum(1, 2, 3) or sum(slice...)

// Function as value
fn := func(x int) int { return x * 2 }
result := fn(5)

// Immediately invoked
func() {
    fmt.Println("executed immediately")
}()
```

---

## 2. DATA STRUCTURES

### Arrays (fixed size)
```go
var arr [5]int                  // [0 0 0 0 0]
arr := [3]string{"a", "b", "c"}
arr := [...]int{1, 2, 3}       // size inferred
len(arr)                        // length
```

### Slices (dynamic size)
```go
// Create
s := []int{1, 2, 3}
s := make([]int, 5)       // length 5, cap 5
s := make([]int, 0, 10)   // length 0, cap 10

// Append
s = append(s, 4)
s = append(s, 5, 6, 7)
s = append(s, otherSlice...)

// Slice from slice
sub := s[1:3]   // index 1 to 2 (exclusive end)
sub := s[:3]    // first 3 elements
sub := s[2:]    // from index 2 to end

// Length and capacity
len(s)  // number of elements
cap(s)  // underlying array capacity

// Copy
dst := make([]int, len(src))
copy(dst, src)

// Delete element at index i
s = append(s[:i], s[i+1:]...)

// Nil slice vs empty slice
var s []int         // nil, len 0, cap 0
s := []int{}        // not nil, len 0, cap 0
s := make([]int, 0) // not nil, len 0, cap 0
```

### Maps
```go
// Create
m := map[string]int{"a": 1, "b": 2}
m := make(map[string]int)

// Set
m["key"] = 42

// Get (returns zero value if missing)
val := m["key"]

// Check existence
val, ok := m["key"]
if !ok {
    fmt.Println("key not found")
}

// Delete
delete(m, "key")

// Iterate
for key, value := range m {
    fmt.Println(key, value)
}

// Length
len(m)
```

### Structs
```go
// Definition
type User struct {
    ID    string
    Name  string
    Email string
    Age   int
}

// Create
u := User{ID: "1", Name: "Habeeb", Email: "h@e.com", Age: 25}
u := User{}              // zero values for all fields
p := &User{Name: "Ali"}  // pointer to struct

// Access fields
u.Name = "Updated"
fmt.Println(u.Name)

// Anonymous struct
point := struct {
    X, Y int
}{10, 20}

// Struct embedding (composition, not inheritance)
type Admin struct {
    User            // embedded — Admin "inherits" User's fields
    Role string
}
a := Admin{User: User{Name: "Habeeb"}, Role: "super"}
a.Name  // promoted field — accessed directly

// Struct tags
type User struct {
    Name  string `json:"name" validate:"required"`
    Email string `json:"email,omitempty"`
    Age   int    `json:"-"`  // excluded from JSON
}
```

### Pointers
```go
x := 42
p := &x         // p is a pointer to x
fmt.Println(*p) // dereference: 42
*p = 100        // x is now 100

// Pointer to struct
u := &User{Name: "Habeeb"}
u.Name  // Go auto-dereferences: same as (*u).Name

// new() allocates and returns pointer
p := new(int)   // *int, points to 0
```

---

## 3. METHODS & INTERFACES

### Methods
```go
// Value receiver (works on a copy)
func (u User) FullName() string {
    return u.Name
}

// Pointer receiver (modifies the original)
func (u *User) SetName(name string) {
    u.Name = name
}

// Rule: if any method needs a pointer receiver, use pointer receivers for all methods on that type
```

### Interfaces
```go
// Definition — a set of method signatures
type Repository interface {
    Save(user User) error
    FindByID(id string) (User, error)
}

// Implicit implementation — no "implements" keyword
// Any type with matching methods satisfies the interface
type MemoryRepo struct{}

func (m *MemoryRepo) Save(user User) error          { return nil }
func (m *MemoryRepo) FindByID(id string) (User, error) { return User{}, nil }

// Usage
var repo Repository = &MemoryRepo{}

// Empty interface (accepts any type)
var anything interface{}
anything = 42
anything = "hello"

// Preferred in Go 1.18+:
var anything any

// Type assertion
s := anything.(string)          // panics if not string
s, ok := anything.(string)      // safe — ok is false if not string

// Interface composition
type ReadWriter interface {
    Reader
    Writer
}
```

---

## 4. ERROR HANDLING

```go
// Errors are values
func doWork() error {
    return errors.New("something failed")
}

// fmt.Errorf with formatting
return fmt.Errorf("user %s not found", id)

// Error wrapping (Go 1.13+)
return fmt.Errorf("database query failed: %w", err)

// Check errors
if err != nil {
    return err
}

// Custom error types
type NotFoundError struct {
    ID string
}
func (e *NotFoundError) Error() string {
    return fmt.Sprintf("resource %s not found", e.ID)
}

// errors.Is — check if error matches a value
if errors.Is(err, ErrUserNotFound) {
    // handle not found
}

// errors.As — check if error matches a type
var nfe *NotFoundError
if errors.As(err, &nfe) {
    fmt.Println("missing ID:", nfe.ID)
}

// Sentinel errors
var (
    ErrNotFound  = errors.New("not found")
    ErrForbidden = errors.New("forbidden")
)

// Panic and recover (use sparingly — not for normal error handling)
func risky() {
    panic("something terrible")
}

func safe() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered:", r)
        }
    }()
    risky()
}
```

---

## 5. CONCURRENCY

### Goroutines
```go
// Launch a goroutine
go doWork()

// Anonymous goroutine
go func() {
    fmt.Println("running concurrently")
}()

// Goroutines with WaitGroup
var wg sync.WaitGroup

for i := 0; i < 5; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        fmt.Println("worker", id)
    }(i)  // pass i to avoid closure capture bug
}
wg.Wait()
```

### Channels
```go
// Unbuffered channel (blocks until both sides are ready)
ch := make(chan int)

// Buffered channel (blocks when buffer is full)
ch := make(chan int, 10)

// Send and receive
ch <- 42        // send
val := <-ch     // receive

// Directional channels (for function parameters)
func producer(out chan<- int)  {}  // send only
func consumer(in <-chan int)   {}  // receive only

// Close a channel
close(ch)

// Range over channel (loops until closed)
for val := range ch {
    fmt.Println(val)
}

// Check if channel is closed
val, ok := <-ch
if !ok {
    fmt.Println("channel closed")
}
```

### Select
```go
select {
case msg := <-ch1:
    fmt.Println("from ch1:", msg)
case msg := <-ch2:
    fmt.Println("from ch2:", msg)
case ch3 <- 42:
    fmt.Println("sent to ch3")
default:
    fmt.Println("no channel ready")
}

// Timeout pattern
select {
case result := <-ch:
    fmt.Println(result)
case <-time.After(5 * time.Second):
    fmt.Println("timed out")
}
```

### Sync Primitives
```go
// Mutex — mutual exclusion lock
var mu sync.Mutex
mu.Lock()
// critical section
mu.Unlock()

// RWMutex — multiple readers OR one writer
var rw sync.RWMutex
rw.RLock()    // read lock (shared)
rw.RUnlock()
rw.Lock()     // write lock (exclusive)
rw.Unlock()

// Once — execute exactly once
var once sync.Once
once.Do(func() {
    fmt.Println("only runs once")
})

// Map — concurrent-safe map
var m sync.Map
m.Store("key", "value")
val, ok := m.Load("key")
m.Delete("key")
m.Range(func(key, value any) bool {
    fmt.Println(key, value)
    return true // continue iteration
})
```

### Context
```go
// Background context (top-level, never cancelled)
ctx := context.Background()

// With cancel
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// With deadline
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
defer cancel()

// With value (use sparingly)
ctx = context.WithValue(ctx, "userID", "abc-123")
val := ctx.Value("userID").(string)

// Check if context is done
select {
case <-ctx.Done():
    fmt.Println("cancelled:", ctx.Err())
default:
    // still running
}
```

---

## 6. STANDARD LIBRARY ESSENTIALS

### Strings
```go
strings.Contains("hello", "ell")         // true
strings.HasPrefix("hello", "hel")        // true
strings.HasSuffix("hello", "llo")        // true
strings.ToUpper("hello")                 // "HELLO"
strings.ToLower("HELLO")                 // "hello"
strings.TrimSpace("  hello  ")           // "hello"
strings.TrimPrefix("/users/123", "/users/")  // "123"
strings.Split("a,b,c", ",")             // ["a", "b", "c"]
strings.Join([]string{"a", "b"}, "-")   // "a-b"
strings.Replace("aaa", "a", "b", 2)     // "bba"
strings.ReplaceAll("aaa", "a", "b")     // "bbb"
strings.Index("hello", "ll")            // 2
```

### String Conversion
```go
import "strconv"

// Int <-> String
strconv.Itoa(42)              // "42"
strconv.Atoi("42")            // 42, error

// Parse
strconv.ParseBool("true")     // true, error
strconv.ParseFloat("3.14", 64) // 3.14, error
strconv.ParseInt("42", 10, 64) // 42, error

// Format
strconv.FormatBool(true)       // "true"
strconv.FormatFloat(3.14, 'f', 2, 64) // "3.14"
```

### Time
```go
import "time"

now := time.Now()
t := time.Date(2026, time.June, 1, 12, 0, 0, 0, time.UTC)

// Duration
d := 5 * time.Second
d := 100 * time.Millisecond
time.Sleep(d)

// Since / Until
elapsed := time.Since(start)     // time.Now() - start
remaining := time.Until(deadline) // deadline - time.Now()

// Format (Go uses reference time: Mon Jan 2 15:04:05 MST 2006)
now.Format("2006-01-02")             // "2026-06-01"
now.Format("2006-01-02 15:04:05")    // "2026-06-01 12:00:00"
now.Format(time.RFC3339)             // "2026-06-01T12:00:00Z"

// Parse
t, err := time.Parse("2006-01-02", "2026-06-01")
t, err := time.Parse(time.RFC3339, "2026-06-01T12:00:00Z")

// Comparison
t1.Before(t2)
t1.After(t2)
t1.Equal(t2)
t1.Add(24 * time.Hour)
t1.Sub(t2)  // returns Duration
```

### JSON
```go
import "encoding/json"

// Marshal (struct -> JSON bytes)
data, err := json.Marshal(user)
fmt.Println(string(data))

// MarshalIndent (pretty print)
data, err := json.MarshalIndent(user, "", "  ")

// Unmarshal (JSON bytes -> struct)
var user User
err := json.Unmarshal(data, &user)

// Encoder (write JSON to io.Writer)
json.NewEncoder(w).Encode(user)

// Decoder (read JSON from io.Reader)
var user User
json.NewDecoder(r.Body).Decode(&user)

// Struct tags control JSON behavior
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email,omitempty"` // omit if empty
    Pass  string `json:"-"`              // always omit
}
```

### File I/O
```go
import "os"

// Read entire file
data, err := os.ReadFile("file.txt")

// Write entire file
err := os.WriteFile("file.txt", []byte("content"), 0644)

// Open for reading
f, err := os.Open("file.txt")
defer f.Close()

// Create / truncate
f, err := os.Create("file.txt")
defer f.Close()
f.WriteString("hello")

// Open with flags
f, err := os.OpenFile("file.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
defer f.Close()

// Check if file exists
_, err := os.Stat("file.txt")
if os.IsNotExist(err) {
    fmt.Println("file does not exist")
}
```

---

## 7. NET/HTTP

### HTTP Server
```go
// Handler function
func hello(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"msg": "hello"})
}

// Register and serve
http.HandleFunc("/hello", hello)
log.Fatal(http.ListenAndServe(":8080", nil))

// With http.Server (for graceful shutdown)
srv := &http.Server{
    Addr:         ":8080",
    Handler:      router,
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  60 * time.Second,
}
log.Fatal(srv.ListenAndServe())
```

### HTTP Handler Interface
```go
// Any type implementing this is an http.Handler
type Handler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// http.HandlerFunc adapts a function into a Handler
var h http.Handler = http.HandlerFunc(myFunc)
```

### Middleware Pattern
```go
func MyMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // before
        next.ServeHTTP(w, r)
        // after
    })
}
```

### HTTP Client
```go
// Simple GET
resp, err := http.Get("https://api.example.com/data")
defer resp.Body.Close()
body, err := io.ReadAll(resp.Body)

// Custom request
req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer token")

client := &http.Client{Timeout: 10 * time.Second}
resp, err := client.Do(req)

// With context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
req = req.WithContext(ctx)
```

### httptest (Testing)
```go
import "net/http/httptest"

// Create fake request
req := httptest.NewRequest(http.MethodGet, "/users", nil)
req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"test"}`))

// Capture response
rec := httptest.NewRecorder()

// Call handler directly
handler.ServeHTTP(rec, req)

// Assert
fmt.Println(rec.Code)         // status code
fmt.Println(rec.Body.String()) // response body
```

---

## 8. TESTING

### Basic Test
```go
// file: math_test.go
package math

import "testing"

func TestAdd(t *testing.T) {
    result := Add(2, 3)
    if result != 5 {
        t.Errorf("expected 5, got %d", result)
    }
}
```

### Table-Driven Tests
```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("expected %d, got %d", tt.expected, result)
            }
        })
    }
}
```

### Testify Assertions
```go
import "github.com/stretchr/testify/assert"

assert.Equal(t, expected, actual)
assert.NotEqual(t, a, b)
assert.NoError(t, err)
assert.Error(t, err)
assert.Nil(t, obj)
assert.NotNil(t, obj)
assert.Empty(t, slice)
assert.NotEmpty(t, slice)
assert.Contains(t, "hello world", "hello")
assert.True(t, condition)
assert.False(t, condition)
```

### Testify Mocks
```go
import "github.com/stretchr/testify/mock"

type MockRepo struct {
    mock.Mock
}

func (m *MockRepo) Save(user User) error {
    args := m.Called(user)
    return args.Error(0)
}

func (m *MockRepo) FindByID(id string) (User, error) {
    args := m.Called(id)
    return args.Get(0).(User), args.Error(1)
}

// In test
mockRepo := new(MockRepo)
mockRepo.On("Save", mock.Anything).Return(nil)
mockRepo.On("FindByID", "123").Return(User{ID: "123"}, nil)
mockRepo.AssertExpectations(t)
```

### Benchmarks
```go
func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
// Run: go test -bench=. -benchmem
```

### Test Commands
```bash
go test ./...              # run all tests
go test -v ./...           # verbose output
go test -run TestAdd ./... # run specific test
go test -cover ./...       # show coverage
go test -coverprofile=c.out ./... && go tool cover -html=c.out  # HTML report
go test -race ./...        # detect race conditions
go test -bench=. ./...     # run benchmarks
go test -count=1 ./...     # disable test caching
```

---

## 9. GENERICS (Go 1.18+)

```go
// Generic function
func Map[T any, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Usage
doubled := Map([]int{1, 2, 3}, func(n int) int { return n * 2 })

// Type constraints
type Number interface {
    ~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Sum[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}

// Generic struct
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

// Comparable constraint (supports == and !=)
func Contains[T comparable](slice []T, target T) bool {
    for _, v := range slice {
        if v == target {
            return true
        }
    }
    return false
}
```

---

## 10. ADVANCED PATTERNS

### Defer
```go
// Runs when surrounding function returns (LIFO order)
f, _ := os.Open("file.txt")
defer f.Close()  // guaranteed cleanup

// Defer evaluates arguments immediately
x := 10
defer fmt.Println(x) // prints 10, not whatever x is later
x = 20

// Stacked defers (LIFO)
defer fmt.Println("first")  // prints third
defer fmt.Println("second") // prints second
defer fmt.Println("third")  // prints first
```

### Closures
```go
func counter() func() int {
    count := 0
    return func() int {
        count++
        return count
    }
}

c := counter()
c() // 1
c() // 2
c() // 3
```

### Embedding Interfaces in Structs
```go
type Logger interface {
    Log(msg string)
}

type Service struct {
    Logger  // embedded interface — Service must be given a Logger at creation
}

// The Service can now call s.Log("msg") directly
```

### Functional Options Pattern
```go
type Server struct {
    port    int
    timeout time.Duration
}

type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) { s.port = port }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(opts ...Option) *Server {
    s := &Server{port: 8080, timeout: 30 * time.Second}  // defaults
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
srv := NewServer(WithPort(9090), WithTimeout(60*time.Second))
```

### Builder Pattern with Method Chaining
```go
type QueryBuilder struct {
    table  string
    wheres []string
    limit  int
}

func NewQuery(table string) *QueryBuilder {
    return &QueryBuilder{table: table}
}

func (q *QueryBuilder) Where(cond string) *QueryBuilder {
    q.wheres = append(q.wheres, cond)
    return q
}

func (q *QueryBuilder) Limit(n int) *QueryBuilder {
    q.limit = n
    return q
}

// Usage
q := NewQuery("users").Where("age > 18").Where("active = true").Limit(10)
```

### Graceful Shutdown Pattern
```go
srv := &http.Server{Addr: ":8080", Handler: router}

go srv.ListenAndServe()

quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

srv.Shutdown(ctx)
```

### Worker Pool
```go
func workerPool(jobs <-chan int, results chan<- int, numWorkers int) {
    var wg sync.WaitGroup
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    wg.Wait()
    close(results)
}
```

### Fan-Out / Fan-In
```go
// Fan-out: multiple goroutines reading from one channel
// Fan-in: merge multiple channels into one

func fanIn(channels ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    merged := make(chan int)

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for val := range c {
                merged <- val
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(merged)
    }()

    return merged
}
```

---

## 11. MODULES & TOOLING

### Module Commands
```bash
go mod init github.com/user/project   # initialize module
go mod tidy                            # add missing, remove unused deps
go get github.com/pkg/errors           # add dependency
go get github.com/pkg/errors@v0.9.1    # specific version
go get -u ./...                        # update all dependencies
go mod vendor                          # copy deps to vendor/
go mod download                        # download deps to cache
```

### Build & Run
```bash
go run main.go              # compile and run
go run ./cmd/server         # run a package
go build -o myapp ./cmd/server  # compile binary
go install ./cmd/server     # compile and install to $GOPATH/bin

# Build for different OS/arch
GOOS=linux GOARCH=amd64 go build -o myapp-linux ./cmd/server

# Build flags
go build -ldflags="-s -w"  # strip debug info (smaller binary)
```

### Code Quality
```bash
go fmt ./...       # format code
go vet ./...       # detect suspicious code
golangci-lint run  # comprehensive linter (install separately)
```

### Documentation
```bash
go doc net/http                    # package docs
go doc net/http.Handler            # type docs
go doc net/http.ListenAndServe     # function docs
go doc -all net/http.ResponseWriter # all methods
go doc -src json.NewEncoder        # source code
```

---

## 12. gRPC

### Proto3 Syntax
```protobuf
syntax = "proto3";

package mypackage;

option go_package = "github.com/user/project/proto";

// Messages are like structs. Field tags (= 1, = 2) are permanent —
// they identify fields in the binary encoding. Never change them.
message User {
    string id    = 1;
    string name  = 2;
    string email = 3;
}

message GetUserRequest {
    string id = 1;
}

message CreateUserRequest {
    string name  = 1;
    string email = 2;
}

// repeated = slice
message ListUsersResponse {
    repeated User users = 1;
}

// map field
message Config {
    map<string, string> settings = 1;
}

// enum
enum Status {
    UNKNOWN = 0;  // proto3 enums must start at 0
    ACTIVE  = 1;
    INACTIVE = 2;
}

// Service defines RPC methods
service UserService {
    rpc GetUser    (GetUserRequest)    returns (User);           // unary
    rpc ListUsers  (ListRequest)       returns (stream User);    // server streaming
    rpc Upload     (stream Chunk)      returns (UploadResponse); // client streaming
    rpc Chat       (stream Message)    returns (stream Message); // bidirectional
}
```

### Generate Go Code
```bash
# Install plugins (once)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate — run from project root
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/user.proto

# Output:
# proto/user.pb.go       — message structs
# proto/user_grpc.pb.go  — service interface + registration
```

### Implement the Server
```go
// The generated interface — you must implement every method
type UserServiceServer interface {
    GetUser(context.Context, *GetUserRequest) (*User, error)
    CreateUser(context.Context, *CreateUserRequest) (*User, error)
    mustEmbedUnimplementedUserServiceServer()
}

// Your implementation — embed Unimplemented* by value (not pointer)
// for forward compatibility when new methods are added to the proto
type UserServer struct {
    proto.UnimplementedUserServiceServer
    users map[string]*proto.User
    mu    sync.RWMutex
}

func (s *UserServer) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.User, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()
    u, ok := s.users[req.Id]
    if !ok {
        return nil, status.Errorf(codes.NotFound, "user not found")
    }
    return u, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.User, error) {
    u := &proto.User{Id: uuid.New().String(), Name: req.Name, Email: req.Email}
    s.mu.Lock()
    defer s.mu.Unlock()
    s.users[u.Id] = u
    return u, nil
}
```

### Start the Server
```go
import "google.golang.org/grpc"
import "google.golang.org/grpc/reflection"

lis, err := net.Listen("tcp", ":50051")
if err != nil {
    log.Fatal(err)
}

grpcServer := grpc.NewServer()
proto.RegisterUserServiceServer(grpcServer, &UserServer{})
reflection.Register(grpcServer)  // enables grpcurl and other dev tools

log.Println("gRPC listening on :50051")
log.Fatal(grpcServer.Serve(lis))
```

### gRPC Status Errors
```go
import (
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// Return structured errors from server methods
return nil, status.Errorf(codes.NotFound, "user %s not found", id)
return nil, status.Errorf(codes.InvalidArgument, "name is required")
return nil, status.Errorf(codes.Internal, "database error")
return nil, status.Errorf(codes.AlreadyExists, "email already taken")
return nil, status.Errorf(codes.Unauthenticated, "missing token")
return nil, status.Errorf(codes.PermissionDenied, "forbidden")

// Common codes
// codes.OK              — success (returned implicitly)
// codes.NotFound        — resource does not exist
// codes.InvalidArgument — bad input from client
// codes.AlreadyExists   — duplicate resource
// codes.Internal        — server-side error
// codes.Unauthenticated — missing or invalid credentials
// codes.PermissionDenied — authenticated but not authorized
// codes.Unavailable     — server temporarily down
```

### Test with grpcurl
```bash
# Install
brew install grpcurl

# List services (requires reflection.Register on server)
grpcurl -plaintext localhost:50051 list

# List methods on a service
grpcurl -plaintext localhost:50051 list user.UserService

# Call a method
grpcurl -plaintext -d '{"name": "Habeeb", "email": "h@example.com"}' \
    localhost:50051 user.UserService/CreateUser

grpcurl -plaintext -d '{"id": "abc-123"}' \
    localhost:50051 user.UserService/GetUser
```

### Install Dependencies
```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
go get github.com/google/uuid  # for ID generation
```

### When to Use gRPC vs REST
```
REST                          gRPC
─────────────────────────     ──────────────────────────
Public APIs (browsers)        Internal service-to-service
JSON / human readable         Binary / high performance
No codegen needed             Contract enforced by proto
Flexible schema               Strongly typed schema
Stateless request/response    Supports streaming
```
