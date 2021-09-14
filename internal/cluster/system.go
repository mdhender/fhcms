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

// System represents a single star system in the cluster.
type System struct {
	Coords *Coords `json:"coords"` // location of the star within the cluster
	Color  struct {
		DisplayCode string `json:"display_code"`
		Description string `json:"description"`
	} `json:"color"`
	Message             int      `json:"message,omitempty"`
	Planets             []string `json:"planets"`                         // list of identifiers for each planet, ordered by orbit
	PotentialHomeSystem bool     `json:"potential_home_system,omitempty"` // true if the system is a good potential home system
	Type                struct {
		DisplayCode string `json:"display_code"`
		Description string `json:"description"`
	} `json:"type"`
	VisitedBy []string `json:"visited_by"`         // list of identifiers for every species that has visited the system
	Wormhole  *Coords  `json:"wormhole,omitempty"` // coordinates of other end of wormhole, nil if not a wormhole
}

// starColorTranslate maps star color to display code and description
func starColorTranslate(i int) (string, string) {
	switch i {
	case 1:
		return "O", "BLUE"
	case 2:
		return "B", "BLUE_WHITE"
	case 3:
		return "A", "WHITE"
	case 4:
		return "F", "YELLOW_WHITE"
	case 5:
		return "G", "YELLOW"
	case 6:
		return "K", "ORANGE"
	case 7:
		return "K", "RED"
	}
	panic(fmt.Sprintf("assert(star.color != %d)", i))
}

// starTypeTranslate maps star type to display code and description
func starTypeTranslate(i int) (string, string) {
	switch i {
	case 1:
		return "d", "DWARF"
	case 2:
		return "D", "DEGENERATE"
	case 3:
		return " ", "MAIN_SEQUENCE"
	case 4:
		return "G", "GIANT"
	}
	panic(fmt.Sprintf("assert(star.type != %d)", i))
}
