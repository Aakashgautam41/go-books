package main

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
)

type App struct {
    Store *Store
}

func NewApp(store *Store) *App {
    return &App{Store: store}
}

func (a *App) Routes() http.Handler {
    r := chi.NewRouter()
    r.Post("/books", a.createBookHandler)
    r.Get("/books", a.getAllBooksHandler)
    r.Get("/books/{id}", a.getBookHandler)
    r.Put("/books/{id}", a.updateBookHandler)
    r.Delete("/books/{id}", a.deleteBookHandler)
    return r
}

// simple request validation
func validateBookInput(b *Book) (int, string) {
    if b.Title == "" {
        return http.StatusBadRequest, "title is required"
    }
    if b.Author == "" {
        return http.StatusBadRequest, "author is required"
    }
    if b.Year <= 0 {
        return http.StatusBadRequest, "year must be a positive integer"
    }
    return http.StatusOK, ""
}

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

func (a *App) createBookHandler(w http.ResponseWriter, r *http.Request) {
    var b Book
    if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    if code, msg := validateBookInput(&b); code != http.StatusOK {
        writeError(w, code, msg)
        return
    }
    id, err := a.Store.CreateBook(&b)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to create book")
        return
    }
    b.ID = id
    writeJSON(w, http.StatusCreated, b)
}

func (a *App) getAllBooksHandler(w http.ResponseWriter, r *http.Request) {
    books, err := a.Store.GetAllBooks()
    if err != nil {
        writeError(w, http.StatusInternalServerError, "failed to fetch books")
        return
    }
    writeJSON(w, http.StatusOK, books)
}

func (a *App) getBookHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }
    b, err := a.Store.GetBookByID(id)
    if err != nil {
        if err == sql.ErrNoRows {
            writeError(w, http.StatusNotFound, "book not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "failed to fetch book")
        return
    }
    writeJSON(w, http.StatusOK, b)
}

func (a *App) updateBookHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }
    var b Book
    if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON")
        return
    }
    if code, msg := validateBookInput(&b); code != http.StatusOK {
        writeError(w, code, msg)
        return
    }
    err = a.Store.UpdateBook(id, &b)
    if err != nil {
        if err == sql.ErrNoRows {
            writeError(w, http.StatusNotFound, "book not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "failed to update book")
        return
    }
    // return updated book
    b.ID = id
    writeJSON(w, http.StatusOK, b)
}

func (a *App) deleteBookHandler(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        writeError(w, http.StatusBadRequest, "invalid id")
        return
    }
    err = a.Store.DeleteBook(id)
    if err != nil {
        if err == sql.ErrNoRows {
            writeError(w, http.StatusNotFound, "book not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "failed to delete book")
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
