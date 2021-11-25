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
	"fmt"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"time"
)

func init() {
	rootCmd.AddCommand(defaultOrdersCmd)
}

var defaultOrdersCmd = &cobra.Command{
	Use:   "default-orders",
	Short: "Generate orders for the current turn.",
	Long:  `Generate a set of default orders for a species.`,
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

		started := time.Now().UTC()
		// create default orders for all species in the list
		for _, sp := range spList {
			fmt.Printf("generating orders for species %s, SP %s...\n", sp.Id, sp.Name)
			fmt.Printf(";; %s T%d %v\n", sp.Id, ds.Turn, started)
			fmt.Printf(";; default orders\n")

			// print out ship location and inventory
			if len(sp.Fleet.Ships) != 0 {
				fmt.Printf(";; Fleet Data\n")
				for _, ship := range sp.Fleet.Ships {
					fmt.Printf(";;   %-45s  %3d %3d %3d", ship.Display.Name, ship.Location.X, ship.Location.Y, ship.Location.Z)
					if ship.Location.Orbit == 0 {
						fmt.Printf("   ")
					} else {
						fmt.Printf(" #%d", ship.Location.Orbit)
					}
					if ship.Special != 0 {
						fmt.Printf("  special %d", ship.Special)
					}
					if ship.UnloadingPoint != 0 {
						fmt.Printf("  unload %d", ship.UnloadingPoint)
					}
					if ship.Status.UnderConstruction {
						fmt.Printf("  Under Construction\n")
					} else {
						fmt.Printf("  Age %2d\n", ship.Age)
						for _, item := range ship.Inventory {
							fmt.Printf(";;       %5d %-3s %s\n", item.Quantity, item.Abbr, item.Descr)
						}
					}
				}
			}
			fmt.Printf("\n\n")

			// PRE-DEPARTURE orders
			fmt.Printf("START PRE-DEPARTURE\n")
			fmt.Printf("    ; Place pre-departure orders here.\n\n")

			for _, colony := range sp.Colonies.ById {
				// generate auto-installs for colonies that were loaded via the DEVELOP command
				if colony.Mining.AutoIUs > 0 {
					fmt.Printf("    INSTALL %4d IU  PL %-32s", colony.Mining.AutoIUs, colony.Name.Display.Name)
					if item, ok := colony.Inventory["CU"]; ok && item.Quantity > 0 {
						fmt.Printf(" ;; consume %4d of %5d CU", colony.Mining.AutoIUs, item.Quantity)
						item.Quantity -= colony.Mining.AutoIUs
					}
					fmt.Printf("\n")
				}
				if colony.Manufacturing.AutoAUs > 0 {
					fmt.Printf("    INSTALL %4d AU  PL %-32s", colony.Manufacturing.AutoAUs, colony.Name.Display.Name)
					if item, ok := colony.Inventory["CU"]; ok && item.Quantity > 0 {
						fmt.Printf(" ;; consume %4d of %5d CU", colony.Manufacturing.AutoAUs, item.Quantity)
						item.Quantity -= colony.Manufacturing.AutoAUs
					}
					fmt.Printf("\n")
				}

				// generate auto UNLOAD orders for transports at this colony
				for _, ship := range sp.Fleet.Ships {
					if !ship.Class.Is.Transport {
						continue
					} else if ship.Location == nil || !(ship.Location.X == colony.Planet.Location.X && ship.Location.Y == colony.Planet.Location.Y && ship.Location.Z == colony.Planet.Location.Z && ship.Location.Orbit == colony.Planet.Location.Orbit) {
						continue
					} else if ship.Status != nil && (ship.Status.JumpedInCombat || ship.Status.ForcedJump) {
						continue
					}
					item, ok := colony.Inventory["CU"]
					if !ok || item.Quantity < 1 {
						continue
					}

					n := 0
					// colonies will never be started automatically unless ship was loaded via a DEVELOP order
					if ship.LoadingPoint != 0 {
						// is transport at specified unloading point?
						n = ship.UnloadingPoint
						if n == colony.Name.Index || (n == 9999 && colony.Name.Index == 0) {
							goto unloadShip
						}
					}
					if !colony.Is.Populated {
						continue
					}
					if colony.Mining.Base+colony.Manufacturing.Base >= 2000 {
						continue
					}
					if colony.Planet.Location.X == sp.HomeWorld.Planet.Location.X && colony.Planet.Location.Y == sp.HomeWorld.Planet.Location.Y && colony.Planet.Location.Z == sp.HomeWorld.Planet.Location.Z {
						// don't auto unload in the home sector
						continue
					}

				unloadShip:

					n = ship.LoadingPoint
					if n == 9999 { // home planet
						n = 0
					}
					if n == colony.Name.Index {
						// ship was just loaded here
						continue
					}

					fmt.Printf("    UNLOAD %-45s ;;", ship.Display.Name)
					if shipIUs, ok := ship.Inventory["IU"]; ok && shipIUs.Quantity > 0 {
						fmt.Printf(" %4d IU", shipIUs.Quantity)
					}
					if shipAUs, ok := ship.Inventory["AU"]; ok && shipAUs.Quantity > 0 {
						fmt.Printf(" %4d AU", shipAUs.Quantity)
					}
					if shipCUs, ok := ship.Inventory["CU"]; ok && shipCUs.Quantity > 0 {
						fmt.Printf(" %4d CU", shipCUs.Quantity)
					}
					fmt.Printf("\n")
					item.Quantity = 0 // set CU quantity to zero

					ship.Special = ship.LoadingPoint
					n = colony.Name.Index
					if n == 0 { // home planet
						n = 9999
					}
					ship.UnloadingPoint = n
				}

				if colony.Is.HomePlanet { // never auto install on the home world
					continue
				}
				if item, ok := colony.Inventory["CU"]; !ok || item.Quantity < 1 {
					continue
				}
				if item, ok := colony.Inventory["IU"]; ok && item.Quantity > 0 {
					fmt.Printf("    INSTALL    0 IU  PL %s\n", colony.Name.Display.Name)
				}
				if item, ok := colony.Inventory["AU"]; ok && item.Quantity > 0 {
					fmt.Printf("    INSTALL    0 AU  PL %s\n", colony.Name.Display.Name)
				}
			}
			fmt.Printf("END\n\n")

			// generate jump orders for ships used to develop and scouts
			fmt.Printf("START JUMPS\n")
			fmt.Printf("    ; Place jump orders here.\n\n")
			// initialize to make sure ships are not given more than one JUMP order
			for _, ship := range sp.Fleet.Ships {
				ship.JustJumped = false
			}
			// generate auto-jumps for ships that were loaded via the DEVELOP command or which were UNLOADed because of the AUTO command
			for _, ship := range sp.Fleet.Ships {
				if ship.Status.JumpedInCombat {
					continue
				} else if ship.Status.ForcedJump {
					continue
				} else if ship.Location == nil || ship.Location.Orbit == 99 {
					continue
				} else if ship.JustJumped {
					// how can this be true with the loop just above setting it to false?
					continue
				}

				j := ship.Special
				if j != 0 {
					if j == 9999 { // home planet
						j = 0
					}
					tempNampla := sp.NamedPlanets.Base[j]
					_, mishapChance := cluster.MishapChance(sp, ship, tempNampla.Planet.Location)
					fmt.Printf("    JUMP %s, PL %s  ; age %d  mishap chance = %s  (special)\n", ship.Display.Name, tempNampla.Display.Name, ship.Age, mishapChance)
					ship.JustJumped = true
					continue
				}

				n := ship.UnloadingPoint
				if n != 0 {
					if n == 9999 { // home planet
						n = 0
					}
					tempNampla := sp.NamedPlanets.Base[n]
					if ship.Location.X == tempNampla.Planet.Location.X && ship.Location.Y == tempNampla.Planet.Location.Y && ship.Location.Z == tempNampla.Planet.Location.Z {
						continue
					}
					_, mishapChance := cluster.MishapChance(sp, ship, tempNampla.Planet.Location)
					fmt.Printf("    JUMP %s, PL %s  ; age %d  mishap chance = %s  (unloadingPoint)\n", ship.Display.Name, tempNampla.Display.Name, ship.Age, mishapChance)
					ship.JustJumped = true
					continue
				}
			}
			// generate JUMP orders for all TR1s
			for _, ship := range sp.Fleet.Ships {
				if ship.Location.Orbit == 99 {
					continue
				} else if ship.Status.UnderConstruction {
					continue
				} else if ship.Status.JumpedInCombat {
					continue
				} else if ship.Status.ForcedJump {
					continue
				} else if ship.JustJumped {
					continue
				}
				// todo: calculate delta x, y, or z that moves us closer to that system for sublight ships
				if ship.Class.Is.Transport && ship.Class.Tonnage == 1 && !ship.Class.Is.SubLight {
					closestSystem := ds.ClosestUnvisitedSystem(sp, ship.Location)
					_, mishapChance := cluster.MishapChance(sp, ship, closestSystem)
					fmt.Printf("    JUMP %s, ", ship.Display.Name)
					if closestSystem == nil {
						fmt.Printf("? ? ?  ; scout - no unvisited systems")
					} else {
						fmt.Printf("%3d %3d %3d  ; scout", closestSystem.X, closestSystem.Y, closestSystem.Z)
					}
					fmt.Printf("\n")
					fmt.Printf("            ; Age %d, now at %d %d %d, mishap chance = %s", ship.Age, ship.Location.X, ship.Location.Y, ship.Location.Z, mishapChance)
					fmt.Printf("\n")
					ship.Destination = closestSystem
					ship.JustJumped = true
				}
			}
			fmt.Printf("END\n\n")

			// generate a PRODUCTION order for each planet that can produce
			fmt.Printf("START PRODUCTION\n")
			// run through the colonies in reverse order
			for i := len(sp.NamedPlanets.Base) - 1; i >= 0; i-- {
				nampla := sp.NamedPlanets.Base[i]
				if nampla.Colony == nil || nampla.Planet.Orbit == 99 {
					continue
				} else if nampla.Colony.Mining.Base == 0 && (!nampla.Colony.Is.ResortColony) {
					continue
				} else if nampla.Colony.Manufacturing.Base == 0 && (!nampla.Colony.Is.MiningColony) {
					continue
				}
				fmt.Printf("    PRODUCTION PL %-32s ; %3d %3d %3d #%d\n", nampla.Display.Name, nampla.Planet.Location.X, nampla.Planet.Location.Y, nampla.Planet.Location.Z, nampla.Planet.Location.Orbit)
				if nampla.Colony.Is.MiningColony {
					fmt.Printf("      ; The above PRODUCTION order is required for this mining colony\n")
					fmt.Printf("      ;  even if no other production orders are given for it.\n")
				} else if nampla.Colony.Is.ResortColony {
					fmt.Printf("      ; The above PRODUCTION order is required for this resort colony\n")
					fmt.Printf("      ;  even though no other production orders can be given for it.\n")
				} else if nampla.Planet != sp.HomeWorld.Planet {
					fmt.Printf("      ; Place production orders here for colony.\n")
				} else {
					fmt.Printf("      ; Place production orders here for homeworld.\n")
				}
				for _, item := range nampla.Colony.Inventory {
					if item.Quantity > 0 {
						fmt.Printf("      ; %-3s %-30s %9d\n", item.Abbr, item.Descr, item.Quantity)
					}
				}
				// build IUs and AUs for incoming ships with CUs
				if nampla.Colony.Mining.Needed > 0 {
					fmt.Printf("      BUILD %5d IU\n", nampla.Colony.Mining.Needed)
				}
				if nampla.Colony.Manufacturing.Needed > 0 {
					fmt.Printf("      BUILD %5d AU\n", nampla.Colony.Manufacturing.Needed)
				}
				if nampla.Colony.Is.MiningColony || nampla.Colony.Is.ResortColony {
					continue
				}
				// see if there are any RMs to recycle
				if n := nampla.Colony.Special / 5; n > 0 {
					fmt.Printf("      RECYCLE %5d RM  ; special != 0\n", 5*n)
				} else if item, ok := nampla.Colony.Inventory["RM"]; ok && item.Quantity > 5 {
					fmt.Printf("      RECYCLE %5d RM  ; of %d total\n", (item.Quantity/5)*5, item.Quantity)
				}
				// generate DEVELOP commands for ships arriving here because of	AUTO command
				for _, ship := range sp.Fleet.Ships {
					if ship.Location == nil || ship.Location.Orbit == 99 {
						continue
					}
					// k wants to be a relative planet index
					k := ship.Special
					if k == 0 {
						continue
					}
					if k == 9999 { // home planet
						k = 0
					}
					if nampla.Index != k { // nampla != nampla_base + k
						continue
					}
					k = ship.UnloadingPoint
					if k == 9999 { // home planet?
						k = 0
					}
					temp_nampla := sp.NamedPlanets.Base[k]
					fmt.Printf("      DEVELOP PL %s, %s  ; ship arriving because of auto\n", temp_nampla.Display.Name, ship.Display.Name)
				}
				// give orders to continue construction of unfinished ships and starbases
				for _, ship := range sp.Fleet.Ships {
					if !(ship.Status.UnderConstruction || ship.Class.Is.Starbase) {
						continue
					} else if ship.Location == nil || ship.Location.Orbit == 99 {
						continue
					} else if !(ship.Location.X != nampla.Planet.Location.X && ship.Location.Y != nampla.Planet.Location.Y && ship.Location.Z != nampla.Planet.Location.Z && ship.Location.Orbit != nampla.Planet.Location.Orbit) {
						continue
					}
					if ship.Status.UnderConstruction {
						fmt.Printf("      CONTINUE %s, %d\t; Left to pay = %d\n", ship.Display.Name, ship.RemainingCost, ship.RemainingCost)
					} else if j := (sp.MA.Level / 2) - ship.Class.Tonnage; j > 0 {
						// ship is a starbase that is not already as large as the tech level allows
						fmt.Printf("      CONTINUE BAS %s, %d\t; Current tonnage = %s\n", ship.Name, 100*j, commas(10_000*ship.Class.Tonnage))
					}
				}
				// generate DEVELOP command if this is a colony with an economic base less than 200
				n := nampla.Colony.Mining.Base + nampla.Colony.Mining.Needed +
					nampla.Colony.Manufacturing.Base + nampla.Colony.Manufacturing.Needed
				if nampla.Colony.Is.Colony && n < 2000 && nampla.Colony.Population > 0 {
					nn := nampla.Colony.Population
					if nn > 2000-n {
						nn = 2000 - n
					}
					fmt.Printf("      DEVELOP %d  ; colony econ base %d\n", 2*nn, n/10)
					nampla.Colony.Mining.Needed += nn
				}
				// for home planets and any colonies that have an economic base of at least 200,
				// check if there are other colonized planets in the same sector that are not
				// self-sufficient. if so, DEVELOP them.
				if n >= 2000 || nampla.Colony.Is.HomePlanet {
					// loop skips index zero since it is the home planet.
					// it makes sense because we will never target the home planet for development.
					for i := 1; i < len(sp.NamedPlanets.Base); i++ {
						if i == nampla.Index { // skip self
							continue
						}
						temp_nampla := sp.NamedPlanets.Base[i]

						if temp_nampla.Planet.Location.Orbit == 99 {
							continue
						} else if !(temp_nampla.Planet.Location.X != nampla.Planet.Location.X && temp_nampla.Planet.Location.Y != nampla.Planet.Location.Y && temp_nampla.Planet.Location.Z != nampla.Planet.Location.Z) {
							continue
						}
						n = temp_nampla.Colony.Mining.Base + temp_nampla.Colony.Mining.Needed +
							temp_nampla.Colony.Manufacturing.Base + temp_nampla.Colony.Manufacturing.Needed
						if n == 0 {
							continue
						}
						numberNeeded := 0
						if item, ok := temp_nampla.Colony.Inventory["IU"]; ok && item.Quantity > 0 {
							numberNeeded += item.Quantity
						}
						if item, ok := temp_nampla.Colony.Inventory["IU"]; ok && item.Quantity > 0 {
							numberNeeded += item.Quantity
						}
						if item, ok := temp_nampla.Colony.Inventory["CU"]; ok && numberNeeded > item.Quantity {
							numberNeeded = item.Quantity
						}
						n += numberNeeded
						if n >= 2000 {
							continue
						}
						numberNeeded = 2000 - n
						if numberNeeded > nampla.Colony.Population {
							numberNeeded = nampla.Colony.Population
						}
						fmt.Printf("      DEVELOP %d  PL %s ; develop siblings\n", 2*numberNeeded, temp_nampla.Display.Name)
						temp_nampla.Colony.Manufacturing.Needed += numberNeeded
					}
				}
			}
			fmt.Printf("END\n\n")

			fmt.Printf("START POST-ARRIVAL\n")
			fmt.Printf("    ; Place post-arrival orders here.\n")
			fmt.Printf("    AUTO\n") // generate an AUTO command
			// generate SCAN orders for all TR1s in sectors that current species does not inhabit
			for _, ship := range sp.Fleet.Ships {
				if ship.Location == nil || ship.Location.Orbit == 99 {
					continue
				} else if ship.Destination == nil || ship.Destination.X == -1 { // not jumping anywhere
					continue
				} else if ship.Status.UnderConstruction {
					continue
				} else if !(ship.Class.Is.Transport && ship.Class.Tonnage == 1) {
					continue
				} else if ship.Class.Is.SubLight {
					continue
				}
				found := false
				for _, nampla := range sp.NamedPlanets.Base[1:] { // start at 1 to skip home sector
					if nampla.Planet.Location == nil || nampla.Planet.Location.Orbit == 99 {
						continue
					} else if nampla.Planet.Location.X != ship.Destination.X {
						continue
					} else if nampla.Planet.Location.Y != ship.Destination.Y {
						continue
					} else if nampla.Planet.Location.Z != ship.Destination.Z {
						continue
					} else if nampla.Colony != nil && nampla.Colony.Is.Populated {
						found = true
					}
				}
				if !found {
					fmt.Printf("    SCAN %s\n", ship.Display.Name)
				}
			}
			fmt.Printf("END\n\n")
		}
	},
}

func init() {
	defaultOrdersCmd.Flags().Int("species-no", 0, "species number to generate orders for")
}
