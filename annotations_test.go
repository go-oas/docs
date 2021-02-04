package docs

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"testing/quick"
)

const (
	examplesDir = "./examples"
)

func TestUnitMapAnnotationsInPath(t *testing.T) {
	t.Parallel()

	o := prepForInitCallStack(t)

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

	o := prepForInitCallStack(t)
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

type triggerNil struct {
	trigger bool
}

func TestQuickUnitGetWD(t *testing.T) {
	t.Parallel()

	config := quick.Config{
		Values: func(values []reflect.Value, rand *rand.Rand) {
			ca := configAnnotation{}
			rndNm := rand.Int()
			tn := triggerNil{trigger: false}
			if rndNm%2 == 0 {
				tn.trigger = true
				ca.getWD = func() (dir string, err error) {
					return "", nil
				}
			} else {
				ca.getWD = os.Getwd
			}

			values[0] = reflect.ValueOf(ca)
			values[1] = reflect.ValueOf(tn)
		},
	}

	gwdFetcher := func(ca configAnnotation, tn triggerNil) bool {
		got := ca.getCurrentDirFetcher()

		return reflect.TypeOf(got) == reflect.TypeOf(ca.getWD)
	}

	if err := quick.Check(gwdFetcher, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func TestUnitGWD(t *testing.T) {
	t.Parallel()

	ca := configAnnotation{
		getWD: func() (dir string, err error) {
			return "", nil
		},
	}

	got := ca.getCurrentDirFetcher()
	if reflect.TypeOf(got) != reflect.TypeOf(ca.getWD) {
		t.Error("functions differ")
	}

	caNil := configAnnotation{}.getCurrentDirFetcher()

	if gFnName(t, caNil) != gFnName(t, os.Getwd) {
		t.Error("functions differ")
	}
}

func gFnName(t *testing.T, i interface{}) string {
	t.Helper()

	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
