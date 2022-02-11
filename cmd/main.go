package main

import (
	"context"
	"log"
	"net/http"
	"storage/handler"
	"storage/store"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	ctx := context.Background()
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	mongoClient.Connect(context.Background())
	store.NewDataStorage(mongoClient)
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.IndexHandler)
	mux.HandleFunc("/upload", handler.Upload)
	mux.HandleFunc("/file/", handler.GetDocumentHandler)
	if err := http.ListenAndServe(":4500", mux); err != nil {
		log.Fatal(err)
	}
}
