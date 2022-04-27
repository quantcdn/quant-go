package quant

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddFormFull(t *testing.T) {
	setup()
	mux.HandleFunc("/v1/form", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method, "Expected method 'POST', got %s", r.Method)
		fmt.Fprintf(w, `{
	"attachments": {
		"forms": [],
		"js": [],
		"css": [],
		"media": {
			"video": [],
			"audio": [],
			"images": [],
			"documents": []
		}
	},
	"quant_filename": "index.html",
	"form_enabled": true,
	"date_timestamp": 0,
	"md5": "00000000000000000000000000000000",
	"published": true,
	"byte_length": 14,
	"revision_number": 1,
	"url": "\/test-content",
	"form_config": {
		"error_message_mandatory": "Not much success",
		"target_url": "\/test-content",
		"notifications": {
			"email": {
				"cc": "test@test.com",
				"to": "test@test.com",
				"from": "from@test.com",
				"enabled": true,
				"options": {
					"include_results": true,
					"text_only": true
				},
				"subject": "test"
			},
			"slack": {
				"enabled": true,
				"webhook": "http:\/\/slack.webhook\/test"
			}
		},
		"error_message_generic": "Not much success",
		"mandatory_fields": [
			"must"
		],
		"remove_fields": [
			"remove"
		],
		"success_message": "Much success",
		"honeypot_fields": [
			"hide"
		]
	},
	"highest_revision_number": 3,
	"type": "content"
}`)
	})

	want := Form{
		Url:     "/test-content",
		Enabled: true,
		Config: FormConfig{
			Target:                "/test-content",
			HoneypotFields:        []string{"hide"},
			RemoveFields:          []string{"remove"},
			MandatoryFields:       []string{"must"},
			SuccessMessage:        "Much success",
			ErrorMessageMandatory: "Not much success",
			ErrorMessageGeneric:   "Not much success",
			Notifications: FormNotification{
				Slack: FormNotificationSlack{
					Webhook: "http://slack.webhook/test",
					Enabled: true,
				},
				Email: FormNotificationEmail{
					To:      "test@test.com",
					Cc:      "test@test.com",
					From:    "from@test.com",
					Subject: "test",
					Enabled: true,
					Options: FormNotificationEmailOpts{
						TextOnly:       true,
						IncludeResults: true,
					},
				},
			},
		},
	}

	f, err := client.AddForm(want)

	teardown()

	assert.Empty(t, err)
	assert.NotEmpty(t, f)
	assert.Equal(t, want, f)
}

func TestGetForm(t *testing.T) {
	setup()
	mux.HandleFunc("/v1/revisions", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method, "Expected method 'GET', got %s", r.Method)
		assert.Equal(t, "/test-content", r.Header.Get("Quant-Url"))
		fmt.Fprintf(w, `{
	"published": true,
	"published_revision": 2,
	"url": "\/test-content",
	"revisions": {
		"1": {
			"md5": "00000000000000000000000000000000",
			"type": "content",
			"revision_number": 1,
			"byte_length": 14,
			"date_timestamp": 0
		},
		"2": {
			"form_config": {
				"honeypot_fields": [
					"hide"
				],
				"target_url": "\/test-content",
				"notifications": {
					"slack": {
						"webhook": "http:\/\/slack.webhook\/test",
						"enabled": true
					},
					"email": {
						"cc": "test@test.com",
						"to": "test@test.com",
						"from": "from@test.com",
						"enabled": true,
						"options": {
							"text_only": true,
							"include_results": true
						},
						"subject": "test"
					}
				},
				"error_message_generic": "Not much success",
				"remove_fields": [
					"remove"
				],
				"error_message_mandatory": "Not much success",
				"success_message": "Much success",
				"mandatory_fields": [
					"must"
				]
			},
			"md5": "00000000000000000000000000000000",
			"type": "content",
			"form_enabled": true,
			"revision_number": 2,
			"byte_length": 14,
			"date_timestamp": 0
		}
	},
	"transitions": [],
	"highest_revision_number": 2,
	"seq_num": 2
}`)
	})
	f, err := client.GetForm(RevisionQuery{
		Url: "/test-content",
	})

	assert.Empty(t, err)
	assert.NotEmpty(t, f)

	want := Form{
		Url:     "/test-content",
		Enabled: true,
		Config: FormConfig{
			Target:                "/test-content",
			HoneypotFields:        []string{"hide"},
			RemoveFields:          []string{"remove"},
			MandatoryFields:       []string{"must"},
			SuccessMessage:        "Much success",
			ErrorMessageMandatory: "Not much success",
			ErrorMessageGeneric:   "Not much success",
			Notifications: FormNotification{
				Slack: FormNotificationSlack{
					Webhook: "http://slack.webhook/test",
					Enabled: true,
				},
				Email: FormNotificationEmail{
					To:      "test@test.com",
					Cc:      "test@test.com",
					From:    "from@test.com",
					Subject: "test",
					Enabled: true,
					Options: FormNotificationEmailOpts{
						TextOnly:       true,
						IncludeResults: true,
					},
				},
			},
		},
	}

	assert.Equal(t, want, f)
}
func TestDeleteForm(t *testing.T) {
	setup()
	mux.HandleFunc("/v1/form", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method, "Expected method 'DELETE', got %s", r.Method)
	})
}
