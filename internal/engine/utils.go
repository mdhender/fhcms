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

/* The following routine will return a score indicating how closely two
 * strings match.  If the score is exactly 10000, then the strings are
 * identical.  Otherwise, the value returned is the number of character
 * matches, allowing for accidental transpositions, insertions, and
 * deletions.  Excess characters in either string will subtract from
 * the score.  Thus, it's possible for a score to be negative.
 *
 * In general, if the strings are at least 7 characters each, then you can
 * assume the strings are the same if the highest score equals the length of
 * the correct string, length-1, or length-2, AND if the score of the next
 * best match is less than the highest score.  A non-10000 score will never
 * be higher than the length of the correct string. */
func agrep_score(correct_string, unknown_string []byte) int {
	if strcmp(correct_string, unknown_string) == 0 {
		return 10000
	}

	score, p1, p2 := 0, correct_string, unknown_string
	for len(p1) != 0 && len(p2) != 0 {
		c1, c2 := p1[0], p2[0]
		if c1 == 0 || c2 == 0 {
			break
		}
		if len(p1) == 1 {
			p1 = []byte{0}
		} else {
			p1 = p1[1:]
		}
		if len(p2) == 1 {
			p2 = []byte{0}
		} else {
			p2 = p2[1:]
		}
		if c1 == c2 {
			score++
		} else if c1 == p2[0] && c2 == p1[0] { /* Transposed. */
			score += 2
			p1, p2 = p1[1:], p2[1:]
		} else if c1 == p2[0] { /* Unneeded character. */
			score++
			p2 = p2[1:]
		} else if c2 == p1[0] { /* Missing character. */
			score++
			p1 = p1[1:]
		}
	}

	// reduce score by excess characters, if any
	return score - strlen(p1) - strlen(p2)
}

func b2s(s []byte) string {
	return string(s[:strlen(s)])
}

// check_population will set or clear the POPULATED bit for a nampla.
// (also sometimes the MINING_COLONY and RESORT_COLONY bits, too.)
// it will return true if the nampla is currently populated, false if not.
// it will also check if a message associated with this planet should be logged.
func (e *Engine) check_population(nampla *nampla_data) bool {
	total_pop := nampla.mi_base + nampla.ma_base + nampla.IUs_to_install + nampla.AUs_to_install + nampla.item_quantity[PD] + nampla.item_quantity[CU] + nampla.pop_units

	if total_pop <= 0 {
		nampla.status &= (^(POPULATED | MINING_COLONY | RESORT_COLONY))
	} else {
		if was_already_populated := (nampla.status & POPULATED) != 0; !was_already_populated {
			if nampla.message != 0 {
				// there is a message that must be logged whenever this planet becomes populated for the first time
				filename := fmt.Sprintf("message%ld.txt", nampla.message)
				e.log_message(filename)
			}
		}
		nampla.status |= POPULATED
	}

	return total_pop > 0
}

func (e *Engine) closest_unvisited_star_orders(ship *ship_data) {
	e.x = -1 // signals that we have found no destination star

	var closest_star *star_data
	closest_distance, found := 999999, FALSE
	for i := 0; i < e.num_stars; i++ {
		star := e.star_base[i]
		if star.visited_by[e.species_index] != FALSE {
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

	if found != FALSE {
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

// delete_ship will delete a ship record.
// it's more of a logical delete than physical.
func (e *Engine) delete_ship(ship *ship_data) {
	// update all fields to their zero values. mostly sort of
	ship.name = "Unused"
	ship.x, ship.y, ship.z = 0, 0, 0
	ship.pn = 99 // todo: this is a flag for 'deleted'
	ship.status = 0
	ship._type = 0
	ship.dest_x, ship.dest_y, ship.dest_z = 0, 0, 0
	ship.just_jumped = FALSE
	ship.arrived_via_wormhole = FALSE
	ship.reserved1 = 0
	ship.reserved2 = 0
	ship.reserved3 = 0
	ship.class = 0
	ship.tonnage = 0
	for i := 0; i < MAX_ITEMS; i++ {
		ship.item_quantity[i] = 0

	}
	ship.age = 0
	ship.remaining_cost = 0
	ship.reserved4 = 0
	ship.loading_point = 0
	ship.unloading_point = 0
	ship.special = 0
}

// distorted provides the 'distorted' species number used to identify a species that uses field distortion units.
// The input variable 'species_number' is the same number used in filename creation for the species.
func (e *Engine) distorted(species_number int) int {
	// We must use the LS tech level at the start of the turn because
	// the distorted species number must be the same throughout the
	// turn, even if the tech level changes during production.
	ls := e.spec_data[species_number-1].init_tech_level[LS]
	i := species_number & 0x000F        // lower four bits
	j := (species_number >> 4) & 0x000F // upper four bits
	return (ls%5+3)*(4*i+j) + (ls%11 + 7)
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

func (e *Engine) rnd(max int) int {
	return e.defaultPRNG.Roll(max)
}
func (e *Engine) rndGetSeed() uint64 {
	return e.defaultPRNG.GetSeed()
}
func (e *Engine) rndSetSeed(s uint64) {
	e.defaultPRNG.Seed(s)
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

func (e *Engine) undistorted(distorted_species_number int) int {
	for i := 0; i < MAX_SPECIES; i++ {
		species_number := i + 1
		if e.distorted(species_number) == distorted_species_number {
			return species_number
		}
	}
	return 0 // not a legitimate species
}
