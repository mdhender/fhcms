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
	"github.com/mdhender/fhcms/internal/adapters"
	"github.com/mdhender/fhcms/internal/domain"
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/mpa"
	"github.com/mdhender/fhcms/internal/repos/accounts"
	"github.com/mdhender/fhcms/internal/repos/site"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
)

var host string
var port string
var debugDumpRequests bool

var mpaCmd = &cobra.Command{
	Use:   "mpa",
	Short: "Serve multi-page app",
	Long:  `Provide a multi-page application for the game.`,
	Run: func(cmd *cobra.Command, args []string) {
		siteFile := "D:\\GoLand\\fhcms\\testdata\\site.json"
		siteStore, err := site.New(siteFile)
		cobra.CheckErr(err)

		accountsFile := "D:\\GoLand\\fhcms\\testdata\\accounts.json"
		accts, err := accounts.Load(accountsFile)
		cobra.CheckErr(err)

		authSecret, ok := viper.Get("server.secret").(string)
		if !ok || len(authSecret) < 1 {
			log.Fatal("server.secret must be at least 1 character long")
		}
		authSecret = mkkey("fhapp", authSecret)
		log.Printf("mkkey(secretKey) %q\n", authSecret)
		fSigner, err := jot.NewHS256Signer([]byte(authSecret))
		cobra.CheckErr(err)

		ds, err := domain.New(domain.WithAccounts("D:\\GoLand\\fhcms\\testdata\\accounts.json"), domain.WithGames("D:\\GoLand\\fhcms\\testdata\\games.json"))
		cobra.CheckErr(err)

		s, err := mpa.New(host, port, mpa.WithTemplates("D:\\GoLand\\fhcms\\templates"), mpa.WithDomain(ds), mpa.WithAccounts(accts), mpa.WithJotFactory(jot.NewFactory("raven", fSigner)), mpa.WithSite(siteStore))
		cobra.CheckErr(err)

		fmt.Printf("listening on %q serving multi-page router\n", net.JoinHostPort(host, port))
		if debugDumpRequests {
			log.Fatal(http.ListenAndServe(net.JoinHostPort(host, port), adapters.Logger(adapters.DumpRequest(s))))
		} else {
			log.Fatal(http.ListenAndServe(net.JoinHostPort(host, port), adapters.Logger(s)))
		}
	},
}

func init() {
	rootCmd.AddCommand(mpaCmd)
	mpaCmd.Flags().StringVar(&host, "host", "", "interface to run server on")
	_ = viper.BindPFlag("host", mpaCmd.Flags().Lookup("host"))
	mpaCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to run server on")
	_ = viper.BindPFlag("port", mpaCmd.Flags().Lookup("port"))
}
