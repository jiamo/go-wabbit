package wvm

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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

type OpFunc func(args interface{}) interface{}

type Instruction struct {
	opcode string
	args   interface{}
}

func NewWVM() *WVM {
	wvm := WVM{
		globals: make(map[int]interface{}),
	}

	return &wvm
}

// At now we don't have bytecode just instruction
func (vm *WVM) run(instructions []Instruction) {
	vm.pc = 0
	vm.running = true

	// Should update Instruction to Thread code
	opMap := vm.getOpcodeMap()

	for vm.running {
		op := instructions[vm.pc].opcode
		args := instructions[vm.pc].args
		vm.pc++

		if op == "LABEL" {
			continue
		}
		fn, exists := opMap[op]
		if !exists {
			// 错误处理：未知的字节码
			panic(fmt.Sprintf("no such opcode %v", op))
		}
		fn(args)

	}
}

func (vm *WVM) IPUSH(value interface{}) interface{} {
	vm.istack = append(vm.istack, value.(int)) // reflect vs cast
	return nil
}

func (vm *WVM) IPOP(value interface{}) interface{} {
	index := len(vm.istack) - 1
	// Get the top element of the stack
	element := vm.istack[index]
	// Remove the top element of the stack
	vm.istack = vm.istack[:index]
	return element
}

func (vm *WVM) IDUP(value interface{}) interface{} {
	vm.istack = append(vm.istack, vm.istack[len(vm.istack)-1])
	return nil
}

func (vm *WVM) FDUP(value interface{}) interface{} {
	vm.fstack = append(vm.fstack, vm.fstack[len(vm.fstack)-1])
	return nil
}

func (vm *WVM) IADD(value interface{}) interface{} {
	right := (vm.IPOP(nil)).(int)
	left := (vm.IPOP(nil)).(int)
	vm.IPUSH(left + right)
	return nil
}

func (vm *WVM) FADD(value interface{}) interface{} {
	right := (vm.FPOP(nil)).(float64)
	left := (vm.FPOP(nil)).(float64)
	vm.FPUSH(left + right)
	return nil
}

func (vm *WVM) ISUB(value interface{}) interface{} {
	right := (vm.IPOP(nil)).(int)
	left := (vm.IPOP(nil)).(int)
	vm.IPUSH(left - right)
	return nil
}

func (vm *WVM) FSUB(value interface{}) interface{} {
	right := (vm.FPOP(nil)).(float64)
	left := (vm.FPOP(nil)).(float64)
	vm.FPUSH(left - right)
	return nil
}

func (vm *WVM) IMUL(value interface{}) interface{} {
	right := (vm.IPOP(nil)).(int)
	left := (vm.IPOP(nil)).(int)
	vm.IPUSH(left * right)
	return nil
}

func (vm *WVM) FMUL(value interface{}) interface{} {
	right := (vm.FPOP(nil)).(float64)
	left := (vm.FPOP(nil)).(float64)
	vm.FPUSH(left * right)
	return nil
}

func (vm *WVM) IDIV(value interface{}) interface{} {
	right := (vm.IPOP(nil)).(int)
	left := (vm.IPOP(nil)).(int)
	vm.IPUSH(left / right)
	return nil
}

func (vm *WVM) FDIV(value interface{}) interface{} {
	right := (vm.FPOP(nil)).(float64)
	left := (vm.FPOP(nil)).(float64)
	vm.FPUSH(left / right)
	return nil
}

func (vm *WVM) AND(value interface{}) interface{} {
	right := (vm.IPOP(nil)).(int)
	left := (vm.IPOP(nil)).(int)
	vm.IPUSH(left & right)
	return nil
}

func (vm *WVM) OR(value interface{}) interface{} {
	right := vm.IPOP(nil).(int)
	left := vm.IPOP(nil).(int)
	vm.IPUSH(left | right)
	return nil
}

func (vm *WVM) XOR(value interface{}) interface{} {
	right := (vm.IPOP(nil)).(int)
	left := (vm.IPOP(nil)).(int)
	vm.IPUSH(left ^ right)
	return nil
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (vm *WVM) ICMP(value interface{}) interface{} {
	right := (vm.IPOP(nil)).(int)
	left := (vm.IPOP(nil)).(int)
	op := value.(string)
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
	return nil
}

func (vm *WVM) FCMP(value interface{}) interface{} {
	right := (vm.FPOP(nil)).(float64)
	left := (vm.FPOP(nil)).(float64)
	op := value.(string)
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
	return nil
}

func (vm *WVM) FPUSH(value interface{}) interface{} {
	vm.fstack = append(vm.fstack, value.(float64))
	return nil
}

func (vm *WVM) INEG(value interface{}) interface{} {
	vm.IPUSH(-(vm.IPOP(nil).(int)))
	return nil
}

func (vm *WVM) FPOP(value interface{}) interface{} {
	index := len(vm.fstack) - 1
	// Get the top element of the stack
	element := vm.fstack[index]
	// Remove the top element of the stack
	vm.fstack = vm.fstack[:index]
	return element
}

func (vm *WVM) FNEG(value interface{}) interface{} {
	vm.FPUSH(-(vm.FPOP(nil)).(float64))
	return nil
}

func (vm *WVM) PRINTI(value interface{}) interface{} {
	fmt.Println((vm.IPOP(nil)).(int))
	return nil
}

func (vm *WVM) PRINTF(value interface{}) interface{} {
	fmt.Println((vm.FPOP(nil)).(float64))
	return nil
}

func (vm *WVM) PRINTB(value interface{}) interface{} {
	if (vm.IPOP(nil)).(int) == 0 {
		fmt.Println("false")
	} else {
		fmt.Println("true")
	}
	return nil
}

func (vm *WVM) PRINTC(value interface{}) interface{} {
	fmt.Printf("%c", rune((vm.IPOP(nil)).(int)))
	return nil
}

func (vm *WVM) FTOI(value interface{}) interface{} {
	vm.IPUSH(int((vm.FPOP(nil)).(float64)))
	return nil
}

func (vm *WVM) ITOF(value interface{}) interface{} {
	vm.FPUSH(float64((vm.IPOP(nil)).(int)))
	return nil
}

func (vm *WVM) ISTORE_LOCAL(value interface{}) interface{} {
	slot := value.(int)
	vm.frame.locals[slot] = vm.IPOP(nil).(int)
	return nil
}

func (vm *WVM) ILOAD_LOCAL(value interface{}) interface{} {
	slot := value.(int)
	vm.IPUSH(vm.frame.locals[slot].(int))
	return nil
}

func (vm *WVM) FLOAD_LOCAL(value interface{}) interface{} {
	slot := value.(int)
	vm.FPUSH(vm.frame.locals[slot].(float64))
	return nil
}

func (vm *WVM) FSTORE_LOCAL(value interface{}) interface{} {
	slot := value.(int)
	vm.frame.locals[slot] = vm.FPOP(nil).(float64)
	return nil
}

func (vm *WVM) ISTORE_GLOBAL(value interface{}) interface{} {
	slot := value.(int)
	vm.globals[slot] = vm.IPOP(nil).(int)
	return nil
}

func (vm *WVM) ILOAD_GLOBAL(value interface{}) interface{} {
	slot := value.(int)
	vm.IPUSH(vm.globals[slot].(int))
	return nil
}

func (vm *WVM) FLOAD_GLOBAL(value interface{}) interface{} {
	slot := value.(int)
	vm.FPUSH(vm.globals[slot].(float64))
	return nil
}

func (vm *WVM) GOTO(value interface{}) interface{} {
	name := value.(int)
	vm.pc = vm.labels[name]
	return nil
}

func (vm *WVM) BZ(value interface{}) interface{} {
	name := value.(int)
	if vm.IPOP(nil).(int) == 0 {
		vm.pc = vm.labels[name]
	}
	return nil
}

func (vm *WVM) HALT(value interface{}) interface{} {
	vm.running = false
	return nil
}

func (vm *WVM) LABEL(value interface{}) interface{} {
	return nil
}

func (vm *WVM) FSTORE_GLOBAL(value interface{}) interface{} {
	slot := value.(int)
	vm.globals[slot] = vm.FPOP(nil).(float64)
	return nil
}

func (vm *WVM) CALL(value interface{}) interface{} {
	label := value.(int)
	vm.frame = &Frame{vm.pc, make(map[int]interface{}), vm.frame}
	vm.pc = vm.labels[label]
	return nil
}

func (vm *WVM) RETURN(value interface{}) interface{} {
	vm.pc = vm.frame.returnPc
	vm.frame = vm.frame.prevFrame
	return nil
}

func (vm *WVM) getOpcodeMap() map[string]OpFunc {
	return map[string]OpFunc{
		"IPUSH":         vm.IPUSH,
		"IPOP":          vm.IPOP,
		"IDUP":          vm.IDUP,
		"FDUP":          vm.FDUP,
		"IADD":          vm.IADD,
		"FADD":          vm.FADD,
		"ISUB":          vm.ISUB,
		"FSUB":          vm.FSUB,
		"IMUL":          vm.IMUL,
		"FMUL":          vm.FMUL,
		"IDIV":          vm.IDIV,
		"FDIV":          vm.FDIV,
		"AND":           vm.AND,
		"OR":            vm.OR,
		"XOR":           vm.XOR,
		"ICMP":          vm.ICMP,
		"FCMP":          vm.FCMP,
		"FPUSH":         vm.FPUSH,
		"INEG":          vm.INEG,
		"FPOP":          vm.FPOP,
		"FNEG":          vm.FNEG,
		"PRINTI":        vm.PRINTI,
		"PRINTF":        vm.PRINTF,
		"PRINTB":        vm.PRINTB,
		"PRINTC":        vm.PRINTC,
		"FTOI":          vm.FTOI,
		"ITOF":          vm.ITOF,
		"ISTORE_LOCAL":  vm.ISTORE_LOCAL,
		"ILOAD_LOCAL":   vm.ILOAD_LOCAL,
		"FLOAD_LOCAL":   vm.FLOAD_LOCAL,
		"FSTORE_LOCAL":  vm.FSTORE_LOCAL,
		"ISTORE_GLOBAL": vm.ISTORE_GLOBAL,
		"ILOAD_GLOBAL":  vm.ILOAD_GLOBAL,
		"FLOAD_GLOBAL":  vm.FLOAD_GLOBAL,
		"GOTO":          vm.GOTO,
		"BZ":            vm.BZ,
		"HALT":          vm.HALT,
		"LABEL":         vm.LABEL,
		"FSTORE_GLOBAL": vm.FSTORE_GLOBAL,
		"CALL":          vm.CALL,
		"RETURN":        vm.RETURN,
	}
}

type Context struct {
	env       *common.ChainMap
	code      []Instruction
	labels    map[int]int
	nglobals  int
	nlocals   int
	nlabels   int
	scope     string
	haveMain  bool
	parentEnv *map[string]interface{}
}

func NewWVMContext() *Context {
	return &Context{
		env:    common.NewChainMap(),
		scope:  "global",
		code:   make([]Instruction, 0),
		labels: make(map[int]int),
	}
}

type WVMVar struct {
	Type  string
	Scope string
	Slot  int
}

func (ctx *Context) Define(name string, value *WVMVar) {
	ctx.env.SetValue(name, value)
}

func (ctx *Context) Lookup(name string) *WVMVar {
	v, e := ctx.env.GetValue(name)
	if e == true {
		return v.(*WVMVar)
	} else {
		return nil
	}
}

func (ctx *Context) NewVariable() (string, int) {
	if ctx.scope == "global" {
		ctx.nglobals++
		return "global", ctx.nglobals - 1
	} else {
		ctx.nlocals++
		return "local", ctx.nlocals - 1
	}
}

func (ctx *Context) NewLabel() int {
	ctx.nlabels++
	return ctx.nlabels - 1
}

func (ctx *Context) NewScope(do func()) {
	oldEnv := ctx.env
	ctx.env = ctx.env.NewChild()
	defer func() {
		ctx.env = oldEnv
	}()
	do()
}

func (ctx *Context) NewInstruction(instruction Instruction) {
	ctx.code = append(ctx.code, instruction)
	if instruction.opcode == "LABEL" {
		ctx.labels[instruction.args.(int)] = len(ctx.code) - 1
	}
}

func Wvm(program *model.Program) error {
	wctx := NewWVMContext()
	_ = InterpretNode(program.Model, wctx) // generate is InterpretNode in the same meaning
	wvm := &WVM{
		globals: make(map[int]interface{}),
	}

	//wctx.code = append(wctx.code, Instruction{"HALT", nil})
	wctx.NewInstruction(Instruction{"HALT", nil})
	log.Debug(wctx.code)
	wvm.labels = wctx.labels
	wvm.run(wctx.code)

	return nil
}

func InterpretNode(node model.Node, context *Context) string {
	switch v := node.(type) {
	case *model.Integer:
		context.NewInstruction(Instruction{"IPUSH", v.Value})
		return "int"
	case *model.Float:
		context.NewInstruction(Instruction{"FPUSH", v.Value})
		return "float"
	case *model.Character:
		unquoted, err := strconv.Unquote(v.Value)
		if err != nil {
			panic(err)
		}
		context.NewInstruction(Instruction{"IPUSH", int(rune(unquoted[0]))})
		return "char"
	case *model.Name:
		log.Debugf("*model.Name %v env %v code %v", v.Text, context.env, context.code)
		value := context.Lookup(v.Text) // somethings we may need exist
		// bool int char using int
		if value.Scope == "global" {
			if value.Type == "float" {
				context.NewInstruction(Instruction{"FLOAD_GLOBAL", value.Slot})
			} else {
				context.NewInstruction(Instruction{"ILOAD_GLOBAL", value.Slot})
			}
		} else if value.Scope == "local" {
			if value.Type == "float" {
				context.NewInstruction(Instruction{"FLOAD_LOCAL", value.Slot})
			} else {
				context.NewInstruction(Instruction{"ILOAD_LOCAL", value.Slot})
			}
		}
		return value.Type

	case *model.NameType:
		return v.Name
	case *model.NameBool:
		context.NewInstruction(Instruction{"IPUSH", BoolToInt(v.Name == "true")})
		return "bool"
	case *model.IntegerType:
		return "int"
	case *model.FloatType:
		return "float"
	case *model.Add:
		left := InterpretNode(v.Left, context)
		right := InterpretNode(v.Right, context)

		// we should check the type of left and right go we can't make interface + interface
		if left == "int" && right == "int" {
			context.NewInstruction(Instruction{"IADD", nil})
			return "int"
		} else if left == "float" && right == "float" {
			context.NewInstruction(Instruction{"FADD", nil})
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
			context.NewInstruction(Instruction{"IMUL", nil})
			return "int"
		} else if left == "float" && right == "float" {
			context.NewInstruction(Instruction{"FMUL", nil})
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
			context.NewInstruction(Instruction{"ISUB", nil})
			return "int"
		} else if left == "float" && right == "float" {
			context.NewInstruction(Instruction{"FSUB", nil})
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
			context.NewInstruction(Instruction{"IDIV", nil})
			return "int"
		} else if left == "float" && right == "float" {
			context.NewInstruction(Instruction{"FDIV", nil})
			return "float"
		} else {
			// we think it's a type error
			panic("type different")
		}
		return left

	case *model.Neg:
		right := InterpretNode(v.Operand, context)
		if right == "int" {
			context.NewInstruction(Instruction{"INEG", nil})
		} else if right == "float" {
			context.NewInstruction(Instruction{"FNEG", nil})
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
			context.NewInstruction(Instruction{"IPUSH", 1})
			context.NewInstruction(Instruction{"XOR", nil})
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
				context.NewInstruction(Instruction{"FPUSH", 0.0})
			} else {
				context.NewInstruction(Instruction{"IPUSH", 0})
			}
		}
		scope, slot := context.NewVariable()
		context.Define(v.Name.Text, &WVMVar{Type: valtype, Scope: scope, Slot: slot}) // this is for context rember
		if scope == "global" {
			if valtype == "float" {
				context.NewInstruction(Instruction{"FSTORE_GLOBAL", slot})
			} else {
				context.NewInstruction(Instruction{"ISTORE_GLOBAL", slot})
			}
		} else if scope == "local" {
			if valtype == "float" {
				context.NewInstruction(Instruction{"FSTORE_LOCAL", slot})
			} else {
				context.NewInstruction(Instruction{"ISTORE_LOCAL", slot})
			}
		}
		return ""

	case *model.ConstDeclaration:
		valtype := InterpretNode(v.Value, context)
		scope, slot := context.NewVariable()
		context.Define(v.Name.Text, &WVMVar{Type: valtype, Scope: scope, Slot: slot}) // this is for context rember
		if scope == "global" {
			if valtype == "float" {
				context.NewInstruction(Instruction{"FSTORE_GLOBAL", slot})
			} else {
				context.NewInstruction(Instruction{"ISTORE_GLOBAL", slot})
			}
		} else if scope == "local" {
			if valtype == "float" {
				context.NewInstruction(Instruction{"FSTORE_LOCAL", slot})
			} else {
				context.NewInstruction(Instruction{"ISTORE_LOCAL", slot})
			}
		}
		return ""

	case *model.Lt:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			context.NewInstruction(Instruction{"ICMP", "<"})
		} else if left == "float" {
			context.NewInstruction(Instruction{"FCMP", "<"})
		} else if left == "char" {
			context.NewInstruction(Instruction{"ICMP", "<"})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.Le:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			context.NewInstruction(Instruction{"ICMP", "<="})
		} else if left == "float" {
			context.NewInstruction(Instruction{"FCMP", "<="})
		} else if left == "char" {
			context.NewInstruction(Instruction{"ICMP", "<="})
		} else if left == "bool" {
			context.NewInstruction(Instruction{"ICMP", "<="})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.Gt:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			context.NewInstruction(Instruction{"ICMP", ">"})
		} else if left == "float" {
			context.NewInstruction(Instruction{"FCMP", ">"})
		} else if left == "char" {
			context.NewInstruction(Instruction{"ICMP", ">"})
		} else if left == "bool" {
			context.NewInstruction(Instruction{"ICMP", ">"})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.Ge:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			context.NewInstruction(Instruction{"ICMP", ">="})
		} else if left == "float" {
			context.NewInstruction(Instruction{"FCMP", ">="})
		} else if left == "char" {
			context.NewInstruction(Instruction{"ICMP", ">="})
		} else if left == "bool" {
			context.NewInstruction(Instruction{"ICMP", ">="})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.Eq:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			context.NewInstruction(Instruction{"ICMP", "=="})
		} else if left == "float" {
			context.NewInstruction(Instruction{"FCMP", "=="})
		} else if left == "char" {
			context.NewInstruction(Instruction{"ICMP", "=="})
		} else if left == "bool" {
			context.NewInstruction(Instruction{"ICMP", "=="})
		} else {
			panic("type different") // we can just using else For simple
		}
		return "bool"
	case *model.Ne:
		left := InterpretNode(v.Left, context)
		_ = InterpretNode(v.Right, context)
		if left == "int" {
			context.NewInstruction(Instruction{"ICMP", "!="})
		} else if left == "float" {
			context.NewInstruction(Instruction{"FCMP", "!="})
		} else if left == "char" {
			context.NewInstruction(Instruction{"ICMP", "!="})
		} else if left == "bool" {
			context.NewInstruction(Instruction{"ICMP", "!="})
		} else {
			panic("type different")
		}
		return "bool"
	case *model.LogOr:
		// TODO short eval
		done_label := context.NewLabel()
		or_continue_label := context.NewLabel()

		_ = InterpretNode(v.Left, context)
		context.NewInstruction(Instruction{"BZ", or_continue_label})
		context.NewInstruction(Instruction{"IPUSH", 1})
		context.NewInstruction(Instruction{"GOTO", done_label})
		context.NewInstruction(Instruction{"LABEL", or_continue_label})

		_ = InterpretNode(v.Right, context)
		context.NewInstruction(Instruction{"LABEL", done_label})
		return "bool"

	case *model.LogAnd:
		done_label := context.NewLabel()
		and_false_label := context.NewLabel()
		_ = InterpretNode(v.Left, context)
		context.NewInstruction(Instruction{"BZ", and_false_label})

		_ = InterpretNode(v.Right, context)
		context.NewInstruction(Instruction{"GOTO", done_label})
		context.NewInstruction(Instruction{"LABEL", and_false_label})
		context.NewInstruction(Instruction{"IPUSH", 0})
		context.NewInstruction(Instruction{"LABEL", done_label})

		return "bool" // no need or and any more

	case *model.Assignment:
		val := InterpretNode(v.Value, context)
		// assign the value to the name
		wvmvar := context.Lookup(v.Location.(*model.Name).Text)
		if wvmvar.Scope == "global" {
			if wvmvar.Type == "float" {
				context.NewInstruction(Instruction{"FDUP", nil})
				context.NewInstruction(Instruction{"FSTORE_GLOBAL", wvmvar.Slot})
			} else {
				context.NewInstruction(Instruction{"IDUP", nil})
				context.NewInstruction(Instruction{"ISTORE_GLOBAL", wvmvar.Slot})
			}
		} else {
			if wvmvar.Type == "float" {
				context.NewInstruction(Instruction{"FDUP", nil})
				context.NewInstruction(Instruction{"FSTORE_LOCAL", wvmvar.Slot})
			} else {
				context.NewInstruction(Instruction{"IDUP", nil})
				context.NewInstruction(Instruction{"ISTORE_LOCAL", wvmvar.Slot})
			}
		}
		return val

	case *model.PrintStatement:
		value := InterpretNode(v.Value, context)
		switch value {
		case "char":
			context.NewInstruction(Instruction{"PRINTC", nil})
		case "bool":
			context.NewInstruction(Instruction{"PRINTB", nil})
		case "int":
			context.NewInstruction(Instruction{"PRINTI", nil})
		case "float":
			context.NewInstruction(Instruction{"PRINTF", nil})
		default:
			log.Debugf("%v:%v", context.code, value)
			panic("wrong type")
		}
	case *model.Statements:

		var result string
		for _, statement := range v.Statements {
			// do we need pop to keep stack blance?
			if result == "float" {
				context.NewInstruction(Instruction{"FPOP", nil})
			} else if result == "int" {
				context.NewInstruction(Instruction{"IPOP", nil})
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

		context.NewInstruction(Instruction{"BZ", else_label})
		context.NewInstruction(Instruction{"GOTO", then_label})
		context.NewInstruction(Instruction{"LABEL", then_label})

		context.NewScope(
			func() {
				InterpretNode(&v.Consequence, context)
				context.NewInstruction(Instruction{"GOTO", merge_label})
				context.NewInstruction(Instruction{"LABEL", else_label})
			},
		)
		if v.Alternative != nil {
			context.NewScope(
				func() {
					InterpretNode(v.Alternative, context)
					context.NewInstruction(Instruction{"GOTO", merge_label})
				},
			)
		}
		context.NewInstruction(Instruction{"GOTO", merge_label})
		context.NewInstruction(Instruction{"LABEL", merge_label})
	case *model.BreakStatement:
		// we need scope for level break
		val := context.Lookup("break") // fake using type as label
		context.NewInstruction(Instruction{"GOTO", val.Slot})
	case *model.ContinueStatement:
		val := context.Lookup("continue") // fake using type as label
		context.NewInstruction(Instruction{"GOTO", val.Slot})

	case *model.ReturnStatement:
		value := InterpretNode(v.Value, context)
		context.NewInstruction(Instruction{"RETURN", nil})
		return value

	case *model.WhileStatement:
		test_label := context.NewLabel()
		body_label := context.NewLabel()
		exit_label := context.NewLabel()

		context.NewInstruction(Instruction{"GOTO", test_label})
		context.NewInstruction(Instruction{"LABEL", test_label})
		InterpretNode(v.Test, context)
		context.NewInstruction(Instruction{"BZ", exit_label})
		context.NewInstruction(Instruction{"GOTO", body_label})
		context.NewInstruction(Instruction{"LABEL", body_label})

		context.NewScope(func() {
			context.Define("break", &WVMVar{"", "", exit_label})
			context.Define("continue", &WVMVar{"", "", test_label})
			InterpretNode(&v.Body, context)
			context.NewInstruction(Instruction{"GOTO", test_label})

		})
		context.NewInstruction(Instruction{"LABEL", exit_label})

	case *model.FunctionDeclaration:

		start_label := context.NewLabel()
		end_label := context.NewLabel()

		context.NewInstruction(Instruction{"GOTO", end_label})
		context.NewInstruction(Instruction{"LABEL", start_label})

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
					context.NewInstruction(Instruction{"FSTORE_LOCAL", val.Slot})
				} else if val.Type == "int" {
					context.NewInstruction(Instruction{"ISTORE_LOCAL", val.Slot})
				}
			}
			InterpretNode(&v.Body, context)
		})
		context.NewInstruction(Instruction{"LABEL", end_label})
		context.scope = "global"
		if v.Name.Text == "main" {
			context.haveMain = true
		}

	case *model.FunctionApplication:
		argType := "int"
		//value := InterpretNode(v.Func, context) // while lookup
		for _, arg := range v.Arguments {
			argType = InterpretNode(arg, context) // arg eval in current context
		}

		name := v.Func.(*model.Name).Text
		funcVar := context.Lookup(v.Func.(*model.Name).Text) // define in
		log.Debugf("name %v", name)
		if name == "int" {
			// only float need to cast
			if argType == "float" {
				context.NewInstruction(Instruction{"FTOI", nil})
			}
			return "int"
		}
		if name == "float" {
			if argType != "float" {
				context.NewInstruction(Instruction{"ITOF", nil})
			}
			return "float"
		}
		if name == "char" {
			return "char"
		}
		if name == "bool" {
			return "bool"
		}
		context.NewInstruction(Instruction{"CALL", funcVar.Slot})
		return funcVar.Type

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
