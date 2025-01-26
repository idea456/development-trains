package graph

import (
	"fmt"
	"strconv"
	"strings"
)

type StationId = int

type Train struct {
	Name              string
	Capacity          int
	StartingStationId StationId
}

type Package struct {
	Name              string
	Weight            int
	StartingStationId StationId
	EndingStationId   StationId
}

type Station struct {
	Id   StationId
	Name string
}

type Route struct {
	Name       string
	TravelTime int
}

type Graph struct {
	Stations         map[StationId]*Station
	StationNames     map[StationId]string
	Routes           map[StationId]map[StationId]*Route
	Deliveries       []Package
	Trains           []Train
	TravelTimeMatrix map[StationId]map[StationId]int
}

func NewGraph(stationNames []string, rawRoutes []string, rawDeliveries []string, rawTrains []string) (*Graph, error) {
	stations := make(map[int]*Station, 0)
	stationNamesToIdMap := make(map[string]int, 0)
	stationNamesMap := make(map[int]string, 0)
	for i, stationName := range stationNames {
		stations[i] = &Station{
			Id:   i,
			Name: stationName,
		}
		stationNamesToIdMap[stationName] = i
		stationNamesMap[i] = stationName
	}

	routes := make(map[int]map[int]*Route, 0)
	for _, rawRoute := range rawRoutes {
		route := strings.Split(rawRoute, ",")
		routeName := route[0]
		fromStation := stationNamesToIdMap[route[1]]
		toStation := stationNamesToIdMap[route[2]]
		travelTime, err := strconv.Atoi(route[3])
		if err != nil {
			return nil, fmt.Errorf("Route %s is not in integer format", routeName)
		}

		if _, exists := routes[fromStation]; !exists {
			routes[fromStation] = make(map[int]*Route, 0)
		}
		if _, exists := routes[toStation]; !exists {
			routes[toStation] = make(map[int]*Route, 0)
		}

		// bidirectional
		routes[fromStation][toStation] = &Route{
			Name:       routeName,
			TravelTime: travelTime,
		}
		routes[toStation][fromStation] = &Route{
			Name:       routeName,
			TravelTime: travelTime,
		}
	}

	deliveries := make([]Package, 0)
	for _, rawDelivery := range rawDeliveries {
		delivery := strings.Split(rawDelivery, ",")
		packageName := delivery[0]
		weight, err := strconv.Atoi(delivery[1])
		if err != nil {
			return nil, fmt.Errorf("Package %s weight is not in integer format", packageName)
		}
		fromStationName := delivery[2]
		toStationName := delivery[3]

		fromStationId := stationNamesToIdMap[fromStationName]
		toStationId := stationNamesToIdMap[toStationName]

		deliveries = append(deliveries, Package{
			Name:              packageName,
			Weight:            weight,
			StartingStationId: fromStationId,
			EndingStationId:   toStationId,
		})
	}

	trains := make([]Train, 0)
	for _, rawTrain := range rawTrains {
		train := strings.Split(rawTrain, ",")
		trainName := train[0]
		capacity, err := strconv.Atoi(train[1])
		if err != nil {
			return nil, fmt.Errorf("Train %s capacity is not an integer", trainName)
		}
		startingStationName := train[2]

		startingStationId := stationNamesToIdMap[startingStationName]

		trains = append(trains, Train{
			Name:              trainName,
			Capacity:          capacity,
			StartingStationId: startingStationId,
		})
	}

	return &Graph{
		Stations:     stations,
		StationNames: stationNamesMap,
		Routes:       routes,
		Deliveries:   deliveries,
		Trains:       trains,
	}, nil
}

func (g *Graph) PrintRoutes() {
	for startionStationId, startingStation := range g.Routes {
		startingStationName := g.StationNames[startionStationId]
		for endingStationId, route := range startingStation {
			endingStationName := g.StationNames[endingStationId]
			fmt.Printf("Station %s to %s: route %s with %d minutes\n", startingStationName, endingStationName, route.Name, route.TravelTime)
		}
	}
}

func (g *Graph) PrintShortestRoutes() {
	for startionStationId, startingStation := range g.TravelTimeMatrix {
		startingStationName := g.StationNames[startionStationId]
		for endingStationId, travelTime := range startingStation {
			if travelTime != 9999 {
				endingStationName := g.StationNames[endingStationId]
				fmt.Printf("Station %s to %s: %d minutes\n", startingStationName, endingStationName, travelTime)
			}
		}
	}
}

func (g *Graph) BuildTravelTimeMatrix() {
	// Run Floyd-Warshall to get all-pairs shortest path first for all stations
	travelTimeMatrix := make(map[StationId]map[StationId]int, 0)

	// for startingStation := range travelTimeMatrix {
	// 	travelTimeMatrix[startingStation] = make(map[StationId]int, 0)
	// }

	// initialise base cases, A-A, B-B travel time is 0 minutes
	for stationId := range g.Stations {
		travelTimeMatrix[stationId] = make(map[StationId]int, 0)
		// if _, exists := travelTimeMatrix[stationId]; !exists {
		// 	travelTimeMatrix[stationId] = make(map[StationId]int, 0)
		// }
		for adjacentStationId := range g.Stations {
			if existingRoute, exists := g.Routes[stationId][adjacentStationId]; exists {
				fmt.Println("got")
				travelTimeMatrix[stationId][adjacentStationId] = existingRoute.TravelTime
			} else {
				travelTimeMatrix[stationId][adjacentStationId] = 9999
			}
		}
		travelTimeMatrix[stationId][stationId] = 0
	}

	for kStation := range g.Stations {
		for iStation := range g.Stations {
			for jStation := range g.Stations {
				if travelTimeMatrix[iStation][jStation] > travelTimeMatrix[iStation][kStation]+travelTimeMatrix[kStation][jStation] {
					travelTimeMatrix[iStation][jStation] = travelTimeMatrix[iStation][kStation] + travelTimeMatrix[kStation][jStation]
				}
			}
		}
	}

	g.TravelTimeMatrix = travelTimeMatrix
}
