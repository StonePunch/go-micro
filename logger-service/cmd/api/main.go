package main

import (
	"context"
	"fmt"
	"log"
	"logger-service/data"
	"net/http"
	"time"
	"tools"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	webPort  = 80
	mongoURL = "mongodb://mongo:27017"
)

type Config struct {
	tools.Tools
	Models data.Models
}

func main() {
	// connect to mongo
	client, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	// close connection
	defer func() {
		log.Println("Disconnecting from mongo...")

		if err = client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting from mongo:", err)
			panic(err)
		}

		log.Println("Successfully disconnected from mongo")
	}()

	app := Config{
		Tools:  tools.New(),
		Models: data.New(client),
	}

	log.Println("Starting logger-service on port:", webPort)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", webPort),
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	log.Println("Connecting to mongo...")

	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect to mongo
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo:", err)
		return nil, err
	}

	log.Println("Successfully connected to mongo")

	return client, nil
}
