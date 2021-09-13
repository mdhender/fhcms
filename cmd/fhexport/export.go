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

package main

type Export struct {
	TurnNumber int          `json:"turn_number"`
	Species    *SpeciesData `json:"species"`
	Systems    Systems      `json:"systems"`
}

type Coords struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Z     int `json:"z"`
	Orbit int `json:"orbit,omitempty"`
}

type PlanetData struct {
	Coords Coords `json:"coords"`
}

type SpeciesData struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	GovtName   string      `json:"govt_name"`
	HomeSystem *SystemData `json:"-"`
}

type Systems []*SystemData

func (s Systems) Len() int {
	return len(s)
}

func (s Systems) Less(i, j int) bool {
	return s[i].Coords.LessThan(s[j].Coords)
}

func (s Systems) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type SystemData struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	Coords   Coords        `json:"coords"`
	Color    string        `json:"color"`
	Size     int           `json:"size"`
	Type     string        `json:"type"`
	Wormhole *Coords       `json:"wormhole,omitempty"`
	Planets  []*PlanetData `json:"planets,omitempty"`
}

func (c Coords) LessThan(t Coords) bool {
	if c.X < t.X {
		return true
	} else if c.X == t.X {
		if c.Y < t.Y {
			return true
		} else if c.Y == t.Y {
			if c.Z < t.Z {
				return true
			} else if c.Z == t.Z {
				return c.Orbit < t.Orbit
			}
		}
	}
	return false
}
