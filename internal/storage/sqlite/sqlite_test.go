package sqlite

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/carloscastrojumo/remindme/internal/storage"
)

func newStore(t *testing.T, notes ...storage.Note) (*Store, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "notes.db")
	store, err := Initialize(&Config{Path: path})
	if err != nil {
		t.Fatalf("Initialize: %v", err)
	}
	t.Cleanup(func() { store.db.Close() })
	for _, n := range notes {
		if err := store.Insert(n); err != nil {
			t.Fatalf("Insert(%q): %v", n.Command, err)
		}
	}
	return store, path
}

func commands(notes []storage.Note) []string {
	var result []string
	for _, n := range notes {
		result = append(result, n.Command)
	}
	return result
}

var fixture = []storage.Note{
	{Command: "kubectl get pods", Description: "List pods", Tags: []string{"k8s"}},
	{Command: "ip a", Description: "Show addresses", Tags: []string{"net", "linux"}},
	{Command: "df -h", Description: "Disk usage 100%", Tags: []string{"linux"}},
}

func TestInsertAndGetAll(t *testing.T) {
	store, _ := newStore(t, fixture...)

	notes, err := store.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if got, want := commands(notes), []string{"kubectl get pods", "ip a", "df -h"}; !slices.Equal(got, want) {
		t.Errorf("commands = %q, want %q", got, want)
	}
	if got, want := notes[1].Tags, []string{"net", "linux"}; !slices.Equal(got, want) {
		t.Errorf("tags = %q, want %q", got, want)
	}
}

func TestInsertExistingCommandUpdatesNote(t *testing.T) {
	store, _ := newStore(t, fixture...)

	if err := store.Insert(storage.Note{Command: "ip a", Description: "Addresses", Tags: []string{"network"}}); err != nil {
		t.Fatalf("Insert: %v", err)
	}

	notes, err := store.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(notes) != 3 {
		t.Fatalf("stored %d notes, want 3", len(notes))
	}
	if notes[1].Description != "Addresses" || !slices.Equal(notes[1].Tags, []string{"network"}) {
		t.Errorf("updated note = %+v", notes[1])
	}
}

func TestGet(t *testing.T) {
	store, _ := newStore(t, fixture...)

	note, err := store.Get("2")
	if err != nil {
		t.Fatalf("Get(2): %v", err)
	}
	if note.Command != "ip a" {
		t.Errorf("Get(2).Command = %q, want %q", note.Command, "ip a")
	}

	for _, id := range []string{"99", "abc"} {
		if _, err := store.Get(id); err == nil {
			t.Errorf("Get(%q) returned nil error", id)
		}
	}
}

func TestGetByTags(t *testing.T) {
	store, _ := newStore(t, fixture...)

	notes, err := store.GetByTags([]string{"linux", "net"})
	if err != nil {
		t.Fatalf("GetByTags: %v", err)
	}
	if got, want := commands(notes), []string{"ip a", "df -h"}; !slices.Equal(got, want) {
		t.Errorf("commands = %q, want %q", got, want)
	}
}

func TestGetTags(t *testing.T) {
	store, _ := newStore(t, fixture...)

	tags, err := store.GetTags()
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if want := []string{"k8s", "net", "linux"}; !slices.Equal(tags, want) {
		t.Errorf("tags = %q, want %q", tags, want)
	}
}

func TestDelete(t *testing.T) {
	store, _ := newStore(t, fixture...)

	if err := store.Delete("2"); err != nil {
		t.Fatalf("Delete(2): %v", err)
	}
	if err := store.Delete("2"); err == nil {
		t.Error("deleting a missing note returned nil error")
	}

	tags, err := store.GetTags()
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if want := []string{"k8s", "linux"}; !slices.Equal(tags, want) {
		t.Errorf("tags after delete = %q, want %q", tags, want)
	}
}

func TestDeleteByTags(t *testing.T) {
	store, _ := newStore(t, fixture...)

	if err := store.DeleteByTags([]string{"linux"}); err != nil {
		t.Fatalf("DeleteByTags: %v", err)
	}
	if err := store.DeleteByTags([]string{"linux"}); err == nil {
		t.Error("DeleteByTags with no matches returned nil error")
	}

	notes, err := store.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if got, want := commands(notes), []string{"kubectl get pods"}; !slices.Equal(got, want) {
		t.Errorf("remaining commands = %q, want %q", got, want)
	}
	tags, err := store.GetTags()
	if err != nil {
		t.Fatalf("GetTags: %v", err)
	}
	if want := []string{"k8s"}; !slices.Equal(tags, want) {
		t.Errorf("tags after delete = %q, want %q", tags, want)
	}
}

func TestSearch(t *testing.T) {
	store, _ := newStore(t, fixture...)

	tests := []struct {
		name      string
		words     []string
		locations []string
		want      []string
	}{
		{name: "command", words: []string{"get"}, locations: []string{"command"}, want: []string{"kubectl get pods"}},
		{name: "case sensitive", words: []string{"GET"}, locations: []string{"command"}, want: nil},
		{name: "description", words: []string{"Disk"}, locations: []string{"description"}, want: []string{"df -h"}},
		{name: "percent is literal", words: []string{"%"}, locations: []string{"description"}, want: []string{"df -h"}},
		{name: "underscore is literal", words: []string{"_"}, locations: []string{"command", "description"}, want: nil},
		{name: "tags substring", words: []string{"lin"}, locations: []string{"tags"}, want: []string{"ip a", "df -h"}},
		{name: "any word any location", words: []string{"pods", "net"}, locations: []string{"command", "tags"}, want: []string{"kubectl get pods", "ip a"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notes, err := store.Search(tt.words, tt.locations)
			if err != nil {
				t.Fatalf("Search: %v", err)
			}
			if got := commands(notes); !slices.Equal(got, tt.want) {
				t.Errorf("commands = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDatabaseFileIsOwnerOnlyAndPersists(t *testing.T) {
	store, path := newStore(t, fixture...)
	store.db.Close()

	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := fi.Mode().Perm(); perm&0077 != 0 {
		t.Errorf("database file permissions = %v, want no group/other access", perm)
	}

	reopened, err := Initialize(&Config{Path: path})
	if err != nil {
		t.Fatalf("Initialize after close: %v", err)
	}
	defer reopened.db.Close()
	notes, err := reopened.GetAll()
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(notes) != len(fixture) {
		t.Errorf("reopened database has %d notes, want %d", len(notes), len(fixture))
	}
}
