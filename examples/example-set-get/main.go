package main

import "github.com/go-oas/docs"

func main() {
	apiDoc := docs.New()

	apiDoc.SetOASVersion("3.0.1")

	info := apiDoc.GetInfo()
	info.Title = "Build OAS3.0.1"
	info.Description = "Description - Builder Testing for OAS3.0.1"
	info.TermsOfService = "ToS LoremIpsum"
	info.SetContact("aleksandar.nesovic@protonmail.com")
	info.SetLicense("MIT", "https://github.com/go-oas/docs/blob/main/LICENSE")
	info.Version = "0.0.1"

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
		dummyRouteFn,
	})

	err := apiDoc.MapAnnotationsInPath("./examples/users_example")
	if err != nil {
		panic(err)
	}

	err = apiDoc.BuildDocs("")
	if err != nil {
		panic(err)
	}

	err = docs.ServeSwaggerUI(&docs.ConfigSwaggerUI{
		Route: "/docs/api/",
		Port:  "3005",
	})
	if err != nil {
		panic(err)
	}
}

func dummyRouteFn() {}
