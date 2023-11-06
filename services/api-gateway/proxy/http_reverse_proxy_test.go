package my_proxy

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	_mock "github.com/hingew/hsfl-master-ai-cloud-engineering/api-gateway/_mocks"
	"gotest.tools/v3/assert"
)

func TestController(t *testing.T) {
	ctrl := gomock.NewController(t)
	client := _mock.NewMockHttpClient(ctrl)

	t.Run("Add Routes with Map function", func(t *testing.T) {
		// arrange
		proxy := NewHttpReverseProxy(client)

		// act
		proxy.Map("/test", "http://endpoint:3000/test")

		// assert
		assert.Equal(t, 1, len(proxy.routes))
		assert.Assert(t, proxy.routes["/test"] == "http://endpoint:3000/test")
	})

	t.Run("Route not supported", func(t *testing.T) {

		r := httptest.NewRequest(http.MethodGet, "/test", nil)

		response := &http.Response{
			StatusCode: http.StatusNotFound,
			Body:       ioutil.NopCloser(strings.NewReader("")),
		}
		client.EXPECT().Do(r).Return(response, nil).Times(0)

		t.Run("No route supported", func(t *testing.T) {
			// arrange
			proxy := NewHttpReverseProxy(client)

			w := httptest.NewRecorder()

			// act
			proxy.ServeHTTP(w, r)

			// assert
			errorMsg := "Could not found: /test\nSupported URLs:\n"
			assert.Equal(t, http.StatusNotFound, w.Code)
			assert.Equal(t, errorMsg, w.Body.String())
		})

		t.Run("Wrong routes supported", func(t *testing.T) {
			// arrange
			proxy := NewHttpReverseProxy(client)
			proxy.routes["/test2"] = "http://newEndpoint:3000/test2"
			proxy.routes["/test3"] = "http://newEndpoint2:3000/test3"

			w := httptest.NewRecorder()

			// act
			proxy.ServeHTTP(w, r)

			// assert
			errorMsg := "Could not found: /test\nSupported URLs:\n\t/test2\n\t/test3\n"
			assert.Equal(t, http.StatusNotFound, w.Code)
			assert.Equal(t, errorMsg, w.Body.String())
		})
	})

	t.Run("Route supported", func(t *testing.T) {
		proxy := NewHttpReverseProxy(client)
		proxy.routes["/test2"] = "http://newEndpoint:3000/test2"
		proxy.routes["/test3"] = "http://newEndpoint2:3000/test3"

		t.Run("Client answers with http.StatusAccepted", func(t *testing.T) {
			// arrange
			proxyReq := httptest.NewRequest(http.MethodGet, "/test3", nil)
			serverReq := httptest.NewRequest(http.MethodGet, "http://newEndpoint2:3000/test3", nil)
			serverReq.RequestURI = ""

			w := httptest.NewRecorder()

			response := &http.Response{
				StatusCode: http.StatusOK,
				Body:       ioutil.NopCloser(strings.NewReader("")),
			}
			client.EXPECT().Do(serverReq).Return(response, nil).Times(1)

			// act
			proxy.ServeHTTP(w, proxyReq)

			// assert
			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "", w.Body.String())
		})

		t.Run("Client answers with http.StatusInternalServerError", func(t *testing.T) {
			// arrange
			proxyReq := httptest.NewRequest(http.MethodGet, "/test2", nil)
			serverReq := httptest.NewRequest(http.MethodGet, "http://newEndpoint:3000/test2", nil)
			serverReq.RequestURI = ""

			w := httptest.NewRecorder()

			err := fmt.Errorf("ErrorMsg")
			client.EXPECT().Do(serverReq).Return(nil, err).Times(1)

			// act
			proxy.ServeHTTP(w, proxyReq)

			// assert
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, err.Error(), w.Body.String())
		})
	})
}
