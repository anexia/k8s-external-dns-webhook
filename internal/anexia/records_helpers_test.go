package anexia

import (
	"testing"

	"github.com/stretchr/testify/assert"
	anxcloudDns "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
)

func TestCreateRecord(t *testing.T) {
	zoneName := "a.de"
	target := "1.2.3.4"
	ttl := 300
	typeRecord := "A"

	testCases := []struct {
		name           string
		givenRecord    *anxcloudDns.Record
		expectedRecord *anxcloudDns.Record
	}{
		{
			name:        "should return a record where the dns name is '@' when endpointName is the same as zoneName",
			givenRecord: CreateRecord(zoneName, "a.de", target, ttl, typeRecord),
			expectedRecord: &anxcloudDns.Record{
				ZoneName: zoneName,
				Name:     "@",
				RData:    target,
				TTL:      ttl,
				Type:     typeRecord,
			},
		},
		{
			name:        "should return the subdomain whend endpointName is a subdomain of zoneName",
			givenRecord: CreateRecord(zoneName, "sub.a.de", target, ttl, typeRecord),
			expectedRecord: &anxcloudDns.Record{
				ZoneName: zoneName,
				Name:     "sub",
				RData:    target,
				TTL:      ttl,
				Type:     typeRecord,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedRecord, tc.givenRecord)
		})
	}
}

func TestDomainNameFor(t *testing.T) {
	zoneName := "a.de"

	testCases := []struct {
		name                 string
		givenEndpointName    string
		expectedEndpointName string
	}{
		{
			name:                 "should return '@' when the zoneName is the same as the endpointName",
			givenEndpointName:    DomainNameFor(zoneName, "a.de"),
			expectedEndpointName: "@",
		},
		{
			name:                 "should return 'sub' when the endpointName is a subdomain of zoneName",
			givenEndpointName:    DomainNameFor(zoneName, "sub.a.de"),
			expectedEndpointName: "sub",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedEndpointName, tc.givenEndpointName)
		})
	}
}
