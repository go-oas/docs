package main

import (
	docs "github.com/go-oas/docs"
)

func main() {
	apiDoc := docs.New()

	apiDoc.OASVersion = "3.0.1"

	apiDoc.Info = docs.Info{
		Title:          "Build OAS3.0.1",
		Description:    "Description - Builder Testing for OAS3.0.1",
		TermsOfService: "ToS LoremIpsum",
		Contact: docs.Contact{
			Email: "nesovic@protonmail.com",
		},
		License: docs.License{
			Name: "MIT",
			URL:  "https://github.com/go-oas/docs/blob/main/LICENSE",
		},
		Version: "0.0.1",
	}

	apiDoc.Tags = docs.Tags{
		docs.Tag{
			Name:        "user",
			Description: "Operations about the User",
		}, docs.Tag{
			Name:        "pet",
			Description: "Pet Store example",
			ExternalDocs: docs.ExternalDocs{
				Description: "Find out more about our store",
				URL:         "http://swagger.io",
			},
		},
	}

	apiDoc.Servers = docs.Servers{
		docs.Server{
			URL: "https://petstore.swagger.io/v2",
		},
		docs.Server{
			URL: "http://petstore.swagger.io/v2",
		},
	}

	apiDoc.ExternalDocs = docs.ExternalDocs{
		Description: "External documentation",
		URL:         "https://kaynetik.com",
	}

	apiDoc.Components = docs.Components{
		docs.Component{
			Schemas: docs.Schemas{
				docs.Schema{
					Name: "User",
					Type: "object",
					Properties: docs.SchemaProperties{
						docs.SchemaProperty{
							Type:        "integer",
							Format:      "int64",
							Description: "UserID",
						},
						docs.SchemaProperty{
							Type:        "bool",
							Description: "isActive",
							Default:     true,
						},
					},
					XML: docs.XMLEntry{},
				},
			},
			SecuritySchemes: docs.SecuritySchemes{
				docs.SecurityScheme{
					Name: "users_auth",
					Type: "apiKey",
					In:   "header",
				},
			},
		},
	}

	apiDoc.AttachRoutes([]interface{}{
		handleCreateUserRoute,
		handleGetUserRoute,
	})

	err := apiDoc.MapAnnotationsInPath("./examples/users_example")
	if err != nil {
		panic(err)
	}

	err = apiDoc.BuildDocs("")
	if err != nil {
		panic(err)
	}

	// Serve Swagger //fixme: pass struct for config.
	docs.ServeSwaggerUI("/docs/api/", "3005")
}
