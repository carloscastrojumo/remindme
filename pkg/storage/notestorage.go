package storage

import (
	"errors"
	"strings"

	mongo "github.com/carloscastrojumo/remindme/pkg/storage/mongo"
	yaml "github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
)

type NoteStorage interface {
	Insert(item interface{}) error
	Get(id string) (interface{}, error)
	GetByTags(tags []string) (interface{}, error)
	GetAll() (interface{}, error)
	Delete(id string) error
	DeleteByTags(tags []string) error
	Search(searchWord string, searchLocations []string) (interface{}, error)
}

type NoteService struct {
	store NoteStorage
}

type StorageConfig struct {
	StorageType   string
	StorageConfig interface{}
}

type Note struct {
	Id          string
	Tags        []string
	Command     string
	Description string
}

var storageType string

func GetStorage(config *StorageConfig) NoteStorage {
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

func NewNoteService(store NoteStorage) *NoteService {
	return &NoteService{store: store}
}

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

func (s *NoteService) Get(id string) (interface{}, error) {
	return s.store.Get(id)
}

func (s *NoteService) GetAll() (interface{}, error) {
	return s.store.GetAll()
}

func (s *NoteService) GetByTags(tags []string) (interface{}, error) {
	return s.store.GetByTags(tags)
}

func (s *NoteService) Remove(id string) error {
	return s.store.Delete(id)
}

func (s *NoteService) RemoveByTags(tags []string) error {
	return s.store.DeleteByTags(tags)
}

func (s *NoteService) Search(searchWord string, searchLocations []string) (interface{}, error) {
	color.Blue("Searching: %s\n", color.GreenString(searchWord))
	color.Blue("In: %s\n", color.GreenString(strings.Join(searchLocations, " ")))
	return s.store.Search(searchWord, searchLocations)
}
