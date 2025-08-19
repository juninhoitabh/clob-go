package rest

import (
	"io"
	"net/http"
)

type IHttpClient interface {
	Get(url string, data io.Reader, headers http.Header) (*http.Response, error)
	Post(url string, data io.Reader, headers http.Header) (*http.Response, error)
	Patch(url string, data io.Reader, headers http.Header) (*http.Response, error)
	Delete(url string, data io.Reader, headers http.Header) (*http.Response, error)
}
