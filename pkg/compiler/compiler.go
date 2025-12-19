package compiler

import (
	"fmt"

	"github.com/warriorguo/gopy/pkg/ast"
	"github.com/warriorguo/gopy/pkg/object"
	"github.com/warriorguo/gopy/pkg/runtime"
)

type Compiler struct {
	instructions []Instruction
	consts       []object.Object
	names        []string
	varnames     []string
	constMap     map[string]int
	nameMap      map[string]int
	varnameMap   map[string]int
	loopStack    []int
	scopeDepth   int
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions: []Instruction{},
		consts:       []object.Object{},
		names:        []string{},
		varnames:     []string{},
		constMap:     make(map[string]int),
		nameMap:      make(map[string]int),
		varnameMap:   make(map[string]int),
		loopStack:    []int{},
		scopeDepth:   0,
	}
}

func (c *Compiler) emit(op OpCode, arg int) int {
	pos := len(c.instructions)
	c.instructions = append(c.instructions, Instruction{Op: op, Arg: arg})
	return pos
}

func (c *Compiler) changeOperand(pos int, arg int) {
	c.instructions[pos].Arg = arg
}

func (c *Compiler) addConstant(obj object.Object) int {
	key := fmt.Sprintf("%T:%s", obj, obj.String())
	if idx, exists := c.constMap[key]; exists {
		return idx
	}
	idx := len(c.consts)
	c.consts = append(c.consts, obj)
	c.constMap[key] = idx
	return idx
}

func (c *Compiler) addName(name string) int {
	if idx, exists := c.nameMap[name]; exists {
		return idx
	}
	idx := len(c.names)
	c.names = append(c.names, name)
	c.nameMap[name] = idx
	return idx
}

func (c *Compiler) addVarname(name string) int {
	if idx, exists := c.varnameMap[name]; exists {
		return idx
	}
	idx := len(c.varnames)
	c.varnames = append(c.varnames, name)
	c.varnameMap[name] = idx
	return idx
}

func (c *Compiler) Compile(node ast.Node) (*CodeObject, error) {
	switch n := node.(type) {
	case *ast.Module:
		return c.compileModule(n)
	case *ast.FuncDef:
		return c.compileFuncDef(n)
	default:
		return nil, fmt.Errorf("cannot compile node type %T", node)
	}
}

func (c *Compiler) compileModule(module *ast.Module) (*CodeObject, error) {
	for i, stmt := range module.Body {
		isLastStmt := i == len(module.Body)-1
		
		// Special handling for the last statement if it's an expression
		if isLastStmt {
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				// Compile the expression but don't pop it - return its value
				if err := c.compileExpr(exprStmt.Expr); err != nil {
					return nil, err
				}
				c.emit(OpReturnValue, 0)
				
				return &CodeObject{
					Instructions: c.instructions,
					Consts:       c.consts,
					Names:        c.names,
					Varnames:     c.varnames,
					Argcount:     0,
					Filename:     "<module>",
					Name:         "<module>",
					Firstlineno:  1,
				}, nil
			}
		}
		
		if err := c.compileStmt(stmt); err != nil {
			return nil, err
		}
	}

	c.emit(OpLoadConst, c.addConstant(&runtime.PyNone{}))
	c.emit(OpReturnValue, 0)

	return &CodeObject{
		Instructions: c.instructions,
		Consts:       c.consts,
		Names:        c.names,
		Varnames:     c.varnames,
		Argcount:     0,
		Filename:     "<module>",
		Name:         "<module>",
		Firstlineno:  1,
	}, nil
}

func (c *Compiler) compileFuncDef(funcDef *ast.FuncDef) (*CodeObject, error) {
	for _, arg := range funcDef.Args {
		c.addVarname(arg)
	}

	for _, stmt := range funcDef.Body {
		if err := c.compileStmt(stmt); err != nil {
			return nil, err
		}
	}

	c.emit(OpLoadConst, c.addConstant(&runtime.PyNone{}))
	c.emit(OpReturnValue, 0)

	return &CodeObject{
		Instructions: c.instructions,
		Consts:       c.consts,
		Names:        c.names,
		Varnames:     c.varnames,
		Argcount:     len(funcDef.Args),
		Filename:     "<function>",
		Name:         funcDef.Name,
		Firstlineno:  funcDef.Position.Line,
	}, nil
}

func (c *Compiler) compileStmt(stmt ast.Stmt) error {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		return c.compileAssignStmt(s)
	case *ast.AugAssignStmt:
		return c.compileAugAssignStmt(s)
	case *ast.ExprStmt:
		return c.compileExprStmt(s)
	case *ast.PrintStmt:
		return c.compilePrintStmt(s)
	case *ast.IfStmt:
		return c.compileIfStmt(s)
	case *ast.WhileStmt:
		return c.compileWhileStmt(s)
	case *ast.ForStmt:
		return c.compileForStmt(s)
	case *ast.FuncDef:
		return c.compileFuncDefStmt(s)
	case *ast.ReturnStmt:
		return c.compileReturnStmt(s)
	case *ast.PassStmt:
		return c.compilePassStmt(s)
	default:
		return fmt.Errorf("unsupported statement type: %T", stmt)
	}
}

func (c *Compiler) compileAssignStmt(stmt *ast.AssignStmt) error {
	if err := c.compileExpr(stmt.Value); err != nil {
		return err
	}

	switch target := stmt.Target.(type) {
	case *ast.Name:
		if c.scopeDepth == 0 {
			c.emit(OpStoreName, c.addName(target.Id))
		} else {
			c.emit(OpStoreFast, c.addVarname(target.Id))
		}
	case *ast.Subscript:
		if err := c.compileExpr(target.Value); err != nil {
			return err
		}
		if err := c.compileExpr(target.Slice); err != nil {
			return err
		}
		c.emit(OpStoreSubscr, 0)
	default:
		return fmt.Errorf("unsupported assignment target: %T", target)
	}

	return nil
}

func (c *Compiler) compileAugAssignStmt(stmt *ast.AugAssignStmt) error {
	// For augmented assignment like x += 5, we need to:
	// 1. Load the current value of x
	// 2. Load the value 5
	// 3. Perform the addition
	// 4. Store the result back to x
	
	// Load the current value of the target variable
	switch target := stmt.Target.(type) {
	case *ast.Name:
		if c.scopeDepth == 0 {
			c.emit(OpLoadName, c.addName(target.Id))
		} else {
			c.emit(OpLoadFast, c.addVarname(target.Id))
		}
	default:
		return fmt.Errorf("unsupported augmented assignment target: %T", target)
	}
	
	// Compile the right-hand side value
	if err := c.compileExpr(stmt.Value); err != nil {
		return err
	}
	
	// Emit the appropriate binary operation
	switch stmt.Op {
	case "+=":
		c.emit(OpBinaryAdd, 0)
	case "-=":
		c.emit(OpBinarySub, 0)
	default:
		return fmt.Errorf("unsupported augmented assignment operator: %s", stmt.Op)
	}
	
	// Store the result back to the target variable
	switch target := stmt.Target.(type) {
	case *ast.Name:
		if c.scopeDepth == 0 {
			c.emit(OpStoreName, c.addName(target.Id))
		} else {
			c.emit(OpStoreFast, c.addVarname(target.Id))
		}
	}
	
	return nil
}

func (c *Compiler) compileExprStmt(stmt *ast.ExprStmt) error {
	if err := c.compileExpr(stmt.Expr); err != nil {
		return err
	}
	c.emit(OpPopTop, 0)
	return nil
}

func (c *Compiler) compilePrintStmt(stmt *ast.PrintStmt) error {
	for i, value := range stmt.Values {
		if err := c.compileExpr(value); err != nil {
			return err
		}
		c.emit(OpPrintExpr, 0)
		if i < len(stmt.Values)-1 {
			c.emit(OpLoadConst, c.addConstant(&runtime.PyString{Value: " "}))
			c.emit(OpPrintExpr, 0)
		}
	}
	c.emit(OpPrintNewline, 0)
	return nil
}

func (c *Compiler) compileIfStmt(stmt *ast.IfStmt) error {
	if err := c.compileExpr(stmt.Test); err != nil {
		return err
	}

	jumpIfFalse := c.emit(OpPopJumpIfFalse, 0)

	for _, s := range stmt.Body {
		if err := c.compileStmt(s); err != nil {
			return err
		}
	}

	var jumpEnd int
	if len(stmt.Orelse) > 0 {
		jumpEnd = c.emit(OpJumpForward, 0)
	}

	c.changeOperand(jumpIfFalse, len(c.instructions))

	if len(stmt.Orelse) > 0 {
		for _, s := range stmt.Orelse {
			if err := c.compileStmt(s); err != nil {
				return err
			}
		}
		c.changeOperand(jumpEnd, len(c.instructions))
	}

	return nil
}

func (c *Compiler) compileWhileStmt(stmt *ast.WhileStmt) error {
	loopStart := len(c.instructions)
	c.loopStack = append(c.loopStack, loopStart)

	if err := c.compileExpr(stmt.Test); err != nil {
		return err
	}

	jumpIfFalse := c.emit(OpPopJumpIfFalse, 0)

	for _, s := range stmt.Body {
		if err := c.compileStmt(s); err != nil {
			return err
		}
	}

	c.emit(OpJumpAbsolute, loopStart)
	c.changeOperand(jumpIfFalse, len(c.instructions))

	c.loopStack = c.loopStack[:len(c.loopStack)-1]
	return nil
}

func (c *Compiler) compileForStmt(stmt *ast.ForStmt) error {
	if err := c.compileExpr(stmt.Iter); err != nil {
		return err
	}

	c.emit(OpGetIter, 0)
	loopStart := len(c.instructions)
	c.loopStack = append(c.loopStack, loopStart)

	forIter := c.emit(OpForIter, 0)

	switch target := stmt.Target.(type) {
	case *ast.Name:
		if c.scopeDepth == 0 {
			c.emit(OpStoreName, c.addName(target.Id))
		} else {
			c.emit(OpStoreFast, c.addVarname(target.Id))
		}
	default:
		return fmt.Errorf("unsupported for target: %T", target)
	}

	for _, s := range stmt.Body {
		if err := c.compileStmt(s); err != nil {
			return err
		}
	}

	c.emit(OpJumpAbsolute, loopStart)
	c.changeOperand(forIter, len(c.instructions))

	c.loopStack = c.loopStack[:len(c.loopStack)-1]
	return nil
}

func (c *Compiler) compileFuncDefStmt(stmt *ast.FuncDef) error {
	compiler := NewCompiler()
	compiler.scopeDepth = c.scopeDepth + 1

	codeObj, err := compiler.compileFuncDef(stmt)
	if err != nil {
		return err
	}

	pyFunc := &PyFunction{
		Code:    codeObj,
		Name:    stmt.Name,
		Globals: nil, // Will be set by VM at runtime
	}
	constIdx := c.addConstant(pyFunc)

	c.emit(OpLoadConst, constIdx)
	c.emit(OpStoreName, c.addName(stmt.Name))

	return nil
}

func (c *Compiler) compileReturnStmt(stmt *ast.ReturnStmt) error {
	if stmt.Value != nil {
		if err := c.compileExpr(stmt.Value); err != nil {
			return err
		}
	} else {
		c.emit(OpLoadConst, c.addConstant(&runtime.PyNone{}))
	}
	c.emit(OpReturnValue, 0)
	return nil
}

func (c *Compiler) compilePassStmt(stmt *ast.PassStmt) error {
	// Pass statement is a no-op, emit nothing
	return nil
}

func (c *Compiler) compileExpr(expr ast.Expr) error {
	switch e := expr.(type) {
	case *ast.BinaryOp:
		return c.compileBinaryOp(e)
	case *ast.UnaryOp:
		return c.compileUnaryOp(e)
	case *ast.BoolOp:
		return c.compileBoolOp(e)
	case *ast.Compare:
		return c.compileCompare(e)
	case *ast.Call:
		return c.compileCall(e)
	case *ast.Subscript:
		return c.compileSubscript(e)
	case *ast.Name:
		return c.compileName(e)
	case *ast.Num:
		return c.compileNum(e)
	case *ast.Str:
		return c.compileStr(e)
	case *ast.NameConstant:
		return c.compileNameConstant(e)
	case *ast.List:
		return c.compileList(e)
	case *ast.Dict:
		return c.compileDict(e)
	default:
		return fmt.Errorf("unsupported expression type: %T", expr)
	}
}

func (c *Compiler) compileBinaryOp(expr *ast.BinaryOp) error {
	if err := c.compileExpr(expr.Left); err != nil {
		return err
	}
	if err := c.compileExpr(expr.Right); err != nil {
		return err
	}

	switch expr.Op {
	case "+":
		c.emit(OpBinaryAdd, 0)
	case "-":
		c.emit(OpBinarySub, 0)
	case "*":
		c.emit(OpBinaryMul, 0)
	case "/":
		c.emit(OpBinaryDiv, 0)
	case "%":
		c.emit(OpBinaryMod, 0)
	default:
		return fmt.Errorf("unsupported binary operator: %s", expr.Op)
	}

	return nil
}

func (c *Compiler) compileUnaryOp(expr *ast.UnaryOp) error {
	if err := c.compileExpr(expr.Expr); err != nil {
		return err
	}

	switch expr.Op {
	case "+":
		c.emit(OpUnaryPos, 0)
	case "-":
		c.emit(OpUnaryNeg, 0)
	case "not":
		c.emit(OpUnaryNot, 0)
	default:
		return fmt.Errorf("unsupported unary operator: %s", expr.Op)
	}

	return nil
}

func (c *Compiler) compileBoolOp(expr *ast.BoolOp) error {
	if len(expr.Values) < 2 {
		return fmt.Errorf("boolean operation needs at least 2 values")
	}

	if err := c.compileExpr(expr.Values[0]); err != nil {
		return err
	}

	var jumpPositions []int

	for i := 1; i < len(expr.Values); i++ {
		if expr.Op == "and" {
			c.emit(OpDupTop, 0)
			jumpPositions = append(jumpPositions, c.emit(OpPopJumpIfFalse, 0))
			c.emit(OpPopTop, 0)
		} else {
			c.emit(OpDupTop, 0)
			jumpPositions = append(jumpPositions, c.emit(OpPopJumpIfTrue, 0))
			c.emit(OpPopTop, 0)
		}

		if err := c.compileExpr(expr.Values[i]); err != nil {
			return err
		}
	}

	for _, pos := range jumpPositions {
		c.changeOperand(pos, len(c.instructions))
	}

	return nil
}

func (c *Compiler) compileCompare(expr *ast.Compare) error {
	if err := c.compileExpr(expr.Left); err != nil {
		return err
	}

	for i, op := range expr.Ops {
		if err := c.compileExpr(expr.Right[i]); err != nil {
			return err
		}

		switch op {
		case "==":
			c.emit(OpCompareEq, 0)
		case "!=":
			c.emit(OpCompareNe, 0)
		case "<":
			c.emit(OpCompareLt, 0)
		case "<=":
			c.emit(OpCompareLe, 0)
		case ">":
			c.emit(OpCompareGt, 0)
		case ">=":
			c.emit(OpCompareGe, 0)
		case "in":
			c.emit(OpCompareIn, 0)
		default:
			return fmt.Errorf("unsupported comparison operator: %s", op)
		}
	}

	return nil
}

func (c *Compiler) compileCall(expr *ast.Call) error {
	if err := c.compileExpr(expr.Func); err != nil {
		return err
	}

	for _, arg := range expr.Args {
		if err := c.compileExpr(arg); err != nil {
			return err
		}
	}

	c.emit(OpCallFunction, len(expr.Args))
	return nil
}

func (c *Compiler) compileSubscript(expr *ast.Subscript) error {
	if err := c.compileExpr(expr.Value); err != nil {
		return err
	}
	if err := c.compileExpr(expr.Slice); err != nil {
		return err
	}
	c.emit(OpBinarySubscr, 0)
	return nil
}

func (c *Compiler) compileName(expr *ast.Name) error {
	if c.scopeDepth == 0 {
		c.emit(OpLoadName, c.addName(expr.Id))
	} else {
		if _, exists := c.varnameMap[expr.Id]; exists {
			c.emit(OpLoadFast, c.addVarname(expr.Id))
		} else {
			c.emit(OpLoadGlobal, c.addName(expr.Id))
		}
	}
	return nil
}

func (c *Compiler) compileNum(expr *ast.Num) error {
	obj := runtime.ToPyObject(expr.N)
	c.emit(OpLoadConst, c.addConstant(obj))
	return nil
}

func (c *Compiler) compileStr(expr *ast.Str) error {
	obj := &runtime.PyString{Value: expr.S}
	c.emit(OpLoadConst, c.addConstant(obj))
	return nil
}

func (c *Compiler) compileNameConstant(expr *ast.NameConstant) error {
	obj := runtime.ToPyObject(expr.Value)
	c.emit(OpLoadConst, c.addConstant(obj))
	return nil
}

func (c *Compiler) compileList(expr *ast.List) error {
	for _, elt := range expr.Elts {
		if err := c.compileExpr(elt); err != nil {
			return err
		}
	}
	c.emit(OpBuildList, len(expr.Elts))
	return nil
}

func (c *Compiler) compileDict(expr *ast.Dict) error {
	for i := range expr.Keys {
		if err := c.compileExpr(expr.Keys[i]); err != nil {
			return err
		}
		if err := c.compileExpr(expr.Values[i]); err != nil {
			return err
		}
	}
	c.emit(OpBuildDict, len(expr.Keys))
	return nil
}

func Compile(module *ast.Module) (*CodeObject, error) {
	compiler := NewCompiler()
	return compiler.Compile(module)
}