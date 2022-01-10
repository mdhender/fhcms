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
	"bytes"
	"fmt"
	"log"
	"time"
)

/* This program will generate default orders for a species if no explicit orders have been provided. */
func (e *Engine) no_orders() {
	//int i, j, k, ship_index, locations_fd, my_loc_index, nampla_index, its_loc_index, tonnage, found, alien_number, alien_index, array_index, bit_number, ls_needed, production_penalty;
	var i, j, k, n, nampla_index, nn int
	//char *random_name(), message_line[132];
	//long n, nn, raw_material_units, production_capacity, balance, current_base, CUs_needed, IUs_needed, AUs_needed, EUs, bit_mask;
	//
	//FILE *message_file, *log_file;
	//
	//struct species_data *alien;
	//struct e.nampla_data * e.nampla, *temp_nampla;
	var home_nampla *nampla_data
	//struct ship_data *   ship;
	//struct sp_loc_data * locations_base, *my_loc, *its_loc;

	e.truncate_name = TRUE

	/* Major loop. Check each species in the game. */
	for e.species_number = 1; e.species_number <= e.galaxy.num_species; e.species_number++ {
		e.species_index = e.species_number - 1

		/* Check if we have orders. */
		if e.spec_orders[e.species_index] != nil {
			continue
		}

		e.species = e.spec_data[e.species_index]
		e.nampla_base = e.namp_data[e.species_index]
		e.ship_base = e.ship_data[e.species_index]
		home_nampla = e.nampla_base[0]
		e.home_planet = e.planet_base[home_nampla.planet_index]

		for i = 0; i < e.species.num_ships; i++ {
			e.ship = e.ship_base[i]
			e.ship.special = 0
		}

		/* Print message for gamemaster. */
		log.Printf("[engine] generating orders for species #%02d, SP %s...\n", e.species_number, e.species.name)

		// copy message to log file
		e.spec_logs[e.species_index].WriteString(`
Note from the gamemaster:

    Unfortunately, I did not receive your orders for the current turn on
    time.  Please make an effort to be on time for the next turn.  If you
    DID send me orders but they arrived late, I will use them for the next
    turn (but keep in mind that many of them may not work properly).  You
    may also change and re-submit them to me based on the information in
    THIS status report.  However, if you do re-submit them, please make
    an effort to submit them by the deadline.  Thank you.

`)

		// create bytes.Buffer for the default orders we're about to create
		e.orders_file = &bytes.Buffer{}

		e.orders_file.WriteString(fmt.Sprintf(";; SP%02d %s -- Default Orders %s\n\n", e.species_number, e.species.name, time.Now().UTC().Format(time.RFC3339)))

		/* Issue PRE-DEPARTURE orders. */
		e.orders_file.WriteString("START PRE-DEPARTURE\n")
		e.orders_file.WriteString("; Place pre-departure orders here.\n\n")

		for nampla_index = 0; nampla_index < e.species.num_namplas; nampla_index++ {
			e.nampla = e.nampla_base[nampla_index]
			if e.nampla.pn == 99 {
				continue
			}

			/* Generate auto-installs for colonies that were loaded via
			 *  the DEVELOP command. */
			if e.nampla.auto_IUs != 0 {
				e.orders_file.WriteString(fmt.Sprintf("\tInstall\t%d IU\tPL %s\n", e.nampla.auto_IUs, e.nampla.name))
			}
			if e.nampla.auto_AUs != 0 {
				e.orders_file.WriteString(fmt.Sprintf("\tInstall\t%d AU\tPL %s\n", e.nampla.auto_AUs, e.nampla.name))
			}
			if e.nampla.auto_IUs != 0 || e.nampla.auto_AUs != 0 {
				e.orders_file.WriteString("\n")
			}

			e.nampla.item_quantity[CU] -= e.nampla.auto_IUs + e.nampla.auto_AUs

			/* Generate auto UNLOAD orders for transports at this e.nampla. */
			for j = 0; j < e.species.num_ships; j++ {
				e.ship = e.ship_base[j]
				if e.ship.pn == 99 {
					continue
				} else if e.ship.x != e.nampla.x {
					continue
				} else if e.ship.y != e.nampla.y {
					continue
				} else if e.ship.z != e.nampla.z {
					continue
				} else if e.ship.pn != e.nampla.pn {
					continue
				} else if e.ship.status == JUMPED_IN_COMBAT {
					continue
				} else if e.ship.status == FORCED_JUMP {
					continue
				} else if e.ship.class != TR {
					continue
				} else if e.ship.item_quantity[CU] < 1 {
					continue
				}

				/* New colonies will never be started automatically unless ship was loaded via a DEVELOP order. */
				if e.ship.loading_point != 0 {
					/* Check if transport is at specified unloading point. */
					n = e.ship.unloading_point
					if n == nampla_index || (n == 9999 && nampla_index == 0) {
						goto unload_ship
					}
				}

				if (e.nampla.status & POPULATED) == 0 {
					continue
				}

				if (e.nampla.mi_base + e.nampla.ma_base) >= 2000 {
					continue
				}

				if e.nampla.x == e.nampla_base[0].x && e.nampla.y == e.nampla_base[0].y && e.nampla.z == e.nampla_base[0].z {
					continue /* Home sector. */
				}
			unload_ship:

				n = e.ship.loading_point
				if n == 9999 {
					n = 0 /* Home planet. */
				}
				if n == nampla_index {
					continue /* Ship was just loaded here. */
				}
				e.orders_file.WriteString(fmt.Sprintf("\tUnload\tTR%d%s %s\n\n", e.ship.tonnage, ship_type[e.ship._type], e.ship.name))

				e.nampla.item_quantity[CU] = 0

				e.ship.special = e.ship.loading_point

				// n = e.nampla - e.nampla_base
				for n = 0; n < len(e.nampla_base); n++ {
					if e.nampla == e.nampla_base[n] {
						break
					}
				}
				if n == 0 {
					n = 9999
				}
				e.ship.unloading_point = n
			}

			if (e.nampla.status & HOME_PLANET) != 0 {
				continue
			} else if e.nampla.item_quantity[CU] == 0 {
				continue
			} else if e.nampla.item_quantity[IU] == 0 && e.nampla.item_quantity[AU] == 0 {
				continue
			}

			if e.nampla.item_quantity[IU] > 0 {
				e.orders_file.WriteString(fmt.Sprintf("\tInstall\t0 IU\tPL %s\n", e.nampla.name))
			}
			if e.nampla.item_quantity[AU] > 0 {
				e.orders_file.WriteString(fmt.Sprintf("\tInstall\t0 AU\tPL %s\n\n", e.nampla.name))
			}
		}

		e.orders_file.WriteString("END\n\n")

		e.orders_file.WriteString("START JUMPS\n")
		e.orders_file.WriteString("; Place jump orders here.\n\n")

		/* Initialize to make sure ships are not given more than one JUMP order. */
		for i = 0; i < e.species.num_ships; i++ {
			e.ship = e.ship_base[i]
			e.ship.just_jumped = FALSE
		}

		/* Generate auto-jumps for ships that were loaded via the DEVELOP
		 * command or which were UNLOADed because of the AUTO command. */
		for i = 0; i < e.species.num_ships; i++ {
			e.ship = e.ship_base[i]

			if e.ship.status == JUMPED_IN_COMBAT {
				continue
			} else if e.ship.status == FORCED_JUMP {
				continue
			} else if e.ship.pn == 99 {
				continue
			} else if e.ship.just_jumped != 0 {
				continue
			}

			j = e.ship.special
			if j != 0 {
				if j == 9999 {
					j = 0 /* Home planet. */
				}
				temp_nampla := e.nampla_base[j]

				e.orders_file.WriteString(fmt.Sprintf("\tJump\t%s, PL %s\t; ", e.ship_name(e.ship), temp_nampla.name))

				e.print_mishap_chance_orders(e.ship, temp_nampla.x, temp_nampla.y, temp_nampla.z)

				e.orders_file.WriteString("\n\n")

				e.ship.just_jumped = TRUE

				continue
			}

			n = e.ship.unloading_point
			if n != 0 {
				if n == 9999 {
					n = 0 /* Home planet. */
				}
				temp_nampla := e.nampla_base[n]

				if temp_nampla.x == e.ship.x && temp_nampla.y == e.ship.y && temp_nampla.z == e.ship.z {
					continue
				}

				e.orders_file.WriteString(fmt.Sprintf("\tJump\t%s, PL %s\t; ", e.ship_name(e.ship), temp_nampla.name))

				e.print_mishap_chance_orders(e.ship, temp_nampla.x, temp_nampla.y, temp_nampla.z)

				e.orders_file.WriteString("\n\n")

				e.ship.just_jumped = TRUE
			}
		}

		/* Generate JUMP orders for all TR1s. */
		for i = 0; i < e.species.num_ships; i++ {
			e.ship = e.ship_base[i]
			if e.ship.pn == 99 {
				continue
			} else if e.ship.status == UNDER_CONSTRUCTION {
				continue
			} else if e.ship.status == JUMPED_IN_COMBAT {
				continue
			} else if e.ship.status == FORCED_JUMP {
				continue
			} else if e.ship.just_jumped != 0 {
				continue
			}

			if e.ship.class == TR && e.ship.tonnage == 1 && e.ship._type == FTL {
				e.orders_file.WriteString(fmt.Sprintf("\tJump\tTR1 %s, ", e.ship.name))
				e.closest_unvisited_star_orders(e.ship)
				e.orders_file.WriteString(fmt.Sprintf("\n\t\t\t; Age %d, now at %d %d %d, ", e.ship.age, e.ship.x, e.ship.y, e.ship.z))

				e.print_mishap_chance_orders(e.ship, e.x, e.y, e.z)

				e.ship.dest_x = e.x
				e.ship.dest_y = e.y
				e.ship.dest_z = e.z

				e.orders_file.WriteString("\n\n")

				e.ship.just_jumped = TRUE
			}
		}

		e.orders_file.WriteString("END\n\n")

		e.orders_file.WriteString("START PRODUCTION\n")

		/* Generate a PRODUCTION order for each planet that can produce. */
		for nampla_index = e.species.num_namplas - 1; nampla_index >= 0; nampla_index-- {
			e.nampla = e.nampla_base[nampla_index]
			if e.nampla.pn == 99 {
				continue
			}

			if e.nampla.mi_base == 0 && (e.nampla.status&RESORT_COLONY) == 0 {
				continue
			}
			if e.nampla.ma_base == 0 && (e.nampla.status&MINING_COLONY) == 0 {
				continue
			}

			e.orders_file.WriteString(fmt.Sprintf("    PRODUCTION PL %s\n", e.nampla.name))

			if (e.nampla.status & MINING_COLONY) != 0 {
				e.orders_file.WriteString("    ; The above PRODUCTION order is required for this mining colony, even\n")
				e.orders_file.WriteString("    ;  if no other production orders are given for it.\n")
			} else if (e.nampla.status & RESORT_COLONY) != 0 {
				e.orders_file.WriteString("    ; The above PRODUCTION order is required for this resort colony, even\n")
				e.orders_file.WriteString("    ;  though no other production orders can be given for it.\n")
			} else {
				e.orders_file.WriteString(fmt.Sprintf("    ; Place production orders here for planet %s.\n\n", e.nampla.name))
			}

			/* Build IUs and AUs for incoming ships with CUs. */
			if e.nampla.IUs_needed != 0 {
				e.orders_file.WriteString(fmt.Sprintf("\tBuild\t%d IU\n", e.nampla.IUs_needed))
			}
			if e.nampla.AUs_needed != 0 {
				e.orders_file.WriteString(fmt.Sprintf("\tBuild\t%d AU\n", e.nampla.AUs_needed))
			}
			if e.nampla.IUs_needed != 0 || e.nampla.AUs_needed != 0 {
				e.orders_file.WriteString("\n")
			}

			if (e.nampla.status & MINING_COLONY) != 0 {
				continue
			}
			if (e.nampla.status & RESORT_COLONY) != 0 {
				continue
			}

			/* See if there are any RMs to recycle. */
			n = e.nampla.special / 5
			if n > 0 {
				e.orders_file.WriteString(fmt.Sprintf("\tRecycle\t%d RM\n\n", 5*n))
			}

			/* Generate DEVELOP commands for ships arriving here because of
			 *  AUTO command. */
			for i = 0; i < e.species.num_ships; i++ {
				e.ship = e.ship_base[i]
				if e.ship.pn == 99 {
					continue
				}

				k = e.ship.special
				if k == 0 {
					continue
				}
				if k == 9999 {
					k = 0 /* Home planet. */
				}
				if e.nampla != e.nampla_base[k] {
					continue
				}

				k = e.ship.unloading_point
				if k == 9999 {
					k = 0
				}
				temp_nampla := e.nampla_base[k]

				e.orders_file.WriteString(fmt.Sprintf("\tDevelop\tPL %s, TR%d%s %s\n\n", temp_nampla.name, e.ship.tonnage, ship_type[e.ship._type], e.ship.name))
			}

			/* Give orders to continue construction of unfinished ships and
			 *  starbases. */
			for i = 0; i < e.species.num_ships; i++ {
				e.ship = e.ship_base[i]
				if e.ship.pn == 99 {
					continue
				} else if e.ship.x != e.nampla.x {
					continue
				} else if e.ship.y != e.nampla.y {
					continue
				} else if e.ship.z != e.nampla.z {
					continue
				} else if e.ship.pn != e.nampla.pn {
					continue
				} else if e.ship.status == UNDER_CONSTRUCTION {
					e.orders_file.WriteString(fmt.Sprintf("\tContinue\t%s, %d\t; Left to pay = %d\n\n", e.ship_name(e.ship), e.ship.remaining_cost, e.ship.remaining_cost))
					continue
				} else if e.ship._type != STARBASE {
					continue
				}
				j = (e.species.tech_level[MA] / 2) - e.ship.tonnage
				if j < 1 {
					continue
				}

				e.orders_file.WriteString(fmt.Sprintf("\tContinue\tBAS %s, %d\t; Current tonnage = %s\n\n", e.ship.name, 100*j, commas(10000*e.ship.tonnage)))
			}

			/* Generate DEVELOP command if this is a colony with an economic base less than 200. */
			n = e.nampla.mi_base + e.nampla.ma_base + e.nampla.IUs_needed + e.nampla.AUs_needed
			if (e.nampla.status&COLONY) != 0 && n < 2000 && e.nampla.pop_units > 0 {
				if e.nampla.pop_units > (2000 - n) {
					nn = 2000 - n
				} else {
					nn = e.nampla.pop_units
				}
				e.orders_file.WriteString(fmt.Sprintf("\tDevelop\t%d\n\n", 2*nn))
				e.nampla.IUs_needed += nn
			}

			/* For home planets and any colonies that have an economic base of
			 *  at least 200, check if there are other colonized planets in
			 *  the same sector that are not self-sufficient.  If so, DEVELOP
			 *  them. */
			if n >= 2000 || (e.nampla.status&HOME_PLANET) != 0 {
				for i = 1; i < e.species.num_namplas; i++ { /* Skip HP. */
					if i == nampla_index {
						continue
					}
					temp_nampla := e.nampla_base[i]
					if temp_nampla.pn == 99 {
						continue
					} else if temp_nampla.x != e.nampla.x {
						continue
					} else if temp_nampla.y != e.nampla.y {
						continue
					} else if temp_nampla.z != e.nampla.z {
						continue
					}
					n = temp_nampla.mi_base + temp_nampla.ma_base + temp_nampla.IUs_needed + temp_nampla.AUs_needed
					if n == 0 {
						continue
					}
					nn = temp_nampla.item_quantity[IU] + temp_nampla.item_quantity[AU]
					if nn > temp_nampla.item_quantity[CU] {
						nn = temp_nampla.item_quantity[CU]
					}
					n += nn
					if n >= 2000 {
						continue
					}
					nn = 2000 - n
					if nn > e.nampla.pop_units {
						nn = e.nampla.pop_units
					}
					e.orders_file.WriteString(fmt.Sprintf("\tDevelop\t%d\tPL %s\n\n", 2*nn, temp_nampla.name))
					temp_nampla.AUs_needed += nn
				}
			}
		}

		e.orders_file.WriteString("END\n\n")

		e.orders_file.WriteString("START POST-ARRIVAL\n")
		e.orders_file.WriteString("; Place post-arrival orders here.\n\n")

		/* Generate an AUTO command. */
		e.orders_file.WriteString("\tAuto\n\n")

		/* Generate SCAN orders for all TR1s in sectors that current species does not inhabit. */
		for i = 0; i < e.species.num_ships; i++ {
			e.ship = e.ship_base[i]
			if e.ship.pn == 99 {
				continue
			} else if e.ship.status == UNDER_CONSTRUCTION {
				continue
			} else if e.ship.class != TR {
				continue
			} else if e.ship.tonnage != 1 {
				continue
			} else if e.ship._type != FTL {
				continue
			} else if e.ship.dest_x == -1 {
				continue /* Not jumping anywhere. */
			}
			found := FALSE
			for j = 1; j < e.species.num_namplas; j++ { /* Skip home sector. */
				e.nampla = e.nampla_base[j]
				if e.nampla.pn == 99 {
					continue
				} else if e.nampla.x != e.ship.dest_x {
					continue
				} else if e.nampla.y != e.ship.dest_y {
					continue
				} else if e.nampla.z != e.ship.dest_z {
					continue
				} else if (e.nampla.status & POPULATED) != 0 {
					found = TRUE
					break
				}
			}
			if found == FALSE {
				e.orders_file.WriteString(fmt.Sprintf("\tScan\tTR1 %s\n", e.ship.name))
			}
		}

		e.orders_file.WriteString("END\n\n")

		// replace any existing orders with these generated orders
		e.spec_orders[e.species_index] = e.orders_file.Bytes()
	}
}
