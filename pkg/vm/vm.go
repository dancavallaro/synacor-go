package vm

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"dancavallaro.com/synacor-go/pkg/op"
	"errors"
	"fmt"
	"log"
)

type opRef struct {
	opcode   uint16
	nArgs    int
	execute  func(registers *memory.Registers, args []uint16)
	mnemonic string
}

var ops = []opRef{
	{0, 0, op.Halt, "halt"},
	{6, 1, op.Jmp, "jmp"},
	{7, 2, op.Jt, "jt"},
	{8, 2, op.Jf, "jf"},
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

func Execute(bin []byte, trace bool) error {
	r := memory.Registers{}
	for r.PC = 0; r.PC < len(bin)-1; {
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
		if trace {
			log.Printf("%d (%s): %v\n", o.opcode, o.mnemonic, args)
		}
		o.execute(&r, args)
	}
	return nil
}
