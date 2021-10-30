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

// System represents a single star system in the cluster.
type System struct {
	Id        string // unique identifier for the system
	Color     *Code
	Location  *Coords // location of the star within the cluster
	Message   int
	Planets   []*Planet // planets in system, ordered by orbit
	Type      *Code
	VisitedBy map[string]*Species // list of species that have visited the system
	Wormhole  *System             // other end of wormhole, nil if not a wormhole
}
