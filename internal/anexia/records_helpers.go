package anexia

import (
	"strings"

	log "github.com/sirupsen/logrus"
	anxcloudDns "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
)

func CreateRecord(zoneName string, endpointName string, target string, recordTTL int, recordType string) *anxcloudDns.Record {
	dName := DomainNameFor(zoneName, endpointName)
	return &anxcloudDns.Record{
		ZoneName: zoneName,
		Name:     dName,
		RData:    target,
		TTL:      recordTTL,
		Type:     recordType,
	}
}

func DomainNameFor(zoneName string, endpointName string) string {
	name := strings.TrimSuffix(endpointName, "."+zoneName)

	if endpointName == zoneName {
		log.Debugf("Found a root domain(%s), converting %s to @", endpointName, name)
		name = "@"
	}

	return name
}
