package debugger

import (
	"dancavallaro.com/synacor-go/pkg/env"
	"dancavallaro.com/synacor-go/pkg/op"
	"dancavallaro.com/synacor-go/pkg/vm"
	"errors"
	"fmt"
	"github.com/awesome-gocui/gocui"
	"log"
	"time"
)

type State int

const (
	Paused State = iota
	Running
	StepOnce
	Halted
)

func (s State) String() string {
	if s == Paused {
		return "PAUSED"
	} else if s == Running {
		return "RUNNING"
	} else if s == StepOnce {
		return "STEP"
	} else if s == Halted {
		return "HALTED"
	}
	panic("unknown value of State!")
}

type Debugger struct {
	VM             *vm.VM
	viewsToRefresh map[*gocui.View]Frame
	state          State
	inCharCh       chan uint16
	inputCh        chan string
	displayBase    base
}

func NewDebugger(VM *vm.VM, g *gocui.Gui) *Debugger {
	debug := &Debugger{VM, make(map[*gocui.View]Frame), Paused, make(chan uint16), make(chan string), hex}
	go debug.refreshUI(g)
	go debug.executeWhenRunning()
	env.Config.ReadChar = requestChar(g, debug)
	env.Config.Halt = debug.halt
	return debug
}

func (d *Debugger) refreshUI(g *gocui.Gui) {
	for {
		select {
		case <-time.After(100 * time.Millisecond):
			g.Update(func(g *gocui.Gui) error {
				for v, f := range d.viewsToRefresh {
					f.Draw(&View{v})
				}
				return nil
			})
		}
	}
}

func (d *Debugger) executeWhenRunning() {
	for {
		state := d.state
		if state == Running || state == StepOnce {
			if err := d.VM.Step(); err != nil {
				// An error here from the VM includes things like invalid opcodes and
				// may not be recoverable, so just panic.
				panic(err)
			}
		}
		if state == StepOnce && d.state != Halted {
			d.state = Paused
		}
	}
}

func (d *Debugger) InitKeybindings(gui *gocui.Gui) error {
	if err := gui.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, d.pause); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, d.execute); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlS, gocui.ModNone, d.step); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlX, gocui.ModNone, d.toggleBase); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlBackslash, gocui.ModNone, d.restart); err != nil {
		return err
	}
	if err := gui.SetKeybinding("", gocui.KeyCtrlB, gocui.ModNone, d.runUntil); err != nil {
		return err
	}
	return nil
}

func requestChar(g *gocui.Gui, debugger *Debugger) func() (uint16, error) {
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
			if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, readChar(debugger.inCharCh)); err != nil {
				return 0, err
			}
		}
		v.Visible = true
		if _, err := g.SetCurrentView("msg"); err != nil {
			return 0, err
		}
		g.Cursor = true

		ch := <-debugger.inCharCh
		if ch != op.CancelInput {
			log.Printf("IN read '%s' (%d) from stdin\n", str(rune(ch)), rune(ch))
		}
		return ch, nil
	}
}

func requestInput(g *gocui.Gui, debugger *Debugger) (string, error) {
	var v *gocui.View
	var err error
	maxX, maxY := g.Size()
	if v, err = g.SetView("input", maxX/2-30, maxY/2, maxX/2+30, maxY/2+2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return "", err
		}
		v.Editable = true
		v.Title = "Enter breakpoint:"
		if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, readInput(debugger.inputCh)); err != nil {
			return "", err
		}
	}
	v.Visible = true
	if _, err := g.SetCurrentView("input"); err != nil {
		return "", err
	}
	g.Cursor = true

	//return <-debugger.inputCh, nil
	return "", nil
}

func str(ch rune) string {
	if ch == '\n' {
		return "\\n"
	} else {
		return string(ch)
	}
}

func readChar(input chan<- uint16) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		l := v.Buffer()
		if len(l) == 0 {
			input <- uint16('\n')
		} else {
			input <- uint16(l[0])
		}
		return closeModal(g, v)
	}
}

func readInput(input chan<- string) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		//input <- v.Buffer()
		return closeModal(g, v)
	}
}

func closeModal(g *gocui.Gui, v *gocui.View) error {
	v.Visible = false
	v.Clear()
	g.Cursor = false
	if _, err := g.SetCurrentView("output"); err != nil {
		return err
	}
	return nil
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

	if err := d.drawView(g, MemoryView{&d.VM.M, &d.VM.OriginalPC, &d.displayBase}, "memory", 0, midY+1, mult(maxX, 0.75), maxY-7, true); err != nil {
		return err
	}

	if err := d.drawView(g, LogView{}, "log", int(0.75*float32(maxX))+1, 0, maxX-1, midY, true); err != nil {
		return err
	}

	if err := d.drawView(g, DisassemblyView{d, &d.displayBase}, "disassembly", int(0.75*float32(maxX))+1, midY+1, maxX-1, maxY-7, true); err != nil {
		return err
	}

	if err := d.drawView(g, StackView{&d.VM.M, &d.displayBase}, "stack", -1, maxY-6, maxX, maxY-4, true); err != nil {
		return err
	}

	if err := d.drawView(g, RegisterView{&d.VM.M, &d.VM.OriginalPC, &d.displayBase}, "registers", -1, maxY-4, maxX, maxY-2, true); err != nil {
		return err
	}

	if err := d.drawView(g, HelpView{}, "help", -1, maxY-2, maxX-35, maxY, false); err != nil {
		return err
	}

	if err := d.drawView(g, OpsView{d.VM}, "ops", maxX-35, maxY-2, maxX-10, maxY, true); err != nil {
		return err
	}

	if err := d.drawView(g, StateView{&d.state}, "state", maxX-10, maxY-2, maxX, maxY, true); err != nil {
		return err
	}

	return nil
}

func (d *Debugger) drawView(g *gocui.Gui, f Frame, name string, x0, y0, x1, y1 int, refresh bool) error {
	var v *gocui.View
	var err error
	if v, err = g.SetView(name, x0, y0, x1, y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		f.Init(&View{v})
		if refresh {
			d.viewsToRefresh[v] = f
		}
	}
	f.Draw(&View{v})

	return nil
}

func (d *Debugger) pause(_ *gocui.Gui, _ *gocui.View) error {
	if d.state == Halted {
		return nil
	}
	d.state = Paused
	return nil
}

func (d *Debugger) execute(_ *gocui.Gui, _ *gocui.View) error {
	if d.state == Halted {
		return nil
	}
	d.state = Running
	return nil
}

func (d *Debugger) step(_ *gocui.Gui, _ *gocui.View) error {
	if d.state == Halted {
		return nil
	}
	d.state = StepOnce
	return nil
}

func (d *Debugger) restart(g *gocui.Gui, _ *gocui.View) error {
	log.Println("Resetting state and restarting VM...")
	d.state = Paused

	o, err := g.View("output")
	if err != nil {
		return err
	}
	o.Clear()

	if v := g.CurrentView(); v != nil && v.Name() == "msg" {
		d.inCharCh <- op.CancelInput
		if err := closeModal(g, v); err != nil {
			return err
		}
	}

	d.VM.Restart()

	return nil
}

func (d *Debugger) runUntil(g *gocui.Gui, _ *gocui.View) error {
	// TODO: this is causing a deadlock for some reason
	input, err := requestInput(g, d)
	if err != nil {
		return err
	}
	log.Printf("Got input: '%s'\n", input)
	return nil
}

func (d *Debugger) halt() {
	log.Println("Execution halted.")
	d.state = Halted
}

type base int

const (
	hex base = iota
	dec
)

func (b base) str(n int) string {
	if b == hex {
		return fmt.Sprintf("%04x", n)
	} else if b == dec {
		return fmt.Sprintf("%06d", n)
	}
	panic("unknown base!")
}

func (b base) strSym(n int) string {
	return fmt.Sprintf("%s%s", b.prefix(), b.str(n))
}

func (b base) prefix() string {
	if b == hex {
		return "0x"
	}
	return ""
}

func (d *Debugger) toggleBase(_ *gocui.Gui, _ *gocui.View) error {
	if d.displayBase == hex {
		d.displayBase = dec
	} else {
		d.displayBase = hex
	}
	return nil
}
