package index

import "sort"


type Index2 struct {
	lessF func(a,b BlockType) bool
	data []BlockType
}

func NewIndex2(lessF func(a,b BlockType) bool) Index2 {
	return Index2{lessF:lessF}
}

func (i *Index2) Add(item BlockType) {
	loc := sort.Search (len(i.data), func(n int) bool {
		v := i.GetAt(n)
		return !i.lessF(v, item)
	})
	i.insertAt(loc, item)
}

func (i *Index2) GetAt(n int) BlockType {
	return i.data[n]
}

func (i *Index2) insertAt(n int, item BlockType) {
	i.data = append(i.data, item)
	copy(i.data[n+1:], i.data[n:])
	i.data[n] = item
}
