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
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Router returns a router that will serve static files for the EmberJS front end.
func Router() http.Handler {
	r := chi.NewRouter()

	r.Get("/", home)
	r.Get("/contact", contact)
	r.Get("/{slug}", widgetGet)
	r.Get("/{slug}/admin", widgetAdmin)
	r.Post("/{slug}/image", widgetImage)

	return r
}

func home(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "home\n")
}

func contact(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "contact\n")
}

func widgetGet(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	_, _ = fmt.Fprintf(w, "widget %s\n", slug)
}

func widgetAdmin(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	_, _ = fmt.Fprintf(w, "widgetAdmin %s\n", slug)
}

func widgetImage(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	_, _ = fmt.Fprintf(w, "widgetImage %s\n", slug)
}
