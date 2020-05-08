package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"

	"github.com/aristanetworks/glog"
)

const (
	contentTypeJSON   = "application/json"
	contentTypeNDJSON = "application/x-ndjson"
	docType           = "_doc"
	update            = "_update"
)

// Client is an elasticsearch Client.
type Client struct {
	*http.Client
	host string
}

func (c *Client) do(req *http.Request, ignoreStatusCode int) ([]byte, error) {
	glog.Infof("\n...\nRequest is %v", req)
	resp, err := c.Do(req)
	glog.Infof("Resp in do = %v", resp)
	glog.Infof("Error in do = %v", err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 && resp.StatusCode != ignoreStatusCode {
		return nil, fmt.Errorf("%s", resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

// CreateIndex creates a new index with the specified body.
func (c *Client) CreateIndex(ctx context.Context, index string, body map[string]interface{},
	ignoreStatusCode int) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, c.host+"/"+index, bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", contentTypeJSON)
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// CreateIndexNoTemplate creates a new index without specified body.
func (c *Client) CreateIndexNoTemplate(ctx context.Context, index string,
	ignoreStatusCode int) error {
	req, err := http.NewRequest(http.MethodPut, c.host+"/"+index, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// CreateIndexTemplate creates a new index template with the specified body.
func (c *Client) CreateIndexTemplate(ctx context.Context, name string, body map[string]interface{},
	ignoreStatusCode int) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, c.host+"/_template/"+name,
		bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", contentTypeJSON)
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// DeleteIndex deletes an index.
func (c *Client) DeleteIndex(ctx context.Context, index string, ignoreStatusCode int) error {
	req, err := http.NewRequest(http.MethodDelete, c.host+"/"+index, nil)
	if err != nil {
		return err
	}
	if _, err = c.do(req, ignoreStatusCode); err != nil {
		return err
	}
	// TODO: update cached indices
	return nil
}

// DeleteDoc deletes a document with id from the specific index
func (c *Client) DeleteDoc(ctx context.Context, index string, id string,
	ignoreStatusCode int) error {
	req, err := http.NewRequest(http.MethodDelete,
		c.host+"/"+index+"/"+docType+"/"+id+"?refresh", nil)
	if err != nil {
		return err
	}
	if _, err = c.do(req, ignoreStatusCode); err != nil {
		return err
	}
	return nil
}

// UpdateDoc updates a document with id from the specific index
func (c *Client) UpdateDoc(ctx context.Context, index string, id string,
	body map[string]interface{}) ([]byte, error) {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost,
		c.host+"/"+index+"/"+update+"/"+id, bytes.NewReader(bodyJSON))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", contentTypeJSON)
	return c.do(req, 400)
}

// BulkRequest is a request for either index or delete operation, that can be bulked with others.
type BulkRequest struct {
	Delete  bool
	Index   string
	ID      string
	Body    interface{}
	Version int64
}

// Bulk allows to perform bulk indexing and/or delete operations without
// embedding the index or the operation.
func (c *Client) Bulk(ctx context.Context, request []BulkRequest) error {
	body := &bytes.Buffer{}
	encoder := json.NewEncoder(body)
	for _, request := range request {
		action := "index"
		if request.Delete {
			action = "delete"
		}
		operation := map[string]interface{}{
			action: map[string]interface{}{
				"_index": request.Index,
				"_id":    request.ID,
				"_type":  docType,
			},
		}
		if request.Version > 0 {
			// include a _version for lastest state bulk index request
			operation["index"].(map[string]interface{})["version"] = request.Version
			// set the version_type to be external when _version is set explicitly
			operation["index"].(map[string]interface{})["version_type"] = "external"
		}
		err := encoder.Encode(operation)
		if err != nil {
			return err
		}
		if !request.Delete {
			err = encoder.Encode(request.Body)
			if err != nil {
				return err
			}
		}
	}
	req, err := http.NewRequest(http.MethodPost, c.host+"/_bulk", body)
	if err != nil {
		return err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", contentTypeNDJSON)
	_, err = c.do(req, 0)
	return err
}

// Search send the query dsl body to the elasticsearch /index/_doc/_search endpoint
func (c *Client) Search(ctx context.Context, indices []string,
	body map[string]interface{}, ignoreStatusCode int) ([]byte, error) {
	bodyJSON := new(bytes.Buffer)
	enc := json.NewEncoder(bodyJSON)
	// aggregation query in some cases has special character >
	// these characters get escaped by json encoder, hence setting it false.
	enc.SetEscapeHTML(false)
	err := enc.Encode(body)
	if err != nil {
		return nil, err
	}

	glog.V(5).Infof("Elasticsearch query: %s", bodyJSON.String())

	index := strings.Join(indices[:], ",")
	req, err := http.NewRequest(http.MethodGet,
		c.host+"/"+path.Join(index, docType, "/_search"),
		bytes.NewReader(bodyJSON.Bytes()))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", contentTypeJSON)
	return c.do(req, ignoreStatusCode)
}

// RefreshIndex refreshes the index
func (c *Client) RefreshIndex(ctx context.Context, index string,
	ignoreStatusCode int) error {
	req, err := http.NewRequest(http.MethodPost, c.host+"/"+index+"/_refresh", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// Flush all the indices to disk
func (c *Client) Flush(ctx context.Context, ignoreStatusCode int) error {
	req, err := http.NewRequest(http.MethodPost, c.host+"/_flush/synced", nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// Get is a generic get method for testing, can perform a wide range of get requests to check
// results with provided url
func (c *Client) Get(ctx context.Context, url string, ignoreStatusCode int) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+"/"+url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, ignoreStatusCode)
}

// DeleteByQuery performs a deletion on every document that match a query for the given index
func (c *Client) DeleteByQuery(ctx context.Context, index string, body map[string]interface{},
	ignoreStatusCode int) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.host+"/"+path.Join(index,
		"_delete_by_query?conflicts=proceed"),
		bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", contentTypeJSON)
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// CreateLifecyclePolicy creates a new index lifecycle management policy with the
// policy specification in the body
func (c *Client) CreateLifecyclePolicy(ctx context.Context, name string,
	body map[string]interface{}, ignoreStatusCode int) error {
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, c.host+"/_ilm/policy/"+name,
		bytes.NewReader(bodyJSON))
	if err != nil {
		return err
	}
	defer req.Body.Close()
	req.Header.Set("Content-Type", contentTypeJSON)
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// DeleteLifecyclePolicy deletes an index lifecycle management policy
func (c *Client) DeleteLifecyclePolicy(ctx context.Context, name string,
	ignoreStatusCode int) error {
	req, err := http.NewRequest(http.MethodDelete, c.host+"/_ilm/policy/"+name, nil)
	if err != nil {
		return err
	}
	_, err = c.do(req, ignoreStatusCode)
	return err
}

// NewClient instantiates a new client to the given host.
func NewClient(host string) *Client {
	c := &Client{
		Client: &http.Client{
			Transport: &http.Transport{},
		},
		host: host,
	}
	return c
}
