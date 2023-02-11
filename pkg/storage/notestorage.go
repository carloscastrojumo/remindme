package storage

import (
	"errors"
	"strings"

	mongo "github.com/carloscastrojumo/remindme/pkg/storage/mongo"
	yaml "github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
)

// NoteStorage is the interface that wraps the basic storage methods.
type NoteStorage interface {
	Insert(item interface{}) error
	Get(id string) (interface{}, error)
	GetByTags(tags []string) (interface{}, error)
	GetAll() (interface{}, error)
	GetTags() ([]string, error)
	Delete(id string) error
	DeleteByTags(tags []string) error
	Search(searchWord string, searchLocations []string) (interface{}, error)
}

// NoteService is the service that handles the storage
type NoteService struct {
	store NoteStorage
}

// Config is the configuration for the storage
type Config struct {
	StorageType   string
	StorageConfig interface{}
}

// Note is the struct that represents a note
type Note struct {
	ID          string
	Tags        []string
	Command     string
	Description string
}

var storageType string

// GetStorage returns the storage type
func GetStorage(config *Config) NoteStorage {
	switch config.StorageType {
	case "yaml":
		storageType = "yaml"
		return yaml.Initialize(config.StorageConfig.(*yaml.Config))
	case "mongo":
		storageType = "mongo"
		return mongo.Initialize(config.StorageConfig.(*mongo.Config))
	}
	return nil
}

// NewNoteService returns a new note service
func NewNoteService(store NoteStorage) *NoteService {
	return &NoteService{store: store}
}

// Add adds a new note
func (s *NoteService) Add(note interface{}) error {
	switch storageType {
	case "yaml":
		return s.store.Insert(yaml.Note{
			Tags:        note.(Note).Tags,
			Command:     note.(Note).Command,
			Description: note.(Note).Description,
		})
	case "mongo":
		return s.store.Insert(mongo.Note{
			Tags:        note.(Note).Tags,
			Command:     note.(Note).Command,
			Description: note.(Note).Description,
		})
	}
	return errors.New("storage type not supported")
}

// Get returns a note by id
func (s *NoteService) Get(id string) (interface{}, error) {
	return s.store.Get(id)
}

// GetByTags returns all the notes that match the tags
func (s *NoteService) GetByTags(tags []string) (interface{}, error) {
	return s.store.GetByTags(tags)
}

// GetAll returns all the notes
func (s *NoteService) GetAll() (interface{}, error) {
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

// Search returns all the notes that match the search word
func (s *NoteService) Search(searchWord string, searchLocations []string) (interface{}, error) {
	color.Blue("Searching: %s\n", color.GreenString(searchWord))
	color.Blue("In: %s\n", color.GreenString(strings.Join(searchLocations, " ")))
	return s.store.Search(searchWord, searchLocations)
}
