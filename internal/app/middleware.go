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

package app

import (
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/models"
	"log"
	"net/http"
)

// CORS adds a non-functional CORS handler
func (s *Server) CORS(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[app] cors\n")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, DELETE, HEAD, OPTIONS, POST, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Accept Accept-Encoding Accept-Language Authorization Cache-Control Connection Content-Type DNT Host Origin Pragma Referer User-Agent")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	}
}

// WithUser extracts the cookie or bearer token from the request
// and injects it into the request context.
func (s *Server) WithUser(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[app] withUser\n")
		var a models.Account
		if j := jot.FromRequest(r); s.jf.Validate(j) == nil {
			if id, ok := j.UserID(); ok {
				if x, ok := s.accountStore.FetchById(id); ok {
					a.Id, a.UserName, a.IsActive, a.IsAuthenticated = x.Id, x.UserName, x.IsActive, true
					// prevent inactive unauthenticated or users from having admin role
					a.IsAdmin = a.IsActive && a.IsAuthenticated && x.IsAdmin
				}
			}
		}
		h(w, r.WithContext(a.NewContext(r.Context())))
	}
}
