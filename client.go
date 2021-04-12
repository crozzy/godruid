package godruid

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strconv"
)

// NewClient gives you a new client to play with, all for the migre price of
// two instanciating args, url (of the brokers), and path (/v2/query).
func NewClient(brokerURLs []string, path string) *Client {
	return &Client{
		brokerURLs: brokerURLs,
		path:       path,
		httpClient: &http.Client{},
	}
}

func NewClientWithQueryContext(brokerURLs []string, path string, context *QueryContext) *Client {
	return &Client{
		brokerURLs:   brokerURLs,
		path:         path,
		queryContext: context,
		httpClient:   &http.Client{},
	}
}

// Client is what you initiate when you want to compose and execute a query.
type Client struct {
	brokerURLs []string
	path       string
	// Allow callers to specify a default context that will be used unless
	// setting the Context explicity in the query.
	queryContext *QueryContext
	httpClient   *http.Client
	lastUsed     int
}

func (c *Client) getHost() string {
	if len(c.brokerURLs) < c.lastUsed+2 { //offset + 0 indexed
		c.lastUsed = 0
		return c.brokerURLs[0]
	}
	c.lastUsed++
	return c.brokerURLs[c.lastUsed]
}

type ErrorQueryingDruid struct {
	errorCode   int
	query, body []byte
}

func NewErrorQueryingDruid(status int, body, query []byte) *ErrorQueryingDruid {
	return &ErrorQueryingDruid{
		errorCode: status,
		query:     query,
		body:      body,
	}
}

func (eqd *ErrorQueryingDruid) Error() string {
	return "Got a non-200 response from Druid: " + strconv.Itoa(eqd.errorCode) +
		"\nResponse: " + string(eqd.body) +
		"\nQuery: " + string(eqd.query)
}

func (c *Client) Run(q query) ([]byte, error) {
	queryReq, err := q.ToJSON()
	if err != nil {
		return nil, err
	}
	// Do actual querying
	req, err := http.NewRequest("POST", c.getHost()+c.path, bytes.NewReader(queryReq))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, NewErrorQueryingDruid(resp.StatusCode, body, queryReq)
	}
	return body, nil
}
