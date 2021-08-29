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
	"github.com/mdhender/fhcms/parser"
)

type action_data struct {
	bomb_damage            [MAX_SHIPS]int
	fighting_species_index [MAX_SHIPS]int
	//fighting_unit          [MAX_SHIPS]*int // warning: code mixes pointer types here
	fighting_unit [MAX_SHIPS]struct {
		nampla *nampla_data
		ship   *ship_data_
	}
	original_age_or_PDs  [MAX_SHIPS]int
	num_shots            [MAX_SHIPS]int
	num_units_fighting   int
	shield_strength      [MAX_SHIPS]int
	shield_strength_left [MAX_SHIPS]int
	shots_left           [MAX_SHIPS]int
	surprised            [MAX_SHIPS]int
	unit_type            [MAX_SHIPS]int
	weapon_damage        [MAX_SHIPS]int
}

type battle_data struct {
	ambush_amount             [MAX_SPECIES]int
	can_be_surprised          [MAX_SPECIES]int
	enemy_mine                [MAX_SPECIES][MAX_SPECIES]int // will be 0 for no enemy, 1 for ATTACK, 2 for HIJACK
	engage_option             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]int
	engage_planet             [MAX_SPECIES][MAX_ENGAGE_OPTIONS]int
	fleet_withdraw_percentage [MAX_SPECIES]int
	haven_x                   [MAX_SPECIES]int
	haven_y                   [MAX_SPECIES]int
	haven_z                   [MAX_SPECIES]int
	hijacker                  [MAX_SPECIES]bool
	num_engage_options        [MAX_SPECIES]int
	num_species_here          int
	spec_num                  [MAX_SPECIES]int
	special_target            [MAX_SPECIES]int
	summary_only              [MAX_SPECIES]bool
	transport_withdraw_age    [MAX_SPECIES]int
	warship_withdraw_age      [MAX_SPECIES]int
	x                         int
	y                         int
	z                         int
}

type class_abbr struct {
	abbr_index int
	abbr_type  int
	name       string // argh. why did i add this?
	sub_light  bool
	tonnage    int
}

type galaxy_data struct {
	d_num_species int /* Design number of species in galaxy. */
	num_species   int /* Actual number of species allocated. */
	radius        int /* Galactic radius in parsecs. */
	turn_number   int /* Current turn number. */
}

type gas_code struct {
	code int
	abbr string
	name string
}

type gas_data struct {
	code       *gas_code
	percentage int
}

type nampla_data struct {
	name           string         /* Name of planet. */ // warning: code expects [32]byte
	x, y, z, pn    int            /* Coordinates. */
	AUs_needed     int            /* Incoming ship with only CUs on board. */
	auto_AUs       int            /* Number of AUs to be automatically installed. */
	auto_IUs       int            /* Number of IUs to be automatically installed. */
	hidden         int            /* Colony is hidden. */
	hiding         int            /* HIDE order given. */
	AUs_to_install int            /* Colonial manufacturing units to be installed. */
	IUs_to_install int            /* Colonial mining units to be installed. */
	IUs_needed     int            /* Incoming ship with only CUs on board. */
	item_quantity  [MAX_ITEMS]int /* Quantity of each item available. */
	ma_base        int            /* Manufacturing base times 10. */
	message        int            /* Message associated with this planet, if any. */
	mi_base        int            /* Mining base times 10. */
	planet_index   int            /* Index (starting at zero) into the file "planets.dat" of this planet. */
	pop_units      int            /* Number of available population units. */
	shipyards      int            /* Number of shipyards on planet. */
	siege_eff      int            /* Siege effectiveness - a percentage between  0 and 99. */
	special        int            /* Different for each application. */
	status         int            /* Status of planet. */
	use_on_ambush  int            /* Amount to use on ambush. */
}

type planet_data struct {
	x, y, z, orbit    int
	diameter          int    /* Diameter in thousands of kilometers. */
	gas               [4]int /* Gas in atmosphere. Zero if none. */
	gas_percent       [4]int /* Percentage of gas in atmosphere. */
	gravity           int    /* Surface gravity. Multiple of Earth gravity times 100. */
	econ_efficiency   int    /* Economic efficiency. Always 100 for a  home planet. */
	md_increase       int    /* Increase in mining difficulty. */
	message           int    /* Message associated with this planet,  if any. */
	mining_difficulty int    /* Mining difficulty times 100. */
	pressure_class    int    /* Pressure class, 0-29. */
	special           int    /* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	temperature_class int    /* Temperature class, 1-30. */
	// mdhender: added the following for convenience
	atmosphere []*gas_data // atmospheric gases sorted by percentage descending
}

type ship_data_ struct {
	name                   string         /* Name of ship. */ // warning: code expects [32]byte
	x, y, z, pn            int            /* Current coordinates. */
	age                    int            /* Ship age. */
	arrived_via_wormhole   bool           /* Ship arrived via wormhole in the PREVIOUS turn. */
	class                  int            /* Ship class. */
	dest_x, dest_y, dest_z int            /* Destination if ship was forced to jump from combat.  Also used by TELESCOPE command. And for SHIELDS in combat. */
	item_quantity          [MAX_ITEMS]int /* Quantity of each item carried. */
	just_jumped            int            /* Set if ship jumped this turn. */ // 50 means MOVEd?
	loading_point          int            /* Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet. */
	remaining_cost         int            /* The cost needed to complete the ship if still under construction. */
	special                int            /* Different for each application. */
	status                 int            /* Current status of ship. */
	tonnage                int            /* Ship tonnage divided by 10,000. */
	ttype                  int            /* Ship type. */
	unloading_point        int            /* Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet. */
}

type sp_loc_data struct {
	s, x, y, z int /* Species number, x, y, and z. */
}

type species_data struct {
	name               string /* Name of species. */ // warning: code expects [32]byte
	x, y, z, pn        int    /* Coordinates of home planet. */
	auto_orders        bool   /* AUTO command was issued. */
	econ_units         int    /* Number of economic units. */
	fleet_cost         int    /* Total fleet maintenance cost. */
	fleet_percent_cost int    /* Fleet maintenance cost as a percentage times one hundred. */
	govt_name          string /* Name of government. */ // warning: code expects [32]byte
	govt_type          string /* Type of government. */ // warning: code expects [32]byte
	hp_original_base   int    /* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	init_tech_level    [6]int /* Tech levels at start of turn. */
	neutral_gas        [6]int /* Gases neutral to species. */
	num_namplas        int    /* Number of named planets, including home planet and colonies. */
	num_ships          int    /* Number of ships. */
	poison_gas         [6]int /* Gases poisonous to species. */
	required_gas       int    /* Gas required by species. */
	required_gas_max   int    /* Maximum allowed percentage. */
	required_gas_min   int    /* Minimum needed percentage. */
	tech_eps           [6]int /* Experience points for tech levels. */
	tech_knowledge     [6]int /* Unapplied tech level knowledge. */
	tech_level         [6]int /* Actual tech levels. */
	// warning: code expects [NUM_CONTACT_WORDS]int for ally, contact, and enemy
	ally    [MAX_SPECIES]bool /* A bit is set if corresponding species is considered an ally. */
	contact [MAX_SPECIES]bool /* A bit is set if corresponding species has been met. */
	enemy   [MAX_SPECIES]bool /* A bit is set if corresponding species is considered an enemy. */
	// mdhender: added the following for convenience
	id      int
	namplas []*nampla_data
	orders  struct {
		data     *parser.Node
		filename string
		errors   []error
	}
	ships []*ship_data_
}

type star_color_code struct {
	code int
	abbr string
	name string
}

type star_data struct {
	x, y, z                int               /* Coordinates. */
	color                  int               /* Star color. Blue, blue-white, etc. */ // TODO: this should be a *star_color_code
	home_system            bool              /* true if this is a good potential home system. */
	message                int               /* Message associated with this star system, if any. */
	num_planets            int               /* Number of usable planets in star system. */
	planet_index           int               /* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	size                   int               /* Star size, from 0 thru 9 inclusive. */
	ttype                  int               /* Dwarf, degenerate, main sequence or giant. */
	visited_by             [MAX_SPECIES]bool /* A bit is set if corresponding species has  been here. */ // warning: was [NUM_CONTACT_WORDS]int
	worm_here              bool              /* true if wormhole entry/exit. */
	worm_x, worm_y, worm_z int
}

type star_type_code struct {
	code int
	abbr string
	name string
}

type trans_data struct {
	ttype            int /* Transaction type. */
	donor, recipient int
	value            int    /* Value of transaction. */
	x, y, z, pn      int    /* Location associated with transaction. */
	number1          int    /* Other items associated with transaction.*/
	name1            string // warning: was [40]byte
	number2          int
	name2            string // warning: was [40]byte
	number3          int
	name3            string // warning: was [40]byte
}

// transfer_point can be either a ship or a named planet
type transfer_point struct {
	nampla *nampla_data
	ship   *ship_data_
}
