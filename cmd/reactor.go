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
	"github.com/mdhender/fhcms/internal/adapters"
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/reactor"
	"github.com/mdhender/fhcms/internal/repos/cdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"
)

var globalReactor struct {
	host  string
	port  string
	debug struct {
		dumpRequests bool
	}
}

var reactorCmd = &cobra.Command{
	Use:   "reactor",
	Short: "Serve multi-page app",
	Long:  `Provide a multi-page application for the game.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootDir := viper.Get("files.path").(string)
		log.Printf("[reactor] rootDir %q\n", rootDir)
		templatesDir := viper.Get("templates").(string)
		log.Printf("[reactor] templatesDir %q\n", templatesDir)

		authSecret, ok := viper.Get("server.secret").(string)
		if !ok || len(authSecret) < 1 {
			log.Fatal("server.secret must be at least 1 character long")
		}
		authSecret = mkkey("fhapp", authSecret)
		log.Printf("[reactor] mkkey(secretKey) %q\n", authSecret)
		fSigner, err := jot.NewHS256Signer([]byte(authSecret))
		cobra.CheckErr(err)

		//ds, err := domain.New(domain.WithAccounts(filepath.Join(rootDir, "accounts.json")), domain.WithGames(filepath.Join(rootDir, "games.json")))
		//cobra.CheckErr(err)

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

		s, err := reactor.New(globalReactor.host, globalReactor.port, reactor.WithAuthStore(db), reactor.WithGamesStore(db), reactor.WithJotFactory(jot.NewFactory("raven", fSigner)), reactor.WithProfileStore(db), reactor.WithSiteStore(db), reactor.WithTemplates(templatesDir))
		cobra.CheckErr(err)

		log.Printf("[reactor] listening on %q\n", net.JoinHostPort(globalReactor.host, globalReactor.port))
		//if globalReactor.debug.dumpRequests {
		log.Fatal(http.ListenAndServe(net.JoinHostPort(globalReactor.host, globalReactor.port), adapters.Logger(adapters.DumpRequest(s))))
		//} else {
		//	log.Fatal(http.ListenAndServe(net.JoinHostPort(globalReactor.host, globalReactor.port), adapters.Logger(s)))
		//}
	},
}

func init() {
	rootCmd.AddCommand(reactorCmd)
	reactorCmd.Flags().StringVar(&globalReactor.host, "host", "", "interface to run server on")
	_ = viper.BindPFlag("host", reactorCmd.Flags().Lookup("host"))
	reactorCmd.Flags().StringVarP(&globalReactor.port, "port", "p", "8080", "port to run server on")
	_ = viper.BindPFlag("port", reactorCmd.Flags().Lookup("port"))
}
