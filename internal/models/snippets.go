package models

import (
	"database/sql"
	"errors"
	"time"
)

// snippet fields should correspond to the columns in the snippets table
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

// Insert a new snippet into the database
// db.exec() prepares the statement and executes it
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// pointer to a new zero-valued Snippet struct that we'll fill up
	s := &Snippet{}

	// QueryRow() returns a pointer to a sql.Row object, never nil
	// Scan() scans the row into the fields of the Snippet struct
	// all the fields in the row must match the fields in the Snippet struct
	// Scan() will also convert raw data into proper Go types as needed
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		} else {
			return nil, err
		}
	}
	return s, nil
}

// Latest returns the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]*Snippet, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// make sure to close after defer otherwise it could panic
	// MUST do this, otherwise connection to db stays open and
	// can use up the entire connection pool
	defer rows.Close()

	// initialize an empty slice of snippets types
	snippets := []*Snippet{}

	// iterate over the rows and append each snippet into the slice
	for rows.Next() {
		s := &Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	// check for errors, one of them could cause errors even after iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
