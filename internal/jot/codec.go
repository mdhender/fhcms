/*******************************************************************************
Far Horizons Engine
Copyright (C) 2021  Michael D Henderson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
******************************************************************************/

package jot

import "encoding/base64"

// decode is a helper function for converting base64 representation to raw bytes
func decode(raw string) (b []byte, err error) {
	return base64.RawURLEncoding.DecodeString(raw)
}

// encode is a helper function for converting a base64 representation to a string
func encode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}
