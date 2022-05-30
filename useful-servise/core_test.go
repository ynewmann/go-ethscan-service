package useful_servise

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUsefulService(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		us := NewUsefulService()
		require.NotNil(t, us)
		assert.NotNil(t, us.cfg)
		assert.NotNil(t, us.server)
		assert.NotNil(t, us.api)
		assert.Nil(t, us.totalAmountCache) // because cfg.UseCache is false
	})

	t.Run("WithConfig", func(t *testing.T) {
		cfg := &Config{
			UseCache:              true,
			MemoryCacheBackupPath: "test",
			CacheSize:             1,
			Address:               ":3000",
			ApiUrl:                "apiUrl",
			ApiKey:                "apiKey",
		}
		us := NewUsefulService(WithConfig(cfg))
		require.NotNil(t, us)
		assert.EqualValues(t, cfg, us.cfg)
		assert.NotNil(t, us.server)
		assert.NotNil(t, us.api)
		assert.NotNil(t, us.totalAmountCache)
	})
}
