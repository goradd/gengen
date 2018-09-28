package maps

import (
	"fmt"
	"sort"
	"testing"
)

func TestMap(t *testing.T) {
	var v interface{}

	m := NewMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 5)

	if v = m.Get("B"); v != "This" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "This", v)
	}

	if v = m.Get("C"); v != 5 {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "Other", v)
	}

	m.Remove("A")

	if m.Len() != 2 {
		t.Error("Len Failed.")
	}

	if m.Has("NOT THERE") {
		t.Error("Getting non-existant value did not return false")
	}

	v = m.Get("B")
	if v != "This" {
		t.Error("Get failed")
	}

	if !m.Has("B") {
		t.Error("Existance test failed.")
	}

	// Can set non-string values

	m.Set("E", 15.5)
	if m.Get("E") != 15.5 {
		t.Error("Setting non-string value failed.")
	}

	// Verify it satisfies the MapI interface
	var i MapI = m
	if i2 := i.Get("B"); i2 != "This" {
		t.Error("MapI interface test failed.")
	}

	m.Clear()
	v = m.Get("B")
	if v != nil {
		t.Error("Clear failed")
	}

	m.Set("E", 15.5)
	if m.Get("E") != 15.5 {
		t.Error("Set after clear failed.")
	}

 	m.Clear()
    m.SetChanged("E", 15.5)
    if m.Get("E") != 15.5 {
        t.Error("SetChanged after clear failed.")
    }

    n := m.Copy()
    if n.Get("E") != 15.5 {
        t.Error("Copy failed.")
    }

}

func TestEmpty(t *testing.T) {
    var m *Map
    var n = new(Map)

    for _, o := range ([]*Map{m, n}) {
        i := o.Get("A")
        if i != nil {
            t.Error("Empty Get failed")
        }
        if o.Has("A") {
            t.Error("Empty Has failed")
        }
        o.Remove("E")
        o.Clear()

        if len(o.Values()) != 0 {
            t.Error("Empty Values() failed")
        }

        if len(o.Keys()) != 0 {
            t.Error("Empty Keys() failed")
        }
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

func TestMapChange(t *testing.T) {
	m := NewMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 5)

	if changed := m.SetChanged("D", 6); !changed {
		t.Error("Set did not produce a change flag")
	}

	if changed := m.SetChanged("D", 6); changed {
		t.Error("Set again erroneously produced a change flag")
	}

    if changed := m.SetChanged("D", "That"); !changed {
        t.Error("Set again did not produce a change flag")
    }
}

func TestMapNotEqual(t *testing.T) {
	m := NewMap()
	m.Set("A", "This")
	m.Set("B","That")
	n := NewMap()
	n.Set("B", "This")
	n.Set("A","That")
	if m.Equals(n) {
		t.Error("Equals test failed")
	}
}

func ExampleMap_Set() {
	m := NewMap()
	m.Set("a", "Here")
	fmt.Println(m.Get("a"))
	// Output Here
}

func ExampleMap_Values() {
	m := NewMap()
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 5)

	values := m.Values()
	var values2 []string
	for _,value := range values {
	    values2 = append(values2, fmt.Sprintf("%v", value))
	}
	sort.Sort(sort.StringSlice(values2))
	fmt.Println(values2)
	//Output: [5 That This]
}

func ExampleMap_Keys() {
	m := NewMap()
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Keys()
	sort.Sort(sort.StringSlice(values))
	fmt.Println(values)
	//Output: [A B C]
}

func ExampleMap_Range() {
	m := NewMap()
	a := []string{}

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 5)

	m.Range(func(key string, val interface{}) bool {
		a = append(a, fmt.Sprintf("%v", val))
		return true // keep iterating to the end
	})
	fmt.Println()

	sort.Sort(sort.StringSlice(a)) // Unordered maps cannot be guaranteed to range in a particular order. Sort it so we can compare it.
	fmt.Println(a)
	//Output: [5 That This]
}

func ExampleMap_Merge() {
	m := NewMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

    n := NewMap()
    n.Set("D",5)
	m.Merge(n)

	fmt.Println(m.Get("D"))
	//Output: 5
}

func ExampleNewMapFrom() {
    n := NewMap()
    n.Set("a", "this")
    n.Set("b", 5)
	m := NewMapFrom(n)
	fmt.Println(m.Get("b"))
	//Output: 5
}

func ExampleMap_Equals() {
	m := NewMap()
	m.Set("A","This")
	m.Set("B", "That")
	n := NewMap()
	n.Set("B", "That")
	n.Set("A", "This")
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}



