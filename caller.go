package docs

import (
	"reflect"
)

// routePostfix will get exported by v1.3.
const routePostfix = "Route"

// Call is used init registered functions that already exist in the *OAS, and return results if there are any.
func (o *OAS) Call(name string, params ...interface{}) (result []reflect.Value) {
	f := reflect.ValueOf(o.RegisteredRoutes[name])

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)

	return result
}

func (o *OAS) initCallStackForRoutes() {
	for oasPathIndex := range o.Paths {
		o.Call(o.Paths[oasPathIndex].HandlerFuncName+routePostfix, oasPathIndex, o)
	}
}
