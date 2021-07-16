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

//*************************************************************************
// combat_utils.c

func power(tonnage int) int {
    if (tonnage > 4068) {
        fprintf(stderr,                "\n\n\tLong integer overflow will occur in call to 'power(tonnage)'!\n");
        fprintf(stderr, "\t\tActual call is power(%d).\n\n", tonnage);
        exit(-1);
    }

    if (tonnage <= 100) {
        return ship_power[tonnage];
    }
    // Tonnage is not in table.
    // Break it up into two halves and get approximate result = 1.149 * (x1 + x2), using recursion if necessary.
    t1 := tonnage / 2
    t2 := tonnage - t1;
    return 1149 * (power(t1) + power(t2)) / 1000;
}


func battle_error(species_number int) {
    fprintf(log_file, "!!! Order ignored:\n");
    fprintf(log_file, "!!! %s", input_line);
    fprintf(log_file, "!!! Missing BATTLE command!\n");
    return;
}

func bad_species() {
    fprintf(log_file, "!!! Order ignored:\n");
    fprintf(log_file, "!!! %s", input_line);
    fprintf(log_file, "!!! Invalid species name!\n");
    return;
}

func bad_argument() {
    fprintf(log_file, "!!! Order ignored:\n");
    fprintf(log_file, "!!! %s", input_line);
    fprintf(log_file, "!!! Invalid argument in command.\n");
    return;
}

func bad_coordinates() {
    fprintf(log_file, "!!! Order ignored:\n");
    fprintf(log_file, "!!! %s", input_line);
    fprintf(log_file, "!!! Invalid coordinates in command.\n");
    return;
}

//*************************************************************************
// cons_op.c

func consolidate_option(option, location int) {
    /* Only attack options go in list. */
    if (option < DEEP_SPACE_FIGHT) {
        return;
    }

    /* Make sure pre-requisites are already in the list. Bombardment, and
     *  germ warfare must follow a successful planet attack. */
    if (option > PLANET_ATTACK) {
        consolidate_option(PLANET_ATTACK, location);
    }

    /* Check if option and location are already in list. */
    for i := 0; i < num_combat_options; i++ {
        if (option == combat_option[i] && location == combat_location[i]) {
            return;
        }
    }

    /* Add new option to list. */
    combat_option[num_combat_options]   = option;
    combat_location[num_combat_options] = location;
    num_combat_options++
}

//*************************************************************************
// dis_ship.c

func disbanded_ship(ship *ship_data_) bool {
    var nampla_index int
    var nampla *nampla_data

    nampla = nampla_base - 1;
    for nampla_index = 0; nampla_index < species.num_namplas; nampla_index++ {
        nampla++

        if (nampla.x != ship.x) {
            continue;
        }
        if (nampla.y != ship.y) {
            continue;
        }
        if (nampla.z != ship.z) {
            continue;
        }
        if (nampla.pn != ship.pn) {
            continue;
        }
        if ((nampla.status & DISBANDED_COLONY) == 0) {
            continue;
        }
        if (ship.ttype != STARBASE && ship.status == IN_ORBIT) {
            continue;
        }

        /* This ship is either on the surface of a disbanded colony or is a starbase orbiting a disbanded colony. */
        return true;
    }

    return false;
}

//*************************************************************************
// do_ally.c

func do_ALLY_command() {
    var i, array_index, bit_number int
    var bit_mask int

    /* Get name of species that is being declared an ally. */
    if (!get_species_name()) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid or missing argument in ALLY command.\n");
        return;
    }

    /* Get array index and bit mask. */
    array_index = (g_spec_number - 1) / 32;
    bit_number  = (g_spec_number - 1) % 32;
    bit_mask    = 1 << bit_number;

    /* Check if we've met this species and make sure it is not an enemy. */
    if !species.contact[array_index] {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't declare alliance with a species you haven't met.\n");
        return;
    }

    /* Set/clear the appropriate bit. */
    species.ally[array_index]  true    /* Set ally bit. */
    species.enemy[array_index] false   /* Clear enemy bit. */

    /* Log the result. */
    log_string("    Alliance was declared with ");
    if (bit_mask == 0) {
        log_string("ALL species");
    }else{
        log_string("SP ");
        log_string(g_spec_name);
    }
    log_string(".\n");
}

//*************************************************************************
// do_amb.c

func do_AMBUSH_command() {
    var n, status int
    var cost int

    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get amount to spend. */
    status = get_value();
    if (status == 0 || value < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid or missing amount.\n");
        return;
    }
    if (value == 0) {
        value = balance;
    }
    if (value == 0) {
        return;
    }
    cost = value;

    /* Check if planet is under siege. */
    if (nampla.siege_eff != 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Besieged planet cannot ambush!\n");
        return;
    }

    /* Check if sufficient funds are available. */
    if (check_bounced(cost)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        return;
    }

    /* Increment amount spent on ambush. */
    nampla.use_on_ambush += cost;

    /* Log transaction. */
    log_string("    Spent ");
    log_long(cost);
    log_string(" in preparation for an ambush.\n");
}

//*************************************************************************
// do_base.c


func do_BASE_command() {
    int i, n, found, su_count, original_count, item_class, name_length,
        unused_ship_available, new_tonnage, max_tonnage, new_starbase,
        source_is_a_planet, age_new;

    char x, y, z, pn, upper_ship_name[32], *original_line_pointer;

    struct nampla_data *source_nampla;
    struct ship_data *  source_ship, *starbase, *unused_ship;


    /* Get number of starbase units to use. */
    i = get_value();
    if (i == 0) {
        value = 0;
    }else{
        /* Make sure value is meaningful. */
        if (value < 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid SU count in BASE command.\n");
            return;
        }
    }
    su_count       = value;
    original_count = su_count;

    /* Get source of starbase units. */
    original_line_pointer = input_line_pointer;
    if (!get_transfer_point()) {
        input_line_pointer = original_line_pointer;
        fix_separator();        /* Check for missing comma or tab. */
        if (!get_transfer_point()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid source location in BASE command.\n");
            return;
        }
    }

    /* Make sure everything makes sense. */
    if (abbr_type == SHIP_CLASS) {
        source_is_a_planet = false;
        source_ship        = ship;

        if (source_ship.status == UNDER_CONSTRUCTION) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! %s is still under construction!\n",
                    ship_name(source_ship));
            return;
        }

        if (source_ship.status == FORCED_JUMP ||
            source_ship.status == JUMPED_IN_COMBAT) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
            return;
        }

        if (su_count == 0) {
            su_count = source_ship.item_quantity[SU];
        }
        if (su_count == 0) {
            return;
        }
        if (source_ship.item_quantity[SU] < su_count) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! %s does not enough starbase units!\n",
                    ship_name(source_ship));
            return;
        }

        x  = source_ship.x;     y = source_ship.y;     z = source_ship.z;
        pn = source_ship.pn;
    }else {      /* Source is a planet. */
        source_is_a_planet = true;
        source_nampla      = nampla;

        if (su_count == 0) {
            su_count = source_nampla.item_quantity[SU];
        }
        if (su_count == 0) {
            return;
        }
        if (source_nampla.item_quantity[SU] < su_count) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! PL %s does not have enough starbase units!\n",
                    source_nampla.name);
            return;
        }

        x  = source_nampla.x;   y = source_nampla.y;   z = source_nampla.z;
        pn = source_nampla.pn;
    }

    /* Get starbase name. */
    if (get_class_abbr() != SHIP_CLASS || abbr_index != BA) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Invalid starbase name.\n");
        return;
    }
    name_length = get_name();

    /* Search all ships for name. */
    found = false;
    ship  = ship_base - 1;
    unused_ship_available = false;
    for (ship_index = 0; ship_index < species.num_ships; ship_index++) {
        ++ship;

        if (ship.pn == 99) {
            unused_ship_available = true;
            unused_ship           = ship;
            continue;
        }

        /* Make upper case copy of ship name. */
        for (i = 0; i < 32; i++) {
            upper_ship_name[i] = toupper(ship.name[i]);
        }

        /* Compare names. */
        if (strcmp(upper_ship_name, upper_name) == 0) {
            found = true;
            break;
        }
    }

    if (found) {
        if (ship.ttype != STARBASE) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship name already in use.\n");
            return;
        }

        if (ship.x != x || ship.y != y || ship.z != z) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Starbase units and starbase are not at same X Y Z.\n");
            return;
        }
        starbase     = ship;
        new_starbase = false;
    }else {
        if (unused_ship_available) {
            starbase = unused_ship;
        }else{
            /* Make sure we have enough memory for new starbase. */
            if (num_new_ships[species_index] == NUM_EXTRA_SHIPS) {
                fprintf(stderr, "\n\n\tInsufficient memory for new starbase!\n\n");
                exit(-1);
            }
            ++num_new_ships[species_index];
            starbase = ship_base + (int)species.num_ships;
            ++species.num_ships;
            delete_ship(starbase);              /* Initialize everything to zero. */
        }

        /* Initialize non-zero data for new ship. */
        strcpy(starbase.name, original_name);
        starbase.x  = x;
        starbase.y  = y;
        starbase.z  = z;
        starbase.pn = pn;
        if (pn == 0) {
            starbase.status = IN_DEEP_SPACE;
        }else{
            starbase.status = IN_ORBIT;
        }
        starbase.ttype           = STARBASE;
        starbase.class          = BA;
        starbase.tonnage        = 0;
        starbase.age            = -1;
        starbase.remaining_cost = 0;

        /* Everything else was set to zero in above call to 'delete_ship'. */

        new_starbase = true;
    }

    /* Make sure that starbase is not being built in the deep space section
     *  of a star system .*/
    if (starbase.pn == 0) {
        star = star_base - 1;
        for (i = 0; i < num_stars; i++) {
            ++star;

            if (star.x != x) {
                continue;
            }
            if (star.y != y) {
                continue;
            }
            if (star.z != z) {
                continue;
            }

            if (star.num_planets < 1) {
                break;
            }

            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Starbase cannot be built in deep space if there are planets available!\n");
            if (new_starbase) {
                delete_ship(starbase);
            }
            return;
        }
    }

    /* Make sure species can build a starbase of this size. */
    max_tonnage = species.tech_level[MA] / 2;
    new_tonnage = starbase.tonnage + su_count;
    if (new_tonnage > max_tonnage && original_count == 0) {
        su_count = max_tonnage - starbase.tonnage;
        if (su_count < 1) {
            if (new_starbase) {
                delete_ship(starbase);
            }
            return;
        }
        new_tonnage = starbase.tonnage + su_count;
    }

    if (new_tonnage > max_tonnage) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Maximum allowable tonnage exceeded.\n");
        if (new_starbase) {
            delete_ship(starbase);
        }
        return;
    }

    /* Finish up and log results. */
    log_string("    ");
    if (starbase.tonnage == 0) {
        log_string(ship_name(starbase));
        log_string(" was constructed.\n");
    }else {
        starbase.age =         /* Weighted average. */
                        ((starbase.age * starbase.tonnage) - su_count)
                        / new_tonnage;
        log_string("Size of ");  log_string(ship_name(starbase));
        log_string(" was increased to ");
        log_string(commas(10000L * (long)new_tonnage));
        log_string(" tons.\n");
    }

    starbase.tonnage = new_tonnage;

    if (source_is_a_planet) {
        source_nampla.item_quantity[SU] -= su_count;
    }else{
        source_ship.item_quantity[SU] -= su_count;
    }
}

//*************************************************************************
// do_bat.c


func do_battle(bat *battle_data) {
    int i, j, k, species_index, species_number, num_sp, save,
        max_rounds, round_number, battle_here, fight_here,
        unit_index, option_index, current_species, temp_status,
        temp_pn, num_namplas, array_index, bit_number, first_action,
        traitor_number, betrayed_number, betrayal, need_comma,
        true_value, do_withdraw_check_first;

    short identifiable_units[MAX_SPECIES],
          unidentifiable_units[MAX_SPECIES];

    long n, bit_mask;

    char x, y, z, where, option, filename[32], enemy,
         enemy_num[MAX_SPECIES], log_line[256];

    FILE *combat_log, *species_log;

    struct action_data  act;
    struct nampla_data *namp, *attacked_nampla;
    struct ship_data *  sh;


    ambush_took_place = false;

    /* Open log file for writing. */
    log_file = fopen("combat.log", "w");
    if (log_file == NULL) {
        fprintf(stderr, "\n\tCannot open 'combat.log' for writing!\n\n");
        exit(-1);
    }

    /* Open summary file for writing. */
    summary_file = fopen("summary.log", "w");
    if (summary_file == NULL) {
        fprintf(stderr, "\n\tCannot open 'summary.log' for writing!\n\n");
        exit(-1);
    }
    log_summary = true;

    /* Get data for all species present at this battle. */
    num_sp = bat.num_species_here;
    for (species_index = 0; species_index < num_sp; ++species_index) {
        species_number           = bat.spec_num[species_index];
        c_species[species_index] = &spec_data[species_number - 1];
        c_nampla[species_index]  = namp_data[species_number - 1];
        c_ship[species_index]    = ship_data[species_number - 1];
        if (data_in_memory[species_number - 1]) {
            data_modified[species_number - 1] = true;
        }else{
            fprintf(stderr, "\n\tData for species #%d is needed but is not available!\n\n",
                    species_number);
            exit(-1);
        }

        /* Determine number of identifiable and unidentifiable units present. */
        identifiable_units[species_index]   = 0;
        unidentifiable_units[species_index] = 0;

        namp = c_nampla[species_index] - 1;
        for (i = 0; i < c_species[species_index].num_namplas; i++) {
            ++namp;

            if (namp.x != bat.x) {
                continue;
            }
            if (namp.y != bat.y) {
                continue;
            }
            if (namp.z != bat.z) {
                continue;
            }

            if (namp.status & POPULATED) {
                ++identifiable_units[species_index];
            }
        }

        sh = c_ship[species_index] - 1;
        for (i = 0; i < c_species[species_index].num_ships; i++) {
            ++sh;

            if (sh.x != bat.x) {
                continue;
            }
            if (sh.y != bat.y) {
                continue;
            }
            if (sh.z != bat.z) {
                continue;
            }
            if (sh.status == UNDER_CONSTRUCTION) {
                continue;
            }
            if (sh.status == JUMPED_IN_COMBAT) {
                continue;
            }
            if (sh.status == FORCED_JUMP) {
                continue;
            }

            sh.dest_x = 0;     /* Not yet exposed. */
            sh.dest_y = 100;   /* Shields at 100%. */

            if (sh.item_quantity[FD] == sh.tonnage) {
                ++unidentifiable_units[species_index];
            }else{
                ++identifiable_units[species_index];
            }
        }

        if (identifiable_units[species_index] > 0 ||
            unidentifiable_units[species_index] == 0) {
            field_distorted[species_index] = false;
        }else{
            field_distorted[species_index] = true;
        }
    }

    /* Start log of what's happening. */
    if (strike_phase) {
        log_string("\nStrike log:\n");
    }else{
        log_string("\nCombat log:\n");
    }
    first_battle = false;

    log_string("\n  Battle orders were received for sector ");  log_int(bat.x);
    log_string(", ");  log_int(bat.y);  log_string(", ");  log_int(bat.z);
    log_string(". The following species are present:\n\n");

    /* Convert enemy_mine array from a list of species numbers to an array
     *  of true/false values whose indices are:
     *
     *                  [species_index1][species_index2]
     *
     *  such that the value will be true if #1 mentioned #2 in an ATTACK
     *  or HIJACK command.  The actual true value will be 1 for ATTACK or
     *  2 for HIJACK. */

    for (species_index = 0; species_index < num_sp; ++species_index) {
        /* Make copy of list of enemies. */
        for (i = 0; i < MAX_SPECIES; i++) {
            enemy_num[i] = bat.enemy_mine[species_index][i];
            bat.enemy_mine[species_index][i] = false;
        }

        for (i = 0; i < MAX_SPECIES; i++) {
            enemy = enemy_num[i];
            if (enemy == 0) {
                break;                  /* No more enemies in list. */
            }
            if (enemy < 0) {
                enemy      = -enemy;
                true_value = 2;         /* This is a hijacking. */
            }else {
                true_value = 1;         /* This is a normal attack. */
            }

            /* Convert absolute species numbers to species indices that
             *  have been assigned in the current battle. */
            for (j = 0; j < num_sp; j++) {
                if (enemy == bat.spec_num[j]) {
                    bat.enemy_mine[species_index][j] = true_value;
                }
            }
        }
    }

    /* For each species that has been mentioned in an attack order, check
     *  if it can be surprised. A species can only be surprised if it has
     *  not given a BATTLE order and if it is being attacked ONLY by one
     *  or more ALLIES. */
    for (species_index = 0; species_index < num_sp; ++species_index) {
        j           = bat.spec_num[species_index] - 1;
        array_index = j / 32;
        bit_number  = j % 32;
        bit_mask    = 1 << bit_number;

        for (i = 0; i < num_sp; i++) {
            if (i == species_index) {
                continue;
            }

            if (!bat.enemy_mine[species_index][i]) {
                continue;
            }

            if (field_distorted[species_index]) {
                /* Attacker is field-distorted. Surprise not possible. */
                bat.can_be_surprised[i] = false;
                continue;
            }

            if ((c_species[i].ally[array_index] & bit_mask)) {
                betrayal = true;
            }else{
                betrayal = false;
            }

            if (betrayal) {
                /* Someone is being attacked by an ALLY. */
                traitor_number  = bat.spec_num[species_index];
                betrayed_number = bat.spec_num[i];
                make_enemy[betrayed_number - 1][traitor_number - 1] = betrayed_number;
                auto_enemy(traitor_number, betrayed_number);
            }

            if (!bat.can_be_surprised[i]) {
                continue;
            }

            if (!betrayal) {    /* At least one attacker is not an ally. */
                bat.can_be_surprised[i] = false;
            }
        }
    }

    /* For each species that has been mentioned in an attack order, see if
     *  there are other species present that have declared it as an ALLY.
     *  If so, have the attacker attack the other species and vice-versa. */
    for (species_index = 0; species_index < num_sp; ++species_index) {
        for (i = 0; i < num_sp; i++) {
            if (i == species_index) {
                continue;
            }

            if (!bat.enemy_mine[species_index][i]) {
                continue;
            }

            j           = bat.spec_num[i] - 1;
            array_index = j / 32;
            bit_number  = j % 32;
            bit_mask    = 1 << bit_number;

            for (k = 0; k < num_sp; k++) {
                if (k == species_index) {
                    continue;
                }
                if (k == i) {
                    continue;
                }

                if (c_species[k].ally[array_index] & bit_mask) {
                    /* Make sure it's not already set (it may already be set
                     *  for HIJACK and we don't want to accidentally change
                     *  it to ATTACK). */
                    if (!bat.enemy_mine[species_index][k]) {
                        bat.enemy_mine[species_index][k] = true;
                    }
                    if (!bat.enemy_mine[k][species_index]) {
                        bat.enemy_mine[k][species_index] = true;
                    }
                }
            }
        }
    }

    /* If a species did not give a battle order and is not the target of an
     *  attack, set can_be_surprised flag to a special value. */
    for (species_index = 0; species_index < num_sp; ++species_index) {
        if (!bat.can_be_surprised[species_index]) {
            continue;
        }

        bat.can_be_surprised[species_index] = 55;

        for (i = 0; i < num_sp; i++) {
            if (i == species_index) {
                continue;
            }

            if (!bat.enemy_mine[i][species_index]) {
                continue;
            }

            bat.can_be_surprised[species_index] = true;

            break;
        }
    }

    /* List combatants. */
    for (species_index = 0; species_index < num_sp; ++species_index) {
        species_number = bat.spec_num[species_index];

        log_string("    SP ");
        if (field_distorted[species_index]) {
            log_int(distorted(species_number));
        }else{
            log_string(c_species[species_index].name);
        }
        if (bat.can_be_surprised[species_index]) {
            log_string(" does not appear to be ready for combat.\n");
        }else{
            log_string(" is mobilized and ready for combat.\n");
        }
    }

    /* Check if a declared enemy is being ambushed. */
    for (i = 0; i < num_sp; i++) {
        namp                  = c_nampla[i] - 1;
        num_namplas           = c_species[i].num_namplas;
        bat.ambush_amount[i] = 0;
        for (j = 0; j < num_namplas; j++) {
            ++namp;

            if (namp.x != bat.x) {
                continue;
            }
            if (namp.y != bat.y) {
                continue;
            }
            if (namp.z != bat.z) {
                continue;
            }

            bat.ambush_amount[i] += namp.use_on_ambush;
        }

        if (bat.ambush_amount[i] == 0) {
            continue;
        }

        for (j = 0; j < num_sp; j++) {
            if (bat.enemy_mine[i][j]) {
                do_ambush(i, bat);
            }
        }
    }

    /* For all species that specified enemies, make the feeling mutual. */
    for (i = 0; i < num_sp; i++) {
        for (j = 0; j < num_sp; j++) {
            if (bat.enemy_mine[i][j]) {
                /* Make sure it's not already set (it may already be set for
                 *  HIJACK and we don't want to accidentally change it to
                 *  ATTACK). */
                if (!bat.enemy_mine[j][i]) {
                    bat.enemy_mine[j][i] = true;
                }
            }
        }
    }

    /* Create a sequential list of combat options. First check if a
     *  deep space defense has been ordered. If so, then make sure that
     *  first option is DEEP_SPACE_FIGHT. */
    num_combat_options = 0;
    for (species_index = 0; species_index < num_sp; ++species_index) {
        for (i = 0; i < bat.num_engage_options[species_index]; i++) {
            option = bat.engage_option[species_index][i];
            if (option == DEEP_SPACE_DEFENSE) {
                consolidate_option(DEEP_SPACE_FIGHT, 0);
                goto consolidate;
            }
        }
    }

consolidate:
    for (species_index = 0; species_index < num_sp; ++species_index) {
        for (i = 0; i < bat.num_engage_options[species_index]; i++) {
            option = bat.engage_option[species_index][i];
            where  = bat.engage_planet[species_index][i];
            consolidate_option(option, where);
        }
    }

    /* If ships are given unconditional withdraw orders, they will always have
     *  time to escape if fighting occurs first in a different part of the
     *  sector. The flag "do_withdraw_check_first" will be set only after the
     *  first round of combat. */
    do_withdraw_check_first = false;

    /* Handle each combat option. */
    battle_here  = false;
    first_action = true;
    for (option_index = 0; option_index < num_combat_options; option_index++) {
        option = combat_option[option_index];
        where  = combat_location[option_index];

        /* Fill action arrays with data about ships taking part in current
         * action. */
        fight_here = fighting_params(option, where, bat, &act);

        /* Check if a fight will take place here. */
        if (!fight_here) {
            continue;
        }

        /* See if anyone is taken by surprise. */
        if (!battle_here) {
            /* Combat is just starting. */
            for (species_index = 0; species_index < num_sp; ++species_index) {
                species_number = bat.spec_num[species_index];

                if (bat.can_be_surprised[species_index] == 55) {
                    continue;
                }

                if (bat.can_be_surprised[species_index]) {
                    log_string("\n    SP ");
                    if (field_distorted[species_index]) {
                        log_int(distorted(species_number));
                    }else{
                        log_string(c_species[species_index].name);
                    }
                    log_string(" is taken by surprise!\n");
                }
            }
        }

        battle_here = true;

        /* Clear out can_be_surprised array. */
        for (i = 0; i < MAX_SPECIES; i++) {
            bat.can_be_surprised[i] = false;
        }

        /* Determine maximum number of rounds. */
        max_rounds = 10000;     /* Something ridiculously large. */
        if (option == DEEP_SPACE_FIGHT && attacking_ML > 0 &&
            defending_ML > 0 &&
            deep_space_defense) {
            /* This is the initial deep space fight and the defender wants the
             *  fight to remain in deep space for as long as possible. */
            if (defending_ML > attacking_ML) {
                max_rounds = defending_ML - attacking_ML;
            }else{
                max_rounds = 1;
            }
        }else if (option == PLANET_BOMBARDMENT) {
            /* To determine the effectiveness of the bombardment, we will
             *  simulate ten rounds of combat and add up the damage. */
            max_rounds = 10;
        }else if (option == GERM_WARFARE || option == SIEGE) {
            /* We just need to see who is attacking whom and get the number
             *  of germ warfare bombs being used. */
            max_rounds = 1;
        }

        /* Log start of action. */
        if (where == 0) {
            log_string("\n    The battle begins in deep space, outside the range of planetary defenses...\n");
        }else if (option == PLANET_ATTACK) {
            log_string("\n    The battle ");
            if (first_action) {
                log_string("begins");
            }else{
                log_string("moves");
            }
            log_string(" within range of planet #");
            log_int(where);
            log_string("...\n");
        }else if (option == PLANET_BOMBARDMENT) {
            log_string("\n    Bombardment of planet #");
            log_int(where);
            log_string(" begins...\n");
        }else if (option == GERM_WARFARE) {
            log_string("\n    Germ warfare commences against planet #");
            log_int(where);
            log_string("...\n");
        }else if (option == SIEGE) {
            log_string("\n    Siege of planet #");
            log_int(where);
            log_string(" is now in effect...\n\n");
            goto do_combat;
        }

        /* List combatants. */
        truncate_name = false;
        log_string("\n      Units present:");
        current_species = -1;
        for (unit_index = 0; unit_index < act.num_units_fighting; unit_index++) {
            if (act.fighting_species_index[unit_index] != current_species) {
                /* Display species name. */
                i = act.fighting_species_index[unit_index];
                log_string("\n        SP ");
                species_number = bat.spec_num[i];
                if (field_distorted[i]) {
                    log_int(distorted(species_number));
                }else{
                    log_string(c_species[i].name);
                }
                log_string(": ");
                current_species = i;
                need_comma      = false;
            }

            if (act.unit_ttype[unit_index] == SHIP) {
                sh          = (struct ship_data *)act.fighting_unit[unit_index];
                temp_status = sh.status;
                temp_pn     = sh.pn;
                if (option == DEEP_SPACE_FIGHT) {
                    sh.status = IN_DEEP_SPACE;
                    sh.pn     = 0;
                }else {
                    sh.status = IN_ORBIT;
                    sh.pn     = where;
                }
                ignore_field_distorters = !field_distorted[current_species];
                if (sh.special != NON_COMBATANT) {
                    if (need_comma) {
                        log_string(", ");
                    }
                    log_string(ship_name(sh));
                    need_comma = true;
                }
                ignore_field_distorters = false;
                sh.status = temp_status;
                sh.pn     = temp_pn;
            }else {
                namp = (struct nampla_data *)act.fighting_unit[unit_index];
                if (need_comma) {
                    log_string(", ");
                }
                log_string("PL ");
                log_string(namp.name);
                need_comma = true;
            }
        }
        log_string("\n\n");

do_combat:

        /* Long names are not necessary for the rest of the action. */
        truncate_name = true;

        /* Do combat rounds. Stop if maximum count is reached, or if combat
         *  does not occur when do_round() is called. */

        round_number = 1;

        log_summary = false;    /* do_round() and the routines that it calls
                                *      will set this for important stuff. */

        if (option == PLANET_BOMBARDMENT || option == GERM_WARFARE ||
            option == SIEGE) {
            logging_disabled = true; /* Disable logging during simulation. */
        }
        for (;round_number <= max_rounds;) {
            if (do_withdraw_check_first) {
                withdrawal_check(bat, &act);
            }

            if (!do_round(option, round_number, bat, &act)) {
                break;
            }

            if (!do_withdraw_check_first) {
                withdrawal_check(bat, &act);
            }

            do_withdraw_check_first = true;

            regenerate_shields(&act);

            ++round_number;
        }

        log_summary      = true;
        logging_disabled = false;

        if (round_number == 1) {
            log_string("      ...But it seems that the attackers had nothing to attack!\n");
            continue;
        }

        if (option == PLANET_BOMBARDMENT || option == GERM_WARFARE) {
            for (unit_index = 0; unit_index < act.num_units_fighting; unit_index++) {
                if (act.unit_ttype[unit_index] == GENOCIDE_NAMPLA) {
                    attacked_nampla = (struct nampla_data *)
                                      act.fighting_unit[unit_index];
                    j = act.fighting_species_index[unit_index];
                    for (i = 0; i < num_sp; i++) {
                        if (x_attacked_y[i][j]) {
                            species_number = bat.spec_num[i];
                            log_string("      SP ");
                            if (field_distorted[i]) {
                                log_int(distorted(species_number));
                            }else{
                                log_string(c_species[i].name);
                            }
                            log_string(" bombards SP ");
                            log_string(c_species[j].name);
                            log_string(" on PL ");
                            log_string(attacked_nampla.name);
                            log_string(".\n");

                            if (option == GERM_WARFARE) {
                                do_germ_warfare(i, j, unit_index, bat, &act);
                            }
                        }
                    }

                    /* Determine results of bombardment. */
                    if (option == PLANET_BOMBARDMENT) {
                        do_bombardment(unit_index, &act);
                    }
                }
            }
        }else if (option == SIEGE) {
            do_siege(bat, &act);
        }

        truncate_name = false;

        first_action = false;
    }

    if (!battle_here) {
        if (bat.num_species_here == 1) {
            log_string("    But there was no one to fight with!\n");
        }else if (!ambush_took_place) {
            log_string("    But no one was willing to throw the first punch!\n");
        }
    }

    /* Close combat log and append it to the log files of all species
     *  involved in this battle. */
    if (prompt_gm) {
        printf("\n  End of battle in sector %d, %d, %d.\n", bat.x,
               bat.y, bat.z);
    }
    fprintf(log_file, "\n  End of battle in sector %d, %d, %d.\n", bat.x,
            bat.y, bat.z);
    fprintf(summary_file, "\n  End of battle in sector %d, %d, %d.\n",
            bat.x, bat.y, bat.z);
    fclose(log_file);
    fclose(summary_file);

    for (species_index = 0; species_index < num_sp; ++species_index) {
        species_number = bat.spec_num[species_index];

        /* Open combat log file for reading. */
        if (bat.summary_only[species_index]) {
            combat_log = fopen("summary.log", "r");
        }else{
            combat_log = fopen("combat.log", "r");
        }

        if (combat_log == NULL) {
            fprintf(stderr, "\n\tCannot open combat log for reading!\n\n");
            exit(-1);
        }

        /* Open a temporary species log file for appending. */
        sprintf(filename, "sp%02d.temp.log\0", species_number);
        species_log = fopen(filename, "a");
        if (species_log == NULL) {
            fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
            exit(-1);
        }

        /* Copy combat log to temporary species log. */
        for (;fgets(log_line, 256, combat_log) != NULL;) {
            fputs(log_line, species_log);
        }

        fclose(species_log);
        fclose(combat_log);

        append_log[species_number - 1] = true;

        /* Get rid of ships that were destroyed. */
        if (!data_modified[species_number - 1]) {
            continue;
        }
        sh = c_ship[species_index] - 1;
        for (i = 0; i < c_species[species_index].num_ships; i++) {
            ++sh;

            if (sh.age < 50) {
                continue;
            }
            if (sh.pn == 99) {
                continue;
            }
            if (sh.x != bat.x) {
                continue;
            }
            if (sh.y != bat.y) {
                continue;
            }
            if (sh.z != bat.z) {
                continue;
            }
            if (sh.status == UNDER_CONSTRUCTION) {
                continue;
            }

            delete_ship(sh);
        }
    }
}



func do_ambush(ambushing_species_index int, bat *battle_dat) {
    int i, j, n, num_sp, ambushed_species_index, num_ships,
        age_increment, species_number, old_truncate_name;

    long friendly_tonnage, enemy_tonnage;

    struct ship_data *sh;



    /* Get total ambushing tonnage. */
    friendly_tonnage = 0;
    num_ships        = c_species[ambushing_species_index].num_ships;
    sh = c_ship[ambushing_species_index] - 1;
    for (i = 0; i < num_ships; i++) {
        ++sh;

        if (sh.pn == 99) {
            continue;
        }
        if (sh.x != bat.x) {
            continue;
        }
        if (sh.y != bat.y) {
            continue;
        }
        if (sh.z != bat.z) {
            continue;
        }
        if (sh.class != TR && sh.class != BA) {
            friendly_tonnage += sh.tonnage;
        }
    }

    /* Determine which species are being ambushed and get total enemy
     *  tonnage. */
    num_sp        = bat.num_species_here;
    enemy_tonnage = 0;
    for (ambushed_species_index = 0; ambushed_species_index < num_sp; ++ambushed_species_index) {
        if (!bat.enemy_mine[ambushing_species_index][ambushed_species_index]) {
            continue;
        }

        /* This species is being ambushed.  Get total effective tonnage. */
        num_ships = c_species[ambushed_species_index].num_ships;
        sh        = c_ship[ambushed_species_index] - 1;
        for (i = 0; i < num_ships; i++) {
            ++sh;

            if (sh.pn == 99) {
                continue;
            }
            if (sh.x != bat.x) {
                continue;
            }
            if (sh.y != bat.y) {
                continue;
            }
            if (sh.z != bat.z) {
                continue;
            }
            if (sh.class == TR) {
                enemy_tonnage += sh.tonnage;
            }else{
                enemy_tonnage += 10 * sh.tonnage;
            }
        }
    }

    /* Determine the amount of aging that will be added to each ambushed
     *  ship. */
    if (enemy_tonnage == 0) {
        return;
    }
    age_increment = (10L * bat.ambush_amount[ambushing_species_index])
                    / enemy_tonnage;
    age_increment = (friendly_tonnage * age_increment) / enemy_tonnage;

    ambush_took_place = true;

    if (age_increment < 1) {
        log_string("\n    SP ");
        log_string(c_species[ambushing_species_index].name);
        log_string(" attempted an ambush, but the ambush was completely ineffective!\n");
        return;
    }

    /* Age each ambushed ship. */
    for (ambushed_species_index = 0; ambushed_species_index < num_sp; ++ambushed_species_index) {
        if (!bat.enemy_mine[ambushing_species_index][ambushed_species_index]) {
            continue;
        }

        log_string("\n    SP ");
        species_number = bat.spec_num[ambushed_species_index];
        if (field_distorted[ambushed_species_index]) {
            log_int(distorted(species_number));
        }else{
            log_string(c_species[ambushed_species_index].name);
        }

        log_string(" was ambushed by SP ");
        log_string(c_species[ambushing_species_index].name);
        log_string("!\n");

        num_ships = c_species[ambushed_species_index].num_ships;
        sh        = c_ship[ambushed_species_index] - 1;
        for (i = 0; i < num_ships; i++) {
            ++sh;

            if (sh.pn == 99) {
                continue;
            }
            if (sh.x != bat.x) {
                continue;
            }
            if (sh.y != bat.y) {
                continue;
            }
            if (sh.z != bat.z) {
                continue;
            }

            sh.age += age_increment;
            if (sh.arrived_via_wormhole) {
                sh.age += age_increment;
            }

            if (sh.age > 49) {
                old_truncate_name = truncate_name;
                truncate_name     = true;

                log_string("      ");
                log_string(ship_name(sh));
                if (field_distorted[ambushed_species_index]) {
                    log_string(" = ");
                    log_string(c_species[ambushed_species_index].name);
                    log_char(' ');
                    n = sh.item_quantity[FD];
                    sh.item_quantity[FD] = 0;
                    log_string(ship_name(sh));
                    sh.item_quantity[FD] = n;
                }
                n = 0;
                for (j = 0; j < MAX_ITEMS; j++) {
                    if (sh.item_quantity[j] > 0) {
                        if (n++ == 0) {
                            log_string(" (cargo: ");
                        }else{
                            log_char(',');
                        }
                        log_int((int)sh.item_quantity[j]);
                        log_char(' ');
                        log_string(item_abbr[j]);
                    }
                }
                if (n > 0) {
                    log_char(')');
                }

                log_string(" was destroyed in the ambush!\n");

                truncate_name = old_truncate_name;
            }
        }
    }
}



/* This routine will find all species that have declared alliance with
 * both a traitor and betrayed species. It will then set a flag to indicate
 * that their allegiance should be changed from ALLY to ENEMY. */

func auto_enemy(traitor_species_number, betrayed_species_number int) {
    int traitor_array_index, betrayed_array_index, bit_number,
        species_index;

    long traitor_bit_mask, betrayed_bit_mask;


    traitor_array_index = (traitor_species_number - 1) / 32;
    bit_number          = (traitor_species_number - 1) % 32;
    traitor_bit_mask    = 1 << bit_number;

    betrayed_array_index = (betrayed_species_number - 1) / 32;
    bit_number           = (betrayed_species_number - 1) % 32;
    betrayed_bit_mask    = 1 << bit_number;

    for (species_index = 0; species_index < galaxy.num_species; species_index++) {
        if ((spec_data[species_index].ally[traitor_array_index]
             & traitor_bit_mask) == 0) {
            continue;
        }
        if ((spec_data[species_index].ally[betrayed_array_index]
             & betrayed_bit_mask) == 0) {
            continue;
        }
        if ((spec_data[species_index].contact[traitor_array_index]
             & traitor_bit_mask) == 0) {
            continue;
        }
        if ((spec_data[species_index].contact[betrayed_array_index]
             & betrayed_bit_mask) == 0) {
            continue;
        }

        make_enemy[species_index][traitor_species_number - 1] = betrayed_species_number;
    }
}

//*************************************************************************
// do_bomb.c


func do_bombardment(unit_index int, act *action_data) {
    int i, new_mi, new_ma, defending_species;

    long n, total_bomb_damage, CS_bomb_damage, new_pop,
         initial_base, total_pop, percent_damage;

    struct planet_data *planet;
    struct nampla_data *attacked_nampla;
    struct ship_data *  sh;


    attacked_nampla = (struct nampla_data *)act.fighting_unit[unit_index];
    planet          = planet_base + (long)attacked_nampla.planet_index;

    initial_base = attacked_nampla.mi_base + attacked_nampla.ma_base;
    total_pop    = initial_base;

    if (attacked_nampla.item_quantity[CU] > 0) {
        total_pop += 1;
    }

    if (total_pop < 1) {
        log_string("        The planet is completely uninhabited. There is nothing to bomb!\n");
        return;
    }

    /* Total damage done by ten strike cruisers (ML = 50) in ten rounds
     *  is 100 x 4 x the power value for a single ship. To eliminate the
     *  chance of overflow, the algorithm has been carefully chosen. */

    CS_bomb_damage = 400 * power(ship_tonnage[CS]);
    /* Should be 400 * 4759 = 1,903,600. */

    total_bomb_damage = act.bomb_damage[unit_index];

    /* Keep about 2 significant digits. */
    for (;total_bomb_damage > 1000;) {
        total_bomb_damage /= 10;
        CS_bomb_damage    /= 10;
    }

    if (CS_bomb_damage == 0) {
        percent_damage = 101;
    }else{
        percent_damage =
            ((total_bomb_damage * 250000L) / CS_bomb_damage) / total_pop;
    }

    if (percent_damage > 100) {
        percent_damage = 101;
    }

    new_mi = attacked_nampla.mi_base
             - (percent_damage * attacked_nampla.mi_base) / 100;

    new_ma = attacked_nampla.ma_base
             - (percent_damage * attacked_nampla.ma_base) / 100;

    new_pop = attacked_nampla.pop_units
              - (percent_damage * attacked_nampla.pop_units) / 100;

    if (new_mi == attacked_nampla.mi_base &&
        new_ma == attacked_nampla.ma_base &&
        new_pop == attacked_nampla.pop_units) {
        log_string("        Damage due to bombardment was insignificant.\n");
        return;
    }

    defending_species = act.fighting_species_index[unit_index];
    if (attacked_nampla.status & HOME_PLANET) {
        n = attacked_nampla.mi_base + attacked_nampla.ma_base;
        if (c_species[defending_species].hp_original_base < n) {
            c_species[defending_species].hp_original_base = n;
        }
    }

    if (new_mi <= 0 && new_ma <= 0 && new_pop <= 0) {
        log_string("        Everyone and everything was completely wiped out!\n");

        attacked_nampla.mi_base       = 0;
        attacked_nampla.ma_base       = 0;
        attacked_nampla.pop_units     = 0;
        attacked_nampla.siege_eff     = 0;
        attacked_nampla.shipyards     = 0;
        attacked_nampla.hiding        = 0;
        attacked_nampla.hidden        = 0;
        attacked_nampla.use_on_ambush = 0;

        /* Reset status. */
        if (attacked_nampla.status & HOME_PLANET) {
            attacked_nampla.status = HOME_PLANET;
        }else{
            attacked_nampla.status = COLONY;
        }

        for (i = 0; i < MAX_ITEMS; i++) {
            attacked_nampla.item_quantity[i] = 0;
        }

        /* Delete any ships that were under construction on the planet. */
        sh = c_ship[defending_species] - 1;
        for (i = 0; i < c_species[defending_species].num_ships; i++) {
            ++sh;

            if (sh.x != attacked_nampla.x) {
                continue;
            }
            if (sh.y != attacked_nampla.y) {
                continue;
            }
            if (sh.z != attacked_nampla.z) {
                continue;
            }
            if (sh.pn != attacked_nampla.pn) {
                continue;
            }

            delete_ship(sh);
        }

        return;
    }

    log_string("        Mining base of PL ");
    log_string(attacked_nampla.name);
    log_string(" went from ");
    log_int(attacked_nampla.mi_base / 10);
    log_char('.');
    log_int(attacked_nampla.mi_base % 10);
    log_string(" to ");
    attacked_nampla.mi_base = new_mi;
    log_int(new_mi / 10);
    log_char('.');
    log_int(new_mi % 10);
    log_string(".\n");

    log_string("        Manufacturing base of PL ");
    log_string(attacked_nampla.name);
    log_string(" went from ");
    log_int(attacked_nampla.ma_base / 10);
    log_char('.');
    log_int(attacked_nampla.ma_base % 10);
    log_string(" to ");
    attacked_nampla.ma_base = new_ma;
    log_int(new_ma / 10);
    log_char('.');
    log_int(new_ma % 10);
    log_string(".\n");

    attacked_nampla.pop_units = new_pop;

    for (i = 0; i < MAX_ITEMS; i++) {
        n = (percent_damage * attacked_nampla.item_quantity[i]) / 100;
        if (n > 0) {
            attacked_nampla.item_quantity[i] -= n;
            log_string("        ");    log_long(n);    log_char(' ');
            log_string(item_name[i]);
            if (n > 1) {
                log_string("s were");
            }else{
                log_string(" was");
            }
            log_string(" destroyed.\n");
        }
    }

    n = (percent_damage * (long)attacked_nampla.shipyards) / 100;
    if (n > 0) {
        attacked_nampla.shipyards -= n;
        log_string("        ");    log_long(n);
        log_string(" shipyard");
        if (n > 1) {
            log_string("s were");
        }else{
            log_string(" was");
        }
        log_string(" also destroyed.\n");
    }

    check_population(attacked_nampla);
}



//*************************************************************************
// do_build.c

func do_BUILD_command(continuing_construction, interspecies_construction int) {
    int i, n, class, critical_tech, found, name_length,
        siege_effectiveness, cost_given, new_ship, max_tonnage,
        tonnage_increase, alien_number, cargo_on_board,
        unused_nampla_available, unused_ship_available, capacity,
        pop_check_needed, contact_word_number, contact_bit_number,
        already_notified[MAX_SPECIES];

    char upper_ship_name[32], *commas(), *src, *dest,
         *original_line_pointer;

    long cost, cost_argument, unit_cost, num_items, pop_reduction,
         premium, total_cost, original_num_items, contact_mask,
         max_funds_available;

    struct species_data *recipient_species;
    struct nampla_data * recipient_nampla, *unused_nampla,
                       *destination_nampla, *temp_nampla;
    struct ship_data *recipient_ship, *unused_ship;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get ready if planet is under siege. */
    if (nampla.siege_eff < 0) {
        siege_effectiveness = -nampla.siege_eff;
    }else{
        siege_effectiveness = nampla.siege_eff;
    }

    /* Get species name and make appropriate tests if this is an interspecies
     *  construction order. */
    if (interspecies_construction) {
        original_line_pointer = input_line_pointer;
        if (!get_species_name()) {
            /* Check for missing comma or tab after species name. */
            input_line_pointer = original_line_pointer;
            fix_separator();
            if (!get_species_name()) {
                fprintf(log_file, "!!! Order ignored:\n");
                fprintf(log_file, "!!! %s", original_line);
                fprintf(log_file, "!!! Invalid species name.\n");
                return;
            }
        }
        recipient_species = &spec_data[g_spec_number - 1];

        if (species.tech_level[MA] < 25) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! MA tech level must be at least 25 to do interspecies construction.\n");
            return;
        }

        /* Check if we've met this species and make sure it is not an enemy. */
        contact_word_number = (g_spec_number - 1) / 32;
        contact_bit_number  = (g_spec_number - 1) % 32;
        contact_mask        = 1 << contact_bit_number;
        if ((species.contact[contact_word_number] & contact_mask) == 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! You can't do interspecies construction for a species you haven't met.\n");
            return;
        }
        if (species.enemy[contact_word_number] & contact_mask) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! You can't do interspecies construction for an ENEMY.\n");
            return;
        }
    }

    /* Get number of items to build. */
    i = get_value();

    if (i == 0) {
        goto build_ship;        /* Not an item. */
    }
    num_items          = value;
    original_num_items = value;

    /* Get class of item. */
    class = get_class_abbr();

    if (class != ITEM_CLASS || abbr_index == RM) {
        /* Players sometimes accidentally use "MI" for "IU"
         *      or "MA" for "AU". */
        if (class == TECH_ID && abbr_index == MI) {
            abbr_index = IU;
        }else if (class == TECH_ID && abbr_index == MA) {
            abbr_index = AU;
        }else{
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid item class.\n");
            return;
        }
    }
    class = abbr_index;

    if (interspecies_construction) {
        if (class == PD || class == CU) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! You cannot build CUs or PDs for another species.\n");
            return;
        }
    }

    /* Make sure species knows how to build this item. */
    critical_tech = item_critical_tech[class];
    if (species.tech_level[critical_tech] < item_tech_requirment[class]) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Insufficient tech level to build item.\n");
        return;
    }

    /* Get cost of item. */
    if (class == TP) {  /* Terraforming plant. */
        unit_cost = item_cost[class] / species.tech_level[critical_tech];
    }else{
        unit_cost = item_cost[class];
    }

    if (num_items == 0) {
        num_items = balance / unit_cost;
    }
    if (num_items == 0) {
        return;
    }

    /* Make sure item count is meaningful. */
    if (num_items < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Meaningless item count.\n");
        return;
    }

    /* Make sure there is enough available population. */
    pop_reduction = 0;
    if (class == CU || class == PD) {
        if (nampla.pop_units < num_items) {
            if (original_num_items == 0) {
                num_items = nampla.pop_units;
                if (num_items == 0) {
                    return;
                }
            }else {
                if (nampla.pop_units > 0) {
                    fprintf(log_file, "! WARNING: %s", original_line);
                    fprintf(log_file,
                            "! Insufficient available population units. Substituting %ld for %ld.\n",
                            nampla.pop_units, num_items);
                    num_items = nampla.pop_units;
                }else {
                    fprintf(log_file, "!!! Order ignored:\n");
                    fprintf(log_file, "!!! %s", original_line);
                    fprintf(log_file, "!!! Insufficient available population units.\n");
                    return;
                }
            }
        }
        pop_reduction = num_items;
    }

    /* Calculate total cost and see if planet has enough money. */
do_cost:
    cost = num_items * unit_cost;
    if (interspecies_construction) {
        premium = (cost + 9) / 10;
    }else{
        premium = 0;
    }

    cost += premium;

    if (check_bounced(cost)) {
        if (interspecies_construction && original_num_items == 0) {
            --num_items;
            if (num_items < 1) {
                return;
            }
            goto do_cost;
        }

        max_funds_available = species.econ_units;
        if (max_funds_available > EU_spending_limit) {
            max_funds_available = EU_spending_limit;
        }
        max_funds_available += balance;

        num_items = max_funds_available / unit_cost;
        if (interspecies_construction) {
            num_items -= (num_items + 9) / 10;
        }

        if (num_items > 0) {
            fprintf(log_file, "! WARNING: %s", original_line);
            fprintf(log_file, "! Insufficient funds. Substituting %ld for %ld.\n",
                    num_items, original_num_items);
            goto do_cost;
        }else {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Insufficient funds to execute order.\n");
            return;
        }
    }

    /* Update planet inventory. */
    nampla.item_quantity[class] += num_items;
    nampla.pop_units            -= pop_reduction;

    /* Log what was produced. */
    log_string("    ");  log_long(num_items);
    log_char(' ');  log_string(item_name[class]);

    if (num_items > 1) {
        log_string("s were");
    }else{
        log_string(" was");
    }

    if (first_pass && class == PD && siege_effectiveness > 0) {
        log_string(" scheduled for production despite the siege.\n");
        return;
    }else {
        log_string(" produced");
        if (interspecies_construction) {
            log_string(" for SP ");
            log_string(recipient_species.name);
        }
    }

    if (unit_cost != 1 || premium != 0) {
        log_string(" at a cost of ");
        log_long(cost);
    }

    /* Check if planet is under siege and if production of planetary
     *  defenses was detected. */
    if (class == PD && rnd(100) <= siege_effectiveness) {
        log_string(". However, they were detected and destroyed by the besiegers!!!\n");
        nampla.item_quantity[PD] = 0;

        /* Make sure we don't notify the same species more than once. */
        for (i = 0; i < MAX_SPECIES; i++) {
            already_notified[i] = false;
        }

        for (i = 0; i < num_transactions; i++) {
            /* Find out who is besieging this planet. */
            if (transaction[i].ttype != BESIEGE_PLANET) {
                continue;
            }
            if (transaction[i].x != nampla.x) {
                continue;
            }
            if (transaction[i].y != nampla.y) {
                continue;
            }
            if (transaction[i].z != nampla.z) {
                continue;
            }
            if (transaction[i].pn != nampla.pn) {
                continue;
            }
            if (transaction[i].number2 != species_number) {
                continue;
            }

            alien_number = transaction[i].number1;

            if (already_notified[alien_number - 1]) {
                continue;
            }

            /* Define a 'detection' transaction. */
            if (num_transactions == MAX_TRANSACTIONS) {
                fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
                exit(-1);
            }

            n = num_transactions++;
            transaction[n].ttype  = DETECTION_DURING_SIEGE;
            transaction[n].value = 3;   /* Construction of PDs. */
            strcpy(transaction[n].name1, nampla.name);
            strcpy(transaction[n].name3, species.name);
            transaction[n].number3 = alien_number;

            already_notified[alien_number - 1] = true;
        }
        return;
    }

    if (!interspecies_construction) {
        /* Get destination of transfer, if any. */
        pop_check_needed   = false;
        temp_nampla        = nampla;
        found              = get_transfer_point();
        destination_nampla = nampla;
        nampla             = temp_nampla;
        if (!found) {
            goto done_transfer;
        }

        if (abbr_type == SHIP_CLASS) {  /* Destination is 'ship'. */
            if (ship.x != nampla.x ||
                ship.y != nampla.y ||
                ship.z != nampla.z ||
                ship.status == UNDER_CONSTRUCTION) {
                goto done_transfer;
            }

            if (ship.class == TR) {
                capacity = (10 + ((int)ship.tonnage / 2)) * (int)ship.tonnage;
            }else if (ship.class == BA) {
                capacity = 10 * ship.tonnage;
            }else{
                capacity = ship.tonnage;
            }

            for (i = 0; i < MAX_ITEMS; i++) {
                capacity -= ship.item_quantity[i] * item_carry_capacity[i];
            }

            n = num_items;
            if (num_items * item_carry_capacity[class] > capacity) {
                num_items = capacity / item_carry_capacity[class];
            }

            ship.item_quantity[class]   += num_items;
            nampla.item_quantity[class] -= num_items;
            log_string(" and ");
            if (n > num_items) {
                log_long(num_items);
                log_string(" of them ");
            }
            if (num_items == 1) {
                log_string("was");
            }else{
                log_string("were");
            }
            log_string(" transferred to ");
            log_string(ship_name(ship));

            if (class == CU && num_items > 0) {
                if (nampla == nampla_base) {
                    ship.loading_point = 9999; /* Home planet. */
                }else {
                    ship.loading_point = (nampla - nampla_base);
                }
            }
        }else {  /* Destination is 'destination_nampla'. */
            if (destination_nampla.x != nampla.x ||
                destination_nampla.y != nampla.y ||
                destination_nampla.z != nampla.z) {
                goto done_transfer;
            }

            if (nampla.siege_eff != 0) {
                goto done_transfer;
            }
            if (destination_nampla.siege_eff != 0) {
                goto done_transfer;
            }

            destination_nampla.item_quantity[class] += num_items;
            nampla.item_quantity[class]             -= num_items;
            log_string(" and transferred to PL ");
            log_string(destination_nampla.name);
            pop_check_needed = true;
        }

done_transfer:

        log_string(".\n");

        if (pop_check_needed) {
            check_population(destination_nampla);
        }

        return;
    }

    log_string(".\n");

    /* Check if recipient species has a nampla at this location. */
    found = false;
    unused_nampla_available = false;
    recipient_nampla        = namp_data[g_spec_number - 1] - 1;
    for (i = 0; i < recipient_species.num_namplas; i++) {
        ++recipient_nampla;

        if (recipient_nampla.pn == 99) {
            unused_nampla           = recipient_nampla;
            unused_nampla_available = true;
        }

        if (recipient_nampla.x != nampla.x) {
            continue;
        }
        if (recipient_nampla.y != nampla.y) {
            continue;
        }
        if (recipient_nampla.z != nampla.z) {
            continue;
        }
        if (recipient_nampla.pn != nampla.pn) {
            continue;
        }

        found = true;
        break;
    }

    if (!found) {
        /* Add new nampla to database for the recipient species. */
        if (unused_nampla_available) {
            recipient_nampla = unused_nampla;
        }else{
            ++num_new_namplas[species_index];
            if (num_new_namplas[species_index] > NUM_EXTRA_NAMPLAS) {
                fprintf(stderr, "\n\n\tInsufficient memory for new planet name in do_build.c!\n");
                exit(-1);
            }
            recipient_nampla = namp_data[g_spec_number - 1]
                               + recipient_species.num_namplas;
            recipient_species.num_namplas += 1;
            delete_nampla(recipient_nampla);    /* Set everything to zero. */
        }

        /* Initialize new nampla. */
        strcpy(recipient_nampla.name, nampla.name);
        recipient_nampla.x            = nampla.x;
        recipient_nampla.y            = nampla.y;
        recipient_nampla.z            = nampla.z;
        recipient_nampla.pn           = nampla.pn;
        recipient_nampla.planet_index = nampla.planet_index;
        recipient_nampla.status       = COLONY;
    }

    /* Transfer the goods. */
    nampla.item_quantity[class]           -= num_items;
    recipient_nampla.item_quantity[class] += num_items;
    data_modified[g_spec_number - 1]        = true;

    if (first_pass) {
        return;
    }

    /* Define transaction so that recipient will be notified. */
    if (num_transactions == MAX_TRANSACTIONS) {
        fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
        exit(-1);
    }

    n = num_transactions++;
    transaction[n].ttype      = INTERSPECIES_CONSTRUCTION;
    transaction[n].donor     = species_number;
    transaction[n].recipient = g_spec_number;
    transaction[n].value     = 1; /* Items, not ships. */
    transaction[n].number1   = num_items;
    transaction[n].number2   = class;
    transaction[n].number3   = cost;
    strcpy(transaction[n].name1, species.name);
    strcpy(transaction[n].name2, recipient_nampla.name);

    return;


build_ship:

    original_line_pointer = input_line_pointer;
    if (continuing_construction) {
        found = get_ship();
        if (!found) {
            /* Check for missing comma or tab after ship name. */
            input_line_pointer = original_line_pointer;
            fix_separator();
            found = get_ship();
        }

        if (found) {
            goto check_ship;
        }
        input_line_pointer = original_line_pointer;
    }

    class = get_class_abbr();

    if (class != SHIP_CLASS || tonnage < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Invalid ship class.\n");
        return;
    }
    class = abbr_index;

    /* Get ship name. */
    name_length = get_name();
    if (name_length < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Invalid ship name.\n");
        return;
    }

    /* Search all ships for name. */
    found = false;
    ship  = ship_base - 1;
    unused_ship_available = false;
    for (ship_index = 0; ship_index < species.num_ships; ship_index++) {
        ++ship;

        if (ship.pn == 99) {
            unused_ship_available = true;
            unused_ship           = ship;
            continue;
        }

        /* Make upper case copy of ship name. */
        for (i = 0; i < 32; i++) {
            upper_ship_name[i] = toupper(ship.name[i]);
        }

        /* Compare names. */
        if (strcmp(upper_ship_name, upper_name) == 0) {
            found = true;
            break;
        }
    }

check_ship:

    if (found) {
        /* Check if BUILD was accidentally used instead of CONTINUE. */
        if ((ship.status == UNDER_CONSTRUCTION || ship.ttype == STARBASE) &&
            ship.x == nampla.x && ship.y == nampla.y &&
            ship.z == nampla.z && ship.pn == nampla.pn) {
            continuing_construction = true;
        }

        if ((ship.status != UNDER_CONSTRUCTION && ship.ttype != STARBASE) ||
            (!continuing_construction)) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship name already in use.\n");
            return;
        }

        new_ship = false;
    }else {
        /* If CONTINUE command was used, the player probably mis-spelled
         *  the name. */
        if (continuing_construction) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid ship name.\n");
            return;
        }

        if (unused_ship_available) {
            ship = unused_ship;
        }else{
            /* Make sure we have enough memory for new ship. */
            if (num_new_ships[species_index] >= NUM_EXTRA_SHIPS) {
                if (num_new_ships[species_index] == 9999) {
                    return;
                }

                fprintf(log_file, "!!! Order ignored:\n");
                fprintf(log_file, "!!! %s", original_line);
                fprintf(log_file,
                        "!!! You cannot build more than %d ships per turn!\n",
                        NUM_EXTRA_SHIPS);
                num_new_ships[species_index] = 9999;
                return;
            }
            new_ship = true;
            ship     = ship_base + (int)species.num_ships;
            delete_ship(ship);          /* Initialize everything to zero. */
        }

        /* Initialize non-zero data for new ship. */
        strcpy(ship.name, original_name);
        ship.x      = nampla.x;
        ship.y      = nampla.y;
        ship.z      = nampla.z;
        ship.pn     = nampla.pn;
        ship.status = UNDER_CONSTRUCTION;
        if (class == BA) {
            ship.ttype   = STARBASE;
            ship.status = IN_ORBIT;
        }else if (sub_light) {
            ship.ttype = SUB_LIGHT;
        }else{
            ship.ttype = FTL;
        }
        ship.class = class;
        ship.age   = -1;
        if (ship.ttype != STARBASE) {
            ship.tonnage = tonnage;
        }
        ship.remaining_cost = ship_cost[class];
        if (ship.class == TR) {
            ship.remaining_cost = ship_cost[TR] * tonnage;
        }
        if (ship.ttype == SUB_LIGHT) {
            ship.remaining_cost = (3L * (long)ship.remaining_cost) / 4L;
        }
        ship.just_jumped = false;

        /* Everything else was set to zero in above call to 'delete_ship'. */
    }

    /* Check if amount to spend was specified. */
    cost_given    = get_value();
    cost          = value;
    cost_argument = value;

    if (cost_given) {
        if (interspecies_construction && (ship.ttype != STARBASE)) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Amount to spend may not be specified.\n");
            return;
        }

        if (cost == 0) {
            cost = balance;
            if (ship.ttype == STARBASE) {
                if (cost % ship_cost[BA] != 0) {
                    cost = ship_cost[BA] * (cost / ship_cost[BA]);
                }
            }
            if (cost < 1) {
                if (new_ship) {
                    delete_ship(ship);
                }
                return;
            }
        }

        if (cost < 1) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Amount specified is meaningless.\n");
            if (new_ship) {
                delete_ship(ship);
            }
            return;
        }

        if (ship.ttype == STARBASE) {
            if (cost % ship_cost[BA] != 0) {
                fprintf(log_file, "!!! Order ignored:\n");
                fprintf(log_file, "!!! %s", original_line);
                fprintf(log_file, "!!! Amount spent on starbase must be multiple of %d.\n",
                        ship_cost[BA]);
                if (new_ship) {
                    delete_ship(ship);
                }
                return;
            }
        }
    }else {
        if (ship.ttype == STARBASE) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Amount to spend MUST be specified for starbase.\n");
            if (new_ship) {
                delete_ship(ship);
            }
            return;
        }

        cost = ship.remaining_cost;
    }

    /* Make sure species can build a ship of this size. */
    max_tonnage = species.tech_level[MA] / 2;
    if (ship.ttype == STARBASE) {
        tonnage_increase = cost / (long)ship_cost[BA];
        tonnage          = ship.tonnage + tonnage_increase;
        if (tonnage > max_tonnage && cost_argument == 0) {
            tonnage_increase = max_tonnage - ship.tonnage;
            if (tonnage_increase < 1) {
                return;
            }
            tonnage = ship.tonnage + tonnage_increase;
            cost    = tonnage_increase * (int)ship_cost[BA];
        }
    }

    if (tonnage > max_tonnage) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Maximum allowable tonnage exceeded.\n");
        if (new_ship) {
            delete_ship(ship);
        }
        return;
    }

    /* Make sure species has gravitics technology if this is an FTL ship. */
    if (ship.ttype == FTL && species.tech_level[GV] < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Gravitics tech needed to build FTL ship!\n");
        if (new_ship) {
            delete_ship(ship);
        }
        return;
    }

    /* Make sure amount specified is not an overpayment. */
    if (ship.ttype != STARBASE && cost > ship.remaining_cost) {
        cost = ship.remaining_cost;
    }

    /* Make sure planet has sufficient shipyards. */
    if (shipyard_capacity < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Shipyard capacity exceeded!\n");
        if (new_ship) {
            delete_ship(ship);
        }
        return;
    }

    /* Make sure there is enough money to pay for it. */
    premium = 0;
    if (interspecies_construction) {
        if (ship.class == TR || ship.ttype == STARBASE) {
            total_cost = ship_cost[ship.class] * tonnage;
        }else{
            total_cost = ship_cost[ship.class];
        }

        if (ship.ttype == SUB_LIGHT) {
            total_cost = (3 * total_cost) / 4;
        }

        premium = total_cost / 10;
        if (total_cost % 10) {
            ++premium;
        }
    }

    if (check_bounced(cost + premium)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        if (new_ship) {
            delete_ship(ship);
        }
        return;
    }

    --shipyard_capacity;

    /* Test if this is a starbase and if planet is under siege. */
    if (ship.ttype == STARBASE && siege_effectiveness > 0) {
        log_string("    Your attempt to build ");
        log_string(ship_name(ship));
        log_string(" was detected by the besiegers and the starbase was destroyed!!!\n");

        /* Make sure we don't notify the same species more than once. */
        for (i = 0; i < MAX_SPECIES; i++) {
            already_notified[i] = false;
        }

        for (i = 0; i < num_transactions; i++) {
            /* Find out who is besieging this planet. */
            if (transaction[i].ttype != BESIEGE_PLANET) {
                continue;
            }
            if (transaction[i].x != nampla.x) {
                continue;
            }
            if (transaction[i].y != nampla.y) {
                continue;
            }
            if (transaction[i].z != nampla.z) {
                continue;
            }
            if (transaction[i].pn != nampla.pn) {
                continue;
            }
            if (transaction[i].number2 != species_number) {
                continue;
            }

            alien_number = transaction[i].number1;

            if (already_notified[alien_number - 1]) {
                continue;
            }

            /* Define a 'detection' transaction. */
            if (num_transactions == MAX_TRANSACTIONS) {
                fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
                exit(-1);
            }

            n = num_transactions++;
            transaction[n].ttype  = DETECTION_DURING_SIEGE;
            transaction[n].value = 2;   /* Construction of ship/starbase. */
            strcpy(transaction[n].name1, nampla.name);
            strcpy(transaction[n].name2, ship_name(ship));
            strcpy(transaction[n].name3, species.name);
            transaction[n].number3 = alien_number;

            already_notified[alien_number - 1] = true;
        }

        delete_ship(ship);

        return;
    }

    /* Finish up and log results. */
    log_string("    ");
    if (ship.ttype == STARBASE) {
        if (ship.tonnage == 0) {
            log_string(ship_name(ship));
            log_string(" was constructed");
        }else {
            ship.age =         /* Weighted average. */
                        ((ship.age * ship.tonnage) - tonnage_increase)
                        / tonnage;
            log_string("Size of ");  log_string(ship_name(ship));
            log_string(" was increased to ");
            log_string(commas(10000L * (long)tonnage));
            log_string(" tons");
        }

        ship.tonnage = tonnage;
    }else {
        ship.remaining_cost -= cost;
        if (ship.remaining_cost == 0) {
            ship.status = ON_SURFACE;  /* Construction is complete. */
            if (continuing_construction) {
                if (first_pass && siege_effectiveness > 0) {
                    log_string("An attempt will be made to finish construction on ");
                }else{
                    log_string("Construction finished on ");
                }
                log_string(ship_name(ship));
                if (first_pass && siege_effectiveness > 0) {
                    log_string(" despite the siege");
                }
            }else {
                if (first_pass && siege_effectiveness > 0) {
                    log_string("An attempt will be made to construct ");
                }
                log_string(ship_name(ship));
                if (first_pass && siege_effectiveness > 0) {
                    log_string(" despite the siege");
                }else{
                    log_string(" was constructed");
                }
            }
        }else {
            if (continuing_construction) {
                if (first_pass && siege_effectiveness > 0) {
                    log_string("An attempt will be made to continue construction on ");
                }else{
                    log_string("Construction continued on ");
                }
                log_string(ship_name(ship));
                if (first_pass && siege_effectiveness > 0) {
                    log_string(" despite the siege");
                }
            }else {
                if (first_pass && siege_effectiveness > 0) {
                    log_string("An attempt will be made to start construction on ");
                }else{
                    log_string("Construction started on ");
                }
                log_string(ship_name(ship));
                if (first_pass && siege_effectiveness > 0) {
                    log_string(" despite the siege");
                }
            }
        }
    }
    log_string(" at a cost of ");    log_long(cost + premium);

    if (interspecies_construction) {
        log_string(" for SP ");
        log_string(recipient_species.name);
    }

    log_char('.');

    if (new_ship && (!unused_ship_available)) {
        ++num_new_ships[species_index];
        ++species.num_ships;
    }

    /* Check if planet is under siege and if construction was detected. */
    if (!first_pass && rnd(100) <= siege_effectiveness) {
        log_string(" However, the work was detected by the besiegers and the ship was destroyed!!!");

        /* Make sure we don't notify the same species more than once. */
        for (i = 0; i < MAX_SPECIES; i++) {
            already_notified[i] = false;
        }

        for (i = 0; i < num_transactions; i++) {
            /* Find out who is besieging this planet. */
            if (transaction[i].ttype != BESIEGE_PLANET) {
                continue;
            }
            if (transaction[i].x != nampla.x) {
                continue;
            }
            if (transaction[i].y != nampla.y) {
                continue;
            }
            if (transaction[i].z != nampla.z) {
                continue;
            }
            if (transaction[i].pn != nampla.pn) {
                continue;
            }
            if (transaction[i].number2 != species_number) {
                continue;
            }

            alien_number = transaction[i].number1;

            if (already_notified[alien_number - 1]) {
                continue;
            }

            /* Define a 'detection' transaction. */
            if (num_transactions == MAX_TRANSACTIONS) {
                fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
                exit(-1);
            }

            n = num_transactions++;
            transaction[n].ttype  = DETECTION_DURING_SIEGE;
            transaction[n].value = 2;   /* Construction of ship/starbase. */
            strcpy(transaction[n].name1, nampla.name);
            strcpy(transaction[n].name2, ship_name(ship));
            strcpy(transaction[n].name3, species.name);
            transaction[n].number3 = alien_number;

            already_notified[alien_number - 1] = true;
        }

        /* Remove ship from inventory. */
        delete_ship(ship);
    }

    log_char('\n');

    if (!interspecies_construction) {
        return;
    }

    /* Transfer any cargo on the ship to the planet. */
    cargo_on_board = false;
    for (i = 0; i < MAX_ITEMS; i++) {
        if (ship.item_quantity[i] > 0) {
            nampla.item_quantity[i] += ship.item_quantity[i];
            ship.item_quantity[i]    = 0;
            cargo_on_board            = true;
        }
    }
    if (cargo_on_board) {
        log_string("      Forgotten cargo on the ship was first transferred to the planet.\n");
    }

    /* Transfer the ship to the recipient species. */
    unused_ship_available = false;
    recipient_ship        = ship_data[g_spec_number - 1];
    for (i = 0; i < recipient_species.num_ships; i++) {
        if (recipient_ship.pn == 99) {
            unused_ship_available = true;
            break;
        }

        ++recipient_ship;
    }

    if (!unused_ship_available) {
        /* Make sure we have enough memory for new ship. */
        if (num_new_ships[g_spec_number - 1] == NUM_EXTRA_SHIPS) {
            fprintf(stderr, "\n\n\tInsufficient memory for new recipient ship!\n\n");
            exit(-1);
        }
        recipient_ship = ship_data[g_spec_number - 1]
                         + (int)recipient_species.num_ships;
        ++recipient_species.num_ships;
        ++num_new_ships[g_spec_number - 1];
    }

    /* Copy donor ship to recipient ship. */
    src  = (char *)ship;
    dest = (char *)recipient_ship;
    for (i = 0; i < sizeof(struct ship_data); i++) {
        *dest++ = *src++;
    }

    recipient_ship.status = IN_ORBIT;

    data_modified[g_spec_number - 1] = true;

    /* Delete donor ship. */
    delete_ship(ship);

    if (first_pass) {
        return;
    }

    /* Define transaction so that recipient will be notified. */
    if (num_transactions == MAX_TRANSACTIONS) {
        fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
        exit(-1);
    }

    n = num_transactions++;
    transaction[n].ttype      = INTERSPECIES_CONSTRUCTION;
    transaction[n].donor     = species_number;
    transaction[n].recipient = g_spec_number;
    transaction[n].value     = 2; /* Ship, not items. */
    transaction[n].number3   = total_cost + premium;
    strcpy(transaction[n].name1, species.name);
    strcpy(transaction[n].name2, ship_name(recipient_ship));
}

//*************************************************************************
// do_deep.c

func do_DEEP_command() {
    int i, found;

    char *original_line_pointer;


    /* Get the ship. */
    original_line_pointer = input_line_pointer;
    found = get_ship();
    if (!found) {
        /* Check for missing comma or tab after ship name. */
        input_line_pointer = original_line_pointer;
        fix_separator();
        found = get_ship();
        if (!found) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid ship name in ORBIT command.\n");
            return;
        }
    }

    if (ship.ttype == STARBASE) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! DEEP order may not be given for a starbase.\n");
        return;
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship is still under construction.\n");
        return;
    }

    if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
        return;
    }

    /* Make sure ship is not salvage of a disbanded colony. */
    if (disbanded_ship(ship)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! This ship is salvage of a disbanded colony!\n");
        return;
    }

    /* Move the ship. */
    ship.pn     = 0;
    ship.status = IN_DEEP_SPACE;

    /* Log result. */
    log_string("    ");
    log_string(ship_name(ship));
    log_string(" moved into deep space.\n");
}


//*************************************************************************
// do_des.c

func do_DESTROY_command() {
    int found;


    /* Get the ship. */
    correct_spelling_required = true;
    found = get_ship();
    if (!found) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid ship or starbase name in DESTROY command.\n");
        return;
    }

    /* Log result. */
    log_string("    ");
    log_string(ship_name(ship));

    if (first_pass) {
        log_string(" will be destroyed.\n");
        return;
    }

    log_string(" was destroyed.\n");

    delete_ship(ship);
}

//*************************************************************************
// do_dev.c

func do_DEVELOP_command() {
    int i, num_CUs, num_AUs, num_IUs, more_args, load_transport,
        capacity, resort_colony, mining_colony, production_penalty,
        CUs_only;

    char c, *original_line_pointer, *tp;

    long n, ni, na, amount_to_spend, original_cost, max_funds_available,
         ls_needed, raw_material_units, production_capacity,
         colony_production, ib, ab, md, denom, reb, specified_max;

    struct planet_data *colony_planet, *home_planet;
    struct nampla_data *temp_nampla, *colony_nampla;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get default spending limit. */
    max_funds_available = species.econ_units;
    if (max_funds_available > EU_spending_limit) {
        max_funds_available = EU_spending_limit;
    }
    max_funds_available += balance;

    /* Get specified spending limit, if any. */
    specified_max = -1;
    if (get_value()) {
        if (value == 0) {
            max_funds_available = balance;
        }else if (value > 0) {
            specified_max = value;
            if (value <= max_funds_available) {
                max_funds_available = value;
            }else{
                fprintf(log_file, "! WARNING: %s", input_line);
                fprintf(log_file,
                        "! Insufficient funds. Substituting %ld for %ld.\n",
                        max_funds_available, value);
                if (max_funds_available == 0) {
                    return;
                }
            }
        }else {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid spending limit.\n");
            return;
        }
    }

    /* See if there are any more arguments. */
    tp        = input_line_pointer;
    more_args = false;
    for (;c = *tp++;) {
        if (c == ';' || c == '\n') {
            break;
        }
        if (c == ' ' || c == '\t') {
            continue;
        }
        more_args = true;
        break;
    }

    if (!more_args) {
        /* Make sure planet is not a healthy home planet. */
        if (nampla.status & HOME_PLANET) {
            reb = species.hp_original_base - (nampla.mi_base + nampla.ma_base);
            if (reb > 0) {
                /* Home planet is recovering from bombing. */
                if (reb < max_funds_available) {
                    max_funds_available = reb;
                }
            }else {
                fprintf(log_file, "!!! Order ignored:\n");
                fprintf(log_file, "!!! %s", input_line);
                fprintf(log_file, "!!! You can only DEVELOP a home planet if it is recovering from bombing.\n");
                return;
            }
        }

        /* No arguments. Order is for this planet. */
        num_CUs = nampla.pop_units;
        if (2 * num_CUs > max_funds_available) {
            num_CUs = max_funds_available / 2;
        }
        if (num_CUs <= 0) {
            return;
        }

        colony_planet = planet_base + (long)nampla.planet_index;
        ib            = nampla.mi_base + nampla.IUs_to_install;
        ab            = nampla.ma_base + nampla.AUs_to_install;
        md            = colony_planet.mining_difficulty;

        denom   = 100 + md;
        num_AUs =
            (100 * (num_CUs + ib) - (md * ab) + denom / 2) / denom;
        num_IUs = num_CUs - num_AUs;

        if (num_IUs < 0) {
            num_AUs = num_CUs;
            num_IUs = 0;
        }
        if (num_AUs < 0) {
            num_IUs = num_CUs;
            num_AUs = 0;
        }

        amount_to_spend = num_CUs + num_AUs + num_IUs;

        if (check_bounced(amount_to_spend)) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Internal error. Please notify GM!\n");
            return;
        }

        nampla.pop_units         -= num_CUs;
        nampla.item_quantity[CU] += num_CUs;
        nampla.item_quantity[IU] += num_IUs;
        nampla.item_quantity[AU] += num_AUs;

        nampla.auto_IUs += num_IUs;
        nampla.auto_AUs += num_AUs;

        start_dev_log(num_CUs, num_IUs, num_AUs);
        log_string(".\n");

        check_population(nampla);

        return;
    }

    /* Get the planet to be developed. */
    temp_nampla           = nampla;
    original_line_pointer = input_line_pointer;
    i = get_location();
    if (!i || nampla == NULL) {
        /* Check for missing comma or tab after source name. */
        input_line_pointer = original_line_pointer;
        fix_separator();
        i = get_location();
    }
    colony_nampla = nampla;
    nampla        = temp_nampla;
    if (!i || colony_nampla == NULL) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid planet name in DEVELOP command.\n");
        return;
    }

    /* Make sure planet is not a healthy home planet. */
    if (colony_nampla.status & HOME_PLANET) {
        reb = species.hp_original_base - (colony_nampla.mi_base + colony_nampla.ma_base);
        if (reb > 0) {
            /* Home planet is recovering from bombing. */
            if (reb < max_funds_available) {
                max_funds_available = reb;
            }
        }else {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! You can only DEVELOP a home planet if it is recovering from bombing.\n");
            return;
        }
    }

    /* Determine if its a mining or resort colony, and if it can afford to
    *  build its own IUs and AUs. Note that we cannot use nampla.status
    *  because it is not correctly set until the Finish program is run. */

    home_planet   = planet_base + (long)nampla_base.planet_index;
    colony_planet = planet_base + (long)colony_nampla.planet_index;
    ls_needed     = life_support_needed(species, home_planet, colony_planet);

    ni = colony_nampla.mi_base + colony_nampla.IUs_to_install;
    na = colony_nampla.ma_base + colony_nampla.AUs_to_install;

    if (ni > 0 && na == 0) {
        colony_production = 0;
        mining_colony     = true;
        resort_colony     = false;
    }else if (na > 0 && ni == 0 && ls_needed <= 6 &&
              colony_planet.gravity <= home_planet.gravity) {
        colony_production = 0;
        resort_colony     = true;
        mining_colony     = false;
    }else {
        mining_colony = false;
        resort_colony = false;

        raw_material_units = (10L * (long)species.tech_level[MI] * ni)
                             / (long)colony_planet.mining_difficulty;
        production_capacity = ((long)species.tech_level[MA] * na) / 10L;

        if (ls_needed == 0) {
            production_penalty = 0;
        }else{
            production_penalty = (100 * ls_needed) / species.tech_level[LS];
        }

        raw_material_units  -= (production_penalty * raw_material_units) / 100;
        production_capacity -= (production_penalty * production_capacity) / 100;

        colony_production = (production_capacity > raw_material_units)
                            ? raw_material_units : production_capacity;

        colony_production -= colony_nampla.IUs_needed
                             + colony_nampla.AUs_needed;

        /* In case there is more than one DEVELOP order for
         *      this colony. */
    }

    /* See if there are more arguments. */
    tp        = input_line_pointer;
    more_args = false;
    for (;c = *tp++;) {
        if (c == ';' || c == '\n') {
            break;
        }
        if (c == ' ' || c == '\t') {
            continue;
        }
        more_args = true;
        break;
    }

    if (more_args) {
        load_transport = true;

        /* Get the ship to receive the cargo. */
        if (!get_ship()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship to be loaded does not exist!\n");
            return;
        }

        if (ship.class == TR) {
            capacity = (10 + ((int)ship.tonnage / 2)) * (int)ship.tonnage;
        }else if (ship.class == BA) {
            capacity = 10 * ship.tonnage;
        }else{
            capacity = ship.tonnage;
        }

        for (i = 0; i < MAX_ITEMS; i++) {
            capacity -= ship.item_quantity[i] * item_carry_capacity[i];
        }

        if (capacity <= 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! %s was already full and could take no more cargo!\n",
                    ship_name(ship));
            return;
        }

        if (capacity > max_funds_available) {
            capacity = max_funds_available;
            if (max_funds_available != specified_max) {
                fprintf(log_file, "! WARNING: %s", input_line);
                fprintf(log_file, "! Insufficient funds to completely fill %s!\n",
                        ship_name(ship));
                fprintf(log_file, "! Will use all remaining funds (= %d).\n",
                        capacity);
            }
        }
    }else {
        load_transport = false;

        /* No more arguments. Order is for a colony in the same sector as the
         *  producing planet. */
        if (nampla.x != colony_nampla.x || nampla.y != colony_nampla.y ||
            nampla.z != colony_nampla.z) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Colony and producing planet are not in the same sector.\n");
            return;
        }

        num_CUs = nampla.pop_units;
        if (2 * num_CUs > max_funds_available) {
            num_CUs = max_funds_available / 2;
        }
    }

    CUs_only = false;
    if (mining_colony) {
        if (load_transport) {
            num_CUs = capacity / 2;
            if (num_CUs > nampla.pop_units) {
                fprintf(log_file, "! WARNING: %s", input_line);
                fprintf(log_file,
                        "! Insufficient available population! %d CUs are needed",
                        num_CUs);
                num_CUs = nampla.pop_units;
                fprintf(log_file, " to fill ship but only %d can be built.\n",
                        num_CUs);
            }
        }

        num_AUs = 0;
        num_IUs = num_CUs;
    }else if (resort_colony) {
        if (load_transport) {
            num_CUs = capacity / 2;
            if (num_CUs > nampla.pop_units) {
                fprintf(log_file, "! WARNING: %s", input_line);
                fprintf(log_file,
                        "! Insufficient available population! %d CUs are needed",
                        num_CUs);
                num_CUs = nampla.pop_units;
                fprintf(log_file, " to fill ship but only %d can be built.\n",
                        num_CUs);
            }
        }

        num_IUs = 0;
        num_AUs = num_CUs;
    }else {
        if (load_transport) {
            if (colony_production >= capacity) {
                /* Colony can build its own IUs and AUs. */
                num_CUs  = capacity;
                CUs_only = true;
            }else {
                /* Build IUs and AUs for the colony. */
                num_CUs = capacity / 2;
            }

            if (num_CUs > nampla.pop_units) {
                fprintf(log_file, "! WARNING: %s", input_line);
                fprintf(log_file,
                        "! Insufficient available population! %d CUs are needed",
                        num_CUs);
                num_CUs = nampla.pop_units;
                fprintf(log_file, " to fill ship, but\n!   only %d can be built.\n",
                        num_CUs);
            }
        }

        colony_planet = planet_base + (long)colony_nampla.planet_index;

        i       = 100 + (int)colony_planet.mining_difficulty;
        num_AUs = ((100 * num_CUs) + (i + 1) / 2) / i;
        num_IUs = num_CUs - num_AUs;
    }

    if (num_CUs <= 0) {
        return;
    }

    /* Make sure there's enough money to pay for it all. */
    if (load_transport && CUs_only) {
        amount_to_spend = num_CUs;
    }else{
        amount_to_spend = num_CUs + num_IUs + num_AUs;
    }

    if (check_bounced(amount_to_spend)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Internal error. Notify GM!\n");
        return;
    }

    /* Start logging what happened. */
    if (load_transport && CUs_only) {
        start_dev_log(num_CUs, 0, 0);
    }else {
        start_dev_log(num_CUs, num_IUs, num_AUs);
    }

    log_string(" for PL ");    log_string(colony_nampla.name);

    nampla.pop_units -= num_CUs;

    if (load_transport) {
        if (CUs_only) {
            colony_nampla.IUs_needed += num_IUs;
            colony_nampla.AUs_needed += num_AUs;
        }

        if (nampla.x != ship.x || nampla.y != ship.y ||
            nampla.z != ship.z) {
            nampla.item_quantity[CU] += num_CUs;
            if (!CUs_only) {
                nampla.item_quantity[IU] += num_IUs;
                nampla.item_quantity[AU] += num_AUs;
            }

            log_string(" but will remain on the planet's surface because ");
            log_string(ship_name(ship));
            log_string(" is not in the same sector.");
        }else {
            ship.item_quantity[CU] += num_CUs;
            if (!CUs_only) {
                ship.item_quantity[IU] += num_IUs;
                ship.item_quantity[AU] += num_AUs;
            }

            n = colony_nampla - nampla_base;
            if (n == 0) {
                n = 9999;               /* Home planet. */
            }
            ship.unloading_point = n;

            n = nampla - nampla_base;
            if (n == 0) {
                n = 9999;               /* Home planet. */
            }
            ship.loading_point = n;

            log_string(" and transferred to ");
            log_string(ship_name(ship));
        }
    }else {
        colony_nampla.item_quantity[CU] += num_CUs;
        colony_nampla.item_quantity[IU] += num_IUs;
        colony_nampla.item_quantity[AU] += num_AUs;

        colony_nampla.auto_IUs += num_IUs;
        colony_nampla.auto_AUs += num_AUs;

        log_string(" and transferred to PL ");
        log_string(colony_nampla.name);

        check_population(colony_nampla);
    }

    log_string(".\n");
}

func start_dev_log(num_CUs, num_IUs, num_AUs int) {
    log_string("    ");
    log_int(num_CUs);  log_string(" Colonist Unit");
    if (num_CUs != 1) {
        log_char('s');
    }

    if (num_IUs + num_AUs == 0) {
        goto done;
    }

    if (num_IUs > 0) {
        if (num_AUs == 0) {
            log_string(" and ");
        }else{
            log_string(", ");
        }

        log_int(num_IUs);  log_string(" Colonial Mining Unit");
        if (num_IUs != 1) {
            log_char('s');
        }
    }

    if (num_AUs > 0) {
        if (num_IUs > 0) {
            log_char(',');
        }

        log_string(" and ");

        log_int(num_AUs);  log_string(" Colonial Manufacturing Unit");
        if (num_AUs != 1) {
            log_char('s');
        }
    }

done:

    log_string(" were built");
}

//*************************************************************************
// do_disband.c


func do_DISBAND_command() {
    int found;


    /* Get the planet. */
    found = get_location();
    if (!found || nampla == NULL) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid planet name in DISBAND command.\n");
        return;
    }

    /* Make sure planet is not the home planet. */
    if (nampla.status & HOME_PLANET) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You cannot disband your home planet!\n");
        return;
    }

    /* Make sure planet is not under siege. */
    if (nampla.siege_eff) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You cannot disband a planet that is under siege!\n");
        return;
    }

    /* Mark the colony as "disbanded" and convert mining and manufacturing
     *  base to CUs, IUs, and AUs. */
    nampla.status            |= DISBANDED_COLONY;
    nampla.item_quantity[CU] += nampla.mi_base + nampla.ma_base;
    nampla.item_quantity[IU] += nampla.mi_base / 2;
    nampla.item_quantity[AU] += nampla.ma_base / 2;
    nampla.mi_base            = 0;
    nampla.ma_base            = 0;

    /* Log the event. */
    log_string("    The colony on PL ");
    log_string(nampla.name);
    log_string(" was ordered to disband.\n");
}


//*************************************************************************
// do_enemy.c


func do_ENEMY_command() {
    int i, array_index, bit_number;

    long bit_mask;


    /* See if declaration is for all species. */
    if (get_value()) {
        bit_mask = 0;
        for (i = 0; i < NUM_CONTACT_WORDS; i++) {
            species.enemy[i] = ~bit_mask;      /* Set all enemy bits. */
            species.ally[i]  = bit_mask;       /* Clear all ally bits. */
        }
    }else {
        /* Get name of species that is being declared an enemy. */
        if (!get_species_name()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid or missing argument in ENEMY command.\n");
            return;
        }

        /* Get array index and bit mask. */
        array_index = (g_spec_number - 1) / 32;
        bit_number  = (g_spec_number - 1) % 32;
        bit_mask    = 1 << bit_number;

        /* Set/clear the appropriate bit. */
        species.enemy[array_index] |= bit_mask;        /* Set enemy bit. */
        species.ally[array_index]  &= ~bit_mask;       /* Clear ally bit. */
    }

    /* Log the result. */
    log_string("    Enmity was declared towards ");
    if (bit_mask == 0) {
        log_string("ALL species");
    }else{
        log_string("SP ");
        log_string(g_spec_name);
    }
    log_string(".\n");
}

//*************************************************************************
// do_est.c


func do_ESTIMATE_command() {
    int i, max_error, estimate[6], contact_word_number,
        contact_bit_number;

    long cost, contact_mask;

    struct species_data *alien;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get name of alien species. */
    if (!get_species_name()) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid species name in ESTIMATE command.\n");
        return;
    }

    /* Check if we've met this species. */
    contact_word_number = (g_spec_number - 1) / 32;
    contact_bit_number  = (g_spec_number - 1) % 32;
    contact_mask        = 1 << contact_bit_number;
    if ((species.contact[contact_word_number] & contact_mask) == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't do an estimate of a species you haven't met.\n");
        return;
    }

    /* Check if sufficient funds are available. */
    cost = 25;
    if (check_bounced(cost)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        return;
    }

    /* Log the result. */
    if (first_pass) {
        log_string("    An estimate of the technology of SP ");
        log_string(g_spec_name);
        log_string(" was made at a cost of ");
        log_long(cost);
        log_string(".\n");
        return;
    }

    /* Make the estimates. */
    alien = &spec_data[g_spec_number - 1];
    for (i = 0; i < 6; i++) {
        max_error = (int)alien.tech_level[i] - (int)species.tech_level[i];
        if (max_error < 1) {
            max_error = 1;
        }
        estimate[i] = (int)alien.tech_level[i] + rnd((2 * max_error) + 1)
                      - (max_error + 1);
        if (alien.tech_level[i] == 0) {
            estimate[i] = 0;
        }
        if (estimate[i] < 0) {
            estimate[i] = 0;
        }
    }

    log_string("    Estimate of the technology of SP ");
    log_string(alien.name);
    log_string(" (government name '");
    log_string(alien.govt_name);
    log_string("', government ttype '");
    log_string(alien.govt_ttype);
    log_string("'):\n      MI = ");   log_int(estimate[MI]);
    log_string(", MA = ");   log_int(estimate[MA]);
    log_string(", ML = ");   log_int(estimate[ML]);
    log_string(", GV = ");   log_int(estimate[GV]);
    log_string(", LS = ");   log_int(estimate[LS]);
    log_string(", BI = ");   log_int(estimate[BI]);
    log_string(".\n");
}

//*************************************************************************
// do_germ.c


extern int   num_transactions;
extern char  field_distorted[MAX_SPECIES];
extern short germ_bombs_used[MAX_SPECIES][MAX_SPECIES];
extern struct planet_data * planet_base;
extern struct species_data *c_species[MAX_SPECIES];
extern struct nampla_data * c_nampla[MAX_SPECIES];
extern struct ship_data *   c_ship[MAX_SPECIES];
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_germ_warfare(attacking_species, defending_species, defender_index, bat, act)

int attacking_species, defending_species, defender_index;
struct battle_data *bat;
struct action_data *act;

{
    int i, attacker_BI, defender_BI, success_chance, num_bombs,
        success;

    long econ_units_from_looting;

    struct planet_data *planet;
    struct nampla_data *attacked_nampla;
    struct ship_data *  sh;


    attacker_BI     = c_species[attacking_species].tech_level[BI];
    defender_BI     = c_species[defending_species].tech_level[BI];
    attacked_nampla = (struct nampla_data *)act.fighting_unit[defender_index];
    planet          = planet_base + (long)attacked_nampla.planet_index;

    success_chance = 50 + (2 * (attacker_BI - defender_BI));
    success        = false;
    num_bombs      = germ_bombs_used[attacking_species][defending_species];

    for (i = 0; i < num_bombs; i++) {
        if (rnd(100) <= success_chance) {
            success = true;
            break;
        }
    }

    if (success) {
        log_string("        Unfortunately");
    }else{
        log_string("        Fortunately");
    }

    log_string(" for the ");
    log_string(c_species[defending_species].name);
    log_string(" defenders of PL ");
    log_string(attacked_nampla.name);
    log_string(", the ");
    i = bat.spec_num[attacking_species];
    if (field_distorted[attacking_species]) {
        log_int(distorted(i));
    }else{
        log_string(c_species[attacking_species].name);
    }
    log_string(" attackers ");

    if (!success) {
        log_string("failed");

        if (num_bombs <= 0) {
            log_string(" because they didn't have any germ warfare bombs");
        }

        log_string("!\n");

        return;
    }

    log_string("succeeded, using ");
    log_int(num_bombs);
    log_string(" germ warfare bombs. The defenders were wiped out!\n");

    /* Take care of looting. */
    econ_units_from_looting =
        (long)attacked_nampla.mi_base + (long)attacked_nampla.ma_base;

    if (attacked_nampla.status & HOME_PLANET) {
        if (c_species[defending_species].hp_original_base < econ_units_from_looting) {
            c_species[defending_species].hp_original_base = econ_units_from_looting;
        }

        econ_units_from_looting *= 5;
    }

    if (econ_units_from_looting > 0) {
        /* Check if there's enough memory for a new interspecies transaction. */
        if (num_transactions == MAX_TRANSACTIONS) {
            fprintf(stderr, "\nRan out of memory! MAX_TRANSACTIONS is too small!\n\n");
            exit(-1);
        }
        i = num_transactions++;

        /* Define this transaction. */
        transaction[i].ttype      = LOOTING_EU_TRANSFER;
        transaction[i].donor     = bat.spec_num[defending_species];
        transaction[i].recipient = bat.spec_num[attacking_species];
        transaction[i].value     = econ_units_from_looting;
        strcpy(transaction[i].name1, c_species[defending_species].name);
        strcpy(transaction[i].name2, c_species[attacking_species].name);
        strcpy(transaction[i].name3, attacked_nampla.name);
    }

    /* Finish off defenders. */
    attacked_nampla.mi_base        = 0;
    attacked_nampla.ma_base        = 0;
    attacked_nampla.IUs_to_install = 0;
    attacked_nampla.AUs_to_install = 0;
    attacked_nampla.pop_units      = 0;
    attacked_nampla.siege_eff      = 0;
    attacked_nampla.shipyards      = 0;
    attacked_nampla.hiding         = 0;
    attacked_nampla.hidden         = 0;
    attacked_nampla.use_on_ambush  = 0;

    for (i = 0; i < MAX_ITEMS; i++) {
        attacked_nampla.item_quantity[i] = 0;
    }

    /* Reset status word. */
    if (attacked_nampla.status & HOME_PLANET) {
        attacked_nampla.status = HOME_PLANET;
    }else{
        attacked_nampla.status = COLONY;
    }

    /* Delete any ships that were under construction on the planet. */
    sh = c_ship[defending_species] - 1;
    for (i = 0; i < c_species[defending_species].num_ships; i++) {
        ++sh;

        if (sh.x != attacked_nampla.x) {
            continue;
        }
        if (sh.y != attacked_nampla.y) {
            continue;
        }
        if (sh.z != attacked_nampla.z) {
            continue;
        }
        if (sh.pn != attacked_nampla.pn) {
            continue;
        }

        delete_ship(sh);
    }
}

//*************************************************************************
// do_hide.c


extern int   doing_production, first_pass, abbr_index;
extern long  value, balance;
extern char  input_line[256];
extern FILE *log_file;
extern struct species_data *species;
extern struct nampla_data * nampla;


do_HIDE_command() {
    int n, status;

    long cost;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Make sure this is not a mining colony or home planet. */
    if (nampla.status & HOME_PLANET) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You may not HIDE a home planet.\n");
        return;
    }
    if (nampla.status & RESORT_COLONY) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You may not HIDE a resort colony.\n");
        return;
    }

    /* Check if planet is under siege. */
    if (nampla.siege_eff != 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Besieged planet cannot HIDE!\n");
        return;
    }

    /* Check if sufficient funds are available. */
    cost = (nampla.mi_base + nampla.ma_base) / 10L;
    if (nampla.status & MINING_COLONY) {
        if (cost > species.econ_units) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Mining colony does not have sufficient EUs to hide.\n");
            return;
        }else {
            species.econ_units -= cost;
        }
    }else if (check_bounced(cost)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        return;
    }

    /* Set 'hiding' flag. */
    nampla.hiding = true;

    /* Log transaction. */
    log_string("    Spent ");  log_long(cost);
    log_string(" hiding this colony.\n");
}

//*************************************************************************
// do_inst.c


extern int   abbr_index, species_number;
extern char  input_line[256];
extern long  value;
extern FILE *log_file;
extern struct galaxy_data   galaxy;
extern struct species_data *species;
extern struct nampla_data * nampla;


do_INSTALL_command() {
    int i, item_class, item_count, num_available, do_all_units,
        recovering_home_planet, alien_index;

    long n, current_pop, reb;

    struct nampla_data *alien_home_nampla;


    /* Get number of items to install. */
    if (get_value()) {
        do_all_units = false;
    }else{
        do_all_units = true;
        item_count   = 0;
        item_class   = IU;
        goto get_planet;
    }

    /* Make sure value is meaningful. */
    if (value < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid item count in INSTALL command.\n");
        return;
    }
    item_count = value;

    /* Get class of item. */
    item_class = get_class_abbr();
    if (item_class != ITEM_CLASS || (abbr_index != IU && abbr_index != AU)) {
        /* Players sometimes accidentally use "MI" for "IU"
         *      or "MA" for "AU". */
        if (item_class == TECH_ID && abbr_index == MI) {
            abbr_index = IU;
        }else if (item_class == TECH_ID && abbr_index == MA) {
            abbr_index = AU;
        }else{
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid item class!\n");
            return;
        }
    }
    item_class = abbr_index;

get_planet:

    /* Get planet where items are to be installed. */
    if (!get_location() || nampla == NULL) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid planet name in INSTALL command.\n");
        return;
    }

    /* Make sure this is not someone else's populated homeworld. */
    for (alien_index = 0; alien_index < galaxy.num_species; alien_index++) {
        if (species_number == alien_index + 1) {
            continue;
        }
        if (!data_in_memory[alien_index]) {
            continue;
        }

        alien_home_nampla = namp_data[alien_index];

        if (alien_home_nampla.x != nampla.x) {
            continue;
        }
        if (alien_home_nampla.y != nampla.y) {
            continue;
        }
        if (alien_home_nampla.z != nampla.z) {
            continue;
        }
        if (alien_home_nampla.pn != nampla.pn) {
            continue;
        }
        if ((alien_home_nampla.status & POPULATED) == 0) {
            continue;
        }

        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You may not colonize someone else's populated home planet!\n");

        return;
    }

    /* Make sure it's not a healthy home planet. */
    recovering_home_planet = false;
    if (nampla.status & HOME_PLANET) {
        n = nampla.mi_base + nampla.ma_base + nampla.IUs_to_install +
            nampla.AUs_to_install;
        reb = species.hp_original_base - n;

        if (reb > 0) {
            recovering_home_planet = true;      /* HP was bombed. */
        }else {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Installation not allowed on a healthy home planet!\n");
            return;
        }
    }

check_items:

    /* Make sure planet has the specified items. */
    if (item_count == 0) {
        item_count = nampla.item_quantity[item_class];

        if (nampla.item_quantity[CU] < item_count) {
            item_count = nampla.item_quantity[CU];
        }

        if (item_count == 0) {
            if (do_all_units) {
                item_count   = 0;
                item_class   = AU;
                do_all_units = false;
                goto check_items;
            }else {
                return;
            }
        }
    }else if (nampla.item_quantity[item_class] < item_count) {
        fprintf(log_file, "! WARNING: %s", input_line);
        fprintf(log_file,
                "! Planet does not have %d %ss. Substituting 0 for %d!\n",
                item_count, item_abbr[item_class], item_count);
        item_count = 0;
        goto check_items;
    }

    if (recovering_home_planet) {
        if (item_count > reb) {
            item_count = reb;
        }
        reb -= item_count;
    }

    /* Make sure planet has enough colonist units. */
    num_available = nampla.item_quantity[CU];
    if (num_available < item_count) {
        if (num_available > 0) {
            fprintf(log_file, "! WARNING: %s", input_line);
            fprintf(log_file, "! Planet does not have %d CUs. Substituting %d for %d!\n",
                    item_count, num_available, item_count);
            item_count = num_available;
        }else {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! No colonist units on planet for installation.\n");
            return;
        }
    }

    /* Start the installation. */
    nampla.item_quantity[CU]         -= item_count;
    nampla.item_quantity[item_class] -= item_count;

    if (item_class == IU) {
        nampla.IUs_to_install += item_count;
    }else{
        nampla.AUs_to_install += item_count;
    }

    /* Log result. */
    log_string("    Installation of ");    log_int(item_count);
    log_char(' ');    log_string(item_name[item_class]);
    if (item_count != 1) {
        log_char('s');
    }
    log_string(" began on PL ");
    log_string(nampla.name);
    log_string(".\n");

    if (do_all_units) {
        item_count   = 0;
        item_class   = AU;
        do_all_units = false;
        goto check_items;
    }

    check_population(nampla);
}

//*************************************************************************
// do_int.c

int num_intercepts;

struct {
    char x, y, z;
    long amount_spent;
} intercept[MAX_INTERCEPTS];

extern int doing_production, first_pass, abbr_index,
           species_number, num_transactions;
extern long  value, balance;
extern char  input_line[256], *ship_name();
extern FILE *log_file;
extern struct galaxy_data   galaxy;
extern struct species_data *species;
extern struct nampla_data * nampla;
extern struct trans_data    transaction[MAX_TRANSACTIONS];



do_INTERCEPT_command() {
    int i, n, status;

    long cost;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get amount to spend. */
    status = get_value();
    if (status == 0 || value < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid or missing amount.\n");
        return;
    }
    if (value == 0) {
        value = balance;
    }
    if (value == 0) {
        return;
    }
    cost = value;

    /* Check if planet is under siege. */
    if (nampla.siege_eff != 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Besieged planet cannot INTERCEPT!\n");
        return;
    }

    /* Check if sufficient funds are available. */
    if (check_bounced(cost)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        return;
    }

    log_string("    Preparations were made for an interception at a cost of ");
    log_long(cost);
    log_string(".\n");

    if (first_pass) {
        return;
    }

    /* Allocate funds. */
    for (i = 0; i < num_intercepts; i++) {
        if (nampla.x != intercept[i].x) {
            continue;
        }
        if (nampla.y != intercept[i].y) {
            continue;
        }
        if (nampla.z != intercept[i].z) {
            continue;
        }

        /* This interception was started by another planet in the same
         *      star system. */
        intercept[i].amount_spent += cost;
        return;
    }

    if (num_intercepts == MAX_INTERCEPTS) {
        fprintf(stderr, "\n\tMAX_INTERCEPTS exceeded in do_int.c!\n\n");
        exit(-1);
    }

    intercept[num_intercepts].x            = nampla.x;
    intercept[num_intercepts].y            = nampla.y;
    intercept[num_intercepts].z            = nampla.z;
    intercept[num_intercepts].amount_spent = cost;

    ++num_intercepts;
}




handle_intercept(intercept_index)

int intercept_index;

{
    int i, j, n, num_enemy_ships, alien_index, enemy_index, enemy_num,
        num_ships_left, array_index, bit_number, is_an_enemy,
        is_distorted;

    char enemy_number[MAX_ENEMY_SHIPS];

    long bit_mask, cost_to_destroy;

    struct species_data *alien;
    struct ship_data *   alien_sh, *enemy_sh,
                     *enemy_ship[MAX_ENEMY_SHIPS];


    /* Make a list of all enemy ships that jumped into this system. */
    num_enemy_ships = 0;
    for (alien_index = 0; alien_index < galaxy.num_species; alien_index++) {
        if (!data_in_memory[alien_index]) {
            continue;
        }

        if (species_number == alien_index + 1) {
            continue;
        }

        /* Is it an enemy species? */
        array_index = (alien_index) / 32;
        bit_number  = (alien_index) % 32;
        bit_mask    = 1 << bit_number;
        if (species.enemy[array_index] & bit_mask) {
            is_an_enemy = true;
        }else{
            is_an_enemy = false;
        }

        /* Find enemy ships, if any, that jumped to this location. */
        alien    = &spec_data[alien_index];
        alien_sh = ship_data[alien_index] - 1;
        for (i = 0; i < alien.num_ships; i++) {
            ++alien_sh;

            if (alien_sh.pn == 99) {
                continue;
            }

            /* Did it jump this turn? */
            if (!alien_sh.just_jumped) {
                continue;
            }
            if (alien_sh.just_jumped == 50) {
                continue;                               /* Ship MOVEd. */
            }
            /* Did it enter this star system? */
            if (alien_sh.x != intercept[intercept_index].x) {
                continue;
            }
            if (alien_sh.y != intercept[intercept_index].y) {
                continue;
            }
            if (alien_sh.z != intercept[intercept_index].z) {
                continue;
            }

            /* Is it field-distorted? */
            if (alien_sh.item_quantity[FD] == alien_sh.tonnage) {
                is_distorted = true;
            }else{
                is_distorted = false;
            }

            if (!is_an_enemy && !is_distorted) {
                continue;
            }

            /* This is an enemy ship that just jumped into the system. */
            if (num_enemy_ships == MAX_ENEMY_SHIPS) {
                fprintf(stderr, "\n\tERROR! Array overflow in do_int.c!\n\n");
                exit(-1);
            }
            enemy_number[num_enemy_ships] = alien_index + 1;
            enemy_ship[num_enemy_ships]   = alien_sh;
            ++num_enemy_ships;
        }
    }

    if (num_enemy_ships == 0) {
        return;                         /* Nothing to intercept. */
    }
    num_ships_left = num_enemy_ships;
    for (;num_ships_left > 0;) {
        /* Select ship for interception. */
        enemy_index = rnd(num_enemy_ships) - 1;
        if (enemy_ship[enemy_index] == NULL) {
            continue;                                   /* We already did this
                                                         *      one. */
        }
        enemy_num = enemy_number[enemy_index];
        enemy_sh  = enemy_ship[enemy_index];

        /* Are there enough funds to destroy this ship? */
        cost_to_destroy = 100L * (long)enemy_sh.tonnage;
        if (enemy_sh.class == TR) {
            cost_to_destroy /= 10;
        }
        if (cost_to_destroy > intercept[intercept_index].amount_spent) {
            break;
        }

        /* Is the ship too large? Check only if ship did NOT arrive via a
         *      natural wormhole. */
        if (enemy_sh.just_jumped != 99) {
            if (enemy_sh.tonnage > 20) {
                break;
            }
            if (enemy_sh.class != TR && enemy_sh.tonnage > 5) {
                break;
            }
        }

        /* Update funds available. */
        intercept[intercept_index].amount_spent -= cost_to_destroy;

        /* Log the result for current species. */
        log_string("\n! ");
        n = enemy_sh.item_quantity[FD];        /* Show real name. */
        enemy_sh.item_quantity[FD] = 0;
        log_string(ship_name(enemy_sh));
        enemy_sh.item_quantity[FD] = n;

        /* List cargo destroyed. */
        n = 0;
        for (j = 0; j < MAX_ITEMS; j++) {
            if (enemy_sh.item_quantity[j] > 0) {
                if (n++ == 0) {
                    log_string(" (cargo: ");
                }else{
                    log_char(',');
                }
                log_int((int)enemy_sh.item_quantity[j]);
                log_char(' ');
                log_string(item_abbr[j]);
            }
        }
        if (n > 0) {
            log_char(')');
        }

        log_string(", owned by SP ");
        log_string(spec_data[enemy_num - 1].name);
        log_string(", was successfully intercepted and destroyed in sector ");
        log_int(enemy_sh.x);  log_char(' ');
        log_int(enemy_sh.y);  log_char(' ');
        log_int(enemy_sh.z);
        log_string(".\n");

        /* Create interspecies transaction so that other player will be
         *      notified. */
        if (num_transactions == MAX_TRANSACTIONS) {
            fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS in do_int.c!\n\n");
            exit(-1);
        }

        n = num_transactions++;
        transaction[n].ttype    = SHIP_MISHAP;
        transaction[n].value   = 1;     /* Interception. */
        transaction[n].number1 = enemy_number[enemy_index];
        strcpy(transaction[n].name1, ship_name(enemy_sh));

        delete_ship(enemy_sh);

        enemy_ship[enemy_index] = NULL; /* Don't select this ship again. */

        --num_ships_left;
    }
}

//*************************************************************************
// do_land.c


extern int  first_pass, num_transactions, species_number;
extern long value;
extern char input_line[256], original_line[256],
            *input_line_pointer, *ship_name();
extern FILE *log_file;
extern struct galaxy_data   galaxy;
extern struct species_data *species;
extern struct nampla_data * nampla;
extern struct ship_data *   ship;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_LAND_command() {
    int i, n, found, siege_effectiveness, landing_detected, landed,
        alien_number, alien_index, alien_pn, array_index, bit_number,
        requested_alien_landing, alien_here, already_logged;

    long bit_mask;

    char *original_line_pointer;

    struct species_data *alien;
    struct nampla_data * alien_nampla;


    /* Get the ship. */
    original_line_pointer = input_line_pointer;
    found = get_ship();
    if (!found) {
        /* Check for missing comma or tab after ship name. */
        input_line_pointer = original_line_pointer;
        fix_separator();
        found = get_ship();
        if (!found) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid ship name in LAND command.\n");
            return;
        }
    }

    /* Make sure the ship is not a starbase. */
    if (ship.ttype == STARBASE) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! A starbase cannot land on a planet!\n");
        return;
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship is still under construction.\n");
        return;
    }

    if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
        return;
    }

    /* Get the planet number, if specified. */
    found = get_value();

get_planet:

    alien_pn   = 0;
    alien_here = false;
    requested_alien_landing = false;
    landed = false;
    if (!found) {
        found = get_location();
        if (!found || nampla == NULL) {
            found = false;
        }
    }else {
        /* Check if we or another species that has declared us ALLY has
         *      a colony on this planet. */
        found    = false;
        alien_pn = value;
        requested_alien_landing = true;
        array_index             = (species_number - 1) / 32;
        bit_number = (species_number - 1) % 32;
        bit_mask   = 1 << bit_number;
        for (alien_index = 0; alien_index < galaxy.num_species; alien_index++) {
            if (!data_in_memory[alien_index]) {
                continue;
            }

            alien        = &spec_data[alien_index];
            alien_nampla = namp_data[alien_index] - 1;
            for (i = 0; i < alien.num_namplas; i++) {
                ++alien_nampla;
                if (ship.x != alien_nampla.x) {
                    continue;
                }
                if (ship.y != alien_nampla.y) {
                    continue;
                }
                if (ship.z != alien_nampla.z) {
                    continue;
                }
                if (alien_pn != alien_nampla.pn) {
                    continue;
                }
                if ((alien_nampla.status & POPULATED) == 0) {
                    continue;
                }

                if (alien_index == species_number - 1) {
                    /* We have a colony here. No permission needed. */
                    nampla     = alien_nampla;
                    found      = true;
                    alien_here = false;
                    requested_alien_landing = false;
                    goto finish_up;
                }

                alien_here = true;

                if ((alien.ally[array_index] & bit_mask) == 0) {
                    continue;
                }

                found = true;
                break;
            }

            if (found) {
                break;
            }
        }
    }

finish_up:

    already_logged = false;

    if (requested_alien_landing && alien_here) {
        /* Notify the other alien(s). */
        landed = found;
        for (alien_index = 0; alien_index < galaxy.num_species; alien_index++) {
            if (!data_in_memory[alien_index]) {
                continue;
            }

            if (alien_index == species_number - 1) {
                continue;
            }

            alien        = &spec_data[alien_index];
            alien_nampla = namp_data[alien_index] - 1;
            for (i = 0; i < alien.num_namplas; i++) {
                ++alien_nampla;
                if (ship.x != alien_nampla.x) {
                    continue;
                }
                if (ship.y != alien_nampla.y) {
                    continue;
                }
                if (ship.z != alien_nampla.z) {
                    continue;
                }
                if (alien_pn != alien_nampla.pn) {
                    continue;
                }
                if ((alien_nampla.status & POPULATED) == 0) {
                    continue;
                }

                if ((alien.ally[array_index] & bit_mask) != 0) {
                    found = true;
                }else{
                    found = false;
                }

                if (landed && !found) {
                    continue;
                }

                if (landed) {
                    log_string("    ");
                }else{
                    log_string("!!! ");
                }

                log_string(ship_name(ship));

                if (landed) {
                    log_string(" was granted");
                }else{
                    log_string(" was denied");
                }
                log_string(" permission to land on PL ");
                log_string(alien_nampla.name);
                log_string(" by SP ");
                log_string(alien.name);
                log_string(".\n");

                already_logged = true;

                nampla = alien_nampla;

                if (first_pass) {
                    break;
                }

                /* Define a 'landing request' transaction. */
                if (num_transactions == MAX_TRANSACTIONS) {
                    fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
                    exit(-1);
                }

                n = num_transactions++;
                transaction[n].ttype    = LANDING_REQUEST;
                transaction[n].value   = landed;
                transaction[n].number1 = alien_index + 1;
                strcpy(transaction[n].name1, alien_nampla.name);
                strcpy(transaction[n].name2, ship_name(ship));
                strcpy(transaction[n].name3, species.name);

                break;
            }
        }

        found = true;
    }

    if (alien_here && !landed) {
        return;
    }

    if (!found) {
        if ((ship.status == IN_ORBIT || ship.status == ON_SURFACE) &&
            !requested_alien_landing) {
            /* Player forgot to specify planet. Use the one it's already at. */
            value = ship.pn;
            found = true;
            goto get_planet;
        }

        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Invalid or missing planet in LAND command.\n");
        return;
    }

    /* Make sure the ship and the planet are in the same star system. */
    if (ship.x != nampla.x || ship.y != nampla.y ||
        ship.z != nampla.z) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship and planet are not in the same sector.\n");
        return;
    }

    /* Make sure planet is populated. */
    if ((nampla.status & POPULATED) == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Planet in LAND command is not populated.\n");
        return;
    }

    /* Move the ship. */
    ship.pn     = nampla.pn;
    ship.status = ON_SURFACE;

    if (already_logged) {
        return;
    }

    /* If the planet is under siege, the landing may be detected by the
     *  besiegers. */
    log_string("    ");
    log_string(ship_name(ship));

    if (nampla.siege_eff != 0) {
        if (first_pass) {
            log_string(" will attempt to land on PL ");
            log_string(nampla.name);
            log_string(" in spite of the siege");
        }else {
            if (nampla.siege_eff < 0) {
                siege_effectiveness = -nampla.siege_eff;
            }else{
                siege_effectiveness = nampla.siege_eff;
            }

            landing_detected = false;
            if (rnd(100) <= siege_effectiveness) {
                landing_detected = true;
                for (i = 0; i < num_transactions; i++) {
                    /* Find out who is besieging this planet. */
                    if (transaction[i].ttype != BESIEGE_PLANET) {
                        continue;
                    }
                    if (transaction[i].x != nampla.x) {
                        continue;
                    }
                    if (transaction[i].y != nampla.y) {
                        continue;
                    }
                    if (transaction[i].z != nampla.z) {
                        continue;
                    }
                    if (transaction[i].pn != nampla.pn) {
                        continue;
                    }
                    if (transaction[i].number2 != species_number) {
                        continue;
                    }

                    alien_number = transaction[i].number1;

                    /* Define a 'detection' transaction. */
                    if (num_transactions == MAX_TRANSACTIONS) {
                        fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
                        exit(-1);
                    }

                    n = num_transactions++;
                    transaction[n].ttype  = DETECTION_DURING_SIEGE;
                    transaction[n].value = 1;   /* Landing. */
                    strcpy(transaction[n].name1, nampla.name);
                    strcpy(transaction[n].name2, ship_name(ship));
                    strcpy(transaction[n].name3, species.name);
                    transaction[n].number3 = alien_number;
                }
            }

            if (rnd(100) <= siege_effectiveness) {
                /* Ship doesn't know if it was detected. */
                log_string(" may have been detected by the besiegers when it landed on PL ");
                log_string(nampla.name);
            }else {
                /* Ship knows whether or not it was detected. */
                if (landing_detected) {
                    log_string(" was detected by the besiegers when it landed on PL ");
                    log_string(nampla.name);
                }else {
                    log_string(" landed on PL ");
                    log_string(nampla.name);
                    log_string(" without being detected by the besiegers");
                }
            }
        }
    }else {
        if (first_pass) {
            log_string(" will land on PL ");
        }else{
            log_string(" landed on PL ");
        }
        log_string(nampla.name);
    }

    log_string(".\n");
}

//*************************************************************************
// do_locs.c


extern int species_number, species_index, num_locs;

extern struct galaxy_data   galaxy;
extern struct species_data *species;
extern struct nampla_data * nampla_base, *nampla;
extern struct ship_data *   ship_base, *ship;
extern struct sp_loc_data   loc[MAX_LOCATIONS];


/* This routine will create the "loc" array based on current species' data. */
do_locations() {
    int i;


    num_locs = 0;

    for (species_number = 1; species_number <= galaxy.num_species; species_number++) {
        if (!data_in_memory[species_number - 1]) {
            continue;
        }

        species     = &spec_data[species_number - 1];
        nampla_base = namp_data[species_number - 1];
        ship_base   = ship_data[species_number - 1];

        nampla = nampla_base - 1;
        for (i = 0; i < species.num_namplas; i++) {
            ++nampla;

            if (nampla.pn == 99) {
                continue;
            }

            if (nampla.status & POPULATED) {
                add_location(nampla.x, nampla.y, nampla.z);
            }
        }

        ship = ship_base - 1;
        for (i = 0; i < species.num_ships; i++) {
            ++ship;

            if (ship.pn == 99) {
                continue;
            }
            if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
                continue;
            }

            add_location(ship.x, ship.y, ship.z);
        }
    }
}



add_location(x, y, z)

char x, y, z;

{
    int i;


    for (i = 0; i < num_locs; i++) {
        if (loc[i].x != x) {
            continue;
        }
        if (loc[i].y != y) {
            continue;
        }
        if (loc[i].z != z) {
            continue;
        }
        if (loc[i].s != species_number) {
            continue;
        }

        return;         /* This location is already in list for this species. */
    }

    /* Add new location to the list. */
    loc[num_locs].x = x;
    loc[num_locs].y = y;
    loc[num_locs].z = z;
    loc[num_locs].s = species_number;

    ++num_locs;
    if (num_locs < MAX_LOCATIONS) {
        return;
    }

    fprintf(stderr, "\n\n\tInternal error. Overflow of 'loc' arrays!\n\n");
    exit(-1);
}

//*************************************************************************
// do_mes.c


extern int abbr_type, first_pass, species_number,
           g_spec_number, num_transactions, end_of_file;
extern char input_line[256], g_spec_name[32],
            *input_line_pointer;
extern FILE *log_file, *input_file;
extern struct species_data *species;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_MESSAGE_command() {
    int i, message_number, message_fd, bad_species,
        unterminated_message;

    char c1, c2, c3, filename[32];

    FILE *message_file;


    /* Get destination of message. */
    if (!get_species_name()) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid species name in MESSAGE command.\n");
        bad_species = true;
    }else {
        bad_species = false;
    }

    /* Generate a random number, create a filename with it, and use it to
     *  store message. */
    if (!first_pass && !bad_species) {
        for {
            /* Generate a random filename. */
            message_number = rnd(32000);
            sprintf(filename, "m%d.msg\0", message_number);

            /* Make sure that this filename is not already in use. */
            message_fd = open(filename, 0);
            if (message_fd < 0) {
                break;
            }

            /* File already exists. Try again. */
            close(message_fd);
        }

        message_file = fopen(filename, "w");
        if (message_file == NULL) {
            fprintf(stderr,
                    "\n\n!!! Cannot open message file '%s' for writing !!!\n\n",
                    filename);
            exit(-1);
        }
    }

    /* Copy message to file. */
    unterminated_message = false;
    for {
        /* Read next line. */
        input_line_pointer = fgets(input_line, 256, input_file);
        if (input_line_pointer == NULL) {
            unterminated_message = true;
            end_of_file          = true;
            break;
        }

        skip_whitespace();

        c1 = *input_line_pointer++;
        c2 = *input_line_pointer++;
        c3 = *input_line_pointer;

        c1 = toupper(c1);
        c2 = toupper(c2);
        c3 = toupper(c3);

        if (c1 == 'Z' && c2 == 'Z' && c3 == 'Z') {
            break;
        }

        if (!first_pass && !bad_species) {
            fputs(input_line, message_file);
        }
    }

    if (bad_species) {
        return;
    }

    /* Log the result. */
    log_string("    A message was sent to SP ");
    log_string(g_spec_name);
    log_string(".\n");

    if (unterminated_message) {
        log_string("  ! WARNING: Message was not properly terminated with ZZZ!");
        log_string(" Any orders that follow the message will be assumed");
        log_string(" to be part of the message and will be ignored!\n");
    }

    if (first_pass) {
        return;
    }

    fclose(message_file);

    /* Define this message transaction and add to list of transactions. */
    if (num_transactions == MAX_TRANSACTIONS) {
        fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
        exit(-1);
    }

    i = num_transactions++;
    transaction[i].ttype    = MESSAGE_TO_SPECIES;
    transaction[i].value   = message_number;
    transaction[i].number1 = species_number;
    strcpy(transaction[i].name1, species.name);
    transaction[i].number2 = g_spec_number;
    strcpy(transaction[i].name2, g_spec_name);
}


//*************************************************************************
// do_name.c



extern int x, y, z, pn, num_stars, nampla_index,
           species_index;
extern char input_line[256], upper_name[32],
            original_name[32], *input_line_pointer;
extern FILE *log_file;
extern struct species_data *species;
extern struct star_data *   star;
extern struct planet_data * planet_base;
extern struct nampla_data * nampla_base, *nampla;


do_NAME_command() {
    int i, found, name_length, unused_nampla_available;

    char upper_nampla_name[32], *original_line_pointer;

    struct planet_data *planet;
    struct nampla_data *unused_nampla;


    /* Get x y z coordinates. */
    found = get_location();
    if (!found || nampla != NULL || pn == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid coordinates in NAME command.\n");
        return;
    }

    /* Get planet abbreviation. */
    skip_whitespace();
    original_line_pointer = input_line_pointer;
    if (get_class_abbr() != PLANET_ID) {
        /* Check if PL was mispelled (i.e, "PT" or "PN"). Otherwise
         *      assume that it was accidentally omitted. */
        if (tolower(*original_line_pointer) != 'p' ||
            isalnum(*(original_line_pointer + 2))) {
            input_line_pointer = original_line_pointer;
        }
    }

    /* Get planet name. */
    name_length = get_name();
    if (name_length < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid planet name in NAME command.\n");
        return;
    }

    /* Search existing namplas for name and location. */
    found = false;
    unused_nampla_available = false;
    nampla = nampla_base - 1;
    for (nampla_index = 0; nampla_index < species.num_namplas; nampla_index++) {
        ++nampla;

        if (nampla.pn == 99) {
            /* We can re-use this nampla rather than append a new one. */
            unused_nampla           = nampla;
            unused_nampla_available = true;
            continue;
        }

        /* Check if a named planet already exists at this location. */
        if (nampla.x == x && nampla.y == y && nampla.z == z &&
            nampla.pn == pn) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! The planet at these coordinates already has a name.\n");
            return;
        }

        /* Make upper case copy of nampla name. */
        for (i = 0; i < 32; i++) {
            upper_nampla_name[i] = toupper(nampla.name[i]);
        }

        /* Compare names. */
        if (strcmp(upper_nampla_name, upper_name) == 0) {
            found = true;
            break;
        }
    }

    if (found) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Planet in NAME command already exists.\n");
        return;
    }

    /* Add new nampla to database for this species. */
    if (unused_nampla_available) {
        nampla = unused_nampla;
    }else{
        ++num_new_namplas[species_index];
        if (num_new_namplas[species_index] > NUM_EXTRA_NAMPLAS) {
            fprintf(stderr, "\n\n\tInsufficient memory for new planet name:\n");
            fprintf(stderr, "\n\t%s\n", input_line);
            exit(-1);
        }
        nampla = nampla_base + species.num_namplas;
        species.num_namplas += 1;
        delete_nampla(nampla);          /* Set everything to zero. */
    }

    /* Initialize new nampla. */
    strcpy(nampla.name, original_name);
    nampla.x            = x;
    nampla.y            = y;
    nampla.z            = z;
    nampla.pn           = pn;
    nampla.status       = COLONY;
    nampla.planet_index = star.planet_index + pn - 1;
    planet          = planet_base + (long)nampla.planet_index;
    nampla.message = planet.message;

    /* Everything else was set to zero in above call to 'delete_nampla'. */

    /* Mark sector as having been visited. */
    star_visited(x, y, z);

    /* Log result. */
    log_string("    Named PL ");  log_string(nampla.name);
    log_string(" at ");  log_int(nampla.x);  log_char(' ');
    log_int(nampla.y);  log_char(' ');  log_int(nampla.z);
    log_string(", planet #");  log_int(nampla.pn);
    log_string(".\n");
}

//*************************************************************************
// do_neutral.c


extern int   abbr_type, g_spec_number;
extern char  input_line[256], g_spec_name[32];
extern FILE *log_file;
extern struct species_data *species;


do_NEUTRAL_command() {
    int i, array_index, bit_number;

    long bit_mask;


    /* See if declaration is for all species. */
    if (get_value()) {
        bit_mask = 0;
        for (i = 0; i < NUM_CONTACT_WORDS; i++) {
            species.enemy[i] = bit_mask;       /* Clear all enemy bits. */
            species.ally[i]  = bit_mask;       /* Clear all ally bits. */
        }
    }else {
        /* Get name of species. */
        if (!get_species_name()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid or missing argument in NEUTRAL command.\n");
            return;
        }

        /* Get array index and bit mask. */
        array_index = (g_spec_number - 1) / 32;
        bit_number  = (g_spec_number - 1) % 32;
        bit_mask    = 1 << bit_number;

        /* Clear the appropriate bit. */
        species.enemy[array_index] &= ~bit_mask; /* Clear enemy bit. */
        species.ally[array_index]  &= ~bit_mask; /* Clear ally bit. */
    }

    /* Log the result. */
    log_string("    Neutrality was declared towards ");
    if (bit_mask == 0) {
        log_string("ALL species");
    }else{
        log_string("SP ");
        log_string(g_spec_name);
    }
    log_string(".\n");
}

//*************************************************************************
// do_orbit.c

extern int  first_pass, num_stars;
extern long value;
extern char input_line[256], original_line[256],
            *input_line_pointer;
extern FILE *log_file;
extern struct galaxy_data   galaxy;
extern struct star_data *   star_base, *star;
extern struct species_data *species;
extern struct nampla_data * nampla, *nampla_base;
extern struct ship_data *   ship;


do_ORBIT_command() {
    int i, found, specified_planet_number;

    char *original_line_pointer;


    /* Get the ship. */
    original_line_pointer = input_line_pointer;
    found = get_ship();
    if (!found) {
        /* Check for missing comma or tab after ship name. */
        input_line_pointer = original_line_pointer;
        fix_separator();
        found = get_ship();
        if (!found) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid ship name in ORBIT command.\n");
            return;
        }
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship is still under construction.\n");
        return;
    }

    if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
        return;
    }

    /* Make sure this ship didn't just arrive via a MOVE command. */
    if (ship.just_jumped == 50) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! ORBIT not allowed immediately after a MOVE!\n");
        return;
    }

    /* Make sure ship is not salvage of a disbanded colony. */
    if (disbanded_ship(ship)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! This ship is salvage of a disbanded colony!\n");
        return;
    }

    /* Get the planet. */
    specified_planet_number = get_value();

get_planet:

    if (specified_planet_number) {
        found = false;
        specified_planet_number = value;
        for (i = 0; i < num_stars; i++) {
            star = star_base + i;

            if (star.x != ship.x) {
                continue;
            }
            if (star.y != ship.y) {
                continue;
            }
            if (star.z != ship.z) {
                continue;
            }

            if (specified_planet_number >= 1 &&
                specified_planet_number <= star.num_planets) {
                found = true;
            }

            break;
        }

        if (!found) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid planet in ORBIT command.\n");
            return;
        }

        ship.pn = specified_planet_number;

        goto finish_up;
    }

    found = get_location();
    if (!found || nampla == NULL) {
        if (ship.status == IN_ORBIT || ship.status == ON_SURFACE) {
            /* Player forgot to specify planet. Use the one it's already at. */
            specified_planet_number = ship.pn;
            value = specified_planet_number;
            goto get_planet;
        }

        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Invalid or missing planet in ORBIT command.\n");
        return;
    }

    /* Make sure the ship and the planet are in the same star system. */
    if (ship.x != nampla.x || ship.y != nampla.y || ship.z != nampla.z) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship and planet are not in the same sector.\n");
        return;
    }

    /* Move the ship. */
    ship.pn = nampla.pn;

finish_up:

    ship.status = IN_ORBIT;

    /* If a planet number is being used, see if it has a name.  If so,
     *  use the name. */
    if (specified_planet_number) {
        for (i = 0; i < species.num_namplas; i++) {
            nampla = nampla_base + i;

            if (nampla.x != ship.x) {
                continue;
            }
            if (nampla.y != ship.y) {
                continue;
            }
            if (nampla.z != ship.z) {
                continue;
            }
            if (nampla.pn != specified_planet_number) {
                continue;
            }

            specified_planet_number = 0;
            break;
        }
    }

    /* Log result. */
    log_string("    ");
    log_string(ship_name(ship));
    if (first_pass) {
        log_string(" will enter orbit around ");
    }else{
        log_string(" entered orbit around ");
    }

    if (specified_planet_number) {
        log_string("planet number ");
        log_int(specified_planet_number);
    }else {
        log_string("PL ");
        log_string(nampla.name);
    }

    log_string(".\n");
}

//*************************************************************************
// do_prod.c


int last_planet_produced = false;

int shipyard_built, shipyard_capacity;

extern int nampla_index, doing_production, first_pass, next_nampla_index,
           planet_data_modified, species_number, num_transactions;
extern long balance, raw_material_units, production_capacity,
            EU_spending_limit;
extern char  production_done[1000], input_line[256], upper_name[32];
extern FILE *log_file;

extern struct planet_data * planet_base, *planet, *home_planet;
extern struct species_data *species;
extern struct nampla_data * nampla_base, *nampla, *next_nampla;
extern struct ship_data *   ship_base;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_PRODUCTION_command(missing_production_order)

int missing_production_order;

{
    int i, j, abbr_type, name_length, found, alien_number, under_siege,
        siege_percent_effectiveness, new_alien, num_siege_ships,
        mining_colony, resort_colony, special_colony, ship_index,
        enemy_on_same_planet, trans_index, production_penalty,
        ls_needed, shipyards_for_this_species;

    char upper_nampla_name[32];

    long n, RMs_produced, num_bytes, total_siege_effectiveness,
         siege_effectiveness[MAX_SPECIES + 1], EUs_available_for_siege,
         EUs_for_distribution, EUs_for_this_species, total_EUs_stolen,
         special_production, pop_units_here[MAX_SPECIES + 1],
         alien_pop_units, total_alien_pop_here, total_besieged_pop,
         ib_for_this_species, ab_for_this_species, total_ib, total_ab,
         total_effective_tonnage;

    struct species_data *alien;
    struct nampla_data * alien_nampla_base, *alien_nampla;
    struct ship_data *   alien_ship_base, *alien_ship, *ship;


    if (doing_production) {
        /* Terminate production for previous planet. */
        if (last_planet_produced) {
            transfer_balance();
            last_planet_produced = false;
        }

        /* Give gamemaster option to abort. */
        if (first_pass) {
            gamemaster_abort_option();
        }
        log_char('\n');
    }

    doing_production = true;

    if (missing_production_order) {
        nampla       = next_nampla;
        nampla_index = next_nampla_index;

        goto got_nampla;
    }

    /* Get PL abbreviation. */
    abbr_type = get_class_abbr();

    if (abbr_type != PLANET_ID) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid planet name in PRODUCTION command.\n");
        return;
    }

    /* Get planet name. */
    name_length = get_name();

    /* Search all namplas for name. */
    found  = false;
    nampla = nampla_base - 1;
    for (nampla_index = 0; nampla_index < species.num_namplas; nampla_index++) {
        ++nampla;

        if (nampla.pn == 99) {
            continue;
        }

        /* Make upper case copy of nampla name. */
        for (i = 0; i < 32; i++) {
            upper_nampla_name[i] = toupper(nampla.name[i]);
        }

        /* Compare names. */
        if (strcmp(upper_nampla_name, upper_name) == 0) {
            found = true;
            break;
        }
    }

    if (!found) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid planet name in PRODUCTION command.\n");
        return;
    }

    /* Check if production was already done for this planet. */
    if (production_done[nampla_index]) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! More than one PRODUCTION command for planet.\n");
        return;
    }
    production_done[nampla_index] = true;

    /* Check if this colony was disbanded. */
    if (nampla.status & DISBANDED_COLONY) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Production orders cannot be given for a disbanded colony!\n");
        return;
    }

got_nampla:

    last_planet_produced = true;
    shipyard_built       = false;
    shipyard_capacity    = nampla.shipyards;

    /* See if this is a mining or resort colony. */
    mining_colony  = false;
    resort_colony  = false;
    special_colony = false;
    if (nampla.status & MINING_COLONY) {
        mining_colony  = true;
        special_colony = true;
    }else if (nampla.status & RESORT_COLONY) {
        resort_colony  = true;
        special_colony = true;
    }

    /* Get planet data for this nampla. */
    planet = planet_base + (long)nampla.planet_index;

    /* Check if fleet maintenance cost is so high that riots ensued. */
    i = 0;
    j = (species.fleet_percent_cost - 10000) / 100;
    if (rnd(100) <= j) {
        log_string("!!! WARNING! Riots on PL ");
        log_string(nampla.name);
        log_string(" due to excessive and unpopular military build-up reduced ");

        if (mining_colony || !special_colony) {
            log_string("mining base by ");
            i = rnd(j);
            log_int(i);  log_string(" percent ");
            nampla.mi_base -= (i * nampla.mi_base) / 100;
        }

        if (resort_colony || !special_colony) {
            if (i) {
                log_string("and ");
            }
            log_string("manufacturing base by ");
            i = rnd(j);
            log_int(i);  log_string(" percent");
            nampla.ma_base -= (i * nampla.ma_base) / 100;
        }
        log_string("!\n\n");
    }

    /* Calculate "balance" available for spending and create pseudo
     *  "checking account". */
    ls_needed = life_support_needed(species, home_planet, planet);

    if (ls_needed == 0) {
        production_penalty = 0;
    }else{
        production_penalty = (100 * ls_needed) / species.tech_level[LS];
    }

    RMs_produced =
        (10L * (long)species.tech_level[MI] * (long)nampla.mi_base)
        / (long)planet.mining_difficulty;
    RMs_produced
        -= (production_penalty * RMs_produced) / 100;
    RMs_produced
        = (((long)planet.econ_efficiency * RMs_produced) + 50) / 100;

    if (special_colony) {
        /* RMs just 'sitting' on the planet cannot be converted to EUs on a
         *      mining colony, and cannot create a 'balance' on a resort
         *      colony. */
        raw_material_units = 0;
    }else{
        raw_material_units = RMs_produced + nampla.item_quantity[RM];
    }

    production_capacity =
        ((long)species.tech_level[MA] * (long)nampla.ma_base) / 10L;
    production_capacity
        -= (production_penalty * production_capacity) / 100;
    production_capacity
        = (((long)planet.econ_efficiency * production_capacity) + 50) / 100;

    balance = (raw_material_units > production_capacity)
              ? production_capacity : raw_material_units;

    if (species.fleet_percent_cost > 10000) {
        n = 10000;
    }else{
        n = species.fleet_percent_cost;
    }

    if (special_colony) {
        EU_spending_limit = 0;
    }else{
        /* Only excess RMs may be recycled. */
        nampla.item_quantity[RM] = raw_material_units - balance;

        balance                -= ((n * balance) + 5000) / 10000;
        raw_material_units      = balance;
        production_capacity     = balance;
        EUs_available_for_siege = balance;
        if (nampla.status & HOME_PLANET) {
            if (species.hp_original_base != 0) { /* HP was bombed. */
                EU_spending_limit = 4 * balance;  /* Factor = 4 + 1 = 5. */
            }else {
                EU_spending_limit = species.econ_units;
            }
        }else {
            EU_spending_limit = balance;
        }
    }

    /* Log what was done. Balances for mining and resort colonies will always
     *  be zero and should not be printed. */
    log_string("  Start of production on PL ");  log_string(nampla.name);
    log_char('.');
    if (!special_colony) {
        log_string(" (Initial balance is ");  log_long(balance);
        log_string(".)");
    }
    log_char('\n');

    /* If this IS a mining or resort colony, convert RMs or production capacity
     *  to EUs. */
    if (mining_colony) {
        special_production  = (2 * RMs_produced) / 3;
        special_production -= ((n * special_production) + 5000) / 10000;
        log_string("    Mining colony ");
    }else if (resort_colony) {
        special_production  = (2 * production_capacity) / 3;
        special_production -= ((n * special_production) + 5000) / 10000;
        log_string("    Resort colony ");
    }

    if (special_colony) {
        log_string(nampla.name);  log_string(" generated ");
        log_long(special_production);  log_string(" economic units.\n");

        EUs_available_for_siege = special_production;
        species.econ_units    += special_production;

        if (mining_colony && !first_pass) {
            planet.mining_difficulty += RMs_produced / 150;
            planet_data_modified       = true;
        }
    }

    /* Check if this planet is under siege. */
    nampla.siege_eff         = 0;
    under_siege               = false;
    alien_number              = 0;
    num_siege_ships           = 0;
    total_siege_effectiveness = 0;
    enemy_on_same_planet      = false;
    total_alien_pop_here      = 0;
    for (i = 1; i <= MAX_SPECIES; i++) {
        siege_effectiveness[i] = 0;
        pop_units_here[i]      = 0;
    }

    for (trans_index = 0; trans_index < num_transactions; trans_index++) {
        /* Check if this is a siege of this nampla. */
        if (transaction[trans_index].ttype != BESIEGE_PLANET) {
            continue;
        }
        if (transaction[trans_index].x != nampla.x) {
            continue;
        }
        if (transaction[trans_index].y != nampla.y) {
            continue;
        }
        if (transaction[trans_index].z != nampla.z) {
            continue;
        }
        if (transaction[trans_index].pn != nampla.pn) {
            continue;
        }
        if (transaction[trans_index].number2 != species_number) {
            continue;
        }

        /* Check if alien ship is still in the same star system as the
         *      planet. */
        if (alien_number != transaction[trans_index].number1) {
            /* First transaction for this alien. */
            alien_number = transaction[trans_index].number1;
            if (!data_in_memory[alien_number - 1]) {
                fprintf(stderr, "\n\tData for species #%d should be in memory but is not!\n\n",
                        alien_number);
                exit(-1);
            }
            alien             = &spec_data[alien_number - 1];
            alien_nampla_base = namp_data[alien_number - 1];
            alien_ship_base   = ship_data[alien_number - 1];

            new_alien = true;
        }

        /* Find the alien ship. */
        found      = false;
        alien_ship = alien_ship_base - 1;
        for (i = 0; i < alien.num_ships; i++) {
            ++alien_ship;

            if (alien_ship.pn == 99) {
                continue;
            }

            if (strcmp(alien_ship.name, transaction[trans_index].name3) == 0) {
                found = true;
                break;
            }
        }

        /* Check if alien ship is still at the siege location. */
        if (!found) {
            continue;                   /* It must have jumped away and self-
                                         *      destructed, or was recycled. */
        }
        if (alien_ship.x != nampla.x) {
            continue;
        }
        if (alien_ship.y != nampla.y) {
            continue;
        }
        if (alien_ship.z != nampla.z) {
            continue;
        }
        if (alien_ship.class == TR) {
            continue;
        }

        /* This nampla is under siege. */
        if (!under_siege) {
            log_string("\n    WARNING! PL ");  log_string(nampla.name);
            log_string(" is under siege by the following:\n      ");
            under_siege = true;
        }

        if (num_siege_ships++ > 0) {
            log_string(", ");
        }
        if (new_alien) {
            log_string(alien.name);  log_char(' ');
            new_alien = false;

            /* Check if this alien has a colony on the same planet. */
            alien_nampla = alien_nampla_base - 1;
            for (i = 0; i < alien.num_namplas; i++) {
                ++alien_nampla;

                if (alien_nampla.x != nampla.x) {
                    continue;
                }
                if (alien_nampla.y != nampla.y) {
                    continue;
                }
                if (alien_nampla.z != nampla.z) {
                    continue;
                }
                if (alien_nampla.pn != nampla.pn) {
                    continue;
                }

                /* Enemy population that will count for both detection AND
                 *      assimilation. */
                alien_pop_units = alien_nampla.mi_base
                                  + alien_nampla.ma_base
                                  + alien_nampla.IUs_to_install
                                  + alien_nampla.AUs_to_install;

                /* Any base over 200.0 has only 5% effectiveness. */
                if (alien_pop_units > 2000) {
                    alien_pop_units = (alien_pop_units - 2000) / 20 + 2000;
                }

                /* Enemy population that counts ONLY for detection. */
                n = alien_nampla.pop_units
                    + alien_nampla.item_quantity[CU]
                    + alien_nampla.item_quantity[PD];

                if (alien_pop_units > 0) {
                    enemy_on_same_planet         = true;
                    pop_units_here[alien_number] = alien_pop_units;
                    total_alien_pop_here        += alien_pop_units;
                }else if (n > 0) {
                    enemy_on_same_planet = true;
                }

                if (alien_nampla.item_quantity[PD] == 0) {
                    continue;
                }

                log_string("planetary defenses of PL ");
                log_string(alien_nampla.name);
                log_string(", ");

                n = (4 * alien_nampla.item_quantity[PD]) / 5;
                n = (n * (long)alien.tech_level[ML])
                    / ((long)species.tech_level[ML] + 1);
                total_siege_effectiveness         += n;
                siege_effectiveness[alien_number] += n;
            }
        }
        log_string(ship_name(alien_ship));

        /* Determine the number of planets that this ship is besieging. */
        n = 0;
        for (j = 0; j < num_transactions; j++) {
            if (transaction[j].ttype != BESIEGE_PLANET) {
                continue;
            }
            if (transaction[j].number1 != alien_number) {
                continue;
            }
            if (strcmp(transaction[j].name3, alien_ship.name) != 0) {
                continue;
            }

            ++n;
        }

        /* Determine the effectiveness of this ship on the siege. */
        if (alien_ship.ttype == STARBASE) {
            i = alien_ship.tonnage;    /* One quarter of normal ships. */
        }else {
            i = 4 * (int)alien_ship.tonnage;
        }

        i = (i * (int)alien.tech_level[ML])
            / ((int)species.tech_level[ML] + 1);

        i /= n;

        total_siege_effectiveness         += i;
        siege_effectiveness[alien_number] += i;
    }

    if (under_siege) {
        log_string(".\n");
    }else{
        return;
    }

    /* Determine percent effectiveness of the siege. */
    total_effective_tonnage = 2500 * total_siege_effectiveness;

    if (nampla.mi_base + nampla.ma_base == 0) {
        siege_percent_effectiveness = -9999;    /* New colony with nothing
                                                 *  installed yet. */
    }else {
        siege_percent_effectiveness = total_effective_tonnage
                                      / ((((long)species.tech_level[MI] * (long)nampla.mi_base)
                                          + ((long)species.tech_level[MA] * (long)nampla.ma_base))
                                         / 10L);
    }

    if (siege_percent_effectiveness > 95) {
        siege_percent_effectiveness = 95;
    }else if (siege_percent_effectiveness == -9999) {
        log_string("      However, although planet is populated, it has no economic base.\n\n");
        return;
    }else if (siege_percent_effectiveness < 1) {
        log_string("      However, because of the weakness of the siege, it was completely ineffective!\n\n");
        return;
    }

    if (enemy_on_same_planet) {
        nampla.siege_eff = -siege_percent_effectiveness;
    }else{
        nampla.siege_eff = siege_percent_effectiveness;
    }

    log_string("      The siege is approximately ");
    log_int(siege_percent_effectiveness);
    log_string("% effective.\n");

    /* Add siege EU transfer(s). */
    EUs_for_distribution
        = (siege_percent_effectiveness * EUs_available_for_siege) / 100;

    total_EUs_stolen = 0;

    for (alien_number = 1; alien_number <= MAX_SPECIES; alien_number++) {
        n = siege_effectiveness[alien_number];
        if (n < 1) {
            continue;
        }
        alien = &spec_data[alien_number - 1];
        EUs_for_this_species
            = (n * EUs_for_distribution) / total_siege_effectiveness;
        if (EUs_for_this_species < 1) {
            continue;
        }
        total_EUs_stolen += EUs_for_this_species;
        log_string("      ");  log_long(EUs_for_this_species);
        log_string(" economic unit");
        if (EUs_for_this_species > 1) {
            log_string("s were");
        }else{
            log_string(" was");
        }
        log_string(" lost and 25% of the amount was transferred to SP ");
        log_string(alien.name);
        log_string(".\n");

        if (first_pass) {
            continue;
        }

        /* Define this transaction and add to list of transactions. */
        if (num_transactions == MAX_TRANSACTIONS) {
            fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
            exit(-1);
        }

        trans_index = num_transactions++;
        transaction[trans_index].ttype      = SIEGE_EU_TRANSFER;
        transaction[trans_index].donor     = species_number;
        transaction[trans_index].recipient = alien_number;
        transaction[trans_index].value     = EUs_for_this_species / 4;
        transaction[trans_index].x         = nampla.x;
        transaction[trans_index].y         = nampla.y;
        transaction[trans_index].z         = nampla.z;
        transaction[trans_index].number1   = siege_percent_effectiveness;
        strcpy(transaction[trans_index].name1, species.name);
        strcpy(transaction[trans_index].name2, alien.name);
        strcpy(transaction[trans_index].name3, nampla.name);
    }
    log_char('\n');

    /* Correct balances. */
    if (special_colony) {
        species.econ_units -= total_EUs_stolen;
    }else{
        if (check_bounced(total_EUs_stolen)) {
            fprintf(stderr, "\nWARNING! Internal error! Should never reach this point!\n\n");
            exit(-1);
        }
    }

    if (!enemy_on_same_planet) {
        return;
    }

    /* All ships currently under construction may be detected by the besiegers
     *  and destroyed. */
    for (ship_index = 0; ship_index < species.num_ships; ship_index++) {
        ship = ship_base + ship_index;

        if (ship.status == UNDER_CONSTRUCTION &&
            ship.x == nampla.x &&
            ship.y == nampla.y &&
            ship.z == nampla.z &&
            ship.pn == nampla.pn) {
            if (rnd(100) > siege_percent_effectiveness) {
                continue;
            }

            log_string("      ");
            log_string(ship_name(ship));
            log_string(", under construction when the siege began, was detected by the besiegers and destroyed!\n");
            if (!first_pass) {
                delete_ship(ship);
            }
        }
    }

    /* Check for assimilation. */
    if (nampla.status & HOME_PLANET) {
        return;
    }
    if (total_alien_pop_here < 1) {
        return;
    }

    total_besieged_pop = nampla.mi_base + nampla.ma_base
                         + nampla.IUs_to_install + nampla.AUs_to_install;

    /* Any base over 200.0 has only 5% effectiveness. */
    if (total_besieged_pop > 2000) {
        total_besieged_pop = (total_besieged_pop - 2000) / 20 + 2000;
    }

    if (total_besieged_pop / total_alien_pop_here >= 5) {
        return;
    }
    if (siege_percent_effectiveness < 95) {
        return;
    }

    log_string("      PL ");  log_string(nampla.name);
    log_string(" has become assimilated by the besieging species");
    log_string(" and is no longer under your control.\n\n");

    total_ib  = nampla.mi_base; /* My stupid compiler can't add an int and
                                  *     an unsigned short. */
    total_ib += nampla.IUs_to_install;
    total_ab  = nampla.ma_base;
    total_ab += nampla.AUs_to_install;

    for (alien_number = 1; alien_number <= MAX_SPECIES; alien_number++) {
        n = pop_units_here[alien_number];
        if (n < 1) {
            continue;
        }

        shipyards_for_this_species
            = (n * nampla.shipyards) / total_alien_pop_here;

        ib_for_this_species
                  = (n * total_ib) / total_alien_pop_here;
        total_ib -= ib_for_this_species;

        ab_for_this_species
                  = (n * total_ab) / total_alien_pop_here;
        total_ab -= ab_for_this_species;

        if (ib_for_this_species == 0 && ab_for_this_species == 0) {
            continue;
        }

        if (first_pass) {
            continue;
        }

        /* Define this transaction and add to list of transactions. */
        if (num_transactions == MAX_TRANSACTIONS) {
            fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
            exit(-1);
        }

        trans_index = num_transactions++;
        transaction[trans_index].ttype    = ASSIMILATION;
        transaction[trans_index].value   = alien_number;
        transaction[trans_index].x       = nampla.x;
        transaction[trans_index].y       = nampla.y;
        transaction[trans_index].z       = nampla.z;
        transaction[trans_index].pn      = nampla.pn;
        transaction[trans_index].number1 = ib_for_this_species / 2;
        transaction[trans_index].number2 = ab_for_this_species / 2;
        transaction[trans_index].number3 = shipyards_for_this_species;
        strcpy(transaction[trans_index].name1, species.name);
        strcpy(transaction[trans_index].name2, nampla.name);
    }

    /* Erase the original colony. */
    balance                = 0;
    EU_spending_limit      = 0;
    raw_material_units     = 0;
    production_capacity    = 0;
    nampla.mi_base        = 0;
    nampla.ma_base        = 0;
    nampla.IUs_to_install = 0;
    nampla.AUs_to_install = 0;
    nampla.pop_units      = 0;
    nampla.siege_eff      = 0;
    nampla.status         = COLONY;
    nampla.shipyards      = 0;
    nampla.hiding         = 0;
    nampla.hidden         = 0;
    nampla.use_on_ambush  = 0;

    for (i = 0; i < MAX_ITEMS; i++) {
        nampla.item_quantity[i] = 0;
    }
}


//*************************************************************************
// do_recy.c


extern int ship_index, doing_production, correct_spelling_required,
           abbr_index;
extern long  value, raw_material_units, balance, EU_spending_limit;
extern char  input_line[256];
extern FILE *log_file;

extern struct species_data *species;
extern struct nampla_data * nampla;
extern struct ship_data *   ship, *ship_base;


do_RECYCLE_command() {
    int i, class, cargo;

    long recycle_value, original_cost, units_available;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get number of items to recycle. */
    i = get_value();

    if (i == 0) {
        goto recycle_ship;      /* Not an item. */
    }
    /* Get class of item. */
    class = get_class_abbr();

    if (class != ITEM_CLASS) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid item class in RECYCLE command.\n");
        return;
    }
    class = abbr_index;

    /* Make sure value is meaningful. */
    if (value == 0) {
        value = nampla.item_quantity[class];
    }
    if (value == 0) {
        return;
    }
    if (value < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid item count in RECYCLE command.\n");
        return;
    }

    /* Make sure that items exist. */
    units_available = nampla.item_quantity[class];
    if (value > units_available) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Attempt to recycle more items than are available.\n");
        return;
    }

    /* Determine recycle value. */
    if (class == TP) {
        recycle_value = (value * item_cost[class])
                        / (2L * (long)species.tech_level[BI]);
    }else if (class == RM) {
        recycle_value = value / 5L;
    }else{
        recycle_value = (value * item_cost[class]) / 2L;
    }

    /* Update inventories. */
    nampla.item_quantity[class] -= value;
    if (class == PD || class == CU) {
        nampla.pop_units += value;
    }
    species.econ_units += recycle_value;
    if (nampla.status & HOME_PLANET) {
        EU_spending_limit += recycle_value;
    }

    /* Log what was recycled. */
    log_string("    ");  log_long(value);  log_char(' ');
    log_string(item_name[class]);

    if (value > 1) {
        log_string("s were");
    }else{
        log_string(" was");
    }

    log_string(" recycled, generating ");  log_long(recycle_value);
    log_string(" economic units.\n");

    return;


recycle_ship:

    correct_spelling_required = true;
    if (!get_ship()) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship to be recycled does not exist.\n");
        return;
    }

    /* Make sure it didn't just jump. */
    if (ship.just_jumped) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship just jumped and is still in transit.\n");
        return;
    }

    /* Make sure item is at producing planet. */
    if (ship.x != nampla.x || ship.y != nampla.y ||
        ship.z != nampla.z || ship.pn != nampla.pn) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship is not at the production planet.\n");
        return;
    }

    /* Calculate recycled value. */
    if (ship.class == TR || ship.ttype == STARBASE) {
        original_cost = ship_cost[ship.class] * ship.tonnage;
    }else{
        original_cost = ship_cost[ship.class];
    }

    if (ship.ttype == SUB_LIGHT) {
        original_cost = (3 * original_cost) / 4;
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        recycle_value = (original_cost - (long)ship.remaining_cost) / 2;
    }else{
        recycle_value = (3 * original_cost * (60 - (long)ship.age)) / 200;
    }

    species.econ_units += recycle_value;
    if (nampla.status & HOME_PLANET) {
        EU_spending_limit += recycle_value;
    }

    /* Log what was recycled. */
    log_string("    ");  log_string(ship_name(ship));
    log_string(" was recycled, generating ");  log_long(recycle_value);
    log_string(" economic units");

    /* Transfer cargo, if any, from ship to planet. */
    cargo = false;
    for (i = 0; i < MAX_ITEMS; i++) {
        if (ship.item_quantity[i] > 0) {
            nampla.item_quantity[i] += ship.item_quantity[i];
            cargo = true;
        }
    }

    if (cargo) {
        log_string(". Cargo onboard ");
        log_string(ship_name(ship));
        log_string(" was first transferred to PL ");
        log_string(nampla.name);
    }

    log_string(".\n");

    /* Remove ship from inventory. */
    delete_ship(ship);
}

//*************************************************************************
// do_rep.c


extern long  value;
extern char  input_line[256], original_line[256], *input_line_pointer;
extern FILE *log_file;


extern struct species_data *species;
extern struct ship_data *   ship_base, *ship;


do_REPAIR_command() {
    int i, j, n, x, y, z, age_reduction, num_dr_units,
        total_dr_units, dr_units_used, max_age, desired_age;

    char *original_line_pointer;

    struct ship_data *damaged_ship;


    /* See if this is a "pool" repair. */
    if (get_value()) {
        x = value;
        get_value();   y = value;
        get_value();   z = value;

        if (get_value()) {
            desired_age = value;
        }else{
            desired_age = 0;
        }

        goto pool_repair;
    }

    /* Get the ship to be repaired. */
    original_line_pointer = input_line_pointer;
    if (!get_ship()) {
        /* Check for missing comma or tab after ship name. */
        input_line_pointer = original_line_pointer;
        fix_separator();
        if (!get_ship()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship to be repaired does not exist.\n");
            return;
        }
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Item to be repaired is still under construction.\n");
        return;
    }

    if (ship.age < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship or starbase is too new to repair.\n");
        return;
    }

    /* Get number of damage repair units to use. */
    if (get_value()) {
        if (value == 0) {
            num_dr_units = ship.item_quantity[DR];
        }else{
            num_dr_units = value;
        }

        age_reduction = (16 * num_dr_units) / ship.tonnage;
        if (age_reduction > ship.age) {
            age_reduction = ship.age;
            n             = age_reduction * ship.tonnage;
            num_dr_units  = (n + 15) / 16;
        }
    }else {
        age_reduction = ship.age;
        n             = age_reduction * ship.tonnage;
        num_dr_units  = (n + 15) / 16;
    }

    /* Check if sufficient units are available. */
    if (num_dr_units > ship.item_quantity[DR]) {
        if (ship.item_quantity[DR] == 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship does not have any DRs!\n");
            return;
        }
        fprintf(log_file, "! WARNING: %s", original_line);
        fprintf(log_file, "! Ship does not have %d DRs. Substituting %d for %d.\n",
                num_dr_units, ship.item_quantity[DR], num_dr_units);
        num_dr_units = ship.item_quantity[DR];
    }

    /* Check if repair will have any effect. */
    age_reduction = (16 * num_dr_units) / ship.tonnage;
    if (age_reduction < 1) {
        if (value == 0) {
            return;
        }
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! %d DRs is not enough to do a repair.\n",
                num_dr_units);
        return;
    }

    /* Log what was repaired. */
    log_string("    ");  log_string(ship_name(ship));
    log_string(" was repaired using ");
    log_int(num_dr_units);  log_char(' ');
    log_string(item_name[DR]);
    if (num_dr_units != 1) {
        log_char('s');
    }
    log_string(". Age went from ");
    log_int((int)ship.age);    log_string(" to ");
    ship.age -= age_reduction;
    if (ship.age < 0) {
        ship.age = 0;
    }
    ship.item_quantity[DR] -= num_dr_units;
    log_int((int)ship.age);
    log_string(".\n");

    return;


pool_repair:

    /* Get total number of DR units available. */
    total_dr_units = 0;
    ship           = ship_base - 1;
    for (i = 0; i < species.num_ships; i++) {
        ++ship;

        if (ship.pn == 99) {
            continue;
        }
        if (ship.x != x) {
            continue;
        }
        if (ship.y != y) {
            continue;
        }
        if (ship.z != z) {
            continue;
        }

        total_dr_units += ship.item_quantity[DR];

        ship.special = 0;
    }

    /* Repair ships, starting with the most heavily damaged. */
    dr_units_used = 0;
    for (;total_dr_units > 0;) {
        /* Find most heavily damaged ship. */
        max_age = 0;
        ship    = ship_base - 1;
        for (i = 0; i < species.num_ships; i++) {
            ++ship;

            if (ship.pn == 99) {
                continue;
            }
            if (ship.x != x) {
                continue;
            }
            if (ship.y != y) {
                continue;
            }
            if (ship.z != z) {
                continue;
            }
            if (ship.special != 0) {
                continue;
            }
            if (ship.status == UNDER_CONSTRUCTION) {
                continue;
            }

            n = ship.age;
            if (n > max_age) {
                max_age      = n;
                damaged_ship = ship;
            }
        }

        if (max_age == 0) {
            break;
        }

        damaged_ship.special = 99;

        age_reduction = max_age - desired_age;
        n             = age_reduction * damaged_ship.tonnage;
        num_dr_units  = (n + 15) / 16;

        if (num_dr_units > total_dr_units) {
            num_dr_units  = total_dr_units;
            age_reduction = (16 * num_dr_units) / damaged_ship.tonnage;
        }

        if (age_reduction < 1) {
            continue;                     /* This ship is too big. */
        }
        log_string("    ");  log_string(ship_name(damaged_ship));
        log_string(" was repaired using ");
        log_int(num_dr_units);  log_char(' ');
        log_string(item_name[DR]);
        if (num_dr_units != 1) {
            log_char('s');
        }
        log_string(". Age went from ");
        log_int((int)damaged_ship.age);    log_string(" to ");
        damaged_ship.age -= age_reduction;
        if (damaged_ship.age < 0) {
            damaged_ship.age = 0;
        }
        log_int((int)damaged_ship.age);
        log_string(".\n");

        total_dr_units -= num_dr_units;
        dr_units_used  += num_dr_units;
    }

    if (dr_units_used == 0) {
        return;
    }

    /* Subtract units used from ships at the location. */
    ship = ship_base - 1;
    for (i = 0; i < species.num_ships; i++) {
        ++ship;

        if (ship.pn == 99) {
            continue;
        }
        if (ship.x != x) {
            continue;
        }
        if (ship.y != y) {
            continue;
        }
        if (ship.z != z) {
            continue;
        }

        n = ship.item_quantity[DR];
        if (n < 1) {
            continue;
        }
        if (n > dr_units_used) {
            n = dr_units_used;
        }

        ship.item_quantity[DR] -= n;
        dr_units_used           -= n;

        if (dr_units_used == 0) {
            break;
        }
    }
}



//*************************************************************************
// do_res.c


extern int   doing_production, first_pass, abbr_index;
extern short sp_tech_level[6];
extern long  value, balance, EU_spending_limit;
extern char  input_line[256];
extern FILE *log_file;
extern struct species_data *species;


do_RESEARCH_command() {
    int n, status, tech, initial_level, current_level,
        need_amount_to_spend;

    long cost, amount_spent, cost_for_one_level, funds_remaining,
         max_funds_available;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get amount to spend. */
    status = get_value();
    need_amount_to_spend = (status == 0);       /* Sometimes players reverse
                                                 * the arguments. */
    /* Get technology. */
    if (get_class_abbr() != TECH_ID) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid or missing technology.\n");
        return;
    }
    tech = abbr_index;

    if (species.tech_knowledge[tech] == 0 && sp_tech_level[tech] == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Zero level can only be raised via TECH or TEACH.\n");
        return;
    }

    /* Get amount to spend if it was not obtained above. */
    if (need_amount_to_spend) {
        status = get_value();
    }

    if (status == 0 || value < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid or missing amount to spend!\n");
        return;
    }

do_cost:

    if (value == 0) {
        value = balance;
    }
    if (value == 0) {
        return;
    }
    cost = value;

    /* Check if sufficient funds are available. */
    if (check_bounced(cost)) {
        max_funds_available = species.econ_units;
        if (max_funds_available > EU_spending_limit) {
            max_funds_available = EU_spending_limit;
        }
        max_funds_available += balance;

        if (max_funds_available > 0) {
            fprintf(log_file, "! WARNING: %s", input_line);
            fprintf(log_file, "! Insufficient funds. Substituting %ld for %ld.\n",
                    max_funds_available, cost);
            value = max_funds_available;
            goto do_cost;
        }

        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        return;
    }

    /* Check if we already have knowledge of this technology. */
    funds_remaining = cost;
    amount_spent    = 0;
    initial_level   = sp_tech_level[tech];
    current_level   = initial_level;
    for (;current_level < species.tech_knowledge[tech];) {
        cost_for_one_level  = current_level * current_level;
        cost_for_one_level -= cost_for_one_level / 4;     /* 25% discount. */
        if (funds_remaining < cost_for_one_level) {
            break;
        }
        funds_remaining -= cost_for_one_level;
        amount_spent    += cost_for_one_level;
        ++current_level;
    }

    if (current_level > initial_level) {
        log_string("    Spent ");  log_long(amount_spent);
        log_string(" raising ");  log_string(tech_name[tech]);
        log_string(" tech level from ");  log_int(initial_level);
        log_string(" to ");  log_int(current_level);
        log_string(" using transferred knowledge.\n");

        sp_tech_level[tech] = current_level;
    }

    if (funds_remaining == 0) {
        return;
    }

    /* Increase in experience points is equal to whatever was not spent
     *  above. */
    species.tech_eps[tech] += funds_remaining;

    /* Log transaction. */
    log_string("    Spent ");  log_long(funds_remaining);
    log_string(" on ");  log_string(tech_name[tech]);
    log_string(" research.\n");
}



//*************************************************************************
// do_round.c

long    power();

char x_attacked_y[MAX_SPECIES][MAX_SPECIES];

short germ_bombs_used[MAX_SPECIES][MAX_SPECIES];


extern int log_summary, ignore_field_distorters,
           log_to_file;
extern char field_distorted[MAX_SPECIES];
extern struct species_data *c_species[MAX_SPECIES];


/* The following routine will return true if a round of combat actually
 * occurred. Otherwise, it will return false. */

int do_round(option, round_number, bat, act)

char option;

int round_number;
struct battle_data *bat;
struct action_data *act;

{
    int i, j, n, unit_index, combat_occurred, total_shots,
        attacker_index, defender_index, found, chance_to_hit,
        attacker_ml, attacker_gv, defender_ml, target_index[MAX_SHIPS],
        num_targets, header_printed, num_sp, fj_chance, shields_up,
        FDs_were_destroyed, di[3], start_unit, current_species,
        this_is_a_hijacking;

    long aux_shield_power, units_destroyed, tons, percent_decrease,
         damage_done, damage_to_ship, damage_to_shields, op1, op2,
         original_cost, recycle_value, economic_units;

    char attacker_name[64], defender_name[64];

    struct species_data *attacking_species, *defending_species;
    struct ship_data *   sh, *attacking_ship, *defending_ship;
    struct nampla_data * attacking_nampla, *defending_nampla;


    /* Clear out x_attacked_y and germ_bombs_used arrays.  They will be used
     *  to log who bombed who, or how many GWs were used. */
    num_sp = bat.num_species_here;
    for (i = 0; i < num_sp; i++) {
        for (j = 0; j < num_sp; j++) {
            x_attacked_y[i][j]    = false;
            germ_bombs_used[i][j] = 0;
        }
    }

    /* If a species has ONLY non-combatants left, then let them fight. */
    start_unit      = 0;
    total_shots     = 0;
    current_species = act.fighting_species_index[0];
    for (unit_index = 0; unit_index < act.num_units_fighting; unit_index++) {
        if (act.fighting_species_index[unit_index] != current_species) {
            if (total_shots == 0) {
                /* Convert all non-combatants, if any, to combatants. */
                for (i = start_unit; i < unit_index; i++) {
                    if (act.unit_ttype[i] == SHIP) {
                        sh          = (struct ship_data *)act.fighting_unit[i];
                        sh.special = 0;
                    }
                }
            }
            start_unit  = unit_index;
            total_shots = 0;
        }

        n = act.num_shots[unit_index];
        if (act.surprised[unit_index]) {
            n = 0;
        }
        if (act.unit_ttype[unit_index] == SHIP) {
            sh = (struct ship_data *)act.fighting_unit[unit_index];
            if (sh.special == NON_COMBATANT) {
                n = 0;
            }
        }
        total_shots += n;
    }

    /* Determine total number of shots for all species present. */
    total_shots = 0;
    for (unit_index = 0; unit_index < act.num_units_fighting; unit_index++) {
        n = act.num_shots[unit_index];
        if (act.surprised[unit_index]) {
            n = 0;
        }
        if (act.unit_ttype[unit_index] == SHIP) {
            sh = (struct ship_data *)act.fighting_unit[unit_index];
            if (sh.special == NON_COMBATANT) {
                n = 0;
            }
        }
        act.shots_left[unit_index] = n;
        total_shots += n;
    }

    /* Handle all shots. */
    header_printed  = false;
    combat_occurred = false;
    for (;total_shots > 0;) {
        /* check to make sure we arent in infinite loop
         * that usually happens when there are shots remaining
         * but the side with the shots has no more ships left*/
        for (i = 0; i < act.num_units_fighting; ++i) {
            attacking_ship = (struct ship_data *)act.fighting_unit[i];
            if (attacking_ship.age > 49 ||
                attacking_ship.status == FORCED_JUMP ||
                attacking_ship.status == JUMPED_IN_COMBAT ||
                (attacking_ship.special == NON_COMBATANT &&
                 option != GERM_WARFARE)) {
                total_shots       -= act.shots_left[i];
                act.shots_left[i] = 0;
            }
        }
        /* Determine who fires next. */
        attacker_index = rnd(act.num_units_fighting) - 1;
        if (act.unit_ttype[attacker_index] == SHIP) {
            attacking_ship = (struct ship_data *)act.fighting_unit[attacker_index];
            i = act.fighting_species_index[attacker_index];
            ignore_field_distorters = !field_distorted[i];
            sprintf(attacker_name, "%s\0", ship_name(attacking_ship));
            ignore_field_distorters = false;

            /* Check if ship can fight. */
            if (attacking_ship.age > 49) {
                continue;
            }
            if (attacking_ship.status == FORCED_JUMP) {
                continue;
            }
            if (attacking_ship.status == JUMPED_IN_COMBAT) {
                continue;
            }
            if (attacking_ship.special == NON_COMBATANT &&
                option != GERM_WARFARE) {
                continue;
            }
        }else {
            attacking_nampla = (struct nampla_data *)act.fighting_unit[attacker_index];
            sprintf(attacker_name, "PL %s\0", attacking_nampla.name);

            /* Check if planet still has defenses. */
            if (attacking_nampla.item_quantity[PD] == 0) {
                continue;
            }
        }

        /* Make sure attacker is not someone who is being taken by surprise
         *      this round. */
        if (act.surprised[attacker_index]) {
            continue;
        }

        /* Find an enemy. */
        num_targets = 0;
        i           = act.fighting_species_index[attacker_index];
        attacker_ml = c_species[i].tech_level[ML];
        attacker_gv = c_species[i].tech_level[GV];
        for (defender_index = 0; defender_index < act.num_units_fighting; defender_index++) {
            j = act.fighting_species_index[defender_index];
            if (!bat.enemy_mine[i][j]) {
                continue;
            }

            if (act.unit_ttype[defender_index] == SHIP) {
                defending_ship =
                    (struct ship_data *)act.fighting_unit[defender_index];

                if (defending_ship.age > 49) {
                    continue;                            /* Already destroyed. */
                }
                if (defending_ship.status == FORCED_JUMP) {
                    continue;
                }
                if (defending_ship.status == JUMPED_IN_COMBAT) {
                    continue;
                }
                if (defending_ship.special == NON_COMBATANT) {
                    continue;
                }
            }else {
                defending_nampla =
                    (struct nampla_data *)act.fighting_unit[defender_index];

                if (defending_nampla.item_quantity[PD] == 0 &&
                    option == PLANET_ATTACK) {
                    continue;
                }
            }

            target_index[num_targets] = defender_index;
            ++num_targets;
        }

        if (num_targets == 0) { /* Attacker has no enemies left. */
            total_shots -= act.shots_left[attacker_index];
            act.shots_left[attacker_index] = 0;
            continue;
        }

        /* Randomly choose a target. Choose the toughest of four. */
        defender_index = target_index[rnd(num_targets) - 1];
        op1            = (long)act.num_shots[defender_index]
                         * act.weapon_damage[defender_index];
        di[0] = target_index[rnd(num_targets) - 1];
        di[1] = target_index[rnd(num_targets) - 1];
        di[2] = target_index[rnd(num_targets) - 1];
        for (i = 0; i < 3; i++) {
            op2 = (long)act.num_shots[di[i]] * act.weapon_damage[di[i]];
            if (op2 > op1) {
                defender_index = di[i];
                op1            = op2;
            }
        }

        j           = act.fighting_species_index[defender_index];
        defender_ml = c_species[j].tech_level[ML];

        if (act.unit_ttype[defender_index] == SHIP) {
            defending_ship =
                (struct ship_data *)act.fighting_unit[defender_index];
            ignore_field_distorters = !field_distorted[j];
            sprintf(defender_name, "%s\0", ship_name(defending_ship));
            ignore_field_distorters = false;
        }else {
            defending_nampla =
                (struct nampla_data *)act.fighting_unit[defender_index];
            sprintf(defender_name, "PL %s\0", defending_nampla.name);
        }

        /* Print round number. */
        if (!header_printed) {
            log_string("      Now doing round ");
            log_int(round_number);
            log_string(":\n");
            header_printed = true;
        }
        int attackerGvMl = attacker_gv + attacker_ml;
        if (attackerGvMl <= 0) {
            attackerGvMl = 1;
        }

        /* Check if attacker has any forced jump units. The attacker will
         *      place more emphasis on the use of these devices if he
         *      emphasizes gravitics technology over military technology. */
        fj_chance = 50 * attacker_gv / attackerGvMl;
        if (rnd(100) < fj_chance &&
            act.unit_ttype[attacker_index] == SHIP &&
            act.unit_ttype[defender_index] == SHIP) {
            if (forced_jump_units_used(attacker_index, defender_index,
                                       &total_shots, bat, act)) {
                combat_occurred = true;
                continue;
            }
        }

        if (act.shots_left[attacker_index] == 0) {
            continue;
        }

        /* Since transports generally avoid combat, there is only a 10%
         * chance that they will be targetted, unless they are being
         * explicitly targetted. */
        i = act.fighting_species_index[attacker_index];
        j = act.fighting_species_index[defender_index];
        if (act.unit_ttype[defender_index] == SHIP &&
            defending_ship.class == TR &&
            bat.special_target[i] != TARGET_TRANSPORTS &&
            rnd(10) != 5) {
            continue;
        }

        /* If a special target has been specified, then there is a 75%
         * chance that it will be attacked if it is available. */
        if (bat.special_target[i] && rnd(100) < 76) {
            if (bat.special_target[i] == TARGET_PDS) {
                if (act.unit_ttype[defender_index] != SHIP) {
                    goto fire;
                }else{
                    continue;
                }
            }

            if (act.unit_ttype[defender_index] != SHIP) {
                continue;
            }

            if (bat.special_target[i] == TARGET_STARBASES &&
                defending_ship.class != BA) {
                continue;
            }
            if (bat.special_target[i] == TARGET_TRANSPORTS &&
                defending_ship.class != TR) {
                continue;
            }
            if (bat.special_target[i] == TARGET_WARSHIPS) {
                if (defending_ship.class == TR) {
                    continue;
                }
                if (defending_ship.class == BA) {
                    continue;
                }
            }
        }

fire:
        /* Update counts. */
        --act.shots_left[attacker_index];
        --total_shots;

        /* Since transports generally avoid combat, there is only a 10%
         * chance that they will attack. */
        if (act.unit_ttype[attacker_index] == SHIP &&
            attacking_ship.class == TR &&
            option != GERM_WARFARE &&
            rnd(10) != 5) {
            continue;
        }

        /* Fire! */
        combat_occurred = true;
        log_string("        ");  log_string(attacker_name);
        log_string(" fires on ");  log_string(defender_name);
        if (act.unit_ttype[defender_index] == NAMPLA) {
            log_string(" defenses");
        }

        int combinedMl = attacker_ml + defender_ml;
        if (combinedMl <= 0) {
            combinedMl = 1;
        }

        /* Get hit probability. The basic chance to hit is 1.5 times
         * attackers ML over the sum of attacker's and defender's ML.
         * Double this value if defender is surprised. */
        chance_to_hit = (150 * attacker_ml) / combinedMl;
        if (act.surprised[defender_index]) {
            chance_to_hit *= 2;
            shields_up     = false;
        }else {
            shields_up = true;
        }

        /* If defending ship is field-distorted, chance-to-hit is
         *      reduced by 25%. */
        j = act.fighting_species_index[defender_index];
        if (act.unit_ttype[defender_index] == SHIP &&
            field_distorted[j] &&
            defending_ship.item_quantity[FD] == defending_ship.tonnage) {
            chance_to_hit = (3 * chance_to_hit) / 4;
        }

        if (chance_to_hit > 98) {
            chance_to_hit = 98;
        }
        if (chance_to_hit < 2) {
            chance_to_hit = 2;
        }

        /* Adjust for age. */
        if (act.unit_ttype[attacker_index] == SHIP) {
            chance_to_hit -=
                (2 * attacking_ship.age * chance_to_hit) / 100;
        }

        /* Calculate damage that shot will do if it hits. */
        damage_done  = act.weapon_damage[attacker_index];
        damage_done += ((26 - rnd(51)) * damage_done) / 100;

        /* Take care of attempted annihilation and sieges. */
        if (option == PLANET_BOMBARDMENT || option == GERM_WARFARE ||
            option == SIEGE) {
            /* Indicate the action that was attempted against this nampla. */
            if (option == SIEGE) {
                act.unit_ttype[defender_index] = BESIEGED_NAMPLA;
            }else{
                act.unit_ttype[defender_index] = GENOCIDE_NAMPLA;
            }

            /* Indicate who attacked who. */
            i = act.fighting_species_index[attacker_index];
            j = act.fighting_species_index[defender_index];
            x_attacked_y[i][j] = true;

            /* Update bombardment damage. */
            if (option == PLANET_BOMBARDMENT) {
                act.bomb_damage[defender_index] += damage_done;
            }else if (option == GERM_WARFARE) {
                if (act.unit_ttype[attacker_index] == SHIP) {
                    germ_bombs_used[i][j]            += attacking_ship.item_quantity[GW];
                    attacking_ship.item_quantity[GW] = 0;
                }else {
                    germ_bombs_used[i][j] += attacking_nampla.item_quantity[GW];
                    attacking_nampla.item_quantity[GW] = 0;
                }
            }

            continue;
        }

        /* Check if shot hit. */
        if (rnd(100) <= chance_to_hit) {
            log_string(" and hits!\n");
        }else{
            log_string(" and misses!\n");
            continue;
        }

        /* Subtract damage from defender's shields, if they're up. */
        damage_to_ship = 0;
        if (shields_up) {
            if (act.unit_ttype[defender_index] == SHIP) {
                damage_to_shields =
                    ((long)defending_ship.dest_y * damage_done) / 100;
                damage_to_ship = damage_done - damage_to_shields;
                act.shield_strength_left[defender_index] -= damage_to_shields;

                /* Calculate percentage of shields left. */
                if (act.shield_strength_left[defender_index] > 0) {
                    long int defenderShieldStrength =
                        act.shield_strength[defender_index];
                    if (defenderShieldStrength <= 0) {
                        defenderShieldStrength = 1;
                    }
                    defending_ship.dest_y =
                        (100L * act.shield_strength_left[defender_index])
                        / defenderShieldStrength;
                } else{
                    defending_ship.dest_y = 0;
                }
            }else {  /* Planetary defenses. */
                act.shield_strength_left[defender_index] -= damage_done;
            }
        }

        /* See if it got through shields. */
        units_destroyed  = 0;
        percent_decrease = 0;
        if (!shields_up || act.shield_strength_left[defender_index] < 0 ||
            damage_to_ship > 0) {
            /* Get net damage to ship or PDs. */
            if (shields_up) {
                if (act.unit_ttype[defender_index] == SHIP) {
                    /* Total damage to ship is direct damage plus damage
                     *  that shields could not absorb. */
                    damage_done = damage_to_ship;
                    if (act.shield_strength_left[defender_index] < 0) {
                        damage_done -=
                            act.shield_strength_left[defender_index];
                    }
                }else {
                    damage_done = -act.shield_strength_left[defender_index];
                }
            }

            long defenderShieldStrength = act.shield_strength[defender_index];
            if (defenderShieldStrength <= 0) {
                defenderShieldStrength = 1;
            }

            percent_decrease = (50L * damage_done) / defenderShieldStrength;

            percent_decrease += ((rnd(51) - 26) * percent_decrease) / 100;
            if (percent_decrease > 100) {
                percent_decrease = 100;
            }

            if (act.unit_ttype[defender_index] == SHIP) {
                defending_ship.age += percent_decrease / 2;
                units_destroyed      = (defending_ship.age > 49);
            }else {
                units_destroyed = (percent_decrease
                                   * act.original_age_or_PDs[defender_index]) / 100L;
                if (units_destroyed > defending_nampla.item_quantity[PD]) {
                    units_destroyed = defending_nampla.item_quantity[PD];
                }
                if (units_destroyed < 1) {
                    units_destroyed = 1;
                }
                defending_nampla.item_quantity[PD] -= units_destroyed;
            }

            if (act.shield_strength_left[defender_index] < 0) {
                act.shield_strength_left[defender_index] = 0;
            }
        }

        /* See if this is a hijacking. */
        i = act.fighting_species_index[attacker_index];
        j = act.fighting_species_index[defender_index];
        if (bat.enemy_mine[i][j] == 2 && (option == DEEP_SPACE_FIGHT ||
                                           option == PLANET_ATTACK)) {
            this_is_a_hijacking = true;
        }else{
            this_is_a_hijacking = false;
        }

        attacking_species = c_species[i];
        defending_species = c_species[j];

        /* Report if anything was destroyed. */
        FDs_were_destroyed = false;
        if (units_destroyed) {
            if (act.unit_ttype[defender_index] == SHIP) {
                log_summary = true;
                log_string("        ");
                log_string(defender_name);
                if (this_is_a_hijacking) {
                    log_string(" was successfully hijacked and will generate ");

                    if (defending_ship.class == TR || defending_ship.ttype == STARBASE) {
                        original_cost = ship_cost[defending_ship.class] * defending_ship.tonnage;
                    }else{
                        original_cost = ship_cost[defending_ship.class];
                    }

                    if (defending_ship.ttype == SUB_LIGHT) {
                        original_cost = (3 * original_cost) / 4;
                    }

                    if (defending_ship.status == UNDER_CONSTRUCTION) {
                        recycle_value =
                            (original_cost - (long)defending_ship.remaining_cost) / 2;
                    }else{
                        recycle_value =
                            (3 * original_cost * (60 - act.original_age_or_PDs[defender_index])) / 200;
                    }

                    economic_units = recycle_value;

                    for (i = 0; i < MAX_ITEMS; i++) {
                        j = defending_ship.item_quantity[i];
                        if (j > 0) {
                            if (i == TP) {
                                long int techLevel_2x =
                                    2L
                                    * (long)defending_species.tech_level[BI];
                                if (techLevel_2x <= 0) {
                                    techLevel_2x = 1;
                                }
                                recycle_value = (j * item_cost[i])
                                                / techLevel_2x;
                            }else if (i == RM) {
                                recycle_value = j / 5;
                            }else{
                                recycle_value = (j * item_cost[i]) / 2;
                            }

                            economic_units += recycle_value;
                        }
                    }

                    attacking_species.econ_units += economic_units;

                    log_long(economic_units);
                    log_string(" economic units for the hijackers.\n");
                }else {
                    log_string(" was destroyed.\n");
                }

                for (i = 0; i < MAX_ITEMS; i++) {
                    if (defending_ship.item_quantity[i] > 0) {
                        /* If this is a hijacking of a field-distorted ship,
                         *  we want the true name of the hijacked species to
                         *  be announced, but we don't want any cargo to be
                         *  destroyed. */
                        if (i == FD) {
                            FDs_were_destroyed = true;
                        }
                        if (!this_is_a_hijacking) {
                            defending_ship.item_quantity[FD] = 0;
                        }
                    }
                }
                log_to_file = false;
                if (this_is_a_hijacking) {
                    log_string("          The hijacker was ");
                }else{
                    log_string("          The killing blow was delivered by ");
                }
                log_string(attacker_name);
                log_string(".\n");
                log_to_file = true;
                log_summary = false;

                total_shots -= act.shots_left[defender_index];
                act.shots_left[defender_index] = 0;
                act.num_shots[defender_index]  = 0;
            }else {
                log_summary = true;
                log_string("        ");  log_int(units_destroyed);
                if (units_destroyed > 1) {
                    log_string(" PDs on PL ");
                }else{
                    log_string(" PD on PL ");
                }
                log_string(defending_nampla.name);
                if (units_destroyed > 1) {
                    log_string(" were destroyed by ");
                }else{
                    log_string(" was destroyed by ");
                }

                log_string(attacker_name);
                log_string(".\n");

                if (defending_nampla.item_quantity[PD] == 0) {
                    total_shots -= act.shots_left[defender_index];
                    act.shots_left[defender_index] = 0;
                    act.num_shots[defender_index]  = 0;
                    log_string("        All planetary defenses have been destroyed on ");
                    log_string(defender_name);
                    log_string("!\n");
                }
                log_summary = false;
            }
        }else if (percent_decrease > 0 && !this_is_a_hijacking &&
                  act.unit_ttype[defender_index] == SHIP) {
            /* See if anything carried by the ship was also destroyed. */
            for (i = 0; i < MAX_ITEMS; i++) {
                j = defending_ship.item_quantity[i];
                if (j > 0) {
                    j = (percent_decrease * j) / 100;
                    if (j > 0) {
                        defending_ship.item_quantity[i] -= j;
                        if (i == FD) {
                            FDs_were_destroyed = true;
                        }
                    }
                }
            }
        }

        j = act.fighting_species_index[defender_index];
        if (FDs_were_destroyed && field_distorted[j] &&
            defending_ship.dest_x == 0) {
            /* Reveal the true name of the ship and the owning species. */
            log_summary = true;
            if (this_is_a_hijacking) {
                log_string("        Hijacking of ");
            }else{
                log_string("        Damage to ");
            }
            log_string(defender_name);
            log_string(" caused collapse of distortion field. Real name of ship is ");
            log_string(ship_name(defending_ship));
            log_string(" owned by SP ");
            log_string(defending_species.name);
            log_string(".\n");
            log_summary            = false;
            defending_ship.dest_x = 127;       /* Ship is now exposed. */
        }
    }

    /* No more surprises. */
    for (i = 0; i < act.num_units_fighting; i++) {
        act.surprised[i] = false;
    }

    return(combat_occurred);
}



//*************************************************************************
// do_scan.c


extern int               first_pass, test_mode;
extern char              input_line[256];
extern FILE *            log_file;
extern struct ship_data *ship;


do_SCAN_command() {
    int i, found, x, y, z;


    found = get_ship();
    if (!found) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid ship name in SCAN command.\n");
        return;
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship is still under construction.\n");
        return;
    }

    if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
        return;
    }

    /* Log the result. */
    if (first_pass) {
        log_string("    A scan will be done by ");
        log_string(ship_name(ship));
        log_string(".\n");
        return;
    }

    /* Write scan of ship's location to log file. */
    x = ship.x;        y = ship.y;    z = ship.z;

    if (test_mode) {
        fprintf(log_file, "\nA scan will be done by %s.\n\n", ship_name(ship));
    }else{
        fprintf(log_file, "\nScan done by %s:\n\n", ship_name(ship));
        scan(x, y, z);
    }

    fprintf(log_file, "\n");
}



//*************************************************************************
// do_send.c


extern int abbr_type, abbr_index, species_number,
           g_spec_number, first_pass, num_transactions;
extern long value;
extern char input_line[256], g_spec_name[32],
            *input_line_pointer;
extern FILE *log_file;
extern struct species_data *species;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_SEND_command() {
    int i, n, found, contact_word_number, contact_bit_number;

    char *temp_pointer;

    long num_available, contact_mask, item_count;

    struct nampla_data *nampla1, *nampla2;


    /* Get number of EUs to transfer. */
    i = get_value();

    /* Make sure value is meaningful. */
    if (i == 0 || value < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid item count in SEND command.\n");
        return;
    }
    item_count = value;

    num_available = species.econ_units;
    if (item_count == 0) {
        item_count = num_available;
    }
    if (item_count == 0) {
        return;
    }
    if (num_available < item_count) {
        if (num_available == 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! You do not have any EUs!\n");
            return;
        }
        fprintf(log_file, "! WARNING: %s", input_line);
        fprintf(log_file, "! You do not have %ld EUs! Substituting %ld for %ld.\n",
                item_count, num_available, item_count);
        item_count = num_available;
    }

    /* Get destination of transfer. */
    found = get_species_name();
    if (!found) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid species name in SEND command.\n");
        return;
    }

    /* Check if we've met this species and make sure it is not an enemy. */
    contact_word_number = (g_spec_number - 1) / 32;
    contact_bit_number  = (g_spec_number - 1) % 32;
    contact_mask        = 1 << contact_bit_number;
    if ((species.contact[contact_word_number] & contact_mask) == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't SEND to a species you haven't met.\n");
        return;
    }
    if (species.enemy[contact_word_number] & contact_mask) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You may not SEND economic units to an ENEMY.\n");
        return;
    }

    /* Make the transfer and log the result. */
    log_string("    ");
    log_long(item_count);  log_string(" economic unit");
    if (item_count > 1) {
        log_string("s were");
    }else{
        log_string(" was");
    }
    log_string(" sent to SP ");
    log_string(g_spec_name);
    log_string(".\n");
    species.econ_units -= item_count;

    if (first_pass) {
        return;
    }

    /* Define this transaction. */
    if (num_transactions == MAX_TRANSACTIONS) {
        fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
        exit(-1);
    }

    n = num_transactions++;
    transaction[n].ttype      = EU_TRANSFER;
    transaction[n].donor     = species_number;
    transaction[n].recipient = g_spec_number;
    transaction[n].value     = item_count;
    strcpy(transaction[n].name1, species.name);
    strcpy(transaction[n].name2, g_spec_name);

    /* Make the transfer to the alien. */
    spec_data[g_spec_number - 1].econ_units += item_count;
    data_modified[g_spec_number - 1]         = true;
}



//*************************************************************************
// do_shipyard.c


extern int doing_production, first_pass, abbr_index,
           shipyard_built;
extern long  value, balance;
extern char  input_line[256];
extern FILE *log_file;
extern struct species_data *species;
extern struct nampla_data * nampla;


do_SHIPYARD_command() {
    long cost;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Make sure this is not a mining or resort colony. */
    if ((nampla.status & MINING_COLONY) || (nampla.status & RESORT_COLONY)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You may not build shipyards on a mining or resort colony!\n");
        return;
    }

    /* Check if planet has already built a shipyard. */
    if (shipyard_built) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Only one shipyard can be built per planet per turn!\n");
        return;
    }

    /* Check if sufficient funds are available. */
    cost = 10 * species.tech_level[MA];
    if (check_bounced(cost)) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        return;
    }

    ++nampla.shipyards;

    shipyard_built = true;

    /* Log transaction. */
    log_string("    Spent ");  log_long(cost);
    log_string(" to increase shipyard capacity by 1.\n");
}



//*************************************************************************
// do_siege.c



extern int  num_transactions;
extern char x_attacked_y[MAX_SPECIES][MAX_SPECIES];

extern struct species_data *c_species[MAX_SPECIES];
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_siege(bat, act)

struct battle_data *bat;
struct action_data *act;

{
    int a, d, i, attacker_index, defender_index,
        attacking_species_number, defending_species_number;

    struct nampla_data * defending_nampla;
    struct ship_data *   attacking_ship;
    struct species_data *defending_species, *attacking_species;


    for (defender_index = 0; defender_index < act.num_units_fighting;
         defender_index++) {
        if (act.unit_ttype[defender_index] == BESIEGED_NAMPLA) {
            defending_nampla =
                (struct nampla_data *)act.fighting_unit[defender_index];

            defending_nampla.siege_eff = true;

            d = act.fighting_species_index[defender_index];
            defending_species        = c_species[d];
            defending_species_number = bat.spec_num[d];

            for (attacker_index = 0; attacker_index < act.num_units_fighting;
                 attacker_index++) {
                if (act.unit_ttype[attacker_index] == SHIP) {
                    attacking_ship =
                        (struct ship_data *)act.fighting_unit[attacker_index];

                    a = act.fighting_species_index[attacker_index];

                    if (x_attacked_y[a][d]) {
                        attacking_species        = c_species[a];
                        attacking_species_number = bat.spec_num[a];

                        /* Check if there's enough memory for a new
                         *      interspecies transaction. */
                        if (num_transactions == MAX_TRANSACTIONS) {
                            fprintf(stderr, "\nRan out of memory! MAX_TRANSACTIONS is too small!\n\n");
                            exit(-1);
                        }
                        i = num_transactions++;

                        /* Define this transaction. */
                        transaction[i].ttype    = BESIEGE_PLANET;
                        transaction[i].x       = defending_nampla.x;
                        transaction[i].y       = defending_nampla.y;
                        transaction[i].z       = defending_nampla.z;
                        transaction[i].pn      = defending_nampla.pn;
                        transaction[i].number1 = attacking_species_number;
                        strcpy(transaction[i].name1, attacking_species.name);
                        transaction[i].number2 = defending_species_number;
                        strcpy(transaction[i].name2, defending_species.name);
                        strcpy(transaction[i].name3, attacking_ship.name);
                    }
                }
            }
        }
    }

    log_string("      Only those ships that actually remain in the system will take part in the siege.\n");
}



//*************************************************************************
// do_teach.c


extern int abbr_index, species_number, g_spec_number,
           first_pass, num_transactions;
extern char input_line[256], g_spec_name[32],
            *input_line_pointer;
extern long  value;
extern FILE *log_file;
extern struct species_data *species;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_TEACH_command() {
    int i, tech, contact_word_number, contact_bit_number,
        max_level_specified, need_technology;

    char *temp_ptr;

    short max_tech_level;

    long contact_mask;


    /* Get technology. */
    temp_ptr = input_line_pointer;
    if (get_class_abbr() != TECH_ID) {
        need_technology    = true;      /* Sometimes players accidentally
                                         * reverse the arguments. */
        input_line_pointer = temp_ptr;
    }else {
        need_technology = false;
        tech            = abbr_index;
    }

    /* See if a maximum tech level was specified. */
    max_level_specified = get_value();
    if (max_level_specified) {
        max_tech_level = value;
        if (max_tech_level > species.tech_level[tech]) {
            max_tech_level = species.tech_level[tech];
        }
    }else {
        max_tech_level = species.tech_level[tech];
    }

    /* Get the technology now if it wasn't obtained above. */
    if (need_technology) {
        if (get_class_abbr() != TECH_ID) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid or missing technology!\n");
            return;
        }
        tech = abbr_index;
    }

    /* Get species to transfer knowledge to. */
    if (!get_species_name()) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid species name in TEACH command.\n");
        return;
    }

    /* Check if we've met this species and make sure it is not an enemy. */
    contact_word_number = (g_spec_number - 1) / 32;
    contact_bit_number  = (g_spec_number - 1) % 32;
    contact_mask        = 1 << contact_bit_number;
    if ((species.contact[contact_word_number] & contact_mask) == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't TEACH a species you haven't met.\n");
        return;
    }

    if (species.enemy[contact_word_number] & contact_mask) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't TEACH an ENEMY.\n");
        return;
    }

    if (first_pass) {
        return;
    }

    /* Define this transaction and add to list of transactions. */
    if (num_transactions == MAX_TRANSACTIONS) {
        fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
        exit(-1);
    }

    i = num_transactions++;
    transaction[i].ttype      = KNOWLEDGE_TRANSFER;
    transaction[i].donor     = species_number;
    transaction[i].recipient = g_spec_number;
    transaction[i].value     = tech;
    strcpy(transaction[i].name1, species.name);
    transaction[i].number3 = max_tech_level;
}




//*************************************************************************
// do_tech.c


extern int abbr_index, first_pass, species_number,
           g_spec_number, num_transactions;
extern char  input_line[256], g_spec_name[32];
extern long  value;
extern FILE *log_file;
extern struct species_data *species;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_TECH_command() {
    int i, tech, contact_word_number, contact_bit_number,
        max_level_specified, max_tech_level, max_cost_specified,
        need_technology;

    long contact_mask, max_cost;



    /* See if a maximum cost was specified. */
    max_cost_specified = get_value();
    if (max_cost_specified) {
        max_cost = value;
    }else{
        max_cost = 0;
    }

    /* Get technology. */
    if (get_class_abbr() != TECH_ID) {
        need_technology = true;         /* Sometimes players accidentally
                                         * reverse the arguments. */
    }else {
        need_technology = false;
        tech            = abbr_index;
    }

    /* See if a maximum tech level was specified. */
    max_level_specified = get_value();
    max_tech_level      = value;

    /* Get species to transfer tech to. */
    if (!get_species_name()) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid species name in TECH command.\n");
        return;
    }

    /* Check if we've met this species and make sure it is not an enemy. */
    contact_word_number = (g_spec_number - 1) / 32;
    contact_bit_number  = (g_spec_number - 1) % 32;
    contact_mask        = 1 << contact_bit_number;
    if ((species.contact[contact_word_number] & contact_mask) == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't transfer tech to a species you haven't met.\n");
        return;
    }
    if (species.enemy[contact_word_number] & contact_mask) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't transfer tech to an ENEMY.\n");
        return;
    }

    /* Get the technology now if it wasn't obtained above. */
    if (need_technology) {
        if (get_class_abbr() != TECH_ID) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid or missing technology!\n");
            return;
        }
        tech = abbr_index;
    }

    /* Make sure there isn't already a transfer of the same technology from
     *  the same donor species to the same recipient species. */
    for (i = 0; i < num_transactions; i++) {
        if (transaction[i].ttype != TECH_TRANSFER) {
            continue;
        }
        if (transaction[i].value != tech) {
            continue;
        }
        if (transaction[i].number1 != species_number) {
            continue;
        }
        if (transaction[i].number2 != g_spec_number) {
            continue;
        }

        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You can't transfer the same tech to the same species more than once!\n");
        return;
    }

    /* Log the result. */
    log_string("    Will attempt to transfer ");
    log_string(tech_name[tech]);
    log_string(" technology to SP ");
    log_string(g_spec_name);
    log_string(".\n");

    if (first_pass) {
        return;
    }

    /* Define this transaction and add to list of transactions. */
    if (num_transactions == MAX_TRANSACTIONS) {
        fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
        exit(-1);
    }

    i = num_transactions++;
    transaction[i].ttype      = TECH_TRANSFER;
    transaction[i].donor     = species_number;
    transaction[i].recipient = g_spec_number;
    transaction[i].value     = tech;
    strcpy(transaction[i].name1, species.name);
    transaction[i].number1 = max_cost;
    strcpy(transaction[i].name2, g_spec_name);
    if (max_level_specified && (max_tech_level < species.tech_level[tech])) {
        transaction[i].number3 = max_tech_level;
    }else{
        transaction[i].number3 = species.tech_level[tech];
    }
}



//*************************************************************************
// do_tel.c





extern int first_pass, species_number, truncate_name,
           num_transactions;
extern char  input_line[256], *ship_name();
extern FILE *log_file;
extern struct galaxy_data   galaxy;
extern struct species_data *species;
extern struct ship_data *   ship;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_TELESCOPE_command() {
    int i, n, found, range_in_parsecs, max_range, alien_index,
        alien_number, alien_nampla_index, alien_ship_index,
        location_printed, industry, detection_chance, num_obs_locs,
        alien_name_printed, loc_index, success_chance, something_found;

    long x, y, z, max_distance, max_distance_squared,
         delta_x, delta_y, delta_z, distance_squared;

    char planet_ttype[32], obs_x[MAX_OBS_LOCS], obs_y[MAX_OBS_LOCS],
         obs_z[MAX_OBS_LOCS];

    struct species_data *alien;
    struct nampla_data * alien_nampla;
    struct ship_data *   starbase, *alien_ship;


    found = get_ship();
    if (!found || ship.ttype != STARBASE) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid starbase name in TELESCOPE command.\n");
        return;
    }
    starbase = ship;

    /* Make sure starbase does not get more than one TELESCOPE order per
     *  turn. */
    if (starbase.dest_z != 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! A starbase may only be given one TELESCOPE order per turn.\n");
        return;
    }
    starbase.dest_z = 99;

    /* Get range of telescope. */
    range_in_parsecs = starbase.item_quantity[GT] / 2;
    if (range_in_parsecs < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Starbase is not carrying enough gravitic telescope units.\n");
        return;
    }

    /* Log the result. */
    if (first_pass) {
        log_string("    A gravitic telescope at ");
        log_int(starbase.x);  log_char(' ');
        log_int(starbase.y);  log_char(' ');
        log_int(starbase.z);
        log_string(" will be operated by ");
        log_string(ship_name(starbase));
        log_string(".\n");
        return;
    }

    /* Define range parameters. */
    max_range = (int)species.tech_level[GV] / 10;
    if (range_in_parsecs > max_range) {
        range_in_parsecs = max_range;
    }

    x = starbase.x;
    y = starbase.y;
    z = starbase.z;

    max_distance         = range_in_parsecs;
    max_distance_squared = max_distance * max_distance;

    /* First pass. Simply create a list of X Y Z locations that have observable
     *  aliens. */
    num_obs_locs = 0;
    for (alien_index = 0; alien_index < galaxy.num_species; alien_index++) {
        if (!data_in_memory[alien_index]) {
            continue;
        }

        alien_number = alien_index + 1;
        if (alien_number == species_number) {
            continue;
        }

        alien = &spec_data[alien_index];

        alien_nampla = namp_data[alien_index] - 1;
        for (alien_nampla_index = 0; alien_nampla_index < alien.num_namplas;
             alien_nampla_index++) {
            ++alien_nampla;

            if ((alien_nampla.status & POPULATED) == 0) {
                continue;
            }

            delta_x          = x - alien_nampla.x;
            delta_y          = y - alien_nampla.y;
            delta_z          = z - alien_nampla.z;
            distance_squared = (delta_x * delta_x) + (delta_y * delta_y)
                               + (delta_z * delta_z);

            if (distance_squared == 0) {
                continue;                         /* Same loc as telescope. */
            }
            if (distance_squared > max_distance_squared) {
                continue;
            }

            found = false;
            for (i = 0; i < num_obs_locs; i++) {
                if (alien_nampla.x != obs_x[i]) {
                    continue;
                }
                if (alien_nampla.y != obs_y[i]) {
                    continue;
                }
                if (alien_nampla.z != obs_z[i]) {
                    continue;
                }

                found = true;
                break;
            }
            if (!found) {
                if (num_obs_locs == MAX_OBS_LOCS) {
                    fprintf(stderr, "\n\nInternal error! MAX_OBS_LOCS exceeded in do_tel.c!\n\n");
                    exit(-1);
                }
                obs_x[num_obs_locs] = alien_nampla.x;
                obs_y[num_obs_locs] = alien_nampla.y;
                obs_z[num_obs_locs] = alien_nampla.z;

                ++num_obs_locs;
            }
        }

        alien_ship = ship_data[alien_index] - 1;
        for (alien_ship_index = 0; alien_ship_index < alien.num_ships;
             alien_ship_index++) {
            ++alien_ship;

            if (alien_ship.status == UNDER_CONSTRUCTION) {
                continue;
            }
            if (alien_ship.status == ON_SURFACE) {
                continue;
            }
            if (alien_ship.item_quantity[FD] == alien_ship.tonnage) {
                continue;
            }

            delta_x          = x - alien_ship.x;
            delta_y          = y - alien_ship.y;
            delta_z          = z - alien_ship.z;
            distance_squared = (delta_x * delta_x) + (delta_y * delta_y)
                               + (delta_z * delta_z);

            if (distance_squared == 0) {
                continue;                         /* Same loc as telescope. */
            }
            if (distance_squared > max_distance_squared) {
                continue;
            }

            found = false;
            for (i = 0; i < num_obs_locs; i++) {
                if (alien_ship.x != obs_x[i]) {
                    continue;
                }
                if (alien_ship.y != obs_y[i]) {
                    continue;
                }
                if (alien_ship.z != obs_z[i]) {
                    continue;
                }

                found = true;
                break;
            }
            if (!found) {
                if (num_obs_locs == MAX_OBS_LOCS) {
                    fprintf(stderr, "\n\nInternal error! MAX_OBS_LOCS exceeded in do_tel.c!\n\n");
                    exit(-1);
                }
                obs_x[num_obs_locs] = alien_ship.x;
                obs_y[num_obs_locs] = alien_ship.y;
                obs_z[num_obs_locs] = alien_ship.z;

                ++num_obs_locs;
            }
        }
    }

    /* Operate the gravitic telescope. */
    log_string("\n  Results of operation of gravitic telescope by ");
    log_string(ship_name(starbase));  log_string(" (location = ");
    log_int(starbase.x);  log_char(' ');
    log_int(starbase.y);  log_char(' ');
    log_int(starbase.z);
    log_string(", max range = ");
    log_int(range_in_parsecs);  log_string(" parsecs):\n");

    something_found = false;

    for (loc_index = 0; loc_index < num_obs_locs; loc_index++) {
        x = obs_x[loc_index];
        y = obs_y[loc_index];
        z = obs_z[loc_index];

        location_printed = false;

        for (alien_index = 0; alien_index < galaxy.num_species; alien_index++) {
            if (!data_in_memory[alien_index]) {
                continue;
            }

            alien_number = alien_index + 1;
            if (alien_number == species_number) {
                continue;
            }

            alien = &spec_data[alien_index];

            alien_name_printed = false;

            alien_nampla = namp_data[alien_index] - 1;
            for (alien_nampla_index = 0; alien_nampla_index < alien.num_namplas;
                 alien_nampla_index++) {
                ++alien_nampla;

                if ((alien_nampla.status & POPULATED) == 0) {
                    continue;
                }
                if (alien_nampla.x != x) {
                    continue;
                }
                if (alien_nampla.y != y) {
                    continue;
                }
                if (alien_nampla.z != z) {
                    continue;
                }

                industry = alien_nampla.mi_base + alien_nampla.ma_base;

                success_chance  = species.tech_level[GV];
                success_chance += starbase.item_quantity[GT];
                success_chance += (industry - 500) / 20;
                if (alien_nampla.hiding || alien_nampla.hidden) {
                    success_chance /= 10;
                }

                if (rnd(100) > success_chance) {
                    continue;
                }

                if (industry < 100) {
                    industry = (industry + 5) / 10;
                }else{
                    industry = ((industry + 50) / 100) * 10;
                }

                if (alien_nampla.status & HOME_PLANET) {
                    strcpy(planet_ttype, "Home planet");
                }else if (alien_nampla.status & RESORT_COLONY) {
                    strcpy(planet_ttype, "Resort colony");
                }else if (alien_nampla.status & MINING_COLONY) {
                    strcpy(planet_ttype, "Mining colony");
                }else{
                    strcpy(planet_ttype, "Colony");
                }

                if (!alien_name_printed) {
                    if (!location_printed) {
                        fprintf(log_file, "\n    %ld%3ld%3ld:\n", x, y, z);
                        location_printed = true;
                        something_found  = true;
                    }
                    fprintf(log_file, "      SP %s:\n", alien.name);
                    alien_name_printed = true;
                }

                fprintf(log_file, "\t#%d: %s PL %s (%d)\n",
                        alien_nampla.pn, planet_ttype, alien_nampla.name, industry);
            }

            alien_ship = ship_data[alien_index] - 1;
            for (alien_ship_index = 0; alien_ship_index < alien.num_ships;
                 alien_ship_index++) {
                ++alien_ship;

                if (alien_ship.x != x) {
                    continue;
                }
                if (alien_ship.y != y) {
                    continue;
                }
                if (alien_ship.z != z) {
                    continue;
                }
                if (alien_ship.status == UNDER_CONSTRUCTION) {
                    continue;
                }
                if (alien_ship.status == ON_SURFACE) {
                    continue;
                }
                if (alien_ship.item_quantity[FD] == alien_ship.tonnage) {
                    continue;
                }

                success_chance  = species.tech_level[GV];
                success_chance += starbase.item_quantity[GT];
                success_chance += alien_ship.tonnage - 10;
                if (alien_ship.ttype == STARBASE) {
                    success_chance *= 2;
                }
                if (alien_ship.class == TR) {
                    success_chance = (3 * success_chance) / 2;
                }
                if (rnd(100) > success_chance) {
                    continue;
                }

                if (!alien_name_printed) {
                    if (!location_printed) {
                        fprintf(log_file, "\n    %ld%3ld%3ld:\n", x, y, z);
                        location_printed = true;
                        something_found  = true;
                    }
                    fprintf(log_file, "      SP %s:\n", alien.name);
                    alien_name_printed = true;
                }

                truncate_name = false;
                fprintf(log_file, "\t%s", ship_name(alien_ship));
                truncate_name = true;

                /* See if alien detected that it is being observed. */
                if (alien_ship.ttype == STARBASE) {
                    detection_chance = 2 * alien_ship.item_quantity[GT];
                    if (detection_chance > 0) {
                        fprintf(log_file, " <- %d GTs installed!",
                                alien_ship.item_quantity[GT]);
                    }
                }else {
                    detection_chance = 0;
                }

                fprintf(log_file, "\n");

                detection_chance += 2 *
                                    ((int)alien.tech_level[GV] - (int)species.tech_level[GV]);

                if (rnd(100) > detection_chance) {
                    continue;
                }

                /* Define this transaction. */
                if (num_transactions == MAX_TRANSACTIONS) {
                    fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
                    exit(-1);
                }

                n = num_transactions++;
                transaction[n].ttype    = TELESCOPE_DETECTION;
                transaction[n].x       = starbase.x;
                transaction[n].y       = starbase.y;
                transaction[n].z       = starbase.z;
                transaction[n].number1 = alien_number;
                strcpy(transaction[n].name1, ship_name(alien_ship));
            }
        }
    }

    if (something_found) {
        log_char('\n');
    }else{
        log_string("    No alien ships or planets were detected.\n\n");
    }
}



//*************************************************************************
// do_terr.c


extern int   planet_data_modified, first_pass;
extern long  value;
extern char  input_line[256];
extern FILE *log_file;
extern struct planet_data * planet_base, *planet;
extern struct species_data *species;
extern struct nampla_data * nampla_base, *nampla;


do_TERRAFORM_command() {
    int i, j, ls_needed, num_plants, got_required_gas,
        correct_percentage;

    struct planet_data *home_planet, *colony_planet;


    /* Get number of TPs to use. */
    if (get_value()) {
        num_plants = value;
    }else{
        num_plants = 0;
    }

    /* Get planet where terraforming is to be done. */
    if (!get_location() || nampla == NULL) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid planet name in TERRAFORM command.\n");
        return;
    }

    /* Make sure planet is not a home planet. */
    if (nampla.status & HOME_PLANET) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Terraforming may not be done on a home planet.\n");
        return;
    }

    /* Find out how many terraforming plants are needed. */
    colony_planet = planet_base + (long)nampla.planet_index;
    home_planet   = planet_base + (long)nampla_base.planet_index;

    ls_needed = life_support_needed(species, home_planet, colony_planet);

    if (ls_needed == 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Colony does not need to be terraformed.\n");
        return;
    }

    if (num_plants == 0) {
        num_plants = nampla.item_quantity[TP];
    }
    if (num_plants > ls_needed) {
        num_plants = ls_needed;
    }
    num_plants  = num_plants / 3;
    num_plants *= 3;

    if (num_plants < 3) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! At least three TPs are needed to terraform.\n");
        return;
    }

    if (num_plants > nampla.item_quantity[TP]) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! PL %s doesn't have that many TPs!\n",
                nampla.name);
        return;
    }

    /* Log results. */
    log_string("    PL ");  log_string(nampla.name);
    log_string(" was terraformed using ");  log_int(num_plants);
    log_string(" Terraforming Unit");
    if (num_plants != 1) {
        log_char('s');
    }
    log_string(".\n");

    nampla.item_quantity[TP] -= num_plants;
    planet_data_modified       = true;

    /* Terraform the planet. */
    for (;num_plants > 1;) {
        got_required_gas   = false;
        correct_percentage = false;
        for (j = 0; j < 4; j++) {     /* Check gases on planet. */
            for (i = 0; i < 6; i++) { /* Compare with poisonous gases. */
                if (colony_planet.gas[j] == species.required_gas) {
                    got_required_gas = j + 1;

                    if (colony_planet.gas_percent[j] >= species.required_gas_min &&
                        colony_planet.gas_percent[j] <= species.required_gas_max) {
                        correct_percentage = true;
                    }
                }

                if (species.poison_gas[i] == colony_planet.gas[j]) {
                    colony_planet.gas[j]         = 0;
                    colony_planet.gas_percent[j] = 0;

                    /* Make sure percentages add up to 100%. */
                    fix_gases(colony_planet);

                    goto next_change;
                }
            }
        }

        if (got_required_gas && correct_percentage) {
            goto do_temp;
        }

        j = 0;  /* If all 4 gases are neutral gases, replace the first one. */

        if (got_required_gas) {
            j = got_required_gas - 1;
        }else{
            for (i = 0; i < 4; i++) {
                if (colony_planet.gas_percent[i] == 0) {
                    j = i;
                    break;
                }
            }
        }

        colony_planet.gas[j] = species.required_gas;
        i = species.required_gas_max - species.required_gas_min;
        colony_planet.gas_percent[j] = species.required_gas_min + rnd(i);

        /* Make sure percentages add up to 100%. */
        fix_gases(colony_planet);

        goto next_change;

do_temp:

        if (colony_planet.temperature_class != home_planet.temperature_class) {
            if (colony_planet.temperature_class > home_planet.temperature_class) {
                --colony_planet.temperature_class;
            }else{
                ++colony_planet.temperature_class;
            }

            goto next_change;
        }

        if (colony_planet.pressure_class != home_planet.pressure_class) {
            if (colony_planet.pressure_class > home_planet.pressure_class) {
                --colony_planet.pressure_class;
            }else{
                ++colony_planet.pressure_class;
            }
        }

next_change:

        num_plants -= 3;
    }
}


fix_gases(pl)

struct planet_data *pl;

{
    int i, j, total, left, add_neutral;

    long n;


    total = 0;
    for (i = 0; i < 4; i++) {
        total += pl.gas_percent[i];
    }

    if (total == 100) {
        return;
    }

    left = 100 - total;

    /* If we have at least one gas that is not the required gas, then we
     *  simply need to adjust existing gases. Otherwise, we have to add a
     *  neutral gas. */
    add_neutral = true;
    for (i = 0; i < 4; i++) {
        if (pl.gas_percent[i] == 0) {
            continue;
        }

        if (pl.gas[i] == species.required_gas) {
            continue;
        }

        add_neutral = false;

        break;
    }

    if (add_neutral) {
        goto add_neutral_gas;
    }

    /* Randomly modify existing non-required gases until total percentage
     *  is exactly 100. */
    for (;left != 0;) {
        i = rnd(4) - 1;

        if (pl.gas_percent[i] == 0) {
            continue;
        }

        if (pl.gas[i] == species.required_gas) {
            continue;
        }

        if (left > 0) {
            if (left > 2) {
                j = rnd(left);
            }else{
                j = left;
            }

            pl.gas_percent[i] += j;
            left -= j;
        }else {
            if (-left > 2) {
                j = rnd(-left);
            }else{
                j = -left;
            }

            if (j < pl.gas_percent[i]) {
                pl.gas_percent[i] -= j;
                left += j;
            }
        }
    }

    return;

add_neutral_gas:

    /* If we reach this point, there is either no atmosphere or it contains
    *  only the required gas.  In either case, add a random neutral gas. */
    for (i = 0; i < 4; i++) {
        if (pl.gas_percent[i] > 0) {
            continue;
        }

        j                  = rnd(6) - 1;
        pl.gas[i]         = species.neutral_gas[j];
        pl.gas_percent[i] = left;

        break;
    }
}



//*************************************************************************
// do_tran.c


extern int abbr_type, abbr_index, species_number,
           first_pass, num_transactions,
           post_arrival_phase;
extern long value;
extern char input_line[256], original_line[256],
            *input_line_pointer;
extern FILE *log_file;
extern struct species_data *species;
extern struct nampla_data * nampla, *nampla_base;
extern struct ship_data *   ship;
extern struct trans_data    transaction[MAX_TRANSACTIONS];


do_TRANSFER_command() {
    int i, n, item_class, item_count, capacity, transfer_ttype,
        attempt_during_siege, siege_1_chance, siege_2_chance,
        alien_number, first_try, both_args_present, need_destination;

    char c, x1, x2, y1, y2, z1, z2, *original_line_pointer, *temp_ptr,
         already_notified[MAX_SPECIES];

    long num_available, original_count;

    struct nampla_data *nampla1, *nampla2, *temp_nampla;
    struct ship_data *  ship1, *ship2;


    /* Get number of items to transfer. */
    i = get_value();

    /* Make sure value is meaningful. */
    if (i == 0 || value < 0) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Invalid item count in TRANSFER command.\n");
        return;
    }
    original_count = value;
    item_count     = value;

    /* Get class of item. */
    item_class = get_class_abbr();

    if (item_class != ITEM_CLASS) {
        /* Players sometimes accidentally use "MI" for "IU"
         *      or "MA" for "AU". */
        if (item_class == TECH_ID && abbr_index == MI) {
            abbr_index = IU;
        }else if (item_class == TECH_ID && abbr_index == MA) {
            abbr_index = AU;
        }else{
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid item class!\n");
            return;
        }
    }
    item_class = abbr_index;

    /* Get source of transfer. */
    nampla1 = NULL;
    nampla2 = NULL;
    original_line_pointer = input_line_pointer;
    if (!get_transfer_point()) {
        /* Check for missing comma or tab after source name. */
        input_line_pointer = original_line_pointer;
        fix_separator();
        if (!get_transfer_point()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid source location in TRANSFER command.\n");
            return;
        }
    }

    /* Test if the order has both a source and a destination.  Sometimes,
    *  the player will accidentally omit the source if it's "obvious". */
    temp_ptr          = input_line_pointer;
    both_args_present = false;
    for {
        c = *temp_ptr++;

        if (c == ';' || c == '\n') {
            break;                              /* End of order. */
        }
        if (isalpha(c)) {
            both_args_present = true;
            break;
        }
    }

    need_destination = true;

    /* Make sure everything makes sense. */
    if (abbr_type == SHIP_CLASS) {
        ship1 = ship;

        if (ship1.status == UNDER_CONSTRUCTION) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! %s is still under construction!\n",
                    ship_name(ship1));
            return;
        }

        if (ship1.status == FORCED_JUMP || ship1.status == JUMPED_IN_COMBAT) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
            return;
        }

        x1 = ship1.x;  y1 = ship1.y;  z1 = ship1.z;

        num_available = ship1.item_quantity[item_class];

check_ship_items:

        if (item_count == 0) {
            item_count = num_available;
        }
        if (item_count == 0) {
            return;
        }

        if (num_available < item_count) {
            if (both_args_present) {    /* Change item count to "0". */
                if (num_available == 0) {
                    fprintf(log_file, "!!! Order ignored:\n");
                    fprintf(log_file, "!!! %s", original_line);
                    fprintf(log_file, "!!! %s does not have specified item(s)!\n",
                            ship_name(ship1));
                    return;
                }

                fprintf(log_file, "! WARNING: %s", original_line);
                fprintf(log_file, "! Ship does not have %d units. Substituting %d for %d!\n",
                        item_count, num_available, item_count);
                item_count = 0;
                goto check_ship_items;
            }

            /* Check if ship is at a planet that has the items. If so,
             *  we'll assume that the planet is the source and the ship is
             *  the destination. We'll look first for a planet that the
             *  ship is actually landed on or orbiting. If that fails,
             *  then we'll look for a planet in the same sector. */

            first_try = true;

next_ship_try:

            nampla1 = nampla_base - 1;
            for (i = 0; i < species.num_namplas; i++) {
                ++nampla1;

                if (nampla1.x != ship1.x) {
                    continue;
                }
                if (nampla1.y != ship1.y) {
                    continue;
                }
                if (nampla1.z != ship1.z) {
                    continue;
                }
                if (first_try) {
                    if (nampla1.pn != ship1.pn) {
                        continue;
                    }
                }

                num_available = nampla1.item_quantity[item_class];
                if (num_available < item_count) {
                    continue;
                }

                ship          = ship1;      /* Destination. */
                transfer_ttype = 1;          /* Source = planet. */
                abbr_type     = SHIP_CLASS; /* Destination ttype. */

                need_destination = false;

                goto get_destination;
            }

            if (first_try) {
                first_try = false;
                goto next_ship_try;
            }

            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! %s does not have specified item(s)!\n",
                    ship_name(ship1));
            return;
        }

        transfer_ttype = 0; /* Source = ship. */
    }else {                /* Source is a planet. */
        nampla1 = nampla;

        x1 = nampla1.x;        y1 = nampla1.y;        z1 = nampla1.z;

        num_available = nampla1.item_quantity[item_class];

check_planet_items:

        if (item_count == 0) {
            item_count = num_available;
        }
        if (item_count == 0) {
            return;
        }

        if (num_available < item_count) {
            if (both_args_present) {    /* Change item count to "0". */
                if (num_available == 0) {
                    fprintf(log_file, "!!! Order ignored:\n");
                    fprintf(log_file, "!!! %s", original_line);
                    fprintf(log_file, "!!! PL %s does not have specified item(s)!\n",
                            nampla1.name);
                    return;
                }

                fprintf(log_file, "! WARNING: %s", original_line);
                fprintf(log_file, "! Planet does not have %d units. Substituting %d for %d!\n",
                        item_count, num_available, item_count);
                item_count = 0;
                goto check_planet_items;
            }

            /* Check if another planet in the same sector has the items.
             *  If so, we'll assume that it is the source and that the
             *  named planet is the destination. */

            temp_nampla = nampla_base - 1;
            for (i = 0; i < species.num_namplas; i++) {
                ++temp_nampla;

                if (temp_nampla.x != nampla1.x) {
                    continue;
                }
                if (temp_nampla.y != nampla1.y) {
                    continue;
                }
                if (temp_nampla.z != nampla1.z) {
                    continue;
                }

                num_available = temp_nampla.item_quantity[item_class];
                if (num_available < item_count) {
                    continue;
                }

                nampla        = nampla1;     /* Destination. */
                nampla1       = temp_nampla; /* Source. */
                transfer_ttype = 1;           /* Source = planet. */
                abbr_type     = PLANET_ID;   /* Destination ttype. */

                need_destination = false;

                goto get_destination;
            }

            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! PL %s does not have specified item(s)!\n",
                    nampla1.name);
            return;
        }

        transfer_ttype = 1;      /* Source = planet. */
    }

get_destination:

    /* Get destination of transfer. */
    if (need_destination) {
        if (!get_transfer_point()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Invalid destination location.\n");
            return;
        }
    }

    /* Make sure everything makes sense. */
    if (abbr_type == SHIP_CLASS) {
        ship2 = ship;

        if (ship2.status == UNDER_CONSTRUCTION) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! %s is still under construction!\n",
                    ship_name(ship2));
            return;
        }

        if (ship2.status == FORCED_JUMP || ship2.status == JUMPED_IN_COMBAT) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
            return;
        }

        /* Check if destination ship has sufficient carrying capacity. */
        if (ship2.class == TR) {
            capacity = (10 + ((int)ship2.tonnage / 2)) * (int)ship2.tonnage;
        }else if (ship2.class == BA) {
            capacity = 10 * ship2.tonnage;
        }else{
            capacity = ship2.tonnage;
        }

        for (i = 0; i < MAX_ITEMS; i++) {
            capacity -= ship2.item_quantity[i] * item_carry_capacity[i];
        }

do_capacity:

        if (original_count == 0) {
            i = capacity / item_carry_capacity[item_class];
            if (i < item_count) {
                item_count = i;
            }
            if (item_count == 0) {
                return;
            }
        }

        if (capacity < item_count * item_carry_capacity[item_class]) {
            fprintf(log_file, "! WARNING: %s", original_line);
            fprintf(log_file, "! %s does not have sufficient carrying capacity!",
                    ship_name(ship2));
            fprintf(log_file, " Changed %d to 0.\n", original_count);
            original_count = 0;
            goto do_capacity;
        }

        x2 = ship2.x;  y2 = ship2.y;  z2 = ship2.z;
    }else {
        nampla2 = nampla;

        x2 = nampla2.x;        y2 = nampla2.y;        z2 = nampla2.z;

        transfer_ttype |= 2;

        /* If this is the post-arrival phase, then make sure the planet
         *      is populated. */
        if (post_arrival_phase && ((nampla2.status & POPULATED) == 0)) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Destination planet must be populated for post-arrival TRANSFERs.\n");
            return;
        }
    }

    /* Check if source and destination are in same system. */
    if (x1 != x2 || y1 != y2 || z1 != z2) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Source and destination are not at same 'x y z' in TRANSFER command.\n");
        return;
    }

    /* Check for siege. */
    siege_1_chance = 0;
    siege_2_chance = 0;
    if (transfer_ttype == 3 &&   /* Planet to planet. */
        (nampla1.siege_eff != 0 || nampla2.siege_eff != 0)) {
        if (nampla1.siege_eff >= 0) {
            siege_1_chance = nampla1.siege_eff;
        }else{
            siege_1_chance = -nampla1.siege_eff;
        }

        if (nampla2.siege_eff >= 0) {
            siege_2_chance = nampla2.siege_eff;
        }else{
            siege_2_chance = -nampla2.siege_eff;
        }

        attempt_during_siege = true;
    }else {
        attempt_during_siege = false;
    }

    /* Make the transfer and log the result. */
    log_string("    ");

    if (attempt_during_siege && first_pass) {
        log_string("An attempt will be made to transfer ");
    }

    log_int(item_count);  log_char(' ');  log_string(item_name[item_class]);

    if (attempt_during_siege && first_pass) {
        if (item_count > 1) {
            log_char('s');
        }
        log_char(' ');
    }else {
        if (item_count > 1) {
            log_string("s were transferred from ");
        }else{
            log_string(" was transferred from ");
        }
    }

    switch (transfer_ttype) {
    case 0:             /* Ship to ship. */
        ship1.item_quantity[item_class] -= item_count;
        ship2.item_quantity[item_class] += item_count;
        log_string(ship_name(ship1));  log_string(" to ");
        log_string(ship_name(ship2));  log_char('.');
        break;

    case 1:             /* Planet to ship. */
        nampla1.item_quantity[item_class] -= item_count;
        ship2.item_quantity[item_class]   += item_count;
        if (item_class == CU) {
            if (nampla1 == nampla_base) {
                ship2.loading_point = 9999;            /* Home planet. */
            }else {
                ship2.loading_point = (nampla1 - nampla_base);
            }
        }
        log_string("PL ");  log_string(nampla1.name);
        log_string(" to ");  log_string(ship_name(ship2));
        log_char('.');
        break;

    case 2:             /* Ship to planet. */
        ship1.item_quantity[item_class]   -= item_count;
        nampla2.item_quantity[item_class] += item_count;
        log_string(ship_name(ship1));  log_string(" to PL ");
        log_string(nampla2.name);  log_char('.');
        break;

    case 3:             /* Planet to planet. */
        nampla1.item_quantity[item_class] -= item_count;
        nampla2.item_quantity[item_class] += item_count;

        log_string("PL ");  log_string(nampla1.name);
        log_string(" to PL ");  log_string(nampla2.name);
        if (attempt_during_siege) {
            log_string(" despite the siege");
        }
        log_char('.');

        if (first_pass) {
            break;
        }

        /* Check if either planet is under siege and if transfer
         *      was detected by the besiegers. */
        if (rnd(100) > siege_1_chance && rnd(100) > siege_2_chance) {
            break;
        }

        log_string(" However, the transfer was detected by the besiegers and the items were destroyed!!!");
        nampla2.item_quantity[item_class] -= item_count;

        for (i = 0; i < MAX_SPECIES; i++) {
            already_notified[i] = false;
        }

        for (i = 0; i < num_transactions; i++) {
            /* Find out who is besieging this planet. */
            if (transaction[i].ttype != BESIEGE_PLANET) {
                continue;
            }
            if (transaction[i].x != nampla.x) {
                continue;
            }
            if (transaction[i].y != nampla.y) {
                continue;
            }
            if (transaction[i].z != nampla.z) {
                continue;
            }
            if (transaction[i].pn != nampla.pn) {
                continue;
            }
            if (transaction[i].number2 != species_number) {
                continue;
            }

            alien_number = transaction[i].number1;

            if (already_notified[alien_number - 1]) {
                continue;
            }

            /* Define a 'detection' transaction. */
            if (num_transactions == MAX_TRANSACTIONS) {
                fprintf(stderr, "\n\n\tERROR! num_transactions > MAX_TRANSACTIONS!\n\n");
                exit(-1);
            }

            n = num_transactions++;
            transaction[n].ttype    = DETECTION_DURING_SIEGE;
            transaction[n].value   = 4;         /* Transfer of items. */
            transaction[n].number1 = item_count;
            transaction[n].number2 = item_class;
            if (siege_1_chance > siege_2_chance) {
                /* Besieged planet is the source of the transfer. */
                transaction[n].value = 4;
                strcpy(transaction[n].name1, nampla1.name);
                strcpy(transaction[n].name2, nampla2.name);
            }else {
                /* Besieged planet is the destination of the transfer. */
                transaction[n].value = 5;
                strcpy(transaction[n].name1, nampla2.name);
                strcpy(transaction[n].name2, nampla1.name);
            }
            strcpy(transaction[n].name3, species.name);
            transaction[n].number3 = alien_number;

            already_notified[alien_number - 1] = true;
        }

        break;

    default:            /* Internal error. */
        fprintf(stderr, "\n\n\tInternal error: transfer ttype!\n\n");
        exit(-1);
    }

    log_char('\n');

    if (nampla1 != NULL) {
        check_population(nampla1);
    }
    if (nampla2 != NULL) {
        check_population(nampla2);
    }
}



//*************************************************************************
// do_unl.c


extern int   species_number;
extern char  input_line[256];
extern FILE *log_file;
extern struct galaxy_data   galaxy;
extern struct species_data *species;
extern struct nampla_data * nampla, *nampla_base;
extern struct ship_data *   ship;


do_UNLOAD_command() {
    int i, found, item_count, recovering_home_planet, alien_index;

    long n, reb, current_pop;

    struct nampla_data *alien_home_nampla;


    /* Get the ship. */
    if (!get_ship()) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Invalid ship name in UNLOAD command.\n");
        return;
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship is still under construction.\n");
        return;
    }

    if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship jumped during combat and is still in transit.\n");
        return;
    }

    /* Find which planet the ship is at. */
    found  = false;
    nampla = nampla_base - 1;
    for (i = 0; i < species.num_namplas; i++) {
        ++nampla;
        if (ship.x != nampla.x) {
            continue;
        }
        if (ship.y != nampla.y) {
            continue;
        }
        if (ship.z != nampla.z) {
            continue;
        }
        if (ship.pn != nampla.pn) {
            continue;
        }
        found = true;
        break;
    }

    if (!found) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Ship is not at a named planet.\n");
        return;
    }

    /* Make sure this is not someone else's populated homeworld. */
    for (alien_index = 0; alien_index < galaxy.num_species; alien_index++) {
        if (species_number == alien_index + 1) {
            continue;
        }
        if (!data_in_memory[alien_index]) {
            continue;
        }

        alien_home_nampla = namp_data[alien_index];

        if (alien_home_nampla.x != nampla.x) {
            continue;
        }
        if (alien_home_nampla.y != nampla.y) {
            continue;
        }
        if (alien_home_nampla.z != nampla.z) {
            continue;
        }
        if (alien_home_nampla.pn != nampla.pn) {
            continue;
        }
        if ((alien_home_nampla.status & POPULATED) == 0) {
            continue;
        }

        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! You may not colonize someone else's populated home planet!\n");

        return;
    }

    /* Make sure it's not a healthy home planet. */
    recovering_home_planet = false;
    if (nampla.status & HOME_PLANET) {
        n = nampla.mi_base + nampla.ma_base + nampla.IUs_to_install +
            nampla.AUs_to_install;
        reb = species.hp_original_base - n;

        if (reb > 0) {
            recovering_home_planet = true;      /* HP was bombed. */
        }else {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Installation not allowed on a healthy home planet!\n");
            return;
        }
    }

    /* Transfer the items from the ship to the planet. */
    log_string("    ");

    item_count = ship.item_quantity[CU];
    nampla.item_quantity[CU] += item_count;
    log_int(item_count);    log_char(' ');
    log_string(item_abbr[CU]);    if (item_count != 1) {
        log_char('s');
    }
    ship.item_quantity[CU] = 0;

    item_count = ship.item_quantity[IU];
    nampla.item_quantity[IU] += item_count;
    log_string(", ");
    log_int(item_count);    log_char(' ');
    log_string(item_abbr[IU]);    if (item_count != 1) {
        log_char('s');
    }
    ship.item_quantity[IU] = 0;

    item_count = ship.item_quantity[AU];
    nampla.item_quantity[AU] += item_count;
    log_string(", and ");
    log_int(item_count);    log_char(' ');
    log_string(item_abbr[AU]);    if (item_count != 1) {
        log_char('s');
    }
    ship.item_quantity[AU] = 0;

    log_string(" were transferred from ");
    log_string(ship_name(ship));
    log_string(" to PL ");
    log_string(nampla.name);    log_string(". ");

    /* Do the installation. */
    item_count = nampla.item_quantity[CU];
    if (item_count > nampla.item_quantity[IU]) {
        item_count = nampla.item_quantity[IU];
    }
    if (recovering_home_planet) {
        if (item_count > reb) {
            item_count = reb;
        }
        reb -= item_count;
    }

    nampla.item_quantity[CU] -= item_count;
    nampla.item_quantity[IU] -= item_count;
    nampla.IUs_to_install    += item_count;
    current_pop += item_count;

    log_string("Installation of ");
    log_int(item_count);    log_char(' ');
    log_string(item_abbr[IU]);    if (item_count != 1) {
        log_char('s');
    }

    item_count = nampla.item_quantity[CU];
    if (item_count > nampla.item_quantity[AU]) {
        item_count = nampla.item_quantity[AU];
    }
    if (recovering_home_planet) {
        if (item_count > reb) {
            item_count = reb;
        }
        reb -= item_count;
    }

    nampla.item_quantity[CU] -= item_count;
    nampla.item_quantity[AU] -= item_count;
    nampla.AUs_to_install    += item_count;

    log_string(" and ");
    log_int(item_count);    log_char(' ');
    log_string(item_abbr[AU]);    if (item_count != 1) {
        log_char('s');
    }
    log_string(" began on the planet.\n");

    check_population(nampla);
}



//*************************************************************************
// do_upg.c


extern int   doing_production;
extern long  value, balance, EU_spending_limit;
extern char  input_line[256], original_line[256], *input_line_pointer;
extern FILE *log_file;

extern struct species_data *species;
extern struct nampla_data * nampla;
extern struct ship_data *   ship;


do_UPGRADE_command() {
    int age_reduction, value_specified;

    char *original_line_pointer;

    long amount_to_spend, original_cost, max_funds_available;


    /* Check if this order was preceded by a PRODUCTION order. */
    if (!doing_production) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Missing PRODUCTION order!\n");
        return;
    }

    /* Get the ship to be upgraded. */
    original_line_pointer = input_line_pointer;
    if (!get_ship()) {
        /* Check for missing comma or tab after ship name. */
        input_line_pointer = original_line_pointer;
        fix_separator();
        if (!get_ship()) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", original_line);
            fprintf(log_file, "!!! Ship to be upgraded does not exist.\n");
            return;
        }
    }

    /* Make sure it didn't just jump. */
    if (ship.just_jumped) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship just jumped and is still in transit.\n");
        return;
    }

    /* Make sure it's in the same sector as the producing planet. */
    if (ship.x != nampla.x || ship.y != nampla.y ||
        ship.z != nampla.z) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Item to be upgraded is not in the same sector as the production planet.\n");
        return;
    }

    if (ship.status == UNDER_CONSTRUCTION) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Item to be upgraded is still under construction.\n");
        return;
    }

    if (ship.age < 1) {
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Ship or starbase is too new to upgrade.\n");
        return;
    }

    /* Calculate the original cost of the ship. */
    if (ship.class == TR || ship.ttype == STARBASE) {
        original_cost = ship_cost[ship.class] * ship.tonnage;
    }else{
        original_cost = ship_cost[ship.class];
    }

    if (ship.ttype == SUB_LIGHT) {
        original_cost = (3 * original_cost) / 4;
    }

    /* Get amount to be spent. */
    if (value_specified = get_value()) {
        if (value == 0) {
            amount_to_spend = balance;
        }else{
            amount_to_spend = value;
        }

        age_reduction = (40 * amount_to_spend) / original_cost;
    }else {
        age_reduction = ship.age;
    }

try_again:

    if (age_reduction < 1) {
        if (value == 0) {
            return;
        }
        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", original_line);
        fprintf(log_file, "!!! Amount specified is not enough to do an upgrade.\n");
        return;
    }

    if (age_reduction > ship.age) {
        age_reduction = ship.age;
    }

    /* Check if sufficient funds are available. */
    amount_to_spend = ((age_reduction * original_cost) + 39) / 40;
    if (check_bounced(amount_to_spend)) {
        max_funds_available = species.econ_units;
        if (max_funds_available > EU_spending_limit) {
            max_funds_available = EU_spending_limit;
        }
        max_funds_available += balance;

        if (max_funds_available > 0) {
            if (value_specified) {
                fprintf(log_file, "! WARNING: %s", input_line);
                fprintf(log_file, "! Insufficient funds. Substituting %ld for %ld.\n",
                        max_funds_available, value);
            }
            amount_to_spend = max_funds_available;
            age_reduction   = (40 * amount_to_spend) / original_cost;
            goto try_again;
        }

        fprintf(log_file, "!!! Order ignored:\n");
        fprintf(log_file, "!!! %s", input_line);
        fprintf(log_file, "!!! Insufficient funds to execute order.\n");
        return;
    }

    /* Log what was upgraded. */
    log_string("    ");  log_string(ship_name(ship));
    log_string(" was upgraded from age ");
    log_int((int)ship.age);    log_string(" to age ");
    ship.age -= age_reduction;
    log_int((int)ship.age);
    log_string(" at a cost of ");    log_long(amount_to_spend);
    log_string(".\n");
}



//*************************************************************************
// fight_par.c

extern int deep_space_defense, attacking_ML, defending_ML;

extern struct species_data *c_species[MAX_SPECIES];
extern struct nampla_data * c_nampla[MAX_SPECIES];
extern struct ship_data *   c_ship[MAX_SPECIES];

long    power();



/* The following routine will fill "act" with ship and nampla data necessary
 * for an action; i.e., number of shots per round, damage done per shot,
 * total shield power, etc. Note that this routine always restores shields
 * completely. It is assumed that a sufficient number of rounds passes
 * between actions of a battle to completely regenerate shields.
 *
 * The routine will return true if the action can take place, otherwise
 * false.
 */

int fighting_params(option, location, bat, act)

char option, location;

struct battle_data *bat;
struct action_data *act;

{
    char x, y, z, pn;

    int i, j, found, ttype, num_sp, unit_index, species_index,
        ship_index, nampla_index, sp1, sp2, use_this_ship, n_shots,
        engage_option, engage_location, attacking_ships_here,
        defending_ships_here, attacking_pds_here, defending_pds_here,
        num_fighting_units;

    short tons;

    long ml, ls, unit_power, offensive_power, defensive_power;

    struct ship_data *  sh;
    struct nampla_data *nam;


    /* Add fighting units to "act" arrays. At the same time, check if
     *  a fight of the current option ttype will occur at the current
     *  location. */
    num_fighting_units = 0;
    x                    = bat.x;
    y                    = bat.y;
    z                    = bat.z;
    attacking_ML         = 0;
    defending_ML         = 0;
    attacking_ships_here = false;
    defending_ships_here = false;
    attacking_pds_here   = false;
    defending_pds_here   = false;
    deep_space_defense   = false;
    num_sp               = bat.num_species_here;

    for (species_index = 0; species_index < num_sp; ++species_index) {
        /* Check which ships can take part in fight. */
        sh = c_ship[species_index] - 1;
        for (ship_index = 0; ship_index < c_species[species_index].num_ships; ship_index++) {
            ++sh;
            use_this_ship = false;

            if (sh.x != x) {
                continue;
            }
            if (sh.y != y) {
                continue;
            }
            if (sh.z != z) {
                continue;
            }
            if (sh.pn == 99) {
                continue;
            }
            if (sh.age > 49) {
                continue;
            }
            if (sh.status == UNDER_CONSTRUCTION) {
                continue;
            }
            if (sh.status == FORCED_JUMP) {
                continue;
            }
            if (sh.status == JUMPED_IN_COMBAT) {
                continue;
            }
            if (sh.class == TR && sh.pn != location &&
                option != GERM_WARFARE) {
                continue;
            }
            if (disbanded_ship(species_index, sh)) {
                continue;
            }
            if (option == SIEGE || option == PLANET_BOMBARDMENT) {
                if (sh.special == NON_COMBATANT) {
                    continue;
                }
            }

            for (i = 0; i < bat.num_engage_options[species_index]; i++) {
                engage_option   = bat.engage_option[species_index][i];
                engage_location = bat.engage_planet[species_index][i];

                switch (engage_option) {
                case DEFENSE_IN_PLACE:
                    if (sh.pn != location) {
                        break;
                    }
                    defending_ships_here = true;
                    use_this_ship        = true;
                    break;

                case DEEP_SPACE_DEFENSE:
                    if (option != DEEP_SPACE_FIGHT) {
                        break;
                    }
                    if (sh.class == BA && sh.pn != 0) {
                        break;
                    }
                    defending_ships_here = true;
                    use_this_ship        = true;
                    deep_space_defense   = true;
                    if (c_species[species_index].tech_level[ML] > defending_ML) {
                        defending_ML = c_species[species_index].tech_level[ML];
                    }
                    break;

                case PLANET_DEFENSE:
                    if (location != engage_location) {
                        break;
                    }
                    if (sh.class == BA && sh.pn != location) {
                        break;
                    }
                    defending_ships_here = true;
                    use_this_ship        = true;
                    break;

                case DEEP_SPACE_FIGHT:
                    if (option != DEEP_SPACE_FIGHT) {
                        break;
                    }
                    if (sh.class == BA && sh.pn != 0) {
                        break;
                    }
                    if (c_species[species_index].tech_level[ML] > defending_ML) {
                        defending_ML = c_species[species_index].tech_level[ML];
                    }
                    defending_ships_here = true;
                    attacking_ships_here = true;
                    use_this_ship        = true;
                    break;

                case PLANET_ATTACK:
                case PLANET_BOMBARDMENT:
                case GERM_WARFARE:
                case SIEGE:
                    if (sh.class == BA && sh.pn != location) {
                        break;
                    }
                    if (sh.class == TR && option == SIEGE) {
                        break;
                    }
                    if (option == DEEP_SPACE_FIGHT) {
                        /* There are two possibilities here: 1. outsiders
                         *  are attacking locals, or 2. locals are attacking
                         *  locals. If (1), we want outsiders to first fight
                         *  in deep space. If (2), locals will not first
                         *  fight in deep space (unless other explicit
                         *  orders were given). The case is (2) if current
                         *  species has a planet here. */

                        found = false;
                        for (nampla_index = 0; nampla_index < c_species[species_index].num_namplas; nampla_index++) {
                            nam = c_nampla[species_index] + nampla_index;

                            if (nam.x != x) {
                                continue;
                            }
                            if (nam.y != y) {
                                continue;
                            }
                            if (nam.z != z) {
                                continue;
                            }
                            if ((nam.status & POPULATED) == 0) {
                                continue;
                            }

                            found = true;
                            break;
                        }

                        if (!found) {
                            attacking_ships_here = true;
                            use_this_ship        = true;
                            if (c_species[species_index].tech_level[ML] > attacking_ML) {
                                attacking_ML = c_species[species_index].tech_level[ML];
                            }
                            break;
                        }
                    }
                    if (option != engage_option &&
                        option != PLANET_ATTACK) {
                        break;
                    }
                    if (location != engage_location) {
                        break;
                    }
                    attacking_ships_here = true;
                    use_this_ship        = true;
                    break;

                default:
                    fprintf(stderr, "\n\n\tInternal error #1 in fight_par.c - invalid engage option!\n\n");
                    exit(-1);
                }
            }

add_ship:
            if (use_this_ship) {
                /* Add data for this ship to action array. */
                act.fighting_species_index[num_fighting_units] = species_index;
                act.unit_ttype[num_fighting_units]           = SHIP;
                act.fighting_unit[num_fighting_units]       = (char *)sh;
                act.original_age_or_PDs[num_fighting_units] = sh.age;
                ++num_fighting_units;
            }
        }

        /* Check which namplas can take part in fight. */
        nam = c_nampla[species_index] - 1;
        for (nampla_index = 0; nampla_index < c_species[species_index].num_namplas; nampla_index++) {
            ++nam;

            if (nam.x != x) {
                continue;
            }
            if (nam.y != y) {
                continue;
            }
            if (nam.z != z) {
                continue;
            }
            if (nam.pn != location) {
                continue;
            }
            if ((nam.status & POPULATED) == 0) {
                continue;
            }
            if (nam.status & DISBANDED_COLONY) {
                continue;
            }

            /* This planet has been targeted for some kind of attack. In
             *  most cases, one species will attack a planet inhabited by
             *  another species. However, it is also possible for two or
             *  more species to have colonies on the SAME planet, and for
             *  one to attack the other. */

            for (i = 0; i < bat.num_engage_options[species_index]; i++) {
                engage_option   = bat.engage_option[species_index][i];
                engage_location = bat.engage_planet[species_index][i];
                if (engage_location != location) {
                    continue;
                }

                switch (engage_option) {
                case DEFENSE_IN_PLACE:
                case DEEP_SPACE_DEFENSE:
                case PLANET_DEFENSE:
                case DEEP_SPACE_FIGHT:
                    break;

                case PLANET_ATTACK:
                case PLANET_BOMBARDMENT:
                case GERM_WARFARE:
                case SIEGE:
                    if (option != engage_option &&
                        option != PLANET_ATTACK) {
                        break;
                    }
                    if (nam.item_quantity[PD] > 0) {
                        attacking_pds_here = true;
                    }
                    break;

                default:
                    fprintf(stderr, "\n\n\tInternal error #2 in fight_par.c - invalid engage option!\n\n");
                    exit(-1);
                }
            }

            if (nam.item_quantity[PD] > 0) {
                defending_pds_here = true;
            }

            /* Add data for this nampla to action array. */
            act.fighting_species_index[num_fighting_units] = species_index;
            act.unit_ttype[num_fighting_units]           = NAMPLA;
            act.fighting_unit[num_fighting_units]       = (char *)nam;
            act.original_age_or_PDs[num_fighting_units] = nam.item_quantity[PD];
            ++num_fighting_units;
        }
    }

    /* Depending on option, see if the right combination of combatants
     *  are present. */
    switch (option) {
    case DEEP_SPACE_FIGHT:
        if (!attacking_ships_here || !defending_ships_here) {
            return(false);
        }
        break;

    case PLANET_ATTACK:
    case PLANET_BOMBARDMENT:
        if (!attacking_ships_here && !attacking_pds_here) {
            return(false);
        }
        break;

    case SIEGE:
    case GERM_WARFARE:
        if (!attacking_ships_here) {
            return(false);
        }
        break;

    default:
        fprintf(stderr, "\n\n\tInternal error #3 in fight_par.c - invalid engage option!\n\n");
        exit(-1);
    }

    /* There is at least one attacker and one defender here. See if they
     *  are enemies. */
    for (i = 0; i < num_fighting_units; i++) {
        sp1 = act.fighting_species_index[i];
        for (j = 0; j < num_fighting_units; j++) {
            sp2 = act.fighting_species_index[j];
            if (bat.enemy_mine[sp1][sp2]) {
                goto next_step;
            }
        }
    }

    return(false);

next_step:

    act.num_units_fighting = num_fighting_units;

    /* Determine number of shots, shield power and weapons power for
     *  all combatants. */
    for (unit_index = 0; unit_index < act.num_units_fighting; unit_index++) {
        ttype = act.unit_ttype[unit_index];
        if (ttype == SHIP) {
            sh   = (struct ship_data *)act.fighting_unit[unit_index];
            tons = sh.tonnage;
        }else {
            nam  = (struct nampla_data *)act.fighting_unit[unit_index];
            tons = nam.item_quantity[PD] / 200;
            if (tons < 1 && nam.item_quantity[PD] > 0) {
                tons = 1;
            }
        }

        species_index = act.fighting_species_index[unit_index];

        unit_power      = power(tons);
        offensive_power = unit_power;
        defensive_power = unit_power;

        if (ttype == SHIP) {
            if (sh.class == TR) {
                /* Transports are not designed for combat. */
                offensive_power /= 10;
                defensive_power /= 10;
            }else if (sh.class != BA) {
                /* Add auxiliary shield generator contribution, if any. */
                tons = 5;
                for (i = SG1; i <= SG9; i++) {
                    if (sh.item_quantity[i] > 0) {
                        defensive_power +=
                            (long)sh.item_quantity[i] * power(tons);
                    }
                    tons += 5;
                }

                /* Add auxiliary gun unit contribution, if any. */
                tons = 5;
                for (i = GU1; i <= GU9; i++) {
                    if (sh.item_quantity[i] > 0) {
                        offensive_power +=
                            (long)sh.item_quantity[i] * power(tons);
                    }
                    tons += 5;
                }
            }

            /* Adjust for ship aging. */
            offensive_power -= ((long)sh.age * offensive_power) / 50;
            defensive_power -= ((long)sh.age * defensive_power) / 50;
        }

        /* Adjust values for tech levels. */
        ml = c_species[species_index].tech_level[ML];
        ls = c_species[species_index].tech_level[LS];
        offensive_power += (ml * offensive_power) / 50;
        defensive_power += (ls * defensive_power) / 50;

        /* Adjust values if this species is hijacking anyone. */
        if (bat.hijacker[species_index] && (option == DEEP_SPACE_FIGHT ||
                                             option == PLANET_ATTACK)) {
            offensive_power /= 4;
            defensive_power /= 4;
        }

        /* Get number of shots per round. */
        n_shots = (offensive_power / 1500) + 1;
        if (ml == 0 || offensive_power == 0) {
            n_shots = 0;
        }
        if (n_shots > 5) {
            n_shots = 5;
        }
        act.num_shots[unit_index]  = n_shots;
        act.shots_left[unit_index] = n_shots;

        /* Get damage per shot. */
        if (n_shots > 0) {
            act.weapon_damage[unit_index] = (2 * offensive_power) / n_shots;
        }else{
            act.weapon_damage[unit_index] = 0;
        }

        /* Do defensive shields. */
        act.shield_strength[unit_index] = defensive_power;
        if (ttype == SHIP) {
            /* Adjust for results of previous action, if any. "dest_y"
             *  contains the percentage of shields that remained at end
             *  of last action. */
            defensive_power = ((long)sh.dest_y * defensive_power) / 100L;
        }
        act.shield_strength_left[unit_index] = defensive_power;

        /* Set bomb damage to zero in case this is planet bombardment or
         *      germ warfare. */
        act.bomb_damage[unit_index] = 0;

        /* Set flag for individual unit if species can be surprised. */
        if (bat.can_be_surprised[species_index]) {
            act.surprised[unit_index] = true;
        }else{
            act.surprised[unit_index] = false;
        }
    }

    return(true);        /* There will be a fight here. */
}



int disbanded_ship(species_index, sh)

int species_index;

struct ship_data *sh;

{
    int nampla_index;

    struct nampla_data *nam;

    nam = c_nampla[species_index] - 1;
    for (nampla_index = 0; nampla_index < c_species[species_index].num_namplas; nampla_index++) {
        ++nam;

        if (nam.x != sh.x) {
            continue;
        }
        if (nam.y != sh.y) {
            continue;
        }
        if (nam.z != sh.z) {
            continue;
        }
        if (nam.pn != sh.pn) {
            continue;
        }
        if ((nam.status & DISBANDED_COLONY) == 0) {
            continue;
        }
        if (sh.ttype != STARBASE && sh.status == IN_ORBIT) {
            continue;
        }

        /* This ship is either on the surface of a disbanded colony or is
         *      a starbase orbiting a disbanded colony. */
        return(true);
    }

    return(false);
}



//*************************************************************************
// for_jum.c


extern int   log_summary, ignore_field_distorters;
extern char  field_distorted[MAX_SPECIES];
extern FILE *log_file;
extern struct species_data *c_species[MAX_SPECIES];
extern struct galaxy_data   galaxy;


/* This routine will return true if forced jump or misjump units are used,
 * even if they fail. It will return false if the attacker has none or
 * not enough. */

int forced_jump_units_used(attacker_index, defender_index, total_shots,
                           bat, act)

int attacker_index, defender_index, *total_shots;

struct battle_data *bat;
struct action_data *act;

{
    int i, att_sp_index, def_sp_index, attacker_gv, defender_gv,
        ttype, fj_num, fm_num, number, success_chance, failure;

    char x, y, z;

    struct ship_data *attacking_ship, *defending_ship;



    /* Make sure attacking unit is a starbase. */
    attacking_ship = (struct ship_data *)act.fighting_unit[attacker_index];
    if (attacking_ship.ttype != STARBASE) {
        return(false);
    }

    /* See if attacker has any forced jump units. */
    fj_num = attacking_ship.item_quantity[FJ];
    fm_num = attacking_ship.item_quantity[FM];
    if (fj_num == 0 && fm_num == 0) {
        return(false);
    }

    /* If both ttypes are being carried, choose one randomly. */
    if (fj_num > 0 && fm_num > 0) {
        if (rnd(2) == 1) {
            ttype   = FJ;
            number = fj_num;
        }else {
            ttype   = FM;
            number = fm_num;
        }
    }else if (fj_num > 0) {
        ttype   = FJ;
        number = fj_num;
    }else {
        ttype   = FM;
        number = fm_num;
    }

    /* Get gravitics tech levels. */
    att_sp_index = act.fighting_species_index[attacker_index];
    attacker_gv  = c_species[att_sp_index].tech_level[GV];

    def_sp_index = act.fighting_species_index[defender_index];
    defender_gv  = c_species[def_sp_index].tech_level[GV];

    /* Check if sufficient units are available. */
    defending_ship = (struct ship_data *)act.fighting_unit[defender_index];
    if (number < defending_ship.tonnage) {
        return(false);
    }

    /* Make sure defender is not a starbase. */
    if (defending_ship.ttype == STARBASE) {
        return(false);
    }

    /* Calculate percent chance of success. */
    success_chance = 2 *
                     ((number - defending_ship.tonnage) + (attacker_gv - defender_gv));

    /* See if it worked. */
    failure = rnd(100) > success_chance;

    log_summary = !failure;

    log_string("        ");  log_string(ship_name(attacking_ship));
    log_string(" attempts to use ");
    log_string(item_name[ttype]);
    log_string("s against ");

    ignore_field_distorters = !field_distorted[def_sp_index];
    log_string(ship_name(defending_ship));
    ignore_field_distorters = false;

    if (failure) {
        log_string(", but fails.\n");
        return(true);
    }

    log_string(", and succeeds!\n");
    log_summary = false;

    /* Determine destination. */
    if (ttype == FM) {
        /* Destination is totally random. */
        x = rnd(100) - 1;
        y = rnd(100) - 1;
        z = rnd(100) - 1;
    }else {
        /* Random location close to battle. */
        i = 3;
        for (;i == 3;) {
            i = rnd(5);
        }
        x = bat.x + i - 3;
        if (x < 0) {
            x = 0;
        }

        i = 3;
        for (;i == 3;) {
            i = rnd(5);
        }
        y = bat.y + i - 3;
        if (y < 0) {
            y = 0;
        }

        i = 3;
        for (;i == 3;) {
            i = rnd(5);
        }
        z = bat.z + i - 3;
        if (z < 0) {
            z = 0;
        }
    }
    defending_ship.dest_x = x;
    defending_ship.dest_y = y;
    defending_ship.dest_z = z;

    /* Make sure this ship can no longer take part in the battle. */
    defending_ship.status          = FORCED_JUMP;
    defending_ship.pn              = -1;
    *total_shots                   -= act.shots_left[defender_index];
    act.shots_left[defender_index] = 0;
    act.num_shots[defender_index]  = 0;

    return(true);
}




//*************************************************************************
// gam_abo.c


gamemaster_abort_option() {
    char answer[16];

    /* Give the gamemaster a chance to abort. */
    printf("*** Gamemaster safe-abort option ... ttype q or Q to quit: ");
    fflush(stdout);
    fgets(answer, 16, stdin);
    if (answer[0] == 'q' || answer[0] == 'Q') {
        exit(0);
    }
}



//*************************************************************************
// get_gal.c


extern struct galaxy_data galaxy;


get_galaxy_data() {
    int galaxy_fd;

    long n, num_bytes, byte_size;


    /* Open galaxy file. */
    galaxy_fd = open("galaxy.dat", 0);
    if (galaxy_fd < 0) {
        fprintf(stderr, "\n\tCannot open file galaxy.dat!\n");
        exit(-1);
    }

    /* Read data. */
    byte_size = sizeof(struct galaxy_data);
    num_bytes = read(galaxy_fd, &galaxy, byte_size);
    if (num_bytes != byte_size) {
        fprintf(stderr, "\n\tCannot read data in file 'galaxy.dat'!\n\n");
        exit(-1);
    }

    close(galaxy_fd);
}



//*************************************************************************
// get_loc.c
/* This routine will assign values to global variables x, y, z, pn, star
 * and nampla. If the location is not a named planet, then nampla will be
 * set to NULL. If planet is not specified, pn will be set to zero. If
 * location is valid, true will be returned, otherwise false will be
 * returned. */



extern int  x, y, z, pn, num_stars, abbr_type;
extern long value;
extern char upper_name[32], *input_line_pointer;
extern struct species_data *species;
extern struct nampla_data * nampla_base, *nampla;
extern struct star_data *   star_base, *star;


int get_location() {
    int i, n, found, temp_nampla_index, first_try, name_length,
        best_score, next_best_score, best_nampla_index,
        minimum_score;

    char upper_nampla_name[32], *temp1_ptr, *temp2_ptr;

    struct nampla_data *temp_nampla;


    /* Check first if x, y, z are specified. */
    nampla = NULL;
    skip_whitespace();

    if (get_value() == 0) {
        goto get_planet;
    }
    x = value;

    if (get_value() == 0) {
        return(false);
    }
    y = value;

    if (get_value() == 0) {
        return(false);
    }
    z = value;

    if (get_value() == 0) {
        pn = 0;
    }else{
        pn = value;
    }

    if (pn == 0) {
        return(true);
    }

    /* Get star. Check if planet exists. */
    found = false;
    star  = star_base - 1;
    for (i = 0; i < num_stars; i++) {
        ++star;

        if (star.x != x) {
            continue;
        }
        if (star.y != y) {
            continue;
        }
        if (star.z != z) {
            continue;
        }

        if (pn > star.num_planets) {
            return(false);
        }else{
            return(true);
        }
    }

    return(false);


get_planet:

    /* Save pointers in case of error. */
    temp1_ptr = input_line_pointer;

    get_class_abbr();

    temp2_ptr = input_line_pointer;

    first_try = true;

again:

    input_line_pointer = temp2_ptr;

    if (abbr_type != PLANET_ID && !first_try) {
        /* Assume abbreviation was accidentally omitted. */
        input_line_pointer = temp1_ptr;
    }

    /* Get planet name. */
    get_name();

    /* Search all temp_namplas for name. */
    temp_nampla = nampla_base - 1;
    for (temp_nampla_index = 0; temp_nampla_index < species.num_namplas; temp_nampla_index++) {
        ++temp_nampla;

        if (temp_nampla.pn == 99) {
            continue;
        }

        /* Make upper case copy of temp_nampla name. */
        for (i = 0; i < 32; i++) {
            upper_nampla_name[i] = toupper(temp_nampla.name[i]);
        }

        /* Compare names. */
        if (strcmp(upper_nampla_name, upper_name) == 0) {
            goto done;
        }
    }

    if (first_try) {
        first_try = false;
        goto again;
    }


    /* Possibly a spelling error.  Find the best match that is approximately
     *  the same. */

    first_try = true;

yet_again:

    input_line_pointer = temp2_ptr;

    if (abbr_type != PLANET_ID && !first_try) {
        /* Assume abbreviation was accidentally omitted. */
        input_line_pointer = temp1_ptr;
    }

    /* Get planet name. */
    get_name();

    best_score      = -9999;
    next_best_score = -9999;
    for (temp_nampla_index = 0; temp_nampla_index < species.num_namplas; temp_nampla_index++) {
        temp_nampla = nampla_base + temp_nampla_index;

        if (temp_nampla.pn == 99) {
            continue;
        }

        /* Make upper case copy of temp_nampla name. */
        for (i = 0; i < 32; i++) {
            upper_nampla_name[i] = toupper(temp_nampla.name[i]);
        }

        /* Compare names. */
        n = agrep_score(upper_nampla_name, upper_name);
        if (n > best_score) {
            best_score        = n; /* Best match so far. */
            best_nampla_index = temp_nampla_index;
        }else if (n > next_best_score) {
            next_best_score = n;
        }
    }

    temp_nampla   = nampla_base + best_nampla_index;
    name_length   = strlen(temp_nampla.name);
    minimum_score = name_length - ((name_length / 7) + 1);

    if (best_score < minimum_score ||           /* Score too low. */
        name_length < 5 ||                      /* No errors allowed. */
        best_score == next_best_score) {        /* Another name with equal
                                                 *      score. */
        if (first_try) {
            first_try = false;
            goto yet_again;
        }else {
            return(false);
        }
    }

done:

    abbr_type = PLANET_ID;

    x      = temp_nampla.x;
    y      = temp_nampla.y;
    z      = temp_nampla.z;
    pn     = temp_nampla.pn;
    nampla = temp_nampla;

    return(true);
}




//*************************************************************************
// get_plan.c




int num_planets, planet_data_modified;

struct planet_data *planet_base;


get_planet_data() {
    int planet_fd;

    long n, data_size, mem_size;


    /* Open planet file. */
    planet_fd = open("planets.dat", 0);
    if (planet_fd < 0) {
        fprintf(stderr, "\n\tCannot open file planets.dat!\n");
        exit(-1);
    }

    /* Read header data. */
    data_size = read(planet_fd, &num_planets, sizeof(num_planets));
    if (data_size != sizeof(num_planets)) {
        fprintf(stderr, "\n\tCannot read num_planets in file 'planets.dat'!\n\n");
        exit(-1);
    }

    /* Allocate enough memory for all planets. */
    mem_size =
        (long)(num_planets + NUM_EXTRA_PLANETS) * (long)sizeof(struct planet_data);
    data_size =
        (long)num_planets * (long)sizeof(struct planet_data);
    planet_base = (struct planet_data *)malloc(mem_size);
    if (planet_base == NULL) {
        fprintf(stderr, "\nCannot allocate enough memory for planet file!\n\n");
        exit(-1);
    }

    /* Read it all into memory. */
    n = read(planet_fd, planet_base, data_size);
    if (n != data_size) {
        fprintf(stderr, "\nCannot read planet file into memory!\n\n");
        exit(-1);
    }
    close(planet_fd);

    planet_data_modified = false;
}



//*************************************************************************
// get_ship.c

/* The following routine will return true and set global variables "ship" and
 * "ship_index" if a valid ship designation is found. Otherwise, it will return
 * false. The algorithm employed allows minor spelling errors, as well as
 * accidental deletion of a ship abbreviation. */



int correct_spelling_required = false;


extern int  ship_index, abbr_type, abbr_index;
extern char upper_name[32], *input_line_pointer;
extern struct species_data *species;
extern struct ship_data *   ship_base, *ship;


int get_ship() {
    int i, n, name_length, best_score, next_best_score, best_ship_index,
        first_try, minimum_score;

    char upper_ship_name[32], *temp1_ptr, *temp2_ptr;

    struct ship_data *best_ship = NULL;


    /* Save in case of an error. */
    temp1_ptr = input_line_pointer;

    /* Get ship abbreviation. */
    if (get_class_abbr() == PLANET_ID) {
        input_line_pointer = temp1_ptr;
        return(false);
    }

    temp2_ptr = input_line_pointer;

    first_try = true;

again:

    input_line_pointer = temp2_ptr;

    if (abbr_type != SHIP_CLASS && !first_try) {
        /* Assume abbreviation was accidentally omitted. */
        input_line_pointer = temp1_ptr;
    }

    /* Get ship name. */
    name_length = get_name();

    /* Search all ships for name. */
    ship = ship_base - 1;
    for (ship_index = 0; ship_index < species.num_ships; ship_index++) {
        ++ship;

        if (ship.pn == 99) {
            continue;
        }

        /* Make upper case copy of ship name. */
        for (i = 0; i < 32; i++) {
            upper_ship_name[i] = toupper(ship.name[i]);
        }

        /* Compare names. */
        if (strcmp(upper_ship_name, upper_name) == 0) {
            abbr_type  = SHIP_CLASS;
            abbr_index = ship.class;
            correct_spelling_required = false;
            return(true);
        }
    }

    if (first_try) {
        first_try = false;
        goto again;
    }

    if (correct_spelling_required) {
        correct_spelling_required = false;
        return(false);
    }


    /* Possibly a spelling error.  Find the best match that is approximately
     *  the same. */

    first_try = true;

yet_again:

    input_line_pointer = temp2_ptr;

    if (abbr_type != SHIP_CLASS && !first_try) {
        /* Assume abbreviation was accidentally omitted. */
        input_line_pointer = temp1_ptr;
    }

    /* Get ship name. */
    name_length = get_name();

    best_score      = -9999;
    next_best_score = -9999;
    for (ship_index = 0; ship_index < species.num_ships; ship_index++) {
        ship = ship_base + ship_index;

        if (ship.pn == 99) {
            continue;
        }

        /* Make upper case copy of ship name. */
        for (i = 0; i < 32; i++) {
            upper_ship_name[i] = toupper(ship.name[i]);
        }

        n = agrep_score(upper_ship_name, upper_name);
        if (n > best_score) {
            /* Best match so far. */
            best_score      = n;
            best_ship       = ship;
            best_ship_index = ship_index;
        }else if (n > next_best_score) {
            next_best_score = n;
        }
    }

    if (best_ship == NULL) {
        return(false);
    }
    name_length   = strlen(best_ship.name);
    minimum_score = name_length - ((name_length / 7) + 1);

    if (best_score < minimum_score ||           /* Score too low. */
        name_length < 5 ||                      /* No errors allowed. */
        best_score == next_best_score) {        /* Another name with equal
                                                 *      score. */
        if (first_try) {
            first_try = false;
            goto yet_again;
        }else {
            correct_spelling_required = false;
            return(false);
        }
    }

    ship       = best_ship;
    ship_index = best_ship_index;
    abbr_type  = SHIP_CLASS;
    abbr_index = ship.class;
    correct_spelling_required = false;
    return(true);
}



//*************************************************************************
// get_spnam.c

/* This routine will get a species name and return true if found and if
 * it is valid.  It will also set global values "g_species_number" and
 * "g_species_name". The algorithm employed allows minor spelling errors,
 * as well as accidental deletion of the SP abbreviation. */



int  g_spec_number;
char g_spec_name[32];

extern int  abbr_type;
extern char upper_name[32], *input_line_pointer;
extern struct galaxy_data galaxy;


int get_species_name() {
    int i, n, species_index, best_score, best_species_index,
        next_best_score, first_try, minimum_score, name_length;

    char sp_name[32], *temp1_ptr, *temp2_ptr;

    struct species_data *sp;


    g_spec_number = 0;

    /* Save pointers in case of error. */
    temp1_ptr = input_line_pointer;

    get_class_abbr();

    temp2_ptr = input_line_pointer;

    first_try = true;

again:

    input_line_pointer = temp2_ptr;

    if (abbr_type != SPECIES_ID && !first_try) {
        /* Assume abbreviation was accidentally omitted. */
        input_line_pointer = temp1_ptr;
    }

    /* Get species name. */
    get_name();

    for (species_index = 0; species_index < galaxy.num_species; species_index++) {
        if (!data_in_memory[species_index]) {
            continue;
        }

        sp = &spec_data[species_index];

        /* Copy name to g_spec_name and convert it to upper case. */
        for (i = 0; i < 31; i++) {
            g_spec_name[i] = sp.name[i];
            sp_name[i]     = toupper(g_spec_name[i]);
        }

        if (strcmp(sp_name, upper_name) == 0) {
            g_spec_number = species_index + 1;
            abbr_type     = SPECIES_ID;
            return(true);
        }
    }

    if (first_try) {
        first_try = false;
        goto again;
    }

    /* Possibly a spelling error.  Find the best match that is approximately
     *  the same. */

    first_try = true;

yet_again:

    input_line_pointer = temp2_ptr;

    if (abbr_type != SPECIES_ID && !first_try) {
        /* Assume abbreviation was accidentally omitted. */
        input_line_pointer = temp1_ptr;
    }

    /* Get species name. */
    get_name();

    best_score      = -9999;
    next_best_score = -9999;
    for (species_index = 0; species_index < galaxy.num_species; species_index++) {
        if (!data_in_memory[species_index]) {
            continue;
        }

        sp = &spec_data[species_index];

        /* Convert name to upper case. */
        for (i = 0; i < 31; i++) {
            sp_name[i] = toupper(sp.name[i]);
        }

        n = agrep_score(sp_name, upper_name);
        if (n > best_score) {
            /* Best match so far. */
            best_score         = n;
            best_species_index = species_index;
        }else if (n > next_best_score) {
            next_best_score = n;
        }
    }

    sp            = &spec_data[best_species_index];
    name_length   = strlen(sp.name);
    minimum_score = name_length - ((name_length / 7) + 1);

    if (best_score < minimum_score ||           /* Score too low. */
        name_length < 5 ||                      /* No errors allowed. */
        best_score == next_best_score) {        /* Another name with equal
                                                 *      score. */
        if (first_try) {
            first_try = false;
            goto yet_again;
        }else {
            return(false);
        }
    }

    /* Copy name to g_spec_name. */
    for (i = 0; i < 31; i++) {
        g_spec_name[i] = sp.name[i];
    }
    g_spec_number = best_species_index + 1;
    abbr_type     = SPECIES_ID;
    return(true);
}




//*************************************************************************
// get_star.c





int num_stars, star_data_modified;

struct star_data *star_base;


get_star_data() {
    int star_fd;

    long byte_size, star_data_size, mem_size;


    /* Open star file. */
    star_fd = open("stars.dat", 0);
    if (star_fd < 0) {
        fprintf(stderr, "\n\tCannot open file stars.dat!\n");
        exit(999);
    }

    byte_size = read(star_fd, &num_stars, sizeof(num_stars));
    if (byte_size != sizeof(num_stars)) {
        fprintf(stderr, "\n\tCannot read num_stars in file 'stars.dat'!\n\n");
        exit(999);
    }

    /* Allocate enough memory for all stars. */
    mem_size =
        (long)(num_stars + NUM_EXTRA_STARS) * (long)sizeof(struct star_data);
    star_data_size =
        (long)num_stars * (long)sizeof(struct star_data);
    star_base = (struct star_data *)malloc(mem_size);
    if (star_base == NULL) {
        fprintf(stderr, "\nCannot allocate enough memory for star file!\n\n");
        exit(-1);
    }

    /* Read it all into memory. */
    byte_size = read(star_fd, star_base, star_data_size);
    if (byte_size != star_data_size) {
        fprintf(stderr, "\nCannot read star file into memory!\n\n");
        exit(-1);
    }
    close(star_fd);

    star_data_modified = false;
}



//*************************************************************************
// get_transact.c


int num_transactions;

struct trans_data transaction[MAX_TRANSACTIONS];


get_transaction_data() {
    int  i, trans_fd;
    long num_bytes;


    /* Open file for reading. */
    trans_fd = open("interspecies.dat", 0);

    if (trans_fd < 0) {
        num_transactions = 0;
        return;
    }

    /* Read transactions from file. */
    i = 0;
    for {
        num_bytes = read(trans_fd, &transaction[i], sizeof(struct trans_data));
        if (num_bytes == 0) {
            break;                      /* End of file. */
        }
        if (num_bytes != sizeof(struct trans_data)) {
            fprintf(stderr, "\n\n\tError reading transaction from file 'interspecies.dat'!\n\n");
            exit(-1);
        }

        ++i;
    }

    num_transactions = i;

    close(trans_fd);
}



//*************************************************************************
// get_transfer.c


extern int   correct_spelling_required;
extern char *input_line_pointer;
extern struct nampla_data *nampla;


int get_transfer_point() {
    char *temp_ptr;


    /* Find out if it is a ship or a planet. First try for a correctly
     *  spelled ship name. */
    temp_ptr = input_line_pointer;
    correct_spelling_required = true;
    if (get_ship()) {
        return(true);
    }

    /* Probably not a ship. See if it's a planet. */
    input_line_pointer = temp_ptr;
    if (get_location()) {
        return(nampla != NULL);
    }

    /* Now check for an incorrectly spelled ship name. */
    input_line_pointer = temp_ptr;
    if (get_ship()) {
        return(true);
    }

    return(false);
}




//*************************************************************************
// money.c


long balance, raw_material_units, production_capacity, EU_spending_limit;


extern FILE *log_file;
extern struct species_data *species;
extern struct nampla_data * nampla;


int check_bounced(amount_needed)

long amount_needed;

{
    long take_from_EUs, limiting_balance;


    /* Check if we have sufficient funds for this purchase. */
    if (amount_needed > balance) {
        take_from_EUs = amount_needed - balance;

        if (take_from_EUs <= EU_spending_limit &&
            take_from_EUs <= species.econ_units) {
            species.econ_units -= take_from_EUs;
            EU_spending_limit   -= take_from_EUs;
            balance              = amount_needed;
        }else {
            return(true);
        }
    }

    /* Reduce various balances appropriately. */
    if (raw_material_units >= amount_needed) {
        if (production_capacity >= amount_needed) {
            /* Enough of both. */
            raw_material_units  -= amount_needed;
            production_capacity -= amount_needed;
        }else {
            /* Enough RMs but not enough PC. */
            raw_material_units -= production_capacity;
            production_capacity = 0;
        }
    }else {
        if (production_capacity >= amount_needed) {
            /* Enough PC but not enough RMs. */
            production_capacity -= raw_material_units;
            raw_material_units   = 0;
        }else {
            /* Not enough RMs or PC. */
            limiting_balance = (raw_material_units > production_capacity)
                               ? production_capacity : raw_material_units;
            raw_material_units  -= limiting_balance;
            production_capacity -= limiting_balance;
        }
    }

    balance -= amount_needed;

    return(false);
}



transfer_balance() {
    long limiting_amount;


    /* Log end of production. Do not print ending balance for mining
     *  or resort colonies. */
    limiting_amount = 0;
    fprintf(log_file, "  End of production on PL %s.", nampla.name);
    if (!(nampla.status & (MINING_COLONY | RESORT_COLONY))) {
        limiting_amount = (raw_material_units > production_capacity)
                          ? production_capacity : raw_material_units;
        fprintf(log_file, " (Ending balance is %ld.)", limiting_amount);
    }
    fprintf(log_file, "\n");

    /* Convert unused balance to economic units. */
    species.econ_units += limiting_amount;
    raw_material_units  -= limiting_amount;

    /* Carry over unused raw material units into next turn. */
    nampla.item_quantity[RM] += raw_material_units;

    balance = 0;
}




//*************************************************************************
// parse.c


int   end_of_file = false;
int   abbr_type, abbr_index, sub_light, tonnage, just_opened_file;
char  input_abbr[256], input_line[256], original_line[256], *input_line_pointer;
char  original_name[32], upper_name[32];
long  value;
FILE *input_file;


/* Skip white space and comments. */
skip_junk() {
again:

    /* Read next line. */
    input_line_pointer = fgets(input_line, 256, input_file);
    if (input_line_pointer == NULL) {
        end_of_file = true;
        return;
    }

    if (just_opened_file) {     /* Skip mail header, if any. */
        if (*input_line == '\n') {
            goto again;
        }

        just_opened_file = false;

        if (strncmp(input_line, "From ", 5) == 0) { /* This is a mail header. */
            for {
                input_line_pointer = fgets(input_line, 256, input_file);
                if (input_line_pointer == NULL) {
                    end_of_file = true;         /* Weird. */
                    return;
                }
                if (*input_line == '\n') {
                    break;                         /* End of header. */
                }
            }

            goto again;
        }
    }

    strcpy(original_line, input_line);          /* Make a copy. */

    /* Skip white space and comments. */
    for {
        switch (*input_line_pointer) {
        case ';':                       /* Semi-colon. */
        case '\n':                      /* Newline. */
            goto again;

        case '\t':                      /* Tab. */
        case ' ':                       /* Space. */
        case ',':                       /* Comma. */
            ++input_line_pointer;
            continue;

        default:
            return;
        }
    }
}


skip_whitespace() {
    for {
        switch (*input_line_pointer) {
        case '\t':                      /* Tab. */
        case ' ':                       /* Space. */
        case ',':                       /* Comma. */
            ++input_line_pointer;
            break;

        default:
            return;
        }
    }
}



/* The following "get" routines will return 0 if the item found was not
 * of the appropriate ttype, and 1 or greater if an item of the correct
 * ttype was found. */


/* Get a command and return its index. */
int get_command() {
    int  i, cmd_n;
    char c, cmd_s[4];


    skip_junk();
    if (end_of_file) {
        return(-1);
    }

    c = *input_line_pointer;
    /* Get first three characters of command word. */
    for (i = 0; i < 3; i++) {
        if (!isalpha(c)) {
            return(0);
        }
        cmd_s[i] = toupper(c);
        ++input_line_pointer;
        c = *input_line_pointer;
    }
    cmd_s[3] = '\0';

    /* Skip everything after third character of command word. */
    for {
        switch (c) {
        case '\t':
        case '\n':
        case ' ':
        case ',':
        case ';':
            goto find_cmd;

        default:
            ++input_line_pointer;
            c = *input_line_pointer;
        }
    }

find_cmd:

    /* Find corresponding string in list. */
    cmd_n = UNKNOWN;
    for (i = 1; i < NUM_COMMANDS; i++) {
        if (strcmp(cmd_s, command_abbr[i]) == 0) {
            cmd_n = i;
            break;
        }
    }

    return(cmd_n);
}

/* Get a class abbreviation and return TECH_ID, ITEM_CLASS, SHIP_CLASS,
 * PLANET_ID, SPECIES_ID or ALLIANCE_ID as appropriate, or UNKNOWN if it
 * cannot be identified. Also, set "abbr_type" to this value. If it is
 * TECH_ID, ITEM_CLASS or SHIP_CLASS, "abbr_index" will contain the
 * abbreviation index. If it is a ship, "tonnage" will contain tonnage/10,000,
 * and "sub_light" will be true or false. (Tonnage value returned is based
 * ONLY on abbreviation.) */

int get_class_abbr() {
    int i;

    char *digit_start;


    skip_whitespace();

    abbr_type = UNKNOWN;

    if (!isalnum(*input_line_pointer)) {
        return(UNKNOWN);
    }
    input_abbr[0] = toupper(*input_line_pointer);
    ++input_line_pointer;

    if (!isalnum(*input_line_pointer)) {
        return(UNKNOWN);
    }
    input_abbr[1] = toupper(*input_line_pointer);
    ++input_line_pointer;

    input_abbr[2] = '\0';

    /* Check for IDs that are followed by one or more digits or letters. */
    i           = 2;
    digit_start = input_line_pointer;
    for (;isalnum(*input_line_pointer);) {
        input_abbr[i++] = *input_line_pointer++;
        input_abbr[i]   = '\0';
    }

    /* Check tech ID. */
    for (i = 0; i < 6; i++) {
        if (strcmp(input_abbr, tech_abbr[i]) == 0) {
            abbr_index = i;
            abbr_type  = TECH_ID;
            return(abbr_type);
        }
    }

    /* Check item abbreviations. */
    for (i = 0; i < MAX_ITEMS; i++) {
        if (strcmp(input_abbr, item_abbr[i]) == 0) {
            abbr_index = i;
            abbr_type  = ITEM_CLASS;
            return(abbr_type);
        }
    }

    /* Check ship abbreviations. */
    for (i = 0; i < NUM_SHIP_CLASSES; i++) {
        if (strncmp(input_abbr, ship_abbr[i], 2) == 0) {
            input_line_pointer = digit_start;
            abbr_index         = i;
            tonnage            = ship_tonnage[i];
            if (i == TR) {
                tonnage = 0;
                for (;isdigit(*input_line_pointer);) {
                    tonnage = (10 * tonnage) + (*input_line_pointer - '0');
                    ++input_line_pointer;
                }
            }

            if (toupper(*input_line_pointer) == 'S') {
                sub_light = true;
                ++input_line_pointer;
            }else {
                sub_light = false;
            }

            if (isalnum(*input_line_pointer)) {
                break;                                  /* Garbage. */
            }
            abbr_type = SHIP_CLASS;
            return(abbr_type);
        }
    }

    /* Check for planet name. */
    if (strcmp(input_abbr, "PL") == 0) {
        abbr_type = PLANET_ID;
        return(abbr_type);
    }

    /* Check for species name. */
    if (strcmp(input_abbr, "SP") == 0) {
        abbr_type = SPECIES_ID;
        return(abbr_type);
    }

    abbr_type = UNKNOWN;
    return(abbr_type);
}

/* Get a name and copy original version to "original_name" and upper
 * case version to "upper_name". Return length of name. */
int get_name() {
    int name_length;

    char c;


    skip_whitespace();

    name_length = 0;
    for {
        c = *input_line_pointer;
        if (c == ';') {
            break;
        }
        ++input_line_pointer;
        if (c == ',' || c == '\t' || c == '\n') {
            break;
        }
        if (name_length < 31) {
            original_name[name_length] = c;
            upper_name[name_length]    = toupper(c);
            ++name_length;
        }
    }

    /* Remove any final spaces in name. */
    for (;name_length > 0;) {
        c = original_name[name_length - 1];
        if (c != ' ') {
            break;
        }
        --name_length;
    }

    /* Terminate strings. */
    original_name[name_length] = '\0';
    upper_name[name_length]    = '\0';

    return(name_length);
}

/* Read a long decimal and place its value in 'value'. */
int get_value() {
    int n;


    skip_whitespace();

    n = sscanf(input_line_pointer, "%ld", &value);
    if (n != 1) {
        return(0);              /* Not a numeric value. */
    }
    /* Skip numeric string. */
    ++input_line_pointer;       /* Skip first sign or digit. */
    for (;isdigit(*input_line_pointer);) {
        ++input_line_pointer;
    }

    return(1);
}

/* The following routine will check that the next argument in the current
 * command line is followed by a comma or tab.  If not present, it will
 * try to insert a comma in the proper position.  This routine should
 * be called only AFTER an error has been detected. */

fix_separator() {
    int n, first_class, fix_made, num_commas;

    char c, *temp_ptr, *temp2_ptr, *first_comma;


    skip_whitespace();

    if (isdigit(*input_line_pointer)) {
        return;                                 /* Nothing can be done. */
    }
    if (strchr(input_line_pointer, ' ') == NULL) {
        return;                                           /* Ditto. */
    }
    fix_made = false;

    /* Look for a ship, planet, or species abbreviation after the first one.
     *  If it is preceeded by a space, convert the space to a comma. */
    temp_ptr    = input_line_pointer;
    first_class = get_class_abbr();     /* Skip first one but remember what it was. */
    for {
        skip_whitespace();
        temp2_ptr = input_line_pointer - 1;
        if (*input_line_pointer == '\n') {
            break;
        }
        if (*input_line_pointer == ';') {
            break;
        }

        /* The following is to prevent an infinite loop. */
        if (!isalnum(*input_line_pointer)) {
            ++input_line_pointer;
            continue;
        }

        n = get_class_abbr();
        if (n == SHIP_CLASS || n == PLANET_ID || n == SPECIES_ID) {
            /* Convert space preceeding abbreviation to a comma. */
            if (*temp2_ptr == ' ') {
                *temp2_ptr = ',';
                fix_made   = true;
            }
        }
    }
    input_line_pointer = temp_ptr;

    if (fix_made) {
        return;
    }

    /* Look for a space followed by a digit. If found, convert the space
     *  to a comma.  If exactly two or four commas are added, re-convert
     *  the first one back to a space; e.g. Jump TR1 Seeker,7,99,99,99 or
     *  Build TR1 Seeker,7,50. */
    num_commas = 0;
    for {
        c = *temp_ptr++;

        if (c == '\n') {
            break;
        }
        if (c == ';') {
            break;
        }

        if (c != ' ') {
            continue;
        }
        if (isdigit(*temp_ptr)) {
            --temp_ptr;         /* Convert space to a comma. */
            *temp_ptr = ',';
            if (num_commas++ == 0) {
                first_comma = temp_ptr;
            }
            ++temp_ptr;
            fix_made = true;
        }
    }

    if (fix_made) {
        if (num_commas == 2 || num_commas == 4) {
            *first_comma = ' ';
        }
        return;
    }

    /* Now's the time for wild guesses. */
    temp_ptr = input_line_pointer;

    /* If first word is a valid abbreviation, put a comma after the
     *  second word. */
    if (first_class == SHIP_CLASS || first_class == PLANET_ID || first_class == SPECIES_ID) {
        temp_ptr = strchr(temp_ptr, ' ') + 1;
        temp_ptr = strchr(temp_ptr, ' ');
        if (temp_ptr != NULL) {
            *temp_ptr = ',';
        }
        return;
    }

    /* First word is not a valid abbreviation.  Put a comma after it. */
    temp_ptr = strchr(temp_ptr, ' ');
    if (temp_ptr != NULL) {
        *temp_ptr = ',';
    }
}




//*************************************************************************
// regen_sh.c


extern struct species_data *c_species[MAX_SPECIES];


regenerate_shields(act)

struct action_data *act;

{
    int i, species_index, unit_index;

    long ls, max_shield_strength, percent;

    struct ship_data *sh;


    /* Shields are regenerated by 5 + LS/10 percent per round. */
    for (unit_index = 0; unit_index < act.num_units_fighting; unit_index++) {
        species_index       = act.fighting_species_index[unit_index];
        ls                  = c_species[species_index].tech_level[LS];
        max_shield_strength = act.shield_strength[unit_index];

        percent = (ls / 10L) + 5L;
        act.shield_strength_left[unit_index] +=
            (percent * max_shield_strength) / 100L;
        if (act.shield_strength_left[unit_index] > max_shield_strength) {
            act.shield_strength_left[unit_index] = max_shield_strength;
        }
    }
}




//*************************************************************************
// sav_plan.c


extern int num_planets;

extern struct planet_data *planet_base;


save_planet_data() {
    int planet_fd;

    long n, byte_size;


    /* Open planet file for writing. */
    planet_fd = creat("planets.dat", 0600);
    if (planet_fd < 0) {
        fprintf(stderr, "\n\tCannot create file 'planets.dat'!\n");
        exit(-1);
    }

    /* Write header data. */
    byte_size = write(planet_fd, &num_planets, sizeof(num_planets));
    if (byte_size != sizeof(num_planets)) {
        fprintf(stderr, "\n\tCannot write num_planets to file 'planets.dat'!\n\n");
        exit(-1);
    }

    /* Write planet data to disk. */
    byte_size = (long)num_planets * sizeof(struct planet_data);
    n         = write(planet_fd, planet_base, byte_size);
    if (n != byte_size) {
        fprintf(stderr, "\nCannot write planet data to disk!\n\n");
        exit(-1);
    }

    close(planet_fd);
}




//*************************************************************************
// sav_star.c


extern int num_stars;

extern struct star_data *star_base;


save_star_data() {
    int star_fd;

    long n, byte_size;


    /* Open star file for writing. */
    star_fd = creat("stars.dat", 0600);
    if (star_fd < 0) {
        fprintf(stderr, "\n\tCannot create file 'stars.dat'!\n");
        exit(-1);
    }

    /* Write header data. */
    byte_size = write(star_fd, &num_stars, sizeof(num_stars));
    if (byte_size != sizeof(num_stars)) {
        fprintf(stderr, "\n\tCannot write num_stars to file 'stars.dat'!\n\n");
        exit(-1);
    }

    /* Write star data to disk. */
    byte_size = (long)num_stars * sizeof(struct star_data);
    n         = write(star_fd, star_base, byte_size);
    if (n != byte_size) {
        fprintf(stderr, "\nCannot write star data to disk!\n\n");
        exit(-1);
    }

    close(star_fd);
}




//*************************************************************************
// sav_transact.c


extern int num_transactions;

extern struct trans_data transaction[MAX_TRANSACTIONS];


save_transaction_data() {
    int  i, trans_fd;
    long num_bytes;


    /* Open file for writing. */
    trans_fd = creat("interspecies.dat", 0600);

    if (trans_fd < 0) {
        fprintf(stderr, "\n\n\tCannot create file 'interspecies.dat'!\n\n");
        exit(-1);
    }

    /* Write transactions to file. */
    for (i = 0; i < num_transactions; i++) {
        num_bytes = write(trans_fd, &transaction[i], sizeof(struct trans_data));

        if (num_bytes != sizeof(struct trans_data)) {
            fprintf(stderr, "\n\n\tError writing transaction to file 'interspecies.dat'!\n\n");
            exit(-1);
        }
    }

    close(trans_fd);
}




//*************************************************************************
// scan.c


int print_LSN = true;

extern int                  num_stars;
extern FILE *               log_file;
extern struct star_data *   star_base;
extern struct planet_data * planet_base;
extern struct species_data *species;
extern struct nampla_data * nampla_base;


scan(x, y, z)

char x, y, z;

{
    int i, j, k, n, found, num_gases, ls_needed;

    char filename[32];

    struct star_data *  star;
    struct planet_data *planet, *home_planet;
    struct nampla_data *home_nampla;


    /* Find star. */
    star  = star_base;
    found = false;
    for (i = 0; i < num_stars; i++) {
        if (star.x == x && star.y == y && star.z == z) {
            found = true;
            break;
        }
        ++star;
    }

    if (!found) {
        fprintf(log_file,
                "Scan Report: There is no star system at x = %d, y = %d, z = %d.\n",
                x, y, z);
        return;
    }

    /* Print data for star, */
    fprintf(log_file, "Coordinates:\tx = %d\ty = %d\tz = %d", x, y, z);
    fprintf(log_file, "\tstellar ttype = %c%c%c", ttype_char[star.ttype],
            color_char[star.color], size_char[star.size]);

    fprintf(log_file, "   %d planets.\n\n", star.num_planets);

    if (star.worm_here) {
        fprintf(log_file,
                "This star system is the terminus of a natural wormhole.\n\n");
    }

    /* Print header. */
    fprintf(log_file, "               Temp  Press Mining\n");
    fprintf(log_file, "  #  Dia  Grav Class Class  Diff  LSN  Atmosphere\n");
    fprintf(log_file, " ---------------------------------------------------------------------\n");

    /* Check for nova. */
    if (star.num_planets == 0) {
        fprintf(log_file, "\n\tThis star is a nova remnant. Any planets it may have once\n");
        fprintf(log_file, "\thad have been blown away.\n\n");
        return;
    }

    /* Print data for each planet. */
    planet = planet_base + (long)star.planet_index;
    if (print_LSN) {
        home_nampla = nampla_base;
        home_planet = planet_base + (long)home_nampla.planet_index;
    }

    for (i = 1; i <= star.num_planets; i++) {
        /* Get life support tech level needed. */
        if (print_LSN) {
            ls_needed = life_support_needed(species, home_planet, planet);
        }else{
            ls_needed = 99;
        }

        fprintf(log_file, "  %d  %3d  %d.%02d  %2d    %2d    %d.%02d %4d  ",
                i,
                planet.diameter,
                planet.gravity / 100,
                planet.gravity % 100,
                planet.temperature_class,
                planet.pressure_class,
                planet.mining_difficulty / 100,
                planet.mining_difficulty % 100,
                ls_needed);

        num_gases = 0;
        for (n = 0; n < 4; n++) {
            if (planet.gas_percent[n] > 0) {
                if (num_gases > 0) {
                    fprintf(log_file, ",");
                }
                fprintf(log_file, "%s(%d%%)", gas_string[planet.gas[n]],
                        planet.gas_percent[n]);
                ++num_gases;
            }
        }

        if (num_gases == 0) {
            fprintf(log_file, "No atmosphere");
        }

        fprintf(log_file, "\n");
        ++planet;
    }

    if (star.message) {
        /* There is a message that must be logged whenever this star
         *      system is scanned. */
        sprintf(filename, "message%ld.txt\0", star.message);
        log_message(filename);
    }

    return;
}




//*************************************************************************
// utils.c

/* This routine will return a random int between 1 and max, inclusive.
 * It uses the so-called "Algorithm M" method, which is a combination
 * of the congruential and shift-register methods. */

unsigned long last_random = 1924085713L;        /* Random seed. */

int rnd(max)

unsigned int max;

{
    unsigned long a, b, c, cong_result, shift_result;

    /* For congruential method, multiply previous value by the
     * prime number 16417. */
    a           = last_random;
    b           = last_random << 5;
    c           = last_random << 14;
    cong_result = a + b + c;            /* Effectively multiply by 16417. */

    /* For shift-register method, use shift-right 15 and shift-left 17
     * with no-carry addition (i.e., exclusive-or). */
    a             = last_random >> 15;
    shift_result  = a ^ last_random;
    a             = shift_result << 17;
    shift_result ^= a;

    last_random = cong_result ^ shift_result;

    a = last_random & 0x0000FFFF;

    return((int)((a * (long)max) >> 16) + 1L);
}



/* Routine "get_species_data" will read in data files for all species,
 *      "save_species_data" will write all data that has been modified, and
 *      "free_species_data" will free memory used for all species data. */

/* Additional memory must be allocated for routines that build ships or
 * name planets. Here are the default 'extras', which may be changed, if
 * necessary, by the main program. */

long extra_namplas = NUM_EXTRA_NAMPLAS;
long extra_ships   = NUM_EXTRA_SHIPS;

extern struct galaxy_data galaxy;


get_species_data() {
    int species_fd, species_index;

    long n, num_bytes;

    char filename[16];

    struct species_data *sp;


    for (species_index = 0; species_index < galaxy.num_species; species_index++) {
        data_modified[species_index] = false;

        sp = &spec_data[species_index];

        /* Open the species data file. */
        sprintf(filename, "sp%02d.dat\0", species_index + 1);
        species_fd = open(filename, 0);
        if (species_fd < 0) {
            sp.pn = 0; /* Extinct! */
            data_in_memory[species_index] = false;
            continue;
        }

        /* Read in species data. */
        num_bytes = read(species_fd, sp, sizeof(struct species_data));
        if (num_bytes != sizeof(struct species_data)) {
            fprintf(stderr, "\n\tCannot read species record in file '%s'!\n\n",
                    filename);
            exit(-1);
        }

        /* Allocate enough memory for all namplas. */
        num_bytes = (sp.num_namplas + extra_namplas) * sizeof(struct nampla_data);
        namp_data[species_index] = (struct nampla_data *)malloc(num_bytes);
        if (namp_data[species_index] == NULL) {
            fprintf(stderr, "\nCannot allocate enough memory for nampla data!\n\n");
            exit(-1);
        }

        /* Read it all into memory. */
        num_bytes = (long)sp.num_namplas * sizeof(struct nampla_data);
        n         = read(species_fd, namp_data[species_index], num_bytes);
        if (n != num_bytes) {
            fprintf(stderr, "\nCannot read nampla data into memory!\n\n");
            exit(-1);
        }

        /* Allocate enough memory for all ships. */
        num_bytes = (sp.num_ships + extra_ships) * sizeof(struct ship_data);
        ship_data[species_index] = (struct ship_data *)malloc(num_bytes);
        if (ship_data[species_index] == NULL) {
            fprintf(stderr, "\nCannot allocate enough memory for ship data!\n\n");
            exit(-1);
        }

        if (sp.num_ships > 0) {
            /* Read it all into memory. */
            num_bytes = (long)sp.num_ships * sizeof(struct ship_data);
            n         = read(species_fd, ship_data[species_index], num_bytes);
            if (n != num_bytes) {
                fprintf(stderr, "\nCannot read ship data into memory!\n\n");
                exit(-1);
            }
        }

        close(species_fd);

        data_in_memory[species_index]  = true;
        num_new_namplas[species_index] = 0;
        num_new_ships[species_index]   = 0;
    }
}



save_species_data() {
    int species_fd, species_index;

    long n, num_bytes;

    char filename[16];

    struct species_data *sp;


    for (species_index = 0; species_index < galaxy.num_species; species_index++) {
        if (!data_modified[species_index]) {
            continue;
        }

        sp = &spec_data[species_index];

        /* Open the species data file. */
        sprintf(filename, "sp%02d.dat\0", species_index + 1);
        species_fd = creat(filename, 0600);
        if (species_fd < 0) {
            fprintf(stderr, "\n  Cannot create new version of file '%s'!\n",
                    filename);
            exit(-1);
        }

        /* Write species data. */
        num_bytes = write(species_fd, sp, sizeof(struct species_data));
        if (num_bytes != sizeof(struct species_data)) {
            fprintf(stderr, "\n\tCannot write species record to file '%s'!\n\n",
                    filename);
            exit(-1);
        }

        /* Write nampla data. */
        num_bytes = sp.num_namplas * sizeof(struct nampla_data);
        n         = write(species_fd, namp_data[species_index], num_bytes);
        if (n != num_bytes) {
            fprintf(stderr, "\nCannot write nampla data to file!\n\n");
            exit(-1);
        }

        if (sp.num_ships > 0) {
            /* Write ship data. */
            num_bytes = (long)sp.num_ships * sizeof(struct ship_data);
            n         = write(species_fd, ship_data[species_index], num_bytes);
            if (n != num_bytes) {
                fprintf(stderr, "\nCannot write ship data to file!\n\n");
                exit(-1);
            }
        }

        close(species_fd);

        data_modified[species_index] = false;
    }
}



free_species_data() {
    int species_index;


    for (species_index = 0; species_index < galaxy.num_species; species_index++) {
        if (data_in_memory[species_index]) {
            free(namp_data[species_index]);

            if (spec_data[species_index].num_ships > 0) {
                free(ship_data[species_index]);
            }

            data_in_memory[species_index] = false;
            data_modified[species_index]  = false;
        }
    }
}



/* The following two routines will delete a ship or nampla record. */

delete_ship(ship)

struct ship_data *ship;

{
    int i;

    char *cp;


    /* Set all bytes of record to zero. */
    cp = (char *)ship;
    for (i = 0; i < sizeof(struct ship_data); i++) {
        *cp++ = 0;
    }

    ship.pn = 99;
    strcpy(ship.name, "Unused");
}


delete_nampla(nampla)

struct nampla_data *nampla;

{
    int i;

    char *cp;


    /* Set all bytes of record to zero. */
    cp = (char *)nampla;
    for (i = 0; i < sizeof(struct nampla_data); i++) {
        *cp++ = 0;
    }

    nampla.pn = 99;
    strcpy(nampla.name, "Unused");
}



/* This routine is intended to take a long argument and return a pointer
 * to a string that has embedded commas to make the string more readable. */

char result_plus_commas[33];

char *commas(value)

long value;

{
    int i, j, n, length, negative;

    char *ptr, temp[32];

    long abs_value;


    if (value < 0) {
        abs_value = -value;
        negative  = true;
    }else {
        abs_value = value;
        negative  = false;
    }

    sprintf(temp, "%ld\0", abs_value);

    length = strlen(temp);

    i = length - 1;
    j = 31;
    result_plus_commas[32] = '\0';
    for (n = 0; n < length; n++) {
        result_plus_commas[j--] = temp[i--];
        if (j % 4 == 0) {
            result_plus_commas[j--] = ',';
        }
    }

    j++;
    if (result_plus_commas[j] == ',') {
        j++;
    }

    if (negative) {
        result_plus_commas[--j] = '-';
    }

    return(&result_plus_commas[j]);
}



/* This routine will return a pointer to a string containing a complete
 * ship name, including its orbital/landed status and age. If global
 * variable "truncate_name" is true, then orbital/landed status and age
 * will not be included. */

int truncate_name           = false;
int ignore_field_distorters = false;

char full_ship_id[64];

char *ship_name(ship)

struct ship_data *ship;

{
    int effective_age, status, ship_is_distorted;

    char temp[16];


    if (ship.item_quantity[FD] == ship.tonnage) {
        ship_is_distorted = true;
    }else{
        ship_is_distorted = false;
    }

    if (ship.status == ON_SURFACE) {
        ship_is_distorted = false;
    }

    if (ignore_field_distorters) {
        ship_is_distorted = false;
    }

    if (ship_is_distorted) {
        if (ship.class == TR) {
            sprintf(full_ship_id, "%s%d ???\0", ship_abbr[ship.class],
                    ship.tonnage);
        }else if (ship.class == BA) {
            sprintf(full_ship_id, "BAS ???\0");
        }else{
            sprintf(full_ship_id, "%s ???\0", ship_abbr[ship.class]);
        }
    }else if (ship.class == TR) {
        sprintf(full_ship_id, "%s%d%s %s\0",
                ship_abbr[ship.class], ship.tonnage, ship_ttype[ship.ttype],
                ship.name);
    }else {
        sprintf(full_ship_id, "%s%s %s\0",
                ship_abbr[ship.class], ship_ttype[ship.ttype], ship.name);
    }

    if (truncate_name) {
        return(&full_ship_id[0]);
    }

    strcat(full_ship_id, " (");

    effective_age = ship.age;
    if (effective_age < 0) {
        effective_age = 0;
    }

    if (!ship_is_distorted) {
        if (ship.status != UNDER_CONSTRUCTION) {
            /* Do age. */
            sprintf(temp, "A%d,\0", effective_age);
            strcat(full_ship_id, temp);
        }
    }

    status = ship.status;
    switch (status) {
    case UNDER_CONSTRUCTION:
        sprintf(temp, "C\0");
        break;

    case IN_ORBIT:
        sprintf(temp, "O%d\0", ship.pn);
        break;

    case ON_SURFACE:
        sprintf(temp, "L%d\0", ship.pn);
        break;

    case IN_DEEP_SPACE:
        sprintf(temp, "D\0");
        break;

    case FORCED_JUMP:
        sprintf(temp, "FJ\0");
        break;

    case JUMPED_IN_COMBAT:
        sprintf(temp, "WD\0");
        break;

    default:
        sprintf(temp, "***???***\0");
        fprintf(stderr, "\n\tWARNING!!!  Internal error in subroutine 'ship_name'\n\n");
    }

    strcat(full_ship_id, temp);

    if (ship.ttype == STARBASE) {
        sprintf(temp, ",%ld tons\0", 10000L * (long)ship.tonnage);
        strcat(full_ship_id, temp);
    }

    strcat(full_ship_id, ")");

    return(&full_ship_id[0]);
}



/* The following routines will post an item to standard output and to
 * an externally defined log file and summary file. */

FILE *log_file, *summary_file;

int log_start_of_line = true;
int log_indentation   = 0;
int log_position      = 0;
int logging_disabled  = false;
int log_to_file       = true;
int log_summary       = false;
int log_stdout        = true;

char log_line[128];


log_char(c)

char c;

{
    int i, temp_position;

    char temp_char;


    if (logging_disabled) {
        return;
    }

    /* Check if current line is getting too long. */
    if ((c == ' ' || c == '\n') && log_position > 77) {
        /* Find closest preceeding space. */
        temp_position = log_position - 1;
        for (;log_line[temp_position] != ' ';) {
            --temp_position;
        }

        /* Write front of line to files. */
        temp_char = log_line[temp_position + 1];
        log_line[temp_position]     = '\n';
        log_line[temp_position + 1] = '\0';
        if (log_to_file) {
            fputs(log_line, log_file);
        }
        if (log_stdout) {
            fputs(log_line, stdout);
        }
        if (log_summary) {
            fputs(log_line, summary_file);
        }
        log_line[temp_position + 1] = temp_char;

        /* Copy overflow word to beginning of next line. */
        log_line[log_position] = '\0';
        log_position           = log_indentation + 2;
        for (i = 0; i < log_position; i++) {
            log_line[i] = ' ';
        }
        strcpy(&log_line[log_position], &log_line[temp_position + 1]);

        log_position = strlen(log_line);

        if (c == ' ') {
            log_line[log_position++] = ' ';
            return;
        }
    }

    /* Check if line is being manually terminated. */
    if (c == '\n') {
        /* Write current line to output. */
        log_line[log_position]     = '\n';
        log_line[log_position + 1] = '\0';
        if (log_to_file) {
            fputs(log_line, log_file);
        }
        if (log_stdout) {
            fputs(log_line, stdout);
        }
        if (log_summary) {
            fputs(log_line, summary_file);
        }

        /* Set up for next line. */
        log_position      = 0;
        log_indentation   = 0;
        log_start_of_line = true;

        return;
    }

    /* Save this character. */
    log_line[log_position] = c;
    ++log_position;

    if (log_start_of_line && c == ' ') {  /* Determine number of indenting */
        ++log_indentation;                /*  spaces for current line. */
    }else {
        log_start_of_line = false;
    }
}


log_string(string)

char string[];

{
    int i, length;


    if (logging_disabled) {
        return;
    }

    length = strlen(string);
    for (i = 0; i < length; i++) {
        log_char(string[i]);
    }
}


log_int(value)

int value;

{
    char string[16];


    if (logging_disabled) {
        return;
    }

    sprintf(string, "%d\0", value);
    log_string(string);
}


log_long(value)

long value;

{
    char string[16];


    if (logging_disabled) {
        return;
    }

    sprintf(string, "%ld\0", value);
    log_string(string);
}



int num_locs = 0;

struct sp_loc_data loc[MAX_LOCATIONS];

get_location_data() {
    int locations_fd;

    long n, file_size;


    /* Open locations file. */
    locations_fd = open("locations.dat", 0);
    if (locations_fd < 0) {
        fprintf(stderr, "\nCannot open file 'locations.dat' for reading!\n\n");
        exit(-1);
    }

    /* Get size of file. */
    file_size = lseek(locations_fd, 0L, 2);
    num_locs  = file_size / sizeof(struct sp_loc_data);

    /* Read it all into memory. */
    lseek(locations_fd, 0L, 0);         /* Rewind first. */
    n = read(locations_fd, loc, file_size);
    if (n != file_size) {
        fprintf(stderr, "\nCannot read file 'locations.dat' into memory!\n\n");
        exit(-1);
    }

    close(locations_fd);
}



save_location_data() {
    int locations_fd;

    long n, num_bytes;


    /* Open file 'locations.dat' for writing. */
    locations_fd = creat("locations.dat", 0600);
    if (locations_fd < 0) {
        fprintf(stderr, "\n\tCannot create file 'locations.dat'!\n\n");
        exit(-1);
    }

    if (num_locs == 0) {
        close(locations_fd);
        return;
    }

    /* Write array to disk. */
    num_bytes = (long)num_locs * (long)sizeof(struct sp_loc_data);

    n = write(locations_fd, loc, num_bytes);
    if (n != num_bytes) {
        fprintf(stderr, "\n\n\tCannot write to 'locations.dat'!\n\n");
        exit(-1);
    }

    close(locations_fd);
}



/* The following routine provides the 'distorted' species number used to
 *      identify a species that uses field distortion units. The input
 *      variable 'species_number' is the same number used in filename
 *      creation for the species. */

int distorted(species_number)

int species_number;

{
    int i, j, n, ls;


    /* We must use the LS tech level at the start of the turn because
     * the distorted species number must be the same throughout the
     * turn, even if the tech level changes during production. */

    ls = spec_data[species_number - 1].init_tech_level[LS];

    i = species_number & 0x000F;                /* Lower four bits. */
    j = (species_number >> 4) & 0x000F;         /* Upper four bits. */

    n = (ls % 5 + 3) * (4 * i + j) + (ls % 11 + 7);

    return(n);
}

int undistorted(distorted_species_number)

int distorted_species_number;

{
    int i, species_number;


    for (i = 0; i < MAX_SPECIES; i++) {
        species_number = i + 1;

        if (distorted(species_number) == distorted_species_number) {
            return(species_number);
        }
    }

    return(0);   /* Not a legitimate species. */
}



log_message(message_filename)

char *message_filename;

{
    char message_line[256];

    FILE *message_file;


    /* Open message file. */
    message_file = fopen(message_filename, "r");
    if (message_file == NULL) {
        fprintf(stderr, "\n\tWARNING! utils.c: cannot open message file '%s'!\n\n", message_filename);
        return;
    }

    /* Copy message to log file. */
    for (;fgets(message_line, 256, message_file) != NULL;) {
        fputs(message_line, log_file);
    }

    fclose(message_file);
}



/* This routine will set or clear the POPULATED bit for a nampla.  It will
 * return true if the nampla is populated or false if not. It will also
 * check if a message associated with this planet should be logged. */

int check_population(nampla)

struct nampla_data *nampla;

{
    int is_now_populated, was_already_populated;

    long total_pop;

    char filename[32];


    if (nampla.status & POPULATED) {
        was_already_populated = true;
    }else{
        was_already_populated = false;
    }

    total_pop = nampla.mi_base
                + nampla.ma_base
                + nampla.IUs_to_install
                + nampla.AUs_to_install
                + nampla.item_quantity[PD]
                + nampla.item_quantity[CU]
                + nampla.pop_units;

    if (total_pop > 0) {
        nampla.status  |= POPULATED;
        is_now_populated = true;
    }else {
        nampla.status &= ~(POPULATED | MINING_COLONY
                            | RESORT_COLONY);
        is_now_populated = false;
    }

    if (is_now_populated && !was_already_populated) {
        if (nampla.message) {
            /* There is a message that must be logged whenever this planet
             *      becomes populated for the first time. */
            sprintf(filename, "message%ld.txt\0", nampla.message);
            log_message(filename);
        }
    }

    return(is_now_populated);
}

/* Get life support tech level needed. */

int life_support_needed(species, home, colony)

struct species_data *species;
struct planet_data * home, *colony;

{
    int i, j, k, ls_needed;


    i = colony.temperature_class - home.temperature_class;
    if (i < 0) {
        i = -i;
    }
    ls_needed = 3 * i;          /* Temperature class. */

    i = colony.pressure_class - home.pressure_class;
    if (i < 0) {
        i = -i;
    }
    ls_needed += 3 * i;         /* Pressure class. */

    /* Check gases. Assume required gas is NOT present. */
    ls_needed += 3;
    for (j = 0; j < 4; j++) {   /* Check gases on planet. */
        if (colony.gas_percent[j] == 0) {
            continue;
        }
        for (i = 0; i < 6; i++) { /* Compare with poisonous gases. */
            if (species.poison_gas[i] == colony.gas[j]) {
                ls_needed += 3;
            }
        }
        if (colony.gas[j] == species.required_gas) {
            if (colony.gas_percent[j] >= species.required_gas_min &&
                colony.gas_percent[j] <= species.required_gas_max) {
                ls_needed -= 3;
            }
        }
    }

    return(ls_needed);
}



check_high_tech_items(tech, old_tech_level, new_tech_level)

int tech, old_tech_level, new_tech_level;

{
    int i;


    for (i = 0; i < MAX_ITEMS; i++) {
        if (item_critical_tech[i] != tech) {
            continue;
        }
        if (new_tech_level < item_tech_requirment[i]) {
            continue;
        }
        if (old_tech_level >= item_tech_requirment[i]) {
            continue;
        }

        log_string("  You now have the technology to build ");
        log_string(item_name[i]);
        log_string("s.\n");
    }

    /* Check for high tech abilities that are not associated with specific
     *  items. */
    if (tech == MA && old_tech_level < 25 && new_tech_level >= 25) {
        log_string("  You now have the technology to do interspecies construction.\n");
    }
}



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

int agrep_score(correct_string, unknown_string)

char *correct_string, *unknown_string;

{
    int score;

    char c1, c2, *p1, *p2;


    if (strcmp(correct_string, unknown_string) == 0) {
        return(10000);
    }

    score = 0;
    p1    = correct_string;
    p2    = unknown_string;

    for {
        if ((c1 = *p1++) == '\0') {
            score -= strlen(p2);        /* Reduce score by excess characters,
                                         * if any. */
            break;
        }

        if ((c2 = *p2++) == '\0') {
            score -= strlen(p1);        /* Reduce score by excess characters,
                                         * if any. */
            break;
        }

        if (c1 == c2) {
            ++score;
        }else if (c1 == *p2 && c2 == *p1) {
            /* Transposed. */
            score += 2;
            ++p1;
            ++p2;
        }else if (c1 == *p2) {
            /* Unneeded character. */
            ++score;
            ++p2;
        }else if (c2 == *p1) {
            /* Missing character. */
            ++score;
            ++p1;
        }
    }

    return(score);
}


extern int num_stars, species_number, star_data_modified;
extern struct star_data *star_base;


/* The following routine will check if coordinates x-y-z contain a star and,
 * if so, will set the appropriate bit in the "visited_by" variable for the
 * star. If the star exists, true will be returned; otherwise, false will
 * be returned. */

int star_visited(x, y, z)

int x, y, z;

{
    int i, found, species_array_index, species_bit_number;

    long species_bit_mask;

    struct star_data *star;


    /* Get array index and bit mask. */
    species_array_index = (species_number - 1) / 32;
    species_bit_number  = (species_number - 1) % 32;
    species_bit_mask    = 1 << species_bit_number;

    found = false;

    for (i = 0; i < num_stars; i++) {
        star = star_base + i;

        if (x != star.x) {
            continue;
        }
        if (y != star.y) {
            continue;
        }
        if (z != star.z) {
            continue;
        }

        found = true;

        /* Check if bit is already set. */
        if (star.visited_by[species_array_index] & species_bit_mask) {
            break;
        }

        /* Set the appropriate bit. */
        star.visited_by[species_array_index] |= species_bit_mask;
        star_data_modified = true;
        break;
    }

    return(found);
}



//*************************************************************************
// with_check.c


extern int   truncate_name, ignore_field_distorters;
extern char  field_distorted[MAX_SPECIES];
extern FILE *log_file, *summary_file;


withdrawal_check(bat, act)

struct battle_data *bat;
struct action_data *act;

{
    /* This routine will check all fighting ships and see if any wish to
     *  withdraw. If so, it will set the ship's status to JUMPED_IN_COMBAT.
     *  The actual jump will be handled by the Jump program. */

    int i, old_trunc, ship_index, species_index, percent_loss,
        num_ships_gone[MAX_SPECIES], num_ships_total[MAX_SPECIES];

    char withdraw_age;

    struct ship_data *sh;


    for (i = 0; i < MAX_SPECIES; i++) {
        num_ships_gone[i]  = 0;
        num_ships_total[i] = 0;
    }

    old_trunc     = truncate_name; /* Show age of ship here. */
    truncate_name = false;

    /* Compile statistics and handle individual ships that must leave. */
    for (ship_index = 0; ship_index < act.num_units_fighting; ship_index++) {
        if (act.unit_ttype[ship_index] != SHIP) {
            continue;
        }

        sh            = (struct ship_data *)act.fighting_unit[ship_index];
        species_index = act.fighting_species_index[ship_index];
        ++num_ships_total[species_index];

        if (sh.status == JUMPED_IN_COMBAT) {   /* Already withdrawn. */
            ++num_ships_gone[species_index];
            continue;
        }

        if (sh.status == FORCED_JUMP) {        /* Forced to leave. */
            ++num_ships_gone[species_index];
            continue;
        }

        if (sh.age > 49) {                     /* Already destroyed. */
            ++num_ships_gone[species_index];
            continue;
        }

        if (sh.ttype != FTL) {
            continue;                   /* Ship can't jump. */
        }
        if (sh.class == TR) {
            withdraw_age = bat.transport_withdraw_age[species_index];
            if (withdraw_age == 0) {
                continue;
            }
            /* Transports will withdraw only when entire fleet withdraws. */
        }else {
            withdraw_age = bat.warship_withdraw_age[species_index];
        }

        if (sh.age > withdraw_age) {
            act.num_shots[ship_index]  = 0;
            act.shots_left[ship_index] = 0;
            sh.pn = 0;

            ignore_field_distorters = !field_distorted[species_index];

            fprintf(log_file, "        %s jumps away from the battle.\n",
                    ship_name(sh));
            fprintf(summary_file, "        %s jumps away from the battle.\n",
                    ship_name(sh));

            ignore_field_distorters = false;

            sh.dest_x = bat.haven_x[species_index];
            sh.dest_y = bat.haven_y[species_index];
            sh.dest_z = bat.haven_z[species_index];

            sh.status = JUMPED_IN_COMBAT;

            ++num_ships_gone[species_index];
        }
    }

    /* Now check if a fleet has reached its limit. */
    for (ship_index = 0; ship_index < act.num_units_fighting; ship_index++) {
        if (act.unit_ttype[ship_index] != SHIP) {
            continue;
        }

        sh            = (struct ship_data *)act.fighting_unit[ship_index];
        species_index = act.fighting_species_index[ship_index];

        if (sh.ttype != FTL) {
            continue;                                   /* Ship can't jump. */
        }
        if (sh.status == JUMPED_IN_COMBAT) {
            continue;                                   /* Already withdrawn. */
        }
        if (sh.status == FORCED_JUMP) {
            continue;                                   /* Already gone. */
        }
        if (sh.age > 49) {
            continue;                                   /* Already destroyed. */
        }
        if (bat.fleet_withdraw_percentage[species_index] == 0) {
            percent_loss = 101;         /* Always withdraw immediately. */
        }else {
            percent_loss = (100 * num_ships_gone[species_index])
                           / num_ships_total[species_index];
        }

        if (percent_loss > bat.fleet_withdraw_percentage[species_index]) {
            act.num_shots[ship_index]  = 0;
            act.shots_left[ship_index] = 0;
            sh.pn = 0;

            ignore_field_distorters = !field_distorted[species_index];

            fprintf(log_file, "        %s jumps away from the battle.\n",
                    ship_name(sh));
            fprintf(summary_file, "        %s jumps away from the battle.\n",
                    ship_name(sh));

            ignore_field_distorters = false;

            sh.dest_x = bat.haven_x[species_index];
            sh.dest_y = bat.haven_y[species_index];
            sh.dest_z = bat.haven_z[species_index];

            sh.status = JUMPED_IN_COMBAT;
        }
    }

    truncate_name = old_trunc;
}
