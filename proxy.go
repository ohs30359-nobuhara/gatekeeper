package main

import (
	"bytes"
	"fmt"
	"github.com/my/repo/cache"
	"github.com/my/repo/rateLimit"
	"github.com/my/repo/util"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"unsafe"
)

type DataSet struct{
	Status int
	Body   []byte
	Header http.Header
}

type ProxyConfig struct {
	schema string
	host string
}

var cache *customCache.Redis
var config *util.Config
var rateLimiter *rateLimit.RateLimit
var proxyConfig *ProxyConfig

// proxyHandler handler
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// rate limit
	if arrow := rateLimiter.Allow(); !arrow {
		http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
		return
	}

	var data DataSet

	key := r.RequestURI

	cr := cloneRequest(r)

	if r.Method == "GET" {
		// cache read only get request
		if v, hit := cache.Get(key); hit {
			writeBody(w, *(*[]byte)(unsafe.Pointer(&v)))
			return
		}

		data = curl(cr)
		cache.Set(key, *(*string)(unsafe.Pointer(&data.Body)), time.Second * 1)
	} else {
		data = curl(cr)
	}

	if data.Status >= http.StatusBadRequest {
		http.Error(w, http.StatusText(data.Status), data.Status)
		return
	}
	writeBody(w, data.Body)
}

// overwrite client request host
func cloneRequest(r *http.Request) *http.Request {
	b, e := ioutil.ReadAll(r.Body)

	if e != nil {
		log.Fatal(e.Error())
	}

	url := fmt.Sprintf("%s://%s%s", proxyConfig.schema, proxyConfig.host, r.RequestURI)
	req, e:= http.NewRequest(r.Method, url, bytes.NewReader(b))

	if e != nil {
		log.Print(e.Error())
	}

	req.Header = make(http.Header)
	for h, val := range r.Header {
		req.Header[h] = val
	}

	return req
}

// http request
func curl(r *http.Request) DataSet {
	var dataSet DataSet

	res, e := (&http.Client {
		// timeout is 1 sec
		Timeout: time.Second * 1,
	}).Do(r)

	if e != nil {
		log.Print(e.Error())

		// TODO: this case return 400 ... ?
		dataSet.Status = http.StatusBadRequest
		return dataSet
	}

	defer res.Body.Close()

	dataSet.Status = res.StatusCode
	dataSet.Header = res.Header

	var buf bytes.Buffer
	_, e = io.Copy(&buf, res.Body)

	if e != nil {
		dataSet.Status = http.StatusInternalServerError
		return dataSet
	}

	dataSet.Body = buf.Bytes()

	return dataSet
}

func writeBody(w http.ResponseWriter, body []byte) {
	w.WriteHeader(http.StatusOK)

	if _, e := w.Write(body); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
}

func main() {
	// load config
	config = util.ConfigLoad()

	if config == nil {
		panic("can't read config file...")
	}

	// set up redis
	c, e := customCache.Init(customCache.Option {
		Host: config.Redis.Host,
		Port: config.Redis.Port,
		Password: config.Redis.Password,
		Db: config.Redis.No,
		TimeoutMs: struct {
			Read  int
			Write int
		}{Read: config.Redis.TimeoutMs.Read, Write: config.Redis.TimeoutMs.Write},
	})

	if e != nil {
		panic("can't listen redis server")
	}

	cache = c

	// setup rate limit
	rateLimiter = rateLimit.Init(config.Proxy.RateLimit)

	// listen server
	p := ":" + config.Server.Port
	log.Println("server listen localhost" + p)

	// setup proxy config
	proxyConfig = &ProxyConfig{
		schema: config.Proxy.Target.Schema,
		host: config.Proxy.Target.Host + ":" + config.Proxy.Target.Port,
	}

	if err := http.ListenAndServe(p, http.HandlerFunc(proxyHandler)); err != nil {
		panic(err.Error())
	}
}