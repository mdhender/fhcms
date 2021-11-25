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

// Coords represents a location in the cluster.
type Coords struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Z     int `json:"z"`
	Orbit int `json:"orbit,omitempty"`
}

func (c *Coords) Delta(to *Coords) int {
	return (c.X-to.X)*(c.X-to.X) + (c.Y-to.Y)*(c.Y-to.Y) + (c.Z-to.Z)*(c.Z-to.Z)
}

// Id returns a unique representation of the location.
func (c *Coords) Id() string {
	if c.Orbit != 0 {
		return fmt.Sprintf("%d.%d.%d.%d", c.X, c.Y, c.Z, c.Orbit)
	}
	return fmt.Sprintf("%d.%d.%d", c.X, c.Y, c.Z)
}

func NewCoords(x, y, z, orbit int) *Coords {
	return &Coords{X: x, Y: y, Z: z, Orbit: orbit}
}
