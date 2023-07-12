package model

import (
	"fmt"
)

//type Node struct {
//	Id     int
//	Indent int
//}

type Node interface {
	Id() int
}

type NodeInfo struct {
	Id     int
	Indent int
}

type Expression interface {
	Node
	expressionNode()
}

//type Statement struct {
//	Node
//}

type Statement interface {
	Node
	statementNode()
}

var staticID = 0

func NewNodeInfo(args map[string]int) *NodeInfo {
	self := &NodeInfo{}
	staticID++
	self.Id = staticID
	self.Indent += args["indent"]
	return self
}

type Boolean interface {
	Node
}

type NameBool struct {
	Name string
}

type Float struct {
	Value string
}

func (n *Float) expressionNode() {}
func (n *Float) Id() int         { return 0 }

type Integer struct {
	Value string
}

func (n *Integer) expressionNode() {}
func (n *Integer) Id() int         { return 0 }

// Type should be interface ?
type Type interface {
	Node
	Type() string
}

type NameType struct {
	Name string
}

type IntegerType struct {
}

type FloatType struct {
}

func (n *FloatType) Id() int      { return 0 }
func (n *FloatType) Type() string { return "" }

type Op struct {
	Value string
}

type UnaryOp struct {
	Operand Expression
}

type BinOpWithOp struct {
	Op    string
	Left  Expression
	Right Expression
}

func (n *BinOpWithOp) expressionNode() {}
func (n *BinOpWithOp) Id() int         { return 0 }

type BinOp struct {
	Left  Expression
	Right Expression
}

type RelOp struct {
	Left  Expression
	Right Expression
	Chain *RelOp
}

type LogicalOp struct {
	Left  Expression
	Right Expression
}

type CompareExp struct {
	Left   Expression
	Ops    []Op
	Values []Expression
}

type NumOp struct {
	BinOp
}

type PrintStatement struct {
	Value Expression
}

func (n *PrintStatement) Id() int { return 0 }

type Statements struct {
	Statements []Statement
}

type Name struct {
	Value string
}

type CompoundExpression struct {
	Statements []Statement
}

type ExpressionAsStatement struct {
	Expression Expression
}

type Grouping struct {
	Expression Expression
}

type Declaration struct {
}

// type is type
// if have field it should be declare

type ConstDeclaration struct {
	Name  Name
	Type  Type
	Value Expression
}

type VarDeclaration struct {
	Name  Name
	Type  Type
	Value Expression
}

type Assignment struct {
	Location Expression
	Value    Expression
}

type Character struct {
	Value rune
}

type IfStatement struct {
	Test        Expression
	Consequence Statements
	Alternative *Statements
}

type WhileStatement struct {
	Test Expression
	Body Statements
}

type BreakStatement struct {
}

type ContinueStatement struct {
}

type FunctionDeclaration struct {
	Name       Name
	Parameters []Parameter
	ReturnType Type
	Body       Statements
}

type FunctionApplication struct {
	Func      Expression
	Arguments []Expression
}

type ReturnStatement struct {
	Value Expression
}

type Parameter struct {
	Name Name
	Type Type
}

type Context struct {
	Indent string
}

func NewContext() *Context {
	return &Context{Indent: ""}
}

func (c *Context) NewBlock() *Context {
	return &Context{Indent: c.Indent + "    "}
}

func NodeAsSource(node Node, context *Context) string {
	switch v := node.(type) {
	case *Integer:
		return v.Value
	case *Float:
		return v.Value
	case *FloatType:
		return "float"
	case *BinOpWithOp:
		return fmt.Sprintf(" %s %s %s",
			NodeAsSource(v.Left, context), v.Op,
			NodeAsSource(v.Right, context))
	// ... 其他类型
	case *PrintStatement:
		return fmt.Sprintf("print %s", NodeAsSource(v.Value, context))
	default:
		panic(fmt.Sprintf("Can't convert %v to source", v))
	}
}
