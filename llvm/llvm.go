package llvm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"strconv"
	"strings"
	"wabbit-go/common"
	"wabbit-go/model"
)

type LValue struct {
	WType  string
	LValue string
	Scope  string
}

type Parameter struct {
	PType  string
	PValue string
}

type Function struct {
	name       string
	parameters []model.Parameter
	retType    string
	code       []string
	locals     map[string]string // may be map name and regster
}

var _typemap = map[string]string{
	"":      "void",
	"int":   "i32",
	"float": "double",
	"bool":  "i1",
	"char":  "i8",
}

var _zero = map[string]string{
	"int":   "0",
	"float": "0x0",
	"bool":  "0",
	"char":  "0",
}

func (f *Function) String() string {
	var parms []string
	for index, parm := range f.parameters {
		argname := fmt.Sprintf("%%\".%d\"", index+1)
		parms = append(parms, fmt.Sprintf("%s %s", _typemap[parm.Type.Type()], argname))
	}
	parmstr := strings.Join(parms, ", ")
	out := fmt.Sprintf("define %s @\"%s\"(%s)\n{\n", _typemap[f.retType], f.name, parmstr)
	out += strings.Join(f.code, "\n  ")
	out += "\n}"
	return out
}

type Context struct {
	N        int
	nlabels  int
	globals  []string
	code     []string
	function Function
	scope    string
	env      *common.ChainMap
}

func (ctx *Context) NewRegister() string {
	ctx.N++
	return fmt.Sprintf("%%\".%d\"", ctx.N)
}
func (ctx *Context) NewLabel() string {
	ctx.N++
	return fmt.Sprintf("\".%d\"", ctx.N)
}

func (ctx *Context) Define(name string, value *LValue) {
	ctx.env.SetValue(name, value)
}

func (ctx *Context) Lookup(name string) *LValue {
	v, e := ctx.env.GetValue(name)
	if e == true {
		return v.(*LValue)
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

func LLVM(program *model.Program) string {
	context := &Context{
		N: 0,
		globals: []string{
			"declare void @\"_printi\"(i32 %\".1\")",
			"declare void @\"_printf\"(double %\".1\")",
			"declare void @\"_printb\"(i1 %\".1\")",
			"declare void @\"_printc\"(i8 %\".1\")",
		},
		function: Function{name: "main", retType: "int"},

		//code: []string{
		//	"define void @main()",
		//	"{",
		//},
		env:   common.NewChainMap(),
		scope: "global",
	}
	context.function.code = append(context.function.code, "entry:")
	_ = InterpretNode(program.Model, context) // generate is InterpretNode in the same meaning
	context.function.code = append(context.function.code, "ret i32 0")
	begin := "; ModuleID = \"wabbit\"\ntarget triple = \"unknown-unknown-unknown\"\ntarget datalayout = \"\"\n\n"
	return begin + context.function.String() + "\n\n" +
		strings.Join(context.globals, "\n\n") + "\n\n"

}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func InterpretNode(node model.Node, context *Context) *LValue {
	switch v := node.(type) {
	case *model.Integer:
		//context.function.code = append(context.function.code, fmt.Sprintf("i32.const %v", v.Value))
		return &LValue{"int", fmt.Sprintf("%v", v.Value), ""} // why Need Stringint
	case *model.Float:
		bits := math.Float64bits(v.Value)

		// OR using %f
		//fmt.Printf("0x%x\n", bits)
		//context.function.code = append(context.function.code, fmt.Sprintf("f64.const %v", v.Value))
		return &LValue{"float", fmt.Sprintf("0x%x", bits), ""}
	case *model.Character:
		unquoted, err := strconv.Unquote(v.Value)
		if err != nil {
			panic(err)
		}
		//context.code = append(context.code, Instruction{"IPUSH", int(rune(unquoted[0]))})
		//context.function.code = append(context.function.code, fmt.Sprintf("i32.const %v", int(rune(unquoted[0]))))
		return &LValue{"char", fmt.Sprintf("%v", int(rune(unquoted[0]))), ""}
	case *model.Name:
		//func square(x int) int {
		//    return x*x;  x is not declared
		//}
		value := context.Lookup(v.Text)
		r := context.NewRegister()
		ltype := _typemap[value.WType]
		context.function.code = append(context.function.code,
			fmt.Sprintf("%s = load %s, %s* %s", r, ltype, ltype, value.LValue))
		return &LValue{value.WType, r, ""}

	//case *model.NameType:
	//	return v.Name
	// may be not need
	case *model.NameBool:
		//context.function.code = append(context.function.code, fmt.Sprintf("i32.const %v", BoolToInt(v.Name == "true")))
		return &LValue{"bool", fmt.Sprintf("%v", BoolToInt(v.Name == "true")), ""}
	//case *model.IntegerType:
	//	return "int"
	//return &WabbitValue{Type: "type", Value: "int"}
	//case *model.FloatType:
	//	//return &WabbitValue{Type: "type", Value: "float"}
	//	return "float"
	case *model.Add:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		// we should check the type of left and right go we can't make interface + interface
		if left.WType == "int" && right.WType == "int" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = add i32 %s, %s", val, left.LValue, right.LValue))

		} else if left.WType == "float" && right.WType == "float" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = fadd double %s, %s", val, left.LValue, right.LValue))
		} else {
			panic("type different")
		}
		return &LValue{left.WType, val, ""}
	case *model.Mul:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		// we should check the type of left and right go we can't make interface + interface
		if left.WType == "int" && right.WType == "int" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = mul i32 %s, %s", val, left.LValue, right.LValue))

		} else if left.WType == "float" && right.WType == "float" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = fmul double %s, %s", val, left.LValue, right.LValue))
		} else {
			panic("type different")
		}
		return &LValue{left.WType, val, ""}
	case *model.Sub:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		// we should check the type of left and right go we can't make interface + interface
		if left.WType == "int" && right.WType == "int" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = sub i32 %s, %s", val, left.LValue, right.LValue))

		} else if left.WType == "float" && right.WType == "float" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = fsub double %s, %s", val, left.LValue, right.LValue))
		} else {
			panic("type different")
		}
		return &LValue{left.WType, val, ""}
	case *model.Div:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		// we should check the type of left and right go we can't make interface + interface
		if left.WType == "int" && right.WType == "int" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = sdiv i32 %s, %s", val, left.LValue, right.LValue))

		} else if left.WType == "float" && right.WType == "float" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = fdiv double %s, %s", val, left.LValue, right.LValue))
		} else {
			panic("type different")
		}
		return &LValue{left.WType, val, ""}

	case *model.Neg:
		right := InterpretNode(v.Operand, context)
		val := context.NewRegister()
		if right.WType == "int" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = sub i32 0, %s", val, right.LValue))

		} else if right.WType == "float" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = fsub double 0x0, %s", val, right.LValue)) // or fneg double
		} else {
			// we think it's a type error
			//return &WabbitValue{"error", "type error"}
			panic("type different")
		}
		return &LValue{right.WType, val, context.scope}

	case *model.Pos:
		right := InterpretNode(v.Operand, context)
		return right
	case *model.Not:
		right := InterpretNode(v.Operand, context)
		val := context.NewRegister()
		if right.WType == "bool" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = xor i1 1, %s", val, right.LValue))
		} else {
			panic("type different")
		}
		return &LValue{right.WType, val, ""}
	case *model.VarDeclaration:
		// make value number like llvmlite
		context.N++
		var val *LValue
		var valtype string
		if v.Value != nil {
			val = InterpretNode(v.Value, context) // store in stack
			valtype = val.WType
		} else {
			valtype = v.Type.Type()
		}
		if context.scope == "global" {
			context.globals = append(context.globals,
				fmt.Sprintf("@\"%s\" = global %s %s", v.Name.Text, _typemap[valtype], _zero[valtype]))
			if val != nil {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store %s %s, %s* @\"%s\"",
						_typemap[valtype], val.LValue, _typemap[valtype], v.Name.Text))
				//context.N++
			}
			context.Define(v.Name.Text, &LValue{valtype,
				fmt.Sprintf("@\"%s\"", v.Name.Text), context.scope}) // TODO

		} else if context.scope == "local" {
			// local using function
			// need first alloc
			context.function.code = append(context.function.code,
				fmt.Sprintf("%%\"%s\" = alloca %s", v.Name.Text, _typemap[valtype]))
			if val != nil {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store %s %s, %s* %%\"%s\"",
						_typemap[valtype], val.LValue, _typemap[valtype], v.Name.Text))
				//context.N++
			} else {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store %s %s, %s* %%\"%s\"",
						_typemap[valtype], _zero[valtype], _typemap[valtype], v.Name.Text))
			}
			context.Define(v.Name.Text, &LValue{valtype,
				fmt.Sprintf("%%\"%s\"", v.Name.Text), context.scope}) // TODO

		}

	case *model.ConstDeclaration:
		//var val *WVMVar
		context.N++
		var val *LValue
		var valtype string
		if v.Value != nil {
			val = InterpretNode(v.Value, context) // store in stack
			valtype = val.WType
		} else {
			valtype = v.Type.Type()
			//val = &LValue{valtype, "0", ""}
		}
		if context.scope == "global" {
			context.globals = append(context.globals,
				fmt.Sprintf("@\"%s\" = global %s %s", v.Name.Text, _typemap[valtype], _zero[valtype]))
			if val != nil {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store %s %s, %s* @\"%s\"",
						_typemap[valtype], val.LValue, _typemap[valtype], v.Name.Text))
			}
			context.Define(v.Name.Text, &LValue{valtype,
				fmt.Sprintf("@\"%s\"", v.Name.Text), context.scope}) // TODO

		} else if context.scope == "local" {
			// local using function
			// need first alloc
			context.function.code = append(context.function.code,
				fmt.Sprintf("%%\"%s\" = alloca %s", v.Name.Text, _typemap[valtype]))
			if val != nil {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store %s %s, %s* %%\"%s\"",
						_typemap[valtype], val.LValue, _typemap[valtype], v.Name.Text))
				//context.N++
			} else {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store %s %s, %s* @\"%s\"",
						_typemap[valtype], _zero[valtype], _typemap[valtype], v.Name.Text))
			}
			context.Define(v.Name.Text, &LValue{valtype,
				fmt.Sprintf("%%\"%s\"", v.Name.Text), context.scope}) // TO
		}

	case *model.Lt:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		if left.WType == "float" {
			context.function.code = append(context.function.code, fmt.Sprintf("%s = fcmp olt double %s, %s", val, left.LValue, right.LValue))
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, fmt.Sprintf("%s = icmp slt %s %s, %s", val, ltype, left.LValue, right.LValue))
		}
		return &LValue{"bool", val, context.scope}
	case *model.Le:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		if left.WType == "float" {
			context.function.code = append(context.function.code, fmt.Sprintf("%s = fcmp ole double %s, %s", val, left.LValue, right.LValue))
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, fmt.Sprintf("%s = icmp sle %s %s, %s", val, ltype, left.LValue, right.LValue))
		}
		return &LValue{"bool", val, context.scope}
	case *model.Gt:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		if left.WType == "float" {
			context.function.code = append(context.function.code, fmt.Sprintf("%s = fcmp ogt double %s, %s", val, left.LValue, right.LValue))
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, fmt.Sprintf("%s = icmp sgt %s %s, %s", val, ltype, left.LValue, right.LValue))
		}
		return &LValue{"bool", val, context.scope}
	case *model.Ge:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		if left.WType == "float" {
			context.function.code = append(context.function.code, fmt.Sprintf("%s = fcmp oge double %s, %s", val, left.LValue, right.LValue))
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, fmt.Sprintf("%s = icmp sge %s %s, %s", val, ltype, left.LValue, right.LValue))
		}
		return &LValue{"bool", val, context.scope}
	case *model.Eq:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		if left.WType == "float" {
			context.function.code = append(context.function.code, fmt.Sprintf("%s = fcmp oeq double %s, %s", val, left.LValue, right.LValue))
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, fmt.Sprintf("%s = icmp eq %s %s, %s", val, ltype, left.LValue, right.LValue))
		}
		return &LValue{"bool", val, context.scope}
	case *model.Ne:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		if left.WType == "float" {
			context.function.code = append(context.function.code, fmt.Sprintf("%s = fcmp one double %s, %s", val, left.LValue, right.LValue))
		} else {
			//return &WabbitValue{Type: "bool", Value: left.Value.(int) < right.Value.(int)}
			context.function.code = append(context.function.code, fmt.Sprintf("%s = icmp ne %s %s, %s", val, ltype, left.LValue, right.LValue))
		}
		return &LValue{"bool", val, context.scope}
	case *model.LogOr:
		// TODO short eval
		// logOr like If
		consequence := context.NewLabel()
		alternative := context.NewLabel()
		merge := context.NewLabel()
		var right *LValue
		left := InterpretNode(v.Left, context)
		context.function.code = append(context.function.code,
			fmt.Sprintf("br i1 %s, label %%%s, label %%%s", left.LValue, consequence, alternative))
		context.function.code = append(context.function.code, fmt.Sprintf("%s:", consequence))
		context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", merge))

		context.NewScope(func() {
			context.function.code = append(context.function.code, fmt.Sprintf("%s:", alternative))
			right = InterpretNode(v.Right, context)
			context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", merge))
		})

		context.function.code = append(context.function.code, fmt.Sprintf("%s:", merge))

		ret := context.NewRegister()
		context.function.code = append(context.function.code,
			fmt.Sprintf("%s = phi i1 [ %s, %%%s ], [ %s, %%%s ]",
				ret, left.LValue, consequence, right.LValue, alternative))
		return &LValue{"bool", ret, context.scope}

	case *model.LogAnd:
		consequence := context.NewLabel()
		alternative := context.NewLabel()
		merge := context.NewLabel()
		var right *LValue
		left := InterpretNode(v.Left, context)
		context.function.code = append(context.function.code,
			fmt.Sprintf("br i1 %s, label %%%s, label %%%s", left.LValue, consequence, alternative))
		context.function.code = append(context.function.code, fmt.Sprintf("%s:", alternative))
		context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", merge))

		context.NewScope(func() {
			context.function.code = append(context.function.code, fmt.Sprintf("%s:", consequence))
			right = InterpretNode(v.Right, context)
			context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", merge))
		})

		context.function.code = append(context.function.code, fmt.Sprintf("%s:", merge))

		ret := context.NewRegister()
		context.function.code = append(context.function.code,
			fmt.Sprintf("%s = phi i1 [ %s, %%%s ], [ %s, %%%s ] ",
				ret, left.LValue, alternative, right.LValue, consequence))
		return &LValue{"bool", ret, context.scope}

	case *model.Assignment:
		//context.N++ // why assignment need ++
		val := InterpretNode(v.Value, context)
		context.N++
		// assign the value to the name
		ltype := _typemap[val.WType]
		decl := context.Lookup(v.Location.(*model.Name).Text)
		context.function.code = append(context.function.code,
			fmt.Sprintf("store %s %s, %s* %s", ltype, val.LValue, ltype, decl.LValue))

		return val

	case *model.PrintStatement:
		value := InterpretNode(v.Value, context)
		context.N++
		ltype := value.WType
		switch ltype {
		case "char":
			context.function.code = append(context.function.code,
				fmt.Sprintf("call void @\"_printc\"(%s %s)", _typemap[ltype], value.LValue))
		case "bool":
			context.function.code = append(context.function.code,
				fmt.Sprintf("call void @\"_printb\"(%s %s)", _typemap[ltype], value.LValue))
		case "int":
			context.function.code = append(context.function.code,
				fmt.Sprintf("call void @\"_printi\"(%s %s)", _typemap[ltype], value.LValue))
		case "float":
			context.function.code = append(context.function.code,
				fmt.Sprintf("call void @\"_printf\"(%s %s)", _typemap[ltype], value.LValue))
		default:
			panic("wrong type")
		}
	case *model.Statements:
		// llvm no need drop....
		var result *LValue
		for _, statement := range v.Statements {
			result = InterpretNode(statement, context)
		}
		return result

	case *model.ExpressionAsStatement:
		result := InterpretNode(v.Expression, context)
		return result

	case *model.Grouping:
		return InterpretNode(v.Expression, context)

	case *model.IfStatement:

		cons_label := context.NewLabel()
		alt_label := context.NewLabel()
		merge_label := context.NewLabel()
		test := InterpretNode(v.Test, context)
		context.function.code = append(context.function.code,
			fmt.Sprintf("br i1 %s, label %%%s, label %%%s", test.LValue, cons_label, alt_label))
		context.function.code = append(context.function.code, fmt.Sprintf("%s:", cons_label))
		InterpretNode(&v.Consequence, context)
		context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", merge_label))
		context.function.code = append(context.function.code, fmt.Sprintf("%s:", alt_label))
		if v.Alternative != nil {
			InterpretNode(v.Alternative, context)
		}
		// we should using block......
		context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", merge_label))
		context.function.code = append(context.function.code, fmt.Sprintf("%s:", merge_label))

		// llvmlite have is_terminated we have have too

	case *model.BreakStatement:
		// we need scope for level break
		val := context.Lookup("break") // fake using type as label
		context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", val.LValue))
	case *model.ContinueStatement:
		val := context.Lookup("continue") // fake using type as label
		context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", val.LValue))

	case *model.ReturnStatement:
		value := InterpretNode(v.Value, context)
		context.function.code = append(context.function.code,
			fmt.Sprintf("ret %s %s", _typemap[value.WType], value.LValue))
		return value
	//
	case *model.WhileStatement:
		test_label := context.NewLabel()
		body_label := context.NewLabel()
		exit_label := context.NewLabel()

		context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", test_label))
		context.function.code = append(context.function.code, fmt.Sprintf("%s:", test_label))
		test := InterpretNode(v.Test, context)
		context.function.code = append(context.function.code,
			fmt.Sprintf("br i1 %s, label %%%s, label %%%s", test.LValue, body_label, exit_label))

		context.NewScope(func() {
			context.Define("break", &LValue{"", exit_label, context.scope}) // we only fake using scope..
			context.Define("continue", &LValue{"", test_label, context.scope})
			context.function.code = append(context.function.code, fmt.Sprintf("%s:", body_label))
			InterpretNode(&v.Body, context)
			context.function.code = append(context.function.code, fmt.Sprintf("br label %%%s", test_label))
		})

		context.function.code = append(context.function.code, fmt.Sprintf("%s:", exit_label))

	case *model.FunctionDeclaration:
		oldfunc := context.function
		context.function = Function{
			name:       v.Name.Text,
			parameters: v.Parameters,
			retType:    v.ReturnType.Type(),
		}
		context.Define(v.Name.Text, &LValue{
			WType:  v.ReturnType.Type(),
			LValue: "", // TODO
			Scope:  context.scope,
		})
		context.NewScope(func() {
			context.scope = "local"
			for index, param := range v.Parameters {
				context.N++
				context.Define(param.Name.Text, &LValue{
					WType:  param.Type.Type(),
					LValue: fmt.Sprintf("%%\"%s\"", param.Name.Text), // TODO
					Scope:  context.scope,
				})
				pname := fmt.Sprintf("%%\"%s\"", param.Name.Text)
				argName := fmt.Sprintf("%%\".%d\"", index+1)
				ptype := _typemap[param.Type.Type()]
				// need define

				context.function.code = append(context.function.code,
					fmt.Sprintf("%s  = alloca %s", pname, ptype))
				context.function.code = append(context.function.code,
					fmt.Sprintf("store %s %s, %s* %s ", ptype, argName, ptype, pname))

			}
			InterpretNode(&v.Body, context)
			context.function.code = append(context.function.code,
				fmt.Sprintf("ret %s %s", _typemap[v.ReturnType.Type()], _zero[v.ReturnType.Type()]))

		})
		log.Debug("begining function", context.function.String())
		context.globals = append(context.globals, context.function.String())
		context.function = oldfunc
		context.scope = "global"

	//
	case *model.FunctionApplication:

		name := v.Func.(*model.Name).Text
		if name == "int" {
			// only float need to cast
			argVal := InterpretNode(v.Arguments[0], context)
			var result string
			if argVal.WType == "float" {
				result = context.NewRegister()
				ltype := _typemap[argVal.WType]
				context.function.code = append(context.function.code,
					fmt.Sprintf("%s = fptosi %s %s to i32", result, ltype, argVal.LValue))
			} else {
				result = argVal.LValue
			}
			return &LValue{
				WType:  "int",
				LValue: result,
			}
		}
		if name == "float" {
			argVal := InterpretNode(v.Arguments[0], context)
			var result string
			if argVal.WType != "float" {
				result = context.NewRegister()
				ltype := _typemap[argVal.WType]
				context.function.code = append(context.function.code,
					fmt.Sprintf("%s = sitofp %s %s to double", result, ltype, argVal.LValue))
			} else {
				result = argVal.LValue
			}
			return &LValue{
				WType:  "float",
				LValue: result,
			}
		}
		if name == "char" {
			argVal := InterpretNode(v.Arguments[0], context)
			var result string
			if argVal.WType != "char" {
				result = context.NewRegister()
				ltype := _typemap[argVal.WType]
				context.function.code = append(context.function.code,
					fmt.Sprintf("%s = trunc %s %s to i8", result, ltype, argVal.LValue))
			} else {
				result = argVal.LValue
			}
			return &LValue{
				WType:  "char",
				LValue: result,
			}
		}

		var result = context.NewRegister()
		var argValues []*LValue
		var argValuesStr []string
		for _, arg := range v.Arguments {
			argVal := InterpretNode(arg, context)
			argValues = append(argValues)
			argValuesStr = append(argValuesStr, fmt.Sprintf("%s %s",
				_typemap[argVal.WType], argVal.LValue))
		}
		funcVar := context.Lookup(v.Func.(*model.Name).Text)
		funcName := fmt.Sprintf("@\"%s\"", v.Func.(*model.Name).Text)
		context.function.code = append(context.function.code,
			fmt.Sprintf("%s = call %s %s(%s)", result, _typemap[funcVar.WType], funcName,
				strings.Join(argValuesStr, ", ")))

		return &LValue{
			WType:  funcVar.WType,
			LValue: result,
		}

	case *model.CompoundExpression:
		var val *LValue
		context.NewScope(func() {
			val = InterpretNode(&v.Statements, context)
		})
		return val

	default:
		panic(fmt.Sprintf("Can't intepre %#v to source", v))
	}

	return nil
}
