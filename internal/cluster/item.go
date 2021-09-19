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

type Item struct {
	Code     string
	Descr    string
	Quantity int
	Storage  int
}

func itemTranslate(i, qty int) *Item {
	switch i {
	case 0:
		return &Item{Code: "RM", Descr: "Raw Material Unit", Quantity: qty, Storage: qty * 1}
	case 1:
		return &Item{Code: "PD", Descr: "Planetary Defense Unit", Quantity: qty, Storage: qty * 1}
	case 2:
		return &Item{Code: "SU", Descr: "Starbase Unit", Quantity: qty, Storage: qty * 1}
	case 3:
		return &Item{Code: "DR", Descr: "Damage Repair Unit", Quantity: qty, Storage: qty * 1}
	case 4:
		return &Item{Code: "CU", Descr: "Colonist Unit", Quantity: qty, Storage: qty * 1}
	case 5:
		return &Item{Code: "IU", Descr: "Colonial Mining Unit", Quantity: qty, Storage: qty * 1}
	case 6:
		return &Item{Code: "AU", Descr: "Colonial Manufacturing Unit", Quantity: qty, Storage: qty * 1}
	case 7:
		return &Item{Code: "FS", Descr: "Fail-Safe Jump Unit", Quantity: qty, Storage: qty * 1}
	case 8:
		return &Item{Code: "JP", Descr: "Jump Portal Unit", Quantity: qty, Storage: qty * 1}
	case 9:
		return &Item{Code: "FM", Descr: "Forced Misjump Unit", Quantity: qty, Storage: qty * 1}
	case 10:
		return &Item{Code: "FJ", Descr: "Forced Jump Unit", Quantity: qty, Storage: qty * 1}
	case 11:
		return &Item{Code: "GT", Descr: "Gravitic Telescope Unit", Quantity: qty, Storage: qty * 1}
	case 12:
		return &Item{Code: "FD", Descr: "Field Distortion Unit", Quantity: qty, Storage: qty * 1}
	case 13:
		return &Item{Code: "TP", Descr: "Terraforming Plant", Quantity: qty, Storage: qty * 1}
	case 14:
		return &Item{Code: "GW", Descr: "Germ Warfare Bomb", Quantity: qty, Storage: qty * 1}
	case 15:
		return &Item{Code: "SG1", Descr: "Mark-1 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 16:
		return &Item{Code: "SG2", Descr: "Mark-2 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 17:
		return &Item{Code: "SG3", Descr: "Mark-3 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 18:
		return &Item{Code: "SG4", Descr: "Mark-4 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 19:
		return &Item{Code: "SG5", Descr: "Mark-5 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 20:
		return &Item{Code: "SG6", Descr: "Mark-6 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 21:
		return &Item{Code: "SG7", Descr: "Mark-7 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 22:
		return &Item{Code: "SG8", Descr: "Mark-8 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 23:
		return &Item{Code: "SG9", Descr: "Mark-9 Shield Generator", Quantity: qty, Storage: qty * 1}
	case 24:
		return &Item{Code: "GU1", Descr: "Mark-1 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 25:
		return &Item{Code: "GU2", Descr: "Mark-2 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 26:
		return &Item{Code: "GU3", Descr: "Mark-3 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 27:
		return &Item{Code: "GU4", Descr: "Mark-4 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 28:
		return &Item{Code: "GU5", Descr: "Mark-5 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 29:
		return &Item{Code: "GU6", Descr: "Mark-6 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 30:
		return &Item{Code: "GU7", Descr: "Mark-7 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 31:
		return &Item{Code: "GU8", Descr: "Mark-8 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 32:
		return &Item{Code: "GU9", Descr: "Mark-9 Gun Unit", Quantity: qty, Storage: qty * 1}
	case 33:
		return &Item{Code: "X1", Descr: "X1 Unit", Quantity: qty, Storage: qty * 1}
	case 34:
		return &Item{Code: "X2", Descr: "X2 Unit", Quantity: qty, Storage: qty * 1}
	case 35:
		return &Item{Code: "X3", Descr: "X3 Unit", Quantity: qty, Storage: qty * 1}
	case 36:
		return &Item{Code: "X4", Descr: "X4 Unit", Quantity: qty, Storage: qty * 1}
	case 37:
		return &Item{Code: "X5", Descr: "X5 Unit", Quantity: qty, Storage: qty * 1}
	}
	panic(fmt.Sprintf("assert(item != %d)", i))
}
