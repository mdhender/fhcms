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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/fhcms/internal/admin"
	"github.com/mdhender/fhcms/internal/api"
	"github.com/mdhender/fhcms/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net"
	"net/http"
)

var host string
var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve data for the game",
	Long:  `Provide a REST-ish API for game data.`,
	Run: func(cmd *cobra.Command, args []string) {
		// For debugging/example purposes, we generate and printsample jwt tokens here
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"species": 18})
		fmt.Printf("jwt: player 18 %s\n\n", tokenString)
		_, adminTokenString, _ := tokenAuth.Encode(map[string]interface{}{"species": 18, "admin": true})
		fmt.Printf("jst: admin %s\n\n", adminTokenString)

		r, routerName := chi.NewRouter(), "chi"
		r.Use(middleware.Logger)

		// public routes
		r.Mount("/", app.Router()) // mount the static sub-router

		// protected routes
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth)) // extract, verify, validate JWT
			r.Use(jwtauth.Authenticator)       // handle valid and invalid JWT

			r.Mount("/admin", admin.Router()) // mount the admin sub-router
			r.Mount("/api", api.Router())     // mount the api sub-router
		})

		fmt.Printf("listening on %q using %s router\n", net.JoinHostPort(host, port), routerName)
		log.Fatal(http.ListenAndServe(net.JoinHostPort(host, port), r))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&host, "host", "", "interface to run server on")
	_ = viper.BindPFlag("host", serveCmd.Flags().Lookup("host"))
	serveCmd.Flags().StringVarP(&port, "port", "p", "8080", "port to run server on")
	_ = viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}
