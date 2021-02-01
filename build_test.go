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

	err := oasPrep.BuildDocs(ConfigBuilder{customPath: "./testing_out.yaml"})
	if err != nil {
		t.Errorf("unexpected error for OAS builder: %v", err)
	}
}

// TODO: MAJORITY OF TESTING SHOULD HAPPEN HERE
// BOTH WITH AND WITHOUT QUICK-CHECK
