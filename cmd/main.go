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
			return nil, fmt.Errorf("expected %d stations to be in the input, found less stations than expected", count)
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
			return nil, fmt.Errorf("expected %d routes to be in the input, found less routes than expected", count)
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
			return nil, fmt.Errorf("expected %d packages to be in the input, found less routes than expected", count)
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
			return nil, fmt.Errorf("expected %d trains to be in the input, found less routes than expected", count)
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

func ScanInputFromPrompt() (*RawInput, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var rawStations, rawRoutes, rawPackages, rawTrains []string

	// Read stations
	fmt.Println("Enter the number of stations:")
	scanner.Scan()
	stationCount := 0
	fmt.Sscanf(scanner.Text(), "%d", &stationCount)
	fmt.Println("Enter the stations (one per line):")
	for i := 0; i < stationCount; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d stations, found less stations than expected", stationCount)
		}
		rawStations = append(rawStations, text)
	}

	// Read routes
	fmt.Println("Enter the number of routes:")
	scanner.Scan()
	routeCount := 0
	fmt.Sscanf(scanner.Text(), "%d", &routeCount)
	fmt.Println("Enter the routes (one per line, format: E1,A,B,10):")
	for i := 0; i < routeCount; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d routes, but found less routes than expected", routeCount)
		}
		rawRoutes = append(rawRoutes, text)
	}

	// Read packages
	fmt.Println("Enter the number of packages:")
	scanner.Scan()
	packageCount := 0
	fmt.Sscanf(scanner.Text(), "%d", &packageCount)
	fmt.Println("Enter the packages (one per line, format: K1,3,A,E):")
	for i := 0; i < packageCount; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d packages, but found less packages than expected", packageCount)
		}
		rawPackages = append(rawPackages, text)
	}

	// Read trains
	fmt.Println("Enter the number of trains:")
	scanner.Scan()
	trainCount := 0
	fmt.Sscanf(scanner.Text(), "%d", &trainCount)
	fmt.Println("Enter the trains (one per line, format: Q1,3,A):")
	for i := 0; i < trainCount; i++ {
		scanner.Scan()
		text := scanner.Text()
		if text == "" {
			return nil, fmt.Errorf("expected %d trains, but found less trains than expected", trainCount)
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
	prompt := flag.Bool("prompt", false, "Prompt for input instead")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	summary := flag.Bool("summary", false, "Enable summary output")
	flag.Parse()

	var rawInput *RawInput
	if *prompt {
		rawInputFromPrompt, err := ScanInputFromPrompt()
		if err != nil {
			slog.Error(fmt.Sprintf("unable to process input fil:e %v", err))
			os.Exit(1)
		}
		rawInput = rawInputFromPrompt
	} else {
		if *inputFilePath == "" {
			fmt.Println("Error: Input file path is required, e.g. ./tests/sample.txt")
			flag.PrintDefaults()
			os.Exit(1)
		}

		rawInputFile, err := ScanInputFile(*inputFilePath)
		if err != nil {
			slog.Error(fmt.Sprintf("unable to process input fil:e %v", err))
			os.Exit(1)
		}
		rawInput = rawInputFile
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
