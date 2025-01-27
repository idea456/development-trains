package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/idea456/development-trains/pkg/graph"
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("Please specify the input file path")
		os.Exit(1)
	}
	inputFilePath := os.Args[1]

	file, err := os.ReadFile(inputFilePath)
	if err != nil {
		slog.Error(fmt.Sprintf("There was an issue in reading the file: %v", err))
		os.Exit(1)
	}

	scanner := bufio.NewScanner(strings.NewReader(string(file)))

	var rawStations, rawRoutes, rawPackages, rawTrains []string

	scanner.Scan()
	count := 0
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		rawStations = append(rawStations, scanner.Text())
	}

	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		rawRoutes = append(rawRoutes, scanner.Text())
	}

	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		rawPackages = append(rawPackages, scanner.Text())
	}

	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		rawTrains = append(rawTrains, scanner.Text())
	}

	fmt.Println(rawStations, rawRoutes, rawPackages, rawTrains)
	g, err := graph.NewGraph(rawStations, rawRoutes, rawPackages, rawTrains)
	if err != nil {
		slog.Error(fmt.Sprintf("There was an issue in building the graph: %v", err))
		os.Exit(1)
	}

	g.BuildTravelTimeMatrix()
	// graph.TrackCommonDestinationPackages()
	// graph.PrintShortestRoutes()
	if err := g.Deliver(); err != nil {
		slog.Error(fmt.Sprintf("unable to deliver all packages: %v", err))
		os.Exit(1)
	}
	// graph.PrintMoves()

	printer := graph.NewPrinter(g.Moves, g.StationNames)
	// printer.PrintMoves()
	printer.PrintMovesVerbose()
	printer.PrintSummary()
}
