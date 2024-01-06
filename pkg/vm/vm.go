package vm

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"dancavallaro.com/synacor-go/pkg/op"
	"errors"
	"fmt"
	"log"
)

type opRef struct {
	opcode   int16
	nArgs    int
	execute  func(args []int16, registers *memory.Registers)
	mnemonic string
}

var ops = []opRef{
	{0, 0, op.Halt, "halt"},
	{6, 1, op.Jmp, "jmp"},
	{19, 1, op.Out, "out"},
	{21, 0, op.Noop, "noop"},
}
var opRefs = map[int16]opRef{}

func init() {
	for _, o := range ops {
		opRefs[o.opcode] = o
	}
}

func readWord(bin []byte, address int) int16 {
	return (int16(bin[address+1]) << 8) + int16(bin[address])
}

func Execute(bin []byte, trace bool) error {
	r := memory.Registers{}
	for r.PC = 0; r.PC < len(bin)-1; {
		w := readWord(bin, r.PC)
		o, ok := opRefs[w]
		if !ok {
			return errors.New(fmt.Sprintf("invalid opcode %d", w))
		}

		var args []int16
		for arg := 1; arg <= o.nArgs; arg++ {
			w := readWord(bin, r.PC+1+arg)
			args = append(args, w)
			r.PC += 2
		}
		if trace {
			log.Printf("%d (%s): %v\n", o.opcode, o.mnemonic, args)
		}
		r.PC += 2
		o.execute(args, &r)
	}
	return nil
}
