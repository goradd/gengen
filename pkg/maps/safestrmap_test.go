package maps

import (
	"fmt"
	"sort"
	"testing"
	"bytes"
	"encoding/gob"
)

func TestSafeStringMap(t *testing.T) {
	var s string

	m := NewSafeStringMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	if s = m.Get("B"); s != "This" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "This", s)
	}

	if s = m.Get("C"); s != "Other" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "Other", s)
	}

	m.Delete("A")

	if m.Len() != 2 {
		t.Error("Len Failed.")
	}

	if m.Has("NOT THERE") {
		t.Error("Getting non-existant value did not return false")
	}

	s = m.Get("B")
	if s != "This" {
		t.Error("Get failed")
	}

	if !m.Has("B") {
		t.Error("Existance test failed.")
	}

	// Can set non-string values

	m.Set("E", "8")
	if m.Get("E") != "8" {
		t.Error("Setting non-string value failed.")
	}

	// Verify it satisfies the StringMapI interface
	var i StringMapI = m
	if s := i.Get("B"); s != "This" {
		t.Error("StringMapI interface test failed.")
	}

	m.Clear()
	s = m.Get("B")
	if s != "" {
		t.Error("Clear failed")
	}
}

func TestSafeStringMapChange(t *testing.T) {
	m := NewSafeStringMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	if changed := m.SetChanged("D", "And another"); !changed {
		t.Error("Set did not produce a change flag")
	}

	if changed := m.SetChanged("D", "And another"); changed {
		t.Error("Set again erroneously produced a change flag")
	}

    if changed := m.SetChanged("D", "That"); !changed {
        t.Error("Set again did not produce a change flag")
    }

}

func TestSafeStringMapNotEqual(t *testing.T) {
	m := NewSafeStringMap()
	m.Set("A", "This")
	m.Set("B","That")
	n := NewSafeStringMap()
	n.Set("B", "This")
	n.Set("A","That")
	if m.Equals(n) {
		t.Error("Equals test failed")
	}
}

func ExampleSafeStringMap_Set() {
	m := NewSafeStringMap()
	m.Set("a", "Here")
	fmt.Println(m.Get("a"))
	// Output Here
}


func ExampleSafeStringMap_Values() {
	m := NewSafeStringMap()
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Values()
	sort.Sort(sort.StringSlice(values))
	fmt.Println(values)
	//Output: [Other That This]
}

func ExampleSafeStringMap_Keys() {
	m := NewSafeStringMap()
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Keys()
	sort.Sort(sort.StringSlice(values))
	fmt.Println(values)
	//Output: [A B C]
}

func ExampleSafeStringMap_Range() {
	m := NewSafeStringMap()
	var a []string

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	m.Range(func(key string, val string) bool {
		a = append(a, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	sort.Sort(sort.StringSlice(a)) // unordered maps cannot be guaranteed to range in a particular order. Sort it so we can compare it.
	fmt.Println(a)
	//Output: [Other That This]

}

func ExampleSafeStringMap_Merge() {
	m := NewSafeStringMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

    n := NewSafeStringMap()
    n.Set("D","Last")
	n.Merge(m)

	fmt.Println(n.Get("C"))
	fmt.Println(n.Get("D"))
	// Output: Other
	// Last
}

func ExampleSafeStringMap_MergeMap() {
	m := map[string]string {
	    "B": "This",
	    "A": "That",
	    "C": "Other",
	}

    n := NewSafeStringMap()
    n.Set("D","Last")
	n.MergeMap(m)

	fmt.Println(n.Get("C"))
	fmt.Println(n.Get("D"))
	// Output: Other
	// Last
}


func ExampleNewSafeStringMapFrom() {
    n:= NewSafeStringMap()
    n.Set("a", "this")
    n.Set("b", "that")
	m := NewSafeStringMapFrom(n)

	fmt.Println(m.Get("b"))
	//Output: that
}

func ExampleNewSafeStringMapFromMap() {
    n:= map[string]string{"a":"this","b":"that"}
	m := NewSafeStringMapFromMap(n)

	fmt.Println(m.String())
	// Output: {"a":"this","b":"that"}
}


func ExampleSafeStringMap_Equals() {
	m := NewSafeStringMap()
	m.Set("A","This")
	m.Set("B", "That")
	n := NewSafeStringMap()
	n.Set("B", "That")
	n.Set("A", "This")
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}

func TestSafeStringMapCopy(t *testing.T) {
    n:= map[string]string{"a":"this","b":"that","c":"other"}
	m := NewSafeStringMapFromMap(n)
	c := m.Copy()
	m.Delete("b")
	if !c.Has("b") {
	    t.Error("Underlying data did not copy")
	}
    if c.String() != `{"a":"this","b":"that","c":"other"}` {
	    t.Error("Did not copy")
    }
}

func TestSafeStringMapEmpty(t *testing.T) {
    var m *SafeStringMap
    var n = new(SafeStringMap)

    if !m.IsNil() {
        t.Error("Empty Nil test failed")
    }

    if n.IsNil() {
        t.Error("Empty Nil test failed")
    }

    for _, o := range ([]*SafeStringMap{m, n}) {
        i := o.Get("A")
        if i != "" {
            t.Error("Empty Get failed")
        }
        if o.Has("A") {
            t.Error("Empty Has failed")
        }
        o.Delete("E")
        o.Clear()

        if len(o.Values()) != 0 {
            t.Error("Empty Values() failed")
        }

        if len(o.Keys()) != 0 {
            t.Error("Empty Keys() failed")
        }

        var j int
        o.Range(func (k string, v string) bool {
            j = 1
            return false
        })
        if j == 1 {
            t.Error("Empty Range failed")
        }

        o.Merge(nil)

    }

    if !m.Equals(n) {
        t.Error("Empty Equals() failed")
    }
    n.Set("a","b")
    if m.Equals(n) {
       t.Error("Empty Equals() failed")
    }
    if n.Equals(m) {
       t.Error("Empty Equals() failed")
    }


}

func TestSafeStringMap_MarshalBinary(t *testing.T) {
	m := new (SafeStringMap)
	var m2 SafeStringMap

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // Will write
	dec := gob.NewDecoder(&buf) // Will read

	enc.Encode(m)
	dec.Decode(&m2)
	if s := m2.Get("A"); s != "That" {
	    t.Error("MarshalBinary failed")
	}
	if s := m2.Get("B"); s != "This" {
	    t.Error("MarshalBinary failed")
	}
}
