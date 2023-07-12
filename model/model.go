package model

import (
	"fmt"
	"strings"
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
	ExpressionNode()
}

//type Statement struct {
//	Node
//}

type Statement interface {
	Node
	StatementNode()
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

func (n *Float) ExpressionNode() {}
func (n *Float) Id() int         { return 0 }

type Integer struct {
	Value string
}

func (n *Integer) ExpressionNode() {}
func (n *Integer) Id() int         { return 0 }

// Type should be interface ?
type Type interface {
	Node
	Type() string
}

type NameType struct {
	Name string
}

func (n *NameType) Id() int      { return 0 }
func (n *NameType) Type() string { return "" }

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

type Neg struct {
	Operand Expression
}

func (n *Neg) ExpressionNode() {}
func (n *Neg) Id() int         { return 0 }

type Pos struct {
	Operand Expression
}

func (n *Pos) ExpressionNode() {}
func (n *Pos) Id() int         { return 0 }

type BinOpWithOp struct {
	Op    string
	Left  Expression
	Right Expression
}

func (n *BinOpWithOp) ExpressionNode() {}
func (n *BinOpWithOp) Id() int         { return 0 }

type BinOp struct {
	Left  Expression
	Right Expression
}

func (n *BinOp) ExpressionNode() {}
func (n *BinOp) Id() int         { return 0 }

type Add struct {
	Left  Expression
	Right Expression
}

func (n *Add) ExpressionNode() {}
func (n *Add) Id() int         { return 0 }

type Sub struct {
	Left  Expression
	Right Expression
}

func (n *Sub) ExpressionNode() {}
func (n *Sub) Id() int         { return 0 }

type Mul struct {
	Left  Expression
	Right Expression
}

func (n *Mul) ExpressionNode() {}
func (n *Mul) Id() int         { return 0 }

type Div struct {
	Left  Expression
	Right Expression
}

func (n *Div) ExpressionNode() {}
func (n *Div) Id() int         { return 0 }

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

func (n *PrintStatement) StatementNode() {}
func (n *PrintStatement) Id() int        { return 0 }

type Statements struct {
	Statements []Statement
}

func (n *Statements) Id() int { return 0 }

type Name struct {
	Text string
}

func (n *Name) ExpressionNode() {}
func (n *Name) Id() int         { return 0 }

type CompoundExpression struct {
	Statements []Statement
}

func (n *CompoundExpression) Id() int { return 0 }

type ExpressionAsStatement struct {
	Expression Expression
}

func (n *ExpressionAsStatement) StatementNode() {}
func (n *ExpressionAsStatement) Id() int        { return 0 }

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

func (n *ConstDeclaration) StatementNode() {}
func (n *ConstDeclaration) Id() int        { return 0 }

type VarDeclaration struct {
	Name  Name
	Type  Type
	Value Expression
}

func (n *VarDeclaration) StatementNode() {}
func (n *VarDeclaration) Id() int        { return 0 }

type Assignment struct {
	Location Expression
	Value    Expression
}

func (n *Assignment) ExpressionNode() {}
func (n *Assignment) Id() int         { return 0 }

type Character struct {
	Value string
}

func (n *Character) ExpressionNode() {}
func (n *Character) Id() int         { return 0 }

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
	indent_str := context.Indent
	switch v := node.(type) {
	case *Integer:
		return v.Value
	case *Float:
		return v.Value
	case *Name:
		return v.Text
	case *NameType:
		return v.Name
	case *FloatType:
		return "float"
	case *BinOpWithOp:
		return fmt.Sprintf(" %s %s %s",
			NodeAsSource(v.Left, context), v.Op,
			NodeAsSource(v.Right, context))
	case *Add:
		return fmt.Sprintf(" %s %s %s",
			NodeAsSource(v.Left, context), "+",
			NodeAsSource(v.Right, context))
	case *Sub:
		return fmt.Sprintf(" %s %s %s",
			NodeAsSource(v.Left, context), "-",
			NodeAsSource(v.Right, context))
	case *Mul:
		return fmt.Sprintf(" %s %s %s",
			NodeAsSource(v.Left, context), "*",
			NodeAsSource(v.Right, context))
	case *Div:
		return fmt.Sprintf(" %s %s %s",
			NodeAsSource(v.Left, context), "/",
			NodeAsSource(v.Right, context))
	case *ConstDeclaration:
		if v.Type == nil {
			return fmt.Sprintf("%sconst %s = %s;",
				indent_str, NodeAsSource(&v.Name, context),
				NodeAsSource(v.Value, context))
		} else {
			return fmt.Sprintf("%sconst %s %s = %s;",
				indent_str, NodeAsSource(&v.Name, context),
				NodeAsSource(v.Type, context),
				NodeAsSource(v.Value, context))
		}
	case *VarDeclaration:
		if v.Value != nil {
			if v.Type == nil {
				return fmt.Sprintf("%sconst %s = %s;",
					indent_str, NodeAsSource(&v.Name, context),
					NodeAsSource(v.Value, context))
			} else {
				return fmt.Sprintf("%sconst %s %s = %s;",
					indent_str, NodeAsSource(&v.Name, context),
					NodeAsSource(v.Type, context),
					NodeAsSource(v.Value, context))
			}
		} else {
			return fmt.Sprintf("%sconst %s %s;",
				indent_str, NodeAsSource(&v.Name, context),
				NodeAsSource(v.Type, context))
		}
	case *Assignment: // Assignment is exp not statement
		return fmt.Sprintf("%s = %s",
			NodeAsSource(v.Location, context),
			NodeAsSource(v.Value, context))
	case *ExpressionAsStatement:
		return fmt.Sprintf("%s;", NodeAsSource(v.Expression, context))
	case *PrintStatement:
		return fmt.Sprintf("print %s;", NodeAsSource(v.Value, context))
	case *Statements:
		var ts []string
		for _, s := range v.Statements {
			ts = append(ts, NodeAsSource(s, context))
		}
		return strings.Join(ts, "\n")
	default:
		panic(fmt.Sprintf("Can't convert %v to source", v))
	}
}
