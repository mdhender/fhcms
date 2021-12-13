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

package games

type Game struct {
	Id      string         `json:"id"`
	Name    string         `json:"name,omitempty"`
	Files   string         `json:"files,omitempty"`
	Players map[string]int `json:"players,omitempty"` // map of user id to species number
}

type GameList struct {
	ById     map[string]*Game   // map of game id to game
	ByPlayer map[string][]*Game // map of user id to slice of games
}
