package debugger

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"dancavallaro.com/synacor-go/pkg/op"
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type HelpView struct{}

func (h HelpView) Init(v *gocui.View) {
	v.Title = "Help"
	fmt.Fprint(v, "(r) resume execution\t(x) toggle hex/dec")
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
	op.Output = v
}

func (h OutputView) Draw(_ *gocui.View) {}
