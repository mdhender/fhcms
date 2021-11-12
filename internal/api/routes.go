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

package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"net/http"
	"strconv"
)

type Server http.Handler

// Router returns a router that will serve REST-ish data for the EmberJS front end.
func Router( /*ta *jwtauth.JWTAuth*/ ) http.Handler {
	r := chi.NewRouter()

	r.Get("/widgets", apiGetWidgets)
	r.Post("/widgets", apiCreateWidget)
	r.Post("/widgets/{slug}", apiUpdateWidget)
	r.Post("/widgets/{slug}/parts", apiCreateWidgetPart)
	r.Post("/widgets/{slug}/parts/{id:[0-9]+}/update", apiUpdateWidgetPart)
	r.Post("/widgets/{slug}/parts/{id:[0-9]+}/delete", apiDeleteWidgetPart)

	var s Server
	s = r
	return s
}

func apiGetWidgets(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	_, _ = fmt.Fprintf(w, "apiGetWidgets: claims %v\n", claims["user_id"])
}

func apiCreateWidget(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	_, _ = fmt.Fprintf(w, "apiCreateWidget: claims %v\n", claims["user_id"])
}

func apiUpdateWidget(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	_, _ = fmt.Fprintf(w, "apiUpdateWidget: claims %v: slug %s\n", claims["user_id"], slug)
}

func apiCreateWidgetPart(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	slug := chi.URLParam(r, "slug")
	_, _ = fmt.Fprintf(w, "apiCreateWidgetPart: claims %v: slug %s\n", claims["user_id"], slug)
}

func apiUpdateWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, _ = fmt.Fprintf(w, "apiUpdateWidgetPart %s %d\n", slug, id)
}

func apiDeleteWidgetPart(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	_, _ = fmt.Fprintf(w, "apiDeleteWidgetPart %s %d\n", slug, id)
}
