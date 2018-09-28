package gengen

// StringStringMap maps a string to a string.
// This version is fast, but not safe for concurrent use.
type StringStringMap map[string]string

func NewStringStringMap() StringStringMap {
	return make(StringStringMap)
}

func (o StringStringMap) Copy() StringStringMap {
	m := NewStringStringMap()
	m.Merge(o)
	return m
}

// NewStringStringMapFrom creates a new StringStringMap from a
// StringStringMapI interface object
func NewStringStringMapFrom(i StringStringMapI) StringStringMap {
	m := NewStringStringMap()
	m.Merge(i)
	return m
}

// SetChanged sets the key to the value and returns a boolean indicating whether doing this caused
// the map to change.
func (o StringStringMap) SetChanged(key string, val string) (changed bool) {
	var ok bool
	var oldVal string

	if o == nil {
		panic("StringStringMap is not initialized.")
	}

	if oldVal, ok = o[key]; !ok || oldVal != val {
		o[key] = val
		changed = true
	}
	return
}

// Set will set the key to the value and return itself for easy chaining
func (o StringStringMap) Set(key string, val string) StringStringMap {
	o[key] = val
	return o
}

// Get returns the string based on its key. If it does not exist, an empty string will be returned.
func (o StringStringMap) Get(key string) (val string) {
	return o[key]
}

// Has returns true if the give key exists in the map
func (o StringStringMap) Has(key string) (exists bool) {
	_, exists = o[key]
	return
}


// Remove deletes the given key from the map
func (o StringStringMap) Remove(key string) {
	delete(o, key)
}

// RemoveAll resets the map to an empty map
func (o *StringStringMap) RemoveAll() {
	*o = NewStringStringMap()
}

// Values returns a slice of the values
func (o StringStringMap) Values() []string {
	vals := make([]string, 0, len(o))

	for _, v := range o {
		vals = append(vals, v)
	}
	return vals
}

// Keys returns a slice of the keys
func (o StringStringMap) Keys() []string {
	keys := make([]string, 0, len(o))

	for k := range o {
		keys = append(keys, k)
	}
	return keys
}

// Len returns the number of items in the map
func (o StringStringMap) Len() int {
	return len(o)
}

// Range will call the given function with every key and value.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o StringStringMap) Range(f func(key string, value string) bool) {
	if o == nil {
		return
	}
	for k, v := range o {
		if !f(k, v) {
			break
		}
	}
}

// Merge merges the given string map with the current one. The given one takes precedent on collisions.
func (o StringStringMap) Merge(i StringStringMapI) {
	if i == nil {
		return
	}
	i.Range(func(k string, v string) bool {
		o[k] = i.Get(k)
		return true
	})
}

// Equals returns true if all the keys in the given map exist in this map, and the values are the same
func (o StringStringMap) Equals(i StringStringMapI) bool {
	if i == nil {
		return false
	}
	if i.Len() != o.Len() {
		return false
	}
	var ret bool = true

	i.Range(func(k string, v string) bool {
		if !o.Has(k) || o[k] != i.Get(k) {
			ret = false
			return false // stop iterating
		}
		return true
	})

	return ret
}
