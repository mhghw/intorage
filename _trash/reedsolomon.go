package controller

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"log"
	"os"

	"github.com/ory/viper"
)

type DataInfo struct {
	FileName    string
	Size        int64
	HashedBlobs map[string][]byte
}

func DataSplitter(filename string) (*DataInfo, error) {
	loc := fmt.Sprintf("../%v/%v", viper.GetString("dir.tmp"), filename)

	file, err := os.Open(loc)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	di := &DataInfo{
		FileName: filename,
		Size:     fileInfo.Size(),
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 256<<10)
	sha1 := sha1.New()

	for {
		_, err = reader.Read(buffer)
		if err != nil {
			return nil, err
		}
	}

}
