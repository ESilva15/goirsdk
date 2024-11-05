package goirsdk

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestParseTelemetrySubHeader_WithGoodBuffer
// Given a well structured buffer it will output the expected
// DiskSubHeader struct
func TestParseTelemetrySubHeader_WithGoodBuffer(t *testing.T) {
	// Arrange
	header := [32]byte{
		0x54, 0x1e, 0x14, 0x67, 0x00, 0x00, 0x00, 0x00, 0x54, 0x09, 0x00, 0xf0,
		0xee, 0x7e, 0x6b, 0x40, 0x53, 0x74, 0x88, 0x44, 0x44, 0x86, 0x8f, 0x40,
		0x08, 0x00, 0x00, 0x00, 0xe1, 0xb8, 0x00, 0x00,
	}

	expectedHeader := DiskSubHeader{
		StartDate:   1729371732,
		StartTime:   219.96666717536084,
		EndTime:     1008.7833338413715,
		LapCount:    8,
		RecordCount: 47329,
	}

	// Act
	headers, err := parseTelemetrySubHeader(header)

	// Assert
	if err != nil {
		t.Fatalf("Error parsing buffer: %v", err)
	}
	if !cmp.Equal(&expectedHeader, headers) {
		t.Fatalf("Expected:\n%#v\nGot:\n%#v\n", expectedHeader, headers)
	}
}
