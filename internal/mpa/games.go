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
	"github.com/mdhender/fhcms/internal/way"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// fetch specific game for the current user
func (s *Server) gameGetIndex(sf models.SiteFetcher, gf models.GameFetcher, templates string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		u := s.currentUser(r)
		gameId := way.Param(r.Context(), "gameId")
		log.Printf("mpa: gameGetIndex: u.id %q gameId %q\n", u.Id, gameId)

		t, err := template.ParseFiles(filepath.Join(templates, "game.index.gohtml"))
		if err != nil {
			log.Printf("mpa: gameGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var payload struct {
			Site *models.Site
			Game *models.Game
		}
		payload.Site = sf.FetchSite()
		payload.Game = gf.FetchGame(u.Id, gameId)

		b := &bytes.Buffer{}
		if err := t.Execute(b, payload); err != nil {
			log.Printf("mpa: gameGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b.Bytes())
	}
}

// fetch list of games for the current user
func (s *Server) gamesGetIndex(sf models.SiteFetcher, glf models.GamesFetcher, templates string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		u := s.currentUser(r)
		log.Printf("mpa: gamesGetIndex: u.id %q\n", u.Id)

		t, err := template.ParseFiles(filepath.Join(templates, "games.index.gohtml"))
		if err != nil {
			log.Printf("mpa: gamesGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var payload struct {
			Site  *models.Site
			Games models.Games
		}
		payload.Site = sf.FetchSite()
		payload.Games = glf.FetchGames(u.Id)

		b := &bytes.Buffer{}
		if err := t.Execute(b, payload); err != nil {
			log.Printf("mpa: gamesGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b.Bytes())
	}
}
