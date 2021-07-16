/*
 * Far Horizons Engine
 * Copyright (C) 2021  Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

type battle_data struct {
}
type galaxy_data struct {
	d_num_species int          /* Design number of species in galaxy. */
	num_species int            /* Actual number of species allocated. */
	radius int                 /* Galactic radius in parsecs. */
	turn_number int            /* Current turn number. */
}
type nampla_data struct {
	name[32]byte                 /* Name of planet. */
	x, y, z, pn int              /* Coordinates. */
	status int;                   /* Status of planet. */
	hiding bool                   /* HIDE order given. */
	hidden bool                   /* Colony is hidden. */
	planet_index int             /* Index (starting at zero) into the file "planets.dat" of this planet. */
	siege_eff int                /* Siege effectiveness - a percentage between  0 and 99. */
	shipyards int                /* Number of shipyards on planet. */
	IUs_needed int               /* Incoming ship with only CUs on board. */
	AUs_needed int               /* Incoming ship with only CUs on board. */
	auto_IUs int                 /* Number of IUs to be automatically installed. */
	auto_AUs int                 /* Number of AUs to be automatically installed. */
	IUs_to_install int           /* Colonial mining units to be installed. */
	AUs_to_install int           /* Colonial manufacturing units to be installed. */
	mi_base int                  /* Mining base times 10. */
	ma_base int                  /* Manufacturing base times 10. */
	pop_units int                /* Number of available population units. */
	item_quantity[MAX_ITEMS]int /* Quantity of each item available. */
	use_on_ambush int            /* Amount to use on ambush. */
	message int                  /* Message associated with this planet, if any. */
	special int                  /* Different for each application. */
}
type planet_data struct {
	temperature_class int;        /* Temperature class, 1-30. */
	pressure_class int;           /* Pressure class, 0-29. */
	special int;                  /* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	gas [4]int;                   /* Gas in atmosphere. Zero if none. */
	gas_percent[4]int;           /* Percentage of gas in atmosphere. */
	diameter int                 /* Diameter in thousands of kilometers. */
	gravity int                  /* Surface gravity. Multiple of Earth gravity times 100. */
	mining_difficulty int        /* Mining difficulty times 100. */
	econ_efficiency int          /* Economic efficiency. Always 100 for a  home planet. */
	md_increase int              /* Increase in mining difficulty. */
	message int                  /* Message associated with this planet,  if any. */
}
type ship_data_ struct {
	name[32]byte                     /* Name of ship. */
	x, y, z, pn int                  /* Current coordinates. */
	status int                       /* Current status of ship. */
	ttype int                         /* Ship type. */
	dest_x, dest_y, dest_z int               /* Destination if ship was forced to jump from combat.  Also used by TELESCOPE command. */
	just_jumped bool                  /* Set if ship jumped this turn. */
	arrived_via_wormhole bool         /* Ship arrived via wormhole in the PREVIOUS turn. */
	class int                        /* Ship class. */
	tonnage int                      /* Ship tonnage divided by 10,000. */
	item_quantity[MAX_ITEMS] int     /* Quantity of each item carried. */
	age int                          /* Ship age. */
	remaining_cost int               /* The cost needed to complete the ship if still under construction. */
	reserved4 int                    /* Unused. Zero for now. */
	loading_point int                /* Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet. */
	unloading_point int              /* Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet. */
	special int                      /* Different for each application. */
}
type sp_loc_data struct {
	s, x, y, z int            /* Species number, x, y, and z. */
}
type species_data struct {
	name[32]byte                     /* Name of species. */
	govt_name[32]byte                /* Name of government. */
	govt_ttype[32]byte                /* Type of government. */
	char  x, y, z, pn int                  /* Coordinates of home planet. */
	required_gas int                 /* Gas required by species. */
	required_gas_min int             /* Minimum needed percentage. */
	required_gas_max int             /* Maximum allowed percentage. */
	reserved5 int                    /* Zero for now. */
	neutral_gas[6]int               /* Gases neutral to species. */
	poison_gas[6]int                /* Gases poisonous to species. */
	auto_orders bool                  /* AUTO command was issued. */
	tech_level[6]int                /* Actual tech levels. */
	init_tech_level[6]int           /* Tech levels at start of turn. */
	tech_knowledge[6]int            /* Unapplied tech level knowledge. */
	num_namplas int                  /* Number of named planets, including home planet and colonies. */
	num_ships int                    /* Number of ships. */
	tech_eps[6]int                  /* Experience points for tech levels. */
	hp_original_base int             /* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	econ_units int                   /* Number of economic units. */
	fleet_cost int                   /* Total fleet maintenance cost. */
	fleet_percent_cost int           /* Fleet maintenance cost as a percentage times one hundred. */
	contact[NUM_CONTACT_WORDS]int /* A bit is set if corresponding species has been met. */
	ally[NUM_CONTACT_WORDS]int /* A bit is set if corresponding species is considered an ally. */
	enemy[NUM_CONTACT_WORDS]int /* A bit is set if corresponding species is considered an enemy. */
}
type star_data struct {
	x, y, z int              /* Coordinates. */
	ttype int                 /* Dwarf, degenerate, main sequence or giant. */
	color int                /* Star color. Blue, blue-white, etc. */
	size int                 /* Star size, from 0 thru 9 inclusive. */
	num_planets int          /* Number of usable planets in star system. */
	home_system int          /* true if this is a good potential home system. */
	worm_here int            /* true if wormhole entry/exit. */
	worm_x, worm_y, worm_z int
	planet_index int         /* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	message int              /* Message associated with this star system, if any. */
	visited_by[NUM_CONTACT_WORDS] int /* A bit is set if corresponding species has  been here. */
}
type trans_data struct {
	ttype int;                 /* Transaction type. */
	donor, recipient int
	value int                /* Value of transaction. */
	char  x, y, z, pn int          /* Location associated with transaction. */
	number1 int              /* Other items associated with transaction.*/
	name1[40] byte
	number2 int
	name2[40] byte
	number3 int
	name3[40] byte
}
