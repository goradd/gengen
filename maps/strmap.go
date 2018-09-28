package maps

// StringMap maps a string to a string.
// This version is not safe for concurrent use.
// A zero value is ready for use
type StringMap struct {
    items map[string]string
}

// NewStringMap creates a new map that maps string's to string's.
func NewStringMap() *StringMap {
	return new(StringMap)
}

// NewStringMapFrom creates a new StringMap from a
// StringMapI interface object
func NewStringMapFrom(i StringMapI) *StringMap {
	m := NewStringMap()
	m.Merge(i)
	return m
}

// SetChanged sets the key to the value and returns a boolean indicating whether doing this caused
// the map to change.
func (o *StringMap) SetChanged(key string, val string) (changed bool) {
	var ok bool
	var oldVal string

	if o == nil {
		panic("The map must be created before being used.")
	}

	if o.items == nil {
	    o.items = make(map[string]string)
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
		o.items[key] = val
		changed = true
	}
	return
}

// Set will set the key to the value
func (o *StringMap) Set(key string, val string) {
	if o == nil {
		panic("The map must be initialized before being used.")
	}

    if o.items == nil {
        o.items = make(map[string]string)
    }

	o.items[key] = val
}

// Get returns the string based on its key. If it does not exist, an empty string will be returned.
func (o *StringMap) Get(key string) (val string) {
    if o == nil || o.items == nil {
		return
	}
	val,_ = o.items[key]
	return
}

// Has returns true if the given key exists in the map
func (o *StringMap) Has(key string) (exists bool) {
    if o == nil || o.items == nil {
		return
	}
	_, exists = o.items[key]
	return
}


// Remove deletes the given key from the map
func (o *StringMap) Remove(key string) {
    if o == nil || o.items == nil {
		return
	}
	delete(o.items, key)
}

// Clear resets the map to an empty map
func (o *StringMap) Clear() {

    if o == nil || o.items == nil {
		return
	}
	o.items = nil
}

// Values returns a slice of the values
func (o *StringMap) Values() []string {
    if o == nil {
        return nil
    }

	vals := make([]string, 0, len(o.items))

    for _, v := range o.items {
        vals = append(vals, v)
    }

	return vals
}

// Keys returns a slice of the keys
func (o *StringMap) Keys() []string {
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
func (o *StringMap) Len() int {
    if o == nil {
		return 0
	}
	return len(o.items)
}

// Range will call the given function with every key and value.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *StringMap) Range(f func(key string, value string) bool) {
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
func (o *StringMap) Merge(i StringMapI) {
	if o == nil {
		panic("The map must be created before being used.")
	}

	if o.items == nil {
	    o.items = make(map[string]string)
	}
	i.Range(func(k string, v string) bool {
		o.items[k] = i.Get(k)
		return true
	})
}

// Equals returns true if all the keys in the given map exist in this map, and the values are the same
func (o *StringMap) Equals(i StringMapI) bool {
	if i.Len() != o.Len() {
		return false
	}
	var ret = true

	i.Range(func(k string, v string) bool {
		if v2,ok := o.items[k]; !ok || v2 != v {
			ret = false
			return false // stop iterating
		}
		return true
	})

	return ret
}

func (o *StringMap) Copy() *StringMap {
	m := NewStringMap()
	m.Merge(o)
	return m
}

