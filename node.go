package wabbit_go

type Node struct {
	Id     int
	Indent int
}

var staticID = 0

func NewNode(args map[string]int) *Node {
	self := &Node{}
	staticID++
	self.Id = staticID
	self.Indent += args["indent"]
	return self
}

type Expression struct {
	Node
}

type Statement struct {
	Node
}

type Boolean struct {
	Expression
}

type NameBool struct {
	Expression
	Name string
}

type Float struct {
	Expression
	Value float64
}

type Integer struct {
	Expression
	Value int
}

type Type struct {
	Node
}

type NameType struct {
	Type
	Name string
}

type IntegerType struct {
	Type
}

type FloatType struct {
	Type
}

type Op struct {
	Node
	Value string
}

type UnaryOp struct {
	Expression
	Operand Expression
}

type BinOpWithOp struct {
	Node
	Op    string
	Left  Expression
	Right Expression
}

type BinOp struct {
	Expression
	Left  Expression
	Right Expression
}

type RelOp struct {
	Expression
	Left  Expression
	Right Expression
	Chain *RelOp
}

type LogicalOp struct {
	Expression
	Left  Expression
	Right Expression
}

type CompareExp struct {
	Expression
	Left   Expression
	Ops    []Op
	Values []Expression
}

type NumOp struct {
	BinOp
}

type PrintStatement struct {
	Statement
	Value Expression
}

type Statements struct {
	Node
	Statements []Statement
}

type Name struct {
	Expression
	Value string
}

type CompoundExpression struct {
	Expression
	Statements []Statement
}

type ExpressionAsStatement struct {
	Statement
	Expression Expression
}

type Grouping struct {
	Expression
}

type Declaration struct {
	Statement
}

type ConstDeclaration struct {
	Declaration
	Name  Name
	Type  Type
	Value Expression
}

type VarDeclaration struct {
	Declaration
	Name  Name
	Type  Type
	Value Expression
}

type Assignment struct {
	Expression
	Location Expression
	Value    Expression
}

type Character struct {
	Expression
	Value rune
}

type IfStatement struct {
	Statement
	Test        Expression
	Consequence Statements
	Alternative *Statements
}

type WhileStatement struct {
	Statement
	Test Expression
	Body Statements
}

type BreakStatement struct {
	Statement
}

type ContinueStatement struct {
	Statement
}

type FunctionDeclaration struct {
	Declaration
	Name       Name
	Parameters []Parameter
	ReturnType Type
	Body       Statements
}

type FunctionApplication struct {
	Expression
	Func      Expression
	Arguments []Expression
}

type ReturnStatement struct {
	Statement
	Value Expression
}

type Parameter struct {
	Node
	Name Name
	Type Type
}

type Context struct {
	Indent string
}
