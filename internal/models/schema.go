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

package models

type Site struct {
	Title string
}

type SiteFetcher interface {
	FetchSite() *Site
}

type Game struct {
	Id           string
	Name         string
	Files        string
	PlayerCount  int
	SpeciesCount int
}

type Games []*Game

// Len implements the Sorter interface
func (g Games) Len() int {
	return len(g)
}

// Less implements the Sorter interface
func (g Games) Less(i, j int) bool {
	return g[i].Name < g[j].Name
}

// Swap implements the Sorter interface
func (g Games) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

type GameFetcher interface {
	FetchGame(uid, gid string) *Game
}

type GamesFetcher interface {
	FetchGames(uid string) Games
}

type Specie struct {
	Id   string
	Name string
}

type Species []Specie

// Len implements the Sorter interface
func (s Species) Len() int {
	return len(s)
}

// Less implements the Sorter interface
func (s Species) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// Swap implements the Sorter interface
func (s Species) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SpecieFetcher interface {
	FetchSpecie(gid, spid string) *Specie
}

type SpeciesFetcher interface {
	FetchSpecies(gid string) Species
}

type User struct {
	Id      string
	Name    string
	IsAdmin bool
}

type UserFetcher interface {
	FetchUser(uid string) *User
}
