package useful_servise

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_SaveAndLoad(t *testing.T) {
	cfgFolder, err := ioutil.TempDir("", "cfg-")
	require.NoError(t, err, err)
	defer os.RemoveAll(cfgFolder)

	cfg := &Config{
		UseCache:              true,
		MemoryCacheBackupPath: path.Join(cfgFolder, "backup"),
		CacheSize:             333,
		Address:               ":3333",
		ApiUrl:                "test.com",
		ApiKey:                "key",
	}

	cfgFile := path.Join(cfgFolder, "config")
	err = cfg.SaveToFile(cfgFile)
	require.NoError(t, err, err)

	loadedCfg := &Config{}
	err = loadedCfg.LoadFromFile(cfgFile)
	require.NoError(t, err, err)
	assert.EqualValues(t, cfg, loadedCfg)
}
