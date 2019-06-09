package wetsuit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Client struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	Origin       string
	HTTPClient   *http.Client
}

func NewClient(or string, cid string, cs string, at string) *Client {
	return &Client{
		ClientID:     cid,
		ClientSecret: cs,
		Origin:       or,
		AccessToken:  at,
		HTTPClient:   &http.Client{},
	}
}

func (c *Client) getAPIURL(path string) string {
	return strings.Join([]string{c.Origin, path}, "/api")
}

func (c *Client) newRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := c.getAPIURL(path)
	fmt.Println(url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", strings.Join([]string{"Bearer", c.AccessToken}, " "))
	return req, nil
}

func (c *Client) read(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (c *Client) handleRequestError(resp *http.Response) error {
	if resp.StatusCode > 400 {
		return errors.New(fmt.Sprintf("HTTP Error: %s", resp.Status))
	}
	return nil
}

func (c *Client) Patch(path string, bodyToJSON interface{}) ([]byte, error) {
	jb, err := json.Marshal(bodyToJSON)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(jb)

	req, err := c.newRequest(http.MethodPatch, path, reader)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := c.handleRequestError(resp); err != nil {
		return nil, err
	}

	return c.read(resp)
}

func (c *Client) Post(path string, bodyToJSON interface{}) ([]byte, error) {
	jb, err := json.Marshal(bodyToJSON)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(jb)

	req, err := c.newRequest(http.MethodPost, path, reader)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := c.handleRequestError(resp); err != nil {
		return nil, err
	}

	return c.read(resp)
}

func (c *Client) Get(path string) ([]byte, error) {
	req, err := c.newRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err := c.handleRequestError(resp); err != nil {
		return nil, err
	}

	return c.read(resp)
}
