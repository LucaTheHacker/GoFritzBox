/*
 * GoFritzBox
 * Copyright (C) 2020-2020 Dametto Luca <https://damettoluca.com>
 *
 * functions.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0
 * along with GoFritzBox. If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf16"
)

// Login does a login on the Fritz!Box using infos of Connection
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
	resp.Body.Close()

	resp, err = http.Get(endpoint + "/login_sid.lua?response=" + prelogin.Challenge + "-" + preparePassword(prelogin.Challenge, password) + "&username=" + username)
	if err != nil {
		return SessionInfo{}, err
	}

	body, err = ioutil.ReadAll(resp.Body)
	var login SessionInfo
	err = xml.Unmarshal(body, &login)
	if err != nil {
		return SessionInfo{}, err
	}
	resp.Body.Close()

	if login.SID != "0000000000000000" {
		return login, nil
	} else {
		return login, errors.New("failed to login, try again in " + strconv.Itoa(login.BlockTime) + " second(s)")
	}
}

func (s *SessionInfo) LoadInfos() {

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
