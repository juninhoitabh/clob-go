package httpClient

import (
	"context"
	"io"
	"net/http"
)

type HttpResponse struct {
	RawBody    io.ReadCloser
	Headers    http.Header
	Body       []byte
	StatusCode int
}

type HttpClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error)
}
