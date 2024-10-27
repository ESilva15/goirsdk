package ibtReader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"strings"
)

const (
	VarHeaderSize  = 144
	IRSDK_char     = 0
	IRSDK_bool     = 1
	IRSDK_int      = 2
	IRSDK_bitField = 3
	IRSDK_float    = 4
	IRSDK_double   = 5
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
	tickCount int
	bufOffset int
}

type TelemetryVars struct {
	LastVersion int
	Vars        map[string]Var
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

func (i *IBT) readData() error {
	// I think that we can add one extra check or verification here
	// The file headers tells us how many data frames there are, we can probably
	// cap it at that instead of waiting for the read to fail
	// Probably wouldn't work on live data tho
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
