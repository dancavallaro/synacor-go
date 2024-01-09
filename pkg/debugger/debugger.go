package debugger

import (
	"dancavallaro.com/synacor-go/pkg/env"
	"dancavallaro.com/synacor-go/pkg/vm"
	"errors"
	"github.com/awesome-gocui/gocui"
	"time"
)

type State int

const (
	Paused State = iota
	Running
	Pausing
	StepOnce
)

type Debugger struct {
	VM             *vm.VM
	viewsToRefresh map[*gocui.View]Frame
	state          State
	inputCh        chan uint16
}

func NewDebugger(VM *vm.VM, g *gocui.Gui) *Debugger {
	debug := &Debugger{VM, make(map[*gocui.View]Frame), Paused, make(chan uint16)}
	go debug.refreshUI(g)
	go debug.executeWhenRunning()
	env.Config.ReadChar = requestInput(g, debug)
	return debug
}

func (d *Debugger) refreshUI(g *gocui.Gui) {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			if d.state != Paused {
				g.Update(func(g *gocui.Gui) error {
					for v, f := range d.viewsToRefresh {
						f.Draw(v)
					}
					return nil
				})
			}
			if d.state == Pausing {
				d.state = Paused
			}
		}
	}
}

func (d *Debugger) executeWhenRunning() {
	for {
		if d.state == Running || d.state == StepOnce {
			if err := d.VM.Step(); err != nil {
				// TODO: Don't panic, this should return an error to the debugger
				panic(err)
			}
		}
		if d.state == StepOnce {
			d.state = Pausing
		}
	}
}

func (d *Debugger) InitKeybindings(gui *gocui.Gui) error {
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

func requestInput(g *gocui.Gui, debugger *Debugger) func() (uint16, error) {
	return func() (uint16, error) {
		output, err := g.View("output")
		if err != nil {
			return 0, err
		}

		var v *gocui.View
		maxX, maxY := output.Size()
		if v, err = g.SetView("msg", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return 0, err
			}
			v.Editable = true
			v.Title = "Enter a character:"
			if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, readInput(debugger.inputCh)); err != nil {
				return 0, err
			}
		}
		v.Visible = true
		if _, err := g.SetCurrentView("msg"); err != nil {
			return 0, err
		}

		return <-debugger.inputCh, nil
	}
}

func readInput(input chan<- uint16) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		var l string
		var err error
		_, cy := v.Cursor()
		if l, err = v.Line(cy); err != nil {
			return err
		}
		input <- uint16(l[0])
		v.Visible = false
		v.Clear()
		if _, err := g.SetCurrentView("output"); err != nil {
			return err
		}
		return nil
	}
}

func mult(base int, fraction float32) int {
	return int(fraction * float32(base))
}

func (d *Debugger) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	midY := mult(maxY-7, 0.5)

	if err := d.drawView(g, OutputView{}, "output", 0, 0, mult(maxX, 0.75), midY, false); err != nil {
		return err
	}

	if err := d.drawView(g, MemoryView{&d.VM.M}, "memory", 0, midY+1, mult(maxX, 0.75), maxY-7, true); err != nil {
		return err
	}

	if err := d.drawView(g, LogView{}, "log", int(0.75*float32(maxX))+1, 0, maxX-1, midY, true); err != nil {
		return err
	}

	if err := d.drawView(g, DisassemblyView{}, "disassembly", int(0.75*float32(maxX))+1, midY+1, maxX-1, maxY-7, true); err != nil {
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

func (d *Debugger) drawView(g *gocui.Gui, f Frame, name string, x0, y0, x1, y1 int, refresh bool) error {
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

func (d *Debugger) execute(_ *gocui.Gui, _ *gocui.View) error {
	d.state = Running
	return nil
}

func (d *Debugger) step(_ *gocui.Gui, _ *gocui.View) error {
	d.state = StepOnce
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
