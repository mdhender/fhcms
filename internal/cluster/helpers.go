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

import (
	"fmt"
	"strings"
)

/* This routine is intended to take a long argument and return a pointer
 * to a string that has embedded commas to make the string more readable. */
func commas(value int) string {
	if value == 0 {
		return "0"
	} else if value < 0 {
		return "-" + commas(-1*value)
	}
	s := fmt.Sprintf("%024d", value)
	return strings.TrimLeft(s[0:3]+","+s[3:6]+","+s[6:9]+","+s[9:12]+","+s[12:15]+","+s[15:18]+","+s[18:21]+","+s[21:], "0,")
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

func lifeSupportNeeded(species *Species, planet *Planet) int {
	lsn := 0
	deltaTemperature := species.HomeWorld.Planet.TemperatureClass - planet.TemperatureClass
	if deltaTemperature < 0 {
		deltaTemperature *= -1
	}
	lsn = 3 * deltaTemperature

	deltaPressure := species.HomeWorld.Planet.PressureClass - planet.PressureClass
	if deltaPressure < 0 {
		deltaPressure *= -1
	}
	lsn += (3 * deltaPressure)

	/* Check gases. Assume required gas is NOT present. */
	lsn += 3                                      // assumes required gas is not present
	for j := 0; j < len(planet.Atmosphere); j++ { /* Check gases on planet. */
		if species.Gases.Required.Code == planet.Atmosphere[j].Gas.Code {
			if species.Gases.RequiredMinPct <= planet.Atmosphere[j].Percentage && planet.Atmosphere[j].Percentage <= species.Gases.RequiredMaxPct {
				lsn -= 3
			}
			break
		}
	}
	for j := 0; j < len(planet.Atmosphere); j++ { /* Check gases on planet. */
		if planet.Atmosphere[j].Percentage != 0 {
			for i := 0; i < len(species.Gases.Poison); i++ { /* Compare with poisonous gases. */
				if species.Gases.Poison[i].Code == planet.Atmosphere[j].Gas.Code {
					lsn += 3
				}
			}
		}
	}

	return lsn
}

/* Look-up table for ship defensive/offensive power uses ship->tonnage
 * as an index. Each value is equal to 100 * (ship->tonnage)^1.2. The
 * 'power' subroutine uses recursion to calculate values for tonnages
 * over 100. */
var ship_power = []int{0, /* Zeroth element not used. */
	100, 230, 374, 528, 690, 859, 1033, 1213, 1397, 1585,
	1777, 1973, 2171, 2373, 2578, 2786, 2996, 3209, 3424, 3641,
	3861, 4082, 4306, 4532, 4759, 4988, 5220, 5452, 5687, 5923,
	6161, 6400, 6641, 6883, 7127, 7372, 7618, 7866, 8115, 8365,
	8617, 8870, 9124, 9379, 9635, 9893, 10151, 10411, 10672, 10934,
	11197, 11461, 11725, 11991, 12258, 12526, 12795, 13065, 13336, 13608,
	13880, 14154, 14428, 14703, 14979, 15256, 15534, 15813, 16092, 16373,
	16654, 16936, 17218, 17502, 17786, 18071, 18356, 18643, 18930, 19218,
	19507, 19796, 20086, 20377, 20668, 20960, 21253, 21547, 21841, 22136,
	22431, 22727, 23024, 23321, 23619, 23918, 24217, 24517, 24818, 25119,
}

func power(tonnage int) int {
	if tonnage < 1 {
		return 0
	} else if tonnage < len(ship_power) {
		return ship_power[tonnage]
	} else if tonnage > 4068 {
		panic(fmt.Sprintf("assert(tonnage(%d <= 4068))", tonnage))
	}
	/* Tonnage is not in table. Break it up into two halves and get
	 * approximate result = 1.149 * (x1 + x2), using recursion. */
	t1 := tonnage / 2
	t2 := tonnage - t1
	return 1149 * (power(t1) + power(t2)) / 1000
}
