package maps

type Getter interface {
	Get(key string) (val interface{})
}

// The MapI interface provides a common interface to the many kinds of similar map objects.
type MapI interface {
	SetChanged(key string, val interface{}) (changed bool)
	Set(key string, val interface{})
	Get(key string) (val interface{})
	Has(key string) (exists bool)
	Remove(key string)
	Values() []interface{}
	Keys() []string
	Len() int
	Clear()
	// Range will iterate over the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value interface{}) bool)
	Merge(i MapI)
}

