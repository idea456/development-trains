package graph

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"text/tabwriter"
)

// Printer represents a helper struct to print out moves and information for each train's moves and packages
type Printer struct {
	Moves            []Move
	StationNames     map[StationId]StationName
	TravelTimeMatrix map[StationId]map[StationId]int
}

func NewPrinter(moves []Move, stationNames map[StationId]string, travelTimeMatrix map[StationId]map[StationId]int) *Printer {
	return &Printer{
		Moves:            moves,
		StationNames:     stationNames,
		TravelTimeMatrix: travelTimeMatrix,
	}
}

// Prints out the list of moves as specified by assignment requirements in the format of:
// W=0, T=Q1, N1=B, P1=[], N2=A, P2=[]
func (printer *Printer) PrintMoves() {
	slices.SortStableFunc(printer.Moves, func(a Move, b Move) int {
		return strings.Compare(a.Train.Name, b.Train.Name)
	})
	for _, move := range printer.Moves {
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
	fmt.Println()
}

// Prints out the list of moves each train takes in a more detailed format
func (printer *Printer) PrintMovesVerbose() {
	// sort by train to easily track moves per train
	slices.SortStableFunc(printer.Moves, func(a Move, b Move) int {
		return strings.Compare(a.Train.Name, b.Train.Name)
	})
	for _, move := range printer.Moves {
		fmt.Printf("[%d minutes] Train %s moving from station %s to station %s\n", move.TimeTaken, move.Train.Name, move.StartingStation.Name, move.EndingStation.Name)
		if len(move.PackagesCarried) > 0 {
			fmt.Println("Carried packages:")
			for _, carriedPackage := range move.PackagesCarried {
				fmt.Printf("	- %s package with weight %d heading to %s station\n", carriedPackage.Name, carriedPackage.Weight, printer.StationNames[carriedPackage.EndingStationId])
			}
		}
		if len(move.PackagesDropped) > 0 {
			fmt.Println("Dropped packages:")
			for _, dropppedPackage := range move.PackagesDropped {
				fmt.Printf("	- %s package with weight %d at %s station\n", dropppedPackage.Name, dropppedPackage.Weight, printer.StationNames[dropppedPackage.EndingStationId])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// Prints an overall summary for each package's delivery time as well as which train delivered it
func (printer *Printer) PrintSummary() {
	// sort by train to easily track moves per train
	slices.SortStableFunc(printer.Moves, func(a Move, b Move) int {
		return strings.Compare(a.Train.Name, b.Train.Name)
	})
	movesWithDeliveredPackages := make([]Move, 0)
	for _, move := range printer.Moves {
		if len(move.PackagesDropped) > 0 {
			movesWithDeliveredPackages = append(movesWithDeliveredPackages, move)
		}
	}
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Name\tWeight\tDeliveredAt\tTrain\t")

	for _, move := range movesWithDeliveredPackages {
		for _, deliveredPackage := range move.PackagesDropped {
			travelTime := printer.TravelTimeMatrix[move.StartingStation.Id][move.EndingStation.Id]
			fmt.Fprintf(w, "%s\t%dkg\t%dm\t%s\t\n", deliveredPackage.Name, deliveredPackage.Weight, move.TimeTaken+travelTime, move.Train.Name)
		}
	}
	w.Flush()
}
