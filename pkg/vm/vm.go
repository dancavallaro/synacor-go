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

func Execute(bin []byte, opts *ExecutionOptions) error {
	m := &memory.Memory{}
	loadProgram(m, bin)
	opsExecuted := 0
	stepDebugging := false
	for m.PC = 0; m.PC < len(bin)-1; {
		pc := m.PC // Save original PC value
		w := m.ReadWord(uint16(m.PC))
		m.PC += 1
		o, ok := opRefs[w]
		if !ok {
			return errors.New(fmt.Sprintf("invalid opcode %d", w))
		}

		var args []uint16
		for arg := 1; arg <= o.nArgs; arg++ {
			w := m.ReadWord(uint16(m.PC))
			m.PC += 1
			args = append(args, w)
		}

		if opts.Breakpoint >= 0 && pc == opts.Breakpoint {
			opts.Trace = true
			stepDebugging = true
		}
		if opts.Trace {
			prefix := ""
			if stepDebugging {
				prefix = "(step, r to resume) "
			}
			log.Printf("%s[PC=%d (0x%x)] %d (%s): %v", prefix, pc, pc, o.opcode, o.mnemonic, args)
		}
		if stepDebugging {
			ch := util.ReadChar()
			if ch == 'r' {
				stepDebugging = false
			}
		}

		o.execute(m, args)
		opsExecuted++

		if opts.Delay > 0 {
			time.Sleep(time.Duration(opts.Delay) * time.Millisecond)
		}
	}
	fmt.Printf("Executed %d instructions\n", opsExecuted)
	return nil
}

func loadProgram(r *memory.Memory, bin []byte) {
	for i, b := range bin {
		r.Mem[i] = uint16(b)
	}
}
