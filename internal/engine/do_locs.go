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

// add_location will create a new location if there is not already one for the current species
func (e *Engine) add_location(x, y, z int) {
	for _, loc := range e.loc {
		if loc.s == e.species_number && loc.x == x && loc.y == y && loc.z == z {
			return /* This location is already in list for this species. */
		}
	}
	/* Add new location to the list. */
	e.loc = append(e.loc, &sp_loc_data{s: e.species_number, x: x, y: y, z: z})
	e.num_locs = len(e.loc)
}

// do_locations will populate the "loc" array based on all species' data
func (e *Engine) do_locations() {
	e.loc, e.num_locs = nil, 0
	for e.species_number = 1; e.species_number <= e.galaxy.num_species; e.species_number++ {
		e.species = e.spec_data[e.species_number-1]
		e.nampla_base = e.namp_data[e.species_number-1]
		e.ship_base = e.ship_data[e.species_number-1]

		for i := 0; i < e.species.num_namplas; i++ {
			e.nampla = e.nampla_base[i]
			if e.nampla.pn == 99 {
				continue
			}
			if (e.nampla.status & POPULATED) != 0 {
				e.add_location(e.nampla.x, e.nampla.y, e.nampla.z)
			}
		}

		for i := 0; i < e.species.num_ships; i++ {
			e.ship = e.ship_base[i]
			if e.ship.pn == 99 {
				continue
			}
			if e.ship.status == FORCED_JUMP || e.ship.status == JUMPED_IN_COMBAT {
				continue
			}
			e.add_location(e.ship.x, e.ship.y, e.ship.z)
		}
	}
}
