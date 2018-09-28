package maps

// Map maps a string to a interface{}.
// This version is not safe for concurrent use.
// A zero value is ready for use
type Map struct {
    items map[string]interface{}
}

// NewMap creates a new map that maps string's to interface{}'s.
func NewMap() *Map {
	return new(Map)
}

// NewMapFrom creates a new Map from a
// MapI interface object
func NewMapFrom(i MapI) *Map {
	m := NewMap()
	m.Merge(i)
	return m
}

// SetChanged sets the key to the value and returns a boolean indicating whether doing this caused
// the map to change.
func (o *Map) SetChanged(key string, val interface{}) (changed bool) {
	var ok bool
	var oldVal interface{}

	if o == nil {
		panic("The map must be created before being used.")
	}

	if o.items == nil {
	    o.items = make(map[string]interface{})
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
		o.items[key] = val
		changed = true
	}
	return
}

// Set will set the key to the value
func (o *Map) Set(key string, val interface{}) {
	if o == nil {
		panic("The map must be initialized before being used.")
	}

    if o.items == nil {
        o.items = make(map[string]interface{})
    }

	o.items[key] = val
}

// Get returns the string based on its key. If it does not exist, an empty string will be returned.
func (o *Map) Get(key string) (val interface{}) {
    if o == nil || o.items == nil {
		return
	}
	val,_ = o.items[key]
	return
}

// Has returns true if the given key exists in the map
func (o *Map) Has(key string) (exists bool) {
    if o == nil || o.items == nil {
		return
	}
	_, exists = o.items[key]
	return
}


// Remove deletes the given key from the map
func (o *Map) Remove(key string) {
    if o == nil || o.items == nil {
		return
	}
	delete(o.items, key)
}

// Clear resets the map to an empty map
func (o *Map) Clear() {

    if o == nil || o.items == nil {
		return
	}
	o.items = nil
}

// Values returns a slice of the values
func (o *Map) Values() []interface{} {
    if o == nil {
        return nil
    }

	vals := make([]interface{}, 0, len(o.items))

    for _, v := range o.items {
        vals = append(vals, v)
    }

	return vals
}

// Keys returns a slice of the keys
func (o *Map) Keys() []string {
    if o == nil {
        return nil
    }

	keys := make([]string, 0, len(o.items))

    for k := range o.items {
        keys = append(keys, k)
    }
	return keys
}

// Len returns the number of items in the map
func (o *Map) Len() int {
    if o == nil {
		return 0
	}
	return len(o.items)
}

// Range will call the given function with every key and value.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *Map) Range(f func(key string, value interface{}) bool) {
	if o == nil {
		return
	}
	for k, v := range o.items {
		if !f(k, v) {
			break
		}
	}
}

// Merge merges the given string map with the current one. The given one takes precedent on collisions.
func (o *Map) Merge(i MapI) {
	if o == nil {
		panic("The map must be created before being used.")
	}

	if o.items == nil {
	    o.items = make(map[string]interface{})
	}
	i.Range(func(k string, v interface{}) bool {
		o.items[k] = i.Get(k)
		return true
	})
}

// Equals returns true if all the keys in the given map exist in this map, and the values are the same
func (o *Map) Equals(i MapI) bool {
	if i.Len() != o.Len() {
		return false
	}
	var ret = true

	i.Range(func(k string, v interface{}) bool {
		if v2,ok := o.items[k]; !ok || v2 != v {
			ret = false
			return false // stop iterating
		}
		return true
	})

	return ret
}

func (o *Map) Copy() *Map {
	m := NewMap()
	m.Merge(o)
	return m
}

