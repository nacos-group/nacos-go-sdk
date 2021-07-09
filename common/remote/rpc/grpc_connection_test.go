package rpc

import (
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/common/remote/rpc/rpc_request"
	"github.com/nacos-group/nacos-sdk-go/common/remote/rpc/rpc_response"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertRequest(t *testing.T) {
	r := rpc_request.NewConnectionSetupRequest()
	p := convertRequest(r)
	assert.NotNil(t, p)
}

func TestConvertResponse(t *testing.T) {
	response := rpc_response.Response{
		RequestId:  "xxx",
		ResultCode: constant.RESPONSE_CODE_SUCCESS,
	}
	r := rpc_response.HealthCheckResponse{
		Response: &response,
	}
	p := convertResponse(&r)
	assert.NotNil(t, p)
}
