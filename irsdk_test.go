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

	i, err := Init(input, "", "")
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
		if _, ok := i.Vars.Vars["Gear"]; !ok {
			log.Fatal("Field `Gear` doesn't exist")
		}

		gear := int32(i.Vars.Vars["Gear"].Value.(int))

		fmt.Printf("\033[?25l\033[2J\033[H")
		fmt.Printf("Gear: %d\n", gear)

		<-mainLoopTicker.C
	}
}
