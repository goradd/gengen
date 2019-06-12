package index

import (
	"sort"
	"strconv"
	"testing"
)

func TestNodeInsert(t *testing.T) {
	n := 100000
	index := makeSampleStringIndex(n)
	index2 := makeSampleStringSlice(n)

	for i := 0; i < n; i++ {
		if string(GetAt(i)) != index2[i] {
			t.Fail()
		}
	}
}

func makeSampleStringIndex(n int) Index {
	index := NewIndex(func(a,b BlockType) bool {
		r := a < b
		return r
	})

	for i := 0; i < n; i++ {
		s := strconv.Itoa(i)
		Add(BlockType(s))
	}
	return index
}

func makeSampleStringIndex2(n int) Index2 {
	index := NewIndex2(func(a,b BlockType) bool {
		r := a < b
		return r
	})

	for i := 0; i < n; i++ {
		s := strconv.Itoa(i)
		Add(BlockType(s))
	}
	return index
}


func makeSampleStringSlice(n int) []string {
	var a []string

	for i := 0; i < n; i++ {
		s := strconv.Itoa(i)
		a = append(a, (s))
	}
	sort.Strings(a)
	return a
}

func makeSampleStringSliceInserts(n int) []string {
	var a []string

	for i := 0; i < n; i++ {
		s := strconv.Itoa(i)
		loc := sort.Search(len(a), func(n int) bool {
			return s < a[n]
		})
		a = append(a, s)
		copy(a[loc+1:], a[loc:])
		a[loc] = s
	}
	return a
}

var indexSize = 200000

func BenchmarkIndexInsert(b *testing.B) {
	blockSize = 5000
	for i := 0; i < b.N; i++ {
		_ = makeSampleStringIndex(indexSize)
	}
}

func BenchmarkIndexInsert2(b *testing.B) {
	blockSize = 10000
	for i := 0; i < b.N; i++ {
		_ = makeSampleStringIndex(indexSize)
	}
}

func BenchmarkIndexInsert25(b *testing.B) {
	blockSize = 15000
	for i := 0; i < b.N; i++ {
		_ = makeSampleStringIndex(indexSize)
	}
}


func BenchmarkIndexInsert3(b *testing.B) {
	blockSize = 20000

	for i := 0; i < b.N; i++ {
		_ = makeSampleStringIndex(indexSize)
	}
}

func BenchmarkIndexInsert4(b *testing.B) {
	blockSize = 50000

	for i := 0; i < b.N; i++ {
		_ = makeSampleStringIndex(indexSize)
	}
}

func BenchmarkIndexInsertNoIndex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = makeSampleStringIndex2(indexSize)
	}
}

/*
func BenchmarkSliceInsert(b *testing.B) {
	n := 10000
	//o := 1000

	for i := 0; i < b.N; i++ {
		_ = makeSampleStringSlice(n)

	}
}

func BenchmarkSliceInsertSort(b *testing.B) {
	n := 10000
	//o := 1000

	a := makeSampleStringSlice(n)
	for i := 0; i < b.N; i++ {
		sort.Strings(a)
	}
}
*/

