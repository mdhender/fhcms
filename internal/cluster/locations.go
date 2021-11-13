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

package cluster

type SpeciesLocationData struct {
	Species  *Species
	Location *Coords
}

// RunLocations will create the file locations.dat and will update the
// economic efficiencies of all planets.  These functions are also performed
// by Finish.c.  This program should be run before the strike phase or whenever
// manual changes are made to the species data files that resulted in something
// not being where it was or something being where it was not. It should also
// be run if you run Finish on fewer than all species and decide to keep the
// resulting planets.dat file.
func (ds *Store) RunLocations(verbose, test bool) []*SpeciesLocationData {
	// update economic efficiencies of all planets
	ds.UpdateEconEfficiency()

	// populate new locations array
	var locations []*SpeciesLocationData
	for _, sp := range ds.Species {
		for _, np := range sp.NamedPlanets.ById {
			if np.Colony != nil && np.Colony.Is.Populated {
				locations = addLocation(locations, sp, np.Planet.Location)
			}
		}

		for _, ship := range sp.Fleet.Ships {
			if ship.Location != nil && ship.Status != nil && !(ship.Status.ForcedJump || ship.Status.JumpedInCombat) {
				locations = addLocation(locations, sp, ship.Location)
			}
		}
	}

	return locations
}

func addLocation(locations []*SpeciesLocationData, sp *Species, c *Coords) []*SpeciesLocationData {
	for _, l := range locations {
		if l.Species.Id == sp.Id && l.Location.X == c.X && l.Location.Y == c.Y && l.Location.Z == c.Z {
			return locations // already present
		}
	}
	return append(locations, &SpeciesLocationData{Species: sp, Location: &Coords{X: c.X, Y: c.Y, Z: c.Z}})
}
