package docs

import "testing"

const (
	examplesDir = "./examples"
)

func TestUnitMapAnnotationsInPath(t *testing.T) {
	o := prepForInitCallStack(t, false)

	// TODO: Finish this test.
	_ = o.MapAnnotationsInPath(examplesDir)
}
