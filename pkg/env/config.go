package env

import (
	"bufio"
	"io"
	"log"
	"os"
)

var Config = struct {
	Output   io.Writer
	ReadChar func() (uint16, error)
	Halt     func()
}{
	os.Stdin,
	readChar,
	halt,
}

func halt() {
	log.Println()
	log.Println("Execution halted.")
	os.Exit(0)
}

func readChar() (uint16, error) {
	reader := bufio.NewReader(os.Stdin)
	ch, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return uint16(ch), nil
}
