package cloud

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/blackburn29/cloud/model"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// A client represents a client connection to a {own|next}cloud
type Client struct {
	Url      *url.URL
	Username string
	Password string
}

type Header struct {
	Name	string
	Value	string
}

// Error type encapsulates the returned error messages from the
// server.
type Error struct {
	// Exception contains the type of the exception returned by
	// the server.
	Exception string `xml:"exception"`

	// Message contains the error message string from the server.
	Message string `xml:"message"`
}

// Dial connects to an {own|next}Cloud instance at the specified
// address using the given credentials.
func Dial(host, username, password string) (*Client, error) {
	uri, err := url.Parse(host)

	if err != nil {
		return nil, err
	}
	return &Client{
		Url:      uri,
		Username: username,
		Password:  password,
	}, nil
}

func (c *Client) ListDirectory(path string, depth int) (*model.MultiStatusResponse, error) {
	properties := model.CreateProperties()

	destUrl, err := url.Parse(fmt.Sprintf("files/%s/%s", c.Username, path))
	if err != nil {
		return nil, err
	}

	headers := []Header{
		{Name: "Depth", Value: fmt.Sprintf("%d", depth)},
	}
	body, _ := xml.Marshal(properties)
	dir := model.MultiStatusResponse{}

	resp, _, err := c.sendRequest("PROPFIND", c.Url.ResolveReference(destUrl).String(), body, &dir, headers)

	if resp != nil && resp.StatusCode == 404 {
		return nil, fmt.Errorf("file not found")
	}

	if err != nil {
		return nil, err
	}

	return &dir, nil
}

// Mkdir creates a new directory on the cloud with the specified name.
func (c *Client) Mkdir(path string) error {
	resp, _, err := c.sendRequest("MKCOL", fmt.Sprintf("files/%s/%s", c.Username, path), nil, nil, nil)

	//405 Status code is returned when the directory already exists
	if resp.StatusCode == 405 {
		return nil
	}

	return err
}

// Delete removes the specified folder from the cloud.
func (c *Client) Delete(path string) error {
	_, _, err := c.sendRequest("DELETE", fmt.Sprintf("files/%s/%s", c.Username, path), nil, nil, nil)
	return err
}

// Upload uploads the specified source to the specified destination
// path on the cloud.
func (c *Client) Upload(src []byte, dest string) error {

	destUrl, err := url.Parse(fmt.Sprintf("files/%s/%s", c.Username, dest))
	if err != nil {
		return err
	}

	_, _, err = c.sendRequest("PUT", c.Url.ResolveReference(destUrl).String(), src, nil, nil)

	return err
}

// Download downloads a file from the specified path.
func (c *Client) Download(path string) ([]byte, error) {

	pathUrl, err := url.Parse(fmt.Sprintf("files/%s/%s", c.Username, path))
	if err != nil {
		return nil, err
	}

	// Create the https request

	_, content, err := c.sendRequest("GET", c.Url.ResolveReference(pathUrl).String(), nil, nil, nil)

	if err != nil {
		return nil, err
	}

	return content, nil
}

func (c *Client) Exists(path string) bool {
	_, _, err := c.sendRequest("PROPFIND", fmt.Sprintf("files/%s/%s", c.Username, path), nil, nil, nil)
	return err == nil
}

func (c *Client) AddTag(path string, tag *model.Tag) (bool, error) {
	_, err := c.AddSystemTag(tag)

	if err != nil {
		return false, err
	}

	file, err := c.ListDirectory(path, 0)

	if file == nil || file.Responses == nil {
		return false, err
	}

	body, _ := json.Marshal(tag)
	headers := []Header{
		{
			Name: "Content-Type",
			Value: "application/json",
		},
	}

	tags, _ := c.GetSystemTags()

	systemTag := model.SystemTagProperty{}

	for _, response := range tags.Responses {
		for _, prop := range response.Properties {
			if prop.DisplayName == tag.Name {
				systemTag = prop
			}
		}
	}

	if systemTag.Id == "" {
		return false, fmt.Errorf("could not find system tag with name %s", tag.Name)
	}

	fileId := file.Responses[0].Properties[0].FileId
	resp, _, err := c.sendRequest("PUT", fmt.Sprintf("systemtags-relations/files/%s/%s", fileId, systemTag.Id), body, nil, headers)

	if resp.StatusCode == 409 {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) FindSystemTag(displayName string) (*model.TagPropResponse, error) {
	tags, err := c.GetSystemTags()

	if err != nil {
		return nil, err
	}

	for _, tag := range tags.Responses {
		for _, properties := range tag.Properties {
			if properties.DisplayName == displayName {
				return &tag, nil
			}
		}
	}

	return nil, nil
}

func (c *Client) GetSystemTags() (*model.MultiStatusTagResponse, error) {
	pathUrl, err := url.Parse("systemtags/")
	if err != nil {
		return nil, err
	}

	// Create the https request

	props, _ := xml.Marshal(model.CreateTagProperties())

	resp := model.MultiStatusTagResponse{}

	_, _, err = c.sendRequest("PROPFIND", c.Url.ResolveReference(pathUrl).String(), props, &resp, nil)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) AddSystemTag(tag *model.Tag) (bool, error) {

	pathUrl, err := url.Parse("systemtags/")
	if err != nil {
		return false, err
	}

	// Create the https request

	body, _ := json.Marshal(tag)
	headers := []Header{
		{
			Name: "Content-Type",
			Value: "application/json",
		},
	}

	resp, _, err := c.sendRequest("POST", c.Url.ResolveReference(pathUrl).String(), body, nil, headers)

	if resp.StatusCode == 409 {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (c *Client) sendRequest(request string, path string, body []byte, returnValue interface{}, headers []Header) (*http.Response, []byte, error) {
	// Create the https request

	folderUrl, err := url.Parse(path)
	if err != nil {
		return nil, nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest(request, c.Url.ResolveReference(folderUrl).String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)

	if headers != nil {
		for _, header := range headers {
			req.Header.Add(header.Name, header.Value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	resp.Body.Close()
	req.Body = ioutil.NopCloser(bytes.NewBuffer(content))

	if len(content) > 0 {
		if returnValue != nil {
			err = xml.Unmarshal(content, &returnValue)

			if err == nil {
				return resp, content, nil
			}
		}

		errorXml := Error{}
		err = xml.Unmarshal(content, &errorXml)
		if err != nil {
			if err == io.EOF {
				return resp, content, nil
			}

			return resp, content, fmt.Errorf("error during XML Unmarshal for response %s. the error was %s", content, err)
		}
		if errorXml.Exception != "" {
			return resp, content, fmt.Errorf("exception: %s, message: %s", errorXml.Exception, errorXml.Message)
		}
	}

	return resp, content, nil
}
