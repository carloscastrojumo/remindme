package storage

import (
	"path/filepath"
	"slices"
	"testing"

	"github.com/carloscastrojumo/remindme/pkg/storage/yaml"
)

func newYAMLService(t *testing.T) *NoteService {
	t.Helper()
	store, err := GetStorage(&Config{
		StorageType:   "yaml",
		StorageConfig: &yaml.Config{Name: filepath.Join(t.TempDir(), "data.yaml")},
	})
	if err != nil {
		t.Fatalf("GetStorage: %v", err)
	}
	return NewNoteService(store)
}

func TestAddRejectsMissingFields(t *testing.T) {
	tests := []struct {
		name string
		note Note
	}{
		{name: "empty note", note: Note{}},
		{name: "blank command", note: Note{Command: "  ", Tags: []string{"x"}}},
		{name: "no tags", note: Note{Command: "ls"}},
		{name: "blank tags", note: Note{Command: "ls", Tags: []string{"", " "}}},
		{name: "description only", note: Note{Description: "list files"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := newYAMLService(t)

			if err := service.Add(tt.note); err == nil {
				t.Fatal("Add returned nil error")
			}

			all, err := service.GetAll()
			if err != nil {
				t.Fatalf("GetAll: %v", err)
			}
			if notes := all.([]yaml.Note); len(notes) != 0 {
				t.Errorf("invalid note was stored: %+v", notes)
			}
		})
	}
}

func TestAddTrimsFields(t *testing.T) {
	service := newYAMLService(t)

	if err := service.Add(Note{Command: " ls -la ", Tags: []string{" shell", "", "files "}}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	all, err := service.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	notes := all.([]yaml.Note)
	if len(notes) != 1 {
		t.Fatalf("stored %d notes, want 1", len(notes))
	}
	if notes[0].Command != "ls -la" {
		t.Errorf("Command = %q, want %q", notes[0].Command, "ls -la")
	}
	if want := []string{"shell", "files"}; !slices.Equal(notes[0].Tags, want) {
		t.Errorf("Tags = %q, want %q", notes[0].Tags, want)
	}
}
