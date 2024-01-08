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
	fmt.Fprint(v, "(r) resume execution")
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

	pc := h.m.PC
	fmt.Fprintf(v, "PC: %#04x\t", pc)

	gp := h.m.GP
	for i := 0; i < memory.NumRegisters; i++ {
		fmt.Fprintf(v, "R%d: %#04x\t", i, gp[i])
	}
}

type OutputView struct{}

func (h OutputView) Init(v *gocui.View) {
	v.Title = "Output"
	op.Output = v
}

func (h OutputView) Draw(_ *gocui.View) {}
