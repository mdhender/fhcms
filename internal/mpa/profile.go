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
	"bytes"
	"github.com/mdhender/fhcms/internal/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func (s *Server) profileGetHandler(sf models.SiteFetcher, uf models.UserFetcher, templates string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		t, err := template.ParseFiles(filepath.Join(templates, "site.layout.gohtml"), filepath.Join(templates, "fragments", "navbar.gohtml"), filepath.Join(templates, "fragments", "footer.gohtml"), filepath.Join(templates, "profile.index.gohtml"))
		if err != nil {
			log.Printf("mpa: profileGetHandler: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var payload struct {
			Site *models.Site
			User *models.User
		}
		payload.Site = sf.FetchSite()
		payload.User = uf.FetchUser(s.currentUser(r).Id)

		b := &bytes.Buffer{}
		if err = t.ExecuteTemplate(b, "layout", payload); err != nil {
			log.Printf("mpa: profileGetHandler: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b.Bytes())
	}
}
