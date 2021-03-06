package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
	"strings"
	"fmt"
    "sync"
)

// A SafeStringSliceMap combines a map with a slice so that you can range over a
// map in a predictable order. By default, the order will be the same order that items were inserted,
// i.e. a FIFO list. This is similar to how PHP arrays work.
// SafeStringSliceMap implements the sort interface so you can change the order
// before ranging over the values if desired.
// It is  safe for concurrent use.
// The zero of this is usable immediately.
// The SafeStringSliceMap satisfies the StringMapI interface.
type SafeStringSliceMap struct {
    sync.RWMutex
	items map[string]string
	order []string
	lessF func(key1,key2 string, val1, val2 string) bool
}

// NewSafeStringSliceMap creates a new map that maps string's to string's.
func NewSafeStringSliceMap() *SafeStringSliceMap {
	return new (SafeStringSliceMap)
}

// NewSafeStringSliceMapFrom creates a new SafeStringMap from a
// StringMapI interface object
func NewSafeStringSliceMapFrom(i StringMapI) *SafeStringSliceMap {
	m := new (SafeStringSliceMap)
	m.Merge(i)
	return m
}

// NewSafeStringSliceMapFromMap creates a new SafeStringSliceMap from a
// GO map[string]string object. Note that this will pass control of the given map to the
// new object. After you do this, DO NOT change the original map.
func NewSafeStringSliceMapFromMap(i map[string]string) *SafeStringSliceMap {
	m := NewSafeStringSliceMap()
	m.items = i
	m.order = make([]string, len(m.items), len(m.items))
	j := 0
	for k := range m.items {
	    m.order[j] = k
	    j++
	}
	return m
}

// SetSortFunc sets the sort function which will determine the order of the items in the map
// on an ongoing basis. Normally, items will iterate in the order they were added.
// The sort function is a Less function, that returns true when item 1 is "less" than item 2.
// The sort function receives both the keys and values, so it can use either to decide how to sort.
func (o *SafeStringSliceMap) SetSortFunc(f func(key1,key2 string, val1, val2 string) bool) *SafeStringSliceMap {
        o.Lock()
    o.lessF = f
    if f != nil && len(o.order) > 0 {
        sort.Slice(o.order, func(i,j int) bool {
            return f(o.order[i], o.order[j], o.items[o.order[i]], o.items[o.order[j]])
        })
    }
        o.Unlock()

    return o
}

// SortByKeys sets up the map to have its sort order sort by keys, lowest to highest
func (o *SafeStringSliceMap) SortByKeys() *SafeStringSliceMap {
    o.SetSortFunc(keySortSafeStringSliceMap)
    return o
}

func keySortSafeStringSliceMap(key1, key2 string, val1, val2 string) bool {
    return key1 < key2
}


// SortByValues sets up the map to have its sort order sort by values, lowest to highest
func (o *SafeStringSliceMap) SortByValues() {
    o.SetSortFunc(valueSortSafeStringSliceMap)
}

func valueSortSafeStringSliceMap(key1, key2 string, val1, val2 string) bool {
    return val1 < val2
}



// SetChanged sets the value.
// It returns true if something in the map changed. If the key
// was already in the map, and you have not provided a sort function,
// the order will not change, but the value will be replaced. If you wanted the
// order to change, you must Delete then call SetChanged. If you have previously set a sort function,
// the order will be updated.
func (o *SafeStringSliceMap) SetChanged(key string, val string) (changed bool) {
	var ok bool
	var oldVal string

	if o == nil {
	    panic("You must initialize the map before using it.")
	}
    o.Lock()

	if o.items == nil {
	    o.items = make(map[string]string)
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
        if o.lessF != nil {
            if ok {
                // delete old key location
                loc := sort.Search (len(o.items), func(n int) bool {
                    return !o.lessF(o.order[n], key, o.items[o.order[n]], oldVal)
                })
                o.order = append(o.order[:loc], o.order[loc+1:]...)
            }

            loc := sort.Search (len(o.order), func(n int) bool {
                return o.lessF(key, o.order[n], val, o.items[o.order[n]])
            })
            // insert
            o.order = append(o.order, key)
            copy(o.order[loc+1:], o.order[loc:])
            o.order[loc] = key
        } else {
		    if !ok {
			    o.order = append(o.order, key)
		    }
		}
		o.items[key] = val
		changed = true
	}
    o.Unlock()

	return
}


// Set sets the given key to the given value.
// If the key already exists, the range order will not change.
func (o *SafeStringSliceMap) Set(key string, val string) {
	var ok bool
	var oldVal string

	if o == nil {
	    panic("You must initialize the map before using it.")
	}
    o.Lock()

	if o.items == nil {
	    o.items = make(map[string]string)
	}

	_, ok = o.items[key]
    if o.lessF != nil {
        if ok {
            // delete old key location
            loc := sort.Search (len(o.items), func(n int) bool {
                return !o.lessF(o.order[n], key, o.items[o.order[n]], oldVal)
            })
            o.order = append(o.order[:loc], o.order[loc+1:]...)
        }

        loc := sort.Search (len(o.order), func(n int) bool {
            return o.lessF(key, o.order[n], val, o.items[o.order[n]])
        })
        // insert
        o.order = append(o.order, key)
        copy(o.order[loc+1:], o.order[loc:])
        o.order[loc] = key
    } else {
        if !ok {
            o.order = append(o.order, key)
        }
    }
    o.items[key] = val
    o.Unlock()

	return
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, it puts it at the end. Negative indexes are backwards from the end.
func (o *SafeStringSliceMap) SetAt(index int, key string, val string)  {
    if o == nil {
        panic("You must initialize the map before using it.")
    }

    if o.lessF != nil {
        panic("You cannot use SetAt if you are also using a sort function.")
    }

	if index >= len(o.order) {
		o.Set(key, val)
		return
	}

	var ok bool
	var emptyKey string
    o.Lock()

	if _, ok = o.items[key]; !ok {
		if index <= -len(o.items) {
			index = 0
		}
		if index < 0 {
			index = len(o.items) + index
		}

		o.order = append(o.order, emptyKey)
		copy(o.order[index+1:], o.order[index:])
		o.order[index] = key
	}
	o.items[key] = val
    o.Unlock()
    return
}

// Delete removes the item with the given key.
func (o *SafeStringSliceMap) Delete(key string) {
    if o == nil {
        return
    }
    o.Lock()

    if _,ok := o.items[key]; ok {
        if o.lessF != nil {
            oldVal := o.items[key]
            loc := sort.Search (len(o.items), func(n int) bool {
                return !o.lessF(o.order[n], key, o.items[o.order[n]], oldVal)
            })
            o.order = append(o.order[:loc], o.order[loc+1:]...)
        } else {
            for i, v := range o.order {
                if v == key {
                    o.order = append(o.order[:i], o.order[i+1:]...)
                    break
                }
            }
        }
        delete(o.items, key)
    }
    o.Unlock()
}

// Get returns the value based on its key. If the key does not exist, an empty value is returned.
func (o *SafeStringSliceMap) Get(key string) (val string) {
    val,_ = o.Load(key)
    return
}

// Load returns the value based on its key, and a boolean indicating whether it exists in the map.
// This is the same interface as sync.Map.Load()
func (o *SafeStringSliceMap) Load(key string) (val string, ok bool) {
    if o == nil {
        return
    }
    o.RLock()
    if o.items != nil {
    	val, ok = o.items[key]
    }
    o.RUnlock()
	return
}



// Has returns true if the given key exists in the map.
func (o *SafeStringSliceMap) Has(key string) (ok bool) {
    if o == nil {
        return false
    }
    o.RLock()
    if o.items != nil {
	    _, ok = o.items[key]
    }
    o.RUnlock()
	return
}


// Is returns true if the given key exists in the map and has the given value.
func (o *SafeStringSliceMap) Is(key string, val string) (is bool) {
    if o == nil {
		return
	}

    var v string
    o.RLock()
    if o.items != nil {
 	    v, is = o.items[key]
    }
    o.RUnlock()
	return is && v == val
}


// GetAt returns the value based on its position. If the position is out of bounds, an empty value is returned.
func (o *SafeStringSliceMap) GetAt(position int) (val string) {
    if o == nil {
        return
    }
    o.RLock()
	if position < len(o.order) && position >= 0 {
		val, _ = o.items[o.order[position]]
	}
    o.RUnlock()
	return
}

// GetKeyAt returns the key based on its position. If the position is out of bounds, an empty value is returned.
func (o *SafeStringSliceMap) GetKeyAt(position int) (key string) {
    if o == nil {
        return
    }
    o.RLock()
	if position < len(o.order) && position >= 0 {
		key = o.order[position]
	}
    o.RUnlock()
	return
}

// Values returns a slice of the values in the order they were added or sorted.
func (o *SafeStringSliceMap) Values() (vals []string) {
    if o == nil {
        return
    }
    o.RLock()

    if o.items != nil {
  	    vals = make([]string, len(o.order))
        for i, v := range o.order {
            vals[i] = o.items[v]
        }
    }
    o.RUnlock()

	return
}

// Keys returns the keys of the map, in the order they were added or sorted
func (o *SafeStringSliceMap) Keys() (keys []string) {
    if o == nil {
        return
    }
    o.RLock()

    if len(o.order) != 0 {
 	    keys = make([]string, len(o.order))
        for i, v := range o.order {
            keys[i] = v
        }
    }
    o.RUnlock()

	return
}

// Len returns the number of items in the map
func (o *SafeStringSliceMap) Len() int {
    if o == nil {
        return 0
    }
    o.RLock()
    l := len(o.order)
    o.RUnlock()
	return l
}


// Copy will make a copy of the map and a copy of the underlying data.
func (o *SafeStringSliceMap) Copy() *SafeStringSliceMap {
	cp := NewSafeStringSliceMap()

	o.Range(func(key string, value string) bool {

		cp.Set(key, value)
		return true
	})
	cp.lessF = o.lessF
	return cp
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
// If you are using a sort function, you must save and restore the sort function in a separate operation
// since functions are not serializable.
func (o *SafeStringSliceMap) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)
	o.RLock()
	defer o.RUnlock()

	err = encoder.Encode(o.items)
	if err == nil {
		err = encoder.Encode(o.order)
	}
	data = buf.Bytes()
	return
}

// UnmarshalBinary implements the BinaryUnmarshaler interface to convert a byte stream to a
// SafeStringSliceMap
func (o *SafeStringSliceMap) UnmarshalBinary(data []byte) (err error) {
    var items map[string]string
	var order []string

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	if err = dec.Decode(&items); err == nil {
		err = dec.Decode(&order)
	}

	if err == nil {
        o.Lock()
        o.items = items
        o.order = order
        o.Unlock()
	}
	return err
}

// MarshalJSON implements the json.Marshaler interface to convert the map into a JSON object.
func (o *SafeStringSliceMap) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
    o.RLock()
    defer o.RUnlock()
	data, err = json.Marshal(o.items)
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a SafeStringMap.
// The JSON must start with an object.
func (o *SafeStringSliceMap) UnmarshalJSON(data []byte) (err error) {
    var items map[string]string

	if err = json.Unmarshal(data, &items); err == nil {
        o.Lock()
        o.items = items
        // Create a default order, since these are inherently unordered
        o.order = make([]string, len(o.items))
        i := 0
        for k := range o.items {
            o.order[i] = k
            i++
        }
        o.Unlock()
	}
	return
}


// Merge the given map into the current one
func (o *SafeStringSliceMap) Merge(i StringMapI) {
	if i != nil {
		i.Range(func(k string, v string) bool {
			o.Set(k, v)
			return true
		})
	}
}

// MergeMap merges the given standard map with the current one. The given one takes precedent on collisions.
func (o *SafeStringSliceMap) MergeMap(m map[string]string) {
	if m == nil {
		return
	}

	for k,v := range m {
		o.Set(k, v)
	}
}


// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *SafeStringSliceMap) Range(f func(key string, value string) bool) {
	if o == nil {
		return
	}
    o.Lock()
    defer o.Unlock()
	if  o.items != nil {
        for _, k := range o.order {
            if !f(k, o.items[k]) {
                break
            }
        }
    }
}

// Equals returns true if the map equals the given map, paying attention only to the content of the
// map and not the order.
func (o *SafeStringSliceMap) Equals(i StringMapI) bool {
	l := i.Len()
	if l == 0 {
		return o == nil
	}
	if l != o.Len() {
		return false
	}
	var ret = true

	o.Range(func(k string, v string) bool {
		if !i.Has(k) || v != i.Get(k) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

func (o *SafeStringSliceMap) Clear() {
    if o == nil {return}
    o.Lock()
	o.items = nil
	o.order = nil
    o.Unlock()

}

func (o *SafeStringSliceMap) IsNil() bool {
	return o == nil
}

func (o *SafeStringSliceMap) String() string {
	var s string

	s = "{"
	o.Range(func(k string, v string) bool {
		s += fmt.Sprintf(`%#v:%#v,`, k, v)
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}


// Join is just like strings.Join
func (o *SafeStringSliceMap) Join(glue string) string {
	return strings.Join(o.Values(), glue)
}


func init() {
	gob.Register(new (SafeStringSliceMap))
}