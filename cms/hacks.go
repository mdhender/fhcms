/*****************************************************************************
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
****************************************************************************/

package main

func byte4_to_byte(k [4]byte) (b []byte) {
	for i := 0; i < 4 && k[i] != 0; i++ {
		b = append(b, k[i])
	}
	return b
}

func gas_abbr_to_code(abbr string) *gas_code {
	for _, code := range gas_codes {
		if abbr == code.abbr {
			return code
		}
	}
	return nil
}

func star_color_to_code(color string) *star_color_code {
	for _, code := range star_color_codes {
		if color == code.name {
			return code
		}
	}
	return nil
}
