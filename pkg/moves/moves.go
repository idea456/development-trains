package moves

import "github.com/idea456/development-trains/pkg/graph"

type Move struct {
	StartedSeconds  int
	Train           *graph.Train
	StartingStation *graph.Station
	PickedPackages  []*graph.Package
	EndingStation   *graph.Station
	DroppedPackages *graph.Package
}
type MovesScheduler struct {
	Moves    []Move
	Stations map[StationId]*Station
}

func NewMovesScheduler() *MovesScheduler {
	return &MovesScheduler{
		Moves: make([]Move, 0),
	}
}

func (m *MovesScheduler) MoveToStation(fromStationId int, toStationId int) {

}
