package udp

import (
	"testing"

	"gotest.tools/assert"
)

func TestRandUDP(t *testing.T) {
	ln, err := randUDPListener()
	assert.Assert(t, err == nil && ln != nil)
	ln.Close()

}
