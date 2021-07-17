/*****************************************************************************
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
****************************************************************************/

package main

import (
	"fmt"
	"github.com/mdhender/fhcms/agrep"
	"io"
)

// Combat.c
/*
 *      Terminology: A "battle" consists of one or more "actions", and each
 *      action can take place at a different location in a star system.
 *      In turn, each action can be broken down into one or more "rounds",
 *      in which actual combat occurs.
 *
 *      A battle is defined by a variable of ttype "struct battle_data",
 *      and a pointer to a variable of this ttype, called "bat", is used
 *      throughout the combat routines.
 *
 *      An action is defined by a variable of ttype "struct action_data",
 *      and a pointer to a variable of this ttype, called "act", is used
 *      throughout the combat routines.
 */

func CombatMain(argc int, argv []string) {
	var (
		i, j, k                                                                                  int
		found                                                                                    bool
		command, species_number, sp_index, num_battles                                           int
		location_index, num_enemies, battle_index                                                int
		option_index, num_species                                                                int
		sp_num                                                                                   [MAX_SPECIES]int
		arg_index                                                                                int
		do_all_species                                                                           bool
		default_summary                                                                          bool
		at_index                                                                                 int
		really_hidden                                                                            bool
		num_pls                                                                                  int
		pl_num                                                                                   [9]int
		enemy_word_number                                                                        int
		enemy_bit_number                                                                         int
		log_open                                                                                 bool
		distorted_name                                                                           bool
		save, best_score                                                                         int
		next_best_score, best_species_index, betrayed_species_number, name_length, minimum_score int
		n, enemy_mask                                                                            int
		x, y, z, option                                                                          int
		keyword                                                                                  [4]byte
		answer                                                                                   [16]byte
		log_line                                                                                 [256]byte
		filename                                                                                 string
		sp_name                                                                                  [MAX_SPECIES][32]byte
		temp_ptr                                                                                 *byte
		temp_species_log, species_log                                                            io.Writer
		sp                                                                                       *species_data
		namp                                                                                     *nampla_data
		sh                                                                                       *ship_data_
		bat                                                                                      *battle_data
		locations_base, location                                                                 *sp_loc_data
		// at_number, n_bytes, species_fd int
		// at_namp *nampla_data
		// at_sh *ship_data_
		// at_sp *species_data
		// first_line bool
	)

	/* Seed random number generator. */
	// TODO: fix last_random = time(NULL);
	n = rnd(100) + rnd(200) + rnd(300)
	for i = 0; i < n; i++ {
		rnd(10)
	}

	/* Get commonly used data. */
	get_galaxy_data()
	get_planet_data()
	get_transaction_data()
	get_location_data()
	locations_base = &loc[0]

	/* Allocate memory for battle data. */
	battle_base = make([]*battle_data, MAX_BATTLES, MAX_BATTLES)
	if battle_base == nil {
		fprintf(stderr, "\nCannot allocate enough memory for battle data!\n\n")
		exit(-1)
	}

	/* Check arguments. If an argument is -s, then set SUMMARY mode for
	 *  everyone. The default is for players to receive a detailed report
	 *  of the battles. If an argument is -p, then prompt the GM before
	 *  saving results; otherwise, operate quietly; i.e, do not prompt GM
	 *  before saving results and do not display anything except errors.
	 *  Any additional arguments must be species numbers. If no species
	 *  numbers are specified, then do all species. */
	num_species = 0
	default_summary = false
	prompt_gm = false
	test_mode = false
	verbose_mode = false

	if strstr(argv[0], "Strike") != -1 {
		strike_phase = true
	}

	for i = 1; i < argc; i++ {
		if strcmp(argv[i], "-s") == 0 {
			default_summary = true
		} else if strcmp(argv[i], "-p") == 0 {
			prompt_gm = true
		} else if strcmp(argv[i], "-t") == 0 {
			test_mode = true
		} else if strcmp(argv[i], "-v") == 0 {
			verbose_mode = true
		} else {
			n = atoi(argv[i])
			sp_num[num_species] = n
			num_species++
		}
	}

	log_stdout = prompt_gm

	if num_species == 0 {
		num_species = galaxy.num_species
		for i = 0; i < num_species; i++ {
			sp_num[i] = i + 1
		}
		do_all_species = true
	} else {
		do_all_species = false
	}

	if default_summary && prompt_gm {
		printf("\nSUMMARY mode is in effect for all species.\n\n")
	}

	/* Read in species data and make an uppercase copy of each name for
	 *  comparison purposes later. Also do some initializations. */
	get_species_data()
	for sp_index = 0; sp_index < galaxy.num_species; sp_index++ {
		sp_name[sp_index][0] = 0 /* Initialize. */

		if !data_in_memory[sp_index] {
			continue /* No longer in game. */
		}
		sp = spec_data[sp_index]
		ship_base = ship_data[sp_index]

		/* Convert name to upper case. */
		for i = 0; i < 31; i++ {
			sp_name[sp_index][i] = toupper(sp.name[i])
		}

		for i = 0; i < sp.num_ships; i++ {
			ship = ship_base[i]
			ship.special = 0
		}
	}

	/* Main loop. For each species, take appropriate action. */
	num_battles = 0
	for arg_index = 0; arg_index < num_species; arg_index++ {
		species_number = sp_num[arg_index]
		if !data_in_memory[species_number-1] {
			continue
		}

		sp = &spec_data[species_number-1]

		/* The following two items are needed by get_ship(). */
		species = sp
		ship_base = ship_data[species_number-1]

		/* Open orders file for this species. */
		filename = fmt.Sprintf("sp%02d.ord", species_number)
		input_file = fopen(filename, "r")
		if input_file == nil {
			if do_all_species {
				if prompt_gm {
					printf("\nNo orders for species #%d, SP %s.\n",
						species_number, sp.name)
				}
				continue
			} else {
				fprintf(stderr, "\n\tCannot open '%s' for reading!\n\n", filename)
				exit(-1)
			}
		}

		end_of_file = false

		just_opened_file = true /* Tell parse.c to skip mail header,
		 *      if any. */
	find_start:

		/* Search for START COMBAT order. */
		found = false
		for !found {
			command = get_command()
			if command == MESSAGE {
				/* Skip MESSAGE text. It may contain a line that starts
				 *      with "start". */
				for {
					command = get_command()
					if command < 0 {
						fprintf(stderr,
							"WARNING: Unterminated MESSAGE command in file %s!\n",
							filename)
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

			/* Get the first three letters of the keyword and convert to
			 *  upper case. */
			skip_whitespace()
			for i = 0; i < 3; i++ {
				keyword[i] = toupper(input_line_pointer.get())
				input_line_pointer.incr()
			}
			keyword[3] = 0

			if strike_phase {
				if strcmp(string(byte4_to_byte(keyword)), "STR") == 0 {
					found = true
				}
			} else {
				if strcmp(string(byte4_to_byte(keyword)), "COM") == 0 {
					found = true
				}
			}
		}

		if found {
			if prompt_gm {
				if strike_phase {
					printf("\nStrike orders for species #%d, SP %s...\n", species_number, sp.name)
				} else {
					printf("\nCombat orders for species #%d, SP %s...\n", species_number, sp.name)
				}
			}
		} else {
			if prompt_gm {
				if strike_phase {
					printf("\nNo strike orders for species #%d, SP %s...\n", species_number, sp.name)
				} else {
					printf("\nNo combat orders for species #%d, SP %s...\n", species_number, sp.name)
				}
			}
			goto done_orders
		}

		/* Open temporary log file for appending. */
		filename = fmt.Sprintf("sp%02d.temp.log", species_number)
		log_file = fopen(filename, "a")
		if log_file == nil {
			fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename)
			exit(-1)
		}

		append_log[species_number-1] = true

		log_stdout = false
		if strike_phase {
			log_string("\nStrike orders:\n")
		} else {
			log_string("\nCombat orders:\n")
		}
		log_stdout = prompt_gm

		/* Parse all combat commands for this species and save results
		 * for later use. */
		battle_index = -1
		for {
			command = get_command()
			if end_of_file {
				break
			}

			if command == END {
				break
			}

			if command == BATTLE {
				num_enemies = 0 /* No enemies specified yet. */

				if get_value() == 0 {
					bad_coordinates()
					continue
				}
				x = value

				if get_value() == 0 {
					bad_coordinates()
					continue
				}
				y = value

				if get_value() == 0 {
					bad_coordinates()
					continue
				}
				z = value

				/* Make sure that species is present at battle location. */
				found = false
				location = locations_base - 1
				for i = 0; i < num_locs; i++ {
					location++
					if location.s != species_number {
						continue
					}
					if location.x != x {
						continue
					}
					if location.y != y {
						continue
					}
					if location.z != z {
						continue
					}

					found = true
					break
				}
				if !found {
					fprintf(log_file, "!!! Order ignored:\n")
					fprintf(log_file, "!!! %s", input_line)
					fprintf(log_file, "!!! Your species is not at this location!\n")
					continue
				}

				log_string("  A battle order was issued for sector ")
				log_int(x)
				log_char(' ')
				log_int(y)
				log_char(' ')
				log_int(z)
				log_string(".\n")

				/* Add coordinates to list if not already there. */
				found = false
				bat = battle_base
				for i = 0; i < num_battles; i++ {
					if x == bat.x && y == bat.y && z == bat.z {
						found = true
						battle_index = i
						break
					}

					bat++
				}

				if !found {
					/* This is a new battle location. */
					if num_battles == MAX_BATTLES {
						fprintf(stderr, "\n\n\tMAX_BATTLES exceeded! Edit file 'combat.h' and re-compile!\n\n")
						exit(-1)
					}
					battle_index = num_battles
					sp_index = 0
					bat.x = x
					bat.y = y
					bat.z = z
					bat.spec_num[0] = species_number
					bat.special_target[0] = 0              /* Default. */
					bat.transport_withdraw_age[0] = 0      /* Default. */
					bat.warship_withdraw_age[0] = 100      /* Default. */
					bat.fleet_withdraw_percentage[0] = 100 /* Default. */
					bat.haven_x[0] = 127
					/* 127 means not yet specified. */
					bat.engage_option[sp_index][0] = DEFENSE_IN_PLACE
					bat.num_engage_options[0] = 1
					bat.can_be_surprised[0] = false
					bat.hijacker[0] = false
					bat.summary_only[0] = default_summary
					bat.num_species_here = 1
					for i = 0; i < MAX_SPECIES; i++ {
						bat.enemy_mine[0][i] = 0
					}
					num_battles++
				} else {
					/* Add another species to existing battle location. */
					sp_index = bat.num_species_here
					bat.spec_num[sp_index] = species_number
					bat.special_target[sp_index] = 0              /* Default. */
					bat.transport_withdraw_age[sp_index] = 0      /* Default. */
					bat.warship_withdraw_age[sp_index] = 100      /* Default. */
					bat.fleet_withdraw_percentage[sp_index] = 100 /* Default. */
					bat.haven_x[sp_index] = 127
					/* 127 means not yet specified. */
					bat.engage_option[sp_index][0] = DEFENSE_IN_PLACE
					bat.num_engage_options[sp_index] = 1
					bat.can_be_surprised[sp_index] = false
					bat.hijacker[sp_index] = false
					bat.summary_only[sp_index] = default_summary
					bat.num_species_here++
					for i = 0; i < MAX_SPECIES; i++ {
						bat.enemy_mine[sp_index][i] = 0
					}
				}
				continue
			}

			if command == SUMMARY {
				if battle_index < 0 {
					battle_error(species_number)
					continue
				}

				bat.summary_only[sp_index] = true

				log_string("    Summary mode was specified.\n")

				continue
			}

			if command == WITHDRAW {
				if battle_index < 0 {
					battle_error(species_number)
					continue
				}

				if get_value() == 0 || value < 0 || value > 100 {
					bad_argument()
					continue
				}
				i = value
				bat.transport_withdraw_age[sp_index] = i

				if get_value() == 0 || value < 0 || value > 100 {
					bad_argument()
					continue
				}
				j = value
				bat.warship_withdraw_age[sp_index] = j

				if get_value() == 0 || value < 0 || value > 100 {
					bad_argument()
					continue
				}
				k = value
				bat.fleet_withdraw_percentage[sp_index] = k

				log_string("    Withdrawal conditions were set to ")
				log_int(i)
				log_char(' ')
				log_int(j)
				log_char(' ')
				log_int(k)
				log_string(".\n")

				continue
			}

			if command == HAVEN {
				if battle_index < 0 {
					battle_error(species_number)
					continue
				}

				if get_value() == 0 {
					bad_coordinates()
					continue
				}
				i = value
				bat.haven_x[sp_index] = value

				if get_value() == 0 {
					bad_coordinates()
					continue
				}
				j = value
				bat.haven_y[sp_index] = value

				if get_value() == 0 {
					bad_coordinates()
					continue
				}
				k = value
				bat.haven_z[sp_index] = value

				log_string("    Haven location set to sector ")
				log_int(i)
				log_char(' ')
				log_int(j)
				log_char(' ')
				log_int(k)
				log_string(".\n")

				continue
			}

			if command == ENGAGE {
				if battle_index < 0 {
					battle_error(species_number)
					continue
				}

				option_index = bat.num_engage_options[sp_index]
				if option_index >= MAX_ENGAGE_OPTIONS {
					fprintf(log_file, "!!! Order ignored:\n")
					fprintf(log_file, "!!! %s", input_line)
					fprintf(log_file, "!!! Too many ENGAGE orders!\n")
					continue
				}

				if get_value() == 0 || value < 0 || value > 7 {
					fprintf(log_file, "!!! Order ignored:\n")
					fprintf(log_file, "!!! %s", input_line)
					fprintf(log_file, "!!! Invalid ENGAGE option!\n")
					continue
				}
				option = value

				if strike_phase && (option > 4) {
					fprintf(log_file, "!!! Order ignored:\n")
					fprintf(log_file, "!!! %s", input_line)
					fprintf(log_file, "!!! Invalid ENGAGE option for strike phase!\n")
					continue
				}

				bat.engage_option[sp_index][option_index] = option

				/* Get planet to attack/defend, if any. */
				if option == PLANET_DEFENSE || (option >= PLANET_ATTACK &&
					option <= SIEGE) {
					if get_value() == 0 {
						fprintf(log_file, "!!! Order ignored:\n")
						fprintf(log_file, "!!! %s", input_line)
						fprintf(log_file, "!!! Missing planet argument in ENGAGE order!\n")
						continue
					}

					if value < 1 || value > 9 {
						fprintf(log_file, "!!! Order ignored:\n")
						fprintf(log_file, "!!! %s", input_line)
						fprintf(log_file, "!!! Invalid planet argument in ENGAGE order!\n")
						continue
					}

					bat.engage_planet[sp_index][option_index] = value
				} else {
					value = 0
					bat.engage_planet[sp_index][option_index] = 0
				}

				bat.num_engage_options[sp_index]++

				log_string("    Engagement order ")
				log_int(option)
				if value != 0 {
					log_char(' ')
					log_long(value)
				}
				log_string(" was specified.\n")

				continue
			}

			if command == HIDE {
				if battle_index < 0 {
					battle_error(species_number)
					continue
				}

				if !get_ship() {
					fprintf(log_file, "!!! Order ignored:\n")
					fprintf(log_file, "!!! %s", input_line)
					fprintf(log_file, "!!! Invalid or missing ship name!\n")
					continue
				}

				if ship.status != ON_SURFACE {
					fprintf(log_file, "!!! Order ignored:\n")
					fprintf(log_file, "!!! %s", input_line)
					fprintf(log_file, "!!! Ship must be landed to HIDE!\n")
					continue
				}

				ship.special = NON_COMBATANT

				log_string("    ")
				log_string(ship_name(ship))
				log_string(" will attempt to stay out of the battle.\n")

				continue
			}

			if command == TARGET {
				if battle_index < 0 {
					battle_error(species_number)
					continue
				}

				if get_value() == 0 || value < 1 || value > 4 {
					fprintf(log_file, "!!! Order ignored:\n")
					fprintf(log_file, "!!! %s", input_line)
					fprintf(log_file, "!!! Invalid TARGET option!\n")
					continue
				}
				bat.special_target[sp_index] = value

				log_string("    Strategic target ")
				log_long(value)
				log_string(" was specified.\n")

				continue
			}

			if command == ATTACK || command == HIJACK {
				if battle_index < 0 {
					battle_error(species_number)
					continue
				}

				if command == HIJACK {
					bat.hijacker[sp_index] = true
				}

				/* Check if this is an order to attack all declared enemies. */
				if get_value() && value == 0 {
					for i = 0; i < galaxy.num_species; i++ {
						if species_number == i+1 {
							continue
						}

						if !data_in_memory[i] {
							continue
						}

						enemy_word_number = i / 32
						enemy_bit_number = i % 32
						enemy_mask = 1 << enemy_bit_number

						if (sp.enemy[enemy_word_number] & enemy_mask) != 0 {
							if num_enemies == MAX_SPECIES {
								fprintf(stderr,
									"\n\n\tToo many enemies to ATTACK or HIJACK!\n\n")
								exit(-1)
							}
							if command == HIJACK {
								bat.enemy_mine[sp_index][num_enemies] = -(i + 1)
							} else {
								bat.enemy_mine[sp_index][num_enemies] = i + 1
							}
							num_enemies++
						}
					}

					if command == HIJACK {
						log_string("    An order was given to hijack all declared enemies.\n")
					} else {
						log_string("    An order was given to attack all declared enemies.\n")
					}

					continue
				}

				if num_enemies == MAX_SPECIES {
					fprintf(stderr, "\n\n\tToo many enemies to ATTACK or HIJACK!\n\n")
					exit(-1)
				}

				/* Set 'n' to the species number of the named enemy. */
				temp_ptr = input_line_pointer
				if get_class_abbr() != SPECIES_ID {
					/* Check if SP abbreviation was accidentally omitted. */
					if isdigit(*temp_ptr) {
						input_line_pointer = temp_ptr
					} else if *input_line_pointer != ' ' &&
						*input_line_pointer != '\t' {
						input_line_pointer = temp_ptr
					}
				}

				distorted_name = false
				if get_value() {
					if !isalpha(input_line_pointer.get()) {
						if n = undistorted(value); n != 0 {
							distorted_name = true
							goto att1
						}
					}
				} else if get_name() < 5 {
					bad_species()
					continue
				}

				/* Check for spelling error. */
				best_score = -9999
				next_best_score = -9999
				for i = 0; i < galaxy.num_species; i++ {
					if *sp_name[i] == 0 {
						continue
					}

					n = agrep.Score(sp_name[i], upper_name)
					if n > best_score {
						best_score = n
						best_species_index = i
					} else if n > next_best_score {
						next_best_score = n
					}
				}

				name_length = strlen(sp_name[best_species_index])
				minimum_score = name_length - ((name_length / 7) + 1)

				if best_score < minimum_score ||
					best_score == next_best_score {
					/* Score too low or another name with equal score. */
					bad_species()
					continue
				}

				n = best_species_index + 1

			att1:

				/* Make sure the named species is at the battle location. */
				found = false
				location = locations_base - 1
				for i = 0; i < num_locs; i++ {
					location++
					if location.s != n {
						continue
					}
					if location.x != bat.x {
						continue
					}
					if location.y != bat.y {
						continue
					}
					if location.z != bat.z {
						continue
					}

					found = true
					break
				}

				/* Save species number temporarily in enemy_mine array. */
				if found {
					if command == HIJACK {
						bat.enemy_mine[sp_index][num_enemies] = -n
					} else {
						bat.enemy_mine[sp_index][num_enemies] = n
					}
					num_enemies++
				}

				if command == HIJACK {
					log_string("    An order was given to hijack SP ")
				} else {
					log_string("    An order was given to attack SP ")
				}

				if distorted_name {
					log_int(distorted(n))
				} else {
					log_string(spec_data[n-1].name)
				}
				log_string(".\n")

				continue
			}

			fprintf(log_file, "!!! Order ignored:\n")
			fprintf(log_file, "!!! %s", input_line)
			fprintf(log_file, "!!! Invalid combat command.\n")
		}

		fclose(log_file)

	done_orders:

		fclose(input_file)
	}

	/* Check each battle.  If a species specified a BATTLE command but did
	 *   not specify any engage options, then add a DEFENSE_IN_PLACE option. */
	bat = battle_base
	for battle_index = 0; battle_index < num_battles; battle_index++ {
		for i = 0; i < bat.num_species_here; i++ {
			if bat.num_engage_options[i] == 0 {
				bat.num_engage_options[i] = 1
				bat.engage_option[i][0] = DEFENSE_IN_PLACE
			}
		}

		bat++
	}

	/* Initialize make_enemy array. */
	for i = 0; i < galaxy.num_species; i++ {
		for j = 0; j < galaxy.num_species; j++ {
			make_enemy[i][j] = 0
		}
	}

	/* Check each battle location. If a species is at the location
	 *  but has no combat orders, add it to the list of species at that
	 *  battle, and apply defaults. After all species are accounted for
	 *  at the current battle location, do battle. */
	bat = battle_base - 1
	for battle_index = 0; battle_index < num_battles; battle_index++ {
		bat++

		x = bat.x
		y = bat.y
		z = bat.z

		/* Check file 'locations.dat' for other species at this location. */
		location = locations_base - 1
		for location_index = 0; location_index < num_locs; location_index++ {
			location++
			if location.x != x {
				continue
			}
			if location.y != y {
				continue
			}
			if location.z != z {
				continue
			}

			/* Check if species is already accounted for. */
			found = false
			species_number = location.s
			for sp_index = 0; sp_index < bat.num_species_here; sp_index++ {
				if bat.spec_num[sp_index] == species_number {
					found = true
					break
				}
			}

			if found {
				continue
			}

			/* Species is present but did not give any combat orders.
			 *  This species will be included in the battle ONLY if it has
			 *  ships in deep space or in orbit or if it has an unhidden,
			 *  populated planet in this sector or if it has a hidden
			 *  planet that is being explicitly attacked. */
			found = false

			sp = &spec_data[species_number-1]

			num_pls = 0

			namp = namp_data[species_number-1] - 1
			for i = 0; i < sp.num_namplas; i++ {
				namp++

				if namp.pn == 99 {
					continue
				}
				if namp.x != x {
					continue
				}
				if namp.y != y {
					continue
				}
				if namp.z != z {
					continue
				}
				if (namp.status & POPULATED) == 0 {
					continue
				}

				really_hidden = false
				if namp.hidden {
					/* If this species and planet is explicitly mentioned in
					 *  ATTACK/ENGAGE orders, then the planet cannot hide
					 *  during the battle. */

					really_hidden = true

					for at_index = 0; at_index < bat.num_species_here; at_index++ {
						for j = 0; j < MAX_SPECIES; j++ {
							k = bat.enemy_mine[at_index][j]
							if k < 0 {
								k = -k
							}
							if k == species_number {
								for k = 0; k < bat.num_engage_options[at_index]; k++ {
									if bat.engage_option[at_index][k] >= PLANET_ATTACK && bat.engage_option[at_index][k] <= SIEGE && bat.engage_planet[at_index][k] == namp.pn {
										really_hidden = false
										break
									}
								}
								if !really_hidden {
									break
								}
							}
						}
						if !really_hidden {
							break
						}
					}
				}

				if really_hidden {
					continue
				}

				found = true
				pl_num[num_pls] = namp.pn
				num_pls++
			}

			sh = ship_data[species_number-1] - 1
			for i = 0; i < sp.num_ships; i++ {
				sh++

				if sh.pn == 99 {
					continue
				}
				if sh.x != x {
					continue
				}
				if sh.y != y {
					continue
				}
				if sh.z != z {
					continue
				}
				if sh.status == UNDER_CONSTRUCTION {
					continue
				}
				if sh.status == ON_SURFACE {
					continue
				}
				if sh.status == JUMPED_IN_COMBAT {
					continue
				}
				if sh.status == FORCED_JUMP {
					continue
				}
				found = true

				break
			}

			if !found {
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
				/* Provide default Engage 2 options. */
				for i = 0; i < num_pls; i++ {
					bat.engage_option[sp_index][i+1] = PLANET_DEFENSE
					bat.engage_planet[sp_index][i+1] = pl_num[i]
				}
				bat.num_engage_options[sp_index] = num_pls + 1
			}
			bat.can_be_surprised[sp_index] = true
			bat.hijacker[sp_index] = false
			bat.summary_only[sp_index] = default_summary
			for i = 0; i < MAX_SPECIES; i++ {
				bat.enemy_mine[sp_index][i] = 0
			}
			bat.num_species_here++
		}

		/* If haven locations have not been specified, provide random
		 * locations nearby. */
		for sp_index = 0; sp_index < bat.num_species_here; sp_index++ {
			if bat.haven_x[sp_index] != 127 {
				continue
			}

			for {
				i = x + 2 - rnd(3)
				j = y + 2 - rnd(3)
				k = z + 2 - rnd(3)

				if i != x || j != y || k != z {
					break
				}
			}

			bat.haven_x[sp_index] = i
			bat.haven_y[sp_index] = j
			bat.haven_z[sp_index] = k
		}

		/* Do battle at this battle location. */
		do_battle(bat)

		if prompt_gm {
			printf("Hit RETURN to continue...")
			fflush(stdout)
			fgets(answer, 16, stdin)
		}
	}

	/* Declare new enmities. */
	for i = 0; i < galaxy.num_species; i++ {
		log_open = false

		for j = 0; j < galaxy.num_species; j++ {
			if i == j {
				continue
			}

			betrayed_species_number = make_enemy[i][j]
			if betrayed_species_number == 0 {
				continue
			}

			enemy_word_number = j / 32
			enemy_bit_number = j % 32
			enemy_mask = 1 << enemy_bit_number

			spec_data[i].ally[j] = false /* Clear ally bit. */

			/* Set enemy and contact bits (in case this is first encounter). */
			spec_data[i].enemy[j] = true
			spec_data[i].contact[j] = true

			data_modified[i] = true

			if !log_open {
				/* Open temporary species log file for appending. */
				filename = fmt.Sprintf("sp%02d.temp.log", i+1)
				log_file = fopen(filename, "a")
				if log_file == NULL {
					fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename)
					exit(-1)
				}

				append_log[i] = true
				log_open = true
			}

			log_string("\n!!! WARNING: Enmity has been automatically declared towards SP ")
			log_string(spec_data[j].name)
			log_string(" because they surprise-attacked SP ")
			log_string(spec_data[betrayed_species_number-1].name)
			log_string("!\n")
		}

		if log_open {
			fclose(log_file)
		}
	}

	save = true
	if prompt_gm {
		printf("\n*** Gamemaster safe-abort option ... ttype q or Q to quit: ")
		fflush(stdout)
		fgets(answer, 16, stdin)
		if answer[0] == 'q' || answer[0] == 'Q' {
			save = false
		}
	}

	/* If results are to be saved, append temporary logs to actual species
	 *  logs. In either case, delete temporary logs. */
	for i = 0; i < galaxy.num_species; i++ {
		if !append_log[i] {
			continue
		}

		if save {
			filename = fmt.Sprintf("sp%02d.log", i+1)
			species_log = fopen(filename, "a")
			if species_log == NULL {
				fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename)
				exit(-1)
			}
		}

		filename = fmt.Sprintf("sp%02d.temp.log", i+1)

		if save {
			temp_species_log = fopen(filename, "r")
			if temp_species_log == NULL {
				fprintf(stderr, "\n\tCannot open '%s' for reading!\n\n", filename)
				exit(-1)
			}

			/* Copy temporary log to permanent species log. */
			for fgets(log_line, 256, temp_species_log) != nil {
				fputs(log_line, species_log)
			}

			fclose(temp_species_log)
			fclose(species_log)
		}

		/* Delete temporary log file. */
		unlink(filename)
	}

	if save {
		save_planet_data()
		save_species_data()
		save_transaction_data()
	}

	free_species_data()

	exit(0)
}
