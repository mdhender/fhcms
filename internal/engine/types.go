/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published
by the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY int without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
******************************************************************************/

package engine

import (
	"bytes"
	"github.com/mdhender/fhcms/cms/prng"
)

type Engine struct {
	galaxy galaxy_data

	star_base []*star_data
	num_stars int

	planet_base []*planet_data
	home_planet *planet_data

	// named planet globals
	namp_data   [][]*nampla_data // zero-based index by species of nampla
	nampla_base []*nampla_data   // single species nampla
	nampla      *nampla_data     // single nampla

	// ship globals
	ship_data                 [][]*ship_data // zero-based index by species of ship
	ship_base                 []*ship_data   // single species ship
	abbr_index                int
	abbr_type                 int
	correct_spelling_required int // TRUE or FALSE
	input_abbr                []byte
	original_name             []byte
	ship                      *ship_data // single ship
	ship_index                int        // zero-based
	sub_light                 int        // TRUE or FALSE
	tonnage                   int
	upper_name                []byte

	// species globals
	spec_data      []*species_data // zero-based index by species
	spec_logs      []*bytes.Buffer // zero-based index by species
	spec_orders    [][]byte        // zero-based index by species
	species        *species_data
	species_index  int
	species_number int

	// location globals
	loc            []*sp_loc_data
	locations_base []*sp_loc_data
	num_locs       int

	// combat globals
	ambush_took_place  int // TRUE or FALSE
	attacking_ML       int
	battle_base        []*battle_data
	c_nampla           [MAX_SPECIES][]*nampla_data // indexed by species_no
	c_ship             [MAX_SPECIES][]*ship_data   // indexed by species_no
	c_species          [MAX_SPECIES]*species_data  // indexed by species_no
	combat_location    [1000]int
	combat_log         *FILE
	combat_option      [1000]int
	defending_ML       int
	deep_space_defense int              // TRUE or FALSE, maybe?
	field_distorted    [MAX_SPECIES]int // indexed by species, TRUE or FALSE
	first_battle       int              // TRUE or FALSE
	germ_bombs_used    [MAX_SPECIES][MAX_SPECIES]int
	make_enemy         [MAX_SPECIES][MAX_SPECIES]int // zero-based index, matrix of species that are enemies, content is one-based species_number
	num_combat_options int
	num_transactions   int
	strike_phase       int
	temp_log           [MAX_SPECIES]*bytes.Buffer // zero-based indexed by species_no
	transaction        [MAX_TRANSACTIONS]trans_data
	x_attacked_y       [MAX_SPECIES][MAX_SPECIES]int

	// input and output hacks
	append_log         [MAX_SPECIES]int // zero-based index by species
	end_of_file        int
	input_file         *FILE
	just_opened_file   int
	input_line         []byte
	input_line_pointer []byte
	log_indentation    int // number of spaces to indent
	log_file           *FILE
	log_line           []byte
	log_position       int
	log_to_file        int // TRUE or FALSE
	log_start_of_line  int // TRUE or FALSE
	log_stdout         int // TRUE or FALSE
	log_summary        int // TRUE or FALSE
	logging_disabled   int
	original_line      []byte
	stderr             *FILE
	stdout             *FILE
	summary_file       *FILE

	// miscellaneous globals
	defaultPRNG             *prng.PRNG
	ignore_field_distorters int
	orders_file             *bytes.Buffer
	prompt_gm               bool
	test_mode               int
	truncate_name           int
	value                   int // set by get_value()
	verbose_mode            int
	x, y, z                 int
}

type Coords struct{}
type Species struct{}

type galaxy_data struct {
	d_num_species int /* Design number of species in galaxy. */
	num_species   int /* Actual number of species allocated. */
	radius        int /* Galactic radius in parsecs. */
	turn_number   int /* Current turn number. */
}

type nampla_data struct {
	name           string         /* Name of planet. */
	x, y, z, pn    int            /* Coordinates. */
	status         int            /* Status of planet. */
	reserved1      int            /* Zero for now. */
	hiding         int            /* HIDE order given. */
	hidden         int            /* Colony is hidden. */
	reserved2      int            /* Zero for now. */
	planet_index   int            /* Index (starting at zero) into the file "planets.dat" of this planet. */
	siege_eff      int            /* Siege effectiveness - a percentage between 0 and 99. */
	shipyards      int            /* Number of shipyards on planet. */
	reserved4      int            /* Zero for now. */
	IUs_needed     int            /* Incoming ship with only CUs on board. */
	AUs_needed     int            /* Incoming ship with only CUs on board. */
	auto_IUs       int            /* Number of IUs to be automatically installed. */
	auto_AUs       int            /* Number of AUs to be automatically installed. */
	reserved5      int            /* Zero for now. */
	IUs_to_install int            /* Colonial mining units to be installed. */
	AUs_to_install int            /* Colonial manufacturing units to be installed. */
	mi_base        int            /* Mining base times 10. */
	ma_base        int            /* Manufacturing base times 10. */
	pop_units      int            /* Number of available population units. */
	item_quantity  [MAX_ITEMS]int /* Quantity of each item available. */
	reserved6      int            /* Zero for now. */
	use_on_ambush  int            /* Amount to use on ambush. */
	message        int            /* Message associated with this planet, if any. */
	special        int            /* Different for each application. */
	padding        [28]int        /* Use for expansion. Initialized to all zeroes. */
}

type planet_data struct {
	temperature_class int    /* Temperature class, 1-30. */
	pressure_class    int    /* Pressure class, 0-29. */
	special           int    /* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	reserved1         int    /* Reserved for future use. Zero for now. */
	gas               [4]int /* Gas in atmosphere. Zero if none. */
	gas_percent       [4]int /* Percentage of gas in atmosphere. */
	reserved2         int    /* Reserved for future use. Zero for now. */
	diameter          int    /* Diameter in thousands of kilometers. */
	gravity           int    /* Surface gravity. Multiple of Earth gravity times 100. */
	mining_difficulty int    /* Mining difficulty times 100. */
	econ_efficiency   int    /* Economic efficiency. Always 100 for a home planet. */
	md_increase       int    /* Increase in mining difficulty. */
	message           int    /* Message associated with this planet, if any. */
	reserved3         int    /* Reserved for future use. Zero for now. */
	reserved4         int    /* Reserved for future use. Zero for now. */
	reserved5         int    /* Reserved for future use. Zero for now. */
}

type ship_data struct {
	name                 string         /* Name of ship. */
	x, y, z, pn          int            /* Current coordinates. */
	status               int            /* Current status of ship. */
	_type                int            /* Ship type. */
	dest_x, dest_y       int            /* Destination if ship was forced to jump from combat. */
	dest_z               int            /* Ditto. Also used by TELESCOPE command. */
	just_jumped          int            /* Set if ship jumped this turn. */
	arrived_via_wormhole int            /* Ship arrived via wormhole in the PREVIOUS turn. */
	reserved1            int            /* Unused. Zero for now. */
	reserved2            int            /* Unused. Zero for now. */
	reserved3            int            /* Unused. Zero for now. */
	class                int            /* Ship class. */
	tonnage              int            /* Ship tonnage divided by 10,000. */
	item_quantity        [MAX_ITEMS]int /* Quantity of each item carried. */
	age                  int            /* Ship age. */
	remaining_cost       int            /* The cost needed to complete the ship if still under construction. */
	reserved4            int            /* Unused. Zero for now. */
	loading_point        int            /* Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet. */
	unloading_point      int            /* Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet. */
	special              int            /* Different for each application. */
	padding              [28]int        /* Use for expansion. Initialized to all zeroes. */
}

type sp_loc_data struct {
	s, x, y, z int /* Species number, x, y, and z. */
}

type species_data struct {
	name               string  /* Name of species. */
	govt_name          string  /* Name of government. */
	govt_type          string  /* Type of government. */
	x, y, z, pn        int     /* Coordinates of home planet. */
	required_gas       int     /* Gas required by species. */
	required_gas_min   int     /* Minimum needed percentage. */
	required_gas_max   int     /* Maximum allowed percentage. */
	reserved5          int     /* Zero for now. */
	neutral_gas        [6]int  /* Gases neutral to species. */
	poison_gas         [6]int  /* Gases poisonous to species. */
	auto_orders        int     /* AUTO command was issued. */
	reserved3          int     /* Zero for now. */
	reserved4          int     /* Zero for now. */
	tech_level         [6]int  /* Actual tech levels. */
	init_tech_level    [6]int  /* Tech levels at start of turn. */
	tech_knowledge     [6]int  /* Unapplied tech level knowledge. */
	num_namplas        int     /* Number of named planets, including home planet and colonies. */
	num_ships          int     /* Number of ships. */
	tech_eps           [6]int  /* Experience points for tech levels. */
	hp_original_base   int     /* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	econ_units         int     /* Number of economic units. */
	fleet_cost         int     /* Total fleet maintenance cost. */
	fleet_percent_cost int     /* Fleet maintenance cost as a percentage times one hundred. */
	contact            []int   /* zero based, true if corresponding species has been met. */
	ally               []int   /* zero based, true if corresponding species is considered an ally. */
	enemy              []int   /* zero based, true if corresponding species is considered an enemy. */
	padding            [12]int /* Use for expansion. Initialized to all zeroes. */
}

type star_data struct {
	x, y, z                int /* Coordinates. */
	_type                  int /* Dwarf, degenerate, main sequence or giant. */
	color                  int /* Star color. Blue, blue-white, etc. */
	size                   int /* Star size, from 0 thru 9 inclusive. */
	num_planets            int /* Number of usable planets in star system. */
	home_system            int /* TRUE if this is a good potential home system. */
	worm_here              int /* TRUE if wormhole entry/exit. */
	worm_x, worm_y, worm_z int
	reserved1              int   /* Reserved for future use. Zero for now. */
	reserved2              int   /* Reserved for future use. Zero for now. */
	planet_index           int   /* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	message                int   /* Message associated with this star system, if any. */
	visited_by             []int /* A bit is set if corresponding species has been here. */
	reserved3              int   /* Reserved for future use. Zero for now. */
	reserved4              int   /* Reserved for future use. Zero for now. */
	reserved5              int   /* Reserved for future use. Zero for now. */
}

type action_data struct {
	num_units_fighting     int
	fighting_species_index [MAX_SHIPS]int
	num_shots              [MAX_SHIPS]int
	shots_left             [MAX_SHIPS]int
	weapon_damage          [MAX_SHIPS]int
	shield_strength        [MAX_SHIPS]int
	shield_strength_left   [MAX_SHIPS]int
	original_age_or_PDs    [MAX_SHIPS]int
	bomb_damage            [MAX_SHIPS]int
	surprised              [MAX_SHIPS]int
	unit_type              [MAX_SHIPS]int
	fighting_unit          [MAX_SHIPS]interface{} // either *ship_data or *nampla_data
}

type battle_data struct {
	x, y, z                   int
	num_species_here          int
	spec_num                  [MAX_SPECIES]int
	summary_only              [MAX_SPECIES]int
	transport_withdraw_age    [MAX_SPECIES]int
	warship_withdraw_age      [MAX_SPECIES]int
	fleet_withdraw_percentage [MAX_SPECIES]int
	haven_x                   [MAX_SPECIES]int
	haven_y                   [MAX_SPECIES]int
	haven_z                   [MAX_SPECIES]int
	special_target            [MAX_SPECIES]int
	hijacker                  [MAX_SPECIES]int
	can_be_surprised          [MAX_SPECIES]int
	enemy_mine                [MAX_SPECIES][MAX_SPECIES]int
	num_engage_options        [MAX_SPECIES]int
	engage_option             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]int
	engage_planet             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]int
	ambush_amount             [MAX_SPECIES]int
}

type trans_data struct {
	_type            int /* Transaction type. */
	donor, recipient int
	value            int /* Value of transaction. */
	x, y, z, pn      int /* Location associated with transaction. */
	number1          int /* Other items associated with transaction.*/
	name1            string
	number2          int
	name2            string
	number3          int
	name3            string
}
