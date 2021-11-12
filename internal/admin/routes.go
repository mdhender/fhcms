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

package admin

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
)

// Router returns a router that will handle all administrator routes
func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(adminOnly)
	r.Get("/", adminIndex)
	r.Get("/accounts", adminListAccounts)
	return r
}

// adminOnly is middleware that forwards only if the request has admin rights.
func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, _ := jwtauth.FromContext(r.Context())
		isAdmin, ok := claims["admin"].(bool)
		if !ok || !isAdmin {
			http.Error(w, http.StatusText(403), 403)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func adminIndex(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "adminIndex\n")
}

func adminListAccounts(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "adminListAccounts\n")
}
