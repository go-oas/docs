package docs

import (
	"reflect"
	"testing"
)

func TestUnitNew(t *testing.T) {
	t.Parallel()

	got := New()

	want := OAS{
		RegisteredRoutes: RegRoutes{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, but want %+v", got, want)
	}
}

func TestUnitEDIsEmpty(t *testing.T) {
	t.Parallel()

	edNil := (*ExternalDocs)(nil)
	if !edNil.isEmpty() {
		t.Errorf("expected an empty external docs struct")
	}

	edEmpty := ExternalDocs{}
	if !edEmpty.isEmpty() {
		t.Error("expected an empty external docs struct")
	}

	edEmptyURL := ExternalDocs{
		Description: "description",
	}

	if !edEmptyURL.isEmpty() && !isStrEmpty(string(edEmptyURL.URL)) {
		t.Errorf("expected an empty URL")
	}

	edEmptyDesc := ExternalDocs{
		URL: "description",
	}

	if !edEmptyDesc.isEmpty() && !isStrEmpty(edEmptyDesc.Description) {
		t.Errorf("expected an empty Description")
	}

	ed := ExternalDocs{
		Description: "to many to describe",
		URL:         "Gagarin URI",
	}

	if ed.isEmpty() &&
		!isStrEmpty(ed.Description) &&
		!isStrEmpty(string(ed.URL)) {
		t.Errorf("expected complete ExternalDocs struct")
	}
}
