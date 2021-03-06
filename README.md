# docs

### Automatically generate RESTful API documentation for GO projects - aligned with [Open API Specification standard](https://github.com/OAI/OpenAPI-Specification/blob/master/versions/3.0.3.md).

<img align="right" width="180px" src="https://raw.githubusercontent.com/kaynetik/dotfiles/master/svg-resources/go-grpc-web.svg">

![GolangCI](https://github.com/go-oas/docs/workflows/golangci/badge.svg?branch=main)
![Build](https://github.com/go-oas/docs/workflows/Build/badge.svg?branch=main)
[![Version](https://img.shields.io/badge/version-v1.0.5-success.svg)](https://github.com/go-oas/docs/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-oas/docs)](https://goreportcard.com/report/github.com/go-oas/docs)
[![Coverage Status](https://coveralls.io/repos/github/go-oas/docs/badge.svg?branch=main)](https://coveralls.io/github/go-oas/docs?branch=main)
[![codebeat badge](https://codebeat.co/badges/32b86556-84e3-4db9-9f11-923d12994f90)](https://codebeat.co/projects/github-com-go-oas-docs-main)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-oas/docs.svg)](https://pkg.go.dev/github.com/go-oas/docs)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://awesome-go.com)

go-OAS Docs converts structured OAS3 (Swagger3) objects into the Open API Specification & automatically serves it on
chosen route and port. It's extremely flexible and simple, so basically it can be integrated into any framework or
existing project.

We invite anyone interested to join our **[GH Discussions board](https://github.com/go-oas/docs/discussions)**. Honest
feedback will enable us to build better product and at the same time not waste valuable time and effort on something
that might not fit intended usage. So if you can, please spare few minutes to give your opinion of what should be done
next, or what should be the priority for our roadmap. :muscle: :tada:

----

## Table of Contents

- [Getting Started](#getting-started)
- [How to use](#how-to-use)
    * [Examples](#examples)
- [Contact](#contact)
- [The current roadmap (planned features)](#roadmap)

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
	"github.com/go-oas/docs"
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

## Contact

Check out the current [Project board](https://github.com/go-oas/docs/projects/1) or
our **[GH Discussions board](https://github.com/go-oas/docs/discussions)** for more information.

You can join our Telegram group at: [https://t.me/go_oas](https://t.me/go_oas)

## Roadmap

| Feature (GH issues)                                             | Description                                                                              | Release |
| --------------------------------------------------------------- | ---------------------------------------------------------------------------------------- | ------- |
| [Validation](https://github.com/go-oas/docs/issues/17)          | Add validation to all structures based on OAS3.0.3                                       | v1.1.0  |
| [CLI](https://github.com/go-oas/docs/issues/18)                 | Add CLI support - make it CLI friendly                                                   | v1.2.0  |
| [Postman](https://github.com/go-oas/docs/issues/19)             | Add postman support via PM API                                                           | v1.3.0  |
| [ReDoc](https://github.com/go-oas/docs/issues/20)               | Add ReDoc support as an alternative to SwaggerUI                                         | v1.4.0  |
| [E2E Auto-generation](https://github.com/go-oas/docs/issues/21) | Go tests conversion to Cypress/Katalon suites (convert mocked unit tests into e2e tests) | v1.5.0  |

