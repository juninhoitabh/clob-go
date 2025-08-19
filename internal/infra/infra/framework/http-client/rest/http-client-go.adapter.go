package rest

import (
	"io"
	"net/http"

	"github.com/juninhoitabh/clob-go/internal/shared"
)

type HttpClientGoAdapter struct {
	Client http.Client
}

func (r *HttpClientGoAdapter) Get(url string, data io.Reader, headers http.Header) (*http.Response, error) {
	return r.handleRequest(http.MethodGet, url, data, headers)
}

func (r *HttpClientGoAdapter) Post(url string, data io.Reader, headers http.Header) (*http.Response, error) {
	return r.handleRequest(http.MethodPost, url, data, headers)
}

func (r *HttpClientGoAdapter) Patch(url string, data io.Reader, headers http.Header) (*http.Response, error) {
	return r.handleRequest(http.MethodPatch, url, data, headers)
}

func (r *HttpClientGoAdapter) Delete(url string, data io.Reader, headers http.Header) (*http.Response, error) {
	return r.handleRequest(http.MethodDelete, url, data, headers)
}

func (r *HttpClientGoAdapter) handleRequest(methodType, url string, data io.Reader, headers http.Header) (*http.Response, error) {
	request, err := http.NewRequest(methodType, url, data)
	if err != nil {
		return nil, err
	}

	request.Header = headers
	response, err := r.Client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return response, shared.ErrExternalApi
	}

	return response, nil
}

func NewHttpClientGoAdapter() *HttpClientGoAdapter {
	return &HttpClientGoAdapter{
		Client: http.Client{},
	}
}
