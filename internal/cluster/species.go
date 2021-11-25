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

type Species struct {
	Id         string // unique identifier for species
	MI         *Technology
	MA         *Technology
	ML         *Technology
	GV         *Technology
	LS         *Technology
	BI         *Technology
	Ally       map[string]*Species
	AutoOrders bool
	Colonies   struct {
		ById       map[string]*Colony // key is name of planet, converted to upper case
		ByLocation map[string]*Colony // key is location of planet, including orbit
	}
	Contact         map[string]*Species
	DistortedNumber int // number if field distortion units are working
	EconUnits       int
	Enemy           map[string]*Species
	Fleet           struct {
		Base              []*Ship // indexed from original data file
		Cost              int
		DefensiveStrength int
		MaintenancePct    int // percentage of production applied to fleet maintenance
		OffensiveStrength int
		Ships             map[string]*Ship // key is name of ship, converted to upper case
		Starbases         []*Ship
		Transports        []*Ship
		Warships          []*Ship
	}
	Gases struct {
		Neutral        []*Code
		Poison         []*Code
		Required       *Code
		RequiredMaxPct int
		RequiredMinPct int
	}
	Government struct {
		Name string
		Type string
	}
	HomeWorld struct {
		Name         *NamedPlanet
		OriginalBase int // something something home world
		Planet       *Planet
		System       *System
	}
	Name         string // name of the species
	NamedPlanets struct {
		Base       []*NamedPlanet          // indexed from original data file
		ById       map[string]*NamedPlanet // key is name of planet, converted to upper case
		ByLocation map[string]*NamedPlanet // key is location, including orbit
	}
	No      int       // number of the species (for reporting?)
	Scanned []*System // all systems that have been scanned
	Visited []*System // all systems that have been visited
}
