package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

)

// Map maps a string to a interface{}.
// This version is not safe for concurrent use.
// A zero value is ready for use, but you may not copy it after first using it.
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

// NewMapFromMap creates a new Map from a
// GO map[string]interface{} object. Note that this will pass control of the given map to the
// new object. After you do this, DO NOT change the original map.
func NewMapFromMap(i map[string]interface{}) *Map {
	m := NewMap()
	m.items = i
	return m
}

// Clear resets the map to an empty map
func (o *Map) Clear() {
    if o == nil {
		return
	}
	o.items = nil
}

// SetChanged sets the key to the value and returns a boolean indicating whether doing this caused
// the map to change. It will return true if the key did not first exist, or if the value associated
// with the key was different than the new value.
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

// Set sets the key to the given value
func (o *Map) Set(key string, val interface{}) {
	if o == nil {
		panic("The map must be initialized before being used.")
	}
    if o.items == nil {
        o.items = make(map[string]interface{})
    }

	o.items[key] = val
}

// Get returns the value based on its key. If it does not exist, an empty string will be returned.
func (o *Map) Get(key string) (val interface{}) {
    val,_ = o.Load(key)
	return
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (o *Map) Load(key string) (val interface{}, ok bool) {
    if o == nil {
		return
	}
	if o.items != nil {
	    val,ok = o.items[key]
	}
	return
}


func (o *Map) LoadString(key string) (val string, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(string)
    }
    return
}

func (o *Map) LoadInt(key string) (val int, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(int)
    }
    return
}

func (o *Map) LoadBool(key string) (val bool, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(bool)
    }
    return
}

func (o *Map) LoadFloat64(key string) (val float64, ok bool) {
    var v interface{}
    v,ok = o.Load(key)
    if ok {
        val,ok = v.(float64)
    }
    return
}



// Delete removes the key from the map. If the key does not exist, nothing happens.
func (o *Map) Delete(key string) {
    if o == nil {
		return
	}
 	if o.items != nil {
	    delete(o.items, key)
	}
}


// Has returns true if the given key exists in the map.
func (o *Map) Has(key string) (exists bool) {
    if o == nil {
		return
	}
    if o.items != nil {
 	    _, exists = o.items[key]
    }
	return
}

// Values returns a slice of the values. It will return a nil slice if the map is empty.
// Multiple calls to Values will result in the same list of values, but may be in a different order.
func (o *Map) Values() (vals []interface{}) {
    if o == nil {
        return
    }
    if len(o.items) > 0 {
        vals = make([]interface{}, len(o.items))

        var i int
        for _, v := range o.items {
            vals[i] = v
            i++
        }
    }

	return
}

// Keys returns a slice of the keys. It will return a nil slice if the map is empty.
// Multiple calls to Keys will result in the same list of keys, but may be in a different order.
func (o *Map) Keys() (keys []string) {
    if o == nil {
        return nil
    }
    if len(o.items) > 0 {
        keys = make([]string, len(o.items))

        var i int
        for k := range o.items {
            keys[i] = k
            i++
        }
    }
	return
}

// Len returns the number of items in the map
func (o *Map) Len() (l int) {
    if o == nil {
		return
	}
    l = len(o.items)
	return
}

// Range will call the given function with every key and value in the map.
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

// Merge merges the given  map with the current one. The given one takes precedent on collisions.
func (o *Map) Merge(i MapI) {
	if i == nil {
		return
	}

	if o == nil {
		panic("The map must be created before being used.")
	}

	if o.items == nil {
	    o.items = make(map[string]interface{}, i.Len())
	}
	i.Range(func(k string, v interface{}) bool {
		o.items[k] = v
		return true
	})
}

// Equals returns true if all the keys in the given map exist in this map, and the values are the same
func (o *Map) Equals(i MapI) bool {
    len := o.Len()
	if i.Len() != len {
		return false
	} else if len == 0 { // both are zero
	    return true
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

// Copy will make a copy of the map and a copy of the underlying data.
func (o *Map) Copy() MapI {
	cp := NewMap()

	o.Range(func(key string, value interface{}) bool {




		cp.Set(key, value)
		return true
	})
	return cp
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (o *Map) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer

 	enc := gob.NewEncoder(&b)
	err := enc.Encode(o.items)
	return b.Bytes(), err
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// Map
func (o *Map) UnmarshalBinary(data []byte) (err error) {
    var v map[string]interface{}

	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	if err = dec.Decode(&v); err == nil {
        o.items = v
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (o *Map) MarshalJSON() (out []byte, err error) {
    out,err = json.Marshal(o.items)
    return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a Map.
// The JSON must start with an object.
func (o *Map) UnmarshalJSON(in []byte) (err error) {
    var v map[string]interface{}
    if err = json.Unmarshal(in, &v); err == nil {
        o.items = v
    }
    return
}

func (o *Map) IsNil() bool {
	return o == nil
}

func (o *Map) String() string {
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
	gob.Register(new (Map))
}
