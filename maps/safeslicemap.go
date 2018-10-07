package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
	"strings"
    "sync"

	"fmt"

)

// A SafeSliceMap combines a map with a slice so that you can range over a
// map in a predictable order. By default, the order will be the same order that items were inserted,
// i.e. a FIFO list. This is similar to how PHP arrays work.
// SafeSliceMap implements the sort interface so you can change the order
// before ranging over the values if desired.
// It is  safe for concurrent use.
// The zero of this is usable immediately.
// The SafeSliceMap satisfies the MapI interface.
type SafeSliceMap struct {
    sync.RWMutex
	items map[string]interface{}
	order []string
}

func NewSafeSliceMap() *SafeSliceMap {
	return new (SafeSliceMap)
}

func NewSafeSliceMapFrom(i MapI) *SafeSliceMap {
	m := new (SafeSliceMap)
	m.Merge(i)
	return m
}


// SetChanged sets the value, but also appends the value to the end of the list.
// It returns true if something in the map changed. If the key
// was already in the map, the order will not change, but the value will be replaced. If you want the
// order to change, you must Delete then call SetChanged.
func (o *SafeSliceMap) SetChanged(key string, val interface{}) (changed bool) {
	var ok bool
	var oldVal interface{}

	if o == nil {
	    panic("You must initialize the map before using it.")
	}
    o.Lock()

	if o.items == nil {
	    o.items = make(map[string]interface{})
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
		if !ok {
			o.order = append(o.order, key)
		}
		o.items[key] = val
		changed = true
	}
    o.Unlock()

	return
}

// Set sets the given key to the given value.
// If the key already exists, the range order will not change.
func (o *SafeSliceMap) Set(key string, val interface{}) {
	o.SetChanged(key, val)
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, or -1, it is the same as Set, in that it puts it at the end. Negative indexes are backwards from the
// end, if smaller than the negative length, just inserts at the beginning.
func (o *SafeSliceMap) SetAt(index int, key string, val interface{})  {
    if o == nil {
        panic("You must initialize the map before using it.")
    }

	if index == -1 || index >= len(o.order) {
		o.Set(key, val)
		return
	}

	var ok bool
	var emptyKey string
    o.Lock()

	if _, ok = o.items[key]; !ok {
		if index < -len(o.items) {
			index = 0
		}
		if index < 0 {
			index = len(o.items) + index + 1
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
func (o *SafeSliceMap) Delete(key string) {
    if o == nil {
        return
    }
    o.Lock()
	for i, v := range o.order {
		if v == key {
			o.order = append(o.order[:i], o.order[i+1:]...)
			break
		}
	}
	delete(o.items, key)
    o.Unlock()
}

// Get returns the value based on its key. If the key does not exist, it will return an empty value.
func (o *SafeSliceMap) Get(key string) (val interface{}) {
    if o == nil {
        return
    }
    o.RLock()
    if o.items != nil {
    	val, _ = o.items[key]
    }
    o.RUnlock()
	return
}

// Has returns true if the given key exists in the map.
func (o *SafeSliceMap) Has(key string) (ok bool) {
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

// GetAt returns the value based on its position. If the position is out of bounds, an empty value is returned.
func (o *SafeSliceMap) GetAt(position int) (val interface{}) {
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


// Values returns a slice of the values in the order they were added or sorted.
func (o *SafeSliceMap) Values() (vals []interface{}) {
    if o == nil {
        return
    }
    o.RLock()

    if o.items != nil {
  	    vals = make([]interface{}, len(o.order))
        for i, v := range o.order {
            vals[i] = o.items[v]
        }
    }
    o.RUnlock()

	return
}

// Keys returns the keys of the map, in the order they were added or sorted
func (o *SafeSliceMap) Keys() (keys []string) {
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
func (o *SafeSliceMap) Len() int {
    if o == nil {
        return 0
    }
    o.RLock()
    l := len(o.order)
    o.RUnlock()
	return l
}

// Less is part of the interface that allows the map to be sorted by values.
// It returns true if the value at position i should be sorted before the value at position j.
func (o *SafeSliceMap) Less(i, j int) bool {
	o.RLock()
	defer o.RUnlock()


    switch v := o.items[o.order[i]].(type) {
    case Comparer:
        return v.Compare(o.items[o.order[j]]) < 0
    default:
    	panic ("Values are not sortable")
    	return false
    }

}

// Swap is part of the interface that allows the slice to be sorted. It swaps the positions
// of the items and position i and j.
func (o *SafeSliceMap) Swap(i, j int) {
	o.Lock()
	o.order[i], o.order[j] = o.order[j], o.order[i]
    o.Unlock()
}

// Sort by keys interface
type sortSafebykeys struct {
	// This embedded interface permits Reverse to use the methods of
	// another interface implementation.
	sort.Interface
}

// A helper function to allow SafeSliceMaps to be sorted by keys
// To sort the map by keys, call:
//   sort.Sort(OrderSafeStringSliceMapByKeys(m))
func OrderSafeSliceMapByKeys(o *SafeSliceMap) sort.Interface {
	return &sortSafebykeys{o}
}

// A helper function to allow SafeSliceMaps to be sorted by keys
func (r sortSafebykeys) Less(i, j int) bool {
	var o *SafeSliceMap = r.Interface.(*SafeSliceMap)
	o.RLock()
	defer o.RUnlock()


	return o.order[i] < o.order[j]

}

// Copy will make a copy of the map and a copy of the underlying data.
// If the values implement the Copier interface, the value's Copy function will be called to deep copy the items.
func (o *SafeSliceMap) Copy() *SafeSliceMap {
	cp := NewSafeSliceMap()

	o.Range(func(key string, value interface{}) bool {
		if copier, ok := value.(Copier); ok {
			value = copier.Copy()
		}

		cp.Set(key, value)
		return true
	})
	return cp
}

// MarshalBinary implements the BinaryMarshaler interface to convert the map to a byte stream.
func (o *SafeSliceMap) MarshalBinary() (data []byte, err error) {
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
// SafeSliceMap
func (o *SafeSliceMap) UnmarshalBinary(data []byte) (err error) {
    var items map[string]interface{}
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
func (o *SafeSliceMap) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
    o.RLock()
    defer o.RUnlock()
	data, err = json.Marshal(o.items)
	return
}

// UnmarshalJSON implements the json.Unmarshaler interface to convert a json object to a SafeMap.
// The JSON must start with an object.
func (o *SafeSliceMap) UnmarshalJSON(data []byte) (err error) {
    var items map[string]interface{}

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
func (o *SafeSliceMap) Merge(i MapI) {
	if i != nil {
		i.Range(func(k string, v interface{}) bool {
			o.Set(k, v)
			return true
		})
	}
}

// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *SafeSliceMap) Range(f func(key string, value interface{}) bool) {
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
func (o *SafeSliceMap) Equals(i MapI) bool {
	l := i.Len()
	if l == 0 {
		return o == nil
	}
	if l != o.Len() {
		return false
	}
	var ret = true

	o.Range(func(k string, v interface{}) bool {
		if !i.Has(k) || v != i.Get(k) {
			ret = false
			return false
		}
		return true
	})
	return ret
}

func (o *SafeSliceMap) Clear() {
    if o == nil {return}
    o.Lock()
	o.items = nil
	o.order = nil
    o.Unlock()

}

func (o *SafeSliceMap) IsNil() bool {
	return o == nil
}

func (o *SafeSliceMap) String() string {
	var s string

	s = "{"
	o.Range(func(k string, v interface{}) bool {
		s += `"` + k + `":"` +
fmt.Sprintf("%v", v)+ `",`
		return true
	})
	s = strings.TrimRight(s, ",")
	s += "}"
	return s
}



func init() {
	gob.Register(new (SafeSliceMap))
}