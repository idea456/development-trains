package graph

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

type StationId = int
type PackageName = string

// // to represent Infinity in integer format
// const MaxUint = ^uint(0)
const MaxInt = 999999

type Train struct {
	Name              string
	Capacity          int
	StartingStationId StationId
	CurrentStationId  StationId
	PackagesCarried   []*Package
}

type Package struct {
	Name              PackageName
	Weight            int
	StartingStationId StationId
	EndingStationId   StationId
}

type Station struct {
	Id                StationId
	Name              string
	PotentialPackages map[PackageName]*Package
}

type Route struct {
	Name       string
	TravelTime int
}

type Move struct {
	StartedSeconds  int
	Train           Train
	StartingStation Station
	PickedPackages  []Package
	EndingStation   Station
	DroppedPackages Package
}

type Graph struct {
	Stations         map[StationId]*Station
	StationNames     map[StationId]string
	Routes           map[StationId]map[StationId]*Route
	Deliveries       []Package
	Trains           []Train
	TravelTimeMatrix map[StationId]map[StationId]int
	TravelPathMatrix map[StationId]map[StationId]StationId
	Moves            []Move
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
			CurrentStationId:  startingStationId,
			PackagesCarried:   make([]*Package, 0),
		})
	}

	return &Graph{
		Stations:     stations,
		StationNames: stationNamesMap,
		Routes:       routes,
		Deliveries:   deliveries,
		Trains:       trains,
		Moves:        make([]Move, 0),
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
			if travelTime != MaxInt {
				endingStationName := g.StationNames[endingStationId]
				fmt.Printf("Station %s to %s: %d minutes\n", startingStationName, endingStationName, travelTime)
			}
		}
	}
}

func (g *Graph) BuildTravelTimeMatrix() {
	// Run Floyd-Warshall to get all-pairs shortest path first for all stations
	travelTimeMatrix := make(map[StationId]map[StationId]int, 0)
	travelPathMatrix := make(map[StationId]map[StationId]StationId, 0)

	// initialise base cases, A-A, B-B travel time is 0 minutes
	for stationId := range g.Stations {
		travelTimeMatrix[stationId] = make(map[StationId]int, 0)
		travelPathMatrix[stationId] = make(map[StationId]StationId, 0)

		for adjacentStationId := range g.Stations {
			travelPathMatrix[stationId][adjacentStationId] = -1

			if existingRoute, exists := g.Routes[stationId][adjacentStationId]; exists {
				travelTimeMatrix[stationId][adjacentStationId] = existingRoute.TravelTime

				travelPathMatrix[stationId][adjacentStationId] = stationId
				if _, exists := travelPathMatrix[adjacentStationId]; !exists {
					travelPathMatrix[adjacentStationId] = make(map[StationId]StationId, 0)
				}
				travelPathMatrix[adjacentStationId][stationId] = adjacentStationId
			} else {
				travelTimeMatrix[stationId][adjacentStationId] = MaxInt
			}
		}
		travelTimeMatrix[stationId][stationId] = 0
	}

	for kStation := range g.Stations {
		for iStation := range g.Stations {
			for jStation := range g.Stations {
				if travelTimeMatrix[iStation][jStation] > travelTimeMatrix[iStation][kStation]+travelTimeMatrix[kStation][jStation] {
					travelTimeMatrix[iStation][jStation] = travelTimeMatrix[iStation][kStation] + travelTimeMatrix[kStation][jStation]
					travelTimeMatrix[jStation][iStation] = travelTimeMatrix[iStation][kStation] + travelTimeMatrix[kStation][jStation]

					travelPathMatrix[iStation][jStation] = travelPathMatrix[kStation][jStation]
					travelPathMatrix[jStation][iStation] = travelPathMatrix[jStation][kStation]
				}
			}
		}
	}

	g.TravelTimeMatrix = travelTimeMatrix
	g.TravelPathMatrix = travelPathMatrix
}

func (g *Graph) MoveToStation(train *Train, fromStationId StationId, toStationId StationId) {
	paths := make([]string, 0)

	start := fromStationId
	end := toStationId
	for start != end {
		paths = append(paths, g.StationNames[end])
		end = g.TravelPathMatrix[start][end]
	}
	paths = append(paths, g.StationNames[start])

	slices.Reverse(paths)
	fmt.Println("for train", train.Name)
	fmt.Println(paths)
}

func (g *Graph) Deliver() {
	undeliveredPackages := make([]Package, 0)
	undeliveredPackages = append(undeliveredPackages, g.Deliveries...)

	assignedTrains := make([]Train, 0)
	unassignedTrains := make([]Train, 0)
	unassignedTrains = append(unassignedTrains, g.Trains...)

	for len(assignedTrains) <= len(unassignedTrains) {
		// undeliveredPackage, _ := undeliveredPackages[0], undeliveredPackages[1:]
		train := unassignedTrains[0]

		var nearestPackage *Package
		var nearestPackageIndex int
		for i, undeliveredPackage := range undeliveredPackages {
			if undeliveredPackage.Weight > train.Capacity {
				continue
			}

			if nearestPackage == nil {
				nearestPackage = &undeliveredPackage
				continue
			}

			currentTravelTimeToPickupPackage := g.TravelTimeMatrix[train.CurrentStationId][undeliveredPackage.StartingStationId]
			nearestTravelTimeToPickupPackage := g.TravelTimeMatrix[nearestPackage.StartingStationId][train.CurrentStationId]
			if currentTravelTimeToPickupPackage < nearestTravelTimeToPickupPackage {
				nearestPackage = &undeliveredPackage
				nearestPackageIndex = i
			} else if currentTravelTimeToPickupPackage == nearestTravelTimeToPickupPackage {
				// NOTE: There are 2 nearest trains that can pickup the package
				// TODO: Handle this case
				nearestPackage = &undeliveredPackage
				nearestPackageIndex = i
			}
		}

		if nearestPackage == nil {
			// NOTE: this train cannot pick up anymore packages, might be too heavy or its already filled with pickups
			// TODO: Handle this case
			unassignedTrains = unassignedTrains[1:]
			continue
		}

		fmt.Println("delivering...", train.CurrentStationId, nearestPackage.StartingStationId)
		g.MoveToStation(&train, train.CurrentStationId, nearestPackage.StartingStationId)

		unassignedTrains[0].CurrentStationId = nearestPackage.StartingStationId
		unassignedTrains[0].Capacity = train.Capacity - nearestPackage.Weight
		unassignedTrains[0].PackagesCarried = append(unassignedTrains[0].PackagesCarried, nearestPackage)

		var isAdded bool
		for j := 0; j < len(assignedTrains); j++ {
			if assignedTrains[j].Name == train.Name {
				isAdded = true
				assignedTrains[j].CurrentStationId = nearestPackage.StartingStationId
				assignedTrains[j].Capacity = train.Capacity - nearestPackage.Weight
				assignedTrains[j].PackagesCarried = append(assignedTrains[j].PackagesCarried, nearestPackage)
				break
			}
		}
		if !isAdded {
			assignedTrains = append(assignedTrains, unassignedTrains[0])
		}

		undeliveredPackages = append(undeliveredPackages[:nearestPackageIndex], undeliveredPackages[nearestPackageIndex+1:]...)
	}

	// at this point, all packages should have been picked up by some trains
	// now we can deliver them to their destinations!
	// NOTE: But we could have some other pakcages that has not been assigned yet, so we need the trains to drop them first :(
	fmt.Println(assignedTrains)
}
