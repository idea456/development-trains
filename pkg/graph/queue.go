package graph

type PackagesQueue []Package

func (q *PackagesQueue) Push(newPackage interface{}) {
	delivery := newPackage.(Package)
	*q = append(*q, delivery)
}

func (q *PackagesQueue) Pop() interface{} {
	old := *q
	size := len(old)
	last := old[size-1]
	*q = old[:size-1]

	return last
}

// NOTE: prioritze packages with same destinations, AND nearest amongst one another
// - check their destinations, if same, check their shortest distances using travel time matrix packageX[start][destination] - packageY[start][destination]
// - sort first by destination, then by distance
// - or by weight?
func (q PackagesQueue) Less(i, j int) bool {
	// if q[i].EndingStationId == q[j].EndingStationId {
	// 	return q[i].DistanceToDestination <= q[j].DistanceToDestination
	// }
	return q[i].EndingStationId <= q[j].EndingStationId
}

func (q PackagesQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q PackagesQueue) Len() int {
	return len(q)
}

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
