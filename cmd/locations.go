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
	rootCmd.AddCommand(locationsCmd)
}

var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "Create locations data file for current game.",
	Long: `This command creates the file locations.dat and updates the economic
efficiencies of all planets. (These functions are also performed by
the finish command.) This command should run before the strike phase
and whenever manual changes are made to the species data files that
change the location of ships, planets, or inventories. It should also
run if you run finish on fewer than all species and decide to keep
the resulting planets.dat file.`,
	Run: func(cmd *cobra.Command, args []string) {
		ds, err := loader(viper.GetString("files.path"), viper.GetBool("files.big_endian"))
		cobra.CheckErr(err)
		locations := ds.RunLocations(verboseFlag, testFlag)
		if verboseFlag {
			for _, loc := range locations {
				fmt.Printf("location: %q %s\n", loc.Species.Name, loc.Location.Id())
			}
		}
	},
}
