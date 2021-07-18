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
	"github.com/mdhender/fhcms/orders"
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

	for i := 1; i <= numSpecies; i++ {
		turnData.Species[i] = &SpeciesTurnData{Id: i}
		td := turnData.Species[i]
		td.Species = jdb.Species[fmt.Sprintf("SP%02d", td.Id)]
		td.OrderFile = filepath.Join(cfg.Data.Orders, fmt.Sprintf("sp%02d.ord", i))

		log.Printf("orders: loading %q\n", td.OrderFile)
		o := orders.Parse(td.OrderFile)
		if o.Errors == nil {
			if verbose {
				fmt.Printf(";; SP%02d TURN %3d\n", i, jdb.Galaxy.TurnNumber)
				for _, section := range []*orders.Section{o.Combat, o.PreDeparture, o.Jumps, o.Production, o.PostArrival, o.Strikes} {
					if section != nil {
						fmt.Printf("START %q\n", section.Name)
						for _, command := range section.Commands {
							fmt.Printf("    %-18s", command.Name)
							for _, arg := range command.Args {
								fmt.Printf(" %q", arg)
							}
							fmt.Printf("\n")
						}
					}
				}
			}
		} else {
			errors = append(errors, o.Errors...)
		}
	}
	return errors
}
