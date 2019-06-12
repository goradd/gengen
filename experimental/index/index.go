// The index package is an experiment to see at what point a single block of memory for the index begins to
// slow down the system. It uses a form of memory management to break the single block into a binary tree
// It appears that once an index reaches a size of about 100,000 items, the speed of inserting 1 item starts to
// significantly slow down. Using a block size of 15,000 items per node in the binary tree restores its
// performance to log(n).
//
// The index is not implemented yet. Once someone has a need for a sorted map with more than 100,000 items, we
// can add that. It does not appear to significantly slow down smaller maps.
// This could be the basis for a multi-index system as well.
package index

import "sort"

var blockSize = 15000

type BlockType string // could be anything usable to index a map

type Index struct {
	itemCount int
	lessF     func(a,b BlockType) bool
	root      iNode
}

func NewIndex(lessF func(a,b BlockType) bool) Index {
	return Index{lessF:lessF, root: iNode{data: make([]BlockType, 0, 32)}}
}

func (i *Index) Add(item BlockType) {
	loc := sort.Search (i.itemCount, func(n int) bool {
		v := i.GetAt(n)
		return !i.lessF(v, item)
	})
	i.insertAt(loc, item)
}

func (i *Index) GetAt(n int) BlockType {
	return getAt(n)
}

func (i *Index) insertAt(loc int, item BlockType) {
	insertAt(loc, item)
	i.itemCount ++
}
