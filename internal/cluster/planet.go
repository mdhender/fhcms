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
	Id                       string // unique identifier for the planet
	Atmosphere               []*AtmosphericGas
	Colonies                 map[string]*Colony // key is species id
	Diameter                 int
	EconEfficiency           int
	Gravity                  int
	Location                 *Coords // location of the planet in the cluster
	Message                  int
	MiningDifficulty         int
	MiningDifficultyIncrease int
	Orbit                    int
	PressureClass            int
	TemperatureClass         int
	Special                  string
	System                   *System             // system containing the planet
	VisitedBy                map[string]*Species // all species that have visited the planet
}

// AtmosphericGas represents information about a specific gas in the planet's atmospher
type AtmosphericGas struct {
	Gas        *Code
	Percentage int
}

func atmosphereTranslate(i, pct int) *AtmosphericGas {
	return &AtmosphericGas{
		Gas:        gasTranslate(i),
		Percentage: pct,
	}
}
