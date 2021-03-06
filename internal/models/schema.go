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

type Galaxy struct {
	Id     int
	Name   string
	Specie struct {
		Id         int
		Name       string
		Government struct {
			Name string
		}
	}
	CurrentTurn  int
	TurnNo       int
	Files        string
	PlayerCount  int
	SpeciesCount int
	Turns        GameTurns
	Display      struct {
		Deadline bool
	}
}

type Galaxies []*Galaxy

// Len implements the Sorter interface
func (g Galaxies) Len() int {
	return len(g)
}

// Less implements the Sorter interface
func (g Galaxies) Less(i, j int) bool {
	return g[i].Name < g[j].Name
}

// Swap implements the Sorter interface
func (g Galaxies) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

type GalaxyFetcher interface {
	FetchGalaxy(uid string, gid int) *Galaxy
}

type GalaxiesFetcher interface {
	FetchGalaxies(uid string) Galaxies
}

type GameTurns struct {
	Current int `json:"current"`
	Files   GameTurnFiles
}

type GameTurnFile struct {
	Turn  int
	Files string
}

type GameTurnFiles []*GameTurnFile

// Len implements the Sorter interface
func (g GameTurnFiles) Len() int {
	return len(g)
}

// Less implements the Sorter interface
func (g GameTurnFiles) Less(i, j int) bool {
	return g[i].Turn < g[j].Turn
}

// Swap implements the Sorter interface
func (g GameTurnFiles) Swap(i, j int) {
	g[i], g[j] = g[j], g[i]
}

type Specie struct {
	Id         int
	Name       string
	Government struct {
		Name string
	}
	Technology [6]Technology
	Statistics []*SpecieStatistic
}

type Species []Specie

// Len implements the Sorter interface
func (s Species) Len() int {
	return len(s)
}

// Less implements the Sorter interface
func (s Species) Less(i, j int) bool {
	return s[i].Government.Name < s[j].Government.Name
}

// Swap implements the Sorter interface
func (s Species) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SpecieFetcher interface {
	FetchSpecie(uid string, gid, spNo int, turnNo int) *Specie
}

type SpeciesFetcher interface {
	FetchSpecies(gid string) Species
}

type SpecieStatistic struct {
	Label string
	Value float64
	Units string
}

type Technology struct {
	Name             string
	Code             string
	InitialLevel     int
	CurrentLevel     int
	KnowledgeLevel   int
	ExperiencePoints int
}

type User struct {
	Id      string
	Name    string
	IsAdmin bool
}

type UserFetcher interface {
	FetchUser(uid string) *User
}
