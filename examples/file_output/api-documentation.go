package main

import "github.com/go-oas/docs"

func handleCreateUserRoute(oasPathIndex int, oas *docs.OAS) {
	path := oas.GetPathByIndex(oasPathIndex)

	path.Summary = "Create a new User"
	path.OperationID = "createUser"

	path.RequestBody = docs.RequestBody{
		Description: "Create a new User",
		Content: docs.ContentTypes{
			getContentApplicationJSON("#/components/schemas/User"),
		},
		Required: true,
	}

	path.Responses = docs.Responses{
		getResponseNotFound(),
		getResponseOK(),
	}

	path.Security = docs.SecurityEntities{
		docs.Security{
			AuthName:  "petstore_auth",
			PermTypes: []string{"write:users", "read:users"},
		},
	}

	path.Tags = append(path.Tags, "user")
}

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
