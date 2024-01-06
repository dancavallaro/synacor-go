package vm

import "fmt"

func Execute(bin []byte) {
	for i := 0; i < 100; i++ {
		fmt.Printf("%d ", bin[i])
	}
}
