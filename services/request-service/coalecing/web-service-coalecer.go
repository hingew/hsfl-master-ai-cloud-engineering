package coalecing

import (
	"io"
	"net/http"
	"os"

	"golang.org/x/sync/singleflight"
)

type WebServiceCoalecer struct {
	sfGroup  *singleflight.Group
	client   http.Client
	endpoint string
}

func NewWebServiceCoalecer() *WebServiceCoalecer {
	endpoint := os.Getenv("WEB_SERVICE_ENDPOINT")
	g := &singleflight.Group{}

	return &WebServiceCoalecer{
		g,
		http.Client{},
		endpoint,
	}
}

func (server *WebServiceCoalecer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res, err, _ := server.sfGroup.Do(r.URL.Path, func() (interface{}, error) {
		req, _ := http.NewRequest(r.Method, server.endpoint+r.URL.Path, r.Body)
		req.Header = r.Header

		resp, err := server.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res != nil {
		io.WriteString(w, string(res.([]byte)))
	}
}
