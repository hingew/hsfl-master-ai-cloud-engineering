package main

import (
	"log"
	"net/http"
	"os"

	my_proxy "github.com/hingew/hsfl-master-ai-cloud-engineering/api-gateway/proxy"
)

func main() {
	// templateing_service_endpoint := os.Getenv("TEMPLATE_ENDPOINT")
	web_service_endpoint := os.Getenv("WEB_ENDPOINT")
	// creation_service_endpoint := os.Getenv("CREATION_ENDPOINT")
	// user_service_endpoint := os.Getenv("USER_ENDPOINT")

	proxy, err := my_proxy.NewReverseProxy("0.0.0.0:3000")
	if err != nil {
		log.Fatalf("could not create reverse proxy: %s", err.Error())
	}

	proxy.Map("/", web_service_endpoint)

	if err := http.ListenAndServe(":3000", proxy); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
