package dnsprovider

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.anx.io/external-dns-webhook/cmd/webhook/init/configuration"
	"go.anx.io/external-dns-webhook/internal/anexia"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/provider"
)

func Init(config configuration.Config) (provider.Provider, error) {
	var domainFilter endpoint.DomainFilter
	createMsg := "Creating anexia provider with "

	if config.RegexDomainFilter != "" {
		createMsg += fmt.Sprintf("regexp domain filter: '%s', ", config.RegexDomainFilter)
		if config.RegexDomainExclusion != "" {
			createMsg += fmt.Sprintf("with exclusion: '%s', ", config.RegexDomainExclusion)
		}
		domainFilter = endpoint.NewRegexDomainFilter(
			regexp.MustCompile(config.RegexDomainFilter),
			regexp.MustCompile(config.RegexDomainExclusion),
		)
	} else {
		if len(config.DomainFilter) > 0 {
			createMsg += fmt.Sprintf("domain filter: '%s', ", strings.Join(config.DomainFilter, ","))
		}
		if len(config.ExcludeDomains) > 0 {
			createMsg += fmt.Sprintf("exclude domain filter: '%s', ", strings.Join(config.ExcludeDomains, ","))
		}
		domainFilter = endpoint.NewDomainFilterWithExclusions(config.DomainFilter, config.ExcludeDomains)
	}

	createMsg = strings.TrimSuffix(createMsg, ", ")
	if strings.HasSuffix(createMsg, "with ") {
		createMsg += "no kind of domain filters"
	}
	log.Info(createMsg)

	anexiaConfig, err := anexia.InitConfiguration()
	if err != nil {
		return nil, fmt.Errorf("creating configuration failed: %w", err)
	}

	anexiaProvider, err := anexia.NewProvider(anexiaConfig, domainFilter)
	if err != nil {
		return nil, fmt.Errorf("creating anexia provider failed: %w", err)
	}

	return anexiaProvider, nil
}
