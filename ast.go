package wg

import "io"

type Type string
type Emitter interface {
	wat(io.Writer)
}

const (
	V   Type = ""
	I32      = "i32"
	U32      = "i32"
	I64      = "i64"
	U64      = "i64"
	F32      = "f32"
	F64      = "f64"
)

type Module struct {
	Memory  string
	Globals []Assign
	Funcs   []Func
	Table   []TableEntries
	current *Func
}
type TableEntries struct {
	Off   int
	Names []string
}
type Func struct {
	Name string
	Args []Arg
	Rets []Type
	Locs []Local
	Body Stmts
	Doc  string
}
type Arg struct {
	Name string
	Type Type
}
type Local struct {
	Name string
	Type Type
}
type Stmt Emitter
type Stmts []Stmt
type Assign struct { //Stmt
	Name []string
	Expr []Expr
	Glob []bool
	Typs []Type
	Mod  string
}
type Return struct { //Stmt
	Last bool
	List []Expr
}
type Drop struct {
	Expr
}
type Nop struct{}
type Expr Emitter
type Unary struct { //Expr
	X  Expr
	Op Op
}
type Binary struct { //Expr
	X, Y Expr
	Op   Op
}
type Literal struct {
	Type  Type
	Value string
}
type GlobalGet string    //Expr
type GlobalGets []string //Expr
type LocalGet string     //Expr
type LocalGets []string  //Expr (struct)
type Op struct {
	Name string
	Type Type
}
type Call struct { //Expr
	Func string
	Args []Expr
}
type Cast struct {
	Dst, Src Type
	Arg      Expr
}
type CallIndirect struct {
	Func    Expr
	Args    []Expr
	ArgType []Type
	ResType []Type
}
type If struct {
	If         Expr
	Then, Else Stmts
}
