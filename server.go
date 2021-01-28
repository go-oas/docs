package docs

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	defaultRoute     = "/api"
	defaultDirectory = "./internal/dist"
	defaultIndexPath = "/index.html"
	fwSlashSuffix    = "/"
)

func ServeSwaggerUI(route, port string) {
	if route == "" {
		route = defaultRoute
	}

	fileServer := http.FileServer(FileSystem{http.Dir(defaultDirectory)})
	http.Handle(route, http.StripPrefix(strings.TrimRight(route, fwSlashSuffix), fileServer))

	log.Printf("Serving SwaggerIU on HTTP port: %s\n", port)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		// TODO: Add graceful shutdown/handling with err chan.
		panic(err)
	}
}

type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return f, err
	}

	if fileInfo.IsDir() {
		index := strings.TrimSuffix(path, fwSlashSuffix) + defaultIndexPath
		if _, err = fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
