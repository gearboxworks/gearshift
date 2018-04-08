package main

import (
	"net/http"
	"os"
	"fmt"
	"os/exec"
	"log"
	"io/ioutil"
	"io"
)

var barnacleInterface *BarnacleInterface

func main() {
	var err error
	barnacleInterface = requestedBarnacleInterface()
	// test to see if our command's path is value
	_,err = barnacleInterface.CommandsPath()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", RequestHandler)
	http.ListenAndServe(":9970", mux)

}
// RequestHandler handles an HTTP-ish request made to Barnacle server
func RequestHandler(response http.ResponseWriter, request *http.Request) {
	var jr *JsonResponse
	var sc int
	var msg []byte
	var stderr io.ReadCloser
	var stdout io.ReadCloser
	methodPath,err := barnacleInterface.CommandsPath(request.Method)
	cfp, err := getCommandFilepath(request.RequestURI,methodPath)
	if err != nil {
		sc = 404
	} else {
		for{
			cmd := exec.Command(cfp)
			stderr, err = cmd.StderrPipe()
			if err != nil {
				log.Fatal(err)
			}
			stdout, err = cmd.StdoutPipe()
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Start()
			if err != nil {
				log.Fatal(err)
			}
			err = cmd.Wait()
			if err != nil {
				sc = 500
				break
			}
			sc = 200
		}
	}
	if err != nil {
		msg, _ = ioutil.ReadAll(stderr)
		jr = NewJsonResponse("fail",string(msg))
	} else {
		msg, _ = ioutil.ReadAll(stdout)
		jr = NewJsonResponse("ok",string(msg))
	}
	response.WriteHeader(sc)
	response.Write([]byte(jr.toString()))
}


