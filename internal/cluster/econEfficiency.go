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

// UpdateEconEfficiency will recalculate the economic efficiencies of all planets.
func (ds *Store) UpdateEconEfficiency() {
	// calculate the total economic base for each planet from named planet data
	totalEconBase := make(map[string]int)
	for _, sp := range ds.Species {
		for _, np := range sp.NamedPlanets.ById {
			if np.Colony != nil {
				totalEconBase[np.Planet.Id] = totalEconBase[np.Planet.Id] + np.Colony.Mining.Base + np.Colony.Manufacturing.Base
			}
		}
	}
	// recalculate economic efficiencies of all planets
	for id, base := range totalEconBase {
		econEfficiency := 100
		if base > 2000 {
			econEfficiency = (100 * ((base-2000)/20 + 2000)) / base
		}
		ds.Planets[id].EconEfficiency = econEfficiency
	}
}
