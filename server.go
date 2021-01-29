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
	// fs is wrapped to avoid unwanted dir traversal.
	fs http.FileSystem
}

// Open opens file. Returns http.File, and error if there is any.
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file in path %s :%w", path, err)
	}

	fileInfo, err := f.Stat()
	if err != nil {
		return f, fmt.Errorf("failed to fetch file info :%w", err)
	}

	if fileInfo.IsDir() {
		index := strings.TrimSuffix(path, fwSlashSuffix) + defaultIndexPath
		if _, err = fs.fs.Open(index); err != nil {
			return nil, fmt.Errorf("failed trimming path sufix :%w", err)
		}
	}

	return f, nil
}
