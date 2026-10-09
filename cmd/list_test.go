package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/atotto/clipboard"
	"github.com/carloscastrojumo/remindme/pkg/storage"
	"github.com/carloscastrojumo/remindme/pkg/storage/yaml"
	"github.com/fatih/color"
)

type fakeStore struct {
	storage.NoteStorage
	note yaml.Note
}

func (f fakeStore) Get(id string) (interface{}, error) {
	return f.note, nil
}

func captureOutput(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	prevOutput, prevNoColor, prevClipboard := color.Output, color.NoColor, clipboard.Unsupported
	color.Output, color.NoColor, clipboard.Unsupported = &buf, true, true
	t.Cleanup(func() {
		color.Output, color.NoColor, clipboard.Unsupported = prevOutput, prevNoColor, prevClipboard
	})
	return &buf
}

func TestListByIDPrintsNote(t *testing.T) {
	out := captureOutput(t)
	noteService = storage.NewNoteService(fakeStore{note: yaml.Note{ID: "42", Tags: []string{"k8s"}, Command: "kubectl get pods"}})
	t.Cleanup(func() { noteService = nil })
	if err := listCmd.Flags().Set("id", "42"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listCmd.Flags().Set("id", "") })

	listCmd.Run(listCmd, nil)

	if !strings.Contains(out.String(), "kubectl get pods") {
		t.Errorf("output does not contain the note command:\n%s", out.String())
	}
}
