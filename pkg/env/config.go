package env

import (
	"bufio"
	"io"
	"os"
)

// TODO: add halt
var Config = struct {
	Output   io.Writer
	ReadChar func() (uint16, error)
}{
	os.Stdin,
	readChar,
}

func readChar() (uint16, error) {
	reader := bufio.NewReader(os.Stdin)
	ch, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return uint16(ch), nil
}
