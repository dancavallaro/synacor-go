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
	binPath     = flag.String("bin", required, "Path to executable (.bin)")
	enableDebug = flag.Bool("debug", false, "Enable Go debug server on port 6060")
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

	vm := vm.NewVM(binary, &vm.ExecutionOptions{})
	debug := debugger.NewDebugger(vm, gui)
	if err != nil {
		log.Fatalln(err)
	}
	gui.SetManagerFunc(debug.Layout)
	err = debug.InitKeybindings(gui)
	if err != nil {
		log.Fatalln(err)
	}

	if *enableDebug {
		go func() {
			http.ListenAndServe("localhost:6060", nil)
		}()
	}

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
