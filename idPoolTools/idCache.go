package IdPoolTools

type idCache struct {
	// ids is a slice of IDs available in this idCache.
	ids map[ID]struct{}

	// leased is the set of IDs that are leased in this idCache.
	leased map[ID]struct{}
}

func newIDCache(start_from ID, end_to ID) *idCache {
	n := int(end_to - start_from + 1)
	if n < 0 {
		n = 0
	}

	c := &idCache{
		ids:    make(map[ID]struct{}, n),
		leased: make(map[ID]struct{}),
	}

	for id := start_from; id < end_to+1; id++ {
		c.ids[id] = struct{}{}
	}

	return c
}

// allocateID returns a random available ID without leasing it
func (c *idCache) allocateID() ID {
	for id := range c.ids {
		delete(c.ids, id)
		return id
	}

	return NoID
}

// leaseAvailableID returns a random available ID.
func (c *idCache) leaseAvailableID() ID {
	id := c.allocateID()
	if id == NoID {
		return NoID
	}

	// Mark as leased
	c.leased[id] = struct{}{}

	return id
}

// release makes the ID available again if it is currently
// leased and has no effect otherwise. Returns true if the
// ID was made available as a result of this call.
func (c *idCache) release(id ID) bool {
	if _, exists := c.leased[id]; !exists {
		return false
	}

	delete(c.leased, id)
	c.insert(id)

	return true
}

// use makes the ID unavailable if it is currently
// leased and has no effect otherwise. Returns true if the
// ID was made unavailable as a result of this call.
func (c *idCache) use(id ID) bool {
	if _, exists := c.leased[id]; !exists {
		return false
	}

	delete(c.leased, id)
	return true
}

// insert adds the ID into the cache if it is currently unavailable.
// Returns true if the ID was added to the cache.
func (c *idCache) insert(id ID) bool {
	if _, ok := c.ids[id]; ok {
		return false
	}

	if _, exists := c.leased[id]; exists {
		return false
	}

	c.ids[id] = struct{}{}
	return true
}

// remove removes the ID from the cache.
// Returns true if the ID was available in the cache.
func (c *idCache) remove(id ID) bool {
	delete(c.leased, id)

	if _, ok := c.ids[id]; ok {
		delete(c.ids, id)
		return true
	}

	return false
}
