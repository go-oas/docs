package main

import "github.com/go-oas/docs"

// Args packer could be used to improve readability of such functions,
// and to make the more flexible in order to avoid empty string comparisons.
func getResponseOK(args ...interface{}) docs.Response {
	description := "OK"
	if args != nil {
		description = args[0].(string)
	}

	return docs.Response{
		Code:        200,
		Description: description,
	}
}

func getResponseNotFound() docs.Response {
	return docs.Response{
		Code:        404,
		Description: "Not Found",
		Content: docs.ContentTypes{
			getContentApplicationJSON("#/components/schemas/Pet"),
		},
	}
}

func getContentApplicationJSON(refSchema string) docs.ContentType {
	return docs.ContentType{
		Name:   "application/json",
		Schema: refSchema,
	}
}
