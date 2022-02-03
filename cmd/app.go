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
	"github.com/mdhender/fhcms/internal/app"
	"github.com/mdhender/fhcms/internal/repos/accounts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
)

func init() {
	rootCmd.AddCommand(appCmd)
}

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Serve player app",
	Long:  `Allow players to view game data.`,
	Run: func(cmd *cobra.Command, args []string) {
		var options []app.Option
		options = append(options, app.WithMaxBodyLength(32*1024))

		host, ok := viper.Get("server.host").(string)
		if !ok {
			host = ""
		}
		options = append(options, app.WithHost(host))

		port, ok := viper.Get("server.port").(string)
		if !ok {
			port = "8080"
		}
		options = append(options, app.WithPort(port))

		acctRepo, err := accounts.Load("D:\\GoLand\\fhcms\\testdata\\accounts.json")
		if err != nil {
			log.Fatal(err)
		}
		options = append(options, app.WithAccountStore(acctRepo))

		s, err := app.New(options...)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("[app] listening on %q\n", net.JoinHostPort(globalReactor.host, globalReactor.port))
		if err = http.ListenAndServe(s.Addr, s); err != nil {
			log.Fatal(err)
		}
	},
}
