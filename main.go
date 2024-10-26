// Package IbtParser is all you need for you iRacing telemetry parsing
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

const (
	ibtFile = "./telemetryFiles/mx5_2016Okayama_full_2024_10_19_22_02_12.ibt"
)

// Reader is an interface (??? very useful)
type Reader interface {
	io.Reader
	io.ReaderAt
	io.ReadCloser
}

// IBT struct will hold the relevant data for a given IBT file
type IBT struct {
	File          Reader            // Source of the data
	Headers       *TelemetryHeaders // IBT file Headers
	SubHeaders    *DiskSubHeader    // IBT file Sub Headers
	SessionInfo   *SessionInfoYAML  // IBT file Session Info
	Vars          *TelemetryVars
	LastValidData int64
	Tick          int32
}

// Init serves to initialize and get a hold of a IBT struct
func Init(f Reader) (*IBT, error) {
	// Read the header of the file
	var err error
	ibt := IBT{
		File: f,
		Vars: &TelemetryVars{},
	}

	// Read the file headers
	var headerRaw [FileHeaderSize]byte
	_, err = ibt.File.ReadAt(headerRaw[:], 0)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from file: %v", err)
	}
	ibt.Headers, err = parseTelemetryHeader(headerRaw)
	if err != nil {
		return nil, fmt.Errorf("Unable to read headers from file: %v", err)
	}

	// Read the disk sub headers
	var subheaderRaw [SubHeaderSize]byte
	_, err = ibt.File.ReadAt(subheaderRaw[:], FileHeaderSize)
	if err != nil {
		return nil, fmt.Errorf("Failed to read subheader from file: %v", err)
	}
	ibt.SubHeaders, err = parseTelemetrySubHeader(subheaderRaw)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse subheaders from file: %v", err)
	}

	// Read session info string
	sessionInfoStringRaw := make([]byte, ibt.Headers.SessionInfoLength)
	_, err = ibt.File.ReadAt(sessionInfoStringRaw, int64(ibt.Headers.SessionInfoOffset))
	if err != nil {
		return nil, fmt.Errorf("Failed to read sessionInfoString from file: %v", err)
	}
	ibt.SessionInfo, err = parseSessionInfo(sessionInfoStringRaw)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse Session Info from file: %v", err)
	}

	return &ibt, nil
}

func (i *IBT) readVariablerHeaders() {
	i.Vars = &TelemetryVars{Vars: make(map[string]Var, i.Headers.NumVars)}

	var k int32
	for k = 0; k < i.Headers.NumVars; k++ {
		rbuf := make([]byte, VarSize)

		_, err := i.File.ReadAt(rbuf, int64(i.Headers.VarHeaderOffset+k*VarSize))
		if err != nil {
			log.Fatal(err)
		}

		var dst IBTVar
		err = binary.Read(bytes.NewBuffer(rbuf[:]), binary.LittleEndian, &dst)
		if err != nil {
			log.Fatal(err)
		}

		v := Var{
			Type:        dst.Type,
			Offset:      dst.Offset,
			Count:       dst.Count,
			CountAsTime: dst.CountAsTime,
			Name:        strings.TrimLeft(strings.TrimRight(string(dst.Name[:]), "\x00"), "\x00"),
			Description: strings.TrimLeft(strings.TrimRight(string(dst.Description[:]), "\x00"), "\x00"),
			Unit:        strings.TrimLeft(strings.TrimRight(string(dst.Unit[:]), "\x00"), "\x00"),
			Value:       nil,
		}

		// fmt.Println(dst.Name)
		// fmt.Println(v.Name)

		i.Vars.Vars[v.Name] = v
	}
}

func (i *IBT) readData() error {
	start := i.Headers.BufOffset + i.Tick*i.Headers.BufLen
	buf := make([]byte, i.Headers.BufLen)
	_, err := i.File.ReadAt(buf, int64(start))
	if err != nil {
		return err
	}

	for k, v := range i.Vars.Vars {
		rbuf := buf[v.Offset : v.Offset+int32(VarTypes[int(v.Type)].Size)]

		// Read the value
    switch v.Type {
    case IRSDK_char:
      v.Value = string(rbuf[0])
    case IRSDK_bool:
      v.Value = int(rbuf[0]) > 0
    case IRSDK_int:
      v.Value = int(binary.LittleEndian.Uint32(rbuf))
    case IRSDK_bitField:
      v.Value = fmt.Sprintf("0x%x", int(binary.LittleEndian.Uint32(rbuf)))
    case IRSDK_float:
      v.Value = math.Float32frombits(uint32(binary.LittleEndian.Uint32(rbuf)))
    case IRSDK_double:
      v.Value = math.Float64frombits(uint64(binary.LittleEndian.Uint64(rbuf)))
    }
		// --------------

		i.Vars.Vars[k] = v
	}

	i.Tick++

	return nil
}

func (i *IBT) Update() bool {
	err := i.readData()
	if err != nil && err != io.EOF {
		log.Fatalf("What happened?\n%v\n", err)
	}

	if err == io.EOF {
		return false
	}

	return true
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
	fmt.Printf("%s\n", ibt.Headers.ToString())
	fmt.Printf("%s", ibt.SubHeaders.ToString())
	//  fmt.Println(ibt.SessionInfo.ToString())

	// last := time.Now().Unix()
	ibt.readVariablerHeaders()
	ibt.Update()

	last := time.Now().UnixMilli()
	for {
    time.Sleep(time.Second / 60)
		res := ibt.Update()

		curTime := time.Now().UnixMilli()

		if curTime - last > 250{
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
}
