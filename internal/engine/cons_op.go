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

func (e *Engine) consolidate_option(option, location int) {
	// only attack options go in list
	if option < DEEP_SPACE_FIGHT {
		return
	}

	// make sure pre-requisites are already in the list.
	// bombardment, and germ warfare must follow a successful planet attack.
	if option > PLANET_ATTACK {
		e.consolidate_option(PLANET_ATTACK, location)
	}

	/* Check if option and location are already in list. */
	for i := 0; i < e.num_combat_options; i++ {
		if option == e.combat_option[i] && location == e.combat_location[i] {
			return
		}
	}

	// add new option to list
	e.combat_option[e.num_combat_options] = option
	e.combat_location[e.num_combat_options] = location
	e.num_combat_options++
}
