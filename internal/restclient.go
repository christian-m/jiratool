package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type RestClient interface {
	GetProject(prjKey string) (*Project, error)
	CreateVersion(version Version) error
	UpdateVersion(version Version) error
}

type JiraRestClient struct {
	BaseURL    *url.URL
	HttpClient *http.Client
}

func CreateRestClient(userinfo *url.Userinfo, u *url.URL) (*JiraRestClient, error) {
	if u == nil {
		return nil, fmt.Errorf("url not specified")
	}
	u.User = userinfo
	httpClient := http.DefaultClient
	return &JiraRestClient{HttpClient: httpClient, BaseURL: u}, nil
}

func (c *JiraRestClient) GetProject(prjKey string) (*Project, error) {
	rel := &url.URL{Path: fmt.Sprintf("/rest/api/3/project/%s", prjKey)}
	req, err := c.createGetRequest(rel)
	if err != nil {
		return nil, err
	}
	prj := &Project{}
	_, err = c.call(req, prj)
	return prj, err
}

func (c *JiraRestClient) CreateVersion(version Version) error {
	rel := &url.URL{Path: fmt.Sprintf("/rest/api/3/version")}
	req, err := c.createRestRequest(rel, "POST", version)
	if err != nil {
		return err
	}
	_, err = c.call(req, &version)
	t, ok := err.(RestError)
	if ok && t.Status() == http.StatusBadRequest {
		return fmt.Errorf("Version %s kann nicht angelegt werden", version.Name)
	}
	return err
}

func (c *JiraRestClient) UpdateVersion(version Version) error {
	rel := &url.URL{Path: fmt.Sprintf("/rest/api/3/version/%s", version.Id)}
	req, err := c.createRestRequest(rel, "PUT", version)
	if err != nil {
		return err
	}
	_, err = c.call(req, &version)
	t, ok := err.(RestError)
	if ok && t.Status() == http.StatusBadRequest {
		return fmt.Errorf("Version %s kann nicht aktualisiert werden", version.Name)
	}
	return err
}

func (c *JiraRestClient) createGetRequest(url *url.URL) (*http.Request, error) {
	u := c.BaseURL.ResolveReference(url)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *JiraRestClient) createRestRequest(url *url.URL, method string, body interface{}) (*http.Request, error) {
	u := c.BaseURL.ResolveReference(url)
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *JiraRestClient) call(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	err = json.NewDecoder(resp.Body).Decode(v)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		restErr := RestError{resp.Status, resp.StatusCode}
		return resp, restErr
	}
	return resp, err
}
