package store

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var Ds DataStorage

type HashObject struct {
	Hash string `bson:"_id"`
	Data []byte `bson:"data"`
}

type Document struct {
	Name   string   `bson:"_id" json:"name"`
	Hashes []string `bson:"hashes"`
}

type DataStorage struct {
	client *mongo.Client
}

func NewDataStorage(client *mongo.Client) {
	Ds.client = client
}

func (s DataStorage) hashCollection() *mongo.Collection {
	return s.client.Database("pr").Collection("hash")
}
func (s DataStorage) docCollection() *mongo.Collection {
	return s.client.Database("pr").Collection("document")
}

func (s DataStorage) ForceWriteHashedObject(ctx context.Context, ho *HashObject) {
	s.hashCollection().InsertOne(ctx, ho)
}

func (s DataStorage) InsertDocument(ctx context.Context, doc *Document) error {
	_, err := s.docCollection().InsertOne(ctx, doc)
	return err
}

func (s DataStorage) GetAndWriteFile(ctx context.Context, name string) error {
	filter := bson.M{"_id": name}
	var doc Document
	err := s.docCollection().FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		log.Println(err)
		return err
	}

	err = os.MkdirAll("./fileTemp", os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	destFile, err := os.Create(fmt.Sprintf("./fileTemp/%v", name))
	if err != nil {
		log.Println(err)
		return err
	}
	defer destFile.Close()

	for _, hash := range doc.Hashes {
		err = s.ReadHashData(ctx, hash, destFile)
		if err != nil {
			log.Println(err)
			os.Remove(destFile.Name())
			return err
		}
	}

	return nil
}

func (s DataStorage) ReadHashData(ctx context.Context, hash string, wr io.Writer) error {
	filter := bson.M{"_id": hash}
	var hashObj HashObject
	err := s.hashCollection().FindOne(ctx, filter).Decode(&hashObj)
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = wr.Write(hashObj.Data)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
