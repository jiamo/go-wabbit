package interpreter

import (
	"errors"
	"fmt"
	"strconv"
	"wabbit-go/common"
	"wabbit-go/model"
)

type Context struct {
	env   *common.ChainMap
	level int
}

type WabbitValue struct {
	Type  string
	Value interface{}
}

type FunctionClosure struct {
	Node    *model.FunctionDeclaration
	Context *Context
}

func (w *WabbitValue) String() string {
	return fmt.Sprintf("%v", w.Value)
}

type WabbitVar struct {
	Kind  string
	Value *WabbitValue
}

func NewContext(env *common.ChainMap, level int) *Context {
	if env == nil {
		env = common.NewChainMap()
	}
	return &Context{env: env, level: level}
}

func (c *Context) Define(name string, value *WabbitValue) error {
	if _, ok := c.env.GetValue(name); ok {
		return errors.New(fmt.Sprintf("%s already defined", name))
	}
	c.env.SetValue(name, value)
	return nil
}

func (c *Context) Assign(name string, val *WabbitValue) error {
	value, _ := c.Lookup(name)
	_ = value.Value.(*WabbitVar).Store(val) // if value store val it will return error
	return nil
}

func (c *Context) Lookup(name string) (*WabbitValue, bool) {
	value, ok := c.env.GetValue(name)
	if !ok {
		return nil, ok
	}
	return value.(*WabbitValue), ok
}

func (c *Context) NewBlock() *Context {
	return NewContext(c.env.NewChild(), c.level+1)
}

func NewWabbitVar(kind string, value *WabbitValue) WabbitVar {
	return WabbitVar{Kind: kind, Value: value}
}

func (w *WabbitVar) Load() *WabbitValue {
	return w.Value
}

func (w *WabbitVar) Store(value *WabbitValue) error {
	if w.Kind == "const" {
		return errors.New("Can't store const")
	}
	if value.Type != w.Value.Type {
		return errors.New("Type error in assignment")
	}
	w.Value = value
	return nil
}

func InterpretProgram(program model.Node) interface{} {
	// 创建一个新的上下文环境实例，用于存储和管理变量
	context := NewContext(nil, 0)

	// 在上下文环境中预定义基本类型
	_ = context.Define("int", &WabbitValue{Type: "cast", Value: "int"})
	_ = context.Define("float", &WabbitValue{Type: "cast", Value: "float"})
	_ = context.Define("char", &WabbitValue{Type: "cast", Value: "char"})
	_ = context.Define("bool", &WabbitValue{Type: "cast", Value: "bool"})

	// 解释输入的程序模型，并返回解释结果
	return InterpretNode(program, context)
}

func InterpretNode(node model.Node, context *Context) *WabbitValue {
	switch v := node.(type) {
	case *model.Integer:
		return &WabbitValue{"int", v.Value}
	case *model.Float:
		return &WabbitValue{"float", v.Value}
	case *model.Character:
		unquoted, err := strconv.Unquote(v.Value)
		if err != nil {
			panic(err)
		}
		return &WabbitValue{"char", rune(unquoted[0])}
	case *model.Name:
		value, _ := context.Lookup(v.Text) // somethings we may need exist
		// check value.Value is WabbitVar if success then return  its load
		// else return value
		value_var, ok := value.Value.(*WabbitVar)
		//fmt.Println("hello world", "value ", value, " ", value_var, " ", v.Text)
		if ok {
			return value_var.Load()
		} else {
			return value
		}
	case *model.NameType:
		return &WabbitValue{Type: "type", Value: v.Name}
	case *model.NameBool:
		return &WabbitValue{Type: "bool", Value: v.Name == "true"}
	case *model.IntegerType:
		return &WabbitValue{Type: "type", Value: "int"}
	case *model.FloatType:
		return &WabbitValue{Type: "type", Value: "float"}

	case *model.Add:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left.Type == "int" && right.Type == "int" {
			return &WabbitValue{"int", left.Value.(int) + right.Value.(int)}
		} else if left.Type == "float" && right.Type == "float" {
			return &WabbitValue{"float", left.Value.(float64) + right.Value.(float64)}
		} else {
			// we think it's a type error
			return &WabbitValue{"error", "type error"}
		}
	case *model.Mul:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left.Type == "int" && right.Type == "int" {
			return &WabbitValue{"int", left.Value.(int) * right.Value.(int)}
		} else if left.Type == "float" && right.Type == "float" {
			return &WabbitValue{"float", left.Value.(float64) * right.Value.(float64)}
		} else {
			// we think it's a type error
			return &WabbitValue{"error", "type error"}
		}
	case *model.Sub:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left.Type == "int" && right.Type == "int" {
			return &WabbitValue{"int", left.Value.(int) - right.Value.(int)}
		} else if left.Type == "float" && right.Type == "float" {
			return &WabbitValue{"float", left.Value.(float64) - right.Value.(float64)}
		} else {
			// we think it's a type error
			return &WabbitValue{"error", "type error"}
		}
	case *model.Div:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left.Type == "int" && right.Type == "int" {
			return &WabbitValue{"int", left.Value.(int) / right.Value.(int)}
		} else if left.Type == "float" && right.Type == "float" {
			return &WabbitValue{"float", left.Value.(float64) / right.Value.(float64)}
		} else {
			// we think it's a type error
			return &WabbitValue{"error", "type error"}
		}

	case *model.Neg:
		right := InterpretNode(v.Operand, context)
		if right.Type == "int" {
			return &WabbitValue{"int", -right.Value.(int)}
		} else if right.Type == "float" {
			return &WabbitValue{"float", -right.Value.(float64)}
		} else {
			// we think it's a type error
			return &WabbitValue{"error", "type error"}
		}
	case *model.Pos:
		right := InterpretNode(v.Operand, context)
		if right.Type == "int" {
			return &WabbitValue{"int", +right.Value.(int)}
		} else if right.Type == "float" {
			return &WabbitValue{"float", +right.Value.(float64)}
		} else {
			// we think it's a type error
			return &WabbitValue{"error", "type error"}
		}
	case *model.Not:
		right := InterpretNode(v.Operand, context)
		if right.Type == "bool" {
			return &WabbitValue{"int", !right.Value.(bool)}
		} else {
			// we think it's a type error
			return &WabbitValue{"error", "type error"}
		}
	case *model.VarDeclaration:
		var val *WabbitValue
		if v.Value != nil {
			val = InterpretNode(v.Value, context)
		} else {
			val = &WabbitValue{Type: v.Type.Type(), Value: nil}
		}
		context.Define(v.Name.Text, &WabbitValue{Type: val.Type, Value: &WabbitVar{"var", val}})
	case *model.ConstDeclaration:
		var val *WabbitValue
		if v.Value != nil {
			val = InterpretNode(v.Value, context)
		} else {
			val = &WabbitValue{Type: v.Type.Type(), Value: nil}
		}
		context.Define(v.Name.Text, &WabbitValue{Type: val.Type, Value: &WabbitVar{"constant", val}})

	case *model.Lt:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		if left.Type == "int" {
			return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
		} else if left.Type == "float" {
			return &WabbitValue{Type: "bool", Value: left.Value.(float64) < right.Value.(float64)}
		} else {
			return &WabbitValue{Type: "error", Value: "type error"}
		}
	case *model.Le:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		if left.Type == "int" {
			return &WabbitValue{Type: "bool", Value: left.Value.(int) <= right.Value.(int)}
		} else if left.Type == "float" {
			return &WabbitValue{Type: "bool", Value: left.Value.(float64) <= right.Value.(float64)}
		} else {
			return &WabbitValue{Type: "error", Value: "type error"}
		}
	case *model.Gt:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		if left.Type == "int" {
			return &WabbitValue{Type: "bool", Value: left.Value.(int) > right.Value.(int)}
		} else if left.Type == "float" {
			return &WabbitValue{Type: "bool", Value: left.Value.(float64) > right.Value.(float64)}
		} else {
			return &WabbitValue{Type: "error", Value: "type error"}
		}
	case *model.Ge:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		if left.Type == "int" {
			return &WabbitValue{Type: "bool", Value: left.Value.(int) >= right.Value.(int)}
		} else if left.Type == "float" {
			return &WabbitValue{Type: "bool", Value: left.Value.(float64) >= right.Value.(float64)}
		} else {
			return &WabbitValue{Type: "error", Value: "type error"}
		}
	case *model.Eq:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		if left.Type == "int" {
			return &WabbitValue{Type: "bool", Value: left.Value.(int) == right.Value.(int)}
		} else if left.Type == "float" {
			return &WabbitValue{Type: "bool", Value: left.Value.(float64) == right.Value.(float64)}
		} else if left.Type == "bool" {
			return &WabbitValue{Type: "bool", Value: left.Value.(bool) == right.Value.(bool)}
		} else {
			return &WabbitValue{Type: "error", Value: "type error"}
		}
	case *model.Ne:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		if left.Type == "int" {
			return &WabbitValue{Type: "bool", Value: left.Value.(int) != right.Value.(int)}
		} else if left.Type == "float" {
			return &WabbitValue{Type: "bool", Value: left.Value.(float64) != right.Value.(float64)}
		} else if left.Type == "bool" {
			return &WabbitValue{Type: "bool", Value: left.Value.(bool) != right.Value.(bool)}
		} else {
			return &WabbitValue{Type: "error", Value: "type error"}
		}
	case *model.LogOr:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		return &WabbitValue{Type: "bool", Value: left.Value.(bool) || right.Value.(bool)}
	case *model.LogAnd:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		return &WabbitValue{Type: "bool", Value: left.Value.(bool) && right.Value.(bool)}
	case *model.Assignment:
		val := InterpretNode(v.Value, context)
		// assign the value to the name
		_ = context.Assign(v.Location.(*model.Name).Text, val)
		// print (c=4) + 4
		// just return the val ?
		return val
	case *model.PrintStatement:
		value := InterpretNode(v.Value, context)
		switch value.Type {
		case "char":
			fmt.Printf("%c", value.Value.(rune)) // we may need change
		case "bool":
			if value.Value.(bool) {
				fmt.Println("true")
			} else {
				fmt.Println("false")
			}
		default:
			fmt.Println(value)
		}
	case *model.Statements:

		var result *WabbitValue
		for _, statement := range v.Statements {
			result = InterpretNode(statement, context)

			// need check break return too
			if result != nil {
				if result.Type == "break" || result.Type == "return" || result.Type == "continue" {
					return result
				}
			}
		}
		return result

	case *model.ExpressionAsStatement:
		InterpretNode(v.Expression, context)

	case *model.Grouping:
		return InterpretNode(v.Expression, context)

	case *model.IfStatement:
		condition := InterpretNode(v.Test, context)
		if condition.Value.(bool) {
			return InterpretNode(&v.Consequence, context)
		} else if v.Alternative != nil {
			return InterpretNode(v.Alternative, context)
		}

	case *model.BreakStatement:
		return &WabbitValue{"break", nil}
	case *model.ContinueStatement:
		return &WabbitValue{"continue", nil}
	case *model.ReturnStatement:
		value := InterpretNode(v.Value, context)
		return &WabbitValue{"return", value}

	case *model.WhileStatement:
		for true {
			condtion := InterpretNode(v.Test, context)
			if condtion.Type == "error" { // fix it
				return condtion
			}
			if condtion.Value.(bool) {
				result := InterpretNode(&v.Body, context)
				//fmt.Println("result %v", result)
				if result == nil {
					continue
				}
				if result.Type == "break" {
					break
				} else if result.Type == "continue" {
					continue
				} else if result.Type == "return" {
					return result
				}
				continue
			} else {

				break
			}
		}
	case *model.FunctionDeclaration:
		// we should check the function name is not defined
		context.Define(v.Name.Text, &WabbitValue{"func", &FunctionClosure{v, context}})

	case *model.FunctionApplication:
		value := InterpretNode(v.Func, context) // while lookup

		//savedContext = funtionClosure.Value.Context
		// TODO make it as builtin function...
		if value.Type == "cast" { // this is conversion function
			if value.Value == "int" {
				result := InterpretNode(v.Arguments[0], context) // # type conversion must only have one args ?
				if result.Type == "float" {
					return &WabbitValue{"int", int(result.Value.(float64))}
				}
				return result
			}
		}
		// custom function and it should be....
		if value.Type == "func" {
			funtionClosure := value.Value.(*FunctionClosure)
			eval_context := funtionClosure.Context.NewBlock()
			for i, arg := range v.Arguments {
				// TODO check the type
				//fmt.Println("arg %v", arg)
				argValue := InterpretNode(arg, context) // arg eval in current context
				eval_context.Define(
					funtionClosure.Node.Parameters[i].Name.Text,
					&WabbitValue{
						funtionClosure.Node.Parameters[i].Type.Type(),
						&WabbitVar{"var", argValue}})
			}

			bodyValue := InterpretNode(&funtionClosure.Node.Body, eval_context)
			if bodyValue != nil {
				if bodyValue.Type == "return" {
					return bodyValue.Value.(*WabbitValue)
				}
				return bodyValue
			}
		}

	default:
		panic(fmt.Sprintf("Can't intepre %#v to source", v))
	}

	return nil
}
