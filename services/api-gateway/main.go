package main

import (
	"log"
	"net/http"
)

func main() {
	proxy := httpproxy.New()
	proxy.Map("/", "http://web-service:3000")
	log.Fatal(http.ListenAndServe("0.0.0.0:3000", proxy))
}
