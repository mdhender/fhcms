/*******************************************************************************
fhapp - a player aid for Far Horizons
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

type StatsData struct {
	SpeciesId          string `json:"species_id"` // unique identifier for the species
	SpeciesNo          int    `json:"species_no"` // number of the species (for reporting and sorting?)
	Name               string `json:"name"`
	MiningLevel        int    `json:"mi"`
	ManufacturingLevel int    `json:"ma"`
	MilitaryLevel      int    `json:"ml"`
	GraviticsLevel     int    `json:"gv"`
	LifeSupportLevel   int    `json:"ls"`
	BiologyLevel       int    `json:"bi"`
	Production         int    `json:"production"`
	Colonies           int    `json:"colonies"`
	Ships              int    `json:"ships"`
	Shipyards          int    `json:"shipyards"`
	OffensivePower     int    `json:"offensive_power"`
	DefensivePower     int    `json:"defensive_power"`
	BankedEu           int    `json:"banked_eu"`
}
