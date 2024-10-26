package main

import (
	"fmt"
)

const (
	VarSize        = 144
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
	Padding     [0]byte
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
	Value       interface{}
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

// func findLatestBuffer(i *IBT) varBuffer {
// 	var vb varBuffer
// 	foundTickCount := 0
// 	for k := 0; k < int(i.Headers.NumBuf); k++ {
// 		rbuf := make([]byte, 16)
// 		_, err := i.File.ReadAt(rbuf, int64(48+k*16))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		currentVb := varBuffer{
// 			int(binary.LittleEndian.Uint32(rbuf[0:4])),
// 			int(binary.LittleEndian.Uint32(rbuf[4:8])),
// 		}
//
// 		if foundTickCount < currentVb.tickCount {
// 			foundTickCount = currentVb.tickCount
// 			vb = currentVb
// 		}
// 	}
//
//   return vb
// }

// func (i *IBT) parseVariableHeaders(offset int32) (int32, error) {
//   if i.Vars.Vars == nil {
//     i.Vars.Vars = make(map[string]*Variable, i.Headers.NumVars)
//   }
//
//   var size int32 = 0
// 	for k := range i.Headers.NumVars {
// 		start := k * VarSize + offset
//     size += start
//
// 		buf := make([]byte, VarSize)
// 		_, err := i.File.ReadAt(buf, int64(start))
// 		if err != nil {
// 			return 0, err
// 		}
//
// 		newVar, err := parseVariable(buf)
// 		i.Vars.Vars[string(newVar.Name[:])] = newVar
// 	}
//
//   return size, nil
// }

// This function will read a single variable
// func parseVariable(buf []byte) (*Variable, error) {
// 	if len(buf)%VarSize != 0 {
// 		return nil, fmt.Errorf("buffer must be multiple of size: %d", VarSize)
// 	}
//
// 	dst := Variable{}
// 	err := binary.Read(bytes.NewBuffer(buf[:]), binary.LittleEndian, &dst)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &dst, nil
// }

// this function will read the variables
// func parseVariables(i *IBT) error {
// 	i.parseVariableHeaders(0)
//   vb := findLatestBuffer(i)
//   fmt.Printf("%+v\n", vb)
//   if i.Vars.LastVersion < vb.tickCount {
//     // Then we have new data
//     i.Vars.LastVersion = vb.tickCount
//     i.LastValidData = time.Now().Unix()
// 	  for _, v := range i.Vars.Vars {
// 	  	// fmt.Printf("%s\n", v.ToString())
//       rbuf := make([]byte, VarTypes[int(v.Type)].Size)
//       _, err := i.File.ReadAt(rbuf, int64(vb.bufOffset + int(v.Offset)))
//       if err != nil {
//         log.Fatalf("Reading values: %v\n", err)
//       }
// 	  }
//   }
//
// 	return nil
// }
