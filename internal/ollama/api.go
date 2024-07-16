package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type API struct {
	addr   string
	client *http.Client
}

func NewApi(addr string) *API {
	return &API{
		addr: addr,
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

func (a *API) Prompt(ctx context.Context, prompt PromptReq) (reply PromptResp, err error) {
	buf := bytes.NewBuffer(nil)
	if err := json.NewEncoder(buf).Encode(prompt); err != nil {
		return reply, err
	}
	resp, err := a.client.Post(
		a.addr+"/api/generate",
		"application/json; charset=UTF-8",
		buf,
	)
	if err != nil {
		return reply, err
	}
	if err := json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return reply, err
	}
	resp.Body.Close()
	return reply, nil
}
