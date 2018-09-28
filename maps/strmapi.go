package maps

type StringGetter interface {
	Get(key string) (val string)
}

// The StringMapI interface provides a common interface to the many kinds of similar map objects.
type StringMapI interface {
	SetChanged(key string, val string) (changed bool)
	Set(key string, val string)
	Get(key string) (val string)
	Has(key string) (exists bool)
	Remove(key string)
	Values() []string
	Keys() []string
	Len() int
	Clear()
	// Range will iterate over the keys and values in the map. Pattern is taken from sync.Map
	Range(f func(key string, value string) bool)
	Merge(i StringMapI)
}

