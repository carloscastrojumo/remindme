package yaml

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func newStore(t *testing.T, notes ...Note) (*Yaml, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "data.yaml")
	store, err := Initialize(&Config{Name: path})
	if err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	for _, n := range notes {
		if err := store.Insert(n); err != nil {
			t.Fatalf("Insert(%q): %v", n.Command, err)
		}
	}
	return store, path
}

func commands(notes []Note) []string {
	var result []string
	for _, n := range notes {
		result = append(result, n.Command)
	}
	return result
}

func TestDeleteByTags(t *testing.T) {
	notes := []Note{
		{Command: "a", Tags: []string{"x"}},
		{Command: "b", Tags: []string{"x"}},
		{Command: "c", Tags: []string{"y"}},
		{Command: "d", Tags: []string{"y", "x"}},
		{Command: "e", Tags: []string{"z"}},
	}

	tests := []struct {
		name    string
		tags    []string
		want    []string
		wantErr bool
	}{
		{name: "adjacent matches", tags: []string{"x"}, want: []string{"c", "e"}},
		{name: "multiple tags", tags: []string{"x", "z"}, want: []string{"c"}},
		{name: "no match", tags: []string{"nope"}, want: []string{"a", "b", "c", "d", "e"}, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store, path := newStore(t, notes...)

			if err := store.DeleteByTags(tt.tags); (err != nil) != tt.wantErr {
				t.Fatalf("DeleteByTags error = %v, wantErr %v", err, tt.wantErr)
			}

			reloaded, err := Initialize(&Config{Name: path})
			if err != nil {
				t.Fatalf("Initialize after delete: %v", err)
			}
			if got := commands(reloaded.Notes); !slices.Equal(got, tt.want) {
				t.Errorf("remaining commands = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeleteUnknownIDReturnsError(t *testing.T) {
	store, _ := newStore(t, Note{Command: "a", Tags: []string{"x"}})

	if err := store.Delete("missing"); err == nil {
		t.Fatal("Delete(missing) returned nil error")
	}
	if len(store.Notes) != 1 {
		t.Errorf("Delete(missing) changed notes: %+v", store.Notes)
	}
}

func TestInitializeRejectsCorruptFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "data.yaml")
	if err := os.WriteFile(path, []byte("::: not yaml ["), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := Initialize(&Config{Name: path}); err == nil {
		t.Fatal("Initialize returned nil error for a corrupt file")
	}
}

func TestGet(t *testing.T) {
	store, _ := newStore(t, Note{Command: "a", Tags: []string{"x"}})
	id := store.Notes[0].ID

	got, err := store.Get(id)
	if err != nil {
		t.Fatalf("Get(%q): %v", id, err)
	}
	if note := got.(Note); note.Command != "a" {
		t.Errorf("Get(%q).Command = %q, want %q", id, note.Command, "a")
	}

	if _, err := store.Get("missing"); err == nil {
		t.Error("Get(missing) returned nil error")
	}
}
