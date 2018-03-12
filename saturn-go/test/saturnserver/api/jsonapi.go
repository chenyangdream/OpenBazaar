package api

import (
	"net/http"
	"net/url"
	"net/http/httputil"
	"runtime/debug"
	"strings"
	"encoding/json"
	"fmt"
	"os/exec"
	"os"
	"path"
	"path/filepath"
	"time"
	"github.com/Saturn/saturn-go"
	"github.com/Saturn/saturn-go/cmd"
	"github.com/Saturn/saturn-go/core"
)

type jsonAPIHandler struct {
	node *core.SaturnNode
}

func newJsonAPIHandler(node *core.SaturnNode)(*jsonAPIHandler, error) {
	i := &jsonAPIHandler{
		node: node,
	}

	return i, nil
}

func (i *jsonAPIHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.Path)
	if err != nil {
		log.Error(err)
		return
	}

	dump, err := httputil.DumpRequest(r, false)

	if err != nil {
		log.Error("Error reading http request:", err)
	}

	log.Debugf("%s", dump)

	defer func() {
		if r := recover(); r != nil {
			log.Error("A panic occurred in the rest api handler!")
			log.Error(r)
			debug.PrintStack()
		}
	}()

	w.Header().Add("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		get(i, u.String(), w, r)
	case "POST":
		post(i, u.String(), w, r)
	}
}

func (i *jsonAPIHandler) POSTAdd(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ipfs add reqeust")
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dirPath := filepath.Dir(path)
	testFilePath := dirPath + "/resource/helloworld.txt"
	//fileHash, err := ipfscmd.AddFile(i.node.Context, testFilePath)
	fileHash, err := ipfs.AddFile(testFilePath)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		fmt.Println("add file ", path, " failed! error ", err)
		return
	}

	fmt.Println("add file ", testFilePath, " hash ", fileHash)
}

func (i *jsonAPIHandler) POSTPin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ipfs pin request")
	_, fileHash := path.Split(r.URL.Path)
	if fileHash == "" {
		fileHash = "zb2rhnWiqWEjwxAvdtx5V2j4cgG7BoFAmsGHNcKdYN7hb8Lqr"
	}
	err := ipfs.PinFile(fileHash)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		fmt.Println("pin file ", fileHash, " failed! error ", err)
		return
	}

	fmt.Println("pin file hash ", fileHash, " ok!")
}

func (i *jsonAPIHandler) POSTUnpin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ipfs unpin request")
	_, fileHash := path.Split(r.URL.Path)
	if fileHash == "" {
		fileHash = "zb2rhnWiqWEjwxAvdtx5V2j4cgG7BoFAmsGHNcKdYN7hb8Lqr"
	}

	err := ipfs.UnpinFile(fileHash)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		fmt.Println("unpin file ", fileHash, " failed! error ", err)
		return
	}

	fmt.Println("unpin file hash ", fileHash, " ok!")
}

func (i *jsonAPIHandler) GETPeers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ipfs swarm peers request")
	peers, err := ipfscmd.ConnectedPeers(i.node.Context)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	peerJson, err := json.MarshalIndent(peers, "", "    ")
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	SanitizeResponse(w, string(peerJson))
}

func (i *jsonAPIHandler) GETPeerId(w http.ResponseWriter, r *http.Request) {
	peerInfo, err := ipfs.GetPeerInfo()
	peerid := peerInfo.PeerId
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("PeerId = ", peerid)
}

func (i *jsonAPIHandler) GETFileContent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ipfs cat request")
	_, fileHash := path.Split(r.URL.Path)
	if fileHash == "" {
		fileHash = "zb2rhnWiqWEjwxAvdtx5V2j4cgG7BoFAmsGHNcKdYN7hb8Lqr"
	}
	dataText, err := ipfscmd.Cat(i.node.Context, fileHash, time.Second*120)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("file hash %s filedata %s", fileHash, dataText)
}

func (i *jsonAPIHandler) GETPinLs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ipfs pin ls request")
	output, err := ipfscmd.PinLs(i.node.Context)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("pin ls output %s", output)
}

func (i *jsonAPIHandler) GETGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ifps get request")
	_, fileHash := path.Split(r.URL.Path)
	if fileHash == "" {
		fileHash = "zb2rhnWiqWEjwxAvdtx5V2j4cgG7BoFAmsGHNcKdYN7hb8Lqr"
	}
	outDir := i.node.RepoPath
	err := ipfs.GetFile(fileHash, outDir, nil)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("get %s to dir %s\n", fileHash, outDir)
}

func ErrorResponse(w http.ResponseWriter, errorCode int, reason string) {
	type ApiError struct {
		Success bool `json:"success"`
		Reason string `json:"reason"`
	}

	reason = strings.Replace(reason, `"`, `'`, -1)
	err := ApiError{false, reason}

	resp, _ := json.MarshalIndent(err, "", "    ")
	w.WriteHeader(errorCode)
	fmt.Fprint(w, string(resp))
}

func SanitizeResponse(w http.ResponseWriter, response string) {
	ret, err := SanitizeJSON([]byte(response))

	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Fprint(w, string(ret))
}