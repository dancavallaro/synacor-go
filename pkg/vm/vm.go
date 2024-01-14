package vm

import (
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

	binary        []byte
	opsExecuted   int
	stepDebugging bool
}

func NewVM(bin []byte, opts *ExecutionOptions) *VM {
	vm := &VM{}
	vm.Options = opts
	vm.binary = bin
	vm.loadProgram()
	return vm
}

func (vm *VM) Step() error {
	pc := vm.M.PC // Save original PC value
	w := vm.M.Mem[vm.M.PC]
	vm.M.PC += 1
	o, ok := opRefs[w]
	if !ok {
		return errors.New(fmt.Sprintf("invalid opcode %d", w))
	}

	var args []uint16
	for arg := 1; arg <= o.nArgs; arg++ {
		w := vm.M.Mem[vm.M.PC]
		vm.M.PC += 1
		args = append(args, w)
	}

	if vm.Options.Breakpoint >= 0 && pc == vm.Options.Breakpoint {
		vm.Options.Trace = true
	}
	if vm.Options.Trace {
		log.Printf("[PC=%d (0x%x)] %d (%s): %v", pc, pc, o.opcode, o.mnemonic, args)
	}

	o.execute(&vm.M, args)
	vm.opsExecuted++

	if vm.Options.Delay > 0 {
		time.Sleep(time.Duration(vm.Options.Delay) * time.Millisecond)
	}
	return nil
}

func (vm *VM) Execute() error {
	for vm.M.PC < len(vm.M.Mem) {
		if err := vm.Step(); err != nil {
			return err
		}
	}
	fmt.Printf("Executed %d instructions\n", vm.opsExecuted)
	return nil
}

func (vm *VM) Restart() {
	vm.M = memory.Memory{}
	vm.loadProgram()
}

func (vm *VM) loadProgram() {
	for i := 0; i < len(vm.binary)-1; {
		low, high := uint16(vm.binary[i]), uint16(vm.binary[i+1])
		vm.M.Mem[i/2] = (high << 8) + low
		i += 2
	}
}
