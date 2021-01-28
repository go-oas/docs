package docs

import (
	"errors"
	"reflect"
)

func (o *OAS) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(o.registeredRoutes[name])
	if len(params) != f.Type().NumIn() {
		err = errors.New("the number of params is not adapted")
		return //nolint: nakedret //implemetation speed. fixme: upgrade.
	}

	in := make([]reflect.Value, len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)

	return //nolint: nakedret //implemetation speed. fixme: upgrade.
}

const routePostfix = "Route"

func (o *OAS) initCallStackForRoutes() error {
	for oasPathIndex, oasPath := range o.Paths { //nolint:gocritic //fixme: troubleshoot if this will be an issue.
		_, err := o.Call(oasPath.handlerFuncName+routePostfix, oasPathIndex, o)
		if err != nil {
			return err
		}
	}

	return nil
}
