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

// Technology represents all the data for a single technology.
type Technology struct {
	Code             string `json:"code"`
	Level            int    `json:"level"`             // current level
	KnowledgeLevel   int    `json:"knowledge_level"`   // un-applied technology knowledge
	ExperiencePoints int    `json:"experience_points"` // experience points for technologies
}

func (t *Technology) DisplayName() string {
	switch t.Code {
	case "BI":
		return "Biology"
	case "GV":
		return "Gravitics"
	case "LS":
		return "Life Support"
	case "MA":
		return "Manufacturing"
	case "MI":
		return "Mining"
	case "ML":
		return "Military"
	}
	panic(fmt.Sprintf("assert(tech != %q)", t.Code))
}

func TechName(i int) string {
	switch i {
	case 0:
		return "Mining"
	case 1:
		return "Manufacturing"
	case 2:
		return "Military"
	case 3:
		return "Gravitics"
	case 4:
		return "Life Support"
	case 5:
		return "Biology"
	}
	panic(fmt.Sprintf("assert(tech != %d)", i))
}
