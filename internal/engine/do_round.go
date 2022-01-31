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

import "fmt"

// do_round will return TRUE if a round of combat actually occurred.
// Otherwise, it will return false.
func (e *Engine) do_round(option, round_number int, bat *battle_data, act *action_data) int {
	var (
		i, j, n, unit_index, combat_occurred, total_shots          int
		attacker_index, defender_index, chance_to_hit              int
		attacker_ml, attacker_gv, defender_ml                      int
		target_index                                               [MAX_SHIPS]int
		num_targets, header_printed, num_sp, fj_chance, shields_up int
		FDs_were_destroyed                                         int
		di                                                         [3]int
		start_unit, current_species                                int
		this_is_a_hijacking                                        int
		units_destroyed, percent_decrease                          int
		damage_done, damage_to_ship, damage_to_shields, op1, op2   int
		original_cost, recycle_value, economic_units               int
		attacker_name, defender_name                               string
		attacking_species, defending_species                       *species_data
		sh, attacking_ship, defending_ship                         *ship_data
		attacking_nampla, defending_nampla                         *nampla_data
		ok                                                         bool
		//aux_shield_power, tons, found                              int
	)

	/* Clear out x_attacked_y and germ_bombs_used arrays.  They will be used to log who bombed who, or how many GWs were used. */
	num_sp = bat.num_species_here
	for i = 0; i < num_sp; i++ {
		for j = 0; j < num_sp; j++ {
			e.x_attacked_y[i][j] = FALSE
			e.germ_bombs_used[i][j] = 0
		}
	}

	/* If a species has ONLY non-combatants left, then let them fight. */
	start_unit = 0
	total_shots = 0
	current_species = act.fighting_species_index[0]
	for unit_index = 0; unit_index < act.num_units_fighting; unit_index++ {
		if act.fighting_species_index[unit_index] != current_species {
			if total_shots == 0 {
				/* Convert all non-combatants, if any, to combatants. */
				for i = start_unit; i < unit_index; i++ {
					if act.unit_type[i] == SHIP {
						sh, ok = act.fighting_unit[i].(*ship_data)
						if !ok {
							panic("act.fighting_unit[i].(*ship_data); !ok")
						}
						sh.special = 0
					}
				}
			}
			start_unit = unit_index
			total_shots = 0
		}

		n = act.num_shots[unit_index]
		if act.surprised[unit_index] != FALSE {
			n = 0
		}
		if act.unit_type[unit_index] == SHIP {
			sh, ok = act.fighting_unit[unit_index].(*ship_data)
			if !ok {
				panic("act.fighting_unit[i].(*ship_data); !ok")
			}
			if sh.special == NON_COMBATANT {
				n = 0
			}
		}
		total_shots += n
	}

	/* Determine total number of shots for all species present. */
	total_shots = 0
	for unit_index = 0; unit_index < act.num_units_fighting; unit_index++ {
		n = act.num_shots[unit_index]
		if act.surprised[unit_index] != FALSE {
			n = 0
		}
		if act.unit_type[unit_index] == SHIP {
			sh, ok = act.fighting_unit[unit_index].(*ship_data)
			if !ok {
				panic("act.fighting_unit[i].(*ship_data); !ok")
			}
			if sh.special == NON_COMBATANT {
				n = 0
			}
		}
		act.shots_left[unit_index] = n
		total_shots += n
	}

	/* Handle all shots. */
	header_printed = FALSE
	combat_occurred = FALSE
	for total_shots > 0 {
		// check to make sure we aren't in infinite loop.
		// that can happen when there are shots remaining but the side with the shots has no more ships left.
		for i = 0; i < act.num_units_fighting; i++ {
			attacking_ship, ok = act.fighting_unit[i].(*ship_data)
			if !ok {
				panic("act.fighting_unit[i].(*ship_data); !ok")
			}
			if attacking_ship.age > 49 || attacking_ship.status == FORCED_JUMP || attacking_ship.status == JUMPED_IN_COMBAT || (attacking_ship.special == NON_COMBATANT && option != GERM_WARFARE) {
				total_shots -= act.shots_left[i]
				act.shots_left[i] = 0
			}
		}
		/* Determine who fires next. */
		attacker_index = e.rnd(act.num_units_fighting) - 1
		if act.unit_type[attacker_index] == SHIP {
			attacking_ship = act.fighting_unit[attacker_index].(*ship_data)
			i = act.fighting_species_index[attacker_index]
			// e.ignore_field_distorters = !e.field_distorted[i]
			if e.field_distorted[i] == FALSE {
				e.ignore_field_distorters = TRUE
			} else {
				e.ignore_field_distorters = FALSE
			}
			attacker_name = e.ship_name(attacking_ship)
			e.ignore_field_distorters = FALSE

			/* Check if ship can fight. */
			if attacking_ship.age > 49 {
				continue
			}
			if attacking_ship.status == FORCED_JUMP {
				continue
			}
			if attacking_ship.status == JUMPED_IN_COMBAT {
				continue
			}
			if attacking_ship.special == NON_COMBATANT && option != GERM_WARFARE {
				continue
			}
		} else {
			attacking_nampla, ok = act.fighting_unit[attacker_index].(*nampla_data)
			if !ok {
				panic("act.fighting_unit[attacker_index].(*nampla_data); !ok")
			}
			attacker_name = fmt.Sprintf("PL %s", attacking_nampla.name)

			/* Check if planet still has defenses. */
			if attacking_nampla.item_quantity[PD] == 0 {
				continue
			}
		}

		/* Make sure attacker is not someone who is being taken by surprise this round. */
		if act.surprised[attacker_index] != FALSE {
			continue
		}

		/* Find an enemy. */
		num_targets = 0
		i = act.fighting_species_index[attacker_index]
		attacker_ml = e.c_species[i].tech_level[ML]
		attacker_gv = e.c_species[i].tech_level[GV]
		for defender_index = 0; defender_index < act.num_units_fighting; defender_index++ {
			j = act.fighting_species_index[defender_index]
			if bat.enemy_mine[i][j] == FALSE {
				continue
			}

			if act.unit_type[defender_index] == SHIP {
				defending_ship, ok = act.fighting_unit[defender_index].(*ship_data)
				if !ok {
					panic("act.fighting_unit[defender_index].(*ship_data); !ok")
				}
				if defending_ship.age > 49 { /* Already destroyed. */
					continue
				} else if defending_ship.status == FORCED_JUMP {
					continue
				} else if defending_ship.status == JUMPED_IN_COMBAT {
					continue
				} else if defending_ship.special == NON_COMBATANT {
					continue
				}
			} else {
				defending_nampla, ok = act.fighting_unit[defender_index].(*nampla_data)
				if !ok {
					panic("act.fighting_unit[defender_index].(*nampla_data); !ok")
				}
				if defending_nampla.item_quantity[PD] == 0 && option == PLANET_ATTACK {
					continue
				}
			}

			target_index[num_targets] = defender_index
			num_targets++
		}

		if num_targets == 0 { /* Attacker has no enemies left. */
			total_shots -= act.shots_left[attacker_index]
			act.shots_left[attacker_index] = 0
			continue
		}

		/* Randomly choose a target. Choose the toughest of four. */
		defender_index = target_index[e.rnd(num_targets)-1]
		op1 = act.num_shots[defender_index] * act.weapon_damage[defender_index]
		di[0] = target_index[e.rnd(num_targets)-1]
		di[1] = target_index[e.rnd(num_targets)-1]
		di[2] = target_index[e.rnd(num_targets)-1]
		for i = 0; i < 3; i++ {
			op2 = act.num_shots[di[i]] * act.weapon_damage[di[i]]
			if op2 > op1 {
				defender_index = di[i]
				op1 = op2
			}
		}

		j = act.fighting_species_index[defender_index]
		defender_ml = e.c_species[j].tech_level[ML]

		if act.unit_type[defender_index] == SHIP {
			defending_ship, ok = act.fighting_unit[defender_index].(*ship_data)
			if !ok {
				panic("act.fighting_unit[defender_index].(*ship_data); !ok")
			}
			//e.ignore_field_distorters = !e.field_distorted[j]
			if e.field_distorted[j] == FALSE {
				e.ignore_field_distorters = TRUE
			} else {
				e.ignore_field_distorters = FALSE
			}
			defender_name = e.ship_name(defending_ship)
			e.ignore_field_distorters = FALSE
		} else {
			defending_nampla, ok = act.fighting_unit[defender_index].(*nampla_data)
			if !ok {
				panic("act.fighting_unit[defender_index].(*nampla_data); !ok")
			}
			defender_name = fmt.Sprintf("PL %s", defending_nampla.name)
		}

		/* Print round number. */
		if header_printed == FALSE {
			e.log_string("      Now doing round ")
			e.log_int(round_number)
			e.log_string(":\n")
			header_printed = TRUE
		}
		attackerGvMl := attacker_gv + attacker_ml
		if attackerGvMl <= 0 {
			attackerGvMl = 1
		}
		// Check if attacker has any forced jump units.
		// The attacker will place more emphasis on the use of these devices if he emphasizes gravitics technology over military technology.
		fj_chance = 50 * attacker_gv / attackerGvMl
		if e.rnd(100) < fj_chance && act.unit_type[attacker_index] == SHIP && act.unit_type[defender_index] == SHIP {
			if e.forced_jump_units_used(attacker_index, defender_index, &total_shots, bat, act) != FALSE {
				combat_occurred = TRUE
				continue
			}
		}

		if act.shots_left[attacker_index] == 0 {
			continue
		}

		/* Since transports generally avoid combat, there is only a 10%
		   chance that they will be targeted, unless they are being
		   explicitly targeted. */
		i = act.fighting_species_index[attacker_index]
		j = act.fighting_species_index[defender_index]
		if act.unit_type[defender_index] == SHIP && defending_ship.class == TR && bat.special_target[i] != TARGET_TRANSPORTS && e.rnd(10) != 5 {
			continue
		}

		/* If a special target has been specified, then there is a 75%
		   chance that it will be attacked if it is available. */
		if bat.special_target[i] != FALSE && e.rnd(100) < 76 {
			if bat.special_target[i] == TARGET_PDS {
				if act.unit_type[defender_index] != SHIP {
					goto fire
				} else {
					continue
				}
			}

			if act.unit_type[defender_index] != SHIP {
				continue
			} else if bat.special_target[i] == TARGET_STARBASES && defending_ship.class != BA {
				continue
			} else if bat.special_target[i] == TARGET_TRANSPORTS && defending_ship.class != TR {
				continue
			} else if bat.special_target[i] == TARGET_WARSHIPS && (defending_ship.class == TR || defending_ship.class == BA) {
				continue
			}
		}

	fire:
		/* Update counts. */
		act.shots_left[attacker_index]--
		total_shots--

		/* Since transports generally avoid combat, there is only a 10% chance that they will attack. */
		if act.unit_type[attacker_index] == SHIP && attacking_ship.class == TR && option != GERM_WARFARE && e.rnd(10) != 5 {
			continue
		}

		/* Fire! */
		combat_occurred = TRUE
		e.log_string("        ")
		e.log_string(attacker_name)
		e.log_string(" fires on ")
		e.log_string(defender_name)
		if act.unit_type[defender_index] == NAMPLA {
			e.log_string(" defenses")
		}

		combinedMl := attacker_ml + defender_ml
		if combinedMl <= 0 {
			combinedMl = 1
		}
		/* Get hit probability.
		 * The basic chance to hit is 1.5 times attackers ML over the sum of attacker's and defender's ML.
		 * Double this value if defender is surprised. */
		chance_to_hit = (150 * attacker_ml) / combinedMl
		if act.surprised[defender_index] == FALSE {
			shields_up = TRUE
		} else {
			chance_to_hit *= 2
			shields_up = FALSE
		}

		/* If defending ship is field-distorted, chance-to-hit is reduced by 25%. */
		j = act.fighting_species_index[defender_index]
		if act.unit_type[defender_index] == SHIP && e.field_distorted[j] != FALSE && defending_ship.item_quantity[FD] == defending_ship.tonnage {
			chance_to_hit = (3 * chance_to_hit) / 4
		}

		if chance_to_hit < 2 {
			chance_to_hit = 2
		} else if chance_to_hit > 98 {
			chance_to_hit = 98
		}

		/* Adjust for age. */
		if act.unit_type[attacker_index] == SHIP {
			chance_to_hit -= (2 * attacking_ship.age * chance_to_hit) / 100
		}

		/* Calculate damage that shot will do if it hits. */
		damage_done = act.weapon_damage[attacker_index]
		damage_done += ((26 - e.rnd(51)) * damage_done) / 100

		/* Take care of attempted annihilation and sieges. */
		if option == PLANET_BOMBARDMENT || option == GERM_WARFARE || option == SIEGE {
			/* Indicate the action that was attempted against this nampla. */
			if option == SIEGE {
				act.unit_type[defender_index] = BESIEGED_NAMPLA
			} else {
				act.unit_type[defender_index] = GENOCIDE_NAMPLA
			}

			/* Indicate who attacked who. */
			i = act.fighting_species_index[attacker_index]
			j = act.fighting_species_index[defender_index]
			e.x_attacked_y[i][j] = TRUE

			/* Update bombardment damage. */
			if option == PLANET_BOMBARDMENT {
				act.bomb_damage[defender_index] += damage_done
			} else if option == GERM_WARFARE {
				if act.unit_type[attacker_index] == SHIP {
					e.germ_bombs_used[i][j] += attacking_ship.item_quantity[GW]
					attacking_ship.item_quantity[GW] = 0
				} else {
					e.germ_bombs_used[i][j] += attacking_nampla.item_quantity[GW]
					attacking_nampla.item_quantity[GW] = 0
				}
			}

			continue
		}

		/* Check if shot hit. */
		if e.rnd(100) <= chance_to_hit {
			e.log_string(" and hits!\n")
		} else {
			e.log_string(" and misses!\n")
			continue
		}

		/* Subtract damage from defender's shields, if they're up. */
		damage_to_ship = 0
		if shields_up != FALSE {
			if act.unit_type[defender_index] == SHIP {
				damage_to_shields = (defending_ship.dest_y * damage_done) / 100
				damage_to_ship = damage_done - damage_to_shields
				act.shield_strength_left[defender_index] -= damage_to_shields

				/* Calculate percentage of shields left. */
				if act.shield_strength_left[defender_index] > 0 {
					defenderShieldStrength := act.shield_strength[defender_index]
					if defenderShieldStrength <= 0 {
						defenderShieldStrength = 1
					}
					defending_ship.dest_y = (100 * act.shield_strength_left[defender_index]) / defenderShieldStrength
				} else {
					defending_ship.dest_y = 0
				}
			} else { /* Planetary defenses. */
				act.shield_strength_left[defender_index] -= damage_done
			}
		}

		/* See if it got through shields. */
		units_destroyed = 0
		percent_decrease = 0
		if shields_up == FALSE || act.shield_strength_left[defender_index] < 0 || damage_to_ship > 0 {
			/* Get net damage to ship or PDs. */
			if shields_up != FALSE {
				if act.unit_type[defender_index] == SHIP {
					/* Total damage to ship is direct damage plus damage that shields could not absorb. */
					damage_done = damage_to_ship
					if act.shield_strength_left[defender_index] < 0 {
						damage_done -= act.shield_strength_left[defender_index]
					}
				} else {
					damage_done = -act.shield_strength_left[defender_index]
				}
			}

			defenderShieldStrength := act.shield_strength[defender_index]
			if defenderShieldStrength <= 0 {
				defenderShieldStrength = 1
			}

			percent_decrease = (50 * damage_done) / defenderShieldStrength

			percent_decrease += ((e.rnd(51) - 26) * percent_decrease) / 100
			if percent_decrease > 100 {
				percent_decrease = 100
			}

			if act.unit_type[defender_index] == SHIP {
				defending_ship.age += percent_decrease / 2
				if defending_ship.age > 49 {
					units_destroyed = 1
				}
			} else {
				units_destroyed = (percent_decrease * act.original_age_or_PDs[defender_index]) / 100
				if units_destroyed > defending_nampla.item_quantity[PD] {
					units_destroyed = defending_nampla.item_quantity[PD]
				}
				if units_destroyed < 1 {
					units_destroyed = 1
				}
				defending_nampla.item_quantity[PD] -= units_destroyed
			}

			if act.shield_strength_left[defender_index] < 0 {
				act.shield_strength_left[defender_index] = 0
			}
		}

		/* See if this is a hijacking. */
		i = act.fighting_species_index[attacker_index]
		j = act.fighting_species_index[defender_index]
		if bat.enemy_mine[i][j] == 2 && (option == DEEP_SPACE_FIGHT || option == PLANET_ATTACK) {
			this_is_a_hijacking = TRUE
		} else {
			this_is_a_hijacking = FALSE
		}

		attacking_species = e.c_species[i]
		defending_species = e.c_species[j]

		/* Report if anything was destroyed. */
		FDs_were_destroyed = FALSE
		if units_destroyed != 0 {
			if act.unit_type[defender_index] == SHIP {
				e.log_summary = TRUE
				e.log_string("        ")
				e.log_string(defender_name)
				if this_is_a_hijacking != FALSE {
					e.log_string(" was successfully hijacked and will generate ")

					if defending_ship.class == TR || defending_ship._type == STARBASE {
						original_cost = ship_cost[defending_ship.class] * defending_ship.tonnage
					} else {
						original_cost = ship_cost[defending_ship.class]
					}

					if defending_ship._type == SUB_LIGHT {
						original_cost = (3 * original_cost) / 4
					}

					if defending_ship.status == UNDER_CONSTRUCTION {
						recycle_value = (original_cost - defending_ship.remaining_cost) / 2
					} else {
						recycle_value = (3 * original_cost * (60 - act.original_age_or_PDs[defender_index])) / 200
					}

					economic_units = recycle_value

					for i = 0; i < MAX_ITEMS; i++ {
						j = defending_ship.item_quantity[i]
						if j > 0 {
							if i == TP {
								techLevel_2x := 2 * defending_species.tech_level[BI]
								if techLevel_2x <= 0 {
									techLevel_2x = 1
								}
								recycle_value = (j * item_cost[i]) / techLevel_2x
							} else if i == RM {
								recycle_value = j / 5
							} else {
								recycle_value = (j * item_cost[i]) / 2
							}

							economic_units += recycle_value
						}
					}

					attacking_species.econ_units += economic_units

					e.log_long(economic_units)
					e.log_string(" economic units for the hijackers.\n")
				} else {
					e.log_string(" was destroyed.\n")
				}

				for i = 0; i < MAX_ITEMS; i++ {
					if defending_ship.item_quantity[i] > 0 {
						/* If this is a hijacking of a field-distorted ship,
						   we want the true name of the hijacked species to
						   be announced, but we don't want any cargo to be
						   destroyed. */
						if i == FD {
							FDs_were_destroyed = TRUE
						}
						if this_is_a_hijacking == FALSE {
							defending_ship.item_quantity[FD] = 0
						}
					}
				}
				e.log_to_file = FALSE
				if this_is_a_hijacking != FALSE {
					e.log_string("          The hijacker was ")
				} else {
					e.log_string("          The killing blow was delivered by ")
				}
				e.log_string(attacker_name)
				e.log_string(".\n")
				e.log_to_file = TRUE
				e.log_summary = FALSE

				total_shots -= act.shots_left[defender_index]
				act.shots_left[defender_index] = 0
				act.num_shots[defender_index] = 0
			} else {
				e.log_summary = TRUE
				e.log_string("        ")
				e.log_int(units_destroyed)
				if units_destroyed > 1 {
					e.log_string(" PDs on PL ")
				} else {
					e.log_string(" PD on PL ")
				}
				e.log_string(defending_nampla.name)
				if units_destroyed > 1 {
					e.log_string(" were destroyed by ")
				} else {
					e.log_string(" was destroyed by ")
				}

				e.log_string(attacker_name)
				e.log_string(".\n")

				if defending_nampla.item_quantity[PD] == 0 {
					total_shots -= act.shots_left[defender_index]
					act.shots_left[defender_index] = 0
					act.num_shots[defender_index] = 0
					e.log_string("        All planetary defenses have been destroyed on ")
					e.log_string(defender_name)
					e.log_string("!\n")
				}
				e.log_summary = FALSE
			}
		} else if percent_decrease > 0 && this_is_a_hijacking == FALSE && act.unit_type[defender_index] == SHIP {
			/* See if anything carried by the ship was also destroyed. */
			for i = 0; i < MAX_ITEMS; i++ {
				j = defending_ship.item_quantity[i]
				if j > 0 {
					j = (percent_decrease * j) / 100
					if j > 0 {
						defending_ship.item_quantity[i] -= j
						if i == FD {
							FDs_were_destroyed = TRUE
						}
					}
				}
			}
		}

		j = act.fighting_species_index[defender_index]
		if FDs_were_destroyed != FALSE && e.field_distorted[j] != FALSE && defending_ship.dest_x == 0 {
			/* Reveal the true name of the ship and the owning species. */
			e.log_summary = TRUE
			if this_is_a_hijacking != FALSE {
				e.log_string("        Hijacking of ")
			} else {
				e.log_string("        Damage to ")
			}
			e.log_string(defender_name)
			e.log_string(" caused collapse of distortion field. Real name of ship is ")
			e.log_string(e.ship_name(defending_ship))
			e.log_string(" owned by SP ")
			e.log_string(defending_species.name)
			e.log_string(".\n")
			e.log_summary = FALSE
			defending_ship.dest_x = 127 /* Ship is now exposed. */
		}
	}

	/* No more surprises. */
	for i = 0; i < act.num_units_fighting; i++ {
		act.surprised[i] = FALSE
	}

	return combat_occurred
}
