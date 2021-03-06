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

import "strings"

// NamedPlanet represents the name applied to a planet by a species.
// The ID is usually the name of the planet converted to uppercase.
type NamedPlanet struct {
	Id      string  // unique identifier for named planet
	Colony  *Colony // optional
	Display struct {
		Name string // original name of the planet
	}
	Index   int // nampla_index
	Planet  *Planet
	Special struct { // this is the hellish field that is used for multiple purposes
		// captures the value
		Value int
		// and these capture what I hope is the intent
		ExcessRawMaterialUnitsThatMayBeRecylcedInAutoMode int
	}
}

func newNamedPlanet(name string, planet *Planet, index int) *NamedPlanet {
	np := &NamedPlanet{
		Id:     strings.ToUpper(strings.TrimSpace(name)),
		Planet: planet,
		Index:  index,
	}
	np.Display.Name = strings.TrimSpace(name)
	return np
}
