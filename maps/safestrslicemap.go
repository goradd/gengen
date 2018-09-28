package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
	"strings"
	"sync"
)

// A SafeStringSliceMap combines a map with a slice so that you can range over a
// map in a predictable order. By default, the order will be the same order that items were inserted,
// i.e. a FIFO list. This is similar to how PHP arrays work. You can change this order by providing a
// sorting mechanism. However, this object is NOT safe for concurrent use.
// The zero of this is usable immediately.
// The SafeStringSliceMap satisfies the StringMapI interface.
type SafeStringSliceMap struct {
	sync.RWMutex
	items map[string]string
	order []string
}

func NewSafeStringSliceMap() *SafeStringSliceMap {
	return new (SafeStringSliceMap)
}

func NewSafeStringSliceMapFrom(i StringMapI) *SafeStringSliceMap {
	m := new (SafeStringSliceMap)
	m.Merge(i)
	return m
}

func (o *SafeStringSliceMap) Copy() *SafeStringSliceMap {
	return NewSafeStringSliceMapFrom(o)
}

// SetChanged sets the value, but also appends the value to the end of the list for when you
// iterate over the list. Returns whether something changed, and if an error occurred. If the key
// was already in the map, the order will not change, but the value will be replaced. If you want the
// order to change, you must Remove then SetChanged
func (o *SafeStringSliceMap) SetChanged(key string, val string) (changed bool) {
	o.Lock()
	defer o.Unlock()

	var ok bool
	var oldVal string

	if o.items == nil {
	    o.items = make(map[string]string)
	}

	if oldVal, ok = o.items[key]; !ok || oldVal != val {
		if !ok {
			o.order = append(o.order, key)
		}
		o.items[key] = val
		changed = true
	}
	return
}

// Set sets the given key to the given value
// If the key already exists, the range order will not change.
func (o *SafeStringSliceMap) Set(key string, val string) {
	o.SetChanged(key, val)
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, or -1, it is the same as Set, in that it puts it at the end. Negative indexes are backwards from the
// end, if smaller than the negative length, just inserts at the beginning.
func (o *SafeStringSliceMap) SetAt(index int, key string, val string) {
	var l int
	o.RLock()
	l = len(o.items)
	o.RUnlock()

	if index == -1 || index >= l {
		o.Set(key, val)
		return
	}
	
	var ok bool
	var emptyKey string

	o.Lock()
	defer o.Unlock()

	if _, ok = o.items[key]; !ok {
		if index < -l {
			index = 0
		}
		if index < 0 {
			index = l + index + 1
		}

		o.order = append(o.order, emptyKey)
		copy(o.order[index+1:], o.order[index:])
		o.order[index] = key
	}
	o.items[key] = val
	return
}

// Remove removes the item with the given key.
func (o *SafeStringSliceMap) Remove(key string) {
	o.Lock()
	defer o.Unlock()

	for i, v := range o.order {
		if v == key {
			o.order = append(o.order[:i], o.order[i+1:]...)
			break
		}
	}
	delete(o.items, key)
}

// Get returns the value based on its key. If the key does not exist, it will return an empty value.
func (o *SafeStringSliceMap) Get(key string) (val string) {
	o.RLock()
	defer o.RUnlock()

    if o.items != nil {
    	val, _ = o.items[key]
    }
	return
}

// Has returns true if the given key exists in the map.
func (o *SafeStringSliceMap) Has(key string) (ok bool) {
 	o.RLock()
 	defer o.RUnlock()

   if o.items == nil {
        return false
    }
	_, ok = o.items[key]
	return
}

// GetAt returns the value based on its position.
func (o *SafeStringSliceMap) GetAt(position int) (val string) {
	o.RLock()
	defer o.RUnlock()
	if position < len(o.order) && position >= 0 {
		val, _ = o.items[o.order[position]]
	}
	return
}


// Strings returns a slice of the strings in the order they were added
func (o *SafeStringSliceMap) Values() []string {
	o.RLock()
	defer o.RUnlock()
	vals := make([]string, len(o.order))

    if o.items != nil {
        for i, v := range o.order {
            vals[i] = o.items[v]
        }
    }
	return vals
}

// Keys are the keys of the strings, in the order they were added
func (o *SafeStringSliceMap) Keys() []string {
	o.RLock()
	defer o.RUnlock()
	vals := make([]string, len(o.order))

    if len(o.order) != 0 {
        for i, v := range o.order {
            vals[i] = v
        }
    }
	return vals
}

func (o *SafeStringSliceMap) Len() int {
	o.RLock()
	defer o.RUnlock()
	return len(o.order)
}

func (o *SafeStringSliceMap) Less(i, j int) bool {
	o.RLock()
	defer o.RUnlock()
	return o.items[o.order[i]] < o.items[o.order[j]]
}

func (o *SafeStringSliceMap) Swap(i, j int) {
	o.Lock()
	defer o.Unlock()
	o.order[i], o.order[j] = o.order[j], o.order[i]
}

// Sort by keys interface
type safestringstringsortedbykeys struct {
	// This embedded interface permits Reverse to use the methods of
	// another interface implementation.
	sort.Interface
}

// A helper function to allow SafeStringSliceMaps to be sorted by keys
// To sort the map by keys, call:
//   sort.Sort(OrderStringStringSliceMapByKeys(m))
func OrderSafeStringSliceMapByKeys(o *SafeStringSliceMap) sort.Interface {
	return &safestringstringsortedbykeys{o}
}

// A helper function to allow SafeStringSliceMaps to be sorted by keys
func (r safestringstringsortedbykeys) Less(i, j int) bool {
	var o *SafeStringSliceMap = r.Interface.(*SafeStringSliceMap)
	o.RLock()
	defer o.RUnlock()
	return o.order[i] < o.order[j]
}

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

func (o *SafeStringSliceMap) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	o.Lock()
	defer o.Unlock()
	err := dec.Decode(&o.items)
	if err == nil {
		err = dec.Decode(&o.order)
	}
	return err
}

func (o *SafeStringSliceMap) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
	o.RLock()
	defer o.RUnlock()
	data, err = json.Marshal(o.items)
	return
}

func (o *SafeStringSliceMap) UnmarshalJSON(data []byte) error {
	o.Lock()
	defer o.Unlock()
	err := json.Unmarshal(data, &o.items)
	if err == nil {
		// Create a default order, since these are inherently unordered
		o.order = make([]string, len(o.items))
		i := 0
		for k := range o.items {
			o.order[i] = k
			i++
		}
	}
	return err
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

// Range will call the given function with every key and value in the order
// they were placed in the map, or in if you sorted the map, in your custom order.
// If f returns false, it stops the iteration. This pattern is taken from sync.Map.
func (o *SafeStringSliceMap) Range(f func(key string, value string) bool) {
	if o == nil || o.items == nil {
		return
	}
	o.RLock()
	defer o.RUnlock()
	for _, k := range o.order {
		if !f(k, o.items[k]) {
			break
		}
	}
}

// Equals returns true if the map equals the given map, paying attention only to the content of the
// map and not the order.
func (o *SafeStringSliceMap) Equals(i StringMapI) bool {
	if i == nil {
		return o == nil
	}
	if i.Len() != o.Len() {
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
	o.Lock()
	defer o.Unlock()
	o.items = nil
	o.order = nil
}

func (o *SafeStringSliceMap) IsNil() bool {
	return o == nil
}

func (o *SafeStringSliceMap) String() string {
	var s string

	s = "{"
	o.Range(func(k, v string) bool {
		s += `"` + k + `":"` +
v+ `",`
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


