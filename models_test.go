package docs

import (
	"reflect"
	"testing"
)

func TestUnitNew(t *testing.T) {
	got := New()

	want := OAS{
		RegisteredRoutes: RegRoutes{},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, but want %+v", got, want)
	}
}
