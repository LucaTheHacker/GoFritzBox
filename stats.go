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
	"io/ioutil"
	"net/http"
	"time"
)

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
