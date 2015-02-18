package goSimpleHttp

import (
	"errors"
	"fmt"
	"github.com/REAANDREW/gopubsubio"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	STOPPED_EVENT_KEY string = "stopped"
)

type SimpleHttpError struct {
	code int
}

func (instance SimpleHttpError) Error() string {
	return fmt.Sprintf("Code : %d", instance.code)
}

type HttpHandler func(w http.ResponseWriter, r *http.Request)

const (
	SimpleHttpServerHandler_NoHandler          int    = 405
	SimpleHttpServerHandler_MethodNotSupported int    = 404
	Any                                        string = "*"
)

type SimpleHttpServerHandler struct {
	handlers map[string]HttpHandler
}

func newSimpleHttpServerHandler() SimpleHttpServerHandler {
	return SimpleHttpServerHandler{map[string]HttpHandler{}}
}

func (instance *SimpleHttpServerHandler) addAnyHandler(handler HttpHandler) {
	instance.handlers[Any] = handler
}

func (instance *SimpleHttpServerHandler) addHandler(path string, method string, handler HttpHandler) {
	lowerPath := strings.ToLower(path)
	lowerMethod := strings.ToLower(method)
	key := lowerPath + "_" + lowerMethod
	instance.handlers[key] = handler
}

func (instance *SimpleHttpServerHandler) handlerFor(path string, method string) (HttpHandler, error) {
	if handler, ok := instance.handlers[Any]; ok {
		return handler, nil
	}
	lowerPath := strings.ToLower(path)
	lowerMethod := strings.ToLower(method)
	key := lowerPath + "_" + lowerMethod
	if handler, ok := instance.handlers[key]; !ok {
		return nil, SimpleHttpError{http.StatusNotFound}
	} else {
		return handler, nil
	}

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
	listener  net.Listener
	handler   *SimpleHttpServerHandler
	publisher gopubsubio.Publisher
	mux       *http.ServeMux
	server    *http.Server
}

func NewSimpleHttpServer(port int, host string) *SimpleHttpServer {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Printf("error %v\n", err)
		panic(errors.New("A listener cannot be setup"))
	}
	handler := newSimpleHttpServerHandler()
	mux := http.NewServeMux()
	mux.Handle("/", &handler)
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", host, port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return &SimpleHttpServer{ln, &handler, gopubsubio.NewPublisher(), mux, server}
}

func (instance *SimpleHttpServer) Start() {
	go func() {
		err := instance.server.Serve(instance.listener)
		if err != nil {
			fmt.Errorf("Error encountered here starting the http server: %v")
		}
	}()
}

func (instance *SimpleHttpServer) publishOnStopped() {
	instance.publisher.Publish(STOPPED_EVENT_KEY, nil)
}

func (instance *SimpleHttpServer) OnStopped(delegate func()) {
	subscriber := gopubsubio.NewSubscriber(func(message interface{}) {
		delegate()
	})
	instance.publisher.Subscribe(STOPPED_EVENT_KEY, subscriber)
}

func createHandlerFor(method string, handler HttpHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ToLower(r.Method) == strings.ToLower(method) {
			handler(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
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

func (instance *SimpleHttpServer) Any(handler HttpHandler) {
	instance.handler.addAnyHandler(handler)
}

func (instance *SimpleHttpServer) Stop() {
	instance.listener.Close()
	instance.publishOnStopped()
}
