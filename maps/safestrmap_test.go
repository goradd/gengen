package maps

import (
	"fmt"
	"sort"
	"testing"
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

	m.Remove("A")

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
	m := NewStringMap()
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
	m := NewStringMap()
	a := []string{}

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
	m.Merge(n)

	fmt.Println(m.Get("D"))
	//Output: Last
}

func ExampleNewSafeStringMapFrom() {
    n:= NewSafeStringMap()
    n.Set("a", "this")
    n.Set("b", "that")
	m := NewSafeStringMapFrom(n)
	fmt.Println(m.Get("b"))
	//Output: that
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
