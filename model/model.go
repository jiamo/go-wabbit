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

//func NewNodeInfo(args map[string]int) *NodeInfo {
//	self := &NodeInfo{}
//	staticID++
//	self.Id = staticID
//	self.Indent += args["indent"]
//	return self
//}

type NodeManager struct {
	NodeMap map[Node]NodeInfo
}

var TheNodeManager NodeManager

func init() {
	// Ugly Now
	TheNodeManager = NodeManager{
		make(map[Node]NodeInfo),
	}
}

func RegisterNode(node Node) {
	staticID++
	n := NodeInfo{
		Id: staticID,
	}
	TheNodeManager.NodeMap[node] = n
}

func GetNodeInfo(node Node) NodeInfo {
	// if node in TheNodeManager return it
	// else put it with create a nodeinfo
	if _, ok := TheNodeManager.NodeMap[node]; ok {
		// return it
		return TheNodeManager.NodeMap[node]
	} else {
		// create a nodeinfo
		// put it
		staticID++
		n := NodeInfo{
			Id: staticID,
		}
		TheNodeManager.NodeMap[node] = n
		return n
	}
}

type Boolean interface {
	Node
	Bool() bool
}

type TrueBool struct {
}

func (n *TrueBool) Bool()   {}
func (n *TrueBool) Id() int { return GetNodeInfo(n).Id }

type FalseBool struct {
}

func (n *FalseBool) Bool()   {}
func (n *FalseBool) Id() int { return GetNodeInfo(n).Id }

type NameBool struct {
	Name string
}

func (n *NameBool) ExpressionNode() {}
func (n *NameBool) Bool()           {}
func (n *NameBool) Id() int         { return GetNodeInfo(n).Id }

type Float struct {
	Value float64
}

func (n *Float) ExpressionNode() {}
func (n *Float) Id() int         { return GetNodeInfo(n).Id }

type Integer struct {
	Value int
}

func (n *Integer) ExpressionNode() {}
func (n *Integer) Id() int         { return GetNodeInfo(n).Id }

// Type should be interface ?
type Type interface {
	Node
	Type() string
}

type NameType struct {
	Name string
}

func (n *NameType) Id() int      { return GetNodeInfo(n).Id }
func (n *NameType) Type() string { return n.Name }

type IntegerType struct {
}

func (n *IntegerType) Id() int      { return GetNodeInfo(n).Id }
func (n *IntegerType) Type() string { return "int" }

type FloatType struct {
}

func (n *FloatType) Id() int      { return GetNodeInfo(n).Id }
func (n *FloatType) Type() string { return "float" }

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
func (n *Neg) Id() int         { return GetNodeInfo(n).Id }

//func (n *Neg) String() int     { return GetNodeInfo(n).Id }

type Pos struct {
	Operand Expression
}

func (n *Pos) ExpressionNode() {}
func (n *Pos) Id() int         { return GetNodeInfo(n).Id }

type Not struct {
	Operand Expression
}

func (n *Not) ExpressionNode() {}
func (n *Not) Id() int         { return GetNodeInfo(n).Id }

type BinOpWithOp struct {
	Op    string
	Left  Expression
	Right Expression
}

func (n *BinOpWithOp) ExpressionNode() {}
func (n *BinOpWithOp) Id() int         { return GetNodeInfo(n).Id }

type BinOp struct {
	Left  Expression
	Right Expression
}

func (n *BinOp) ExpressionNode() {}
func (n *BinOp) Id() int         { return GetNodeInfo(n).Id }

type Add struct {
	Left  Expression
	Right Expression
}

func (n *Add) ExpressionNode() {}
func (n *Add) Id() int         { return GetNodeInfo(n).Id }

type Sub struct {
	Left  Expression
	Right Expression
}

func (n *Sub) ExpressionNode() {}
func (n *Sub) Id() int         { return GetNodeInfo(n).Id }

type Mul struct {
	Left  Expression
	Right Expression
}

func (n *Mul) ExpressionNode() {}
func (n *Mul) Id() int         { return GetNodeInfo(n).Id }

type Div struct {
	Left  Expression
	Right Expression
}

func (n *Div) ExpressionNode() {}
func (n *Div) Id() int         { return GetNodeInfo(n).Id }

type Lt struct {
	Left  Expression
	Right Expression
}

func (n *Lt) ExpressionNode() {}
func (n *Lt) Id() int         { return GetNodeInfo(n).Id }

type Le struct {
	Left  Expression
	Right Expression
}

func (n *Le) ExpressionNode() {}
func (n *Le) Id() int         { return GetNodeInfo(n).Id }

type Gt struct {
	Left  Expression
	Right Expression
}

func (n *Gt) ExpressionNode() {}
func (n *Gt) Id() int         { return GetNodeInfo(n).Id }

type Ge struct {
	Left  Expression
	Right Expression
}

func (n *Ge) ExpressionNode() {}
func (n *Ge) Id() int         { return GetNodeInfo(n).Id }

type Eq struct {
	Left  Expression
	Right Expression
}

func (n *Eq) ExpressionNode() {}
func (n *Eq) Id() int         { return GetNodeInfo(n).Id }

type Ne struct {
	Left  Expression
	Right Expression
}

func (n *Ne) ExpressionNode() {}
func (n *Ne) Id() int         { return GetNodeInfo(n).Id }

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
func (n *LogOr) Id() int         { return GetNodeInfo(n).Id }

type LogAnd struct {
	Left  Expression
	Right Expression
}

func (n *LogAnd) ExpressionNode() {}
func (n *LogAnd) Id() int         { return GetNodeInfo(n).Id }

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
func (n *PrintStatement) Id() int        { return GetNodeInfo(n).Id }

type Statements struct {
	Statements []Statement
}

func (n *Statements) Id() int { return GetNodeInfo(n).Id }

type Name struct {
	Text string
}

func (n *Name) ExpressionNode() {}
func (n *Name) Id() int         { return GetNodeInfo(n).Id }

type CompoundExpression struct {
	Statements Statements
}

func (n *CompoundExpression) ExpressionNode() {}
func (n *CompoundExpression) Id() int         { return GetNodeInfo(n).Id }

type ExpressionAsStatement struct {
	Expression Expression
}

func (n *ExpressionAsStatement) StatementNode() {}
func (n *ExpressionAsStatement) Id() int        { return GetNodeInfo(n).Id }

type Grouping struct {
	Expression Expression
}

func (n *Grouping) Id() int         { return GetNodeInfo(n).Id }
func (n *Grouping) ExpressionNode() {}

type Declaration interface {
	Statement
}

type ConstDeclaration struct {
	Name  Name
	Type  Type
	Value Expression
}

func (n *ConstDeclaration) StatementNode() {}
func (n *ConstDeclaration) Id() int        { return GetNodeInfo(n).Id }

type VarDeclaration struct {
	Name  Name
	Type  Type
	Value Expression
}

func (n *VarDeclaration) StatementNode() {}
func (n *VarDeclaration) Id() int        { return GetNodeInfo(n).Id }

type Assignment struct {
	Location Expression
	Value    Expression
}

func (n *Assignment) ExpressionNode() {}
func (n *Assignment) Id() int         { return GetNodeInfo(n).Id }

type Character struct {
	Value string
}

func (n *Character) ExpressionNode() {}
func (n *Character) Id() int         { return GetNodeInfo(n).Id }

type IfStatement struct {
	Test        Expression
	Consequence Statements
	Alternative *Statements
}

func (n *IfStatement) StatementNode() {}
func (n *IfStatement) Id() int        { return GetNodeInfo(n).Id }

type WhileStatement struct {
	Test Expression
	Body Statements
}

func (n *WhileStatement) StatementNode() {}
func (n *WhileStatement) Id() int        { return GetNodeInfo(n).Id }

type BreakStatement struct {
}

func (n *BreakStatement) StatementNode() {}
func (n *BreakStatement) Id() int        { return GetNodeInfo(n).Id }

type ContinueStatement struct {
}

func (n *ContinueStatement) StatementNode() {}
func (n *ContinueStatement) Id() int        { return GetNodeInfo(n).Id }

type FunctionDeclaration struct {
	Name       Name
	Parameters []Parameter
	ReturnType Type
	Body       Statements
}

func (n *FunctionDeclaration) StatementNode() {}
func (n *FunctionDeclaration) Id() int        { return GetNodeInfo(n).Id }

type FunctionApplication struct {
	Func      Expression
	Arguments []Expression
}

func (n *FunctionApplication) ExpressionNode() {}
func (n *FunctionApplication) Id() int         { return GetNodeInfo(n).Id }

type ReturnStatement struct {
	Value Expression
}

func (n *ReturnStatement) StatementNode() {}
func (n *ReturnStatement) Id() int        { return GetNodeInfo(n).Id }

type Parameter struct {
	Name Name
	Type Type
}

func (n *Parameter) Id() int { return GetNodeInfo(n).Id }

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
		return fmt.Sprintf("%v", v.Value)
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
	case *FunctionApplication:
		var ts []string
		for _, p := range v.Arguments {
			ts = append(ts, fmt.Sprintf("%s", NodeAsSource(p, context)))
		}
		return fmt.Sprintf("%s(%s)", NodeAsSource(v.Func, context), strings.Join(ts, " "))
	case *VarDeclaration:
		if v.Value != nil {
			if v.Type == nil {
				return fmt.Sprintf("%svar %s = %s;",
					indent_str, NodeAsSource(&v.Name, context),
					NodeAsSource(v.Value, context))
			} else {
				return fmt.Sprintf("%svar %s %s = %s;",
					indent_str, NodeAsSource(&v.Name, context),
					NodeAsSource(v.Type, context),
					NodeAsSource(v.Value, context))
			}
		} else {
			return fmt.Sprintf("%sconst %s %s;",
				indent_str, NodeAsSource(&v.Name, context),
				NodeAsSource(v.Type, context))
		}
	case *CompoundExpression:
		var ts []string
		for _, s := range v.Statements.Statements {
			ts = append(ts, fmt.Sprintf("%s%s", indent_str, NodeAsSource(s, context)))
		}
		return fmt.Sprintf("{ %s }", strings.Join(ts, " "))
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
	case *ReturnStatement:
		return fmt.Sprintf("return %s;", NodeAsSource(v.Value, context))
	case *FunctionDeclaration:
		var ts []string
		for _, p := range v.Parameters {
			ts = append(ts, fmt.Sprintf("%s %s", NodeAsSource(&p.Name, context), NodeAsSource(p.Type, context)))
		}

		return fmt.Sprintf("func %s(%s) %s{\n%s\n}",
			NodeAsSource(&v.Name, context), strings.Join(ts, ", "), NodeAsSource(v.ReturnType, context),
			NodeAsSource(&v.Body, context.NewBlock()))
	default:
		panic(fmt.Sprintf("Can't convert %v to source", v))
	}

}
