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
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	"github.com/mdhender/fhcms/internal/api"
	"github.com/mdhender/fhcms/internal/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

var host string
var port string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve data for the game",
	Long:  `Provide a REST-ish API for game data.`,
	Run: func(cmd *cobra.Command, args []string) {
		var accounts struct {
			Accounts map[string]struct {
				Email    string `json:"email"`
				Username string `json:"username"`
				Password string `json:"password"`
			}
		}
		if data, err := ioutil.ReadFile("D:\\GoLand\\fhcms\\testdata\\accounts.json"); err != nil {
			cobra.CheckErr(err)
		} else if err := json.Unmarshal(data, &accounts); err != nil {
			cobra.CheckErr(err)
		}

		authSecret, ok := viper.Get("server.secret").(string)
		if !ok || len(authSecret) < 12 {
			log.Fatal("server.secret must be at least 12 characters long")
		}
		authSecret = mkkey("fhapp", authSecret)
		log.Printf("mkkey(secretKey) %q\n", authSecret)

		// For debugging/example purposes, we generate and print sample jwt tokens here
		tokenAuth := jwtauth.New("HS256", []byte("secret"), nil)
		_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"species": 18})
		fmt.Printf("jwt: player 18 %s\n\n", tokenString)
		_, adminTokenString, _ := tokenAuth.Encode(map[string]interface{}{"species": 18, "admin": true})
		fmt.Printf("jst: admin %s\n\n", adminTokenString)

		r, routerName := chi.NewRouter(), "chi"
		r.Use(middleware.Logger)
		r.Use(middleware.Recoverer)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				"GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS",
			},
			AllowedHeaders: []string{
				"Accept", "Accept-Encoding", "Accept-Language", "Authorization", "Cache-Control", "Connection", "Content-Type", "DNT", "Host", "Origin", "Pragma", "Referer", "User-Agent",
			},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

		// auth/auth routes
		tokenTimeToLive := 14 * 24 * 60 * 60 * time.Second
		r.Post("/auth/jwt", func(w http.ResponseWriter, r *http.Request) {
			var input struct {
				Username string `json:"username,omitempty"`
				Password string `json:"password,omitempty"`
			}

			if contentType := r.Header.Get("Content-type"); contentType != "application/json" {
				http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
				return
			}
			log.Printf("app: %s %q\n", r.Method, r.URL.Path)
			r.Body = http.MaxBytesReader(w, r.Body, 1024) // enforce a maximum read of 1kb from the response body
			// create a json decoder that will accept only our specific fields
			dec := json.NewDecoder(r.Body)
			dec.DisallowUnknownFields()
			if err := dec.Decode(&input); err != nil {
				log.Printf("app: %s %q: decode %+v\n", r.Method, r.URL.Path, err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			// call decode again to confirm that the request contained only a single JSON object
			if err := dec.Decode(&struct{}{}); err != io.EOF {
				log.Printf("app: %s %q: eof check %+v\n", r.Method, r.URL.Path, err)
				http.Error(w, "request body must only contain a single JSON object", http.StatusBadRequest)
				return
			}
			log.Printf("app: json: %s %q: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)

			for id, acct := range accounts.Accounts {
				if acct.Username == input.Username && acct.Password == input.Password {
					log.Printf("app: json: %s %q: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
					var response struct {
						Token string `json:"token"`
					}
					iat := time.Now()
					_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
						"iat": iat.UTC().Unix(),
						"exp": iat.Add(tokenTimeToLive).UTC().Unix(),
						"sub": id,
					})
					if err != nil {
						log.Printf("app: %s %q: gentok %+v\n", r.Method, r.URL.Path, err)
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}
					response.Token = tokenString
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(response)
					return
				}
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		})

		// public routes
		r.Mount("/", app.Router()) // mount the static sub-router

		// protected routes
		r.Route("/api", func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth)) // extract, verify, validate JWT
			r.Use(jwtauth.Authenticator)       // handle valid and invalid JWT
			r.Mount("/", api.Router())         // mount the api sub-router
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
