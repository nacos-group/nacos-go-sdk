package rpc

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/keepalive"
	"testing"
	"time"
)

func TestGetMaxCallRecvMsgSize(t *testing.T) {
	size := getMaxCallRecvMsgSize()
	assert.Equal(t, size, 10*1024*1024)
}

func TestGetKeepAliveTimeMillis(t *testing.T) {
	result := getKeepAliveTimeMillis()
	expect := keepalive.ClientParameters{
		Time:                6 * 60 * 1000 * time.Millisecond,
		Timeout:             3 * time.Second,
		PermitWithoutStream: true,
	}
	assert.Equal(t, expect, result)
}
