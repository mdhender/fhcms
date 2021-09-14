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

// GasType represents information about a specific gas.
type GasType struct {
	Code       string `json:"code"`
	Percentage int    `json:"percentage"`
}

// gasToType maps gas to code and description
func gasToType(i int) (string, string) {
	switch i {
	case 1:
		return "H2", "Hydrogen"
	case 2:
		return "CH4", "Methane"
	case 3:
		return "HE", "Helium"
	case 4:
		return "NH3", "Ammonia"
	case 5:
		return "N2", "Nitrogen"
	case 6:
		return "CO2", "Carbon Dioxide"
	case 7:
		return "O2", "Oxygen"
	case 8:
		return "HCL", "Hydrogen Chloride"
	case 9:
		return "CL2", "Chlorine"
	case 10:
		return "F2", "Fluorine"
	case 11:
		return "H2O", "Steam"
	case 12:
		return "SO2", "Sulfur Dioxide"
	case 13:
		return "H2S", "Hydrogen Sulfide"
	}
	panic(fmt.Sprintf("assert(gas != %d)", i))
}
