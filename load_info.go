/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * load_info.go is part of GoFritzBox
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
)

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

	return result.Data, nil
}
