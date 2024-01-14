package router

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/hingew/hsfl-master-ai-cloud-engineering/lib/health"
)

type Router struct{}

const DIR = "./public/"

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/health/web" {
		health.Check(w, r)
	} else {
		fs := http.FileServer(http.Dir(DIR))

		// If the requested file exists then return if; otherwise return index.html (fileserver default page)
		if r.URL.Path != "/" {
			fullPath := DIR + strings.TrimPrefix(path.Clean(r.URL.Path), "/")
			_, err := os.Stat(fullPath)
			if err != nil {
				if !os.IsNotExist(err) {
					panic(err)
				}
				// Requested file does not exist so we return the default (resolves to index.html)
				r.URL.Path = "/"
			}
		}

		fs.ServeHTTP(w, r)
	}
}
