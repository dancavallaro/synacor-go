package debugger

import (
	"dancavallaro.com/synacor-go/pkg/env"
	"dancavallaro.com/synacor-go/pkg/memory"
	"dancavallaro.com/synacor-go/pkg/vm"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
	"strings"
)

type View struct {
	*gocui.View
}

func (v *View) Print(a ...any) {
	if _, err := fmt.Fprint(v, a...); err != nil {
		panic(err)
	}
}

func (v *View) Printf(format string, a ...any) {
	if _, err := fmt.Fprintf(v, format, a...); err != nil {
		panic(err)
	}
}

func (v *View) Println(a ...any) {
	if _, err := fmt.Fprintln(v, a...); err != nil {
		panic(err)
	}
}

type Frame interface {
	Init(v *View)
	Draw(v *View)
}

type HelpView struct{}

func (h HelpView) Init(v *View) {
	v.Title = "Help"
	v.Print("(^p) pause execution\t(^r) resume execution\t(^s) step execution\t" +
		"(^b) run until breakpoint\t(^x) toggle hex/dec\t(^\\) reset state and restart")
}

func (h HelpView) Draw(_ *View) {}

type RegisterView struct {
	m  *memory.Memory
	pc *int
	b  *base
}

func (r RegisterView) Init(v *View) {
	v.Title = "Registers"
}

func (r RegisterView) Draw(v *View) {
	v.Clear()
	gp := r.m.GP

	v.Printf("PC: %s\t", r.b.strSym(*r.pc))
	for i := 0; i < memory.NumRegisters; i++ {
		v.Printf("R%d: %s\t", i, r.b.strSym(int(gp[i])))
	}
}

type StackView struct {
	m *memory.Memory
	b *base
}

func (s StackView) Init(v *View) {
	v.Title = "Stack"
}

func (s StackView) Draw(v *View) {
	v.Clear()

	for i := 0; i < len(s.m.Stack); i++ {
		v.Print(s.b.strSym(int(s.m.Stack[i])))
		if i < len(s.m.Stack)-1 {
			v.Printf("\t")
		} else {
			v.Printf(" ")
		}
	}

	v.Print("◄SP")
}

type OutputView struct{}

func (o OutputView) Init(v *View) {
	v.Title = "Output"
	v.Autoscroll = true
	env.Config.Output = v
}

func (o OutputView) Draw(_ *View) {}

type LogView struct{}

func (l LogView) Init(v *View) {
	v.Title = "System Log"
	v.Autoscroll = true
	log.Default().SetOutput(v)
}

func (l LogView) Draw(_ *View) {}

type DisassemblyView struct {
	d *Debugger
	b *base
}

func (d DisassemblyView) Init(v *View) {
	v.Title = "Disassembly"
}

func numStr(num uint16, b base) string {
	if num >= memory.RegOffset && num < memory.RegOffset+memory.NumRegisters {
		regNum := num - memory.RegOffset
		return fmt.Sprintf("R%d", regNum)
	} else if isChar(num) || num == '\n' {
		return fmt.Sprintf("'%s'", str(rune(num)))
	} else {
		return b.strSym(int(num))
	}
}

func argStr(args []uint16, b base) string {
	if len(args) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(" ")
	for _, arg := range args {
		sb.WriteString(fmt.Sprintf("%v ", numStr(arg, b)))
	}
	return sb.String()
}

func (d DisassemblyView) Draw(v *View) {
	v.Clear()
	v.Println()
	pc := d.d.VM.OriginalPC
	_, y := v.Size()
	for line := 0; line < y-2; line++ {
		o, args, err := d.d.VM.DecodeOp(pc)
		if err != nil {
			panic(err)
		}
		v.Printf("    %s: %s%s\n", d.b.strSym(pc), o.Mnemonic, argStr(args, *d.b))
		pc += 1 + len(args)
	}
}

type StateView struct {
	state *State
}

func (s StateView) Init(_ *View) {}

func (s StateView) Draw(v *View) {
	v.Clear()
	v.Printf("%v", *s.state)
}

type OpsView struct {
	vm *vm.VM
}

func (o OpsView) Init(_ *View) {}

func (o OpsView) Draw(v *View) {
	v.Clear()
	v.Printf("%d ops executed", o.vm.OpsExecuted())
}

type MemoryView struct {
	m  *memory.Memory
	pc *int
	b  *base
}

func (m MemoryView) Init(v *View) {
	v.Title = "Memory"
}

const memoryLineLength = 16

func (m MemoryView) Draw(v *View) {
	v.Clear()
	if m.m != nil {
		x, y := v.Size()
		v.Println()
		for line := 0; line < y-2; line++ {
			startAddr := *m.pc + memoryLineLength*line
			endAddr := startAddr + memoryLineLength
			s := drawMemLine(*m.b, startAddr, m.m.Mem[startAddr:endAddr])
			spaces := (x - len(s)) / 2
			for i := 0; i < spaces; i++ {
				v.Print(" ")
			}
			v.Println(s)
		}
	}
}

func drawMemLine(b base, startAddr int, mem []uint16) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%s:  ", b.strSym(startAddr)))
	for _, w := range mem {
		sb.WriteString(fmt.Sprintf("%s ", b.str(int(w))))
	}
	sb.WriteString("  ")
	for _, w := range mem {
		ch := '.'
		if isChar(w) {
			ch = rune(w)
		}
		sb.WriteString(string(ch))
	}
	return sb.String()
}

func isChar(w uint16) bool {
	if w >= 32 && w <= 126 {
		return true
	}
	return false
}
