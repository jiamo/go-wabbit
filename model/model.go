package model

import (
	"fmt"
	"strconv"
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
	Bool() bool
}

type TrueBool struct {
}

func (n *TrueBool) Bool()   {}
func (n *TrueBool) Id() int { return 0 }

type FalseBool struct {
}

func (n *FalseBool) Bool()   {}
func (n *FalseBool) Id() int { return 0 }

type NameBool struct {
	Name string
}

func (n *NameBool) ExpressionNode() {}
func (n *NameBool) Bool()           {}
func (n *NameBool) Id() int         { return 0 }

type Float struct {
	Value string
}

func (n *Float) ExpressionNode() {}
func (n *Float) Id() int         { return 0 }

type Integer struct {
	Value int
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

func (n *IntegerType) Id() int      { return 0 }
func (n *IntegerType) Type() string { return "" }

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

type Lt struct {
	Left  Expression
	Right Expression
}

func (n *Lt) ExpressionNode() {}
func (n *Lt) Id() int         { return 0 }

type Le struct {
	Left  Expression
	Right Expression
}

func (n *Le) ExpressionNode() {}
func (n *Le) Id() int         { return 0 }

type Gt struct {
	Left  Expression
	Right Expression
}

func (n *Gt) ExpressionNode() {}
func (n *Gt) Id() int         { return 0 }

type Ge struct {
	Left  Expression
	Right Expression
}

func (n *Ge) ExpressionNode() {}
func (n *Ge) Id() int         { return 0 }

type Eq struct {
	Left  Expression
	Right Expression
}

func (n *Eq) ExpressionNode() {}
func (n *Eq) Id() int         { return 0 }

type Ne struct {
	Left  Expression
	Right Expression
}

func (n *Ne) ExpressionNode() {}
func (n *Ne) Id() int         { return 0 }

type RelOp struct {
	Left  Expression
	Right Expression
	Chain *RelOp
}

type LogicalOp struct {
	Left  Expression
	Right Expression
}

type LogOr struct {
	Left  Expression
	Right Expression
}

func (n *LogOr) ExpressionNode() {}
func (n *LogOr) Id() int         { return 0 }

type LogAnd struct {
	Left  Expression
	Right Expression
}

func (n *LogAnd) ExpressionNode() {}
func (n *LogAnd) Id() int         { return 0 }

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

func (g *Grouping) Id() int         { return 0 }
func (g *Grouping) ExpressionNode() {}

type Declaration struct {
}

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

func (n *IfStatement) StatementNode() {}
func (n *IfStatement) Id() int        { return 0 }

type WhileStatement struct {
	Test Expression
	Body Statements
}

func (n *WhileStatement) StatementNode() {}
func (n *WhileStatement) Id() int        { return 0 }

type BreakStatement struct {
}

func (n *BreakStatement) StatementNode() {}
func (n *BreakStatement) Id() int        { return 0 }

type ContinueStatement struct {
}

func (n *ContinueStatement) StatementNode() {}
func (n *ContinueStatement) Id() int        { return 0 }

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
		return strconv.Itoa(v.Value)
	case *Float:
		return v.Value
	case *Name:
		return v.Text
	case *NameType:
		return v.Name
	case *NameBool:
		return v.Name
	case *IntegerType:
		return "int"
	case *FloatType:
		return "float"
	case *BinOpWithOp:
		return fmt.Sprintf(" %s %s %s",
			NodeAsSource(v.Left, context), v.Op,
			NodeAsSource(v.Right, context))
	case *Add:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "+",
			NodeAsSource(v.Right, context))
	case *Sub:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "-",
			NodeAsSource(v.Right, context))
	case *Mul:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "*",
			NodeAsSource(v.Right, context))
	case *Div:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "/",
			NodeAsSource(v.Right, context))
	case *Lt:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "<",
			NodeAsSource(v.Right, context))
	case *Le:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "<=",
			NodeAsSource(v.Right, context))
	case *Gt:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), ">",
			NodeAsSource(v.Right, context))
	case *Ge:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), ">=",
			NodeAsSource(v.Right, context))
	case *Eq:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "==",
			NodeAsSource(v.Right, context))
	case *Ne:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "!=",
			NodeAsSource(v.Right, context))
	case *Neg:
		return fmt.Sprintf("-%s",
			NodeAsSource(v.Operand, context))

	case *LogOr:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "or",
			NodeAsSource(v.Right, context))
	case *LogAnd:
		return fmt.Sprintf("%s %s %s",
			NodeAsSource(v.Left, context), "and",
			NodeAsSource(v.Right, context))
	case *Grouping:
		return fmt.Sprintf("(%s)", NodeAsSource(v.Expression, context))
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
			ts = append(ts, fmt.Sprintf("%s%s", indent_str, NodeAsSource(s, context)))
		}
		return strings.Join(ts, "\n")
	case *BreakStatement:
		return "break;"
	case *ContinueStatement:
		return "continue;"
	case *IfStatement:
		var ts []string
		ts = append(ts, fmt.Sprintf("if %s {", NodeAsSource(v.Test, context)))
		ts = append(ts, fmt.Sprintf("%s", NodeAsSource(&v.Consequence, context.NewBlock())))

		if v.Alternative != nil {
			ts = append(ts, fmt.Sprintf("%s} else {", context.Indent))
			ts = append(ts, fmt.Sprintf("%s", NodeAsSource(v.Alternative, context.NewBlock())))
		} else {

		}
		ts = append(ts, fmt.Sprintf("%s}", context.Indent))
		return strings.Join(ts, "\n")
	case *WhileStatement:
		var ts []string
		ts = append(ts, fmt.Sprintf("while %s {", NodeAsSource(v.Test, context)))
		ts = append(ts, NodeAsSource(&v.Body, context.NewBlock()))
		ts = append(ts, "}")
		return strings.Join(ts, "\n")
	default:
		panic(fmt.Sprintf("Can't convert %v to source", v))
	}
}
