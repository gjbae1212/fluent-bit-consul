package main

import (
	"testing"

	consul_api "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func TestWithAddrOption(t *testing.T) {
	assert := assert.New(t)

	c := &Container{}
	opt := WithAddrOption("allan")
	opt.apply(c)
	assert.Equal(c.addr, "allan")
}

func TestWithPortOption(t *testing.T) {
	assert := assert.New(t)

	c := &Container{}
	opt := WithPortOption("1010")
	opt.apply(c)
	assert.Equal(c.port, "1010")
}

func TestWithConfigOption(t *testing.T) {
	assert := assert.New(t)

	c := &Container{}
	config := &consul_api.Config{
		Address: "testtest",
	}
	opt := WithConfigOption(config)
	opt.apply(c)
	assert.Equal(c.config.Address, "testtest")
}

func TestNewKeeper(t *testing.T) {
	assert := assert.New(t)

	keeper, err := NewKeeper()
	assert.NoError(err)

	assert.Equal(keeper.(*Container).addr, "localhost")
	assert.Equal(keeper.(*Container).port, "8500")
	assert.NotEmpty(keeper.(*Container).config)
	assert.NotEmpty(keeper.(*Container).agent)

	opts := []Option{
		WithAddrOption("127.0.0.1"),
		WithPortOption("8080"),
	}

	keeper, err = NewKeeper(opts...)
	assert.NoError(err)
	assert.Equal(keeper.(*Container).addr, "127.0.0.1")
	assert.Equal(keeper.(*Container).port, "8080")
	assert.NotEmpty(keeper.(*Container).config)
	assert.NotEmpty(keeper.(*Container).agent)
}

//func TestContainer_DeRegister(t *testing.T) {
//	assert := assert.New(t)
//}
//
//func TestContainer_Register(t *testing.T) {
//	assert := assert.New(t)
//}
//
//func TestContainer_Private(t *testing.T) {
//
//}
