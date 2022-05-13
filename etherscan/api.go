package etherscan

type (
	Api interface {
		Proxy() Proxy
	}

	ApiOption func(*api)

	api struct {
		client *apiClient

		proxy Proxy
	}
)

func WithClient(c *apiClient) ApiOption {
	return func(client *api) {
		client.client = c
	}
}

func NewApi(options ...ApiOption) Api {
	api := &api{client: NewApiClient()}
	for _, opt := range options {
		opt(api)
	}

	api.proxy = NewProxyModule(api.client)

	return api
}

func (a *api) Proxy() Proxy {
	return a.proxy
}
