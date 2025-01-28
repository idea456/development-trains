package graph

// TrainsQueue represents a max priority queue to enqueue trains to be scheduled with packages
type TrainsQueue []Train

func (q *TrainsQueue) Push(train interface{}) {
	*q = append(*q, train.(Train))
}

func (q *TrainsQueue) Pop() interface{} {
	old := *q
	size := len(old)
	last := old[size-1]
	*q = old[:size-1]

	return last
}

// NOTE: trains with higher capacity has higher priority
func (q TrainsQueue) Less(i, j int) bool {
	// if q[i].Capacity == q[j].Capacity {
	// 	return len(q[i].PackagesCarried) < len(q[j].PackagesCarried)
	// }
	return q[i].Capacity >= q[j].Capacity
}

func (q TrainsQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q TrainsQueue) Len() int {
	return len(q)
}
