package goSimpleHttp

import (
	"fmt"
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_SimpleHttpServer(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("SimpleHttpServer", func() {

		var (
			port   int
			host   string
			url    func(path string) string
			server *SimpleHttpServer
			client *http.Client
		)

		g.BeforeEach(func() {
			port = 5000
			host = "127.0.0.1"
			url = func(path string) string {
				return fmt.Sprintf("http://%s:%d%s", host, port, path)
			}
			server = NewSimpleHttpServer(port, host)
			client = &http.Client{}
		})

		g.It("Supports GET", func() {
			r, _ := http.NewRequest("GET", url("/"), nil)
			server.Get("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			resp, _ := client.Do(r)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			server.Stop()
			_, err := client.Do(r)
			assert.True(t, err != nil)
		})

		g.It("Supports POST", func() {
			body := strings.NewReader("{}")
			r, _ := http.NewRequest("POST", url("/"), body)
			server.Post("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			resp, _ := client.Do(r)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			server.Stop()
			_, err := client.Do(r)
			assert.True(t, err != nil)
		})

		g.It("Supports PUT", func() {
			body := strings.NewReader("{}")
			r, _ := http.NewRequest("PUT", url("/"), body)
			server.Put("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			resp, _ := client.Do(r)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			server.Stop()
			_, err := client.Do(r)
			assert.True(t, err != nil)
		})

		g.It("Supports DELETE", func() {
			body := strings.NewReader("{}")
			r, _ := http.NewRequest("DELETE", url("/"), body)
			server.Delete("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			resp, _ := client.Do(r)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			server.Stop()
			_, err := client.Do(r)
			assert.True(t, err != nil)
		})

		g.It("Returns Not Found when path not configured", func() {
			r, _ := http.NewRequest("GET", url("/"), nil)
			server.Start()
			resp, _ := client.Do(r)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
			server.Stop()
		})

		g.It("Returns Method Not Implemented when handler for path not configured")

	})
}
