package ast

type Position struct {
	Line   int
	Column int
}

type Node interface {
	Pos() Position
	String() string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type Module struct {
	Body []Stmt
	Position Position
}

func (m *Module) Pos() Position { return m.Position }
func (m *Module) String() string { return "Module" }

type AssignStmt struct {
	Target Expr
	Value  Expr
	Position    Position
}

func (a *AssignStmt) Pos() Position { return a.Position }
func (a *AssignStmt) String() string { return "AssignStmt" }
func (a *AssignStmt) stmtNode() {}

type AugAssignStmt struct {
	Target Expr
	Op     string // "+=", "-=", etc.
	Value  Expr
	Position Position
}

func (a *AugAssignStmt) Pos() Position { return a.Position }
func (a *AugAssignStmt) String() string { return "AugAssignStmt" }
func (a *AugAssignStmt) stmtNode() {}

type ExprStmt struct {
	Expr Expr
	Position  Position
}

func (e *ExprStmt) Pos() Position { return e.Position }
func (e *ExprStmt) String() string { return "ExprStmt" }
func (e *ExprStmt) stmtNode() {}

type PrintStmt struct {
	Values []Expr
	Position    Position
}

func (p *PrintStmt) Pos() Position { return p.Position }
func (p *PrintStmt) String() string { return "PrintStmt" }
func (p *PrintStmt) stmtNode() {}

type IfStmt struct {
	Test   Expr
	Body   []Stmt
	Orelse []Stmt
	Position    Position
}

func (i *IfStmt) Pos() Position { return i.Position }
func (i *IfStmt) String() string { return "IfStmt" }
func (i *IfStmt) stmtNode() {}

type WhileStmt struct {
	Test Expr
	Body []Stmt
	Position  Position
}

func (w *WhileStmt) Pos() Position { return w.Position }
func (w *WhileStmt) String() string { return "WhileStmt" }
func (w *WhileStmt) stmtNode() {}

type ForStmt struct {
	Target Expr
	Iter   Expr
	Body   []Stmt
	Position    Position
}

func (f *ForStmt) Pos() Position { return f.Position }
func (f *ForStmt) String() string { return "ForStmt" }
func (f *ForStmt) stmtNode() {}

type FuncDef struct {
	Name string
	Args []string
	Body []Stmt
	Position  Position
}

func (f *FuncDef) Pos() Position { return f.Position }
func (f *FuncDef) String() string { return "FuncDef" }
func (f *FuncDef) stmtNode() {}

type ReturnStmt struct {
	Value Expr
	Position   Position
}

func (r *ReturnStmt) Pos() Position { return r.Position }
func (r *ReturnStmt) String() string { return "ReturnStmt" }
func (r *ReturnStmt) stmtNode() {}

type PassStmt struct {
	Position Position
}

func (p *PassStmt) Pos() Position { return p.Position }
func (p *PassStmt) String() string { return "PassStmt" }
func (p *PassStmt) stmtNode() {}

type BinaryOp struct {
	Left  Expr
	Op    string
	Right Expr
	Position   Position
}

func (b *BinaryOp) Pos() Position { return b.Position }
func (b *BinaryOp) String() string { return "BinaryOp" }
func (b *BinaryOp) exprNode() {}

type UnaryOp struct {
	Op   string
	Expr Expr
	Position  Position
}

func (u *UnaryOp) Pos() Position { return u.Position }
func (u *UnaryOp) String() string { return "UnaryOp" }
func (u *UnaryOp) exprNode() {}

type BoolOp struct {
	Op     string
	Values []Expr
	Position    Position
}

func (b *BoolOp) Pos() Position { return b.Position }
func (b *BoolOp) String() string { return "BoolOp" }
func (b *BoolOp) exprNode() {}

type Compare struct {
	Left  Expr
	Ops   []string
	Right []Expr
	Position   Position
}

func (c *Compare) Pos() Position { return c.Position }
func (c *Compare) String() string { return "Compare" }
func (c *Compare) exprNode() {}

type Call struct {
	Func Expr
	Args []Expr
	Position  Position
}

func (c *Call) Pos() Position { return c.Position }
func (c *Call) String() string { return "Call" }
func (c *Call) exprNode() {}

type Subscript struct {
	Value Expr
	Slice Expr
	Position   Position
}

func (s *Subscript) Pos() Position { return s.Position }
func (s *Subscript) String() string { return "Subscript" }
func (s *Subscript) exprNode() {}

type Name struct {
	Id  string
	Position Position
}

func (n *Name) Pos() Position { return n.Position }
func (n *Name) String() string { return "Name" }
func (n *Name) exprNode() {}

type Num struct {
	N   interface{}
	Position Position
}

func (n *Num) Pos() Position { return n.Position }
func (n *Num) String() string { return "Num" }
func (n *Num) exprNode() {}

type Str struct {
	S   string
	Position Position
}

func (s *Str) Pos() Position { return s.Position }
func (s *Str) String() string { return "Str" }
func (s *Str) exprNode() {}

type NameConstant struct {
	Value interface{}
	Position   Position
}

func (n *NameConstant) Pos() Position { return n.Position }
func (n *NameConstant) String() string { return "NameConstant" }
func (n *NameConstant) exprNode() {}

type List struct {
	Elts []Expr
	Position  Position
}

func (l *List) Pos() Position { return l.Position }
func (l *List) String() string { return "List" }
func (l *List) exprNode() {}

type Dict struct {
	Keys   []Expr
	Values []Expr
	Position    Position
}

func (d *Dict) Pos() Position { return d.Position }
func (d *Dict) String() string { return "Dict" }
func (d *Dict) exprNode() {}