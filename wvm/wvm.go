package wvm

import (
	"fmt"
	"reflect"
	"strconv"
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
	labels  map[string]int
	frame   *Frame
	running bool
}

type Instruction struct {
	opcode string
	args   []interface{}
}

func (vm *WVM) run(instructions []Instruction) {
	vm.pc = 0
	vm.running = true
	vm.labels = make(map[string]int)
	for i, instruction := range instructions {
		if instruction.opcode == "LABEL" {
			vm.labels[instruction.args[0].(string)] = i
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
		argValues := make([]reflect.Value, len(args))
		for i, arg := range args {
			argValues[i] = reflect.ValueOf(arg)
		}
		method := reflect.ValueOf(vm).MethodByName(op)
		method.Call(argValues)
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
	return vm.istack[len(vm.istack)-1]
}

func (vm *WVM) DUP() {
	vm.istack = append(vm.istack, vm.istack[len(vm.istack)-1])
}

func (vm *WVM) IADD() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left + right)
}

func (vm *WVM) ISUB() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left - right)
}

func (vm *WVM) IMUL() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left * right)
}

func (vm *WVM) IDIV() {
	right := vm.IPOP()
	left := vm.IPOP()
	vm.IPUSH(left / right)
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

func (vm *WVM) FPUSH(value float64) {
	vm.fstack = append(vm.fstack, value)
}

func (vm *WVM) ITOF() {
	value := vm.IPOP()
	vm.FPUSH(float64(value))
}

func (vm *WVM) INEG() {
	vm.IPUSH(-vm.IPOP())
}

func (vm *WVM) FPOP() float64 {
	return vm.fstack[len(vm.fstack)-1]
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
	fmt.Print(rune(vm.IPOP()))
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

func (vm *WVM) FSTORE_GLOBAL(slot int) {
	vm.globals[slot] = vm.FPOP()
}

func (vm *WVM) CALL(label string) {
	vm.frame = &Frame{vm.pc, make(map[int]interface{}), vm.frame}
	vm.pc = vm.labels[label]
}

func (vm *WVM) RETURN() {
	vm.pc = vm.frame.returnPc
	vm.frame = vm.frame.prevFrame
}

type WVMContext struct {
	env       map[string]interface{}
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
		env:   make(map[string]interface{}),
		scope: "global",
		code:  make([]interface{}, 0),
	}
}

func (ctx *WVMContext) Define(name string, value interface{}) {
	ctx.env[name] = value
}

func (ctx *WVMContext) Lookup(name string) interface{} {
	return ctx.env[name]
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
	ctx.env = make(map[string]interface{})
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
	wvm.run(wctx.code)
	return nil
}

func InterpretNode(node model.Node, context *WVMContext) string {
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
		} else if left.Type == "char" {
			return &WabbitValue{Type: "bool", Value: left.Value.(rune) < right.Value.(rune)}
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
		} else if left.Type == "char" {
			return &WabbitValue{Type: "bool", Value: left.Value.(rune) <= right.Value.(rune)}
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
		} else if left.Type == "char" {
			return &WabbitValue{Type: "bool", Value: left.Value.(rune) > right.Value.(rune)}
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
		} else if left.Type == "char" {
			return &WabbitValue{Type: "bool", Value: left.Value.(rune) >= right.Value.(rune)}
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
		} else if left.Type == "char" {
			return &WabbitValue{Type: "bool", Value: left.Value.(rune) == right.Value.(rune)}
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
		} else if left.Type == "char" {
			return &WabbitValue{Type: "bool", Value: left.Value.(rune) != right.Value.(rune)}
		} else {
			return &WabbitValue{Type: "error", Value: "type error"}
		}
	case *model.LogOr:
		left := InterpretNode(v.Left, context)
		if left.Value.(bool) {
			return &WabbitValue{Type: "bool", Value: true}
		}
		right := InterpretNode(v.Right, context)
		return &WabbitValue{Type: "bool", Value: left.Value.(bool) || right.Value.(bool)}
	case *model.LogAnd:
		left := InterpretNode(v.Left, context)
		if !left.Value.(bool) {
			return &WabbitValue{Type: "bool", Value: false}
		}
		right := InterpretNode(v.Right, context)
		return &WabbitValue{Type: "bool", Value: left.Value.(bool) && right.Value.(bool)}
	case *model.Assignment:
		val := InterpretNode(v.Value, context)
		// assign the value to the name
		_ = context.Assign(v.Location.(*model.Name).Text, val)
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
		return InterpretNode(v.Expression, context)
		// should return a = {a , b , c} must return one
		// when c as ExpressionAsStatement

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
				} else if result.Type == "char" {
					return &WabbitValue{"int", int(result.Value.(rune))}
				}
				return result
			}
			if value.Value == "float" {
				result := InterpretNode(v.Arguments[0], context) // # type conversion must only have one args ?
				if result.Type == "int" {
					return &WabbitValue{"float", float64(result.Value.(int))}
				} else if result.Type == "char" {
					return &WabbitValue{"float", float64(result.Value.(rune))}
				}
				return result
			}
			if value.Value == "char" {
				result := InterpretNode(v.Arguments[0], context) // # type conversion must only have one args ?
				if result.Type == "int" {
					return &WabbitValue{"char", rune(result.Value.(int))}
				} else if result.Type == "float64" {
					return &WabbitValue{"float", rune(result.Value.(float64))}
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
	case *model.CompoundExpression:
		return InterpretNode(&v.Statements, context)
	default:
		panic(fmt.Sprintf("Can't intepre %#v to source", v))
	}

	return nil
}
