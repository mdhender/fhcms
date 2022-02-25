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

import "fmt"

// speciesBitIsSet returns true if the bit is set for the species.
// note: the species number must be 1 based!
// sp01       65536                       1 0000 0000 0000 0000
// sp09    16777216             1 0000 0000 0000 0000 0000 0000
// sp18  8589934592  10 0000 0000 0000 0000 0000 0000 0000 0000

func speciesBitIsSet(set [2]uint64, sp int) bool {
	return (set[0] & (1 << (sp + 15))) != 0
}

// gasToType maps gas to code and description
func gasToType(i int) (string, string) {
	switch i {
	case 1:
		return "H2", "Hydrogen"
	case 2:
		return "CH4", "Methane"
	case 3:
		return "He", "Helium"
	case 4:
		return "NH3", "Ammonia"
	case 5:
		return "N2", "Nitrogen"
	case 6:
		return "CO2", "Carbon Dioxide"
	case 7:
		return "O2", "Oxygen"
	case 8:
		return "HCl", "Hydrogen Chloride"
	case 9:
		return "Cl2", "Chlorine"
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

func gasToGasType(g int) (gt GasType) {
	switch g {
	case 1:
		gt.Code = "H2" // Hydrogen
	case 2:
		gt.Code = "CH4" // Methane
	case 3:
		gt.Code = "He" // Helium
	case 4:
		gt.Code = "NH3" // Ammonia
	case 5:
		gt.Code = "N2" // Nitrogen
	case 6:
		gt.Code = "CO2" // Carbon Dioxide
	case 7:
		gt.Code = "O2" // Oxygen
	case 8:
		gt.Code = "HCl" // Hydrogen Chloride
	case 9:
		gt.Code = "Cl2" // Chlorine
	case 10:
		gt.Code = "F2" // Fluorine
	case 11:
		gt.Code = "H2O" // Steam
	case 12:
		gt.Code = "SO2" // Sulfur Dioxide
	case 13:
		gt.Code = "H2S" // Hydrogen Sulfide
	default:
		panic(fmt.Sprintf("assert(gas != %d)", g))
	}
	return gt
}

func itemToCode(i int) (string, string) {
	switch i {
	case 0:
		return "RM", "Raw Material Unit"
	case 1:
		return "PD", "Planetary Defense Unit"
	case 2:
		return "SU", "Starbase Unit"
	case 3:
		return "DR", "Damage Repair Unit"
	case 4:
		return "CU", "Colonist Unit"
	case 5:
		return "IU", "Colonial Mining Unit"
	case 6:
		return "AU", "Colonial Manufacturing Unit"
	case 7:
		return "FS", "Fail-Safe Jump Unit"
	case 8:
		return "JP", "Jump Portal Unit"
	case 9:
		return "FM", "Forced Misjump Unit"
	case 10:
		return "FJ", "Forced Jump Unit"
	case 11:
		return "GT", "Gravitic Telescope Unit"
	case 12:
		return "FD", "Field Distortion Unit"
	case 13:
		return "TP", "Terraforming Plant"
	case 14:
		return "GW", "Germ Warfare Bomb"
	case 15:
		return "SG1", "Mark-1 Shield Generator"
	case 16:
		return "SG2", "Mark-2 Shield Generator"
	case 17:
		return "SG3", "Mark-3 Shield Generator"
	case 18:
		return "SG4", "Mark-4 Shield Generator"
	case 19:
		return "SG5", "Mark-5 Shield Generator"
	case 20:
		return "SG6", "Mark-6 Shield Generator"
	case 21:
		return "SG7", "Mark-7 Shield Generator"
	case 22:
		return "SG8", "Mark-8 Shield Generator"
	case 23:
		return "SG9", "Mark-9 Shield Generator"
	case 24:
		return "GU1", "Mark-1 Gun Unit"
	case 25:
		return "GU2", "Mark-2 Gun Unit"
	case 26:
		return "GU3", "Mark-3 Gun Unit"
	case 27:
		return "GU4", "Mark-4 Gun Unit"
	case 28:
		return "GU5", "Mark-5 Gun Unit"
	case 29:
		return "GU6", "Mark-6 Gun Unit"
	case 30:
		return "GU7", "Mark-7 Gun Unit"
	case 31:
		return "GU8", "Mark-8 Gun Unit"
	case 32:
		return "GU9", "Mark-9 Gun Unit"
	case 33:
		return "X1", "X1 Unit"
	case 34:
		return "X2", "X2 Unit"
	case 35:
		return "X3", "X3 Unit"
	case 36:
		return "X4", "X4 Unit"
	case 37:
		return "X5", "X5 Unit"
	}
	panic(fmt.Sprintf("assert(item != %d)", i))
}

// shipToClass maps ship class to class
func shipToClass(i int) (string, string, int, int) {
	switch i {
	case 0:
		return "PB", "Picketboat", 1, 100
	case 1:
		return "CT", "Corvette", 2, 200
	case 2:
		return "ES", "Escort", 5, 500
	case 3:
		return "FF", "Frigate", 10, 1000
	case 4:
		return "DD", "Destroyer", 15, 1500
	case 5:
		return "CL", "Light Cruiser", 20, 2000
	case 6:
		return "CS", "Strike Cruiser", 25, 2500
	case 7:
		return "CA", "Heavy Cruiser", 30, 3000
	case 8:
		return "CC", "Command Cruiser", 35, 3500
	case 9:
		return "BC", "Battlecruiser", 40, 4000
	case 10:
		return "BS", "Battleship", 45, 4500
	case 11:
		return "DN", "Dreadnought", 50, 5000
	case 12:
		return "SD", "Super Dreadnought", 55, 5500
	case 13:
		return "BM", "Battlemoon", 60, 6000
	case 14:
		return "BW", "Battleworld", 65, 6500
	case 15:
		return "BR", "Battlestar", 70, 7000
	case 16:
		return "BA", "Starbase", 1, 100
	case 17:
		return "TR", "Transport", 1, 100
	}

	panic(fmt.Sprintf("assert(ship.class != %d)", i))
}

// shipToStatus maps ship status code to text
func shipToStatus(i int) string {
	switch i {
	case 0:
		return "UNDER_CONSTRUCTION"
	case 1:
		return "ON_SURFACE"
	case 2:
		return "IN_ORBIT"
	case 3:
		return "IN_DEEP_SPACE"
	case 4:
		return "JUMPED_IN_COMBAT"
	case 5:
		return "FORCED_JUMP"
	}
	panic(fmt.Sprintf("assert(ship.status != %d)", i))
}

// starColorTranslate maps star color to display code and description
func starColorTranslate(i int) (string, string) {
	switch i {
	case 1:
		return "O", "BLUE"
	case 2:
		return "B", "BLUE_WHITE"
	case 3:
		return "A", "WHITE"
	case 4:
		return "F", "YELLOW_WHITE"
	case 5:
		return "G", "YELLOW"
	case 6:
		return "K", "ORANGE"
	case 7:
		return "K", "RED"
	}
	panic(fmt.Sprintf("assert(star.color != %d)", i))
}

// starTypeTranslate maps star type to display code and description
func starTypeTranslate(i int) (string, string) {
	switch i {
	case 1:
		return "d", "DWARF"
	case 2:
		return "D", "DEGENERATE"
	case 3:
		return " ", "MAIN_SEQUENCE"
	case 4:
		return "G", "GIANT"
	}
	panic(fmt.Sprintf("assert(star.type != %d)", i))
}

// techToCode maps tech level id to the code
func techToCode(i int) string {
	switch i {
	case 0:
		return "MI" // Mining
	case 1:
		return "MA" // Manufacturing
	case 2:
		return "ML" // Military
	case 3:
		return "GV" // Gravitics
	case 4:
		return "LS" // Life Support
	case 5:
		return "BI" // Biology}
	}
	panic(fmt.Sprintf("assert(tech.code != %d)", i))
}
