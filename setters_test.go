package docs

import (
	"reflect"
	"testing"
)

func TestUnitInfoSetLicense(t *testing.T) {
	t.Parallel()

	info := Info{
		Title: "tester",
	}
	lic := "MIT"
	urlTest := "https://github.com/kaynetik"
	info.SetLicense(lic, urlTest)

	if info.License.Name != lic {
		t.Errorf("license name not set correctly")
	}

	if info.License.URL != URL(urlTest) {
		t.Error("license URL not set correctly")
	}
}

func TestUnitInfoGetInfo(t *testing.T) {
	t.Parallel()

	oasPrep := OAS{Info: Info{
		Title: testingPostfix,
	}}

	got := oasPrep.GetInfo()

	if !reflect.DeepEqual(got, &oasPrep.Info) {
		t.Errorf("failed getting OAS.Info reference")
	}
}

func TestUnitSetOASVersion(t *testing.T) {
	t.Parallel()

	oasPrep := OAS{Info: Info{
		Title: testingPostfix,
	}}

	verToSet := "3.1.1"
	oasPrep.SetOASVersion(verToSet)

	if oasPrep.OASVersion != OASVersion(verToSet) {
		t.Error("failed setting OAS.OASVersion")
	}
}

func TestUnitSetTag(t *testing.T) {
	t.Parallel()

	tags := Tags{}
	tName := "pettag"
	tDesc := "Everything about your Pets"
	tED := ExternalDocs{
		Description: "Find out more about our store (Swagger UI Example)",
		URL:         "http://swagger.io",
	}
	tags.SetTag(tName, tDesc, tED)

	firstTag := tags[0]
	if firstTag.Name != tName ||
		firstTag.Description != tDesc ||
		firstTag.ExternalDocs != tED {
		t.Error("tag not set properly")
	}
}
