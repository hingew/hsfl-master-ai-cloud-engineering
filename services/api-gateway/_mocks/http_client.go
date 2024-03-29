// Code generated by MockGen. DO NOT EDIT.
// Source: proxy/reverse_proxy.go
//
// Generated by this command:
//
//	mockgen -package=mocks -destination=_mocks/http_client.go -source=proxy/reverse_proxy.go
//
// Package mocks is a generated GoMock package.
package mocks

import (
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockReverseProxy is a mock of ReverseProxy interface.
type MockReverseProxy struct {
	ctrl     *gomock.Controller
	recorder *MockReverseProxyMockRecorder
}

// MockReverseProxyMockRecorder is the mock recorder for MockReverseProxy.
type MockReverseProxyMockRecorder struct {
	mock *MockReverseProxy
}

// NewMockReverseProxy creates a new mock instance.
func NewMockReverseProxy(ctrl *gomock.Controller) *MockReverseProxy {
	mock := &MockReverseProxy{ctrl: ctrl}
	mock.recorder = &MockReverseProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReverseProxy) EXPECT() *MockReverseProxyMockRecorder {
	return m.recorder
}

// Map mocks base method.
func (m *MockReverseProxy) Map(sourcePath, destinationPath string, coalescing bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Map", sourcePath, destinationPath, coalescing)
}

// Map indicates an expected call of Map.
func (mr *MockReverseProxyMockRecorder) Map(sourcePath, destinationPath, coalescing any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Map", reflect.TypeOf((*MockReverseProxy)(nil).Map), sourcePath, destinationPath, coalescing)
}

// ServeHTTP mocks base method.
func (m *MockReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ServeHTTP", w, r)
}

// ServeHTTP indicates an expected call of ServeHTTP.
func (mr *MockReverseProxyMockRecorder) ServeHTTP(w, r any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServeHTTP", reflect.TypeOf((*MockReverseProxy)(nil).ServeHTTP), w, r)
}

// MockHttpClient is a mock of HttpClient interface.
type MockHttpClient struct {
	ctrl     *gomock.Controller
	recorder *MockHttpClientMockRecorder
}

// MockHttpClientMockRecorder is the mock recorder for MockHttpClient.
type MockHttpClientMockRecorder struct {
	mock *MockHttpClient
}

// NewMockHttpClient creates a new mock instance.
func NewMockHttpClient(ctrl *gomock.Controller) *MockHttpClient {
	mock := &MockHttpClient{ctrl: ctrl}
	mock.recorder = &MockHttpClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHttpClient) EXPECT() *MockHttpClientMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockHttpClient) Do(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockHttpClientMockRecorder) Do(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockHttpClient)(nil).Do), arg0)
}
