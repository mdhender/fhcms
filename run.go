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
	"github.com/mdhender/fhcms/config"
	"github.com/mdhender/fhcms/store/jsondb"
	"log"
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

	verbose_mode = cfg.Log.Verbose

	/* Get commonly used data. */
	get_galaxy_data()
	get_planet_data()
	get_species_data()
	for species_number := 1; species_number <= num_species; species_number++ {
		spec_data[species_number-1].orders.filename = filepath.Join(cfg.Data.Orders, fmt.Sprintf("sp%02d.ord", spec_data[species_number-1].id))
	}
	get_transaction_data()
	if errors := get_order_data(verbose_mode, false); errors != nil {
		for _, err := range errors {
			log.Printf("error: %+v\n", err)
		}
	}
	num_locs = len(__jdb.Locations)

	// no-orders if not the first turn
	log.Printf("[orders] skipping NoOrders\n")

	log.Printf("[orders] skipping Combat\n")
	log.Printf("[orders] skipping PreDeparture\n")
	log.Printf("[orders] skipping Jumps\n")
	log.Printf("[orders] skipping Production\n")
	log.Printf("[orders] skipping PostArrival\n")
	__jdb.SetLocations()
	log.Printf("[orders] skipping Strike\n")
	log.Printf("[orders] skipping Finish\n")
	log.Printf("[orders] skipping Reports\n")
	log.Printf("[orders] skipping Stats\n") //	Stats(ds)

	return errors
}
