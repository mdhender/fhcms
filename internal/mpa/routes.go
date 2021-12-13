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
	"net/http"
)

func (s *Server) routes(reports, uploads string) {
	s.router.HandleFunc("GET", "/", s.authOnly(s.gamesGetIndex(s.site, s.ds, s.templates)))
	s.router.HandleFunc("GET", "/about", s.authOnly(s.aboutGetHandler(s.site, s.templates)))
	s.router.HandleFunc("GET", "/favicon.ico", http.NotFound)
	s.router.HandleFunc("GET", "/game/:gameId", s.authOnly(s.gameGetIndex(s.site, s.ds, s.templates)))
	s.router.HandleFunc("GET", "/game/:gameId/turn/:turnId", s.notImplemented)
	s.router.HandleFunc("GET", "/game/:gameId/turn/:turnId/orders", s.notImplemented)
	s.router.HandleFunc("GET", "/game/:gameId/turn/:turnId/reports", s.notImplemented)
	s.router.HandleFunc("GET", "/games", s.authOnly(s.gamesGetIndex(s.site, s.ds, s.templates)))
	s.router.HandleFunc("GET", "/logo192.png", http.NotFound)
	s.router.HandleFunc("GET", "/logout", s.handleLogout)
	s.router.HandleFunc("GET", "/manifest.json", http.NotFound)
	s.router.HandleFunc("GET", "/profile", s.authOnly(s.profileGetHandler(s.site, s.ds, s.templates)))

	s.router.HandleFunc("POST", "/login", s.handlePostLogin)
	s.router.HandleFunc("POST", "/logout", s.handleLogout)
	s.router.HandleFunc("POST", "/game/:gameId/turn/:turnId/orders", s.notImplemented)

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

func (s *Server) notImplemented(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
}
