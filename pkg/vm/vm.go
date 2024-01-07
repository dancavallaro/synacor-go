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
	execute  func(registers *memory.Registers, args []uint16)
	mnemonic string
}

var ops = []opRef{
	{0, 0, op.Halt, "halt"},
	{1, 2, op.Set, "set"},
	{6, 1, op.Jmp, "jmp"},
	{7, 2, op.Jt, "jt"},
	{8, 2, op.Jf, "jf"},
	{9, 3, op.Add, "add"},
	{19, 1, op.Out, "out"},
	{21, 0, op.Noop, "noop"},
}
var opRefs = map[uint16]opRef{}

func init() {
	for _, o := range ops {
		opRefs[o.opcode] = o
	}
}

func readWord(bin []byte, address int) uint16 {
	return (uint16(bin[address+1]) << 8) + uint16(bin[address])
}

type ExecutionOptions struct {
	Trace bool
	Delay int
}

func Execute(bin []byte, opts *ExecutionOptions) error {
	r := memory.Registers{}
	for r.PC = 0; r.PC < len(bin)-1; {
		pc := r.PC
		w := readWord(bin, r.PC)
		o, ok := opRefs[w]
		if !ok {
			return errors.New(fmt.Sprintf("invalid opcode %d", w))
		}
		r.PC += 2

		var args []uint16
		for arg := 1; arg <= o.nArgs; arg++ {
			w := readWord(bin, r.PC)
			args = append(args, w)
			r.PC += 2
		}
		if opts.Trace {
			log.Printf("[PC=%d (0x%x)] %d (%s): %v\n", pc, pc, o.opcode, o.mnemonic, args)
		}
		o.execute(&r, args)

		if opts.Delay != -1 {
			time.Sleep(time.Duration(opts.Delay) * time.Millisecond)
		}
	}
	return nil
}
