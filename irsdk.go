// Package goirsdk is all you need for you iRacing telemetry parsing
package goirsdk

import (
	"fmt"
	"log"
	"os"
	// "os"
	"io"
	// "log"
	// "time"

	// conv "ibtReader/conversions"
	"github.com/ESilva15/goirsdk/winutils"
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
	File         Reader                    // Source of the data
	FileToExport string                    // If set, it will export the IBT data to the file
	Headers      *TelemetryHeaders         // IBT file Headers
	SubHeaders   *DiskSubHeader            // IBT file Sub Headers
	SessionInfo  *SessionInfoYAML          // IBT file Session Info
	Vars         *TelemetryVars            // Vars will hold the telemetry data
	winUtils     *winutils.IRacingWinUtils // WinUtils gives access to the system utilities
}

func (i *IBT) IsConnected() bool {
	if i.Headers != nil {
		if sessionStatusOK(int(i.Headers.Status)) {
			return true
		}
		// if sessionStatusOK(int(i.Headers.Status)) && (sdk.lastValidData+connTimeout > time.Now().Unix()) {
		// 	return true
		// }
	}

	return false
}

func (i *IBT) ExportToIBT(filepath string) {
	rbuf := make([]byte, fileMapSize)

	_, err := i.File.ReadAt(rbuf, 0)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filepath, rbuf, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

// Init serves to initialize and get a hold of a IBT struct
func Init(f Reader) (*IBT, error) {
	// Read the header of the file
	var err error
	ibt := IBT{
		File:     f,
		Vars:     &TelemetryVars{},
		winUtils: nil,
	}

	if ibt.File == nil {
		// User is requesting us to read live data - present in the mem map file
		ibt.File, err = winutils.OpenMemMap(IRSDK_MEMMAPFILENAME, fileMapSize)
		if err != nil {
			return nil, fmt.Errorf("Failed to open memory mapped file: %v", err)
		}

		// To use our windows interface we need to initialize it first
		// it will return a struct with a pointer to the windows handles
		// if, for some reason, we need to stub out this to run in on Linux its easier
		ibt.winUtils, err = winutils.Init()
		if err != nil {
			return nil, err
		}

		// We need to open the windows event thing
		err = ibt.winUtils.OpenWinEvent(IRSDK_DATAVALIDEVENTNAME)
		if err != nil {
			return nil, err
		}

		// We need to open the broadcast channel
		err = ibt.winUtils.OpenBroadcastChannel(IRSDK_BROADCASTMSGNAME)
		if err != nil {
			return nil, err
		}
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

func (i *IBT) Close() {
	if i.winUtils != nil {
		// If its not live data, the user is the one with ownership of the handle
		i.File.Close()
		i.winUtils.Close()
	}
}

// func main() {
// 	fmt.Println("================== IBT FILE PARSER ==================")
//
// 	file, err := os.Open(ibtFile)
// 	if err != nil {
// 		log.Fatalf("Failed to open IBT file: %v", err)
// 	}
//
// 	ibt, err := Init(file)
// 	if err != nil {
// 		log.Fatalf("Failed to create irsdk instance: %v", err)
// 	}
// 	// fmt.Printf("%s\n", ibt.Headers.ToString())
// 	// fmt.Printf("%s\n", ibt.SubHeaders.ToString())
// 	// fmt.Printf("%s\n", ibt.SessionInfo.ToString())
//
// 	// Display the human readable start date
// 	// unixStartDate := time.Unix(ibt.SubHeaders.StartDate, 0)
// 	// startDate := unixStartDate.Format("2006/01/02 15:04:05 -0700 MST")
// 	// fmt.Println("StartDate:", startDate)
//
// 	// Display the human readable version of start time
// 	// unixStartTime := time.Unix(ibt.SubHeaders.StartDate+int64(ibt.SubHeaders.StartTime), 0)
// 	// startTime := unixStartTime.Format("2006/01/02 15:04:05 -0700 MST")
// 	// fmt.Println("StartTime:", startTime)
//
// 	// Display the human readable version of end time
// 	// unixEndTime := time.Unix(ibt.SubHeaders.StartDate+int64(ibt.SubHeaders.EndTime), 0)
// 	// endTime := unixEndTime.Format("2006/01/02 15:04:05 -0700 MST")
// 	// fmt.Println("EndTime:  ", endTime)
//
// 	last := time.Now().UnixMilli()
// 	for {
// 		time.Sleep(time.Second / 60)
// 		res, err := ibt.Update(100 * time.Millisecond)
// 		if res == Unknown {
// 			log.Fatalf("Some unknown error occurred: %v\n", err)
// 		}
//
// 		if res == Paused {
// 			fmt.Printf("\r                                                                    \r")
// 			fmt.Println("GAME IS PAUSED")
// 			continue
// 		}
//
// 		curTime := time.Now().UnixMilli()
//
// 		if curTime-last > 250 {
// 			fmt.Printf("                                                           \r")
// 			if val, ok := ibt.Vars.Vars["Speed"]; ok {
// 				fmt.Printf("\r%d %d", ibt.Vars.Tick/60, conv.MsToKph(val.Value.(float32)))
// 			} else {
// 				fmt.Printf("\r%d %s", ibt.Vars.Tick/60, "KEY DOESN'T EXIST")
// 			}
// 		}
//
// 		if res == Ended {
// 			fmt.Println("\nEnd of file found...")
// 			break
// 		}
// 	}
// 	fmt.Printf("%d\n", ibt.Vars.Tick)
//
// 	ibt.Close()
// }
