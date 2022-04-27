package quant

import (
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
)

type Revision struct {
	Published         bool                    `json:"published"`
	PublishedRevision int                     `json:"published_revision"`
	Sequence          int                     `json:"seq_num"`
	Revisions         map[string]RevisionItem `json:"revisions"`
	Url               string                  `json:"url"`
}

type RevisionItem struct {
	MD5            string     `json:"md5"`
	Type           string     `json:"type"`
	RevisionNumber int        `json:"revision_number"`
	ByteLength     int        `json:"byte_length"`
	DateTimestamp  int        `json:"date_timestamp"`
	FormConfig     FormConfig `json:"form_config,omitempty"`
	FormEnabled    bool       `json:"form_enabled,omitempty"`
}

type RevisionInfo struct {
	Log           string `json:"log"`
	Author        string `json:"admin"`
	DateTimestamp int    `json:"date_timestamp"`
}

type RevisionQuery struct {
	Url string
}

type MarkupRevision struct {
	Url             string `json:"url"`
	FindAttachments bool   `json:"find_attachments"`
	Content         []byte `json:"content"`
	Published       bool   `json:"published"`
}

func (c *Client) AddMarkupRevision(revision MarkupRevision, skip bool) (r Revision, err error) {
	j, err := json.Marshal(revision)
	req, err := c.NewRequest("", "POST", bytes.NewBuffer(j))
	if skip {
		req.Header.Set("Quant-Skip-Purge", "true")
	}

	res, err := c.doRequest(req)

	var apiError ApiError
	json.Unmarshal(res, &apiError)
	if apiError.Message != "" {
		err = errors.New(apiError.Message)
		return
	}

	err = json.Unmarshal(res, &r)
	return
}

// Get route revision information.
func (c *Client) GetRevision(query RevisionQuery) (r Revision, err error) {
	request, err := c.NewRequest("revisions", "GET", nil)
	request.Header.Set("Quant-Url", query.Url)
	res, err := c.doRequest(request)

	var apiError ApiError
	json.Unmarshal(res, &apiError)
	if apiError.Message != "" {
		err = errors.New(apiError.Message)
		return
	}

	json.Unmarshal(res, &r)
	return
}

// Get the latest revision for a URL.
func (c *Client) GetRevisionLatest(query RevisionQuery) (r RevisionItem, err error) {
	revision, err := c.GetRevision(query)
	key := strconv.FormatInt(int64(revision.PublishedRevision), 10)
	r = revision.Revisions[key]
	return
}
