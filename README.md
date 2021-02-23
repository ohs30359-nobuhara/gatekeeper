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

access to http://localhost:8000/proxy/${origin_path}

## EntryPoint 
### /health
health check entry point.
return status code 200 & body "OK"

### /metrics
prometheus entry point.

### /proxy/${origin_path}
proxy entry point.

return proxied response body & header.

## Option
