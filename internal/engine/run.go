/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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
******************************************************************************/

package engine

import "log"

func (e *Engine) Run() error {
	log.Printf("[engine] running turn      %5d\n", e.galaxy.turn_number)
	log.Printf("[engine] running Locations......\n")
	e.do_locations()
	log.Printf("[engine] created %d/%d locations\n", e.num_locs, len(e.loc))
	log.Printf("[engine] running NoOrders.......\n")
	e.no_orders()
	for i, b := range e.spec_orders {
		if i == 17 && b != nil {
			log.Printf("orders: SP%02d\n%s\n", i+1, string(b))
		}
	}
	log.Printf("[engine] running Combat.........\n")
	log.Printf("[engine] running PreDeparture...\n")
	log.Printf("[engine] running Jump...........\n")
	log.Printf("[engine] running Production.....\n")
	log.Printf("[engine] running PostArrival....\n")
	log.Printf("[engine] running Locations......\n")
	log.Printf("[engine] running Strike.........\n")
	log.Printf("[engine] running Finish.........\n")
	log.Printf("[engine] running Report.........\n")
	log.Printf("[engine] success!\n")
	return nil
}
