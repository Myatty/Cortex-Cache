package models

import (
	"database/sql"
	"errors"
	"time"
)

// define a Snippet type to hold the data for an individual snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel wraps a sql.DB connection pool
type SnippetModel struct {
	DB *sql.DB
}

// inserts new snippet to database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// to get id of newly inserted record in snippet table
	// LastInsertId() is not supported by PostgreSQL
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// id is type int64, so convert it to int type
	return int(id), nil

}

// returns snippet using id
func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
		    WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	// return Snippet object
	return s, nil

}

// returns recent 10 snippet
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
