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
	"context"
	"encoding/json"
	"fmt"
	"github.com/mdhender/fhcms/internal/repos/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"path/filepath"
)

var dbCreateSchema bool

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "run database commands",
	Long:  `Create, backup, upgrade, and restore the database tables.`,
	Run: func(cmd *cobra.Command, args []string) {
		dataPath := viper.Get("files.path").(string)
		log.Printf("[db] dataPath %q\n", dataPath)

		var owner struct {
			Username string
			Password string
			Email    string
		}
		if dbCreateSchema {
			var err error
			var ok bool
			if owner.Username, ok = viper.Get("owner.username").(string); !ok || owner.Username == "" {
				err = fmt.Errorf("missing owner.username in configuration")
			} else if owner.Password, ok = viper.Get("owner.password").(string); !ok || owner.Password == "" {
				err = fmt.Errorf("missing owner.password in configuration")
			} else if owner.Email, ok = viper.Get("owner.email").(string); !ok || owner.Email == "" {
				err = fmt.Errorf("missing owner.email in configuration")
			}
			cobra.CheckErr(err)
		}

		dbConfig := &cdb.DBConfig{}
		log.Printf("[db] loading database configuration %q\n", filepath.Join(dataPath, "database.json"))
		data, err := ioutil.ReadFile(filepath.Join(dataPath, "database.json"))
		cobra.CheckErr(err)
		err = json.Unmarshal(data, dbConfig)
		cobra.CheckErr(err)

		log.Printf("[db] opening database connection\n")
		ctx := context.Background()
		db, err := cdb.New(ctx, dbConfig)
		cobra.CheckErr(err)
		defer func(db *cdb.DB) {
			log.Printf("[db] closing database connection\n")
			db.Close()
		}(db)

		if dbCreateSchema {
			cobra.CheckErr(db.CreateSchema(ctx, owner.Username, owner.Password, owner.Email))
		}
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.Flags().BoolVar(&dbCreateSchema, "create", false, "create database tables")
}
