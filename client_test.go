package quant

import (
	"net/http"
	"net/http/httptest"
	"os"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	token := os.Getenv("QUANT_TOKEN")
	project := os.Getenv("QUANT_PROJECT")
	client_id := os.Getenv("QUANT_CLIENT_ID")
	client = NewClient(token, client_id, project)
	client.Host = server.URL
}

func teardown() {
	server.Close()
}
