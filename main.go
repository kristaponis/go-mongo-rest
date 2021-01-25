package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Person type contains some basic info
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, opts)

	router := mux.NewRouter()

	log.Println("Starting application...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
