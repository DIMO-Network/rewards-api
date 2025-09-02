package vinvc

import (
	"iter"
	"sync"

	"github.com/DIMO-Network/attestation-api/pkg/types"
)

// VINStore is a concurrent safe store for VIN subjects.
type VINStore struct {
	subjectsByVIN map[string][]*types.VINSubject
	readLock      sync.RWMutex
}

// Add adds a subject to the VINStore for a given VIN.
func (v *VINStore) Add(vin string, subject *types.VINSubject) {
	v.readLock.Lock()
	defer v.readLock.Unlock()
	if v.subjectsByVIN == nil {
		v.subjectsByVIN = make(map[string][]*types.VINSubject)
	}
	v.subjectsByVIN[vin] = append(v.subjectsByVIN[vin], subject)
}

// Get returns the subjects associated with a VIN.
func (v *VINStore) Get(vin string) []*types.VINSubject {
	v.readLock.RLock()
	defer v.readLock.RUnlock()
	return v.subjectsByVIN[vin]
}

// Items returns a concurrent safe iterator over the VINStore.
func (v *VINStore) Items() iter.Seq2[string, []*types.VINSubject] {
	return func(yield func(string, []*types.VINSubject) bool) {
		v.readLock.RLock()
		defer v.readLock.RUnlock()
		for vin, subjects := range v.subjectsByVIN {
			if !yield(vin, subjects) {
				return
			}
		}
	}
}
