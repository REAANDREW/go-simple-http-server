package go_Simple_Http_Server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

type HttpHandler func(w http.ResponseWriter, r *http.Request)

const (
	SimpleHttpServerHandler_NoHandler          int = 1
	SimpleHttpServerHandler_MethodNotSupported int = 2
)

type SimpleHttpServerHandler struct {
	handlers map[string]map[string]HttpHandler
}

func NewSimpleHttpServerHandler() *SimpleHttpServerHandler {
	return &SimpleHttpServerHandler{map[string]map[string]HttpHandler{}}
}

func (instance *SimpleHttpServerHandler) AddHandler(path string, method string, handler HttpHandler) {
	lowerPath := strings.ToLower(path)
	lowerMethod := strings.ToLower(method)
	if _, ok := instance.handlers[lowerPath]; !ok {
		instance.handlers[lowerPath] = map[string]HttpHandler{}
	}
	instance.handlers[lowerPath][lowerMethod] = handler
}

func (instance *SimpleHttpServerHandler) HandlerFor(path string, method string) (HttpHandler, error) {
	lowerPath := strings.ToLower(path)
	lowerMethod := strings.ToLower(method)
	if _, ok := instance.handlers[lowerPath]; !ok {
		return nil, errors.New("no handler for path")
	}
	handler, ok := instance.handlers[lowerPath][lowerMethod]
	if !ok {
		return nil, errors.New("no handler for method")
	}
	return handler, nil

}

func (instance *SimpleHttpServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	handler, err := instance.HandlerFor(path, method)
	if err != nil {
		fmt.Errorf("error encountered %v\n", err)
	} else {
		handler(w, r)
	}
}

type SimpleHttpServer struct {
	listener net.Listener
	handler  *SimpleHttpServerHandler
}

func NewSimpleHttpServer(port int, host string) *SimpleHttpServer {
	handler := NewSimpleHttpServerHandler()
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	fmt.Errorf("error encountered %v", err)
	return &SimpleHttpServer{ln, handler}
}

func (instance *SimpleHttpServer) Start() {
	go func() {
		http.Serve(instance.listener, instance.handler)
	}()
}

func (instance *SimpleHttpServer) Get(path string, handler HttpHandler) {
	instance.handler.AddHandler(path, "get", handler)
}

func (instance *SimpleHttpServer) Stop() {
	instance.listener.Close()
}
