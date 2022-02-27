package quant

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const FormsUrl = apiHost + apiBase + "/forms"

type Form struct {
	FormRoute         string   `json:"form_route"`
	Enabled           bool     `json:"form_enabled,omitempty"`
	SuccessMessage    string   `json:"success_message,omitempty"`
	MandatoryMessage  string   `json:"mandatory_message,omitempty"`
	FailureMessage    string   `json:"generic_message,omitempty"`
	RequiredFields    []string `json:"form_mandatory,omitempty"`
	HoneypotFields    []string `json:"form_honeypot,omitempty"`
	RemoveFields      []string `json:"form_remove,omitempty"`
	NotificationEmail FormNotificationEmail
	NotificationSlack FormNotificationSlack
}

type NewForm struct {
	FormRoute             string // required
	Enabled               bool
	SuccessMessage        string
	FailureMessage        string
	RequiredFields        string
	HoneypotFields        string
	RemoveFields          string
	DisableHtmlEmails     bool
	IncludeSubmissionData bool
	ToAddress             string
	CcAddress             string
	FromAddress           string
	Subject               string
	SlackWebhook          string
}

type FormNotificationEmail struct {
	DisableHtml           bool   `json:"email_text_only,omitempty"`
	IncludeSubmissionData bool   `json:"email_include_results,omitempty"`
	ToAddress             string `json:"notification_email_to,omitempty"`
	CcAddress             string `json:"notification_email_cc,omitempty"`
	FromAddress           string `json:"notification_email_from,omitempty"`
	Subject               string `json:"notification_email_subject,omitempty"`
}

type FormNotificationSlack struct {
	Webhook string `json:"slack_webhook,omitempty"`
}

func (c *Client) ListForms() ([]Form, error) {
	request, err := c.newRequest("forms", "GET")

	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	var forms []Form
	json.Unmarshal(res, &forms)
	return forms, nil
}

func (c *Client) GetForm(formRoute string) (*Form, error) {
	request, err := c.newRequest("forms", "GET")

	if err != nil {
		return nil, err
	}

	request.Header.Set("Quant-Form", formRoute)
	res, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	var form *Form
	json.Unmarshal(res, &form)

	return form, nil
}

func (c *Client) AddForm(form NewForm) ([]byte, error) {
	j, err := json.Marshal(form)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", FormsUrl, bytes.NewBuffer(j))

	if err != nil {
		return nil, err
	}

	req.Header.Set("Quant-Form", form.FormRoute)

	res, err := c.doRequest(req)

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) DeleteForm(formRoute string) ([]byte, error) {
	request, err := c.newRequest("forms", "DELETE")

	if err != nil {
		return nil, err
	}

	request.Header.Set("Quant-Form", formRoute)
	res, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	return res, nil
}
