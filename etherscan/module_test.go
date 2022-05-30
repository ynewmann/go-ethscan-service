package etherscan

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewModule(t *testing.T) {
	api := &api{}
	name := "test"

	m := NewModule(api, name)
	require.NotNil(t, m)
	assert.Equal(t, api, m.api)
	assert.EqualValues(t, name, m.Name)
	assert.EqualValues(t, name, m.GetName())
}
