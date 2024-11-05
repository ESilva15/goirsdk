package goirsdk

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	FileHeaderSize = 112 // FileHeaderSize is the size of the headers
	HeaderSize     = 4   // HeaderSize is the size of a single header
)

// TelemetryHeaders struct to hold an IBT file's headers
type TelemetryHeaders struct {
	Version int32
	// Status of 1 indicates a completed session and status of 0 a live session
	Status int32
	// TickRate indicates the frequency of writes (usually 60)
	TickRate int32
	// SessionInfoUpdate indicates the number of times the SessionInfo was
	// updated. 0 for finished sessions and >1 for active sessions
	SessionInfoUpdate int32
	// SessionInfoLength is the length of the session info buffer
	SessionInfoLength int32
	// SessionInfoOffset is the offset of the session info in the buffer
	SessionInfoOffset int32
	// NumVars is the number of variables in each input
	NumVars int32
	// VarHeaderOffset is the offset of the VarHeader
	VarHeaderOffset int32
	// NumBuf will be 1 for static files and 3 for live telemetry files
	NumBuf int32
	// BufLen is the length for parsing VarHeader values
	BufLen int32
	// Padding
	Padding [12]byte
	// I still don't know what this is:
	BufOffset int32
}

// ToString renders a string showing the values of the struct
func (th *TelemetryHeaders) ToString() string {
	return fmt.Sprintf(
		"Version:                               %5d (0x%04x)\n"+
			"Status:                                %5d (0x%04x)\n"+
			"TickRate:                              %5d (0x%04x)\n"+
			"SIUpdate:                              %5d (0x%04x)\n"+
			"SILength:                              %5d (0x%04x)\n"+
			"SIOffset:                              %5d (0x%04x)\n"+
			"NumVars:                               %5d (0x%04x)\n"+
			"VarHeaderOffset:                       %5d (0x%04x)\n"+
			"NumBuf:                                %5d (0x%04x)\n"+
			"BufLen:                                %5d (0x%04x)\n"+
			"BufOffset:                             %5d (0x%04x)\n",
		th.Version, th.Version, th.Status, th.Status, th.TickRate, th.TickRate,
		th.SessionInfoUpdate, th.SessionInfoUpdate,
		th.SessionInfoLength, th.SessionInfoLength,
		th.SessionInfoOffset, th.SessionInfoOffset,
		th.NumVars, th.NumVars, th.VarHeaderOffset, th.VarHeaderOffset,
		th.NumBuf, th.NumBuf, th.BufLen, th.BufLen, th.BufOffset, th.BufOffset,
	)
}

// parseTelemetryHeader will read the IBT file headers from a correctly sized
// buffer.
// You need to pass a the first FILE_HEADER_SIZE bytes of the buffer
func parseTelemetryHeader(buf [FileHeaderSize]byte) (*TelemetryHeaders, error) {
	// utils.HexDump(buf[:])
	// fmt.Printf("Len: %d\n", len(buf))

	dst := TelemetryHeaders{}
	err := binary.Read(bytes.NewBuffer(buf[:]), binary.LittleEndian, &dst)
	if err != nil {
		return nil, fmt.Errorf("unable to unpack data: %v", err)
	}

	return &dst, nil
}
