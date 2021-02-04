package docs

import (
	"errors"
	"path/filepath"
	"testing"
)

const (
	examplesDir = "./examples"
)

func TestUnitMapAnnotationsInPath(t *testing.T) {
	t.Parallel()

	o := prepForInitCallStack(t, false)

	_ = o.MapAnnotationsInPath(examplesDir)
}

func TestUnitMapAnnotationsInPathErr(t *testing.T) {
	t.Parallel()

	o := (*OAS)(nil)

	err := o.MapAnnotationsInPath(examplesDir)
	if err == nil {
		t.Error("expected an error, got none")
	}
}

func TestUnitScanForChangesInPathErrFnWD(t *testing.T) {
	t.Parallel()

	wd := func() (dir string, err error) {
		return "", errors.New("test")
	}

	_, err := scanForChangesInPath("", wd, walkFilepath)
	if err == nil {
		t.Error("expected an error, got none")
	}
}

func TestUnitScanForChangesInPathErrWalk(t *testing.T) {
	t.Parallel()

	wd := func() (dir string, err error) {
		return "~@!", nil
	}
	wdErr := func() (dir string, err error) {
		return "~@!", errors.New("walkDir error")
	}

	pathWalkerErr := func(path string, walker walkerFn) ([]string, error) {
		return []string{}, errors.New("triggered")
	}

	_, err := scanForChangesInPath("", wd, pathWalkerErr)
	if err == nil {
		t.Error("expected an error, got none")
	}

	o := prepForInitCallStack(t, false)
	errConfig := configAnnotation{
		getWD: wdErr,
	}

	err = o.MapAnnotationsInPath(".", errConfig)
	if err == nil {
		t.Error("expected an error, got none")
	}

	walkFnErr := func(root string, walkFn filepath.WalkFunc) error {
		return errors.New("triggered")
	}

	_, err = walkFilepath("", walkFnErr)
	if err == nil {
		t.Error("expected an error, got none")
	}
}
