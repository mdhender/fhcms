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
