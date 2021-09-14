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

package dat32

func nameToString(name [32]uint8) string {
	var b []byte
	for _, ch := range name {
		if ch == 0 {
			break
		}
		b = append(b, byte(ch))
	}
	return string(b)
}

// speciesBitIsSet returns true if the bit is set for the species.
// note: the species number must be 1 based!
// sp01       65536                       1 0000 0000 0000 0000
// sp09    16777216             1 0000 0000 0000 0000 0000 0000
// sp18  8589934592  10 0000 0000 0000 0000 0000 0000 0000 0000
func speciesBitIsSet(set [2]uint64, sp int) bool {
	return (set[0] & (1 << (sp + 15))) != 0
}
