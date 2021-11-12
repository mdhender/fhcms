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
	rootCmd.AddCommand(turnCmd)
}

var turnCmd = &cobra.Command{
	Use:   "turn",
	Short: "Print the turn number for the current game",
	Long:  `Load game data and print the current turn number.`,
	Run: func(cmd *cobra.Command, args []string) {
		ds, err := loader(viper.GetString("files.path"), viper.GetBool("files.big_endian"))
		cobra.CheckErr(err)
		fmt.Printf("%d\n", ds.Turn)
	},
}
