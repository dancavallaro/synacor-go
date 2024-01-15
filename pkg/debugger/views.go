package debugger

import (
	"dancavallaro.com/synacor-go/pkg/env"
	"dancavallaro.com/synacor-go/pkg/memory"
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
	v.Print("(^p) pause execution\t(^r) resume execution\t(^s) step execution\t(^x) toggle hex/dec\t(^\\) reset state and restart")
}

func (h HelpView) Draw(_ *View) {}

type RegisterView struct {
	m  *memory.Memory
	pc *int
	b  *base
}

func (h RegisterView) Init(v *View) {
	v.Title = "Registers"
}

func (h RegisterView) Draw(v *View) {
	v.Clear()
	gp := h.m.GP

	v.Printf("PC: %s\t", h.b.strSym(*h.pc))
	for i := 0; i < memory.NumRegisters; i++ {
		v.Printf("R%d: %s\t", i, h.b.strSym(int(gp[i])))
	}
}

type StackView struct {
	m *memory.Memory
	b *base
}

func (h StackView) Init(v *View) {
	v.Title = "Stack"
}

func (h StackView) Draw(v *View) {
	v.Clear()

	for i := 0; i < len(h.m.Stack); i++ {
		v.Print(h.b.strSym(int(h.m.Stack[i])))
		if i < len(h.m.Stack)-1 {
			v.Printf("\t")
		} else {
			v.Printf(" ")
		}
	}

	v.Print("◄SP")
}

type OutputView struct{}

func (h OutputView) Init(v *View) {
	v.Title = "Output"
	v.Autoscroll = true
	env.Config.Output = v
}

func (h OutputView) Draw(_ *View) {}

type LogView struct{}

func (h LogView) Init(v *View) {
	v.Title = "System Log"
	v.Autoscroll = true
	log.Default().SetOutput(v)
}

func (h LogView) Draw(_ *View) {}

type DisassemblyView struct {
	d *Debugger
	b *base
}

func (h DisassemblyView) Init(v *View) {
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

func (h DisassemblyView) Draw(v *View) {
	v.Clear()
	v.Println()
	pc := h.d.VM.OriginalPC
	_, y := v.Size()
	for line := 0; line < y-2; line++ {
		o, args, err := h.d.VM.DecodeOp(pc)
		if err != nil {
			panic(err)
		}
		v.Printf("    %s: %s%s\n", h.b.strSym(pc), o.Mnemonic, argStr(args, *h.b))
		pc += 1 + len(args)
	}
}

type StateView struct {
	state *State
}

func (h StateView) Init(_ *View) {}

func (h StateView) Draw(v *View) {
	v.Clear()
	v.Printf("%v", *h.state)
}

type MemoryView struct {
	m  *memory.Memory
	pc *int
	b  *base
}

func (h MemoryView) Init(v *View) {
	v.Title = "Memory"
}

const memoryLineLength = 16

func (h MemoryView) Draw(v *View) {
	v.Clear()
	if h.m != nil {
		x, y := v.Size()
		v.Println()
		for line := 0; line < y-2; line++ {
			startAddr := *h.pc + memoryLineLength*line
			endAddr := startAddr + memoryLineLength
			s := drawMemLine(*h.b, startAddr, h.m.Mem[startAddr:endAddr])
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
