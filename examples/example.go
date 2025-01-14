package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ESilva15/goirsdk"
)

func msToKph(v float32) int {
	return int((3600 * v) / 1000)
}

func main() {
	// Open the data source file
	file, err := os.Open("/path/to/ibtFile")
	if err != nil {
		log.Fatalf("Failed to open IBT file: %v", err)
	}

	// Instantiate our iRacing SDK instance
	irsdk, err := goirsdk.Init(file, "", "")
	if err != nil {
		log.Fatalf("Failed to create iRacing interface: %v", err)
	}
	defer irsdk.Close()

	// Set up a loop to iterate our data
	mainLoopTicker := time.NewTicker(time.Second / 60)
	defer mainLoopTicker.Stop()

	for {
		// Update the data that the SDK is holding with the next tick
		_, err := irsdk.Update(100 * time.Millisecond)
		if err != nil {
			log.Printf("could not update data: %v", err)
			continue
		}

		// Vehicle Movement data gathered from the names we can find on the
		// telemetry_docs.pdf file
		// - I wish to make this less verbose if possible
		if _, ok := irsdk.Vars.Vars["Gear"]; !ok {
			log.Fatal("Field `Gear` doesn't exist")
		}

		if _, ok := irsdk.Vars.Vars["RPM"]; !ok {
			log.Fatal("Field `RPM` doesn't exist")
		}

		if _, ok := irsdk.Vars.Vars["Speed"]; !ok {
			log.Fatal("Field `Speed` doesn't exist")
		}

		gear := int32(irsdk.Vars.Vars["Gear"].Value.(int))
		rpm := int32(irsdk.Vars.Vars["RPM"].Value.(float32))
		speed := int32(msToKph(irsdk.Vars.Vars["Speed"].Value.(float32)))

		fmt.Printf("\033[?25l\033[2J\033[H")
		fmt.Printf("Gear: %d, RPM: %d, Speed: %d", gear, rpm, speed)

		<-mainLoopTicker.C
	}
}
