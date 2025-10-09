package anexia

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitConfiguration(t *testing.T) {
	conf, err := InitConfiguration()
	require.Error(t, err)
	assert.Nil(t, conf)

	const apiToken = "tokenvalue"

	t.Setenv("ANEXIA_API_TOKEN", apiToken)
	conf, err = InitConfiguration()
	require.NoError(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, apiToken, conf.APIToken)
}
