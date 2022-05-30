package etherscan

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultApiUrl = "https://api.etherscan.io/api"
	ApikeyArg     = "apikey"

	ModuleArg = "module"
)

type (
	Api interface {
		Proxy() Proxy
	}

	api struct {
		client *http.Client
		url    string
		apiKey string

		proxy Proxy
	}

	ApiOption func(*api)
	ArgOption func(*url.Values)
)

var ErrBadRequest = errors.New("BadRequest")

func WithUrl(url string) ApiOption {
	return func(client *api) {
		client.url = url
	}
}

func WithApiKey(apiKey string) ApiOption {
	return func(client *api) {
		client.apiKey = apiKey
	}
}

func NewApi(options ...ApiOption) Api {
	api := &api{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		url: DefaultApiUrl,
	}
	for _, opt := range options {
		opt(api)
	}

	api.proxy = NewProxyModule(api)

	return api
}

func (a *api) Proxy() Proxy {
	return a.proxy
}

func withArg(arg, val string) ArgOption {
	return func(values *url.Values) {
		values.Add(arg, val)
	}
}

func (a *api) newRequest(ctx context.Context, method string, body interface{}, args ...ArgOption) (*http.Request, error) {
	u, err := url.Parse(a.url)
	if err != nil {
		return nil, err
	}

	values := u.Query()
	if a.apiKey != "" {
		values.Add(ApikeyArg, a.apiKey)
	}
	for _, arg := range args {
		arg(&values)
	}
	u.RawQuery = values.Encode()

	buf := new(bytes.Buffer)
	if body != nil {
		enc := json.NewEncoder(buf)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (a *api) do(req *http.Request, v interface{}) (*http.Response, error) {
	if req == nil {
		return nil, ErrBadRequest
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if v != nil {
		err = json.Unmarshal(body, v)
		if err != nil && err != io.EOF {
			ep := &ErrorPage{}
			err = json.Unmarshal(body, ep)
			if err != nil && err != io.EOF {
				return nil, err
			}

			return nil, errors.New(ep.Result)
		}
	}

	return resp, nil
}

func (a *api) doNewRequest(ctx context.Context, method string, body, v interface{}, args ...ArgOption) (*http.Response, error) {
	req, err := a.newRequest(ctx, method, body, args...)
	if err != nil {
		return nil, err
	}

	return a.do(req, v)
}
