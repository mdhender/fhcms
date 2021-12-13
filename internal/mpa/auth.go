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
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/users"
	"log"
	"net/http"
)

// we must serve only the login page to unauthenticated visitors
func (s *Server) authOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u users.User
		if j := jot.FromRequest(r); s.jf.Validate(j) == nil {
			u.Id = j.UserID()
			if acct, ok := s.accts.ById[u.Id]; ok {
				u.Name = acct.Username
				u.IsAuthenticated = u.Id != ""
			}
		}
		if !u.IsAuthenticated {
			log.Printf("mw: authOnly: %s %s: !authenticated\n", r.Method, r.URL.Path)
			s.handleGetLogin(w, r)
			return
		}
		h(w, r.WithContext(u.NewContext(r.Context())))
	}
}
