// Package IbtReader is all you need for you iRacing telemetry parsing
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"ibtReader/sharedMem"
)

const (
	ibtFile = "./telemetryFiles/mx5_2016Okayama_full_2024_10_19_22_02_12.ibt"
)

// Reader is an interface to represent the readable data that can be either
// a .ibt file (or live data, hopefully)
type Reader interface {
	io.Reader
	io.ReaderAt
	io.ReadCloser
}

// IBT struct will hold the relevant data for a given IBT file
type IBT struct {
	File        Reader            // Source of the data
	Headers     *TelemetryHeaders // IBT file Headers
	SubHeaders  *DiskSubHeader    // IBT file Sub Headers
	SessionInfo *SessionInfoYAML  // IBT file Session Info
	Vars        *TelemetryVars    // Vars will hold the telemetry data
	Tick        int32             // Tick holds the cound of the reads
	LiveData    bool
}

// Init serves to initialize and get a hold of a IBT struct
func Init(f Reader) (*IBT, error) {
	// Read the header of the file
	var err error
	ibt := IBT{
		File:     f,
		Vars:     &TelemetryVars{},
		LiveData: false,
	}

	if ibt.File == nil {
		// User is requesting us to read live data - present in the mem map file
		ibt.File, err = sharedMem.Open(IRSDK_MEMMAPFILENAME, fileMapSize)
		if err != nil {
			return nil, fmt.Errorf("Failed to open memory mapped file: %v", err)
		}
		ibt.LiveData = true
	}

	// Read the file headers
	var headerRaw [FileHeaderSize]byte
	_, err = ibt.File.ReadAt(headerRaw[:], 0)
	if err != nil {
		return nil, fmt.Errorf("Failed to read headers from file: %v", err)
	}
	ibt.Headers, err = parseTelemetryHeader(headerRaw)
	if err != nil {
		return nil, fmt.Errorf("Unable to read headers from file: %v", err)
	}
	fmt.Println(ibt.Headers.ToString())

	// Read the disk sub headers
	var subheaderRaw [SubHeaderSize]byte
	_, err = ibt.File.ReadAt(subheaderRaw[:], 112)
	if err != nil {
		return nil, fmt.Errorf("Failed to read disk subheaders from file: %v", err)
	}
	ibt.SubHeaders, err = parseTelemetrySubHeader(subheaderRaw)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse disk subheaders from file: %v", err)
	}
	fmt.Println(ibt.SubHeaders.ToString())

	// Read session info string
	sessionInfoStringRaw := make([]byte, ibt.Headers.SessionInfoLength)
	_, err = ibt.File.ReadAt(sessionInfoStringRaw, int64(ibt.Headers.SessionInfoOffset))
	if err != nil {
		return nil, fmt.Errorf("Failed to read sessionInfoString from file: %v", err)
	}
	ibt.SessionInfo, err = parseSessionInfo(sessionInfoStringRaw, ibt.Headers.SessionInfoLength)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse SessionInfoString from file: %v", err)
	}

	// Read the telemetry vars info
	err = ibt.readVariablerHeaders()
	if err != nil {
		return nil, fmt.Errorf("Unable to parser variable headers from file: %v", err)
	}

	return &ibt, nil
}

func msToKph(v float32) int {
	return int((3600 * v) / 1000)
}

func main() {
	fmt.Println("================== IBT FILE PARSER ==================")

	file, err := os.Open(ibtFile)
	if err != nil {
		log.Fatalf("Failed to open IBT file: %v", err)
	}

	ibt, err := Init(file)
	if err != nil {
		log.Fatalf("Failed to create irsdk instance: %v", err)
	}
	fmt.Printf("%s\n", ibt.Headers.ToString())
	fmt.Printf("%s\n", ibt.SubHeaders.ToString())
	fmt.Printf("%s\n", ibt.SessionInfo.ToString())

	// Display the human readable start date
	unixStartDate := time.Unix(ibt.SubHeaders.StartDate, 0)
	startDate := unixStartDate.Format("2006/01/02 15:04:05 -0700 MST")
	fmt.Println("StartDate:", startDate)

	// Display the human readable version of start time
	unixStartTime := time.Unix(ibt.SubHeaders.StartDate+int64(ibt.SubHeaders.StartTime), 0)
	startTime := unixStartTime.Format("2006/01/02 15:04:05 -0700 MST")
	fmt.Println("StartTime:", startTime)

	// Display the human readable version of end time
	unixEndTime := time.Unix(ibt.SubHeaders.StartDate+int64(ibt.SubHeaders.EndTime), 0)
	endTime := unixEndTime.Format("2006/01/02 15:04:05 -0700 MST")
	fmt.Println("EndTime:  ", endTime)

	last := time.Now().UnixMilli()
	for {
		time.Sleep(time.Second / 60)
		res := ibt.Update()

		curTime := time.Now().UnixMilli()

		if curTime-last > 250 {
			fmt.Printf("                                                           \r")
			if val, ok := ibt.Vars.Vars["Speed"]; ok {
				fmt.Printf("\r%d %d", ibt.Tick/60, msToKph(val.Value.(float32)))
			} else {
				fmt.Printf("\r%d %s", ibt.Tick/60, "KEY DOESN'T EXIST")
			}
		}

		if !res {
			fmt.Println("\nEnd of file found...")
			break
		}
	}

	fmt.Printf("%d\n", ibt.Tick)
}
