package docs

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

const testingPostfix = "Testing"

func TestQuickUnitCaller(t *testing.T) {
	t.Parallel()

	successParamNumber := func(i int, oas *OAS) {}
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

			args[0] = reflect.ValueOf(oas)
		},
	}

	callerCorrectParamNumber := func(oas OAS) bool {
		for i, oasPath := range oas.Paths {
			res := oas.Call(oasPath.Route+testingPostfix, i, &oas)

			if len(res) > 0 {
				t.Error("failed executing (OAS).Call() with")

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

	successParamNumber := func(i int, oas *OAS) {}
	routeName := "testRouteTesting"
	rr := make(RegRoutes)
	rr[routeName] = successParamNumber

	o := OAS{
		RegisteredRoutes: rr,
	}

	_ = o.Call(routeName, 0, &o)
}

func TestUnitInitCallStack(t *testing.T) {
	t.Parallel()

	o := prepForInitCallStack(t)

	o.initCallStackForRoutes()
}

func prepForInitCallStack(t *testing.T) OAS {
	t.Helper()

	routeName := "testRoute" + routePostfix
	rr := make(RegRoutes)

	rr[routeName] = getSuccessParamNumber

	path := Path{
		HandlerFuncName: "testRoute",
	}
	o := OAS{
		Paths:            Paths{path},
		RegisteredRoutes: rr,
	}

	return o
}

func getSuccessParamNumber(_ int, _ *OAS) {}

// func getFailureParamNumber(t *testing.T, _ int, _ *OAS) { t.Helper() }
