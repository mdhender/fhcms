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
	"fmt"
	"github.com/mdhender/fhcms/internal/adapters"
	"github.com/mdhender/fhcms/internal/app"
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/repos/accounts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.Flags().StringVar(&globalApp.data, "data", "D:\\GoLand\\fhcms\\testdata\\raven", "path to data files")
	_ = viper.BindPFlag("data", appCmd.Flags().Lookup("data"))
	appCmd.Flags().StringVar(&globalApp.templates, "templates", "D:\\GoLand\\fhcms\\templates", "path to template files")
	_ = viper.BindPFlag("data", appCmd.Flags().Lookup("data"))
	appCmd.Flags().StringVar(&globalApp.host, "host", "", "interface to run server on")
	_ = viper.BindPFlag("host", appCmd.Flags().Lookup("host"))
	appCmd.Flags().StringVarP(&globalApp.port, "port", "p", "8080", "port to run server on")
	_ = viper.BindPFlag("port", appCmd.Flags().Lookup("port"))
}

var globalApp struct {
	data      string // path to data files
	templates string // path to template files
	host      string
	port      string
	debug     struct {
		dumpRequests bool
	}
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Serve player app",
	Long:  `Allow players to view game data.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("[app] data %q\n", globalApp.data)
		log.Printf("[app] templates %q\n", globalApp.templates)

		authSecret, ok := viper.Get("server.secret").(string)
		if !ok || len(authSecret) == 0 {
			cobra.CheckErr(fmt.Errorf("server.secret is required and must not be empty"))
		}
		authSalt, ok := viper.Get("server.salt").(string)
		if !ok || len(authSalt) == 0 {
			cobra.CheckErr(fmt.Errorf("server.salt is required and must not be empty"))
		}

		acctRepo, err := accounts.Load("D:\\GoLand\\fhcms\\testdata\\accounts.json")
		if err != nil {
			log.Fatal(err)
		}

		fSigner, err := jot.NewHS256Signer([]byte(mkkey(authSalt, authSecret)))
		if err != nil {
			log.Fatal(err)
		}
		jf := jot.NewFactory("raven", fSigner)

		var options []app.Option
		options = append(options, app.WithMaxBodyLength(32*1024))
		options = append(options, app.WithHost(globalApp.host))
		options = append(options, app.WithPort(globalApp.port))
		options = append(options, app.WithAccountStore(acctRepo))
		options = append(options, app.WithJotFactory(jf))
		options = append(options, app.WithData(globalApp.data))
		options = append(options, app.WithTemplates(globalApp.templates))

		s, err := app.New(options...)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[app] listening on %q\n", s.Addr)
		if err = http.ListenAndServe(s.Addr, adapters.Logger(s.WithUser(s.CORS(s)))); err != nil {
			log.Fatal(err)
		}
	},
}
