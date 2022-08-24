package docs

import (
	"reflect"
)

// routePostfix will get exported by v1.3.
const routePostfix = "Route"

// Call is used init registered functions that already exist in the *OAS, and return results if there are any.
func (oas *OAS) Call(name string, params ...interface{}) (result []reflect.Value) {
	f := reflect.ValueOf(oas.RegisteredRoutes[name])

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)

	return result
}

func (oas *OAS) initCallStackForRoutes() {
	for oasPathIndex := range oas.Paths {
		oas.Call(oas.Paths[oasPathIndex].HandlerFuncName+routePostfix, oasPathIndex, oas)
	}
}
