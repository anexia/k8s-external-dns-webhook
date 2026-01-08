package anexia

import (
	"strings"

	log "github.com/sirupsen/logrus"
	anxcloudDns "go.anx.io/go-anxcloud/pkg/apis/clouddns/v1"
)

func CreateRecord(zoneName string, endpointName string, target string, recordTTL int, recordType string) *anxcloudDns.Record {
	epName := strings.TrimSuffix(endpointName, "."+zoneName)
	if endpointName == zoneName {
		log.Debugf("Found a root domain(%s), converting %s to @", endpointName, epName)
		epName = "@"
	}
	return &anxcloudDns.Record{
		ZoneName: zoneName,
		Name:     epName,
		RData:    target,
		TTL:      recordTTL,
		Type:     recordType,
	}
}
