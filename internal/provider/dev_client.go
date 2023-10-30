package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	devops_resource "github.com/liatrio/devops-bootcamp/examples/ch6/devops-resources"
)

// Function to create a dev
func (c *Client) CreateDev(dev *devops_resource.Dev) (*devops_resource.Dev, error) {
	reqBody, err := json.Marshal(dev)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/dev", c.HostURL), bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	newDev := devops_resource.Dev{}

	err = json.Unmarshal(res, &newDev)

	if err != nil {
		return nil, err
	}

	return &newDev, nil
}

func (c *Client) GetDevByName(name string) (*devops_resource.Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev/name/%s", c.HostURL, name), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	dev := devops_resource.Dev{}

	err = json.Unmarshal(body, &dev)

	if err != nil {
		return nil, err
	}

	return &dev, nil
}

func (c *Client) GetDevById(id string) (*devops_resource.Dev, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/dev/id/%s", c.HostURL, id), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	dev := devops_resource.Dev{}

	err = json.Unmarshal(body, &dev)

	if err != nil {
		return nil, err
	}

	return &dev, nil
}

func (c *Client) DeleteDev(dev *devops_resource.Dev) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/dev/%s", c.HostURL, dev.Id), nil)

	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)

	return err
}

func (c *Client) UpdateDev(dev *devops_resource.Dev) (*devops_resource.Dev, error) {
	reqBody, err := json.Marshal(dev)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/dev/%s", c.HostURL, dev.Id), bytes.NewBuffer(reqBody))

	if err != nil {
		return nil, err
	}

	res, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	newDev := devops_resource.Dev{}

	err = json.Unmarshal(res, &newDev)

	if err != nil {
		return nil, err
	}

	return &newDev, nil
}
