package docs

import (
	"log"
	"net/http"
	"strings"
)

func ServeSwaggerUI(route, port string) {
	directory := "./internal/dist"

	if route == "" {
		route = "/api/"
	}

	fileServer := http.FileServer(FileSystem{http.Dir(directory)})
	http.Handle(route, http.StripPrefix(strings.TrimRight(route, "/"), fileServer))

	log.Printf("Serving %s on HTTP port: %s\n", directory, port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
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

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
