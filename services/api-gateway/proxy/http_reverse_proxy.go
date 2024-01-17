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

	"golang.org/x/sync/singleflight"
)

type Config struct {
	path       *regexp.Regexp
	upstream   *url.URL
	coalescing bool
}

type HttpReverseProxy struct {
	client  HttpClient
	configs []Config
	sfGroup *singleflight.Group
}

func NewHttpReverseProxy(client HttpClient) *HttpReverseProxy {
	return &HttpReverseProxy{client: client, configs: make([]Config, 0), sfGroup: &singleflight.Group{}}
}

func (reverseProxy *HttpReverseProxy) Map(path string, upstream string, coalescing bool) {
	pathRegex, err := regexp.Compile(path)
	if err != nil {
		log.Fatalf("The path %s is not a valid regular expression", path)
	}

	upstreamUrl, err := url.Parse(upstream)
	if err != nil {
		log.Fatalf("The url %s is not a valid url", upstream)
	}

	reverseProxy.configs = append(reverseProxy.configs, Config{pathRegex, upstreamUrl, false})
	log.Printf("Requests at %s will be proxied to %s", path, upstream)
}

func (reverseProxy *HttpReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	inputRoute := r.URL.Path

	config, err := reverseProxy.evaluateEndpointServer(inputRoute)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "%s", err.Error())
		log.Print(err.Error())
		return
	}

	reverseProxy.modifyRequest(r, config.upstream)
	log.Println(fmt.Sprintf("%s -> %s", inputRoute, r.URL.String()))

	if config.coalescing {
		r.URL.Path = removeCoalescingPrefix(r.URL.Path)
		reverseProxy.sfGroup.Do(config.upstream.String(), func() (interface{}, error) {
			reverseProxy.doRequest(r, w)
			return nil, nil
		})
	} else {
		reverseProxy.doRequest(r, w)
	}
}

func (reverseProxy *HttpReverseProxy) doRequest(req *http.Request, rw http.ResponseWriter) {
	response, err := reverseProxy.client.Do(req)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%s", err)
		return
	}

	reverseProxy.copyResponseHeader(response, rw)
	reverseProxy.copyResponseStreamBody(response, rw)
}

func (p *HttpReverseProxy) evaluateEndpointServer(sourceUrl string) (*Config, error) {
	for _, config := range p.configs {
		if config.path.MatchString(sourceUrl) {
			return &config, nil
		}
	}

	errorMsg := fmt.Sprintf("Could not found: %s\n", sourceUrl)
	errorMsg += "Supported URLs:\n"
	for _, config := range p.configs {
		errorMsg += fmt.Sprintf("\t%s\n", config.upstream.String())
	}
	return nil, fmt.Errorf(errorMsg)

}

func (reverseProxy *HttpReverseProxy) containsId(str string) bool {
	regex := regexp.MustCompile(`/\d+`)
	return regex.MatchString(str)
}

func (reverseProxy *HttpReverseProxy) isGatewayCoalescingRoute(str string) bool {
	return strings.Contains(str, "/gateway_coalescing")
}

func removeCoalescingPrefix(str string) string {
	return strings.Replace(str, "/gateway_coalescing", "", 1)
}

func (reverseProxy *HttpReverseProxy) modifyRequest(req *http.Request, endpoint *url.URL) {
	req.Host = endpoint.Host
	req.URL.Host = endpoint.Host
	req.URL.Scheme = endpoint.Scheme
	req.RequestURI = ""
	//req.URL.Path = req.UR.Path
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
