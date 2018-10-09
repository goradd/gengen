package maps

type Getter interface {
	Get(key string) (val interface{})
}

type Loader interface {
	Load(key string) (val interface{}, ok bool)
}

type Setter interface {
	Set(string, interface{})
}


// The MapI interface provides a common interface to the many kinds of similar map objects.
//
// Most functions that change the map are omitted so that you can wrap the map in additional functionality that might
// use Set or SetChanged. If you want to use them in an interface setting, you can create your own interface
// that includes them.
type MapI interface {
	Get(key string) (val interface{})
	Has(key string) (exists bool)
	Values() []interface{}
	Keys() []string
	Len() int
	// Range will iterate over the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value interface{}) bool)
	Merge(i MapI)
}
