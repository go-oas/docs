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

	ir := initRoutes(t, count)

	for i := 0; i < count; i++ {
		routes[randStr] = ir[i]
	}

	return routes
}

func initRoutes(t *testing.T, count int) []RouteFn {
	t.Helper()

	var routes []RouteFn

	for i := 0; i < count; i++ {
		tmpFn := func(i int, o *OAS) {}
		routes = append(routes, tmpFn)
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

	type st struct {
		name   string
		fields fields
		want   RegRoutes
	}

	tests := []st{
		{
			name: "success getting registered routes",
			fields: fields{
				registeredRoutes: regRoutes,
			},
			want: regRoutes,
		},
	}
	for _, tt := range tests {
		stin := tt

		t.Run(stin.name, func(t *testing.T) {
			t.Parallel()

			o := &OAS{
				RegisteredRoutes: stin.fields.registeredRoutes,
			}
			got := o.GetRegisteredRoutes()
			want := stin.want

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

	type st struct {
		name   string
		fields fields
		want   *Path
	}

	paths := Paths{
		Path{
			Route:      "/test",
			HTTPMethod: http.MethodGet,
		},
	}

	v := rand.Intn(1000-1) + 1 //nolint:gosec //ignored in tests.
	regRoutes := fetchRegRoutes(t, v)

	tests := []st{
		{
			name: "success get paths",
			fields: fields{
				Paths:            paths,
				registeredRoutes: regRoutes,
			},
			want: &paths[0],
		},
	}
	for _, tt := range tests {
		stin := tt
		t.Run(stin.name, func(t *testing.T) {
			t.Parallel()

			o := &OAS{
				Paths:            stin.fields.Paths,
				RegisteredRoutes: stin.fields.registeredRoutes,
			}
			if got := o.GetPathByIndex(0); !reflect.DeepEqual(got, stin.want) {
				t.Errorf("GetPathByIndex() = %v, want %v", got, stin.want)
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

	routes := initRoutes(t, 5)

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
				RegisteredRoutes: map[string]RouteFn{},
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

			routes := initRoutes(t, 5)

			args[0] = reflect.ValueOf(oas)
			args[1] = reflect.ValueOf(routes)
		},
	}

	gotRegRoutes := func(oas OAS, routes []RouteFn) bool {
		oas.AttachRoutes(routes)
		got := oas.GetRegisteredRoutes()

		return reflect.DeepEqual(got, oas.RegisteredRoutes)
	}

	if err := quick.Check(gotRegRoutes, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func TestOAS_AddRoute(t *testing.T) {
	t.Parallel()

	var (
		respose200         = Response{Code: 200, Description: "Ok"}
		respose404         = Response{Code: 404, Description: "Not Found"}
		contentTypeUser    = ContentType{Name: "application/json", Schema: "#/components/schemas/User"}
		requestBodyGetUser = RequestBody{
			Description: "Get a User",
			Content:     ContentTypes{contentTypeUser},
			Required:    true,
		}
		requestBodyCreateUser = RequestBody{
			Description: "Create a new User",
			Content:     ContentTypes{contentTypeUser},
			Required:    true,
		}
		pathGetUser = Path{
			Route:       "/users",
			HTTPMethod:  "GET",
			OperationID: "getUser",
			Summary:     "Get a User",
			Responses:   Responses{respose200},
			RequestBody: requestBodyGetUser,
		}
		pathCreateUser = Path{
			Route:       "/users",
			HTTPMethod:  "POST",
			OperationID: "createUser",
			Summary:     "Create a new User",
			Responses:   Responses{respose200, respose404},
			RequestBody: requestBodyCreateUser,
		}
	)

	tests := []struct {
		name      string
		oas       *OAS
		path      *Path
		wantPaths Paths
	}{
		{
			name:      "success-no-existing-paths",
			oas:       &OAS{},
			path:      &pathGetUser,
			wantPaths: Paths{pathGetUser},
		},
		{
			name:      "success-existing-paths",
			oas:       &OAS{Paths: Paths{pathGetUser}},
			path:      &pathCreateUser,
			wantPaths: Paths{pathGetUser, pathCreateUser},
		},
	}

	for _, tt := range tests {
		trn := tt

		t.Run(trn.name, func(t *testing.T) {
			t.Parallel()

			trn.oas.AddRoute(trn.path)
			if !reflect.DeepEqual(trn.wantPaths, trn.oas.Paths) {
				t.Errorf("OAS.AddRoute() = [%v], want {%v}", trn.oas.Paths, trn.wantPaths)
			}
		})
	}
}
