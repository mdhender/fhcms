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

import "io"

/* Global data used in most or all programs. */

var (
	abbr_index   int
	abbr_type    int
	ambush_took_place int
	append_log   [MAX_SPECIES]bool
	attacking_ML int
	balance      int
	battle_base  []*battle_data // warning: code expects *battle_data
	c_nampla     [MAX_SPECIES]*nampla_data
	color_char   = []byte(" OBAFGKM")
	command_abbr = []string{ // warning: code assumes [NUM_COMMANDS][4]byte
		"   ", "ALL", "AMB", "ATT", "AUT", "BAS", "BAT", "BUI", "CON",
		"DEE", "DES", "DEV", "DIS", "END", "ENE", "ENG", "EST", "HAV",
		"HID", "HIJ", "IBU", "ICO", "INS", "INT", "JUM", "LAN", "MES",
		"MOV", "NAM", "NEU", "ORB", "PJU", "PRO", "REC", "REP", "RES",
		"SCA", "SEN", "SHI", "STA", "SUM", "SUR", "TAR", "TEA", "TEC",
		"TEL", "TER", "TRA", "UNL", "UPG", "VIS", "WIT", "WOR", "ZZZ",
	}
	command_name = []string{ // warning: code assumes [NUM_COMMANDS][16]byte
		"Undefined", "Ally", "Ambush", "Attack", "Auto", "Base",
		"Battle", "Build", "Continue", "Deep", "Destroy", "Develop",
		"Disband", "End", "Enemy", "Engage", "Estimate", "Haven",
		"Hide", "Hijack", "Ibuild", "Icontinue", "Install", "Intercept",
		"Jump", "Land", "Message", "Move", "Name", "Neutral", "Orbit",
		"Pjump", "Production", "Recycle", "Repair", "Research", "Scan",
		"Send", "Shipyard", "Start", "Summary", "Surrender", "Target",
		"Teach", "Tech", "Telescope", "Terraform", "Transfer", "Unload",
		"Upgrade", "Visited", "Withdraw", "Wormhole", "ZZZ",
	}
	c_ship             [MAX_SPECIES]*ship_data_
	c_species          [MAX_SPECIES]*species_data
	combat_option      [1000]int
	combat_location    [1000]int
	data_in_memory     [MAX_SPECIES]bool
	data_modified      [MAX_SPECIES]bool
	deep_space_defense int
	defending_ML int
	doing_production   int
	end_of_file        bool
	field_distorted [MAX_SPECIES]bool
	first_battle bool = true
	first_pass         bool
	fleet_percent_cost int
	galaxy             galaxy_data
	gas_string         = []string{ // warning: code assumes [14][4]byte
		"   ", "H2", "CH4", "He", "NH3", "N2", "CO2", "O2", "HCl", "Cl2", "F2", "H2O", "SO2", "H2S",
	}
	g_spec_name             [32]byte
	g_spec_number           int
	header_printed          int
	home_planet             *planet_data
	ignore_field_distorters int
	input_file              io.Reader
	input_line              [256]byte
	input_line_pointer      *cstring
	item_abbr               = []string{ // warning: code assumes [MAX_ITEMS][4]byte
		"RM", "PD", "SU", "DR", "CU", "IU", "AU", "FS",
		"JP", "FM", "FJ", "GT", "FD", "TP", "GW", "SG1",
		"SG2", "SG3", "SG4", "SG5", "SG6", "SG7", "SG8", "SG9",
		"GU1", "GU2", "GU3", "GU4", "GU5", "GU6", "GU7", "GU8",
		"GU9", "X1", "X2", "X3", "X4", "X5",
	}
	item_carry_capacity = []int{ // warning: code assumes [MAX_ITEMS]int
		1, 3, 20, 1, 1, 1, 1, 1,
		10, 5, 5, 20, 1, 100, 100, 5,
		10, 15, 20, 25, 30, 35, 40, 45,
		5, 10, 15, 20, 25, 30, 35, 40,
		45, 9999, 9999, 9999, 9999, 9999,
	}
	item_cost = []int{ // warning: code assumes [MAX_ITEMS]int
		1, 1, 110, 50, 1, 1, 1, 25,
		100, 100, 125, 500, 50, 50000, 1000, 250,
		500, 750, 1000, 1250, 1500, 1750, 2000, 2250,
		250, 500, 750, 1000, 1250, 1500, 1750, 2000,
		2250, 9999, 9999, 9999, 9999, 9999,
	}
	item_critical_tech = []int{ // warning: code assumes [MAX_ITEMS]int
		MI, ML, MA, MA, LS, MI, MA, GV,
		GV, GV, GV, GV, LS, BI, BI, LS,
		LS, LS, LS, LS, LS, LS, LS, LS,
		ML, ML, ML, ML, ML, ML, ML, ML,
		ML, 99, 99, 99, 99, 99,
	}
	item_name = []string{ // warning: code assumes [MAX_ITEMS][32]byte
		"Raw Material Unit", "Planetary Defense Unit", "Starbase Unit", "Damage Repair Unit",
		"Colonist Unit", "Colonial Mining Unit", "Colonial Manufacturing Unit", "Fail-Safe Jump Unit",
		"Jump Portal Unit", "Forced Misjump Unit", "Forced Jump Unit", "Gravitic Telescope Unit",
		"Field Distortion Unit", "Terraforming Plant", "Germ Warfare Bomb", "Mark-1 Shield Generator",
		"Mark-2 Shield Generator", "Mark-3 Shield Generator", "Mark-4 Shield Generator", "Mark-5 Shield Generator",
		"Mark-6 Shield Generator", "Mark-7 Shield Generator", "Mark-8 Shield Generator", "Mark-9 Shield Generator",
		"Mark-1 Gun Unit", "Mark-2 Gun Unit", "Mark-3 Gun Unit", "Mark-4 Gun Unit",
		"Mark-5 Gun Unit", "Mark-6 Gun Unit", "Mark-7 Gun Unit", "Mark-8 Gun Unit",
		"Mark-9 Gun Unit", "X1 Unit", "X2 Unit", "X3 Unit",
		"X4 Unit", "X5 Unit",
	}
	item_tech_requirment = []int{ // warning [MAX_ITEMS]int
		1, 1, 20, 30, 1, 1, 1, 20,
		25, 30, 40, 50, 20, 40, 50, 10,
		20, 30, 40, 50, 60, 70, 80, 90,
		10, 20, 30, 40, 50, 60, 70, 80,
		90, 999, 999, 999, 999, 999,
	}
	just_opened_file     int
	last_random          int
	loc                  [MAX_LOCATIONS]sp_loc_data
	log_file             io.Writer
	log_indentation      int
	log_line             [1028]byte
	log_position         int
	log_start_of_line    bool
	log_stdout           bool
	log_summary          bool
	log_to_file          bool
	logging_disabled     bool
	make_enemy           [MAX_SPECIES][MAX_SPECIES]bool
	namp_data            [MAX_SPECIES]*nampla_data
	nampla1_base         *nampla_data
	nampla2_base         *nampla_data
	nampla_base          *nampla_data
	nampla_index         int
	nampla               *nampla_data
	next_nampla_index    int
	next_nampla          *nampla_data
	num_combat_options   int
	num_intercepts       int
	num_locs             int
	num_new_namplas      [MAX_SPECIES]int
	num_new_ships        [MAX_SPECIES]int
	num_planets          int
	num_stars            int
	num_transactions     int
	orders_file          io.Writer
	original_line        [256]byte
	planet_base          *planet_data
	planet_data_modified int
	planet               *planet_data
	pn                   int
	post_arrival_phase   bool // warning: must be initialized to true for PostArrivalMain
	printing_alien       int
	production_done      [1000]byte
	prompt_gm            bool
	report_file          io.Writer
	ship1_base           *ship_data_
	ship2_base           *ship_data_
	ship_abbr            = []string{ // warning: code assumes [NUM_SHIP_CLASSES][4]byte
		"PB", "CT", "ES", "FF", "DD", "CL", "CS",
		"CA", "CC", "BC", "BS", "DN", "SD", "BM",
		"BW", "BR", "BA", "TR",
	}
	ship_already_listed [5000]byte
	ship_base           []*ship_data_ // warning: code expects *ship_data_
	ship_cost           = []int{      // warning: code assumes [NUM_SHIP_CLASSES]int
		100, 200, 500, 1000, 1500, 2000, 2500,
		3000, 3500, 4000, 4500, 5000, 5500, 6000,
		6500, 7000, 100, 100,
	}
	ship_data  [MAX_SPECIES][]*ship_data_ // warning: code expects [MAX_SPECIES]*ship_data_
	ship_index int
	ship       *ship_data_
	// Look-up table for ship defensive/offensive power uses ship.tonnage as an index.
	// Each value is equal to 100 * (ship.tonnage)^1.2.
	// The 'power' subroutine uses recursion to calculate values for tonnages over 100.
	ship_power = []int{ // warning: code expects [101]int
		0, // zeroth element required but not used
		100, 230, 374, 528, 690, 859, 1033, 1213, 1397, 1585,
		1777, 1973, 2171, 2373, 2578, 2786, 2996, 3209, 3424, 3641,
		3861, 4082, 4306, 4532, 4759, 4988, 5220, 5452, 5687, 5923,
		6161, 6400, 6641, 6883, 7127, 7372, 7618, 7866, 8115, 8365,
		8617, 8870, 9124, 9379, 9635, 9893, 10151, 10411, 10672, 10934,
		11197, 11461, 11725, 11991, 12258, 12526, 12795, 13065, 13336, 13608,
		13880, 14154, 14428, 14703, 14979, 15256, 15534, 15813, 16092, 16373,
		16654, 16936, 17218, 17502, 17786, 18071, 18356, 18643, 18930, 19218,
		19507, 19796, 20086, 20377, 20668, 20960, 21253, 21547, 21841, 22136,
		22431, 22727, 23024, 23321, 23619, 23918, 24217, 24517, 24818, 25119,
	}
	ship_tonnage = []int{ // warning: code assumes [NUM_SHIP_CLASSES]int
		1, 2, 5, 10, 15, 20, 25,
		30, 35, 40, 45, 50, 55, 60,
		65, 70, 1, 1,
	}
	ship_ttype = []string{ // warning: code assumes [3][2]byte
		"", "S", "S",
	}
	size_char          []byte = []byte("0123456789")
	sp_tech_level      [6]int
	spec_data          [MAX_SPECIES]*species_data
	species            *species_data
	species_index      int
	species_number     int
	star_base          *star_data
	star_data_modified int
	star               *star_data
	strike_phase       bool
	summary_file       io.Writer
	tech_abbr          = []string{ // warning: code assumes [6][4]byte
		"MI", "MA", "ML", "GV", "LS", "BI",
	}
	tech_name = []string{ // warning: code assumes [6][16]byte
		"Mining", "Manufacturing", "Military", "Gravitics", "Life Support", "Biology",
	}
	test_mode     bool
	transaction   [MAX_TRANSACTIONS]trans_data
	truncate_name bool
	ttype_char    []byte = []byte(" dD g")
	upper_name    []byte
	value         int
	verbose_mode  bool
	x             int
	y             int
	z             int
)

/*
extern int abbr_type, abbr_index, species_number,
           species_index, ship_index, num_stars;
extern long value;
extern char input_line[256], original_line[256], original_name[32],
            upper_name[32], *input_line_pointer;

extern FILE *               log_file;
extern struct star_data *   star_base, *star;
extern struct species_data *species;
extern struct nampla_data * nampla;
extern struct ship_data *   ship, *ship_base;

extern int attacker_here, defender_here, logging_disabled,
           strike_phase, prompt_gm;
extern int   log_summary, num_combat_options;
extern int   truncate_name, num_transactions;
extern int   ignore_field_distorters;
extern char  combat_option[1000], combat_location[1000];
extern char  x_attacked_y[MAX_SPECIES][MAX_SPECIES];
extern FILE *log_file, *summary_file;
extern struct galaxy_data   galaxy;
extern struct species_data *c_species[MAX_SPECIES];
extern struct nampla_data * nampla_base, *c_nampla[MAX_SPECIES];
extern struct ship_data *   ship_base, *c_ship[MAX_SPECIES];




*/