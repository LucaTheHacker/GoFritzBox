/*
 * GoFritzBox
 * Copyright (C) 2020-2020 Dametto Luca <https://damettoluca.com>
 *
 * types.go is part of GoFritzBox
 *
 * You should have received a copy of the GNU Affero General Public License v3.0
 * along with GoFritzBox. If not, see <https://github.com/LucaTheHacker/GoFritzBox/blob/main/LICENSE>.
 */

package GoFritzBox

type SessionInfo struct {
	SID       string `xml:"SID"`
	Challenge string `xml:"Challenge"`
	BlockTime int    `xml:"BlockTime"`
	// Rights interface{} Not implemented
}
