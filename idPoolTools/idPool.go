package IdPoolTools

import (
	"strconv"

	"sync"
)

// IDPool represents a pool of IDs that can be managed concurrently
// via local usage and external events.
//
// An intermediate state (leased) is introduced to the life cycle
// of an ID in the pool, in order to prevent lost updates to the
// pool that can occur as a result of employing both management schemes
// simultaneously.
// Local usage of an ID becomes a two stage process of leasing
// the ID from the pool, and later, Use()ing or Release()ing the ID on
// the pool upon successful or unsuccessful usage respectively,
//
// The table below shows the state transitions in the ID's life cycle.
// In the case of LeaseAvailableID() the ID is returned rather
// than provided as an input to the operation.
// All ID's begin in the available state.
/*
---------------------------------------------------------------------
|state\event   | LeaseAvailableID | Release | Use | Insert | Remove |
---------------------------------------------------------------------
|1 available   |        2         |    *    |  *  |   *    |   3    |
---------------------------------------------------------------------
|2 leased      |        **        |    1    |  3  |   *    |   3    |
---------------------------------------------------------------------
|3 unavailable |        **        |    *    |  *  |   1    |   *    |
---------------------------------------------------------------------
*  The event has no effect.
** This is guaranteed never to occur.
*/

type IDPool struct {
	mutex     sync.Mutex    `json:"-"`
	StartFrom ID            `json:"start_from"`
	EndTo     ID            `json:"end_to"`
	Members   map[string]ID `json:"members"`
	IdCache   *IdCache      `json:"cache"`
}

func IsValid(id ID) bool {
	return ((id > 0) && (id <= 9223372036854775807))
}

func (p *IDPool) IsValid() bool {
	return IsValid(p.StartFrom) && IsValid(p.EndTo) && (p.StartFrom < p.EndTo)
}

// String returns the string representation of an allocated ID
func (i ID) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

// NewIDPool returns a new ID pool
func NewIDPool(StartFrom ID, EndTo ID) *IDPool {
	return &IDPool{
		StartFrom: StartFrom,
		EndTo:     EndTo,
		IdCache:   newIDCache(StartFrom, EndTo),
		Members:   make(map[string]ID),
	}
}

// LeaseAvailableID returns an available ID at random from the pool.
// Returns an ID or NoID if no there is no available ID in the pool.
func (p *IDPool) LeaseAvailableID() ID {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.IdCache.leaseAvailableID()
}

// AllocateID returns a random available ID. Unlike LeaseAvailableID, the ID is
// immediately marked for use and there is no need to call Use().
func (p *IDPool) AllocateID(name string) ID {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	id := p.IdCache.allocateID()
	p.Members[name] = id
	return id
}

// Release returns a leased ID back to the pool.
// This operation accounts for IDs that were previously leased
// from the pool but were unused, e.g if allocation was unsuccessful.
// Thus, it has no effect if the ID is not currently leased in the
// pool, or the pool has since been refreshed.
//
// Returns true if the ID was returned back to the pool as
// a result of this call.
func (p *IDPool) Release(id ID) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for k, v := range p.Members {
		if v == id {
			delete(p.Members, k)
      return p.IdCache.insert(id)
		}
	}
	return p.IdCache.release(id)
}

// Use makes a leased ID unavailable in the pool and has no effect
// otherwise. Returns true if the ID was made unavailable
// as a result of this call.
func (p *IDPool) Use(id ID) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.IdCache.use(id)
}

// Insert makes an unavailable ID available in the pool
// and has no effect otherwise. Returns true if the ID
// was added back to the pool.
func (p *IDPool) Insert(id ID) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.IdCache.insert(id)
}

// Remove makes an ID unavailable in the pool.
// Returns true if the ID was previously available in the pool.
func (p *IDPool) Remove(id ID) bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.IdCache.remove(id)
}
