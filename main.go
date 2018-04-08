package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"bytes"
)

var commandsPath string

func main() {
	var err error
	// test to see if our command's path is value
	commandsPath, err = RequestedBarnacleInterface().CommandsPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	srv := &http.Server{Addr: ":9970"}
	http.HandleFunc("/", RequestHandler)

	fmt.Println("Listening on port 9970...")
	if err = srv.ListenAndServe(); err != nil {
		// Don't panic, this probably is an intentional close
		log.Printf("Httpserver: ListenAndServe() error: %s", err)
	}

	err = srv.Shutdown(nil)
	if err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}

}

// RequestHandler handles an HTTP-ish request made to Barnacle server
func RequestHandler(response http.ResponseWriter, request *http.Request) {
	var jr JsonResponse
	var sc int
	var errbuf bytes.Buffer
	var out []byte
	methodPath := commandsPath + "/" + request.Method
	cfp, err := getCommandFilepath(request.RequestURI, methodPath)
	if err != nil {
		sc = 404
	} else {
		for {
			cmd := exec.Command(cfp)
			cmd.Stderr = &errbuf
			out, err = cmd.Output()
			if err != nil {
				sc = 500
				break
			}
			sc = 200
			break
		}
	}
	if err != nil {
		errMsg := strings.Trim( errbuf.String(), "\n")
		jr = NewJsonResponse("fail", errMsg, "")
	} else {
		body := strings.Trim( string(out), "\n")
		jr = NewJsonResponse("ok", "", body)
	}

	response.WriteHeader(sc)
	json.NewEncoder(response).Encode(jr)

}

type BarnacleInterface struct {
	namespace string
	name      string
}

func NewBarnacleInterface(sbi string) *BarnacleInterface {
	bi := new(BarnacleInterface)
	bi.parse(sbi)
	return bi
}

func (bi *BarnacleInterface) QualifiedName() string {
	return bi.namespace + "/" + bi.name
}

func (bi BarnacleInterface) CommandsPath(leaf ...string) (path string, err error) {
	l := ""
	if len(leaf) == 1 {
		l = "/" + strings.TrimLeft(leaf[0], "/")
	}
	d, err := os.Getwd()
	return d + "/files/etc/barnacle/" + bi.QualifiedName() + l, err
}

func (bi *BarnacleInterface) parse(sbi string) {
	var parts []string = strings.Split(sbi, "/")
	bi.namespace = parts[0]
	bi.name = parts[1]
}

// isDirEntry checks the filepath fp to determine if the file exists
func isDirEntry(fp string) bool {
	_, err := os.Stat(fp)
	return !os.IsNotExist(err)
}

func RequestedBarnacleInterface() *BarnacleInterface {
	sbi := os.Args[1]
	return NewBarnacleInterface(sbi)
}

// getCommandFilepath translates the URI u into Command FilePath cfp
// for the given methodPath mp. It does this but parsing the URI and
// comparing to the directory entries.
func getCommandFilepath(u string, mp string) (cfp string, err error) {
	// Ensure URI has a leading slash
	u = "/" + strings.TrimLeft(u, "/")
	// Split into URL parts
	up := strings.Split(u, "/")
	// Create a filepath fp starting with methodPath mp
	fp := mp
	for i := 1; i < len(up); i++ {
		// Add a path segment from the requested URI to see if it validates
		fp = mp + "/" + up[i]
		isde := isDirEntry(fp)
		if !isde {
			// Grab all files but ignore "." and ".." and ".whatever"
			files, err := filepath.Glob("^[^.]")
			if err != nil {
				// If is error, this will fall through to assign err below
			} else {
				for _, file := range files {
					// Check for a template "variable", e.g. "{domain}"
					if '{' != getFirstChar(file) {
						continue
					}
					// There should only be one variable per directory
					fp = mp + "/" + file
					// Set isDirEntry to true to see if we have a match
					isde = true
					break
				}
			}

		}
		if isde {
			// So the current loop is a Directory Entry...
			if i+1 < len(up) {
				// ...but we've got more URL parts to process
				continue
			}
			// ...we are done processing.
			// We have a match, let's assign and exit
			cfp = fp
			break
		}
		// Not a Directory Entry so we have an error
		msg := fmt.Sprintf("%s is not a valid resource", u)
		err = errors.New(msg)
		// Exit now as no more need to process remaining URL parts
		break
	}
	return cfp, err
}

func getFirstChar(s string) byte {
	var ch byte
	for _, ch = range []byte(s) {
		break
	}
	return ch
}

type JsonResponse interface {
}

type OkJsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Body    interface{} `json:"body"`
}
type FailJsonResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
}

// NewJsonResponse() is a constructor for JsonResponse
func NewJsonResponse(s, m, b string) JsonResponse {
	var jr JsonResponse
	js := new(interface{})
	json.Unmarshal([]byte(b),&js)
	if 0 < len(b) {
		jr = &OkJsonResponse{
			Status:  s,
			Message: m,
			Body:    js,
		}
	} else {
		jr = &FailJsonResponse{
			Status:  s,
			Message: m,
		}
	}
	return jr
}

