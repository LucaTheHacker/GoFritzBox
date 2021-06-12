/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * stats.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

// GetStats returns Stats to build the usage graph
func (s *SessionInfo) GetStats() (Stats, error) {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	request.SetRequestURI(fmt.Sprintf(
		"%s/internet/inetstat_monitor.lua?sid=%s&myXhr=1&action=get_graphic&useajax=1&xhr=1&t%d=nocache",
		s.EndPoint, s.SID, time.Now().Unix(),
	))

	err := client.Do(request, response)
	if err != nil {
		return Stats{}, err
	}

	var result []Stats
	err = json.Unmarshal(response.Body(), &result)
	if err != nil {
		return Stats{}, err
	}
	return result[0], nil
}
