package docs

import (
	"fmt"
	"reflect"
)

func (o *OAS) Call(name string, params ...interface{}) (result []reflect.Value, err error) {
	f := reflect.ValueOf(o.RegisteredRoutes[name])
	paramNum := len(params)
	fnParamNum := f.Type().NumIn()

	if paramNum != fnParamNum {
		return result, fmt.Errorf(
			"param number differs -> expected %d, got %d",
			paramNum, fnParamNum,
		)
	}

	in := make([]reflect.Value, paramNum)
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}

	result = f.Call(in)

	return result, nil
}

// should this be flexible for change?
const routePostfix = "Route"

func (o *OAS) initCallStackForRoutes() error {
	for oasPathIndex, oasPath := range o.Paths { //nolint:gocritic //fixme: troubleshoot if this will be an issue.
		_, err := o.Call(oasPath.HandlerFuncName+routePostfix, oasPathIndex, o)
		if err != nil {
			return fmt.Errorf(" :%w", err)
		}
	}

	return nil
}
