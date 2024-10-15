package sql

import (
	"sync"
)

type registry struct {
	mu      sync.RWMutex
	drivers map[string]struct{}
}

var driverRegistry *registry //nolint:gochecknoglobals

func init() {
	driverRegistry = &registry{
		drivers: make(map[string]struct{}),
	}
}

func (r *registry) register(driverName string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := driverName

	r.drivers[id] = struct{}{}
}

func (r *registry) lookup(id string) (bool, string) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, found := r.drivers[id]
	if !found {
		return false, ""
	}

	return true, id
}

// RegisterDriver registers an SQL database driver.
func RegisterDriver(driverName string) {
	driverRegistry.register(driverName)
}

// lookupDriver search SQL database driver by id in the registry.
// Returns true and the driver name if the driver is registered,
// otherwise returns false and an empty string.
func lookupDriver(id string) (bool, string) {
	return driverRegistry.lookup(id)
}
