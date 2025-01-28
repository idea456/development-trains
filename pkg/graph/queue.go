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

func (q TrainsQueue) Less(i, j int) bool {
	// // NOTE: first prioritize trains that carries less packages first to encourage and diversify other trains to pick this up instead
	// // then prioritize trains that has a bigger capacity to carry more packages
	if len(q[i].PackagesCarried) == len(q[j].PackagesCarried) {
		// prioritize trains that has a bigger capacity to carry more packages
		return q[i].Capacity > q[j].Capacity
	}
	// prioritize trains that carries less packages first to encourage and diversify other trains to pick this up instead
	return len(q[i].PackagesCarried) < len(q[j].PackagesCarried)
}

func (q TrainsQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q TrainsQueue) Len() int {
	return len(q)
}
