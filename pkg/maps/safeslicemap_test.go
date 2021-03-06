package maps

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestSafeSliceMap(t *testing.T) {
	var s string

	m := new (SafeSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 1)

	if m.Values()[1] != "That" {
		t.Errorf("Strings test failed. Expected  (%q) got (%q).", "That", m.Values()[1])
	}

	if m.Keys()[1] != "A" {
		t.Errorf("Keys test failed. Expected  (%q) got (%q).", "A", m.Keys()[1])
	}

	if i := m.GetAt(2); i != 1 {
		t.Errorf("GetAt test failed. Expected  (%q) got (%q).", 1, s)
	}

    if k := m.GetKeyAt(2); k != "C" {
        t.Errorf("GetAt test failed. Expected  (%q) got (%q).", 1, s)
    }

	if m.GetAt(3) != nil {
		t.Errorf("GetAt test failed. Expected no response, got %q", s)
	}

	m.Delete("A")

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

	// Test that it satisfies the MapI interface
	var i MapI = m
	if i := i.Get("B"); i != "This" {
		t.Error("MapI interface test failed.")
	}

	m.Set("F", 9)
	if m.Get("F") != 9 {
		t.Error("Add non-string value failed.")
	}

	n := m.Copy()
    if n.Get("F") != 9 {
        t.Error("Copy failed.")
    }

}


func ExampleSafeSliceMap_Range() {
	m := new (SafeSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	// Iterate by insertion order
	m.Range(func(key string, val interface{}) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()


	// Iterate after sorting keys
	m.SortByKeys()
	m.Set("D", "Other2")

	m.Range(func(key string, val interface{}) bool {
		fmt.Printf("%s:%s,", key, val)
		return true // keep iterating to the end
	})
	fmt.Println()

	// Output: B:This,A:That,C:Other,
	// A:That,B:This,C:Other,D:Other2,
}

func ExampleSafeSliceMap_MarshalBinary() {
	// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding

	m := new (SafeSliceMap)
	var m2 SafeSliceMap

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

func ExampleSafeSliceMap_MarshalJSON() {
	// You don't normally call MarshallJSON directly, but rather use the Marshall and Unmarshall json commands
	m := new (SafeSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	s, _ := json.Marshal(m)

	// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
	os.Stdout.Write(s)
	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleSafeSliceMap_UnmarshalJSON() {
	b := []byte(`{"A":"That","B":"This","C":"Other"}`)
	var m SafeSliceMap

	json.Unmarshal(b, &m)
	m.SortByKeys()
	fmt.Println(&m)

	// Output: {"A":"That","B":"This","C":"Other"}
}

func ExampleSafeSliceMap_Merge() {
	m := new (SafeSliceMap)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 5)

    n := new (SafeSliceMap)
    n.SortByKeys()
    n.Set("D", "Last")
	n.Merge(m)
	values := n.Values()
	fmt.Println(values)
	//Output: [That This 5 Last]
}

func ExampleSafeSliceMap_MergeMap() {
	m := map[string]interface{} {
	    "B": "This",
	    "A": "That",
	    "C": 5,
	}

    n := NewSafeSliceMap()
    n.SortByKeys()
    n.Set("D","Last")
	n.MergeMap(m)
	values := n.Values()
	fmt.Println(values)
	// Output: [That This 5 Last]
}


func ExampleSafeSliceMap_Values() {
	m := new (SafeSliceMap)
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Values()
	fmt.Println(values)
	//Output: [This That Other]
}

func ExampleSafeSliceMap_Keys() {
	m := new (SafeSliceMap)
	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", "Other")

	values := m.Keys()
	fmt.Println(values)
	//Output: [B A C]
}

func ExampleNewSafeSliceMapFrom() {
    n := new (Map)
    n.Set("a", "this")
    n.Set("b", "that")
	m := NewSafeSliceMapFrom(n)
	fmt.Println(m.Get("b"))
	//Output: that
}

func ExampleSafeSliceMap_Equals() {
    n := new (Map)
    n.Set("A", "This")
    n.Set("B", "That")
	m := NewSafeSliceMapFrom(n)
	if m.Equals(n) {
		fmt.Print("Equal")
	} else {
		fmt.Print("Not Equal")
	}
	//Output: Equal
}

func TestSafeSliceMap_SetAt(t *testing.T) {
	m := NewSafeSliceMap()

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

func TestSafeSliceMapLoaders(t *testing.T) {
    n := map[string]interface{}{"a":1,"b":"2","c":3.0, "d":true}
    m := NewSafeSliceMapFromMap(n)

    if i,ok := m.LoadInt("a"); i != 1 || !ok {
        t.Error("LoadInt failed")
    }
    if j,ok := m.LoadString("b"); j != "2" || !ok {
        t.Error("LoadString failed")
    }
    if k,ok := m.LoadFloat64("c"); k != 3.0 || !ok {
        t.Error("LoadFloat failed")
    }
    if l,ok := m.LoadBool("d"); l != true || !ok {
        t.Error("LoadBool failed")
    }

    if _,ok := m.LoadFloat64("d"); ok {
        t.Error("Type check failed")
    }

}



func TestSafeSliceMapEmpty(t *testing.T) {
    var m *SafeSliceMap
    var n = new(SafeSliceMap)

    if !m.IsNil() {
        t.Error("Empty Nil test failed")
    }

    if n.IsNil() {
        t.Error("Empty Nil test failed")
    }


    for _, o := range ([]*SafeSliceMap{m, n}) {
        i := o.Get("A")
        if i != nil {
            t.Error("Empty Get failed")
        }

        i = o.GetAt(5)
        if i != nil {
            t.Error("Empty GetAt failed")
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
        o.Range(func (k string, v interface{}) bool {
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
