package main

import (
    "database/sql"
    "errors"
)

type Store struct {
    DB *sql.DB
}

func NewStore(db *sql.DB) *Store {
    return &Store{DB: db}
}

func (s *Store) CreateBook(b *Book) (int64, error) {
    var id int64
    err := s.DB.QueryRow(
        "INSERT INTO books (title, author, year) VALUES ($1, $2, $3) RETURNING id",
        b.Title, b.Author, b.Year,
    ).Scan(&id)
    if err != nil {
        return 0, err
    }
    return id, nil
}

func (s *Store) GetBookByID(id int64) (*Book, error) {
    row := s.DB.QueryRow("SELECT id, title, author, year FROM books WHERE id=$1", id)
    book, err := scanBook(row)
    if err != nil {
        return nil, err
    }
    if book == nil {
        return nil, sql.ErrNoRows
    }
    return book, nil
}

func (s *Store) GetAllBooks() ([]*Book, error) {
    rows, err := s.DB.Query("SELECT id, title, author, year FROM books ORDER BY id")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var books []*Book
    for rows.Next() {
        b, err := scanBook(rows)
        if err != nil {
            return nil, err
        }
        books = append(books, b)
    }
    return books, rows.Err()
}

func (s *Store) UpdateBook(id int64, b *Book) error {
    res, err := s.DB.Exec("UPDATE books SET title=$1, author=$2, year=$3 WHERE id=$4",
        b.Title, b.Author, b.Year, id)
    if err != nil {
        return err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if n == 0 {
        return sql.ErrNoRows
    }
    return nil
}

func (s *Store) DeleteBook(id int64) error {
    res, err := s.DB.Exec("DELETE FROM books WHERE id=$1", id)
    if err != nil {
        return err
    }
    n, err := res.RowsAffected()
    if err != nil {
        return err
    }
    if n == 0 {
        return sql.ErrNoRows
    }
    return nil
}

// For migrations convenience, export a small helper to exec SQL
func (s *Store) Exec(statement string) error {
    _, err := s.DB.Exec(statement)
    return err
}

// small helper to check connection
func (s *Store) Ping() error {
    return s.DB.Ping()
}

// helpful error conversion
var ErrNotFound = errors.New("not found")
