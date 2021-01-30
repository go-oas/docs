package docs

import (
	"reflect"
	"runtime"
	"strings"
)

type RegRoutes map[string]interface{}

func (o *OAS) AttachRoutes(fns []interface{}) {
	for _, fn := range fns {
		// TODO: Benchmark performance of this function with 1-1.5k routes.
		fnDeclaration := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		fields := strings.SplitAfter(fnDeclaration, ".")
		fnName := fields[len(fields)-1]

		o.RegisteredRoutes[fnName] = fn
	}
}

func (o *OAS) GetRegisteredRoutes() RegRoutes {
	return o.RegisteredRoutes
}

func (o *OAS) GetPathByIndex(index int) *Path {
	return &o.Paths[index]
}
