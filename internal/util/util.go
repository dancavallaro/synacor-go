package util

import (
	"fmt"
	"io"
	"os"
)

func ReadChar() uint16 {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(fmt.Sprintf("error reading from stdin: %s", err))
	}
	return uint16(b[0])
}

func WaitForEnter() {
	ReadChar()
}
