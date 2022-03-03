package quant

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const FormsUrl = apiHost + apiBase + "/forms"

type Form struct {
	Url     string     `json:"url"`
	Enabled string     `json:"form_enabled"`
	Config  FormConfig `json:"form_config"`
}

type FormConfig struct {
	Target          string           `json:"target_url"`
	HoneypotFields  []string         `json:"honeypot_fields"`
	MandatoryFields []string         `json:"mandatory_fields"`
	RemoveFields    []string         `json:"remove_fields"`
	Notifications   FormNotification `json:"notifications,omitempty"`
}

type FormNotification struct {
	Slack FormNotificationSlack `json:"slack,omitempty"`
	Email FormNotificationEmail `json:"email,omitempty"`
}

type FormNotificationEmail struct {
	To      string                    `json:"to"`
	Cc      string                    `json:"cc"`
	From    string                    `json:"from"`
	Subject string                    `json:"subject"`
	Enabled bool                      `json:"enabled"`
	Options FormNotificationEmailOpts `json:"options,omitempty"`
}

type FormNotificationEmailOpts struct {
	TextOnly       bool `json:"text_only,omitempty"`
	IncludeResults bool `json:"include_results,omitempty"`
}

type FormNotificationSlack struct {
	Webhook string `json:"webhook,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

type FormQuery struct {
	FormRoute string
}

// ListForms gathers all the forms for a given project.
func (c *Client) ListForms() ([]Form, error) {
	panic("Not implemented.")
}

// Return a form revision from the API.
// v1 of the API treats forms as content revisions, this information
// is bundled with a routes configuration. We need to collect the revision
// information and then see if the latest revision has form configuration attached.
func (c *Client) GetForm(query FormQuery) (FormConfig, error) {
	revision, err := c.GetRevisionLatest(RevisionQuery{
		Url: query.FormRoute,
	})

	if err != nil {
		return nil, err
	}

	return revision.FormConfig
}

// Add form configuration to a specific route.
func (c *Client) AddForm(form Form) ([]byte, error) {
	j, err := json.Marshal(form)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", FormsUrl, bytes.NewBuffer(j))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Quant-Url", form.Url)

	res, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) DeleteForm(formRoute string) ([]byte, error) {
	request, err := c.NewRequest("forms", "DELETE")

	if err != nil {
		return nil, err
	}

	request.Header.Set("Quant-Url", formRoute)
	res, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	return res, nil
}
