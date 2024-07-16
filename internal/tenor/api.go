package tenor

import (
	"net/http"
	"time"
)

type API struct {
	client *http.Client
}

func NewApi() *API {
	return &API{
		client: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true,
			},
		},
	}
}

func (a *API) 
