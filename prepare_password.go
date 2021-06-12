/*
 * GoFritzBox
 *
 * Copyright (C) 2016-2021 Dametto Luca <https://damettoluca.com>
 *
 * prepare_password.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0 along with GoFritzBox.
 * If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"unicode/utf16"
)

// preparePassword hashes with MD5 the UTF16LE conversion of the parameters
func preparePassword(challenge, password string) string {
	converted := utf16.Encode([]rune(challenge + "-" + password))
	b := make([]byte, 2*len(converted))
	for i, v := range converted {
		binary.LittleEndian.PutUint16(b[i*2:], v)
	}

	hasher := md5.New()
	hasher.Write(b)
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
