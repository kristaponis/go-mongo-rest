package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

var client *mongo.Client

// getPersonHandler gets all items in the collection
// /person GET
func getPersonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var persons []person
	coll := client.Database("personsdb").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var p person
		cursor.Decode(&p)
		persons = append(persons, p)
	}
	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(persons)
}

// createPersonHandler creates new item in the collection
// /person POST
func createPersonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var p person
	json.NewDecoder(r.Body).Decode(&p)
	coll := client.Database("personsdb").Collection("people")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, _ := coll.InsertOne(ctx, p)
	json.NewEncoder(w).Encode(result)
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	opts := options.Client().ApplyURI("mongodb://localhost:27017")
	client, _ = mongo.Connect(ctx, opts)

	router := mux.NewRouter()
	router.HandleFunc("/person", getPersonHandler).Methods(http.MethodGet)
	router.HandleFunc("/person", createPersonHandler).Methods(http.MethodPost)

	log.Println("Starting application...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
