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
	Lap        int32
	DriverName string
	EstTime    float32
	TimeBehind float32
}

func lapTimeRepresentation(t float32) string {
	if t < 0 {
		t = 0
	}

	wholeSeconds := int64(t)
	lapTime := time.Unix(wholeSeconds, int64((t-float32(wholeSeconds))*1e9))

	return lapTime.Format("04:05.000")
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
		driversLapDistPct := i.Vars.Vars["CarIdxLapDistPct"].Value.([]float32)
		driversEstTime := i.Vars.Vars["CarIdxEstTime"].Value.([]float32)
		driversLap := i.Vars.Vars["CarIdxLap"].Value.([]int32)
		// driversBehind := i.Vars.Vars["CarIdxF2Time"].Value.([]float32)

		drivers := i.SessionInfo.DriverInfo.Drivers
    myIdx := i.SessionInfo.DriverInfo.DriverCarIdx

		standings := make([]StandingsLine, len(drivers))

		fmt.Printf("\033[?25l\033[2J\033[H")
		for k := range len(drivers) {
			if drivers[k].CarIsPaceCar == 1 || drivers[k].IsSpectator == 1 {
				continue
			}

			standings[k] = StandingsLine{
				CarIdx:     k,
				LapPct:     driversLapDistPct[k],
				DriverName: drivers[k].UserName,
				EstTime:    driversEstTime[k],
				Lap:        driversLap[k],
				TimeBehind: driversEstTime[myIdx],
			}
		}

		sort.Slice(standings, func(i int, j int) bool {
			if standings[i].Lap > int32(standings[j].Lap) {
				return true
			}

			return standings[i].LapPct >= standings[j].LapPct
		})

    fmt.Printf("%v\n", driversEstTime)
		// for p, v := range standings {
		// 	fmt.Printf("[%2d] %-30s %13f %13f\n",
		// 		p+1, v.DriverName, v.LapPct, driversEstTime[p] - driversEstTime[myIdx])
		// }

		<-mainLoopTicker.C
	}
}
