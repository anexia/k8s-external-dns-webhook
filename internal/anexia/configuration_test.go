package anexia

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfiguration(t *testing.T) {
	os.Unsetenv("ANEXIA_API_TOKEN")
	conf, err := InitConfiguration()
	assert.Error(t, err)
	assert.Nil(t, conf)

	const apiToken = "tokenvalue"

	os.Setenv("ANEXIA_API_TOKEN", apiToken)
	conf, err = InitConfiguration()
	assert.NoError(t, err)
	assert.NotNil(t, conf)
	assert.Equal(t, apiToken, conf.APIToken)
}
