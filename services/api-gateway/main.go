package main

import (
	"log"
	"net/http"
	"os"

	my_proxy "github.com/hingew/hsfl-master-ai-cloud-engineering/api-gateway/proxy"
)

type RoutesConfig struct {
	TemplateingService []string
	UserService        []string
	WebService         []string
	CreationService    []string
}

func main() {
	templateing_service_endpoint := os.Getenv("TEMPLATE_ENDPOINT")
	web_service_endpoint := os.Getenv("WEB_ENDPOINT")
	creation_service_endpoint := os.Getenv("CREATION_ENDPOINT")
	user_service_endpoint := os.Getenv("USER_ENDPOINT")

	proxy := my_proxy.NewHttpReverseProxy(http.DefaultClient)

	// Add configurations to the proxy
	proxy.Map("/auth/*", user_service_endpoint, false)
	proxy.Map("/api/templates/*", templateing_service_endpoint, false)
	proxy.Map("/api/render/*", creation_service_endpoint, false)
	proxy.Map("/*", web_service_endpoint, false)

	if err := http.ListenAndServe(":3000", proxy); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
