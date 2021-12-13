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

package adapters

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// DumpRequest returns an adapter that dumps the request header.
func DumpRequest(next http.Handler) http.HandlerFunc {
	log.Printf("[dump] initializing\n")
	defer log.Printf("[dump] initialized\n")

	return func(w http.ResponseWriter, r *http.Request) {
		buf, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Printf("[dump] %s: %q\n", r.URL.Path, err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		log.Printf("[dump] %s\n%s\n", r.URL.Path, string(buf))

		if next != nil {
			next.ServeHTTP(w, r)
		}
	}
}
