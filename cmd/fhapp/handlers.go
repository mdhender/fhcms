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
	"fmt"
	"github.com/mdhender/fhcms/cmd/fhapp/internal/way"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

func (s *Server) handleAdminIndex() http.HandlerFunc {
	format := "Hello, %s!"
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.DS.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, format, "World")
	}
}

func (s *Server) handleAuthenticate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("server: %s %q: %+v\n", r.Method, r.URL.Path, err)
			http.SetCookie(w, &http.Cookie{
				Name:     s.sessions.cookieName,
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
			})
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Printf("server: %s %q: %v\n", r.Method, r.URL.Path, r.PostForm)
		var input struct {
			username string
			password string
		}
		for k, v := range r.Form {
			switch k {
			case "username":
				if len(v) != 1 || !utf8.ValidString(v[0]) {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.username = v[0]
			case "password":
				if len(v) != 1 || !utf8.ValidString(v[0]) {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.password = v[0]
			}
		}
		if input.username == "" || input.password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Printf("server: %s %q: %v\n", r.Method, r.URL.Path, input)
		var user UserData
		for _, p := range s.data.Players {
			if p.Password == input.password && p.User == input.username {
				for _, sp := range s.data.DS.Species {
					if sp.Id != p.Species {
						continue
					}
					user.Player = p.User // confusing, I know
					user.Species = sp
					user.SpeciesId = sp.Id
					user.IsAuthenticated = true
				}
			}
		}
		if user.Species == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		log.Printf("server: %s %q: %v\n", r.Method, r.URL.Path, user)
		sess := s.sessions.SessionStart(user)
		http.SetCookie(w, &http.Cookie{
			Name:     s.sessions.cookieName,
			Path:     "/",
			Value:    sess.Uuid,
			MaxAge:   7 * 24 * 60 * 60,
			HttpOnly: true,
		})
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// handleCachedFile will return a document along with caching instructions.
func (s *Server) handleCachedFile(name string) http.HandlerFunc {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Printf("server: cache: %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
	etag := mketag(b)
	maxAge := fmt.Sprintf("max-age=%d", 4*7*24*60*60) // 4 weeks
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleCachedFile\n", r.Method, r.URL.Path)
		if match := r.Header.Get("If-None-Match"); match != "" {
			if strings.Contains(match, etag) {
				log.Printf("server: %s %q: handleCachedFile: cached\n", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotModified)
				return
			}
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.DS.Semver)
		w.Header().Set("Etag", etag)
		w.Header().Set("Cache-Control", maxAge)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Server) handleGetCookie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(s.sessions.cookieName)
		if err != nil {
			_, _ = fmt.Fprintf(w, "GET cookie(%q): %+v", s.sessions.cookieName, err)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.DS.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "GET cookie(%q): %q", s.sessions.cookieName, c.Value)
	}
}

func (s *Server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if cookie, err := r.Cookie(s.sessions.cookieName); err == nil {
			s.sessions.SessionDelete(cookie.Value)
		}
		http.SetCookie(w, &http.Cookie{
			Name:     s.sessions.cookieName,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		w.Header().Set("Far-Horizons", s.data.DS.Semver)
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *Server) handleSetCookie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     s.sessions.cookieName,
			Path:     "/",
			Value:    "thisIsMyCookie",
			HttpOnly: true,
		})
	}
}

func (s *Server) handleTurnOrders() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleTurnOrders\n", r.Method, r.URL.Path)
		u := currentUser(r)
		if u.SpeciesId == 0 || u.Species == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		turnNumber, err := strconv.Atoi(way.Param(r.Context(), "turn"))
		if err != nil || turnNumber < 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		data := struct {
			Engine     *Engine
			DS         *JDB
			Site       *Site
			TurnNumber int
			Date       string
			Orders     string
			User       UserData
		}{
			Engine:     s.data.Engine,
			DS:         s.data.DS,
			Site:       s.data.Site,
			User:       u,
			TurnNumber: turnNumber,
		}
		for _, f := range s.data.Files[u.SpeciesId] {
			if f.Turn == turnNumber {
				data.Date = f.Date
				if f.Orders != "" {
					if b, err := ioutil.ReadFile(filepath.Join("testdata", "files", f.Orders)); err != nil {
						log.Printf("server: %s %q: handleTurnOrders: %+v\n", r.Method, r.URL.Path, err)
					} else {
						data.Orders = string(b)
					}
				}
				break
			}
		}
		b, err := s.render("turnOrders", data)
		if err != nil {
			log.Printf("server: %s %q: handleTurnOrders: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.DS.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Server) handleTurnReport() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleTurnReport\n", r.Method, r.URL.Path)
		u := currentUser(r)
		if u.SpeciesId == 0 || u.Species == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		turnNumber, err := strconv.Atoi(way.Param(r.Context(), "turn"))
		if err != nil || turnNumber < 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		data := struct {
			Engine     *Engine
			DS         *JDB
			Site       *Site
			TurnNumber int
			Date       string
			Report     string
			User       UserData
		}{
			Engine:     s.data.Engine,
			DS:         s.data.DS,
			Site:       s.data.Site,
			User:       u,
			TurnNumber: turnNumber,
		}
		for _, f := range s.data.Files[u.SpeciesId] {
			if f.Turn == turnNumber {
				data.Date = f.Date
				if f.Report != "" {
					if b, err := ioutil.ReadFile(filepath.Join("testdata", "files", f.Report)); err != nil {
						log.Printf("server: %s %q: handleTurnOrders: %+v\n", r.Method, r.URL.Path, err)
					} else {
						data.Report = string(b)
					}
				}
				break
			}
		}
		b, err := s.render("turnReport", data)
		if err != nil {
			log.Printf("server: %s %q: handleTurnReport: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.DS.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Server) handleUI() http.HandlerFunc {
	type turnFile struct {
		Turn   int
		Report string
		Orders string
		Date   string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleUI\n", r.Method, r.URL.Path)
		u := currentUser(r)
		data := struct {
			Engine *Engine
			DS     *JDB
			Game   struct {
				Title     string
				Turn      int
				NextTurn  int
				OrdersDue string
			}
			Site  Site
			Files []*turnFile
			User  UserData
			Stats *StatsData
		}{
			Engine: s.data.Engine,
			DS:     s.data.DS,
			User:   u,
			Stats:  s.data.Stats[u.SpeciesId],
		}
		data.Game.Title = "Raven's Beta"
		data.Game.Turn = s.data.Cluster.Turn
		data.Game.NextTurn = s.data.Cluster.Turn + 1
		data.Game.OrdersDue = "Monday, September 20th by 7PM MDT. MDT is 6 hours behind London."
		data.Site.Slug = s.data.Site.Slug
		data.Site.Copyright = s.data.Site.Copyright
		data.Site.Title = s.data.Site.Title

		for _, f := range s.data.Files[u.SpeciesId] {
			tf := &turnFile{Turn: f.Turn, Date: f.Date, Report: f.Report, Orders: f.Orders}
			data.Files = append(data.Files, tf)
		}
		b, err := s.render("index", data)
		if err != nil {
			log.Printf("server: %s %q: handleUI: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.DS.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}
