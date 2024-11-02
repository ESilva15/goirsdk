package main

type Msg struct {
	Cmd int
	P1  int32
	P2  int32
	P3  int32
}

const (
	SimStatusUrl             string = "http://127.0.0.1:32034/get_sim_status?object=simStatus"
	IRSDK_DATAVALIDEVENTNAME string = "Local\\IRSDKDataValidEvent"
	IRSDK_MEMMAPFILENAME     string = "Local\\IRSDKMemMapFileName"
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
const (
// CamNose
)


// Some other constants that I need to get trough
// // bit fields
// enum irsdk_EngineWarnings 
// {
// 	irsdk_waterTempWarning		= 0x01,
// 	irsdk_fuelPressureWarning	= 0x02,
// 	irsdk_oilPressureWarning	= 0x04,
// 	irsdk_engineStalled			= 0x08,
// 	irsdk_pitSpeedLimiter		= 0x10,
// 	irsdk_revLimiterActive		= 0x20,
// };
//
// // global flags
// enum irsdk_Flags
// {
// 	// global flags
// 	irsdk_checkered				= 0x00000001,
// 	irsdk_white					= 0x00000002,
// 	irsdk_green					= 0x00000004,
// 	irsdk_yellow				= 0x00000008,
// 	irsdk_red					= 0x00000010,
// 	irsdk_blue					= 0x00000020,
// 	irsdk_debris				= 0x00000040,
// 	irsdk_crossed				= 0x00000080,
// 	irsdk_yellowWaving			= 0x00000100,
// 	irsdk_oneLapToGreen			= 0x00000200,
// 	irsdk_greenHeld				= 0x00000400,
// 	irsdk_tenToGo				= 0x00000800,
// 	irsdk_fiveToGo				= 0x00001000,
// 	irsdk_randomWaving			= 0x00002000,
// 	irsdk_caution				= 0x00004000,
// 	irsdk_cautionWaving			= 0x00008000,
//
// 	// drivers black flags
// 	irsdk_black					= 0x00010000,
// 	irsdk_disqualify			= 0x00020000,
// 	irsdk_servicible			= 0x00040000, // car is allowed service (not a flag)
// 	irsdk_furled				= 0x00080000,
// 	irsdk_repair				= 0x00100000,
//
// 	// start lights
// 	irsdk_startHidden			= 0x10000000,
// 	irsdk_startReady			= 0x20000000,
// 	irsdk_startSet				= 0x40000000,
// 	irsdk_startGo				= 0x80000000,
// };
//
//
// // status 
// enum irsdk_TrkLoc
// {
// 	irsdk_NotInWorld = -1,
// 	irsdk_OffTrack,
// 	irsdk_InPitStall,
// 	irsdk_AproachingPits,
// 	irsdk_OnTrack
// };
//
// enum irsdk_TrkSurf
// {
// 	irsdk_SurfaceNotInWorld = -1,
// 	irsdk_UndefinedMaterial = 0,
//
// 	irsdk_Asphalt1Material,
// 	irsdk_Asphalt2Material,
// 	irsdk_Asphalt3Material,
// 	irsdk_Asphalt4Material,
// 	irsdk_Concrete1Material,
// 	irsdk_Concrete2Material,
// 	irsdk_RacingDirt1Material,
// 	irsdk_RacingDirt2Material,
// 	irsdk_Paint1Material,
// 	irsdk_Paint2Material,
// 	irsdk_Rumble1Material,
// 	irsdk_Rumble2Material,
// 	irsdk_Rumble3Material,
// 	irsdk_Rumble4Material,
//
// 	irsdk_Grass1Material,
// 	irsdk_Grass2Material,
// 	irsdk_Grass3Material,
// 	irsdk_Grass4Material,
// 	irsdk_Dirt1Material,
// 	irsdk_Dirt2Material,
// 	irsdk_Dirt3Material,
// 	irsdk_Dirt4Material,
// 	irsdk_SandMaterial,
// 	irsdk_Gravel1Material,
// 	irsdk_Gravel2Material,
// 	irsdk_GrasscreteMaterial,
// 	irsdk_AstroturfMaterial,
// };
//
// enum irsdk_SessionState
// {
// 	irsdk_StateInvalid,
// 	irsdk_StateGetInCar,
// 	irsdk_StateWarmup,
// 	irsdk_StateParadeLaps,
// 	irsdk_StateRacing,
// 	irsdk_StateCheckered,
// 	irsdk_StateCoolDown
// };
//
// enum irsdk_CameraState
// {
// 	irsdk_IsSessionScreen          = 0x0001, // the camera tool can only be activated if viewing the session screen (out of car)
// 	irsdk_IsScenicActive           = 0x0002, // the scenic camera is active (no focus car)
//
// 	//these can be changed with a broadcast message
// 	irsdk_CamToolActive            = 0x0004,
// 	irsdk_UIHidden                 = 0x0008,
// 	irsdk_UseAutoShotSelection     = 0x0010,
// 	irsdk_UseTemporaryEdits        = 0x0020,
// 	irsdk_UseKeyAcceleration       = 0x0040,
// 	irsdk_UseKey10xAcceleration    = 0x0080,
// 	irsdk_UseMouseAimMode          = 0x0100
// };
//
// enum irsdk_PitSvFlags
// {
// 	irsdk_LFTireChange		= 0x0001,
// 	irsdk_RFTireChange		= 0x0002,
// 	irsdk_LRTireChange		= 0x0004,
// 	irsdk_RRTireChange		= 0x0008,
//
// 	irsdk_FuelFill			= 0x0010,
// 	irsdk_WindshieldTearoff	= 0x0020,
// 	irsdk_FastRepair		= 0x0040
// };

//----
//
