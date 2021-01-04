# gatekeeper
gatekeeper is an api gateway who controls access to server

## feature
provides the following features

- RateLimit (Leaky bucket) 
- cache (redis)

## Quick Start
### before setup

setting [config.yaml](https://github.com/ohs30359-nobuhara/gatekeeper/blob/master/config.yaml) & 
start redis server
```shell
$ redis-server
```

### start server

```shell
$ go run proxy.go
```
or

```shell
$ go build && ./gatekeeper 
```

access to http://localhost:8000/


## Option
