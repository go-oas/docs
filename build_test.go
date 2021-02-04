package docs

import "testing"

func TestUnitBuild(t *testing.T) {
	t.Parallel()

	oasPrep := prepForInitCallStack(t)

	setInfoForTest(t, &oasPrep.Info)
	setPathForTest(t, &oasPrep.Paths[0])

	components := Components{}
	component := Component{
		Schemas: Schemas{Schema{
			Name: "schema_testing",
			Properties: SchemaProperties{
				SchemaProperty{
					Name:        "EnumProp",
					Type:        "enum",
					Description: "short desc",
					Enum:        []string{"enum", "test", "strSlc"},
				},
				SchemaProperty{
					Name:        "intProp",
					Type:        "integer",
					Format:      "int64",
					Description: "short desc",
					Default:     1337,
				},
			},
			XML: XMLEntry{Name: "XML entry test"},
		}},
		SecuritySchemes: SecuritySchemes{SecurityScheme{
			Name: "ses_scheme_testing",
			In:   "not empty",
			Flows: SecurityFlows{SecurityFlow{
				Type:    "implicit",
				AuthURL: "http://petstore.swagger.io/oauth/dialog",
				Scopes: SecurityScopes{
					SecurityScope{
						Name:        "write:pets",
						Description: "Write to Pets",
					},
					SecurityScope{
						Name:        "read:pets",
						Description: "Read Pets",
					},
				},
			}},
		}},
	}
	components = append(components, component)
	oasPrep.Components = components

	err := oasPrep.BuildDocs(ConfigBuilder{CustomPath: "./testing_out.yaml"})
	if err != nil {
		t.Errorf("unexpected error for OAS builder: %v", err)
	}
}

func setInfoForTest(t *testing.T, info *Info) {
	t.Helper()

	info.Title = "Info Testing"
	info.Description = "Not Mandatory"
	info.SetContact("aleksandar.nesovic@protonmail.com")
	info.SetLicense("MIT", "https://en.wikipedia.org/wiki/MIT_License")
	info.Version = "0.0.1"
}

func setPathForTest(t *testing.T, path *Path) {
	t.Helper()

	cts := ContentTypes{ContentType{
		Name:   "testingType",
		Schema: "schema_testing",
	}}
	response := Response{
		Code:        200,
		Description: "OK",
		Content:     cts,
	}
	responses := Responses{}
	responses = append(responses, response)

	path.HTTPMethod = "GET"
	path.Tags = []string{}
	path.Summary = "TestingSummary"
	path.OperationID = "TestingOperationID"
	path.RequestBody = RequestBody{
		Description: "testReq",
		Content:     cts,
		Required:    true,
	}
	path.Responses = responses
	path.Security = SecurityEntities{Security{AuthName: "sec"}}
}

func TestUnitGetPathFromFirstElem(t *testing.T) {
	t.Parallel()

	cbs := make([]ConfigBuilder, 0)
	got := getPathFromFirstElement(cbs)

	if got != defaultDocsOutPath {
		t.Error("default docs path not set correctly")
	}
}

// QUICK CHECK TESTS ARE COMING WITH NEXT RELEASE.
