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

/* This routine will return TRUE if forced jump or misjump units are used,
   even if they fail. It will return FALSE if the attacker has none or
   not enough. */

func (e *Engine) forced_jump_units_used(attacker_index, defender_index int, total_shots *int, bat *battle_data, act *action_data) int {
	//var i, att_sp_index, def_sp_index, attacker_gv, defender_gv, _type, fj_num, fm_num, number, success_chance int

	/* Make sure attacking unit is a starbase. */
	attacking_ship, ok := act.fighting_unit[attacker_index].(*ship_data)
	if !ok {
		panic("act.fighting_unit[attacker_index].(*ship_data); !ok")
	}
	if attacking_ship._type != STARBASE {
		return FALSE
	}

	/* See if attacker has any forced jump units. */
	fj_num, fm_num := attacking_ship.item_quantity[FJ], attacking_ship.item_quantity[FM]
	if fj_num == 0 && fm_num == 0 {
		return FALSE
	}

	/* If both types are being carried, choose one randomly. */
	var jumpType, quantity int
	if fj_num > 0 && fm_num > 0 {
		if e.rnd(2) == 1 {
			jumpType, quantity = FJ, fj_num
		} else {
			jumpType, quantity = FM, fm_num
		}
	} else if fj_num > 0 {
		jumpType, quantity = FJ, fj_num
	} else {
		jumpType, quantity = FM, fm_num
	}

	/* Get gravitics tech levels. */
	att_sp_index := act.fighting_species_index[attacker_index]
	attacker_gv := e.c_species[att_sp_index].tech_level[GV]

	def_sp_index := act.fighting_species_index[defender_index]
	defender_gv := e.c_species[def_sp_index].tech_level[GV]

	/* Check if sufficient units are available. */
	defending_ship, ok := act.fighting_unit[defender_index].(*ship_data)
	if !ok {
		panic("act.fighting_unit[defender_index].(*ship_data); !ok")
	}
	if quantity < defending_ship.tonnage {
		return FALSE
	}

	/* Make sure defender is not a starbase. */
	if defending_ship._type == STARBASE {
		return FALSE
	}

	/* Calculate percent chance of success. */
	success_chance := 2 * ((quantity - defending_ship.tonnage) + (attacker_gv - defender_gv))

	/* See if it worked. */
	var failure int
	if e.rnd(100) > success_chance {
		failure, e.log_summary = TRUE, FALSE
	} else {
		failure, e.log_summary = FALSE, TRUE
	}
	e.log_string("        ")
	e.log_string(e.ship_name(attacking_ship))
	e.log_string(" attempts to use ")
	e.log_string(item_name[jumpType])
	e.log_string("s against ")

	//e.ignore_field_distorters = !e.field_distorted[def_sp_index]
	if e.field_distorted[def_sp_index] != FALSE {
		e.ignore_field_distorters = FALSE
	} else {
		e.ignore_field_distorters = TRUE
	}
	e.log_string(e.ship_name(defending_ship))
	e.ignore_field_distorters = FALSE

	if failure != FALSE {
		e.log_string(", but fails.\n")
		return TRUE
	}

	e.log_string(", and succeeds!\n")
	e.log_summary = FALSE

	/* Determine destination. */
	if jumpType == FM { /* Destination is totally random. */
		defending_ship.dest_x, defending_ship.dest_y, defending_ship.dest_z = e.rnd(100)-1, e.rnd(100)-1, e.rnd(100)-1
	} else { /* Random location close to battle. */
		var i int
		for i = 3; i == 3; {
			i = e.rnd(5)
		}
		if defending_ship.dest_x = bat.x + i - 3; defending_ship.dest_x < 0 {
			defending_ship.dest_x = 0
		}

		for i = 3; i == 3; {
			i = e.rnd(5)
		}
		if defending_ship.dest_y = bat.y + i - 3; defending_ship.dest_y < 0 {
			defending_ship.dest_y = 0
		}

		for i = 3; i == 3; {
			i = e.rnd(5)
		}
		if defending_ship.dest_z = bat.z + i - 3; defending_ship.dest_z < 0 {
			defending_ship.dest_z = 0
		}
	}

	/* Make sure this ship can no longer take part in the battle. */
	*total_shots -= act.shots_left[defender_index]
	defending_ship.status = FORCED_JUMP
	defending_ship.pn = -1
	act.shots_left[defender_index] = 0
	act.num_shots[defender_index] = 0

	return TRUE
}
