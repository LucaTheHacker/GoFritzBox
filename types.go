/*
 * GoFritzBox
 * Copyright (C) 2020-2020 Dametto Luca <https://damettoluca.com>
 *
 * types.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0
 * along with GoFritzBox. If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SessionInfo struct {
	SID       string `xml:"SID"`
	Challenge string `xml:"Challenge"`
	BlockTime int    `xml:"BlockTime"`
	// Rights interface{} Not implemented
}

type RequestData struct {
	PID  string `json:"pid"`
	Data Data   `json:"data"`
	Hide Hide   `json:"hide"`
	Time Time   `json:"time"`
	SID  string `json:"sid"`
}

type Data struct {
	NasLink    string   `json:"naslink"`
	FritzOS    FritzOS  `json:"fritzos"`
	Webdav     string   `json:"webdav"`
	Manual     string   `json:"MANUAL_URL"`
	Language   string   `json:"language"`
	AVM        string   `json:"AVM_URL"`
	USBConnect string   `json:"usbconnect"`
	Foncalls   Foncalls `json:"foncalls"`
	VPN        VPN      `json:"vpn"`
	Internet   Internet `json:"internet"`
	DSL        DSL      `json:"dsl"`
}

type Hide struct {
}

type Time struct {
}

type FritzOS struct {
	Name           string `json:"Productname"`
	NoPWD          bool   `json:"NoPwd"`
	Defaults       bool   `json:"ShowDefaults"`
	Expert         string `json:"expert_mode"`
	FBName         string `json:"fb_name"`
	Version        string `json:"nspver"`
	Labor          bool   `json:"isLabor"`
	TFADisabled    bool   `json:"twofactor_disabled"`
	FirmwareSigned bool   `json:"FirmwareSigned"`
	ShowUpdate     bool   `json:"showUpdate"`
	Updatable      bool   `json:"isUpdateAvail"`
	Energy         int    `json:"energy"`
	BoxDate        string `json:"boxDate"`
}

type Foncalls struct {
	ActiveCalls string `json:"activecalls"`
	Calls       string `json:"calls"`
	CallsToday  string `json:"callsToday"`
	Count       int    `json:"count_all"`
	CountToday  int    `json:"count_today"`
}

type VPN struct {
	Elements []interface{} `json:"elements"`
	Title    string        `json:"title"`
	Link     string        `json:"link"`
}

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
	case "Tbit/s":
		uploadMultiplier = 1000000000000
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
	case "Tbit/s":
		downloadMultiplier = 1000000000000
	default:
		downloadMultiplier = 1000000
	}
	downloadSpeed, _ := strconv.Atoi(strings.ReplaceAll(download[0], ",", ""))
	i.Download = int64(downloadSpeed / 10 * downloadMultiplier)
}

type DSL struct {
	Txt string `json:""`
}
