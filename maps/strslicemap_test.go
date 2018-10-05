package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"testing"
)

func TestStringSliceMap(t *testing.T) {
	var s string
	var ok bool

	m := new (StringSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	if m.Values()[1] != "That" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "That", m.Values()[1])
	}

	if m.Keys()[1] != "A" {
		t.Errorf("Keys test failed. Expected  (%q) got (%q).", "A", m.Keys()[1])
	}

	if s = m.GetAt(2); s != "Other" {
		t.Errorf("GetAt test failed. Expected  (%q) got (%q).", "Other", s)
	}

	if s = m.GetAt(3); ok {
		t.Errorf("GetAt test failed. Expected no response, got %q", s)
	}

	s = m.Join("+")

	if s != "This+That+Other" {
		t.Error("Failed Join.")
	}

	m.Delete("A")

	s = m.Join("-")

	if s != "This-Other" {
		t.Error("Delete Failed.")
	}

	if m.Len() != 2 {
		t.Error("Len Failed.")
	}

	if m.Has("NOT THERE") {
		t.Error("Getting non-existant value did not return false")
	}

	val := m.Get("B")
	if val != "This" {
		t.Error("Get failed")
	}

	// Test that it satisfies the StringMapI interface
	var i StringMapI = m
	if s = i.Get("B"); s != "This" {
		t.Error("StringMapI interface test failed.")
	}

	// Test that it satisfies the Sort interface
	sort.Sort(m)
	if s = m.GetAt(0); s != "Other" {
		t.Error("Sort interface test failed.")
	}

	if changed := m.SetChanged("F", "9"); !changed {
		t.Error("Add non-string value failed.")
	}
	if m.Get("F") != "9" {
		t.Error("Add non-string value failed.")
	}
}

func TestStringSliceMapChange(t *testing.T) {
	m := new (StringSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	if changed := m.SetChanged("D", "And another"); !changed {
		t.Error("Set did not produce a change flag")
	}

	if changed := m.SetChanged("D", "And another"); changed {
		t.Error("Set again erroneously produced a change flag")
	}
}

func ExampleStringSliceMap_Range() {
	m := new (StringSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	// Iterate by insertion order
	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Iterate after sorting values
	sort.Sort(m)
	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Iterate after sorting keys
	sort.Sort(OrderStringSliceMapByKeys(m))
	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Output: B:This,A:That,C:Other,
	// C:Other,A:That,B:This,
	// A:That,B:This,C:Other,
}

func ExampleStringSliceMap_MarshalBinary() {
	// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding

	m := new (StringSliceMap)
	var m2 StringSliceMap

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // Will write
	dec := gob.NewDecoder(&buf) // Will read

	enc.Encode(m)
	dec.Decode(&m2)
	s := m2.Get("A")
	fmt.Println(s)
	s = m2.GetAt(2)
	fmt.Println(s)
	// Output: That
	// Other
}

func ExampleStringSliceMap_MarshalJSON() {
	// You don't normally call MarshallJSON directly, but rather use the Marshall and Unmarshall json commands
	m := new (StringSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	s, _ := json.Marshal(m)
	os.Stdout.Write(s)

	// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleStringSliceMap_UnmarshalJSON() {
	b := []byte(`{"A":"That","B":"This","C":"Other"}`)
	var m StringSliceMap

	json.Unmarshal(b, &m)
	sort.Sort(OrderStringSliceMapByKeys(&m))

	fmt.Println(&m)

	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleStringSliceMap_Merge() {
	m := new (StringSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

    n := new (StringMap)
    n.Set("D", "Last")
	m.Merge(n)

	fmt.Println(m.GetAt(3))
	//Output: Last
}

func ExampleStringSliceMap_Values() {
	m := new (StringSliceMap)
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Values()
	fmt.Println(values)
	//Output: [This That Other]
}

func ExampleStringSliceMap_Keys() {
	m := new (StringSliceMap)
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Keys()
	fmt.Println(values)
	//Output: [B A C]
}

func ExampleNewStringSliceMapFrom() {
    n := new (StringMap)
    n.Set("a", "this")
    n.Set("b", "that")
	m := NewStringSliceMapFrom(n)
	fmt.Println(m.Get("b"))
	//Output: that
}

func ExampleStringSliceMap_Equals() {
    n := new (StringMap)
    n.Set("A", "This")
    n.Set("B", "That")
	m := NewStringSliceMapFrom(n)
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}

func TestStringSliceMap_SetAt(t *testing.T) {
	m := NewStringSliceMap()

	m.Set("a", "A")
	m.Set("b", "B")

	// Test middle inserts
	m.SetAt(1, "c", "C")
	if "C" != m.GetAt(1) {
	    t.Errorf("Middle insert failed. Expected C and got %s", m.GetAt(1))
	}

	m.SetAt(-2, "d", "D")
    if "D" != m.GetAt(2) {
        t.Errorf("Middle insert failed. Expected D and got %s", m.GetAt(2))
    }
    if "B" != m.GetAt(3) {
        t.Errorf("Middle insert failed. Expected B and got %s", m.GetAt(3))
    }

	// Test end inserts
	m.SetAt(-1, "e", "E")
	m.SetAt(1000, "f", "F")
    if "E" != m.GetAt(4) {
        t.Errorf("End insert failed. Expected E and got %s", m.GetAt(4))
    }
    if "F" != m.GetAt(5) {
        t.Errorf("End insert failed. Expected F and got %s", m.GetAt(5))
    }

	// Test beginning inserts
	m.SetAt(0, "g", "G")
	m.SetAt(-1000, "h", "H")
    if "H" != m.GetAt(0) {
        t.Errorf("Beginning insert failed. Expected H and got %s", m.GetAt(0))
    }
    if "G" != m.GetAt(1) {
        t.Errorf("Beginning insert failed. Expected G and got %s", m.GetAt(1))
    }
}
