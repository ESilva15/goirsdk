package goirsdk

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
	"time"
)

const (
	VarHeaderSize               = 144
	IRSDK_char                  = 0
	IRSDK_bool                  = 1
	IRSDK_int                   = 2
	IRSDK_bitField              = 3
	IRSDK_float                 = 4
	IRSDK_double                = 5
	Running        IRacingState = iota
	Paused
	Ended
	Failed
	Unknown
)

// I think I can make an interface if IRSDK types with available types and
// that they need a parser (reads and type coerces I guess)
var (
	VarTypes = map[int]VarType{
		IRSDK_char:     {1, "irsdk_char"},
		IRSDK_bool:     {1, "irsdk_bool"},
		IRSDK_int:      {4, "irsdk_int"},
		IRSDK_bitField: {4, "irsdk_bitField"},
		IRSDK_float:    {4, "irsdk_float"},
		IRSDK_double:   {8, "irsdk_double"},
	}
)

type IRacingState int
type VarType struct {
	Size int    // Size is the var type size in bytes
	Name string // Name is the irsdk var name
}

type IBTVar struct {
	Type        int32
	Offset      int32
	Count       int32
	CountAsTime bool
	Padding     [3]byte
	Name        [32]byte
	Description [64]byte
	Unit        [32]byte
}

type Var struct {
	Type        int32
	Offset      int32
	Count       int32
	CountAsTime bool
	Name        string
	Description string
	Unit        string
	// TODO
	// Create an interface for this value
	// Represent the IRSDK var types with a struct each that implements the Parse
	// method or something like that I guess
	Value interface{}
}

func (v *IBTVar) ToString() string {
	return fmt.Sprintf(
		"Type:        %5d (0x%08x)\n"+
			"Offset:      %5d (0x%08x)\n"+
			"Count:       %5d (0x%08x)\n"+
			"CountAsTime: %5t\n"+
			"Name:        %s\n"+
			"Description: %s\n"+
			"Unit:        %s",
		v.Type, v.Type, v.Offset, v.Offset, v.Count, v.Count,
		v.CountAsTime, v.Name, v.Description, v.Unit,
	)
}

type varBuffer struct {
	TickCount int32
	BufOffset int32
}

type TelemetryVars struct {
	Tick         int32          // Keeps track of the current data buffer tick
	RecorderTick int32          // Counts from 0 when creating a telemetry file from a replay or live data
	Vars         map[string]Var // Variables content
}

func (i *IBT) readVariablerHeaders() error {
	i.Vars = &TelemetryVars{Vars: make(map[string]Var, i.Headers.NumVars)}

	var k int32
	for k = 0; k < i.Headers.NumVars; k++ {
		rbuf := make([]byte, VarHeaderSize)

		_, err := i.File.ReadAt(rbuf, int64(i.Headers.VarHeaderOffset+k*VarHeaderSize))
		if err != nil {
			return err
		}

		_, err = i.FileToExport.WriteAt(rbuf, int64(i.Headers.VarHeaderOffset+k*VarHeaderSize))
		if err != nil {
			log.Fatal(err)
		}

		var dst IBTVar
		err = binary.Read(bytes.NewBuffer(rbuf[:]), binary.LittleEndian, &dst)
		if err != nil {
			return err
		}

		v := Var{
			Type:        dst.Type,
			Offset:      dst.Offset,
			Count:       dst.Count,
			CountAsTime: dst.CountAsTime,
			Name:        strings.TrimRight(string(dst.Name[:]), "\x00"),
			Description: strings.TrimRight(string(dst.Description[:]), "\x00"),
			Unit:        strings.TrimRight(string(dst.Unit[:]), "\x00"),
			Value:       nil,
		}

		i.Vars.Vars[v.Name] = v
	}

	return nil
}

func (i *IBT) readData(buf []byte) error {
	for k, v := range i.Vars.Vars {
		// Slice of the variable value in the buffer
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

	return nil
}

func (i *IBT) Update(timeout time.Duration) (IRacingState, error) {
	if i.winUtils != nil {
		// Put a way to check if the sim is active here
		// fmt.Println("NOT CHECKING IF SIM IS ACTIVE - ADD ME")

		// WORKING HERE
		// Need to figure out how to grab the latest buffer with data
		var vb varBuffer
		foundTickCount := 0
		for k := 0; k < int(i.Headers.NumBuf); k++ {
			rbuf := make([]byte, 16)
			// Read 16 bytes, I don't know why, but do need to understand this
			_, err := i.File.ReadAt(rbuf, int64(48+k*16))
			if err != nil {
				return Failed, err
			}

			var curVb varBuffer
			err = binary.Read(bytes.NewBuffer(rbuf[:]), binary.LittleEndian, &curVb)
			if err != nil {
				return Failed, err
			}

			if foundTickCount < int(curVb.TickCount) {
				foundTickCount = int(curVb.TickCount)
				vb = curVb
			}
		}

		i.Vars.Tick = vb.TickCount

		start := vb.BufOffset
		buf := make([]byte, i.Headers.BufLen)

		_, err := i.File.ReadAt(buf, int64(start))
		if err != nil {
			return Failed, err
		}

		_, err = i.FileToExport.WriteAt(buf, int64(i.Headers.BufOffset+i.Vars.RecorderTick*i.Headers.BufLen))
		if err != nil {
			log.Fatalf("Failed to write to file [2]: %v", err)
		}

		err = i.readData(buf)
		if err != nil && err != io.EOF {
			return Unknown, err
		}

		if err == io.EOF {
			return Ended, nil
		}

		i.Vars.RecorderTick++
	} else {
		// This will get the dataframe corresponding to a given tick
		start := i.Headers.BufOffset + i.Vars.Tick*i.Headers.BufLen
		buf := make([]byte, i.Headers.BufLen)
		_, err := i.File.ReadAt(buf, int64(start))

		// Make this happen in a different thread, or have this send to a queue that has a thread
		// writing to a file
		_, err = i.FileToExport.WriteAt(buf, int64(start))
		if err != nil {
			log.Fatalf("Failed to write to file [1]: %v\n", err)
		}

		if err == io.EOF {
			return Ended, nil
		}
		if err != nil {
			return Unknown, err
		}

		err = i.readData(buf)
		if err != nil && err != io.EOF {
			log.Fatalf("What happened?\n%v\n", err)
		}

		// This was previously in the read data method, but it probably fits here better
		i.Vars.Tick++
	}

	return Running, nil
}
