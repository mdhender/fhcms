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

package main

// argh on endianess
type bit_flags [16]uint8 // this is the same size as [NUM_CONTACTS]long in the C source

func (b bit_flags) IsSet(species_number int) bool {
	/* Get array index and bit mask. */
	//species_array_index := (species_number - 1) / 32
	//species_bit_number := (species_number - 1) % 32
	//var species_bit_mask uint32 = 1 << species_bit_number
	switch (species_number - 1) / 8 {
	case 0: // 1..8
		return (b[2] & (1 << ((species_number - 1) % 8))) != 0
	case 1: // 9..16
		return (b[3] & (1 << ((species_number - 1) % 8))) != 0
	case 2: // 17..24
		return (b[4] & (1 << ((species_number - 1) % 8))) != 0
	default: // 25..MAX_SPECIES
		return false // we have not thought this through!
	}
}
