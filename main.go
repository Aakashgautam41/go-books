package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    _ "github.com/lib/pq"
)

func main() {
    // read DB configuration from env
    dbHost := getEnv("DB_HOST", "localhost")
    dbPort := getEnv("DB_PORT", "5432")
    dbUser := getEnv("DB_USER", "postgres")
    dbPassword := getEnv("DB_PASSWORD", "postgres")
    dbName := getEnv("DB_NAME", "booksdb")
    addr := getEnv("ADDR", ":8080")

    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("failed to open db: %v", err)
    }
    defer db.Close()

    // Wait and retry until DB is ready (useful in docker-compose)
    for i := 0; i < 15; i++ {
        if err := db.Ping(); err == nil {
            break
        }
        log.Println("waiting for database...")
        time.Sleep(2 * time.Second)
    }

    store := NewStore(db)
    app := NewApp(store)

    log.Printf("starting server on %s", addr)
    log.Fatal(http.ListenAndServe(addr, app.Routes()))
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}

//https://chatgpt.com/share/68c95990-7cf0-8006-bb86-2c7814886ca5
