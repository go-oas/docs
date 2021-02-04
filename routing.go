package docs

import (
	"reflect"
	"runtime"
	"strings"
)

type (
	// RouteFn represents a typeFunc which needs to be satisfied in order to use default routes attaching method.
	RouteFn func(index int, oas *OAS)

	// RegRoutes represent a map of RouteFn's.
	//
	// Note: Considering to un-export it.
	RegRoutes map[string]RouteFn
)

func (o *OAS) AttachRoutes(fns []RouteFn) {
	for _, fn := range fns {
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
