package graph

import (
	"container/heap"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

const MaxInt = 999999999999

type StationId = int
type StationName = string
type PackageName = string

// Package struct represents the package to be delivered
type Package struct {
	Name              PackageName
	Weight            int
	StartingStationId StationId
	EndingStationId   StationId
	DeliveredAt       int
}

// Station represents a station node
type Station struct {
	Id              StationId
	Name            string
	InitialPackages map[PackageName]*Package
}

// Route represents the weighted edges between stations with travel time
type Route struct {
	Name       string
	TravelTime int
}

// Move represents a train's movement and pickup/dropoff actions
type Move struct {
	TimeTaken       int
	Train           Train
	StartingStation Station
	EndingStation   Station
	PackagesCarried []Package
	PackagesDropped []Package
}

// Graph represents the transit network
// Edges are represented with a 'hashmap' adjancency matrix to optimise space for non-existing edges
type Graph struct {
	Stations         map[StationId]*Station
	StationNames     map[StationId]StationName
	Routes           map[StationId]map[StationId]*Route
	Deliveries       []Package
	Trains           map[string]*Train
	TravelTimeMatrix map[StationId]map[StationId]int       // Stores shortest travel time between all stations
	TravelPathMatrix map[StationId]map[StationId]StationId // Stores references of previous nodes to backtrack shortest path
	Moves            []Move                                // Tracks list of moves performed by the trains
}

// Creates a new Graph instance, receives the raw input strings of the stations, routes, deliveries and trains
func NewGraph(stationNames []StationName, rawRoutes []string, rawDeliveries []string, rawTrains []string) (*Graph, error) {
	stations := make(map[int]*Station, 0)
	stationNamesToIdMap := make(map[string]int, 0)
	stationNamesMap := make(map[int]string, 0)
	for i, stationName := range stationNames {
		stations[i] = &Station{
			Id:              i,
			Name:            stationName,
			InitialPackages: make(map[PackageName]*Package, 0),
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

		newDelivery := Package{
			Name:              packageName,
			Weight:            weight,
			StartingStationId: fromStationId,
			EndingStationId:   toStationId,
		}

		// keep track of which stations is initially holding the packages
		stations[fromStationId].InitialPackages[newDelivery.Name] = &newDelivery

		deliveries = append(deliveries, newDelivery)
	}

	trains := make(map[string]*Train, 0)
	for _, rawTrain := range rawTrains {
		train := strings.Split(rawTrain, ",")
		trainName := train[0]
		capacity, err := strconv.Atoi(train[1])
		if err != nil {
			return nil, fmt.Errorf("Train %s capacity is not an integer", trainName)
		}
		startingStationName := train[2]

		startingStationId := stationNamesToIdMap[startingStationName]

		trains[trainName] = &Train{
			Name:             trainName,
			Capacity:         capacity,
			CurrentStationId: startingStationId,
			PackagesCarried:  make([]Package, 0),
		}

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

// BuildTravelTimeMatrix creates a distance matrix for every shortest path between every stations using Floyd-Warshall algorithm
// Allows for O(1) lookup for every shortest-path between stations, but costs O(V^3) preprocessing time
func (g *Graph) BuildTravelTimeMatrix() {
	// Run Floyd-Warshall to get all-pairs shortest path first for all stations
	travelTimeMatrix := make(map[StationId]map[StationId]int, 0)
	// Store references of the previous paths to backtrack and reconstruct the shortest path
	travelPathMatrix := make(map[StationId]map[StationId]StationId, 0)

	stationIds := make([]StationId, 0)
	for stationId := range g.Stations {
		stationIds = append(stationIds, stationId)
	}
	slices.Sort(stationIds)
	// initialise base cases, A-A, B-B travel time is 0 minutes
	for stationId := range stationIds {
		travelTimeMatrix[stationId] = make(map[StationId]int, 0)
		travelPathMatrix[stationId] = make(map[StationId]StationId, 0)

		// Initialise base cases, allocate memory and initial travel times if connection exists
		for adjacentStationId := range stationIds {
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
		travelPathMatrix[stationId][stationId] = stationId
	}

	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Pseudocode
	for kStation := range stationIds {
		for iStation := range stationIds {
			for jStation := range stationIds {
				if travelTimeMatrix[iStation][jStation] > travelTimeMatrix[iStation][kStation]+travelTimeMatrix[kStation][jStation] {
					travelTimeMatrix[iStation][jStation] = travelTimeMatrix[iStation][kStation] + travelTimeMatrix[kStation][jStation]
					travelTimeMatrix[jStation][iStation] = travelTimeMatrix[jStation][kStation] + travelTimeMatrix[kStation][iStation]

					// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Path_reconstruction
					travelPathMatrix[iStation][jStation] = travelPathMatrix[kStation][jStation]
					travelPathMatrix[jStation][iStation] = travelPathMatrix[kStation][iStation]

				}
			}
		}
	}

	g.TravelTimeMatrix = travelTimeMatrix
	g.TravelPathMatrix = travelPathMatrix
}

// GetShortestPath returns the backtracked shortest path between 2 stations
// Time complexity: O(E) where E is the number of routes
func (g *Graph) GetShortestPath(startingStationId StationId, endingStationId StationId) []StationId {
	paths := make([]StationId, 0)
	start := startingStationId
	end := endingStationId

	// Reconstruct the shortest path and tracks it to the Moves slice
	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Path_reconstruction
	paths = append(paths, end)
	for start != end {
		end = g.TravelPathMatrix[start][end]

		paths = append(paths, end)
	}
	slices.Reverse(paths)
	return paths
}

// MoveToPickupPackage moves a train to pick up a package using the shortest path and updates its location and capacity
// Tracks the move and adds it to the Moves slice
func (g *Graph) MoveToPickupPackage(train Train, nearestPackage Package) {
	// CASE: If the package to pickup is already at the train's current location
	if train.CurrentStationId == nearestPackage.StartingStationId {
		// Add the package and no need to update the time since no time is taken to pickup packages
		g.Trains[train.Name].AddPackage(nearestPackage)
		// CASE: if we are also already in the same station that we can drop off the package
		droppedPackages := g.Trains[train.Name].DropPackages()
		g.Moves = append(g.Moves, Move{
			TimeTaken:       g.Trains[train.Name].TravelTime, // no time taken to pickup package since the train is already there
			Train:           train,
			StartingStation: *g.Stations[train.CurrentStationId],
			EndingStation:   *g.Stations[train.CurrentStationId],
			PackagesCarried: g.Trains[train.Name].PackagesCarried,
			PackagesDropped: droppedPackages,
		})
		return
	}

	// get the list of shortest path and adds it as moves
	paths := g.GetShortestPath(train.CurrentStationId, nearestPackage.StartingStationId)
	moves := make([]Move, 0)
	currentTravelTime := g.Trains[train.Name].TravelTime
	for i := 0; i < len(paths)-1; i++ {
		currentStationId := paths[i]
		nextStationId := paths[i+1]

		// CASE: if along the way, we passed by a station that we can drop by packages
		g.Trains[train.Name].UpdatePosition(currentStationId)
		droppedPackages := g.Trains[train.Name].DropPackages()

		moves = append(moves, Move{
			TimeTaken:       currentTravelTime,
			Train:           train,
			StartingStation: *g.Stations[currentStationId],
			EndingStation:   *g.Stations[nextStationId],
			PackagesCarried: g.Trains[train.Name].PackagesCarried,
			PackagesDropped: droppedPackages,
		})
		currentTravelTime += g.TravelTimeMatrix[currentStationId][nextStationId]
	}
	g.Trains[train.Name].TravelTime = currentTravelTime
	g.Trains[train.Name].UpdatePosition(nearestPackage.StartingStationId)
	g.Trains[train.Name].AddPackage(nearestPackage)
	g.Moves = append(g.Moves, moves...)
}

/*
MoveToDropPackage drops a package to its destination station using the specified train
Tracks the move and adds it to the Moves slice
*/
func (g *Graph) MoveToDropPackage(trainName string, packages []Package, destinationStationId int) {
	train := g.Trains[trainName]
	// CASE: If the train is alerady at the drop station
	// CASE: If the package to pickup is already at the train's current location
	if train.CurrentStationId == destinationStationId {
		// Add the package and no need to update the time since no time is taken to pickup packages
		g.Trains[train.Name].RemovePackages(packages)
		g.Moves = append(g.Moves, Move{
			TimeTaken:       g.Trains[train.Name].TravelTime, // no time taken to dropoff package since the train is already there
			Train:           *train,
			StartingStation: *g.Stations[train.CurrentStationId],
			EndingStation:   *g.Stations[train.CurrentStationId],
			PackagesCarried: g.Trains[train.Name].PackagesCarried,
			PackagesDropped: packages,
		})
		return
	}
	paths := make([]StationId, 0)
	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Path_reconstruction
	start := g.Trains[trainName].CurrentStationId
	end := destinationStationId

	// Reconstruct the shortest path and tracks it to the Moves slice
	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Path_reconstruction
	paths = append(paths, end)
	for start != end {
		end = g.TravelPathMatrix[start][end]
		paths = append(paths, end)
	}
	slices.Reverse(paths)

	totalDroppedWeight := 0
	for _, delivery := range packages {
		totalDroppedWeight += delivery.Weight
	}

	moves := make([]Move, 0)
	currentTravelTime := g.Trains[train.Name].TravelTime
	for i := 0; i < len(paths)-1; i++ {
		currentStationId := paths[i]
		nextStationId := paths[i+1]
		moves = append(moves, Move{
			TimeTaken:       currentTravelTime,
			Train:           *train,
			StartingStation: *g.Stations[currentStationId],
			EndingStation:   *g.Stations[nextStationId],
			PackagesCarried: g.Trains[train.Name].PackagesCarried,
		})
		currentTravelTime += g.TravelTimeMatrix[currentStationId][nextStationId]
	}

	g.Trains[train.Name].TravelTime = currentTravelTime
	g.Trains[train.Name].UpdatePosition(destinationStationId)
	g.Trains[train.Name].RemovePackages(packages)
	// have to update again, since RemoveDroppedPackages will filter out some of the carried packages that are dropped
	moves[len(moves)-1].PackagesCarried = g.Trains[train.Name].PackagesCarried
	moves[len(moves)-1].PackagesDropped = packages
	g.Moves = append(g.Moves, moves...)
}

/*
Deliver is the main entry function for delivering all packages to their destinations
It manages assignment for each train (using a max heap) a set of packages based on the package's weight and distance and attempts to deliver all packages to their destinations
It has 2 phases: the pickup phase and dropoff phase
Pickup phase assigns each train a set of packages to pick up
Dropoff phase directs each train to dropoff its set of packages to their destinations before it can pick up again
*/
func (g *Graph) Deliver() error {
	undeliveredPackages := make([]Package, 0)
	undeliveredPackages = append(undeliveredPackages, g.Deliveries...)

	// trainsQueue is a max heap which prioritizes trains with bigger capacity
	trainsQueue := &TrainsQueue{}
	heap.Init(trainsQueue)

	// go returns the keys in random order, to make it determinstic we sort it ahead of time
	trainNames := make([]string, 0, len(g.Trains))
	for trainName := range g.Trains {
		trainNames = append(trainNames, trainName)
	}
	slices.Sort(trainNames)
	for _, trainName := range trainNames {
		heap.Push(trainsQueue, *g.Trains[trainName])
	}

	// if we haven't delivered all the packages yet
	for len(undeliveredPackages) > 0 {
		// NOTE: PICKUP phase
		// assign all the trains (if possible) first
		for len(*trainsQueue) > 0 {
			if len(undeliveredPackages) == 0 {
				break
			}
			assignableTrain := heap.Pop(trainsQueue).(Train)
			train := g.Trains[assignableTrain.Name]

			// TODO: Can use a custom heap here
			slices.SortFunc(undeliveredPackages, func(packageX, packageY Package) int {
				// first sort by their package pickup distance from the train
				packageXDistanceToTrain := g.TravelTimeMatrix[train.CurrentStationId][packageX.StartingStationId]
				packageYDistanceToTrain := g.TravelTimeMatrix[train.CurrentStationId][packageY.StartingStationId]
				if packageXDistanceToTrain != packageYDistanceToTrain {
					return packageXDistanceToTrain - packageYDistanceToTrain
				} else {
					// otherwise, sort them by their weight, higher weight goes up first
					if packageX.Weight < packageY.Weight {
						return 1
					} else {
						return -1
					}
				}

			})
			nearestPackage := undeliveredPackages[0]

			if nearestPackage.Weight > g.Trains[train.Name].Capacity {
				// NOTE: this train cannot pick up anymore packages, package might be too heavy or the train is already filled with packages
				continue
			} else {
				heap.Push(trainsQueue, *g.Trains[train.Name])
			}

			g.MoveToPickupPackage(*train, nearestPackage)
			// this package has been picked up and can be delivered, update the undeliveredPackages
			undeliveredPackages = undeliveredPackages[1:]

			// if this train can still pick up more packages, we can be greedy and try to take up more packages
			// NOTE: This solves the issue of assigning packages with common destinations since undeliveredPackages is already sorted by destinations
			// TODO: Find a way to optimize this, this is a very greedy solution, what if there are other trains that can pick this up?
			for g.Trains[train.Name].Capacity > 0 && len(undeliveredPackages) > 0 {
				nearestPackage = undeliveredPackages[0]
				if nearestPackage.Weight > g.Trains[train.Name].Capacity {
					break
				}

				g.MoveToPickupPackage(*train, nearestPackage)
				undeliveredPackages = undeliveredPackages[1:]
			}
		}

		// at this point, all packages should have been picked up by some trains OR all trains have been packed
		// we need the trains to deliver the packages to their destinations first, then the train can pick up more packages if needed
		// NOTE: But we could have some other pakcages that has not been assigned yet, so we need the trains to drop them first :(
		// DROPOFF phase
		assignedTrains := make([]Train, 0)
		trainNames := make([]string, 0, len(g.Trains))
		for trainName := range g.Trains {
			trainNames = append(trainNames, trainName)
		}
		slices.Sort(trainNames)
		for _, trainName := range trainNames {
			train := g.Trains[trainName]
			if len(train.PackagesCarried) > 0 {
				assignedTrains = append(assignedTrains, *train)
			}
		}

		for len(assignedTrains) > 0 {
			assignedTrain := assignedTrains[0]
			assignedTrains = assignedTrains[1:]

			// track common destination packages
			packagesByDestinationMap := make(map[StationId][]Package, 0)
			for _, packageCarried := range assignedTrain.PackagesCarried {
				if _, exists := packagesByDestinationMap[packageCarried.EndingStationId]; !exists {
					packagesByDestinationMap[packageCarried.EndingStationId] = make([]Package, 0)
				}
				packagesByDestinationMap[packageCarried.EndingStationId] = append(packagesByDestinationMap[packageCarried.EndingStationId], packageCarried)
			}

			// for each package to be delivered for this train, choose the package that can be delivered earliest (use the matrix)
			for packageDestinationStationId, carriedPackages := range packagesByDestinationMap {
				g.MoveToDropPackage(assignedTrain.Name, carriedPackages, packageDestinationStationId)
			}

			assignedTrain.PackagesCarried = []Package{}
			// CASE: If there's still packages to pick up after all the trains have been assigned
			if !assignedTrain.HasPackagesToDeliver() {
				heap.Push(trainsQueue, assignedTrain)
			}
		}

		// CASE: There are still packages to deliver, but no trains can deliver them
		// Because they might not have enough capacity
		if len(*trainsQueue) == 0 && len(undeliveredPackages) > 0 {
			return fmt.Errorf("there are still packages to deliver, but no trains can deliver them :(")
		}
	}

	return nil
}
