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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		//StatsMain(ds)

		game := ds.Stats()
		if game.TotalSpecies == 0 {
			return
		}

		// m attempts to compensate for rounding integer values down
		m := game.TotalSpecies / 2

		fmt.Printf("SP Species               Tech Levels        Total  Num Num  Num  Offen.  Defen.  Econ\n")
		fmt.Printf(" # Name             MI  MA  ML  GV  LS  BI  Prod.  Pls Shps Yrds  Power   Power  Units\n")
		fmt.Printf("----------------------------------------------------------------------------------------\n")
		for _, sp := range game.Stats {
			fmt.Printf("%2d %-15.15s%4d%4d%4d%4d%4d%4d%7.0f%4.0f%5.0f%5.0f%8.0f%8.0f%9.0f\n", sp.No, sp.Name, sp.MI, sp.MA, sp.ML, sp.GV, sp.LS, sp.BI, sp.Production, sp.PopulatedPlanets, sp.Ships, sp.Shipyards, sp.OffensivePower, sp.DefensivePower, sp.BankedEconUnits)
		}
		fmt.Println("")
		fmt.Printf("Average %-13s tech level = %3d (min = %3d, max = %3d)\n", game.MI.Name, (game.MI.Total+m)/game.TotalSpecies, game.MI.Min, game.MI.Max)
		fmt.Printf("Average %-13s tech level = %3d (min = %3d, max = %3d)\n", game.MA.Name, (game.MA.Total+m)/game.TotalSpecies, game.MA.Min, game.MA.Max)
		fmt.Printf("Average %-13s tech level = %3d (min = %3d, max = %3d)\n", game.ML.Name, (game.ML.Total+m)/game.TotalSpecies, game.ML.Min, game.ML.Max)
		fmt.Printf("Average %-13s tech level = %3d (min = %3d, max = %3d)\n", game.GV.Name, (game.GV.Total+m)/game.TotalSpecies, game.GV.Min, game.GV.Max)
		fmt.Printf("Average %-13s tech level = %3d (min = %3d, max = %3d)\n", game.LS.Name, (game.LS.Total+m)/game.TotalSpecies, game.LS.Min, game.LS.Max)
		fmt.Printf("Average %-13s tech level = %3d (min = %3d, max = %3d)\n", game.BI.Name, (game.BI.Total+m)/game.TotalSpecies, game.BI.Min, game.BI.Max)
		fmt.Println("")
		fmt.Printf("Average number of warships per species          %9.1f (min = %6.0f max = %6.0f)\n", game.Warships.Average, game.Warships.Min, game.Warships.Max)
		fmt.Printf("Average total warship tonnage per species       %9s tons\n", commas(int(game.WarshipTonnage.Average)))
		fmt.Printf("Average warship size                            %9s tons\n", commas(int(game.WarshipTonnage.Total/game.Warships.Total)))
		fmt.Println("")
		fmt.Printf("Average number of starbases per species         %9.1f (min = %6.0f max = %6.0f)\n", game.Starbases.Average, game.Starbases.Min, game.Starbases.Max)
		fmt.Printf("Average total starbase tonnage per species      %9s tons\n", commas(int(game.StarbaseTonnage.Average)))
		fmt.Printf("Average starbase size                           %9s tons\n", commas(int(game.StarbaseTonnage.Total/game.Starbases.Total)))
		fmt.Println("")
		fmt.Printf("Average number of transports per species        %9.1f (min = %6.0f max = %6.0f)\n", game.Transports.Average, game.Transports.Min, game.Transports.Max)
		fmt.Printf("Average total transport tonnage per species     %9s tons\n", commas(int(game.TransportTonnage.Average)))
		fmt.Printf("Average transport size                          %9s tons\n", commas(int(game.TransportTonnage.Total/game.Transports.Total)))
		fmt.Println("")
		fmt.Printf("Average number of shipyards per species         %9.1f (min = %6.0f max = %6.0f)\n", game.Shipyards.Average, game.Shipyards.Min, game.Shipyards.Max)
		fmt.Printf("Average number of populated planets per species %9.1f (min = %6.0f max = %6.0f)\n", game.PopulatedPlanets.Average, game.PopulatedPlanets.Min, game.PopulatedPlanets.Max)
		fmt.Printf("Average total production per species            %9.1f (min = %6.0f max = %6.0f)\n", game.Production.Average, game.Production.Min, game.Production.Max)
		fmt.Printf("Average banked economic units per species       %9.1f (min = %6.0f max = %6.0f)\n", game.BankedEconUnits.Average, game.BankedEconUnits.Min, game.BankedEconUnits.Max)
	},
}
