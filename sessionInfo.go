package goirsdk

import (
	"github.com/ESilva15/goirsdk/logger"

	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"gopkg.in/yaml.v3"
)

// SessionInfoYAML is a string with session info in the IBT file
type SessionInfoYAML struct {
	WeekendInfo struct {
		TrackName              string `yaml:"TrackName"`
		TrackID                int    `yaml:"TrackID"`
		TrackLength            string `yaml:"TrackLength"`
		TrackDisplayName       string `yaml:"TrackDisplayName"`
		TrackDisplayShortName  string `yaml:"TrackDisplayShortName"`
		TrackConfigName        string `yaml:"TrackConfigName"`
		TrackCity              string `yaml:"TrackCity"`
		TrackCountry           string `yaml:"TrackCountry"`
		TrackAltitude          string `yaml:"TrackAltitude"`
		TrackLatitude          string `yaml:"TrackLatitude"`
		TrackLongitude         string `yaml:"TrackLongitude"`
		TrackNorthOffset       string `yaml:"TrackNorthOffset"`
		TrackNumTurns          int    `yaml:"TrackNumTurns"`
		TrackPitSpeedLimit     string `yaml:"TrackPitSpeedLimit"`
		TrackType              string `yaml:"TrackType"`
		TrackDirection         string `yaml:"TrackDirection"`
		TrackWeatherType       string `yaml:"TrackWeatherType"`
		TrackSkies             string `yaml:"TrackSkies"`
		TrackSurfaceTemp       string `yaml:"TrackSurfaceTemp"`
		TrackAirTemp           string `yaml:"TrackAirTemp"`
		TrackAirPressure       string `yaml:"TrackAirPressure"`
		TrackWindVel           string `yaml:"TrackWindVel"`
		TrackWindDir           string `yaml:"TrackWindDir"`
		TrackRelativeHumidity  string `yaml:"TrackRelativeHumidity"`
		TrackFogLevel          string `yaml:"TrackFogLevel"`
		TrackCleanup           int    `yaml:"TrackCleanup"`
		TrackDynamicTrack      int    `yaml:"TrackDynamicTrack"`
		TrackVersion           string `yaml:"TrackVersion"`
		SeriesID               int    `yaml:"SeriesID"`
		SeasonID               int    `yaml:"SeasonID"`
		SessionID              int    `yaml:"SessionID"`
		SubSessionID           int    `yaml:"SubSessionID"`
		LeagueID               int    `yaml:"LeagueID"`
		Official               int    `yaml:"Official"`
		RaceWeek               int    `yaml:"RaceWeek"`
		EventType              string `yaml:"EventType"`
		Category               string `yaml:"Category"`
		SimMode                string `yaml:"SimMode"`
		TeamRacing             int    `yaml:"TeamRacing"`
		MinDrivers             int    `yaml:"MinDrivers"`
		MaxDrivers             int    `yaml:"MaxDrivers"`
		DCRuleSet              string `yaml:"DCRuleSet"`
		QualifierMustStartRace int    `yaml:"QualifierMustStartRace"`
		NumCarClasses          int    `yaml:"NumCarClasses"`
		NumCarTypes            int    `yaml:"NumCarTypes"`
		HeatRacing             int    `yaml:"HeatRacing"`
		BuildType              string `yaml:"BuildType"`
		BuildTarget            string `yaml:"BuildTarget"`
		BuildVersion           string `yaml:"BuildVersion"`
		WeekendOptions         struct {
			NumStarters                int    `yaml:"NumStarters"`
			StartingGrid               string `yaml:"StartingGrid"`
			QualifyScoring             string `yaml:"QualifyScoring"`
			CourseCautions             string `yaml:"CourseCautions"`
			StandingStart              int    `yaml:"StandingStart"`
			ShortParadeLap             int    `yaml:"ShortParadeLap"`
			Restarts                   string `yaml:"Restarts"`
			WeatherType                string `yaml:"WeatherType"`
			Skies                      string `yaml:"Skies"`
			WindDirection              string `yaml:"WindDirection"`
			WindSpeed                  string `yaml:"WindSpeed"`
			WeatherTemp                string `yaml:"WeatherTemp"`
			RelativeHumidity           string `yaml:"RelativeHumidity"`
			FogLevel                   string `yaml:"FogLevel"`
			TimeOfDay                  string `yaml:"TimeOfDay"`
			Date                       string `yaml:"Date"`
			EarthRotationSpeedupFactor int    `yaml:"EarthRotationSpeedupFactor"`
			Unofficial                 int    `yaml:"Unofficial"`
			CommercialMode             string `yaml:"CommercialMode"`
			NightMode                  string `yaml:"NightMode"`
			IsFixedSetup               int    `yaml:"IsFixedSetup"`
			StrictLapsChecking         string `yaml:"StrictLapsChecking"`
			HasOpenRegistration        int    `yaml:"HasOpenRegistration"`
			HardcoreLevel              int    `yaml:"HardcoreLevel"`
			NumJokerLaps               int    `yaml:"NumJokerLaps"`
			IncidentLimit              string `yaml:"IncidentLimit"`
			FastRepairsLimit           string `yaml:"FastRepairsLimit"`
			GreenWhiteCheckeredLimit   int    `yaml:"GreenWhiteCheckeredLimit"`
		} `yaml:"WeekendOptions"`
		TelemetryOptions struct {
			TelemetryDiskFile string `yaml:"TelemetryDiskFile"`
		} `yaml:"TelemetryOptions"`
	} `yaml:"WeekendInfo"`
	SessionInfo struct {
		Sessions []struct {
			SessionNum              int         `yaml:"SessionNum"`
			SessionLaps             string      `yaml:"SessionLaps"`
			SessionTime             string      `yaml:"SessionTime"`
			SessionNumLapsToAvg     int         `yaml:"SessionNumLapsToAvg"`
			SessionType             string      `yaml:"SessionType"`
			SessionTrackRubberState string      `yaml:"SessionTrackRubberState"`
			SessionName             string      `yaml:"SessionName"`
			SessionSubType          interface{} `yaml:"SessionSubType"`
			SessionSkipped          int         `yaml:"SessionSkipped"`
			SessionRunGroupsUsed    int         `yaml:"SessionRunGroupsUsed"`
			ResultsPositions        interface{} `yaml:"ResultsPositions"`
			ResultsFastestLap       []struct {
				CarIdx      int `yaml:"CarIdx"`
				FastestLap  int `yaml:"FastestLap"`
				FastestTime int `yaml:"FastestTime"`
			} `yaml:"ResultsFastestLap"`
			ResultsAverageLapTime  int `yaml:"ResultsAverageLapTime"`
			ResultsNumCautionFlags int `yaml:"ResultsNumCautionFlags"`
			ResultsNumCautionLaps  int `yaml:"ResultsNumCautionLaps"`
			ResultsNumLeadChanges  int `yaml:"ResultsNumLeadChanges"`
			ResultsLapsComplete    int `yaml:"ResultsLapsComplete"`
			ResultsOfficial        int `yaml:"ResultsOfficial"`
		} `yaml:"Sessions"`
	} `yaml:"SessionInfo"`
	CameraInfo struct {
		Groups []struct {
			GroupNum  int    `yaml:"GroupNum"`
			GroupName string `yaml:"GroupName"`
			Cameras   []struct {
				CameraNum  int    `yaml:"CameraNum"`
				CameraName string `yaml:"CameraName"`
			} `yaml:"Cameras"`
			IsScenic bool `yaml:"IsScenic,omitempty"`
		} `yaml:"Groups"`
	} `yaml:"CameraInfo"`
	RadioInfo struct {
		SelectedRadioNum int `yaml:"SelectedRadioNum"`
		Radios           []struct {
			RadioNum            int `yaml:"RadioNum"`
			HopCount            int `yaml:"HopCount"`
			NumFrequencies      int `yaml:"NumFrequencies"`
			TunedToFrequencyNum int `yaml:"TunedToFrequencyNum"`
			ScanningIsOn        int `yaml:"ScanningIsOn"`
			Frequencies         []struct {
				FrequencyNum  int    `yaml:"FrequencyNum"`
				FrequencyName string `yaml:"FrequencyName"`
				Priority      int    `yaml:"Priority"`
				CarIdx        int    `yaml:"CarIdx"`
				EntryIdx      int    `yaml:"EntryIdx"`
				ClubID        int    `yaml:"ClubID"`
				CanScan       int    `yaml:"CanScan"`
				CanSquawk     int    `yaml:"CanSquawk"`
				Muted         int    `yaml:"Muted"`
				IsMutable     int    `yaml:"IsMutable"`
				IsDeletable   int    `yaml:"IsDeletable"`
			} `yaml:"Frequencies"`
		} `yaml:"Radios"`
	} `yaml:"RadioInfo"`
	DriverInfo struct {
		DriverCarIdx              int      `yaml:"DriverCarIdx"`
		DriverUserID              int      `yaml:"DriverUserID"`
		PaceCarIdx                int      `yaml:"PaceCarIdx"`
		DriverHeadPosX            float64  `yaml:"DriverHeadPosX"`
		DriverHeadPosY            float64  `yaml:"DriverHeadPosY"`
		DriverHeadPosZ            float64  `yaml:"DriverHeadPosZ"`
		DriverCarIdleRPM          float64  `yaml:"DriverCarIdleRPM"`
		DriverCarRedLine          float64  `yaml:"DriverCarRedLine"`
		DriverCarEngCylinderCount int      `yaml:"DriverCarEngCylinderCount"`
		DriverCarFuelKgPerLtr     float64  `yaml:"DriverCarFuelKgPerLtr"`
		DriverCarFuelMaxLtr       float64  `yaml:"DriverCarFuelMaxLtr"`
		DriverCarMaxFuelPct       float64  `yaml:"DriverCarMaxFuelPct"`
		DriverCarGearNumForward   int      `yaml:"DriverCarGearNumForward"`
		DriverCarGearNeutral      int      `yaml:"DriverCarGearNeutral"`
		DriverCarGearReverse      int      `yaml:"DriverCarGearReverse"`
		DriverCarSLFirstRPM       float64  `yaml:"DriverCarSLFirstRPM"`
		DriverCarSLShiftRPM       float64  `yaml:"DriverCarSLShiftRPM"`
		DriverCarSLLastRPM        float64  `yaml:"DriverCarSLLastRPM"`
		DriverCarSLBlinkRPM       float64  `yaml:"DriverCarSLBlinkRPM"`
		DriverCarVersion          string   `yaml:"DriverCarVersion"`
		DriverPitTrkPct           float64  `yaml:"DriverPitTrkPct"`
		DriverCarEstLapTime       float64  `yaml:"DriverCarEstLapTime"`
		DriverSetupName           string   `yaml:"DriverSetupName"`
		DriverSetupIsModified     int      `yaml:"DriverSetupIsModified"`
		DriverSetupLoadTypeName   string   `yaml:"DriverSetupLoadTypeName"`
		DriverSetupPassedTech     int      `yaml:"DriverSetupPassedTech"`
		DriverIncidentCount       int      `yaml:"DriverIncidentCount"`
		Drivers                   []Driver `yaml:"Drivers"`
	} `yaml:"DriverInfo"`
	SplitTimeInfo struct {
		Sectors []struct {
			SectorNum      int     `yaml:"SectorNum"`
			SectorStartPct float64 `yaml:"SectorStartPct"`
		} `yaml:"Sectors"`
	} `yaml:"SplitTimeInfo"`
	CarSetup struct {
		UpdateCount int `yaml:"UpdateCount"`
		TiresAero   struct {
			LeftFront struct {
				StartingPressure string `yaml:"StartingPressure"`
				LastHotPressure  string `yaml:"LastHotPressure"`
				LastTempsOMI     string `yaml:"LastTempsOMI"`
				TreadRemaining   string `yaml:"TreadRemaining"`
			} `yaml:"LeftFront"`
			LeftRear struct {
				StartingPressure string `yaml:"StartingPressure"`
				LastHotPressure  string `yaml:"LastHotPressure"`
				LastTempsOMI     string `yaml:"LastTempsOMI"`
				TreadRemaining   string `yaml:"TreadRemaining"`
			} `yaml:"LeftRear"`
			RightFront struct {
				StartingPressure string `yaml:"StartingPressure"`
				LastHotPressure  string `yaml:"LastHotPressure"`
				LastTempsIMO     string `yaml:"LastTempsIMO"`
				TreadRemaining   string `yaml:"TreadRemaining"`
			} `yaml:"RightFront"`
			RightRear struct {
				StartingPressure string `yaml:"StartingPressure"`
				LastHotPressure  string `yaml:"LastHotPressure"`
				LastTempsIMO     string `yaml:"LastTempsIMO"`
				TreadRemaining   string `yaml:"TreadRemaining"`
			} `yaml:"RightRear"`
		} `yaml:"TiresAero"`
		Chassis struct {
			Front struct {
				ArbSetting  int    `yaml:"ArbSetting"`
				ToeIn       string `yaml:"ToeIn"`
				FuelLevel   string `yaml:"FuelLevel"`
				CrossWeight string `yaml:"CrossWeight"`
			} `yaml:"Front"`
			LeftFront struct {
				CornerWeight      string `yaml:"CornerWeight"`
				RideHeight        string `yaml:"RideHeight"`
				SpringPerchOffset string `yaml:"SpringPerchOffset"`
				Camber            string `yaml:"Camber"`
			} `yaml:"LeftFront"`
			LeftRear struct {
				CornerWeight      string `yaml:"CornerWeight"`
				RideHeight        string `yaml:"RideHeight"`
				SpringPerchOffset string `yaml:"SpringPerchOffset"`
				Camber            string `yaml:"Camber"`
				ToeIn             string `yaml:"ToeIn"`
			} `yaml:"LeftRear"`
			InCarDials struct {
				DisplayPage       string `yaml:"DisplayPage"`
				BrakePressureBias string `yaml:"BrakePressureBias"`
			} `yaml:"InCarDials"`
			RightFront struct {
				CornerWeight      string `yaml:"CornerWeight"`
				RideHeight        string `yaml:"RideHeight"`
				SpringPerchOffset string `yaml:"SpringPerchOffset"`
				Camber            string `yaml:"Camber"`
			} `yaml:"RightFront"`
			RightRear struct {
				CornerWeight      string `yaml:"CornerWeight"`
				RideHeight        string `yaml:"RideHeight"`
				SpringPerchOffset string `yaml:"SpringPerchOffset"`
				Camber            string `yaml:"Camber"`
				ToeIn             string `yaml:"ToeIn"`
			} `yaml:"RightRear"`
			Rear struct {
				ArbSetting  int `yaml:"ArbSetting"`
				WingSetting int `yaml:"WingSetting"`
			} `yaml:"Rear"`
		} `yaml:"Chassis"`
	} `yaml:"CarSetup"`
}

// Driver ...
type Driver struct {
	CarIdx                  int     `yaml:"CarIdx"`
	UserName                string  `yaml:"UserName"`
	AbbrevName              string  `yaml:"AbbrevName"`
	Initials                string  `yaml:"Initials"`
	UserID                  int     `yaml:"UserID"`
	TeamID                  int     `yaml:"TeamID"`
	TeamName                string  `yaml:"TeamName"`
	CarNumber               string  `yaml:"CarNumber"`
	CarNumberRaw            int     `yaml:"CarNumberRaw"`
	CarPath                 string  `yaml:"CarPath"`
	CarClassID              int     `yaml:"CarClassID"`
	CarID                   int     `yaml:"CarID"`
	CarIsPaceCar            int     `yaml:"CarIsPaceCar"`
	CarIsAI                 int     `yaml:"CarIsAI"`
	CarScreenName           string  `yaml:"CarScreenName"`
	CarScreenNameShort      string  `yaml:"CarScreenNameShort"`
	CarClassShortName       string  `yaml:"CarClassShortName"`
	CarClassRelSpeed        int     `yaml:"CarClassRelSpeed"`
	CarClassLicenseLevel    int     `yaml:"CarClassLicenseLevel"`
	CarClassMaxFuelPct      string  `yaml:"CarClassMaxFuelPct"`
	CarClassWeightPenalty   string  `yaml:"CarClassWeightPenalty"`
	CarClassPowerAdjust     string  `yaml:"CarClassPowerAdjust"`
	CarClassDryTireSetLimit string  `yaml:"CarClassDryTireSetLimit"`
	CarClassColor           int     `yaml:"CarClassColor"`
	CarClassEstLapTime      float64 `yaml:"CarClassEstLapTime"`
	IRating                 int     `yaml:"IRating"`
	LicLevel                int     `yaml:"LicLevel"`
	LicSubLevel             int     `yaml:"LicSubLevel"`
	LicString               string  `yaml:"LicString"`
	LicColor                string  `yaml:"LicColor"`
	IsSpectator             int     `yaml:"IsSpectator"`
	CarDesignStr            string  `yaml:"CarDesignStr"`
	HelmetDesignStr         string  `yaml:"HelmetDesignStr"`
	SuitDesignStr           string  `yaml:"SuitDesignStr"`
	CarNumberDesignStr      string  `yaml:"CarNumberDesignStr"`
	CarSponsor1             int     `yaml:"CarSponsor_1"`
	CarSponsor2             int     `yaml:"CarSponsor_2"`
	CurDriverIncidentCount  int     `yaml:"CurDriverIncidentCount"`
	TeamIncidentCount       int     `yaml:"TeamIncidentCount"`
}

// readSessionInfo will read the session info yaml out of the telemetry data
func (i *IBT) readSessionInfo() error {
	log := logger.GetInstance()

	sessionInfoStringRaw := make([]byte, i.Headers.SessionInfoLength)
	_, err := i.File.ReadAt(sessionInfoStringRaw, int64(i.Headers.SessionInfoOffset))
	if err != nil {
		return fmt.Errorf("Failed to read sessionInfoString from file: %v", err)
	}

	// Write to the output file
	if i.IBTExport != nil {
		err := i.exportIBT(sessionInfoStringRaw[:], int64(i.Headers.SessionInfoOffset))
		if err != nil {
			log.Printf("Failed to export offline telemetry data: %v", err)
		}
	}

	i.SessionInfo, err = parseSessionInfo(sessionInfoStringRaw, i.Headers.SessionInfoLength)
	if err != nil {
		return fmt.Errorf("Unable to parse SessionInfoString from file: %v", err)
	}

	// Write to YAML output file
	if i.YAMLExportPath != "" {
		err := i.exportYAML()
		if err != nil {
			log.Printf("Failed to export YAML string: %v\n", err)
		}
	}

	return nil
}

// parseSessionInfo will parse the sessionInfo buffer into the SessionInfoYAML
// struct
func parseSessionInfo(buf []byte, len int32) (*SessionInfoYAML, error) {
	// this seems not to work on windows
	// windows := true
	var sessionInfo SessionInfoYAML
	dataBuffer := buf

	// FOR WINDOWS
	// if windows {
	decoder := charmap.Windows1252.NewDecoder()
	buf, err := decoder.Bytes(buf)
	if err != nil {
		log.Fatal(err)
	}
	dataBuffer = []byte(strings.TrimRight(string(buf[:len]), "\x00"))
	// }

	err = yaml.Unmarshal(dataBuffer, &sessionInfo)
	if err != nil {
		return nil, err
	}

	return &sessionInfo, nil
}

// sessionStatusOK will tell us if we are connected to the live data
func sessionStatusOK(status int) bool {
	return (status & stConnected) > 0
}

// ToString will return a readable string of the struct
func (s *SessionInfoYAML) ToString() string {
	stringified, _ := json.MarshalIndent(s, "", "  ")
	return string(stringified)
}
