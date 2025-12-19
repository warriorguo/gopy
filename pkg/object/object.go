package object

type Object interface {
	String() string
	Type() string
	IsTruthy() bool
	Equal(other Object) bool
}