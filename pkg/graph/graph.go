package graph

import (
	"fmt"
	"strconv"
	"strings"
)

type Train struct {
	Name            string
	Capacity        int
	StartingStation *Station
}

type Package struct {
	Name            string
	Weight          int
	StartingStation *Station
	EndingStation   *Station
}

type Station struct {
	Name   string
	Routes []Route
}

type Route struct {
	Name       string
	TravelTime int
	Stations   [2]string
}

type Graph struct {
	Stations   map[string]*Station
	Routes     []Route
	Deliveries []Package
	Trains     []Train
}

func NewGraph(stationNames []string, rawRoutes []string, rawDeliveries []string, rawTrains []string) (*Graph, error) {
	stations := make(map[string]*Station, 0)
	for _, stationName := range stationNames {
		stations[stationName] = &Station{
			Name: stationName,
		}
	}

	routes := make([]Route, 0)
	for _, rawRoute := range rawRoutes {
		route := strings.Split(rawRoute, ",")
		routeName := route[0]
		fromStationName := route[1]
		toStationName := route[2]
		travelTime, err := strconv.Atoi(route[3])
		if err != nil {
			return nil, fmt.Errorf("Route %s is not in integer format", routeName)
		}

		routes = append(routes, Route{
			Name:       routeName,
			TravelTime: travelTime,
			Stations:   [2]string{fromStationName, toStationName},
		})
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

		fromStation := stations[fromStationName]
		toStation := stations[toStationName]

		deliveries = append(deliveries, Package{
			Name:            packageName,
			Weight:          weight,
			StartingStation: fromStation,
			EndingStation:   toStation,
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

		startingStation := stations[startingStationName]

		trains = append(trains, Train{
			Name:            trainName,
			Capacity:        capacity,
			StartingStation: startingStation,
		})
	}

	return &Graph{
		Stations:   stations,
		Routes:     routes,
		Deliveries: deliveries,
		Trains:     trains,
	}, nil
}
