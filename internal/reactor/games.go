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

package reactor

import (
	"bytes"
	"github.com/mdhender/fhcms/internal/models"
	"github.com/mdhender/fhcms/internal/way"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
)

type GamesStore interface {
	FetchGalaxies(uid int) (models.Galaxies, bool)
}

// fetch specific game for the current user
func (s *Server) gameGetIndex(sf SiteStore, gf models.GalaxyFetcher, spf models.SpecieFetcher, templates string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		u := s.currentUser(r)
		gameId := way.Param(r.Context(), "gameId")
		gid, err := strconv.Atoi(gameId)
		if err != nil || gid < 1 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		spNo := way.Param(r.Context(), "spNo")
		spid, err := strconv.Atoi(spNo)
		if err != nil || spid < 1 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var turnNo int
		if qParm := way.Param(r.Context(), "turnNo"); qParm != "" {
			if val, err := strconv.Atoi(qParm); err != nil {
				log.Printf("mpa: gameGetIndex: u.id %q gameId %q spNo %q turnNo %q: %+v\n", u.Id, gameId, spNo, way.Param(r.Context(), "turnNo"), err)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			} else {
				turnNo = val
			}
		}
		log.Printf("mpa: gameGetIndex: u.id %q gameId %q spNo %q turnNo %d\n", u.Id, gameId, spNo, turnNo)
		if turnNo < 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		t, err := template.ParseFiles(filepath.Join(templates, "site.layout.gohtml"), filepath.Join(templates, "fragments", "navbar.gohtml"), filepath.Join(templates, "fragments", "footer.gohtml"), filepath.Join(templates, "game.index.gohtml"))
		if err != nil {
			log.Printf("mpa: gameGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var payload struct {
			Site   models.Site
			Game   *models.Galaxy
			Specie *models.Specie
		}
		payload.Site, _ = sf.FetchSite()
		payload.Game = gf.FetchGalaxy(u.Id, gid)
		if turnNo == 0 {
			turnNo = payload.Game.CurrentTurn
			log.Printf("mpa: gameGetIndex: u.id %q gameId %q spNo %q turnNo 0 => %d\n", u.Id, gameId, spNo, turnNo)
		}
		payload.Game.TurnNo = turnNo
		payload.Game.Display.Deadline = payload.Game.TurnNo == payload.Game.CurrentTurn
		payload.Specie = spf.FetchSpecie(u.Id, gid, spid, turnNo)

		b := &bytes.Buffer{}
		if err = t.ExecuteTemplate(b, "layout", payload); err != nil {
			log.Printf("mpa: gameGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b.Bytes())
	}
}

// fetch list of games for the current user
func (s *Server) gamesGetIndex(sf SiteStore, glf GamesStore, templates string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		u := s.currentUser(r)
		log.Printf("mpa: gamesGetIndex: u.id %d\n", u.Id)

		t, err := template.ParseFiles(filepath.Join(templates, "site.layout.gohtml"), filepath.Join(templates, "fragments", "navbar.gohtml"), filepath.Join(templates, "fragments", "footer.gohtml"), filepath.Join(templates, "games.index.gohtml"))
		if err != nil {
			log.Printf("mpa: gamesGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var payload struct {
			Site  models.Site
			Games models.Galaxies
		}
		payload.Site, _ = sf.FetchSite()
		payload.Games, _ = glf.FetchGalaxies(u.Id)

		b := &bytes.Buffer{}
		if err = t.ExecuteTemplate(b, "layout", payload); err != nil {
			log.Printf("mpa: gamesGetIndex: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b.Bytes())
	}
}

// fetch specific game turn for the current user
func (s *Server) gamesSpecieTurnGetIndex(sf SiteStore, gf models.GalaxyFetcher, spf models.SpecieFetcher, templates string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		u := s.currentUser(r)
		gameId := way.Param(r.Context(), "gameId")
		gid, err := strconv.Atoi(gameId)
		if err != nil || gid < 1 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		spNo := way.Param(r.Context(), "spNo")
		spid, err := strconv.Atoi(spNo)
		if err != nil || spid < 1 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		var turnNo int
		if qParm := way.Param(r.Context(), "turnNo"); qParm != "" {
			if val, err := strconv.Atoi(qParm); err != nil {
				log.Printf("mpa: gamesSpecieTurnGetIndex: u.id %q gameId %q spNo %q turnNo %q: %+v\n", u.Id, gameId, spNo, way.Param(r.Context(), "turnNo"), err)
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			} else {
				turnNo = val
			}
		}
		log.Printf("mpa: gamesSpecieTurnGetIndex: u.id %q gameId %q spNo %q turnNo %d\n", u.Id, gameId, spNo, turnNo)
		if turnNo < 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		t, err := template.ParseFiles(filepath.Join(templates, "site.layout.gohtml"), filepath.Join(templates, "fragments", "navbar.gohtml"), filepath.Join(templates, "fragments", "footer.gohtml"), filepath.Join(templates, "game.turn.gohtml"))
		if err != nil {
			log.Printf("mpa: gamesSpecieTurnGetIndex: u.id %q gameId %q spNo %q turnNo %q: %+v\n", u.Id, gameId, spNo, way.Param(r.Context(), "turnNo"), err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var payload struct {
			Site   models.Site
			Game   *models.Galaxy
			Specie *models.Specie
		}
		payload.Site, _ = sf.FetchSite()
		payload.Game = gf.FetchGalaxy(u.Id, gid)
		if turnNo == 0 {
			turnNo = payload.Game.CurrentTurn
			log.Printf("mpa: gamesSpecieTurnGetIndex: u.id %q gameId %q spNo %q turnNo 0 => %d\n", u.Id, gameId, spNo, turnNo)
		}
		payload.Game.TurnNo = turnNo
		payload.Game.Display.Deadline = payload.Game.TurnNo == payload.Game.CurrentTurn
		payload.Specie = spf.FetchSpecie(u.Id, gid, spid, turnNo)

		b := &bytes.Buffer{}
		if err = t.ExecuteTemplate(b, "layout", payload); err != nil {
			log.Printf("mpa: gamesSpecieTurnGetIndex: u.id %q gameId %q spNo %q turnNo %d: %+v\n", u.Id, gameId, spNo, turnNo, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b.Bytes())
	}
}
