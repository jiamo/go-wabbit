package llvm

import (
	"fmt"
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

func NewLValue(wtype string, lvalue string, scope string) *LValue {
	return &LValue{
		WType:  wtype,
		LValue: lvalue,
		Scope:  scope,
	}
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
	parms := []string{}
	for _, parm := range f.parameters {
		parms = append(parms, fmt.Sprintf("$%s %%\"%s\")", _typemap[parm.Type.Type()],
			f.locals[parm.Name.Text])) // TODO.....
	}
	parmstr := strings.Join(parms, ", ")
	out := fmt.Sprintf("define %s @\"%s\"(%s)\n{\n", _typemap[f.retType], f.name, parmstr)

	//if f.retType != "" {
	//	out += fmt.Sprintf("(result %s)\n", _typemap[f.retType])
	//	out += fmt.Sprintf("(local $return %s)\n", _typemap[f.retType])
	//}
	//out += "\n" + strings.Join(f.locals, "\n")
	//out += "\nblock $return\n"
	out += strings.Join(f.code, "\n  ")
	//out += "\nend\n"
	//if f.retType != "" {
	//	out += "local.get $return\n"
	//}
	out += "\n}"
	return out
}

type LLVMContext struct {
	N        int
	nlabels  int
	globals  []string
	code     []string
	function Function
	scope    string
	env      *common.ChainMap
}

func (l *LLVMContext) NewRegister() string {
	l.N++
	return fmt.Sprintf("%%\".%d\"", l.N) // we using llvmlite format
}
func (l *LLVMContext) NewLabel() string {
	l.N++
	return fmt.Sprintf("\".%d\"", l.N) // we using llvmlite format
}

// may be
func (ctx *LLVMContext) Define(name string, value *LValue) {
	ctx.env.SetValue(name, value)
}

func (ctx *LLVMContext) Lookup(name string) *LValue {
	v, e := ctx.env.GetValue(name)
	if e == true {
		return v.(*LValue)
	} else {
		return nil
	}
}

func (ctx *LLVMContext) NewScope(do func()) {
	oldEnv := ctx.env
	ctx.env = ctx.env.NewChild()
	defer func() {
		ctx.env = oldEnv
	}()
	do()
}

func LLVM(program *model.Program) string {
	context := &LLVMContext{
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
	//context.function.code = append(context.function.code, "}")

	// new function was put in global d
	// this function is Main
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

func insert(slice []string, value string, position int) []string {
	slice = append(slice[:position], append([]string{value}, slice[position:]...)...)
	return slice
}

// may be not same like wasm generate the code directly
// we just using code

func InterpretNode(node model.Node, context *LLVMContext) *LValue {
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
		return &LValue{"int", fmt.Sprintf("%v", int(rune(unquoted[0]))), ""}
	case *model.Name:
		// may be the scope is not Important...!
		// %".4" = load double, double* @"xmin"
		//value := context.Lookup(v.Text)
		//if value.Scope == "global" {
		//	context.function.code = append(context.function.code, fmt.Sprintf("global.get $%s", v.Text))
		//} else if value.Scope == "local" {
		//	context.function.code = append(context.function.code, fmt.Sprintf("local.get $%s", v.Text))
		//}
		// just load it
		value := context.Lookup(v.Text)
		r := context.NewRegister()
		ltype := _typemap[value.WType]
		context.function.code = append(context.function.code,
			fmt.Sprintf("%s = load %s, %s* @\"%s\"", r, ltype, ltype, v.Text))
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
			//context.function.code = append(context.function.code, "i32.const 1")
			//context.function.code = append(context.function.code, "i32.xor")
			context.function.code = append(context.function.code,
				fmt.Sprintf("%s = xor i1 1, %s", val, right.LValue))
		} else {
			// we think it's a type error
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
			//val = &LValue{valtype, "0", ""}
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
			if valtype == "float" {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store double %s, double* @\"%s\"", val.LValue, v.Name.Text))
			} else {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store i32 %s, i32* @\"%s\"", val.LValue, v.Name.Text))
			}
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
			if valtype == "float" {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store double %s, double* @\"%s\"", val.LValue, v.Name.Text))
			} else {
				context.function.code = append(context.function.code,
					fmt.Sprintf("store i32 %s, i32* @\"%s\"", val.LValue, v.Name.Text))
			}
		}
		// the local
		// context.Define(v.Name.Text, &LValue{valtype, val,context.scope})

		return nil
	//
	//case *model.ConstDeclaration:
	//	valtype := InterpretNode(v.Value, context)
	//	if context.scope == "global" {
	//		if valtype == "float" {
	//			context.module = append(context.module, fmt.Sprintf("(global $%s (mut f64) (f64.const 0.0))", v.Name.Text))
	//		} else {
	//			context.module = append(context.module, fmt.Sprintf("(global $%s (mut i32) (i32.const 0))", v.Name.Text))
	//		}
	//		context.function.code = append(context.function.code, fmt.Sprintf("global.set $%s", v.Name.Text))
	//	} else if context.scope == "local" {
	//		if valtype == "float" {
	//			context.function.locals = append(context.function.locals, fmt.Sprintf("(local $%s (mut f64) (f64.const 0.0))", v.Name.Text))
	//		} else {
	//			context.function.locals = append(context.function.locals, fmt.Sprintf("(local $%s (mut i32) (i32.const 0))", v.Name.Text))
	//		}
	//		context.function.code = append(context.function.code, fmt.Sprintf("local.set $%s", v.Name.Text))
	//	}
	//	context.Define(v.Name.Text, &WASMVar{Type: valtype, Scope: context.scope})
	//	return ""

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
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		context.function.code = append(context.function.code, fmt.Sprintf("%s = and %s %s, %s", val, ltype, left.LValue, right.LValue))
		return &LValue{"bool", val, context.scope}
		//begin := context.NewLabel("begin")
		//context.function.code = append(context.function.code, fmt.Sprintf("block $%s (result i32)", begin))
		//
		//or_block := context.NewLabel("or_block")
		//context.function.code = append(context.function.code, fmt.Sprintf("block $%s", or_block))
		//_ = InterpretNode(v.Left, context)
		//context.function.code = append(context.function.code, fmt.Sprintf("br_if $%s", or_block))
		//_ = InterpretNode(v.Right, context)
		//context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		//context.function.code = append(context.function.code, "end")
		//context.function.code = append(context.function.code, "i32.const 1")
		//context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		//context.function.code = append(context.function.code, "end")
		//return "bool"

	case *model.LogAnd:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)
		val := context.NewRegister()
		ltype := _typemap[left.WType]
		context.function.code = append(context.function.code, fmt.Sprintf("%s = or %s %s, %s", val, ltype, left.LValue, right.LValue))
		return &LValue{"bool", val, context.scope}
		//begin := context.NewLabel("begin")
		//context.function.code = append(context.function.code, fmt.Sprintf("block $%s (result i32)", begin))
		//
		//and_block := context.NewLabel("and_block")
		//context.function.code = append(context.function.code, fmt.Sprintf("block $%s", and_block))
		//_ = InterpretNode(v.Left, context)
		//context.function.code = append(context.function.code, "i32.const 1")
		//context.function.code = append(context.function.code, "i32.xor")
		//context.function.code = append(context.function.code, fmt.Sprintf("br_if $%s", and_block))
		//_ = InterpretNode(v.Right, context)
		//context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		//context.function.code = append(context.function.code, "end")
		//context.function.code = append(context.function.code, "i32.const 0")
		//context.function.code = append(context.function.code, fmt.Sprintf("br $%s", begin))
		//context.function.code = append(context.function.code, "end")
		//return "bool" // no need or and any more

	case *model.Assignment:
		//context.N++ // why assignment need ++
		val := InterpretNode(v.Value, context)
		context.N++
		// assign the value to the name
		ltype := _typemap[val.WType]
		decl := context.Lookup(v.Location.(*model.Name).Text)
		if decl.Scope == "global" {
			context.function.code = append(context.function.code,
				fmt.Sprintf("store %s %s, %s* %s", ltype, val.LValue, ltype, decl.LValue))
		} else {
			context.function.code = append(context.function.code, fmt.Sprintf("local.tee $%s", v.Location.(*model.Name).Text))
		}
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

		//var result string
		for _, statement := range v.Statements {
			// do we need pop to keep stack blance?
			//if result != "" {
			//	context.function.code = append(context.function.code, "drop")
			//}
			InterpretNode(statement, context)
			// need check break return too

		}
		//return result

	case *model.ExpressionAsStatement:
		InterpretNode(v.Expression, context)
		// why need drop... . Don't need drop I think
		//context.function.code = append(context.function.code, "drop")
		//return val

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
	//
	//case *model.WhileStatement:
	//	test_label := context.NewLabel()
	//	exit_label := context.NewLabel()
	//
	//	context.function.code = append(context.function.code, fmt.Sprintf("block $%s", exit_label))
	//	context.function.code = append(context.function.code, fmt.Sprintf("loop $%s", test_label))
	//	InterpretNode(v.Test, context)
	//	context.function.code = append(context.function.code, "i32.const 1")
	//	context.function.code = append(context.function.code, "i32.xor")
	//	context.function.code = append(context.function.code, fmt.Sprintf("br_if $%s", exit_label))
	//	context.NewScope(func() {
	//		context.Define("break", &WASMVar{"", exit_label}) // we only fake using scope..
	//		context.Define("continue", &WASMVar{"", test_label})
	//		InterpretNode(&v.Body, context)
	//		context.function.code = append(context.function.code, fmt.Sprintf("br $%s", test_label))
	//		context.function.code = append(context.function.code, "end")
	//
	//	})
	//	context.function.code = append(context.function.code, "end")
	//
	//case *model.FunctionDeclaration:
	//	// we should check the function name is not defined
	//	// we can keep function into another position // that's what I am doing in 2022
	//	// and put function in the end....
	//
	//	oldfuc := context.function
	//	context.function = WasmFunction{
	//		name:       v.Name.Text,
	//		parameters: v.Parameters,
	//		retType:    v.ReturnType.Type(),
	//	}
	//	context.Define(v.Name.Text, &WASMVar{v.ReturnType.Type(), ""}) //
	//	context.NewScope(func() {
	//		context.scope = "local"
	//		for _, param := range v.Parameters {
	//			context.Define(param.Name.Text, &WASMVar{param.Type.Type(), "local"})
	//		}
	//		InterpretNode(&v.Body, context)
	//	})
	//	context.module = append(context.module, context.function.String())
	//	context.function = oldfuc
	//	context.scope = "global"
	//
	//	if v.Name.Text == "main" {
	//		context.haveMain = true
	//	}
	//
	//case *model.FunctionApplication:
	//	argType := "int"
	//	//value := InterpretNode(v.Func, context) // while lookup
	//	for _, arg := range v.Arguments {
	//		argType = InterpretNode(arg, context) // arg eval in current context
	//	}
	//	//
	//	name := v.Func.(*model.Name).Text
	//	//funcVar := context.Lookup(v.Func.(*model.Name).Text) // define in
	//	log.Debugf("name %v", name)
	//	if name == "int" {
	//		// only float need to cast
	//		if argType == "float" {
	//			context.function.code = append(context.function.code, "i32.trunc_f64_s")
	//		}
	//		return "int"
	//	}
	//	if name == "float" {
	//		if argType != "float" {
	//			context.function.code = append(context.function.code, "f64.convert_i32_s")
	//		}
	//		return "float"
	//	}
	//	if name == "char" {
	//		return "char"
	//	}
	//	if name == "bool" {
	//		return "bool"
	//	}
	//	context.function.code = append(context.function.code, fmt.Sprintf("call $%s", name))
	//	val := context.Lookup(name)
	//	return val.Type
	//	// custom function and it should be....
	//
	//case *model.CompoundExpression:
	//	var val string
	//	context.NewScope(func() {
	//		for _, statement := range v.Statements.Statements[:len(v.Statements.Statements)-1] {
	//			// do we need pop to keep stack blance?
	//
	//			InterpretNode(statement, context)
	//			// need check break return too
	//
	//		}
	//		// here must using expression not statment
	//		// it is expression as Statment
	//		val = InterpretNode(
	//			v.Statements.Statements[len(v.Statements.Statements)-1].(*model.ExpressionAsStatement).Expression,
	//			context)
	//		log.Debugf("CompoundExpression1 %v", val)
	//	})
	//	log.Debugf("CompoundExpression %v", val)
	//	return val
	default:
		panic(fmt.Sprintf("Can't intepre %#v to source", v))
	}

	return nil
}
