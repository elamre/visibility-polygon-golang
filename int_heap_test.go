package visibility_polygon_golang

import (
	"container/heap"
	"fmt"
)

func ExampleNewIntHeap() {
	h := &IntHeap{}
	heap.Init(h)
	heap.Push(h, 5)
	heap.Push(h, 1)
	heap.Push(h, 7)
	heap.Push(h, 3)
	for ; h.Len() > 0; {
		fmt.Println(heap.Pop(h))
	}
	// Output:
	// 1
	// 3
	// 5
	// 7
}
