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
    "io"
)

//*************************************************************************
// TurnNumber Locations NoOrders Combat PreDeparture Jump Production PostArrival Strike Finish Report Stats



//*************************************************************************
// Finish.c
/*
 * This program should be run immediately before running the Report program;
 * i.e. immediately after the last run of AddSpecies in the very first turn, or
 * immediately after running PostArrival on all subsequent turns.  This program
 * will create the file 'locations.dat' (via the do_locations subroutine),
 * update populations, handle interspecies transactions, and do some other
 * housekeeping chores.
 */



func FinishMain(argc int, argv []string) {
    var i, j, n, rec, don, nampla_index, ship_index, ls_needed int
    var ls_actual, tech, turn_number, percent_increase, old_tech_level int
    var new_tech_level, experience_points, their_level, my_level int
    var new_level, orders_received, contact_bit_number int
    var contact_word_number, alien_number, galaxy_fd int
    var production_penalty, max_tech_level int
    var ns int
    var change, total_pop_units, contact_mask, salvage_EUs int
    var salvage_value, original_cost, ib, ab, increment, old_base int
    var max_cost, actual_cost, one_point_cost, working_pop_units int
    var ib_increment, ab_increment, md, growth_factor, denom int
    var fleet_maintenance_cost, balance, total_species_production int
    var RMs_produced, production_capacity, diff, total, eb int
    var filename[32]byte
    var dest, src *byte
    var home_planet *planet_data
    var donor_species *species_data
    var home_nampla *nampla_data

    /* Check for options, if any. */
    test_mode    = false;
    verbose_mode = false;
    for i = 1; i < argc; i++ {
        if (strcmp(argv[i], "-t") == 0) {
            test_mode = true;
        }
        if (strcmp(argv[i], "-v") == 0) {
            verbose_mode = true;
        }
    }

    /* Seed random number generator. */
    last_random = time(nil);
    n           = 907;
    for i = 0; i < n; i++ {
        rnd(100);
    }

    /* Get commonly used data. */
    get_galaxy_data();
    get_planet_data();
    get_species_data();
    get_transaction_data();
    num_locs = 0;

    /* Allocate memory for array "total_econ_base". */
    total           = num_planets * sizeof("something?");
    total_econ_base := make([]int, num_planets, num_planets) // warning: was (long *)malloc(total);
    if (total_econ_base == nil) {
        fprintf(stderr, "\nCannot allocate enough memory for total_econ_base!\n\n");
        exit(-1);
    }

    /* Handle turn number. */
    galaxy.turn_number++
    turn_number = galaxy.turn_number;
    galaxy_fd   = creat("galaxy.dat", 0600);
    if (galaxy_fd < 0) {
        fprintf(stderr, "\n  Cannot create new version of file galaxy.dat!\n");
        exit(-1);
    }

    n = write(galaxy_fd, &galaxy, sizeof("struct galaxy_data"));
    if (n != sizeof("struct galaxy_data")) {
        fprintf(stderr, "\n\tCannot write data to file 'galaxy.dat'!\n\n");
        exit(-1);
    }
    close(galaxy_fd);

    /* Do mining difficulty increases and initialize total economic base for each planet. */
    for i = 0; i < num_planets; i++ {
        planet = planet_base[i]
        planet.mining_difficulty += planet.md_increase;
        planet.md_increase        = 0;
        total_econ_base[i] = 0;
    }

    /* Main loop. For each species, take appropriate action. */
    if (verbose_mode) {
        printf("\nFinishing up for all species...\n");
    }
    for species_number = 1; species_number <= galaxy.num_species; species_number++ {
        if (!data_in_memory[species_number - 1]) {
            continue;
        }

        data_modified[species_number - 1] = true;

        species     = spec_data[species_number - 1];
        nampla_base = namp_data[species_number - 1];
        ship_base   = ship_data[species_number - 1];

        /* Check if player submitted orders for this turn. */
        filename := fmt.Sprintf( "sp%02d.ord", species_number);
        i = open(filename, 0);
        if (i < 0) {
            orders_received = false;
        }else{
            orders_received = true;
            close(i);
        }
        if (turn_number == 1) {
            orders_received = true;
        }

        /* Display name of species. */
        if (verbose_mode) {
            printf("  Now doing SP %s...", species.name);
            if (!orders_received) {
                printf(" WARNING: player did not submit orders this turn!");
            }
            printf("\n");
        }

        /* Open log file for appending. */
        filename := fmt.Sprintf("sp%02d.log", species_number);
        log_file = fopen(filename, "a");
        if (log_file == nil) {
            fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
            exit(-1);
        }
        log_stdout     = false;
        header_printed = false;

        if (turn_number == 1) {
            goto check_for_message;
        }

        /* Check if any ships of this species experienced mishaps. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == SHIP_MISHAP &&
                transaction[i].number1 == species_number) {
                if (!header_printed) {
                    print_header();
                }
                log_string("  !!! ");
                log_string(transaction[i].name1);
                if (transaction[i].value < 3) {
                    /* Intercepted or self-destructed. */
                    log_string(" disappeared without a trace, cause unknown!\n");
                }else if (transaction[i].value == 3) {
                    /* Mis-jumped. */
                    log_string(" mis-jumped to ");
                    log_int(transaction[i].x);  log_char(' ');
                    log_int(transaction[i].y);  log_char(' ');
                    log_int(transaction[i].z);  log_string("!\n");
                }else {
                    /* One fail-safe jump unit used. */
                    log_string(" had a jump mishap! A fail-safe jump unit was expended.\n");
                }
            }
        }

        /* Take care of any disbanded colonies. */
        home_nampla = nampla_base[0];
        for nampla_index = 0; nampla_index < species.num_namplas; nampla_index++ {
            nampla      = nampla_base[nampla_index]

            if ((nampla.status & DISBANDED_COLONY) == 0) {
                continue;
            }

            /* Salvage ships on the surface and starbases in orbit. */
            salvage_EUs = 0;
            for ship_index = 0; ship_index < species.num_ships; ship_index++ {
                ship        = ship_base[ship_index]
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
                if (ship.ttype != STARBASE && ship.status == IN_ORBIT) {
                    continue;
                }

                /* Transfer cargo to planet. */
                for i = 0; i < MAX_ITEMS; i++ {
                    nampla.item_quantity[i] += ship.item_quantity[i];
                }

                /* Salvage the ship. */
                if (ship.class == TR || ship.ttype == STARBASE) {
                    original_cost = ship_cost[ship.class] * ship.tonnage;
                }else{
                    original_cost = ship_cost[ship.class];
                }

                if (ship.ttype == SUB_LIGHT) {
                    original_cost = (3 * original_cost) / 4;
                }

                if (ship.status == UNDER_CONSTRUCTION) {
                    salvage_value =
                        (original_cost - ship.remaining_cost) / 4;
                }else{
                    salvage_value =
                        (3 * original_cost * (60 - ship.age)) / 400;
                }

                salvage_EUs += salvage_value;

                /* Destroy the ship. */
                delete_ship(ship);
            }

            /* Salvage items on the planet. */
            for i = 0; i < MAX_ITEMS; i++ {
                if (i == RM) {
                    salvage_value = nampla.item_quantity[RM] / 10;
                }else if (nampla.item_quantity[i] > 0) {
                    original_cost = nampla.item_quantity[i] * item_cost[i];
                    if (i == TP) {
                        if (species.tech_level[BI] > 0) {
                            original_cost /= species.tech_level[BI];
                        }else{
                            original_cost /= 100;
                        }
                    }
                    salvage_value = original_cost / 4;
                }else {
                    salvage_value = 0;
                }

                salvage_EUs += salvage_value;
            }

            /* Transfer EUs to species. */
            species.econ_units += salvage_EUs;

            /* Log what happened. */
            if (!header_printed) {
                print_header();
            }
            log_string("  PL ");
            log_string(nampla.name);
            log_string(" was disbanded, generating ");
            log_long(salvage_EUs);
            log_string(" economic units in salvage.\n");

            /* Destroy the colony. */
            delete_nampla(nampla);
        }

        /* Check if this species is the recipient of a transfer of economic
         *      units from another species. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].recipient == species_number && (transaction[i].ttype == EU_TRANSFER || transaction[i].ttype == SIEGE_EU_TRANSFER || transaction[i].ttype == LOOTING_EU_TRANSFER)) {
                /* Transfer EUs to attacker if this is a siege or looting
                 * transfer. If this is a normal transfer, then just log
                 * the result since the actual transfer was done when the
                 * order was processed. */
                if (transaction[i].ttype != EU_TRANSFER) {
                    species.econ_units += transaction[i].value;
                }

                if (!header_printed) {
                    print_header();
                }
                log_string("  ");
                log_long(transaction[i].value);
                log_string(" economic units were received from SP ");
                log_string(transaction[i].name1);
                if (transaction[i].ttype == SIEGE_EU_TRANSFER) {
                    log_string(" as a result of your successful siege of their PL ");
                    log_string(transaction[i].name3);
                    log_string(". The siege was ");
                    log_long(transaction[i].number1);
                    log_string("% effective");
                }else if (transaction[i].ttype == LOOTING_EU_TRANSFER) {
                    log_string(" as a result of your looting their PL ");
                    log_string(transaction[i].name3);
                }
                log_string(".\n");
            }
        }

        /* Check if any jump portals of this species were used by aliens. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == ALIEN_JUMP_PORTAL_USAGE &&
                transaction[i].number1 == species_number) {
                if (!header_printed) {
                    print_header();
                }
                log_string("  ");
                log_string(transaction[i].name1);
                log_char(' ');
                log_string(transaction[i].name2);
                log_string(" used jump portal ");
                log_string(transaction[i].name3);
                log_string(".\n");
            }
        }

        /* Check if any starbases of this species detected the use of gravitic
         *      telescopes by aliens. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == TELESCOPE_DETECTION &&
                transaction[i].number1 == species_number) {
                if (!header_printed) {
                    print_header();
                }
                log_string("! ");
                log_string(transaction[i].name1);
                log_string(" detected the operation of an alien gravitic telescope at x = ");
                log_int(transaction[i].x);
                log_string(", y = ");
                log_int(transaction[i].y);
                log_string(", z = ");
                log_int(transaction[i].z);
                log_string(".\n");
            }
        }

        /* Check if this species is the recipient of a tech transfer from another species. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == TECH_TRANSFER &&
                transaction[i].recipient == species_number) {
                rec = transaction[i].recipient - 1;
                don = transaction[i].donor - 1;

                /* Try to transfer technology. */
                if (!header_printed) {
                    print_header();
                }
                log_string("  ");
                tech = transaction[i].value;
                log_string(tech_name[tech]);
                log_string(" tech transfer from SP ");
                log_string(transaction[i].name1);
                their_level = transaction[i].number3;
                my_level    = species.tech_level[tech];

                if (their_level <= my_level) {
                    log_string(" failed.\n");
                    transaction[i].number1 = -1;
                    continue;
                }

                new_level     = my_level;
                max_cost      = transaction[i].number1;
                donor_species = spec_data[don];
                if (max_cost == 0) {
                    max_cost = donor_species.econ_units;
                }else if (donor_species.econ_units < max_cost) {
                    max_cost = donor_species.econ_units;
                }
                actual_cost = 0;
                for new_level < their_level {
                    one_point_cost  = new_level * new_level;
                    one_point_cost -= one_point_cost / 4;  /* 25% discount. */
                    if ((actual_cost + one_point_cost) > max_cost) {
                        break;
                    }
                    actual_cost += one_point_cost;
                    new_level++
                }

                if (new_level == my_level) {
                    log_string(" failed due to lack of funding.\n");
                    transaction[i].number1 = -2;
                }else {
                    log_string(" raised your tech level from ");
                    log_int(my_level);
                    log_string(" to ");
                    log_int(new_level);
                    log_string(" at a cost to them of ");
                    log_long(actual_cost);
                    log_string(".\n");
                    transaction[i].number1 = actual_cost;
                    transaction[i].number2 = my_level;
                    transaction[i].number3 = new_level;

                    species.tech_level[tech]  = new_level;
                    donor_species.econ_units -= actual_cost;
                }
            }
        }

        /* Calculate tech level increases. */
        for tech = MI; tech <= BI; tech++ {
            old_tech_level = species.tech_level[tech];
            new_tech_level = old_tech_level;

            experience_points = species.tech_eps[tech];
            if (experience_points == 0) {
                goto check_random;
            }

            /* Determine increase as if there were NO randomness in the
             *  process. */
            i = experience_points;
            j = old_tech_level;
            for i >= j * j {
                i -= j * j;
                j++
            }

            /* When extremely large amounts are spent on research, tech
             *  level increases are sometimes excessive.  Set a limit. */
            if (old_tech_level > 50) {
                max_tech_level = j + 1;
            }else{
                max_tech_level = 9999;
            }

            /* Allocate half of the calculated increase NON-RANDOMLY. */
            n = (j - old_tech_level) / 2;
            for i = 0; i < n; i++ {
                experience_points -= new_tech_level * new_tech_level;
                new_tech_level++
            }

            /* Allocate the rest randomly. */
            for experience_points >= new_tech_level {
                experience_points -= new_tech_level;
                n = new_tech_level;

                /* The chance of success is 1 in n. At this point, n is always at least 1. */
                i = rnd(16 * n);
                if (i >= 8 * n && i <= 8 * n + 15) {
                    new_tech_level = n + 1;
                }
            }

            /* Save unused experience points. */
            species.tech_eps[tech] = experience_points;

check_random:

            /* See if any random increase occurred. Odds are 1 in 6. */
            if (old_tech_level > 0 && rnd(6) == 6) {
                new_tech_level++
            }

            if (new_tech_level > max_tech_level) {
                new_tech_level = max_tech_level;
            }

            /* Report result only if tech level went up. */
            if (new_tech_level > old_tech_level) {
                if (!header_printed) {
                    print_header();
                }
                log_string("  ");
                log_string(tech_name[tech]);
                log_string(" tech level rose from ");
                log_int(old_tech_level);
                log_string(" to ");
                log_int(new_tech_level);
                log_string(".\n");

                species.tech_level[tech] = new_tech_level;
            }
        }

        /* Notify of any new high tech items. */
        for tech = MI; tech <= BI; tech++ {
            old_tech_level = species.init_tech_level[tech];
            new_tech_level = species.tech_level[tech];

            if (new_tech_level > old_tech_level) {
                check_high_tech_items(tech, old_tech_level, new_tech_level);
            }

            species.init_tech_level[tech] = new_tech_level;
        }

        /* Check if this species is the recipient of a knowledge transfer
         *      from another species. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == KNOWLEDGE_TRANSFER && transaction[i].recipient == species_number) {
                rec = transaction[i].recipient - 1;
                don = transaction[i].donor - 1;

                /* Try to transfer technology. */
                tech        = transaction[i].value;
                their_level = transaction[i].number3;
                my_level    = species.tech_level[tech];
                n           = species.tech_knowledge[tech];
                if (n > my_level) {
                    my_level = n;
                }

                if (their_level <= my_level) {
                    continue;
                }

                species.tech_knowledge[tech] = their_level;

                if (!header_printed) {
                    print_header();
                }
                log_string("  SP ");
                log_string(transaction[i].name1);
                log_string(" transferred knowledge of ");
                log_string(tech_name[tech]);
                log_string(" to you up to tech level ");
                log_long(their_level);
                log_string(".\n");
            }
        }

        /* Loop through each nampla for this species. */
        home_nampla = nampla_base[0];
        home_planet = planet_base[home_nampla.planet_index]
        for nampla_index = 0; nampla_index < species.num_namplas; nampla_index++ {
            nampla      = nampla_base[nampla_index]

            if (nampla.pn == 99) {
                continue;
            }

            /* Get planet pointer. */
            planet = planet_base[nampla.planet_index]

            /* Clear any amount spent on ambush. */
            nampla.use_on_ambush = 0;

            /* Handle HIDE order. */
            nampla.hidden = nampla.hiding;
            nampla.hiding = false;

            /* Check if any IUs or AUs were installed. */
            if (nampla.IUs_to_install > 0) {
                nampla.mi_base       += nampla.IUs_to_install;
                nampla.IUs_to_install = 0;
            }

            if (nampla.AUs_to_install > 0) {
                nampla.ma_base       += nampla.AUs_to_install;
                nampla.AUs_to_install = 0;
            }

            /* Check if another species on the same planet has become
             *  assimilated. */
            for i = 0; i < num_transactions; i++ {
                if (transaction[i].ttype == ASSIMILATION && transaction[i].value == species_number && transaction[i].x == nampla.x && transaction[i].y == nampla.y && transaction[i].z == nampla.z && transaction[i].pn == nampla.pn) {
                    ib = transaction[i].number1;
                    ab = transaction[i].number2;
                    ns = transaction[i].number3;
                    nampla.mi_base   += ib;
                    nampla.ma_base   += ab;
                    nampla.shipyards += ns;

                    if (!header_printed) {
                        print_header();
                    }

                    log_string("  Assimilation of ");
                    log_string(transaction[i].name1);
                    log_string(" PL ");
                    log_string(transaction[i].name2);
                    log_string(" increased mining base of ");
                    log_string(species.name);
                    log_string(" PL ");
                    log_string(nampla.name);
                    log_string(" by ");
                    log_long(ib / 10);
                    log_char('.');
                    log_long(ib % 10);
                    log_string(", and manufacturing base by ");
                    log_long(ab / 10);
                    log_char('.');
                    log_long(ab % 10);
                    if (ns > 0) {
                        log_string(". Number of shipyards was also increased by ");
                        log_int(ns);
                    }
                    log_string(".\n");
                }
            }

            /* Calculate available population for this turn. */
            nampla.pop_units = 0;

            eb = nampla.mi_base + nampla.ma_base;
            total_pop_units = eb + nampla.item_quantity[CU] + nampla.item_quantity[PD];

            if (nampla.status & HOME_PLANET) != 0 {
                if (nampla.status & POPULATED) != 0 {
                    nampla.pop_units = HP_AVAILABLE_POP;
                    if (species.hp_original_base != 0) {  /* HP was bombed. */
                        if (eb >= species.hp_original_base) {
                            species.hp_original_base = 0; /* Fully recovered. */
                        }else {
                            nampla.pop_units = (eb * HP_AVAILABLE_POP) / species.hp_original_base;
                        }
                    }
                }
            }else if (nampla.status & POPULATED) != 0 {
                /* Get life support tech level needed. */
                ls_needed = life_support_needed(species, home_planet, planet);

                /* Basic percent increase is 10*(1 - ls_needed/ls_actual). */
                ls_actual        = species.tech_level[LS];
                percent_increase =
                    10 * (100 - ((100 * ls_needed) / ls_actual));

                if (percent_increase < 0) {     /* Colony wiped out! */
                    if (!header_printed) {
                        print_header();
                    }

                    log_string("  !!! Life support tech level was too low to support colony on PL ");
                    log_string(nampla.name);
                    log_string(". Colony was destroyed.\n");

                    nampla.status            = COLONY; /* No longer populated or self-sufficient. */
                    nampla.mi_base           = 0;
                    nampla.ma_base           = 0;
                    nampla.pop_units         = 0;
                    nampla.item_quantity[PD] = 0;
                    nampla.item_quantity[CU] = 0;
                    nampla.siege_eff         = 0;
                }else {
                    percent_increase /= 100;

                    /* Add a small random variation. */
                    percent_increase += rnd(percent_increase / 4) - rnd(percent_increase / 4);

                    /* Add bonus for Biology technology. */
                    percent_increase += species.tech_level[BI] / 20;

                    /* Calculate and apply the change. */
                    change = (percent_increase * total_pop_units) / 100;

                    if (nampla.mi_base > 0 && nampla.ma_base == 0) {
                        nampla.status |= MINING_COLONY;
                        change          = 0;
                    }else if (nampla.status & MINING_COLONY) != 0 {
                        /* A former mining colony has been converted to a normal colony. */
                        nampla.status &= ^MINING_COLONY;
                        change          = 0;
                    }

                    if (nampla.ma_base > 0 && nampla.mi_base == 0 && ls_needed <= 6 && planet.gravity <= home_planet.gravity) {
                        nampla.status |= RESORT_COLONY;
                        change          = 0;
                    }else if (nampla.status & RESORT_COLONY) != 0 {
                        /* A former resort colony has been converted to a normal colony. */
                        nampla.status &= ^RESORT_COLONY;
                        change          = 0;
                    }

                    if (total_pop_units == nampla.item_quantity[PD]) {
                        change = 0;     /* Probably an invasion force. */
                    }
                    nampla.pop_units = change;
                }
            }

            /* Handle losses due to attrition and update location array if
             *  planet is still populated. */
            if (nampla.status & POPULATED) != 0 {
                total_pop_units = nampla.pop_units + nampla.mi_base + nampla.ma_base + nampla.item_quantity[CU] + nampla.item_quantity[PD];

                if (total_pop_units > 0 && total_pop_units < 50) {
                    if (nampla.pop_units > 0) {
                        nampla.pop_units--
                        goto do_auto_increases;
                    }else if (nampla.item_quantity[CU] > 0) {
                        nampla.item_quantity[CU]--
                        if (!header_printed) {
                            print_header();
                        }
                        log_string("  Number of colonist units on PL ");
                        log_string(nampla.name);
                        log_string(" was reduced by one unit due to normal attrition.");
                    }else if (nampla.item_quantity[PD] > 0) {
                        nampla.item_quantity[PD]--
                        if (!header_printed) {
                            print_header();
                        }
                        log_string("  Number of planetary defense units on PL ");
                        log_string(nampla.name);
                        log_string(" was reduced by one unit due to normal attrition.");
                    }else if (nampla.ma_base > 0) {
                        nampla.ma_base--
                        if (!header_printed) {
                            print_header();
                        }
                        log_string("  Manufacturing base of PL ");
                        log_string(nampla.name);
                        log_string(" was reduced by 0.1 due to normal attrition.");
                    }else {
                        nampla.mi_base--
                        if (!header_printed) {
                            print_header();
                        }
                        log_string("  Mining base of PL ");
                        log_string(nampla.name);
                        log_string(" was reduced by 0.1 due to normal attrition.");
                    }

                    if (total_pop_units == 1) {
                        if (!header_printed) {
                            print_header();
                        }
                        log_string(" The colony is dead!");
                    }

                    log_char('\n');
                }
            }

do_auto_increases:

            /* Apply automatic 2% increase to mining and manufacturing bases of home planets. */
            if (nampla.status & HOME_PLANET) != 0 {
                growth_factor = 20;
                ib            = nampla.mi_base;
                ab            = nampla.ma_base;
                old_base      = ib + ab;
                increment     = (growth_factor * old_base) / 1000;
                md            = planet.mining_difficulty;

                denom        = 100 + md;
                ab_increment = (100 * (increment + ib) - (md * ab) + denom / 2) / denom;
                ib_increment = increment - ab_increment;

                if (ib_increment < 0) {
                    ab_increment = increment;
                    ib_increment = 0;
                }
                if (ab_increment < 0) {
                    ib_increment = increment;
                    ab_increment = 0;
                }
                nampla.mi_base += ib_increment;
                nampla.ma_base += ab_increment;
            }

check_pop:

            check_population(nampla);

            /* Update total economic base for colonies. */
            if ((nampla.status & HOME_PLANET) == 0) {
                total_econ_base[nampla.planet_index] += nampla.mi_base + nampla.ma_base;
            }
        }

        /* Loop through all ships for this species. */
        for ship_index = 0; ship_index < species.num_ships; ship_index++ {
            ship = ship_base[ship_index]

            if (ship.pn == 99) {
                continue;
            }

            /* Set flag if ship arrived via a natural wormhole. */
            if (ship.just_jumped == 99) {
                ship.arrived_via_wormhole = true;
            }else{
                ship.arrived_via_wormhole = false;
            }

            /* Clear 'just-jumped' flag. */
            ship.just_jumped = false;

            /* Increase age of ship. */
            if (ship.status != UNDER_CONSTRUCTION) {
                ship.age += 1;
                if (ship.age > 49) {
                    ship.age = 49;
                }
            }
        }

        /* Check if this species has a populated planet that another species
         *      tried to land on. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == LANDING_REQUEST && transaction[i].number1 == species_number) {
                if (!header_printed) {
                    print_header();
                }
                log_string("  ");
                log_string(transaction[i].name2);
                log_string(" owned by SP ");
                log_string(transaction[i].name3);
                if (transaction[i].value) {
                    log_string(" was granted");
                }else{
                    log_string(" was denied");
                }
                log_string(" permission to land on PL ");
                log_string(transaction[i].name1);
                log_string(".\n");
            }
        }

        /* Check if this species is the recipient of interspecies
         *      construction. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == INTERSPECIES_CONSTRUCTION && transaction[i].recipient == species_number) {
                /* Simply log the result. */
                if (!header_printed) {
                    print_header();
                }
                log_string("  ");
                if (transaction[i].value == 1) {
                    log_long(transaction[i].number1);
                    log_char(' ');
                    log_string(item_name[transaction[i].number2]);
                    if (transaction[i].number1 == 1) {
                        log_string(" was");
                    }else{
                        log_string("s were");
                    }
                    log_string(" constructed for you by SP ");
                    log_string(transaction[i].name1);
                    log_string(" on PL ");
                    log_string(transaction[i].name2);
                }else {
                    log_string(transaction[i].name2);
                    log_string(" was constructed for you by SP ");
                    log_string(transaction[i].name1);
                }
                log_string(".\n");
            }
        }

        /* Check if this species is besieging another species and detects forbidden construction, landings, etc. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == DETECTION_DURING_SIEGE && transaction[i].number3 == species_number) {
                /* Log what was detected and/or destroyed. */
                if (!header_printed) {
                    print_header();
                }
                log_string("  ");
                log_string("During the siege of ");
                log_string(transaction[i].name3);
                log_string(" PL ");
                log_string(transaction[i].name1);
                log_string(", your forces detected the ");

                if (transaction[i].value == 1) {
                    /* Landing of enemy ship. */
                    log_string("landing of ");
                    log_string(transaction[i].name2);
                    log_string(" on the planet.\n");
                }else if (transaction[i].value == 2) {
                    /* Enemy ship or starbase construction. */
                    log_string("construction of ");
                    log_string(transaction[i].name2);
                    log_string(", but you destroyed it before it");
                    log_string(" could be completed.\n");
                }else if (transaction[i].value == 3) {
                    /* Enemy PD construction. */
                    log_string("construction of planetary defenses, but you");
                    log_string(" destroyed them before they could be completed.\n");
                }else if (transaction[i].value == 4 || transaction[i].value == 5) {
                    /* Enemy item construction. */
                    log_string("transfer of ");
                    log_int(transaction[i].number1);
                    log_char(' ');
                    log_string(item_name[transaction[i].number2]);
                    if (transaction[i].number1 > 1) {
                        log_char('s');
                    }
                    if (transaction[i].value == 4) {
                        log_string(" to PL ");
                    }else{
                        log_string(" from PL ");
                    }
                    log_string(transaction[i].name2);
                    log_string(", but you destroyed them in transit.\n");
                }else {
                    fprintf(stderr, "\n\tInternal error!  Cannot reach this point!\n\n");
                    exit(-1);
                }
            }
        }

check_for_message:

        /* Check if this species is the recipient of a message from another
         *      species. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == MESSAGE_TO_SPECIES && transaction[i].number2 == species_number) {
                if (!header_printed) {
                    print_header();
                }
                log_string("\n  You received the following message from SP ");
                log_string(transaction[i].name1);
                log_string(":\n\n");

                filename := fmt.Sprintf("m%d.msg", transaction[i].value);

                log_message(filename);

                log_string("\n  *** End of Message ***\n\n");
            }
        }

        /* Close log file. */
        fclose(log_file);
    }

    /* Calculate economic efficiency for each planet. */
    for i = 0; i < num_planets; i++ {
        planet = planet_base[i]
        total = total_econ_base[i];
        diff  = total - 2000;

        if (diff <= 0) {
            planet.econ_efficiency = 100;
        }else{
            planet.econ_efficiency = (100 * (diff / 20 + 2000)) / total;
        }
    }

    /* Create new locations array. */
    do_locations();

    if (turn_number == 1) {
        goto clean_up;
    }

    /* Go through all species one more time to update alien contact masks,
     *  report tech transfer results to donors, and calculate fleet
     *  maintenance costs. */
    if (verbose_mode) {
        printf("\nNow updating contact masks et al.\n");
    }
    for species_index = 0; species_index < galaxy.num_species; species_index++ {
        if (!data_in_memory[species_index]) {
            continue;
        }

        species        = spec_data[species_index];
        nampla_base    = namp_data[species_index];
        ship_base      = ship_data[species_index];
        species_number = species_index + 1;

        home_nampla = nampla_base[0]
        home_planet = planet_base[home_nampla.planet_index]

        /* Update contact mask in species data if this species has met a new alien. */
        for i = 0; i < num_locs; i++ {
            if (loc[i].s != species_number) {
                continue;
            }

            for j = 0; j < num_locs; j++ {
                if (loc[j].s == species_number) {
                    continue;
                }
                if (loc[j].x != loc[i].x) {
                    continue;
                }
                if (loc[j].y != loc[i].y) {
                    continue;
                }
                if (loc[j].z != loc[i].z) {
                    continue;
                }

                /* We are in contact with an alien. Make sure it is not hidden from us. */
                alien_number = loc[j].s;
                if (alien_is_visible(loc[j].x, loc[j].y, loc[j].z, species_number, alien_number)) {
                    species.contact[loc[j].s] = true
                }
            }
        }

        /* Report results of tech transfers to donor species. */
        for i = 0; i < num_transactions; i++ {
            if (transaction[i].ttype == TECH_TRANSFER && transaction[i].donor == species_number) {
                /* Open log file for appending. */
                filename := fmt.Sprintf("sp%02d.log", species_number);
                log_file = fopen(filename, "a");
                if (log_file == nil) {
                    fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
                    exit(-1);
                }
                log_stdout = false;

                log_string("  ");
                tech = transaction[i].value;
                log_string(tech_name[tech]);
                log_string(" tech transfer to SP ");
                log_string(transaction[i].name2);

                if (transaction[i].number1 < 0) {
                    log_string(" failed");
                    if (transaction[i].number1 == -2) {
                        log_string(" due to lack of funding");
                    }
                }else {
                    log_string(" raised their tech level from ");
                    log_long(transaction[i].number2);
                    log_string(" to ");
                    log_long(transaction[i].number3);
                    log_string(" at a cost to you of ");
                    log_long(transaction[i].number1);
                }

                log_string(".\n");

                fclose(log_file);
            }
        }

        /* Calculate fleet maintenance cost and its percentage of total
         * production. */
        fleet_maintenance_cost = 0;
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base[i]

            if (ship.pn == 99) {
                continue;
            }

            if (ship.class == TR) {
                n = 4 * ship.tonnage;
            }else if (ship.class == BA) {
                n = 10 * ship.tonnage;
            }else{
                n = 20 * ship.tonnage;
            }

            if (ship.ttype == SUB_LIGHT) {
                n -= (25 * n) / 100;
            }

            fleet_maintenance_cost += n;
        }

        /* Subtract military discount. */
        i = species.tech_level[ML] / 2;
        fleet_maintenance_cost -= (i * fleet_maintenance_cost) / 100;

        /* Calculate total production. */
        total_species_production = 0;
        for i = 0; i < species.num_namplas; i++ {
            nampla = nampla_base[i]

            if (nampla.pn == 99) {
                continue;
            }
            if (nampla.status & DISBANDED_COLONY) != 0 {
                continue;
            }

            planet = planet_base[nampla.planet_index]

            ls_needed = life_support_needed(species, home_planet, planet);

            if (ls_needed == 0) {
                production_penalty = 0;
            }else{
                production_penalty = (100 * ls_needed) / species.tech_level[LS];
            }

            RMs_produced = (10 * species.tech_level[MI] * nampla.mi_base) / planet.mining_difficulty;
            RMs_produced -= (production_penalty * RMs_produced) / 100;

            production_capacity = (species.tech_level[MA] * nampla.ma_base) / 10;
            production_capacity -= (production_penalty * production_capacity) / 100;

            if (nampla.status & MINING_COLONY) != 0 {
                balance = (2 * RMs_produced) / 3;
            }else if (nampla.status & RESORT_COLONY) != 0 {
                balance = (2 * production_capacity) / 3;
            }else{
                RMs_produced += nampla.item_quantity[RM];
                if RMs_produced > production_capacity {
                    balance       = production_capacity
                } else {
                    balance       = RMs_produced;
                }
            }

            balance = ((planet.econ_efficiency * balance) + 50) / 100;

            total_species_production += balance;
        }

        /* If cost is greater than production, take as much as possible from EUs in treasury.
         * if (fleet_maintenance_cost > total_species_production)
         * {
         *  if (fleet_maintenance_cost > species.econ_units)
         *  {
         *      fleet_maintenance_cost -= species.econ_units;
         *      species.econ_units = 0;
         *  {
         *  else
         *  {
         *      species.econ_units -= fleet_maintenance_cost;
         *      fleet_maintenance_cost = 0;
         *  }
         * }
         */

        /* Save fleet maintenance results. */
        species.fleet_cost = fleet_maintenance_cost;
        if (total_species_production > 0) {
            species.fleet_percent_cost = (10000 * fleet_maintenance_cost) / total_species_production;
        }else{
            species.fleet_percent_cost = 10000;
        }
    }

clean_up:

    /* Clean up and exit. */
    save_planet_data();
    save_location_data();
    save_species_data();
    free_species_data();
    free(planet_base);
    free(total_econ_base);
    exit(0);
}

func print_header() {
    log_string("\nOther events:\n");
    header_printed = true;
}

func alien_is_visible(x, y, z, species_number, alien_number int) bool {
    var i, j int
    var species, alien *species_data
    var nampla, alien_nampla *nampla_data
    var alien_ship *ship_data

    /* Check if the alien has a ship or starbase here that is in orbit or in deep space. */
    alien      = spec_data[alien_number - 1];
    for i = 0; i < alien.num_ships; i++ {
        alien_ship = ship_data[alien_number - 1][i]
        if (alien_ship.x != x) {
            continue;
        }
        if (alien_ship.y != y) {
            continue;
        }
        if (alien_ship.z != z) {
            continue;
        }
        if (alien_ship.item_quantity[FD] == alien_ship.tonnage) {
            continue;
        }

        if (alien_ship.status == IN_ORBIT || alien_ship.status == IN_DEEP_SPACE) {
            return(true);
        }
    }

    /* Check if alien has a planet that is not hidden. */
    for i = 0; i < alien.num_namplas; i++ {
        alien_nampla = namp_data[alien_number - 1][i]
        if (alien_nampla.x != x) {
            continue;
        }
        if (alien_nampla.y != y) {
            continue;
        }
        if (alien_nampla.z != z) {
            continue;
        }
        if ((alien_nampla.status & POPULATED) == 0) {
            continue;
        }

        if (!alien_nampla.hidden) {
            return(true);
        }

        /* The colony is hidden. See if we have population on the same planet. */
        species = spec_data[species_number - 1];
        for j = 0; j < species.num_namplas; j++ {
            nampla  = namp_data[species_number - 1][j]
            if (nampla.x != x) {
                continue;
            }
            if (nampla.y != y) {
                continue;
            }
            if (nampla.z != z) {
                continue;
            }
            if (nampla.pn != alien_nampla.pn) {
                continue;
            }
            if ((nampla.status & POPULATED) == 0) {
                continue;
            }

            /* We have population on the same planet, so the alien cannot hide. */
            return(true);
        }
    }

    return(false);
}

//*************************************************************************
// Locations.c
/* This program will create the file locations.dat and will update the
 * economic efficiencies of all planets.  These functions are also performed
 * by Finish.c.  This program should be run before the strike phase or whenever
 * manual changes are made to the species data files that resulted in something
 * not being where it was or something being where it was not. It should also
 * be run if you run Finish on fewer than all species and decide to keep the
 * resulting planets.dat file. */

func LocationsMain(argc int, argv []string) {
    var i, nampla_index, diff, total int
    var total_econ_base []int

    /* Check for options, if any. */
    test_mode    = false;
    verbose_mode = false;
    for i = 1; i < argc; i++ {
        if (strcmp(argv[i], "-t") == 0) {
            test_mode = true;
        }
        if (strcmp(argv[i], "-v") == 0) {
            verbose_mode = true;
        }
    }

    /* Get commonly used data. */
    get_galaxy_data();
    get_planet_data();
    get_species_data();

    /* Allocate memory for array "total_econ_base". */
    total           = num_planets * sizeof("something");
    total_econ_base = make([]int, num_planets, num_planets) // warning: was (long *)malloc(total);
    if (total_econ_base == nil) {
        fprintf(stderr, "\nCannot allocate enough memory for total_econ_base!\n\n");
        exit(-1);
    }

    /* Initialize total econ base for each planet. */
    for i = 0; i < num_planets; i++ {
        planet = planet_base[i] // TODO: why is this here?
        total_econ_base[i] = 0;
    }

    /* Get total economic base for each planet from nampla data. */
    for species_number = 1; species_number <= galaxy.num_species; species_number++ {
        if (!data_in_memory[species_number - 1]) {
            continue;
        }

        data_modified[species_number - 1] = true;

        species     = spec_data[species_number - 1];
        nampla_base = namp_data[species_number - 1];

        for nampla_index = 0; nampla_index < species.num_namplas; nampla_index++ {
            nampla = nampla_base[nampla_index]
            if (nampla.pn == 99) {
                continue;
            }

            if ((nampla.status & HOME_PLANET) == 0) {
                total_econ_base[nampla.planet_index] +=
                    nampla.mi_base + nampla.ma_base;
            }
        }
    }

    /* Update economic efficiencies of all planets. */
    for i = 0; i < num_planets; i++ {
        planet = planet_base[i]
        total = total_econ_base[i];
        diff  = total - 2000;

        if (diff <= 0) {
            planet.econ_efficiency = 100;
        }else{
            planet.econ_efficiency = (100 * (diff / 20 + 2000)) / total;
        }
    }

    /* Create new locations array. */
    do_locations();

    /* Clean up and exit. */
    save_location_data();
    save_planet_data();
    free_species_data();
    free(planet_base);
    exit(0);
}

//*************************************************************************
// Jump.c

func JumpMain(argc int, argv []string) {
    var i, n, found, num_species int
    var sp_num[MAX_SPECIES]int
    var sp_index, command, log_file_open, do_all_species int
    var filename[32]byte
    var species_jumped[MAX_SPECIES]byte
    var keyword[4]byte

    /* Seed random number generator. */
    last_random = time(nil);
    n           = rnd(100) + rnd(200) + rnd(300);
    for i = 0; i < n; i++ {
        rnd(10);
    }

    /* Get commonly used data. */
    get_galaxy_data();
    get_transaction_data();

    ignore_field_distorters = true;

    /* Check arguments. If an argument is -p, then do two passes. In the
     *  first pass, display results and prompt the GM, allowing him to
     *  abort if necessary before saving results to disk. All other
     *  arguments must be species numbers. If no species numbers are
     *  specified, then do all species. */
    num_species  = 0;
    first_pass   = false;
    test_mode    = false;
    verbose_mode = false;
    for i = 1; i < argc; i++ {
        if (strcmp(argv[i], "-p") == 0) {
            first_pass = true;
        }else if (strcmp(argv[i], "-t") == 0) {
            test_mode = true;
        }else if (strcmp(argv[i], "-v") == 0) {
            verbose_mode = true;
        }else{
            n = atoi(argv[i]);
            if (n < 1 || n > galaxy.num_species) {
                fprintf(stderr, "\n    '%s' is not a valid argument!\n", argv[i]);
                exit(-1);
            }
            sp_num[num_species] = n;
            num_species++
        }
    }

    if (num_species == 0) {
        num_species = galaxy.num_species;
        for i = 0; i < num_species; i++ {
            sp_num[i] = i + 1;
        }
        do_all_species = true;
    }else {
        do_all_species = false;
    }

    /* For these commands, do not display age or landed/orbital status of ships. */
    truncate_name = true;
    log_stdout    = false; /* We will control value of log_file from here. */

    /* Initialize array that will indicate which species provided jump
     *  orders. If ships of a species withdrew or were forced from combat
     *  and there were no jump orders for that species, then combat jumps
     *  will not take place. This array will allow us to handle them
     *  separately. */
    for i = 0; i < galaxy.num_species; i++ {
        species_jumped[i] = false;
    }

    /* Two passes through all orders will be done. The first pass will
     *  check for errors and abort if any are found. Results will be written
     *  to disk only on the second pass. */

start_pass:

    if (first_pass) {
        printf("\nStarting first pass...\n\n");
    }

    get_species_data();
    get_star_data();
    get_planet_data();

    /* Main loop. For each species, take appropriate action. */
    for sp_index = 0; sp_index < num_species; sp_index++ {
        species_number = sp_num[sp_index];

        found = data_in_memory[species_number - 1];
        if (!found) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    Skipping species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n    Cannot get data for species #%d!\n",
                        species_number);
                exit(-1);
            }
        }

        species     = spec_data[species_number - 1];
        nampla_base = namp_data[species_number - 1];
        ship_base   = ship_data[species_number - 1];

        /* Open orders file for this species. */
        filename := fmt.Sprintf("sp%02d.ord", species_number);
        input_file = fopen(filename, "r");
        if (input_file == nil) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    No orders for species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n\tCannot open '%s' for reading!\n\n", filename);
                exit(-1);
            }
        }

        /* Open log file. Use stdout for first pass. */
        if (first_pass) {
            log_file = stdout;
        }else {
            /* Open log file for appending. */
            filename := fmt.Sprintf("sp%02d.log", species_number);
            log_file = fopen(filename, "a");
            if (log_file == nil) {
                fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
                exit(-1);
            }
        }

        end_of_file = false;

        just_opened_file = true;        /* Tell parse.c to skip mail header, if any. */
find_start:

        /* Search for START JUMPS order. */
        found = false;
        for !found {
            command = get_command();
            if (command == MESSAGE) {
                /* Skip MESSAGE text. It may contain a line that starts with "start". */
                for {
                    command = get_command();
                    if (command < 0) {
                        fprintf(stderr, "WARNING: Unterminated MESSAGE command in file %s!\n", filename);
                        break;
                    }

                    if (command == ZZZ) {
                        goto find_start;
                    }
                }
            }

            if (command < 0) {
                break;          /* End of file. */
            }
            if (command != START) {
                continue;
            }

            /* Get the first three letters of the keyword and convert to upper case. */
            skip_whitespace();
            for i = 0; i < 3; i++ {
                keyword[i] = toupper(*input_line_pointer);
                input_line_pointer++
            }
            keyword[3] = 0;

            if (strcmp(keyword, "JUM") == 0) {
                found = true;
            }
        }

        if (!found) {
            if (first_pass) {
                printf("\nNo jump orders for species #%d, SP %s.\n", species_number, species.name);
            }
            goto done_orders;
        }

        /* Handle jump orders for this species. */
        log_string("\nJump orders:\n");
        do_jump_orders();
        species_jumped[species_number - 1] = true;
        data_modified[species_number - 1]  = true;

done_orders:

        fclose(input_file);

        /* Take care of any ships that withdrew or were forced to jump during combat. */
        for ship_index = 0; ship_index < species.num_ships; ship_index++ {
            ship = ship_base[ship_index]
            if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
                do_JUMP_command(true, false);
            }
        }

        /* If this is the second pass, close the log file. */
        if (!first_pass) {
            fclose(log_file);
        }
    }

    if (first_pass) {
        printf("\nFinal chance to abort safely!\n");
        gamemaster_abort_option();
        first_pass = false;
        free_species_data();
        free(star_base);        /* In case data was modified. */
        free(planet_base);      /* In case data was modified. */

        printf("\nStarting second pass...\n\n");

        goto start_pass;
    }

no_jump_orders:

    /* Take care of any ships that withdrew from combat but were not
     *  handled above because no jump orders were received for species. */
    log_stdout    = true;
    log_file_open = false;
    for species_number = 1; species_number <= galaxy.num_species; species_number++ {
        if (species_jumped[species_number - 1]) {
            continue;
        }

        if (!data_in_memory[species_number - 1]) {
            continue;
        }

        species     = spec_data[species_number - 1];
        nampla_base = namp_data[species_number - 1];
        ship_base   = ship_data[species_number - 1];

        for ship_index = 0; ship_index < species.num_ships; ship_index++ {
            ship = ship_base[ship_index]
            if (ship.status == FORCED_JUMP || ship.status == JUMPED_IN_COMBAT) {
                if (!log_file_open) {
                    filename := fmt.Sprintf("sp%02d.log", species_number);
                    log_file = fopen(filename, "a");
                    if (log_file == nil) {
                        fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
                        exit(-1);
                    }
                    log_file_open = true;
                    log_string("\nWithdrawals and forced jumps during combat:\n");
                }

                do_JUMP_command(true, false);
            }
        }

        data_modified[species_number - 1] = log_file_open;

        if (log_file_open) {
            fclose(log_file);
            log_file_open = false;
        }
    }

    save_species_data();
    save_transaction_data();
    if (star_data_modified) {
        save_star_data();
    }
    if (planet_data_modified) {
        save_planet_data();
    }
    free_species_data();
    free(star_base);
    free(planet_base);
    exit(0);
}

func do_jump_orders() {
    var i, command int

    if (first_pass) {
        printf("\nStart of jump orders for species #%d, SP %s...\n", species_number, species.name);
    }

    for {
        command = get_command();

        if (command == 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Unknown or missing command.\n");
            continue;
        }

        if (end_of_file || command == END) {
            if (first_pass) {
                printf("End of jump orders for species #%d, SP %s.\n", species_number, species.name);
            }

            if (first_pass) {
                gamemaster_abort_option();
            }

            break;                      /* END for this species. */
        }

        switch (command) {
        case JUMP:
            do_JUMP_command(false, false);
            break;

        case MOVE:
            do_MOVE_command();
            break;

        case PJUMP:
            do_JUMP_command(false, true);
            break;

        case VISITED:
            do_VISITED_command();
            break;

        case WORMHOLE:
            do_WORMHOLE_command();
            break;

        default:
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid jump command.\n");
        }
    }
}

//*************************************************************************
// NoOrders.c

/* This program will generate default orders for a species if no explicit
 *      orders have been provided. */

func NoOrdersMain(argc int, argv []string) {
    var i, j, k, ship_index, locations_fd, my_loc_index, nampla_index, its_loc_index, tonnage, found, alien_number, alien_index, array_index, bit_number, ls_needed, production_penalty int
    var filename[32]byte
    var message_line[132]byte
    var n, nn, raw_material_units, production_capacity, balance, current_base, CUs_needed, IUs_needed, AUs_needed, EUs, bit_mask int
    var message_file, log_file io.Writer
    var alien *species_data
    var nampla, home_nampla, temp_nampla *nampla_data
    var ship *ship_data_
    var locations_base, my_loc, its_loc *sp_loc_data

    /* Check for valid command line. */
    if (argc != 1) {
        fprintf(stderr, "\n\tUsage: NoOrders\n\n");
        exit(0);
    }

    /* Seed random number generator. */
    last_random = time(nil);
    j           = 907;
    for i = 0; i < j; i++ {
        rnd(100);
    }

    /* Get all necessary data. */
    get_galaxy_data();
    get_star_data();
    get_planet_data();
    get_species_data();
    get_location_data();

    truncate_name = true;

    /* Major loop. Check each species in the game. */
    for species_number = 1; species_number <= galaxy.num_species; species_number++ {
        species_index = species_number - 1;

        /* Check if this species is still in the game. */
        if (!data_in_memory[species_index]) {
            continue;
        }

        /* Check if we have orders. */
        filename := fmt.Sprintf("sp%02d.ord", species_number);
        i = open(filename, 0);
        if (i >= 0) {
            close(i);
            continue;
        }

        species     = spec_data[species_index];
        nampla_base = namp_data[species_index];
        ship_base   = ship_data[species_index];
        home_nampla = nampla_base[0]
        home_planet = planet_base[home_nampla.planet_index]

        for i = 0; i < species.num_ships; i++ {
            ship          = ship_base[i]
            ship.special = 0;
        }

        /* Print message for gamemaster. */
        printf("Generating orders for species #%02d, SP %s...\n", species_number, species.name);

        /* Open message file. */
        filename := fmt.Sprintf("noorders.txt");
        message_file = fopen(filename, "r");
        if (message_file == nil) {
            fprintf(stderr, "\n\tCannot open '%s' for reading!\n\n", filename);
            exit(-1);
        }

        /* Open log file. */
        filename := fmt.Sprintf("sp%02d.log", species_number);
        log_file = fopen(filename, "a");
        if (log_file == nil) {
            fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
            exit(-1);
        }

        /* Copy message to log file. */
        for fgets(message_line, 131, message_file) != nil {
            fputs(message_line, log_file);
        }

        fclose(message_file);
        fclose(log_file);

        /* Open orders file for writing. */
        filename := fmt.Sprintf("sp%02d.ord", species_number);
        orders_file = fopen(filename, "w");
        if (orders_file == nil) {
            fprintf(stderr, "\n\tCannot open '%s' for writing!\n\n", filename);
            exit(-1);
        }

        /* Issue PRE-DEPARTURE orders. */
        fprintf(orders_file, "START PRE-DEPARTURE\n");
        fprintf(orders_file, "; Place pre-departure orders here.\n\n");

        for nampla_index = 0; nampla_index < species.num_namplas; nampla_index++ {
            nampla = nampla_base[nampla_index]
            if (nampla.pn == 99) {
                continue;
            }

            /* Generate auto-installs for colonies that were loaded via the DEVELOP command. */
            if (nampla.auto_IUs!= 0) {
                fprintf(orders_file, "\tInstall\t%d IU\tPL %s\n", nampla.auto_IUs, nampla.name);
            }
            if (nampla.auto_AUs != 0) {
                fprintf(orders_file, "\tInstall\t%d AU\tPL %s\n", nampla.auto_AUs, nampla.name);
            }
            if (nampla.auto_IUs != 0 || nampla.auto_AUs != 0) {
                fprintf(orders_file, "\n");
            }

            nampla.item_quantity[CU] -= nampla.auto_IUs + nampla.auto_AUs;

            /* Generate auto UNLOAD orders for transports at this nampla. */
            for j = 0; j < species.num_ships; j++ {
                ship = ship_base[j]
                if (ship.pn == 99) {
                    continue;
                }
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
                if (ship.status == JUMPED_IN_COMBAT) {
                    continue;
                }
                if (ship.status == FORCED_JUMP) {
                    continue;
                }
                if (ship.class != TR) {
                    continue;
                }
                if (ship.item_quantity[CU] < 1) {
                    continue;
                }

                /* New colonies will never be started automatically unless
                 * ship was loaded via a DEVELOP order. */
                if (ship.loading_point != 0) {
                    /* Check if transport is at specified unloading point. */
                    n = ship.unloading_point;
                    if (n == nampla_index || (n == 9999 && nampla_index == 0)) {
                        goto unload_ship;
                    }
                }

                if ((nampla.status & POPULATED) == 0) {
                    continue;
                }

                if ((nampla.mi_base + nampla.ma_base) >= 2000) {
                    continue;
                }

                if (nampla.x == nampla_base.x && nampla.y == nampla_base.y && nampla.z == nampla_base.z) {
                    continue;                                  /* Home sector. */
                }
unload_ship:

                n = ship.loading_point;
                if (n == 9999) {
                    n = 0;              /* Home planet. */
                }
                if (n == nampla_index) {
                    continue;   /* Ship was just loaded here. */
                }
                fprintf(orders_file, "\tUnload\tTR%d%s %s\n\n", ship.tonnage, ship_type[ship.ttype], ship.name);

                nampla.item_quantity[CU] = 0;

                // TODO: fix this. it's trying to find the offset for the unloading point
                ship.special = ship.loading_point;
                n             = nampla - nampla_base;
                if (n == 0) {
                    n = 9999;
                }
                ship.unloading_point = n;
            }

            if (nampla.status & HOME_PLANET) != 0 {
                continue;
            }
            if (nampla.item_quantity[CU] == 0) {
                continue;
            }
            if (nampla.item_quantity[IU] == 0 &&
                nampla.item_quantity[AU] == 0) {
                continue;
            }

            if (nampla.item_quantity[IU] > 0) {
                fprintf(orders_file, "\tInstall\t0 IU\tPL %s\n", nampla.name);
            }
            if (nampla.item_quantity[AU] > 0) {
                fprintf(orders_file, "\tInstall\t0 AU\tPL %s\n\n", nampla.name);
            }
        }

        fprintf(orders_file, "END\n\n");

        fprintf(orders_file, "START JUMPS\n");
        fprintf(orders_file, "; Place jump orders here.\n\n");

        /* Initialize to make sure ships are not given more than one JUMP order. */
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base[i]
            ship.just_jumped = false;
        }

        /* Generate auto-jumps for ships that were loaded via the DEVELOP
         * command or which were UNLOADed because of the AUTO command. */
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base[i]

            if (ship.status == JUMPED_IN_COMBAT) {
                continue;
            }
            if (ship.status == FORCED_JUMP) {
                continue;
            }
            if (ship.pn == 99) {
                continue;
            }
            if (ship.just_jumped) {
                continue;
            }

            j = ship.special;
            if j != 0 {
                if (j == 9999) {
                    j = 0;              /* Home planet. */
                }
                temp_nampla = nampla_base[j];

                fprintf(orders_file, "\tJump\t%s, PL %s\t; ", ship_name(ship), temp_nampla.name);

                print_mishap_chance(ship, temp_nampla.x, temp_nampla.y, temp_nampla.z);

                fprintf(orders_file, "\n\n");

                ship.just_jumped = true;

                continue;
            }

            n = ship.unloading_point;
            if n!=0 {
                if (n == 9999) {
                    n = 0;              /* Home planet. */
                }
                temp_nampla = nampla_base[n];

                if (temp_nampla.x == ship.x && temp_nampla.y == ship.y &&
                    temp_nampla.z == ship.z) {
                    continue;
                }

                fprintf(orders_file, "\tJump\t%s, PL %s\t; ", ship_name(ship), temp_nampla.name);

                print_mishap_chance(ship, temp_nampla.x, temp_nampla.y, temp_nampla.z);

                fprintf(orders_file, "\n\n");

                ship.just_jumped = true;
            }
        }

        /* Generate JUMP orders for all TR1s. */
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base[i];
            if (ship.pn == 99) {
                continue;
            }
            if (ship.status == UNDER_CONSTRUCTION) {
                continue;
            }
            if (ship.status == JUMPED_IN_COMBAT) {
                continue;
            }
            if (ship.status == FORCED_JUMP) {
                continue;
            }
            if (ship.just_jumped) {
                continue;
            }

            if (ship.class == TR && ship.tonnage == 1 &&
                ship.ttype == FTL) {
                fprintf(orders_file, "\tJump\tTR1 %s, ", ship.name);
                closest_unvisited_star(ship);
                fprintf(orders_file, "\n\t\t\t; Age %d, now at %d %d %d, ", ship.age, ship.x, ship.y, ship.z);
                print_mishap_chance(ship, x, y, z);

                ship.dest_x = x;
                ship.dest_y = y;
                ship.dest_z = z;

                fprintf(orders_file, "\n\n");

                ship.just_jumped = true;
            }
        }

        fprintf(orders_file, "END\n\n");

        fprintf(orders_file, "START PRODUCTION\n");

        /* Generate a PRODUCTION order for each planet that can produce. */
        for nampla_index = species.num_namplas - 1; nampla_index >= 0; nampla_index-- {
            nampla = nampla_base[nampla_index];
            if (nampla.pn == 99) {
                continue;
            }

            if (nampla.mi_base == 0 && (nampla.status & RESORT_COLONY) == 0) {
                continue;
            }
            if (nampla.ma_base == 0 && (nampla.status & MINING_COLONY) == 0) {
                continue;
            }

            fprintf(orders_file, "    PRODUCTION PL %s\n", nampla.name);

            if (nampla.status & MINING_COLONY) != 0 {
                fprintf(orders_file, "    ; The above PRODUCTION order is required for this mining colony, even\n");
                fprintf(orders_file, "    ;  if no other production orders are given for it.\n");
            }else if (nampla.status & RESORT_COLONY) != 0 {
                fprintf(orders_file, "    ; The above PRODUCTION order is required for this resort colony, even\n");
                fprintf(orders_file, "    ;  though no other production orders can be given for it.\n");
            }else {
                fprintf(orders_file, "    ; Place production orders here for planet %s.\n\n", nampla.name);
            }

            /* Build IUs and AUs for incoming ships with CUs. */
            if nampla.IUs_needed != 0 {
                fprintf(orders_file, "\tBuild\t%d IU\n", nampla.IUs_needed);
            }
            if nampla.AUs_needed != 0 {
                fprintf(orders_file, "\tBuild\t%d AU\n", nampla.AUs_needed);
            }
            if nampla.IUs_needed != 0 || nampla.AUs_needed != 0 {
                fprintf(orders_file, "\n");
            }

            if (nampla.status & MINING_COLONY) != 0 {
                continue;
            }
            if (nampla.status & RESORT_COLONY) != 0 {
                continue;
            }

            /* See if there are any RMs to recycle. */
            n = nampla.special / 5;
            if (n > 0) {
                fprintf(orders_file, "\tRecycle\t%d RM\n\n", 5 * n);
            }

            /* Generate DEVELOP commands for ships arriving here because of AUTO command. */
            for i = 0; i < species.num_ships; i++ {
                ship = ship_base[i];
                if (ship.pn == 99) {
                    continue;
                }

                k = ship.special;
                if (k == 0) {
                    continue;
                }
                if (k == 9999) {
                    k = 0;              /* Home planet. */
                }
                if (nampla != nampla_base[k]) {
                    continue;
                }

                k = ship.unloading_point;
                if (k == 9999) {
                    k = 0;
                }
                temp_nampla = nampla_base[k];

                fprintf(orders_file, "\tDevelop\tPL %s, TR%d%s %s\n\n", temp_nampla.name, ship.tonnage, ship_type[ship.ttype], ship.name);
            }

            /* Give orders to continue construction of unfinished ships and starbases. */
            for i = 0; i < species.num_ships; i++ {
                ship = ship_base[i];
                if (ship.pn == 99) {
                    continue;
                }
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

                if (ship.status == UNDER_CONSTRUCTION) {
                    fprintf(orders_file, "\tContinue\t%s, %d\t; Left to pay = %d\n\n", ship_name(ship), ship.remaining_cost, ship.remaining_cost);
                    continue;
                }

                if (ship.ttype != STARBASE) {
                    continue;
                }

                j = (species.tech_level[MA] / 2) - ship.tonnage;
                if (j < 1) {
                    continue;
                }

                fprintf(orders_file, "\tContinue\tBAS %s, %d\t; Current tonnage = %s\n\n", ship.name, 100 * j, commas(10000 * ship.tonnage));
            }

            /* Generate DEVELOP command if this is a colony with an economic base less than 200. */
            n = nampla.mi_base + nampla.ma_base + nampla.IUs_needed + nampla.AUs_needed;
            if ((nampla.status & COLONY) != 0 && n < 2000 && nampla.pop_units > 0) {
                if (nampla.pop_units > (2000 - n)) {
                    nn = 2000 - n;
                }else{
                    nn = nampla.pop_units;
                }

                fprintf(orders_file, "\tDevelop\t%d\n\n",2 * nn);

                nampla.IUs_needed += nn;
            }

            /* For home planets and any colonies that have an economic base of
             *  at least 200, check if there are other colonized planets in
             *  the same sector that are not self-sufficient.  If so, DEVELOP
             *  them. */
            if (n >= 2000 || (nampla.status & HOME_PLANET) != 0) {
                for i = 1; i < species.num_namplas; i++ { /* Skip HP. */
                    if (i == nampla_index) {
                        continue;
                    }

                    temp_nampla = nampla_base[i];

                    if (temp_nampla.pn == 99) {
                        continue;
                    }
                    if (temp_nampla.x != nampla.x) {
                        continue;
                    }
                    if (temp_nampla.y != nampla.y) {
                        continue;
                    }
                    if (temp_nampla.z != nampla.z) {
                        continue;
                    }

                    n = temp_nampla.mi_base + temp_nampla.ma_base + temp_nampla.IUs_needed + temp_nampla.AUs_needed;
                    if (n == 0) {
                        continue;
                    }

                    nn = temp_nampla.item_quantity[IU] + temp_nampla.item_quantity[AU];
                    if (nn > temp_nampla.item_quantity[CU]) {
                        nn = temp_nampla.item_quantity[CU];
                    }
                    n += nn;
                    if (n >= 2000) {
                        continue;
                    }
                    nn = 2000 - n;

                    if (nn > nampla.pop_units) {
                        nn = nampla.pop_units;
                    }

                    fprintf(orders_file, "\tDevelop\t%d\tPL %s\n\n", 2 * nn, temp_nampla.name);

                    temp_nampla.AUs_needed += nn;
                }
            }
        }

        fprintf(orders_file, "END\n\n");

        fprintf(orders_file, "START POST-ARRIVAL\n");
        fprintf(orders_file, "; Place post-arrival orders here.\n\n");

        /* Generate an AUTO command. */
        fprintf(orders_file, "\tAuto\n\n");

        /* Generate SCAN orders for all TR1s in sectors that current species does not inhabit. */
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base[i];

            if (ship.pn == 99) {
                continue;
            }
            if (ship.status == UNDER_CONSTRUCTION) {
                continue;
            }
            if (ship.class != TR) {
                continue;
            }
            if (ship.tonnage != 1) {
                continue;
            }
            if (ship.ttype != FTL) {
                continue;
            }
            if (ship.dest_x == -1) {
                continue;                       /* Not jumping anywhere. */
            }
            found = false;
            for j = 1; j < species.num_namplas; j++ { /* Skip home sector. */
                nampla = nampla_base[j];
                if (nampla.pn == 99) {
                    continue;
                }
                if (nampla.x != ship.dest_x) {
                    continue;
                }
                if (nampla.y != ship.dest_y) {
                    continue;
                }
                if (nampla.z != ship.dest_z) {
                    continue;
                }

                if (nampla.status & POPULATED) !=  0 {
                    found = true;
                    break;
                }
            }
            if (!found) {
                fprintf(orders_file, "\tScan\tTR1 %s\n", ship.name);
            }
        }

        fprintf(orders_file, "END\n\n");

        /* Clean up for this species. */
        fclose(orders_file);
    }

    /* Clean up and exit. */
    free_species_data();
    exit(0);
}

func print_mishap_chance(ship *ship_data_, destx, desty, destz int) {
    var mishap_GV, mishap_age, mishap_chance, success_chance int

    if (destx == -1) {
        fprintf(orders_file, "Mishap chance = ???");
        return;
    }

    mishap_GV  = species.tech_level[GV];
    mishap_age = ship.age;

    x, y, z             := destx, desty, destz;
    mishap_chance = (100 * ( ((x - ship.x) * (x - ship.x)) + ((y - ship.y) * (y - ship.y)) + ((z - ship.z) * (z - ship.z)) )) / mishap_GV;

    if (mishap_age > 0 && mishap_chance < 10000) {
        success_chance  = 10000 - mishap_chance;
        success_chance -= (2 * mishap_age * success_chance) / 100;
        mishap_chance   = 10000 - success_chance;
    }

    if (mishap_chance > 10000) {
        mishap_chance = 10000;
    }

    fprintf(orders_file, "mishap chance = %d.%02d%%", mishap_chance / 100, mishap_chance % 100);
}

func closest_unvisited_star(ship *ship_data_) {
    var i, found, species_array_index, species_bit_number, stx, sty, stz, closest_distance, temp_distance, species_bit_mask int
    var star, closest_star *star_data

    /* Get array index and bit mask. */
    species_array_index = (species_number - 1) / 32;
    species_bit_number  = (species_number - 1) % 32;
    species_bit_mask    = 1 << species_bit_number;

    shx, shy, shz := ship.x, ship.y, ship.z;

    x = -1;
    closest_distance = 999999;

    found = false;
    for i = 0; i < num_stars; i++ {
        star = star_base[i];

        /* Check if bit is already set. */
        if star.visited_by[species_number] {
            continue;
        }

        stx, sty, stz := star.x, star.y, star.z;
        temp_distance = ((shx - stx) * (shx - stx)) + ((shy - sty) * (shy - sty)) + ((shz - stz) * (shz - stz));
        if (temp_distance < closest_distance) {
            x, y, z = stx, sty, stz;
            closest_distance = temp_distance;
            closest_star     = star;
            found            = true;
        }
    }

    if (found) {
        fprintf(orders_file, "%d %d %d", x, y, z);
        closest_star.visited_by[species_number] = true /* So that we don't send more than one ship to the same place. */
    }else {
        fprintf(orders_file, "???");
    }

    return;
}

//*************************************************************************
// PostArrival.c

func PostArrivalMain(argc int, argv []string) {
    int i, n, found, num_species, sp_num[MAX_SPECIES], sp_index,
        command, do_all_species;

    char filename[32], keyword[4];


    /* Seed random number generator. */
    last_random = time(nil);
    n           = rnd(100) + rnd(200) + rnd(300);
    for i = 0; i < n; i++ {
        rnd(10);
    }

    /* Get commonly used data. */
    get_galaxy_data();
    get_transaction_data();

    ignore_field_distorters = true;

    /* Check arguments. If an argument is -p, then do two passes. In the
     *  first pass, display results and prompt the GM, allowing him to
     *  abort if necessary before saving results to disk. All other
     *  arguments must be species numbers. If no species numbers are
     *  specified, then do all species. */
    num_species  = 0;
    first_pass   = false;
    test_mode    = false;
    verbose_mode = false;
    for i = 1; i < argc; i++ {
        if (strcmp(argv[i], "-p") == 0) {
            first_pass = true;
        }else if (strcmp(argv[i], "-t") == 0) {
            test_mode = true;
        }else if (strcmp(argv[i], "-v") == 0) {
            verbose_mode = true;
        }else{
            n = atoi(argv[i]);
            if (n < 1 || n > galaxy.num_species) {
                fprintf(stderr,
                        "\n    '%s' is not a valid argument!\n", argv[i]);
                exit(-1);
            }
            sp_num[num_species++] = n;
        }
    }

    if (num_species == 0) {
        num_species = galaxy.num_species;
        for i = 0; i < num_species; i++ {
            sp_num[i] = i + 1;
        }
        do_all_species = true;
    }else {
        do_all_species = false;
    }

    /* Two passes through all orders will be done. The first pass will
     *  check for errors and abort if any are found. Results will be written
     *  to disk only on the second pass. */

start_pass:

    if (first_pass) {
        printf("\nStarting first pass...\n\n");
    }

    get_star_data();
    get_planet_data();
    get_species_data();

    /* Main loop. For each species, take appropriate action. */
    for sp_index = 0; sp_index < num_species; sp_index++ {
        species_number = sp_num[sp_index];
        species_index  = species_number - 1;

        found = data_in_memory[species_index];
        if (!found) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    Skipping species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n    Cannot get data for species #%d!\n",
                        species_number);
                exit(-1);
            }
        }

        species     = &spec_data[species_index];
        nampla_base = namp_data[species_index];
        ship_base   = ship_data[species_index];

        /* Do some initializations. */
        species.auto_orders = false;

        /* Open orders file for this species. */
        filename := fmt.Sprintf("sp%02d.ord", species_number);
        input_file = fopen(filename, "r");
        if (input_file == nil) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    No orders for species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n\tCannot open '%s' for reading!\n\n", filename);
                exit(-1);
            }
        }

        end_of_file = false;

        just_opened_file = true;        /* Tell parse.c to skip mail header,
                                         *      if any. */
find_start:

        /* Search for START POST-ARRIVAL order. */
        found = false;
        for !found {
            command = get_command();
            if (command == MESSAGE) {
                /* Skip MESSAGE text. It may contain a line that starts
                 *      with "start". */
                for {
                    command = get_command();
                    if (command < 0) {
                        fprintf(stderr,
                                "WARNING: Unterminated MESSAGE command in file %s!\n",
                                filename);
                        break;
                    }

                    if (command == ZZZ) {
                        goto find_start;
                    }
                }
            }

            if (command < 0) {
                break;          /* End of file. */
            }
            if (command != START) {
                continue;
            }

            /* Get the first three letters of the keyword and convert to
             *  upper case. */
            skip_whitespace();
            for i = 0; i < 3; i++ {
                keyword[i] = toupper(*input_line_pointer);
                input_line_pointer++
            }
            keyword[3] = 0;

            if (strcmp(keyword, "POS") == 0) {
                found = true;
            }
        }

        if (!found) {
            if (first_pass) {
                printf("\nNo post-arrival orders for species #%d, SP %s.\n",
                       species_number, species.name);
            }
            goto done_orders;
        }

        /* Open log file. Use stdout for first pass. */
        log_stdout = false;  /* We will control value of log_file from here. */
        if (first_pass) {
            log_file = stdout;
        }else {
            /* Open log file for appending. */
            filename := fmt.Sprintf("sp%02d.log", species_number);
            log_file = fopen(filename, "a");
            if (log_file == nil) {
                fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
                exit(-1);
            }
            log_string("\nPost-arrival orders:\n");
        }

        /* For each ship, set dest_z to zero.  If a starbase is used as a
         *      gravitic telescope, it will be set to non-zero.  This will
         *      prevent more than one TELESCOPE order per turn per starbase. */
        ship = ship_base;
        for i = 0; i < species.num_ships; i++ {
            ship.dest_z = 0;
            ship++
        }

        /* Handle post-arrival orders for this species. */
        do_postarrival_orders();

        data_modified[species_index] = true;

        /* If this is the second pass, close the log file. */
        if (!first_pass) {
            fclose(log_file);
        }

done_orders:

        fclose(input_file);
    }

    if (first_pass) {
        printf("\nFinal chance to abort safely!\n");
        gamemaster_abort_option();
        first_pass = false;
        free_species_data();
        free(star_base);        /* In case data was modified. */
        free(planet_base);      /* In case data was modified. */

        printf("\nStarting second pass...\n\n");

        goto start_pass;
    }

    save_species_data();
    save_transaction_data();
    if (star_data_modified) {
        save_star_data();
    }
    if (planet_data_modified) {
        save_planet_data();
    }
    free_species_data();
    free(star_base);
    free(planet_base);
    exit(0);
}



do_postarrival_orders() {
    int i, command;


    if (first_pass) {
        printf("\nStart of post-arrival orders for species #%d, SP %s...\n",
               species_number, species.name);
    }

    truncate_name = true;       /* For these commands, do not display age
                                 * or landed/orbital status of ships. */

    for {
        command = get_command();

        if (command == 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Unknown or missing command.\n");
            continue;
        }

        if (end_of_file || command == END) {
            if (first_pass) {
                printf("End of post-arrival orders for species #%d, SP %s.\n",
                       species_number, species.name);
            }

            if (first_pass) {
                gamemaster_abort_option();
            }

            break;                      /* END for this species. */
        }

        switch (command) {
        case ALLY:
            do_ALLY_command();
            break;

        case AUTO:
            species.auto_orders = true;
            log_string("    An AUTO order was executed.\n");
            break;

        case DEEP:
            do_DEEP_command();
            break;

        case DESTROY:
            do_DESTROY_command();
            break;

        case ENEMY:
            do_ENEMY_command();
            break;

        case LAND:
            do_LAND_command();
            break;

        case MESSAGE:
            do_MESSAGE_command();
            break;

        case NAME:
            do_NAME_command();
            break;

        case NEUTRAL:
            do_NEUTRAL_command();
            break;

        case ORBIT:
            do_ORBIT_command();
            break;

        case REPAIR:
            do_REPAIR_command();
            break;

        case SCAN:
            do_SCAN_command();
            break;

        case SEND:
            do_SEND_command();
            break;

        case TEACH:
            do_TEACH_command();
            break;

/*	    case TECH:
 *              do_TECH_command ();
 *              break;
 */
        case TRANSFER:
            do_TRANSFER_command();
            break;

        case TELESCOPE:
            do_TELESCOPE_command();
            break;

        case TERRAFORM:
            do_TERRAFORM_command();
            break;

        default:
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid post-arrival command.\n");
        }
    }
}



//*************************************************************************
// PreDeparture.c

func PreDepartureMain(argc int, argv []string) {
    int i, n, found, num_species, sp_num[MAX_SPECIES], sp_index,
        command, do_all_species;

    char filename[32], keyword[4];


    /* Seed random number generator. */
    last_random = time(nil);
    n           = rnd(100) + rnd(200) + rnd(300);
    for i = 0; i < n; i++ {
        rnd(10);
    }

    /* Get commonly used data. */
    get_galaxy_data();
    get_transaction_data();

    ignore_field_distorters = true;

    /* Check arguments. If an argument is -p, then do two passes. In the
     *  first pass, display results and prompt the GM, allowing him to
     *  abort if necessary before saving results to disk. If an argument
     *  is -t, then set test mode. All other arguments must be species
     *  numbers. If no species numbers are specified, then do all species. */
    num_species  = 0;
    first_pass   = false;
    test_mode    = false;
    verbose_mode = false;
    for i = 1; i < argc; i++ {
        if (strcmp(argv[i], "-p") == 0) {
            first_pass = true;
        }else if (strcmp(argv[i], "-t") == 0) {
            test_mode = true;
        }else if (strcmp(argv[i], "-v") == 0) {
            verbose_mode = true;
        }else{
            n = atoi(argv[i]);
            if (n < 1 || n > galaxy.num_species) {
                fprintf(stderr,
                        "\n    '%s' is not a valid argument!\n", argv[i]);
                exit(-1);
            }
            sp_num[num_species++] = n;
        }
    }

    if (num_species == 0) {
        num_species = galaxy.num_species;
        for i = 0; i < num_species; i++ {
            sp_num[i] = i + 1;
        }
        do_all_species = true;
    }else {
        do_all_species = false;
    }

    /* Two passes through all orders will be done. The first pass will
     *  check for errors and abort if any are found. Results will be written
     *  to disk only on the second pass. */

start_pass:

    if (first_pass) {
        printf("\nStarting first pass...\n\n");
    }

    get_species_data();
    get_star_data();
    get_planet_data();

    /* Main loop. For each species, take appropriate action. */
    for sp_index = 0; sp_index < num_species; sp_index++ {
        species_number = sp_num[sp_index];
        species_index  = species_number - 1;

        found = data_in_memory[species_index];
        if (!found) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    Skipping species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n    Cannot get data for species #%d!\n",
                        species_number);
                exit(-1);
            }
        }

        species     = &spec_data[species_index];
        nampla_base = namp_data[species_index];
        ship_base   = ship_data[species_index];

        /* Open orders file for this species. */
        filename := fmt.Sprintf("sp%02d.ord", species_number);
        input_file = fopen(filename, "r");
        if (input_file == nil) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    No orders for species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n\tCannot open '%s' for reading!\n\n", filename);
                exit(-1);
            }
        }

        end_of_file = false;

        just_opened_file = true;        /* Tell parse.c to skip mail header,
                                         *      if any. */
find_start:

        /* Search for START PRE-DEPARTURE order. */
        found = false;
        for !found {
            command = get_command();
            if (command == MESSAGE) {
                /* Skip MESSAGE text. It may contain a line that starts
                 *      with "start". */
                for {
                    command = get_command();
                    if (command < 0) {
                        fprintf(stderr,
                                "WARNING: Unterminated MESSAGE command in file %s!\n",
                                filename);
                        break;
                    }

                    if (command == ZZZ) {
                        goto find_start;
                    }
                }
            }

            if (command < 0) {
                break;          /* End of file. */
            }
            if (command != START) {
                continue;
            }

            /* Get the first three letters of the keyword and convert to
             *  upper case. */
            skip_whitespace();
            for i = 0; i < 3; i++ {
                keyword[i] = toupper(*input_line_pointer);
                input_line_pointer++
            }
            keyword[3] = 0;

            if (strcmp(keyword, "PRE") == 0) {
                found = true;
            }
        }

        if (!found) {
            if (first_pass) {
                printf("\nNo pre-departure orders for species #%d, SP %s.\n",
                       species_number, species.name);
            }
            goto done_orders;
        }

        /* Open log file. Use stdout for first pass. */
        log_stdout = false;  /* We will control value of log_file from here. */
        if (first_pass) {
            log_file = stdout;
        }else {
            /* Open log file for appending. */
            filename := fmt.Sprintf("sp%02d.log", species_number);
            log_file = fopen(filename, "a");
            if (log_file == nil) {
                fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
                exit(-1);
            }
            log_string("\nPre-departure orders:\n");
        }

        /* Handle predeparture orders for this species. */
        do_predeparture_orders();

        data_modified[species_index] = true;

        if (first_pass) {
            goto done_orders;
        }

        /* If this is the second pass, close the log file. */
        if (!first_pass) {
            fclose(log_file);
        }

done_orders:

        fclose(input_file);
    }

    if (first_pass) {
        printf("\nFinal chance to abort safely!\n");
        gamemaster_abort_option();
        first_pass = false;
        free_species_data();
        free(star_base);        /* In case data was modified. */

        printf("\nStarting second pass...\n\n");

        goto start_pass;
    }

    save_species_data();
    save_transaction_data();
    if (star_data_modified) {
        save_star_data();
    }
    if (planet_data_modified) {
        save_planet_data();
    }
    free_species_data();
    free(star_base);
    free(planet_base);
    exit(0);
}



do_predeparture_orders() {
    int i, command, old_test_mode;


    if (first_pass) {
        printf("\nStart of pre-departure orders for species #%d, SP %s...\n",
               species_number, species.name);
    }

    truncate_name = true;       /* For these commands, do not display age
                                 * or landed/orbital status of ships. */

    for {
        command = get_command();

        if (command == 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Unknown or missing command.\n");
            continue;
        }

        if (end_of_file || command == END) {
            if (first_pass) {
                printf("End of pre-departure orders for species #%d, SP %s.\n",
                       species_number, species.name);
            }

            if (first_pass) {
                gamemaster_abort_option();
            }

            break;                      /* END for this species. */
        }

        switch (command) {
        case ALLY:
            do_ALLY_command();
            break;

        case BASE:
            do_BASE_command();
            break;

        case DEEP:
            do_DEEP_command();
            break;

        case DESTROY:
            do_DESTROY_command();
            break;

        case DISBAND:
            do_DISBAND_command();
            break;

        case ENEMY:
            do_ENEMY_command();
            break;

        case INSTALL:
            do_INSTALL_command();
            break;

        case LAND:
            do_LAND_command();
            break;

        case MESSAGE:
            do_MESSAGE_command();
            break;

        case NAME:
            do_NAME_command();
            break;

        case NEUTRAL:
            do_NEUTRAL_command();
            break;

        case ORBIT:
            do_ORBIT_command();
            break;

        case REPAIR:
            do_REPAIR_command();
            break;

        case SCAN:
            /* Scan is okay in test mode for pre-departure. */
            old_test_mode = test_mode;
            test_mode     = false;
            do_SCAN_command();
            test_mode = old_test_mode;
            break;

        case SEND:
            do_SEND_command();
            break;

        case TRANSFER:
            do_TRANSFER_command();
            break;

        case UNLOAD:
            do_UNLOAD_command();
            break;

        default:
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid pre-departure command.\n");
        }
    }
}



//*************************************************************************
// Production.c
func ProductionMain(argc int, argv []string) {
    int i, n, found, num_species, sp_num[MAX_SPECIES], sp_index,
        command, do_all_species;

    char filename[32], keyword[4];


    /* Seed random number generator. */
    last_random = time(nil);
    n           = rnd(100) + rnd(200) + rnd(300);
    for i = 0; i < n; i++ {
        rnd(10);
    }

    /* Get commonly used data. */
    get_galaxy_data();
    get_transaction_data();

    ignore_field_distorters = true;

    /* Check arguments. If an argument is -p, then do two passes. In the
     *  first pass, display results and prompt the GM, allowing him to
     *  abort if necessary before saving results to disk. All other
     *  arguments must be species numbers. If no species numbers are
     *  specified, then do all species. */
    num_species  = 0;
    first_pass   = false;
    test_mode    = false;
    verbose_mode = false;
    for i = 1; i < argc; i++ {
        if (strcmp(argv[i], "-p") == 0) {
            first_pass = true;
        }else if (strcmp(argv[i], "-t") == 0) {
            test_mode = true;
        }else if (strcmp(argv[i], "-v") == 0) {
            verbose_mode = true;
        }else{
            n = atoi(argv[i]);
            if (n < 1 || n > galaxy.num_species) {
                fprintf(stderr,
                        "\n    '%s' is not a valid argument!\n", argv[i]);
                exit(-1);
            }
            sp_num[num_species++] = n;
        }
    }

    if (num_species == 0) {
        num_species = galaxy.num_species;
        for i = 0; i < num_species; i++ {
            sp_num[i] = i + 1;
        }
        do_all_species = true;
    }else {
        do_all_species = false;
    }

    /* Two passes through all orders will be done. The first pass will
     *  check for errors and abort if any are found. Results will be written
     *  to disk only on the second pass. */

start_pass:

    if (first_pass) {
        printf("\nStarting first pass...\n\n");
    }

    get_species_data();
    get_star_data();
    get_planet_data();

    /* Main loop. For each species, take appropriate action. */
    for sp_index = 0; sp_index < num_species; sp_index++ {
        species_number = sp_num[sp_index];
        species_index  = species_number - 1;

        found = data_in_memory[species_index];
        if (!found) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    Skipping species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n    Cannot get data for species #%d!\n",
                        species_number);
                exit(-1);
            }
        }

        species     = &spec_data[species_index];
        nampla_base = namp_data[species_index];
        ship_base   = ship_data[species_index];

        home_planet = planet_base + (int)nampla_base.planet_index;

        /* Open orders file for this species. */
        filename := fmt.Sprintf("sp%02d.ord", species_number);
        input_file = fopen(filename, "r");
        if (input_file == nil) {
            if (do_all_species) {
                if (first_pass) {
                    printf("\n    No orders for species #%d.\n", species_number);
                }
                continue;
            }else {
                fprintf(stderr, "\n\tCannot open '%s' for reading!\n\n", filename);
                exit(-1);
            }
        }

        end_of_file = false;

        just_opened_file = true;        /* Tell parse.c to skip mail header,
                                         *      if any. */
find_start:

        /* Search for START PRODUCTION order. */
        found = false;
        for !found {
            command = get_command();
            if (command == MESSAGE) {
                /* Skip MESSAGE text. It may contain a line that starts
                 *      with "start". */
                for {
                    command = get_command();
                    if (command < 0) {
                        fprintf(stderr,
                                "WARNING: Unterminated MESSAGE command in file %s!\n",
                                filename);
                        break;
                    }

                    if (command == ZZZ) {
                        goto find_start;
                    }
                }
            }

            if (command < 0) {
                break;          /* End of file. */
            }
            if (command != START) {
                continue;
            }

            /* Get the first three letters of the keyword and convert to
             *  upper case. */
            skip_whitespace();
            for i = 0; i < 3; i++ {
                keyword[i] = toupper(*input_line_pointer);
                input_line_pointer++
            }
            keyword[3] = 0;

            if (strcmp(keyword, "PRO") == 0) {
                found = true;
            }
        }

        if (!found) {
            if (first_pass) {
                printf("\nNo production orders for species #%d, SP %s.\n",
                       species_number, species.name);
            }
            goto done_orders;
        }

        /* Open log file. Use stdout for first pass. */
        log_stdout = false;  /* We will control value of log_file from here. */
        if (first_pass) {
            log_file = stdout;
        }else {
            /* Open log file for appending. */
            filename := fmt.Sprintf("sp%02d.log", species_number);
            log_file = fopen(filename, "a");
            if (log_file == nil) {
                fprintf(stderr, "\n\tCannot open '%s' for appending!\n\n", filename);
                exit(-1);
            }
            fprintf(log_file, "\nProduction orders:\n");
            fprintf(log_file, "\n  Number of economic units at start of production: %d\n\n", species.econ_units);
        }

        /* Initialize "done" arrays. They will be used to prevent more
         *  than one corresponding order per planet. */
        if (species.num_namplas > 1000) {
            fprintf(stderr, "\n\n\tInternal error. xxx_done array overflow!/n/n");
            exit(-1);
        }

        for i = 0; i < species.num_namplas; i++ {
            production_done[i] = false;
        }

        /* Do other initializations. */
        for i = 0; i < species.num_namplas; i++ {
            nampla             = nampla_base + i;
            nampla.auto_IUs   = 0;
            nampla.auto_AUs   = 0;
            nampla.IUs_needed = 0;
            nampla.AUs_needed = 0;
        }

        /* Handle production orders for this species. */
        num_intercepts = 0;
        for i = 0; i < 6; i++ {
            sp_tech_level[i] = species.tech_level[i];
        }
        do_production_orders();
        for i = 0; i < 6; i++ {
            species.tech_level[i] = sp_tech_level[i];
        }

        for i = 0; i < num_intercepts; i++ {
            handle_intercept(i);
        }

        data_modified[species_index] = true;

        /* If this is the second pass, close the log file. */
        if (!first_pass) {
            fclose(log_file);
        }

done_orders:

        fclose(input_file);
    }

    if (first_pass) {
        printf("\nFinal chance to abort safely!\n");
        gamemaster_abort_option();
        first_pass = false;
        free_species_data();
        free(star_base);        /* In case data was modified. */
        free(planet_base);      /* In case data was modified. */

        printf("\nStarting second pass...\n\n");

        goto start_pass;
    }

    save_species_data();
    free_species_data();
    if (planet_data_modified) {
        save_planet_data();
    }
    free(planet_base);
    save_transaction_data();
    exit(0);
}



do_production_orders() {
    int i, command;


    truncate_name = true;       /* For these commands, do not display age
                                 * or landed/orbital status of ships. */


    if (first_pass) {
        printf("\nStart of production orders for species #%d, SP %s...\n",
               species_number, species.name);
    }

    doing_production = false;   /* This will be set as soon as production
                                 * actually starts. */
    for {
        command = get_command();

        if (command == 0) {
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Unknown or missing command.\n");
            continue;
        }

        if (end_of_file || command == END) {
            /* Handle planets that were not given PRODUCTION orders. */
            next_nampla = nampla_base - 1;
            for i = 0; i < species.num_namplas; i++ {
                next_nampla++

                if (production_done[i]) {
                    continue;
                }

                production_done[i] = true;

                if (next_nampla.status & DISBANDED_COLONY) {
                    continue;
                }

                if (next_nampla.mi_base + next_nampla.ma_base == 0) {
                    continue;
                }

                next_nampla_index = i;

                do_PRODUCTION_command(true);
            }

            transfer_balance();         /* Terminate production for
                                         * last planet for this species. */

            if (first_pass) {
                gamemaster_abort_option();
                printf("\nEnd of production orders for species #%d, SP %s.\n",
                       species_number, species.name);
            }

            break;                      /* END for this species. */
        }

        switch (command) {
        case ALLY:
            do_ALLY_command();
            break;

        case AMBUSH:
            do_AMBUSH_command();
            break;

        case BUILD:
            do_BUILD_command(false, false);
            break;

        case CONTINUE:
            do_BUILD_command(true, false);
            break;

        case DEVELOP:
            do_DEVELOP_command();
            break;

        case ENEMY:
            do_ENEMY_command();
            break;

        case ESTIMATE:
            do_ESTIMATE_command();
            break;

        case HIDE:
            do_HIDE_command();
            break;

        case IBUILD:
            do_BUILD_command(false, true);
            break;

        case ICONTINUE:
            do_BUILD_command(true, true);
            break;

        case INTERCEPT:
            do_INTERCEPT_command();
            break;

        case NEUTRAL:
            do_NEUTRAL_command();
            break;

        case PRODUCTION:
            do_PRODUCTION_command(false);
            break;

        case RECYCLE:
            do_RECYCLE_command();
            break;

        case RESEARCH:
            do_RESEARCH_command();
            break;

        case SHIPYARD:
            do_SHIPYARD_command();
            break;

        case UPGRADE:
            do_UPGRADE_command();
            break;

        default:
            fprintf(log_file, "!!! Order ignored:\n");
            fprintf(log_file, "!!! %s", input_line);
            fprintf(log_file, "!!! Invalid production command.\n");
        }
    }
}



//*************************************************************************
// Report.c


/* This program will generate reports for all species in the game and
 * write them to separate files. Each report will consist of a log of the
 * events of the previous turn, plus status information for the current
 * turn. */

ReportMain(argc int, argv []string) {
    int i, j, k, ship_index, locations_fd, my_loc_index, its_loc_index,
        industry, turn_number, alien_number, species_fd,
        header_printed, alien_can_hide, do_this_species, sp_index,
        array_index, bit_number, we_have_colony_here, nampla_index,
        we_have_planet_here, found, ls_needed, production_penalty,
        temp_ignore_field_distorters;

    char filename[32], log_line[256], temp1[16], temp2[128];

    long n, nn, bit_mask;

    struct species_data *alien;
    struct nampla_data * nampla, *alien_nampla, *our_nampla,
                       *temp_nampla;
    struct ship_data *  ship, *ship2, *alien_ship;
    struct sp_loc_data *locations_base, *my_loc, *its_loc;


    /* Check for options, if any. */
    test_mode    = false;
    verbose_mode = false;
    for i = 1; i < argc; i++ {
        if (strcmp(argv[i], "-t") == 0) {
            test_mode = true;
        }
        if (strcmp(argv[i], "-v") == 0) {
            verbose_mode = true;
        }
    }

    /* Get all necessary data. */
    get_galaxy_data();
    get_star_data();
    get_planet_data();
    get_species_data();
    get_location_data();

    turn_number = galaxy.turn_number;

    /* Generate a report for each species. */
    alien_number = 0;   /* Pointers to alien data not yet assigned. */
    for species_number = 1; species_number <= galaxy.num_species; species_number++ {
        /* Check if we are doing all species, or just one or more specified
         *      ones. */
        do_this_species = true;
        if (argc > 1) {
            do_this_species = false;
            for i = 1; i < argc; i++ {
                j = atoi(argv[i]);
                if (species_number == j) {
                    do_this_species = true;
                    break;
                }
            }
        }

        if (!do_this_species) {
            continue;
        }

        /* Check if this species is still in the game. */
        if (!data_in_memory[species_number - 1]) {
            if (argc == 1) {
                continue;       /* This species is no longer in the game. */
            }
            fprintf(stderr, "\n\tCannot open data file for species #%d!\n\n",
                    species_number);
            exit(-1);
        }

        species      = &spec_data[species_number - 1];
        nampla_base  = namp_data[species_number - 1];
        nampla1_base = nampla_base;
        ship_base    = ship_data[species_number - 1];
        ship1_base   = ship_base;
        home_planet  = planet_base + nampla1_base.planet_index;

        /* Print message for gamemaster. */
        if (verbose_mode) {
            printf("Generating turn %d report for species #%d, SP %s...\n",
                   turn_number, species_number, species.name);
        }

        /* Open report file for writing. */
        filename := fmt.Sprintf("sp%02d.rpt.t%d", species_number, turn_number);
        report_file = fopen(filename, "w");
        if (report_file == nil) {
            fprintf(stderr, "\n\tCannot open '%s' for writing!\n\n", filename);
            exit(-1);
        }

        /* Copy log file, if any, to output file. */
        filename := fmt.Sprintf("sp%02d.log", species_number);
        log_file = fopen(filename, "r");
        if (log_file != nil) {
            if (turn_number > 1) {
                fprintf(report_file, "\n\n\t\t\tEVENT LOG FOR TURN %d\n",
                        turn_number - 1);
            }

            for fgets(log_line, 256, log_file) != nil {
                fputs(log_line, report_file);
            }

            fprintf(report_file, "\n\n");

            fclose(log_file);
        }

        /* Print header for status report. */
        fprintf(report_file,
                "\n\t\t\t SPECIES STATUS\n\n\t\t\tSTART OF TURN %d\n\n",
                turn_number);

        fprintf(report_file, "Species name: %s\n", species.name);
        fprintf(report_file, "Government name: %s\n", species.govt_name);
        fprintf(report_file, "Government ttype: %s\n", species.govt_ttype);

        fprintf(report_file, "\nTech Levels:\n");
        for i = 0; i < 6; i++ {
            fprintf(report_file, "   %s = %d", tech_name[i],
                    species.tech_level[i]);
            if (species.tech_knowledge[i] > species.tech_level[i]) {
                fprintf(report_file, "/%d", species.tech_knowledge[i]);
            }
            fprintf(report_file, "\n");
        }

        fprintf(report_file, "\nAtmospheric Requirement: %d%%-%d%% %s",
                (int)species.required_gas_min,
                (int)species.required_gas_max,
                gas_string[species.required_gas]);
        fprintf(report_file, "\nNeutral Gases:");
        for i = 0; i < 6; i++ {
            if (i != 0) {
                fprintf(report_file, ",");
            }
            fprintf(report_file, " %s", gas_string[species.neutral_gas[i]]);
        }
        fprintf(report_file, "\nPoisonous Gases:");
        for i = 0; i < 6; i++ {
            if (i != 0) {
                fprintf(report_file, ",");
            }
            fprintf(report_file, " %s", gas_string[species.poison_gas[i]]);
        }
        fprintf(report_file, "\n");

        /* List fleet maintenance cost and its percentage of total
         * production. */
        fleet_percent_cost = species.fleet_percent_cost;

        fprintf(report_file,
                "\nFleet maintenance cost = %d (%d.%02d%% of total production)\n",
                species.fleet_cost, fleet_percent_cost / 100,
                fleet_percent_cost % 100);

        if (fleet_percent_cost > 10000) {
            fleet_percent_cost = 10000;
        }

        /* List species that have been met. */
        n              = 0;
        log_file       = report_file;   /* Use log utils for this. */
        log_stdout     = false;
        header_printed = false;
        for sp_index = 0; sp_index < galaxy.num_species; sp_index++ {
            if (!data_in_memory[sp_index]) {
                continue;
            }

            array_index = (sp_index) / 32;
            bit_number  = (sp_index) % 32;
            bit_mask    = 1 << bit_number;
            if ((species.contact[array_index] & bit_mask) == 0) {
                continue;
            }

            if (!header_printed) {
                log_string("\nSpecies met: ");
                header_printed = true;
            }

            if (n > 0) {
                log_string(", ");
            }
            log_string("SP ");  log_string(spec_data[sp_index].name);
            n++
        }
        if (n > 0) {
            log_char('\n');
        }

        /* List declared allies. */
        n = 0;
        header_printed = false;
        for sp_index = 0; sp_index < galaxy.num_species; sp_index++ {
            if (!data_in_memory[sp_index]) {
                continue;
            }

            array_index = (sp_index) / 32;
            bit_number  = (sp_index) % 32;
            bit_mask    = 1 << bit_number;
            if ((species.ally[array_index] & bit_mask) == 0) {
                continue;
            }
            if ((species.contact[array_index] & bit_mask) == 0) {
                continue;
            }

            if (!header_printed) {
                log_string("\nAllies: ");
                header_printed = true;
            }

            if (n > 0) {
                log_string(", ");
            }
            log_string("SP ");  log_string(spec_data[sp_index].name);
            n++
        }
        if (n > 0) {
            log_char('\n');
        }

        /* List declared enemies that have been met. */
        n = 0;
        header_printed = false;
        for sp_index = 0; sp_index < galaxy.num_species; sp_index++ {
            if (!data_in_memory[sp_index]) {
                continue;
            }

            array_index = (sp_index) / 32;
            bit_number  = (sp_index) % 32;
            bit_mask    = 1 << bit_number;
            if ((species.enemy[array_index] & bit_mask) == 0) {
                continue;
            }
            if ((species.contact[array_index] & bit_mask) == 0) {
                continue;
            }

            if (!header_printed) {
                log_string("\nEnemies: ");
                header_printed = true;
            }

            if (n > 0) {
                log_string(", ");
            }
            log_string("SP ");  log_string(spec_data[sp_index].name);
            n++
        }
        if (n > 0) {
            log_char('\n');
        }

        fprintf(report_file, "\nEconomic units = %d\n", species.econ_units);

        /* Initialize flag. */
        for i = 0; i < species.num_ships; i++ {
            ship_already_listed[i] = false;
        }

        /* Print report for each producing planet. */
        nampla = nampla1_base - 1;
        for i = 0; i < species.num_namplas; i++ {
            nampla++

            if (nampla.pn == 99) {
                continue;
            }
            if (nampla.mi_base == 0 && nampla.ma_base == 0 &&
                (nampla.status & HOME_PLANET) == 0) {
                continue;
            }

            planet = planet_base + nampla.planet_index;
            fprintf(report_file,
                    "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n");
            do_planet_report(nampla, ship1_base, species);
        }

        /* Give only a one-line listing for other planets. */
        printing_alien = false;
        header_printed = false;
        nampla         = nampla1_base - 1;
        for i = 0; i < species.num_namplas; i++ {
            nampla++

            if (nampla.pn == 99) {
                continue;
            }
            if (nampla.mi_base > 0 || nampla.ma_base > 0 ||
                (nampla.status & HOME_PLANET) != 0) {
                continue;
            }

            if (!header_printed) {
                fprintf(report_file,
                        "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n");
                fprintf(report_file, "\n\nOther planets and ships:\n\n");
                header_printed = true;
            }
            fprintf(report_file, "%4d%3d%3d #%d\tPL %s", nampla.x,
                    nampla.y, nampla.z, nampla.pn, nampla.name);

            for j = 0; j < MAX_ITEMS; j++ {
                if (nampla.item_quantity[j] > 0) {
                    fprintf(report_file, ", %d %s",
                            nampla.item_quantity[j], item_abbr[j]);
                }
            }
            fprintf(report_file, "\n");

            /* Print any ships at this planet. */
            ship = ship1_base - 1;
            for ship_index = 0; ship_index < species.num_ships; ship_index++ {
                ship++

                if (ship_already_listed[ship_index]) {
                    continue;
                }

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

                fprintf(report_file, "\t\t%s", ship_name(ship));
                for j = 0; j < MAX_ITEMS; j++ {
                    if (ship.item_quantity[j] > 0) {
                        fprintf(report_file, ", %d %s",
                                ship.item_quantity[j], item_abbr[j]);
                    }
                }
                fprintf(report_file, "\n");

                ship_already_listed[ship_index] = true;
            }
        }

        /* Report ships that are not associated with a planet. */
        ship = ship1_base - 1;
        for ship_index = 0; ship_index < species.num_ships; ship_index++ {
            ship++

            ship.special = 0;

            if (ship_already_listed[ship_index]) {
                continue;
            }

            ship_already_listed[ship_index] = true;

            if (ship.pn == 99) {
                continue;
            }

            if (!header_printed) {
                fprintf(report_file,
                        "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n");
                fprintf(report_file, "\n\nOther planets and ships:\n\n");
                header_printed = true;
            }

            if (ship.status == JUMPED_IN_COMBAT || ship.status == FORCED_JUMP) {
                fprintf(report_file, "  ?? ?? ??\t%s", ship_name(ship));
            }else if (test_mode && ship.arrived_via_wormhole) {
                fprintf(report_file, "  ?? ?? ??\t%s", ship_name(ship));
            }else{
                fprintf(report_file, "%4d%3d%3d\t%s",
                        ship.x, ship.y, ship.z, ship_name(ship));
            }

            for i = 0; i < MAX_ITEMS; i++ {
                if (ship.item_quantity[i] > 0) {
                    fprintf(report_file, ", %d %s",
                            ship.item_quantity[i], item_abbr[i]);
                }
            }
            fprintf(report_file, "\n");

            if (ship.status == JUMPED_IN_COMBAT ||
                ship.status == FORCED_JUMP) {
                continue;
            }

            if (test_mode && ship.arrived_via_wormhole) {
                continue;
            }

            /* Print other ships at the same location. */
            ship2 = ship;
            for i = ship_index + 1; i < species.num_ships; i++ {
                ship2++

                if (ship_already_listed[i]) {
                    continue;
                }
                if (ship2.pn == 99) {
                    continue;
                }
                if (ship2.x != ship.x) {
                    continue;
                }
                if (ship2.y != ship.y) {
                    continue;
                }
                if (ship2.z != ship.z) {
                    continue;
                }

                fprintf(report_file, "\t\t%s", ship_name(ship2));
                for j = 0; j < MAX_ITEMS; j++ {
                    if (ship2.item_quantity[j] > 0) {
                        fprintf(report_file, ", %d %s",
                                ship2.item_quantity[j], item_abbr[j]);
                    }
                }
                fprintf(report_file, "\n");

                ship_already_listed[i] = true;
            }
        }

        fprintf(report_file,
                "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n");

        /* Report aliens at locations where current species has inhabited
         * planets or ships. */
        printing_alien = true;
        locations_base = &loc[0];
        my_loc         = locations_base - 1;
        for my_loc_index = 0; my_loc_index < num_locs; my_loc_index++ {
            my_loc++
            if (my_loc.s != species_number) {
                continue;
            }

            header_printed = false;
            its_loc        = locations_base - 1;
            for its_loc_index = 0; its_loc_index < num_locs; its_loc_index++ {
                its_loc++
                if (its_loc.s == species_number) {
                    continue;
                }
                if (my_loc.x != its_loc.x) {
                    continue;
                }
                if (my_loc.y != its_loc.y) {
                    continue;
                }
                if (my_loc.z != its_loc.z) {
                    continue;
                }

                /* There is an alien here. Check if pointers for data for
                 *      this alien have been assigned yet. */
                if (its_loc.s != alien_number) {
                    alien_number = its_loc.s;
                    if (!data_in_memory[alien_number - 1]) {
                        fprintf(stderr, "\n\nWarning! Data for alien #%d is needed but is not in memory!\n\n",
                                alien_number);
                        continue;
                    }
                    alien        = &spec_data[alien_number - 1];
                    nampla2_base = namp_data[alien_number - 1];
                    ship2_base   = ship_data[alien_number - 1];
                }

                /* Check if we have a named planet in this system. If so,
                 *      use it when you print the header. */
                we_have_planet_here = false;
                nampla = nampla1_base - 1;
                for i = 0; i < species.num_namplas; i++ {
                    nampla++
                    if (nampla.x != my_loc.x) {
                        continue;
                    }
                    if (nampla.y != my_loc.y) {
                        continue;
                    }
                    if (nampla.z != my_loc.z) {
                        continue;
                    }
                    if (nampla.pn == 99) {
                        continue;
                    }

                    we_have_planet_here = true;
                    our_nampla          = nampla;

                    break;
                }

                /* Print all inhabited alien namplas at this location. */
                alien_nampla = nampla2_base - 1;
                for i = 0; i < alien.num_namplas; i++ {
                    alien_nampla++

                    if (my_loc.x != alien_nampla.x) {
                        continue;
                    }
                    if (my_loc.y != alien_nampla.y) {
                        continue;
                    }
                    if (my_loc.z != alien_nampla.z) {
                        continue;
                    }
                    if ((alien_nampla.status & POPULATED) == 0) {
                        continue;
                    }

                    /* Check if current species has a colony on the same
                     *  planet. */
                    we_have_colony_here = false;
                    nampla = nampla1_base - 1;
                    for j = 0; j < species.num_namplas; j++ {
                        nampla++

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
                        if ((nampla.status & POPULATED) == 0) {
                            continue;
                        }

                        we_have_colony_here = true;

                        break;
                    }

                    if (alien_nampla.hidden && !we_have_colony_here) {
                        continue;
                    }

                    if (!header_printed) {
                        fprintf(report_file,
                                "\n\nAliens at x = %d, y = %d, z = %d",
                                my_loc.x, my_loc.y, my_loc.z);

                        if (we_have_planet_here) {
                            fprintf(report_file, " (PL %s star system)",
                                    our_nampla.name);
                        }

                        fprintf(report_file, ":\n");
                        header_printed = true;
                    }

                    industry = alien_nampla.mi_base + alien_nampla.ma_base;

                    if (alien_nampla.status & MINING_COLONY) {
                        temp1 = fmt.Sprintf("%s", "Mining colony");
                    }else if (alien_nampla.status & RESORT_COLONY) {
                        temp1 = fmt.Sprintf("%s", "Resort colony");
                    }else if (alien_nampla.status & HOME_PLANET) {
                        temp1 = fmt.Sprintf("%s", "Home planet");
                    }else if (industry > 0) {
                        temp1 = fmt.Sprintf("%s", "Colony planet");
                    }else{
                        temp1 = fmt.Sprintf("%s", "Uncolonized planet");
                    }

                    temp2 = fmt.Sprintf("  %s PL %s (pl #%d)", temp1, alien_nampla.name, alien_nampla.pn);
                    n = 53 - strlen(temp2);
                    for j = 0; j < n; j++ {
                        strcat(temp2, " ");
                    }
                    fprintf(report_file, "%sSP %s\n", temp2, alien.name);

                    j = industry;
                    if (industry < 100) {
                        industry = (industry + 5) / 10;
                    }else{
                        industry = ((industry + 50) / 100) * 10;
                    }

                    if (j == 0) {
                        fprintf(report_file,
                                "      (No economic base.)\n");
                    }else{
                        fprintf(report_file,
                                "      (Economic base is approximately %d.)\n",
                                industry);
                    }

                    /* If current species has a colony on the same
                     *  planet, report any PDs and any shipyards. */
                    if (we_have_colony_here) {
                        if (alien_nampla.item_quantity[PD] == 1) {
                            fprintf(report_file,
                                    "      (There is 1 %s on the planet.)\n",
                                    item_name[PD]);
                        }else if (alien_nampla.item_quantity[PD] > 1) {
                            fprintf(report_file,
                                    "      (There are %d %ss on the planet.)\n",
                                    alien_nampla.item_quantity[PD],
                                    item_name[PD]);
                        }

                        if (alien_nampla.shipyards == 1) {
                            fprintf(report_file,
                                    "      (There is 1 shipyard on the planet.)\n");
                        }else if (alien_nampla.shipyards > 1) {
                            fprintf(report_file,
                                    "      (There are %d shipyards on the planet.)\n",
                                    alien_nampla.shipyards);
                        }
                    }

                    /* Also report if alien colony is actively hiding. */
                    if (alien_nampla.hidden) {
                        fprintf(report_file,
                                "      (Colony is actively hiding from alien observation.)\n");
                    }
                }

                /* Print all alien ships at this location. */
                alien_ship = ship2_base - 1;
                for i = 0; i < alien.num_ships; i++ {
                    alien_ship++

                    if (alien_ship.pn == 99) {
                        continue;
                    }
                    if (my_loc.x != alien_ship.x) {
                        continue;
                    }
                    if (my_loc.y != alien_ship.y) {
                        continue;
                    }
                    if (my_loc.z != alien_ship.z) {
                        continue;
                    }

                    /* An alien ship cannot hide if it lands on the
                     *  surface of a planet populated by the current
                     *  species. */
                    alien_can_hide = true;
                    nampla         = nampla1_base - 1;
                    for j = 0; j < species.num_namplas; j++ {
                        nampla++

                        if (alien_ship.x != nampla.x) {
                            continue;
                        }
                        if (alien_ship.y != nampla.y) {
                            continue;
                        }
                        if (alien_ship.z != nampla.z) {
                            continue;
                        }
                        if (alien_ship.pn != nampla.pn) {
                            continue;
                        }
                        if (nampla.status & POPULATED) {
                            alien_can_hide = false;
                            break;
                        }
                    }

                    if (alien_can_hide && alien_ship.status == ON_SURFACE) {
                        continue;
                    }

                    if (alien_can_hide && alien_ship.status == UNDER_CONSTRUCTION) {
                        continue;
                    }

                    if (!header_printed) {
                        fprintf(report_file,
                                "\n\nAliens at x = %d, y = %d, z = %d",
                                my_loc.x, my_loc.y, my_loc.z);

                        if (we_have_planet_here) {
                            fprintf(report_file, " (PL %s star system)",
                                    our_nampla.name);
                        }

                        fprintf(report_file, ":\n");
                        header_printed = true;
                    }

                    print_ship(alien_ship, alien, alien_number);
                }
            }
        }

        printing_alien = false;

        if (test_mode) {
            goto done_report;
        }

        /* Generate order section. */
        truncate_name = true;
        temp_ignore_field_distorters = ignore_field_distorters;
        ignore_field_distorters      = true;

        fprintf(report_file,
                "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n");

        fprintf(report_file,
                "\n\nORDER SECTION. Remove these two lines and everything above\n");
        fprintf(report_file, "  them, and submit only the orders below.\n\n");

        fprintf(report_file, "START COMBAT\n");
        fprintf(report_file, "; Place combat orders here.\n\n");
        fprintf(report_file, "END\n\n");

        fprintf(report_file, "START PRE-DEPARTURE\n");
        fprintf(report_file, "; Place pre-departure orders here.\n\n");

        for nampla_index = 0; nampla_index < species.num_namplas; nampla_index++ {
            nampla = nampla_base + nampla_index;
            if (nampla.pn == 99) {
                continue;
            }

            /* Generate auto-installs for colonies that were loaded via
             *  the DEVELOP command. */
            if (nampla.auto_IUs) {
                fprintf(report_file, "\tInstall\t%d IU\tPL %s\n",
                        nampla.auto_IUs, nampla.name);
            }
            if (nampla.auto_AUs) {
                fprintf(report_file, "\tInstall\t%d AU\tPL %s\n",
                        nampla.auto_AUs, nampla.name);
            }
            if (nampla.auto_IUs || nampla.auto_AUs) {
                fprintf(report_file, "\n");
            }

            if (!species.auto_orders) {
                continue;
            }

            /* Generate auto UNLOAD orders for transports at this nampla. */
            for j = 0; j < species.num_ships; j++ {
                ship = ship_base + j;
                if (ship.pn == 99) {
                    continue;
                }
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
                if (ship.status == JUMPED_IN_COMBAT) {
                    continue;
                }
                if (ship.status == FORCED_JUMP) {
                    continue;
                }
                if (ship.class != TR) {
                    continue;
                }
                if (ship.item_quantity[CU] < 1) {
                    continue;
                }

                /* New colonies will never be started automatically unless
                 * ship was loaded via a DEVELOP order. */
                if (ship.loading_point != 0) {
                    /* Check if transport is at specified unloading point. */
                    n = ship.unloading_point;
                    if (n == nampla_index ||
                        (n == 9999 && nampla_index == 0)) {
                        goto unload_ship;
                    }
                }

                if ((nampla.status & POPULATED) == 0) {
                    continue;
                }

                if ((nampla.mi_base + nampla.ma_base) >= 2000) {
                    continue;
                }

                if (nampla.x == nampla_base.x &&
                    nampla.y == nampla_base.y &&
                    nampla.z == nampla_base.z) {
                    continue;                                  /* Home sector. */
                }
unload_ship:

                n = ship.loading_point;
                if (n == 9999) {
                    n = 0;              /* Home planet. */
                }
                if (n == nampla_index) {
                    continue;   /* Ship was just loaded here. */
                }
                fprintf(report_file, "\tUnload\tTR%d%s %s\n\n", ship.tonnage,
                        ship_ttype[ship.ttype], ship.name);

                ship.special = ship.loading_point;
                n             = nampla - nampla_base;
                if (n == 0) {
                    n = 9999;
                }
                ship.unloading_point = n;
            }
        }

        fprintf(report_file, "END\n\n");

        fprintf(report_file, "START JUMPS\n");
        fprintf(report_file, "; Place jump orders here.\n\n");

        /* Generate auto-jumps for ships that were loaded via the DEVELOP
         * command or which were UNLOADed because of the AUTO command. */
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base + i;

            ship.just_jumped = false;

            if (ship.pn == 99) {
                continue;
            }
            if (ship.status == JUMPED_IN_COMBAT) {
                continue;
            }
            if (ship.status == FORCED_JUMP) {
                continue;
            }

            j = ship.special;
            if (j) {
                if (j == 9999) {
                    j = 0;              /* Home planet. */
                }
                temp_nampla = nampla_base + j;

                fprintf(report_file, "\tJump\t%s, PL %s\t; Age %d, ",
                        ship_name(ship), temp_nampla.name, ship.age);

                print_mishap_chance(ship, temp_nampla.x, temp_nampla.y,
                                    temp_nampla.z);

                fprintf(report_file, "\n\n");

                ship.just_jumped = true;

                continue;
            }

            n = ship.unloading_point;
            if (n) {
                if (n == 9999) {
                    n = 0;              /* Home planet. */
                }
                temp_nampla = nampla_base + n;

                fprintf(report_file, "\tJump\t%s, PL %s\t; ", ship_name(ship),
                        temp_nampla.name);

                print_mishap_chance(ship, temp_nampla.x, temp_nampla.y,
                                    temp_nampla.z);

                fprintf(report_file, "\n\n");

                ship.just_jumped = true;
            }
        }

        if (!species.auto_orders) {
            goto jump_end;
        }

        /* Generate JUMP orders for all ships that have not yet been
         *      given orders. */
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base + i;
            if (ship.pn == 99) {
                continue;
            }
            if (ship.just_jumped) {
                continue;
            }
            if (ship.status == UNDER_CONSTRUCTION) {
                continue;
            }
            if (ship.status == JUMPED_IN_COMBAT) {
                continue;
            }
            if (ship.status == FORCED_JUMP) {
                continue;
            }

            if (ship.ttype == FTL) {
                fprintf(report_file, "\tJump\t%s, ", ship_name(ship));
                if (ship.class == TR && ship.tonnage == 1) {
                    closest_unvisited_star(ship);
                    fprintf(report_file,
                            "\n\t\t\t; Age %d, now at %d %d %d, ",
                            ship.age, ship.x, ship.y, ship.z);

                    if (ship.status == IN_ORBIT) {
                        fprintf(report_file, "O%d, ", ship.pn);
                    }else if (ship.status == ON_SURFACE) {
                        fprintf(report_file, "L%d, ", ship.pn);
                    }else{
                        fprintf(report_file, "D, ");
                    }

                    print_mishap_chance(ship, x, y, z);
                }else {
                    fprintf(report_file,
                            "???\t; Age %d, now at %d %d %d",
                            ship.age, ship.x, ship.y, ship.z);

                    if (ship.status == IN_ORBIT) {
                        fprintf(report_file, ", O%d", ship.pn);
                    }else if (ship.status == ON_SURFACE) {
                        fprintf(report_file, ", L%d", ship.pn);
                    }else{
                        fprintf(report_file, ", D");
                    }

                    x = 9999;
                }

                fprintf(report_file, "\n");

                /* Save destination so that we can check later if it needs
                 * to be scanned. */
                if (x == 9999) {
                    ship.dest_x = -1;
                }else{
                    ship.dest_x = x;
                    ship.dest_y = y;
                    ship.dest_z = z;
                }
            }
        }

jump_end:
        fprintf(report_file, "END\n\n");

        fprintf(report_file, "START PRODUCTION\n\n");

        fprintf(report_file, ";   Economic units at start of turn = %d\n\n",
                species.econ_units);

        /* Generate a PRODUCTION order for each planet that can produce. */
        for nampla_index = species.num_namplas - 1; nampla_index >= 0;             nampla_index-- {
            nampla = nampla1_base + nampla_index;
            if (nampla.pn == 99) {
                continue;
            }

            if (nampla.mi_base == 0 && (nampla.status & RESORT_COLONY) == 0) {
                continue;
            }
            if (nampla.ma_base == 0 && (nampla.status & MINING_COLONY) == 0) {
                continue;
            }

            fprintf(report_file, "    PRODUCTION PL %s\n", nampla.name);

            if (nampla.status & MINING_COLONY) {
                fprintf(report_file,
                        "    ; The above PRODUCTION order is required for this mining colony, even\n");
                fprintf(report_file,
                        "    ;  if no other production orders are given for it. This mining colony\n");
                fprintf(report_file,
                        "    ;  will generate %d economic units this turn.\n", nampla.use_on_ambush);
            }else if (nampla.status & RESORT_COLONY) {
                fprintf(report_file,
                        "    ; The above PRODUCTION order is required for this resort colony, even\n");
                fprintf(report_file,
                        "    ;  though no other production orders can be given for it.  This resort\n");
                fprintf(report_file,
                        "    ;  colony will generate %d economic units this turn.\n", nampla.use_on_ambush);
            }else {
                fprintf(report_file,
                        "    ; Place production orders here for planet %s",
                        nampla.name);
                fprintf(report_file, " (sector %d %d %d #%d).\n", nampla.x,
                        nampla.y, nampla.z, nampla.pn);
                fprintf(report_file,
                        "    ;  Avail pop = %d, shipyards = %d, to spend = %d",
                        nampla.pop_units, nampla.shipyards, nampla.use_on_ambush);

                n = nampla.use_on_ambush;
                if (nampla.status & HOME_PLANET) {
                    if (species.hp_original_base != 0) {
                        fprintf(report_file, " (max = %d)", 5 * n);
                    }else{
                        fprintf(report_file, " (max = no limit)");
                    }
                }else {
                    fprintf(report_file, " (max = %d)", 2 * n);
                }

                fprintf(report_file, ".\n\n");
            }

            /* Build IUs and AUs for incoming ships with CUs. */
            if (nampla.IUs_needed) {
                fprintf(report_file, "\tBuild\t%d IU\n", nampla.IUs_needed);
            }
            if (nampla.AUs_needed) {
                fprintf(report_file, "\tBuild\t%d AU\n", nampla.AUs_needed);
            }
            if (nampla.IUs_needed || nampla.AUs_needed) {
                fprintf(report_file, "\n");
            }

            if (!species.auto_orders) {
                continue;
            }
            if (nampla.status & MINING_COLONY) {
                continue;
            }
            if (nampla.status & RESORT_COLONY) {
                continue;
            }

            /* See if there are any RMs to recycle. */
            n = nampla.special / 5;
            if (n > 0) {
                fprintf(report_file, "\tRecycle\t%d RM\n\n", 5 * n);
            }

            /* Generate DEVELOP commands for ships arriving here because of
             *  AUTO command. */
            for i = 0; i < species.num_ships; i++ {
                ship = ship_base + i;
                if (ship.pn == 99) {
                    continue;
                }

                k = ship.special;
                if (k == 0) {
                    continue;
                }
                if (k == 9999) {
                    k = 0;              /* Home planet. */
                }
                if (nampla != nampla_base + k) {
                    continue;
                }

                k = ship.unloading_point;
                if (k == 9999) {
                    k = 0;
                }
                temp_nampla = nampla_base + k;

                fprintf(report_file, "\tDevelop\tPL %s, TR%d%s %s\n\n",
                        temp_nampla.name, ship.tonnage, ship_ttype[ship.ttype],
                        ship.name);
            }

            /* Give orders to continue construction of unfinished ships and
             *  starbases. */
            for i = 0; i < species.num_ships; i++ {
                ship = ship_base + i;
                if (ship.pn == 99) {
                    continue;
                }

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

                if (ship.status == UNDER_CONSTRUCTION) {
                    fprintf(report_file,
                            "\tContinue\t%s, %d\t; Left to pay = %d\n\n",
                            ship_name(ship), ship.remaining_cost,
                            ship.remaining_cost);

                    continue;
                }

                if (ship.ttype != STARBASE) {
                    continue;
                }

                j = (species.tech_level[MA] / 2) - ship.tonnage;
                if (j < 1) {
                    continue;
                }

                fprintf(report_file,
                        "\tContinue\tBAS %s, %d\t; Current tonnage = %s\n\n",
                        ship.name, 100 * j, commas(10000 * ship.tonnage));
            }

            /* Generate DEVELOP command if this is a colony with an economic
             *  base less than 200. */
            n = nampla.mi_base + nampla.ma_base + nampla.IUs_needed
                + nampla.AUs_needed;
            nn = nampla.item_quantity[CU];
            for i = 0; i < species.num_ships; i++ {
                /* Get CUs on transports at planet. */
                ship = ship_base + i;
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
                nn += ship.item_quantity[CU];
            }
            n += nn;
            if ((nampla.status & COLONY) && n < 2000 &&
                nampla.pop_units > 0) {
                if (nampla.pop_units > (2000 - n)) {
                    nn = 2000 - n;
                }else{
                    nn = nampla.pop_units;
                }

                fprintf(report_file, "\tDevelop\t%d\n\n", 2 * nn);

                nampla.IUs_needed += nn;
            }

            /* For home planets and any colonies that have an economic base of
             *  at least 200, check if there are other colonized planets in
             *  the same sector that are not self-sufficient.  If so, DEVELOP
             *  them. */
            if (n >= 2000 || (nampla.status & HOME_PLANET)) {
                /* Skip home planet. */
                for i = 1; i < species.num_namplas; i++ {
                    if (i == nampla_index) {
                        continue;
                    }

                    temp_nampla = nampla_base + i;

                    if (temp_nampla.pn == 99) {
                        continue;
                    }
                    if (temp_nampla.x != nampla.x) {
                        continue;
                    }
                    if (temp_nampla.y != nampla.y) {
                        continue;
                    }
                    if (temp_nampla.z != nampla.z) {
                        continue;
                    }

                    n = temp_nampla.mi_base + temp_nampla.ma_base
                        + temp_nampla.IUs_needed + temp_nampla.AUs_needed;

                    if (n == 0) {
                        continue;
                    }

                    nn = temp_nampla.item_quantity[IU]
                         + temp_nampla.item_quantity[AU];
                    if (nn > temp_nampla.item_quantity[CU]) {
                        nn = temp_nampla.item_quantity[CU];
                    }
                    n += nn;
                    if (n >= 2000) {
                        continue;
                    }
                    nn = 2000 - n;

                    if (nn > nampla.pop_units) {
                        nn = nampla.pop_units;
                    }

                    fprintf(report_file, "\tDevelop\t%d\tPL %s\n\n",
                            2 * nn, temp_nampla.name);

                    temp_nampla.AUs_needed += nn;
                }
            }
        }

        fprintf(report_file, "END\n\n");

        fprintf(report_file, "START POST-ARRIVAL\n");
        fprintf(report_file, "; Place post-arrival orders here.\n\n");

        if (!species.auto_orders) {
            goto post_end;
        }

        /* Generate an AUTO command. */
        fprintf(report_file, "\tAuto\n\n");

        /* Generate SCAN orders for all TR1s that are jumping to
         * sectors which current species does not inhabit. */
        for i = 0; i < species.num_ships; i++ {
            ship = ship_base + i;
            if (ship.pn == 99) {
                continue;
            }
            if (ship.status == UNDER_CONSTRUCTION) {
                continue;
            }
            if (ship.class != TR) {
                continue;
            }
            if (ship.tonnage != 1) {
                continue;
            }
            if (ship.ttype != FTL) {
                continue;
            }

            found = false;
            for j = 0; j < species.num_namplas; j++ {
                if (ship.dest_x == -1) {
                    break;
                }

                nampla = nampla_base + j;
                if (nampla.pn == 99) {
                    continue;
                }
                if (nampla.x != ship.dest_x) {
                    continue;
                }
                if (nampla.y != ship.dest_y) {
                    continue;
                }
                if (nampla.z != ship.dest_z) {
                    continue;
                }

                if (nampla.status & POPULATED) {
                    found = true;
                    break;
                }
            }
            if (!found) {
                fprintf(report_file, "\tScan\tTR1 %s\n", ship.name);
            }
        }

post_end:
        fprintf(report_file, "END\n\n");

        fprintf(report_file, "START STRIKES\n");
        fprintf(report_file, "; Place strike orders here.\n\n");
        fprintf(report_file, "END\n");

        truncate_name           = false;
        ignore_field_distorters = temp_ignore_field_distorters;

done_report:

        /* Clean up for this species. */
        fclose(report_file);
    }

    /* Clean up and exit. */
    free_species_data();
    exit(0);
}



do_planet_report(nampla, s_base, species)

struct species_data *species;
struct nampla_data * nampla;
struct ship_data *   s_base;

{
    int i, j, ship_index, header_printed, ls_needed, production_penalty;

    long n1, n2, n3, raw_material_units, production_capacity,
         available_to_spend, n, ib, ab, current_base, md, denom;

    struct ship_data *ship;


    /* Print ttype of planet, name and coordinates. */
    fprintf(report_file, "\n\n");

    if (nampla.status & HOME_PLANET) {
        fprintf(report_file, "HOME PLANET");
    }else if (nampla.status & MINING_COLONY) {
        fprintf(report_file, "MINING COLONY");
    }else if (nampla.status & RESORT_COLONY) {
        fprintf(report_file, "RESORT COLONY");
    }else if (nampla.status & POPULATED) {
        fprintf(report_file, "COLONY PLANET");
    }else{
        fprintf(report_file, "PLANET");
    }

    fprintf(report_file, ": PL %s", nampla.name);

    fprintf(report_file,
            "\n   Coordinates: x = %d, y = %d, z = %d, planet number %d\n",
            nampla.x, nampla.y, nampla.z, nampla.pn);

    if (nampla.status & HOME_PLANET) {
        ib           = nampla.mi_base;
        ab           = nampla.ma_base;
        current_base = ib + ab;
        if (current_base < species.hp_original_base) {
            n = species.hp_original_base - current_base; /* Number of CUs needed. */

            md = home_planet.mining_difficulty;

            denom = 100 + md;
            j     = (100 * (n + ib) - (md * ab) + denom / 2) / denom;
            i     = n - j;

            if (i < 0) {
                j = n;
                i = 0;
            }
            if (j < 0) {
                i = n;
                j = 0;
            }

            fprintf(report_file,
                    "\nWARNING! Home planet has not yet completely recovered from bombardment!\n");
            fprintf(report_file,
                    "         %d IUs and %d AUs will have to be installed for complete recovery.\n",
                    i, j);
        }
    }

    if (!(nampla.status & POPULATED)) {
        goto do_inventory;
    }

    /* Print available population. */
    if (nampla.status & (MINING_COLONY | RESORT_COLONY)) {
        ;
    }else{
        fprintf(report_file, "\nAvailable population units = %d\n",
                nampla.pop_units);
    }

    if (nampla.siege_eff != 0) {
        fprintf(report_file,
                "\nWARNING!  This planet is currently under siege and will remain\n");
        fprintf(report_file,
                "  under siege until the combat phase of the next turn!\n");
    }

    if (nampla.use_on_ambush > 0) {
        fprintf(report_file,
                "\nIMPORTANT!  This planet has made preparations for an ambush!\n");
    }

    if (nampla.hidden) {
        fprintf(report_file,
                "\nIMPORTANT!  This planet is actively hiding from alien observation!\n");
    }

    /* Print what will be produced this turn. */
    raw_material_units =
        (10 * species.tech_level[MI] * nampla.mi_base)
        / planet.mining_difficulty;
    production_capacity =
        (species.tech_level[MA] * nampla.ma_base) / 10;

    ls_needed = life_support_needed(species, home_planet, planet);

    if (ls_needed == 0) {
        production_penalty = 0;
    }else{
        production_penalty = (100 * ls_needed) / species.tech_level[LS];
    }

    fprintf(report_file, "\nProduction penalty = %d%% (LSN = %d)\n",
            production_penalty, ls_needed);

    fprintf(report_file, "\nEconomic efficiency = %d%%\n",
            planet.econ_efficiency);

    raw_material_units
        -= (production_penalty * raw_material_units) / 100;

    raw_material_units
        = ((planet.econ_efficiency * raw_material_units) + 50) / 100;

    production_capacity
        -= (production_penalty * production_capacity) / 100;

    production_capacity
        = ((planet.econ_efficiency * production_capacity) + 50) / 100;

    if (nampla.mi_base > 0) {
        fprintf(report_file, "\nMining base = %d.%d", nampla.mi_base / 10,
                nampla.mi_base % 10);
        fprintf(report_file, " (MI = %d, MD = %d.%02d)\n",
                species.tech_level[MI], planet.mining_difficulty / 100,
                planet.mining_difficulty % 100);

        /* For mining colonies, print economic units that will be produced. */
        if (nampla.status & MINING_COLONY) {
            n1 = (2 * raw_material_units) / 3;
            n2 = ((fleet_percent_cost * n1) + 5000) / 10000;
            n3 = n1 - n2;
            fprintf(report_file,
                    "   This mining colony will generate %d - %d = %d economic units this turn.\n",
                    n1, n2, n3);

            nampla.use_on_ambush = n3;         /* Temporary use only. */
        }else {
            fprintf(report_file,
                    "   %d raw material units will be produced this turn.\n",
                    raw_material_units);
        }
    }

    if (nampla.ma_base > 0) {
        if (nampla.status & RESORT_COLONY) {
            fprintf(report_file, "\n");
        }

        fprintf(report_file, "Manufacturing base = %d.%d",
                nampla.ma_base / 10, nampla.ma_base % 10);
        fprintf(report_file, " (MA = %d)\n", species.tech_level[MA]);

        /* For resort colonies, print economic units that will be produced. */
        if (nampla.status & RESORT_COLONY) {
            n1 = (2 * production_capacity) / 3;
            n2 = ((fleet_percent_cost * n1) + 5000) / 10000;
            n3 = n1 - n2;
            fprintf(report_file,
                    "   This resort colony will generate %d - %d = %d economic units this turn.\n",
                    n1, n2, n3);

            nampla.use_on_ambush = n3;         /* Temporary use only. */
        }else {
            fprintf(report_file,
                    "   Production capacity this turn will be %d.\n",
                    production_capacity);
        }
    }

    if (nampla.item_quantity[RM] > 0) {
        fprintf(report_file, "\n%ss (%s,C%d) carried over from last turn = %d\n",
                item_name[RM], item_abbr[RM], item_carry_capacity[RM],
                nampla.item_quantity[RM]);
    }

    /* Print what can be spent this turn. */
    raw_material_units += nampla.item_quantity[RM];
    if (raw_material_units > production_capacity) {
        available_to_spend = production_capacity;
        nampla.special    = raw_material_units - production_capacity;
        /* Excess raw material units that may be recycled in AUTO mode. */
    }else {
        available_to_spend = raw_material_units;
        nampla.special    = 0;
    }

    /* Don't print spendable amount for mining and resort colonies. */
    n1 = available_to_spend;
    n2 = ((fleet_percent_cost * n1) + 5000) / 10000;
    n3 = n1 - n2;
    if (!(nampla.status & MINING_COLONY) &&
        !(nampla.status & RESORT_COLONY)) {
        fprintf(report_file,
                "\nTotal available for spending this turn = %d - %d = %d\n",
                n1, n2, n3);
        nampla.use_on_ambush = n3;     /* Temporary use only. */

        fprintf(report_file,
                "\nShipyard capacity = %d\n", nampla.shipyards);
    }

do_inventory:

    header_printed = false;

    for i = 0; i < MAX_ITEMS; i++ {
        if (nampla.item_quantity[i] > 0 && i != RM) {
            if (!header_printed) {
                header_printed = true;
                fprintf(report_file, "\nPlanetary inventory:\n");
            }

            fprintf(report_file, "   %ss (%s,C%d) = %d",
                    item_name[i], item_abbr[i],
                    item_carry_capacity[i], nampla.item_quantity[i]);
            if (i == PD) {
                fprintf(report_file, " (warship equivalence = %d tons)",
                        50 * nampla.item_quantity[PD]);
            }
            fprintf(report_file, "\n");
        }
    }

    /* Print all ships that are under construction on, on the surface of,
     *  or in orbit around this planet. */
    printing_alien = false;
    header_printed = false;
    for ship_index = 0; ship_index < species.num_ships; ship_index++ {
        ship = s_base + ship_index;

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
        if (ship.class != BA) {
            continue;
        }

        if (!header_printed) {
            fprintf(report_file, "\nShips at PL %s:\n", nampla.name);
            print_ship_header();
        }
        header_printed = true;

        print_ship(ship, species, species_number);

        ship_already_listed[ship_index] = true;
    }

    for ship_index = 0; ship_index < species.num_ships; ship_index++ {
        ship = s_base + ship_index;

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
        if (ship.class != TR) {
            continue;
        }

        if (!header_printed) {
            fprintf(report_file, "\nShips at PL %s:\n", nampla.name);
            print_ship_header();
        }
        header_printed = true;

        print_ship(ship, species, species_number);

        ship_already_listed[ship_index] = true;
    }

    for ship_index = 0; ship_index < species.num_ships; ship_index++ {
        ship = s_base + ship_index;

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
        if (ship_already_listed[ship_index]) {
            continue;
        }

        if (!header_printed) {
            fprintf(report_file, "\nShips at PL %s:\n", nampla.name);
            print_ship_header();
        }
        header_printed = true;

        print_ship(ship, species, species_number);

        ship_already_listed[ship_index] = true;
    }
}



print_ship_header() {
    fprintf(report_file, "  Name                          ");
    if (printing_alien) {
        fprintf(report_file, "                     Species\n");
    }else{
        fprintf(report_file, "                 Cap. Cargo\n");
    }
    fprintf(report_file, " ---------------------------------------");
    fprintf(report_file, "-------------------------------------\n");
}


extern char full_ship_id[64];

print_ship(ship, species, species_number)

struct species_data *species;
struct ship_data *   ship;
int species_number;

{
    int i, n, length, capacity, need_comma;


    if (printing_alien) {
        ignore_field_distorters = false;
    }else{
        ignore_field_distorters = true;
    }

    fprintf(report_file, "  %s", ship_name(ship));

    length = strlen(full_ship_id);
    if (printing_alien) {
        n = 50;
    }else{
        n = 46;
    }

    for i = 0; i < (n - length); i++ {
        putc(' ', report_file);
    }

    if (ship.class == BA) {
        capacity = 10 * (int)ship.tonnage;
    }else if (ship.class == TR) {
        capacity = (10 + ((int)ship.tonnage / 2)) * (int)ship.tonnage;
    }else{
        capacity = ship.tonnage;
    }

    if (printing_alien) {
        fprintf(report_file, " ");
    }else{
        fprintf(report_file, "%4d  ", capacity);
        if (ship.status == UNDER_CONSTRUCTION) {
            fprintf(report_file, "Left to pay = %d\n", ship.remaining_cost);
            return;
        }
    }

    if (printing_alien) {
        if (ship.status == ON_SURFACE ||
            ship.item_quantity[FD] != ship.tonnage) {
            fprintf(report_file, "SP %s", species.name);
        }else{
            fprintf(report_file, "SP %d", distorted(species_number));
        }
    }else {
        need_comma = false;
        for i = 0; i < MAX_ITEMS; i++ {
            if (ship.item_quantity[i] > 0) {
                if (need_comma) {
                    putc(',', report_file);
                }
                fprintf(report_file, "%d %s",
                        ship.item_quantity[i], item_abbr[i]);
                need_comma = true;
            }
        }
    }

    putc('\n', report_file);
}


print_mishap_chance(ship, destx, desty, destz)

struct ship_data *ship;
int destx, desty, destz;

{
    int mishap_GV, mishap_age;

    long x, y, z, mishap_chance, success_chance;


    if (destx == 9999) {
        fprintf(report_file, "Mishap chance = ???");
        return;
    }

    mishap_GV  = species.tech_level[GV];
    mishap_age = ship.age;

    x             = destx;
    y             = desty;
    z             = destz;
    mishap_chance = (100 * (
                         ((x - ship.x) * (x - ship.x))
                         + ((y - ship.y) * (y - ship.y))
                         + ((z - ship.z) * (z - ship.z))
                         )) / mishap_GV;

    if (mishap_age > 0 && mishap_chance < 10000) {
        success_chance  = 10000 - mishap_chance;
        success_chance -= (2 * mishap_age * success_chance) / 100;
        mishap_chance   = 10000 - success_chance;
    }

    if (mishap_chance > 10000) {
        mishap_chance = 10000;
    }

    fprintf(report_file, "mishap chance = %d.%02d%%",
            mishap_chance / 100, mishap_chance % 100);
}


closest_unvisited_star(ship)

struct ship_data *ship;

{
    int i, found, species_array_index, species_bit_number;

    long shx, shy, shz, stx, sty, stz, closest_distance, temp_distance,
         species_bit_mask;

    struct star_data *star, *closest_star;


    /* Get array index and bit mask. */
    species_array_index = (species_number - 1) / 32;
    species_bit_number  = (species_number - 1) % 32;
    species_bit_mask    = 1 << species_bit_number;

    shx = ship.x;
    shy = ship.y;
    shz = ship.z;

    x = 9999;
    closest_distance = 999999;

    found = false;
    for i = 0; i < num_stars; i++ {
        star = star_base + i;

        /* Check if bit is already set. */
        if (star.visited_by[species_array_index] & species_bit_mask) {
            continue;
        }

        stx = star.x;
        sty = star.y;
        stz = star.z;

        temp_distance =
            ((shx - stx) * (shx - stx))
            + ((shy - sty) * (shy - sty))
            + ((shz - stz) * (shz - stz));

        if (temp_distance < closest_distance) {
            x = stx;
            y = sty;
            z = stz;
            closest_distance = temp_distance;
            closest_star     = star;
            found            = true;
        }
    }

    if (found) {
        fprintf(report_file, "%d %d %d", x, y, z);
        closest_star.visited_by[species_array_index] |= species_bit_mask;
        /* So that we don't send more than one ship to the same place. */
    }else {
        fprintf(report_file, "???");
    }
}



//*************************************************************************
// Stats.c

func StatsMain(argc int, argv []string) {
    int i, j, m, n, nampla_index, ship_index, num_ships, avg_tech_level,
        all_tech_level[6], n_species, n_warships, n_starbases,
        n_transports, avg_pop_pl, n_pop_pl, ls_needed, num_yards,
        production_penalty, fleet_percent_cost, num_pop_planets,
        min_starbases, max_starbases, min_warships, max_warships,
        min_transports, max_transports, min_tech_level[6],
        max_tech_level[6], min_pop_pl, max_pop_pl, ntr, nba, nwa,
        n_yards, min_yards, max_yards, avg_yards;

    long total_production, raw_material_units, production_capacity,
         total_tonnage, total_offensive_power, total_defensive_power,
         avg_production, all_production, avg_warship_tons,
         all_warship_tons, avg_starbase_tons, all_starbase_tons,
         avg_transport_tons, all_transport_tons, n1, n2, n3,
         min_production, max_production;

    short tons;


    /* Check for valid command line. */
    if (argc != 1) {
        fprintf(stderr, "\n\tUsage: Stats\n\n");
        exit(0);
    }

    /* Get all necessary data. */
    get_galaxy_data();
    get_planet_data();
    get_species_data();

    /* Initialize data. */
    n_species          = 0;
    all_production     = 0;
    min_production     = 1000000000;
    max_production     = 0;
    all_warship_tons   = 0;
    all_starbase_tons  = 0;
    all_transport_tons = 0;
    n_warships         = 0;
    min_warships       = 32000;
    max_warships       = 0;
    n_starbases        = 0;
    min_starbases      = 32000;
    max_starbases      = 0;
    n_transports       = 0;
    min_transports     = 32000;
    max_transports     = 0;
    n_pop_pl           = 0;
    min_pop_pl         = 32000;
    max_pop_pl         = 0;
    n_yards            = 0;
    min_yards          = 32000;
    max_yards          = 0;
    for i = 0; i < 6; i++ {
        all_tech_level[i] = 0;
        min_tech_level[i] = 32000;
        max_tech_level[i] = 0;
    }

    /* Print header. */
    printf("SP Species               Tech Levels        Total  Num Num  Num  Offen.  Defen.\n");
    printf(" # Name             MI  MA  ML  GV  LS  BI  Prod.  Pls Shps Yrds  Power   Power\n");
    printf("-------------------------------------------------------------------------------\n");

    /* Main loop. For each species, take appropriate action. */
    for species_number = 1; species_number <= galaxy.num_species; species_number++ {
        if (!data_in_memory[species_number - 1]) {
            continue;
        }

        n_species++

        species     = &spec_data[species_number - 1];
        nampla_base = namp_data[species_number - 1];
        ship_base   = ship_data[species_number - 1];

        /* Get fleet maintenance cost. */
        fleet_percent_cost = species.fleet_percent_cost;

        if (fleet_percent_cost > 10000) {
            fleet_percent_cost = 10000;
        }

        /* Print species data. */
        printf("%2d", species_number);
        printf(" %-15.15s", species.name);

        for i = 0; i < 6; i++ {
            printf("%4d", species.tech_level[i]);
            all_tech_level[i] += (int)species.tech_level[i];
            if (species.tech_level[i] < min_tech_level[i]) {
                min_tech_level[i] = species.tech_level[i];
            }
            if (species.tech_level[i] > max_tech_level[i]) {
                max_tech_level[i] = species.tech_level[i];
            }
        }

        /* Get stats for namplas. */
        total_production      = 0;
        total_defensive_power = 0;
        num_yards             = 0;
        num_pop_planets       = 0;
        home_planet           = planet_base + (int)nampla_base.planet_index;
        nampla = nampla_base - 1;
        for nampla_index = 0; nampla_index < species.num_namplas; nampla_index++ {
            nampla++

            if (nampla.pn == 99) {
                continue;
            }

            num_yards += nampla.shipyards;
            n_yards   += nampla.shipyards;

            planet = planet_base + (int)nampla.planet_index;

            raw_material_units =
                (10 * species.tech_level[MI] * nampla.mi_base)
                / planet.mining_difficulty;

            production_capacity =
                (species.tech_level[MA] * nampla.ma_base) / 10;

            ls_needed = life_support_needed(species, home_planet, planet);

            if (ls_needed == 0) {
                production_penalty = 0;
            }else{
                production_penalty = (100 * ls_needed) / species.tech_level[LS];
            }

            raw_material_units
                -= (production_penalty * raw_material_units) / 100;

            raw_material_units
                = ((planet.econ_efficiency * raw_material_units) + 50) / 100;

            production_capacity
                -= (production_penalty * production_capacity) / 100;

            production_capacity
                = ((planet.econ_efficiency * production_capacity) + 50) / 100;

            if (nampla.status & MINING_COLONY) {
                n1 = (2 * raw_material_units) / 3;
            }else if (nampla.status & RESORT_COLONY) {
                n1 = (2 * production_capacity) / 3;
            }else{
                n1 = (production_capacity > raw_material_units)
                     ? raw_material_units : production_capacity;
            }

            n2 = ((fleet_percent_cost * n1) + 5000) / 10000;
            n3 = n1 - n2;
            total_production += n3;

            tons = nampla.item_quantity[PD] / 200;
            if (tons < 1 && nampla.item_quantity[PD] > 0) {
                tons = 1;
            }
            total_defensive_power += power(tons);

            if (nampla.status & POPULATED) {
                n_pop_pl++
                num_pop_planets++
            }
        }

        printf("%7ld%4d", total_production, num_pop_planets);

        if (total_production < min_production) {
            min_production = total_production;
        }
        if (total_production > max_production) {
            max_production = total_production;
        }

        if (num_pop_planets < min_pop_pl) {
            min_pop_pl = num_pop_planets;
        }
        if (num_pop_planets > max_pop_pl) {
            max_pop_pl = num_pop_planets;
        }

        if (num_yards < min_yards) {
            min_yards = num_yards;
        }
        if (num_yards > max_yards) {
            max_yards = num_yards;
        }

        all_production += total_production;

        /* Get stats for ships. */
        num_ships             = 0;
        ntr                   = 0;  nba = 0;  nwa = 0;
        total_tonnage         = 0;
        total_offensive_power = 0;
        ship                  = ship_base - 1;
        for ship_index = 0; ship_index < species.num_ships; ship_index++ {
            ship++

            if (ship.pn == 99) {
                continue;
            }

            if (ship.status == UNDER_CONSTRUCTION) {
                continue;
            }

            num_ships++
            total_tonnage += ship.tonnage;

            if (ship.ttype == STARBASE) {
                total_defensive_power += power(ship.tonnage);
                all_starbase_tons     += ship.tonnage;
                n_starbases++
                nba++
            }else if (ship.class == TR) {
                all_transport_tons += ship.tonnage;
                n_transports++
                ntr++
            }else {
                if (ship.ttype == SUB_LIGHT) {
                    total_defensive_power += power(ship.tonnage);
                }else{
                    total_offensive_power += power(ship.tonnage);
                }
                all_warship_tons += ship.tonnage;
                n_warships++
                nwa++
            }
        }

        if (nwa < min_warships) {
            min_warships = nwa;
        }
        if (nwa > max_warships) {
            max_warships = nwa;
        }

        if (nba < min_starbases) {
            min_starbases = nba;
        }
        if (nba > max_starbases) {
            max_starbases = nba;
        }

        if (ntr < min_transports) {
            min_transports = ntr;
        }
        if (ntr > max_transports) {
            max_transports = ntr;
        }

        total_offensive_power +=
            (species.tech_level[ML] * total_offensive_power) / 50;

        total_defensive_power +=
            (species.tech_level[ML] * total_defensive_power) / 50;

        if (species.tech_level[ML] == 0) {
            total_defensive_power = 0;
            total_offensive_power = 0;
        }

        total_offensive_power /= 10;
        total_defensive_power /= 10;

        printf("%5d", num_ships);
        printf("%5d", num_yards);
        printf("%8ld%8ld\n", total_offensive_power, total_defensive_power);
    }

    m = n_species / 2;
    printf("\n");
    for i = 0; i < 6; i++ {
        avg_tech_level = (all_tech_level[i] + m) / n_species;
        printf("Average %s tech level = %d (min = %d, max = %d)\n",
               tech_name[i], avg_tech_level, min_tech_level[i], max_tech_level[i]);
    }

    i = ((10 * n_warships) + m) / n_species;
    printf("\nAverage number of warships per species = %d.%d (min = %d, max = %d)\n",
           i / 10, i % 10, min_warships, max_warships);

    if (n_warships == 0) {
        n_warships = 1;
    }
    avg_warship_tons = (10000 * all_warship_tons) / n_warships;
    avg_warship_tons = 1000 * ((avg_warship_tons + 500) / 1000);
    printf("Average warship size = %s tons\n", commas(avg_warship_tons));

    avg_warship_tons = (10000 * all_warship_tons) / n_species;
    avg_warship_tons = 1000 * ((avg_warship_tons + 500) / 1000);
    printf("Average total warship tonnage per species = %s tons\n",
           commas(avg_warship_tons));

    i = ((10 * n_starbases) + m) / n_species;
    printf("\nAverage number of starbases per species = %d.%d (min = %d, max = %d)\n",
           i / 10, i % 10, min_starbases, max_starbases);

    if (n_starbases == 0) {
        n_starbases = 1;
    }
    avg_starbase_tons = (10000 * all_starbase_tons) / n_starbases;
    avg_starbase_tons = 1000 * ((avg_starbase_tons + 500) / 1000);
    printf("Average starbase size = %s tons\n", commas(avg_starbase_tons));

    avg_starbase_tons = (10000 * all_starbase_tons) / n_species;
    avg_starbase_tons = 1000 * ((avg_starbase_tons + 500) / 1000);
    printf("Average total starbase tonnage per species = %s tons\n",
           commas(avg_starbase_tons));

    i = ((10 * n_transports) + m) / n_species;
    printf("\nAverage number of transports per species = %d.%d (min = %d, max = %d)\n",
           i / 10, i % 10, min_transports, max_transports);

    if (n_transports == 0) {
        n_transports = 1;
    }
    avg_transport_tons = (10000 * all_transport_tons) / n_transports;
    avg_transport_tons = 1000 * ((avg_transport_tons + 500) / 1000);
    printf("Average transport size = %s tons\n", commas(avg_transport_tons));

    avg_transport_tons = (10000 * all_transport_tons) / n_species;
    avg_transport_tons = 1000 * ((avg_transport_tons + 500) / 1000);
    printf("Average total transport tonnage per species = %s tons\n",
           commas(avg_transport_tons));

    avg_yards = ((10 * n_yards) + m) / n_species;
    printf("\nAverage number of shipyards per species = %d.%d (min = %d, max = %d)\n",
           avg_yards / 10, avg_yards % 10, min_yards, max_yards);

    avg_pop_pl = ((10 * n_pop_pl) + m) / n_species;
    printf("\nAverage number of populated planets per species = %d.%d (min = %d, max = %d)\n",
           avg_pop_pl / 10, avg_pop_pl % 10, min_pop_pl, max_pop_pl);

    avg_production = (all_production + m) / n_species;
    printf("Average total production per species = %d (min = %d, max = %d)\n",
           avg_production, min_production, max_production);
}



//*************************************************************************
// TurnNumber.c

func TurnNumberMain(argc int, argv []string) {
    /* Check for valid command line. */
    if (argc != 1) {
        fprintf(stderr, "\n\tUsage: TurnNumber\n\n");
        exit(0);
    }

    /* Get galaxy data. */
    get_galaxy_data();

    /* Print the current turn number. */
    printf("%d\n", galaxy.turn_number);

    exit(0);
}
