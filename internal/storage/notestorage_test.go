package storage

import (
	"slices"
	"testing"
)

type memStore struct {
	NoteStorage
	notes []Note
}

func (m *memStore) Insert(note Note) error {
	m.notes = append(m.notes, note)
	return nil
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
			store := &memStore{}

			if err := NewNoteService(store).Add(tt.note); err == nil {
				t.Fatal("Add returned nil error")
			}
			if len(store.notes) != 0 {
				t.Errorf("invalid note was stored: %+v", store.notes)
			}
		})
	}
}

func TestAddTrimsFields(t *testing.T) {
	store := &memStore{}

	if err := NewNoteService(store).Add(Note{Command: " ls -la ", Tags: []string{" shell", "", "files "}}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	if len(store.notes) != 1 {
		t.Fatalf("stored %d notes, want 1", len(store.notes))
	}
	if store.notes[0].Command != "ls -la" {
		t.Errorf("Command = %q, want %q", store.notes[0].Command, "ls -la")
	}
	if want := []string{"shell", "files"}; !slices.Equal(store.notes[0].Tags, want) {
		t.Errorf("Tags = %q, want %q", store.notes[0].Tags, want)
	}
}
