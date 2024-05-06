package anexia

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"context"

	"github.com/caarlos0/env/v11"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	anxcloudDns "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func TestNewProvider(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	t.Setenv("ANEXIA_API_TOKEN", "1")
	anexiaConfig := &Configuration{}

	err := env.Parse(anexiaConfig)
	require.NoError(t, err)
	domainFilter := endpoint.NewDomainFilter([]string{"a.de."})
	p, err := NewProvider(anexiaConfig, domainFilter)
	require.NoError(t, err)
	require.Equal(t, true, p.domainFilter.IsConfigured())
	require.Equal(t, false, p.domainFilter.Match("b.de."))
}

func TestRecords(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()
	testCases := []struct {
		name              string
		givenRecords      []*anxcloudDns.Record
		givenError        error
		givenDomainFilter endpoint.DomainFilter
		expectedEndpoints []*endpoint.Endpoint
		expectedError     error
	}{
		{
			name:              "no records",
			givenRecords:      []*anxcloudDns.Record{},
			expectedEndpoints: []*endpoint.Endpoint{},
		},
		{
			name:              "error reading records",
			givenRecords:      []*anxcloudDns.Record{},
			givenError:        fmt.Errorf("test error"),
			expectedEndpoints: []*endpoint.Endpoint{},
			expectedError:     fmt.Errorf("test error"),
		},
		{
			name: "multiple A records",
			givenRecords: createRecordSlice(3, func(i int) (string, string, string, int, string) {
				recordName := "a" + fmt.Sprintf("%d", i+1)
				zoneName := "a.de"
				return recordName, zoneName, "A", ((i + 1) * 100), fmt.Sprintf("%d.%d.%d.%d", i+1, i+1, i+1, i+1)
			}),
			expectedEndpoints: createEndpointSlice(3, func(i int) (string, string, endpoint.TTL, []string) {
				return "a" + fmt.Sprintf("%d", i+1) + ".a.de", "A", endpoint.TTL((i + 1) * 100), []string{fmt.Sprintf("%d.%d.%d.%d", i+1, i+1, i+1, i+1)}
			}),
		},
		{
			name: "multiple records filtered by domain",
			givenRecords: createRecordSlice(6, func(i int) (string, string, string, int, string) {
				if i < 3 {
					recordName := "a" + fmt.Sprintf("%d", i+1)
					zoneName := "a.de"
					return recordName, zoneName, "A", ((i + 1) * 100), fmt.Sprintf("%d.%d.%d.%d", i+1, i+1, i+1, i+1)
				}
				recordName := "b" + fmt.Sprintf("%d", i+1)
				zoneName := "b.de"
				return recordName, zoneName, "A", ((i + 1) * 100), fmt.Sprintf("%d.%d.%d.%d", i+1, i+1, i+1, i+1)
			}),
			givenDomainFilter: endpoint.NewDomainFilter([]string{"a.de"}),
			expectedEndpoints: createEndpointSlice(3, func(i int) (string, string, endpoint.TTL, []string) {
				return "a" + fmt.Sprintf("%d", i+1) + ".a.de", "A", endpoint.TTL((i + 1) * 100), []string{fmt.Sprintf("%d.%d.%d.%d", i+1, i+1, i+1, i+1)}
			}),
		},
		{
			name: "records mapped to same endpoint",
			givenRecords: createRecordSlice(3, func(i int) (string, string, string, int, string) {
				if i < 2 {
					return "", "a.de", "A", 300, fmt.Sprintf("%d.%d.%d.%d", i+1, i+1, i+1, i+1)
				}
				return "", "c.de", "A", 300, fmt.Sprintf("%d.%d.%d.%d", i+1, i+1, i+1, i+1)

			}),
			expectedEndpoints: createEndpointSlice(2, func(i int) (string, string, endpoint.TTL, []string) {
				if i == 0 {
					return "a.de", "A", endpoint.TTL(300), []string{"1.1.1.1", "2.2.2.2"}
				}
				return "c.de", "A", endpoint.TTL(300), []string{"3.3.3.3"}
			}),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDNSClient := &mockDNSClient{
				allRecords:  tc.givenRecords,
				returnError: tc.givenError,
			}
			provider := &Provider{client: mockDNSClient, domainFilter: tc.givenDomainFilter}
			endpoints, err := provider.Records(ctx)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Equal(t, tc.expectedError, err)
				return
			}
			require.NoError(t, err)
			require.Len(t, endpoints, len(tc.expectedEndpoints))
			assert.ElementsMatch(t, tc.expectedEndpoints, endpoints)
		})
	}
}

func TestApplyChanges(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	log.SetReportCaller(true)
	deZoneName := "de"
	comZoneName := "com"
	ctx := context.Background()
	testCases := []struct {
		name                   string
		givenRecords           []*anxcloudDns.Record
		givenZones             []*anxcloudDns.Zone
		givenZoneRecords       map[string][]*anxcloudDns.Record
		givenDomainFilter      endpoint.DomainFilter
		whenChanges            *plan.Changes
		expectedRecordsCreated map[string][]*anxcloudDns.Record
		expectedRecordsDeleted map[string][]string
	}{
		{
			name:                   "no changes",
			givenZones:             createZoneSlice(0, nil),
			givenZoneRecords:       map[string][]*anxcloudDns.Record{},
			whenChanges:            &plan.Changes{},
			expectedRecordsCreated: nil,
			expectedRecordsDeleted: nil,
		},
		{
			name: "create one record in a blank zone",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(0, nil),
			},
			whenChanges: &plan.Changes{
				Create: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a." + deZoneName, "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
			},
			expectedRecordsCreated: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "1.2.3.4"
				}),
			},
			expectedRecordsDeleted: nil,
		},
		{
			name: "create a record which is filtered out from the domain filter",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(0, nil),
			},
			givenDomainFilter: endpoint.NewDomainFilter([]string{"b.de"}),
			whenChanges: &plan.Changes{
				Create: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a." + deZoneName, "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
			},
			expectedRecordsCreated: nil,
			expectedRecordsDeleted: nil,
		},
		{
			name: "create 2 records from one endpoint in a blank zone",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(0, nil),
			},
			whenChanges: &plan.Changes{
				Create: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a." + deZoneName, "A", endpoint.TTL(300), []string{"1.2.3.4", "5.6.7.8"}
				}),
			},
			expectedRecordsCreated: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(2, func(i int) (string, string, string, int, string) {
					if i == 0 {
						return "a", deZoneName, "A", 300, "1.2.3.4"
					}
					return "a", deZoneName, "A", 300, "5.6.7.8"
				}),
			},
			expectedRecordsDeleted: nil,
		},
		{
			name: "delete the only record in a zone",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "1.2.3.4"
				}),
			},
			whenChanges: &plan.Changes{
				Delete: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
			},
			expectedRecordsDeleted: map[string][]string{
				deZoneName: {"0"},
			},
		},
		{
			name: "delete a record which is filtered out from the domain filter",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "1.2.3.4"
				}),
			},
			givenDomainFilter: endpoint.NewDomainFilter([]string{"b.de"}),
			whenChanges: &plan.Changes{
				Delete: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
			},
			expectedRecordsDeleted: nil,
		},
		{
			name: "delete multiple records, in different zones",
			givenZones: createZoneSlice(2, func(i int) string {
				if i == 0 {
					return deZoneName
				}
				return comZoneName

			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(2, func(n int) (string, string, string, int, string) {
					if n == 0 {
						return "a", deZoneName, "A", 300, "1.2.3.4"
					}
					return "a", deZoneName, "A", 300, "5.6.7.8"

				}),
				comZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", comZoneName, "A", 300, "11.22.33.44"
				}),
			},
			whenChanges: &plan.Changes{
				Delete: createEndpointSlice(2, func(i int) (string, string, endpoint.TTL, []string) {
					if i == 0 {
						return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4", "5.6.7.8"}
					}
					return "a.com", "A", endpoint.TTL(300), []string{"11.22.33.44"}

				}),
			},
			expectedRecordsDeleted: map[string][]string{
				deZoneName:  {"0", "1"},
				comZoneName: {"0"},
			},
		},
		{
			name: "delete record which is not in the zone, deletes nothing",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(0, nil),
			},
			whenChanges: &plan.Changes{
				Delete: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
			},
			expectedRecordsDeleted: nil,
		},
		{
			name: "delete one record from targets part of endpoint",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "1.2.3.4"
				}),
			},
			whenChanges: &plan.Changes{
				Delete: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4", "5.6.7.8"}
				}),
			},
			expectedRecordsDeleted: map[string][]string{
				deZoneName: {"0"},
			},
		},
		{
			name: "update single record",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "1.2.3.4"
				}),
			},
			whenChanges: &plan.Changes{
				UpdateOld: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
				UpdateNew: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"5.6.7.8"}
				}),
			},
			expectedRecordsDeleted: map[string][]string{
				deZoneName: {"0"},
			},
			expectedRecordsCreated: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "5.6.7.8"
				}),
			},
		},
		{
			name: "update a record which is filtered out by domain filter, does nothing",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "1.2.3.4"
				}),
			},
			givenDomainFilter: endpoint.NewDomainFilter([]string{"b.de"}),

			whenChanges: &plan.Changes{
				UpdateOld: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
				UpdateNew: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"5.6.7.8"}
				}),
			},
			expectedRecordsDeleted: nil,
			expectedRecordsCreated: nil,
		},
		{
			name: "update when old and new endpoint are the same, does nothing",
			givenZones: createZoneSlice(1, func(_ int) string {
				return deZoneName
			}),
			givenZoneRecords: map[string][]*anxcloudDns.Record{
				deZoneName: createRecordSlice(1, func(_ int) (string, string, string, int, string) {
					return "a", deZoneName, "A", 300, "1.2.3.4"
				}),
			},
			whenChanges: &plan.Changes{
				UpdateOld: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
				UpdateNew: createEndpointSlice(1, func(_ int) (string, string, endpoint.TTL, []string) {
					return "a.de", "A", endpoint.TTL(300), []string{"1.2.3.4"}
				}),
			},
			expectedRecordsDeleted: nil,
			expectedRecordsCreated: nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDNSClient := &mockDNSClient{
				allRecords:  tc.givenRecords,
				allZones:    tc.givenZones,
				zoneRecords: tc.givenZoneRecords,
			}
			provider := &Provider{client: mockDNSClient, domainFilter: tc.givenDomainFilter}
			err := provider.ApplyChanges(ctx, tc.whenChanges)

			require.NoError(t, err)
			require.Len(t, mockDNSClient.createdRecords, len(tc.expectedRecordsCreated))

			for zoneName, expectedDeletedRecordIDs := range tc.expectedRecordsDeleted {
				require.Len(t, mockDNSClient.deletedRecords[zoneName], len(expectedDeletedRecordIDs), "deleted records in zone '%s' do not fit", zoneName)
				actualDeletedRecordIDs, ok := mockDNSClient.deletedRecords[zoneName]
				require.True(t, ok)
				assert.ElementsMatch(t, expectedDeletedRecordIDs, actualDeletedRecordIDs)
			}
		})
	}
}

func TestAdjustEndpoints(t *testing.T) {
	provider := &Provider{}
	endpoints := createEndpointSlice(rand.Intn(5), func(_ int) (string, string, endpoint.TTL, []string) {
		return RandStringRunes(10), RandStringRunes(1), endpoint.TTL(300), []string{RandStringRunes(5)}
	})
	actualEndpoints, _ := provider.AdjustEndpoints(endpoints)
	require.Equal(t, endpoints, actualEndpoints)
}

type mockDNSClient struct {
	returnError    error
	allRecords     []*anxcloudDns.Record
	zoneRecords    map[string][]*anxcloudDns.Record
	allZones       []*anxcloudDns.Zone
	createdRecords map[string][]*anxcloudDns.Record // zoneName -> recordCreates
	deletedRecords map[string][]string              // zoneName -> recordIDs
}

func (c *mockDNSClient) GetRecords(_ context.Context) ([]*anxcloudDns.Record, error) {
	log.Debugf("GetAllRecords called")
	return c.allRecords, c.returnError
}

func (c *mockDNSClient) GetZoneRecords(_ context.Context, zoneName string) ([]*anxcloudDns.Record, error) {
	log.Debugf("GetZoneRecords called with zoneName %s", zoneName)
	return c.zoneRecords[zoneName], c.returnError
}

func (c *mockDNSClient) GetRecordsByZoneNameAndName(_ context.Context, zoneName, name string) ([]*anxcloudDns.Record, error) {
	log.Debugf("GetRecordsByzoneNameAndName called with zoneName %s and name %s", zoneName, name)
	result := make([]*anxcloudDns.Record, 0)
	recordsOfZone := c.zoneRecords[zoneName]
	for _, record := range recordsOfZone {
		if record.Name == name {
			result = append(result, record)
		}
	}
	return result, c.returnError
}

func (c *mockDNSClient) GetZones(_ context.Context) ([]*anxcloudDns.Zone, error) {
	log.Debug("GetZones called ")
	if c.allZones != nil {
		for _, zone := range c.allZones {
			log.Debugf("GetZones: zone '%s'", zone.Name)
		}
	} else {
		log.Debug("GetZones: no zones")
	}
	return c.allZones, c.returnError
}

func (c *mockDNSClient) GetZonesByDomainName(_ context.Context, domainName string) ([]*anxcloudDns.Zone, error) {
	log.Debugf("GetZonesByDomainName called with domainName %s", domainName)
	result := make([]*anxcloudDns.Zone, 0)
	for _, zone := range c.allZones {
		if strings.HasSuffix(domainName, zone.Name) {
			result = append(result, zone)
		}
	}
	return result, c.returnError
}

func (c *mockDNSClient) CreateRecord(_ context.Context, zoneName string, record *anxcloudDns.Record) error {
	log.Debugf("CreateRecord called with zoneName %s and record %v", zoneName, record)
	if c.createdRecords == nil {
		c.createdRecords = make(map[string][]*anxcloudDns.Record)
	}
	c.createdRecords[zoneName] = append(c.createdRecords[zoneName], record)
	return c.returnError
}

func (c *mockDNSClient) DeleteRecord(_ context.Context, zoneName string, recordID string) error {
	log.Debugf("DeleteRecord called with zoneName %s and recordID %s", zoneName, recordID)
	if c.deletedRecords == nil {
		c.deletedRecords = make(map[string][]string)
	}
	c.deletedRecords[zoneName] = append(c.deletedRecords[zoneName], recordID)
	return c.returnError
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func createRecordSlice(count int, modifier func(int) (string, string, string, int, string)) []*anxcloudDns.Record {
	records := make([]*anxcloudDns.Record, count)
	for i := 0; i < count; i++ {
		name, zone, typ, ttl, target := modifier(i)
		records[i] = &anxcloudDns.Record{
			Name:       name,
			Type:       typ,
			TTL:        ttl,
			ZoneName:   zone,
			RData:      target,
			Identifier: fmt.Sprintf("%d", i),
		}
	}
	return records
}

func createEndpointSlice(count int, modifier func(int) (string, string, endpoint.TTL, []string)) []*endpoint.Endpoint {
	endpoints := make([]*endpoint.Endpoint, count)
	for i := 0; i < count; i++ {
		name, typ, ttl, targets := modifier(i)
		endpoints[i] = &endpoint.Endpoint{
			DNSName:    name,
			RecordType: typ,
			Targets:    targets,
			RecordTTL:  ttl,
		}
	}
	return endpoints
}

func createZoneSlice(count int, modifier func(int) string) []*anxcloudDns.Zone {
	zones := make([]*anxcloudDns.Zone, count)
	for i := 0; i < count; i++ {
		zoneName := modifier(i)
		zones[i] = &anxcloudDns.Zone{
			Name: zoneName,
		}
	}
	return zones
}
