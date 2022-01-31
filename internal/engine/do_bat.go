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
	"log"
)

// auto_enemy will find all species that have declared alliance with both a traitor and betrayed species.
// it will then set a flag to indicate that their allegiance should be changed from ALLY to ENEMY.
//
// traitorSpeciesNo  is one-based species_number
// betrayedSpeciesNo is one-based species_number
func (e *Engine) auto_enemy(traitor_species_no, betrayed_species_no int) {
	for species_index := 0; species_index < e.galaxy.num_species; species_index++ {
		if e.spec_data[species_index].ally[traitor_species_no-1] == FALSE {
			continue
		} else if e.spec_data[species_index].ally[betrayed_species_no-1] == FALSE {
			continue
		} else if e.spec_data[species_index].contact[traitor_species_no-1] == FALSE {
			continue
		} else if e.spec_data[species_index].contact[betrayed_species_no-1] == FALSE {
			continue
		}
		e.make_enemy[species_index][traitor_species_no-1] = betrayed_species_no
	}
}

func (e *Engine) do_ambush(ambushing_species_index int, bat *battle_data) {
	// get total ambushing tonnage
	num_ships, friendly_tonnage := e.c_species[ambushing_species_index].num_ships, 0
	for i := 0; i < num_ships; i++ {
		sh := e.c_ship[ambushing_species_index][i]
		if sh.pn == 99 || sh.x != bat.x || sh.y != bat.y || sh.z != bat.z {
			continue
		}
		if sh.class != TR && sh.class != BA {
			friendly_tonnage += sh.tonnage
		}
	}

	// determine which species are being ambushed and get total enemy tonnage
	num_sp, enemy_tonnage := bat.num_species_here, 0
	for ambushed_species_index := 0; ambushed_species_index < num_sp; ambushed_species_index++ {
		if bat.enemy_mine[ambushing_species_index][ambushed_species_index] == 0 {
			continue
		}

		// this species is being ambushed; get total effective tonnage
		num_ships = e.c_species[ambushed_species_index].num_ships
		for i := 0; i < num_ships; i++ {
			sh := e.c_ship[ambushed_species_index][i]
			if sh.pn == 99 || sh.x != bat.x || sh.y != bat.y || sh.z != bat.z {
				continue
			}
			if sh.class == TR {
				enemy_tonnage += sh.tonnage
			} else {
				enemy_tonnage += 10 * sh.tonnage
			}
		}
	}

	// determine the amount of aging that will be added to each ambushed ship
	if enemy_tonnage == 0 {
		return
	}
	e.ambush_took_place = TRUE
	age_increment := (10 * bat.ambush_amount[ambushing_species_index]) / enemy_tonnage
	age_increment = (friendly_tonnage * age_increment) / enemy_tonnage
	if age_increment < 1 {
		e.log_printf("\n    SP %s attempted an ambush, but the ambush was completely ineffective!\n", e.c_species[ambushing_species_index].name)
		return
	}

	// age each ambushed ship
	for ambushed_species_index := 0; ambushed_species_index < num_sp; ambushed_species_index++ {
		if bat.enemy_mine[ambushing_species_index][ambushed_species_index] == 0 {
			continue
		}
		e.log_string("\n    SP ")
		species_number := bat.spec_num[ambushed_species_index]
		if e.field_distorted[ambushed_species_index] != 0 {
			e.log_int(e.distorted(species_number))
		} else {
			e.log_string(e.c_species[ambushed_species_index].name)
		}
		e.log_printf(" was ambushed by SP %s!\n", e.c_species[ambushing_species_index].name)

		num_ships = e.c_species[ambushed_species_index].num_ships
		for i := 0; i < num_ships; i++ {
			sh := e.c_ship[ambushed_species_index][i]
			if sh.pn == 99 || sh.x != bat.x || sh.y != bat.y || sh.z != bat.z {
				continue
			}
			sh.age += age_increment
			if sh.arrived_via_wormhole != FALSE { // todo: why the wormhole penalty?
				sh.age += age_increment
			}
			if sh.age > 49 {
				old_truncate_name := e.truncate_name // stash the global value
				e.truncate_name = TRUE
				e.log_printf("      %s", e.ship_name(sh))
				if e.field_distorted[ambushed_species_index] != 0 {
					e.log_printf(" = %s ", e.c_species[ambushed_species_index].name)
					// stash the current value to trick the ship_name() function
					originalFDQuantity := sh.item_quantity[FD]
					sh.item_quantity[FD] = 0
					e.log_string(e.ship_name(sh))
					sh.item_quantity[FD] = originalFDQuantity // restore the quantity
				}
				n := 0
				for j := 0; j < MAX_ITEMS; j++ {
					if sh.item_quantity[j] > 0 {
						if n == 0 {
							e.log_string(" (cargo: ")
						} else {
							e.log_char(',')
						}
						n++
						e.log_int(sh.item_quantity[j])
						e.log_char(' ')
						e.log_string(item_abbr[j])
					}
				}
				if n > 0 {
					e.log_char(')')
				}
				e.log_string(" was destroyed in the ambush!\n")
				e.truncate_name = old_truncate_name // restore the global value
			}
		}
	}

}

func (e *Engine) do_battle(bat *battle_data) {
	const (
		canBeSpeciallySurprised = 55 // todo: what does this mean?
		hijackingAttack         = 2
		normalAttack            = 1
	)
	var (
		current_species      int
		identifiable_units   [MAX_SPECIES]int
		namp                 *nampla_data
		need_comma           int // TRUE or FALSE
		sh                   *ship_data
		species_number       int
		unidentifiable_units [MAX_SPECIES]int
	)
	act := &action_data{}
	e.ambush_took_place = FALSE

	// open log file for writing
	if e.combat_log = fopen("combat.log", &bytes.Buffer{}); e.combat_log == nil {
		fprintf(e.stderr, "\n\tCannot open 'combat.log' for writing!\n\n")
		panic("\n\tCannot open 'combat.log' for writing!\n\n")
	}
	e.log_file = e.combat_log

	// open summary file for writing
	if e.summary_file = fopen("summary.log", &bytes.Buffer{}); e.summary_file == nil {
		fprintf(e.stderr, "\n\tCannot open 'summary.log' for writing!\n\n")
		panic("\n\tCannot open 'summary.log' for writing!\n\n")
	}
	e.log_summary = TRUE

	// get data for all species present at this battle
	num_sp := bat.num_species_here
	for species_index := 0; species_index < num_sp; species_index++ {
		species_number := bat.spec_num[species_index]
		e.c_species[species_index] = e.spec_data[species_number-1]
		e.c_nampla[species_index] = e.namp_data[species_number-1]
		e.c_ship[species_index] = e.ship_data[species_number-1]

		// determine number of identifiable and unidentifiable units present
		identifiable_units[species_index] = 0
		unidentifiable_units[species_index] = 0

		namp = nil
		for i := 0; i < e.c_species[species_index].num_namplas; i++ {
			namp = e.c_nampla[species_index][i]
			if namp.x != bat.x || namp.y != bat.y || namp.z != bat.z {
				continue
			}

			if (namp.status & POPULATED) != 0 {
				identifiable_units[species_index]++
			}
		}

		sh = nil
		for i := 0; i < e.c_species[species_index].num_ships; i++ {
			sh = e.c_ship[species_index][i]
			if sh.x != bat.x || sh.y != bat.y || sh.z != bat.z {
				continue
			} else if sh.status == UNDER_CONSTRUCTION {
				continue
			} else if sh.status == JUMPED_IN_COMBAT {
				continue
			} else if sh.status == FORCED_JUMP {
				continue
			}

			// todo: overloading dest_x and dest_y
			sh.dest_x = 0   // not yet exposed
			sh.dest_y = 100 // shields at 100%

			if sh.item_quantity[FD] == sh.tonnage {
				unidentifiable_units[species_index]++
			} else {
				identifiable_units[species_index]++
			}
		}

		// if any unit is not distorted, all units are not distorted
		e.field_distorted[species_index] = FALSE
		if identifiable_units[species_index] == 0 && unidentifiable_units[species_index] > 0 {
			e.field_distorted[species_index] = TRUE
		}
	}

	/* Start log of what's happening. */
	if e.strike_phase != FALSE {
		e.log_string("\nStrike log:\n")
	} else {
		e.log_string("\nCombat log:\n")
	}
	e.first_battle = FALSE

	e.log_printf("\n  Battle orders were received for sector %d, %d, %d. The following species are present:\n\n", bat.x, bat.y, bat.z)

	// Convert enemy_mine array from a list of species numbers to an array of TRUE/FALSE values whose indices are:
	//
	//		[species_index1][species_index2]
	//
	// such that the value will be TRUE if #1 mentioned #2 in an ATTACK or HIJACK command.
	// The actual TRUE value will be 1 for ATTACK or 2 for HIJACK.

	for species_index := 0; species_index < num_sp; species_index++ {
		/* Make copy of list of enemies. */
		var enemy_num [MAX_SPECIES]int
		for i := 0; i < MAX_SPECIES; i++ {
			enemy_num[i] = bat.enemy_mine[species_index][i]
		}
		// and then reset the enemy_mine matrix
		for i := 0; i < MAX_SPECIES; i++ {
			bat.enemy_mine[species_index][i] = FALSE
		}

		for i := 0; i < MAX_SPECIES; i++ {
			enemy := enemy_num[i]
			if enemy == 0 { // no more enemies in list
				break
			}
			var attackType int // kind of overrides the true/false on the enemy_mine array :(
			if enemy < 0 {
				enemy = -enemy
				attackType = hijackingAttack
			} else {
				attackType = normalAttack
			}

			// convert absolute species numbers to species indices that have been assigned in the current battle.
			// also set the attack type (either normal or hijacking attempt)
			for j := 0; j < num_sp; j++ {
				if enemy == bat.spec_num[j] {
					bat.enemy_mine[species_index][j] = attackType
				}
			}
		}
	}

	// for each species that has been mentioned in an attack order, check if it can be surprised.
	// a species can only be surprised if it has not given a BATTLE order and if it is being attacked ONLY by one or more ALLIES.
	for species_index := 0; species_index < num_sp; species_index++ {
		j := bat.spec_num[species_index] - 1
		log.Printf("[do_bat] bat.spec_num[species_index %d] is %d\n", species_index, j)
		for i := 0; i < num_sp; i++ {
			if i == species_index || bat.enemy_mine[species_index][i] == FALSE {
				continue
			} else if e.field_distorted[species_index] != FALSE {
				// attacker is field-distorted; surprise not possible
				bat.can_be_surprised[i] = FALSE
				continue
			}

			betrayal := FALSE
			if e.c_species[i].ally[j] != FALSE { // todo: is this the right index?
				betrayal = TRUE
			}

			if betrayal != FALSE { // someone is being attacked by an ALLY
				traitor_number, betrayed_number := bat.spec_num[species_index], bat.spec_num[i]
				e.make_enemy[betrayed_number-1][traitor_number-1] = betrayed_number
				e.auto_enemy(traitor_number, betrayed_number)
			}

			if bat.can_be_surprised[i] == FALSE {
				continue
			}

			if betrayal == FALSE { /* At least one attacker is not an ally. */
				bat.can_be_surprised[i] = FALSE
			}
		}
	}

	// for each species that has been mentioned in an attack order, see
	// if there are other species present that have declared it as an ALLY.
	// If so, have the attacker attack the other species and vice-versa.
	for species_index := 0; species_index < num_sp; species_index++ {
		for i := 0; i < num_sp; i++ {
			if i == species_index || bat.enemy_mine[species_index][i] == FALSE {
				continue
			}
			j := bat.spec_num[i] - 1 // todo: what is `j`?
			log.Printf("[do_bat] bat.spec_num[%d] is %d\n", i, j)
			for k := 0; k < num_sp; k++ {
				if k == species_index || k == i {
					continue
				}
				if e.c_species[k].ally[j] != FALSE { // todo: is `j` the right index?
					// make sure it's not already set (it may already be set for HIJACK and we don't want to accidentally change it to ATTACK)
					if bat.enemy_mine[species_index][k] == FALSE {
						bat.enemy_mine[species_index][k] = TRUE
					}
					if bat.enemy_mine[k][species_index] == FALSE {
						bat.enemy_mine[k][species_index] = TRUE
					}
				}
			}
		}
	}

	/* If a species did not give a battle order and is not the target of an
	attack, set can_be_surprised flag to a special value. */
	for species_index := 0; species_index < num_sp; species_index++ {
		if bat.can_be_surprised[species_index] == FALSE {
			continue
		}
		bat.can_be_surprised[species_index] = canBeSpeciallySurprised
		for i := 0; i < num_sp; i++ {
			if i == species_index || bat.enemy_mine[i][species_index] == FALSE {
				continue
			}
			bat.can_be_surprised[species_index] = TRUE
			break
		}
	}

	// list combatants
	for species_index := 0; species_index < num_sp; species_index++ {
		species_number = bat.spec_num[species_index]
		e.log_printf("    SP%02d SP ", species_number)
		if e.field_distorted[species_index] == FALSE {
			e.log_string(e.c_species[species_index].name)
		} else {
			e.log_int(e.distorted(species_number))
		}
		if bat.can_be_surprised[species_index] == FALSE {
			e.log_string(" is mobilized and")
		} else {
			e.log_string(" does not appear to be")
		}
		e.log_string(" ready for combat.\n")
	}

	// check if a declared enemy is being ambushed
	for i := 0; i < num_sp; i++ {
		namp = nil
		bat.ambush_amount[i] = 0
		for j := 0; j < e.c_species[i].num_namplas; j++ {
			namp = e.c_nampla[i][j]
			if namp.x != bat.x || namp.y != bat.y || namp.z != bat.z {
				continue
			}
			bat.ambush_amount[i] += namp.use_on_ambush
		}
		if bat.ambush_amount[i] == 0 {
			continue
		}
		for j := 0; j < num_sp; j++ {
			if bat.enemy_mine[i][j] != FALSE {
				e.do_ambush(i, bat)
			}
		}
	}

	// for all species that specified enemies, make the feeling mutual
	for i := 0; i < num_sp; i++ {
		for j := 0; j < num_sp; j++ {
			if bat.enemy_mine[i][j] != FALSE {
				// make sure it's not already set (it may already be set for HIJACK and we don't want to accidentally change it to ATTACK)
				if bat.enemy_mine[j][i] == FALSE {
					bat.enemy_mine[j][i] = normalAttack
				}
			}
		}
	}

	// create a sequential list of combat options.
	// first check if a deep space defense has been ordered.
	// if so, make sure that first option is DEEP_SPACE_FIGHT.
	e.num_combat_options = 0
	for species_index := 0; species_index < num_sp; species_index++ {
		for i := 0; i < bat.num_engage_options[species_index]; i++ {
			option := bat.engage_option[species_index][i]
			if option == DEEP_SPACE_DEFENSE {
				e.consolidate_option(DEEP_SPACE_FIGHT, 0)
				goto consolidate
			}
		}
	}
consolidate:
	for species_index := 0; species_index < num_sp; species_index++ {
		for i := 0; i < bat.num_engage_options[species_index]; i++ {
			e.consolidate_option(bat.engage_option[species_index][i], bat.engage_planet[species_index][i])
		}
	}

	// if ships are given unconditional withdraw orders, they will always have
	// time to escape if fighting occurs first in a different part of the sector.
	// the flag "do_withdraw_check_first" will be set only after the first round of combat.
	do_withdraw_check_first := FALSE

	// handle each combat option
	battle_here, first_action := FALSE, TRUE
	for option_index := 0; option_index < e.num_combat_options; option_index++ {
		option := e.combat_option[option_index]
		where := e.combat_location[option_index]

		// fill action arrays with data about ships taking part in current action
		fight_here := e.fighting_params(option, where, bat, act)

		/* Check if a fight will take place here. */
		if fight_here == FALSE {
			continue
		}

		/* See if anyone is taken by surprise. */
		if battle_here == FALSE {
			/* Combat is just starting. */
			for species_index := 0; species_index < num_sp; species_index++ {
				species_number = bat.spec_num[species_index]

				if bat.can_be_surprised[species_index] == canBeSpeciallySurprised {
					continue
				}

				if bat.can_be_surprised[species_index] != FALSE {
					e.log_string("\n    SP ")
					if e.field_distorted[species_index] != FALSE {
						e.log_int(e.distorted(species_number))
					} else {
						e.log_string(e.c_species[species_index].name)
					}
					e.log_string(" is taken by surprise!\n")
				}
			}
		}

		battle_here = TRUE

		/* Clear out can_be_surprised array. */
		for i := 0; i < MAX_SPECIES; i++ {
			bat.can_be_surprised[i] = FALSE
		}

		/* Determine maximum number of rounds. */
		max_rounds := 10000 /* Something ridiculously large. */
		if option == DEEP_SPACE_FIGHT && e.attacking_ML > 0 && e.defending_ML > 0 && e.deep_space_defense != FALSE {
			/* This is the initial deep space fight and the defender wants the fight to remain in deep space for as long as possible. */
			max_rounds = e.defending_ML - e.attacking_ML
			if max_rounds < 1 {
				max_rounds = 1
			}
		} else if option == PLANET_BOMBARDMENT {
			/* To determine the effectiveness of the bombardment, we will simulate ten rounds of combat and add up the damage. */
			max_rounds = 10
		} else if option == GERM_WARFARE || option == SIEGE {
			/* We just need to see who is attacking whom and get the number of germ warfare bombs being used. */
			max_rounds = 1
		}

		/* Log start of action. */
		if where == 0 {
			e.log_string("\n    The battle begins in deep space, outside the range of planetary defenses...\n")
		} else if option == PLANET_ATTACK {
			e.log_string("\n    The battle ")
			if first_action != FALSE {
				e.log_string("begins")
			} else {
				e.log_string("moves")
			}
			e.log_printf(" within range of planet #%d...\n", where)
		} else if option == PLANET_BOMBARDMENT {
			e.log_printf("\n    Bombardment of planet #%d begins...\n", where)
		} else if option == GERM_WARFARE {
			e.log_printf("\n    Germ warfare commences against planet #%d...\n", where)
		} else if option == SIEGE {
			e.log_printf("\n    Siege of planet #%d is now in effect...\n\n", where)
			goto do_combat
		}

		/* List combatants. */
		e.truncate_name = FALSE
		e.log_string("\n      Units present:")
		current_species = -1
		for unit_index := 0; unit_index < act.num_units_fighting; unit_index++ {
			if act.fighting_species_index[unit_index] != current_species {
				/* Display species name. */
				i := act.fighting_species_index[unit_index]
				e.log_string("\n        SP ")
				species_number = bat.spec_num[i]
				if e.field_distorted[i] != FALSE {
					e.log_int(e.distorted(species_number))
				} else {
					e.log_string((e.c_species[i].name))
				}
				e.log_string(": ")
				current_species = i
				need_comma = FALSE
			}

			var ok bool
			if act.unit_type[unit_index] == SHIP {
				if sh, ok = act.fighting_unit[unit_index].(*ship_data); !ok {
					panic("act.fighting_unit[unit_index].(*ship_data); !ok")
				}
				temp_status, temp_pn := sh.status, sh.pn
				if option == DEEP_SPACE_FIGHT {
					sh.status, sh.pn = IN_DEEP_SPACE, 0
				} else {
					sh.status, sh.pn = IN_ORBIT, where
				}
				if e.field_distorted[current_species] == FALSE {
					e.ignore_field_distorters = TRUE
				} else {
					e.ignore_field_distorters = FALSE
				}
				if sh.special != NON_COMBATANT {
					if need_comma != FALSE {
						e.log_string(", ")
					}
					e.log_string((e.ship_name(sh)))
					need_comma = TRUE
				}
				e.ignore_field_distorters = FALSE
				sh.status, sh.pn = temp_status, temp_pn
			} else {
				if namp, ok = act.fighting_unit[unit_index].(*nampla_data); !ok {
					panic("act.fighting_unit[unit_index].(*nampla_data); !ok")
				}
				if need_comma != FALSE {
					e.log_string(", ")
				}
				e.log_printf("PL %s", namp.name)
				need_comma = TRUE
			}
		}
		e.log_string("\n\n")

	do_combat:

		/* Long names are not necessary for the rest of the action. */
		e.truncate_name = TRUE

		/* Do combat rounds. Stop if maximum count is reached, or if combat does not occur when do_round() is called. */
		round_number := 1
		e.log_summary = FALSE /* do_round() and the routines that it calls will set this for important stuff. */

		if option == PLANET_BOMBARDMENT || option == GERM_WARFARE || option == SIEGE {
			/* Disable logging during simulation. */
			e.logging_disabled = TRUE
		}

		for round_number <= max_rounds {
			if do_withdraw_check_first != FALSE {
				e.withdrawal_check(bat, act)
			}
			if e.do_round(option, round_number, bat, act) == FALSE {
				break
			}
			if do_withdraw_check_first == FALSE {
				e.withdrawal_check(bat, act)
			}
			do_withdraw_check_first = TRUE
			e.regenerate_shields(act)
			round_number++
		}

		e.log_summary = TRUE
		e.logging_disabled = FALSE

		if round_number == 1 {
			e.log_string("      ...But it seems that the attackers had nothing to attack!\n")
			continue
		}

		if option == PLANET_BOMBARDMENT || option == GERM_WARFARE {
			for unit_index := 0; unit_index < act.num_units_fighting; unit_index++ {
				if act.unit_type[unit_index] == GENOCIDE_NAMPLA {
					attacked_nampla, ok := act.fighting_unit[unit_index].(*nampla_data)
					if !ok {
						panic("act.fighting_unit[unit_index].(*nampla_data); !ok")
					}
					j := act.fighting_species_index[unit_index]
					for i := 0; i < num_sp; i++ {
						if e.x_attacked_y[i][j] != FALSE {
							species_number = bat.spec_num[i]
							e.log_string("      SP ")
							if e.field_distorted[i] != FALSE {
								e.log_int(e.distorted(species_number))
							} else {
								e.log_string(e.c_species[i].name)
							}
							e.log_printf(" bombards SP %s on PL %s.\n", e.c_species[j].name, attacked_nampla.name)

							if option == GERM_WARFARE {
								e.do_germ_warfare(i, j, unit_index, bat, act)
							}
						}
					}

					/* Determine results of bombardment. */
					if option == PLANET_BOMBARDMENT {
						e.do_bombardment(unit_index, act)
					}
				}
			}
		} else if option == SIEGE {
			e.do_siege(bat, act)
		}

		e.truncate_name = FALSE
		first_action = FALSE
	}

	if battle_here == FALSE {
		if bat.num_species_here == 1 {
			e.log_string("    But there was no one to fight with!\n")
		} else if e.ambush_took_place == FALSE {
			e.log_string("    But no one was willing to throw the first punch!\n")
		}
	}

	// close combat log and append it to the log files of all species involved in this battle.
	if e.prompt_gm {
		log.Printf("  End of battle in sector %d, %d, %d.\n", bat.x, bat.y, bat.z)
	}
	fprintf(e.log_file, "\n  End of battle in sector %d, %d, %d.\n", bat.x, bat.y, bat.z)
	fprintf(e.summary_file, "\n  End of battle in sector %d, %d, %d.\n", bat.x, bat.y, bat.z)
	fclose(e.log_file)
	fclose(e.summary_file)

	for species_index := 0; species_index < num_sp; species_index++ {
		species_number = bat.spec_num[species_index]

		// todo: fix this
		//// open combat or summary log file for reading
		//var logToCopyFrom *FILE
		//if bat.summary_only[species_index] == FALSE {
		//	if e.combat_log == nil {
		//		fprintf(e.stderr, "\n\tCannot open combat log for reading!\n\n")
		//		panic("\n\tCannot open combat log for reading!\n\n")
		//	}
		//	logToCopyFrom = e.combat_log // fopen("combat.log", "r")
		//} else {
		//	if e.summary_file == nil {
		//		fprintf(e.stderr, "\n\tCannot open combat summary log for reading!\n\n")
		//		panic("\n\tCannot open combat summary log for reading!\n\n")
		//	}
		//	logToCopyFrom = e.summary_file // fopen("summary.log", "r")
		//}
		if e.summary_file != nil {
			e.spec_logs[species_number-1].Write(e.summary_file.bytes())
			log.Printf("[engine] do_bat SP%02d combat summary log\n%s\n", species_number, string(e.summary_file.bytes()))
		}
		if e.combat_log != nil {
			log.Printf("[engine] do_bat SP%02d combat log\n%s\n", species_number, string(e.combat_log.bytes()))
			e.spec_logs[species_number-1].Write(e.combat_log.bytes())
		}

		//// open a temporary species log file for appending
		//species_log := e.spec_logs[species_number-1] // zero based?
		//if species_log == nil {
		//	fprintf(e.stderr, "\n\tCannot open SP%02d log for appending!\n\n", species_number)
		//	panic(fmt.Sprintf("\n\tCannot open SP%02d log for appending!\n\n", species_number))
		//}
		//
		//// copy combat log to temporary species log
		//species_log.Write(logToCopyFrom.bytes())

		e.append_log[species_number-1] = TRUE

		// delete the ships that were destroyed
		for i := 0; i < e.c_species[species_index].num_ships; i++ {
			sh := e.c_ship[species_index][i]
			if sh.age < 50 {
				continue
			} else if sh.pn == 99 {
				continue
			} else if sh.x != bat.x {
				continue
			} else if sh.y != bat.y {
				continue
			} else if sh.z != bat.z {
				continue
			} else if sh.status == UNDER_CONSTRUCTION {
				continue
			}
			e.delete_ship(sh)
		}
	}
}
