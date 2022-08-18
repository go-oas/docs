package main

import "github.com/go-oas/docs"

func apiSetInfo(apiDoc *docs.OAS) {
	apiDoc.SetOASVersion("3.0.1")
	apiInfo := apiDoc.GetInfo()
	apiInfo.Title = "Build OAS3.0.1"
	apiInfo.Description = "Builder Testing for OAS3.0.1"
	apiInfo.TermsOfService = "https://smartbear.com/terms-of-use/"
	apiInfo.SetContact("padiazg@gmail.com") // mixed usage of setters ->
	apiInfo.License = docs.License{         // and direct struct usage.
		Name: "MIT",
		URL:  "https://github.com/go-oas/docs/blob/main/LICENSE",
	}
	apiInfo.Version = "1.0.1"
}

func apiSetTags(apiDoc *docs.OAS) {
	// With Tags example you can see usage of direct struct modifications, setter and appender as well.
	apiDoc.Tags = docs.Tags{
		docs.Tag{
			Name:        "user",
			Description: "Operations about the User",
			ExternalDocs: docs.ExternalDocs{
				Description: "User from the Petstore example",
				URL:         "http://swagger.io",
			},
		},
	}
	apiDoc.Tags.SetTag(
		"pet",
		"Everything about your Pets",
		docs.ExternalDocs{
			Description: "Find out more about our store (Swagger UI Example)",
			URL:         "http://swagger.io",
		},
	)

	newTag := &docs.Tag{
		Name:        "petko",
		Description: "Everything about your Petko",
		ExternalDocs: docs.ExternalDocs{
			Description: "Find out more about our store (Swagger UI Example)",
			URL:         "http://swagger.io",
		},
	}
	apiDoc.Tags.AppendTag(newTag)
}

func apiSetServers(apiDoc *docs.OAS) {
	apiDoc.Servers = docs.Servers{
		docs.Server{
			URL: "https://petstore.swagger.io/v2",
		},
		docs.Server{
			URL: "http://petstore.swagger.io/v2",
		},
	}
}

func apiSetExternalDocs(apiDoc *docs.OAS) {
	apiDoc.ExternalDocs = docs.ExternalDocs{
		Description: "External documentation",
		URL:         "https://kaynetik.com",
	}
}

func apiSetComponents(apiDoc *docs.OAS) {
	apiDoc.Components = docs.Components{
		docs.Component{
			Schemas: docs.Schemas{
				docs.Schema{
					Name: "User",
					Type: "object",
					Properties: docs.SchemaProperties{
						docs.SchemaProperty{
							Name:        "id",
							Type:        "integer",
							Format:      "int64",
							Description: "UserID",
						},
						docs.SchemaProperty{
							Name: "username",
							Type: "string",
						},
						docs.SchemaProperty{
							Name: "email",
							Type: "string",
						},
						docs.SchemaProperty{
							Name:        "userStatus",
							Type:        "integer",
							Description: "User Status",
							Format:      "int32",
						},
						docs.SchemaProperty{
							Name: "phForEnums",
							Type: "enum",
							Enum: []string{"placed", "approved"},
						},
					},
					XML: docs.XMLEntry{Name: "User"},
				},
				docs.Schema{
					Name: "Tag",
					Type: "object",
					Properties: docs.SchemaProperties{
						docs.SchemaProperty{
							Name:   "id",
							Type:   "integer",
							Format: "int64",
						},
						docs.SchemaProperty{
							Name: "name",
							Type: "string",
						},
					},
					XML: docs.XMLEntry{Name: "Tag"},
				},
				docs.Schema{
					Name: "ApiResponse",
					Type: "object",
					Properties: docs.SchemaProperties{
						docs.SchemaProperty{
							Name:   "code",
							Type:   "integer",
							Format: "int32",
						},
						docs.SchemaProperty{
							Name: "type",
							Type: "string",
						},
						docs.SchemaProperty{
							Name: "message",
							Type: "string",
						},
					},
					XML: docs.XMLEntry{Name: "ApiResponse"},
				},
			},
			SecuritySchemes: docs.SecuritySchemes{
				docs.SecurityScheme{
					Name: "api_key",
					Type: "apiKey",
					In:   "header",
				},
				docs.SecurityScheme{
					Name: "petstore_auth",
					Type: "oauth2",
					Flows: docs.SecurityFlows{
						docs.SecurityFlow{
							Type:    "implicit",
							AuthURL: "http://petstore.swagger.io/oauth/dialog",
							Scopes: docs.SecurityScopes{
								docs.SecurityScope{
									Name:        "write:users",
									Description: "Modify users",
								},
								docs.SecurityScope{
									Name:        "read:users",
									Description: "Read users",
								},
							},
						},
					},
				},
			},
		},
	}
}
