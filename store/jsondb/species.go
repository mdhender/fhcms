/*****************************************************************************
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
****************************************************************************/

package jsondb

import (
	"fmt"
	"strings"
)

func (ds *Store) xFilterSpecies(fn func(*Species) bool) []*Species {
	var result []*Species
	for _, s := range ds.Species {
		if fn(s) {
			result = append(result, s)
		}
	}
	if result == nil {
		return []*Species{}
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Name > result[j].Name {

			}
		}
	}
	return result
}

func (sp *Species) xLess(sp2 *Species) bool {
	return sp.Key < sp2.Key
}

func xSpeciesById(roles map[string]bool, ids ...int) func(*Species) bool {
	return func(sp *Species) bool {
		for _, id := range ids {
			if id == sp.Id {
				return roles[sp.Key]
			}
		}
		return false
	}
}

func xSpeciesByName(roles map[string]bool, names ...string) func(*Species) bool {
	return func(sp *Species) bool {
		for _, name := range names {
			if name == sp.Name {
				return roles[fmt.Sprintf("SP%02d", sp.Id)]
			}
		}
		return false
	}
}

type Species struct {
	Id         int    `json:"id"`
	Key        string `json:"key"`
	Name       string `json:"name"`
	Government struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"government"`
	Homeworld struct {
		Id     int    `json:"id"` // this is the planet index
		Key    string `json:"key"`
		Coords Coords `json:"coords"`
		Orbit  int    `json:"orbit"`
	} `json:"homeworld"`
	Gases      Gases `json:"gases"`
	AutoOrders bool  `json:"auto_orders"`
	Tech       struct {
		Biology       Technology `json:"BI"`
		Gravitics     Technology `json:"GV"`
		LifeSupport   Technology `json:"LS"`
		Manufacturing Technology `json:"MA"`
		Mining        Technology `json:"MI"`
		Military      Technology `json:"ML"`
	} `json:"tech"`
	BankedEconUnits  int                     `json:"econ_units"`
	HpOriginalBase   int                     `json:"hp_original_base"`
	FleetCost        int                     `json:"fleet_cost"`
	FleetPercentCost int                     `json:"fleet_percent_cost"`
	Contacts         []string                `json:"contacts"`
	Allies           []string                `json:"allies"`
	Enemies          []string                `json:"enemies"`
	NamedPlanets     map[string]*NamedPlanet `json:"namplas"`
	Ships            map[string]*Ship        `json:"ships"`
	// Aliens is a map of SPxx to AlienRelationship
	Aliens map[int]string `json:"aliens"`
}

func (sp *Species) GetShip(name string) *Ship {
	name = strings.ToUpper(name)
	if ship, ok := sp.Ships[name]; ok {
		return ship
	}
	var bestShip *Ship
	bestScore, nextBestScore := -9999, -9999
	for _, ship := range sp.Ships {
		n := ship.FuzzyMatch(name)
		if n > bestScore || bestShip == nil {
			// best match so far
			bestScore, bestShip = n, ship
		} else if n > nextBestScore {
			nextBestScore = n
		}
	}
	if bestShip == nil {
		return nil
	} else if bestScore == nextBestScore { // another name with equal score
		return nil
	} else if nameLength := len(bestShip.Name); nameLength < 5 { // no errors allowed
		// too short to trust comparison
		return nil
	} else if minimumScore := nameLength - ((nameLength / 7) + 1); bestScore < minimumScore { // score too low
		return nil
	}

	return bestShip
}
