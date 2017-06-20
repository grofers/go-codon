package operations

import (
	"errors"
)

func New() *Client {
	return &Client{}
}

type Client struct {
}

type GetMockOK struct {
	Payload interface{}
}

func (a *Client) GetSuccess(all_params *map[string]interface{}) (*GetMockOK, error) {
	p := (*all_params)
	return &GetMockOK {
		Payload: p
	}, nil
}

func (a *Client) GetFailure(all_params *map[string]interface{}) (*GetMockOK, error) {
	return nil, errors.Error("Mock Error")
}

func (a *Client) GetFailurePayload(all_params *map[string]interface{}) (*GetMockOK, error) {
	p := (*all_params)
	return nil, GetMockOK {
		Payload: p
	}
}
