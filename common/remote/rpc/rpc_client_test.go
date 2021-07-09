package rpc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheck(t *testing.T) {

}

func TestIsConnected(t *testing.T) {
	con := ConnectionEvent{eventType: CONNECTED}
	r := con.isConnected()
	assert.True(t, r)
}

func TestIsDisconnected(t *testing.T) {
	con := ConnectionEvent{DISCONNECTED}
	r := con.isDisConnected()
	assert.True(t, r)
}

func TestToString(t *testing.T) {
	con1 := ConnectionEvent{eventType: CONNECTED}
	r1 := con1.toString()
	assert.Equal(t, "connected", r1)

	con2 := ConnectionEvent{DISCONNECTED}
	r2 := con2.toString()
	assert.Equal(t, "disconnected", r2)
}
