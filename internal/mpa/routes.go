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

package mpa

import (
	"github.com/mdhender/fhcms/internal/models"
	"net/http"
)

func (s *Server) routes(reports, uploads string) {
	var sf models.SiteFetcher = s.site
	var glf models.GamesFetcher = s.ds
	var gf models.GameFetcher = s.ds
	var spf models.SpecieFetcher = s.ds

	s.router.HandleFunc("GET", "/", s.authOnly(s.homeGetIndex(sf, s.templates)))
	s.router.HandleFunc("GET", "/manifest.json", s.manifestJsonV3)
	s.router.HandleFunc("GET", "/about", s.authOnly(s.aboutGetHandler(sf, s.templates)))
	s.router.HandleFunc("GET", "/favicon.ico", http.NotFound)
	s.router.HandleFunc("GET", "/games", s.authOnly(s.gamesGetIndex(sf, glf, s.templates)))
	s.router.HandleFunc("GET", "/games/:gameId", s.authOnly(s.gameGetIndex(sf, gf, spf, s.templates)))
	s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo", s.authOnly(s.gameGetIndex(sf, gf, spf, s.templates)))
	s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo/turn/:turnNo", s.authOnly(s.gameGetIndex(sf, gf, spf, s.templates)))
	s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo/turn/:turnNo/orders", s.notImplemented)
	s.router.HandleFunc("GET", "/games/:gameId/specie/:spNo/turn/:turnNo/reports", s.notImplemented)
	s.router.HandleFunc("GET", "/logo192.png", http.NotFound)
	s.router.HandleFunc("GET", "/logout", s.handleLogout)
	s.router.HandleFunc("GET", "/manifest.json", http.NotFound)
	s.router.HandleFunc("GET", "/profile", s.authOnly(s.profileGetHandler(sf, s.ds, s.templates)))

	s.router.HandleFunc("POST", "/login", s.handlePostLogin)
	s.router.HandleFunc("POST", "/logout", s.handleLogout)
	s.router.HandleFunc("POST", "/game/:gameId/specie/:spNo/turn/:turnId/orders", s.notImplemented)

	////s.router.HandleFunc("GET", "/admin", s.adminOnly(s.handleAdminIndex()))
	//s.router.HandleFunc("GET", "/home", s.handleHomePage(reports))
	//s.router.HandleFunc("GET", "/turn/:turn/orders", s.handleTurnOrders(reports))
	//s.router.HandleFunc("GET", "/turn/:turn/report", s.handleTurnReport(reports))
	//s.router.HandleFunc("GET", "/turn/:turn/upload", s.handleTurnUpload(uploads))
	//s.router.HandleFunc("POST", "/api/authenticate", s.handleAuthenticate())
	//s.router.HandleFunc("POST", "/api/turn/:turn/orders", s.postTurnOrders(uploads))
	//
	////s.router.HandleFunc("GET", "/api/get-cookie", s.handleGetCookie())
	////s.router.HandleFunc("GET", "/api/set-cookie", s.handleSetCookie())

	s.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}

func (s *Server) manifestJsonV3(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"manifest_version":3,"name":"My Extension","version":"versionString"}`))
}

func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
