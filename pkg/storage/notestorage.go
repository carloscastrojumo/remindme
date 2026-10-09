package storage

import (
	"errors"
	"os"
	"strings"

	"github.com/fatih/color"
)

// NoteStorage is the interface that wraps the basic storage methods.
type NoteStorage interface {
	Insert(note Note) error
	Get(id string) (Note, error)
	GetByTags(tags []string) ([]Note, error)
	GetAll() ([]Note, error)
	GetTags() ([]string, error)
	Delete(id string) error
	DeleteByTags(tags []string) error
	Search(searchWords []string, searchLocations []string) ([]Note, error)
}

// NoteService is the service that handles the storage
type NoteService struct {
	store NoteStorage
}

// Note is the struct that represents a note
type Note struct {
	ID          string
	Tags        []string
	Command     string
	Description string
}

// NewNoteService returns a new note service
func NewNoteService(store NoteStorage) *NoteService {
	return &NoteService{store: store}
}

// Add validates and adds a new note
func (s *NoteService) Add(note Note) error {
	note.Command = strings.TrimSpace(note.Command)
	var tags []string
	for _, tag := range note.Tags {
		if tag = strings.TrimSpace(tag); tag != "" {
			tags = append(tags, tag)
		}
	}
	note.Tags = tags

	if note.Command == "" {
		return errors.New("command is required")
	}
	if len(note.Tags) == 0 {
		return errors.New("at least one tag is required")
	}

	return s.store.Insert(note)
}

// Get returns a note by id
func (s *NoteService) Get(id string) (Note, error) {
	return s.store.Get(id)
}

// GetByTags returns all the notes that match the tags
func (s *NoteService) GetByTags(tags []string) ([]Note, error) {
	return s.store.GetByTags(tags)
}

// GetAll returns all the notes
func (s *NoteService) GetAll() ([]Note, error) {
	return s.store.GetAll()
}

// GetTags returns all available tags
func (s *NoteService) GetTags() ([]string, error) {
	return s.store.GetTags()
}

// Remove removes a note by id
func (s *NoteService) Remove(id string) error {
	return s.store.Delete(id)
}

// RemoveByTags removes all the notes that match the tags
func (s *NoteService) RemoveByTags(tags []string) error {
	return s.store.DeleteByTags(tags)
}

// Search returns all the notes that match the search words
func (s *NoteService) Search(searchWords []string, searchLocations []string) ([]Note, error) {
	color.New(color.FgBlue).Fprintf(os.Stderr, "Searching: %s\n", color.GreenString(strings.Join(searchWords, " ")))
	color.New(color.FgBlue).Fprintf(os.Stderr, "In: %s\n", color.GreenString(strings.Join(searchLocations, " ")))
	return s.store.Search(searchWords, searchLocations)
}
