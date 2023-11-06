package my_proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type HttpReverseProxy struct {
	client HttpClient
	routes map[string]string
}

func NewHttpReverseProxy(client HttpClient) *HttpReverseProxy {
	return &HttpReverseProxy{client: client, routes: make(map[string]string)}
}

func (reverseProxy *HttpReverseProxy) Map(sourcePath string, destinationPath string) {
	reverseProxy.routes[sourcePath] = destinationPath
}

func (reverseProxy *HttpReverseProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	endpointServerURL, err := reverseProxy.evaluateEndpointServer(req.URL.Path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "%s", err.Error())
		return
	}

	reverseProxy.modifyRequest(req, endpointServerURL)

	response, err := reverseProxy.client.Do(req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", err)
		return
	}

	reverseProxy.copyResponseHeader(response, rw)
	reverseProxy.copyResponseStreamBody(response, rw)
}

func (reverseProxy *HttpReverseProxy) evaluateEndpointServer(sourceUrl string) (*url.URL, error) {
	ok, rawDestinationURL := reverseProxy.matchSupportedRoute(sourceUrl)
	if ok != true {
		errorMsg := fmt.Sprintf("Could not found: %s\n", sourceUrl)
		errorMsg += "Supported URLs:\n"
		for key, _ := range reverseProxy.routes {
			errorMsg += fmt.Sprintf("\t%s\n", key)
		}
		return nil, fmt.Errorf(errorMsg)
	}
	log.Println(rawDestinationURL)

	return url.Parse(*rawDestinationURL)
}

func (reverseProxy *HttpReverseProxy) matchSupportedRoute(source_route string) (bool, *string) {
	if !containsId(source_route) {
		destination_route, ok := reverseProxy.routes[source_route]
		return ok, &destination_route
	}

	for key, value := range reverseProxy.routes {
		if !strings.Contains(key, ":id") {
			continue
		}
		expression := strings.Replace(key, ":id", `(\d+)`, 1)
		reg := regexp.MustCompile(expression)
		matches := reg.FindStringSubmatch(source_route)
		if len(matches) == 2 {
			id := matches[1]
			destinationRoute := strings.Replace(value, ":id", id, 1)
			return true, &destinationRoute
		}
	}

	return false, nil
}

func containsId(str string) bool {
	regex := regexp.MustCompile(`/\d+`)
	return regex.MatchString(str)
}

func (reverseProxy *HttpReverseProxy) modifyRequest(req *http.Request, endpoint *url.URL) {
	req.Host = endpoint.Host
	req.URL.Host = endpoint.Host
	req.URL.Scheme = endpoint.Scheme
	req.RequestURI = ""
}

func (reverseProxy *HttpReverseProxy) copyResponseHeader(response *http.Response, rw http.ResponseWriter) {
	for key, values := range response.Header {
		for _, value := range values {
			rw.Header().Set(key, value)
		}
	}
}

func (reverseProxy *HttpReverseProxy) copyResponseStreamBody(response *http.Response, rw http.ResponseWriter) {
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-time.Tick(10 * time.Millisecond):
				rw.(http.Flusher).Flush()
			case <-done:
				return
			}
		}
	}()

	rw.WriteHeader(response.StatusCode)
	io.Copy(rw, response.Body)

	close(done)
}
