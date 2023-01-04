package mongo

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Note struct {
	Id          primitive.ObjectID `bson:"_id,omitempty"`
	Tags        []string           `bson:"tags"`
	Command     string             `bson:"command"`
	Description string             `bson:"description"`
}

type Config struct {
	Host       string
	Port       int
	Database   string
	Collection string
}

type Store struct {
	db *mongo.Database
}

func Initialize(config *Config) *Store {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB not running", err)
	}
	return &Store{db: client.Database("notes")}
}

func (s *Store) Insert(item interface{}) error {
	_, err := s.db.Collection("notes").InsertOne(context.Background(), item)
	return err
}

func (s *Store) Get(id string) (interface{}, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	return s.db.Collection("notes").FindOne(context.Background(), filter).DecodeBytes()
}

func (s *Store) GetByTags(tags []string) (interface{}, error) {
	notes := make(map[string]Note)
	for _, tag := range tags {
		filter := bson.M{"tags": bson.M{"$in": []string{tag}}}
		cur, err := s.db.Collection("notes").Find(context.Background(), filter)
		if err != nil {
			return nil, err
		}
		for cur.Next(context.Background()) {
			var n Note
			if err := cur.Decode(&n); err != nil {
				return nil, err
			}
			notes[n.Id.String()] = n
		}
	}

	var result []Note
	for _, v := range notes {
		result = append(result, v)
	}

	return result, nil
}

func (s *Store) GetAll() (interface{}, error) {
	notes := make(map[string]Note)
	filter := bson.M{}
	cur, err := s.db.Collection("notes").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		notes[n.Id.String()] = n
	}

	var result []Note
	for _, v := range notes {
		result = append(result, v)
	}

	return result, nil
}

func (s *Store) Delete(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	_, err = s.db.Collection("notes").DeleteOne(context.Background(), filter)
	return err
}

func (s *Store) DeleteByTags(tags []string) error {
	for _, tag := range tags {
		filter := bson.M{"tags": bson.M{"$in": []string{tag}}}
		_, err := s.db.Collection("notes").DeleteMany(context.Background(), filter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Store) Search(searchWord string, searchLocations []string) (interface{}, error) {
	fmt.Println("Searching: ", searchWord)
	fmt.Println("In: ", searchLocations)

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
	cur, err := s.db.Collection("notes").Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		var n Note
		if err := cur.Decode(&n); err != nil {
			return nil, err
		}
		notes[n.Id.String()] = n
	}

	var result []Note
	for _, v := range notes {
		result = append(result, v)
	}

	return result, nil
}
