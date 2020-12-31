package main

import (
	"net/http"
	"log"
	"net/http/httputil"
)

func main() {
	d := func(r *http.Request) {
		r.URL.Scheme = "http"
		r.URL.Host = ":3000"
	}

	server := http.Server{
		Addr: ":8000",
		Handler: &httputil.ReverseProxy{Director: d},
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}