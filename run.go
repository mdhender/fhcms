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
	"fmt"
	"github.com/mdhender/fhcms/config"
	"github.com/mdhender/fhcms/prng"
	"github.com/mdhender/fhcms/store/jsondb"
	"log"
	"os"
	"path/filepath"
)

func run(cfg *config.Config) (errors []error) {
	jdb, err := jsondb.Read(filepath.Join(cfg.Data.JDB, "galaxy.json"))
	if err != nil {
		return []error{err}
	}
	// globals. argh.
	__jdb = jdb
	__jdb.SetLocations()
	g := &globals{
		__defaultPRNG:     prng.New(1924085713),
		log_start_of_line: true,
		num_locs:          len(jdb.Locations),
		stdout:            os.Stdout,
		stderr:            os.Stderr,
		verbose_mode:      cfg.Log.Verbose,
	}

	/* Get commonly used data. */
	if err := g.get_galaxy_data(jdb); err != nil {
		return []error{err}
	}
	if err := g.get_planet_data(jdb); err != nil {
		return []error{err}
	}
	if err := g.get_species_data(jdb); err != nil {
		return []error{err}
	}
	for species_number := 1; species_number <= num_species; species_number++ {
		g.sp_num = append(g.sp_num, species_number)
		g.spec_data[species_number-1].orders.filename = filepath.Join(cfg.Data.Orders, fmt.Sprintf("sp%02d.ord", g.spec_data[species_number-1].id))
		g.spec_data[species_number-1].log_file = &bytes.Buffer{}
	}
	if err := g.get_transaction_data(jdb); err != nil {
		return []error{err}
	}
	for _, err := range g.get_order_data() {
		log.Printf("error: %+v\n", err)
	}

	// no-orders if not the first turn
	log.Printf("[orders] skipping NoOrders\n")

	log.Printf("[orders] skipping Combat\n")
	log.Printf("[orders] skipping PreDeparture\n")
	log.Printf("[orders] skipping Jumps\n")
	for _, err := range g.executeJumpOrders() {
		log.Printf("error: %+v\n", err)
	}
	log.Printf("[orders] skipping Production\n")
	log.Printf("[orders] skipping PostArrival\n")
	__jdb.SetLocations()
	log.Printf("[orders] skipping Strike\n")
	log.Printf("[orders] skipping Finish\n")
	log.Printf("[orders] skipping Reports\n")
	log.Printf("[orders] skipping Stats\n") //	Stats(ds)

	return errors
}
