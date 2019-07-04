package maps

import (
	"fmt"
	"sort"
	"testing"
	"bytes"
	"encoding/gob"
    "encoding/json"
    "os"
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

func TestMapLoaders(t *testing.T) {
    n := map[string]interface{}{"a":1,"b":"2","c":3.0, "d":true}
    m := NewMapFromMap(n)

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

func ExampleMap_Set() {
	m := NewMap()
	m.Set("a", "Here")
	fmt.Println(m.String())
	// Output: {"a":"Here"}
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
	n.Merge(m)

	fmt.Println(n.Get("C"))
	fmt.Println(n.Get("D"))
	// Output: Other
	// 5
}

func ExampleMap_MergeMap() {
	m := map[string]interface{} {
	    "B": "This",
	    "A": "That",
	    "C": 6.1,
	}

    n := NewMap()
    n.Set("D","Last")
	n.MergeMap(m)

	fmt.Println(n.Get("C"))
	fmt.Println(n.Get("D"))
	// Output: 6.1
	// Last
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

func ExampleMap_MarshalBinary() {
	// You would rarely call MarshallBinary directly, but rather would use an encoder, like GOB for binary encoding

	m := new (Map)
	var m2 Map

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 3)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf) // Will write
	dec := gob.NewDecoder(&buf) // Will read

	enc.Encode(m)
	dec.Decode(&m2)
	s := m2.Get("A")
	fmt.Println(s)
	s = m2.Get("C")
	fmt.Println(s)
	// Output: That
	// 3
}

func ExampleMap_MarshalJSON() {
	// You don't normally call MarshallJSON directly, but rather use the Marshall and Unmarshall json commands
	m := new (Map)

	m.Set("B", "This")
	m.Set("A", "That")
	m.Set("C", 3)

	s, _ := json.Marshal(m)

	// Note: The below output is what is produced, but isn't guaranteed. go seems to currently be sorting keys
	os.Stdout.Write(s)
	// Output: {"A":"That","B":"This","C":3}
}

func ExampleMap_UnmarshalJSON() {
	b := []byte(`{"A":"That","B":"This","C":3}`)
	var m Map

	json.Unmarshal(b, &m)

	fmt.Println(m.Get("C"))

	// Output: 3
}

func TestMapEmpty(t *testing.T) {
    var m *Map
    var n = new(Map)

    if !m.IsNil() {
        t.Error("Empty Nil test failed")
    }

    if n.IsNil() {
        t.Error("Empty Nil test failed")
    }

    for _, o := range ([]*Map{m, n}) {
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
