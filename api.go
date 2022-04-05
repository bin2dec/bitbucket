package bitbucket

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	defaultScheme  = "https"
	defaultPrefix  = "rest/api"
	defaultVersion = "1.0"
	defaultTimeout = 60 // seconds
)

type Api struct {
	Scheme   string
	Host     string
	Prefix   string
	Version  string
	Category string
	Auth     string
	Timeout  time.Duration
}

func NewApi(host, category string) *Api {
	return &Api{
		Scheme:   defaultScheme,
		Host:     host,
		Prefix:   defaultPrefix,
		Version:  defaultVersion,
		Category: category,
		Timeout:  defaultTimeout,
	}
}

func (a *Api) Url(pathElems []string, queryParams map[string]string) string {
	for i, v := range pathElems {
		pathElems[i] = url.PathEscape(v)
	}
	pathElems = append([]string{a.Prefix, a.Version, a.Category}, pathElems...)
	u := fmt.Sprintf("%s://%s/%s", a.Scheme, a.Host, path.Join(pathElems...))
	if len(queryParams) == 0 {
		return u
	}

	q := make([]string, 0, len(queryParams))
	for k, v := range queryParams {
		q = append(q, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(v)))
	}
	return fmt.Sprintf("%s?%s", u, strings.Join(q, "&"))
}

const defaultContentType = "application/json"

type Response struct {
	StatusCode int
	Body       []byte
}

func (a *Api) Get(url string) (resp *Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return a.do(req)
}

func (a *Api) Post(url string, body io.Reader) (resp *Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", defaultContentType)
	return a.do(req)
}

func (a *Api) do(req *http.Request) (resp *Response, err error) {
	req.Header.Set("Authorization", a.Auth)

	client := &http.Client{Timeout: a.Timeout * time.Second}

	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	resp = &Response{StatusCode: r.StatusCode}

	defer func() {
		if e := r.Body.Close(); e != nil && err == nil {
			err = e
		}
	}()

	resp.Body, err = io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
