package dnsprovider

import (
	"testing"

	"github.com/probstenhias/external-dns-anexia-webhook/cmd/webhook/init/configuration"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	cases := []struct {
		name          string
		config        configuration.Config
		env           map[string]string
		expectedError string
	}{
		{
			name:   "minimal config for anexia provider",
			config: configuration.Config{},
			env: map[string]string{
				"ANEXIA_API_TOKEN": "token",
			},
		},
		{
			name: "domain filter config for anexia provider",
			config: configuration.Config{
				DomainFilter:   []string{"domain.com"},
				ExcludeDomains: []string{"sub.domain.com"},
			},
			env: map[string]string{
				"ANEXIA_API_TOKEN": "token",
			},
		},
		{
			name:          "empty configuration",
			config:        configuration.Config{},
			expectedError: "reading anexia configuration failed: env: environment variable \"ANEXIA_API_TOKEN\" should not be empty",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			for k, v := range tc.env {
				t.Setenv(k, v)
			}

			dnsProvider, err := Init(tc.config)

			if tc.expectedError != "" {
				assert.EqualError(t, err, tc.expectedError, "expecting error")
				return
			}

			assert.NoErrorf(t, err, "error creating provider")
			assert.NotNil(t, dnsProvider)
		})
	}
}
