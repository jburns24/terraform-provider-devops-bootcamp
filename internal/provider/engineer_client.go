package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	devops_resource "github.com/liatrio/devops-bootcamp/examples/ch6/devops-resources"
)

func (c *Client) GetEngineer(Id string) (*devops_resource.Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/id/%s", c.HostURL, Id), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	engineer := devops_resource.Engineer{}

	err = json.Unmarshal(body, &engineer)

	if err != nil {
		return nil, err
	}
	return &engineer, nil
}

func (c *Client) GetEngineerByName(name string) (*devops_resource.Engineer, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/engineers/name/%s", c.HostURL, name), nil)

	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	engineer := devops_resource.Engineer{}

	err = json.Unmarshal(body, &engineer)

	if err != nil {
		return nil, err
	}

	return &engineer, nil
}

func (c *Client) CreateEngineer(name string, email string) (*devops_resource.Engineer, error) {
	newEngineer := devops_resource.Engineer{
		Name:  name,
		Email: email,
	}
	jsonBody, err := json.Marshal(newEngineer)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/engineers", c.HostURL), bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	engineer := devops_resource.Engineer{}

	err = json.Unmarshal(body, &engineer)

	if err != nil {
		return nil, err
	}

	return &engineer, nil
}

func (c *Client) UpdateEngineer(id string, name string, email string) (*devops_resource.Engineer, error) {
	engineer := devops_resource.Engineer{
		Name:  name,
		Email: email,
		Id:    id,
	}

	jsonBody, err := json.Marshal(engineer)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/engineers/%s", c.HostURL, id), bytes.NewBuffer(jsonBody))

	if err != nil {
		return nil, err
	}

	body, err := c.DoRequest(req)

	if err != nil {
		return nil, err
	}

	updatedEngineer := devops_resource.Engineer{}

	err = json.Unmarshal(body, &updatedEngineer)

	if err != nil {
		return nil, err
	}

	return &updatedEngineer, nil
}

func (c *Client) DeleteEngineer(engineer *devops_resource.Engineer) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/engineers/%s", c.HostURL, engineer.Id), nil)

	if err != nil {
		return err
	}

	_, err = c.DoRequest(req)

	if err != nil {
		return err
	}

	return nil
}
