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

// Species represents a single species.
// TODO: Scan should track the scan results since the planet attributes can change during play.
type Species struct {
	Name       string `json:"name"` // name of the species
	Government struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"government"`
	HomeWorld  string   `json:"home_world"` // coordinates of the home world
	Ally       []string `json:"ally"`
	AutoOrders bool     `json:"auto_orders,omitempty"`
	Colonies   []string `json:"colonies"` // coordinates of the planet containing the colony
	Contact    []string `json:"contact"`
	EconUnits  int      `json:"econ_units"`
	Enemy      []string `json:"enemy"`
	Fleet      struct {
		Cost           int `json:"cost"`
		MaintenancePct int `json:"maintenance_pct"` // percentage of production applied to fleet maintenance
	} `json:"fleet"`
	Gases struct {
		Required struct {
			Code   string `json:"code"`
			MinPct int    `json:"min_pct"`
			MaxPct int    `json:"max_pct"`
		} `json:"required"`
		Neutral []string `json:"neutral"`
		Poison  []string `json:"poison"`
	} `json:"gases"`
	HPOriginalBase int                     `json:"hp_original_base"`
	NamedPlanets   map[string]*NamedPlanet `json:"named_planets"` // key is name of planet, converted to upper case
	Scanned        []string                `json:"scanned"`       // coordinates of all systems that have been scanned
	Ships          map[string]*Ship        `json:"ships"`         // key is name of ship, converted to upper case
	Tech           struct {
		MI Technology `json:"mi"`
		MA Technology `json:"ma"`
		ML Technology `json:"ml"`
		GV Technology `json:"gv"`
		LS Technology `json:"ls"`
		BI Technology `json:"bi"`
	} `json:"tech"`
	Visited []string `json:"visited"` // coordinates of all systems that have been visited
}
