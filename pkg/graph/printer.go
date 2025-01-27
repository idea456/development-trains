package graph

import (
	"fmt"
	"slices"
	"strings"
)

type Printer struct {
	Moves        []Move
	StationNames map[StationId]string
}

func NewPrinter(moves []Move, stationNames map[StationId]string) *Printer {
	return &Printer{
		Moves:        moves,
		StationNames: stationNames,
	}
}

func (printer *Printer) PrintMoves() {
	// W=0, T=Q1, N1=B, P1=[], N2=A, P2=[]
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
}

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
			fmt.Println("Droppped packages:")
			for _, dropppedPackage := range move.PackagesDropped {
				fmt.Printf("	- %s package with weight %d at %s station\n", dropppedPackage.Name, dropppedPackage.Weight, printer.StationNames[dropppedPackage.EndingStationId])
			}
		}
		fmt.Println()
	}
}
