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

// Colony represents a single colony.
type Colony struct {
	Status struct {
		Colony          bool `json:"colony"`
		DisbandedColony bool `json:"disbanded_colony,omitempty"`
		Hiding          bool `json:"hiding,omitempty"`
		Hidden          bool `json:"hidden,omitempty"`
		HomePlanet      bool `json:"home_planet,omitempty"`
		MiningColony    bool `json:"mining_colony,omitempty"`
		Populated       bool `json:"populated,omitempty"`
		ResortColony    bool `json:"resort_colony,omitempty"`
	} `json:"status"`
	Inventory   map[string]int `json:"inventory"` // key is item code, value is quantity
	MaBase      int            `json:"ma_base"`
	Message     int            `json:"message"`
	MiBase      int            `json:"mi_base"`
	PlanetIndex int            `json:"planet_index"`
	PopUnits    int            `json:"pop_units"`
	SiegeEff    int            `json:"siege_eff"`
	Shipyards   int            `json:"shipyards"`
	Special     int            `json:"special,omitempty"`
	Units       struct {
		Colonists struct {
			Auto               bool `json:"auto,omitempty"`
			AvailableToInstall int  `json:"available_to_install,omitempty"`
			Needed             int  `json:"needed,omitempty"`
		} `json:"colonists"`
		Manufacturing struct {
			Auto               bool `json:"auto,omitempty"`
			AvailableToInstall int  `json:"available_to_install,omitempty"`
			Needed             int  `json:"needed,omitempty"`
		} `json:"manufacturing"`
		Mining struct {
			Auto               bool `json:"auto,omitempty"`
			AvailableToInstall int  `json:"available_to_install,omitempty"`
			Needed             int  `json:"needed,omitempty"`
		} `json:"mining"`
	} `json:"units"`
	UseOnAmbush int `json:"use_on_ambush"`
}
