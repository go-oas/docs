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
func (oas *OAS) AttachRoutes(fns []RouteFn) {
	for _, fn := range fns {
		fnDeclaration := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		fields := strings.SplitAfter(fnDeclaration, ".")
		fnName := fields[len(fields)-1]

		oas.RegisteredRoutes[fnName] = fn
	}
}

// GetRegisteredRoutes returns a map of registered RouteFn functions - in layman terms "routes".
func (oas *OAS) GetRegisteredRoutes() RegRoutes {
	return oas.RegisteredRoutes
}

// GetPathByIndex returns ptr to Path structure, by its index in the parent struct of OAS.
func (oas *OAS) GetPathByIndex(index int) *Path {
	return &oas.Paths[index]
}

// AddRoute is used for add API documentation routes.
func (oas *OAS) AddRoute(path Path) {
	oas.Paths = append(oas.Paths, path)
}
