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

var (
	append_log              [MAX_SPECIES]byte
	battle_base             *battle_data
	c_nampla                [MAX_SPECIES]*nampla_data
	c_ship                  [MAX_SPECIES]*ship_data
	c_species               [MAX_SPECIES]*species_data
	doing_production        int
	end_of_file             int
	first_pass              int
	fleet_percent_cost      int
	galaxy                  galaxy_data
	header_printed          int
	home_planet             *planet_data
	ignore_field_distorters int
	input_file              io.Reader
	input_line              [256]byte
	input_line_pointer      *[]byte
	just_opened_file        int
	last_random             int
	loc                     [MAX_LOCATIONS]sp_loc_data
	log_file                io.Writer
	log_stdout              int
	make_enemy              [MAX_SPECIES][MAX_SPECIES]byte
	nampla1_base            *nampla_data
	nampla2_base            *nampla_data
	nampla_base             *nampla_data
	nampla_index            int
	nampla                  *nampla_data
	next_nampla_index       int
	next_nampla             *nampla_data
	num_intercepts          int
	num_locs                int
	num_planets             int
	num_stars               int
	num_transactions        int
	orders_file             io.Writer
	original_line           [256]byte
	planet_base             *planet_data
	planet_data_modified    int
	planet                  *planet_data
	pn                      int
	post_arrival_phase      bool // warning: must be initialized to true for PostArrivalMain
	printing_alien          int
	production_done         [1000]byte
	prompt_gm               int
	report_file             io.Writer
	ship1_base              *ship_data
	ship2_base              *ship_data
	ship_already_listed     [5000]byte
	ship_base               *ship_data
	ship_index              int
	ship                    *ship_data
	species_index           int
	species_number          int
	species                 *species_data
	sp_tech_level           [6]int
	star_base               *star_data
	star_data_modified      int
	star                    *star_data
	strike_phase            bool
	test_mode               bool
	transaction             [MAX_TRANSACTIONS]trans_data
	truncate_name           bool
	upper_name              []byte
	value                   int
	verbose_mode            bool
	x                       int
	y                       int
	z                       int
)
