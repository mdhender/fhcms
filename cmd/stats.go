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
	"log"
)

func init() {
	rootCmd.AddCommand(statsCmd)
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics for the current game",
	Long:  `Load game data and print statistics for species.`,
	Run: func(cmd *cobra.Command, args []string) {
		ds, err := loader(viper.GetString("files.path"), viper.GetBool("big_endian"))
		cobra.CheckErr(err)
		StatsMain(ds)
	},
}

func StatsMain(ds *cluster.Store) {
	// sort by species number to avoid the random output from a map
	var xlat []*cluster.Species
	for _, sp := range ds.Species {
		xlat = append(xlat, sp)
	}
	for i := 0; i < len(xlat); i++ {
		for j := i + 1; j < len(xlat); j++ {
			if xlat[i].No > xlat[j].No {
				xlat[i], xlat[j] = xlat[j], xlat[i]
			}
		}
	}

	var num_ships, avg_tech_level int
	var fleet_percent_cost int
	var ntr, nba, nwa int
	var all_tech_level, min_tech_level, max_tech_level [6]int

	var total_tonnage, total_offensive_power int

	/* Initialize data. */
	n_species := 0
	all_production := 0
	min_production := 1000000000
	max_production := 0
	all_warship_tons := 0
	all_starbase_tons := 0
	all_transport_tons := 0
	n_warships := 0
	min_warships := 32000
	max_warships := 0
	n_starbases := 0
	min_starbases := 32000
	max_starbases := 0
	n_transports := 0
	min_transports := 32000
	max_transports := 0
	n_pop_pl := 0
	min_pop_pl := 32000
	max_pop_pl := 0
	n_yards := 0
	min_yards := 32000
	max_yards := 0
	totalBankedEconUnits, minBankedEconUnits, maxBankedEconUnits := 0, 0, 0

	var tech_name [6]string
	for i := 0; i < 6; i++ {
		all_tech_level[i] = 0
		min_tech_level[i] = 32000
		max_tech_level[i] = 0
		tech_name[i] = cluster.TechName(i)
	}

	fmt.Printf("SP Species               Tech Levels        Total  Num Num  Num  Offen.  Defen.  Econ\n")
	fmt.Printf(" # Name             MI  MA  ML  GV  LS  BI  Prod.  Pls Shps Yrds  Power   Power  Units\n")
	fmt.Printf("----------------------------------------------------------------------------------------\n")

	/* Main loop. For each species, take appropriate action. */
	for _, species := range xlat {
		n_species++

		/* Get fleet maintenance cost. */
		fleet_percent_cost = species.Fleet.MaintenancePct
		if fleet_percent_cost > 10000 {
			fleet_percent_cost = 10000
		}

		/* Print species data. */
		fmt.Printf("%2d", species.No)
		fmt.Printf(" %-15.15s", species.Name)

		for i := 0; i < 6; i++ {
			var level int
			switch i {
			case 0:
				level = species.MI.Level
			case 1:
				level = species.MA.Level
			case 2:
				level = species.ML.Level
			case 3:
				level = species.GV.Level
			case 4:
				level = species.LS.Level
			case 5:
				level = species.BI.Level
			}
			fmt.Printf("%4d", level)
			all_tech_level[i] += level
			if level < min_tech_level[i] {
				min_tech_level[i] = level
			}
			if level > max_tech_level[i] {
				max_tech_level[i] = level
			}
		}

		/* Get stats for namplas. */
		total_production := 0
		total_defensive_power := 0
		num_yards := 0
		num_pop_planets := 0
		for _, nampla := range species.NamedPlanets.ById {
			if nampla.Planet.Orbit == 99 {
				continue
			}
			planet, colony := nampla.Planet, nampla.Colony
			if planet == nil || colony == nil {
				log.Printf("species %q %v %v\n", species.Name, planet == nil, colony == nil)
				continue
			}

			num_yards += colony.Shipyards
			n_yards += colony.Shipyards

			var production_penalty int
			if ls_needed := lifeSupportNeeded(species, planet); ls_needed == 0 {
				production_penalty = 0
			} else if species.LS.Level > 0 {
				production_penalty = (100 * ls_needed) / species.LS.Level
			} else {
				production_penalty = 100
			}

			raw_material_units := (10 * species.MI.Level * colony.Mining.Base) / planet.MiningDifficulty
			raw_material_units -= (production_penalty * raw_material_units) / 100
			raw_material_units = ((planet.EconEfficiency * raw_material_units) + 50) / 100

			production_capacity := (species.MA.Level * colony.Manufacturing.Base) / 10
			production_capacity -= (production_penalty * production_capacity) / 100
			production_capacity = ((planet.EconEfficiency * production_capacity) + 50) / 100

			var n1 int
			if colony.Is.MiningColony {
				n1 = (2 * raw_material_units) / 3
			} else if colony.Is.ResortColony {
				n1 = (2 * production_capacity) / 3
			} else if production_capacity > raw_material_units {
				n1 = raw_material_units
			} else {
				n1 = production_capacity
			}

			n2 := ((fleet_percent_cost * n1) + 5000) / 10000
			n3 := n1 - n2
			total_production += n3

			if item, ok := colony.Inventory["PD"]; ok {
				tons := item.Quantity / 200
				if tons < 1 && item.Quantity > 0 {
					tons = 1
				}
				total_defensive_power += power(tons)
			}

			if colony.Is.Populated {
				n_pop_pl++
				num_pop_planets++
			}
		}

		fmt.Printf("%7d%4d", total_production, num_pop_planets)

		if total_production < min_production {
			min_production = total_production
		}
		if total_production > max_production {
			max_production = total_production
		}

		if num_pop_planets < min_pop_pl {
			min_pop_pl = num_pop_planets
		}
		if num_pop_planets > max_pop_pl {
			max_pop_pl = num_pop_planets
		}

		if num_yards < min_yards {
			min_yards = num_yards
		}
		if num_yards > max_yards {
			max_yards = num_yards
		}

		all_production += total_production

		/* Get stats for ships. */
		num_ships, ntr, nba, nwa, total_tonnage, total_offensive_power = 0, 0, 0, 0, 0, 0
		for _, ship := range species.Fleet.Ships {
			if ship.Location == nil || ship.Location.Orbit == 99 {
				continue
			} else if ship.Status.UnderConstruction {
				continue
			}

			num_ships++
			total_tonnage += ship.Class.Tonnage

			if ship.Class.Is.Starbase {
				total_defensive_power += power(ship.Class.Tonnage)
				all_starbase_tons += ship.Class.Tonnage
				n_starbases++
				nba++
			} else if ship.Class.Is.Transport {
				all_transport_tons += ship.Class.Tonnage
				n_transports++
				ntr++
			} else if ship.Class.Is.Warship {
				if ship.Class.Is.SubLight {
					total_defensive_power += power(ship.Class.Tonnage)
				} else {
					total_offensive_power += power(ship.Class.Tonnage)
				}
				all_warship_tons += ship.Class.Tonnage
				n_warships++
				nwa++
			}
		}

		if nwa < min_warships {
			min_warships = nwa
		}
		if nwa > max_warships {
			max_warships = nwa
		}

		if nba < min_starbases {
			min_starbases = nba
		}
		if nba > max_starbases {
			max_starbases = nba
		}

		if ntr < min_transports {
			min_transports = ntr
		}
		if ntr > max_transports {
			max_transports = ntr
		}

		if species.ML.Level == 0 {
			total_defensive_power = 0
			total_offensive_power = 0
		} else {
			total_offensive_power += (species.ML.Level * total_offensive_power) / 50
			total_offensive_power /= 10
			total_defensive_power += (species.ML.Level * total_defensive_power) / 50
			total_defensive_power /= 10
		}

		fmt.Printf("%5d", num_ships)
		fmt.Printf("%5d", num_yards)
		fmt.Printf("%8d%8d", total_offensive_power, total_defensive_power)

		totalBankedEconUnits += species.EconUnits
		if species.No == 1 {
			minBankedEconUnits = species.EconUnits
			maxBankedEconUnits = species.EconUnits
		} else {
			if minBankedEconUnits > species.EconUnits {
				minBankedEconUnits = species.EconUnits
			}
			if maxBankedEconUnits < species.EconUnits {
				maxBankedEconUnits = species.EconUnits
			}
		}
		fmt.Printf("%9d\n", species.EconUnits)
	}

	m := n_species / 2
	fmt.Printf("\n")
	for i := 0; i < 6; i++ {
		avg_tech_level = (all_tech_level[i] + m) / n_species
		fmt.Printf("Average %s tech level = %d (min = %d, max = %d)\n", tech_name[i], avg_tech_level, min_tech_level[i], max_tech_level[i])
	}

	i := ((10 * n_warships) + m) / n_species
	fmt.Printf("\nAverage number of warships per species = %d.%d (min = %d, max = %d)\n", i/10, i%10, min_warships, max_warships)

	if n_warships == 0 {
		n_warships = 1
	}

	avg_warship_tons := (10000 * all_warship_tons) / n_warships
	avg_warship_tons = 1000 * ((avg_warship_tons + 500) / 1000)
	fmt.Printf("Average warship size = %s tons\n", commas(avg_warship_tons))

	avg_warship_tons = (10000 * all_warship_tons) / n_species
	avg_warship_tons = 1000 * ((avg_warship_tons + 500) / 1000)
	fmt.Printf("Average total warship tonnage per species = %s tons\n", commas(avg_warship_tons))

	i = ((10 * n_starbases) + m) / n_species
	fmt.Printf("\nAverage number of starbases per species = %d.%d (min = %d, max = %d)\n", i/10, i%10, min_starbases, max_starbases)

	if n_starbases == 0 {
		n_starbases = 1
	}

	avg_starbase_tons := (10000 * all_starbase_tons) / n_starbases
	avg_starbase_tons = 1000 * ((avg_starbase_tons + 500) / 1000)
	fmt.Printf("Average starbase size = %s tons\n", commas(avg_starbase_tons))

	avg_starbase_tons = (10000 * all_starbase_tons) / n_species
	avg_starbase_tons = 1000 * ((avg_starbase_tons + 500) / 1000)
	fmt.Printf("Average total starbase tonnage per species = %s tons\n", commas(avg_starbase_tons))

	i = ((10 * n_transports) + m) / n_species
	fmt.Printf("\nAverage number of transports per species = %d.%d (min = %d, max = %d)\n", i/10, i%10, min_transports, max_transports)

	if n_transports == 0 {
		n_transports = 1
	}

	avg_transport_tons := (10000 * all_transport_tons) / n_transports
	avg_transport_tons = 1000 * ((avg_transport_tons + 500) / 1000)
	fmt.Printf("Average transport size = %s tons\n", commas(avg_transport_tons))

	avg_transport_tons = (10000 * all_transport_tons) / n_species
	avg_transport_tons = 1000 * ((avg_transport_tons + 500) / 1000)
	fmt.Printf("Average total transport tonnage per species = %s tons\n", commas(avg_transport_tons))

	avg_yards := ((10 * n_yards) + m) / n_species
	fmt.Printf("\nAverage number of shipyards per species = %d.%d (min = %d, max = %d)\n", avg_yards/10, avg_yards%10, min_yards, max_yards)

	avg_pop_pl := ((10 * n_pop_pl) + m) / n_species
	fmt.Printf("\nAverage number of populated planets per species = %d.%d (min = %d, max = %d)\n", avg_pop_pl/10, avg_pop_pl%10, min_pop_pl, max_pop_pl)

	avg_production := (all_production + m) / n_species
	fmt.Printf("Average total production per species = %d (min = %d, max = %d)\n", avg_production, min_production, max_production)

	avgBankedEconUnits := (totalBankedEconUnits + m) / n_species
	fmt.Printf("\nAverage banked economic units per species = %d (min = %d, max = %d)\n", avgBankedEconUnits, minBankedEconUnits, maxBankedEconUnits)
}
