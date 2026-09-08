package goirsdk

import (
	"log"
	"strconv"
)

type Msg struct {
	Cmd int
	P1  int32
	P2  int32
	P3  int32
}

const (
	MEMMAPFILENAME                  = "IRSDKMemMapFileName"
	SimStatusUrl             string = "http://127.0.0.1:32034/get_sim_status?object=simStatus"
	IRSDK_DATAVALIDEVENTNAME string = "Local\\IRSDKDataValidEvent"
	IRSDK_MEMMAPFILENAME     string = "Local\\" + MEMMAPFILENAME
	IRSDK_BROADCASTMSGNAME   string = "IRSDK_BROADCASTMSG"
	fileMapSize              uint32 = 1164 * 1024
	connTimeout              int64  = 30
	stConnected              int    = 1
)

const (
	BroadcastCamSwitchPos            int = 0  // car position, group, camera
	BroadcastCamSwitchNum            int = 1  // driver #, group, camera
	BroadcastCamSetState             int = 2  // irsdk_CameraState, unused, unused
	BroadcastReplaySetPlaySpeed      int = 3  // speed, slowMotion, unused
	BroadcastReplaySetPlayPosition   int = 4  // irsdk_RpyPosMode, Frame Number (high, low)
	BroadcastReplaySearch            int = 5  // irsdk_RpySrchMode, unused, unused
	BroadcastReplaySetState          int = 6  // irsdk_RpyStateMode, unused, unused
	BroadcastReloadTextures          int = 7  // irsdk_ReloadTexturesMode, carIdx, unused
	BroadcastChatComand              int = 8  // irsdk_ChatCommandMode, subCommand, unused
	BroadcastPitCommand              int = 9  // irsdk_PitCommandMode, parameter
	BroadcastTelemCommand            int = 10 // irsdk_TelemCommandMode, unused, unused
	BroadcastFFBCommand              int = 11 // irsdk_FFBCommandMode, value (float, high, low)
	BroadcastReplaySearchSessionTime int = 12 // sessionNum, sessionTimeMS (high, low)
	BroadcastLast                    int = 13 // unused placeholder
)

const (
	ChatCommandMacro     int = 0 // pass in a number from 1-15 representing the chat macro to launch
	ChatCommandBeginChat int = 1 // Open up a new chat window
	ChatCommandReply     int = 2 // Reply to last private chat
	ChatCommandCancel    int = 3 // Close chat window
)

// this only works when the driver is in the car
const (
	PitCommandClear      int = 0  // Clear all pit checkboxes
	PitCommandWS         int = 1  // Clean the winshield, using one tear off
	PitCommandFuel       int = 2  // Add fuel, optionally specify the amount to add in liters or pass '0' to use existing amount
	PitCommandLF         int = 3  // Change the left front tire, optionally specifying the pressure in KPa or pass '0' to use existing pressure
	PitCommandRF         int = 4  // right front
	PitCommandLR         int = 5  // left rear
	PitCommandRR         int = 6  // right rear
	PitCommandClearTires int = 7  // Clear tire pit checkboxes
	PitCommandFR         int = 8  // Request a fast repair
	PitCommandClearWS    int = 9  // Uncheck Clean the winshield checkbox
	PitCommandClearFR    int = 10 // Uncheck request a fast repair
	PitCommandClearFuel  int = 11 // Uncheck add fuel
)

// You can call this any time, but telemtry only records when driver is in there car
const (
	TelemCommandStop    int = 0 // Turn telemetry recording off
	TelemCommandStart   int = 1 // Turn telemetry recording on
	TelemCommandRestart int = 2 // Write current file to disk and start a new one
)

const (
	RpyStateEraseTape int = 0 // clear any data in the replay tape
	RpyStateLast      int = 1 // unused place holder
)

const (
	ReloadTexturesAll    int = 0 // reload all textuers
	ReloadTexturesCarIdx int = 1 // reload only textures for the specific carIdx
)

// Search replay tape for events
const (
	RpySrchToStart      int = 0
	RpySrchToEnd        int = 1
	RpySrchPrevSession  int = 2
	RpySrchNextSession  int = 3
	RpySrchPrevLap      int = 4
	RpySrchNextLap      int = 5
	RpySrchPrevFrame    int = 6
	RpySrchNextFrame    int = 7
	RpySrchPrevIncident int = 8
	RpySrchNextIncident int = 9
	RpySrchLast         int = 10 // unused placeholder
)

const (
	RpyPosBegin   int = 0
	RpyPosCurrent int = 1
	RpyPosEnd     int = 2
	RpyPosLast    int = 3 // unused placeholder
)

// You can call this any time
const (
	FFBCommandMaxForce int = 0 // Set the maximum force when mapping steering torque force to direct input units (float in Nm)
	FFBCommandLast     int = 1 // unused placeholder
)

// irsdk_BroadcastCamSwitchPos or irsdk_BroadcastCamSwitchNum camera focus defines
// pass these in for the first parameter to select the 'focus at' types in the camera system.
const (
	csFocusAtIncident int = -3
	csFocusAtLeader   int = -2
	csFocusAtExiting  int = -1
	csFocusAtDriver   int = 0 // ctFocusAtDriver + car number...
)

// Camera positions
type bitfieldValue struct {
	Value int
	Name  string
}

// EngineWarnings - Start
// TODO: wtf
var (
	irsdk_WaterTempWarning    int64 = 0x00000001
	irsdk_FuelPressureWarning int64 = 0x00000002
	irsdk_OilPressureWarning  int64 = 0x00000004
	irsdk_EngineStalled       int64 = 0x00000008
	irsdk_PitSpeedLimiter     int64 = 0x00000010
	irsdk_RevLimiterActive    int64 = 0x00000020
	irsdk_AbsActive           int64 = 0x00000100
	// DEPRECATE THESE ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓
	irsdkWaterTempWarning    = bitfieldValue{0x01, "irsdk_waterTempWarning"}
	irsdkFuelPressureWarning = bitfieldValue{0x02, "irsdk_fueldPressureWarning"}
	irsdkOilPressureWarning  = bitfieldValue{0x04, "irsdk_oilPressureWarning"}
	irsdkEngineStalled       = bitfieldValue{0x08, "irsdk_engineStalled"}
	irsdkPitSpeedLimiter     = bitfieldValue{0x10, "irsdk_pitSpeedLimiter"}
	irsdkRevLimiterActive    = bitfieldValue{0x20, "irsdk_revLimiterActive"}
	irsdkAbsActive           = bitfieldValue{0x100, "irsdk_absActive"}
	irsdkEngineWarnings      = []bitfieldValue{
		irsdkWaterTempWarning, irsdkFuelPressureWarning,
		irsdkOilPressureWarning, irsdkEngineStalled, irsdkPitSpeedLimiter, irsdkRevLimiterActive,
		irsdkAbsActive,
	}
)

func (i *IBT) WaterTempWarning() bool {
	val, ok := i.Vars.Vars["EngineWarnings"]
	if !ok {
		log.Fatal("no EngineWarnings")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get EngineWarnings: " + err.Error())
	}

	return bitfield&irsdk_WaterTempWarning != 0
}

func (i *IBT) FuelPressureWarning() bool {
	val, ok := i.Vars.Vars["EngineWarnings"]
	if !ok {
		log.Fatal("no EngineWarnings")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get EngineWarnings: " + err.Error())
	}

	return bitfield&irsdk_FuelPressureWarning != 0
}

func (i *IBT) OilPressureWarning() bool {
	val, ok := i.Vars.Vars["EngineWarnings"]
	if !ok {
		log.Fatal("no EngineWarnings")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get EngineWarnings: " + err.Error())
	}

	return bitfield&irsdk_OilPressureWarning != 0
}

func (i *IBT) EngineStalled() bool {
	val, ok := i.Vars.Vars["EngineWarnings"]
	if !ok {
		log.Fatal("no EngineWarnings")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get EngineWarnings: " + err.Error())
	}

	return bitfield&irsdk_EngineStalled != 0
}

func (i *IBT) PitSpeedLimiter() bool {
	val, ok := i.Vars.Vars["EngineWarnings"]
	if !ok {
		log.Fatal("no EngineWarnings")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get EngineWarnings: " + err.Error())
	}

	return bitfield&irsdk_PitSpeedLimiter != 0
}

func (i *IBT) RevLimiterActive() bool {
	val, ok := i.Vars.Vars["EngineWarnings"]
	if !ok {
		log.Fatal("no EngineWarnings")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get EngineWarnings: " + err.Error())
	}

	return bitfield&irsdk_RevLimiterActive != 0
}

func (i *IBT) AbsActive() bool {
	val, ok := i.Vars.Vars["EngineWarnings"]
	if !ok {
		log.Fatal("no EngineWarnings")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get EngineWarnings: " + err.Error())
	}

	return bitfield&irsdk_AbsActive != 0
}

// EngineWarnings - END

// SessionState - START
const (
	irsdk_StateInvalid    = 0x00
	irsdk_StateGetInCar   = 0x01
	irsdk_StateWarmup     = 0x02
	irsdk_StateParadeLaps = 0x03
	irsdk_StateRacing     = 0x04
	irsdk_StateCheckered  = 0x05
	irsdk_StateCoolDown   = 0x06
)

func (i *IBT) SessionStateInvalid() bool {
	val, ok := i.Vars.Vars["SessionState"]
	if !ok {
		log.Fatal("no SessionState")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get SessionState: " + err.Error())
	}

	return bitfield == irsdk_StateInvalid
}

func (i *IBT) SessionStateGetInCar() bool {
	val, ok := i.Vars.Vars["SessionState"]
	if !ok {
		log.Fatal("no SessionState")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get SessionState: " + err.Error())
	}

	return bitfield == irsdk_StateGetInCar
}

func (i *IBT) SessionStateWarmup() bool {
	val, ok := i.Vars.Vars["SessionState"]
	if !ok {
		log.Fatal("no SessionState")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get SessionState: " + err.Error())
	}

	return bitfield == irsdk_StateWarmup
}

func (i *IBT) SessionStateParadeLaps() bool {
	val, ok := i.Vars.Vars["SessionState"]
	if !ok {
		log.Fatal("no SessionState")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get SessionState: " + err.Error())
	}

	return bitfield == irsdk_StateParadeLaps
}

func (i *IBT) SessionStateRacing() bool {
	val, ok := i.Vars.Vars["SessionState"]
	if !ok {
		log.Fatal("no SessionState")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get SessionState: " + err.Error())
	}

	return bitfield == irsdk_StateRacing
}

func (i *IBT) SessionStateCheckered() bool {
	val, ok := i.Vars.Vars["SessionState"]
	if !ok {
		log.Fatal("no SessionState")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get SessionState: " + err.Error())
	}

	return bitfield == irsdk_StateCheckered
}

func (i *IBT) SessionStateCoolDown() bool {
	val, ok := i.Vars.Vars["SessionState"]
	if !ok {
		log.Fatal("no SessionState")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get SessionState: " + err.Error())
	}

	return bitfield == irsdk_StateCoolDown
}

func SessionStateToString(state int) string {
	switch state {
	case irsdk_StateInvalid:
		return "StateInvalid"
	case irsdk_StateGetInCar:
		return "StateGetInCar"
	case irsdk_StateWarmup:
		return "StateWarmup"
	case irsdk_StateParadeLaps:
		return "StateParadeLaps"
	case irsdk_StateRacing:
		return "StateRacing"
	case irsdk_StateCheckered:
		return "StateCheckered"
	case irsdk_StateCoolDown:
		return "StateCoolDown"
	default:
		return "UknownSessionState"
	}
}

// SessionState - END

// TrkLoc const
const (
	irsdk_NotInWorld     = -1
	irsdk_OffTrack       = 0
	irsdk_InPitStall     = 1
	irsdk_AproachingPits = 2
	irsdk_OnTrack        = 3
)

func TrkLocToString(trkloc int) string {
	switch trkloc {
	case irsdk_NotInWorld:
		return "NotInWorld"
	case irsdk_OffTrack:
		return "OffTrack"
	case irsdk_InPitStall:
		return "InPitStall"
	case irsdk_AproachingPits:
		return "AproachingPits"
	case irsdk_OnTrack:
		return "OnTrack"
	default:
		return "UknownTrackLocation"
	}
}

// Flags
const (
	irsdk_checkered     = 0x00000001
	irsdk_white         = 0x00000002
	irsdk_green         = 0x00000004
	irsdk_yellow        = 0x00000008
	irsdk_red           = 0x00000010
	irsdk_blue          = 0x00000020
	irsdk_debris        = 0x00000040
	irsdk_crossed       = 0x00000080
	irsdk_yellowWaving  = 0x00000100
	irsdk_oneLapToGreen = 0x00000200
	irsdk_greenHeld     = 0x00000400
	irsdk_tenToGo       = 0x00000800
	irsdk_fiveToGo      = 0x00001000
	irsdk_randomWaving  = 0x00002000
	irsdk_caution       = 0x00004000
	irsdk_cautionWaving = 0x00008000

	// drivers black flags
	irsdk_black      = 0x00010000
	irsdk_disqualify = 0x00020000
	irsdk_servicible = 0x00040000 // car is allowed service (not a flag)
	irsdk_furled     = 0x00080000
	irsdk_repair     = 0x00100000

	// start lights
	irsdk_startHidden = 0x10000000
	irsdk_startReady  = 0x20000000
	irsdk_startSet    = 0x40000000
	irsdk_startGo     = 0x80000000
)

func FlagToString(flag int) string {
	switch flag {
	case irsdk_checkered:
		return "irsdk_checkered"
	case irsdk_white:
		return "irsdk_white"
	case irsdk_green:
		return "irsdk_green"
	case irsdk_yellow:
		return "irsdk_yellow"
	case irsdk_red:
		return "irsdk_red"
	case irsdk_blue:
		return "irsdk_blue"
	case irsdk_debris:
		return "irsdk_debris"
	case irsdk_crossed:
		return "irsdk_crossed"
	case irsdk_yellowWaving:
		return "irsdk_yellowWaving"
	case irsdk_oneLapToGreen:
		return "irsdk_oneLapToGreen"
	case irsdk_greenHeld:
		return "irsdk_greenHeld"
	case irsdk_tenToGo:
		return "irsdk_tenToGo"
	case irsdk_fiveToGo:
		return "irsdk_fiveToGo"
	case irsdk_randomWaving:
		return "irsdk_randomWaving"
	case irsdk_caution:
		return "irsdk_caution"
	case irsdk_cautionWaving:
		return "irsdk_cautionWaving"
	case irsdk_black:
		return "irsdk_black"
	case irsdk_disqualify:
		return "irsdk_disqualify"
	case irsdk_servicible:
		return "irsdk_servicible"
	case irsdk_furled:
		return "irsdk_furled"
	case irsdk_repair:
		return "irsdk_repair"
	case irsdk_startHidden:
		return "irsdk_startHidden"
	case irsdk_startReady:
		return "irsdk_startReady"
	case irsdk_startSet:
		return "irsdk_startSet"
	case irsdk_startGo:
		return "irsdk_startGo"
	default:
		return "UknownFlag"
	}
}

// enum irsdk_TrkSurf
const (
	irsdk_SurfaceNotInWorld = iota - 1
	irsdk_UndefinedMaterial
	irsdk_Asphalt1Material
	irsdk_Asphalt2Material
	irsdk_Asphalt3Material
	irsdk_Asphalt4Material
	irsdk_Concrete1Material
	irsdk_Concrete2Material
	irsdk_RacingDirt1Material
	irsdk_RacingDirt2Material
	irsdk_Paint1Material
	irsdk_Paint2Material
	irsdk_Rumble1Material
	irsdk_Rumble2Material
	irsdk_Rumble3Material
	irsdk_Rumble4Material
	irsdk_Grass1Material
	irsdk_Grass2Material
	irsdk_Grass3Material
	irsdk_Grass4Material
	irsdk_Dirt1Material
	irsdk_Dirt2Material
	irsdk_Dirt3Material
	irsdk_Dirt4Material
	irsdk_SandMaterial
	irsdk_Gravel1Material
	irsdk_Gravel2Material
	irsdk_GrasscreteMaterial
	irsdk_AstroturfMaterial
)

func TrkSurfToString(surface int) string {
	switch surface {
	case irsdk_SurfaceNotInWorld:
		return "SurfaceNotInWorld"
	case irsdk_UndefinedMaterial:
		return "UndefinedMaterial"
	case irsdk_Asphalt1Material:
		return "Asphalt1Material"
	case irsdk_Asphalt2Material:
		return "Asphalt2Material"
	case irsdk_Asphalt3Material:
		return "Asphalt3Material"
	case irsdk_Asphalt4Material:
		return "Asphalt4Material"
	case irsdk_Concrete1Material:
		return "Concrete1Material"
	case irsdk_Concrete2Material:
		return "Concrete2Material"
	case irsdk_RacingDirt1Material:
		return "RacingDirt1Material"
	case irsdk_RacingDirt2Material:
		return "RacingDirt2Material"
	case irsdk_Paint1Material:
		return "Paint1Material"
	case irsdk_Paint2Material:
		return "Paint2Material"
	case irsdk_Rumble1Material:
		return "Rumble1Material"
	case irsdk_Rumble2Material:
		return "Rumble2Material"
	case irsdk_Rumble3Material:
		return "Rumble3Material"
	case irsdk_Rumble4Material:
		return "Rumble4Material"
	case irsdk_Grass1Material:
		return "Grass1Material"
	case irsdk_Grass2Material:
		return "Grass2Material"
	case irsdk_Grass3Material:
		return "Grass3Material"
	case irsdk_Grass4Material:
		return "Grass4Material"
	case irsdk_Dirt1Material:
		return "Dirt1Material"
	case irsdk_Dirt2Material:
		return "Dirt2Material"
	case irsdk_Dirt3Material:
		return "Dirt3Material"
	case irsdk_Dirt4Material:
		return "Dirt4Material"
	case irsdk_SandMaterial:
		return "SandMaterial"
	case irsdk_Gravel1Material:
		return "Gravel1Material"
	case irsdk_Gravel2Material:
		return "Gravel2Material"
	case irsdk_GrasscreteMaterial:
		return "GrasscreteMaterial"
	case irsdk_AstroturfMaterial:
		return "AstroturfMaterial"
	default:
		return "UknownTrackSurface"
	}
}

// CameraState
const (
	irsdk_IsSessionScreen = 0x0001 // the camera tool can only be activated if viewing the session screen (out of car)
	irsdk_IsScenicActive  = 0x0002 // the scenic camera is active (no focus car)
	// these can be changed with a broadcast message
	irsdk_CamToolActive         = 0x0004
	irsdk_UIHidden              = 0x0008
	irsdk_UseAutoShotSelection  = 0x0010
	irsdk_UseTemporaryEdits     = 0x0020
	irsdk_UseKeyAcceleration    = 0x0040
	irsdk_UseKey10xAcceleration = 0x0080
	irsdk_UseMouseAimMode       = 0x0100
)

func CameraStateToString(state int) string {
	switch state {
	case irsdk_IsSessionScreen:
		return "IsSessionScreen"
	case irsdk_IsScenicActive:
		return "IsScenicActive"
	case irsdk_CamToolActive:
		return "CamToolActive"
	case irsdk_UIHidden:
		return "UIHidden"
	case irsdk_UseAutoShotSelection:
		return "UseAutoShotSelection"
	case irsdk_UseTemporaryEdits:
		return "UseTemporaryEdits"
	case irsdk_UseKeyAcceleration:
		return "UseKeyAcceleration"
	case irsdk_UseKey10xAcceleration:
		return "UseKey10xAcceleration"
	case irsdk_UseMouseAimMode:
		return "UseMouseAimMode"
	default:
		return "UknownCameraState"
	}
}

// PitSvFlags -- Start
const (
	// Tires
	irsdk_LFTireChange = 0x00000001
	irsdk_RFTireChange = 0x00000002
	irsdk_LRTireChange = 0x00000004
	irsdk_RRTireChange = 0x00000008
	// Fuel
	irsdk_FuelFill = 0x00000010

	irsdk_WindshieldTearoff = 0x00000020
	irsdk_FastRepair        = 0x00000040

	// Other pit service request flags
	irsdk_ClearTires = 0x00000080 // Uncheck tire change
	irsdk_ClearWS    = 0x00000100 // Uncheck windshield tearoff
	irsdk_ClearFR    = 0x00000200 // Uncheck FastRepair
	irsdk_ClearFuel  = 0x00000400 // Uncheck refuelling
)

func (i *IBT) LFTireChange() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_LFTireChange != 0
}

func (i *IBT) RFTireChange() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_RFTireChange != 0
}

func (i *IBT) LRTireChange() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_LRTireChange != 0
}

func (i *IBT) RRTireChange() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_RRTireChange != 0
}

func (i *IBT) FuelFill() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_FuelFill != 0
}

func (i *IBT) WindshieldTearoff() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_WindshieldTearoff != 0
}

func (i *IBT) FastRepair() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_FastRepair != 0
}

func (i *IBT) ClearTires() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_ClearTires != 0
}

func (i *IBT) ClearWS() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_ClearWS != 0
}

func (i *IBT) ClearFR() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_ClearFR != 0
}

func (i *IBT) ClearFuel() bool {
	val, ok := i.Vars.Vars["PitSvFlags"]
	if !ok {
		log.Fatal("no PitSvFlags")
	}

	bitfield, err := strconv.ParseInt(val.Value.(string), 0, 64)
	if err != nil {
		log.Fatal("unable to get PitSvFlags: " + err.Error())
	}

	return bitfield&irsdk_ClearFuel != 0
}

// PitSvFlags -- End
