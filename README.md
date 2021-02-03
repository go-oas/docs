# docs

## _go-OAS/docs_

![GolangCI](https://github.com/go-oas/docs/workflows/golangci/badge.svg?branch=main)
![Build](https://github.com/go-oas/docs/workflows/Golang/badge.svg?branch=main)
[![Version](https://img.shields.io/badge/version-v0.0.1-green.svg)](https://github.com/go-oas/docs/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-oas/docs)](https://goreportcard.com/report/github.com/go-oas/docs)
[![Coverage Status](https://coveralls.io/repos/github/go-oas/docs/badge.svg?branch=main)](https://coveralls.io/github/go-oas/docs?branch=main)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-oas/docs.svg)](https://pkg.go.dev/github.com/go-oas/docs)

go-OAS Docs converts structured OAS3 (Swagger3) objects into OAS3/Swagger documentation & automatically serves it on
chosen route and port. It's extremely flexible and simple, so basically it can be integrated into any framework or
existing project.

## Getting Started

1. Download **_docs_** by using:
   ```sh
   $ go get -u github.com/go-oas/docs
   ``` 
2. Add one line annotation to the handler you wish to use in the following
   format: `// @OAS <FUNC_NAME> <ROUTE> <HTTP_METHOD>`
   Examples:
   ```
   // @OAS handleCreateUser /users POST
   // @OAS handleGetUser /users GET
   ```
3. Declare all required documentation elements that are shared. Or reuse ones that already exist in the examples
   directory.
4. Declare specific docs elements per route.

----

## How to use

For more explicit example, please refer to [docs/examples](https://github.com/go-oas/docs/examples)

Add OAS TAG to your existing handler that handles fetching of a User:

```go
package users

import "net/http"

// @OAS handleGetUser /users GET
func (s *service) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
```

Create a unique API documentation function for that endpoint:

```go
package main

import "github.com/go-oas/docs"

func handleGetUserRoute(oasPathIndex int, oas *docs.OAS) {
	path := oas.GetPathByIndex(oasPathIndex)

	path.Summary = "Get a User"
	path.OperationID = "getUser"
	path.RequestBody = docs.RequestBody{}
	path.Responses = docs.Responses{
		getResponseOK(),
	}

	path.Tags = append(path.Tags, "pet")
}
```

Bear in mind that creating a unique function per endpoint handler is not required, but simply provides good value in
usability of shared documentation elements.

Once you created the function, simply register it for parsing by using `AttachRoutes()` defined upon `OAS` structure.
E.g.:

```go

package main

import (
	docs "github.com/go-oas/docs"
)

func main() {
	apiDoc := docs.New()
	apiDoc.AttachRoutes([]interface{}{
		handleGetUserRoute,
	})
```

If this approach is too flexible for you, you are always left with the possibility to create your own attacher - or any
other parts of the system for that matter.

### Examples

To run examples, and checkout hosted documentation via Swagger UI, issue the following command:

```sh
$ go run ./examples/*.go
```

And navigate to `http://localhost:3005/docs/api/` in case that you didn't change anything before running the example
above.

----

## Roadmap & Project board

Check out the current [Project board](https://github.com/go-oas/docs/projects/1) for more information about the first
alpha release. Note: Board is still in its early phase.