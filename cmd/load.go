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
	"context"
	"encoding/json"
	"github.com/mdhender/fhcms/internal/repos/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
)

var galaxySetupFile string
var galaxyForceInit bool
var galaxyForcePurge bool

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load galaxy data into database",
	Long:  `Load data.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootDir := viper.Get("files.path").(string)
		log.Printf("[mpa] rootDir %q\n", rootDir)

		dbConfig := &cdb.DBConfig{}
		data, err := ioutil.ReadFile(filepath.Join("D:\\GoLand\\fhcms\\testdata", "database.json"))
		cobra.CheckErr(err)
		err = json.Unmarshal(data, dbConfig)
		cobra.CheckErr(err)
		db, err := cdb.New(context.Background(), dbConfig)
		cobra.CheckErr(err)
		defer func(db *cdb.DB) {
			log.Printf("[cdb] closing connection\n")
			db.Close()
		}(db)

		cobra.CheckErr(db.LoadGalaxy(galaxySetupFile, galaxyForceInit, galaxyForcePurge))
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().StringVar(&galaxySetupFile, "setup", "", "json file containing setup data")
	loadCmd.Flags().BoolVar(&galaxyForceInit, "init", false, "load only galaxy configuration tables")
	loadCmd.Flags().BoolVar(&galaxyForcePurge, "purge", false, "purge any existing galaxy with the same name")
}
