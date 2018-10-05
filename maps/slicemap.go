package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"sort"
	"strings"

	"fmt"

)

// A SliceMap combines a map with a slice so that you can range over a
// map in a predictable order. By default, the order will be the same order that items were inserted,
// i.e. a FIFO list. This is similar to how PHP arrays work. You can change this order by providing a
// sorting mechanism. However, this object is NOT safe for concurrent use.
// The zero of this is usable immediately.
// The SliceMap satisfies the MapI interface.
type SliceMap struct {
	items map[string]interface{}
	order []string
}

func NewSliceMap() *SliceMap {
	return new (SliceMap)
}

func NewSliceMapFrom(i MapI) *SliceMap {
	m := new (SliceMap)
	m.Merge(i)
	return m
}

// Copy will make a copy of the map and a copy of the underlying data.
// If the interfaces implement the Copier interface, the Copy function will
// be called to deep copy the items.
func (o *SliceMap) Copy() MapI {
	cp := NewSliceMap()

	o.Range(func(key string, value interface{}) bool {

		if copier, ok := value.(Copier); ok {
			value = copier.Copy()
		}

		cp.Set(key, value)
		return true
	})
	return cp
}


// SetChanged sets the value, but also appends the value to the end of the list for when you
// iterate over the list. Returns whether something changed, and if an error occurred. If the key
// was already in the map, the order will not change, but the value will be replaced. If you want the
// order to change, you must Delete then SetChanged
func (o *SliceMap) SetChanged(key string, val interface{}) (changed bool) {
	var ok bool
	var oldVal interface{}

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
	return
}

// Set sets the given key to the given value
// If the key already exists, the range order will not change.
func (o *SliceMap) Set(key string, val interface{}) {
	o.SetChanged(key, val)
}

// SetAt sets the given key to the given value, but also inserts it at the index specified.  If the index is bigger than
// the length, or -1, it is the same as Set, in that it puts it at the end. Negative indexes are backwards from the
// end, if smaller than the negative length, just inserts at the beginning.
func (o *SliceMap) SetAt(index int, key string, val interface{})  {
	if index == -1 || index >= len(o.items) {
		o.Set(key, val)
		return
	}

	var ok bool
	var emptyKey string

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
}

// Delete removes the item with the given key.
func (o *SliceMap) Delete(key string) {
	for i, v := range o.order {
		if v == key {
			o.order = append(o.order[:i], o.order[i+1:]...)
			break
		}
	}
	delete(o.items, key)
}

// Get returns the value based on its key. If the key does not exist, it will return an empty value.
func (o *SliceMap) Get(key string) (val interface{}) {
    if o.items != nil {
    	val, _ = o.items[key]
    }
	return
}

// Has returns true if the given key exists in the map.
func (o *SliceMap) Has(key string) (ok bool) {
    if o.items == nil {
        return false
    }
	_, ok = o.items[key]
	return
}

// GetAt returns the value based on its position.
func (o *SliceMap) GetAt(position int) (val interface{}) {
	if position < len(o.order) && position >= 0 {
		val, _ = o.items[o.order[position]]
	}
	return
}


// Strings returns a slice of the strings in the order they were added
func (o *SliceMap) Values() []interface{} {
	vals := make([]interface{}, len(o.order))

    if o.items != nil {
        for i, v := range o.order {
            vals[i] = o.items[v]
        }
    }
	return vals
}

// Keys are the keys of the strings, in the order they were added
func (o *SliceMap) Keys() []string {
	vals := make([]string, len(o.order))

    if len(o.order) != 0 {
        for i, v := range o.order {
            vals[i] = v
        }
    }
	return vals
}

// Len returns the number of items in the slice
func (o *SliceMap) Len() int {
	return len(o.order)
}

// Less is part of the interface that allows the map to be sorted by values.
// It returns true if the value at position i should be sorted before the value at position j.
func (o *SliceMap) Less(i, j int) bool {

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
func (o *SliceMap) Swap(i, j int) {
	o.order[i], o.order[j] = o.order[j], o.order[i]
}

// Sort by keys interface
type sortbykeys struct {
	// This embedded interface permits Reverse to use the methods of
	// another interface implementation.
	sort.Interface
}

// A helper function to allow SliceMaps to be sorted by keys
// To sort the map by keys, call:
//   sort.Sort(OrderStringSliceMapByKeys(m))
func OrderSliceMapByKeys(o *SliceMap) sort.Interface {
	return &sortbykeys{o}
}

// A helper function to allow SliceMaps to be sorted by keys
func (r sortbykeys) Less(i, j int) bool {
	var o *SliceMap = r.Interface.(*SliceMap)

	return o.order[i] < o.order[j]

}

func (o *SliceMap) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err = encoder.Encode(o.items)
	if err == nil {
		err = encoder.Encode(o.order)
	}
	data = buf.Bytes()
	return
}

func (o *SliceMap) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&o.items)
	if err == nil {
		err = dec.Decode(&o.order)
	}
	return err
}

func (o *SliceMap) MarshalJSON() (data []byte, err error) {
	// Json objects are unordered
	data, err = json.Marshal(o.items)
	return
}

func (o *SliceMap) UnmarshalJSON(data []byte) error {
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
func (o *SliceMap) Merge(i MapI) {
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
func (o *SliceMap) Range(f func(key string, value interface{}) bool) {
	if o == nil || o.items == nil {
		return
	}
	for _, k := range o.order {
		if !f(k, o.items[k]) {
			break
		}
	}
}

// Equals returns true if the map equals the given map, paying attention only to the content of the
// map and not the order.
func (o *SliceMap) Equals(i MapI) bool {
	if i == nil {
		return o == nil
	}
	if i.Len() != o.Len() {
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

func (o *SliceMap) Clear() {
    if o == nil {return}
	o.items = nil
	o.order = nil
}

func (o *SliceMap) IsNil() bool {
	return o == nil
}

func (o *SliceMap) String() string {
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



