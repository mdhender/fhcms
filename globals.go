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
	"bytes"
	"github.com/mdhender/fhcms/prng"
	"io"
)

// globals are data used in most or all programs
type globals struct {
	__defaultPRNG           *prng.PRNG
	galaxy                  *galaxy_data
	ignore_field_distorters bool
	log_file                *bytes.Buffer // a buffer needs no initialization and acts like an io.Writer
	log_start_of_line       bool
	log_stdout              bool
	log_summary             bool
	log_to_file             bool
	logging_disabled        bool
	namp_data               [MAX_SPECIES][]*nampla_data // warning: code assumes [MAX_SPECIES]*nampla_data
	nampla_base             []*nampla_data              // warning: code assumes *nampla_data // by convention, nampla_base[0] is the species home planet
	num_locs                int
	num_new_namplas         [MAX_SPECIES]int // zero based
	num_new_ships           [MAX_SPECIES]int // zero based
	num_planets             int
	num_species             int
	planet_base             []*planet_data             // warning: code assumes *planet_data
	ship_base               []*ship_data_              // warning: code expects *ship_data_
	ship_data               [MAX_SPECIES][]*ship_data_ // warning: code expects [MAX_SPECIES]*ship_data_
	sp_num                  []int                      // list of species to process
	spec_data               [MAX_SPECIES]*species_data
	species                 *species_data
	stderr                  io.Writer
	stdout                  io.Writer
	summary_file            *bytes.Buffer // a buffer needs no initialization and acts like an io.Writer
	truncate_name           bool
	verbose_mode            bool
}
