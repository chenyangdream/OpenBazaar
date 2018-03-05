package api

import (
	"net/http"
	"fmt"
	"github.com/Saturn/saturn-go/ipfs"
	"github.com/Saturn/saturn-go/core"
	"os/exec"
	"os"
	"path/filepath"
	"time"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Add request")
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dirPath := filepath.Dir(path)
	testFilePath := dirPath + "/testfiles/helloworld.txt"
	fileHash, _ := ipfs.AddFile(core.Node.Context, testFilePath)
	fmt.Println("file ", testFilePath, " hash  = ", fileHash)
}

func CatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Cat request")
	fileHash := "zb2rhnWiqWEjwxAvdtx5V2j4cgG7BoFAmsGHNcKdYN7hb8Lqr"
	dataText, err := ipfs.Cat(core.Node.Context, fileHash, time.Second*120)
	if err != nil {
		fmt.Println(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("file hash %s filedata %s", fileHash, dataText)

}


func PeersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Peers request")
	peers, err := ipfs.ConnectedPeers(core.Node.Context)
	if err != nil {
		fmt.Println(http.StatusInternalServerError, err.Error())
		return
	}

	for i, peer := range peers {
		fmt.Printf("peer[%d] = %s\n", i, peer)
	}
}

