// Package goirsdk is all you need for you iRacing telemetry parsing
package goirsdk

import (
	"fmt"
	"log"
	"os"

	"io"

	"github.com/ESilva15/goirsdk/winutils"
)

const (
	ibtFile = "./telemetryFiles/mx5_2016Okayama_full_2024_10_19_22_02_12.ibt"
)

func msToKph(v float32) int {
	return int((3600 * v) / 1000)
}

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
	FileToExport *os.File                  // If set, it will export the IBT data to the file
	YAMLExport   *os.File                  // If set, it will export the session YAML to the file
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

func (i *IBT) ExportToIBT() {
	rbuf := make([]byte, fileMapSize)

	_, err := i.File.ReadAt(rbuf, 0)
	if err != nil {
		log.Fatal(err)
	}

	_, err = i.FileToExport.Write(rbuf)
	if err != nil {
		log.Fatal(err)
	}
}

// Init serves to initialize and get a hold of a IBT struct
func Init(f Reader, exportTelem string, exportYAML string) (*IBT, error) {
	// Read the header of the file
	var err error
	ibt := IBT{
		File:         f,
		FileToExport: nil,
		YAMLExport:   nil,
		Vars:         &TelemetryVars{},
		winUtils:     nil,
	}

	// If requested to output to a telemetry file
	if exportTelem != "" {
		ibt.FileToExport, err = os.OpenFile(exportTelem, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, err
		}
		ibt.FileToExport.Sync()
	}

	// Do the same for the YAMLExport
	if exportYAML != "" {
		// Do the thing here, yo
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
	// Write to the output file
	_, err = ibt.FileToExport.WriteAt(headerRaw[:], 0)
	if err != nil {
		log.Fatal(err)
	}

	// Read the disk sub headers
	var subheaderRaw [SubHeaderSize]byte
	_, err = ibt.File.ReadAt(subheaderRaw[:], HeaderSize)
	if err != nil {
		return nil, fmt.Errorf("Failed to read disk subheaders from file: %v", err)
	}
	ibt.SubHeaders, err = parseTelemetrySubHeader(subheaderRaw)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse disk subheaders from file: %v", err)
	}
	// Write to the output file
	_, err = ibt.FileToExport.WriteAt(subheaderRaw[:], HeaderSize)
	if err != nil {
		log.Fatal(err)
	}

	// Read session info string
	sessionInfoStringRaw := make([]byte, ibt.Headers.SessionInfoLength)
	_, err = ibt.File.ReadAt(sessionInfoStringRaw, int64(ibt.Headers.SessionInfoOffset))
	if err != nil {
		return nil, fmt.Errorf("Failed to read sessionInfoString from file: %v", err)
	}
	// Write to the output file
	_, err = ibt.FileToExport.WriteAt(sessionInfoStringRaw[:], int64(ibt.Headers.SessionInfoOffset))
	if err != nil {
		log.Fatal(err)
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
