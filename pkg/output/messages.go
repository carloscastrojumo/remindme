package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type Note struct {
	ID          string   `json:"id"`
	Tags        []string `json:"tags"`
	Command     string   `json:"command"`
	Description string   `json:"description"`
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

	orderedNotes := orderNotes(notes)
	fmt.Println(orderedNotes)
	maxLength := getMaxLength(notes)

	numberOfTags := len(orderedNotes)
	for tag, notes := range orderedNotes {
		numberOfTags--
		size := 20
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
		for _, note := range notes {
			color.HiBlue("ID: %s \n", color.WhiteString(note.ID))
			color.HiBlue("Tags: %s \n", color.GreenString(strings.Join(note.Tags, " ")))
			color.HiBlue("Command: %s \n", color.RedString(note.Command))
			color.HiBlue("Description: %s \n", color.WhiteString(note.Description))
			// add full line only if there are more tags
			if numberOfTags == 0 {
				color.Yellow("%s", strings.Repeat("-", 42+maxLength))
			} else {
				fmt.Println()
			}
		}
	}
}

// get the tag with more length
func getMaxLength(notes []Note) int {
	maxLength := 0
	for _, note := range notes {
		for _, tag := range note.Tags {
			if len(tag) > maxLength {
				maxLength = len(tag)
			}
		}
	}
	return maxLength
}

// orderNotes orders the notes by tag
func orderNotes(notes []Note) map[string][]Note {
	orderedNotes := make(map[string][]Note)
	for _, note := range notes {
		for _, tag := range note.Tags {
			orderedNotes[tag] = append(orderedNotes[tag], note)
		}
	}
	return orderedNotes
}
