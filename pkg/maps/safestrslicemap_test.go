package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestSafeStringSliceMap(t *testing.T) {
	var s string

	m := new (SafeStringSliceMap)

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

    if k := m.GetKeyAt(2); k != "C" {
        t.Errorf("GetAt test failed. Expected  (%q) got (%q).", 1, s)
    }

	if s = m.GetAt(3); s != "" {
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

	if changed := m.SetChanged("F", "9"); !changed {
		t.Error("Add non-string value failed.")
	}
	if m.Get("F") != "9" {
		t.Error("Add non-string value failed.")
	}
}

func TestSafeStringSliceMapChange(t *testing.T) {
	m := new (SafeStringSliceMap)

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

func ExampleSafeStringSliceMap_Range() {
	m := new (SafeStringSliceMap)

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
    m.SetSortFunc(func(key1, key2 string, val1, val2 string) bool {
        return val1 < val2
    })

	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Iterate after sorting keys
    m.SetSortFunc(func(key1, key2 string, val1, val2 string) bool {
        return key1 < key2
    })
	m.Range(func(key string, val string) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Output: B:This,A:That,C:Other,
	// C:Other,A:That,B:This,
	// A:That,B:This,C:Other,
}

func ExampleSafeStringSliceMap_MarshalBinary() {
	// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding

	m := new (SafeStringSliceMap)
	var m2 SafeStringSliceMap

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

func ExampleSafeStringSliceMap_MarshalJSON() {
	// You don't normally call MarshallJSON directly, but rather use the Marshall and Unmarshall json commands
	m := new (SafeStringSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	s, _ := json.Marshal(m)
	os.Stdout.Write(s)

	// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleSafeStringSliceMap_UnmarshalJSON() {
	b := []byte(`{"A":"That","B":"This","C":"Other"}`)
	var m SafeStringSliceMap

	json.Unmarshal(b, &m)
	m.SortByKeys()

	fmt.Println(&m)

	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleSafeStringSliceMap_Merge() {
	m := new (SafeStringSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

    n := new (StringMap)
    n.Set("D", "Last")
	m.Merge(n)

	fmt.Println(m.GetAt(3))
	//Output: Last
}

func ExampleSafeStringSliceMap_Values() {
	m := new (SafeStringSliceMap)
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Values()
	fmt.Println(values)
	//Output: [This That Other]
}

func ExampleSafeStringSliceMap_Keys() {
	m := new (SafeStringSliceMap)
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Keys()
	fmt.Println(values)
	//Output: [B A C]
}

func ExampleNewSafeStringSliceMapFrom() {
    n := new (StringMap)
    n.Set("a", "this")
    n.Set("b", "that")
	m := NewSafeStringSliceMapFrom(n)
	fmt.Println(m.Get("b"))
	//Output: that
}

func ExampleSafeStringSliceMap_Equals() {
    n := new (StringMap)
    n.Set("A", "This")
    n.Set("B", "That")
	m := NewSafeStringSliceMapFrom(n)
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}

func TestSafeStringSliceMap_SetAt(t *testing.T) {
	m := NewSafeStringSliceMap()

	m.Set("a", "A")
	m.Set("b", "B")

	// Test middle inserts
	m.SetAt(1, "c", "C")
	if "C" != m.GetAt(1) {
	    t.Errorf("Middle insert failed. Expected C and got %s", m.GetAt(1))
	}

	m.SetAt(-1, "d", "D")
    if "D" != m.GetAt(2) {
        t.Errorf("Middle insert failed. Expected D and got %s", m.GetAt(2))
    }
    if "B" != m.GetAt(3) {
        t.Errorf("Middle insert failed. Expected B and got %s", m.GetAt(3))
    }

	// Test end inserts
	m.SetAt(m.Len(), "e", "E")
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
