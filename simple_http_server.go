package goSimpleHttp

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

type SimpleHttpError struct {
	code int
}

func (instance SimpleHttpError) Error() string {
	return fmt.Sprintf("Code : %d", instance.code)
}

type HttpHandler func(w http.ResponseWriter, r *http.Request)

const (
	SimpleHttpServerHandler_NoHandler          int = 405
	SimpleHttpServerHandler_MethodNotSupported int = 404
)

type SimpleHttpServerHandler struct {
	handlers map[string]map[string]HttpHandler
}

func newSimpleHttpServerHandler() *SimpleHttpServerHandler {
	return &SimpleHttpServerHandler{map[string]map[string]HttpHandler{}}
}

func (instance *SimpleHttpServerHandler) addHandler(path string, method string, handler HttpHandler) {
	lowerPath := strings.ToLower(path)
	lowerMethod := strings.ToLower(method)
	if _, ok := instance.handlers[lowerPath]; !ok {
		instance.handlers[lowerPath] = map[string]HttpHandler{}
	}
	instance.handlers[lowerPath][lowerMethod] = handler
}

func (instance *SimpleHttpServerHandler) handlerFor(path string, method string) (HttpHandler, error) {
	lowerPath := strings.ToLower(path)
	lowerMethod := strings.ToLower(method)
	if _, ok := instance.handlers[lowerPath]; !ok {
		fmt.Println("Cannot find the path", instance.handlers)
		return nil, SimpleHttpError{http.StatusNotFound}
	}
	handler, ok := instance.handlers[lowerPath][lowerMethod]
	if !ok {
		fmt.Println("Cannot find the method for the path")
		return nil, SimpleHttpError{http.StatusMethodNotAllowed}
	}
	return handler, nil

}

func (instance *SimpleHttpServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	method := r.Method
	handler, err := instance.handlerFor(path, method)
	if err != nil {
		httpError := err.(SimpleHttpError)
		fmt.Errorf("error encountered %v\n", err)
		w.WriteHeader(httpError.code)
	} else {
		handler(w, r)
	}
}

type SimpleHttpServer struct {
	listener net.Listener
	handler  *SimpleHttpServerHandler
}

func NewSimpleHttpServer(port int, host string) *SimpleHttpServer {
	handler := newSimpleHttpServerHandler()
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
	instance.handler.addHandler(path, "get", handler)
}

func (instance *SimpleHttpServer) Post(path string, handler HttpHandler) {
	instance.handler.addHandler(path, "post", handler)
}

func (instance *SimpleHttpServer) Put(path string, handler HttpHandler) {
	instance.handler.addHandler(path, "put", handler)
}

func (instance *SimpleHttpServer) Delete(path string, handler HttpHandler) {
	instance.handler.addHandler(path, "delete", handler)
}

func (instance *SimpleHttpServer) Stop() {
	instance.listener.Close()
}
