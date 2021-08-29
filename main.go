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
	"log"
	"os"
)

//*************************************************************************
// TurnNumber Locations NoOrders Combat PreDeparture Jump Production PostArrival Strike Finish Report Stats

func main() {
	cfg := config.DefaultConfig()
	err := cfg.Load()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}

	log.SetFlags(cfg.Log.Flags)

	if errors := run(cfg); errors != nil {
		for _, err := range errors {
			fmt.Printf("%+v\n", err)
		}
		os.Exit(2)
	}
	//
	//b, err := ioutil.ReadFile("D:\\GoLand\\fhcms\\testdata\\t07\\sp18.ord")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//s, err := scanner.NewScanner(b)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//tokens, err := s.Scan()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, t := range tokens {
	//	fmt.Println(*t)
	//}
}
