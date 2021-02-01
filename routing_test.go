package docs

import (
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
)

func fetchRegRoutes(t *testing.T, count int) RegRoutes {
	t.Helper()

	routes := make(RegRoutes)
	randStr := RandomString(t, count)

	for i := 0; i < count; i++ {
		routes[randStr] = initRoutes(t, count)
	}

	return routes
}

func initRoutes(t *testing.T, count int) []interface{} {
	t.Helper()

	var routes []interface{}

	for i := 0; i < count; i++ {
		routes = append(routes, func() {})
	}

	return routes
}

func TestUnitGetRegisteredRoutes(t *testing.T) {
	t.Parallel()

	type fields struct {
		registeredRoutes RegRoutes
	}

	v := rand.Intn(1000-1) + 1 //nolint:gosec //ignored in tests.
	regRoutes := fetchRegRoutes(t, v)

	tests := []struct {
		name   string
		fields fields
		want   RegRoutes
	}{
		{
			name: "success getting registered routes",
			fields: fields{
				registeredRoutes: regRoutes,
			},
			want: regRoutes,
		},
	}
	for _, tt := range tests { //nolint:paralleltest //ignore.
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o := &OAS{
				RegisteredRoutes: tt.fields.registeredRoutes,
			}
			got := o.GetRegisteredRoutes()
			want := tt.want

			if !reflect.DeepEqual(got, want) {
				t.Errorf("GetRegisteredRoutes() = %v, want %v", got, want)
			}
		})
	}
}

func TestUnitGetPathByIndex(t *testing.T) {
	t.Parallel()

	type fields struct {
		Paths            Paths
		registeredRoutes RegRoutes
	}

	paths := Paths{
		Path{
			Route:      "/test",
			HTTPMethod: http.MethodGet,
		},
	}

	v := rand.Intn(1000-1) + 1 //nolint:gosec //ignored in tests.
	regRoutes := fetchRegRoutes(t, v)

	tests := []struct {
		name   string
		fields fields
		want   *Path
	}{
		{
			name: "success get paths",
			fields: fields{
				Paths:            paths,
				registeredRoutes: regRoutes,
			},
			want: &paths[0],
		},
	}
	for _, tt := range tests { //nolint:paralleltest //Range statement for test TestUnitGetPathByIndex
		// does not reinitialise the variable tt -> TODO: Troubleshoot this further
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			o := &OAS{
				Paths:            tt.fields.Paths,
				RegisteredRoutes: tt.fields.registeredRoutes,
			}
			if got := o.GetPathByIndex(0); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPathByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitAttachRoutes(t *testing.T) {
	t.Parallel()

	rr := make(RegRoutes)

	o := OAS{
		RegisteredRoutes: rr,
	}

	routes := []interface{}{
		fetchRegRoutes,
		initRoutes,
	}

	o.AttachRoutes(routes)

	for _, routeFn := range routes {
		fnDeclaration := runtime.FuncForPC(reflect.ValueOf(routeFn).Pointer()).Name()
		fields := strings.SplitAfter(fnDeclaration, ".")
		fnName := fields[len(fields)-1]
		got := o.RegisteredRoutes[fnName]

		if reflect.ValueOf(got) != reflect.ValueOf(routeFn) {
			t.Errorf("invalid route fn attached: got %v, want %v", got, routeFn)
		}
	}
}

func TestQuickUnitGetRegisteredRoutes(t *testing.T) {
	t.Parallel()

	config := quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			oas := OAS{
				RegisteredRoutes: map[string]interface{}{},
			}
			args[0] = reflect.ValueOf(oas)
		},
	}

	gotRegRoutes := func(oas OAS) bool {
		got := oas.GetRegisteredRoutes()

		return reflect.DeepEqual(got, oas.RegisteredRoutes)
	}

	if err := quick.Check(gotRegRoutes, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func getPaths(t *testing.T, count int, rndStr string) Paths {
	t.Helper()

	pt := make(Paths, 0, count+2)

	for i := 0; i < count+1; i++ {
		pt = append(pt, Path{
			Route:           rndStr,
			HTTPMethod:      http.MethodGet,
			HandlerFuncName: rndStr,
		})
	}

	return pt
}

func TestQuickUnitGetPathByIndex(t *testing.T) {
	t.Parallel()

	config := quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			count := rand.Intn(550-1) + 1
			oas := OAS{
				Paths: getPaths(t, count, RandomString(t, count)),
			}

			args[0] = reflect.ValueOf(oas)
		},
	}

	gotRegRoutes := func(oas OAS) bool {
		pathsLen := len(oas.Paths)
		r := 2

		if pathsLen > 3 {
			r = int(uint(len(oas.Paths) - 2))
		}

		upRnd := int(uint(rand.Intn(r))) //nolint:gosec //week rnd generator - ignore in test.

		randIndex := uint(pathsLen - upRnd)

		got := oas.GetPathByIndex(int(randIndex - 1))

		return reflect.DeepEqual(got, &oas.Paths[randIndex-1])
	}

	if err := quick.Check(gotRegRoutes, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func RandomString(t *testing.T, n int) string {
	t.Helper()

	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))] //nolint:gosec //ignored in tests.
	}

	return string(s)
}

func TestQuickUnitAttachRoutes(t *testing.T) {
	t.Parallel()

	config := quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			rr := make(RegRoutes)

			oas := OAS{
				RegisteredRoutes: rr,
			}

			routes := []interface{}{
				fetchRegRoutes,
				initRoutes,
			}

			args[0] = reflect.ValueOf(oas)
			args[1] = reflect.ValueOf(routes)
		},
	}

	gotRegRoutes := func(oas OAS, routes []interface{}) bool {
		oas.AttachRoutes(routes)
		got := oas.GetRegisteredRoutes()

		return reflect.DeepEqual(got, oas.RegisteredRoutes)
	}

	if err := quick.Check(gotRegRoutes, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}
