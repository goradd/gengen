package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"


)

// SafeStringMap maps a string to a string.
// This version is  safe for concurrent use.
// A zero value is ready for use, but you may not copy it after first using it.
type SafeStringMap struct {
	sync.RWMutex
    items map[string]string
}

// NewSafeStringMap creates a new map that maps string's to string's.
func NewSafeStringMap() *SafeStringMap {
	return new(SafeStringMap)
}

// NewSafeStringMapFrom creates a new SafeStringMap from a
// StringMapI interface object
func NewSafeStringMapFrom(i StringMapI) *SafeStringMap {
	m := NewSafeStringMap()
	m.Merge(i)
	return m
}

// NewSafeStringMapFromMap creates a new SafeStringMap from a
// GO map[string]string object. Note that this will pass control of the given map to the
// new object. After you do this, DO NOT change the original map.
func NewSafeStringMapFromMap(i map[string]string) *SafeStringMap {
	m := NewSafeStringMap()
	m.items = i
	return m
}

// Clear resets the map to an empty map
func (o *SafeStringMap) Clear() {
    if o == nil {
		return
	}
 	o.Lock()
	o.items = nil
    o.Unlock()
}

// SetChanged sets the key to the value and returns a boolean indicating whether doing this caused
// the map to change. It will return true if the key did not first exist, or if the value associated
// with the key was different than the new value.
func (o *SafeStringMap) SetChanged(key string, val string) (changed bool) {
	var ok bool
	var oldVal string

	if o == nil {
		panic("The map must be created before being used.")
	}
 	o.Lock()
	if o.items == nil {
	    o.items = make(map[string]string)
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
		o.items[key] = val
		changed = true
	}
    o.Unlock()
	return
}

// Set sets the key to the given value
func (o *SafeStringMap) Set(key string, val string) {
	if o == nil {
		panic("The map must be initialized before being used.")
	}
 	o.Lock()
    if o.items == nil {
        o.items = make(map[string]string)
    }

	o.items[key] = val
    o.Unlock()
}

// Get returns the value based on its key. If it does not exist, an empty string will be returned.
func (o *SafeStringMap) Get(key string) (val string) {
    val,_ = o.Load(key)
	return
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (o *SafeStringMap) Load(key string) (val string, ok bool) {
    if o == nil {
		return
	}
    o.RLock()
	if o.items != nil {
	    val,ok = o.items[key]
	}
    o.RUnlock()
	return
}




// Delete removes the key from the map. If the key does not exist, nothing happens.
func (o *SafeStringMap) Delete(key string) {
    if o == nil {
		return
	}
 	o.Lock()
 	if o.items != nil {
	    delete(o.items, key)
	}
    o.Unlock()
}


// Has returns true if the given key exists in the map.
func (o *SafeStringMap) Has(key string) (exists bool) {
    if o == nil {
		return
	}
    o.RLock()
    if o.items != nil {
 	    _, exists = o.items[key]
    }
    o.RUnlock()
	return
}

// Values returns a slice of the values. It will return a nil slice if the map is empty.
// Multiple calls to Values will result in the same list of values, but may be in a different order.
func (o *SafeStringMap) Values() (vals []string) {
    if o == nil {
        return
    }
    o.RLock()
    if len(o.items) > 0 {
        vals = make([]string, len(o.items))

        var i int
        for _, v := range o.items {
            vals[i] = v
            i++
        }
    }
    o.RUnlock()

	return
}

// Keys returns a slice of the keys. It will return a nil slice if the map is empty.
// Multiple calls to Keys will result in the same list of keys, but may be in a different order.
func (o *SafeStringMap) Keys() (keys []string) {
    if o == nil {
        return nil
    }
    o.RLock()
    if len(o.items) > 0 {
        keys = make([]string, len(o.items))

        var i int
        for k := range o.items {
            keys[i] = k
            i++
        }
    }
    o.RUnlock()
	return
}

// Len returns the number of items in the map
func (o *SafeStringMap) Len() (l int) {
    if o == nil {
		return
	}
    o.RLock()
    l = len(o.items)
    o.RUnlock()
	return
}

// Range will call the given function with every key and value in the map.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
// During this process, the map will be locked, so do not pass a function that will take significant amounts of time.
func (o *SafeStringMap) Range(f func(key string, value string) bool) {
	if o == nil {
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

// Merge merges the given  map with the current one. The given one takes precedent on collisions.
func (o *SafeStringMap) Merge(i StringMapI) {
	if i == nil {
		return
	}

	if o == nil {
		panic("The map must be created before being used.")
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
    len := o.Len()
	if i.Len() != len {
		return false
	} else if len == 0 { // both are zero
	    return true
	}
	var ret = true
    o.RLock()
    defer o.RUnlock()

	i.Range(func(k string, v string) bool {
		if v2,ok := o.items[k]; !ok || v2 != v {
			ret = false
			return false // stop iterating
		}
		return true
	})

	return ret
}

// Copy will make a copy of the map and a copy of the underlying data.
func (o *SafeStringMap) Copy() StringMapI {
	cp := NewSafeStringMap()

	o.Range(func(key string, value string) bool {




		cp.Set(key, value)
		return true
	})
	return cp
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (o *SafeStringMap) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

 	enc := gob.NewEncoder(&b)
    o.RLock()
    defer o.RUnlock()
	err := enc.Encode(o.items)
	return b.Bytes(), err
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// SafeStringMap
func (o *SafeStringMap) UnmarshalBinary(data []byte) (err error) {
    var v map[string]string

	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	if err = dec.Decode(&v); err == nil {
        o.Lock()
        o.items = v
        o.Unlock()
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (o *SafeStringMap) MarshalJSON() (out []byte, err error) {
    o.RLock()
    defer o.RUnlock()
    out,err = json.Marshal(o.items)
    return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a SafeStringMap.
// The JSON must start with an object.
func (o *SafeStringMap) UnmarshalJSON(in []byte) (err error) {
    var v map[string]string
    if err = json.Unmarshal(in, &v); err == nil {
        o.Lock()
        o.items = v
        o.Unlock()
    }
    return
}

func (o *SafeStringMap) IsNil() bool {
	return o == nil
}

func (o *SafeStringMap) String() string {
	var s string

    // sort on keys to stabilize order
	keys := o.Keys()
    sort.Slice(keys, func(a,b int) bool {
        return keys[a] < keys[b]
    })

	s = "{"
	for _,k := range keys {
	    v := o.Get(k)
	    s += fmt.Sprintf(`%#v:%#v,`, k, v)
	}
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}


func init() {
	gob.Register(new (SafeStringMap))
}
