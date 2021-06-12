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
	"io/ioutil"
	"net/http"
	"strings"
)

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
