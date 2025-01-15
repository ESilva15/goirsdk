package goirsdk

import (
  "fmt"
	"log"
	"os"
	"testing"
	"time"
)

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
    driversPositions := i.Vars.Vars["CarIdxPosition"].Value.([]int)

		drivers := i.SessionInfo.DriverInfo.Drivers

		fmt.Printf("\033[?25l\033[2J\033[H")
		for k := range len(drivers) {
			if drivers[k].CarIsPaceCar == 1 || drivers[k].IsSpectator == 1 {
				continue
			}

			fmt.Printf("[%d] [%d] %s\n", k, driversPositions[drivers[k].CarIdx], drivers[k].UserName)
		}

		<-mainLoopTicker.C
	}
}
