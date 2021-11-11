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

import (
	"github.com/mdhender/fhcms/cms/parser"
	"log"
)

//*************************************************************************
// Jump.c

func (g *globals) executeJumpOrders() []error {
	g.ignore_field_distorters = true
	/* For these commands, do not display age or landed/orbital status of ships. */
	g.truncate_name = true
	/* We will control value of log_file from here. */
	g.log_stdout = false

	// process jump orders for each species in the input list
	for sp_index := 0; sp_index < len(g.sp_num); sp_index++ {
		species_number = g.sp_num[sp_index]
		g.setSpecies(species_number - 1)

		// if ships of a species withdrew or were forced from combat and there were no jump orders for that species, then combat jumps will not take place.
		// this variable allows us to handle them separately.
		g.species.species_jumped = false

		// skip if this species has no jump orders
		var hasJumpOrders bool
		if g.species.orders.data != nil {
			for _, o := range g.species.orders.data.Children {
				if o.Token.Type == parser.Jump {
					hasJumpOrders = len(o.Children) != 0
				}
			}
		}
		if !hasJumpOrders {
			log.Printf("jump: species %02d has no jump orders\n", species_number)
			continue
		}

		g.log_file = g.species.log_file
		end_of_file = false

		// handle jump orders for this species
		g.log("\nJump orders:\n")
		g.do_jump_orders()
		g.species.species_jumped = true

		// take care of any ships that withdrew or were forced to jump during combat
		for ship_index = 0; ship_index < g.species.num_ships; ship_index++ {
			ship = g.ship_base[ship_index]
			if ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT {
				do_JUMP_command(true, false)
			}
		}
	}

	// take care of any ships that withdrew from combat but were not handled above because no jump orders were received for species.
	g.log_stdout = true
	for species_number = 1; species_number <= galaxy.num_species; species_number++ {
		g.setSpecies(species_number)
		if g.species.species_jumped {
			continue
		}

		firstLine := true
		g.log_file = g.species.log_file

		for ship_index = 0; ship_index < g.species.num_ships; ship_index++ {
			ship = g.ship_base[ship_index]
			if ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT {
				if firstLine {
					g.log("\nWithdrawals and forced jumps during combat:\n")
					firstLine = false
				}
				do_JUMP_command(true, false)
			}
		}
	}

	return nil
}

func (g *globals) do_jump_orders() {
	log.Printf("\nStart of jump orders for species #%d, SP %s...\n", g.species.id, g.species.name)
	var command int
	for {
		command = get_command()

		if command == 0 {
			g.log("!!! Order ignored:\n")
			g.log("!!! %s", input_line)
			g.log("!!! Unknown or missing command.\n")
			continue
		}

		if end_of_file || command == END {
			log.Printf("End of jump orders for species #%d, SP %s.\n", g.species.id, g.species.name)
			break /* END for this species. */
		}

		switch command {
		case JUMP:
			do_JUMP_command(false, false)
			break

		case MOVE:
			do_MOVE_command()
			break

		case PJUMP:
			do_JUMP_command(false, true)
			break

		case VISITED:
			do_VISITED_command()
			break

		case WORMHOLE:
			do_WORMHOLE_command()
			break

		default:
			fprintf(log_file, "!!! Order ignored:\n")
			fprintf(log_file, "!!! %s", input_line)
			fprintf(log_file, "!!! Invalid jump command.\n")
		}
	}
}
