package goSimpleHttp

import (
	"fmt"
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	port := 7000
	host := "127.0.0.1"
	server := NewSimpleHttpServer(port, host)
	server.Get("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello world!")
	})
	server.Start()
	os.Exit(m.Run())
	server.Stop()
}

func Benchmark_Server(b *testing.B) {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://127.0.0.1:7000", nil)
	for i := 0; i < b.N; i++ {
		resp, _ := client.Do(r)
		_, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
		}
	}
}

func Benchmark_ServerInParallel(b *testing.B) {
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "http://127.0.0.1:7000", nil)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			resp, _ := client.Do(r)
			_, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				resp.Body.Close()
			}
		}
	})
}

func Test_SimpleHttpServer(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("SimpleHttpServer", func() {

		var (
			port int
			host string
			url  func(path string) string
		)

		g.BeforeEach(func() {
			port = 6000
			host = "127.0.0.1"
			url = func(path string) string {
				return fmt.Sprintf("http://%s:%d%s", host, port, path)
			}
		})

		g.It("Supports GET", func() {
			client := http.Client{}
			server := NewSimpleHttpServer(port, host)
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
			client := http.Client{}
			server := NewSimpleHttpServer(port, host)
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
			client := http.Client{}
			server := NewSimpleHttpServer(port, host)
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
			client := http.Client{}
			server := NewSimpleHttpServer(port, host)
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

		g.It("Supports setting a handler for any path or method", func() {
			client := http.Client{}
			server := NewSimpleHttpServer(port, host)
			r, _ := http.NewRequest("GET", url("/"), nil)
			server.Any(func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			resp, _ := client.Do(r)
			server.Stop()
			g.Assert(resp.StatusCode).Equal(http.StatusOK)
		})

		g.It("Returns Not Found when path not configured", func() {
			client := http.Client{}
			server := NewSimpleHttpServer(port, host)
			r, _ := http.NewRequest("GET", url("/"), nil)
			server.Start()
			resp, _ := client.Do(r)
			server.Stop()
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})

	})
}
