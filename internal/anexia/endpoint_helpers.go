package anexia

import (
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
)

func GetCreateDeleteSetsFromChanges(changes *plan.Changes) ([]*endpoint.Endpoint, []*endpoint.Endpoint) {
	toCreate := make([]*endpoint.Endpoint, 0, len(changes.Create)+len(changes.UpdateNew))
	toDelete := make([]*endpoint.Endpoint, 0, len(changes.Delete)+len(changes.UpdateOld))

	toCreate = append(toCreate, changes.Create...)
	toDelete = append(toDelete, changes.Delete...)

	for i, updateOldEndpoint := range changes.UpdateOld {
		updateNewEndpoint := changes.UpdateNew[i]
		if endpointsAreDifferent(*updateOldEndpoint, *updateNewEndpoint) {
			toDelete = append(toDelete, updateOldEndpoint)
			toCreate = append(toCreate, updateNewEndpoint)
		}
	}
	return toCreate, toDelete
}

func endpointsAreDifferent(a endpoint.Endpoint, b endpoint.Endpoint) bool {
	return a.DNSName != b.DNSName || a.RecordType != b.RecordType ||
		a.RecordTTL != b.RecordTTL || !a.Targets.Same(b.Targets)
}
