package debugger

import (
	"dancavallaro.com/synacor-go/pkg/vm"
	"errors"
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
	if err := gui.SetKeybinding("", 'x', gocui.ModNone, toggleBase); err != nil {
		return err
	}
	return nil
}

func (d Debugger) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if err := drawView(g, OutputView{}, "output", 0, 0, maxX-1, maxY-7); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("output"); err != nil {
		return err
	}

	if err := drawView(g, StackView{&d.VM.M}, "stack", -1, maxY-6, maxX, maxY-4); err != nil {
		return err
	}

	if err := drawView(g, RegisterView{&d.VM.M}, "registers", -1, maxY-4, maxX, maxY-2); err != nil {
		return err
	}

	if err := drawView(g, HelpView{}, "help", -1, maxY-2, maxX, maxY); err != nil {
		return err
	}

	return nil
}

type Frame interface {
	Init(v *gocui.View)
	Draw(v *gocui.View)
}

func drawView(g *gocui.Gui, f Frame, name string, x0, y0, x1, y1 int) error {
	var v *gocui.View
	var err error
	if v, err = g.SetView(name, x0, y0, x1, y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		f.Init(v)
	}
	f.Draw(v)

	return nil
}

func (d Debugger) execute(_ *gocui.Gui, _ *gocui.View) error {
	go d.VM.Execute()
	return nil
}

type base int

const (
	hex base = iota
	dec
)

var displayBase = hex

func toggleBase(_ *gocui.Gui, _ *gocui.View) error {
	if displayBase == hex {
		displayBase = dec
	} else {
		displayBase = hex
	}
	return nil
}
