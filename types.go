/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * types.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SessionInfo contains info about the current Session
// SID is the session ID
// Challenge is a part of the hash used for the Login
// EndPoint is the endpoint used to do the HTTP requests
// BlockTime is the cooldown needed after a wrong login, you need to handle it properly,
// Lang is the lang that the Fritz!Box will use. Usually it's a useless parameters for API, but not on Fritz!OS
// otherwise all the logins will fail, with the correct authentication as well
type SessionInfo struct {
	SID       string      `xml:"SID"`
	Challenge string      `xml:"Challenge"`
	EndPoint  string      ``
	BlockTime int         `xml:"BlockTime"`
	Lang      string      ``
	Rights    interface{} `` // Not Implemented
}

// RequestData contains data about the usual Fritz!Box answer
type RequestData struct {
	PID  string `json:"pid"`
	Data Data   `json:"data"`
	// Hide []Hide `json:"hide"` Removed for security reasons, check issue #1
	// Time []Time `json:"time"` Removed for security reasons, check issue #1
	SID string `json:"sid"`
}

// Data contains data about the Fritz!Box
// WLan needs to be decoded manually, check the type before unmarhsal to WLan or WLanBool
type Data struct {
	NasLink          string           `json:"naslink"`
	FritzOS          *FritzOS         `json:"fritzos"`
	Webdav           int              `json:"webdav,string"`
	Manual           string           `json:"MANUAL_URL"`
	Language         string           `json:"language"`
	AVM              string           `json:"AVM_URL"`
	USBConnect       string           `json:"usbconnect"`
	Foncalls         *Foncalls        `json:"foncalls"`
	VPN              *VPN             `json:"vpn"`
	Internet         *Internet        `json:"internet"`
	DSL              *DSL             `json:"dsl"`
	ServicePortalURL string           `json:"SERVICEPORTAL_URL"`
	Comfort          *Comfort         `json:"comfort"`
	Changelog        *Changelog       `json:"changelog"`
	TamCalls         *TamCalls        `json:"tamcalls"`
	Lan              *External        `json:"lan"`
	Log              *Logs            `json:"log"`
	Filter           int              `json:"filter,string"`
	USB              *External        `json:"usb"`
	FonNum           *External        `json:"fonnum"`
	NewsURL          string           `json:"NEWSLETTER_URL"`
	Net              *Net             `json:"net"`
	Dect             *External        `json:"dect"`
	WLanRaw          *json.RawMessage `json:"wlan"`
	WLanBool         bool             ``
	WLan             WLan             ``
	ConnectionData   *ConnectionData  `json:"connectionData"`
}

// FritzOS contains infos about the current Fritz!OS version
type FritzOS struct {
	Name           string  `json:"Productname"`
	NoPWD          bool    `json:"NoPwd"`
	Defaults       bool    `json:"ShowDefaults"`
	Expert         int     `json:"expert_mode,string"`
	FBName         string  `json:"fb_name"`
	Version        float32 `json:"nspver,string"`
	Labor          bool    `json:"isLabor"`
	TFADisabled    bool    `json:"twofactor_disabled"`
	FirmwareSigned bool    `json:"FirmwareSigned"`
	ShowUpdate     bool    `json:"showUpdate"`
	Updatable      bool    `json:"isUpdateAvail"`
	Energy         int     `json:"energy,string"`
	BoxDate        string  `json:"boxDate"`
}

// Foncalls contains info about Calls
type Foncalls struct {
	ActiveCalls string `json:"activecalls"`
	Calls       string `json:"-"` // See issue #1
	CallsToday  int    `json:"callsToday,string"`
	Count       int    `json:"count_all"`
	CountToday  int    `json:"count_today"`
}

// VPN contains infos about the VPN
type VPN struct {
	Elements []interface{} `json:"elements"`
	Title    string        `json:"title"`
	Link     string        `json:"link"`
}

// Internet contains data about the internet connection
type Internet struct {
	Txt            []string  `json:"txt"`
	Led            string    `json:"led"`
	Online         bool      ``
	Title          string    `json:"title"`
	DownloadString string    `json:"down"`
	Download       int64     ``
	UploadString   string    `json:"up"`
	Upload         int64     ``
	SecondaryLink  string    `json:"link2"`
	Link           string    `json:"link"`
	Provider       string    ``
	ConnectionTime time.Time ``
}

// Sanitize sanitizes the Internet struct by cleaning bad values and generating data from other values
func (i *Internet) Sanitize() {
	if i.Led == "globe_online" {
		i.Online = true
	}

	i.Provider = strings.SplitN(i.Txt[0], ": ", 2)[1]

	timeparser := regexp.MustCompile("([0-9]{1,2}\\.[0-9]{1,2}\\.[0-9]{1,4}, [0-9]{1,2}:[0-9]{1,2})")
	result := string(timeparser.Find([]byte(i.Txt[1])))
	readtime, err := time.Parse("02.01.2006, 15:04", result)
	if err != nil {
		i.ConnectionTime = readtime
	}

	speedparser := regexp.MustCompile("([0-9]{1,5},[0-9]) (Mbit/s|Gbit/s|Tbit/s)")
	upload := speedparser.FindAllString(i.UploadString, -1)
	var uploadMultiplier int
	switch upload[1] {
	case "Mbit/s":
		uploadMultiplier = 1000000
	case "Gbit/s":
		uploadMultiplier = 1000000000

	default:
		uploadMultiplier = 1000000
	}
	uploadSpeed, _ := strconv.Atoi(strings.ReplaceAll(upload[0], ",", ""))
	i.Upload = int64(uploadSpeed / 10 * uploadMultiplier)

	download := speedparser.FindAllString(i.DownloadString, -1)
	var downloadMultiplier int
	switch download[1] {
	case "Mbit/s":
		downloadMultiplier = 1000000
	case "Gbit/s":
		downloadMultiplier = 1000000000
	default:
		downloadMultiplier = 1000000
	}
	downloadSpeed, _ := strconv.Atoi(strings.ReplaceAll(download[0], ",", ""))
	i.Download = int64(downloadSpeed / 10 * downloadMultiplier)
}

// DSL contains infos about the connection status
type DSL struct {
	Txt         string `json:"txt"`
	Led         string `json:"led"`
	Title       string `json:"title"`
	DiagStopPID string `json:"diag_stop_pid"`
	DiagActive  int    `json:"diag_active,string"`
	AddDiag     string `json:"addDiag"`
	Link        string `json:"link"`
	Upload      string `json:"up"`
	Download    string `json:"down"`
}

// Comfort contains infos about Fritz!Box ComfortFunc
type Comfort struct {
	Functions []ComfortFunc `json:"func"`
	Any       bool          `json:"anyComfort"`
}

// ComfortFunc contains infos about a Comfort Function
type ComfortFunc struct {
	Name    string `json:"linktxt"`
	Details string `json:"details"`
	Link    string `json:"link"`
}

// Changelog contains infos about the Fritz!Box
type Changelog struct {
	DeviceName       string  `json:"deviceName"`
	FritzOSVersion   float32 `json:"fritzOsVersion,string"`
	ConnectionStatus bool    `json:"connectionStatus"`
	ProductName      string  `json:"productName"`
	IFrame           string  `json:"iframeUrl"`
}

type TamCalls struct {
	Calls      string `json:"calls"`
	Configured bool   `json:"tam_configured"`
	Count      int    `json:"count"`
	CallsToday int    `json:"callsToday,string"`
}

type External struct {
	Txt   string `json:"txt"`
	Led   string `json:"led"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

// Net contains infos about Device connected to the Fritz!Box
type Net struct {
	UnmeshedDevices bool      `json:"anyUnmeshedDevices"`
	Count           int       `json:"count"`
	ActiveCount     int       `json:"active_count"`
	More            string    `json:"more_link"`
	Devices         *[]Device `json:"devices"`
}

// Device contains infos about a device in the LAN
type Device struct {
	Classes string `json:"classes"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	URL     string `json:"url"`
}

// WLan contains infos about the Wireless Lan
type WLan struct {
	Txt     string `json:"txt"`
	Led     string `json:"led"`
	Title   string `json:"title"`
	Link    string `json:"link"`
	Tooltip string `json:"tooltip"`
}

type Logs [][]string

// Filter filters Logs by type
// 0 -> All,
// 1 -> System,
// 2 -> Internet,
// 3 -> Phone,
// 4 -> Wifi,
// 5 -> USB
func (l *Logs) Filter(filter int) Logs {
	if filter == 0 {
		return *l
	}
	if filter > 5 {
		return Logs{}
	}

	// Bad code to do a type conversion, but the Fritz!Box returns the value as string
	check := strconv.Itoa(filter)
	var result Logs
	for _, v := range *l {
		if v[4] == check {
			result = append(result, v)
		}
	}
	return result
}

// Stats contains infos about the current connection usage
// These infos are used to build the connection graph
type Stats struct {
	DownstreamMax       int     `json:"ds_bps_curr_max"`
	UpstreamMax         int     `json:"us_bps_curr_max"`
	DownstreamCapacity  int     `json:"downstream"`
	UpstreamCapacity    int     `json:"upstream"`
	StaticDownstreamMax int     `json:"ds_bps_max"`
	StaticUpstreamMax   int     `json:"us_bps_max"`
	Dynamic             bool    `json:"dynamic"`
	Node                string  `json:"_node"`
	Mode                string  `json:"mode"`
	Name                string  `json:"name"`
	DownstreamInternet  [19]int `json:"ds_bps_curr"`
	DownstreamIPTV      [19]int `json:"ds_mc_bps_curr"`
	DownstreamGuest     [19]int `json:"ds_guest_bps_curr"`
	UpstreamRealTime    [19]int `json:"us_realtime_bps_curr"`
	UpstreamPriority    [19]int `json:"us_important_bps_curr"`
	UpstreamNormal      [19]int `json:"us_default19_bps_curr"`
	UpstreamBackground  [19]int `json:"us_background_bps_curr"`
	UpstreamGuest       [19]int `json:"guest_us_bps"`
	DownstreamTotal     [19]int ``
	UpstreamTotal       [19]int ``
}

// Load adds Total values to Stats, useful to get the Total internet usage
func (s *Stats) Load() {
	for i := 0; i <= 18; i++ {
		s.DownstreamTotal[i] = s.DownstreamInternet[i] + s.DownstreamIPTV[i] + s.DownstreamGuest[i]
		s.UpstreamTotal[i] = s.UpstreamRealTime[i] + s.UpstreamPriority[i] + s.UpstreamNormal[i] + s.UpstreamBackground[i] + s.UpstreamGuest[i]
	}
}

type ConnectionData struct {
	ExternApValue  string `json:"externApValue"`
	Modell         string `json:"modell"`
	IsDebug        bool   `json:"isDebug"`
	LineLength     int64  `json:"lineLength"`
	ExternAPHeader string `json:"externAPHeader"`
	ExternApText   string `json:"externApText"`
	Line           []Line `json:"line"`
	Version        string `json:"version"`
	Versiontext    string `json:"versiontext"`
	DsRate         string `json:"dsRate"`
	UsRate         string `json:"usRate"`
}

type Line struct {
	State            string `json:"state"`
	TimePrefix       string `json:"timePrefix"`
	TrainState       string `json:"trainState"`
	Mode             string `json:"mode"`
	TrainStatePrefix string `json:"trainStatePrefix"`
	Time             string `json:"time"`
}
