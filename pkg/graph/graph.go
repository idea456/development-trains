package graph

type Train struct {
	Name            string
	StartingStation Station
}

type Package struct {
	Name            string
	Weight          int
	StartingStation Station
	EndingStation   Station
}

type Station struct {
	Name   string
	Routes []Route
}

type Route struct {
	Name       string
	TravelTime int
	Stations   [2]Station
}

type Graph struct{}
