package gearshift

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

type CommandMapper struct{
	RequestUri   string
	CommandsPath string
	Method       string
}

// MethodPath returns path to HTTP Method's commands
func (cm *CommandMapper) methodPath() string {
	return cm.CommandsPath + "/" + cm.Method
}

// Filepath translates the request URI into a Command FilePath cfp
// for the given methodPath mp. It does this but parsing the URI and
// comparing to the directory entries.
func (cm *CommandMapper) Filepath() (cfp string, err error) {
	// Ensure URI has a leading slash
	cm.RequestUri = "/" + strings.TrimLeft(cm.RequestUri, "/")
	// Split into URL parts
	up := strings.Split(cm.RequestUri, "/")
	if len(up) == 2 && "" == up[1] {
		// If the request is a root request give it a root filename
		up[1] = "<root>"
	}

	// Create a filepath cfp starting with methodPath
	for i := 1; i < len(up); i++ {
		// Add a path segment from the requested URI to see if it validates
		cfp = cm.methodPath() + "/" + up[i]
		isde := isDirEntry(cfp)
		if !isde {
			// Grab all files but ignore "." and ".." and ".whatever"
			files, err := filepath.Glob("^[^.]")
			if err != nil {
				// If is error, this will fall through to assign err below
			} else {
				for _, file := range files {
					// Check for a template "variable", e.g. "{extensions}"
					if '{' != getFirstChar(file) {
						continue
					}
					// There should only be one variable per directory
					cfp = cm.methodPath() + "/" + file
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
			cfp = cfp
			break
		}

		// Not a Directory Entry so we have an error
		msg := fmt.Sprintf("%s is not a valid resource", cm.RequestUri)
		err = errors.New(msg)

		// Break out now as we no longer need
		// to process remaining URL parts
		break
	}
	return cfp, err
}

