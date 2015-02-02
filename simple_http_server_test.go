package go_Simple_Http_Server

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
			server SimpleHttpServer
		)

		g.Before(func() {
			port = 5000
			host = "127.0.0.1"
			url = func(path string) string {
				return fmt.Sprintf("http://%s:%d%s", host, port, path)
			}
			server = NewSimpleHttpServer(port, host)
		})

		g.It("Supports GET", func() {
			server.Get("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			resp, _ := http.Get(url("/"))
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			server.Stop()
			_, err := http.Get(url("/")
			assert.True(t, err != nil)
		})

		g.It("Supports POST", func() {
			server.Post("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			body := strings.NewReader("{}")
			resp, _ := http.Post(url("/"), "application/json", body)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			server.Stop()
			_, err := http.Post(url("/"), "application/json", body)
			assert.True(t, err != nil)
		})
	})
}
