/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * disconnect.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"errors"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

// Disconnect disconnects your Fritz!Box from the internet
// This is usually used to change your IP address
// The prodecure can require up to 30 seconds, after that the internet connection will be re-enabled
func (s *SessionInfo) Disconnect() error {
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	request.SetRequestURI(fmt.Sprintf(
		"%s/internet/inetstat_monitor.lua?sid=%s&myXhr=1&action=disconnect&useajax=1&xhr=1&t%d=nocache",
		s.EndPoint, s.SID, time.Now().Unix(),
	))

	err := client.Do(request, response)
	if err != nil {
		return err
	}

	if string(response.Body()) == "done:0" {
		return nil
	} else {
		return errors.New("failed to disconnect")
	}
}
