package main

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

// It is a external interface for consul.
type Keeper interface {
	Register(serviceId, serviceName, addr string, port int, tags []string,
		check *consul.AgentServiceCheck) error
	DeRegister(serviceId string) error
}

type Container struct {
	addr   string
	port   string
	config *consul.Config
	agent  *consul.Agent
}

// This struct used to create `keeper` object from use to a way of dependency injection.
type Option interface {
	apply(c *Container)
}

type optionFunc func(c *Container)

func (of optionFunc) apply(c *Container) {
	of(c)
}

// Create Addr Option
func WithAddrOption(addr string) Option {
	return optionFunc(func(c *Container) {
		c.addr = addr
	})
}

// Create Port Option
func WithPortOption(port string) Option {
	return optionFunc(func(c *Container) {
		c.port = port
	})
}

// Create Consul Config Option
func WithConfigOption(config *consul.Config) Option {
	return optionFunc(func(c *Container) {
		c.config = config
	})
}

// Create interface for consul, it
func NewKeeper(opts ...Option) (Keeper, error) {
	// Dependency injection.
	o := []Option{
		WithAddrOption("localhost"),
		WithPortOption("8500"),
		WithConfigOption(consul.DefaultConfig()),
	}
	o = append(o, opts...)
	ct := &Container{}
	for _, opt := range o {
		opt.apply(ct)
	}

	// Connect to consul.
	ct.config.Address = ct.address()
	client, err := consul.NewClient(ct.config)
	if err != nil {
		return nil, errors.Wrap(err, "[err] consul client")
	}
	ct.agent = client.Agent()
	return Keeper(ct), nil
}

func (c *Container) address() string {
	return fmt.Sprintf("%s:%s", c.addr, c.port)
}

// A method naming `Register` could connect to service in consul.
func (c *Container) Register(serviceId, serviceName, addr string, port int, tags []string,
	check *consul.AgentServiceCheck) error {
	reg := &consul.AgentServiceRegistration{
		ID:      serviceId,
		Name:    serviceName,
		Address: addr,
		Port:    port,
		Tags:    tags,
	}
	if check == nil {
		check = &consul.AgentServiceCheck{}
		check.HTTP = fmt.Sprintf("http://%s:%v", addr, port)
		check.Interval = "10s"
		check.Timeout = "5m"
		check.DeregisterCriticalServiceAfter = "50s"
	}
	reg.Check = check
	return c.agent.ServiceRegister(reg)
}

// A method naming `DeRegister` could disconnect to service in consul.
func (c *Container) DeRegister(serviceId string) error {
	return c.agent.ServiceDeregister(serviceId)
}