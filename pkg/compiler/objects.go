package compiler

import (
	"fmt"
	"github.com/warriorguo/gopy/pkg/object"
)

type PyFunction struct {
	Code    *CodeObject
	Name    string
	Globals map[string]object.Object
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