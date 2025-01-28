package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/idea456/development-trains/pkg/graph"
)

type RawInput struct {
	RawStations []string
	RawRoutes   []string
	RawPackages []string
	RawTrains   []string
}

func ScanInputFile(inputFilePath string) (*RawInput, error) {
	file, err := os.ReadFile(inputFilePath)
	if err != nil {
		return nil, err
	}
	var rawStations, rawRoutes, rawPackages, rawTrains []string

	scanner := bufio.NewScanner(strings.NewReader(string(file)))
	scanner.Scan()
	count := 0
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d stations to be in the input, found less stations that expected", count)
		}
		rawStations = append(rawStations, text)
	}

	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d routes to be in the input, found less routes that expected", count)
		}
		rawRoutes = append(rawRoutes, text)
	}

	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d packages to be in the input, found less routes that expected", count)
		}
		rawPackages = append(rawPackages, text)
	}

	scanner.Scan()
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "%d", &count)
	for i := 0; i < count; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d trains to be in the input, found less routes that expected", count)
		}
		rawTrains = append(rawTrains, text)
	}

	return &RawInput{
		RawStations: rawStations,
		RawRoutes:   rawRoutes,
		RawPackages: rawPackages,
		RawTrains:   rawTrains,
	}, nil
}

func main() {
	inputFilePath := flag.String("i", "", "Path to the input file")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	summary := flag.Bool("summary", false, "Enable summary output")
	flag.Parse()

	if *inputFilePath == "" {
		fmt.Println("Error: Input file path is required, e.g. ./tests/sample.txt")
		flag.PrintDefaults()
		os.Exit(1)
	}

	rawInput, err := ScanInputFile(*inputFilePath)
	if err != nil {
		slog.Error(fmt.Sprintf("unable to process input fil:e %v", err))
		os.Exit(1)
	}

	g, err := graph.NewGraph(rawInput.RawStations, rawInput.RawRoutes, rawInput.RawPackages, rawInput.RawTrains)
	if err != nil {
		slog.Error(fmt.Sprintf("There was an issue in building the graph: %v", err))
		os.Exit(1)
	}

	g.BuildTravelTimeMatrix()
	if err := g.Deliver(); err != nil {
		slog.Error(fmt.Sprintf("unable to deliver all packages: %v", err))
		os.Exit(1)
	}

	printer := graph.NewPrinter(g.Moves, g.StationNames, g.TravelTimeMatrix)
	if *verbose {
		printer.PrintMovesVerbose()
	} else {
		printer.PrintMoves()
	}

	if *summary {
		printer.PrintSummary()
	}
}
