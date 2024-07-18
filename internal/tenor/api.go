package tenor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type API struct {
	client *http.Client
	cache  *gifCache
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
		cache: newGifCache(context.Background(), time.Minute*1),
	}
}

func (a *API) FetchGifById(ctx context.Context, id GifId) ([]byte, error) {
	if cachedGif, ok := a.cache.get(id); ok {
		return cachedGif, nil
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		a.makeGifRequestUrl(id),
		nil,
	)
	if err != nil {
		return nil, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	result := bytes.NewBuffer(nil)
	if _, err := io.Copy(result, resp.Body); err != nil {
		return nil, err
	}
	resp.Body.Close()
	a.cache.set(id, result.Bytes(), time.Hour*24)
	return result.Bytes(), nil
}

func (a *API) makeGifRequestUrl(id GifId) string {
	return fmt.Sprintf("https://c.tenor.com/%s/tenor.gif", id)
}

type GifId string

const (
	NonCommentGifId GifId = "P1XVqnCktqMAAAAd"
)

const YaNeBuduCommentirovatUrl = "https://media.tenor.com/P1XVqnCktqMAAAAd/112.gif"
