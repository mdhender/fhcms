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

// Planet represents a single planet.
type Planet struct {
	Coords                   *Coords            `json:"coords"` // location of the planet within the cluster
	System                   string             `json:"system"` // identifier for system containing the planet
	Atmosphere               []GasType          `json:"atmosphere"`
	Colonies                 map[string]*Colony `json:"colonies"` // key is species
	Diameter                 int                `json:"diameter"`
	EconEfficiency           int                `json:"econ_efficiency"`
	Gravity                  int                `json:"gravity"`
	Message                  int                `json:"message"`
	MiningDifficulty         int                `json:"mining_difficulty"`
	MiningDifficultyIncrease int                `json:"mining_difficulty_increase"`
	PressureClass            int                `json:"pressure_class"`
	TemperatureClass         int                `json:"temperature_class"`
	Special                  string             `json:"special,omitempty"`
}
