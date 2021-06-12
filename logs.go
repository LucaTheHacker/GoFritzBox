/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * logs.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

// GetLogs returns Logs of the Fritz!Box activity
func (s *SessionInfo) GetLogs() (Logs, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	request.SetRequestURI(fmt.Sprintf("%s/data.lua", s.EndPoint))
	request.SetBodyString(fmt.Sprintf("sid=%s&page=log&lang=%s&xhr=1&xhrId=all", s.SID, s.Lang))

	err := client.Do(request, response)
	if err != nil {
		return Logs{}, err
	}

	var result RequestData
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return Logs{}, err
	}
	return *result.Data.Log, nil
}
