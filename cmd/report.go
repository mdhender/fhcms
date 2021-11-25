/*******************************************************************************
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
******************************************************************************/

package cmd

import (
	"bytes"
	"fmt"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func init() {
	rootCmd.AddCommand(reportCmd)
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate report for the current turn.",
	Long:  `Generate the report and default orders for the current turn.`,
	Run: func(cmd *cobra.Command, args []string) {
		spNo, err := cmd.Flags().GetInt("species-no")
		cobra.CheckErr(err)
		fmt.Printf("%-30s == %d\n", "FH_SPECIES_NO", spNo)

		ds, err := loader(viper.GetString("files.path"), viper.GetBool("files.big_endian"))
		cobra.CheckErr(err)

		if spNo < 0 || spNo > len(ds.Species) {
			cobra.CheckErr(fmt.Errorf("species-no must be in range 1..%d", len(ds.Species)))
		}
		var spList []*cluster.Species
		for _, sp := range ds.Species {
			if spNo == sp.No || spNo == 0 {
				spList = append(spList, sp)
			}
		}
		if len(spList) == 0 {
			for _, sp := range ds.Species {
				spList = append(spList, sp)
			}
		}

		// bubble sort the species list
		for i := 0; i < len(spList); i++ {
			for j := i + 1; j < len(spList); j++ {
				if spList[j].No < spList[i].No {
					spList[i], spList[j] = spList[j], spList[i]
				}
			}
		}

		if err := DoReport(ds, spList, ds.Turn, verboseFlag, testFlag); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	reportCmd.Flags().Int("species-no", 0, "species number to generate orders for")
}

func fopen(filename, mode string) (io.Writer, error) {
	if mode != "w" {
		panic("on fopen")
	}
	fmt.Printf("fopen: %s for writing\n", filename)
	return os.Stdout, nil
}

func fprintf(w io.Writer, format string, args ...interface{}) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func DoReport(ds *cluster.Store, spList []*cluster.Species, turn_number int, verbose_mode, test_mode bool) error {
	started := time.Now().UTC()
	// generate report (including default orders) for all species in the list
	for _, sp := range spList {
		species := sp           // todo: use sp directly
		species_number := sp.No // todo: use sp.No directly
		nampla_base := species.NamedPlanets.Base[0]
		if nampla_base == nil {
			panic("assert(nampla_base != nil)")
		}
		//var ship_base *cluster.Ship
		//if len(species.Fleet.Base) > 0 {
		//	ship_base = species.Fleet.Base[0]
		//}
		home_planet := species.HomeWorld.Planet
		if home_planet == nil {
			panic("!")
		}

		// print message for gamemaster
		if verbose_mode {
			fmt.Printf("Generating turn %d report for species #%d, SP %s...\n", turn_number, species_number, species.Name)
		}
		// open report file for writing
		reportFileName := fmt.Sprintf("sp%02d.rpt.t%d", species_number, turn_number)
		report_file := &bytes.Buffer{}
		//report_file, err := fopen(fmt.Sprintf("sp%02d.rpt.t%d", species_number, turn_number), "w")
		//if err != nil {
		//	return err
		//}

		// copy log file, if any, to output file
		if turn_number > 1 {
			fprintf(report_file, "\n\n\t\t\tEVENT LOG FOR TURN %d\n", turn_number-1)
			fprintf(report_file, ";; todo: copy log file sp%02d.rpt.t%d contents\n", species_number, turn_number)
		}

		// print header for status report
		fprintf(report_file, "\n\t\t\t SPECIES STATUS\n\n\t\t\tSTART OF TURN %d\n\n", turn_number)
		fprintf(report_file, "Species name: %s\n", species.Name)
		fprintf(report_file, "Government name: %s\n", species.Government.Name)
		fprintf(report_file, "Government type: %s\n", species.Government.Type)

		fprintf(report_file, "\nTech Levels:\n")
		fprintf(report_file, "   Mining = %d", species.MI.Level)
		if species.MI.KnowledgeLevel > species.MI.Level {
			fprintf(report_file, "/%d\n", species.MI.KnowledgeLevel)
		}
		fprintf(report_file, "\n")
		fprintf(report_file, "   Manufacturing = %d", species.MA.Level)
		if species.MA.KnowledgeLevel > species.MA.Level {
			fprintf(report_file, "/%d\n", species.MA.KnowledgeLevel)
		}
		fprintf(report_file, "\n")
		fprintf(report_file, "   Military = %d", species.ML.Level)
		if species.ML.KnowledgeLevel > species.ML.Level {
			fprintf(report_file, "/%d\n", species.ML.KnowledgeLevel)
		}
		fprintf(report_file, "\n")
		fprintf(report_file, "   Gravitics = %d", species.GV.Level)
		if species.GV.KnowledgeLevel > species.GV.Level {
			fprintf(report_file, "/%d\n", species.GV.KnowledgeLevel)
		}
		fprintf(report_file, "\n")
		fprintf(report_file, "   Life Support = %d", species.LS.Level)
		if species.LS.KnowledgeLevel > species.LS.Level {
			fprintf(report_file, "/%d\n", species.LS.KnowledgeLevel)
		}
		fprintf(report_file, "\n")
		fprintf(report_file, "   Biology = %d", species.BI.Level)
		if species.BI.KnowledgeLevel > species.BI.Level {
			fprintf(report_file, "/%d\n", species.BI.KnowledgeLevel)
		}
		fprintf(report_file, "\n")

		fprintf(report_file, "\nAtmospheric Requirement: %d%%-%d%% %s", species.Gases.RequiredMinPct, species.Gases.RequiredMaxPct, species.Gases.Required.Code)
		fprintf(report_file, "\nNeutral Gases:")
		for i, gas := range species.Gases.Neutral {
			if i != 0 {
				fprintf(report_file, ",")
			}
			fprintf(report_file, " %s", gas.Code)
		}
		fprintf(report_file, "\n")
		fprintf(report_file, "Poisonous Gases:")
		for i, gas := range species.Gases.Poison {
			if i != 0 {
				fprintf(report_file, ",")
			}
			fprintf(report_file, " %s", gas.Code)
		}
		fprintf(report_file, "\n")

		// list fleet maintenance cost and its percentage of total production
		fleet_percent_cost := species.Fleet.MaintenancePct
		fprintf(report_file, "\nFleet maintenance cost = %d (%d.%02d%% of total production)\n", species.Fleet.Cost, fleet_percent_cost/100, fleet_percent_cost%100)
		if fleet_percent_cost > 10000 {
			fleet_percent_cost = 10000
		}

		// list species that have been met
		if len(species.Contact) != 0 {
			var aliens []*cluster.Species
			for _, alien := range species.Contact {
				aliens = append(aliens, alien)
			}
			// sort the list of aliens
			for i := 0; i < len(aliens); i++ {
				for j := i + 1; j < len(aliens); j++ {
					if aliens[j].No < aliens[i].No {
						aliens[i], aliens[j] = aliens[j], aliens[i]
					}
				}
			}
			fprintf(report_file, "\nSpecies met: ")
			for i, alien := range aliens {
				if i > 0 {
					fprintf(report_file, ", ")
				}
				fprintf(report_file, "SP %s", alien.Name)
			}
			fprintf(report_file, "\n")
		}

		// list declared allies
		if len(species.Ally) != 0 {
			var aliens []*cluster.Species
			for _, alien := range species.Contact {
				aliens = append(aliens, alien)
			}
			// sort the list of aliens
			for i := 0; i < len(aliens); i++ {
				for j := i + 1; j < len(aliens); j++ {
					if aliens[j].No < aliens[i].No {
						aliens[i], aliens[j] = aliens[j], aliens[i]
					}
				}
			}
			fprintf(report_file, "\nAllies: ")
			for i, alien := range aliens {
				if i > 0 {
					fprintf(report_file, ", ")
				}
				fprintf(report_file, "SP %s", alien.Name)
			}
			fprintf(report_file, "\n")
		}

		// list declared enemies
		if len(species.Enemy) != 0 {
			var aliens []*cluster.Species
			for _, alien := range species.Enemy {
				aliens = append(aliens, alien)
			}
			// sort the list of aliens
			for i := 0; i < len(aliens); i++ {
				for j := i + 1; j < len(aliens); j++ {
					if aliens[j].No < aliens[i].No {
						aliens[i], aliens[j] = aliens[j], aliens[i]
					}
				}
			}
			fprintf(report_file, "\nEnemies: ")
			for i, alien := range aliens {
				if i > 0 {
					fprintf(report_file, ", ")
				}
				fprintf(report_file, "SP %s", alien.Name)
			}
			fprintf(report_file, "\n")
		}

		fprintf(report_file, "\nEconomic units = %d\n", species.EconUnits)

		// create flag set for ships. if true, the ship has already been reported on / listed in the output
		ship_already_listed := make([]bool, len(species.Fleet.Base), len(species.Fleet.Base))

		// print report for each producing planet
		for i := 0; i < len(species.NamedPlanets.Base); i++ {
			nampla := species.NamedPlanets.Base[i]
			if nampla.Planet.Orbit == 99 {
				continue
			} else if nampla.Colony == nil {
				continue
			} else if nampla.Colony.Mining.Base == 0 && nampla.Colony.Manufacturing.Base == 0 && !nampla.Colony.Is.HomePlanet {
				continue
			}
			planet := nampla.Planet
			fprintf(report_file, "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n")
			do_planet_report(ds, report_file, planet, nampla, species.Fleet.Base, species, home_planet, ship_already_listed)
		}

		// give only a one-line listing for other planets
		printing_alien, header_printed := false, false
		for _, nampla := range species.NamedPlanets.Base {
			if nampla.Planet.Location.Orbit == 99 {
				continue
			} else if nampla.Colony.Mining.Base > 0 || nampla.Colony.Manufacturing.Base > 0 || nampla.Colony.Is.HomePlanet {
				continue
			}
			if !header_printed {
				fprintf(report_file, "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n")
				fprintf(report_file, "\n\nOther planets and ships:\n\n")
				header_printed = true
			}
			fprintf(report_file, "%4d%3d%3d #%d\tPL %s", nampla.Planet.Location.X, nampla.Planet.Location.Y, nampla.Planet.Location.Z, nampla.Planet.Location.Orbit, nampla.Display.Name)
			// print a sorted inventory list
			for _, item := range nampla.Colony.SortedInventory() {
				fprintf(report_file, ", %d %s", item.Quantity, item.Abbr)
			}

			// print any ships at this planet
			for ship_index, ship := range species.Fleet.Base {
				if ship_already_listed[ship_index] || !sameOrbit(nampla.Planet.Location, ship.Location) {
					continue
				}
				fprintf(report_file, "\t\t%s", ship.Named(false, !printing_alien))
				// print a sorted inventory list
				for _, item := range ship.SortedInventory() {
					fprintf(report_file, ", %d %s", item.Quantity, item.Abbr)
				}
				fprintf(report_file, "\n")
				ship_already_listed[ship_index] = true
			}

			fprintf(report_file, "\n")
		}

		// list all ships that are not associated with a planet
		for ship_index, ship := range species.Fleet.Base {
			ship.Special = 0 // why?
			if ship_already_listed[ship_index] {
				continue
			}
			ship_already_listed[ship_index] = true
			if ship.Location.Orbit == 99 {
				continue
			}
			if !header_printed {
				fprintf(report_file, "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n")
				fprintf(report_file, "\n\nOther planets and ships:\n\n")
				header_printed = true
			}
			shipName := ship.Named(false, !printing_alien)
			if ship.Status.JumpedInCombat || ship.Status.ForcedJump {
				fprintf(report_file, "  ?? ?? ??\t%s", shipName)
			} else if test_mode && ship.ArrivedViaWormhole {
				fprintf(report_file, "  ?? ?? ??\t%s", shipName)
			} else {
				fprintf(report_file, "%4d%3d%3d\t%s", ship.Location.X, ship.Location.Y, ship.Location.Z, shipName)
			}
			// print a sorted inventory list
			for _, item := range ship.SortedInventory() {
				fprintf(report_file, ", %d %s", item.Quantity, item.Abbr)
			}
			fprintf(report_file, "\n")

			if ship.Status.JumpedInCombat || ship.Status.ForcedJump {
				continue
			} else if test_mode && ship.ArrivedViaWormhole {
				continue
			}

			// print other ships at the same location
			for i := ship_index + 1; i < len(species.Fleet.Base); i++ {
				ship2 := species.Fleet.Base[i]
				if ship_already_listed[i] || ship2.Location.Orbit == 99 || !sameSystem(ship.Location, ship2.Location) {
					continue
				}
				fprintf(report_file, "\t\t%s", ship2.Named(false, !printing_alien))
				for _, item := range ship.SortedInventory() {
					fprintf(report_file, ", %d %s", item.Quantity, item.Abbr)
				}
			}
		}

		fprintf(report_file, "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n")

		// report aliens at locations where current species has inhabited planets or ships
		printing_alien = true
		for _, my_loc := range ds.Locations {
			if my_loc.Species.No != species_number {
				continue
			}
			header_printed = false
			for _, its_loc := range ds.Locations {
				// is this our system and is there an alien here?
				if its_loc.Species.No == species_number || !sameSystem(my_loc.Location, its_loc.Location) {
					continue
				}
				// there is an alien here!
				alien := its_loc.Species
				alien_number := alien.No
				//nampla2_base := alien.NamedPlanets.Base
				//ship2_base := alien.Fleet.Base

				// check if we have a named planet in this system. if so, use it when we print the header
				var our_nampla *cluster.NamedPlanet
				we_have_planet_here := false
				for _, np := range species.NamedPlanets.Base {
					if sameSystem(its_loc.Location, np.Planet.Location) {
						our_nampla, we_have_planet_here = np, true
						break
					}
				}

				// print all inhabited alien namplas at this location
				var alien_nampla *cluster.NamedPlanet
				for i := range alien.NamedPlanets.Base {
					alien_nampla = alien.NamedPlanets.Base[i]
					if alien_nampla.Colony == nil || !alien_nampla.Colony.Is.Populated || !sameSystem(my_loc.Location, alien_nampla.Planet.Location) {
						continue
					}
					// check if current species has a colony on the same planet
					we_have_colony_here := false
					for _, nampla := range species.NamedPlanets.Base {
						if nampla.Colony == nil || !nampla.Colony.Is.Populated || !sameOrbit(alien_nampla.Planet.Location, nampla.Planet.Location) {
							continue
						}
						we_have_colony_here = true
						break
					}
					if alien_nampla.Colony.Is.Hidden && !we_have_colony_here {
						continue
					}
					if !header_printed {
						fprintf(report_file, "\n\nAliens at x = %d, y = %d, z = %d", my_loc.Location.X, my_loc.Location.Y, my_loc.Location.Z)
						if we_have_planet_here {
							fprintf(report_file, " (PL %s star system)", our_nampla.Display.Name)
						}
						fprintf(report_file, ":\n")
						header_printed = true
					}

					industry := alien_nampla.Colony.Mining.Base + alien_nampla.Colony.Manufacturing.Base
					var temp1 string
					if alien_nampla.Colony.Is.MiningColony {
						temp1 = "Mining colony"
					} else if alien_nampla.Colony.Is.ResortColony {
						temp1 = "Resort colony"
					} else if alien_nampla.Colony.Is.HomePlanet {
						temp1 = "Home planet"
					} else if industry > 0 {
						temp1 = "Colony planet"
					} else {
						temp1 = "Uncolonized planet"
					}
					temp2 := fmt.Sprintf("  %s PL %s (pl #%d)", temp1, alien_nampla.Display.Name, alien_nampla.Planet.Location.Orbit)
					fprintf(report_file, "%-53sSP %s\n", temp2, alien.Name)

					if industry <= 0 {
						fprintf(report_file, "      (No economic base.)\n")
					} else if industry < 100 {
						fprintf(report_file, "      (Economic base is approximately %d.)\n", (industry+5)/10)
					} else {
						fprintf(report_file, "      (Economic base is approximately %d.)\n", ((industry+50)/100)*10)
					}

					// if current species has a colony on the same planet, report any PDs and any shipyards
					if we_have_colony_here {
						if item, ok := alien_nampla.Colony.Inventory["PD"]; ok && item.Quantity > 0 {
							if item.Quantity == 1 {
								fprintf(report_file, "      (There is 1 %s on the planet.)\n", item.Descr)
							} else {
								fprintf(report_file, "      (There are %d %ss on the planet.)\n", item.Quantity, item.Descr)
							}
						}
						if alien_nampla.Colony.Shipyards == 1 {
							fprintf(report_file, "      (There is 1 shipyard on the planet.)\n")
						} else if alien_nampla.Colony.Shipyards > 1 {
							fprintf(report_file, "      (There are %d shipyards on the planet.)\n", alien_nampla.Colony.Shipyards)
						}
					}

					// also report if alien colony is actively hiding
					if alien_nampla.Colony.Is.Hidden {
						fprintf(report_file, "      (Colony is actively hiding from alien observation.)\n")
					}
				}

				// print all alien ships at this location
				for _, alien_ship := range alien.Fleet.Base {
					if alien_ship.Location.Orbit == 99 || !sameSystem(my_loc.Location, alien_ship.Location) {
						continue
					}

					// an alien ship cannot hide if it lands on the surface of a planet populated by the current species
					alien_can_hide := false
					for _, nampla := range species.NamedPlanets.Base {
						if !sameOrbit(alien_ship.Location, nampla.Planet.Location) {
							continue
						}
						if nampla.Colony != nil && nampla.Colony.Is.Populated {
							alien_can_hide = false
							break
						}
					}

					if alien_can_hide && (alien_ship.Status.OnSurface || alien_ship.Status.UnderConstruction) {
						continue
					}

					if !header_printed {
						fprintf(report_file, "\n\nAliens at x = %d, y = %d, z = %d", my_loc.Location.X, my_loc.Location.Y, my_loc.Location.Z)
						if we_have_planet_here {
							fprintf(report_file, " (PL %s star system)", our_nampla.Display.Name)
						}
						fprintf(report_file, ":\n")
						header_printed = true
					}

					print_ship(report_file, alien_ship, alien, alien_number, printing_alien)
				}
			}
		}

		printing_alien = false

		if test_mode {
			return nil
		}

		// generating orders section
		fprintf(report_file, "generating orders for species %s, SP %s...\n", species.Id, species.Name)
		fprintf(report_file, ";; %s T%d %v\n", sp.Id, turn_number, started)
		fprintf(report_file, ";; report orders\n")

		// print out ship location and inventory
		if len(sp.Fleet.Ships) != 0 {
			fprintf(report_file, ";; Fleet Data\n")
			for _, ship := range sp.Fleet.Ships {
				fprintf(report_file, ";;   %-45s  %3d %3d %3d", ship.Display.Name, ship.Location.X, ship.Location.Y, ship.Location.Z)
				if ship.Location.Orbit == 0 {
					fprintf(report_file, "   ")
				} else {
					fprintf(report_file, " #%d", ship.Location.Orbit)
				}
				if ship.Special != 0 {
					fprintf(report_file, "  special %d", ship.Special)
				}
				if ship.UnloadingPoint != 0 {
					fprintf(report_file, "  unload %d", ship.UnloadingPoint)
				}
				if ship.Status.UnderConstruction {
					fprintf(report_file, "  Under Construction\n")
				} else {
					fprintf(report_file, "  Age %2d\n", ship.Age)
					for _, item := range ship.Inventory {
						fprintf(report_file, ";;       %5d %-3s %s\n", item.Quantity, item.Abbr, item.Descr)
					}
				}
			}
		}
		fprintf(report_file, "\n\n")

		fprintf(report_file, "\n\n* * * * * * * * * * * * * * * * * * * * * * * * *\n")

		fprintf(report_file, "\n\nORDER SECTION. Remove these two lines and everything above\n")
		fprintf(report_file, "  them, and submit only the orders below.\n\n")

		fprintf(report_file, "START COMBAT\n")
		fprintf(report_file, "; Place combat orders here.\n\n")
		fprintf(report_file, "END\n\n")

		// PRE-DEPARTURE orders
		fprintf(report_file, "START PRE-DEPARTURE\n")
		fprintf(report_file, "; Place pre-departure orders here.\n\n")

		for nampla_index := 0; nampla_index < len(sp.NamedPlanets.Base); nampla_index++ {
			nampla := species.NamedPlanets.Base[nampla_index]
			if nampla == nil || nampla.Colony == nil {
				continue
			}
			colony := nampla.Colony
			// generate auto-installs for colonies that were loaded via the DEVELOP command
			if colony.Mining.AutoIUs > 0 {
				fprintf(report_file, "\tInstall\t%d IU\tPL %s\n", colony.Mining.AutoIUs, colony.Name.Display.Name)
			}
			if colony.Manufacturing.AutoAUs > 0 {
				fprintf(report_file, "\tInstall\t%d AU\tPL %s\n", colony.Manufacturing.AutoAUs, colony.Name.Display.Name)
			}
			if colony.Mining.AutoIUs > 0 || colony.Manufacturing.AutoAUs > 0 {
				fprintf(report_file, "\n")
			}

			if !species.AutoOrders {
				continue
			}

			// generate auto UNLOAD orders for transports at this colony
			for j := 0; j < len(species.Fleet.Base); j++ {
				ship := species.Fleet.Base[j]
				if ship.Location.Orbit == 99 {
					continue
				}
				if ship.Location.X != nampla.Planet.Location.X {
					continue
				}
				if ship.Location.Y != nampla.Planet.Location.Y {
					continue
				}
				if ship.Location.Z != nampla.Planet.Location.Z {
					continue
				}
				if ship.Location.Orbit != nampla.Planet.Location.Orbit {
					continue
				}
				if ship.Status.JumpedInCombat {
					continue
				}
				if ship.Status.ForcedJump {
					continue
				}
				if !ship.Class.Is.Transport {
					continue
				}
				if item, ok := ship.Inventory["CU"]; !ok || item.Quantity < 1 {
					continue
				}

				// new colonies will never be started automatically unless ship was loaded via a DEVELOP order
				var n int
				if ship.LoadingPoint != 0 {
					// check if transport is at specified unloading point
					n = ship.UnloadingPoint
					if n == nampla_index || (n == 9999 && nampla_index == 0) {
						goto unload_ship
					}
				}

				if !colony.Is.Populated {
					continue
				}

				if colony.Mining.Base+colony.Manufacturing.Base >= 2000 {
					continue
				}

				if nampla.Planet.Location.X == nampla_base.Planet.Location.X && nampla.Planet.Location.Y == nampla_base.Planet.Location.Y && nampla.Planet.Location.Z == nampla_base.Planet.Location.Z {
					continue // home sector
				}

			unload_ship:
				n = ship.LoadingPoint
				if n == 9999 {
					n = 0 // home planet
				}
				if n == nampla_index {
					continue // ship was just loaded here
				}

				fprintf(report_file, "\tUnload\t%s\n\n", ship.Display.Name)

				ship.Special = ship.LoadingPoint
				n = nampla.Index - nampla_base.Index
				if n == 0 {
					n = 9999
				}
				ship.UnloadingPoint = n
			}
		}
		fprintf(report_file, "END\n\n")

		// generate jump orders for ships used to develop and scouts
		fprintf(report_file, "START JUMPS\n")
		fprintf(report_file, "; Place jump orders here.\n\n")

		// initialize to make sure ships are not given more than one JUMP order
		for _, ship := range sp.Fleet.Ships {
			ship.JustJumped = false
		}

		// generate auto-jumps for ships that were loaded via the DEVELOP command or which were UNLOADed because of the AUTO command
		for i := 0; i < len(species.Fleet.Base); i++ {
			ship := species.Fleet.Base[i]
			ship.JustJumped = false

			if ship.Location == nil || ship.Location.Orbit == 99 {
				continue
			} else if ship.Status.JumpedInCombat {
				continue
			} else if ship.Status.ForcedJump {
				continue
			}

			// ship.Special was set to the unloading point in one of the previous steps
			j := ship.Special
			if j != 0 {
				if j == 9999 { // home planet
					j = 0
				}
				temp_nampla := sp.NamedPlanets.Base[j]
				_, mishapChance := cluster.MishapChance(sp, ship, temp_nampla.Planet.Location)
				fprintf(report_file, "\tJump\t%s, PL %s\t; Age %d, mishap chance = %s\n\n", ship.Display.Name, temp_nampla.Display.Name, ship.Age, mishapChance)
				ship.JustJumped = true
				continue
			}

			n := ship.UnloadingPoint
			if n != 0 {
				if n == 9999 {
					n = 0 // home planet
				}
				temp_nampla := sp.NamedPlanets.Base[n]
				_, mishapChance := cluster.MishapChance(sp, ship, temp_nampla.Planet.Location)
				fprintf(report_file, "\tJump\t%s, PL %s\t; mishap chance = %s\n\n", ship.Display.Name, temp_nampla.Display.Name, mishapChance)
				ship.JustJumped = true
				continue
			}
		}

		if !species.AutoOrders {
			goto jump_end
		}

		// generate JUMP orders for all ships that have not yet been given orders
		for i := 0; i < len(species.Fleet.Base); i++ {
			ship := sp.Fleet.Base[i]
			if ship.Location.Orbit == 99 {
				continue
			} else if ship.JustJumped {
				continue
			} else if ship.Status.UnderConstruction {
				continue
			} else if ship.Status.JumpedInCombat {
				continue
			} else if ship.Status.ForcedJump {
				continue
			}

			if ship.Class.Is.SubLight {
				// todo: calculate delta x, y, or z that moves us closer to that system for sublight ships
			} else {
				var closestSystem *cluster.Coords
				fprintf(report_file, "\tJump\t%s, ", ship.Display.Name)
				if ship.Class.Is.Transport && ship.Class.Tonnage == 1 {
					closestSystem = ds.ClosestUnvisitedSystem(sp, ship.Location)
					_, mishapChance := cluster.MishapChance(sp, ship, closestSystem)
					fprintf(report_file, "\n\t\t\t; Age %d, now at %d %d %d, ", ship.Age, ship.Location.X, ship.Location.Y, ship.Location.Z)

					if ship.Status.InOrbit {
						fprintf(report_file, "O%d, ", ship.Location.Orbit)
					} else if ship.Status.OnSurface {
						fprintf(report_file, "L%d, ", ship.Location.Orbit)
					} else {
						fprintf(report_file, "D, ")
					}

					fprintf(report_file, "mishap chance = %s\n\n", mishapChance)

				} else {
					fprintf(report_file, "???\t; Age %d, now at %d %d %d", ship.Age, ship.Location.X, ship.Location.Y, ship.Location.Z)

					if ship.Status.InOrbit {
						fprintf(report_file, "O%d, ", ship.Location.Orbit)
					} else if ship.Status.OnSurface {
						fprintf(report_file, "L%d, ", ship.Location.Orbit)
					} else {
						fprintf(report_file, "D, ")
					}

					closestSystem = nil
				}

				fprintf(report_file, "\n")

				// save destination so that we can check later if it needs to be scanned
				if closestSystem == nil {
					ship.Destination = cluster.NewCoords(-1, -1, -1, -1)
				} else {
					ship.Destination = closestSystem
				}
			}
		}
	jump_end:
		fprintf(report_file, "END\n\n")

		fprintf(report_file, "START PRODUCTION\n\n")
		fprintf(report_file, ";   Economic units at start of turn = %d\n\n", species.EconUnits)

		// generate a PRODUCTION order for each planet that can produce, starting with newest colony
		for nampla_index := len(species.NamedPlanets.Base) - 1; nampla_index >= 0; nampla_index-- {
			nampla := species.NamedPlanets.Base[nampla_index]
			colony := nampla.Colony
			if colony == nil || nampla.Planet.Location.Orbit == 99 {
				continue
			}
			if colony.Mining.Base == 0 && !colony.Is.ResortColony {
				continue
			}
			if colony.Manufacturing.Base == 0 && !colony.Is.MiningColony {
				continue
			}

			fprintf(report_file, "    PRODUCTION PL %s\n", colony.Name.Display.Name)

			if colony.Is.MiningColony {
				fprintf(report_file, "    ; The above PRODUCTION order is required for this mining colony, even\n")
				fprintf(report_file, "    ;  if no other production orders are given for it. This mining colony\n")
				fprintf(report_file, "    ;  will generate %d economic units this turn.\n", colony.UseOnAmbush)
			} else if colony.Is.ResortColony {
				fprintf(report_file, "    ; The above PRODUCTION order is required for this resort colony, even\n")
				fprintf(report_file, "    ;  though no other production orders can be given for it.  This resort\n")
				fprintf(report_file, "    ;  colony will generate %d economic units this turn.\n", colony.UseOnAmbush)
			} else {
				fprintf(report_file, "    ; Place production orders here for planet %s", colony.Name.Display.Name)
				fprintf(report_file, " (sector %d %d %d #%d).\n", colony.Planet.Location.X, colony.Planet.Location.Y, colony.Planet.Location.Z, colony.Planet.Location.Orbit)
				fprintf(report_file, "    ;  Avail pop = %d, shipyards = %d, to spend = %d", colony.Population, colony.Shipyards, colony.UseOnAmbush)

				n := colony.UseOnAmbush
				if colony.Is.HomePlanet {
					if species.HomeWorld.OriginalBase != 0 {
						fprintf(report_file, " (max = %d)", 5*n)
					} else {
						fprintf(report_file, " (max = no limit)")
					}
				} else {
					fprintf(report_file, " (max = %d)", 2*n)
				}

				fprintf(report_file, ".\n\n")
			}

			// build IUs and AUs for incoming ships with CUs
			if colony.Mining.Needed > 0 {
				fprintf(report_file, "\tBuild\t%d IU\n", colony.Mining.Needed)
			}
			if colony.Manufacturing.Needed > 0 {
				fprintf(report_file, "\tBuild\t%d AU\n", colony.Manufacturing.Needed)
			}
			if colony.Mining.Needed > 0 || colony.Manufacturing.Needed > 0 {
				fprintf(report_file, "\n")
			}

			if !species.AutoOrders {
				continue
			}
			if colony.Is.MiningColony {
				continue
			}
			if colony.Is.ResortColony {
				continue
			}

			// see if there are any RMs to recycle
			n := colony.Special / 5
			if n > 0 {
				fprintf(report_file, "\tRecycle\t%d RM\n\n", 5*n)
			}

			// generate DEVELOP commands for ships arriving here because of AUTO command
			for i := 0; i < len(species.Fleet.Base); i++ {
				ship := species.Fleet.Base[i]
				if ship.Location.Orbit == 99 {
					continue
				}

				k := ship.Special
				if k == 0 {
					continue
				}
				if k == 9999 {
					k = 0 // home planet
				}

				if colony.Name.Index != nampla_base.Index+k {
					continue
				}

				k = ship.UnloadingPoint
				if k == 9999 {
					k = 0
				}
				temp_nampla := species.NamedPlanets.Base[k]

				fprintf(report_file, "\tDevelop\tPL %s, %s\n\n", temp_nampla.Display.Name, ship.Display.Name)
			}

			// give orders to continue construction of unfinished ships and starbases
			for i := 0; i < len(species.Fleet.Base); i++ {
				ship := species.Fleet.Base[i]
				if ship.Location.Orbit == 99 {
					continue
				}
				if ship.Location.X != colony.Planet.Location.X {
					continue
				}
				if ship.Location.Y != colony.Planet.Location.Y {
					continue
				}
				if ship.Location.Z != colony.Planet.Location.Z {
					continue
				}
				if ship.Location.Orbit != colony.Planet.Location.Orbit {
					continue
				}

				if ship.Status.UnderConstruction {
					fprintf(report_file, "\tContinue\t%s, %d\t; Left to pay = %d\n\n", ship.Display.Name, ship.RemainingCost, ship.RemainingCost)
					continue
				}

				if !ship.Class.Is.Starbase {
					continue
				}

				j := (species.MA.Level / 2) - ship.Class.Tonnage
				if j < 1 {
					continue
				}

				fprintf(report_file, "\tContinue\tBAS %s, %d\t; Current tonnage = %s\n\n", ship.Display.Name, 100*j, commas(10000*ship.Class.Tonnage))
			}

			// generate DEVELOP command if this is a colony with an economic base less than 200
			n = colony.Mining.Base + colony.Manufacturing.Base + colony.Mining.Needed + colony.Manufacturing.Needed
			var nn int
			if item, ok := colony.Inventory["CU"]; ok && item.Quantity > 0 {
				nn = item.Quantity
			}
			// get CUs on transports at planet
			for i := 0; i < len(species.Fleet.Base); i++ {
				ship := species.Fleet.Base[i]
				if ship.Location.Orbit == 99 {
					continue
				}
				if ship.Location.X != colony.Planet.Location.X {
					continue
				}
				if ship.Location.Y != colony.Planet.Location.Y {
					continue
				}
				if ship.Location.Z != colony.Planet.Location.Z {
					continue
				}
				if ship.Location.Orbit != colony.Planet.Location.Orbit {
					continue
				}
				if item, ok := ship.Inventory["CU"]; ok && item.Quantity > 0 {
					nn += item.Quantity
				}
			}
			n += nn
			if colony.Is.Colony && n < 2000 && colony.Population > 0 {
				if colony.Population > (2000 - n) {
					nn = 2000 - n
				} else {
					nn = colony.Population
				}
				fprintf(report_file, "\tDevelop\t%d\n\n", 2*nn)
				colony.Mining.Needed += nn
			}

			// for home planets and any colonies that have an economic base of at least 200,
			// check if there are other colonized planets in the same sector that are not
			// self-sufficient. if so, DEVELOP them.
			if n >= 2000 || colony.Is.HomePlanet {
				// skip home planet
				for i := 1; i < len(species.NamedPlanets.Base); i++ {
					if i == nampla_index {
						continue
					}
					temp_nampla := species.NamedPlanets.Base[i]
					if temp_nampla.Colony == nil || temp_nampla.Planet.Orbit == 99 {
						continue
					}
					if temp_nampla.Planet.Location.X != nampla.Planet.Location.X {
						continue
					}
					if temp_nampla.Planet.Location.Y != nampla.Planet.Location.Y {
						continue
					}
					if temp_nampla.Planet.Location.Z != nampla.Planet.Location.Z {
						continue
					}

					n = temp_nampla.Colony.Mining.Base + temp_nampla.Colony.Manufacturing.Base + temp_nampla.Colony.Mining.Needed + temp_nampla.Colony.Manufacturing.Needed
					if n == 0 {
						continue
					}

					nn := 0
					if item, ok := temp_nampla.Colony.Inventory["IU"]; ok && item.Quantity > 0 {
						nn += item.Quantity
					}
					if item, ok := temp_nampla.Colony.Inventory["AU"]; ok && item.Quantity > 0 {
						nn += item.Quantity
					}
					if item, ok := temp_nampla.Colony.Inventory["CU"]; !ok {
						nn = 0
					} else if nn > item.Quantity {
						nn = item.Quantity
					}
					n += nn
					if n >= 2000 {
						continue
					}
					nn = 2000 - n

					if nn > colony.Population {
						nn = colony.Population
					}

					fprintf(report_file, "\tDevelop\t%d\tPL %s\n\n", 2*nn, temp_nampla.Display.Name)

					temp_nampla.Colony.Manufacturing.Needed += nn
				}
			}
		}

		fprintf(report_file, "END\n\n")

		fprintf(report_file, "START POST-ARRIVAL\n")
		fprintf(report_file, "; Place post-arrival orders here.\n\n")

		if species.AutoOrders {
			/* Generate an AUTO command. */
			fprintf(report_file, "\tAuto\n\n")
			/* Generate SCAN orders for all TR1s that are jumping to sectors which current species does not inhabit. */
			for i := 0; i < len(species.Fleet.Base); i++ {
				ship := species.Fleet.Base[i]
				if ship.Location.Orbit == 99 {
					continue
				}
				if ship.Status.UnderConstruction {
					continue
				}
				if !ship.Class.Is.Transport {
					continue
				}
				if ship.Class.Tonnage != 1 {
					continue
				}
				if ship.Class.Is.SubLight {
					continue
				}

				found := false
				if ship.Destination != nil && ship.Destination.X != -1 {
					for j := 0; j < len(species.NamedPlanets.Base); j++ {
						nampla := species.NamedPlanets.Base[j]
						if nampla.Planet.Location.Orbit == 99 {
							continue
						}
						if nampla.Planet.Location.X != ship.Destination.X {
							continue
						}
						if nampla.Planet.Location.Y != ship.Destination.Y {
							continue
						}
						if nampla.Planet.Location.Z != ship.Destination.Z {
							continue
						}
						if nampla.Colony != nil && nampla.Colony.Is.Populated {
							found = true
							break
						}
					}
				}
				if !found {
					fprintf(report_file, "\tScan\t%s\n", ship.Display.Name)
				}
			}
		}
		fprintf(report_file, "END\n\n")

		fprintf(report_file, "START STRIKES\n")
		fprintf(report_file, "; Place strike orders here.\n\n")
		fprintf(report_file, "END\n")

		fmt.Println(reportFileName)
		if reportFileName == "sp18.rpt.t28" {
			_ = os.WriteFile(filepath.Join("D:\\GoLand\\fhcms\\testdata", reportFileName+".txt"), report_file.Bytes(), 0600)
		}
	}

	return nil
}

func do_planet_report(ds *cluster.Store, report_file io.Writer, planet *cluster.Planet, nampla *cluster.NamedPlanet, s_base []*cluster.Ship, species *cluster.Species, home_planet *cluster.Planet, ship_already_listed []bool) {
	if nampla == nil {
		panic("assert(nampla != nil)")
	}
	colony := nampla.Colony
	if colony == nil {
		panic("assert(colony != nil)")
	}

	var header_printed, ok bool
	var i, j, ship_index, ls_needed, production_penalty int
	var n1, n2, n3, raw_material_units, production_capacity, available_to_spend, n, ib, ab, current_base, md, denom int
	var item *cluster.Item
	var ship *cluster.Ship

	/* Print type of planet, name and coordinates. */
	fprintf(report_file, "\n\n")

	if colony.Is.HomePlanet {
		fprintf(report_file, "HOME PLANET")
	} else if colony.Is.MiningColony {
		fprintf(report_file, "MINING COLONY")
	} else if colony.Is.ResortColony {
		fprintf(report_file, "RESORT COLONY")
	} else if colony.Is.Populated {
		fprintf(report_file, "COLONY PLANET")
	} else {
		fprintf(report_file, "PLANET")
	}
	fprintf(report_file, ": PL %s", colony.Name.Display.Name)

	fprintf(report_file, "\n   Coordinates: x = %d, y = %d, z = %d, planet number %d\n", colony.Planet.Location.X, colony.Planet.Location.Y, colony.Planet.Location.Z, colony.Planet.Location.Orbit)

	if colony.Is.HomePlanet {
		ib = colony.Mining.Base
		ab = colony.Manufacturing.Base
		current_base = ib + ab
		if current_base < species.HomeWorld.OriginalBase {
			n = species.HomeWorld.OriginalBase - current_base /* Number of CUs needed. */
			md = home_planet.MiningDifficulty

			denom = 100 + md
			j = (100*(n+ib) - (md * ab) + denom/2) / denom
			i = n - j

			if i < 0 {
				j = n
				i = 0
			}
			if j < 0 {
				i = n
				j = 0
			}

			fprintf(report_file, "\nWARNING! Home planet has not yet completely recovered from bombardment!\n")
			fprintf(report_file, "         %d IUs and %d AUs will have to be installed for complete recovery.\n", i, j)
		}
	}

	if !(colony.Is.Populated) {
		goto do_inventory
	}

	/* Print available population. */
	if colony.Is.MiningColony || colony.Is.ResortColony {
		// no population to report
	} else {
		fprintf(report_file, "\nAvailable population units = %d\n", colony.Population)
	}

	if colony.SiegeEff != 0 {
		fprintf(report_file, "\nWARNING!  This planet is currently under siege and will remain\n")
		fprintf(report_file, "  under siege until the combat phase of the next turn!\n")
	}

	if colony.UseOnAmbush > 0 {
		fprintf(report_file, "\nIMPORTANT!  This planet has made preparations for an ambush!\n")
	}

	if colony.Is.Hidden {
		fprintf(report_file, "\nIMPORTANT!  This planet is actively hiding from alien observation!\n")
	}

	/* Print what will be produced this turn. */
	raw_material_units = (10 * species.MI.Level * colony.Mining.Base) / planet.MiningDifficulty
	production_capacity = (species.MA.Level * colony.Manufacturing.Base) / 10

	ls_needed = lifeSupportNeeded(species, planet)

	if ls_needed == 0 {
		production_penalty = 0
	} else {
		production_penalty = (100 * ls_needed) / species.LS.Level
	}

	fprintf(report_file, "\nProduction penalty = %d%% (LSN = %d)\n", production_penalty, ls_needed)

	fprintf(report_file, "\nEconomic efficiency = %d%%\n", planet.EconEfficiency)

	raw_material_units -= (production_penalty * raw_material_units) / 100
	raw_material_units = ((planet.EconEfficiency * raw_material_units) + 50) / 100
	production_capacity -= (production_penalty * production_capacity) / 100
	production_capacity = ((planet.EconEfficiency * production_capacity) + 50) / 100

	if colony.Mining.Base > 0 {
		fprintf(report_file, "\nMining base = %d.%d", colony.Mining.Base/10, colony.Mining.Base%10)
		fprintf(report_file, " (MI = %d, MD = %d.%02d)\n", species.MI.Level, planet.MiningDifficulty/100, planet.MiningDifficulty%100)

		/* For mining colonies, print economic units that will be produced. */
		if colony.Is.MiningColony {
			n1 = (2 * raw_material_units) / 3
			n2 = ((species.Fleet.MaintenancePct * n1) + 5000) / 10000
			n3 = n1 - n2
			fprintf(report_file, "   This mining colony will generate %d - %d = %d economic units this turn.\n", n1, n2, n3)
			colony.UseOnAmbush = n3 /* Temporary use only. */
		} else {
			fprintf(report_file, "   %d raw material units will be produced this turn.\n", raw_material_units)
		}
	}

	if colony.Manufacturing.Base > 0 {
		if colony.Is.ResortColony {
			fprintf(report_file, "\n")
		}
		fprintf(report_file, "Manufacturing base = %d.%d", colony.Manufacturing.Base/10, colony.Manufacturing.Base%10)
		fprintf(report_file, " (MA = %d)\n", species.MA.Level)
		/* For resort colonies, print economic units that will be produced. */
		if colony.Is.ResortColony {
			n1 = (2 * production_capacity) / 3
			n2 = ((species.Fleet.MaintenancePct * n1) + 5000) / 10000
			n3 = n1 - n2
			fprintf(report_file, "   This resort colony will generate %d - %d = %d economic units this turn.\n", n1, n2, n3)
			colony.UseOnAmbush = n3 /* Temporary use only. */
		} else {
			fprintf(report_file, "   Production capacity this turn will be %d.\n", production_capacity)
		}
	}

	if item, ok = colony.Inventory["RM"]; ok && item.Quantity > 0 {
		fprintf(report_file, "\n%ss (%s,C%d) carried over from last turn = %d\n", item.Descr, item.Abbr, item.CarryCapacity, item.Quantity)
	}

	/* Print what can be spent this turn. */
	if ok {
		raw_material_units += item.Quantity
	}
	if raw_material_units > production_capacity {
		available_to_spend = production_capacity
		colony.Special = raw_material_units - production_capacity
		/* Excess raw material units that may be recycled in AUTO mode. */
	} else {
		available_to_spend = raw_material_units
		colony.Special = 0
	}

	/* Don't print spendable amount for mining and resort colonies. */
	n1 = available_to_spend
	n2 = ((species.Fleet.MaintenancePct * n1) + 5000) / 10000
	n3 = n1 - n2
	if !(colony.Is.MiningColony) && !(colony.Is.ResortColony) {
		fprintf(report_file, "\nTotal available for spending this turn = %d - %d = %d\n", n1, n2, n3)
		colony.UseOnAmbush = n3 /* Temporary use only. */

		fprintf(report_file, "\nShipyard capacity = %d\n", colony.Shipyards)
	}

do_inventory:

	header_printed = false

	for _, item := range colony.SortedInventory() {
		if item.Quantity > 0 && item.Abbr != "RM" {
			if !header_printed {
				header_printed = true
				fprintf(report_file, "\nPlanetary inventory:\n")
			}
			fprintf(report_file, "   %ss (%s,C%d) = %d", item.Descr, item.Abbr, item.CarryCapacity, item.Quantity)
			if item.Abbr == "PD" {
				fprintf(report_file, " (warship equivalence = %d tons)", 50*item.Quantity)
			}
			fprintf(report_file, "\n")
		}
	}

	/* Print all ships that are under construction on, on the surface of,
	   or in orbit around this planet. */
	printing_alien := false
	header_printed = false
	for ship_index = 0; ship_index < len(species.Fleet.Base); ship_index++ {
		ship = species.Fleet.Base[ship_index]

		if colony.Planet.Location.X != ship.Location.X {
			continue
		}
		if colony.Planet.Location.Y != ship.Location.Y {
			continue
		}
		if colony.Planet.Location.Z != ship.Location.Z {
			continue
		}
		if colony.Planet.Location.Orbit != ship.Location.Orbit {
			continue
		}
		if !ship.Class.Is.Starbase {
			continue
		}

		if !header_printed {
			fprintf(report_file, "\nShips at PL %s:\n", colony.Name.Display.Name)
			fprintf(report_file, "  Name                          ")
			if printing_alien {
				fprintf(report_file, "                     Species\n")
			} else {
				fprintf(report_file, "                 Cap. Cargo\n")
				fprintf(report_file, " ---------------------------------------")
				fprintf(report_file, "-------------------------------------\n")
			}

		}
		header_printed = true

		print_ship(report_file, ship, species, species.No, printing_alien)

		ship_already_listed[ship_index] = true
	}

	for ship_index = 0; ship_index < len(species.Fleet.Base); ship_index++ {
		ship = species.Fleet.Base[ship_index]

		if colony.Planet.Location.X != ship.Location.X {
			continue
		}
		if colony.Planet.Location.Y != ship.Location.Y {
			continue
		}
		if colony.Planet.Location.Z != ship.Location.Z {
			continue
		}
		if colony.Planet.Location.Orbit != ship.Location.Orbit {
			continue
		}
		if !ship.Class.Is.Transport {
			continue
		}

		if !header_printed {
			fprintf(report_file, "\nShips at PL %s:\n", colony.Name.Display.Name)
			fprintf(report_file, "  Name                          ")
			if printing_alien {
				fprintf(report_file, "                     Species\n")
			} else {
				fprintf(report_file, "                 Cap. Cargo\n")
				fprintf(report_file, " ---------------------------------------")
				fprintf(report_file, "-------------------------------------\n")
			}

		}
		header_printed = true

		print_ship(report_file, ship, species, species.No, printing_alien)

		ship_already_listed[ship_index] = true
	}

	for ship_index = 0; ship_index < len(species.Fleet.Base); ship_index++ {
		ship = species.Fleet.Base[ship_index]

		if colony.Planet.Location.X != ship.Location.X {
			continue
		}
		if colony.Planet.Location.Y != ship.Location.Y {
			continue
		}
		if colony.Planet.Location.Z != ship.Location.Z {
			continue
		}
		if colony.Planet.Location.Orbit != ship.Location.Orbit {
			continue
		}
		if ship_already_listed[ship_index] {
			continue
		}

		if !header_printed {
			fprintf(report_file, "\nShips at PL %s:\n", colony.Name.Display.Name)
			fprintf(report_file, "  Name                          ")
			if printing_alien {
				fprintf(report_file, "                     Species\n")
			} else {
				fprintf(report_file, "                 Cap. Cargo\n")
				fprintf(report_file, " ---------------------------------------")
				fprintf(report_file, "-------------------------------------\n")
			}
		}
		header_printed = true

		print_ship(report_file, ship, species, species.No, printing_alien)

		ship_already_listed[ship_index] = true
	}
}

func print_ship(report_file io.Writer, ship *cluster.Ship, species *cluster.Species, species_number int, printing_alien bool) {
	shipName := ship.Named(false, !printing_alien)
	if printing_alien {
		fprintf(report_file, "  %-50s", shipName)
	} else {
		fprintf(report_file, "  %-46s", shipName)
	}
	if printing_alien {
		fprintf(report_file, " ")
	} else {
		var capacity int
		if ship.Class.Is.Starbase {
			capacity = 10 * ship.Class.Tonnage
		} else if ship.Class.Is.Transport {
			capacity = (10 + (ship.Class.Tonnage / 2)) * ship.Class.Tonnage
		} else {
			capacity = ship.Class.Tonnage
		}
		fprintf(report_file, "%4d  ", capacity)
		if ship.Status.UnderConstruction {
			fprintf(report_file, "Left to pay = %d\n", ship.RemainingCost)
			return
		}
	}
	if printing_alien {
		item, ok := ship.Inventory["FD"]
		if ship.Status.OnSurface || !(ok && item.Quantity == ship.Class.Tonnage) {
			fprintf(report_file, "SP %s", species.Name)
		} else {
			fprintf(report_file, "SP %d", ship.Species.Distorted())
		}
	} else {
		var items []*cluster.Item
		for _, item := range ship.Inventory {
			if item.Quantity > 0 {
				items = append(items, item)
			}
		}
		// sort the inventory list
		for i := 0; i < len(items); i++ {
			for j := i + 1; j < len(items); j++ {
				if items[j].Code < items[i].Code {
					items[i], items[j] = items[j], items[i]
				}
			}
		}
		for i, item := range items {
			if i > 0 {
				fprintf(report_file, ",")
			}
			fprintf(report_file, "%d %s", item.Quantity, item.Abbr)
		}
	}
	fprintf(report_file, "\n")
}
