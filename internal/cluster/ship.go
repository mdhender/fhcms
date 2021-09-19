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

import (
	"fmt"
	"github.com/mdhender/fhcms/internal/coords"
)

// Ship represents a single ship.
type Ship struct {
	Id                 string // unique identifier for ship
	Age                int
	ArrivedViaWormhole bool
	Class              *ShipClass
	DisplayName        string // original name of the ship
	Destination        *coords.Coords
	Inventory          map[string]*Item // key is item code, value is quantity
	JustJumped         bool
	LoadingPoint       int
	Location           *coords.Coords
	RemainingCost      int
	Special            int
	Status             *ShipStatus
	UnloadingPoint     int
}

type ShipClass struct {
	Code        string
	Cost        int
	Description string
	FTL         bool
	Tonnage     int
}

type ShipStatus struct {
	ForcedJump        bool
	InDeepSpace       bool
	InOrbit           bool
	JumpedInCombat    bool
	OnSurface         bool
	UnderConstruction bool
}

// shipClassTranslate maps class to ship class
func shipClassTranslate(i int) *ShipClass {
	switch i {
	case BC:
		return &ShipClass{Code: "BC", Description: "Battlecruiser", Tonnage: 40, Cost: 4000}
	case BS:
		return &ShipClass{Code: "BS", Description: "Battleship", Tonnage: 45, Cost: 4500}
	case DN:
		return &ShipClass{Code: "DN", Description: "Dreadnought", Tonnage: 50, Cost: 5000}
	case BA:
		return &ShipClass{Code: "BA", Description: "Starbase", Tonnage: 1, Cost: 100}
	case BM:
		return &ShipClass{Code: "BM", Description: "Battlemoon", Tonnage: 60, Cost: 6000}
	case BR:
		return &ShipClass{Code: "BR", Description: "Battlestar", Tonnage: 70, Cost: 7000}
	case BW:
		return &ShipClass{Code: "BW", Description: "Battleworld", Tonnage: 65, Cost: 6500}
	case CA:
		return &ShipClass{Code: "CA", Description: "Heavy Cruiser", Tonnage: 30, Cost: 3000}
	case CC:
		return &ShipClass{Code: "CC", Description: "Command Cruiser", Tonnage: 35, Cost: 3500}
	case CL:
		return &ShipClass{Code: "CL", Description: "Light Cruiser", Tonnage: 20, Cost: 2000}
	case CS:
		return &ShipClass{Code: "CS", Description: "Strike Cruiser", Tonnage: 25, Cost: 2500}
	case CT:
		return &ShipClass{Code: "CT", Description: "Corvette", Tonnage: 2, Cost: 200}
	case DD:
		return &ShipClass{Code: "DD", Description: "Destroyer", Tonnage: 10, Cost: 1000}
	case ES:
		return &ShipClass{Code: "ES", Description: "Escort", Tonnage: 5, Cost: 500}
	case FG:
		return &ShipClass{Code: "FG", Description: "Frigate", Tonnage: 15, Cost: 1500}
	case PB:
		return &ShipClass{Code: "PB", Description: "Picketboat", Tonnage: 1, Cost: 100}
	case SD:
		return &ShipClass{Code: "SD", Description: "Super Dreadnought", Tonnage: 55, Cost: 5500}
	case TR:
		return &ShipClass{Code: "TR", Description: "Transport", Tonnage: 1, Cost: 100}
	}
	panic(fmt.Sprintf("assert(ship.class != %d)", i))
}

/* Ship classes. */
const PB = 0  /* Picketboat. */
const CT = 1  /* Corvette. */
const ES = 2  /* Escort. */
const DD = 3  /* Destroyer. */
const FG = 4  /* Frigate. */
const CL = 5  /* Light Cruiser. */
const CS = 6  /* Strike Cruiser. */
const CA = 7  /* Heavy Cruiser. */
const CC = 8  /* Command Cruiser. */
const BC = 9  /* Battlecruiser. */
const BS = 10 /* Battleship. */
const DN = 11 /* Dreadnought. */
const SD = 12 /* Super Dreadnought. */
const BM = 13 /* Battlemoon. */
const BW = 14 /* Battleworld. */
const BR = 15 /* Battlestar. */
const BA = 16 /* Starbase. */
const TR = 17 /* Transport. */

// shipStatusTranslate maps status to ship status
func shipStatusTranslate(i int) *ShipStatus {
	switch i {
	case 0:
		return &ShipStatus{UnderConstruction: true}
	case 1:
		return &ShipStatus{OnSurface: true}
	case 2:
		return &ShipStatus{InOrbit: true}
	case 3:
		return &ShipStatus{InDeepSpace: true}
	case 4:
		return &ShipStatus{JumpedInCombat: true}
	case 5:
		return &ShipStatus{ForcedJump: true}
	}
	panic(fmt.Sprintf("assert(ship.status != %d)", i))
}

// String implements the Stringer interface
func (s *ShipStatus) String() string {
	if s.ForcedJump {
		return "FORCED_JUMP"
	} else if s.JumpedInCombat {
		return "JUMPED_IN_COMBAT"
	} else if s.InDeepSpace {
		return "IN_DEEP_SPACE"
	} else if s.InOrbit {
		return "IN_ORBIT"
	} else if s.OnSurface {
		return "ON_SURFACE"
	} else if s.UnderConstruction {
		return "UNDER_CONSTRUCTION"
	}
	return ""
}
