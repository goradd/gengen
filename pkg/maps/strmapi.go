package maps

type StringGetter interface {
	Get(key string) (val string)
}

type StringLoader interface {
	Load(key string) (val string, ok bool)
}

type StringSetter interface {
	Set(string, string)
}


// The StringMapI interface provides a common interface to the many kinds of similar map objects.
//
// Most functions that change the map are omitted so that you can wrap the map in additional functionality that might
// use Set or SetChanged. If you want to use them in an interface setting, you can create your own interface
// that includes them.
type StringMapI interface {
	Get(key string) (val string)
	Has(key string) (exists bool)
	Values() []string
	Keys() []string
	Len() int
	// Range will iterate over the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value string) bool)
	Merge(i StringMapI)
	String() string
}
