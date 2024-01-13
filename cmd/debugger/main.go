package main

import (
	"dancavallaro.com/synacor-go/pkg/debugger"
	"dancavallaro.com/synacor-go/pkg/vm"
	"errors"
	"flag"
	"github.com/awesome-gocui/gocui"
	"log"
	"net/http"
	_ "net/http"
	_ "net/http/pprof"
	"os"
)

const required = "<required>"

var (
	binPath = flag.String("bin", required, "Path to executable (.bin)")
)

func main() {
	flag.Parse()
	if *binPath == required {
		log.Fatalln("-bin is required")
	}
	binary, err := os.ReadFile(*binPath)
	if err != nil {
		log.Fatalln(err)
	}

	gui, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatalln(err)
	}
	defer gui.Close()

	opts := vm.ExecutionOptions{
		Trace:      false,
		Delay:      -1,
		Breakpoint: -1,
	}
	vm := vm.NewVM(binary, &opts)
	debug := debugger.NewDebugger(vm, gui)
	if err != nil {
		log.Fatalln(err)
	}
	gui.SetManagerFunc(debug.Layout)
	gui.Cursor = true
	err = debug.InitKeybindings(gui)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO delete
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, exit); err != nil {
		log.Fatalln(err)
	}

	if err := gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Fatalln(err)
	}
}

func exit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}
