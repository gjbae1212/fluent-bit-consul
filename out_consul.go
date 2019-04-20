package main

import (
	"C"
	"fmt"

	"os"
	"unsafe"

	"net"

	"strconv"

	"github.com/fluent/fluent-bit-go/output"
	consul_api "github.com/hashicorp/consul/api"
)

var (
	plugin      Keeper
	hostname    string
	serviceName string
	serviceId   string
	localIp     string
	checkPort   int
	wrapper     = OutputWrapper(&Output{})
)

type Output struct{}

type OutputWrapper interface {
	Register(ctx unsafe.Pointer, name string, desc string) int
	GetConfigKey(ctx unsafe.Pointer, key string) string
	GetConsulConfig() *consul_api.Config
}

func (o *Output) Register(ctx unsafe.Pointer, name string, desc string) int {
	return output.FLBPluginRegister(ctx, name, desc)
}

func (o *Output) GetConfigKey(ctx unsafe.Pointer, key string) string {
	return output.FLBPluginConfigKey(ctx, key)
}

func (o *Output) GetConsulConfig() *consul_api.Config {
	return consul_api.DefaultConfig()
}

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return wrapper.Register(ctx, "consul", "register consul")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	var err error
	css := wrapper.GetConfigKey(ctx, "ConsulServer") // required
	csp := wrapper.GetConfigKey(ctx, "ConsulPort")   // required
	sn := wrapper.GetConfigKey(ctx, "ServiceName")   // required
	cp := wrapper.GetConfigKey(ctx, "CheckPort")     // required
	si := wrapper.GetConfigKey(ctx, "ServiceId")     // default: hostname

	fmt.Printf("[consul-go] plugin parameter consul server = '%s'\n", css)
	fmt.Printf("[consul-go] plugin parameter consul port = '%s'\n", csp)
	fmt.Printf("[consul-go] plugin parameter service name = '%s'\n", sn)
	fmt.Printf("[consul-go] plugin parameter check port = '%s'\n", cp)
	fmt.Printf("[consul-go] plugin parameter service id = '%s'\n", si)

	hostname, err = os.Hostname()
	if err != nil {
		fmt.Printf("[err][init] %+v\n", err)
		return output.FLB_ERROR
	}
	fmt.Printf("[consul-go] plugin hostname = '%s'\n", hostname)

	localIp = getLocalIP()
	if localIp == "" {
		fmt.Printf("[err][init] not found local ip\n")
		return output.FLB_ERROR
	}

	if css == "" || csp == "" || sn == "" || cp == "" {
		fmt.Printf("[err][init] empty config\n")
		return output.FLB_ERROR
	}

	if si == "" {
		si = hostname
	}

	plugin, err = NewKeeper(WithAddrOption(css), WithPortOption(csp),
		WithConfigOption(wrapper.GetConsulConfig()))
	if err != nil {
		fmt.Printf("[err][init] empty consul server or port\n")
		return output.FLB_ERROR
	}
	serviceId = si
	serviceName = sn

	checkPort, err = strconv.Atoi(cp)
	if err != nil {
		fmt.Printf("[err][init] invalid check port\n")
		return output.FLB_ERROR
	}

	if err := plugin.Register(serviceId, serviceName, localIp, checkPort,
		[]string{}, nil); err != nil {
		fmt.Printf("[err][init] %+v\n", err)
		return output.FLB_ERROR
	}
	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	if err := plugin.DeRegister(serviceId); err != nil {
		fmt.Printf("[err][exit] %+v\n", err)
	}
	return output.FLB_OK
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func main() {}
