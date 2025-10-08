package anexia

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Configuration holds configuration from environmental variables.
type Configuration struct {
	APIToken       string `env:"ANEXIA_API_TOKEN,notEmpty"`
	APIEndpointURL string `env:"ANEXIA_API_URL" envDefault:"https://engine.anexia-it.com/"`
	DryRun         bool   `env:"DRY_RUN" envDefault:"false"`
}

// InitConfiguration sets up configuration by reading set environmental variables.
func InitConfiguration() (*Configuration, error) {
	cfg := &Configuration{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("reading anexia configuration failed: %w", err)
	}
	return cfg, nil
}
