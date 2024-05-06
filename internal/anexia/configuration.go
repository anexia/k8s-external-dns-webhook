package anexia

import (
	"github.com/caarlos0/env/v11"
	log "github.com/sirupsen/logrus"
)

// Configuration holds configuration from environmental variables
type Configuration struct {
	APIToken       string `env:"ANEXIA_API_TOKEN,notEmpty"`
	APIEndpointURL string `env:"ANEXIA_API_URL"`
	DryRun         bool   `env:"DRY_RUN" envDefault:"false"`
}

// Init sets up configuration by reading set environmental variables
func Init() Configuration {
	cfg := Configuration{}
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("error reading configuration from environment: %v", err)
	}
	return cfg
}
