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

// withdrawal_check checks all fighting ships and see if any wish to withdraw.
// If so, it will set the ship's status to JUMPED_IN_COMBAT.
// The actual jump will be handled by the Jump program.
func (e *Engine) withdrawal_check(bat *battle_data, act *action_data) {
	// stash the global truncate_name so that our logs can show the age of the ship
	old_trunc := e.truncate_name
	e.truncate_name = FALSE

	/* Compile statistics and handle individual ships that must leave. */
	var num_ships_gone, num_ships_total [MAX_SPECIES]int
	for ship_index := 0; ship_index < act.num_units_fighting; ship_index++ {
		if act.unit_type[ship_index] != SHIP {
			continue
		}
		sh, ok := act.fighting_unit[ship_index].(*ship_data)
		if !ok {
			panic("act.fighting_unit[ship_index].(*ship_data); !ok")
		}
		species_index := act.fighting_species_index[ship_index]
		num_ships_total[species_index]++

		if sh.status == JUMPED_IN_COMBAT { /* Already withdrawn. */
			num_ships_gone[species_index]++
			continue
		} else if sh.status == FORCED_JUMP { /* Forced to leave. */
			num_ships_gone[species_index]++
			continue
		} else if sh.age > 49 { /* Already destroyed. */
			num_ships_gone[species_index]++
			continue
		} else if sh._type != FTL { /* Ship can't jump. */
			continue
		}

		var withdraw_age int
		if sh.class == TR {
			/* Transports will withdraw only when entire fleet withdraws. */
			withdraw_age = bat.transport_withdraw_age[species_index]
			if withdraw_age == 0 {
				continue
			}
		} else {
			withdraw_age = bat.warship_withdraw_age[species_index]
		}

		if sh.age > withdraw_age {
			act.num_shots[ship_index] = 0
			act.shots_left[ship_index] = 0
			sh.pn = 0 // todo: why?

			//e.ignore_field_distorters = !e.field_distorted[species_index]
			if e.field_distorted[species_index] == FALSE {
				e.ignore_field_distorters = TRUE
			} else {
				e.ignore_field_distorters = FALSE
			}

			fprintf(e.log_file, "        %s jumps away from the battle.\n", e.ship_name(sh))
			fprintf(e.summary_file, "        %s jumps away from the battle.\n", e.ship_name(sh))
			e.ignore_field_distorters = FALSE

			sh.dest_x = bat.haven_x[species_index]
			sh.dest_y = bat.haven_y[species_index]
			sh.dest_z = bat.haven_z[species_index]

			sh.status = JUMPED_IN_COMBAT

			num_ships_gone[species_index]++
		}
	}

	/* Now check if a fleet has reached its limit. */
	for ship_index := 0; ship_index < act.num_units_fighting; ship_index++ {
		if act.unit_type[ship_index] != SHIP {
			continue
		}
		sh, ok := act.fighting_unit[ship_index].(*ship_data)
		if !ok {
			panic("act.fighting_unit[ship_index].(*ship_data); !ok")
		}
		species_index := act.fighting_species_index[ship_index]

		if sh._type != FTL { /* Ship can't jump. */
			continue
		} else if sh.status == JUMPED_IN_COMBAT { /* Already withdrawn. */
			continue
		} else if sh.status == FORCED_JUMP { /* Already gone. */
			continue
		} else if sh.age > 49 { /* Already destroyed. */
			continue
		}

		var percent_loss int
		if bat.fleet_withdraw_percentage[species_index] == 0 {
			percent_loss = 101 /* Always withdraw immediately. */
		} else {
			percent_loss = (100 * num_ships_gone[species_index]) / num_ships_total[species_index]
		}

		if percent_loss > bat.fleet_withdraw_percentage[species_index] {
			act.num_shots[ship_index] = 0
			act.shots_left[ship_index] = 0
			sh.pn = 0

			// e.ignore_field_distorters = !e.field_distorted[species_index]
			if e.field_distorted[species_index] == FALSE {
				e.ignore_field_distorters = TRUE
			} else {
				e.ignore_field_distorters = FALSE
			}

			fprintf(e.log_file, "        %s jumps away from the battle.\n", e.ship_name(sh))
			fprintf(e.summary_file, "        %s jumps away from the battle.\n", e.ship_name(sh))
			e.ignore_field_distorters = FALSE

			sh.dest_x = bat.haven_x[species_index]
			sh.dest_y = bat.haven_y[species_index]
			sh.dest_z = bat.haven_z[species_index]

			sh.status = JUMPED_IN_COMBAT
		}
	}

	// restore the stashed global value
	e.truncate_name = old_trunc
}
