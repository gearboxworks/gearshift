package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"bytes"
)

var commandsPath string

func main() {
	var err error
	// test to see if our command's path is value
	commandsPath, err = RequestedInterface().CommandsPath()
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
	cm := CommandMapper{
		RequestUri:   request.RequestURI,
		CommandsPath: commandsPath,
		Method:       request.Method,
	}
	cfp, err := cm.Filepath()
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
		var errMsg string
		if sc == 404 {
			errMsg = "The URL path %s is not a valid resource for method %s."
			errMsg = fmt.Sprintf(errMsg, request.RequestURI, request.Method)
		} else {
			errMsg = strings.Trim(errbuf.String(), "\n")
		}
		jr = NewFailJsonResponse(errMsg)
	} else {
		body := strings.Trim(string(out), "\n")
		jr = NewOkJsonResponse(body)
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(sc)
	response.Write(jr.ToByteArray())
}

func RequestedInterface() *Interface {
	sbi := os.Args[1]
	return NewInterface(sbi)
}
