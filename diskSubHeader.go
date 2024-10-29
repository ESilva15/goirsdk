package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const (
	SubHeaderSize = 32 // SubHeaderSize is the size of the subheader
)

// DiskSubHeader represents the IBT sub headers
type DiskSubHeader struct {
	StartDate   int64   // StartDate represents the start date of the telemetry
	StartTime   float64 // StartTime of file relative to start of session
	EndTime     float64 // EndTime of file relative to start of session
	LapCount    int32   // LapCount represents the total number laps
	RecordCount int32   // RecordCount holds the number of data frames
}

// parseTelemetrySubHeader will return a pointer to a DiskSubHeader variable
// or nil if an error occurs. In which case the error return value is more
// valuable
func parseTelemetrySubHeader(buf [SubHeaderSize]byte) (*DiskSubHeader, error) {
	// utils.HexDump(buf[:])

	dst := DiskSubHeader{}
	err := binary.Read(bytes.NewBuffer(buf[:]), binary.LittleEndian, &dst)
	if err != nil {
		return nil, err
	}

	return &dst, nil
}

// ToString renders a string showing the values of the struct
func (d *DiskSubHeader) ToString() string {
	return fmt.Sprintf(
		"StartDate:   %13d (0x%04x)\n"+
			"StartTime:   %13f (0x%08x)\n"+
			"EndTime:     %13f (0x%08x)\n"+
			"LapCount:    %13d (0x%04x)\n"+
			"RecordCount: %13d (0x%04x)\n",
		d.StartDate, d.StartDate, d.StartTime, d.StartTime,
		d.EndTime, d.EndTime, d.LapCount, d.LapCount,
		d.RecordCount, d.RecordCount,
	)
}
