package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Note note struct
type Note struct {
	ID          string   `json:"id"`
	Tags        []string `json:"tags"`
	Command     string   `json:"command"`
	Description string   `json:"description"`
}

// orderedNote struct
type orderedNote struct {
	Tags  []string
	Notes []Note
}

// Print print the notes
func Print(note interface{}) {
	notes := []Note{}
	s, _ := json.MarshalIndent(note, "", "\t")
	if err := json.Unmarshal(s, &notes); err != nil {
		color.Red("Error while unmarshalling notes: %s", err)
	}

	if len(notes) == 0 {
		color.Yellow("No notes found")
		return
	}

	orderedNotes := processNotes(notes)
	maxLength := getMaxLength(orderedNotes)
	numberOfNotes := len(orderedNotes)

	for _, orderedNote := range orderedNotes {
		numberOfNotes--
		size := 20
		tag := strings.Join(orderedNote.Tags, " ")

		if len(tag) < maxLength {
			size = size + (maxLength-len(tag))/2
		} else {
			size = size - (len(tag)-maxLength)/2
		}

		rightSize := size
		if size%2 != 0 {
			rightSize = size + 1
		}

		leftPad := strings.Repeat("-", size)
		rightPad := strings.Repeat("-", rightSize)
		color.Yellow("%s %s %s", leftPad, tag, rightPad)

		for _, note := range orderedNote.Notes {
			color.HiBlue("ID: %s \n", color.WhiteString(note.ID))
			color.HiBlue("Tags: %s \n", color.GreenString(strings.Join(note.Tags, " ")))
			color.HiBlue("Command: %s \n", color.RedString(note.Command))
			color.HiBlue("Description: %s \n", color.WhiteString(note.Description))
			// add full line only if there are more tags
			if numberOfNotes == 0 {
				color.Yellow("%s", strings.Repeat("-", 42+maxLength))
			} else {
				fmt.Println()
			}
		}
	}
}

// get the tag with more length
func getMaxLength(notes []orderedNote) int {
	maxLength := 0
	for _, note := range notes {
		curLen := len(strings.Join(note.Tags, " "))
		if curLen > maxLength {
			maxLength = curLen
		}
	}
	return maxLength
}

// processNotes combine tags to group them if they are the same
func processNotes(notes []Note) []orderedNote {
	orderedNotes := []orderedNote{}

	for _, note := range notes {
		id, missingTags := getIDAndMissingTags(note.Tags, orderedNotes)
		if len(orderedNotes) == 0 || id == -1 {
			or := orderedNote{}
			or.Tags = note.Tags
			or.Notes = append(or.Notes, note)

			orderedNotes = append(orderedNotes, or)
		} else {
			orderedNotes[id].Tags = append(orderedNotes[id].Tags, missingTags...)

			if !containsID(note.ID, orderedNotes[id].Notes) {
				orderedNotes[id].Notes = append(orderedNotes[id].Notes, note)
			}
		}
	}
	return orderedNotes
}

// getIDAndMissingTags gets the id of ordered note and missing tags for that id
func getIDAndMissingTags(tags []string, notes []orderedNote) (int, []string) {
	missingTags := []string{}
	id := -1
	breakFor := false

	for i, note := range notes {
		if breakFor {
			break
		}
		for _, tag := range tags {
			if containsTag(tag, note.Tags) {
				id = i
				breakFor = true
				break
			}
		}
	}

	if id > -1 {
		for _, tag := range tags {
			exists := false
			for _, nt := range notes[id].Tags {
				if tag == nt {
					exists = true
				}
			}
			if !exists {
				missingTags = append(missingTags, tag)
			}
		}
	}
	return id, missingTags
}

// containsTag check if tag exists in group of tags
func containsTag(tag string, tags []string) bool {
	for _, t := range tags {
		if tag == t {
			return true
		}
	}
	return false
}

// containsID check if id already exists
func containsID(id string, notes []Note) bool {
	for _, n := range notes {
		if n.ID == id {
			return true
		}
	}
	return false
}
