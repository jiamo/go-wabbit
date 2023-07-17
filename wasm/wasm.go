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

type WasmFunction struct {
	name       string
	parameters []Parameter
	retType    string
	code       []string
	locals     []string
}

func (f *WasmFunction) String() string {
	out := fmt.Sprintf("(func $%s (export \"%s\")\n", f.name, f.name)
	for _, parm := range f.parameters {
		out += fmt.Sprintf("(param $%s %s)\n", parm.Name, _typemap[parm.Type])
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

type WabbitWasmModule struct {
	module   []string
	env      *common.ChainMap
	function WasmFunction
	scope    string
	nlabels  int
	haveMain bool
}

type WASMVar struct {
	Type  string
	Scope string
}

func NewWabbitWasmModule() *WabbitWasmModule {
	w := &WabbitWasmModule{
		module: []string{"(module"},
		env:    common.NewChainMap(),
		function: WasmFunction{
			name: "_init",
		},
		scope: "global",
	}
	w.module = append(w.module, "(import \"env\" \"_printi\" (func $_printi ( param i32 )))")
	w.module = append(w.module, "(import \"env\" \"_printf\" (func $_printf ( param f64 )))")
	w.module = append(w.module, "(import \"env\" \"_printb\" (func $_printb ( param i32 )))")
	w.module = append(w.module, "(import \"env\" \"_printc\" (func $_printc ( param i32 )))")

	return w
}

func (m *WabbitWasmModule) String() string {
	return strings.Join(m.module, "\n") + "\n)\n"
}

func (ctx *WabbitWasmModule) Define(name string, value *WASMVar) {
	ctx.env.SetValue(name, value)
}

func (ctx *WabbitWasmModule) Lookup(name string) *WASMVar {
	v, e := ctx.env.GetValue(name)
	if e == true {
		return v.(*WASMVar)
	} else {
		return nil
	}
}

func (ctx *WabbitWasmModule) NewScope(do func()) {
	oldEnv := ctx.env
	ctx.env = ctx.env.NewChild()
	defer func() {
		ctx.env = oldEnv
	}()
	do()
}

func (m *WabbitWasmModule) NewLabel() string {
	m.nlabels++
	return fmt.Sprintf("label%d", m.nlabels)
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
	// Grow the slice by one element.
	slice = append(slice, "")

	// Use copy to move the upper part of the slice out of the way and open a hole.
	copy(slice[position+1:], slice[position:])

	// Insert the new element.
	slice[position] = value

	return slice
}

func InterpretNode(node model.Node, context *WabbitWasmModule) string {
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
		context.function.code = append(context.function.code, fmt.Sprintf("f32.const %v", int(rune(unquoted[0]))))
		return "char"
	case *model.Name:

		value := context.Lookup(v.Text) // somethings we may need exist
		// bool int char using int
		if value.Scope == "global" {
			// that why value is simple than wvm
			context.function.code = append(context.function.code, fmt.Sprintf("global.get ${%v}", v.Text))

		} else if value.Scope == "local" {
			context.function.code = append(context.function.code, fmt.Sprintf("local.get ${%v}", v.Text))
		}
		return value.Type

	case *model.NameType:
		return v.Name
	case *model.NameBool:
		context.function.code = append(context.function.code, fmt.Sprintf("i32.const %v", BoolToInt(v.Name == "true")))
		//context.code = append(context.code, Instruction{"IPUSH", BoolToInt(v.Name == "true")})
		//return &WabbitValue{Type: "bool", Value: v.Name == "true"}
		return "bool"
	case *model.IntegerType:
		return "int"
		//return &WabbitValue{Type: "type", Value: "int"}
	case *model.FloatType:
		//return &WabbitValue{Type: "type", Value: "float"}
		return "float"
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
			context.function.code = append(context.function.code, "f64.mul")
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
			context.function.code = append(context.function.code, "i32.div")
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
			//return &WabbitValue{"int", -right.Value.(int)}
			//context.function.code = append(context.function.code, "i32.const 0")
			insert(context.function.code, "i32.const 0", pos)
			context.function.code = append(context.function.code, "i32.sub")
		} else if right == "float" {
			//return &WabbitValue{"float", -right.Value.(float64)}
			insert(context.function.code, "f64.const 0", pos)
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
				context.module = append(context.module, "(global (mut f64) (f64.const 0.0))")
			} else {
				context.module = append(context.module, "(global (mut i32) (i32.const 0))")
			}
			if v.Value != nil {
				context.function.code = append(context.function.code, fmt.Sprintf("global.set $%s", v.Name.Text))
			}
		} else if context.scope == "local" {
			if valtype == "float" {
				context.module = append(context.module, "(local (mut f64) (f64.const 0.0))")
			} else {
				context.module = append(context.module, "(local (mut i32) (i32.const 0))")
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
				context.module = append(context.module, "(global (mut f64) (f64.const 0.0))")
			} else {
				context.module = append(context.module, "(global (mut i32) (i32.const 0))")
			}
			context.function.code = append(context.function.code, fmt.Sprintf("global.set $%s", v.Name.Text))
		} else if context.scope == "local" {
			if valtype == "float" {
				context.module = append(context.module, "(local (mut f64) (f64.const 0.0))")
			} else {
				context.module = append(context.module, "(local (mut i32) (i32.const 0))")
			}
			context.function.code = append(context.function.code, fmt.Sprintf("local.set $%s", v.Name.Text))
		}
		context.Define(v.Name.Text, &WASMVar{Type: valtype, Scope: context.scope})
		return ""

	case *model.Lt:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.code = append(context.code, Instruction{"ICMP", "<"})
		} else if left == "float" {
			context.code = append(context.code, Instruction{"FCMP", "<"})
		} else if left == "char" {
			context.code = append(context.code, Instruction{"ICMP", "<"})
		} else {
			panic("type differnt")
		}
		return "bool"
	case *model.Le:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.code = append(context.code, Instruction{"ICMP", "<="})
		} else if left == "float" {
			context.code = append(context.code, Instruction{"FCMP", "<="})
		} else if left == "char" {
			context.code = append(context.code, Instruction{"ICMP", "<="})
		} else if left == "bool" {
			context.code = append(context.code, Instruction{"ICMP", "<="})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.Gt:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.code = append(context.code, Instruction{"ICMP", ">"})
		} else if left == "float" {
			context.code = append(context.code, Instruction{"FCMP", ">"})
		} else if left == "char" {
			context.code = append(context.code, Instruction{"ICMP", ">"})
		} else if left == "bool" {
			context.code = append(context.code, Instruction{"ICMP", ">"})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.Ge:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.code = append(context.code, Instruction{"ICMP", ">="})
		} else if left == "float" {
			context.code = append(context.code, Instruction{"FCMP", ">="})
		} else if left == "char" {
			context.code = append(context.code, Instruction{"ICMP", ">="})
		} else if left == "bool" {
			context.code = append(context.code, Instruction{"ICMP", ">="})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.Eq:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			context.code = append(context.code, Instruction{"ICMP", "=="})
		} else if left == "float" {
			context.code = append(context.code, Instruction{"FCMP", "=="})
		} else if left == "char" {
			context.code = append(context.code, Instruction{"ICMP", "=="})
		} else if left == "bool" {
			context.code = append(context.code, Instruction{"ICMP", "=="})
		} else {
			panic("type different") // we can just using else For simple
		}
		return "bool"
	case *model.Ne:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.code = append(context.code, Instruction{"ICMP", "!="})
		} else if left == "float" {
			context.code = append(context.code, Instruction{"FCMP", "!="})
		} else if left == "char" {
			context.code = append(context.code, Instruction{"ICMP", "!="})
		} else if left == "bool" {
			context.code = append(context.code, Instruction{"ICMP", "!="})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.LogOr:
		// TODO short eval
		done_label := context.NewLabel()
		or_continue_label := context.NewLabel()

		_ = InterpretNode(v.Left, context)
		context.code = append(context.code, Instruction{"BZ", or_continue_label})
		context.code = append(context.code, Instruction{"IPUSH", 1})
		context.code = append(context.code, Instruction{"GOTO", done_label})
		context.code = append(context.code, Instruction{"LABEL", or_continue_label})
		_ = InterpretNode(v.Right, context)
		context.code = append(context.code, Instruction{"LABEL", done_label})
		//context.code = append(context.code, Instruction{"OR", nil})
		return "bool"

	case *model.LogAnd:
		done_label := context.NewLabel()
		and_false_label := context.NewLabel()
		_ = InterpretNode(v.Left, context)
		context.code = append(context.code, Instruction{"BZ", and_false_label})
		_ = InterpretNode(v.Right, context)
		context.code = append(context.code, Instruction{"GOTO", done_label})
		context.code = append(context.code, Instruction{"LABEL", and_false_label})
		context.code = append(context.code, Instruction{"IPUSH", 0})

		context.code = append(context.code, Instruction{"LABEL", done_label})
		return "bool" // no need or and any more

	case *model.Assignment:
		val := InterpretNode(v.Value, context)
		// assign the value to the name
		wvmvar := context.Lookup(v.Location.(*model.Name).Text)
		if wvmvar.Scope == "global" {
			if wvmvar.Type == "float" {
				context.code = append(context.code, Instruction{"FDUP", nil})
				context.code = append(context.code, Instruction{"FSTORE_GLOBAL", wvmvar.Slot})
			} else {
				context.code = append(context.code, Instruction{"IDUP", nil})
				context.code = append(context.code, Instruction{"ISTORE_GLOBAL", wvmvar.Slot})
			}
		} else {
			if wvmvar.Type == "float" {
				context.code = append(context.code, Instruction{"FDUP", nil})
				context.code = append(context.code, Instruction{"FSTORE_LOCAL", wvmvar.Slot})
			} else {
				context.code = append(context.code, Instruction{"IDUP", nil})
				context.code = append(context.code, Instruction{"ISTORE_LOCAL", wvmvar.Slot})
			}
		}
		return val

	case *model.PrintStatement:
		value := InterpretNode(v.Value, context)
		switch value {
		case "char":
			context.code = append(context.code, Instruction{"PRINTC", nil})
		case "bool":
			context.code = append(context.code, Instruction{"PRINTB", nil})
		case "int":
			context.code = append(context.code, Instruction{"PRINTI", nil})
		case "float":
			context.code = append(context.code, Instruction{"PRINTF", nil})
		default:
			log.Debugf("%v:%v", context.code, value)
			panic("wrong type")
		}
	case *model.Statements:

		var result string
		for _, statement := range v.Statements {
			// do we need pop to keep stack blance?
			if result == "float" {
				context.code = append(context.code, Instruction{"FPOP", nil})
			} else if result == "int" {
				context.code = append(context.code, Instruction{"IPOP", nil})
			}
			result = InterpretNode(statement, context)
			// need check break return too

		}
		return result

	case *model.ExpressionAsStatement:
		return InterpretNode(v.Expression, context)

	case *model.Grouping:
		return InterpretNode(v.Expression, context)

	case *model.IfStatement:
		then_label := context.NewLabel()
		else_label := context.NewLabel()
		merge_label := context.NewLabel()
		InterpretNode(v.Test, context)
		context.code = append(context.code, Instruction{"BZ", else_label})
		context.code = append(context.code, Instruction{"GOTO", then_label})
		context.code = append(context.code, Instruction{"LABEL", then_label})

		context.NewScope(
			func() {
				InterpretNode(&v.Consequence, context)
				context.code = append(context.code, Instruction{"GOTO", merge_label})
				context.code = append(context.code, Instruction{"LABEL", else_label})
			},
		)
		if v.Alternative != nil {
			context.NewScope(
				func() {
					InterpretNode(v.Alternative, context)
					context.code = append(context.code, Instruction{"GOTO", merge_label})
				},
			)
		}
		context.code = append(context.code, Instruction{"GOTO", merge_label})
		context.code = append(context.code, Instruction{"LABEL", merge_label})

	case *model.BreakStatement:
		// we need scope for level break
		val := context.Lookup("break") // fake using type as label
		context.code = append(context.code, Instruction{"GOTO", val.Slot})
	case *model.ContinueStatement:
		val := context.Lookup("continue") // fake using type as label
		context.code = append(context.code, Instruction{"GOTO", val.Slot})

	case *model.ReturnStatement:
		value := InterpretNode(v.Value, context)
		context.code = append(context.code, Instruction{"RETURN", nil})
		return value

	case *model.WhileStatement:
		test_label := context.NewLabel()
		body_label := context.NewLabel()
		exit_label := context.NewLabel()

		context.code = append(context.code, Instruction{"GOTO", test_label})
		context.code = append(context.code, Instruction{"LABEL", test_label})
		InterpretNode(v.Test, context)
		context.code = append(context.code, Instruction{"BZ", exit_label})
		context.code = append(context.code, Instruction{"GOTO", body_label})
		context.code = append(context.code, Instruction{"LABEL", body_label})
		context.NewScope(func() {
			context.Define("break", &WVMVar{"", "", exit_label})
			context.Define("continue", &WVMVar{"", "", test_label})
			InterpretNode(&v.Body, context)
			context.code = append(context.code, Instruction{"GOTO", test_label})

		})
		context.code = append(context.code, Instruction{"LABEL", exit_label})

	case *model.FunctionDeclaration:
		// we should check the function name is not defined
		// we can keep function into another position // that's what I am doing in 2022
		// and put function in the end....

		start_label := context.NewLabel()
		end_label := context.NewLabel()
		// we don't put function

		context.code = append(context.code, Instruction{"GOTO", end_label})
		context.code = append(context.code, Instruction{"LABEL", start_label})
		context.Define(v.Name.Text, &WVMVar{v.ReturnType.Type(), "", start_label}) //
		context.NewScope(func() {
			context.scope = "local"
			for _, param := range v.Parameters {
				scope, slot := context.NewVariable()
				context.Define(param.Name.Text, &WVMVar{param.Type.Type(), scope, slot})
			}
			for i := len(v.Parameters) - 1; i >= 0; i-- {
				val := context.Lookup(v.Parameters[i].Name.Text)
				if val.Type == "float" {
					context.code = append(context.code, Instruction{"FSTORE_LOCAL", val.Slot})
				} else if val.Type == "int" {
					context.code = append(context.code, Instruction{"ISTORE_LOCAL", val.Slot})
				}
			}
			InterpretNode(&v.Body, context)
		})
		context.code = append(context.code, Instruction{"LABEL", end_label})
		context.scope = "global"
		if v.Name.Text == "main" {
			context.haveMain = true
		}

	case *model.FunctionApplication:
		argType := "int"
		//value := InterpretNode(v.Func, context) // while lookup
		for _, arg := range v.Arguments {
			// TODO check the type
			//fmt.Println("arg %v", arg)
			argType = InterpretNode(arg, context) // arg eval in current context
		}
		//
		//savedContext = funtionClosure.Value.Context
		// TODO make it as builtin function...
		name := v.Func.(*model.Name).Text
		funcVar := context.Lookup(v.Func.(*model.Name).Text) // define in
		log.Debugf("name %v", name)
		if name == "int" {
			// only float need to cast
			if argType == "float" {
				context.code = append(context.code, Instruction{"FTOI", nil})
			}
			return "int"
		}
		if name == "float" {
			if argType != "float" {
				context.code = append(context.code, Instruction{"ITOF", nil})
			}
			return "float"
		}
		if name == "char" {
			return "char"
		}
		if name == "bool" {
			return "bool"
		}
		context.code = append(context.code, Instruction{"CALL", funcVar.Slot})
		return funcVar.Type
		// custom function and it should be....

	case *model.CompoundExpression:
		var val string
		context.NewScope(func() {
			val = InterpretNode(&v.Statements, context)
		})
		return val
	default:
		panic(fmt.Sprintf("Can't intepre %#v to source", v))
	}

	return ""
}
