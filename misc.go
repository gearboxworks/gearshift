package main

import (
	"os"
	"fmt"
	"strings"
	"errors"
	"path/filepath"
)

// isDirEntry checks the filepath fp to determine if the file exists
func isDirEntry(fp string) bool {
	_,err := os.Stat(fp)
	return ! os.IsNotExist(err)
}

func requestedBarnacleInterface() *BarnacleInterface {
	sbi := os.Args[1]
	return NewBarnacleInterface(sbi)
}

// getCommandFilepath translates the URI u into Command FilePath cfp
// for the given methodPath mp. It does this but parsing the URI and
// comparing to the directory entries.
func getCommandFilepath(u string, mp string) (cfp string,err error) {
	// Ensure URI has a leading slash
	u = "/" + strings.TrimLeft(u,"/")
	// Split into URL parts
	up := strings.Split(u,"/")
	// Create a filepath fp starting with methodPath mp
	fp := mp
	for i:=0; i<len(up); i++ {
		// Add a path segment from the requested URI to see if it validates
		fp = mp + "/" + up[i]
		isde := isDirEntry(fp)
		if ! isde {
			// Grab all files but ignore "." and ".." and ".whatever"
			files,err := filepath.Glob("^[^.]")
			if err != nil {
				// If is error, this will fall through to assign err below
			} else {
				for _, file := range files {
					// Check for a template "variable", e.g. "{domain}"
					if "{" != getFirstChar( file ) {
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
			if i + 1 < len(up) {
				// ...but we've got more URL parts to process
				continue
			}
			// ...we are done processing.
			// We have a match, let's assign and exit
			cfp = fp
			break
		}
		// Not a Directory Entry so we have an error
		msg := fmt.Sprintf("%s is not a valid resource",u)
		err = errors.New(msg)
		// Exit now as no more need to process remaining URL parts
		break
	}
	return cfp,err
}

func getFirstChar(s string) string {
	var c string
	for _,c = range s {
		break
	}
	return c
}
