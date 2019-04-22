# fluent-bit-consul

<p align="left">
  <a href="https://circleci.com/gh/gjbae1212/fluent-bit-consul">
    <img src="https://circleci.com/gh/gjbae1212/fluent-bit-consul.svg?style=svg">
  </a>  
  <a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-GREEN.svg" alt="license" /></a>  
  <a href="https://goreportcard.com/report/github.com/gjbae1212/fluent-bit-consul">
  <img src="https://goreportcard.com/badge/github.com/gjbae1212/fluent-bit-consul" alt="Go Report Card" />
  </a>
  <a href="https://codecov.io/gh/gjbae1212/fluent-bit-consul">
    <img src="https://codecov.io/gh/gjbae1212/fluent-bit-consul/branch/master/graph/badge.svg" />
  </a>  
</p>

## OVERVIEW
This project is output plugin for fluent-bit registered service in consul.
Agent for fluent-bit will start and then it is registering consul with specified name in order to consul is periodically check to fluent-bit agent live.
So it is available with monitoring solution same as prometheus, if fluent-bit agent will start with http server on.
Prometheus with consul could do watch fluent-bit agent by metrics, and it is possibly alert for you when fluent-bit would be a fault.  

## Build
A bin directory already has been made binaries for mac, linux.

If you should directly make binaries for mac, linux
```bash
# local machine binary
$ bash make.sh build

# Your machine is mac, and if you should do to retry cross compiling for linux.
# A command in below is required a docker.  
$ bash make.sh build_linux
```

## Usage
### configuration options for fluent-bit.conf
| Key           | Description                                    | Default        |
| ----------------|------------------------------------------------|----------------|
| Name            | output plugin name | consul |
| ConsulServer    | consul server ip | NONE(required) |
| ConsulPort      | consul server port | NONE(required) |
| ServiceName     | consul service name for register | NONE(required) |
| CheckPort       | health check port for consul | NONE(required) |
| ServiceId       | service id for consul registered  | your hostname(default) |

### Example fluent-bit.conf
```conf
[SERVICE]
    HTTP_Server On
    HTTP_Listen 0.0.0.0
    HTTP_PORT 2020
[INPUT]
    Name random
    Tag process.health_check
    Samples -1
    Interval_Sec 10
    Interval_NSec 0
[OUTPUT]
    Name consul
    Tag process.*
    ConsulServer localhost
    ConsulPort 8500
    ServiceName allan
    CheckPort 2020  
```

### Example exec
```bash
$ fluent-bit -c [your config file] -e consul.so 
```
