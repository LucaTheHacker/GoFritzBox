/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * assistance_data.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"bytes"
	"fmt"
	"mime/multipart"

	"github.com/valyala/fasthttp"
)

// GetAssistanceData returns the firmwarecfg file useful to generate HLog/QLN graphs
func (s *SessionInfo) GetAssistanceData() ([]byte, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("sid", s.SID)
	_ = writer.WriteField("SupportData", "")
	err := writer.Close()
	if err != nil {
		return []byte{}, err
	}

	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	request.SetRequestURI(fmt.Sprintf("%s/cgi-bin/firmwarecfg", s.EndPoint))
	request.Header.SetContentType(writer.FormDataContentType())
	request.Header.SetMethod(fasthttp.MethodPost)

	request.SetBodyRaw(payload.Bytes())

	err = client.Do(request, response)
	if err != nil {
		return []byte{}, err
	}

	return response.Body(), nil
}
