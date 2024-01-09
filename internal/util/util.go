package util

import (
	"io"
	"os"
)

func ReadChar() (uint16, error) {
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return 0, err
	}
	return uint16(b[0]), nil
}
