package my_proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	_mock "github.com/hingew/hsfl-master-ai-cloud-engineering/api-gateway/_mocks"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert"
)

func TestReverseProxy(t *testing.T) {
	ctrl := gomock.NewController(t)
	httpClient := _mock.NewMockHttpClient(ctrl)

	t.Run("Add Routes with Map function", func(t *testing.T) {
		// arrange
		proxy := NewHttpReverseProxy(httpClient)
		url := "http://endpoint:3000/test"

		// act
		proxy.Map("/test", url, false)

		// assert
		assert.Equal(t, 1, len(proxy.configs))

		config := proxy.configs[0]
		assert.Equal(t, config.coalescing, false)
		assert.Equal(t, config.upstream.String(), url)
		assert.Equal(t, config.path.String(), "/test")
	})

	t.Run("Route not supported", func(t *testing.T) {

		// arrange
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		proxy := NewHttpReverseProxy(httpClient)

		response := http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}

		httpClient.EXPECT().Do(r).Return(&response, nil).Times(0)

		// act
		proxy.ServeHTTP(w, r)

		// assert
		errorMsg := "Could not found: /test\nSupported URLs:\n"
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Equal(t, errorMsg, w.Body.String())
	})

	t.Run("Should route requests to the right upstream by exact match", func(t *testing.T) {
		proxy := NewHttpReverseProxy(httpClient)
		proxy.Map("/exact-match", "http://endpoint:3000", false)
		proxy.Map("/second", "http://endpoint:3001", false)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}

		proxyReq := httptest.NewRequest(http.MethodGet, "/exact-match", nil)
		serverReq := httptest.NewRequest(http.MethodGet, "http://endpoint:3000/exact-match", nil)
		serverReq.RequestURI = ""

		w := httptest.NewRecorder()

		httpClient.EXPECT().Do(serverReq).Return(response, nil).Times(1)

		// act
		proxy.ServeHTTP(w, proxyReq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "", w.Body.String())

	})

	t.Run("Should route requests to the right upstream by regex", func(t *testing.T) {
		proxy := NewHttpReverseProxy(httpClient)
		proxy.Map("/*", "http://endpoint:3000", false)
		proxy.Map("/exact-match", "http://endpoint-2:3000", false)

		response := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}

		proxyReq := httptest.NewRequest(http.MethodGet, "/wild/card/match", nil)
		serverReq := httptest.NewRequest(http.MethodGet, "http://endpoint:3000/wild/card/match", nil)
		serverReq.RequestURI = ""

		w := httptest.NewRecorder()

		httpClient.EXPECT().Do(serverReq).Return(response, nil).Times(1)

		// act
		proxy.ServeHTTP(w, proxyReq)

		// assert
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, "", w.Body.String())

	})
	t.Run("Client answers with http.StatusInternalServerError", func(t *testing.T) {
		proxy := NewHttpReverseProxy(httpClient)
		proxy.Map("/test", "http://endpoint:3000", false)

		// arrange
		proxyReq := httptest.NewRequest(http.MethodGet, "/test", nil)
		serverReq := httptest.NewRequest(http.MethodGet, "http://endpoint:3000/test", nil)
		serverReq.RequestURI = ""

		w := httptest.NewRecorder()

		err := fmt.Errorf("ErrorMsg")
		httpClient.EXPECT().Do(serverReq).Return(nil, err).Times(1)

		// act
		proxy.ServeHTTP(w, proxyReq)

		// assert
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, err.Error(), w.Body.String())
	})
}
