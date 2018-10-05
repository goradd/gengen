package maps

import (
	"fmt"
	"sort"
	"testing"
)

func TestSafeMap(t *testing.T) {
	var v interface{}

	m := NewSafeMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 5)

	if v = m.Get("B"); v != "This" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "This", v)
	}

	if v = m.Get("C"); v != 5 {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "Other", v)
	}

	m.Delete("A")

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

func TestSafeEmpty(t *testing.T) {
    var m *SafeMap
    var n = new(SafeMap)

    for _, o := range ([]*SafeMap{m, n}) {
        i := o.Get("A")
        if i != nil {
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

func TestSafeMapChange(t *testing.T) {
	m := NewSafeMap()

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

func TestSafeMapNotEqual(t *testing.T) {
	m := NewSafeMap()
	m.Set("A", "This")
	m.Set("B","That")
	n := NewSafeMap()
	n.Set("B", "This")
	n.Set("A","That")
	if m.Equals(n) {
		t.Error("Equals test failed")
	}
}

func ExampleSafeMap_Set() {
	m := NewSafeMap()
	m.Set("a", "Here")
	fmt.Println(m.Get("a"))
	// Output Here
}

func ExampleSafeMap_Values() {
	m := NewSafeMap()
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

func ExampleSafeMap_Keys() {
	m := NewSafeMap()
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Keys()
	sort.Sort(sort.StringSlice(values))
	fmt.Println(values)
	//Output: [A B C]
}

func ExampleSafeMap_Range() {
	m := NewSafeMap()
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

func ExampleSafeMap_Merge() {
	m := NewSafeMap()

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

    n := NewSafeMap()
    n.Set("D",5)
	m.Merge(n)

	fmt.Println(m.Get("D"))
	//Output: 5
}

func ExampleNewSafeMapFrom() {
    n := NewSafeMap()
    n.Set("a", "this")
    n.Set("b", 5)
	m := NewSafeMapFrom(n)
	fmt.Println(m.Get("b"))
	//Output: 5
}

func ExampleSafeMap_Equals() {
	m := NewSafeMap()
	m.Set("A","This")
	m.Set("B", "That")
	n := NewSafeMap()
	n.Set("B", "That")
	n.Set("A", "This")
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}

// Test the ability of the copy operation to do a deep copy if available

type toCopySafe struct {
    A int
    b string
}

func (c *toCopySafe) Copy() interface{} {   // normally this would be a more descriptive interface
    m := &toCopySafe{}
    m.A = c.A
    m.b = c.b
    return m
}

func TestSafeCopy(t *testing.T) {
    var c = toCopySafe{2,"s"}
    m := NewSafeMap()
    m.Set("this", &c)
    n := m.Copy()
    c.A = 5
    if n.Get("this").(*toCopySafe).A != 2 {
       t.Error(fmt.Sprintf("Simulated copy failed. A = %d", n.Get("this").(*toCopySafe).A ))
    }
}

