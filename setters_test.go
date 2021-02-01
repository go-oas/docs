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

	var oasPrep = OAS{Info: Info{
		Title: testingPostfix,
	}}

	got := oasPrep.GetInfo()

	if !reflect.DeepEqual(got, &oasPrep.Info) {
		t.Errorf("failed getting OAS.Info reference")
	}
}

func TestUnitSetOASVersion(t *testing.T) {
	t.Parallel()

	var oasPrep = OAS{Info: Info{
		Title: testingPostfix,
	}}

	verToSet := "3.1.1"
	oasPrep.SetOASVersion(verToSet)

	if oasPrep.OASVersion != OASVersion(verToSet) {
		t.Error("failed setting OAS.OASVersion")
	}
}
