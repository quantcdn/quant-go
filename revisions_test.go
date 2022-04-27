package quant

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddMarkupRevision(t *testing.T) {
	setup()
	mux.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, `{
	"md5": "00000000000000000000000000000000",
	"attachments": [],
	"published": true,
	"quant_filename": "index.html",
	"byte_length": 14,
	"type": "content",
	"url": "\/test-content",
	"revision_number": 1,
	"highest_revision_number": 1,
	"date_timestamp": 0
}`)
	})
	r, err := client.AddMarkupRevision(MarkupRevision{
		Url:             "/test-content",
		FindAttachments: false,
		Content:         []byte("this is a test"),
		Published:       true,
	}, false)

	teardown()

	assert.Empty(t, err)

	// Cut down revision response on success.
	want := Revision{
		Url:       "/test-content",
		Published: true,
	}

	assert.Equal(t, want, r)
}

func TestCreateMarkupSameMD5(t *testing.T) {
	setup()
	mux.HandleFunc("/v1/", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		fmt.Fprintf(w, `{
	"published": true,
	"error": true,
	"byte_length": 1,
	"md5": "00000000000000000000000000000000",
	"attachments": [],
	"errorMsg": "Published version already has md5: 00000000000000000000000000000000",
	"type": "content"
}`)
	})
	r, err := client.AddMarkupRevision(MarkupRevision{
		Url:             "/test-content",
		FindAttachments: false,
		Content:         []byte("this is a test"),
		Published:       true,
	}, false)

	teardown()

	assert.Empty(t, r)
	assert.Equal(t, err.Error(), "Published version already has md5: 00000000000000000000000000000000")
}

func TestGetRevision(t *testing.T) {
	setup()
	mux.HandleFunc("/v1/revisions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "/test-content", r.Header.Get("Quant-Url"))
		fmt.Fprintf(w, `{
	"published": true,
	"published_revision": 1,
	"url": "\/test-content",
	"revisions": {
		"1": {
			"md5": "00000000000000000000000000000000",
			"type": "content",
			"revision_number": 1,
			"byte_length": 14,
			"date_timestamp": 0
		}
	},
	"transitions": [],
	"highest_revision_number": 1,
	"seq_num": 1
}`)
	})

	r, err := client.GetRevision(RevisionQuery{
		Url: "/test-content",
	})

	s := make(map[string]RevisionItem, 1)
	s["1"] = RevisionItem{
		MD5:            "00000000000000000000000000000000",
		Type:           "content",
		RevisionNumber: 1,
		ByteLength:     14,
		DateTimestamp:  0,
	}

	want := Revision{
		Published:         true,
		PublishedRevision: 1,
		Sequence:          1,
		Revisions:         s,
		Url:               "/test-content",
	}

	assert.Empty(t, err)
	assert.Equal(t, want, r)
}

func TestRevisionNotFound(t *testing.T) {
	setup()
	mux.HandleFunc("/v1/revisions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{
	"errorMsg": "Resource not found.",
	"error": true
}`)
	})
	r, err := client.GetRevision(RevisionQuery{
		Url: "/test-content",
	})
	assert.Empty(t, r)
	assert.Equal(t, err.Error(), "Resource not found.")
}
