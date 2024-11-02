package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// TestParseTelemetryHeader_WithGoodBuffer
// Given a well structured buffer it will output the expected
// TelemetryHeaders struct
func TestParseTelemetryHeader_WithGoodBuffer(t *testing.T) {
	// Arrange
	header := [112]byte{
		0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x3c, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x05, 0x3f, 0x00, 0x00, 0x90, 0x99, 0x00, 0x00,
		0x10, 0x01, 0x00, 0x00, 0x90, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
		0x1d, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x60, 0x2b, 0x00, 0x00, 0x95, 0xd8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	expectedHeader := TelemetryHeaders{
		Version:           2,
		Status:            1,
		TickRate:          60,
		SessionInfoUpdate: 0,
		SessionInfoLength: 16133,
		SessionInfoOffset: 39312,
		NumVars:           272,
		VarHeaderOffset:   144,
		NumBuf:            1,
		BufLen:            1053,
		Padding:           [12]byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x60, 0x2b, 0x0, 0x0},
		BufOffset:         55445,
	}

	// Act
	headers, err := parseTelemetryHeader(header)

	// Assert
	if err != nil {
		t.Fatalf("Error parsing buffer: %v", err)
	}
	if !cmp.Equal(&expectedHeader, headers) {
		t.Fatalf("Expected:\n%#v\nGot:\n%#v\n", expectedHeader, headers)
	}
}
