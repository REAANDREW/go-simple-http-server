package go_Simple_Http_Server

import (
	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func Test_SimpleHttpServer(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("SimpleHttpServer", func() {
		g.It("Supports GET", func() {
			server := NewSimpleHttpServer(5000, "127.0.0.1")
			server.Get("/", func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, "Hello world!")
			})
			server.Start()
			resp, _ := http.Get("http://127.0.0.1:5000")
			assert.Equal(t, http.StatusOK, resp.StatusCode)
			server.Stop()
			_, err := http.Get("http://127.0.0.1:5000")
			assert.True(t, err != nil)
		})
	})
}
