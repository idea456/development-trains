package graph

import (
	"fmt"
	"slices"
)

type Train struct {
	Name             string
	Capacity         int
	TravelTime       int
	CurrentStationId StationId
	PackagesCarried  []Package
}

func (train *Train) AddPackage(delivery Package) {
	train.PackagesCarried = append(train.PackagesCarried, delivery)
	train.Capacity = train.Capacity - delivery.Weight
}

func (train *Train) UpdatePosition(newStationId StationId) {
	fmt.Println("updating train position to station", newStationId)
	train.CurrentStationId = newStationId
}

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

func (train *Train) HasPackagesToDeliver() bool {
	return len(train.PackagesCarried) > 0
}
