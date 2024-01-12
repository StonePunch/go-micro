package main

import (
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	webPort  = 80
	mongoURL = "mongodb://mongo:27017"
)

var client *mongo.Client

type Config struct {
}

func main() {
}
