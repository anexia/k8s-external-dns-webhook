package anexia

import (
	"context"
	"fmt"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	"go.anx.io/go-anxcloud/pkg/api"
	"go.anx.io/go-anxcloud/pkg/api/types"
	anxcloudDns "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
	"go.anx.io/go-anxcloud/pkg/client"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

type DNSClient struct {
	client types.API
	dryRun bool
}

type DNSService interface {
	GetZones(ctx context.Context) ([]*anxcloudDns.Zone, error)
	GetRecords(ctx context.Context) ([]*anxcloudDns.Record, error)
	GetRecordsByZoneNameAndName(ctx context.Context, zoneName, name string) ([]*anxcloudDns.Record, error)
	DeleteRecord(ctx context.Context, record *anxcloudDns.Record) error
	CreateRecord(ctx context.Context, record *anxcloudDns.Record) error
}

func (c *DNSClient) GetZones(ctx context.Context) ([]*anxcloudDns.Zone, error) {
	log.Debugf("get all zones ...")
	channel := make(types.ObjectChannel)

	if err := c.client.List(ctx, &anxcloudDns.Zone{},
		api.ObjectChannel(&channel),
	); err != nil {
		return nil, fmt.Errorf("failed to list zones while getting all zones: %w", err)
	}

	zones := make([]*anxcloudDns.Zone, 0)
	for res := range channel {
		zone := anxcloudDns.Zone{}
		if err := res(&zone); err != nil {
			return nil, fmt.Errorf("failed to parse zone: %w", err)
		}
		zones = append(zones, &zone)
	}

	return zones, nil
}

func (c *DNSClient) GetRecords(ctx context.Context) ([]*anxcloudDns.Record, error) {
	log.Debugf("get all records ...")
	channel := make(types.ObjectChannel)

	allZones, err := c.GetZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get zones: %w", err)
	}

	for _, zone := range allZones {
		log.Debugf("get records for zone %s ...", zone.Name)
		zoneName := zone.Name

		if err := c.client.List(ctx, &anxcloudDns.Record{ZoneName: zoneName},
			api.ObjectChannel(&channel),
		); err != nil {
			return nil, fmt.Errorf("failed to list records for zone %s: %w", zoneName, err)
		}
	}

	records := make([]*anxcloudDns.Record, 0)
	for res := range channel {
		record := anxcloudDns.Record{}
		if err := res(&record); err != nil {
			return nil, fmt.Errorf("failed to parse record: %w", err)
		}
		records = append(records, &record)
	}

	return records, nil
}

func (c *DNSClient) GetRecordsByZoneNameAndName(ctx context.Context, zoneName, name string) ([]*anxcloudDns.Record, error) {
	log.Debugf("get records for zone %s and name %s ...", zoneName, name)
	channel := make(types.ObjectChannel)

	if err := c.client.List(ctx, &anxcloudDns.Record{ZoneName: zoneName, Name: name}, api.ObjectChannel(&channel)); err != nil {
		return nil, fmt.Errorf("failed to list records for zone %s and name %s: %w", zoneName, name, err)
	}

	records := make([]*anxcloudDns.Record, 0)
	for res := range channel {
		record := anxcloudDns.Record{}
		if err := res(&record); err != nil {
			return nil, fmt.Errorf("failed to parse record: %w", err)
		}
		records = append(records, &record)
	}

	return records, nil
}

func (c *DNSClient) DeleteRecord(ctx context.Context, record *anxcloudDns.Record) error {
	if c.dryRun {
		log.Infof("dry run: would delete record %v", record)
		return nil
	}
	log.Debugf("delete record %v ...", record)
	err := c.client.Destroy(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to delete record %v: %w", record, err)
	}
	log.Debug("record deleted")
	return nil
}

func (c *DNSClient) CreateRecord(ctx context.Context, record *anxcloudDns.Record) error {
	if c.dryRun {
		log.Infof("dry run: would create record %v", record)
		return nil
	}
	log.Debugf("create record %v ...", record)
	err := c.client.Create(ctx, record)
	if err != nil {
		return fmt.Errorf("failed to create record %v: %w", record, err)
	}
	log.Debug("record created")
	return nil
}

type Provider struct {
	provider.BaseProvider
	client       DNSService
	domainFilter *endpoint.DomainFilter
}

// NewProvider returns an instance of new provider.
func NewProvider(configuration *Configuration, domainFilter *endpoint.DomainFilter) (*Provider, error) {
	client, err := createClient(configuration)
	if err != nil {
		return nil, fmt.Errorf("failed to create Anexia client: %w", err)
	}
	prov := &Provider{
		client:       &DNSClient{client: client, dryRun: configuration.DryRun},
		domainFilter: domainFilter,
	}
	return prov, nil
}

func createClient(configuration *Configuration) (types.API, error) {
	options := []client.Option{
		client.TokenFromString(configuration.APIToken),
	}

	if configuration.APIEndpointURL == "" {
		log.Warn("API endpoint URL is not set, using default")
	} else {
		log.Debugf("Creating Anexia client with base URL %s", configuration.APIEndpointURL)
		options = append(options, client.BaseURL(configuration.APIEndpointURL))
	}
	apiClient, err := api.NewAPI(
		api.WithClientOptions(
			options...,
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create Anexia client: %w", err)
	}
	if configuration.DryRun {
		log.Warnf("Dry run mode enabled, no changes will be made")
	}
	return apiClient, nil
}

func (p *Provider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	records, err := p.client.GetRecords(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get records: %w", err)
	}

	// we have to do a little merge magic for records that are identical except for the data
	// while externalDNS target can be a list of IPs, in clouddns we have multiple records with different data
	groups := make(map[string][]*endpoint.Endpoint, 0)
	for _, record := range records {
		ep := recordToEndpoint(record)
		if p.domainFilter.IsConfigured() && !p.domainFilter.Match(ep.DNSName) {
			log.Infof("Skipping record %s because it was filtered out by the domain filter", ep.DNSName)
			continue
		}
		key := ep.DNSName + ep.RecordType + fmt.Sprint(ep.RecordTTL)
		groups[key] = append(groups[key], ep)
	}

	mergedEndpoints := make([]*endpoint.Endpoint, 0)
	for _, endpoints := range groups {
		mergedEndpoint := &endpoint.Endpoint{
			DNSName:    endpoints[0].DNSName,
			RecordType: endpoints[0].RecordType,
			RecordTTL:  endpoints[0].RecordTTL,
		}
		for _, ep := range endpoints {
			mergedEndpoint.Targets = append(mergedEndpoint.Targets, ep.Targets...)
		}
		mergedEndpoints = append(mergedEndpoints, mergedEndpoint)
	}
	return mergedEndpoints, nil
}

func recordToEndpoint(record *anxcloudDns.Record) *endpoint.Endpoint {
	return &endpoint.Endpoint{
		DNSName:    record.Name + "." + record.ZoneName,
		RecordTTL:  endpoint.TTL(record.TTL),
		RecordType: record.Type,
		Targets:    []string{record.RData},
	}
}

func (p *Provider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	epToCreate, epToDelete := GetCreateDeleteSetsFromChanges(changes)
	log.Debugf("apply changes, create: %d, delete: %d", len(epToCreate), len(epToDelete))
	allZones, err := p.client.GetZones(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all dns zones during apply: %w", err)
	}
	if len(allZones) == 0 {
		log.Warnf("did not retrieve any dns zones")
		return nil
	}
	recordsToDelete := make([]*anxcloudDns.Record, 0)
	for _, ep := range epToDelete {
		if p.domainFilter.IsConfigured() && !p.domainFilter.Match(ep.DNSName) {
			log.Debugf("Skipping record %s because it was filtered out by the domain filter", ep.DNSName)
			continue
		}
		potentialZones := FilterZonesByDomainName(allZones, ep.DNSName)
		for _, zone := range potentialZones {
			recordName := strings.TrimSuffix(ep.DNSName, "."+zone.Name)
			records, err := p.client.GetRecordsByZoneNameAndName(ctx, zone.Name, recordName)
			if err != nil {
				log.Errorf("failed to get records for zone %s and name %s: %v", zone.Name, recordName, err)
				break
			}
			for _, record := range records {
				if record.Type != ep.RecordType {
					continue
				}
				for _, target := range ep.Targets {
					if record.RData == target {
						recordsToDelete = append(recordsToDelete, record)
						break
					}
				}
			}
		}
	}

	for _, record := range recordsToDelete {
		if err := p.client.DeleteRecord(ctx, record); err != nil {
			return fmt.Errorf("failed to delete record %v: %w", record, err)
		}
	}

	recordsToCreate := make([]*anxcloudDns.Record, 0)
	for _, ep := range epToCreate {
		if p.domainFilter.IsConfigured() && !p.domainFilter.Match(ep.DNSName) {
			log.Debugf("Skipping record %s because it was filtered out by the domain filter", ep.DNSName)
			continue
		}
		zone := FilterZonesByDomainName(allZones, ep.DNSName)

		if len(zone) == 0 {
			log.Warnf("no zone found for domain %s", ep.DNSName)
			continue
		}
		for _, target := range ep.Targets {
			recordsToCreate = append(recordsToCreate, CreateRecord(zone[0].Name, ep.DNSName, target, int(ep.RecordTTL), ep.RecordType))
		}
	}

	for _, record := range recordsToCreate {
		if err := p.client.CreateRecord(ctx, record); err != nil {
			return fmt.Errorf("failed to create record %v: %w", record, err)
		}
	}

	return nil
}

func FilterZonesByDomainName(zones []*anxcloudDns.Zone, domainName string) []*anxcloudDns.Zone {
	possibleZones := make([]*anxcloudDns.Zone, 0)
	for _, zone := range zones {
		if strings.HasSuffix(domainName, zone.Name) {
			possibleZones = append(possibleZones, zone)
		}
	}

	// sort zones by length, longest first
	// this is necessary because the domain name might match multiple zones
	// and we want to use the most specific one
	sort.Slice(possibleZones, func(i, j int) bool {
		return len(possibleZones[i].Name) > len(possibleZones[j].Name)
	})
	return possibleZones
}
