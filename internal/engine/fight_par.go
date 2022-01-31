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

import "log"

func (e *Engine) disbanded_ship(species_index int, sh *ship_data) int {
	for nampla_index := 0; nampla_index < e.c_species[species_index].num_namplas; nampla_index++ {
		nam := e.c_nampla[species_index][nampla_index]
		if nam.x != sh.x || nam.y != sh.y || nam.z != sh.z || nam.pn != sh.pn {
			continue
		} else if (nam.status & DISBANDED_COLONY) == 0 {
			continue
		} else if sh._type != STARBASE && sh.status == IN_ORBIT {
			continue
		}
		/* This ship is either on the surface of a disbanded colony or is a starbase orbiting a disbanded colony. */
		return TRUE
	}

	return FALSE
}

func (e *Engine) fighting_params(option, location int, bat *battle_data, act *action_data) int {
	var (
		found, _type, unit_index                     int
		nampla_index, sp1, sp2, use_this_ship        int
		engage_option, engage_location               int
		tons                                         int
		unit_power, offensive_power, defensive_power int
		sh                                           *ship_data
		nam                                          *nampla_data
	)

	// add fighting units to "act" arrays.
	// at the same time, check if a fight of the current option type will occur at the current location.
	e.attacking_ML, e.defending_ML, e.deep_space_defense = 0, 0, FALSE
	x, y, z := bat.x, bat.y, bat.z
	attacking_ships_here, defending_ships_here := FALSE, FALSE
	attacking_pds_here, defending_pds_here := FALSE, FALSE
	num_sp, num_fighting_units := bat.num_species_here, 0

	for species_index := 0; species_index < num_sp; species_index++ {
		// check which ships can take part in fight
		for ship_index := 0; ship_index < e.c_species[species_index].num_ships; ship_index++ {
			sh := e.c_ship[species_index][ship_index]
			use_this_ship = FALSE

			if sh.pn == 99 || sh.x != x || sh.y != y || sh.z != z {
				continue
			} else if sh.age > 49 {
				continue
			} else if sh.status == UNDER_CONSTRUCTION {
				continue
			} else if sh.status == FORCED_JUMP {
				continue
			} else if sh.status == JUMPED_IN_COMBAT {
				continue
			} else if sh.class == TR && sh.pn != location && option != GERM_WARFARE {
				continue
			} else if e.disbanded_ship(species_index, sh) != 0 {
				continue
			} else if option == SIEGE || option == PLANET_BOMBARDMENT {
				if sh.special == NON_COMBATANT {
					continue
				}
			}

			for i := 0; i < bat.num_engage_options[species_index]; i++ {
				engage_option = bat.engage_option[species_index][i]
				engage_location = bat.engage_planet[species_index][i]

				switch engage_option {
				case DEFENSE_IN_PLACE:
					if sh.pn != location {
						break
					}
					defending_ships_here = TRUE
					use_this_ship = TRUE

				case DEEP_SPACE_DEFENSE:
					if option != DEEP_SPACE_FIGHT {
						break
					}
					if sh.class == BA && sh.pn != 0 {
						break
					}
					defending_ships_here = TRUE
					use_this_ship = TRUE
					e.deep_space_defense = TRUE
					if e.c_species[species_index].tech_level[ML] > e.defending_ML {
						e.defending_ML = e.c_species[species_index].tech_level[ML]
					}

				case PLANET_DEFENSE:
					if location != engage_location {
						break
					}
					if sh.class == BA && sh.pn != location {
						break
					}
					defending_ships_here = TRUE
					use_this_ship = TRUE

				case DEEP_SPACE_FIGHT:
					if option != DEEP_SPACE_FIGHT {
						break
					}
					if sh.class == BA && sh.pn != 0 {
						break
					}
					if e.c_species[species_index].tech_level[ML] > e.defending_ML {
						e.defending_ML = e.c_species[species_index].tech_level[ML]
					}
					defending_ships_here = TRUE
					attacking_ships_here = TRUE
					use_this_ship = TRUE

				case PLANET_ATTACK, PLANET_BOMBARDMENT, GERM_WARFARE, SIEGE:
					if sh.class == BA && sh.pn != location {
						break
					}
					if sh.class == TR && option == SIEGE {
						break
					}
					if option == DEEP_SPACE_FIGHT {
						// There are two possibilities here:
						//    1. outsiders are attacking locals,
						// or 2. locals are attacking locals.
						// If (1), we want outsiders to first fight in deep space.
						// If (2), locals will not first fight in deep space (unless other explicit orders were given).
						// The case is (2) if current species has a planet here.
						found = FALSE
						for nampla_index = 0; nampla_index < e.c_species[species_index].num_namplas; nampla_index++ {
							nam := e.c_nampla[species_index][nampla_index]

							if nam.x != x || nam.y != y || nam.z != z {
								continue
							}
							if (nam.status & POPULATED) == 0 {
								continue
							}
							found = TRUE
							break
						}

						if found == FALSE {
							attacking_ships_here = TRUE
							use_this_ship = TRUE
							if e.c_species[species_index].tech_level[ML] > e.attacking_ML {
								e.attacking_ML = e.c_species[species_index].tech_level[ML]
							}
							break
						}
					}
					if option != engage_option && option != PLANET_ATTACK {
						break
					}
					if location != engage_location {
						break
					}
					attacking_ships_here = TRUE
					use_this_ship = TRUE

				default:
					panic("\n\n\tInternal error #1 in fight_par.c - invalid engage option!\n\n")
				}
			}

			// add_ship:

			if use_this_ship != FALSE {
				// add data for this ship to action array.
				act.fighting_species_index[num_fighting_units] = species_index
				act.unit_type[num_fighting_units] = SHIP
				act.fighting_unit[num_fighting_units] = sh
				act.original_age_or_PDs[num_fighting_units] = sh.age
				num_fighting_units++
			}
		}

		// check which namplas can take part in fight.
		for nampla_index = 0; nampla_index < e.c_species[species_index].num_namplas; nampla_index++ {
			nam := e.c_nampla[species_index][nampla_index]

			if nam.x != x || nam.y != y || nam.z != z {
				continue
			} else if nam.pn != location {
				continue
			} else if (nam.status & POPULATED) == 0 {
				continue
			} else if (nam.status & DISBANDED_COLONY) != 0 {
				continue
			}

			/* This planet has been targeted for some kind of attack. In
			   most cases, one species will attack a planet inhabited by
			   another species. However, it is also possible for two or
			   more species to have colonies on the SAME planet, and for
			   one to attack the other. */

			for i := 0; i < bat.num_engage_options[species_index]; i++ {
				engage_option = bat.engage_option[species_index][i]
				engage_location = bat.engage_planet[species_index][i]
				if engage_location != location {
					continue
				}

				switch engage_option {
				case DEFENSE_IN_PLACE, DEEP_SPACE_DEFENSE, PLANET_DEFENSE, DEEP_SPACE_FIGHT:
					break
				case PLANET_ATTACK, PLANET_BOMBARDMENT, GERM_WARFARE, SIEGE:
					if option != engage_option && option != PLANET_ATTACK {
						break
					}
					if nam.item_quantity[PD] > 0 {
						attacking_pds_here = TRUE
					}
					break

				default:
					panic("\n\n\tInternal error #2 in fight_par.c - invalid engage option!\n\n")
				}
			}

			if nam.item_quantity[PD] > 0 {
				defending_pds_here = TRUE
			}

			/* Add data for this nampla to action array. */
			act.fighting_species_index[num_fighting_units] = species_index
			act.unit_type[num_fighting_units] = NAMPLA
			act.fighting_unit[num_fighting_units] = nam
			act.original_age_or_PDs[num_fighting_units] = nam.item_quantity[PD]
			num_fighting_units++
		}
	}

	/* Depending on option, see if the right combination of combatants are present. */
	switch option {
	case DEEP_SPACE_FIGHT:
		if attacking_ships_here == 0 || defending_ships_here == 0 {
			return FALSE
		}
		break

	case PLANET_ATTACK, PLANET_BOMBARDMENT:
		if attacking_ships_here == 0 && attacking_pds_here == 0 {
			return FALSE
		}
		break

	case SIEGE, GERM_WARFARE:
		if attacking_ships_here == 0 {
			return FALSE
		}
		break

	default:
		panic("\n\n\tInternal error #3 in fight_par.c - invalid engage option!\n\n")
	}

	/* There is at least one attacker and one defender here. See if they are enemies. */
	for i := 0; i < num_fighting_units; i++ {
		sp1 = act.fighting_species_index[i]
		for j := 0; j < num_fighting_units; j++ {
			sp2 = act.fighting_species_index[j]
			if bat.enemy_mine[sp1][sp2] != 0 {
				goto next_step
			}
		}
	}

	return FALSE

next_step:

	act.num_units_fighting = num_fighting_units

	/* Determine number of shots, shield power and weapons power for all combatants. */
	for unit_index = 0; unit_index < act.num_units_fighting; unit_index++ {
		_type = act.unit_type[unit_index]
		if _type == SHIP {
			var ok bool
			if sh, ok = act.fighting_unit[unit_index].(*ship_data); !ok {
				panic("act.fighting_unit[unit_index].(*ship_data); !ok")
			}
			tons = sh.tonnage
		} else {
			var ok bool
			if nam, ok = act.fighting_unit[unit_index].(*nampla_data); !ok {
				panic("act.fighting_unit[unit_index].(*nampla_data); !ok")
			}
			tons = nam.item_quantity[PD] / 200
			if tons < 1 && nam.item_quantity[PD] > 0 {
				tons = 1
			}
		}

		species_index := act.fighting_species_index[unit_index]

		unit_power = power(tons)
		offensive_power = unit_power
		defensive_power = unit_power

		if _type == SHIP {
			if sh.class == TR {
				/* Transports are not designed for combat. */
				offensive_power /= 10
				defensive_power /= 10
			} else if sh.class != BA {
				/* Add auxiliary shield generator contribution, if any. */
				tons = 5
				for i := SG1; i <= SG9; i++ {
					if sh.item_quantity[i] > 0 {
						defensive_power += sh.item_quantity[i] * power(tons)
					}
					tons += 5
				}

				/* Add auxiliary gun unit contribution, if any. */
				tons = 5
				for i := GU1; i <= GU9; i++ {
					if sh.item_quantity[i] > 0 {
						offensive_power += sh.item_quantity[i] * power(tons)
					}
					tons += 5
				}
			}

			/* Adjust for ship aging. */
			offensive_power -= (sh.age * offensive_power) / 50
			defensive_power -= (sh.age * defensive_power) / 50
		}

		/* Adjust values for tech levels. */
		ml := e.c_species[species_index].tech_level[ML]
		ls := e.c_species[species_index].tech_level[LS]
		offensive_power += (ml * offensive_power) / 50
		defensive_power += (ls * defensive_power) / 50

		/* Adjust values if this species is hijacking anyone. */
		if bat.hijacker[species_index] != 0 && (option == DEEP_SPACE_FIGHT || option == PLANET_ATTACK) {
			offensive_power /= 4
			defensive_power /= 4
		}

		/* Get number of shots per round. */
		n_shots := (offensive_power / 1500) + 1
		if ml == 0 || offensive_power == 0 {
			n_shots = 0
		}
		if n_shots > 5 {
			n_shots = 5
		}
		act.num_shots[unit_index] = n_shots
		act.shots_left[unit_index] = n_shots

		/* Get damage per shot. */
		if n_shots > 0 {
			act.weapon_damage[unit_index] = (2 * offensive_power) / n_shots
		} else {
			act.weapon_damage[unit_index] = 0
		}

		/* Do defensive shields. */
		act.shield_strength[unit_index] = defensive_power
		if _type == SHIP {
			/* Adjust for results of previous action, if any. "dest_y"
			   contains the percentage of shields that remained at end
			   of last action. */
			defensive_power = (sh.dest_y * defensive_power) / 100
		}
		act.shield_strength_left[unit_index] = defensive_power

		/* Set bomb damage to zero in case this is planet bombardment or germ warfare. */
		act.bomb_damage[unit_index] = 0

		/* Set flag for individual unit if species can be surprised. */
		if bat.can_be_surprised[species_index] != 0 {
			act.surprised[unit_index] = TRUE
		} else {
			act.surprised[unit_index] = FALSE
		}
	}

	log.Printf("[fight_par] defending_pds_here %d: not used\n", defending_pds_here)
	return TRUE /* There will be a fight here. */
}
