package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	my_proxy "github.com/hingew/hsfl-master-ai-cloud-engineering/api-gateway/proxy"
	"gopkg.in/yaml.v2"
)

type RoutesConfig struct {
	TemplateingService []string `yaml:"templateing-service"`
	UserService        []string `yaml:"user-service"`
	WebService         []string `yaml:"web-service"`
	CreationService    []string `yaml:"creation-service"`
}

func readRoutesConfig() RoutesConfig {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("Couldn't read Routes Config: %v", err)
	}

	var routes RoutesConfig

	err = yaml.Unmarshal(data, &routes)
	if err != nil {
		log.Fatalf("Error during parsing yaml-file: %v", err)
	}

	return routes
}

func addRoutes(proxy my_proxy.ReverseProxy, endpoint string, routes []string) {
	for _, source_route := range routes {
		target_route := fmt.Sprintf("%s%s", endpoint, source_route)
		proxy.Map(source_route, target_route)
	}
}

func main() {
	templateing_service_endpoint := os.Getenv("TEMPLATE_ENDPOINT")
	web_service_endpoint := os.Getenv("WEB_ENDPOINT")
	creation_service_endpoint := os.Getenv("CREATION_ENDPOINT")
	user_service_endpoint := os.Getenv("USER_ENDPOINT")

	proxy := my_proxy.NewHttpReverseProxy(http.DefaultClient)

	routes := readRoutesConfig()
	addRoutes(proxy, templateing_service_endpoint, routes.TemplateingService)
	addRoutes(proxy, web_service_endpoint, routes.WebService)
	addRoutes(proxy, creation_service_endpoint, routes.CreationService)
	addRoutes(proxy, user_service_endpoint, routes.UserService)

	if err := http.ListenAndServe(":3000", proxy); err != nil {
		log.Fatalf("error while listen and serve: %s", err.Error())
	}
}
