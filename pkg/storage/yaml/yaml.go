package yaml

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	yaml "gopkg.in/yaml.v3"
)

type Note struct {
	Id          string   `yaml:"id"`
	Tags        []string `yaml:"tags"`
	Command     string   `yaml:"command"`
	Description string   `yaml:"description"`
}

type Yaml struct {
	File  *os.File
	Notes []Note
}

type Config struct {
	Name string
}

func Initialize(config *Config) *Yaml {
	// check if file exists, if not create it
	f, err := os.OpenFile(config.Name, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Read([]byte{})

	var notes []Note

	// check if file siza > 0, if so read file and unmarshal it to Notes struct
	if fi, _ := f.Stat(); fi.Size() > 0 {
		if err := yaml.NewDecoder(f).Decode(&notes); err != nil {
			return nil
		}
	}

	return &Yaml{File: f, Notes: notes}
}

func (y *Yaml) Insert(note interface{}) error {
	newNote := note.(Note)

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
	rand.Seed(time.Now().UnixNano())
	newNote.Id = strconv.Itoa(rand.Intn(1000000))

	// append new note to notes
	y.Notes = append(y.Notes, newNote)

	return y.save()
}

func (y *Yaml) save() error {
	data, err := yaml.Marshal(y.Notes)
	if err != nil {
		return errors.New("Error while marshalling notes")
	}

	if err := os.WriteFile(y.File.Name(), data, 0644); err != nil {
		return errors.New("Error while writing notes to file")
	}
	return nil
}

func (y *Yaml) Get(id string) (interface{}, error) {
	for _, note := range y.Notes {
		if note.Id == id {
			return note, nil
		}
	}

	return nil, nil
}

func (y *Yaml) GetByTags(tags []string) (interface{}, error) {
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

	return filteredNotes, nil
}

func (y *Yaml) GetAll() (interface{}, error) {
	return y.Notes, nil
}

func (y *Yaml) Delete(id string) error {
	for i, note := range y.Notes {
		if note.Id == id {
			y.Notes = append(y.Notes[:i], y.Notes[i+1:]...)
			return y.save()
		}
	}
	return nil
}

func (y *Yaml) DeleteByTags(tags []string) error {
	for i, note := range y.Notes {
		for _, tag := range tags {
			for _, noteTag := range note.Tags {
				if noteTag == tag {
					y.Notes = append(y.Notes[:i], y.Notes[i+1:]...)
				}
			}
		}
	}
	return y.save()
}

func (y *Yaml) openFile() *os.File {
	f, err := os.Open(y.File.Name())
	if err != nil {
		log.Fatal(err)
	}
	return f
}

func (y *Yaml) closeFile() error {
	return y.File.Close()
}

func (y *Yaml) Search(searchWord string, searchLocations []string) (interface{}, error) {
	fmt.Println("Searching: ", searchWord)
	fmt.Println("In: ", searchLocations)
	var filteredNotes []Note
	var notes []Note
	var err error

	for _, searchLocation := range searchLocations {
		switch searchLocation {
		case "command":
			notes, err = y.SearchInCommand(searchWord)
		case "description":
			notes, err = y.SearchInDescription(searchWord)
		case "tags":
			notes, err = y.SearchInTags(searchWord)
		}

		if err != nil {
			color.Red("Error while getting notes by tags: %s", err)
		}

		filteredNotes = y.AppendSearchResults(filteredNotes, notes)
	}
	return filteredNotes, nil
}

func (y *Yaml) SearchInTags(searchWord string) ([]Note, error) {
	var filteredNotes []Note
	for _, note := range y.Notes {
		for _, noteTag := range note.Tags {
			if strings.Contains(noteTag, searchWord) {
				filteredNotes = append(filteredNotes, note)
			}
		}
	}
	return filteredNotes, nil
}

func (y *Yaml) SearchInCommand(searchWord string) ([]Note, error) {
	var filteredNotes []Note
	for _, note := range y.Notes {
		if strings.Contains(note.Command, searchWord) {
			filteredNotes = append(filteredNotes, note)
		}
	}
	return filteredNotes, nil
}

func (y *Yaml) SearchInDescription(searchWord string) ([]Note, error) {
	var filteredNotes []Note
	for _, note := range y.Notes {
		if strings.Contains(note.Description, searchWord) {
			filteredNotes = append(filteredNotes, note)
		}
	}
	return filteredNotes, nil
}

func (y *Yaml) AppendSearchResults(notes []Note, newNotes []Note) []Note {
	var filteredNotes []Note
	filteredNotes = notes

	for _, newNote := range newNotes {
		if !(ResultContainsId(filteredNotes, newNote.Id)) {
			filteredNotes = append(filteredNotes, newNote)
		}
	}
	return filteredNotes
}

func ResultContainsId(notes []Note, searchId string) bool {
	for _, item := range notes {
		if item.Id == searchId {
			return true
		}
	}
	return false
}
