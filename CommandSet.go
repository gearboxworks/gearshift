package barnacle

import (
	"net/http"
	"bytes"
	"os/exec"
	"fmt"
	"strings"
)

type CommandSet struct {
	Path string
}

func NewCommandSet(p string) *CommandSet {
	cs := &CommandSet{ Path: p }
	return cs
}

func (c *CommandSet) Handler() http.Handler {
	// Handles an HTTP-ish request made to Barnacle server
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		var jr JsonResponse
		var sc int
		var errbuf bytes.Buffer
		var out []byte
		cm := CommandMapper{
			RequestUri:   r.RequestURI,
			CommandsPath: c.Path,
			Method:       r.Method,
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
				errMsg = fmt.Sprintf(errMsg, r.RequestURI, r.Method)
			} else {
				errMsg = strings.Trim(errbuf.String(), "\n")
			}
			jr = NewFailJsonResponse(errMsg)
		} else {
			body := strings.Trim(string(out), "\n")
			jr = NewOkJsonResponse(body)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(sc)
		w.Write(jr.ToByteArray())
	});
}
