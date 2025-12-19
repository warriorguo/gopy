package vm

import (
	"fmt"

	"github.com/warriorguo/gopy/pkg/compiler"
	"github.com/warriorguo/gopy/pkg/object"
	"github.com/warriorguo/gopy/pkg/runtime"
)

type Frame struct {
	Code     *compiler.CodeObject
	IP       int
	Stack    []object.Object
	SP       int
	Locals   []object.Object
	Globals  map[string]object.Object
	Builtins map[string]object.Object
}

func NewFrame(code *compiler.CodeObject, globals, builtins map[string]object.Object) *Frame {
	locals := make([]object.Object, len(code.Varnames))
	for i := range locals {
		locals[i] = &runtime.PyNone{}
	}

	return &Frame{
		Code:     code,
		IP:       0,
		Stack:    make([]object.Object, 1000),
		SP:       0,
		Locals:   locals,
		Globals:  globals,
		Builtins: builtins,
	}
}

func (f *Frame) push(obj object.Object) {
	f.Stack[f.SP] = obj
	f.SP++
}

func (f *Frame) pop() object.Object {
	if f.SP <= 0 {
		panic("stack underflow")
	}
	f.SP--
	return f.Stack[f.SP]
}

func (f *Frame) peek() object.Object {
	if f.SP <= 0 {
		return nil
	}
	return f.Stack[f.SP-1]
}

type VM struct {
	frames   []*Frame
	frameIdx int
	globals  map[string]object.Object
	builtins map[string]object.Object
}

func NewVM() *VM {
	builtins := make(map[string]object.Object)
	builtins["len"] = &compiler.PyBuiltin{
		Name: "len",
		Func: func(args []object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("len() takes exactly one argument (%d given)", len(args))
			}

			switch obj := args[0].(type) {
			case *runtime.PyString:
				return &runtime.PyInt{Value: len(obj.Value)}, nil
			case *runtime.PyList:
				return &runtime.PyInt{Value: len(obj.Elements)}, nil
			case *runtime.PyDict:
				return &runtime.PyInt{Value: len(obj.Pairs)}, nil
			default:
				return nil, fmt.Errorf("object of type '%s' has no len()", obj.Type())
			}
		},
	}

	builtins["range"] = &compiler.PyBuiltin{
		Name: "range",
		Func: func(args []object.Object) (object.Object, error) {
			if len(args) < 1 || len(args) > 3 {
				return nil, fmt.Errorf("range() takes 1 to 3 arguments")
			}

			var start, stop, step int
			var err error

			if len(args) == 1 {
				start = 0
				stop, err = toGoInt(args[0])
				if err != nil {
					return nil, err
				}
				step = 1
			} else if len(args) == 2 {
				start, err = toGoInt(args[0])
				if err != nil {
					return nil, err
				}
				stop, err = toGoInt(args[1])
				if err != nil {
					return nil, err
				}
				step = 1
			} else {
				start, err = toGoInt(args[0])
				if err != nil {
					return nil, err
				}
				stop, err = toGoInt(args[1])
				if err != nil {
					return nil, err
				}
				step, err = toGoInt(args[2])
				if err != nil {
					return nil, err
				}
				if step == 0 {
					return nil, fmt.Errorf("range() step argument must not be zero")
				}
			}

			var elements []object.Object
			if step > 0 {
				for i := start; i < stop; i += step {
					elements = append(elements, &runtime.PyInt{Value: i})
				}
			} else {
				for i := start; i > stop; i += step {
					elements = append(elements, &runtime.PyInt{Value: i})
				}
			}

			return &runtime.PyList{Elements: elements}, nil
		},
	}

	builtins["type"] = &compiler.PyBuiltin{
		Name: "type",
		Func: func(args []object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("type() takes exactly one argument")
			}
			return &runtime.PyString{Value: args[0].Type()}, nil
		},
	}

	builtins["str"] = &compiler.PyBuiltin{
		Name: "str",
		Func: func(args []object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("str() takes exactly one argument")
			}
			return &runtime.PyString{Value: toGoString(args[0])}, nil
		},
	}

	return &VM{
		frames:   make([]*Frame, 1000),
		frameIdx: -1,
		globals:  make(map[string]object.Object),
		builtins: builtins,
	}
}

func (vm *VM) pushFrame(frame *Frame) {
	vm.frameIdx++
	vm.frames[vm.frameIdx] = frame
}

func (vm *VM) popFrame() *Frame {
	if vm.frameIdx < 0 {
		return nil
	}
	frame := vm.frames[vm.frameIdx]
	vm.frameIdx--
	return frame
}

func (vm *VM) currentFrame() *Frame {
	if vm.frameIdx < 0 {
		return nil
	}
	return vm.frames[vm.frameIdx]
}

func (vm *VM) Run(code *compiler.CodeObject) (object.Object, error) {
	frame := NewFrame(code, vm.globals, vm.builtins)
	vm.pushFrame(frame)

	for vm.currentFrame() != nil {
		frame := vm.currentFrame()

		if frame.IP >= len(frame.Code.Instructions) {
			vm.popFrame()
			continue
		}

		instruction := frame.Code.Instructions[frame.IP]
		frame.IP++

		switch instruction.Op {
		case compiler.OpLoadConst:
			frame.push(frame.Code.Consts[instruction.Arg])

		case compiler.OpLoadName:
			name := frame.Code.Names[instruction.Arg]
			if obj, exists := frame.Globals[name]; exists {
				frame.push(obj)
			} else if obj, exists := frame.Builtins[name]; exists {
				frame.push(obj)
			} else {
				return nil, fmt.Errorf("name '%s' is not defined", name)
			}

		case compiler.OpStoreName:
			name := frame.Code.Names[instruction.Arg]
			frame.Globals[name] = frame.pop()

		case compiler.OpLoadGlobal:
			name := frame.Code.Names[instruction.Arg]
			if obj, exists := frame.Globals[name]; exists {
				frame.push(obj)
			} else if obj, exists := frame.Builtins[name]; exists {
				frame.push(obj)
			} else {
				return nil, fmt.Errorf("global name '%s' is not defined", name)
			}

		case compiler.OpStoreGlobal:
			name := frame.Code.Names[instruction.Arg]
			frame.Globals[name] = frame.pop()

		case compiler.OpLoadFast:
			frame.push(frame.Locals[instruction.Arg])

		case compiler.OpStoreFast:
			frame.Locals[instruction.Arg] = frame.pop()

		case compiler.OpBinaryAdd:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.binaryOp(left, right, "+")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpBinarySub:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.binaryOp(left, right, "-")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpBinaryMul:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.binaryOp(left, right, "*")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpBinaryDiv:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.binaryOp(left, right, "/")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpBinaryMod:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.binaryOp(left, right, "%")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpUnaryPos:
			operand := frame.pop()
			result, err := vm.unaryOp(operand, "+")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpUnaryNeg:
			operand := frame.pop()
			result, err := vm.unaryOp(operand, "-")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpUnaryNot:
			operand := frame.pop()
			result := &runtime.PyBool{Value: !operand.IsTruthy()}
			frame.push(result)

		case compiler.OpCompareEq:
			right := frame.pop()
			left := frame.pop()
			result := &runtime.PyBool{Value: left.Equal(right)}
			frame.push(result)

		case compiler.OpCompareNe:
			right := frame.pop()
			left := frame.pop()
			result := &runtime.PyBool{Value: !left.Equal(right)}
			frame.push(result)

		case compiler.OpCompareLt:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.compareOp(left, right, "<")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpCompareLe:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.compareOp(left, right, "<=")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpCompareGt:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.compareOp(left, right, ">")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpCompareGe:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.compareOp(left, right, ">=")
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpCompareIn:
			right := frame.pop()
			left := frame.pop()
			result, err := vm.inOp(left, right)
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpJumpForward:
			frame.IP += instruction.Arg

		case compiler.OpJumpIfFalse:
			if !frame.peek().IsTruthy() {
				frame.IP = instruction.Arg
			}

		case compiler.OpJumpIfTrue:
			if frame.peek().IsTruthy() {
				frame.IP = instruction.Arg
			}

		case compiler.OpJumpAbsolute:
			frame.IP = instruction.Arg

		case compiler.OpPopJumpIfFalse:
			obj := frame.pop()
			if !obj.IsTruthy() {
				frame.IP = instruction.Arg
			}

		case compiler.OpPopJumpIfTrue:
			obj := frame.pop()
			if obj.IsTruthy() {
				frame.IP = instruction.Arg
			}

		case compiler.OpBuildList:
			elements := make([]object.Object, instruction.Arg)
			for i := instruction.Arg - 1; i >= 0; i-- {
				elements[i] = frame.pop()
			}
			frame.push(&runtime.PyList{Elements: elements})

		case compiler.OpBuildDict:
			dict := runtime.NewPyDict()
			for i := 0; i < instruction.Arg; i++ {
				value := frame.pop()
				key := frame.pop()
				dict.Set(key, value)
			}
			frame.push(dict)

		case compiler.OpBinarySubscr:
			index := frame.pop()
			container := frame.pop()
			result, err := vm.subscript(container, index)
			if err != nil {
				return nil, err
			}
			frame.push(result)

		case compiler.OpStoreSubscr:
			index := frame.pop()
			container := frame.pop()
			value := frame.pop()
			err := vm.storeSubscript(container, index, value)
			if err != nil {
				return nil, err
			}

		case compiler.OpCallFunction:
			args := make([]object.Object, instruction.Arg)
			for i := instruction.Arg - 1; i >= 0; i-- {
				args[i] = frame.pop()
			}
			function := frame.pop()

			switch f := function.(type) {
			case *compiler.PyBuiltin:
				result, err := f.Func(args)
				if err != nil {
					return nil, err
				}
				frame.push(result)
			case *compiler.PyFunction:
				if len(args) != f.Code.Argcount {
					return nil, fmt.Errorf("function takes %d arguments but %d were given", f.Code.Argcount, len(args))
				}

				funcFrame := NewFrame(f.Code, vm.globals, vm.builtins)
				copy(funcFrame.Locals, args)
				vm.pushFrame(funcFrame)
				// Continue execution with the new frame - no result pushed yet
			default:
				return nil, fmt.Errorf("'%s' object is not callable", function.Type())
			}

		case compiler.OpReturnValue:
			result := frame.pop()
			vm.popFrame()
			if vm.currentFrame() != nil {
				vm.currentFrame().push(result)
			} else {
				return result, nil
			}

		case compiler.OpPrintExpr:
			obj := frame.pop()
			fmt.Print(obj.String())

		case compiler.OpPrintNewline:
			fmt.Println()

		case compiler.OpPopTop:
			frame.pop()

		case compiler.OpRotTwo:
			a := frame.pop()
			b := frame.pop()
			frame.push(a)
			frame.push(b)

		case compiler.OpRotThree:
			a := frame.pop()
			b := frame.pop()
			c := frame.pop()
			frame.push(a)
			frame.push(c)
			frame.push(b)

		case compiler.OpDupTop:
			frame.push(frame.peek())

		case compiler.OpGetIter:
			iterable := frame.pop()
			switch obj := iterable.(type) {
			case *runtime.PyList:
				frame.push(obj)
			case *runtime.PyString:
				// Convert string to list of characters
				var chars []object.Object
				for _, char := range obj.Value {
					chars = append(chars, &runtime.PyString{Value: string(char)})
				}
				frame.push(&runtime.PyList{Elements: chars})
			default:
				return nil, fmt.Errorf("OpGetIter: '%s' object is not iterable", iterable.Type())
			}

		case compiler.OpForIter:
			// Stack layout: [..., iterator_index, iterable]
			// Pop the iterable from the stack
			iterable := frame.pop()

			// Get or initialize the iterator index
			var iterIndex int
			if frame.SP > 0 {
				if idx, ok := frame.peek().(*runtime.PyInt); ok {
					iterIndex = idx.Value
					frame.pop() // Remove the old index
				} else {
					iterIndex = 0
				}
			} else {
				iterIndex = 0
			}

			if list, ok := iterable.(*runtime.PyList); ok {
				if iterIndex < len(list.Elements) {
					// Push the new index and current element
					frame.push(&runtime.PyInt{Value: iterIndex + 1})
					frame.push(list)
					frame.push(list.Elements[iterIndex])
				} else {
					// End of iteration - jump to end of loop
					frame.IP = instruction.Arg
				}
			} else {
				return nil, fmt.Errorf("OpForIter:'%s' object is not iterable", iterable.Type())
			}

		case compiler.OpNop:

		default:
			return nil, fmt.Errorf("unknown opcode: %d", instruction.Op)
		}
	}

	return &runtime.PyNone{}, nil
}

func (vm *VM) binaryOp(left, right object.Object, op string) (object.Object, error) {
	switch op {
	case "+":
		switch l := left.(type) {
		case *runtime.PyInt:
			switch r := right.(type) {
			case *runtime.PyInt:
				return &runtime.PyInt{Value: l.Value + r.Value}, nil
			case *runtime.PyFloat:
				return &runtime.PyFloat{Value: float64(l.Value) + r.Value}, nil
			}
		case *runtime.PyFloat:
			switch r := right.(type) {
			case *runtime.PyInt:
				return &runtime.PyFloat{Value: l.Value + float64(r.Value)}, nil
			case *runtime.PyFloat:
				return &runtime.PyFloat{Value: l.Value + r.Value}, nil
			}
		case *runtime.PyString:
			if r, ok := right.(*runtime.PyString); ok {
				return &runtime.PyString{Value: l.Value + r.Value}, nil
			}
		}

	case "-":
		switch l := left.(type) {
		case *runtime.PyInt:
			switch r := right.(type) {
			case *runtime.PyInt:
				return &runtime.PyInt{Value: l.Value - r.Value}, nil
			case *runtime.PyFloat:
				return &runtime.PyFloat{Value: float64(l.Value) - r.Value}, nil
			}
		case *runtime.PyFloat:
			switch r := right.(type) {
			case *runtime.PyInt:
				return &runtime.PyFloat{Value: l.Value - float64(r.Value)}, nil
			case *runtime.PyFloat:
				return &runtime.PyFloat{Value: l.Value - r.Value}, nil
			}
		}

	case "*":
		switch l := left.(type) {
		case *runtime.PyInt:
			switch r := right.(type) {
			case *runtime.PyInt:
				return &runtime.PyInt{Value: l.Value * r.Value}, nil
			case *runtime.PyFloat:
				return &runtime.PyFloat{Value: float64(l.Value) * r.Value}, nil
			}
		case *runtime.PyFloat:
			switch r := right.(type) {
			case *runtime.PyInt:
				return &runtime.PyFloat{Value: l.Value * float64(r.Value)}, nil
			case *runtime.PyFloat:
				return &runtime.PyFloat{Value: l.Value * r.Value}, nil
			}
		}

	case "/":
		switch l := left.(type) {
		case *runtime.PyInt:
			switch r := right.(type) {
			case *runtime.PyInt:
				if r.Value == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return &runtime.PyInt{Value: l.Value / r.Value}, nil
			case *runtime.PyFloat:
				if r.Value == 0.0 {
					return nil, fmt.Errorf("division by zero")
				}
				return &runtime.PyFloat{Value: float64(l.Value) / r.Value}, nil
			}
		case *runtime.PyFloat:
			switch r := right.(type) {
			case *runtime.PyInt:
				if r.Value == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return &runtime.PyFloat{Value: l.Value / float64(r.Value)}, nil
			case *runtime.PyFloat:
				if r.Value == 0.0 {
					return nil, fmt.Errorf("division by zero")
				}
				return &runtime.PyFloat{Value: l.Value / r.Value}, nil
			}
		}

	case "%":
		if l, ok := left.(*runtime.PyInt); ok {
			if r, ok := right.(*runtime.PyInt); ok {
				if r.Value == 0 {
					return nil, fmt.Errorf("integer division or modulo by zero")
				}
				return &runtime.PyInt{Value: l.Value % r.Value}, nil
			}
		}
	}

	return nil, fmt.Errorf("unsupported operand type(s) for %s: '%s' and '%s'", op, left.Type(), right.Type())
}

func (vm *VM) unaryOp(operand object.Object, op string) (object.Object, error) {
	switch op {
	case "+":
		switch o := operand.(type) {
		case *runtime.PyInt:
			return o, nil
		case *runtime.PyFloat:
			return o, nil
		default:
			return nil, fmt.Errorf("bad operand type for unary +: '%s'", operand.Type())
		}

	case "-":
		switch o := operand.(type) {
		case *runtime.PyInt:
			return &runtime.PyInt{Value: -o.Value}, nil
		case *runtime.PyFloat:
			return &runtime.PyFloat{Value: -o.Value}, nil
		default:
			return nil, fmt.Errorf("bad operand type for unary -: '%s'", operand.Type())
		}
	}

	return nil, fmt.Errorf("unknown unary operator: %s", op)
}

func (vm *VM) compareOp(left, right object.Object, op string) (object.Object, error) {
	switch l := left.(type) {
	case *runtime.PyInt:
		switch r := right.(type) {
		case *runtime.PyInt:
			switch op {
			case "<":
				return &runtime.PyBool{Value: l.Value < r.Value}, nil
			case "<=":
				return &runtime.PyBool{Value: l.Value <= r.Value}, nil
			case ">":
				return &runtime.PyBool{Value: l.Value > r.Value}, nil
			case ">=":
				return &runtime.PyBool{Value: l.Value >= r.Value}, nil
			}
		case *runtime.PyFloat:
			lf := float64(l.Value)
			switch op {
			case "<":
				return &runtime.PyBool{Value: lf < r.Value}, nil
			case "<=":
				return &runtime.PyBool{Value: lf <= r.Value}, nil
			case ">":
				return &runtime.PyBool{Value: lf > r.Value}, nil
			case ">=":
				return &runtime.PyBool{Value: lf >= r.Value}, nil
			}
		}
	case *runtime.PyFloat:
		switch r := right.(type) {
		case *runtime.PyInt:
			rf := float64(r.Value)
			switch op {
			case "<":
				return &runtime.PyBool{Value: l.Value < rf}, nil
			case "<=":
				return &runtime.PyBool{Value: l.Value <= rf}, nil
			case ">":
				return &runtime.PyBool{Value: l.Value > rf}, nil
			case ">=":
				return &runtime.PyBool{Value: l.Value >= rf}, nil
			}
		case *runtime.PyFloat:
			switch op {
			case "<":
				return &runtime.PyBool{Value: l.Value < r.Value}, nil
			case "<=":
				return &runtime.PyBool{Value: l.Value <= r.Value}, nil
			case ">":
				return &runtime.PyBool{Value: l.Value > r.Value}, nil
			case ">=":
				return &runtime.PyBool{Value: l.Value >= r.Value}, nil
			}
		}
	case *runtime.PyString:
		if r, ok := right.(*runtime.PyString); ok {
			switch op {
			case "<":
				return &runtime.PyBool{Value: l.Value < r.Value}, nil
			case "<=":
				return &runtime.PyBool{Value: l.Value <= r.Value}, nil
			case ">":
				return &runtime.PyBool{Value: l.Value > r.Value}, nil
			case ">=":
				return &runtime.PyBool{Value: l.Value >= r.Value}, nil
			}
		}
	}

	return nil, fmt.Errorf("'%s' not supported between instances of '%s' and '%s'", op, left.Type(), right.Type())
}

func (vm *VM) inOp(left, right object.Object) (object.Object, error) {
	switch container := right.(type) {
	case *runtime.PyList:
		for _, elem := range container.Elements {
			if left.Equal(elem) {
				return &runtime.PyBool{Value: true}, nil
			}
		}
		return &runtime.PyBool{Value: false}, nil
	case *runtime.PyDict:
		_, exists := container.Get(left)
		return &runtime.PyBool{Value: exists}, nil
	case *runtime.PyString:
		if str, ok := left.(*runtime.PyString); ok {
			found := false
			if len(str.Value) <= len(container.Value) {
				for i := 0; i <= len(container.Value)-len(str.Value); i++ {
					if container.Value[i:i+len(str.Value)] == str.Value {
						found = true
						break
					}
				}
			}
			return &runtime.PyBool{Value: found}, nil
		}
	}
	return nil, fmt.Errorf("argument of type '%s' is not iterable", right.Type())
}

func (vm *VM) subscript(container, index object.Object) (object.Object, error) {
	switch c := container.(type) {
	case *runtime.PyList:
		idx, err := runtime.ToGoInt(index)
		if err != nil {
			return nil, err
		}
		if idx < 0 || idx >= len(c.Elements) {
			return nil, fmt.Errorf("list index out of range")
		}
		return c.Elements[idx], nil
	case *runtime.PyDict:
		value, exists := c.Get(index)
		if !exists {
			return nil, fmt.Errorf("KeyError: %s", index.String())
		}
		return value, nil
	case *runtime.PyString:
		idx, err := runtime.ToGoInt(index)
		if err != nil {
			return nil, err
		}
		if idx < 0 || idx >= len(c.Value) {
			return nil, fmt.Errorf("string index out of range")
		}
		return &runtime.PyString{Value: string(c.Value[idx])}, nil
	}
	return nil, fmt.Errorf("'%s' object is not subscriptable", container.Type())
}

func (vm *VM) storeSubscript(container, index, value object.Object) error {
	switch c := container.(type) {
	case *runtime.PyList:
		idx, err := runtime.ToGoInt(index)
		if err != nil {
			return err
		}
		if idx < 0 || idx >= len(c.Elements) {
			return fmt.Errorf("list index out of range")
		}
		c.Elements[idx] = value
		return nil
	case *runtime.PyDict:
		c.Set(index, value)
		return nil
	}
	return fmt.Errorf("'%s' object does not support item assignment", container.Type())
}

