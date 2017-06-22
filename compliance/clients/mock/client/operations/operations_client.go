package operations

import (
	"errors"
	"fmt"
	"time"
	conv "github.com/cstockton/go-conv"
)

func New() *Client {
	return &Client{}
}

type Client struct {
}

type GetMockOK struct {
	Payload interface{}
}

type GetMockError struct {
	Payload interface{}
}

func (o *GetMockError) Error() string {
	return fmt.Sprintf("[GET /{unknown}][%d] %+v", 400, o.Payload)
}

func (a *Client) GetSuccess(all_params *map[string]interface{}) (*GetMockOK, error) {
	p := (*all_params)
	wait, ok := p["wait"]
	if ok {
		wait_val, err := conv.Int64(wait)
		if err != nil {
			wait_val = int64(0)
		}
		time.Sleep(time.Duration(wait_val) * time.Millisecond)
	}
	delete(p, "wait")
	delete(p, "_timeout")
	return &GetMockOK {
		Payload: p,
	}, nil
}

func (a *Client) GetFailure(all_params *map[string]interface{}) (*GetMockOK, error) {
	return nil, errors.New("Mock Error")
}

func (a *Client) GetFailurePayload(all_params *map[string]interface{}) (*GetMockOK, error) {
	p := (*all_params)
	return nil, &GetMockError {
		Payload: p,
	}
}
