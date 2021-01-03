package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/my/repo/cache"
	"github.com/my/repo/rateLimit"
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

var cache, _ = customCache.Init(context.Background())
var rateLimiter = rateLimit.Init(100)


// proxy handler
func proxy(w http.ResponseWriter, r *http.Request) {
	// rate limit
	if block := rateLimiter.Allow(); block {
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
	schema := "http"
	host := "localhost:3000"

	b, e := ioutil.ReadAll(r.Body)

	if e != nil {
		log.Fatal(e.Error())
	}

	url := fmt.Sprintf("%s://%s%s", schema, host, r.RequestURI)
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
	if err := http.ListenAndServe(":8000", http.HandlerFunc(proxy)); err != nil {
		panic(err.Error())
	}
}