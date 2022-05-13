package etherscan

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApiClient(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		assert.NotPanics(t, func() {
			apiClient := NewApiClient()
			assert.NotNil(t, apiClient)
		})
	})

	t.Run("WithUrl", func(t *testing.T) {
		assert.NotPanics(t, func() {
			url := "test"
			apiClient := NewApiClient(WithUrl(url))
			require.NotNil(t, apiClient)
			assert.EqualValues(t, apiClient.url, url)
		})
	})

	t.Run("WithUrl", func(t *testing.T) {
		assert.NotPanics(t, func() {
			apiKey := "ApiKey"
			apiClient := NewApiClient(WithApiKey(apiKey))
			require.NotNil(t, apiClient)
			assert.EqualValues(t, apiClient.apiKey, apiKey)
		})
	})
}

func TestApiClient_DoNewRequest(t *testing.T) {
	apiClient := NewApiClient()
	require.NotNil(t, apiClient)

	t.Run("NilRequest", func(t *testing.T) {
		resp, err := apiClient.do(nil, nil)
		assert.Nil(t, resp)
		assert.EqualError(t, err, ErrBadRequest.Error())
	})

	t.Run("EmptyRequest", func(t *testing.T) {
		url, err := url.Parse(DefaultApiUrl)
		require.NoError(t, err, err)

		resp, err := apiClient.do(&http.Request{URL: url}, nil)
		assert.NoError(t, err, err)
		assert.NotNil(t, resp)
	})
}
