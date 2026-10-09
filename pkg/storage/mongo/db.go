package mongo

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/carloscastrojumo/remindme/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Note struct for storing notes in MongoDB
type Note struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Tags        []string           `bson:"tags"`
	Command     string             `bson:"command"`
	Description string             `bson:"description"`
}

func (n Note) toStorage() storage.Note {
	return storage.Note{ID: n.ID.Hex(), Tags: n.Tags, Command: n.Command, Description: n.Description}
}

// Config struct for storing MongoDB client
type Config struct {
	Host       string
	Port       int
	Database   string
	Collection string
}

// Store struct for storing MongoDB client/collection
type Store struct {
	db *mongo.Collection
}

const connectTimeout = 5 * time.Second

// Initialize MongoDB client
func Initialize(config *Config) (*Store, error) {
	address := config.Host + ":" + strconv.Itoa(config.Port)
	opts := options.Client().
		ApplyURI("mongodb://" + address).
		SetConnectTimeout(connectTimeout).
		SetServerSelectionTimeout(connectTimeout)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB at %s: %w", address, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("MongoDB not reachable at %s: %w", address, err)
	}
	return &Store{db: client.Database(config.Database).Collection(config.Collection)}, nil
}

// Insert a note into MongoDB
func (s *Store) Insert(note storage.Note) error {
	_, err := s.db.InsertOne(context.Background(), Note{Tags: note.Tags, Command: note.Command, Description: note.Description})
	return err
}

// Get a note from MongoDB
func (s *Store) Get(id string) (storage.Note, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return storage.Note{}, fmt.Errorf("invalid note id %q: %w", id, err)
	}
	filter := bson.M{"_id": objID}
	result := Note{}
	err = s.db.FindOne(context.Background(), filter).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return storage.Note{}, fmt.Errorf("note %s not found", id)
	}
	if err != nil {
		return storage.Note{}, err
	}
	return result.toStorage(), nil
}

// GetByTags gets notes by tags from MongoDB
func (s *Store) GetByTags(tags []string) ([]storage.Note, error) {
	return s.find(bson.M{"tags": bson.M{"$in": tags}})
}

// GetAll gets all notes from MongoDB
func (s *Store) GetAll() ([]storage.Note, error) {
	return s.find(bson.M{})
}

// GetTags returns all available tags
func (s *Store) GetTags() ([]string, error) {
	var tags []string
	notes, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	for _, note := range notes {
		for _, tag := range note.Tags {
			if !containsTag(tags, tag) {
				tags = append(tags, tag)
			}
		}
	}
	return tags, nil
}

// Delete a note by ID from MongoDB
func (s *Store) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	result, err := s.db.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("note %s not found", id)
	}
	return nil
}

// DeleteByTags deletes notes by tags from MongoDB
func (s *Store) DeleteByTags(tags []string) error {
	filter := bson.M{"tags": bson.M{"$in": tags}}
	result, err := s.db.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no notes found with tags %v", tags)
	}
	return nil
}

// Search for notes by tags, description or command from MongoDB
func (s *Store) Search(searchWords []string, searchLocations []string) ([]storage.Note, error) {
	filterLocs := []bson.M{}
	for _, searchLocation := range searchLocations {
		for _, searchWord := range searchWords {
			pattern := primitive.Regex{Pattern: regexp.QuoteMeta(searchWord)}
			switch searchLocation {
			case "command":
				filterLocs = append(filterLocs, bson.M{"command": pattern})
			case "description":
				filterLocs = append(filterLocs, bson.M{"description": pattern})
			case "tags":
				filterLocs = append(filterLocs, bson.M{"tags": pattern})
			}
		}
	}

	return s.find(bson.M{"$or": filterLocs})
}

func (s *Store) find(filter bson.M) ([]storage.Note, error) {
	cur, err := s.db.Find(context.Background(), filter, options.Find().SetSort(bson.D{{Key: "_id", Value: 1}}))
	if err != nil {
		return nil, err
	}
	var notes []Note
	if err := cur.All(context.Background(), &notes); err != nil {
		return nil, err
	}
	result := make([]storage.Note, len(notes))
	for i, note := range notes {
		result[i] = note.toStorage()
	}
	return result, nil
}

func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}
