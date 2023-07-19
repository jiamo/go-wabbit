package wasm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"wabbit-go/common"
	"wabbit-go/model"
)

var _typemap = map[string]string{
	"int":   "i32",
	"float": "f64",
	"bool":  "i32",
	"char":  "i32",
}

type Parameter struct {
	Name string
	Type string
}

type Function struct {
	name       string
	parameters []model.Parameter
	retType    string
	code       []string
	locals     []string
}

func (f *Function) String() string {
	out := fmt.Sprintf("(func $%s (export \"%s\")\n", f.name, f.name)
	for _, parm := range f.parameters {
		out += fmt.Sprintf("(param $%s %s)\n", parm.Name.Text, _typemap[parm.Type.Type()])
	}
	if f.retType != "" {
		out += fmt.Sprintf("(result %s)\n", _typemap[f.retType])
		out += fmt.Sprintf("(local $return %s)\n", _typemap[f.retType])
	}
	out += "\n" + strings.Join(f.locals, "\n")
	out += "\nblock $return\n"
	out += "\n" + strings.Join(f.code, "\n")
	out += "\nend\n"
	if f.retType != "" {
		out += "local.get $return\n"
	}
	out += ")\n"
	return out
}

type Context struct {
	module   []string
	env      *common.ChainMap
	function Function
	scope    string
	nlabels  int
	haveMain bool
}

type WASMVar struct {
	Type  string
	Scope string
}

func NewWabbitWasmModule() *Context {
	w := &Context{
		module: []string{"(module"},
		env:    common.NewChainMap(),
		function: Function{
			name: "main",
		},
		scope: "global",
	}
	w.module = append(w.module, "(import \"env\" \"_printi\" (func $_printi ( param i32 )))")
	w.module = append(w.module, "(import \"env\" \"_printf\" (func $_printf ( param f64 )))")
	w.module = append(w.module, "(import \"env\" \"_printb\" (func $_printb ( param i32 )))")
	w.module = append(w.module, "(import \"env\" \"_printc\" (func $_printc ( param i32 )))")

	return w
}

func (m *Context) String() string {
	return strings.Join(m.module, "\n") + "\n)\n"
}

func (ctx *Context) Define(name string, value *WASMVar) {
	ctx.env.SetValue(name, value)
}

func (ctx *Context) Lookup(name string) *WASMVar {
	v, e := ctx.env.GetValue(name)
	if e == true {
		return v.(*WASMVar)
	} else {
		return nil
	}
}

func (ctx *Context) NewScope(do func()) {
	oldEnv := ctx.env
	ctx.env = ctx.env.NewChild()
	defer func() {
		ctx.env = oldEnv
	}()
	do()
}

func (m *Context) NewLabel(name ...string) string {
	newID := m.nlabels
	m.nlabels++
	if name == nil {
		return fmt.Sprintf(".%d", newID)
	} else {
		ls := strings.Split(name[0], ".")
		return ls[0] + "." + strconv.Itoa(newID)
	}

}

func Wasm(program *model.Program) string {
	wctx := NewWabbitWasmModule()
	_ = InterpretNode(program.Model, wctx) // generate is InterpretNode in the same meaning
	// where
	wctx.module = append(wctx.module, wctx.function.String())
	return wctx.String()
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func insert(slice []string, value string, position int) []string {
	slice = append(slice[:position], append([]string{value}, slice[position:]...)...)
	return slice
}

func InterpretNode(node model.Node, context *Context) string {
	switch v := node.(type) {
	case *model.Integer:
		context.function.code = append(context.function.code, fmt.Sprintf("i32.const %v", v.Value))
		return "int"
	case *model.Float:
		context.function.code = append(context.function.code, fmt.Sprintf("f64.const %v", v.Value))
		return "float"
	case *model.Character:
		unquoted, err := strconv.Unquote(v.Value)
		if err != nil {
			panic(err)
		}
		//context.code = append(context.code, Instruction{"IPUSH", int(rune(unquoted[0]))})
		context.function.code = append(context.function.code, fmt.Sprintf("i32.const %v", int(rune(unquoted[0]))))
		return "char"
	case *model.Name:
		value := context.Lookup(v.Text)
		if value.Scope == "global" {
			context.function.code = append(context.function.code, fmt.Sprintf("global.get $%s", v.Text))
		} else if value.Scope == "local" {
			context.function.code = append(context.function.code, fmt.Sprintf("local.get $%s", v.Text))
		}
		return value.Type

	//case *model.NameType:
	//	return v.Name
	case *model.NameBool:
		context.function.code = append(context.function.code, fmt.Sprintf("i32.const %v", BoolToInt(v.Name == "true")))
		return "bool"
	//case *model.IntegerType:
	//	return "int"
	//return &WabbitValue{Type: "type", Value: "int"}
	//case *model.FloatType:
	//	//return &WabbitValue{Type: "type", Value: "float"}
	//	return "float"
	case *model.Add:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left == "int" && right == "int" {
			context.function.code = append(context.function.code, "i32.add")
			return "int"
			//return &WabbitValue{"int", left.Value.(int) + right.Value.(int)}
		} else if left == "float" && right == "float" {
			//return &WabbitValue{"float", left.Value.(float64) + right.Value.(float64)}
			context.function.code = append(context.function.code, "f64.add")
			return "float"
		} else {
			// we think it's a type error
			panic("type different")
			//return &WabbitValue{"error", "type error"}
		}
		return left
	case *model.Mul:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left == "int" && right == "int" {
			context.function.code = append(context.function.code, "i32.mul")
			return "int"
			//return &WabbitValue{"int", left.Value.(int) * right.Value.(int)}
		} else if left == "float" && right == "float" {
			context.function.code = append(context.function.code, "f64.mul")
			return "float"
		} else {
			// we think it's a type error
			panic("type different")
		}
		return left
	case *model.Sub:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left == "int" && right == "int" {
			context.function.code = append(context.function.code, "i32.sub")
			return "int"
			//return &WabbitValue{"int", left.Value.(int) - right.Value.(int)}
		} else if left == "float" && right == "float" {
			context.function.code = append(context.function.code, "f64.sub")
			return "float"
		} else {
			// we think it's a type error
			panic("type different")
		}
		return left
	case *model.Div:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left == "int" && right == "int" {
			context.function.code = append(context.function.code, "i32.div_s")
			return "int"
		} else if left == "float" && right == "float" {
			context.function.code = append(context.function.code, "f64.div")
			return "float"
		} else {
			// we think it's a type error
			panic("type different")
		}
		return left

	case *model.Neg:
		pos := len(context.function.code)
		right := InterpretNode(v.Operand, context)
		if right == "int" {
			context.function.code = insert(context.function.code, "i32.const 0", pos)
			context.function.code = append(context.function.code, "i32.sub")
		} else if right == "float" {
			//return &WabbitValue{"float", -right.Value.(float64)}
			context.function.code = insert(context.function.code, "f64.const 0", pos)
			context.function.code = append(context.function.code, "f64.sub")
		} else {
			// we think it's a type error
			//return &WabbitValue{"error", "type error"}
			panic("type different")
		}
		return right
	case *model.Pos:
		right := InterpretNode(v.Operand, context)
		return right
	case *model.Not:
		right := InterpretNode(v.Operand, context)
		if right == "bool" {
			context.function.code = append(context.function.code, "i32.const 1")
			context.function.code = append(context.function.code, "i32.xor")
		} else {
			// we think it's a type error
			panic("type different")
		}
		return right
	case *model.VarDeclaration:
		//var val *WVMVar
		var valtype string
		if v.Value != nil {
			valtype = InterpretNode(v.Value, context) // store in stack
		} else {
			valtype = v.Type.Type()
		}

		if context.scope == "global" {
			if valtype == "float" {
				// global using module
				context.module = append(context.module, fmt.Sprintf("(global $%s (mut f64) (f64.const 0.0))", v.Name.Text))
			} else {
				context.module = append(context.module, fmt.Sprintf("(global $%s (mut i32) (i32.const 0))", v.Name.Text))
			}
			if v.Value != nil {
				context.function.code = append(context.function.code, fmt.Sprintf("global.set $%s", v.Name.Text))
			}
		} else if context.scope == "local" {
			// local using function
			if valtype == "float" {
				context.function.locals = append(context.function.locals, fmt.Sprintf("(local $%s f64)", v.Name.Text))
			} else {
				context.function.locals = append(context.function.locals, fmt.Sprintf("(local $%s i32)", v.Name.Text))
			}
			if v.Value != nil {
				context.function.code = append(context.function.code, fmt.Sprintf("local.set $%s", v.Name.Text))
			}
		}
		context.Define(v.Name.Text, &WASMVar{Type: valtype, Scope: context.scope})
		return ""

	case *model.ConstDeclaration:
		valtype := InterpretNode(v.Value, context)
		if context.scope == "global" {
			if valtype == "float" {
				context.module = append(context.module, fmt.Sprintf("(global $%s (mut f64) (f64.const 0.0))", v.Name.Text))
			} else {
				context.module = append(context.module, fmt.Sprintf("(global $%s (mut i32) (i32.const 0))", v.Name.Text))
			}
			context.function.code = append(context.function.code, fmt.Sprintf("global.set $%s", v.Name.Text))
		} else if context.scope == "local" {
			if valtype == "float" {
				context.function.locals = append(context.function.locals, fmt.Sprintf("(local $%s (mut f64) (f64.const 0.0))", v.Name.Text))
			} else {
				context.function.locals = append(context.function.locals, fmt.Sprintf("(local $%s (mut i32) (i32.const 0))", v.Name.Text))
			}
			context.function.code = append(context.function.code, fmt.Sprintf("local.set $%s", v.Name.Text))
		}
		context.Define(v.Name.Text, &WASMVar{Type: valtype, Scope: context.scope})
		return ""

	case *model.Lt:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "float" {
			context.function.code = append(context.function.code, "f64.lt")
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, "i32.lt_s")
		}
		return "bool"
	case *model.Le:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "float" {
			context.function.code = append(context.function.code, "f64.le")
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, "i32.le_s")
		}
		return "bool"
	case *model.Gt:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "float" {
			context.function.code = append(context.function.code, "f64.gt")
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, "i32.gt_s")
		}
		return "bool"
	case *model.Ge:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "float" {
			context.function.code = append(context.function.code, "f64.ge")
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, "i32.ge_s")
		}
		return "bool"
	case *model.Eq:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "float" {
			context.function.code = append(context.function.code, "f64.eq")
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, "i32.eq")
		}
		return "bool"
	case *model.Ne:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "float" {
			context.function.code = append(context.function.code, "f64.ne")
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, "i32.ne")
		}
		return "bool"
	case *model.LogOr:
		// TODO short eval
		begin := context.NewLabel("begin")
		context.function.code = append(context.function.code, fmt.Sprintf("block $%s (result i32)", begin))

		or_block := context.NewLabel("or_block")
		context.function.code = append(context.function.code, fmt.Sprintf("block $%s", or_block))
		_ = InterpretNode(v.Left, context)
		context.function.code = append(context.function.code, fmt.Sprintf("br_if $%s", or_block))
		_ = InterpretNode(v.Right, context)
		context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		context.function.code = append(context.function.code, "end")
		context.function.code = append(context.function.code, "i32.const 1")
		context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		context.function.code = append(context.function.code, "end")
		return "bool"

	case *model.LogAnd:
		begin := context.NewLabel("begin")
		context.function.code = append(context.function.code, fmt.Sprintf("block $%s (result i32)", begin))

		and_block := context.NewLabel("and_block")
		context.function.code = append(context.function.code, fmt.Sprintf("block $%s", and_block))
		_ = InterpretNode(v.Left, context)
		context.function.code = append(context.function.code, "i32.const 1")
		context.function.code = append(context.function.code, "i32.xor")
		context.function.code = append(context.function.code, fmt.Sprintf("br_if $%s", and_block))
		_ = InterpretNode(v.Right, context)
		context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		context.function.code = append(context.function.code, "end")
		context.function.code = append(context.function.code, "i32.const 0")
		context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		context.function.code = append(context.function.code, "end")
		return "bool" // no need or and any more

	case *model.Assignment:
		val := InterpretNode(v.Value, context)
		// assign the value to the name
		wasmvar := context.Lookup(v.Location.(*model.Name).Text)
		if wasmvar.Scope == "global" {
			context.function.code = append(context.function.code, fmt.Sprintf("global.set $%s", v.Location.(*model.Name).Text))
			context.function.code = append(context.function.code, fmt.Sprintf("global.get $%s", v.Location.(*model.Name).Text))
		} else {
			context.function.code = append(context.function.code, fmt.Sprintf("local.tee $%s", v.Location.(*model.Name).Text))
		}
		return val

	case *model.PrintStatement:
		value := InterpretNode(v.Value, context)
		switch value {
		case "char":
			context.function.code = append(context.function.code, "call $_printc")
		case "bool":
			context.function.code = append(context.function.code, "call $_printb")
		case "int":
			context.function.code = append(context.function.code, "call $_printi")
		case "float":
			context.function.code = append(context.function.code, "call $_printf")
		default:
			panic("wrong type")
		}
	case *model.Statements:

		var result string
		for _, statement := range v.Statements {
			if result != "" {
				context.function.code = append(context.function.code, "drop")
			}
			result = InterpretNode(statement, context)
		}
		return result

	case *model.ExpressionAsStatement:
		InterpretNode(v.Expression, context)
		context.function.code = append(context.function.code, "drop")

	case *model.Grouping:
		return InterpretNode(v.Expression, context)

	case *model.IfStatement:

		InterpretNode(v.Test, context)
		context.function.code = append(context.function.code, "if")

		context.NewScope(
			func() {
				InterpretNode(&v.Consequence, context)
			},
		)
		if v.Alternative != nil {
			context.function.code = append(context.function.code, "else")
			context.NewScope(
				func() {
					InterpretNode(v.Alternative, context)
				},
			)
		}
		context.function.code = append(context.function.code, "end")

	case *model.BreakStatement:
		// we need scope for level break
		val := context.Lookup("break") // fake using type as label
		context.function.code = append(context.function.code, fmt.Sprintf("br $%s", val.Scope))
	case *model.ContinueStatement:
		val := context.Lookup("continue") // fake using type as label
		context.function.code = append(context.function.code, fmt.Sprintf("br $%s", val.Scope))

	case *model.ReturnStatement:
		value := InterpretNode(v.Value, context)
		context.function.code = append(context.function.code, "return")
		return value

	case *model.WhileStatement:
		test_label := context.NewLabel()
		exit_label := context.NewLabel()

		context.function.code = append(context.function.code, fmt.Sprintf("block $%s", exit_label))
		context.function.code = append(context.function.code, fmt.Sprintf("loop $%s", test_label))
		InterpretNode(v.Test, context)
		context.function.code = append(context.function.code, "i32.const 1")
		context.function.code = append(context.function.code, "i32.xor")
		context.function.code = append(context.function.code, fmt.Sprintf("br_if $%s", exit_label))
		context.NewScope(func() {
			context.Define("break", &WASMVar{"", exit_label}) // we only fake using scope..
			context.Define("continue", &WASMVar{"", test_label})
			InterpretNode(&v.Body, context)
			context.function.code = append(context.function.code, fmt.Sprintf("br $%s", test_label))
			context.function.code = append(context.function.code, "end")

		})
		context.function.code = append(context.function.code, "end")

	case *model.FunctionDeclaration:

		oldfuc := context.function
		context.function = Function{
			name:       v.Name.Text,
			parameters: v.Parameters,
			retType:    v.ReturnType.Type(),
		}
		context.Define(v.Name.Text, &WASMVar{v.ReturnType.Type(), ""}) //
		context.NewScope(func() {
			context.scope = "local"
			for _, param := range v.Parameters {
				context.Define(param.Name.Text, &WASMVar{param.Type.Type(), "local"})
			}
			InterpretNode(&v.Body, context)
		})
		context.module = append(context.module, context.function.String())
		context.function = oldfuc
		context.scope = "global"

		if v.Name.Text == "main" {
			context.haveMain = true
		}

	case *model.FunctionApplication:
		argType := "int"
		for _, arg := range v.Arguments {
			argType = InterpretNode(arg, context) // arg eval in current context
		}
		name := v.Func.(*model.Name).Text
		log.Debugf("name %v", name)
		if name == "int" {
			// only float need to cast
			if argType == "float" {
				context.function.code = append(context.function.code, "i32.trunc_f64_s")
			}
			return "int"
		}
		if name == "float" {
			if argType != "float" {
				context.function.code = append(context.function.code, "f64.convert_i32_s")
			}
			return "float"
		}
		if name == "char" {
			return "char"
		}
		if name == "bool" {
			return "bool"
		}
		context.function.code = append(context.function.code, fmt.Sprintf("call $%s", name))
		val := context.Lookup(name)
		return val.Type
		// custom function and it should be....

	case *model.CompoundExpression:
		var val string
		context.NewScope(func() {
			for _, statement := range v.Statements.Statements[:len(v.Statements.Statements)-1] {
				// do we need pop to keep stack blance?
				InterpretNode(statement, context)
			}
			// return the last expression
			val = InterpretNode(
				v.Statements.Statements[len(v.Statements.Statements)-1].(*model.ExpressionAsStatement).Expression,
				context)
			log.Debugf("CompoundExpression1 %v", val)
		})
		log.Debugf("CompoundExpression %v", val)
		return val
	default:
		panic(fmt.Sprintf("Can't intepre %#v to source", v))
	}

	return ""
}
