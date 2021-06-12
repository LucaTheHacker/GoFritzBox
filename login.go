/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * login.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"

	"github.com/valyala/fasthttp"
)

// Login does a login on the Fritz!Box
// Returns SessionInfo in case of success
func Login(endpoint, username, password string) (SessionInfo, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	request.SetRequestURI(endpoint + "/login_sid.lua")

	err := client.Do(request, response)
	if err != nil {
		return SessionInfo{}, err
	}

	var prelogin SessionInfo
	err = xml.Unmarshal(response.Body(), &prelogin)
	if err != nil {
		return SessionInfo{}, err
	}

	fasthttp.ReleaseRequest(request)
	fasthttp.ReleaseResponse(response)

	request = fasthttp.AcquireRequest()
	response = fasthttp.AcquireResponse()

	request.SetRequestURI(fmt.Sprintf(
		"%s/login_sid.lua?response=%s-%s&username=%s",
		endpoint, prelogin.Challenge, preparePassword(prelogin.Challenge, password), username,
	))
	err = client.Do(request, response)
	if err != nil {
		return SessionInfo{}, err
	}

	var login SessionInfo
	err = xml.Unmarshal(response.Body(), &login)
	if err != nil {
		return SessionInfo{}, err
	}

	if login.SID != "0000000000000000" {
		login.EndPoint = endpoint
		return login, nil
	} else {
		return login, errors.New("failed to login, try again in " + strconv.Itoa(login.BlockTime) + " second(s)")
	}
}
