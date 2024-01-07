package vm

import (
	"dancavallaro.com/synacor-go/internal/util"
	"dancavallaro.com/synacor-go/pkg/memory"
	"dancavallaro.com/synacor-go/pkg/op"
	"errors"
	"fmt"
	"log"
	"time"
)

type opRef struct {
	opcode   uint16
	nArgs    int
	execute  func(registers *memory.Memory, args []uint16)
	mnemonic string
}

var ops = []opRef{
	{0, 0, op.Halt, "halt"},
	{1, 2, op.Set, "set"},
	{2, 1, op.Push, "push"},
	{3, 1, op.Pop, "pop"},
	{4, 3, op.Eq, "eq"},
	{5, 3, op.Gt, "gt"},
	{6, 1, op.Jmp, "jmp"},
	{7, 2, op.Jt, "jt"},
	{8, 2, op.Jf, "jf"},
	{9, 3, op.Add, "add"},
	{10, 3, op.Mult, "mult"},
	{11, 3, op.Mod, "mod"},
	{12, 3, op.And, "and"},
	{13, 3, op.Or, "or"},
	{14, 2, op.Not, "not"},
	{15, 2, op.Rmem, "rmem"},
	{16, 2, op.Wmem, "wmem"},
	{17, 1, op.Call, "call"},
	{18, 0, op.Ret, "ret"},
	{19, 1, op.Out, "out"},
	{20, 1, op.In, "in"},
	{21, 0, op.Noop, "noop"},
}
var opRefs = map[uint16]opRef{}

func init() {
	for _, o := range ops {
		opRefs[o.opcode] = o
	}
}

type ExecutionOptions struct {
	Trace      bool
	Delay      int
	Breakpoint int
}

type VM struct {
	M       memory.Memory
	Options *ExecutionOptions

	opsExecuted   int
	stepDebugging bool
}

func NewVM(bin []byte, opts *ExecutionOptions) *VM {
	vm := &VM{}
	vm.Options = opts
	loadProgram(&vm.M, bin)
	return vm
}

func (vm *VM) Step() error {
	pc := vm.M.PC // Save original PC value
	w := vm.M.ReadWord(uint16(vm.M.PC))
	vm.M.PC += 1
	o, ok := opRefs[w]
	if !ok {
		return errors.New(fmt.Sprintf("invalid opcode %d", w))
	}

	var args []uint16
	for arg := 1; arg <= o.nArgs; arg++ {
		w := vm.M.ReadWord(uint16(vm.M.PC))
		vm.M.PC += 1
		args = append(args, w)
	}

	if vm.Options.Breakpoint >= 0 && pc == vm.Options.Breakpoint {
		vm.Options.Trace = true
		vm.stepDebugging = true
	}
	if vm.Options.Trace {
		prefix := ""
		if vm.stepDebugging {
			prefix = "(step, r to resume) "
		}
		log.Printf("%s[PC=%d (0x%x)] %d (%s): %v", prefix, pc, pc, o.opcode, o.mnemonic, args)
	}
	if vm.stepDebugging {
		ch := util.ReadChar()
		if ch == 'r' {
			vm.stepDebugging = false
		}
	}

	o.execute(&vm.M, args)
	vm.opsExecuted++

	if vm.Options.Delay > 0 {
		time.Sleep(time.Duration(vm.Options.Delay) * time.Millisecond)
	}
	return nil
}

func (vm *VM) Execute() error {
	for vm.M.PC = 0; vm.M.PC < len(vm.M.Mem); {
		if err := vm.Step(); err != nil {
			return err
		}
	}
	fmt.Printf("Executed %d instructions\n", vm.opsExecuted)
	return nil
}

func loadProgram(r *memory.Memory, bin []byte) {
	for i, b := range bin {
		r.Mem[i] = uint16(b)
	}
}
