## 🔹 Go Basics

**1. What are the main features of Go?**

- Statically typed, compiled language.
- Simple syntax, easy to read.
- Garbage collection.
- Built-in concurrency support (goroutines, channels).
- Cross-platform compilation.
- Fast compilation and execution.

---

**2. Difference between `var`, `:=`, and `const`?**

- `var` → declares variables, can specify type or let compiler infer.

  ```go
  var x int = 10
  var y = 20
  ```

- `:=` → short declaration, only inside functions.

  ```go
  z := 30
  ```

- `const` → declares constants, must be known at compile time.

  ```go
  const Pi = 3.14
  ```

---

**3. Difference between `array`, `slice`, and `map`?**

- **Array** → fixed size, `[5]int`.
- **Slice** → dynamic, built on arrays, `[]int{1,2,3}`.
- **Map** → key-value store, `map[string]int{"a":1}`.

---

**4. How does Go handle pointers? Does it support pointer arithmetic?**

- Yes, Go supports pointers (`*int`, `&x`).
- No pointer arithmetic (unlike C) → improves safety.

---

**5. Difference between value types and reference types?**

- **Value types** → copied when assigned (int, float, bool, struct).
- **Reference types** → point to underlying data (slices, maps, channels, interfaces, functions).

---

## 🔹 Functions & Structs

**6. Function vs Method in Go?**

- **Function** → not tied to a type.
- **Method** → attached to a type (struct).

  ```go
  type User struct {name string}
  func (u User) Greet() string { return "Hello " + u.name }
  ```

---

**7. Why use pointer receivers in methods?**

- To modify struct fields.
- To avoid copying large structs.
- Example:

  ```go
  func (u *User) UpdateName(n string) { u.name = n }
  ```

---

**8. Embedding vs Inheritance?**

- Go doesn’t support inheritance.
- **Embedding** allows composition:

  ```go
  type Address struct { City string }
  type User struct {
      Name string
      Address // embedded
  }
  ```

---

## 🔹 Interfaces

**9. How do interfaces work in Go?**

- Implementation is implicit → no `implements` keyword.
- Example:

  ```go
  type Reader interface { Read() string }
  type File struct{}
  func (f File) Read() string { return "file data" }
  ```

---

**10. What happens if two interfaces have same method signature?**

- No conflict. As long as a type implements the method, it satisfies both interfaces.

---

**11. What is `interface{}` used for?**

- Empty interface → accepts any type (like `Object` in Java).
- Often used for generic functions or JSON decoding.

---

## 🔹 Concurrency

**12. How do goroutines work?**

- Lightweight threads managed by Go runtime.
- Created with `go func() { ... }()`.

---

**13. What are channels? How to send/receive?**

- Used to communicate between goroutines safely.

  ```go
  ch := make(chan int)
  go func() { ch <- 5 }()
  val := <- ch
  ```

---

**14. Buffered vs Unbuffered channels?**

- **Unbuffered** → sender blocks until receiver is ready.
- **Buffered** → allows sending up to capacity without blocking.

---

**15. How to avoid race conditions in Go?**

- Use **channels** or **sync.Mutex**.
- Example with mutex:

  ```go
  var mu sync.Mutex
  mu.Lock()
  count++
  mu.Unlock()
  ```

---

**16. Explain `select` in Go.**

- Waits on multiple channels.

  ```go
  select {
    case msg := <-ch1: fmt.Println(msg)
    case msg := <-ch2: fmt.Println(msg)
    default: fmt.Println("no data")
  }
  ```

---

**17. Concurrency vs Parallelism?**

- **Concurrency** → structuring code to handle multiple tasks (not necessarily at same time).
- **Parallelism** → running tasks truly at the same time (multi-core CPU).

---

## 🔹 Error Handling & Context

**18. How is error handling done in Go? Why no try-catch?**

- Go uses explicit error returns.
- Keeps error handling simple and visible.

  ```go
  val, err := someFunc()
  if err != nil { return err }
  ```

---

**19. Purpose of `panic` and `recover`?**

- `panic` → stop normal execution (like exception).
- `recover` → regain control inside `defer`.

---

**20. How is `context.Context` used?**

- For cancellation, timeouts, passing request-scoped values.

  ```go
  ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
  defer cancel()
  ```

---

## 🔹 Web & Database

**21. How to build REST API in Go?**

- Using `net/http`:

  ```go
  http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
      fmt.Fprintln(w, "Hello World")
  })
  http.ListenAndServe(":8080", nil)
  ```

- Or frameworks like **Gin**, **Echo**.

---

**22. How to connect to PostgreSQL/MySQL?**

- Using `database/sql` + driver.

  ```go
  db, err := sql.Open("postgres", "user=... password=...")
  ```

---

**23. Difference between `defer db.Close()` and just `db.Close()`?**

- `defer` ensures DB closes at the end of function.
- Without `defer`, it closes immediately → bad if you still need it.

---

## 🔹 Testing

**24. How to write unit tests in Go?**

- File: `xxx_test.go`, use `testing` package.

  ```go
  func TestAdd(t *testing.T) {
      got := Add(2,3)
      if got != 5 { t.Errorf("got %d, want 5", got) }
  }
  ```

---

**25. What are benchmarks in Go testing?**

- Measure performance with `b.N`.

  ```go
  func BenchmarkAdd(b *testing.B) {
      for i := 0; i < b.N; i++ {
          Add(2,3)
      }
  }
  ```

---

**26. How to mock dependencies in Go tests?**

- Define interfaces → inject mocks.
- Example:

  ```go
  type DB interface { GetUser(id int) string }
  type MockDB struct{}
  func (m MockDB) GetUser(id int) string { return "mock user" }
  ```

---

## 🔹 Advanced

**27. What is Go module system?**

- `go mod init` → creates module.
- `go.mod` → dependency list.
- `go.sum` → dependency checksums for integrity.

---

**28. Difference between goroutine leaks and memory leaks?**

- **Goroutine leak** → goroutines waiting forever (e.g., blocked on channel).
- **Memory leak** → memory not freed (less common in Go due to GC).

---

**29. Best practices for production Go services?**

- Use `context` for timeouts/cancellation.
- Handle errors explicitly.
- Avoid global variables.
- Write unit + integration tests.
- Avoid goroutine leaks.
- Use logging + monitoring.

---

**30. How does garbage collection work in Go?**

- Go uses concurrent, non-generational, tricolor mark-sweep GC.
- Frees unused memory automatically.
