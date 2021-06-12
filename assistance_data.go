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
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

// GetAssistanceData returns the firmwarecfg file useful to generate HLog/QLN graphs
func (s *SessionInfo) GetAssistanceData() ([]byte, error) {
	url := fmt.Sprintf("%s/cgi-bin/firmwarecfg", s.EndPoint)
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("sid", s.SID)
	_ = writer.WriteField("SupportData", "")
	err := writer.Close()
	if err != nil {
		return []byte{}, err
	}

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return []byte{}, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}
