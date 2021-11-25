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

type Code struct {
	Code  string
	Descr string
}

// gasTranslate maps atmospheric gas to code, description, and amount
func gasTranslate(i int) *Code {
	switch i {
	case 1:
		return &Code{Code: "H2", Descr: "Hydrogen"}
	case 2:
		return &Code{Code: "CH4", Descr: "Methane"}
	case 3:
		return &Code{Code: "He", Descr: "Helium"}
	case 4:
		return &Code{Code: "NH3", Descr: "Ammonia"}
	case 5:
		return &Code{Code: "N2", Descr: "Nitrogen"}
	case 6:
		return &Code{Code: "CO2", Descr: "Carbon Dioxide"}
	case 7:
		return &Code{Code: "O2", Descr: "Oxygen"}
	case 8:
		return &Code{Code: "HCl", Descr: "Hydrogen Chloride"}
	case 9:
		return &Code{Code: "Cl2", Descr: "Chlorine"}
	case 10:
		return &Code{Code: "F2", Descr: "Fluorine"}
	case 11:
		return &Code{Code: "H2O", Descr: "Steam"}
	case 12:
		return &Code{Code: "SO2", Descr: "Sulfur Dioxide"}
	case 13:
		return &Code{Code: "H2S", Descr: "Hydrogen Sulfide"}
	}
	panic(fmt.Sprintf("assert(gas != %d)", i))
}

// planetSpecialTranslate maps planet code to display code and description
func planetSpecialTranslate(i int) *Code {
	switch i {
	case 1:
		return &Code{"1", "Ideal Home Planet"}
	case 2:
		return &Code{"2", "Ideal Colony Planet"}
	case 3:
		return &Code{"3", "Radioactive Hellhole"}
	}
	return &Code{}
}

// starColorTranslate maps star color to display code and description
func starColorTranslate(i int) *Code {
	switch i {
	case 1:
		return &Code{"O", "BLUE"}
	case 2:
		return &Code{"B", "BLUE_WHITE"}
	case 3:
		return &Code{"A", "WHITE"}
	case 4:
		return &Code{"F", "YELLOW_WHITE"}
	case 5:
		return &Code{"G", "YELLOW"}
	case 6:
		return &Code{"K", "ORANGE"}
	case 7:
		return &Code{"K", "RED"}
	}
	panic(fmt.Sprintf("assert(star.color != %d)", i))
}

// starTypeTranslate maps star type to display code and description
func starTypeTranslate(i int) *Code {
	switch i {
	case 1:
		return &Code{"d", "DWARF"}
	case 2:
		return &Code{"D", "DEGENERATE"}
	case 3:
		return &Code{" ", "MAIN_SEQUENCE"}
	case 4:
		return &Code{"G", "GIANT"}
	}
	panic(fmt.Sprintf("assert(star.type != %d)", i))
}
