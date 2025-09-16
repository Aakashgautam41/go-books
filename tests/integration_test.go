package main

import (
    "bytes"
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"
    "time"

    tc "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    _ "github.com/lib/pq"
)

func TestIntegration_CRUD(t *testing.T) {
    ctx := context.Background()

    req := tc.ContainerRequest{
        Image:        "postgres:15",
        Env:          map[string]string{"POSTGRES_PASSWORD": "postgres", "POSTGRES_USER": "postgres", "POSTGRES_DB": "testdb"},
        ExposedPorts: []string{"5432/tcp"},
        WaitingFor:   wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
    }
    postgresC, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatalf("failed to start container: %v", err)
    }
    defer func() {
        _ = postgresC.Terminate(ctx)
    }()

    host, err := postgresC.Host(ctx)
    if err != nil {
        t.Fatalf("host: %v", err)
    }
    port, err := postgresC.MappedPort(ctx, "5432")
    if err != nil {
        t.Fatalf("mapped port: %v", err)
    }

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        host, port.Port(), "postgres", "postgres", "testdb")

    // wait and retry for DB readiness (sometimes wait strategy is not fully enough)
    var db *sql.DB
    for i := 0; i < 20; i++ {
        db, err = sql.Open("postgres", dsn)
        if err == nil {
            err = db.Ping()
            if err == nil {
                break
            }
        }
        time.Sleep(500 * time.Millisecond)
    }
    if err != nil {
        t.Fatalf("could not connect to postgres: %v", err)
    }
    defer db.Close()

    store := NewStore(db)
    // run migration
    migration, err := os.ReadFile("migrations/001_create_books_table.sql")
    if err != nil {
        t.Fatalf("could not read migration: %v", err)
    }
    if err := store.Exec(string(migration)); err != nil {
        t.Fatalf("migration failed: %v", err)
    }

    app := NewApp(store)
    server := httptest.NewServer(app.Routes())
    defer server.Close()

    // helper to make requests
    client := server.Client()

    // 1) Create book
    book := &Book{Title: "1984", Author: "George Orwell", Year: 1949}
    buf := new(bytes.Buffer)
    _ = json.NewEncoder(buf).Encode(book)
    resp, err := client.Post(server.URL+"/books", "application/json", buf)
    if err != nil {
        t.Fatalf("create request failed: %v", err)
    }
    if resp.StatusCode != http.StatusCreated {
        b, _ := io.ReadAll(resp.Body)
        t.Fatalf("expected created, got %d: %s", resp.StatusCode, string(b))
    }
    var created Book
    _ = json.NewDecoder(resp.Body).Decode(&created)
    if created.ID == 0 {
        t.Fatalf("expected id assigned")
    }
    if created.Title != book.Title {
        t.Fatalf("title mismatch")
    }

    // 2) Get by id
    resp, err = client.Get(fmt.Sprintf("%s/books/%d", server.URL, created.ID))
    if err != nil {
        t.Fatalf("get by id failed: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("get by id status: %d", resp.StatusCode)
    }
    var got Book
    _ = json.NewDecoder(resp.Body).Decode(&got)
    if got.ID != created.ID {
        t.Fatalf("id mismatch")
    }

    // 3) Update
    updated := &Book{Title: "Nineteen Eighty-Four", Author: "George Orwell", Year: 1949}
    buf = new(bytes.Buffer)
    _ = json.NewEncoder(buf).Encode(updated)
    reqUp, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/books/%d", server.URL, created.ID), buf)
    reqUp.Header.Set("Content-Type", "application/json")
    resp, err = client.Do(reqUp)
    if err != nil {
        t.Fatalf("update failed: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("update status: %d", resp.StatusCode)
    }
    var up Book
    _ = json.NewDecoder(resp.Body).Decode(&up)
    if up.Title != updated.Title {
        t.Fatalf("update didn't apply")
    }

    // 4) Get all
    resp, err = client.Get(server.URL + "/books")
    if err != nil {
        t.Fatalf("get all failed: %v", err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("get all status: %d", resp.StatusCode)
    }
    var list []*Book
    _ = json.NewDecoder(resp.Body).Decode(&list)
    if len(list) != 1 {
        t.Fatalf("expected 1 book, got %d", len(list))
    }

    // 5) Delete
    reqDel, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/books/%d", server.URL, created.ID), nil)
    resp, err = client.Do(reqDel)
    if err != nil {
        t.Fatalf("delete failed: %v", err)
    }
    if resp.StatusCode != http.StatusNoContent {
        t.Fatalf("delete status: %d", resp.StatusCode)
    }

    // 6) Ensure gone
    resp, err = client.Get(fmt.Sprintf("%s/books/%d", server.URL, created.ID))
    if err != nil {
        t.Fatalf("final get failed: %v", err)
    }
    if resp.StatusCode != http.StatusNotFound {
        t.Fatalf("expected 404 after delete, got %d", resp.StatusCode)
    }
}