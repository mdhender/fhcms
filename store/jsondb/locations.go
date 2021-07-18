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

package jsondb

// SetLocations determines the current locations of colonies and ships.
// It updates Locations in the data store.
// TODO: What does Locations really mean? Is it reset each turn and shows only the current locations? Or is it a history of locations visited?
func (ds *Store) SetLocations() {
	ds.Locations = nil
	for species_number := 1; species_number <= ds.Galaxy.NumSpecies; species_number++ {
		species := ds.Species[species_number-1]
		nampla_base := species.Namplas
		ship_base := species.Ships

		for i := 0; i < species.NumNamplas; i++ {
			nampla := nampla_base[i]
			if nampla.Pn == 99 || (nampla.Status & POPULATED) == 0{
				continue
			}
			ds.addLocation(species_number, nampla.X, nampla.Y, nampla.Z)
		}

		for i := 0; i < species.NumShips; i++ {
			ship := ship_base[i]
			if ship.Pn == 99 {
				continue
			} else if ship.Status == FORCED_JUMP || ship.Status == JUMPED_IN_COMBAT {
				continue
			}
			ds.addLocation(species_number, ship.X, ship.Y, ship.Z)
		}
	}
}

func (ds *Store) addLocation(s, x,y,z int) {
	for i := 0; i < len(ds.Locations); i++ {
		if ds.Locations[i].X != x || ds.Locations[i].Y != y || ds.Locations[i].Z != z || ds.Locations[i].S != s {
			continue
		}
		return /* This location is already in list for this species. */
	}

	/* Add new location to the list. */
	ds.Locations = append(ds.Locations, Location{S:s,X:x,Y:y,Z:z})
}

