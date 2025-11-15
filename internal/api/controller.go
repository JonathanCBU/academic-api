package api

import (
	"bytes"
	"net/http"
)

type Controller struct {
	client *http.Client
}

func NewController(client *http.Client) *Controller {
	return &Controller{
		client: client,
	}
}

func (c *Controller) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	testBody := bytes.NewBuffer([]byte(`{"test": "hello world"}`))
	outgoing, err := http.NewRequest("POST", httpBinPostUrl, testBody)
	if err != nil {
		// TODO: handle request create error
	}
	outgoing.Header.Set("Content-Type", ContentTypeHeader)

	resp, err := c.client.Do(outgoing)
	if err != nil {
		// TODO: handle outgoing error
	}
	defer resp.Body.Close()

}
