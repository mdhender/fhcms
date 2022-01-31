/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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

package engine

import "fmt"

func (e *Engine) bad_argument() {
	fprintf(e.log_file, "!!! Order ignored:\n")
	fprintf(e.log_file, "!!! %s", b2s(e.input_line))
	fprintf(e.log_file, "!!! Invalid argument in command.\n")
}

func (e *Engine) bad_coordinates() {
	fprintf(e.log_file, "!!! Order ignored:\n")
	fprintf(e.log_file, "!!! %s", b2s(e.input_line))
	fprintf(e.log_file, "!!! Invalid coordinates in command.\n")
}

func (e *Engine) bad_species() {
	fprintf(e.log_file, "!!! Order ignored:\n")
	fprintf(e.log_file, "!!! %s", b2s(e.input_line))
	fprintf(e.log_file, "!!! Invalid species name!\n")
}

func (e *Engine) battle_error(species_number int) {
	fprintf(e.log_file, "!!! Order ignored:\n")
	fprintf(e.log_file, "!!! %s", b2s(e.input_line))
	fprintf(e.log_file, "!!! Missing BATTLE command!\n")
	return
}

func power(tonnage int) int {
	if tonnage < 0 {
		return 0
	} else if tonnage <= 100 {
		/* Look-up table for ship defensive/offensive power uses ship->tonnage
		 * as an index. Each value is equal to 100 * (ship->tonnage)^1.2. The
		 * 'power' subroutine uses recursion to calculate values for tonnages
		 * over 100. */
		ship_power := [101]int{0, /* Zeroth element not used. */
			100, 230, 374, 528, 690, 859, 1033, 1213, 1397, 1585,
			1777, 1973, 2171, 2373, 2578, 2786, 2996, 3209, 3424, 3641,
			3861, 4082, 4306, 4532, 4759, 4988, 5220, 5452, 5687, 5923,
			6161, 6400, 6641, 6883, 7127, 7372, 7618, 7866, 8115, 8365,
			8617, 8870, 9124, 9379, 9635, 9893, 10151, 10411, 10672, 10934,
			11197, 11461, 11725, 11991, 12258, 12526, 12795, 13065, 13336, 13608,
			13880, 14154, 14428, 14703, 14979, 15256, 15534, 15813, 16092, 16373,
			16654, 16936, 17218, 17502, 17786, 18071, 18356, 18643, 18930, 19218,
			19507, 19796, 20086, 20377, 20668, 20960, 21253, 21547, 21841, 22136,
			22431, 22727, 23024, 23321, 23619, 23918, 24217, 24517, 24818, 25119}
		return ship_power[tonnage]
	} else if tonnage <= 4068 {
		/* Tonnage is not in table. Break it up into two halves and get
		 * approximate result = 1.149 * (x1 + x2), using recursion if
		 * necessary. */
		t1 := tonnage / 2
		t2 := tonnage - t1
		return 1149 * (power(t1) + power(t2)) / 1000
	}
	panic(fmt.Sprintf("\n\n\tLong integer overflow will occur in call to 'power(tonnage)'!\n\t\tActual call is power(%d).\n\n", tonnage))
}
