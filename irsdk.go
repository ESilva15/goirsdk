// Package goirsdk is all you need for you iRacing telemetry parsing
package goirsdk

import (
	"fmt"
	"io"
	"log/slog"
	"os"

	eventutils "github.com/ESilva15/goirsdk/eventutils"
	"github.com/ESilva15/goirsdk/sharedMem"
	"gopkg.in/yaml.v3"
)

// Reader is an interface to represent the readable data that can be either
// a .ibt file (or live data, hopefully)
type Reader interface {
	io.Reader
	io.ReaderAt
	io.ReadCloser
}

type Writer interface {
	io.WriterAt
	io.Closer
}

type TelemetryContainer int

const (
	IBTFile          TelemetryContainer = iota
	SharedMemoryFile TelemetryContainer = iota
)

type Options struct {
	Logger                *slog.Logger
	SourceType            TelemetryContainer // type of source data
	SourcePath            string             // Path to source
	IBTExportType         TelemetryContainer // export type of telemetry: store .ibt or replay in shm
	IBTExportPath         string             // path where to export the data
	IBTExport             bool               // whether to export the telemetry data
	SessionInfoExport     bool               // whether to export the session info data
	SessionInfoExportPath string             // path where to export the session info
}

// IBT struct will hold the relevant data for a given IBT file
type IBT struct {
	File        Reader // Source of the data
	Opts        Options
	IBTExporter Writer
	winUtils    *eventutils.EventUtils // WinUtils gives access to the system utilities

	// TODO: fragment this struct a little bit, for now I want to actually get
	// stuff done so its enough to work as is
	// Actual FILE
	Headers     *TelemetryHeaders // IBT file Headers
	SubHeaders  *DiskSubHeader    // IBT file Sub Headers
	SessionInfo *SessionInfoYAML  // IBT file Session Info
	Vars        *TelemetryVars    // Vars will hold the telemetry data
}

func (i *IBT) IsConnected() bool {
	if i.Headers == nil {
		return false
	}

	if !i.SessionStatusConnected() {
		return false
	}

	if i.SessionStateInvalid() {
		return false
	}

	return true
}

func (i *IBT) exportYAML() error {
	file, err := os.OpenFile(i.Opts.SessionInfoExportPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		i.Opts.Logger.Debug(fmt.Sprintf("Failed to open file for YAML export: %v\n", err))
		return fmt.Errorf("failed to open output file for YAML: %v", err)
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)

	err = enc.Encode(i.SessionInfo)
	if err != nil {
		i.Opts.Logger.Debug(fmt.Sprintf("Failed to write into file for YAML export: %v\n", err))
		return fmt.Errorf("failed to write YAML contents to file: %v", err)
	}

	return nil
}

func (i *IBT) exportIBT(data []byte, offset int64) error {
	nBytes, err := i.IBTExporter.WriteAt(data, offset)
	if err != nil {
		i.IBTExporter.Close()
		i.IBTExporter = nil
		i.Opts.Logger.Debug(fmt.Sprintf("won't attempt to export anymore: %+v", err))
		return err
	}

	if nBytes > 0 {
		// Send the event stating the data has been created
		err = i.winUtils.Utils.SignalEvent()
		if err != nil {
			i.Opts.Logger.Debug("failed to signal event", "err", err)
		} else {
			i.Opts.Logger.Debug("no error signaling: ", "nBytes", nBytes)
		}
	}

	return nil
}

func (i *IBT) openSource() error {
	var err error

	switch i.Opts.SourceType {
	case SharedMemoryFile:
		// User is requesting us to read live data - present in the mem map file
		i.File, err = eventutils.OpenMemMap(MEMMAPFILENAME, fileMapSize)
		if err != nil {
			return fmt.Errorf("failed to open memory mapped file: %+v", err)
		}

		// To use our windows interface we need to initialize it first
		// it will return a struct with a pointer to the windows handles
		// if, for some reason, we need to stub out this to run in on Linux its easier
		i.winUtils, err = eventutils.Init()
		if err != nil {
			return err
		}

		// I don't believe we need this on windows either, but I'll have to check
		// We need to open the windows event thing
		// err = i.winUtils.OpenWinEvent(IRSDK_DATAVALIDEVENTNAME)
		// if err != nil {
		// 	return err
		// }

		// We need to open the broadcast channel
		// err = i.winUtils.OpenBroadcastChannel(IRSDK_BROADCASTMSGNAME)
		// if err != nil {
		// 	return err
		// }
	case IBTFile:
		i.File, err = os.Open(i.Opts.SourcePath)
		if err != nil {
			return fmt.Errorf("failed to open file `%s`: %+v", i.Opts.SourcePath, err)
		}
	default:
		return fmt.Errorf("a source type must be specified")
	}

	return nil
}

func (i *IBT) openExporter() error {
	var err error

	switch i.Opts.IBTExportType {
	case SharedMemoryFile:
		// Lets create a shared memory file!
		shm, err := sharedMem.Create(MEMMAPFILENAME, fileMapSize)
		if err != nil {
			return fmt.Errorf("unable to create memory map file: %+v", err)
		}

		i.IBTExporter = shm
	case IBTFile:
		i.IBTExporter, err = os.OpenFile(i.Opts.IBTExportPath, os.O_CREATE|os.O_RDWR, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open ibt export file: %v", err)
		}
	}

	return nil
}

// Init serves to initialize and get a hold of a IBT struct
// Receives an Options struct with the required configurations
func Init(opts Options) (*IBT, error) {
	// Create our irsdk instance
	var err error
	ibt := IBT{
		Opts: opts,
		Vars: &TelemetryVars{},
	}

	// Set up the event utils
	evutils, err := eventutils.Init()
	if err != nil {
		return nil, err
	}
	evutils.OpenEvent(IRSDK_DATAVALIDEVENTNAME)
	ibt.winUtils = evutils

	// Setup the source
	err = ibt.openSource()
	if err != nil {
		return nil, err
	}

	// Setup the IBT data export - can be either shared memory or data file
	if opts.IBTExport {
		err = ibt.openExporter()
		if err != nil {
			// We log this only, or return some type of message
			// Set the option to false so we won't export
			ibt.Opts.IBTExport = false
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
		return nil, fmt.Errorf("unable to parser variable headers from file: %v", err)
	}

	return &ibt, nil
}

func (i *IBT) ListVariables() map[string]Var {
	return i.Vars.Vars
}

// Close cleans up our irsdk instance
func (i *IBT) Close() {
	if i == nil {
		return
	}

	if i.winUtils != nil {
		// If its not live data, the user is the one with ownership of the handle
		i.File.Close()
		i.winUtils.Close()
	}
}

// LastTick returns the last tick
// func (i *IBT) LastTick() int {
// }
