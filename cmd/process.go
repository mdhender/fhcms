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

package cmd

import (
	"github.com/mdhender/fhcms/internal/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(processCmd)
}

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process the current turn",
	Long:  `Load orders and process the current turn.`,
	Run: func(cmd *cobra.Command, args []string) {
		ds, err := loader(viper.GetString("files.path"), viper.GetBool("files.big_endian"))
		cobra.CheckErr(err)
		e := engine.FromBinary(ds)
		cobra.CheckErr(e.Run())
	},
}
