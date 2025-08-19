package httpClient

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type DefaultHttpClient struct {
	client  *http.Client
	timeout time.Duration
}

func NewDefaultHttpClient(timeout time.Duration) *DefaultHttpClient {
	return &DefaultHttpClient{
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

func (c *DefaultHttpClient) Get(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error) {
	return c.do(ctx, http.MethodGet, url, nil, headers)
}

func (c *DefaultHttpClient) Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error) {
	return c.do(ctx, http.MethodPost, url, body, headers)
}

func (c *DefaultHttpClient) Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error) {
	return c.do(ctx, http.MethodPut, url, body, headers)
}

func (c *DefaultHttpClient) Delete(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error) {
	return c.do(ctx, http.MethodDelete, url, nil, headers)
}

func (c *DefaultHttpClient) do(ctx context.Context, method, url string, body interface{}, headers map[string]string) (*HttpResponse, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}

		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()

		return nil, err
	}

	bodyReader := io.NopCloser(bytes.NewBuffer(respBody))

	httpResponse := &HttpResponse{
		StatusCode: resp.StatusCode,
		Body:       respBody,
		Headers:    resp.Header,
		RawBody:    bodyReader,
	}

	return httpResponse, nil
}
