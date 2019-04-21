package main

import (
	"testing"
	"unsafe"

	"github.com/fluent/fluent-bit-go/output"
	consul_api "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

type testOutput struct {
	inc int
}

func (o *testOutput) Register(ctx unsafe.Pointer, name string, desc string) int {
	return output.FLBPluginRegister(ctx, name, desc)
}

func (o *testOutput) GetConfigKey(ctx unsafe.Pointer, key string) string {
	if key == "ConsulServer" {
		return "localhost"
	}
	if key == "ConsulPort" {
		return "8500"
	}
	if key == "ServiceName" {
		return "allan-service"
	}
	if key == "CheckPort" {
		return "2020"
	}
	if key == "ServiceId" {
		return ""
	}
	return ""
}

func (o *testOutput) GetConsulConfig() *consul_api.Config {
	return &consul_api.Config{
		HttpClient: mockClient(),
	}
}

func TestFLBPluginInit(t *testing.T) {
	assert := assert.New(t)
	wrapper = OutputWrapper(&testOutput{})
	assert.Equal(output.FLB_OK, FLBPluginInit(nil))
	assert.NotEmpty(hostname)
	assert.NotEmpty(localIp)
	assert.NotEmpty(plugin)
	assert.Equal(hostname, serviceId)
	assert.Equal(2020, checkPort)
	assert.Equal("allan-service", serviceName)
}

func TestFLBPluginFlush(t *testing.T) {
	assert := assert.New(t)
	ok := FLBPluginFlush(nil, 0, nil)
	assert.Equal(output.FLB_OK, ok)
}

func TestFLBPluginExit(t *testing.T) {
	assert := assert.New(t)
	ok := FLBPluginExit()
	assert.Equal(output.FLB_OK, ok)
}