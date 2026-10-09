package yaml

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v3"
)

// Note is a struct that represents a note in YAML storage
type Note struct {
	ID          string   `yaml:"id"`
	Tags        []string `yaml:"tags"`
	Command     string   `yaml:"command"`
	Description string   `yaml:"description"`
}

// Yaml is a struct that represents YAML storage
type Yaml struct {
	File  *os.File
	Notes []Note
}

// Config is a struct that represents YAML storage config
type Config struct {
	Name string
}

// Initialize the YAML storage
func Initialize(config *Config) (*Yaml, error) {
	// check if file exists, if not create it
	f, err := os.OpenFile(config.Name, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, fmt.Errorf("open notes file %s: %w", config.Name, err)
	}
	defer f.Close()

	f.Read([]byte{})

	var notes []Note

	fi, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat notes file %s: %w", config.Name, err)
	}

	// check if file siza > 0, if so read file and unmarshal it to Notes struct
	if fi.Size() > 0 {
		if err := yaml.NewDecoder(f).Decode(&notes); err != nil {
			return nil, fmt.Errorf("parse notes file %s: %w", config.Name, err)
		}
	}

	return &Yaml{File: f, Notes: notes}, nil
}

// Insert inserts a new note to YAML storage
func (y *Yaml) Insert(note storage.Note) error {
	newNote := Note(note)

	// check if command already exists
	// if it does, update tags and description
	for i, n := range y.Notes {
		if n.Command == newNote.Command {
			y.Notes[i].Tags = newNote.Tags
			y.Notes[i].Description = newNote.Description
			return y.save()
		}
	}

	// if it doesn't, create new one
	// generate new id
	rand.New(rand.NewSource(time.Now().UnixNano()))
	newNote.ID = strconv.Itoa(rand.Intn(1000000))

	// append new note to notes
	y.Notes = append(y.Notes, newNote)

	return y.save()
}

func (y *Yaml) save() error {
	data, err := yaml.Marshal(y.Notes)
	if err != nil {
		return errors.New("error while marshalling notes")
	}

	if err := os.WriteFile(y.File.Name(), data, 0600); err != nil {
		return errors.New("error while writing notes to file")
	}
	return nil
}

// Get returns a note by id
func (y *Yaml) Get(id string) (storage.Note, error) {
	for _, note := range y.Notes {
		if note.ID == id {
			return storage.Note(note), nil
		}
	}

	return storage.Note{}, fmt.Errorf("note %s not found", id)
}

// GetByTags returns notes by tags
func (y *Yaml) GetByTags(tags []string) ([]storage.Note, error) {
	var filteredNotes []Note
	for _, note := range y.Notes {
		for _, tag := range tags {
			for _, noteTag := range note.Tags {
				if noteTag == tag {
					filteredNotes = append(filteredNotes, note)
				}
			}
		}
	}

	return toStorage(filteredNotes), nil
}

// GetAll returns all notes
func (y *Yaml) GetAll() ([]storage.Note, error) {
	return toStorage(y.Notes), nil
}

// GetTags returns all available tags
func (y *Yaml) GetTags() ([]string, error) {
	var tags []string
	for _, note := range y.Notes {
		for _, tag := range note.Tags {
			if !containsTag(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags, nil
}

// Delete deletes a note by id
func (y *Yaml) Delete(id string) error {
	for i, note := range y.Notes {
		if note.ID == id {
			y.Notes = append(y.Notes[:i], y.Notes[i+1:]...)
			return y.save()
		}
	}
	return fmt.Errorf("note %s not found", id)
}

// DeleteByTags deletes notes by tags
func (y *Yaml) DeleteByTags(tags []string) error {
	count := len(y.Notes)
	y.Notes = slices.DeleteFunc(y.Notes, func(note Note) bool {
		for _, tag := range tags {
			if containsTag(note.Tags, tag) {
				return true
			}
		}
		return false
	})
	if len(y.Notes) == count {
		return fmt.Errorf("no notes found with tags %v", tags)
	}
	return y.save()
}

// Search returns notes by search words
func (y *Yaml) Search(searchWords []string, searchLocations []string) ([]storage.Note, error) {
	var filteredNotes []Note
	var notes []Note
	var err error

	for _, searchLocation := range searchLocations {
		switch searchLocation {
		case "command":
			notes, err = y.SearchInCommand(searchWords)
		case "description":
			notes, err = y.SearchInDescription(searchWords)
		case "tags":
			notes, err = y.SearchInTags(searchWords)
		}

		if err != nil {
			color.Red("Error while getting notes by tags: %s", err)
		}

		filteredNotes = y.appendSearchResults(filteredNotes, notes)
	}
	return toStorage(filteredNotes), nil
}

// SearchInTags returns notes by search word in tags
func (y *Yaml) SearchInTags(searchWords []string) ([]Note, error) {
	var filteredNotes []Note
	for _, note := range y.Notes {
		for _, noteTag := range note.Tags {
			for _, searchWord := range searchWords {
				if strings.Contains(noteTag, searchWord) {
					filteredNotes = append(filteredNotes, note)
				}
			}
		}
	}
	return filteredNotes, nil
}

// SearchInCommand returns notes by search word in command
func (y *Yaml) SearchInCommand(searchWords []string) ([]Note, error) {
	var filteredNotes []Note
	for _, note := range y.Notes {
		for _, searchWord := range searchWords {
			if strings.Contains(note.Command, searchWord) {
				filteredNotes = append(filteredNotes, note)
			}
		}
	}
	return filteredNotes, nil
}

// SearchInDescription returns notes by search word in description
func (y *Yaml) SearchInDescription(searchWords []string) ([]Note, error) {
	var filteredNotes []Note
	for _, note := range y.Notes {
		for _, searchWord := range searchWords {
			if strings.Contains(note.Description, searchWord) {
				filteredNotes = append(filteredNotes, note)
			}
		}
	}
	return filteredNotes, nil
}

func (y *Yaml) appendSearchResults(notes []Note, newNotes []Note) []Note {
	var filteredNotes []Note
	filteredNotes = notes

	for _, newNote := range newNotes {
		if !(resultContainsID(filteredNotes, newNote.ID)) {
			filteredNotes = append(filteredNotes, newNote)
		}
	}
	return filteredNotes
}

func resultContainsID(notes []Note, searchID string) bool {
	for _, item := range notes {
		if item.ID == searchID {
			return true
		}
	}
	return false
}

func toStorage(notes []Note) []storage.Note {
	result := make([]storage.Note, len(notes))
	for i, note := range notes {
		result[i] = storage.Note(note)
	}
	return result
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
