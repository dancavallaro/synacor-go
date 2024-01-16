package vm

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"dancavallaro.com/synacor-go/pkg/op"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

type OpRef struct {
	opcode   uint16
	nArgs    int
	execute  func(registers *memory.Memory, args []uint16)
	Mnemonic string
}

var ops = []OpRef{
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
var opRefs = map[uint16]OpRef{}

func init() {
	for _, o := range ops {
		opRefs[o.opcode] = o
	}
}

type ExecutionOptions struct {
	Trace bool
	Delay int
}

type VM struct {
	M          memory.Memory
	OriginalPC int
	Options    *ExecutionOptions

	binary      []byte
	opsExecuted int
	lock        sync.Mutex
}

func NewVM(bin []byte, opts *ExecutionOptions) *VM {
	vm := &VM{}
	vm.Options = opts
	vm.binary = bin
	vm.loadProgram()
	return vm
}

func (vm *VM) DecodeOp(start int) (OpRef, []uint16, error) {
	w := vm.M.Mem[start]
	o, ok := opRefs[w]
	if !ok {
		return OpRef{}, nil, errors.New(fmt.Sprintf("invalid opcode %d", w))
	}

	var args []uint16
	for arg := 1; arg <= o.nArgs; arg++ {
		w := vm.M.Mem[start+arg]
		args = append(args, w)
	}

	return o, args, nil
}

func (vm *VM) Step() error {
	vm.lock.Lock()
	defer vm.lock.Unlock()

	vm.OriginalPC = vm.M.PC // Save original PC value
	o, args, err := vm.DecodeOp(vm.M.PC)
	if err != nil {
		return err
	}
	vm.M.PC += 1 + len(args)

	if vm.Options.Trace {
		log.Printf("[PC=%d (0x%x)] %d (%s): %v", vm.OriginalPC, vm.OriginalPC, o.opcode, o.Mnemonic, args)
	}

	o.execute(&vm.M, args)
	vm.opsExecuted++
	vm.OriginalPC = vm.M.PC // Now update the display value with the actual PC which might have been updated

	return nil
}

func (vm *VM) Execute() error {
	for vm.M.PC < len(vm.M.Mem) {
		if err := vm.Step(); err != nil {
			return err
		}

		if vm.Options.Delay > 0 {
			time.Sleep(time.Duration(vm.Options.Delay) * time.Millisecond)
		}
	}
	fmt.Printf("Executed %d instructions\n", vm.opsExecuted)
	return nil
}

func (vm *VM) Restart() {
	vm.lock.Lock()
	defer vm.lock.Unlock()

	vm.M = memory.Memory{}
	vm.OriginalPC = 0
	vm.opsExecuted = 0
	vm.loadProgram()
}

func (vm *VM) OpsExecuted() int {
	return vm.opsExecuted
}

func (vm *VM) loadProgram() {
	for i := 0; i < len(vm.binary)-1; {
		low, high := uint16(vm.binary[i]), uint16(vm.binary[i+1])
		vm.M.Mem[i/2] = (high << 8) + low
		i += 2
	}
}
