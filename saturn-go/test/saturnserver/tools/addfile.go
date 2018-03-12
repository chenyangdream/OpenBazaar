package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type filePath struct {
	FilePath string `json:"filepath"`
}

type fileHash struct {
	FileHash string `json:"filehash"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage:%s FileAbsPath\n", os.Args[0])
		return
	}

	path := os.Args[1]
	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("file %s does not exists\n", path)
		return
	}

	var fp filePath
	fp.FilePath = path

	b, err := json.Marshal(fp)
	body := bytes.NewBuffer([]byte(b))

	res, err := http.Post("http://localhost:4002/saturn/add", "application/json;charset=utf-8", body)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

	var fh fileHash
	err = json.Unmarshal(result, &fh)
	if err != nil {
		fmt.Printf("decode server response error %s\n", err.Error())
		return
	}

	fmt.Printf("add file %s ok.\nHash %s\n", path, fh.FileHash)
}
