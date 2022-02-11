package controller

import (
	"fmt"
	"io"
	"log"
	"os"
)

const dir string = "blobs/"

func WriteToBlob(b [][]byte, name string) {

	for i, bytes := range b {
		blobName := dir + name + "." + fmt.Sprint(i)

		err := os.WriteFile(blobName, bytes, 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func SplitBytes(f *os.File) error {
	fileStats, err := f.Stat()
	if err != nil {
		log.Println(err)
		return err
	}
	bs := make([][]byte, 0)
	for i := int64(0); i <= (fileStats.Size() / 1024); i++ {
		bytes := make([]byte, 1024)
		_, err := f.Read(bytes)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println(err)
			return err
		}
		_, err = f.Seek((i*1024)+1, 0)
		if err != nil {
			log.Println(err)
			return err
		}
		bs = append(bs, bytes)
	}

	WriteToBlob(bs, fileStats.Name())

	return nil
}

// func ReadFile(name string) {
// 	blobName := dir + name + "."
// 	dir, err := os.ReadDir(dir)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
