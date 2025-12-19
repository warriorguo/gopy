package ast

import (
	"fmt"
	"strings"
)

// ASTFormatter provides formatted output for AST nodes with proper indentation
type ASTFormatter struct {
	indent       string
	currentLevel int
}

// NewASTFormatter creates a new formatter instance
func NewASTFormatter() *ASTFormatter {
	return &ASTFormatter{
		indent:       "  ", // Two spaces per level
		currentLevel: 0,
	}
}

// SetIndent allows customization of the indentation string
func (f *ASTFormatter) SetIndent(indent string) {
	f.indent = indent
}

// Format returns a formatted string representation of any AST node
func (f *ASTFormatter) Format(node Node) string {
	if node == nil {
		return f.getIndent() + "<nil>"
	}
	return f.formatNode(node)
}

// FormatModule is a convenience function for formatting entire modules
func (f *ASTFormatter) FormatModule(module *Module) string {
	f.currentLevel = 0
	return f.Format(module)
}

// getIndent returns the current indentation string
func (f *ASTFormatter) getIndent() string {
	return strings.Repeat(f.indent, f.currentLevel)
}

// formatNode formats any AST node with appropriate indentation
func (f *ASTFormatter) formatNode(node Node) string {
	switch n := node.(type) {
	// Statements
	case *Module:
		return f.formatModule(n)
	case *AssignStmt:
		return f.formatAssignStmt(n)
	case *AugAssignStmt:
		return f.formatAugAssignStmt(n)
	case *ExprStmt:
		return f.formatExprStmt(n)
	case *PrintStmt:
		return f.formatPrintStmt(n)
	case *IfStmt:
		return f.formatIfStmt(n)
	case *WhileStmt:
		return f.formatWhileStmt(n)
	case *ForStmt:
		return f.formatForStmt(n)
	case *FuncDef:
		return f.formatFuncDef(n)
	case *ReturnStmt:
		return f.formatReturnStmt(n)
	case *PassStmt:
		return f.formatPassStmt(n)
	
	// Expressions
	case *BinaryOp:
		return f.formatBinaryOp(n)
	case *UnaryOp:
		return f.formatUnaryOp(n)
	case *BoolOp:
		return f.formatBoolOp(n)
	case *Compare:
		return f.formatCompare(n)
	case *Call:
		return f.formatCall(n)
	case *Subscript:
		return f.formatSubscript(n)
	case *Name:
		return f.formatName(n)
	case *Num:
		return f.formatNum(n)
	case *Str:
		return f.formatStr(n)
	case *NameConstant:
		return f.formatNameConstant(n)
	case *List:
		return f.formatList(n)
	case *Dict:
		return f.formatDict(n)
	
	default:
		return f.getIndent() + fmt.Sprintf("UnknownNode<%T>", node)
	}
}

// formatModule formats a module node
func (f *ASTFormatter) formatModule(m *Module) string {
	result := f.getIndent() + fmt.Sprintf("Module (pos: %d:%d)\n", m.Position.Line, m.Position.Column)
	f.currentLevel++
	
	if len(m.Body) == 0 {
		result += f.getIndent() + "Body: <empty>\n"
	} else {
		result += f.getIndent() + "Body:\n"
		f.currentLevel++
		for i, stmt := range m.Body {
			result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(stmt))
			if i < len(m.Body)-1 {
				result += "\n"
			}
		}
		f.currentLevel--
	}
	
	f.currentLevel--
	return result
}

// formatAssignStmt formats an assignment statement
func (f *ASTFormatter) formatAssignStmt(a *AssignStmt) string {
	result := fmt.Sprintf("AssignStmt (pos: %d:%d)\n", a.Position.Line, a.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Target: " + f.formatNode(a.Target) + "\n"
	result += f.getIndent() + "Value: " + f.formatNode(a.Value)
	
	f.currentLevel--
	return result
}

// formatAugAssignStmt formats an augmented assignment statement
func (f *ASTFormatter) formatAugAssignStmt(a *AugAssignStmt) string {
	result := fmt.Sprintf("AugAssignStmt (pos: %d:%d)\n", a.Position.Line, a.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Target: " + f.formatNode(a.Target) + "\n"
	result += f.getIndent() + fmt.Sprintf("Op: %q\n", a.Op)
	result += f.getIndent() + "Value: " + f.formatNode(a.Value)
	
	f.currentLevel--
	return result
}

// formatExprStmt formats an expression statement
func (f *ASTFormatter) formatExprStmt(e *ExprStmt) string {
	result := fmt.Sprintf("ExprStmt (pos: %d:%d)\n", e.Position.Line, e.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Expr: " + f.formatNode(e.Expr)
	
	f.currentLevel--
	return result
}

// formatPrintStmt formats a print statement
func (f *ASTFormatter) formatPrintStmt(p *PrintStmt) string {
	result := fmt.Sprintf("PrintStmt (pos: %d:%d)\n", p.Position.Line, p.Position.Column)
	f.currentLevel++
	
	if len(p.Values) == 0 {
		result += f.getIndent() + "Values: <none>"
	} else {
		result += f.getIndent() + "Values:\n"
		f.currentLevel++
		for i, val := range p.Values {
			result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(val))
			if i < len(p.Values)-1 {
				result += "\n"
			}
		}
		f.currentLevel--
	}
	
	f.currentLevel--
	return result
}

// formatIfStmt formats an if statement
func (f *ASTFormatter) formatIfStmt(i *IfStmt) string {
	result := fmt.Sprintf("IfStmt (pos: %d:%d)\n", i.Position.Line, i.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Test: " + f.formatNode(i.Test) + "\n"
	
	// Format body
	result += f.getIndent() + "Body:\n"
	f.currentLevel++
	for j, stmt := range i.Body {
		result += f.getIndent() + fmt.Sprintf("[%d] %s", j, f.formatNode(stmt))
		if j < len(i.Body)-1 {
			result += "\n"
		}
	}
	f.currentLevel--
	
	// Format else clause if present
	if len(i.Orelse) > 0 {
		result += "\n" + f.getIndent() + "Orelse:\n"
		f.currentLevel++
		for j, stmt := range i.Orelse {
			result += f.getIndent() + fmt.Sprintf("[%d] %s", j, f.formatNode(stmt))
			if j < len(i.Orelse)-1 {
				result += "\n"
			}
		}
		f.currentLevel--
	}
	
	f.currentLevel--
	return result
}

// formatWhileStmt formats a while statement
func (f *ASTFormatter) formatWhileStmt(w *WhileStmt) string {
	result := fmt.Sprintf("WhileStmt (pos: %d:%d)\n", w.Position.Line, w.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Test: " + f.formatNode(w.Test) + "\n"
	result += f.getIndent() + "Body:\n"
	f.currentLevel++
	for i, stmt := range w.Body {
		result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(stmt))
		if i < len(w.Body)-1 {
			result += "\n"
		}
	}
	f.currentLevel--
	
	f.currentLevel--
	return result
}

// formatForStmt formats a for statement
func (f *ASTFormatter) formatForStmt(fs *ForStmt) string {
	result := fmt.Sprintf("ForStmt (pos: %d:%d)\n", fs.Position.Line, fs.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Target: " + f.formatNode(fs.Target) + "\n"
	result += f.getIndent() + "Iter: " + f.formatNode(fs.Iter) + "\n"
	result += f.getIndent() + "Body:\n"
	f.currentLevel++
	for i, stmt := range fs.Body {
		result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(stmt))
		if i < len(fs.Body)-1 {
			result += "\n"
		}
	}
	f.currentLevel--
	
	f.currentLevel--
	return result
}

// formatFuncDef formats a function definition
func (f *ASTFormatter) formatFuncDef(fd *FuncDef) string {
	result := fmt.Sprintf("FuncDef (pos: %d:%d)\n", fd.Position.Line, fd.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + fmt.Sprintf("Name: %q\n", fd.Name)
	result += f.getIndent() + fmt.Sprintf("Args: %v\n", fd.Args)
	result += f.getIndent() + "Body:\n"
	f.currentLevel++
	for i, stmt := range fd.Body {
		result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(stmt))
		if i < len(fd.Body)-1 {
			result += "\n"
		}
	}
	f.currentLevel--
	
	f.currentLevel--
	return result
}

// formatReturnStmt formats a return statement
func (f *ASTFormatter) formatReturnStmt(r *ReturnStmt) string {
	result := fmt.Sprintf("ReturnStmt (pos: %d:%d)\n", r.Position.Line, r.Position.Column)
	f.currentLevel++
	
	if r.Value != nil {
		result += f.getIndent() + "Value: " + f.formatNode(r.Value)
	} else {
		result += f.getIndent() + "Value: <none>"
	}
	
	f.currentLevel--
	return result
}

// formatPassStmt formats a pass statement
func (f *ASTFormatter) formatPassStmt(p *PassStmt) string {
	return fmt.Sprintf("PassStmt (pos: %d:%d)", p.Position.Line, p.Position.Column)
}

// formatBinaryOp formats a binary operation
func (f *ASTFormatter) formatBinaryOp(b *BinaryOp) string {
	result := fmt.Sprintf("BinaryOp (pos: %d:%d)\n", b.Position.Line, b.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Left: " + f.formatNode(b.Left) + "\n"
	result += f.getIndent() + fmt.Sprintf("Op: %q\n", b.Op)
	result += f.getIndent() + "Right: " + f.formatNode(b.Right)
	
	f.currentLevel--
	return result
}

// formatUnaryOp formats a unary operation
func (f *ASTFormatter) formatUnaryOp(u *UnaryOp) string {
	result := fmt.Sprintf("UnaryOp (pos: %d:%d)\n", u.Position.Line, u.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + fmt.Sprintf("Op: %q\n", u.Op)
	result += f.getIndent() + "Expr: " + f.formatNode(u.Expr)
	
	f.currentLevel--
	return result
}

// formatBoolOp formats a boolean operation
func (f *ASTFormatter) formatBoolOp(b *BoolOp) string {
	result := fmt.Sprintf("BoolOp (pos: %d:%d)\n", b.Position.Line, b.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + fmt.Sprintf("Op: %q\n", b.Op)
	result += f.getIndent() + "Values:\n"
	f.currentLevel++
	for i, val := range b.Values {
		result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(val))
		if i < len(b.Values)-1 {
			result += "\n"
		}
	}
	f.currentLevel--
	
	f.currentLevel--
	return result
}

// formatCompare formats a comparison operation
func (f *ASTFormatter) formatCompare(c *Compare) string {
	result := fmt.Sprintf("Compare (pos: %d:%d)\n", c.Position.Line, c.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Left: " + f.formatNode(c.Left) + "\n"
	result += f.getIndent() + fmt.Sprintf("Ops: %v\n", c.Ops)
	result += f.getIndent() + "Right:\n"
	f.currentLevel++
	for i, right := range c.Right {
		result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(right))
		if i < len(c.Right)-1 {
			result += "\n"
		}
	}
	f.currentLevel--
	
	f.currentLevel--
	return result
}

// formatCall formats a function call
func (f *ASTFormatter) formatCall(c *Call) string {
	result := fmt.Sprintf("Call (pos: %d:%d)\n", c.Position.Line, c.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Func: " + f.formatNode(c.Func) + "\n"
	if len(c.Args) == 0 {
		result += f.getIndent() + "Args: <none>"
	} else {
		result += f.getIndent() + "Args:\n"
		f.currentLevel++
		for i, arg := range c.Args {
			result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(arg))
			if i < len(c.Args)-1 {
				result += "\n"
			}
		}
		f.currentLevel--
	}
	
	f.currentLevel--
	return result
}

// formatSubscript formats a subscript operation
func (f *ASTFormatter) formatSubscript(s *Subscript) string {
	result := fmt.Sprintf("Subscript (pos: %d:%d)\n", s.Position.Line, s.Position.Column)
	f.currentLevel++
	
	result += f.getIndent() + "Value: " + f.formatNode(s.Value) + "\n"
	result += f.getIndent() + "Slice: " + f.formatNode(s.Slice)
	
	f.currentLevel--
	return result
}

// formatName formats a name (identifier)
func (f *ASTFormatter) formatName(n *Name) string {
	return fmt.Sprintf("Name (pos: %d:%d) Id: %q", n.Position.Line, n.Position.Column, n.Id)
}

// formatNum formats a number literal
func (f *ASTFormatter) formatNum(n *Num) string {
	return fmt.Sprintf("Num (pos: %d:%d) Value: %v", n.Position.Line, n.Position.Column, n.N)
}

// formatStr formats a string literal
func (f *ASTFormatter) formatStr(s *Str) string {
	return fmt.Sprintf("Str (pos: %d:%d) Value: %q", s.Position.Line, s.Position.Column, s.S)
}

// formatNameConstant formats a name constant (True, False, None)
func (f *ASTFormatter) formatNameConstant(nc *NameConstant) string {
	return fmt.Sprintf("NameConstant (pos: %d:%d) Value: %v", nc.Position.Line, nc.Position.Column, nc.Value)
}

// formatList formats a list literal
func (f *ASTFormatter) formatList(l *List) string {
	result := fmt.Sprintf("List (pos: %d:%d)", l.Position.Line, l.Position.Column)
	if len(l.Elts) == 0 {
		result += " Elements: <empty>"
	} else {
		result += "\n"
		f.currentLevel++
		result += f.getIndent() + "Elements:\n"
		f.currentLevel++
		for i, elt := range l.Elts {
			result += f.getIndent() + fmt.Sprintf("[%d] %s", i, f.formatNode(elt))
			if i < len(l.Elts)-1 {
				result += "\n"
			}
		}
		f.currentLevel--
		f.currentLevel--
	}
	return result
}

// formatDict formats a dictionary literal
func (f *ASTFormatter) formatDict(d *Dict) string {
	result := fmt.Sprintf("Dict (pos: %d:%d)", d.Position.Line, d.Position.Column)
	if len(d.Keys) == 0 {
		result += " Pairs: <empty>"
	} else {
		result += "\n"
		f.currentLevel++
		result += f.getIndent() + "Pairs:\n"
		f.currentLevel++
		for i := 0; i < len(d.Keys); i++ {
			result += f.getIndent() + fmt.Sprintf("[%d] Key: %s\n", i, f.formatNode(d.Keys[i]))
			result += f.getIndent() + fmt.Sprintf("    Value: %s", f.formatNode(d.Values[i]))
			if i < len(d.Keys)-1 {
				result += "\n"
			}
		}
		f.currentLevel--
		f.currentLevel--
	}
	return result
}

// Convenience functions for quick formatting

// FormatAST formats an AST node with default settings
func FormatAST(node Node) string {
	formatter := NewASTFormatter()
	return formatter.Format(node)
}

// FormatASTWithIndent formats an AST node with custom indentation
func FormatASTWithIndent(node Node, indent string) string {
	formatter := NewASTFormatter()
	formatter.SetIndent(indent)
	return formatter.Format(node)
}

// PrintAST prints an AST node to stdout
func PrintAST(node Node) {
	fmt.Println(FormatAST(node))
}