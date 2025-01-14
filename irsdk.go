// Package goirsdk is all you need for you iRacing telemetry parsing
package goirsdk

import (
	"fmt"
	"os"

	"io"

	"github.com/ESilva15/goirsdk/logger"
	"github.com/ESilva15/goirsdk/winutils"
	"gopkg.in/yaml.v3"
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
	File           Reader                    // Source of the data
	IBTExport      *os.File                  // If set, it will export the IBT data to the file
	IBTExportPath  string                    // Path for IBT export
	YAMLExport     *os.File                  // If set, it will export the session YAML to the file
	YAMLExportPath string                    // Path for YAML export
	Headers        *TelemetryHeaders         // IBT file Headers
	SubHeaders     *DiskSubHeader            // IBT file Sub Headers
	SessionInfo    *SessionInfoYAML          // IBT file Session Info
	Vars           *TelemetryVars            // Vars will hold the telemetry data
	winUtils       *winutils.IRacingWinUtils // WinUtils gives access to the system utilities
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

func (i *IBT) exportYAML() error {
	log := logger.GetInstance()

	file, err := os.OpenFile(i.YAMLExportPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Failed to open file for YAML export: %v\n", err)
		return fmt.Errorf("failed to open output file for YAML: %v", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)

	err = enc.Encode(i.SessionInfo)
	if err != nil {
		log.Printf("Failed to write into file for YAML export: %v\n", err)
		return fmt.Errorf("failed to write YAML contents to file: %v", err)
	}

	return nil
}

func (i *IBT) exportIBT(data []byte, offset int64) error {
	log := logger.GetInstance()

	_, err := i.IBTExport.WriteAt(data, offset)

	if err != nil {
		i.IBTExport.Close()
		i.IBTExport = nil
		log.Println("Won't attempt to export anymore")
		return err
	}

	return nil
}

// Init serves to initialize and get a hold of a IBT struct
// f -> is the source data, pass nil for the SDK to read live data or a
// *os.File to read from a file
// exportTelem -> is a string with the path to export the telemetry data, pass
// an empty string to not export any data
// exportTelem -> is a string with the path to export the session info data, pass
// an empty string to not export any data
func Init(f Reader, exportTelem string, exportYAML string) (*IBT, error) {
	// log := logger.GetInstance()

	// Read the header of the file
	var err error
	ibt := IBT{
		File:           f,
		IBTExport:      nil,
		IBTExportPath:  exportTelem,
		YAMLExport:     nil,
		YAMLExportPath: exportYAML,
		Vars:           &TelemetryVars{},
		winUtils:       nil,
	}

	// If requested to output to a telemetry file
	if exportTelem != "" {
		ibt.IBTExport, err = os.OpenFile(exportTelem, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return nil, fmt.Errorf("failed to open ibt export file: %v", err)
		}
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
	err = ibt.readHeader()
	if err != nil {
		return nil, err
	}

	// Read the disk sub headers
	err = ibt.readSubheader()
	if err != nil {
		return nil, err
	}

	// Read session info string
	err = ibt.readSessionInfo()
	if err != nil {
		return nil, err
	}

	// Read the telemetry vars info
	err = ibt.readVariablerHeaders()
	if err != nil {
		return nil, fmt.Errorf("Unable to parser variable headers from file: %v", err)
	}

	return &ibt, nil
}

// Close cleans up our irsdk instance
func (i *IBT) Close() {
	if i.winUtils != nil {
		// If its not live data, the user is the one with ownership of the handle
		i.File.Close()
		i.winUtils.Close()
	}
}
