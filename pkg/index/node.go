package index

import "math/rand"

type iNode struct {
	l *iNode
	r *iNode
	data []BlockType
}

func (i iNode) count() (c int) {
	if i.l != nil {
		c += i.l.count()
	}
	c += len(i.data)
	if i.r != nil {
		c += i.r.count()
	}
	return
}

func (i iNode) getAt(n int) BlockType {
	var lCount int
	if i.l != nil {
		lCount = i.l.count()
		if n < lCount {
			return i.l.getAt(n)
		}
	}
	if n < lCount + len(i.data) {
		return i.data[n - lCount]
	}
	if i.r != nil {
		return i.r.getAt(n - lCount - len(i.data))
	}
	var z BlockType
	return z // we are off the end
}

func (i *iNode) insertAt(n int, item BlockType) {
	var lCount int
	if i.l != nil {
		lCount = i.l.count()
		if n < lCount {
			i.l.insertAt(n, item)
			return
		}
	}
	if i.r == nil || n < lCount + len(i.data)  {
		i.insertData(n - lCount, lCount, item)
		return
	}
	i.r.insertAt(n - lCount - len(i.data), item)
}

func (i *iNode) insertData(n int, lCount int, item BlockType) {
	l := len(i.data)
	if n >= l {
		i.data = append(i.data, item)
		return // a very fast operation, no need to break this apart
	}

	if l >= blockSize {
		// break it up and try again
		i.splitNode()
		i.insertAt(n + lCount, item)
		return
	}

	// slice insert
	i.data = append(i.data, item)
	copy(i.data[n+1:], i.data[n:])
	i.data[n] = item
}

func (i *iNode) splitNode() {
	s := len(i.data) / 2

	d := make([]BlockType, s, s)
	copy(d, i.data[s:])
	r := &iNode{data:d} // create a new slice

	l := &iNode{data:i.data[:s]} // reuse slice
	if i.l == nil  || (i.r != nil && rand.Intn(1) == 1) {
		l.l = i.l
		i.l = l
		i.data = r.data
	} else {
		r.r = i.r
		i.r = r
		i.data = l.data
	}
}