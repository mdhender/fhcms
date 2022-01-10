/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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

package engine

import (
	"fmt"
	"log"
	"strings"
)

func (e *Engine) closest_unvisited_star_orders(ship *ship_data) {
	e.x = -1 // signals that we have found no destination star

	var closest_star *star_data
	closest_distance, found := 999999, FALSE
	for i := 0; i < e.num_stars; i++ {
		star := e.star_base[i]
		if star.visited_by[e.species_index] == TRUE {
			continue // already visited by this species
		} else if ship.x == star.x && ship.y == star.y && ship.z == star.z {
			continue // this is here?
		}

		temp_distance := ((ship.x - star.x) * (ship.x - star.x)) + ((ship.y - star.y) * (ship.y - star.y)) + ((ship.z - star.z) * (ship.z - star.z))
		if temp_distance < closest_distance {
			e.x, e.y, e.z = star.x, star.y, star.z
			closest_distance, closest_star = temp_distance, star
			found = TRUE
		}
	}

	if found == TRUE {
		e.orders_file.WriteString(fmt.Sprintf("%d %d %d", e.x, e.y, e.z))
		/* So that we don't send more than one ship to the same place. */
		closest_star.visited_by[e.species_index] = TRUE
	} else {
		e.orders_file.WriteString("???")
	}

	return
}

/* This routine is intended to take a long argument and return a pointer to a string that has embedded commas to make the string more readable. */
func commas(value int) string {
	if value == 0 {
		return "0"
	} else if value < 0 {
		return "-" + commas(-1*value)
	}
	s := fmt.Sprintf("%024d", value)
	return strings.TrimLeft(s[0:3]+","+s[3:6]+","+s[6:9]+","+s[9:12]+","+s[12:15]+","+s[15:18]+","+s[18:21]+","+s[21:], "0,")
}

func (e *Engine) print_mishap_chance_orders(ship *ship_data, destx, desty, destz int) {
	if destx == -1 {
		e.orders_file.WriteString("Mishap chance = ???")
		return
	}
	x, y, z := destx, desty, destz

	mishap_GV := e.species.tech_level[GV]
	mishap_age := ship.age
	mishap_chance := (100 * (((x - ship.x) * (x - ship.x)) + ((y - ship.y) * (y - ship.y)) + ((z - ship.z) * (z - ship.z)))) / mishap_GV
	if mishap_age > 0 && mishap_chance < 10000 {
		success_chance := 10000 - mishap_chance
		success_chance -= (2 * mishap_age * success_chance) / 100
		mishap_chance = 10000 - success_chance
	}

	if mishap_chance > 10000 {
		mishap_chance = 10000
	}

	e.orders_file.WriteString(fmt.Sprintf("mishap chance = %d.%02d%%", mishap_chance/100, mishap_chance%100))
}

/* This routine will return a pointer to a string containing a complete
 * ship name, including its orbital/landed status and age. If global
 * variable "truncate_name" is TRUE, then orbital/landed status and age
 * will not be included. */
func (e *Engine) ship_name(ship *ship_data) string {
	var full_ship_id string

	ship_is_distorted := ship.item_quantity[FD] == ship.tonnage
	if ship.status == ON_SURFACE {
		ship_is_distorted = false
	}
	if e.ignore_field_distorters != 0 {
		ship_is_distorted = false
	}

	if ship_is_distorted {
		if ship.class == TR {
			full_ship_id = fmt.Sprintf("%s%d ???", ship_abbr[ship.class], ship.tonnage)
		} else if ship.class == BA {
			full_ship_id = fmt.Sprintf("BAS ???")
		} else {
			full_ship_id = fmt.Sprintf("%s ???", ship_abbr[ship.class])
		}
	} else if ship.class == TR {
		full_ship_id = fmt.Sprintf("%s%d%s %s", ship_abbr[ship.class], ship.tonnage, ship_type[ship._type], ship.name)
	} else {
		full_ship_id = fmt.Sprintf("%s%s %s", ship_abbr[ship.class], ship_type[ship._type], ship.name)
	}

	if e.truncate_name != 0 {
		return full_ship_id
	}

	full_ship_id += " ("

	effective_age := ship.age
	if effective_age < 0 {
		effective_age = 0
	}

	if !ship_is_distorted && ship.status != UNDER_CONSTRUCTION { /* Do age. */
		full_ship_id += fmt.Sprintf("A%d,", effective_age)
	}

	switch ship.status {
	case UNDER_CONSTRUCTION:
		full_ship_id += "C"
	case IN_ORBIT:
		full_ship_id += fmt.Sprintf("O%d", ship.pn)
	case ON_SURFACE:
		full_ship_id += fmt.Sprintf("L%d", ship.pn)
	case IN_DEEP_SPACE:
		full_ship_id += "D"
	case FORCED_JUMP:
		full_ship_id += "FJ"
	case JUMPED_IN_COMBAT:
		full_ship_id += "WD"
	default:
		full_ship_id += "***???***"
		log.Printf("[engine] WARNING!!!  Internal error in subroutine 'ship_name'\n")
	}

	if ship._type == STARBASE {
		full_ship_id += fmt.Sprintf(",%d tons", 10000*ship.tonnage)
	}

	return full_ship_id + ")"
}
