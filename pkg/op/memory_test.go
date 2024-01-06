package op

import (
	"dancavallaro.com/synacor-go/pkg/memory"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	r := &memory.Registers{}
	Set(r, args(R0, 42))
	assert.Equal(t, uint16(42), r.GP[0])
}

func TestSet_RegVal(t *testing.T) {
	r := &memory.Registers{}
	r.GP[1] = 42
	Set(r, args(R0, R1))
	assert.Equal(t, uint16(42), r.GP[0])
}
