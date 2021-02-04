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

// AttachRoutes if used for attaching pre-defined API documentation routes.
//
// fns param is a slice of functions that satisfy RouteFn signature.
func (o *OAS) AttachRoutes(fns []RouteFn) {
	for _, fn := range fns {
		fnDeclaration := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		fields := strings.SplitAfter(fnDeclaration, ".")
		fnName := fields[len(fields)-1]

		o.RegisteredRoutes[fnName] = fn
	}
}

// GetRegisteredRoutes returns a map of registered RouteFn functions - in layman terms "routes".
func (o *OAS) GetRegisteredRoutes() RegRoutes {
	return o.RegisteredRoutes
}

// GetPathByIndex returns ptr to Path structure, by its index in the parent struct of OAS.
func (o *OAS) GetPathByIndex(index int) *Path {
	return &o.Paths[index]
}
