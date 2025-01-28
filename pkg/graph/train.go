package graph

import (
	"slices"
)

// Train represents the train which tracks its capacity, travel time, current location and packages its holding
type Train struct {
	Name             string
	Capacity         int
	TravelTime       int
	CurrentStationId StationId
	PackagesCarried  []Package
}

// Adds a package to the train
func (train *Train) AddPackage(delivery Package) {
	train.PackagesCarried = append(train.PackagesCarried, delivery)
	train.Capacity = train.Capacity - delivery.Weight
}

// Updates the train's current station
func (train *Train) UpdatePosition(newStationId StationId) {
	train.CurrentStationId = newStationId
}

// Drops a package if possible at its current station location
// If it cannot drop any packages, an empty slice is returned
func (train *Train) DropPackages() []Package {
	droppedPackages := make([]Package, 0)
	for _, carriedPackage := range train.PackagesCarried {
		if train.CurrentStationId == carriedPackage.EndingStationId {
			droppedPackages = append(droppedPackages, carriedPackage)
		}
	}
	if len(droppedPackages) > 0 {
		train.RemovePackages(droppedPackages)
	}
	return droppedPackages
}

// Removes the specified packages from its holding
func (train *Train) RemovePackages(droppedPackages []Package) {
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

// Checks if the train still has leftover packages its holding that needs to be delivered
func (train *Train) HasPackagesToDeliver() bool {
	return len(train.PackagesCarried) > 0
}
