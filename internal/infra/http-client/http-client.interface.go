package httpClient

import (
	"context"
	"io"
	"net/http"
)

type HttpResponse struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
	RawBody    io.ReadCloser
}

type HttpClient interface {
	Get(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error)
	Post(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error)
	Put(ctx context.Context, url string, body interface{}, headers map[string]string) (*HttpResponse, error)
	Delete(ctx context.Context, url string, headers map[string]string) (*HttpResponse, error)
}
