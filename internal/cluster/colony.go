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
	Name   *NamedPlanet
	System *System
	Planet *Planet // planet the colony is on
	Is     struct {
		Colony          bool
		DisbandedColony bool
		Hiding          bool
		Hidden          bool
		HomePlanet      bool
		MiningColony    bool
		Populated       bool
		ResortColony    bool
	}
	Population int // population units - not the same as colonists
	Colonists  struct {
		Auto               bool
		AvailableToInstall int
		Needed             int
	}
	Manufacturing, Mining struct {
		Auto               bool
		AvailableToInstall int
		Base               int
		Needed             int
	}
	Inventory   map[string]*Item // key is item code
	SiegeEff    int
	Shipyards   int
	Special     int
	UseOnAmbush int
	Message     int
}
