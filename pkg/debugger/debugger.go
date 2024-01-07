package debugger

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"dancavallaro.com/synacor-go/pkg/op"
	"dancavallaro.com/synacor-go/pkg/vm"
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
)

type Debugger struct {
	VM *vm.VM
}

func NewDebugger(VM *vm.VM) *Debugger {
	return &Debugger{VM}
}

func (d Debugger) InitKeybindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("", 'r', gocui.ModNone, d.execute); err != nil {
		return err
	}
	return nil
}

func (d Debugger) Layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	if v, err := gui.SetView("output", 0, 0, maxX-1, maxY-5, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		if _, err := gui.SetCurrentView("output"); err != nil {
			return err
		}
		v.Title = "Output"
		op.Output = v
	}

	var regView *gocui.View
	var err error
	if regView, err = gui.SetView("registers", -1, maxY-4, maxX, maxY-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		regView.Title = "Registers"
	}
	regView.Clear()
	d.drawRegisters(regView)

	if v, err := gui.SetView("help", -1, maxY-2, maxX, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		v.Title = "Help"
		fmt.Fprint(v, "r: resume execution")
	}

	return nil
}

func (d Debugger) drawRegisters(v *gocui.View) {
	pc := d.VM.M.PC
	fmt.Fprintf(v, "PC: %#04x\t", pc)

	gp := d.VM.M.GP
	for i := 0; i < memory.NumRegisters; i++ {
		fmt.Fprintf(v, "R%d: %#04x\t", i, gp[i])
	}
}

func (d Debugger) execute(_ *gocui.Gui, _ *gocui.View) error {
	go d.VM.Execute()
	return nil
}
