package docs

import (
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"testing"
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
				registeredRoutes: tt.fields.registeredRoutes,
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
				registeredRoutes: tt.fields.registeredRoutes,
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
		registeredRoutes: rr,
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
		got := o.registeredRoutes[fnName]

		if reflect.ValueOf(got) != reflect.ValueOf(routeFn) {
			t.Errorf("invalid route fn attached: got %v, want %v", got, routeFn)
		}
	}
}
