/*
 * GoFritzBox
 * Copyright (C) 2020-2020 Dametto Luca <https://damettoluca.com>
 *
 * functions.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"
)

// Login does a login on the Fritz!Box
// Returns SessionInfo in case of success
func Login(endpoint, username, password string) (SessionInfo, error) {
	resp, err := http.Get(endpoint + "/login_sid.lua")
	if err != nil {
		return SessionInfo{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	var prelogin SessionInfo
	err = xml.Unmarshal(body, &prelogin)
	if err != nil {
		return SessionInfo{}, err
	}
	_ = resp.Body.Close()

	url := fmt.Sprintf("%s/login_sid.lua?response=%s-%s&username=%s", endpoint, prelogin.Challenge, preparePassword(prelogin.Challenge, password), username)
	resp, err = http.Get(url)
	if err != nil {
		return SessionInfo{}, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return SessionInfo{}, err
	}
	var login SessionInfo
	err = xml.Unmarshal(body, &login)
	if err != nil {
		return SessionInfo{}, err
	}
	_ = resp.Body.Close()

	if login.SID != "0000000000000000" {
		login.EndPoint = endpoint
		return login, nil
	} else {
		return login, errors.New("failed to login, try again in " + strconv.Itoa(login.BlockTime) + " second(s)")
	}
}

// LoadInfo returns general Data about the Fritz!Box
func (s *SessionInfo) LoadInfo() (Data, error) {
	url := fmt.Sprintf("%s/data.lua?sid=%s&xhr=1&lang=it&page=overview&xhrId=first&noMenuRef=1&no_sidrenew=", s.EndPoint, s.SID)
	resp, err := http.Get(url)
	if err != nil {
		return Data{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Data{}, err
	}

	var result RequestData
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Data{}, err
	}
	err = json.Unmarshal(*result.Data.WLanRaw, &result.Data.WLanBool)
	if err != nil {
		return Data{}, err
	}
	return result.Data, nil
}

// GetStats returns Stats to build the usage graph
func (s *SessionInfo) GetStats() (Stats, error) {
	url := fmt.Sprintf("%s/internet/inetstat_monitor.lua?sid=%s&myXhr=1&action=get_graphic&useajax=1&xhr=1&t%d=nocache", s.EndPoint, s.SID, time.Now().Unix())
	resp, err := http.Get(url)
	if err != nil {
		return Stats{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Stats{}, err
	}

	var result []Stats
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Stats{}, err
	}
	return result[0], nil
}

// GetLogs returns Logs of the Fritz!Box activity
func (s *SessionInfo) GetLogs() (Logs, error) {
	url := fmt.Sprintf("%s/data.lua", s.EndPoint)
	payload := fmt.Sprintf("sid=%s&page=log&lang=%s&xhr=1&xhrId=all", s.SID, s.Lang)
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(payload))
	if err != nil {
		return Logs{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Logs{}, err
	}

	var result RequestData
	err = json.Unmarshal(body, &result)
	if err != nil {
		return Logs{}, err
	}
	return *result.Data.Log, nil
}

// Disconnect disconnects your Fritz!Box from the internet
// This is usually used to change your IP address
// The prodecure can require up to 30 seconds, after that the internet connection will be re-enabled
func (s *SessionInfo) Disconnect() error {
	url := fmt.Sprintf("%s/internet/inetstat_monitor.lua?sid=%s&myXhr=1&action=disconnect&useajax=1&xhr=1&t%d=nocache", s.EndPoint, s.SID, time.Now().Unix())
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if string(body) == "done:0" {
		return nil
	} else {
		return errors.New("failed to disconnect")
	}
}

// preparePassword hashes with MD5 the UTF16LE conversion of the parameters
func preparePassword(challenge, password string) string {
	converted := utf16.Encode([]rune(challenge + "-" + password))
	b := make([]byte, 2*len(converted))
	for i, v := range converted {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}

	hasher := md5.New()
	hasher.Write(b)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
