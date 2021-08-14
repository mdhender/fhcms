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

package orders

//// An ALLY will set the diplomatic status of either a single
//// species or all species to "ally."
//// Accepts the following formats
////   ALLY SPECIES
////   ALLY NUMBER
//// Where
////    SPECIES is the name of a species. Note that it must include the
////            "SP" code!
////    NUMBER  is any integer value.
//
//type Ally struct {
//	all     bool   // true only if all species are to be updated
//	species string // name of species to set diplomatic status
//}
//
//// parseAlly accepts the following formats
////   ALLY SPECIES
////   ALLY NUMBER
//func parseAlly(p []byte) *Ally {
//	command := "ally"
//	words := strings.Fields(strings.TrimSpace(string()))
//	if len(words) == 0 {
//		return nil
//	}
//	command := strings.ToLower(words[0])
//	if command == "ally" || command == "all" {
//
//	}
//}
