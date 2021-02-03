package docs

import "testing"

func TestUnitBuild(t *testing.T) {
	t.Parallel()

	oasPrep := prepForInitCallStack(t, false)

	info := &oasPrep.Info
	info.Title = "Info Testing"
	info.Description = "Not Mandatory"
	info.SetContact("aleksandar.nesovic@protonmail.com")
	info.SetLicense("MIT", "https://en.wikipedia.org/wiki/MIT_License")
	info.Version = "0.0.1"

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

	// this should be a ptr - currently not effective...
	path := oasPrep.Paths[0]

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
	path.HandlerFuncName = "testing"

	components := Components{}
	component := Component{
		Schemas: Schemas{Schema{
			Name: "schema_testing",
			XML:  XMLEntry{Name: "XML entry test"},
		}},
		SecuritySchemes: SecuritySchemes{SecurityScheme{
			Name: "ses_scheme_testing",
			In:   "not empty",
		}},
	}
	components = append(components, component)
	oasPrep.Components = components

	err := oasPrep.BuildDocs(ConfigBuilder{customPath: "./testing_out.yaml"})
	if err != nil {
		t.Errorf("unexpected error for OAS builder: %v", err)
	}
}

// TODO: MAJORITY OF TESTING SHOULD HAPPEN HERE
// BOTH WITH AND WITHOUT QUICK-CHECK

func TestUnitGetPathFromFirstElem(t *testing.T) {
	t.Parallel()

	cbs := make([]ConfigBuilder, 0)
	got := getPathFromFirstElement(cbs)

	if got != defaultDocsOutPath {
		t.Error("default docs path not set correctly")
	}
}
