package vm

import (
	"fmt"
	"github.com/warriorguo/gopy/pkg/runtime"
	"github.com/warriorguo/gopy/pkg/object"
)

func toGoInt(obj object.Object) (int, error) {
	switch o := obj.(type) {
	case *runtime.PyInt:
		return o.Value, nil
	case *runtime.PyFloat:
		return int(o.Value), nil
	case *runtime.PyBool:
		if o.Value {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to int", obj.Type())
	}
}

func toGoFloat(obj object.Object) (float64, error) {
	switch o := obj.(type) {
	case *runtime.PyInt:
		return float64(o.Value), nil
	case *runtime.PyFloat:
		return o.Value, nil
	case *runtime.PyBool:
		if o.Value {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to float", obj.Type())
	}
}

func toGoString(obj object.Object) string {
	switch o := obj.(type) {
	case *runtime.PyString:
		return o.Value
	default:
		return obj.String()
	}
}