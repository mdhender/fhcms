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

import "fmt"

// Ship represents a single ship.
type Ship struct {
	DisplayName string // original name of the ship
	Class       struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		FTL         bool   `json:"ftl,omitempty"`
		Tonnage     int    `json:"tonnage"`
		Cost        int    `json:"cost"`
	} `json:"class"`
	Age                int            `json:"age"`
	Location           *Coords        `json:"location,omitempty"`
	ArrivedViaWormhole bool           `json:"arrived_via_wormhole,omitempty"`
	Destination        *Coords        `json:"destination,omitempty"`
	Inventory          map[string]int `json:"inventory"` // key is item code, value is quantity
	JustJumped         bool           `json:"just_jumped,omitempty"`
	LoadingPoint       int            `json:"loading_point,omitempty"`
	UnloadingPoint     int            `json:"unloading_point,omitempty"`
	RemainingCost      int            `json:"remaining_cost,omitempty"`
	Special            int            `json:"special,omitempty"`
	Status             struct {
		ForcedJump        bool `json:"forced_jump,omitempty"`
		InDeepSpace       bool `json:"in_deep_space,omitempty"`
		InOrbit           bool `json:"in_orbit,omitempty"`
		JumpedInCombat    bool `json:"jumped_in_combat,omitempty"`
		OnSurface         bool `json:"on_surface,omitempty"`
		UnderConstruction bool `json:"under_construction,omitempty"`
	} `json:"status"`
}

// shipToClass maps ship class to class
func shipToClass(i int) (string, string, int, int) {
	switch i {
	case 0:
		return "PB", "Picketboat", 1, 100
	case 1:
		return "CT", "Corvette", 2, 200
	case 2:
		return "ES", "Escort", 5, 500
	case 3:
		return "DD", "Destroyer", 10, 1000
	case 4:
		return "FG", "Frigate", 15, 1500
	case 5:
		return "CL", "Light Cruiser", 20, 2000
	case 6:
		return "CS", "Strike Cruiser", 25, 2500
	case 7:
		return "CA", "Heavy Cruiser", 30, 3000
	case 8:
		return "CC", "Command Cruiser", 35, 3500
	case 9:
		return "BC", "Battlecruiser", 40, 4000
	case 10:
		return "BS", "Battleship", 45, 4500
	case 11:
		return "DN", "Dreadnought", 50, 5000
	case 12:
		return "SD", "Super Dreadnought", 55, 5500
	case 13:
		return "BM", "Battlemoon", 60, 6000
	case 14:
		return "BW", "Battleworld", 65, 6500
	case 15:
		return "BR", "Battlestar", 70, 7000
	case 16:
		return "BA", "Starbase", 1, 100
	case 17:
		return "TR", "Transport", 1, 100
	}

	panic(fmt.Sprintf("assert(ship.class != %d)", i))
}

// shipToStatus maps ship status code to text
func shipToStatus(i int) string {
	switch i {
	case 0:
		return "UNDER_CONSTRUCTION"
	case 1:
		return "ON_SURFACE"
	case 2:
		return "IN_ORBIT"
	case 3:
		return "IN_DEEP_SPACE"
	case 4:
		return "JUMPED_IN_COMBAT"
	case 5:
		return "FORCED_JUMP"
	}
	panic(fmt.Sprintf("assert(ship.status != %d)", i))
}
