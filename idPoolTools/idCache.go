package IdPoolTools

import ()

type IdCache struct {
	// Ids is a slice of IDs available in this IdCache.
	Ids map[ID]struct{} `json:"ids"`

	// Leased is the set of IDs that are Leased in this IdCache.
	Leased map[ID]struct{} `json:"leased"`
}

func newIDCache(start_from ID, end_to ID) *IdCache {
	n := int(end_to - start_from + 1)
	if n < 0 {
		n = 0
	}

	c := &IdCache{
		Ids:    make(map[ID]struct{}, n),
		Leased: make(map[ID]struct{}),
	}

	for id := start_from; id < end_to+1; id++ {
		c.Ids[id] = struct{}{}
	}

	return c
}

// allocateID returns a random available ID without leasing it
func (c *IdCache) allocateID() ID {
	for id := range c.Ids {
		delete(c.Ids, id)
		return id
	}

	return NoID
}

// leaseAvailableID returns a random available ID.
func (c *IdCache) leaseAvailableID() ID {
	id := c.allocateID()
	if id == NoID {
		return NoID
	}

	// Mark as Leased
	c.Leased[id] = struct{}{}

	return id
}

// release makes the ID available again if it is currently
// Leased and has no effect otherwise. Returns true if the
// ID was made available as a result of this call.
func (c *IdCache) release(id ID) bool {
	if _, exists := c.Leased[id]; !exists {
		return false
	}

	delete(c.Leased, id)
	c.insert(id)

	return true
}

// use makes the ID unavailable if it is currently
// Leased and has no effect otherwise. Returns true if the
// ID was made unavailable as a result of this call.
func (c *IdCache) use(id ID) bool {
	if _, exists := c.Leased[id]; !exists {
		return false
	}

	delete(c.Leased, id)
	return true
}

// insert adds the ID into the cache if it is currently unavailable.
// Returns true if the ID was added to the cache.
func (c *IdCache) insert(id ID) bool {
	if _, ok := c.Ids[id]; ok {
		return false
	}

	if _, exists := c.Leased[id]; exists {
		return false
	}

	c.Ids[id] = struct{}{}
	return true
}

// remove removes the ID from the cache.
// Returns true if the ID was available in the cache.
func (c *IdCache) remove(id ID) bool {
	delete(c.Leased, id)

	if _, ok := c.Ids[id]; ok {
		delete(c.Ids, id)
		return true
	}

	return false
}
