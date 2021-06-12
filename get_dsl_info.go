/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * get_dsl_info.go is part of GoFritzBox
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

func (s *SessionInfo) GetDSLInfo() (*ConnectionData, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	request.SetRequestURI(fmt.Sprintf("%s/data.lua", s.EndPoint))
	request.SetBodyString(fmt.Sprintf("xhr=1&sid=%s&lang=it&page=dslOv&xhrId=all", s.SID))
	request.Header.SetMethod(fasthttp.MethodPost)

	err := client.Do(request, response)
	if err != nil {
		return &ConnectionData{}, err
	}

	var result RequestData
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return &ConnectionData{}, err
	}

	return result.Data.ConnectionData, nil
}
