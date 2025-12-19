package runtime

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/warriorguo/gopy/pkg/object"
)

type PyInt struct {
	Value int
}

func (p *PyInt) String() string { return strconv.Itoa(p.Value) }
func (p *PyInt) Type() string   { return "int" }
func (p *PyInt) IsTruthy() bool { return p.Value != 0 }
func (p *PyInt) Equal(other object.Object) bool {
	if o, ok := other.(*PyInt); ok {
		return p.Value == o.Value
	}
	return false
}

type PyFloat struct {
	Value float64
}

func (p *PyFloat) String() string { return fmt.Sprintf("%g", p.Value) }
func (p *PyFloat) Type() string   { return "float" }
func (p *PyFloat) IsTruthy() bool { return p.Value != 0.0 }
func (p *PyFloat) Equal(other object.Object) bool {
	if o, ok := other.(*PyFloat); ok {
		return p.Value == o.Value
	}
	return false
}

type PyString struct {
	Value string
}

func (p *PyString) String() string { return fmt.Sprintf("%s", p.Value) }
func (p *PyString) Type() string   { return "str" }
func (p *PyString) IsTruthy() bool { return len(p.Value) > 0 }
func (p *PyString) Equal(other object.Object) bool {
	if o, ok := other.(*PyString); ok {
		return p.Value == o.Value
	}
	return false
}

type PyBool struct {
	Value bool
}

func (p *PyBool) String() string {
	if p.Value {
		return "True"
	}
	return "False"
}
func (p *PyBool) Type() string   { return "bool" }
func (p *PyBool) IsTruthy() bool { return p.Value }
func (p *PyBool) Equal(other object.Object) bool {
	if o, ok := other.(*PyBool); ok {
		return p.Value == o.Value
	}
	return false
}

type PyNone struct{}

func (p *PyNone) String() string { return "None" }
func (p *PyNone) Type() string   { return "NoneType" }
func (p *PyNone) IsTruthy() bool { return false }
func (p *PyNone) Equal(other object.Object) bool {
	_, ok := other.(*PyNone)
	return ok
}

type PyList struct {
	Elements []object.Object
}

func (p *PyList) String() string {
	var elements []string
	for _, elem := range p.Elements {
		elements = append(elements, elem.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(elements, ", "))
}
func (p *PyList) Type() string   { return "list" }
func (p *PyList) IsTruthy() bool { return len(p.Elements) > 0 }
func (p *PyList) Equal(other object.Object) bool {
	if o, ok := other.(*PyList); ok {
		if len(p.Elements) != len(o.Elements) {
			return false
		}
		for i, elem := range p.Elements {
			if !elem.Equal(o.Elements[i]) {
				return false
			}
		}
		return true
	}
	return false
}

type PyDict struct {
	Pairs map[string]object.Object
	Keys  []string
}

func NewPyDict() *PyDict {
	return &PyDict{
		Pairs: make(map[string]object.Object),
		Keys:  []string{},
	}
}

func (p *PyDict) Set(key, value object.Object) {
	keyStr := key.String()
	if _, exists := p.Pairs[keyStr]; !exists {
		p.Keys = append(p.Keys, keyStr)
	}
	p.Pairs[keyStr] = value
}

func (p *PyDict) Get(key object.Object) (object.Object, bool) {
	value, exists := p.Pairs[key.String()]
	return value, exists
}

func (p *PyDict) String() string {
	var pairs []string
	for _, key := range p.Keys {
		value := p.Pairs[key]
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, value.String()))
	}
	return fmt.Sprintf("{%s}", strings.Join(pairs, ", "))
}
func (p *PyDict) Type() string   { return "dict" }
func (p *PyDict) IsTruthy() bool { return len(p.Pairs) > 0 }
func (p *PyDict) Equal(other object.Object) bool {
	if o, ok := other.(*PyDict); ok {
		if len(p.Pairs) != len(o.Pairs) {
			return false
		}
		for key, value := range p.Pairs {
			if otherValue, exists := o.Pairs[key]; !exists || !value.Equal(otherValue) {
				return false
			}
		}
		return true
	}
	return false
}

type CodeObject interface {
	String() string
	Disassemble() string
}

type PyFunction struct {
	Code     CodeObject
	Globals  map[string]object.Object
	Defaults []object.Object
	Name     string
}

func (p *PyFunction) String() string {
	return fmt.Sprintf("<function %s>", p.Name)
}
func (p *PyFunction) Type() string   { return "function" }
func (p *PyFunction) IsTruthy() bool { return true }
func (p *PyFunction) Equal(other object.Object) bool {
	if o, ok := other.(*PyFunction); ok {
		return p.Code == o.Code
	}
	return false
}

type PyBuiltin struct {
	Name string
	Func func(args []object.Object) (object.Object, error)
}

func (p *PyBuiltin) String() string {
	return fmt.Sprintf("<built-in function %s>", p.Name)
}
func (p *PyBuiltin) Type() string   { return "builtin_function_or_method" }
func (p *PyBuiltin) IsTruthy() bool { return true }
func (p *PyBuiltin) Equal(other object.Object) bool {
	if o, ok := other.(*PyBuiltin); ok {
		return p.Name == o.Name
	}
	return false
}

func ToGoInt(obj object.Object) (int, error) {
	switch o := obj.(type) {
	case *PyInt:
		return o.Value, nil
	case *PyFloat:
		return int(o.Value), nil
	case *PyBool:
		if o.Value {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to int", obj.Type())
	}
}

func ToGoFloat(obj object.Object) (float64, error) {
	switch o := obj.(type) {
	case *PyInt:
		return float64(o.Value), nil
	case *PyFloat:
		return o.Value, nil
	case *PyBool:
		if o.Value {
			return 1.0, nil
		}
		return 0.0, nil
	default:
		return 0, fmt.Errorf("cannot convert %s to float", obj.Type())
	}
}

func ToGoString(obj object.Object) string {
	switch o := obj.(type) {
	case *PyString:
		return o.Value
	default:
		return obj.String()
	}
}

func ToPyObject(value interface{}) object.Object {
	switch v := value.(type) {
	case int:
		return &PyInt{Value: v}
	case float64:
		return &PyFloat{Value: v}
	case string:
		return &PyString{Value: v}
	case bool:
		return &PyBool{Value: v}
	case nil:
		return &PyNone{}
	default:
		return &PyString{Value: fmt.Sprintf("%v", v)}
	}
}
