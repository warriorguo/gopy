package compiler

import (
	"encoding/gob"
	"fmt"
	"io"

	"github.com/warriorguo/gopy/pkg/object"
	"github.com/warriorguo/gopy/pkg/runtime"
)

type OpCode byte

const (
	OpLoadConst OpCode = iota
	OpLoadName
	OpStoreName
	OpLoadGlobal
	OpStoreGlobal
	OpLoadFast
	OpStoreFast
	
	OpBinaryAdd
	OpBinarySub
	OpBinaryMul
	OpBinaryDiv
	OpBinaryMod
	
	OpUnaryPos
	OpUnaryNeg
	OpUnaryNot
	
	OpCompareEq
	OpCompareNe
	OpCompareLt
	OpCompareLe
	OpCompareGt
	OpCompareGe
	OpCompareIn
	
	OpJumpForward
	OpJumpIfFalse
	OpJumpIfTrue
	OpJumpAbsolute
	OpPopJumpIfFalse
	OpPopJumpIfTrue
	
	OpBuildList
	OpBuildDict
	OpBuildTuple
	
	OpBinarySubscr
	OpStoreSubscr
	
	OpCallFunction
	OpReturnValue
	
	OpPrintExpr
	OpPrintNewline
	
	OpPopTop
	OpRotTwo
	OpRotThree
	OpDupTop
	
	OpSetupLoop
	OpBreakLoop
	OpContinueLoop
	
	OpGetIter
	OpForIter
	
	OpNop
)

type Instruction struct {
	Op  OpCode
	Arg int
}

func (i Instruction) String() string {
	return fmt.Sprintf("%s %d", i.Op, i.Arg)
}

func (op OpCode) String() string {
	switch op {
	case OpLoadConst:
		return "LOAD_CONST"
	case OpLoadName:
		return "LOAD_NAME"
	case OpStoreName:
		return "STORE_NAME"
	case OpLoadGlobal:
		return "LOAD_GLOBAL"
	case OpStoreGlobal:
		return "STORE_GLOBAL"
	case OpLoadFast:
		return "LOAD_FAST"
	case OpStoreFast:
		return "STORE_FAST"
	case OpBinaryAdd:
		return "BINARY_ADD"
	case OpBinarySub:
		return "BINARY_SUB"
	case OpBinaryMul:
		return "BINARY_MUL"
	case OpBinaryDiv:
		return "BINARY_DIV"
	case OpBinaryMod:
		return "BINARY_MOD"
	case OpUnaryPos:
		return "UNARY_POS"
	case OpUnaryNeg:
		return "UNARY_NEG"
	case OpUnaryNot:
		return "UNARY_NOT"
	case OpCompareEq:
		return "COMPARE_EQ"
	case OpCompareNe:
		return "COMPARE_NE"
	case OpCompareLt:
		return "COMPARE_LT"
	case OpCompareLe:
		return "COMPARE_LE"
	case OpCompareGt:
		return "COMPARE_GT"
	case OpCompareGe:
		return "COMPARE_GE"
	case OpCompareIn:
		return "COMPARE_IN"
	case OpJumpForward:
		return "JUMP_FORWARD"
	case OpJumpIfFalse:
		return "JUMP_IF_FALSE"
	case OpJumpIfTrue:
		return "JUMP_IF_TRUE"
	case OpJumpAbsolute:
		return "JUMP_ABSOLUTE"
	case OpPopJumpIfFalse:
		return "POP_JUMP_IF_FALSE"
	case OpPopJumpIfTrue:
		return "POP_JUMP_IF_TRUE"
	case OpBuildList:
		return "BUILD_LIST"
	case OpBuildDict:
		return "BUILD_DICT"
	case OpBuildTuple:
		return "BUILD_TUPLE"
	case OpBinarySubscr:
		return "BINARY_SUBSCR"
	case OpStoreSubscr:
		return "STORE_SUBSCR"
	case OpCallFunction:
		return "CALL_FUNCTION"
	case OpReturnValue:
		return "RETURN_VALUE"
	case OpPrintExpr:
		return "PRINT_EXPR"
	case OpPrintNewline:
		return "PRINT_NEWLINE"
	case OpPopTop:
		return "POP_TOP"
	case OpRotTwo:
		return "ROT_TWO"
	case OpRotThree:
		return "ROT_THREE"
	case OpDupTop:
		return "DUP_TOP"
	case OpSetupLoop:
		return "SETUP_LOOP"
	case OpBreakLoop:
		return "BREAK_LOOP"
	case OpContinueLoop:
		return "CONTINUE_LOOP"
	case OpGetIter:
		return "GET_ITER"
	case OpForIter:
		return "FOR_ITER"
	case OpNop:
		return "NOP"
	default:
		return fmt.Sprintf("UNKNOWN_OP_%d", op)
	}
}


type CodeObject struct {
	Instructions []Instruction
	Consts       []object.Object
	Names        []string
	Varnames     []string
	Argcount     int
	Filename     string
	Name         string
	Firstlineno  int
}

func (co *CodeObject) String() string {
	return fmt.Sprintf("CodeObject{name=%s, argcount=%d, instructions=%d, consts=%d, names=%d}",
		co.Name, co.Argcount, len(co.Instructions), len(co.Consts), len(co.Names))
}

func (co *CodeObject) Disassemble() string {
	result := fmt.Sprintf("Code object: %s\n", co.Name)
	result += fmt.Sprintf("Args: %d, Consts: %d, Names: %d, Vars: %d\n", 
		co.Argcount, len(co.Consts), len(co.Names), len(co.Varnames))
	result += "\nConstants:\n"
	for i, c := range co.Consts {
		result += fmt.Sprintf("  %d: %s\n", i, c)
	}
	result += "\nNames:\n"
	for i, n := range co.Names {
		result += fmt.Sprintf("  %d: %s\n", i, n)
	}
	result += "\nVarnames:\n"
	for i, v := range co.Varnames {
		result += fmt.Sprintf("  %d: %s\n", i, v)
	}
	result += "\nInstructions:\n"
	for i, instr := range co.Instructions {
		result += fmt.Sprintf("  %3d: %s\n", i, instr)
	}
	return result
}

func (co *CodeObject) Serialize(w io.Writer) error {
	encoder := gob.NewEncoder(w)
	return encoder.Encode(co)
}

func DeserializeCodeObject(r io.Reader) (*CodeObject, error) {
	decoder := gob.NewDecoder(r)
	var co CodeObject
	err := decoder.Decode(&co)
	return &co, err
}

func init() {
	gob.Register(&runtime.PyInt{})
	gob.Register(&runtime.PyFloat{})
	gob.Register(&runtime.PyString{})
	gob.Register(&runtime.PyBool{})
	gob.Register(&runtime.PyNone{})
	gob.Register(&PyFunction{})
	gob.Register(&runtime.PyList{})
	gob.Register(&runtime.PyDict{})
	gob.Register(&PyBuiltin{})
}