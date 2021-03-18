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

func (a *Client) GetSuccess(all_params map[string]interface{}) (*GetMockOK, error) {
	wait, ok := all_params["wait"]
	if ok {
		wait_val, err := conv.Int64(wait)
		if err != nil {
			wait_val = int64(0)
		}
		time.Sleep(time.Duration(wait_val) * time.Millisecond)
	}
	delete(all_params, "wait")
	delete(all_params, "_timeout")
	return &GetMockOK {
		Payload: all_params,
	}, nil
}

func (a *Client) GetFailure(all_params map[string]interface{}) (*GetMockOK, error) {
	return nil, errors.New("Mock Error")
}

func (a *Client) GetFailurePayload(all_params map[string]interface{}) (*GetMockOK, error) {
	return nil, &GetMockError {
		Payload: all_params,
	}
}
