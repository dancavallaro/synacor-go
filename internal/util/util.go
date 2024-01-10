package util

import (
	"fmt"
	"io"
	"os"
)

func ReadChar() (uint16, error) {
	fmt.Println("about to read")
	b, err := io.ReadAll(os.Stdin)
	fmt.Println("just read")
	if err != nil {
		return 0, err
	}
	return uint16(b[0]), nil
}
