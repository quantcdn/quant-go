package quant

import (
	"bytes"
	"encoding/json"
	"errors"
)

const FormsUrl = apiHost + "/" + apiBase + "/form"

type Form struct {
	Url     string     `json:"url"`
	Enabled bool       `json:"form_enabled"`
	Config  FormConfig `json:"form_config"`
}

type FormConfig struct {
	Target                string           `json:"target_url"`
	HoneypotFields        []string         `json:"honeypot_fields"`
	MandatoryFields       []string         `json:"mandatory_fields"`
	RemoveFields          []string         `json:"remove_fields"`
	SuccessMessage        string           `json:"success_message"`
	ErrorMessageMandatory string           `json:"error_message_mandatory"`
	ErrorMessageGeneric   string           `json:"error_message_generic"`
	Notifications         FormNotification `json:"notifications,omitempty"`
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

// ListForms gathers all the forms for a given project.
func (c *Client) ListForms() ([]Form, error) {
	panic("Not implemented.")
}

// Return a form revision from the API.
// v1 of the API treats forms as content revisions, this information
// is bundled with a routes configuration. We need to collect the revision
// information and then see if the latest revision has form configuration attached.
func (c *Client) GetForm(query RevisionQuery) (f Form, err error) {
	r, err := c.GetRevisionLatest(query)
	f = Form{
		Url:     r.FormConfig.Target,
		Enabled: r.FormEnabled,
		Config:  r.FormConfig,
	}
	return
}

// Add form configuration to a specific route.
func (c *Client) AddForm(form Form) (f Form, err error) {
	j, err := json.Marshal(form)
	req, err := c.NewRequest("/form", "POST", bytes.NewBuffer(j))
	req.Header.Set("Quant-Url", form.Url)
	res, err := c.doRequest(req)

	var apiError ApiError
	json.Unmarshal(res, &apiError)

	if apiError.Message != "" {
		err = errors.New(apiError.Message)
		return
	}

	json.Unmarshal(res, &f)
	return
}

// Update form configuration.
func (c *Client) UpdateForm(form Form) (f Form, err error) {
	f, err = c.AddForm(form)
	return
}

// Delete the form configuration.
func (c *Client) DeleteForm(query RevisionQuery) ([]byte, error) {
	request, err := c.NewRequest("form", "DELETE", nil)

	if err != nil {
		return nil, err
	}

	request.Header.Set("Quant-Url", query.Url)
	res, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	return res, nil
}
