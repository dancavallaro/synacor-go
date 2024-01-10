package util

import (
	"bufio"
	"os"
)

func ReadChar() (uint16, error) {
	reader := bufio.NewReader(os.Stdin)
	ch, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return uint16(ch), nil
}
