package debugger

import (
	"dancavallaro.com/synacor-go/pkg/vm"
	"errors"
	"github.com/awesome-gocui/gocui"
	"time"
)

type Debugger struct {
	VM             *vm.VM
	viewsToRefresh map[*gocui.View]Frame
}

func NewDebugger(VM *vm.VM, g *gocui.Gui) *Debugger {
	debug := &Debugger{VM, make(map[*gocui.View]Frame)}
	go debug.refreshUI(g)
	return debug
}

func (d Debugger) refreshUI(g *gocui.Gui) {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			g.Update(func(g *gocui.Gui) error {
				for v, f := range d.viewsToRefresh {
					f.Draw(v)
				}
				return nil
			})
		}
	}
}

func (d Debugger) InitKeybindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("", 'r', gocui.ModNone, d.execute); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", 's', gocui.ModNone, d.step); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", 'x', gocui.ModNone, toggleBase); err != nil {
		return err
	}
	return nil
}

func mult(base int, fraction float32) int {
	return int(fraction * float32(base))
}

func (d Debugger) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if err := d.drawView(g, OutputView{}, "output", 0, 0, mult(maxX, 0.75), mult(maxY-7, 0.5), false); err != nil {
		return err
	}
	if _, err := g.SetCurrentView("output"); err != nil {
		return err
	}

	if err := d.drawView(g, MemoryView{&d.VM.M}, "memory", 0, mult(maxY-7, 0.5)+1, mult(maxX, 0.75), maxY-7, true); err != nil {
		return err
	}

	if err := d.drawView(g, DisassemblyView{}, "disassembly", int(0.75*float32(maxX))+1, 0, maxX-1, maxY-7, true); err != nil {
		return err
	}

	if err := d.drawView(g, StackView{&d.VM.M}, "stack", -1, maxY-6, maxX, maxY-4, true); err != nil {
		return err
	}

	if err := d.drawView(g, RegisterView{&d.VM.M}, "registers", -1, maxY-4, maxX, maxY-2, true); err != nil {
		return err
	}

	if err := d.drawView(g, HelpView{}, "help", -1, maxY-2, maxX, maxY, false); err != nil {
		return err
	}

	return nil
}

type Frame interface {
	Init(v *gocui.View)
	Draw(v *gocui.View)
}

func (d Debugger) drawView(g *gocui.Gui, f Frame, name string, x0, y0, x1, y1 int, refresh bool) error {
	var v *gocui.View
	var err error
	if v, err = g.SetView(name, x0, y0, x1, y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		f.Init(v)
		if refresh {
			d.viewsToRefresh[v] = f
		}
	}
	f.Draw(v)

	return nil
}

func (d Debugger) execute(_ *gocui.Gui, _ *gocui.View) error {
	go d.VM.Execute()
	return nil
}

func (d Debugger) step(_ *gocui.Gui, _ *gocui.View) error {
	return d.VM.Step()
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
