package storage_test

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/carloscastrojumo/remindme/internal/storage"
	"github.com/carloscastrojumo/remindme/internal/storage/sqlite"
	"github.com/carloscastrojumo/remindme/internal/storage/yaml"
)

const benchmarkNotes = 1000

type backend struct {
	name string
	open func(b *testing.B, path string) (storage.NoteStorage, func())
}

var backends = []backend{
	{name: "yaml", open: func(b *testing.B, path string) (storage.NoteStorage, func()) {
		store, err := yaml.Initialize(&yaml.Config{Name: path})
		if err != nil {
			b.Fatal(err)
		}
		return store, func() {}
	}},
	{name: "sqlite", open: func(b *testing.B, path string) (storage.NoteStorage, func()) {
		store, err := sqlite.Initialize(&sqlite.Config{Path: path})
		if err != nil {
			b.Fatal(err)
		}
		return store, func() {
			if err := store.Close(); err != nil {
				b.Fatal(err)
			}
		}
	}},
}

func seed(b *testing.B, be backend) string {
	b.Helper()
	path := filepath.Join(b.TempDir(), "notes")
	store, closeStore := be.open(b, path)
	defer closeStore()
	for i := range benchmarkNotes {
		note := storage.Note{
			Command:     fmt.Sprintf("command-%d --flag value", i),
			Description: fmt.Sprintf("description of command %d", i),
			Tags:        []string{fmt.Sprintf("tag-%d", i%20), "common"},
		}
		if err := store.Insert(note); err != nil {
			b.Fatal(err)
		}
	}
	return path
}

// Each iteration opens the store and runs one operation, like a single rmm invocation.
func benchmarkInvocation(b *testing.B, op func(storage.NoteStorage) error) {
	for _, be := range backends {
		b.Run(be.name, func(b *testing.B) {
			path := seed(b, be)
			b.ResetTimer()
			for b.Loop() {
				store, closeStore := be.open(b, path)
				if err := op(store); err != nil {
					b.Fatal(err)
				}
				closeStore()
			}
		})
	}
}

func BenchmarkGetAll(b *testing.B) {
	benchmarkInvocation(b, func(s storage.NoteStorage) error {
		_, err := s.GetAll()
		return err
	})
}

func BenchmarkGetByTags(b *testing.B) {
	benchmarkInvocation(b, func(s storage.NoteStorage) error {
		_, err := s.GetByTags([]string{"tag-7"})
		return err
	})
}

func BenchmarkSearch(b *testing.B) {
	benchmarkInvocation(b, func(s storage.NoteStorage) error {
		_, err := s.Search([]string{"command-99"}, []string{"command", "description", "tags"})
		return err
	})
}

func BenchmarkUpdate(b *testing.B) {
	benchmarkInvocation(b, func(s storage.NoteStorage) error {
		return s.Insert(storage.Note{Command: "command-500 --flag value", Description: "updated", Tags: []string{"tag-0"}})
	})
}
