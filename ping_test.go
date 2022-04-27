package quant

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPing(t *testing.T) {
	setup()
	mux.HandleFunc(apiBase+"/ping", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, os.Getenv("QUANT_CLIENT_ID"), r.Header.Get("Quant-Customer"))
		assert.Equal(t, os.Getenv("QUANT_PROJECT"), r.Header.Get("Quant-Project"))
	})

	client.Ping()
	teardown()
}
