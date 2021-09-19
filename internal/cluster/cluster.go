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

// Cluster represents the entire cluster. It contains all the systems,
// planets, species, named planets, colonies, and ships in the game.
type Cluster struct {
	Turn    int                 `json:"turn"`   // current turn number
	Radius  int                 `json:"radius"` // radius of the cluster
	Systems map[string]*System  `json:"systems"`
	Planets map[string]*Planet  `json:"planets"`
	Species map[string]*Species `json:"species"`
}
