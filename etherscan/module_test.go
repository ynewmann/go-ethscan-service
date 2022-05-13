package etherscan

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewModule(t *testing.T) {
	client := NewApiClient()
	name := "test"

	m := NewModule(client, name)
	require.NotNil(t, m)
	assert.Equal(t, client, m.client)
	assert.EqualValues(t, name, m.Name)
	assert.EqualValues(t, name, m.GetName())
}
