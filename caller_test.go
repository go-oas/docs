package docs

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

const testingPostfix = "Testing"

func TestQuickUnitCallerErr(t *testing.T) {
	t.Parallel()

	config := quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			rr := make(RegRoutes)
			count := rand.Intn(550-1) + 1
			paths := getPaths(t, count, RandomString(t, count))

			for _, p := range paths {
				rr[p.Route+testingPostfix] = getPaths
			}

			oas := OAS{
				Paths:            getPaths(t, count, RandomString(t, count)),
				RegisteredRoutes: rr,
			}

			var routes []interface{}
			for _, r := range rr {
				routes = append(routes, r)
			}

			args[0] = reflect.ValueOf(oas)
			args[1] = reflect.ValueOf(routes)
		},
	}

	callerWrongParamNumber := func(oas OAS, routes []interface{}) bool {
		oas.AttachRoutes(routes)

		for oasPathIndex := range oas.Paths {
			_, err := oas.Call("getPaths", oasPathIndex, oas)
			if err == nil {
				t.Errorf("failing (OAS).Call() with err : %s", err)

				return false
			}
		}

		return true
	}

	if err := quick.Check(callerWrongParamNumber, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func TestQuickUnitCaller(t *testing.T) {
	t.Parallel()

	successParamNumber := func(name string, oas *OAS) {}
	config := quick.Config{
		Values: func(args []reflect.Value, rand *rand.Rand) {
			rr := make(RegRoutes)
			count := rand.Intn(550-1) + 1
			paths := getPaths(t, count, RandomString(t, count))

			for i := 0; i < count; i++ {
				rr[paths[0].Route+testingPostfix] = successParamNumber
			}

			oas := OAS{
				Paths:            paths,
				RegisteredRoutes: rr,
			}

			var routes []interface{}
			for _, r := range rr {
				routes = append(routes, r)
			}

			args[0] = reflect.ValueOf(oas)
			// args[1] = reflect.ValueOf(routes)
		},
	}

	callerCorrectParamNumber := func(oas OAS) bool {
		for _, oasPath := range oas.Paths {
			_, err := oas.Call(oasPath.Route+testingPostfix, oasPath.Route, &oas)
			if err != nil {
				t.Errorf("failed executing (OAS).Call() with err : %s", err)

				return false
			}
		}

		return true
	}

	if err := quick.Check(callerCorrectParamNumber, &config); err != nil {
		t.Errorf("Check failed: %#v", err)
	}
}

func TestUnitCaller(t *testing.T) {
	t.Parallel()

	successParamNumber := func(name string, oas *OAS) {}
	routeName := "testRouteTesting"
	rr := make(RegRoutes)
	rr[routeName] = successParamNumber

	o := OAS{
		RegisteredRoutes: rr,
	}

	_, err := o.Call(routeName, routeName, &o)
	if err != nil {
		t.Errorf("failed executing (OAS).Call() with err : %s", err)
	}
}

func TestUnitInitCallStack(t *testing.T) {
	t.Parallel()

	o := prepForInitCallStack(t, false)

	err := o.initCallStackForRoutes()
	if err != nil {
		t.Errorf("failed executing (OAS).initCallStackForRoutes() with err : %s", err)
	}
}

func TestUnitInitCallStackErr(t *testing.T) {
	t.Parallel()

	o := prepForInitCallStack(t, true)

	err := o.initCallStackForRoutes()
	if err == nil {
		t.Errorf("failed executing (OAS).initCallStackForRoutes() with err : %s", err)
	}
}

func prepForInitCallStack(t *testing.T, triggerErr bool) OAS {
	t.Helper()

	routeName := "testRoute" + routePostfix
	rr := make(RegRoutes)

	if !triggerErr {
		rr[routeName] = getSuccessParamNumber
	} else {
		rr[routeName] = getFailureParamNumber
	}

	path := Path{
		HandlerFuncName: "testRoute",
	}
	o := OAS{
		Paths:            Paths{path},
		RegisteredRoutes: rr,
	}

	return o
}

func getSuccessParamNumber(_ int, _ *OAS)               {}
func getFailureParamNumber(_ *testing.T, _ int, _ *OAS) {}
