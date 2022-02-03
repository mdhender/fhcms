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

package app

import (
	"net/http"
)

func (s *Server) routes(reports, uploads string) {
	//var sf SiteStore = s.site
	//var glf GamesStore = s.games
	//var gf models.GalaxyFetcher = s.ds
	//var spf models.SpecieFetcher = s.ds
	//var pf ProfileStore = s.profiles

	//s.router.HandleFunc("GET", "/", s.authOnly(s.homeGetIndex(sf, s.templates)))
	s.router.HandleFunc("GET", "/", s.homeHandler(s.renderHome))
	s.router.HandleFunc("GET", "/login", s.loginGetHandler)
	s.router.HandleFunc("POST", "/login", s.loginPostHandler)
	s.router.HandleFunc("*", "/logout", s.logoutHandler)
	s.router.HandleFunc("*", "/manifest.json", s.manifestJsonV3)

	//s.router.HandleFunc("GET", "/about", s.authOnly(s.aboutGetHandler(sf, s.templates)))
	//s.router.HandleFunc("GET", "/favicon.ico", http.NotFound)
	//s.router.HandleFunc("GET", "/games", s.authOnly(s.gamesGetIndex(sf, glf, s.templates)))
	//s.router.HandleFunc("GET", "/games/:gameId", s.authOnly(s.gameGetIndex(sf, gf, spf, s.templates)))
	//s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo", s.authOnly(s.gameGetIndex(sf, gf, spf, s.templates)))
	//s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo/turn/:turnNo", s.authOnly(s.gamesSpecieTurnGetIndex(sf, gf, spf, s.templates)))
	//s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo/turn/:turnNo/orders", s.notImplemented)
	//s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo/turn/:turnNo/reports", s.notImplemented)
	//s.router.HandleFunc("GET", "/logo192.png", http.NotFound)
	//s.router.HandleFunc("GET", "/logout", s.handleLogout)
	//s.router.HandleFunc("GET", "/profile", s.authOnly(s.profileGetHandler(sf, pf, s.templates)))
	//
	//s.router.HandleFunc("POST", "/login", s.handlePostLogin)
	//s.router.HandleFunc("POST", "/game/:gameId/specie/:spNo/turn/:turnId/orders", s.notImplemented)

	//////s.router.HandleFunc("GET", "/admin", s.adminOnly(s.handleAdminIndex()))
	////s.router.HandleFunc("GET", "/home", s.handleHomePage(reports))
	////s.router.HandleFunc("GET", "/turn/:turn/orders", s.handleTurnOrders(reports))
	////s.router.HandleFunc("GET", "/turn/:turn/report", s.handleTurnReport(reports))
	////s.router.HandleFunc("GET", "/turn/:turn/upload", s.handleTurnUpload(uploads))
	////s.router.HandleFunc("POST", "/api/authenticate", s.handleAuthenticate())
	////s.router.HandleFunc("POST", "/api/turn/:turn/orders", s.postTurnOrders(uploads))
	////
	//////s.router.HandleFunc("GET", "/api/get-cookie", s.handleGetCookie())
	//////s.router.HandleFunc("GET", "/api/set-cookie", s.handleSetCookie())

	s.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/go-chi/chi/v5"
//	"github.com/go-chi/chi/v5/middleware"
//	"github.com/go-chi/cors"
//	"github.com/mdhender/fhcms/internal/api"
//	"io"
//	"log"
//	"net/http"
//	"unicode/utf8"
//)
//
//// Routes is routes and returns a chi.Mux.
//func Routes() *chi.Mux {
//	r := chi.NewRouter()
//	r.Use(middleware.Logger)
//	r.Use(middleware.Recoverer)
//	r.Use(cors.Handler(cors.Options{
//		AllowedOrigins: []string{"*"},
//		AllowedMethods: []string{
//			"GET", "PUT", "POST", "DELETE", "HEAD", "OPTIONS",
//		},
//		AllowedHeaders: []string{
//			"Accept", "Accept-Encoding", "Accept-Language", "Authorization", "Cache-Control", "Connection", "Content-Type", "DNT", "Host", "Origin", "Pragma", "Referer", "User-Agent",
//		},
//		ExposedHeaders:   []string{"Link"},
//		AllowCredentials: true,
//		MaxAge:           300, // Maximum value not ignored by any of major browsers
//	}))
//
//	//// auth/auth routes
//	//tokenTimeToLive := 60 * time.Second // 14 * 24 * 60 * 60 * time.Second
//	//r.Post("/auth/jwt")
//
//	// protected api routes
//	r.Route("/api", func(r chi.Router) {
//		//r.Use(jwtauth.Verifier(tokenAuth)) // extract, verify, validate JWT
//		////r.Use(jwtauth.Authenticator)       // handle valid and invalid JWT
//		//r.Use(JWTAuthenticator)    // handle valid and invalid JWT
//		r.Mount("/", api.Router()) // mount the api sub-router
//	})
//
//	// protected ui routes
//	r.Route("/ui", func(r chi.Router) {
//		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
//			_, _ = fmt.Fprint(w, "home\n")
//		})
//	})
//
//	r.Get("/login", func(w http.ResponseWriter, r *http.Request) {
//		_, _ = fmt.Fprint(w, "login\n")
//	})
//	r.Post("/login", PostLogin)
//	r.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
//		http.SetCookie(w, &http.Cookie{Name: "jwt", Path: "/", MaxAge: -1, HttpOnly: true})
//	})
//	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
//		http.SetCookie(w, &http.Cookie{Name: "jwt", Path: "/", MaxAge: -1, HttpOnly: true})
//	})
//
//	return r
//}
