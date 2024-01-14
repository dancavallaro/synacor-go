package debugger

import (
	"dancavallaro.com/synacor-go/pkg/env"
	"dancavallaro.com/synacor-go/pkg/memory"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
)

type HelpView struct{}

func (h HelpView) Init(v *gocui.View) {
	v.Title = "Help"
	fmt.Fprint(v, "(^p) pause execution\t(^r) resume execution\t(^s) step execution\t(^x) toggle hex/dec")
}

func (h HelpView) Draw(_ *gocui.View) {}

type RegisterView struct {
	m *memory.Memory
}

func (h RegisterView) Init(v *gocui.View) {
	v.Title = "Registers"
}

func (h RegisterView) Draw(v *gocui.View) {
	v.Clear()
	pc, gp := h.m.PC, h.m.GP

	if displayBase == hex {
		fmt.Fprintf(v, "PC: %#04x\t", pc)
	} else {
		fmt.Fprintf(v, "PC: %06d\t", pc)
	}
	for i := 0; i < memory.NumRegisters; i++ {
		if displayBase == hex {
			fmt.Fprintf(v, "R%d: %#04x\t", i, gp[i])
		} else {
			fmt.Fprintf(v, "R%d: %06d\t", i, gp[i])
		}
	}
}

type StackView struct {
	m *memory.Memory
}

func (h StackView) Init(v *gocui.View) {
	v.Title = "Stack"
}

func (h StackView) Draw(v *gocui.View) {
	v.Clear()

	for i := 0; i < len(h.m.Stack); i++ {
		if displayBase == hex {
			fmt.Fprintf(v, "%#04x", h.m.Stack[i])
		} else {
			fmt.Fprintf(v, "%06d", h.m.Stack[i])
		}
		if i < len(h.m.Stack)-1 {
			fmt.Fprintf(v, "\t")
		} else {
			fmt.Fprintf(v, " ")
		}
	}

	fmt.Fprint(v, "◄SP")
}

type OutputView struct{}

func (h OutputView) Init(v *gocui.View) {
	v.Title = "Output"
	v.Autoscroll = true
	env.Config.Output = v
}

func (h OutputView) Draw(_ *gocui.View) {}

type LogView struct{}

func (h LogView) Init(v *gocui.View) {
	v.Title = "System Log"
	v.Autoscroll = true
	log.Default().SetOutput(v)
}

func (h LogView) Draw(_ *gocui.View) {}

type DisassemblyView struct{}

func (h DisassemblyView) Init(v *gocui.View) {
	v.Title = "Disassembly"
}

func (h DisassemblyView) Draw(_ *gocui.View) {
	// TODO
}

type MemoryView struct {
	m *memory.Memory
}

func (h MemoryView) Init(v *gocui.View) {
	v.Title = "Memory"
}

func (h MemoryView) Draw(v *gocui.View) {
	v.Clear()
	if h.m != nil {
		// TODO: Make this dynamic based on size of view?
		// TODO: Or at least center it in the panel
		lineLength := 16

		_, y := v.Size()
		fmt.Fprintln(v)
		for line := 0; line < y-2; line++ {
			startAddr := h.m.PC + lineLength*line
			drawMemLine(v, startAddr, h.m.Mem[startAddr:startAddr+16])
		}
	}
}

func drawMemLine(v *gocui.View, startAddr int, mem []uint16) {
	fmt.Fprintf(v, "  %#04x:  ", startAddr)
	for _, w := range mem {
		fmt.Fprintf(v, "%04x ", w) // TODO: Support toggling units, refactor that code
	}
	fmt.Fprintf(v, "  ")
	for _, w := range mem {
		ch := '.'
		if w >= 32 && w <= 126 {
			ch = rune(w)
		}
		fmt.Fprint(v, string(ch))
	}
	fmt.Fprint(v, "\n")
}
