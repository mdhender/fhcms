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
	"strconv"
	"strings"
)

/*
	Terminology: A "battle" consists of one or more "actions", and each
	action can take place at a different location in a star system.
	In turn, each action can be broken down into one or more "rounds",
	in which actual combat occurs.

	A battle is defined by a variable of type "struct battle_data",
	and a pointer to a variable of this type, called "bat", is used
	throughout the combat routines.

	An action is defined by a variable of type "struct action_data",
	and a pointer to a variable of this type, called "act", is used
	throughout the combat routines.
*/

func (e *Engine) combat(argv ...string) {
	e.locations_base = e.loc
	// e.prompt_gm = FALSE
	e.strike_phase = FALSE
	e.test_mode = FALSE
	e.verbose_mode = FALSE

	/* Allocate memory for battle data. */
	log.Printf("[engine] not making room for battle data\n")
	// e.battle_base = make([]*battle_data, MAX_BATTLES, MAX_BATTLES)

	for i := 0; i < MAX_SPECIES; i++ {
		if !(i < len(e.spec_data)) {
			continue
		}
		sp, spNo := e.spec_data[i], i+1
		log.Printf("[combat] %2d: SP%02d %-25s\n", i, spNo, sp.name)
		for alienIndex := 0; alienIndex < MAX_SPECIES; alienIndex++ {
			alienNo := alienIndex + 1
			if sp.contact[alienIndex] != FALSE {
				name := fmt.Sprintf("SP%02d", alienNo)
				if alienIndex < len(e.spec_data) {
					name += " " + e.spec_data[alienIndex].name
				}
				log.Printf("[combat] %2d: SP%02d %-25s  CONTACT %2d %s\n", i, spNo, sp.name, alienIndex, name)
			}
			if sp.ally[alienIndex] != FALSE {
				name := fmt.Sprintf("SP%02d", alienNo)
				if alienIndex < len(e.spec_data) {
					name += " " + e.spec_data[alienIndex].name
				}
				log.Printf("[combat] %2d: SP%02d %-25s  ALLY    %2d %s\n", i, spNo, sp.name, alienIndex, name)
			}
			if sp.enemy[alienIndex] != FALSE {
				name := fmt.Sprintf("SP%02d", alienNo)
				if alienIndex < len(e.spec_data) {
					name += " " + e.spec_data[alienIndex].name
				}
				log.Printf("[combat] %2d: SP%02d %-25s  ENEMY   %2d %s\n", i, spNo, sp.name, alienIndex, name)
			}
		}
	}

	var (
		arg_index          int
		bat                *battle_data
		battle_index       int
		best_species_index int
		best_score         int
		err                error
		i, j, k            int
		location           *sp_loc_data
		minimum_score      int
		n                  int
		name_length        int
		next_best_score    int
		num_battles        int
		num_enemies        int
		pl_num             [9]int // zero based index of orbits
		sp                 *species_data
		sp_index           int
		sp_num             [MAX_SPECIES]int
		sp_name            [MAX_SPECIES]string
		species_number     int
	)

	default_summary := FALSE
	do_all_species := FALSE
	num_species := 0

	// Check arguments.
	// If an argument is -s, then set SUMMARY mode for everyone.
	// The default is for players to receive a detailed report of the battles.
	// If an argument is -p, then prompt the GM before saving results;
	// otherwise, operate quietly; i.e, do not prompt GM before saving results
	// and do not display anything except errors.
	// Any additional arguments must be species numbers.
	// If no species numbers are specified, then do all species.

	if len(argv) > 0 && argv[0] == "Strike" {
		e.strike_phase = TRUE
	}

	argc := len(argv)
	for i = 1; i < argc; i++ {
		if argv[i] == "-s" {
			default_summary = TRUE
		} else if argv[i] == "-p" {
			e.prompt_gm = true
		} else if argv[i] == "-t" {
			e.test_mode = TRUE
		} else if argv[i] == "-v" {
			e.verbose_mode = TRUE
		} else if n, err = strconv.Atoi(argv[i]); err == nil && num_species < MAX_SPECIES && (0 <= n && n < MAX_SPECIES) {
			sp_num[num_species] = n
			num_species++
		}
	}

	if num_species == 0 {
		num_species = e.galaxy.num_species
		for i = 0; i < num_species; i++ {
			sp_num[i] = i + 1
		}
		do_all_species = TRUE
	}

	if default_summary != 0 && e.prompt_gm {
		log.Printf("SUMMARY mode is in effect for all species.\n\n")
	}

	// Loop through species data and make an uppercase copy of each name for comparison purposes later.
	// Also do some initializations.
	for sp_index = 0; sp_index < e.galaxy.num_species; sp_index++ {
		sp = e.spec_data[sp_index]
		e.ship_base = e.ship_data[sp_index]

		// convert name to upper case
		sp_name[sp_index] = strings.ToUpper(sp.name)

		// set the special flag back to zero
		for i = 0; i < sp.num_ships; i++ {
			e.ship = e.ship_base[i]
			e.ship.special = 0
		}
	}

	// create temporary logs for each species
	for i := 0; i < MAX_SPECIES; i++ {
		e.temp_log[i] = &bytes.Buffer{}
	}

	/* Main loop. For each species, take appropriate action. */
	num_battles = 0
	for arg_index = 0; arg_index < num_species; arg_index++ {
		species_number = sp_num[arg_index]
		sp = e.spec_data[species_number-1]

		// the following two items are needed by get_ship()
		e.species = sp
		e.ship_base = e.ship_data[species_number-1]

		// open orders file for this species
		filename := fmt.Sprintf("sp%02d.ord", species_number)
		b := &bytes.Buffer{}
		_, _ = b.ReadFrom(bytes.NewReader(e.spec_orders[species_number-1]))
		if e.input_file = fopen(filename, b); e.input_file == nil {
			if do_all_species != FALSE {
				if e.prompt_gm {
					log.Printf("No orders for species #%d, SP %s.\n", species_number, sp.name)
				}
				continue
			} else {
				panic(fmt.Sprintf("\n\tCannot open 'sp%02d.ord' for reading!\n\n", species_number))
			}
		}

		e.end_of_file = FALSE
		e.just_opened_file = TRUE /* Tell parse.c to skip mail header, if any. */

	find_start:

		/* Search for START COMBAT order. */
		found := FALSE
		for found == FALSE {
			command := e.get_command()
			if command == MESSAGE {
				/* Skip MESSAGE text. It may contain a line that starts with "start". */
				for {
					command = e.get_command()
					if command < 0 {
						fprintf(e.stderr, "WARNING: Unterminated MESSAGE command in file %s!\n", filename)
						break
					}

					if command == ZZZ {
						goto find_start
					}
				}
			}
			if command < 0 {
				break /* End of file. */
			}
			if command != START {
				continue
			}

			/* Get the first three letters of the keyword and convert to upper case. */
			e.skip_whitespace()
			var keyword string
			for i = 0; i < 3 && len(e.input_line_pointer) != 0; i++ {
				keyword += string(e.input_line_pointer[0])
				e.input_line_pointer = e.input_line_pointer[1:]
			}
			keyword = strings.ToUpper(keyword)

			if e.strike_phase != 0 {
				if keyword == "STR" {
					found = TRUE
				}
			} else {
				if keyword == "COM" {
					found = TRUE
				}
			}

			if found != FALSE {
				if e.strike_phase != 0 {
					log.Printf("Strike orders for species #%d, SP %s...\n", species_number, sp.name)
				} else {
					log.Printf("Combat orders for species #%d, SP %s...\n", species_number, sp.name)
				}
			} else {
				if e.strike_phase != 0 {
					log.Printf("No strike orders for species #%d, SP %s...\n", species_number, sp.name)
				} else {
					log.Printf("No combat orders for species #%d, SP %s...\n", species_number, sp.name)
				}
				goto done_orders
			}

			/* Open temporary log file for appending. */
			filename = fmt.Sprintf("sp%02d.temp.log", species_number)
			e.log_file = fopen(filename, e.spec_logs[species_number-1])
			if e.log_file == nil {
				fprintf(e.stderr, "\n\tCannot open '%s' for appending!\n\n", filename)
				panic(fmt.Sprintf("\n\tCannot open '%s' for appending!\n\n", filename))
			}
			e.append_log[species_number-1] = TRUE

			e.log_stdout = FALSE
			if e.strike_phase != 0 {
				e.log_string("\nStrike orders:\n")
			} else {
				e.log_string("\nCombat orders:\n")
			}
			// e.log_stdout = e.prompt_gm
			if e.prompt_gm {
				e.log_stdout = TRUE
			} else {
				e.log_stdout = FALSE
			}
		}

		/* Parse all combat commands for this species and save results for later use. */
		battle_index = -1
		for {
			command := e.get_command()
			if e.end_of_file != FALSE {
				break
			}
			if command == END {
				break
			}
			if command == BATTLE {
				num_enemies = 0 /* No enemies specified yet. */
				if e.get_value() == FALSE {
					e.bad_coordinates()
					continue
				}
				e.x = e.value

				if e.get_value() == FALSE {
					e.bad_coordinates()
					continue
				}
				e.y = e.value

				if e.get_value() == FALSE {
					e.bad_coordinates()
					continue
				}
				e.z = e.value

				// make sure that species is present at battle location
				found = FALSE
				for i = 0; i < e.num_locs; i++ {
					if location = e.locations_base[i]; location.s == species_number && location.x == e.x && location.y == e.y && location.z == e.z {
						found = TRUE
						break
					}
				}
				if found == FALSE {
					fprintf(e.log_file, "!!! Order ignored:\n")
					fprintf(e.log_file, "!!! %s", b2s(e.input_line))
					fprintf(e.log_file, "!!! Your species is not at this location!\n")
					continue
				}
				e.log_printf("  A battle order was issued for sector %d %d %d.\n", e.x, e.y, e.z)

				// add coordinates to list if not already there
				bat, found = nil, FALSE
				for i = 0; found == FALSE && i < len(e.battle_base); i++ {
					bat = e.battle_base[i]
					if e.x == bat.x && e.y == bat.y && e.z == bat.z {
						found, battle_index = TRUE, i
					}
				}

				if bat == nil || found == FALSE {
					// this is a new battle location
					sp_index = 0 // zero because this is the first species at the new location
					// battle_index will be the index of this new battle
					bat, battle_index = &battle_data{}, len(e.battle_base)
					bat.x, bat.y, bat.z = e.x, e.y, e.z
					bat.num_species_here = 1
					e.battle_base = append(e.battle_base, bat)
					num_battles = len(e.battle_base)
				} else {
					// add another species to existing battle location
					sp_index = bat.num_species_here // num_species_here is 1 based, the sp_index is 0 based
					bat.num_species_here++
				}
				bat.spec_num[sp_index] = species_number
				bat.special_target[sp_index] = 0              /* Default. */
				bat.transport_withdraw_age[sp_index] = 0      /* Default. */
				bat.warship_withdraw_age[sp_index] = 100      /* Default. */
				bat.fleet_withdraw_percentage[sp_index] = 100 /* Default. */
				bat.haven_x[sp_index] = 127                   /* 127 means not yet specified. */
				bat.engage_option[sp_index][0] = DEFENSE_IN_PLACE
				bat.num_engage_options[sp_index] = 1
				bat.can_be_surprised[sp_index] = FALSE
				bat.hijacker[sp_index] = FALSE
				bat.summary_only[sp_index] = default_summary
				for alienIndex := 0; alienIndex < MAX_SPECIES; alienIndex++ {
					bat.enemy_mine[sp_index][alienIndex] = 0
				}

				continue
			}

			if command == SUMMARY {
				if battle_index < 0 {
					e.battle_error(species_number)
					continue
				}
				bat.summary_only[sp_index] = TRUE
				e.log_string("    Summary mode was specified.\n")
				continue
			}

			if command == WITHDRAW {
				if battle_index < 0 {
					e.battle_error(species_number)
					continue
				}

				if e.get_value() == FALSE || e.value < 0 || e.value > 100 {
					e.bad_argument()
					continue
				}
				i, bat.transport_withdraw_age[sp_index] = e.value, e.value

				if e.get_value() == FALSE || e.value < 0 || e.value > 100 {
					e.bad_argument()
					continue
				}
				j, bat.warship_withdraw_age[sp_index] = e.value, e.value

				if e.get_value() == FALSE || e.value < 0 || e.value > 100 {
					e.bad_argument()
					continue
				}
				k, bat.fleet_withdraw_percentage[sp_index] = e.value, e.value

				e.log_printf("    Withdrawal conditions were set to %d %d %d.\n", i, j, k)
				continue
			}

			if command == HAVEN {
				if battle_index < 0 {
					e.battle_error(species_number)
					continue
				}

				if e.get_value() == FALSE {
					e.bad_coordinates()
					continue
				}
				i, bat.haven_x[sp_index] = e.value, e.value

				if e.get_value() == FALSE {
					e.bad_coordinates()
					continue
				}
				j, bat.haven_y[sp_index] = e.value, e.value

				if e.get_value() == FALSE {
					e.bad_coordinates()
					continue
				}
				k, bat.haven_z[sp_index] = e.value, e.value

				e.log_printf("    Haven location set to sector %d %d %d.\n", i, j, k)
				continue
			}

			if command == ENGAGE {
				if battle_index < 0 {
					e.battle_error(species_number)
					continue
				}

				option_index := bat.num_engage_options[sp_index]
				if option_index >= MAX_ENGAGE_OPTIONS {
					fprintf(e.log_file, "!!! Order ignored:\n")
					fprintf(e.log_file, "!!! %s", b2s(e.input_line))
					fprintf(e.log_file, "!!! Too many ENGAGE orders!\n")
					continue
				}

				if e.get_value() == FALSE || e.value < 0 || e.value > 7 {
					fprintf(e.log_file, "!!! Order ignored:\n")
					fprintf(e.log_file, "!!! %s", b2s(e.input_line))
					fprintf(e.log_file, "!!! Invalid ENGAGE option!\n")
					continue
				}
				option := e.value

				if e.strike_phase != 0 && (option > 4) {
					fprintf(e.log_file, "!!! Order ignored:\n")
					fprintf(e.log_file, "!!! %s", b2s(e.input_line))
					fprintf(e.log_file, "!!! Invalid ENGAGE option for strike phase!\n")
					continue
				}

				bat.engage_option[sp_index][option_index] = option

				// get planet to attack/defend, if any
				if option == PLANET_DEFENSE || (option >= PLANET_ATTACK && option <= SIEGE) {
					if e.get_value() == FALSE {
						fprintf(e.log_file, "!!! Order ignored:\n")
						fprintf(e.log_file, "!!! %s", b2s(e.input_line))
						fprintf(e.log_file, "!!! Missing planet argument in ENGAGE order!\n")
						continue
					}
					if e.value < 1 || e.value > 9 {
						fprintf(e.log_file, "!!! Order ignored:\n")
						fprintf(e.log_file, "!!! %s", b2s(e.input_line))
						fprintf(e.log_file, "!!! Invalid planet argument in ENGAGE order!\n")
						continue
					}

					bat.engage_planet[sp_index][option_index] = e.value
				} else {
					e.value = 0
					bat.engage_planet[sp_index][option_index] = 0
				}

				bat.num_engage_options[sp_index]++

				if e.value == 0 {
					e.log_printf("    Engagement order %d was specified.\n", option)
				} else {
					e.log_printf("    Engagement order %d %d was specified.\n", option, e.value)
				}

				continue
			}

			if command == HIDE {
				if battle_index < 0 {
					e.battle_error(species_number)
					continue
				}

				if e.get_ship() == FALSE {
					fprintf(e.log_file, "!!! Order ignored:\n")
					fprintf(e.log_file, "!!! %s", b2s(e.input_line))
					fprintf(e.log_file, "!!! Invalid or missing ship name!\n")
					continue
				}

				if e.ship.status != ON_SURFACE {
					fprintf(e.log_file, "!!! Order ignored:\n")
					fprintf(e.log_file, "!!! %s", b2s(e.input_line))
					fprintf(e.log_file, "!!! Ship must be landed to HIDE!\n")
					continue
				}

				e.ship.special = NON_COMBATANT

				e.log_printf("    %s will attempt to stay out of the battle.\n", e.ship_name(e.ship))

				continue
			}

			if command == TARGET {
				if battle_index < 0 {
					e.battle_error(species_number)
					continue
				}

				if e.get_value() == FALSE || e.value < 1 || e.value > 4 {
					fprintf(e.log_file, "!!! Order ignored:\n")
					fprintf(e.log_file, "!!! %s", b2s(e.input_line))
					fprintf(e.log_file, "!!! Invalid TARGET option!\n")
					continue
				}
				bat.special_target[sp_index] = e.value

				e.log_printf("    Strategic target %d was specified.\n", e.value)

				continue
			}

			if command == ATTACK || command == HIJACK {
				if battle_index < 0 {
					e.battle_error(species_number)
					continue
				}
				if command == HIJACK {
					bat.hijacker[sp_index] = TRUE
				}

				// check if this is an order to attack all declared enemies
				if e.get_value() != FALSE && e.value == 0 {
					for i = 0; i < e.galaxy.num_species; i++ {
						if species_number == i+1 {
							continue
						}

						alien_no := i + 1 // one-based species index
						if sp.enemy[alien_no-1] != FALSE {
							if num_enemies == MAX_SPECIES {
								fprintf(e.stderr, "\n\n\tToo many enemies to ATTACK or HIJACK!\n\n")
								panic("\n\n\tToo many enemies to ATTACK or HIJACK!\n\n")
							}
							if command == HIJACK {
								bat.enemy_mine[sp_index][num_enemies] = -alien_no
							} else {
								bat.enemy_mine[sp_index][num_enemies] = alien_no
							}
							num_enemies++
						}
					}

					if command == HIJACK {
						e.log_string("    An order was given to hijack all declared enemies.\n")
					} else {
						e.log_string("    An order was given to attack all declared enemies.\n")
					}

					continue
				}

				if num_enemies == MAX_SPECIES {
					fprintf(e.stderr, "\n\n\tToo many enemies to ATTACK or HIJACK!\n\n")
					panic("\n\n\tToo many enemies to ATTACK or HIJACK!\n\n")
				}

				// set 'n' to the species number of the named enemy
				temp_ptr := e.input_line_pointer
				if e.get_class_abbr() != SPECIES_ID {
					// check if SP abbreviation was accidentally omitted
					if isdigit(temp_ptr[0]) {
						e.input_line_pointer = temp_ptr
					} else if e.input_line_pointer[0] != ' ' && e.input_line_pointer[0] != '\t' {
						e.input_line_pointer = temp_ptr
					}
				}

				distorted_name := FALSE
				if e.get_value() != FALSE && !isalpha(e.input_line_pointer[0]) {
					if n = e.undistorted(e.value); n != 0 {
						distorted_name = TRUE
						goto att1
					}
				}
				if e.get_name() < 5 {
					e.bad_species()
					continue
				}

				// check for spelling error
				best_score, next_best_score = -9999, -9999
				for i = 0; i < e.galaxy.num_species; i++ {
					if sp_name[i][0] == 0 {
						continue
					}
					n = agrep_score([]byte(sp_name[i]), e.upper_name)
					if n > best_score {
						best_score, best_species_index = n, i
					} else if n > next_best_score {
						next_best_score = n
					}
				}

				name_length = len(sp_name[best_species_index])
				minimum_score = name_length - ((name_length / 7) + 1)
				if best_score < minimum_score || best_score == next_best_score {
					// score too low or another name with equal score
					e.bad_species()
					continue
				}

				n = best_species_index + 1

			att1:

				/* Make sure the named species is at the battle location. */
				found = FALSE
				location = nil
				for i = 0; i < e.num_locs; i++ {
					if location = e.locations_base[i]; location.s == n && location.x == bat.x && location.y == bat.y && location.z == bat.z {
						found = TRUE
						break
					}
				}

				// save species number temporarily in enemy_mine array
				if found != FALSE {
					if command == HIJACK {
						bat.enemy_mine[sp_index][num_enemies] = -n
					} else {
						bat.enemy_mine[sp_index][num_enemies] = n
					}
					num_enemies++
				}

				if command == HIJACK {
					e.log_string("    An order was given to hijack SP ")
				} else {
					e.log_string("    An order was given to attack SP ")
				}

				if distorted_name != FALSE {
					e.log_int(e.distorted(n))
				} else {
					e.log_string(e.spec_data[n-1].name)
				}
				e.log_string(".\n")

				continue
			}

			fprintf(e.log_file, "!!! Order ignored:\n")
			fprintf(e.log_file, "!!! %s", b2s(e.input_line))
			fprintf(e.log_file, "!!! Invalid combat command.\n")
		}

		fclose(e.log_file)

	done_orders:

		fclose(e.input_file)
	}

	// check each battle.
	// if a species specified a BATTLE command but did not specify any engage options, then add a DEFENSE_IN_PLACE option
	bat = nil
	for battle_index = 0; battle_index < num_battles; battle_index++ {
		bat = e.battle_base[battle_index]
		for i = 0; i < bat.num_species_here; i++ {
			if bat.num_engage_options[i] == 0 {
				bat.num_engage_options[i] = 1
				bat.engage_option[i][0] = DEFENSE_IN_PLACE
			}
		}
	}

	// initialize make_enemy array
	for i = 0; i < MAX_SPECIES; i++ {
		for j = 0; j < MAX_SPECIES; j++ {
			e.make_enemy[i][j] = 0
		}
	}

	// check each battle location.
	// if a species is at the location but has no combat orders, add it to the list of species at that battle, and apply defaults.
	// after all species are accounted for at the current battle location, do battle.
	bat = nil
	for battle_index = 0; battle_index < num_battles; battle_index++ {
		bat = e.battle_base[battle_index]
		e.x, e.y, e.z = bat.x, bat.y, bat.z

		// check file 'locations.dat' for other species at this location
		location = nil
		for location_index := 0; location_index < e.num_locs; location_index++ {
			location = e.locations_base[location_index]
			if location.x != e.x || location.y != e.y || location.z != e.z {
				continue
			}
			// check if species is already accounted for
			found := FALSE
			species_number = location.s
			for sp_index = 0; sp_index < bat.num_species_here; sp_index++ {
				if bat.spec_num[sp_index] == species_number {
					found = TRUE
					break
				}
			}
			if found != FALSE {
				continue
			}

			// species is present but did not give any combat orders.
			// this species will be included in the battle ONLY
			//      if it has ships in deep space or in orbit
			//   or if it has an unhidden, populated planet in this sector
			//   or if it has a hidden planet that is being explicitly attacked.
			found = FALSE
			sp = e.spec_data[species_number-1]
			num_pls := 0
			var namp *nampla_data
			for i = 0; i < sp.num_namplas; i++ {
				namp = e.namp_data[species_number-1][i]
				if namp.pn == 99 || namp.x != e.x || namp.y != e.y || namp.z != e.z {
					continue
				}
				if (namp.status & POPULATED) == 0 {
					continue
				}
				really_hidden := FALSE
				if namp.hidden != FALSE {
					/* If this species and planet is explicitly mentioned in
					   ATTACK/ENGAGE orders, then the planet cannot hide
					   during the battle. */
					really_hidden = TRUE
					for at_index := 0; at_index < bat.num_species_here; at_index++ {
						for j = 0; j < MAX_SPECIES; j++ {
							k = bat.enemy_mine[at_index][j]
							if k < 0 {
								k = -k
							}
							if k == species_number {
								for k = 0; k < bat.num_engage_options[at_index]; k++ {
									if bat.engage_option[at_index][k] >= PLANET_ATTACK && bat.engage_option[at_index][k] <= SIEGE && bat.engage_planet[at_index][k] == namp.pn {
										really_hidden = FALSE
										break
									}
								}
								if really_hidden == FALSE {
									break
								}
							}
						}
						if really_hidden == FALSE {
							break
						}
					}
				}
				if really_hidden != FALSE {
					continue
				}
				found = TRUE
				pl_num[num_pls] = namp.pn
				num_pls++
			}
			if found == FALSE {
				continue
			}

			var sh *ship_data
			for i = 0; i < sp.num_ships; i++ {
				sh = e.ship_data[species_number-1][i]
				if sh.pn == 99 || sh.x != e.x || sh.y != e.y || sh.z != e.z {
					continue
				} else if sh.status == UNDER_CONSTRUCTION {
					continue
				} else if sh.status == ON_SURFACE {
					continue
				} else if sh.status == JUMPED_IN_COMBAT {
					continue
				} else if sh.status == FORCED_JUMP {
					continue
				}
				found = TRUE
				break
			}
			if found == FALSE {
				continue
			}

			sp_index = bat.num_species_here
			bat.spec_num[sp_index] = location.s
			bat.special_target[sp_index] = 0
			bat.transport_withdraw_age[sp_index] = 0
			bat.warship_withdraw_age[sp_index] = 100
			bat.fleet_withdraw_percentage[sp_index] = 100
			bat.haven_x[sp_index] = 127
			bat.engage_option[sp_index][0] = DEFENSE_IN_PLACE
			bat.num_engage_options[sp_index] = 1
			if num_pls > 0 {
				// provide default Engage 2 options
				for i = 0; i < num_pls; i++ {
					bat.engage_option[sp_index][i+1] = PLANET_DEFENSE
					bat.engage_planet[sp_index][i+1] = pl_num[i]
				}
				bat.num_engage_options[sp_index] = num_pls + 1
			}
			bat.can_be_surprised[sp_index] = TRUE
			bat.hijacker[sp_index] = FALSE
			bat.summary_only[sp_index] = default_summary
			for alienIndex := 0; alienIndex < MAX_SPECIES; alienIndex++ {
				bat.enemy_mine[sp_index][alienIndex] = 0
			}
			bat.num_species_here++
		}

		// if haven locations have not been specified, provide random locations nearby
		for sp_index = 0; sp_index < bat.num_species_here; sp_index++ {
			if bat.haven_x[sp_index] != 127 {
				continue
			}
			for {
				i, j, k = e.x+2-e.rnd(3), e.y+2-e.rnd(3), e.z+2-e.rnd(3)
				if i != e.x || j != e.y || k != e.z {
					break
				}
			}
			bat.haven_x[sp_index] = i
			bat.haven_y[sp_index] = j
			bat.haven_z[sp_index] = k
		}

		// do battle at this battle location
		e.do_battle(bat)

		//if e.prompt_gm != FALSE {
		//	printf("Hit RETURN to continue...")
		//	e.stdout.fflush()
		//	e.stdin.fgets(answer, 16)
		//}
	}

	// declare new enmities
	for i := 0; i < e.galaxy.num_species; i++ {
		log_open := FALSE

		for j := 0; j < e.galaxy.num_species; j++ {
			if i == j {
				continue
			}

			// make_enemy[spNo-1][spNo-1] = spNo
			betrayed_species_number := e.make_enemy[i][j]
			if betrayed_species_number == 0 {
				continue
			}

			// clear ally flag
			e.spec_data[i].ally[betrayed_species_number-1] = FALSE

			/* Set enemy and contact bits (in case this is first encounter). */
			e.spec_data[i].enemy[betrayed_species_number-1] = TRUE
			e.spec_data[i].contact[betrayed_species_number-1] = TRUE

			if log_open == FALSE {
				/* Open temporary species log file for appending. */
				if e.spec_logs[i] == nil {
					e.spec_logs[i] = &bytes.Buffer{}
				}
				e.log_file = fopen("", e.spec_logs[i])
				e.append_log[i] = TRUE
				log_open = TRUE
			}

			e.log_string("\n!!! WARNING: Enmity has been automatically declared towards SP ")
			e.log_string(e.spec_data[j].name)
			e.log_string(" because they surprise-attacked SP ")
			e.log_string(e.spec_data[betrayed_species_number-1].name)
			e.log_string("!\n")
		}

		if log_open != FALSE {
			e.log_file = nil
		}
	}
}

func printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
