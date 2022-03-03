package quant

import (
	"encoding/json"
	"errors"
	"strconv"
)

type Revision struct {
	Published         bool                    `json:"published"`
	PublishedRevision int                     `json:"published_revision"`
	Sequence          int                     `json:"seq_num"`
	Revisions         map[string]RevisionItem `json:"revisions"`
}

type RevisionItem struct {
	MD5            string     `json:"md5"`
	Type           string     `json:"type"`
	RevisionNumber int        `json:"revision_number"`
	ByteLength     int        `json:"byte_length"`
	DateTimestamp  int        `json:"date_timestamp"`
	FormConfig     FormConfig `json:"form_config,omitempty"`
}

type RevisionInfo struct {
	Log           string `json:"log"`
	Author        string `json:"admin"`
	DateTimestamp int    `json:"date_timestamp"`
}

type RevisionQuery struct {
	Url string
}

// Get route revision information.
func (c *Client) GetRevision(query RevisionQuery) (Revision, error) {
	request, err := c.NewRequest("revisions", "GET")

	if err != nil {
		return nil, err
	}
	request.Header.set("Quant-Url", query.Url)
	response, err := c.doRequest(request)

	if err != nil {
		return nil, err
	}

	var revision Revision
	json.Unmarshal(response.Body, &revision)

	return revision, nil
}

// Get the latest revision for a URL.
func (c *Client) GetRevisionLatest(query RevisionQuery) (RevisionItem, error) {
	revision, err := c.GetRevision(query)
	if err != nil {
		return nil, err
	}
	key := strconv.FormatInt(int64(revision.PublishedRevision), 10)
	revision, ok := revision.Revisions[key]
	if !ok {
		return nil, errors.New("No revisions found for %s", key)
	}
	return revision
}