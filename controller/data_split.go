package controller

import (
	"context"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"storage/store"
)

func SplitAndStore(ctx context.Context, fileLocation string) (*store.Document, error) {
	f, err := os.Open(fileLocation)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer f.Close()

	fileInfo, err := f.Stat()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	//sha1 hasher obj
	hasher := sha1.New()

	//new document to store file info in db
	newDoc := &store.Document{
		Name: fileInfo.Name(),
	}

	//main loop for spliting data and storing them
	//buffer limit is 256kB
	for i := 0; ; i++ {
		buf := make([]byte, 256<<10)
		n, err := f.Read(buf)
		if n < 256<<10 {
			buf = buf[:n]
		}
		log.Println(len(buf))

		if err != nil {

			if err != io.EOF {

				log.Fatal(err)
			}

			break
		}
		hasher.Write(buf)
		bs := hasher.Sum(nil)
		hasher.Reset()
		go store.Ds.ForceWriteHashedObject(ctx, &store.HashObject{
			Hash: fmt.Sprintf("%x", bs),
			Data: buf,
		})
		newDoc.Hashes = append(newDoc.Hashes, fmt.Sprintf("%x", bs))
	}

	err = store.Ds.InsertDocument(ctx, newDoc)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return newDoc, nil
}
