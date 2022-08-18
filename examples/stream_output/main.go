package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-oas/docs"
)

const (
	staticRoute     = "/docs/api/"
	streamRoute     = "/docs/oas/"
	staticDirectory = "./internal/dist"
	port            = 3005
)

func main() {
	apiDoc := docs.New()
	apiSetInfo(&apiDoc)
	apiSetTags(&apiDoc)
	apiSetServers(&apiDoc)
	apiSetExternalDocs(&apiDoc)
	apiSetComponents(&apiDoc)

	apiDoc.AddRoute(docs.Path{
		Route:       "/users",
		HTTPMethod:  "POST",
		OperationID: "createUser",
		Summary:     "Create a new User",
		Responses: docs.Responses{
			getResponseOK(),
			getResponseNotFound(),
		},
		// HandlerFuncName: "handleCreateUser",
		RequestBody: docs.RequestBody{
			Description: "Create a new User",
			Content: docs.ContentTypes{
				getContentApplicationJSON("#/components/schemas/User"),
			},
			Required: true,
		},
	})

	apiDoc.AddRoute(docs.Path{
		Route:       "/users",
		HTTPMethod:  "GET",
		OperationID: "getUser",
		Summary:     "Get a User",
		Responses: docs.Responses{
			getResponseOK(),
		},
		// HandlerFuncName: "handleCreateUser",
		RequestBody: docs.RequestBody{
			Description: "Create a new User",
			Content: docs.ContentTypes{
				getContentApplicationJSON("#/components/schemas/User"),
			},
			Required: true,
		},
	})

	mux := http.NewServeMux()

	// serve static files
	fs := http.FileServer(http.Dir(staticDirectory))
	mux.Handle(staticRoute, http.StripPrefix(staticRoute, fs))

	// serve the oas document from a stream
	mux.HandleFunc(streamRoute, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/yaml")
		if err := apiDoc.BuildStream(w); err != nil {
			http.Error(w, "could not write body", http.StatusInternalServerError)
			return
		}
	})

	// print routes
	// hm := reflect.ValueOf(mux).Elem()
	// fl := hm.FieldByIndex([]int{1})
	// fmt.Printf("%+v\n", fl)

	// path, _ := os.Getwd()
	// fmt.Printf("cwd: %s\n", path)

	fmt.Printf("Listening at :%d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), LogginMiddleware(mux)); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server closed\n")
		} else if err != nil {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}
}
