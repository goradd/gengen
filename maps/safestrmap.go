package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sync"

)

// SafeMap is your basic GO map with a read/write mutex so that it can read and write concurrently.
// The zero map is ready for use, but you may not copy it after first using it.
type SafeStringMap struct {
	sync.RWMutex
	items map[string]string
}

// NewSafeStringMap returns a newly created SafeStringMap
func NewSafeStringMap() *SafeStringMap {
	return new (SafeStringMap)
}


// NewSafeStringMapFrom creates a new SafeStringMap from a
// StringMapI interface object
func NewSafeStringMapFrom(i StringMapI) *SafeStringMap {
	m := new(SafeStringMap)
	m.Merge(i)
	return m
}


// Clear resets the map to an empty map
func (o *SafeStringMap) Clear() {
	o.Lock()
	defer o.Unlock()
	o.items = nil
}

// SetChanged sets the key to the value and returns a boolean indicating whether doing this caused
// the map to change.
func (o *SafeStringMap) SetChanged(key string, val string) (changed bool) {
	var ok bool
	var oldVal string

	if o == nil {
		panic("The map must be created before being used.")
	}

	o.Lock()
	defer o.Unlock()

	if o.items == nil {
	    o.items = make(map[string]string)
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
		o.items[key] = val
		changed = true
	}
	return
}


// Set sets the key to the given value
func (o *SafeStringMap) Set(key string, val string) {
	if o == nil {
		panic("The map must be created before being used.")
	}
	o.Lock()
	defer o.Unlock()

	if o.items == nil {
	    o.items = make(map[string]string)
	}

    o.items[key] = val
}

// Remove deletes the key from the map. If the key does not exist, nothing happens.
func (o *SafeStringMap) Remove(key string) {
    if o == nil || o.items == nil {
		return
	}
	o.Lock()
	delete(o.items, key)
	o.Unlock()
}

// Get returns the string based on its key. If it does not exist, will return a zero string.
func (o *SafeStringMap) Get(key string) (val string) {
    if o == nil || o.items == nil {
		return
	}
	o.RLock()
	defer o.RUnlock()
	val, _ = o.items[key]
	return
}

// Has returns true if the given key exists in the map.
func (o *SafeStringMap) Has(key string) (ok bool) {
    if o == nil {
        return
    }
	o.RLock()
	defer o.RUnlock()
	if o.items == nil {
		return false
	}

	_, ok = o.items[key]
	return
}

// Values returns a slice of the values. It will return an empty slice if the map is empty.
func (o *SafeStringMap) Values() []string {
	o.Lock()
	defer o.Unlock()

	if o.items == nil {
	    return make([]string, 0)
	}

	vals := make([]string, 0, len(o.items))

	for _, v := range o.items {
		vals = append(vals, v)
	}
	return vals
}

// Keys returns a slice of the keys. It will return an empty slice if the map is empty.
func (o *SafeStringMap) Keys() []string {
	o.Lock()
	defer o.Unlock()

    if o.items == nil {
        return make([]string, 0)
    }

	vals := make([]string, 0, len(o.items))

	for k := range o.items {
		vals = append(vals, k)
	}
	return vals
}

// Len returns the number of items in the map
func (o *SafeStringMap) Len() int {
	return len(o.items)
}

// Range will call the given function with every key and value in the SafeStringMap
// During this process, the map will be locked, so do not pass a function that will be taking significant amounts of time
// If f returns false, it stops the iteration. This pattern is taken from the golang sync.Map.Range function.
func (o *SafeStringMap) Range(f func(key string, value string) bool) {
	if o == nil || o.items == nil {
		return
	}

	o.RLock()
	defer o.RUnlock()

	for k, v := range o.items {
		if !f(k, v) {
			break
		}
	}
}

// MarshalBinary implements the BinaryMarshaler interface to convert a SafeStringMap
// to a byte stream.
func (o *SafeStringMap) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

	o.RLock()
	defer o.RUnlock()

	enc := gob.NewEncoder(&b)
	err := enc.Encode(o.items)
	return b.Bytes(), err
}


// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// SafeStringMap
func (o *SafeStringMap) UnmarshalBinary(data []byte) error {
	o.Lock()
	defer o.Unlock()

	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	err := dec.Decode(&o.items)
	return err
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (o *SafeStringMap) MarshalJSON() (out []byte, err error) {
	o.Lock()
	defer o.Unlock()

    out,err = json.Marshal(o.items)
    return
}

// UnmarshalJSON implements the json.Unmarshaller interface to convert a json object to a SafeStringMap.
// The JSON must start with an object.
func (o *SafeStringMap) UnmarshalJSON(in []byte) (err error) {
    var v map[string]string
    err = json.Unmarshal(in, &v)
    return
}

func (o *SafeStringMap) Copy() StringMapI {
	cp := NewSafeStringMap()

	o.Range(func(key string, value string) bool {

		cp.Set(key, value)
		return true
	})
	return cp
}

// Merge merges the given string map with the current one. The given one takes precedent on collisions.
func (o *SafeStringMap) Merge(i StringMapI) {
	if i == nil {
		return
	}
	o.Lock()
	defer o.Unlock()

    if o.items == nil {
        o.items = make(map[string]string, i.Len())
    }

	i.Range(func(k string, v string) bool {
		o.items[k] = v
		return true
	})
}

// Equals returns true if all the keys in the given map exist in this map, and the values are the same
func (o *SafeStringMap) Equals(i StringMapI) bool {
	if i.Len() != o.Len() {
		return false
	}
	var ret = true

	o.Lock()
	defer o.Unlock()

	i.Range(func(k string, v string) bool {
		if v2,ok := o.items[k]; !ok || v2 != v {
			ret = false
			return false // stop iterating
		}
		return true
	})

	return ret
}


func init() {
	gob.Register(new (SafeStringMap))
}
