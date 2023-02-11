package mongo

import (
	"context"
	"log"
	"strconv"

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

// Initialize MongoDB client
func Initialize(config *Config) *Store {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://"+config.Host+":"+strconv.Itoa(config.Port)))
	if err != nil {
		log.Fatal("MongoDB not running", err)
	}
	return &Store{db: client.Database(config.Database).Collection(config.Collection)}
}

// Insert a note into MongoDB
func (s *Store) Insert(item interface{}) error {
	_, err := s.db.InsertOne(context.Background(), item)
	return err
}

// Get a note from MongoDB
func (s *Store) Get(id string) (interface{}, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	return s.db.FindOne(context.Background(), filter).DecodeBytes()
}

// GetByTags gets notes by tags from MongoDB
func (s *Store) GetByTags(tags []string) (interface{}, error) {
	notes := make(map[string]Note)
	for _, tag := range tags {
		filter := bson.M{"tags": bson.M{"$in": []string{tag}}}
		cur, err := s.db.Find(context.Background(), filter)
		if err != nil {
			return nil, err
		}
		for cur.Next(context.Background()) {
			var n Note
			if err := cur.Decode(&n); err != nil {
				return nil, err
			}
			notes[n.ID.String()] = n
		}
	}

	var result []Note
	for _, v := range notes {
		result = append(result, v)
	}

	return result, nil
}

// GetAll gets all notes from MongoDB
func (s *Store) GetAll() (interface{}, error) {
	notes := make(map[string]Note)
	filter := bson.M{}
	cur, err := s.db.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		notes[n.ID.String()] = n
	}

	var result []Note
	for _, v := range notes {
		result = append(result, v)
	}

	return result, nil
}

// GetTags returns all available tags
func (s *Store) GetTags() ([]string, error) {
	var tags []string
	notesInt, _ := s.GetAll()
	notes, _ := notesInt.([]Note)
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
	_, err = s.db.DeleteOne(context.Background(), filter)
	return err
}

// DeleteByTags deletes notes by tags from MongoDB
func (s *Store) DeleteByTags(tags []string) error {
	for _, tag := range tags {
		filter := bson.M{"tags": bson.M{"$in": []string{tag}}}
		_, err := s.db.DeleteMany(context.Background(), filter)
		if err != nil {
			return err
		}
	}

	return nil
}

// Search for notes by tags, description or command from MongoDB
func (s *Store) Search(searchWord string, searchLocations []string) (interface{}, error) {
	notes := make(map[string]Note)
	filterLocs := []bson.M{}
	for _, searchLocation := range searchLocations {
		switch searchLocation {
		case "command":
			filterLocs = append(filterLocs, bson.M{"command": primitive.Regex{Pattern: searchWord, Options: ""}})
		case "description":
			filterLocs = append(filterLocs, bson.M{"description": primitive.Regex{Pattern: searchWord, Options: ""}})
		case "tags":
			filterLocs = append(filterLocs, bson.M{"tags": primitive.Regex{Pattern: searchWord, Options: ""}})
		}
	}

	filter := bson.M{"$or": filterLocs}
	cur, err := s.db.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		notes[n.ID.String()] = n
	}

	var result []Note
	for _, v := range notes {
		result = append(result, v)
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
