package docs

import (
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"testing/quick"
	"time"
)

func fetchRegRoutes(count int, t *testing.T) RegRoutes {
	t.Helper()

	var routes = make(RegRoutes)
	rand.Seed(time.Now().Unix())
	//Only lowercase and £
	charSet := []rune("ijklabcdedfpqrst£ghmno")
	var output strings.Builder
	length := 10
	for i := 0; i < length; i++ {
		random := rand.Intn(len(charSet))
		randomChar := charSet[random]
		output.WriteRune(randomChar)
	}

	for i := 0; i < count; i++ {
		routes[output.String()] = initRoutes(count, t)
	}

	return routes
}

func initRoutes(count int, t *testing.T) []interface{} {
	t.Helper()

	var routes []interface{}

	for i := 0; i < count; i++ {
		routes = append(routes, func() {})
	}

	return routes
}

func TestUnitGetRegisteredRoutes(t *testing.T) {
	type fields struct {
		registeredRoutes RegRoutes
	}

	v := rand.Intn(1000-1) + 1
	regRoutes := fetchRegRoutes(v, t)

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OAS{
				RegisteredRoutes: tt.fields.registeredRoutes,
			}
			if got := o.GetRegisteredRoutes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRegisteredRoutes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnitGetPathByIndex(t *testing.T) {
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

	v := rand.Intn(1000-1) + 1
	regRoutes := fetchRegRoutes(v, t)

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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

		if !reflect.DeepEqual(got, oas.RegisteredRoutes) {
			return false
		}

		return true
	}

	if err := quick.Check(gotRegRoutes, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func TestQuickUnitGetPathByIndex(t *testing.T) {
	paths := func(count int, rndStr string) Paths {
		pt := make(Paths, count+2)

		for i := 0; i < count+1; i++ {
			pt = append(pt, Path{
				Route:      rndStr,
				HTTPMethod: http.MethodGet,
			})
		}

		return pt
	}

	config := quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			count := rand.Intn(550-1) + 1
			oas := OAS{
				Paths: paths(count, RandomString(count)),
			}

			args[0] = reflect.ValueOf(oas)
		},
	}

	gotRegRoutes := func(oas OAS) bool {
		randIndex := uint(len(oas.Paths) - rand.Intn(len(oas.Paths)-2))

		got := oas.GetPathByIndex(int(randIndex))

		if !reflect.DeepEqual(got, &oas.Paths[randIndex]) {
			return false
		}

		return true
	}

	if err := quick.Check(gotRegRoutes, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}

func TestQuickUnitAttachRoutes(t *testing.T) {
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

		if !reflect.DeepEqual(got, oas.RegisteredRoutes) {
			return false
		}

		return true
	}
	if err := quick.Check(gotRegRoutes, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}

}
