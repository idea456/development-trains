package graph

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"
)

type StationId = int
type PackageName = string

// // to represent Infinity in integer format
const MaxUint = ^uint(0)

// const MaxInt = int(MaxUint >> 1)
const MaxInt = 9999999999

type Train struct {
	Name              string
	Capacity          int
	TravelTime        int
	StartingStationId StationId
	CurrentStationId  StationId
	PackagesCarried   []Package
}

func (train *Train) RemoveDroppedPackages(droppedPackages []Package) {
	packages := make([]Package, 0)
	totalCapacityRemoved := 0
	for _, carriedPackage := range train.PackagesCarried {
		isDroppedPackage := slices.ContainsFunc(droppedPackages, func(droppedPackage Package) bool {
			return droppedPackage.Name == carriedPackage.Name
		})
		if !isDroppedPackage {
			packages = append(packages, carriedPackage)
		} else {
			totalCapacityRemoved += carriedPackage.Weight
		}
	}
	train.PackagesCarried = packages
	train.Capacity = train.Capacity + totalCapacityRemoved
}

func (train *Train) AddPackage(delivery Package) {
	train.PackagesCarried = append(train.PackagesCarried, delivery)
	train.Capacity = train.Capacity - delivery.Weight
}

func (train *Train) UpdatePosition(newStationId StationId) {
	train.CurrentStationId = newStationId
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
	InitialPackages   map[PackageName]*Package
	PotentialPackages map[PackageName]*Package
}

type Route struct {
	Name       string
	TravelTime int
}

type Move struct {
	TimeTaken       int
	Train           Train
	StartingStation Station
	EndingStation   Station
	PackagesCarried []Package
	PackagesDropped []Package
}

type Graph struct {
	Stations         map[StationId]*Station
	StationNames     map[StationId]string
	Routes           map[StationId]map[StationId]*Route
	Deliveries       []Package
	Trains           map[string]*Train
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
			Id:                i,
			Name:              stationName,
			InitialPackages:   make(map[PackageName]*Package, 0),
			PotentialPackages: make(map[PackageName]*Package, 0),
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
			Name:              trainName,
			Capacity:          capacity,
			StartingStationId: startingStationId,
			CurrentStationId:  startingStationId,
			PackagesCarried:   make([]Package, 0),
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
			// travelPathMatrix[stationId][adjacentStationId] = -1

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

	stationIds := make([]StationId, 0, len(g.Stations))
	for stationId := range g.Stations {
		stationIds = append(stationIds, stationId)
	}
	slices.Sort(stationIds)

	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm
	for kStation := range stationIds {
		for iStation := range stationIds {
			for jStation := range stationIds {
				if travelTimeMatrix[iStation][jStation] > travelTimeMatrix[iStation][kStation]+travelTimeMatrix[kStation][jStation] {
					travelTimeMatrix[iStation][jStation] = travelTimeMatrix[iStation][kStation] + travelTimeMatrix[kStation][jStation]
					travelTimeMatrix[jStation][iStation] = travelTimeMatrix[jStation][kStation] + travelTimeMatrix[kStation][iStation]

					// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Path_reconstruction
					travelPathMatrix[iStation][jStation] = travelPathMatrix[kStation][jStation]
					// travelPathMatrix[jStation][iStation] = travelPathMatrix[jStation][kStation]

				}
			}
		}
	}

	for i := range stationIds {
		for j := range travelTimeMatrix[i] {
			fmt.Printf("Station %s to station %s: %d minutes\n", g.StationNames[i], g.StationNames[j], travelTimeMatrix[i][j])
		}
		fmt.Println()
	}
	fmt.Println()

	g.TravelTimeMatrix = travelTimeMatrix
	g.TravelPathMatrix = travelPathMatrix
}

func (g *Graph) TrackCommonDestinationPackages() {
	// trace the shortest path for each package
	// while tracing, if we find another package that is also going to the same destination, we store the reference of that in potential packages
	// TODO: This is just one direction, source -> destination, optimize this later
	for _, delivery := range g.Deliveries {
		start := delivery.StartingStationId
		end := delivery.EndingStationId
		commonPackages := make(map[PackageName]*Package, 0)
		fmt.Printf("From station %s to station %s\n", g.StationNames[start], g.StationNames[end])
		for start != end {
			// if this station has packages placed on it
			// and that package is also going to the same destination
			for _, initialPackage := range g.Stations[end].InitialPackages {
				if initialPackage.EndingStationId == delivery.EndingStationId {
					commonPackages[initialPackage.Name] = initialPackage
				}
			}
			// g.Stations[end].PotentialPackages = commonPackages
			maps.Copy(g.Stations[end].PotentialPackages, commonPackages)
			end = g.TravelPathMatrix[start][end]
		}
		for _, initialPackage := range g.Stations[end].InitialPackages {
			if initialPackage.EndingStationId == delivery.EndingStationId {
				g.Stations[end].PotentialPackages[initialPackage.Name] = initialPackage
			}
		}
		maps.Copy(g.Stations[end].PotentialPackages, commonPackages)
	}

	for _, station := range g.Stations {
		if len(station.PotentialPackages) == 0 {
			continue
		}
		fmt.Printf("Common destination packages for station %s:\n", station.Name)
		for _, potentialPackage := range station.PotentialPackages {
			fmt.Printf("Package %s (%d kg) with destination to station %s\n", potentialPackage.Name, potentialPackage.Weight, g.StationNames[potentialPackage.EndingStationId])
		}
		fmt.Println()
	}
}

func (g *Graph) MoveToPickupPackage(train Train, nearestPackage Package) {
	paths := make([]StationId, 0)
	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Path_reconstruction
	start := train.CurrentStationId
	end := nearestPackage.StartingStationId

	paths = append(paths, end)
	for start != end {
		end = g.TravelPathMatrix[start][end]

		paths = append(paths, end)
	}

	slices.Reverse(paths)

	fmt.Printf("moving from %s station to %s station\n", g.StationNames[start], g.StationNames[end])
	fmt.Printf("%+v\n", paths)

	moves := make([]Move, 0)
	currentTravelTime := g.Trains[train.Name].TravelTime
	for i := 0; i < len(paths)-1; i++ {
		// fmt.Printf("%s station to %s station => %d minutes\n", move.StartingStation.Name, move.EndingStation.Name, move.TravelTime)
		currentStationId := paths[i]
		nextStationId := paths[i+1]

		moves = append(moves, Move{
			TimeTaken:       currentTravelTime,
			Train:           train,
			StartingStation: *g.Stations[currentStationId],
			EndingStation:   *g.Stations[nextStationId],
			PackagesCarried: g.Trains[train.Name].PackagesCarried,
		})
		currentTravelTime += g.TravelTimeMatrix[currentStationId][nextStationId]
	}
	g.Trains[train.Name].TravelTime = currentTravelTime
	g.Trains[train.Name].UpdatePosition(nearestPackage.StartingStationId)
	g.Trains[train.Name].AddPackage(nearestPackage)
	// moves[len(moves)-1].PackagesCarried = append(g.Trains[train.Name].PackagesCarried, nearestPackage)
	g.Moves = append(g.Moves, moves...)
}

func (g *Graph) MoveToDropPackage(train Train, packages []Package, destinationStationId int) {
	paths := make([]StationId, 0)
	fmt.Printf("dropping from %s station to %s station\n", g.StationNames[train.CurrentStationId], g.StationNames[destinationStationId])
	// https://en.wikipedia.org/wiki/Floyd%E2%80%93Warshall_algorithm#Path_reconstruction
	start := g.Trains[train.Name].CurrentStationId
	end := destinationStationId

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
		// fmt.Printf("%s station to %s station => %d minutes\n", move.StartingStation.Name, move.EndingStation.Name, move.TravelTime)
		currentStationId := paths[i]
		nextStationId := paths[i+1]
		moves = append(moves, Move{
			TimeTaken:       currentTravelTime,
			Train:           train,
			StartingStation: *g.Stations[currentStationId],
			EndingStation:   *g.Stations[nextStationId],
			PackagesCarried: g.Trains[train.Name].PackagesCarried,
		})
		currentTravelTime += g.TravelTimeMatrix[currentStationId][nextStationId]
	}

	g.Trains[train.Name].TravelTime = currentTravelTime
	g.Trains[train.Name].UpdatePosition(destinationStationId)

	fmt.Printf("dropped packages: %+v\n", packages)
	g.Trains[train.Name].RemoveDroppedPackages(packages)
	// have to update again, since RemoveDroppedPackages will filter out some of the carried packages that are dropped
	moves[len(moves)-1].PackagesCarried = g.Trains[train.Name].PackagesCarried
	moves[len(moves)-1].PackagesDropped = packages
	g.Moves = append(g.Moves, moves...)
}

func (g *Graph) PrintMoves() {
	// W=0, T=Q1, N1=B, P1=[], N2=A, P2=[]
	slices.SortStableFunc(g.Moves, func(a Move, b Move) int {
		return strings.Compare(a.Train.Name, b.Train.Name)
	})
	for _, move := range g.Moves {
		packagesCarriedNames := make([]string, 0)
		for _, packageCarried := range move.PackagesCarried {
			packagesCarriedNames = append(packagesCarriedNames, packageCarried.Name)
		}
		packageCarriedStr := fmt.Sprintf("[%s]", strings.Join(packagesCarriedNames, ","))

		packageDroppedNames := make([]string, 0)
		for _, packageDropped := range move.PackagesDropped {
			packageDroppedNames = append(packageDroppedNames, packageDropped.Name)
		}
		packageDroppedStr := fmt.Sprintf("[%s]", strings.Join(packageDroppedNames, ","))

		fmt.Printf("W=%d, T=%s, N1=%s, P1=%s, N2=%s, P2=%s\n", move.TimeTaken, move.Train.Name, move.StartingStation.Name, packageCarriedStr, move.EndingStation.Name, packageDroppedStr)
	}
}

func (g *Graph) PrintMovesVerbose() {
	// sort by train to easily track moves per train
	slices.SortStableFunc(g.Moves, func(a Move, b Move) int {
		return strings.Compare(a.Train.Name, b.Train.Name)
	})
	for _, move := range g.Moves {
		fmt.Printf("[%d minutes] Train %s moving from station %s to station %s\n", move.TimeTaken, move.Train.Name, move.StartingStation.Name, move.EndingStation.Name)
		if len(move.PackagesCarried) > 0 {
			fmt.Println("Carried packages:")
			for _, carriedPackage := range move.PackagesCarried {
				fmt.Printf("	- %s package with weight %d heading to %s station\n", carriedPackage.Name, carriedPackage.Weight, g.StationNames[carriedPackage.EndingStationId])
			}
		}
		if len(move.PackagesDropped) > 0 {
			fmt.Println("Droppped packages:")
			for _, dropppedPackage := range move.PackagesDropped {
				fmt.Printf("	- %s package with weight %d at %s station\n", dropppedPackage.Name, dropppedPackage.Weight, g.StationNames[dropppedPackage.EndingStationId])
			}
		}
		fmt.Println()
	}
}

func (g *Graph) Deliver() {
	deliveredPackages := make([]Package, 0)
	undeliveredPackages := make([]Package, 0)
	undeliveredPackages = append(undeliveredPackages, g.Deliveries...)

	unassignedTrains := make([]Train, 0)

	// go returns the keys in random order, to make it determinstic we sort it ahead of time
	trainNames := make([]string, 0, len(g.Trains))
	for trainName := range g.Trains {
		trainNames = append(trainNames, trainName)
	}
	slices.Sort(trainNames)
	for _, trainName := range trainNames {
		unassignedTrains = append(unassignedTrains, *g.Trains[trainName])
	}

	// if we haven't delivered all the packages yet
	for len(undeliveredPackages) > 0 {
		// NOTE: PICKUP phase
		// assign all the trains (if possible) first
		// TODO: Handle a case where there are unassigned trains, but there are no more packages to pick up
		for len(unassignedTrains) > 0 {
			// undeliveredPackage, _ := undeliveredPackages[0], undeliveredPackages[1:]
			// TODO: Naivse solution, just take the first train, improve this
			train := g.Trains[unassignedTrains[0].Name]

			var nearestPackage *Package
			var nearestPackageIndex int
			// check each package, and find a nearest package (TODO: AND optimal package using potential common package destinations) to assign to the current train
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
				// TODO: Handle this case (this is done already?)
				unassignedTrains = unassignedTrains[1:]
				continue
			}

			g.MoveToPickupPackage(*train, *nearestPackage)

			undeliveredPackages = append(undeliveredPackages[:nearestPackageIndex], undeliveredPackages[nearestPackageIndex+1:]...)
		}

		// at this point, all packages should have been picked up by some trains OR all trains have been packed
		// we need the trains to deliver the packages to their destinations first, then the train can pick up more packages if needed
		// NOTE: But we could have some other pakcages that has not been assigned yet, so we need the trains to drop them first :(
		// NOTE: DROPOFF phase
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
				g.MoveToDropPackage(assignedTrain, carriedPackages, packageDestinationStationId)

				deliveredPackages = append(deliveredPackages, carriedPackages...)
			}
			assignedTrain.PackagesCarried = []Package{}
		}

	}
}
