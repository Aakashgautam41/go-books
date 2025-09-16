package main

import "database/sql"

type Book struct {
    ID     int64  `json:"id"`
    Title  string `json:"title"`
    Author string `json:"author"`
    Year   int    `json:"year"`
}

// helper to scan a book from rows.Scan
func scanBook(row scanner) (*Book, error) {
    b := &Book{}
    err := row.Scan(&b.ID, &b.Title, &b.Author, &b.Year)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    return b, nil
}

// scanner interface to accept *sql.Row and *sql.Rows
type scanner interface {
    Scan(dest ...any) error
}
