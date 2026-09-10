package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/ESilva15/goirsdk"
)

func msToKph(v float32) int {
	return int((3600 * v) / 1000)
}

func main() {
	output, err := os.OpenFile("./output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		log.Fatalf("Failed to open log file: %+v", err)
	}

	logger := slog.New(
		slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}),
	)

	logger.Debug("Starting test")

	// Instantiate our iRacing SDK instance
	irsdk, err := goirsdk.Init(goirsdk.Options{
		Logger:     logger,
		SourceType: goirsdk.SharedMemoryFile,
	})
	if err != nil {
		log.Fatalf("Failed to create iRacing interface: %v", err)
	}
	defer irsdk.Close()

	// Set up a loop to iterate our data
	// TODO: revert this 240 back to 60 because i recorded the thing wrong or whatever
	mainLoopTicker := time.NewTicker(time.Second / 240)
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
		sessionState := irsdk.Vars.Vars["SessionState"].Value.(int)
		trkloc := irsdk.Vars.Vars["PlayerTrackSurface"].Value.(int)
		trksurf := irsdk.Vars.Vars["PlayerTrackSurfaceMaterial"].Value.(int)
		pitsvflags := irsdk.Vars.Vars["PitSvFlags"].Value.(uint32)

		fmt.Printf("\033[?25l\033[2J\033[H")
		fmt.Printf("Gear: %d, RPM: %d, Speed: %d\n", gear, rpm, speed)
		fmt.Printf("SessionState: %s\n", goirsdk.SessionStateToString(sessionState))
		fmt.Printf("TrkLoc: %s\n", goirsdk.TrkLocToString(trkloc))
		fmt.Printf("TrkSurf: %s\n", goirsdk.TrkSurfToString(trksurf))
		fmt.Printf("PitSvFlags: %d\n", pitsvflags)
		fmt.Printf("    FL  FR\n")
		fmt.Printf("    %t  %t\n", irsdk.LFTireChange(), irsdk.RFTireChange())
		fmt.Printf("\n")
		fmt.Printf("    RL  RR\n")
		fmt.Printf("    %t  %t\n", irsdk.LRTireChange(), irsdk.RRTireChange())
		fmt.Printf("    FuelFill:          %t\n", irsdk.FuelFill())
		fmt.Printf("    WindshieldTearoff: %t\n", irsdk.WindshieldTearoff())
		fmt.Printf("    FastRepair:        %t\n", irsdk.FastRepair())
		fmt.Printf("    ClearTires:        %t\n", irsdk.ClearTires())
		fmt.Printf("    ClearWS:           %t\n", irsdk.ClearWS())
		fmt.Printf("    ClearFR:           %t\n", irsdk.ClearFR())
		fmt.Printf("    ClearFuel:         %t\n", irsdk.ClearFuel())

		<-mainLoopTicker.C
	}
}
