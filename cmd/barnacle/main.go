//
// See: https://medium.com/@benbjohnson/structuring-applications-in-go-3b04be4ff091
//
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/pressboxx/barnacle"
)

func main() {
	var err error
	var bi = RequestedInterface()
	var cp = TrapError(bi.CommandsPath()).(string)
	var cs = barnacle.NewCommandSet(cp)
	srv := &http.Server{Addr: ":9970"}
	http.Handle("/", cs.Handler())
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

func TrapError(v interface{},err error) interface{} {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return v
}

func RequestedInterface() *barnacle.Interface {
	sbi := os.Args[1]
	return barnacle.NewInterface(sbi)
}


