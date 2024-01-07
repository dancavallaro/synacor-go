package util

import "os"

func ReadChar() uint16 {
	b := make([]byte, 1)
	n, err := os.Stdin.Read(b)
	if n != 1 || err != nil {
		panic("error reading from stdin!")
	}
	return uint16(b[0])
}

func WaitForEnter() {
	ReadChar()
}
