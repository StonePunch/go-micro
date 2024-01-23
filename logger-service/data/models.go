package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func New(mongo *mongo.Client) Models {
	client = mongo

	return Models{
		LogEntry: LogEntry{},
	}
}

func (l *LogEntry) Insert(entry LogEntry) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(ctx, LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting into logs:", err)
		return err
	}

	return nil
}

func (l *LogEntry) All() ([]*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Error retrieving all logs:", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var item LogEntry

		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log into slice:", err)
			return nil, err
		}

		logs = append(logs, &item)
	}

	return logs, nil
}

func (l *LogEntry) GetOne(id string) (*LogEntry, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	bsonID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting id into bson:", err)
		return nil, err
	}

	filter := bson.M{"_id": bsonID}

	var entry LogEntry
	err = collection.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		log.Printf("Error retrieving log with bsonID of %s: %v\n", bsonID, err)
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	bsonID, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error converting id into bson:", err)
		return nil, err
	}

	filter := bson.M{"_id": bsonID}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: l.Name},
			{Key: "data", Value: l.Data},
			{Key: "updated_at", Value: time.Now()},
		}},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating log with bsonID of %s: %v\n", bsonID, err)
		return nil, err
	}

	return result, nil
}

func (l *LogEntry) DropCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping collection:", err)
		return err
	}

	return nil
}
