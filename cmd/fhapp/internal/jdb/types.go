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

package jdb

import "github.com/mdhender/fhcms/internal/coords"

// Colony represents a single colony.
type Colony struct {
	Name       *NamedPlanet
	System     *System
	Planet     *Planet // planet the colony is on
	Is         Status
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
	Inventory   map[string]int // key is item code, value is quantity
	SiegeEff    int
	Shipyards   int
	Special     int
	UseOnAmbush int
	Message     int
}

// GasType represents information about a specific gas.
type GasType struct {
	Code       string
	Percentage int
}

// NamedPlanet represents the name applied to a planet by a species.
type NamedPlanet struct {
	Id          string // unique identifier for named planet
	DisplayName string // original name of the planet
	Planet      *Planet
}

// NamedSystem represents the name applied to a system by a species.
type NamedSystem struct {
	Id          string // unique identifier for named system
	DisplayName string // original name of the system
	System      *System
}

// Planet represents a single planet.
type Planet struct {
	Id                       string         // unique identifier for the planet
	Coords                   *coords.Coords // location of the planet in the cluster
	System                   *System        // system containing the planet
	Orbit                    int
	Colonies                 map[string]*Colony // key is species id
	Atmosphere               []GasType
	Diameter                 int
	EconEfficiency           int
	Gravity                  int
	MiningDifficulty         int
	MiningDifficultyIncrease int
	PressureClass            int
	TemperatureClass         int
	Special                  string
	Message                  int
}

// Ship represents a single ship.
type Ship struct {
	Id          string // unique identifier for ship
	DisplayName string // original name of the ship
	Class       struct {
		Code        string
		Description string
		FTL         bool
		Tonnage     int
		Cost        int
	}
	Age                int
	Location           *coords.Coords
	ArrivedViaWormhole bool
	Destination        *coords.Coords
	Inventory          map[string]int // key is item code, value is quantity
	JustJumped         bool
	LoadingPoint       int
	UnloadingPoint     int
	RemainingCost      int
	Special            int
	Status             struct {
		ForcedJump        bool
		InDeepSpace       bool
		InOrbit           bool
		JumpedInCombat    bool
		OnSurface         bool
		UnderConstruction bool
	}
}

type Species struct {
	Id         string // unique identifier for species
	No         int    // number of the species (for reporting?)
	Name       string // name of the species
	Government struct {
		Name string
		Type string
	}
	HomeWorld              *Planet // the home world
	MI, MA, ML, GV, LS, BI *Technology
	Colonies               []*Colony // coordinates of the planet containing the colony
	Fleet                  struct {
		Cost           int
		MaintenancePct int              // percentage of production applied to fleet maintenance
		Ships          map[string]*Ship // key is name of ship, converted to upper case
	}
	Named struct {
		Systems map[string]*NamedSystem // key is name of system, converted to upper case
		Planets map[string]*NamedPlanet // key is name of planet, converted to upper case
	}
	Gases struct {
		Required struct {
			Gas    *GasType
			MinPct int
			MaxPct int
		}
		Neutral []*GasType
		Poison  []*GasType
	}
	Ally           map[string]*Species
	Contact        map[string]*Species
	Enemy          map[string]*Species
	Scanned        map[string]*System // all systems that have been scanned
	Visited        map[string]*System // all systems that have been visited
	AutoOrders     bool
	EconUnits      int
	HPOriginalBase int // something something home world
}

type Status struct {
	Colony          bool
	DisbandedColony bool
	Hiding          bool
	Hidden          bool
	HomePlanet      bool
	MiningColony    bool
	Populated       bool
	ResortColony    bool
}

// System represents a single star system in the cluster.
type System struct {
	Id          string         // unique identifier for the system
	Coords      *coords.Coords // location of the star within the cluster
	Color, Type struct {
		DisplayCode string
		Description string
	}
	Planets   []*Planet           // planets in system, ordered by orbit
	VisitedBy map[string]*Species // list of species that have visited the system
	Wormhole  *System             // other end of wormhole, nil if not a wormhole
	Message   int
}

// Technology represents all the data for a single technology.
type Technology struct {
	Code             string
	Level            int // current level
	KnowledgeLevel   int // un-applied technology knowledge
	ExperiencePoints int // experience points for technologies
}
