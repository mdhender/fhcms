/*******************************************************************************
fhapp - a player aid for Far Horizons
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

package main

import (
	"log"
	"net/http"
)

func (s *Server) routes(reports string) {
	s.router.Handle("GET", "/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.router.HandleFunc("GET", "/favicon.ico", s.handleCachedFile("static/favicon.ico"))
	s.router.HandleFunc("GET", "/rules", s.handleCachedFile("static/rules.html"))

	//s.router.HandleFunc("GET", "/login", s.notImplemented())
	//s.router.HandleFunc("GET", "/logout", s.notImplemented())
	//s.router.HandleFunc("GET", "/admin", s.adminOnly(s.handleAdminIndex()))
	s.router.HandleFunc("GET", "/turn/:turn/orders", s.handleTurnOrders(reports))
	s.router.HandleFunc("GET", "/turn/:turn/report", s.handleTurnReport(reports))
	s.router.HandleFunc("GET", "/", s.handleUI())
	s.router.HandleFunc("POST", "/api/authenticate", s.handleAuthenticate())
	s.router.HandleFunc("GET", "/api/logout", s.handleLogout())
	//s.router.HandleFunc("GET", "/api/get-cookie", s.handleGetCookie())
	//s.router.HandleFunc("GET", "/api/set-cookie", s.handleSetCookie())
	s.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: not found\n", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}
