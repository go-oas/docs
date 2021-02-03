package docs

import (
	"errors"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"syscall"
	"testing"
	"testing/quick"
	"time"
)

func TestUnitQuickInitializeStandardFS(t *testing.T) {
	t.Parallel()

	initStandardFS := func() bool {
		fsInit := http.Dir(defaultDirectory)
		want := fileSystem{
			fileSysInit: fsInit,
			fsOpenFn:    newFSOpen(fsInit),
			getStatFn:   newGetStatFn(),
			getIsDir:    newGetIsDirFn(),
		}

		got := *initializeStandardFS()

		if got.fileSysInit != want.fileSysInit {
			return false
		}

		if reflect.TypeOf(got.fsOpenFn) != reflect.TypeOf(want.fsOpenFn) {
			return false
		}

		if reflect.TypeOf(got.getStatFn) != reflect.TypeOf(want.getStatFn) {
			return false
		}

		if reflect.TypeOf(got.getIsDir) != reflect.TypeOf(want.getIsDir) {
			return false
		}

		return true
	}
	if err := quick.Check(initStandardFS, nil); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func TestUnitQuickInitializeStandardFSErr(t *testing.T) {
	t.Parallel()

	config := quick.Config{
		Values: func(values []reflect.Value, rand *rand.Rand) {
			fis := http.Dir(defaultDirectory)
			testFsOpenFn := newFSOpen(fis)
			testGetStatFn := newGetStatFn()
			testGetIsDirFn := newGetIsDirFn()

			values[0] = reflect.ValueOf(testFsOpenFn)
			values[1] = reflect.ValueOf(testGetStatFn)
			values[2] = reflect.ValueOf(testGetIsDirFn)
		},
	}

	initStandardFS := func(fsOpenFn fsOpenFn, getStatFn getStatFn, getIsDirFn getIsDirFn) bool {
		fsInit := http.Dir(defaultDirectory)
		want := fileSystem{
			fileSysInit: fsInit,
			fsOpenFn:    fsOpenFn,
			getStatFn:   getStatFn,
			getIsDir:    getIsDirFn,
		}

		got := *initializeStandardFS()

		if got.fileSysInit != want.fileSysInit {
			return false
		}

		if reflect.TypeOf(got.fsOpenFn) != reflect.TypeOf(want.fsOpenFn) {
			return false
		}

		if reflect.TypeOf(got.getStatFn) != reflect.TypeOf(want.getStatFn) {
			return false
		}

		if reflect.TypeOf(got.getIsDir) != reflect.TypeOf(want.getIsDir) {
			return false
		}

		return true
	}
	if err := quick.Check(initStandardFS, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func TestUnitNewFSOpen(t *testing.T) {
	t.Parallel()

	fsInit := http.Dir(defaultDirectory)

	got := newFSOpen(fsInit)
	want := func(name string) (http.File, error) {
		return fsInit.Open(name)
	}

	gotRes, _ := got("/")
	wantRes, _ := want("/")

	gstat, _ := gotRes.Stat()
	wstat, _ := wantRes.Stat()

	if !reflect.DeepEqual(gstat, wstat) {
		t.Error()
	}
}

func TestUnitIsFSNil(t *testing.T) {
	t.Parallel()

	nilFS := (*fileSystem)(nil)
	if !nilFS.isNil() {
		t.Error()
	}

	fsInit := fileSystem{}
	if !fsInit.isNil() {
		t.Error()
	}

	fsStandard := *initializeStandardFS()
	if fsStandard.isNil() {
		t.Error()
	}
}

func TestUnitFSOpen(t *testing.T) {
	t.Parallel()

	fsStd := *initializeStandardFS()
	if fsStd.isNil() {
		t.Error()
	}

	file, err := fsStd.Open("")
	if err != nil {
		t.Errorf("got an unexpected error: %v", err)
	}

	if file == nil {
		t.Error("file is expected not to be nil")
	}

	fstat, _ := file.Stat()
	if fstat.Name() != "dist" {
		t.Errorf("got an unexpected file name: %v", fstat.Name())
	}
}

func TestUnitFSOpenOpenErr(t *testing.T) {
	t.Parallel()

	fsInit := http.Dir(defaultDirectory)

	fsStd := &fileSystem{
		fileSysInit: fsInit,
		fsOpenFn:    errFSOpen(t, fsInit),
		getStatFn:   newGetStatFn(),
		getIsDir:    newGetIsDirFn(),
	}

	if fsStd.isNil() {
		t.Error()
	}

	_, err := fsStd.Open("")
	if err == nil {
		t.Error("expected an error, got none")
	}
}

func TestUnitFSOpenGetStatErr(t *testing.T) {
	t.Parallel()

	fsInit := http.Dir(defaultDirectory)

	fsStd := &fileSystem{
		fileSysInit: fsInit,
		fsOpenFn:    newFSOpen(fsInit),
		getStatFn:   errGetStatFn(t),
		getIsDir:    newGetIsDirFn(),
	}

	if fsStd.isNil() {
		t.Error()
	}

	_, err := fsStd.Open("")
	if err == nil {
		t.Error("expected an error, got none")
	}
}

func TestUnitFSOpenDirInnerOpenErr(t *testing.T) {
	t.Parallel()

	fsInitCorrect := http.Dir(defaultDirectory)
	fsErr := http.Dir("!!\\@!")

	fsStd := &fileSystem{
		fileSysInit: fsErr,
		fsOpenFn:    newFSOpen(fsInitCorrect),
		getStatFn:   newGetStatFn(),
		getIsDir:    statTrueGetIsDirFn(),
	}

	if fsStd.isNil() {
		t.Error()
	}

	_, err := fsStd.Open("")
	if err == nil {
		t.Error("expected an error, got none")
	}
}

func errFSOpen(t *testing.T, fis http.FileSystem) fsOpenFn {
	t.Helper()

	return func(name string) (http.File, error) {
		file, _ := fis.Open(name)
		return file, errors.New("triggerErr")
	}
}

func errGetStatFn(t *testing.T) getStatFn {
	t.Helper()

	return func(file http.File) fileStatFn {
		return func() (os.FileInfo, error) {
			fi, _ := file.Stat()

			return fi, errors.New("triggerErr")
		}
	}
}

func statTrueGetIsDirFn() getIsDirFn {
	return func(file os.FileInfo) fsIsDirFn {
		return func() bool {
			return true
		}
	}
}

func TestUnitSwaggerUIShutDown(t *testing.T) {
	t.Parallel()

	conf := (*ConfigSwaggerUI)(nil)

	err := ServeSwaggerUI(conf)
	if err == nil {
		t.Error("expected an error, got none")
	}

	osSignal := make(chan os.Signal)

	emptyRoute := &ConfigSwaggerUI{
		Route:   "",
		stopper: osSignal,
	}

	go func() {
		time.Sleep(20 * time.Millisecond)
		osSignal <- syscall.SIGINT
	}()

	_ = ServeSwaggerUI(emptyRoute)

	if emptyRoute.Route != defaultRoute {
		t.Errorf("route wasn't altered from its zero state to default route")
	}
}

func TestUnitSigCont(t *testing.T) {
	t.Parallel()

	confSwg := &ConfigSwaggerUI{
		stopper: nil,
	}

	confSwg.sigCont()

	if confSwg.stopper == nil {
		t.Error("stopper chan is nil, should be os.Signal")
	}
}
