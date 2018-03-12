package api

import (
	"net/http"
	"net/url"
	"net/http/httputil"
	"runtime/debug"
	"strings"
	"encoding/json"
	"fmt"
	"github.com/Saturn/saturn-go/cmd"
	"github.com/Saturn/saturn-go/core"
	"os/exec"
	"os"
	"path/filepath"
	"time"
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
	fmt.Println("add file")
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	dirPath := filepath.Dir(path)
	testFilePath := dirPath + "/resource/helloworld.txt"
	fileHash, err := ipfscmd.AddFile(i.node.Context, testFilePath)
	if err != nil {
		fmt.Println("add file ", path, " failed! error ", err)
	}
	fmt.Println("add file ", testFilePath, " hash ", fileHash)
}

func (i *jsonAPIHandler) GETPeers(w http.ResponseWriter, r *http.Request) {
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
	peerid, err := i.node.GetPeerId()
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("PeerId = ", peerid)
}

func (i *jsonAPIHandler) GETFileContent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Cat request")
	fileHash := "zb2rhnWiqWEjwxAvdtx5V2j4cgG7BoFAmsGHNcKdYN7hb8Lqr"
	dataText, err := ipfscmd.Cat(i.node.Context, fileHash, time.Second*120)
	if err != nil {
		fmt.Println(http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Printf("file hash %s filedata %s", fileHash, dataText)

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
