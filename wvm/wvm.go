package wvm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"reflect"
	"strconv"
	"wabbit-go/common"
	"wabbit-go/model"
)

type Frame struct {
	returnPc  int
	locals    map[int]interface{}
	prevFrame *Frame
}

type WVM struct {
	pc      int
	istack  []int
	fstack  []float64
	globals map[int]interface{}
	labels  map[int]int
	frame   *Frame
	running bool
}

type Instruction struct {
	opcode string
	args   interface{}
}

func (vm *WVM) run(instructions []Instruction) {
	vm.pc = 0
	vm.running = true
	vm.labels = make(map[int]int)

	// TODO we may need optimse this pass In WvmContext
	for i, instruction := range instructions {
		if instruction.opcode == "LABEL" {
			vm.labels[instruction.args.(int)] = i
		}
	}
	for vm.running {
		op := instructions[vm.pc].opcode
		args := instructions[vm.pc].args
		vm.pc++
		//for _, arg := range args {
		//	args = append(args, vm.get(arg))
		//}
		//getattr(vm, op).call(args...)
		argValues := make([]reflect.Value, 1)
		//for i, arg := range args {
		//	argValues[i] = reflect.ValueOf(arg)
		//}
		//log.Debugf("method is %v", op)
		method := reflect.ValueOf(vm).MethodByName(op)
		if args == nil {
			//argValues[0] = reflect.ValueOf(nil)
			method.Call([]reflect.Value{}) // empty array is not about nil
		} else {
			argValues[0] = reflect.ValueOf(args)
			method.Call(argValues)
		}

	}
}

//func (vm *WVM) get(arg interface{}) interface{} {
//	switch arg.(type) {
//	case int:
//		return arg.(int)
//	case float64:
//		return arg.(float64)
//	case string:
//		return vm.globals[arg.(string)]
//	default:
//		panic("invalid argument type")
//	}
//}

func (vm *WVM) IPUSH(value int) {
	vm.istack = append(vm.istack, value)
}

func (vm *WVM) IPOP() int {
	index := len(vm.istack) - 1
	// Get the top element of the stack
	element := vm.istack[index]
	// Remove the top element of the stack
	vm.istack = vm.istack[:index]
	return element
}

func (vm *WVM) IDUP() {
	vm.istack = append(vm.istack, vm.istack[len(vm.istack)-1])
}
func (vm *WVM) FDUP() {
	vm.fstack = append(vm.fstack, vm.fstack[len(vm.fstack)-1])
}

func (vm *WVM) IADD() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left + right)
}

func (vm *WVM) FADD() {
	right := vm.FPOP()
	left := vm.FPOP()
	vm.FPUSH(left + right)
}

func (vm *WVM) ISUB() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left - right)
}

func (vm *WVM) FSUB() {
	right := vm.FPOP()
	left := vm.FPOP()
	vm.FPUSH(left - right)
}

func (vm *WVM) IMUL() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left * right)
}

func (vm *WVM) FMUL() {
	right := vm.FPOP()
	left := vm.FPOP()
	vm.FPUSH(left * right)
}

func (vm *WVM) IDIV() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left / right)
}

func (vm *WVM) FDIV() {
	right := vm.FPOP()
	left := vm.FPOP()
	vm.FPUSH(left / right)
}

func (vm *WVM) AND() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left & right)
}

func (vm *WVM) OR() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left | right)
}

func (vm *WVM) XOR() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left ^ right)
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (vm *WVM) ICMP(op string) {
	right := vm.IPOP()
	left := vm.IPOP()
	switch op {
	case "<":
		vm.IPUSH(BoolToInt(left < right))
	case "<=":
		vm.IPUSH(BoolToInt(left <= right))
	case ">":
		vm.IPUSH(BoolToInt(left > right))
	case ">=":
		vm.IPUSH(BoolToInt(left >= right))
	case "==":
		vm.IPUSH(BoolToInt(left == right))
	case "!=":
		vm.IPUSH(BoolToInt(left != right))
	}
}

func (vm *WVM) FCMP(op string) {
	right := vm.FPOP()
	left := vm.FPOP()
	switch op {
	case "<":
		vm.IPUSH(BoolToInt(left < right))
	case "<=":
		vm.IPUSH(BoolToInt(left <= right))
	case ">":
		vm.IPUSH(BoolToInt(left > right))
	case ">=":
		vm.IPUSH(BoolToInt(left >= right))
	case "==":
		vm.IPUSH(BoolToInt(left == right))
	case "!=":
		vm.IPUSH(BoolToInt(left != right))
	}
}

func (vm *WVM) FPUSH(value float64) {
	vm.fstack = append(vm.fstack, value)
}

//func (vm *WVM) ITOF() {
//	value := vm.IPOP()
//	vm.FPUSH(float64(value))
//}

func (vm *WVM) INEG() {
	vm.IPUSH(-vm.IPOP())
}

func (vm *WVM) FPOP() float64 {
	index := len(vm.fstack) - 1
	// Get the top element of the stack
	element := vm.fstack[index]
	// Remove the top element of the stack
	vm.fstack = vm.fstack[:index]
	return element
}

func (vm *WVM) FNEG() {
	vm.FPUSH(-vm.FPOP())
}

func (vm *WVM) PRINTI() {
	fmt.Println(vm.IPOP())
}

func (vm *WVM) PRINTF() {
	fmt.Println(vm.FPOP())
}

func (vm *WVM) PRINTB() {
	if vm.IPOP() == 0 {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
}

func (vm *WVM) PRINTC() {
	fmt.Printf("%c", rune(vm.IPOP()))
}

func (vm *WVM) FTOI() {
	vm.IPUSH(int(vm.FPOP()))
}

func (vm *WVM) ITOF() {
	vm.FPUSH(float64(vm.IPOP()))
}

func (vm *WVM) ISTORE_LOCAL(slot int) {
	vm.frame.locals[slot] = vm.IPOP()
}

func (vm *WVM) ILOAD_LOCAL(slot int) {
	vm.IPUSH(vm.frame.locals[slot].(int))
}

func (vm *WVM) FLOAD_LOCAL(slot int) {
	vm.FPUSH(vm.frame.locals[slot].(float64))
}

func (vm *WVM) FSTORE_LOCAL(slot int) {
	vm.frame.locals[slot] = vm.FPOP()
}

func (vm *WVM) ISTORE_GLOBAL(slot int) {
	vm.globals[slot] = vm.IPOP()
}

func (vm *WVM) ILOAD_GLOBAL(slot int) {
	vm.IPUSH(vm.globals[slot].(int))
}

func (vm *WVM) FLOAD_GLOBAL(slot int) {
	vm.FPUSH(vm.globals[slot].(float64))
}

func (vm *WVM) GOTO(name int) {
	vm.pc = vm.labels[name]
}

func (vm *WVM) BZ(name int) {
	if vm.IPOP() == 0 {
		vm.pc = vm.labels[name]
	}
}

func (vm *WVM) HALT() {
	vm.running = false
}

func (vm *WVM) LABEL(name int) {

}

func (vm *WVM) FSTORE_GLOBAL(slot int) {
	vm.globals[slot] = vm.FPOP()
}

func (vm *WVM) CALL(label int) {
	vm.frame = &Frame{vm.pc, make(map[int]interface{}), vm.frame}
	vm.pc = vm.labels[label]
}

func (vm *WVM) RETURN() {
	vm.pc = vm.frame.returnPc
	vm.frame = vm.frame.prevFrame
}

type WVMContext struct {
	env       *common.ChainMap
	code      []Instruction
	nglobals  int
	nlocals   int
	nlabels   int
	scope     string
	haveMain  bool
	parentEnv *map[string]interface{}
}

func NewWVMContext() *WVMContext {
	return &WVMContext{
		env:   common.NewChainMap(),
		scope: "global",
		code:  make([]Instruction, 0),
	}
}

type WVMVar struct {
	Type  string
	Scope string
	Slot  int
}

func (ctx *WVMContext) Define(name string, value *WVMVar) {
	ctx.env.SetValue(name, value)
}

func (ctx *WVMContext) Lookup(name string) *WVMVar {
	v, e := ctx.env.GetValue(name)
	if e == true {
		return v.(*WVMVar)
	} else {
		return nil
	}
}

func (ctx *WVMContext) NewVariable() (string, int) {
	if ctx.scope == "global" {
		ctx.nglobals++
		return "global", ctx.nglobals - 1
	} else {
		ctx.nlocals++
		return "local", ctx.nlocals - 1
	}
}

func (ctx *WVMContext) NewLabel() int {
	ctx.nlabels++
	return ctx.nlabels - 1
}

func (ctx *WVMContext) NewScope(do func()) {
	oldEnv := ctx.env
	ctx.env = ctx.env.NewChild()
	defer func() {
		ctx.env = oldEnv
	}()
	do()
}

func Wvm(program *model.Program) error {
	wctx := NewWVMContext()
	_ = InterpretNode(program.Model, wctx) // generate is InterpretNode in the same meaning
	wvm := &WVM{
		globals: make(map[int]interface{}),
	}

	wctx.code = append(wctx.code, Instruction{"HALT", nil})
	log.Debug(wctx.code)
	wvm.run(wctx.code)

	return nil
}

func InterpretNode(node model.Node, context *WVMContext) string {
	switch v := node.(type) {
	case *model.Integer:
		context.code = append(context.code, Instruction{"IPUSH", v.Value})
		return "int"
	case *model.Float:
		context.code = append(context.code, Instruction{"FPUSH", v.Value})
		return "float"
	case *model.Character:
		unquoted, err := strconv.Unquote(v.Value)
		if err != nil {
			panic(err)
		}
		context.code = append(context.code, Instruction{"IPUSH", int(rune(unquoted[0]))})
		return "char"
	case *model.Name:
		log.Debugf("*model.Name %v env %v code %v", v.Text, context.env, context.code)
		value := context.Lookup(v.Text) // somethings we may need exist
		// bool int char using int
		if value.Scope == "global" {
			if value.Type == "float" {
				context.code = append(context.code, Instruction{"FLOAD_GLOBAL", value.Slot})
			} else {
				context.code = append(context.code, Instruction{"ILOAD_GLOBAL", value.Slot})
			}
		} else if value.Scope == "local" {
			if value.Type == "float" {
				context.code = append(context.code, Instruction{"FLOAD_LOCAL", value.Slot})
			} else {
				context.code = append(context.code, Instruction{"ILOAD_LOCAL", value.Slot})
			}
		}
		return value.Type

	case *model.NameType:
		return v.Name
	case *model.NameBool:
		context.code = append(context.code, Instruction{"IPUSH", BoolToInt(v.Name == "true")})
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
			context.code = append(context.code, Instruction{"IADD", nil})
			return "int"
			//return &WabbitValue{"int", left.Value.(int) + right.Value.(int)}
		} else if left == "float" && right == "float" {
			//return &WabbitValue{"float", left.Value.(float64) + right.Value.(float64)}
			context.code = append(context.code, Instruction{"FADD", nil})
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
			context.code = append(context.code, Instruction{"IMUL", nil})
			return "int"
			//return &WabbitValue{"int", left.Value.(int) * right.Value.(int)}
		} else if left == "float" && right == "float" {
			context.code = append(context.code, Instruction{"FMUL", nil})
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
			context.code = append(context.code, Instruction{"ISUB", nil})
			return "int"
			//return &WabbitValue{"int", left.Value.(int) - right.Value.(int)}
		} else if left == "float" && right == "float" {
			context.code = append(context.code, Instruction{"FSUB", nil})
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
			context.code = append(context.code, Instruction{"IDIV", nil})
			return "int"
		} else if left == "float" && right == "float" {
			context.code = append(context.code, Instruction{"FDIV", nil})
			return "float"
		} else {
			// we think it's a type error
			panic("type different")
		}
		return left

	case *model.Neg:
		right := InterpretNode(v.Operand, context)
		if right == "int" {
			//return &WabbitValue{"int", -right.Value.(int)}
			context.code = append(context.code, Instruction{"INEG", nil})
		} else if right == "float" {
			//return &WabbitValue{"float", -right.Value.(float64)}
			context.code = append(context.code, Instruction{"FNEG", nil})
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
			context.code = append(context.code, Instruction{"IPUSH", 1})
			context.code = append(context.code, Instruction{"XOR", nil})
		} else {
			// we think it's a type error
			panic("type different")
		}
		return right
	case *model.VarDeclaration:
		//var val *WVMVar
		var valtype string
		if v.Value != nil {
			valtype = InterpretNode(v.Value, context)
		} else {
			valtype = v.Type.Type()
			// the default value init
			if valtype == "float" {
				context.code = append(context.code, Instruction{"FPUSH", 0.0})
			} else {
				context.code = append(context.code, Instruction{"IPUSH", 0})
			}
		}
		scope, slot := context.NewVariable()
		context.Define(v.Name.Text, &WVMVar{Type: valtype, Scope: scope, Slot: slot}) // this is for context rember
		if scope == "global" {
			if valtype == "float" {
				context.code = append(context.code, Instruction{"FSTORE_GLOBAL", slot})
			} else {
				context.code = append(context.code, Instruction{"ISTORE_GLOBAL", slot})
			}
		} else if scope == "local" {
			if valtype == "float" {
				context.code = append(context.code, Instruction{"FSTORE_LOCAL", slot})
			} else {
				context.code = append(context.code, Instruction{"ISTORE_LOCAL", slot})
			}
		}
		return ""

	case *model.ConstDeclaration:
		valtype := InterpretNode(v.Value, context)
		scope, slot := context.NewVariable()
		context.Define(v.Name.Text, &WVMVar{Type: valtype, Scope: scope, Slot: slot}) // this is for context rember
		if scope == "global" {
			if valtype == "float" {
				context.code = append(context.code, Instruction{"FSTORE_GLOBAL", slot})
			} else {
				context.code = append(context.code, Instruction{"ISTORE_GLOBAL", slot})
			}
		} else if scope == "local" {
			if valtype == "float" {
				context.code = append(context.code, Instruction{"FSTORE_LOCAL", slot})
			} else {
				context.code = append(context.code, Instruction{"ISTORE_LOCAL", slot})
			}
		}
		return "bool"

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
