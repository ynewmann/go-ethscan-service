package etherscan

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewApi(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		assert.NotPanics(t, func() {
			api := NewApi()
			assert.NotNil(t, api)
		})
	})

	t.Run("WithApiClient", func(t *testing.T) {
		assert.NotPanics(t, func() {
			client := NewApiClient()
			newApi := NewApi(WithClient(client))
			require.NotNil(t, newApi)

			a, ok := newApi.(*api)
			require.True(t, ok)
			assert.Equal(t, client, a.client)
		})
	})
}
