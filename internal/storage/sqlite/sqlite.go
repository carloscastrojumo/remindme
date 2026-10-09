package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/carloscastrojumo/remindme/internal/storage"
	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS notes (
	id          INTEGER PRIMARY KEY AUTOINCREMENT,
	command     TEXT NOT NULL UNIQUE,
	description TEXT NOT NULL DEFAULT ''
);
CREATE TABLE IF NOT EXISTS note_tags (
	note_id INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
	tag     TEXT NOT NULL,
	PRIMARY KEY (note_id, tag)
);
CREATE INDEX IF NOT EXISTS note_tags_tag ON note_tags(tag);
`

const selectNotes = `
SELECT n.id, n.command, n.description,
	(SELECT group_concat(t.tag, char(31) ORDER BY t.rowid) FROM note_tags t WHERE t.note_id = n.id)
FROM notes n`

// Config is the SQLite storage config
type Config struct {
	Path string
}

// Store is the SQLite storage
type Store struct {
	db *sql.DB
}

// Initialize opens the SQLite database, creating the file and schema if needed
func Initialize(config *Config) (*Store, error) {
	f, err := os.OpenFile(config.Path, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("open notes database %s: %w", config.Path, err)
	}
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("close notes database %s: %w", config.Path, err)
	}

	db, err := sql.Open("sqlite", config.Path+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("open notes database %s: %w", config.Path, err)
	}
	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("create schema in %s: %w", config.Path, err)
	}
	return &Store{db: db}, nil
}

// Close closes the database
func (s *Store) Close() error {
	return s.db.Close()
}

// Insert adds a note, or updates the tags and description of the note with the same command
func (s *Store) Insert(note storage.Note) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id int64
	err = tx.QueryRow(`INSERT INTO notes (command, description) VALUES (?, ?)
		ON CONFLICT (command) DO UPDATE SET description = excluded.description
		RETURNING id`, note.Command, note.Description).Scan(&id)
	if err != nil {
		return fmt.Errorf("save note: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM note_tags WHERE note_id = ?`, id); err != nil {
		return fmt.Errorf("replace tags: %w", err)
	}
	for _, tag := range note.Tags {
		if _, err := tx.Exec(`INSERT OR IGNORE INTO note_tags (note_id, tag) VALUES (?, ?)`, id, tag); err != nil {
			return fmt.Errorf("save tag %q: %w", tag, err)
		}
	}
	return tx.Commit()
}

// Get returns a note by id
func (s *Store) Get(id string) (storage.Note, error) {
	notes, err := s.query(` WHERE n.id = ?`, id)
	if err != nil {
		return storage.Note{}, err
	}
	if len(notes) == 0 {
		return storage.Note{}, fmt.Errorf("note %s not found", id)
	}
	return notes[0], nil
}

// GetByTags returns the notes that have any of the tags
func (s *Store) GetByTags(tags []string) ([]storage.Note, error) {
	return s.query(` WHERE n.id IN (SELECT note_id FROM note_tags WHERE tag IN (`+placeholders(len(tags))+`))`, toArgs(tags)...)
}

// GetAll returns all notes
func (s *Store) GetAll() ([]storage.Note, error) {
	return s.query(``)
}

// GetTags returns all tags in the order they were first used
func (s *Store) GetTags() ([]string, error) {
	rows, err := s.db.Query(`SELECT tag FROM note_tags GROUP BY tag ORDER BY MIN(rowid)`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

// Delete deletes a note by id
func (s *Store) Delete(id string) error {
	result, err := s.db.Exec(`DELETE FROM notes WHERE id = ?`, id)
	if err != nil {
		return err
	}
	return requireDeleted(result, fmt.Errorf("note %s not found", id))
}

// DeleteByTags deletes the notes that have any of the tags
func (s *Store) DeleteByTags(tags []string) error {
	result, err := s.db.Exec(`DELETE FROM notes WHERE id IN (SELECT note_id FROM note_tags WHERE tag IN (`+placeholders(len(tags))+`))`, toArgs(tags)...)
	if err != nil {
		return err
	}
	return requireDeleted(result, fmt.Errorf("no notes found with tags %v", tags))
}

// Search returns the notes where any search word is a substring of any of the search locations
func (s *Store) Search(searchWords []string, searchLocations []string) ([]storage.Note, error) {
	var conditions []string
	var args []any
	for _, searchLocation := range searchLocations {
		var condition string
		switch searchLocation {
		case "command":
			condition = `instr(n.command, ?) > 0`
		case "description":
			condition = `instr(n.description, ?) > 0`
		case "tags":
			condition = `EXISTS (SELECT 1 FROM note_tags t WHERE t.note_id = n.id AND instr(t.tag, ?) > 0)`
		default:
			continue
		}
		for _, searchWord := range searchWords {
			conditions = append(conditions, condition)
			args = append(args, searchWord)
		}
	}
	if len(conditions) == 0 {
		return nil, nil
	}
	return s.query(` WHERE `+strings.Join(conditions, " OR "), args...)
}

func (s *Store) query(where string, args ...any) ([]storage.Note, error) {
	rows, err := s.db.Query(selectNotes+where+` ORDER BY n.id`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []storage.Note
	for rows.Next() {
		var id int64
		var note storage.Note
		var tags sql.NullString
		if err := rows.Scan(&id, &note.Command, &note.Description, &tags); err != nil {
			return nil, err
		}
		note.ID = strconv.FormatInt(id, 10)
		if tags.Valid {
			note.Tags = strings.Split(tags.String, "\x1f")
		}
		notes = append(notes, note)
	}
	return notes, rows.Err()
}

func requireDeleted(result sql.Result, notFound error) error {
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return notFound
	}
	return nil
}

func placeholders(n int) string {
	return strings.TrimSuffix(strings.Repeat("?,", n), ",")
}

func toArgs(values []string) []any {
	args := make([]any, len(values))
	for i, v := range values {
		args[i] = v
	}
	return args
}
