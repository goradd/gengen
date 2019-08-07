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

// SafeMap maps a string to a interface{}.
// This version is  safe for concurrent use.
// A zero value is ready for use, but you may not copy it after first using it.
type SafeMap struct {
	sync.RWMutex
    items map[string]interface{}
}

// NewSafeMap creates a new map that maps string's to interface{}'s.
func NewSafeMap() *SafeMap {
	return new(SafeMap)
}

// NewSafeMapFrom creates a new SafeMap from a
// MapI interface object
func NewSafeMapFrom(i MapI) *SafeMap {
	m := NewSafeMap()
	m.Merge(i)
	return m
}

// NewSafeMapFromMap creates a new SafeMap from a
// GO map[string]interface{} object. Note that this will pass control of the given map to the
// new object. After you do this, DO NOT change the original map.
func NewSafeMapFromMap(i map[string]interface{}) *SafeMap {
	m := NewSafeMap()
	m.items = i
	return m
}

// Clear resets the map to an empty map
func (o *SafeMap) Clear() {
    if o == nil {
		return
	}
 	o.Lock()
	o.items = nil
    o.Unlock()
}



// Set sets the key to the given value
func (o *SafeMap) Set(key string, val interface{}) {
	if o == nil {
		panic("The map must be initialized before being used.")
	}
 	o.Lock()
    if o.items == nil {
        o.items = make(map[string]interface{})
    }

	o.items[key] = val
    o.Unlock()
}

// Get returns the value based on its key. If it does not exist, an empty string will be returned.
func (o *SafeMap) Get(key string) (val interface{}) {
    val,_ = o.Load(key)
	return
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (o *SafeMap) Load(key string) (val interface{}, ok bool) {
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


func (o *SafeMap) LoadString(key string) (val string, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(string)
    }
    return
}

func (o *SafeMap) LoadInt(key string) (val int, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(int)
    }
    return
}

func (o *SafeMap) LoadBool(key string) (val bool, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(bool)
    }
    return
}

func (o *SafeMap) LoadFloat64(key string) (val float64, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(float64)
    }
    return
}



// Delete removes the key from the map. If the key does not exist, nothing happens.
func (o *SafeMap) Delete(key string) {
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
func (o *SafeMap) Has(key string) (exists bool) {
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
func (o *SafeMap) Values() (vals []interface{}) {
    if o == nil {
        return
    }
    o.RLock()
    if len(o.items) > 0 {
        vals = make([]interface{}, len(o.items))

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
func (o *SafeMap) Keys() (keys []string) {
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
func (o *SafeMap) Len() (l int) {
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
func (o *SafeMap) Range(f func(key string, value interface{}) bool) {
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
func (o *SafeMap) Merge(i MapI) {
	if i == nil {
		return
	}

	if o == nil {
		panic("The map must be created before being used.")
	}
	o.Lock()
	defer o.Unlock()

	if o.items == nil {
	    o.items = make(map[string]interface{}, i.Len())
	}
	i.Range(func(k string, v interface{}) bool {
		o.items[k] = v
		return true
	})
}

// MergeMap merges the given standard map with the current one. The given one takes precedent on collisions.
func (o *SafeMap) MergeMap(m map[string]interface{}) {
	if m == nil {
		return
	}

	if o == nil {
		panic("The map must be created before being used.")
	}
	o.Lock()
	defer o.Unlock()

	if o.items == nil {
	    o.items = make(map[string]interface{}, len(m))
	}
	for k,v := range m {
		o.items[k] = v
	}
}


// Equals returns true if all the keys in the given map exist in this map, and the values are the same
func (o *SafeMap) Equals(i MapI) bool {
    len := o.Len()
	if i.Len() != len {
		return false
	} else if len == 0 { // both are zero
	    return true
	}
	var ret = true
    o.RLock()
    defer o.RUnlock()

	i.Range(func(k string, v interface{}) bool {
		if v2,ok := o.items[k]; !ok || v2 != v {
			ret = false
			return false // stop iterating
		}
		return true
	})

	return ret
}

// Copy will make a copy of the map and a copy of the underlying data.
func (o *SafeMap) Copy() MapI {
	cp := NewSafeMap()

	o.Range(func(key string, value interface{}) bool {




		cp.Set(key, value)
		return true
	})
	return cp
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (o *SafeMap) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

 	enc := gob.NewEncoder(&b)
    o.RLock()
    defer o.RUnlock()
	err := enc.Encode(o.items)
	return b.Bytes(), err
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// SafeMap
func (o *SafeMap) UnmarshalBinary(data []byte) (err error) {
    var v map[string]interface{}

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
func (o *SafeMap) MarshalJSON() (out []byte, err error) {
    o.RLock()
    defer o.RUnlock()
    out,err = json.Marshal(o.items)
    return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a SafeMap.
// The JSON must start with an object.
func (o *SafeMap) UnmarshalJSON(in []byte) (err error) {
    var v map[string]interface{}
    if err = json.Unmarshal(in, &v); err == nil {
        o.Lock()
        o.items = v
        o.Unlock()
    }
    return
}

func (o *SafeMap) IsNil() bool {
	return o == nil
}

func (o *SafeMap) String() string {
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
	gob.Register(new (SafeMap))
}
