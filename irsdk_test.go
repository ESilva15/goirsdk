package goirsdk

import (
	"fmt"
	"log"
	"os"
	"sort"
	"testing"
	"time"
)

type StandingsLine struct {
	CarIdx     int
	LapPct     float32
	DriverName string
}

func TestFunctionality(t *testing.T) {
	input, err := os.Open("../testTelemetry/supercars_race_watkins_glenn.ibt")
	if err != nil {
		t.Fatal("Was unable to prepare telemetry file for testing.")
	}

	i, _ := Init(input, "", "")
	defer i.Close()

	// Set up a loop to iterate our data
	mainLoopTicker := time.NewTicker(time.Second / 60)
	defer mainLoopTicker.Stop()

	for {
		// Update the data that the SDK is holding with the next tick
		_, err := i.Update(100 * time.Millisecond)
		if err != nil {
			log.Printf("could not update data: %v", err)
			continue
		}

		// Vehicle Movement data gathered from the names we can find on the
		// telemetry_docs.pdf file
		// - I wish to make this less verbose if possible
		if _, ok := i.Vars.Vars["CarIdxPosition"]; !ok {
			log.Fatal("Field `CarIdxPosition` doesn't exist")
		}
		driversPositions := i.Vars.Vars["CarIdxLapDistPct"].Value.([]float32)

		drivers := i.SessionInfo.DriverInfo.Drivers

		standings := make([]StandingsLine, len(drivers))

		fmt.Printf("\033[?25l\033[2J\033[H")
		for k := range len(drivers) {
			if drivers[k].CarIsPaceCar == 1 || drivers[k].IsSpectator == 1 {
				continue
			}

			standings[k] = StandingsLine{
				CarIdx:     k,
				LapPct:     float32(driversPositions[k]),
				DriverName: drivers[k].UserName,
			}
		}

		sort.Slice(standings, func(i int, j int) bool {
			return standings[i].LapPct >= standings[j].LapPct
		})

		for p, v := range standings {
			fmt.Printf("[%d] %s %f\n", p + 1, v.DriverName, v.LapPct)
		}

		<-mainLoopTicker.C
	}
}
