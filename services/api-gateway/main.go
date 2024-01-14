package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	my_proxy "github.com/hingew/hsfl-master-ai-cloud-engineering/api-gateway/proxy"
)

type RoutesConfig struct {
	TemplateingService []string
	UserService        []string
	WebService         []string
	CreationService    []string
}

func readRoutesFromEnv() RoutesConfig {
	var routes RoutesConfig

	routes.CreationService = readRouteEnv("CREATION_ROUTES")
	routes.WebService = readRouteEnv("WEB_ROUTES")
	routes.TemplateingService = readRouteEnv("TEMPLATE_ROUTES")
	routes.UserService = readRouteEnv("USER_ROUTES")

	return routes
}

func readRouteEnv(env string) []string {
	value := os.Getenv(env)
	values := strings.Split(value, ";")
	return values
}

//func addRoutes(proxy my_proxy.ReverseProxy, endpoint string, routes []string) {
//	for _, source_route := range routes {
//		target_route := fmt.Sprintf("%s%s", endpoint, source_route)
//		proxy.Map(source_route, target_route)
//	}
//}

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
