package cmd

import (
	"bytes"
	"errors"
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
	err  error
}

func (f fakeStore) Get(id string) (interface{}, error) {
	if f.err != nil {
		return nil, f.err
	}
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

func listByID(t *testing.T, store fakeStore, id string) error {
	t.Helper()
	noteService = storage.NewNoteService(store)
	t.Cleanup(func() { noteService = nil })
	if err := listCmd.Flags().Set("id", id); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { listCmd.Flags().Set("id", "") })
	return listCmd.RunE(listCmd, nil)
}

func TestListByIDPrintsNote(t *testing.T) {
	out := captureOutput(t)

	err := listByID(t, fakeStore{note: yaml.Note{ID: "42", Tags: []string{"k8s"}, Command: "kubectl get pods"}}, "42")

	if err != nil {
		t.Fatalf("list --id: %v", err)
	}
	if !strings.Contains(out.String(), "kubectl get pods") {
		t.Errorf("output does not contain the note command:\n%s", out.String())
	}
}

func TestListByIDReturnsStoreError(t *testing.T) {
	captureOutput(t)
	notFound := errors.New("note 42 not found")

	err := listByID(t, fakeStore{err: notFound}, "42")

	if !errors.Is(err, notFound) {
		t.Errorf("list --id error = %v, want %v", err, notFound)
	}
}

func TestRemoveRequiresIDOrTags(t *testing.T) {
	if err := removeCmd.ValidateFlagGroups(); err == nil {
		t.Error("rm without --id or --tags passed flag validation")
	}
}
