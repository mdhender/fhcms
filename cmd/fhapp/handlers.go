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
	"bytes"
	"fmt"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/mdhender/fhcms/internal/flist"
	"github.com/mdhender/fhcms/internal/way"
	"github.com/mdhender/fhcms/parser"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func (s *Server) handleAdminIndex() http.HandlerFunc {
	format := "Hello, %s!"
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
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
		log.Printf("server: %s %q: username %q password %q\n", r.Method, r.URL.Path, input.username, input.password)
		if input.username == "" || input.password == "" {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Printf("server: %s %q: %v\n", r.Method, r.URL.Path, input)
		var user UserData
		for _, p := range s.data.Players {
			log.Printf("server: %s %q: player %v\n", r.Method, r.URL.Path, *p)
			if p.Password == input.password && p.User == input.username {
				log.Printf("server: %s %q: player user %q: password matched\n", r.Method, r.URL.Path, p.User)
				sp, ok := s.data.Store.Species[p.SpeciesId]
				if !ok {
					log.Printf("server: %s %q: player user %q species %q: no such species\n", r.Method, r.URL.Path, p.User, p.SpeciesId)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					continue
				}
				user.Player = p.User // confusing, I know
				user.Species = sp
				user.SpeciesId = p.SpeciesId
				user.IsAuthenticated = true
			}
		}
		if user.Species == nil {
			log.Printf("server: %s %q: %v: use.Species is nil\n", r.Method, r.URL.Path, input)
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
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
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
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
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
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
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

func (s *Server) handleTurnOrders(files string) http.HandlerFunc {
	log.Printf("serving turn orders from %q\n", files)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleTurnOrders\n", r.Method, r.URL.Path)
		u := currentUser(r)
		if u.SpeciesId == "" || u.Species == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		turnNumber, err := strconv.Atoi(way.Param(r.Context(), "turn"))
		if err != nil || turnNumber < 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		data := struct {
			Engine *Engine
			Semver string
			Site   struct {
				Title     string
				Slug      string
				Copyright struct {
					Year   int
					Author string
				}
			}
			Game struct {
				Title     string
				Turn      int
				LastTurn  int
				OrdersDue string
			}
			User   UserData
			Player struct {
				Name            string
				Data            string // folder on web server containing this player's data
				IsAdmin         bool
				IsAuthenticated bool
				Species         *cluster.Species
			}
			Stats      *StatsData
			TurnNumber int
			Date       string
			Orders     string
		}{
			Engine:     s.data.Engine,
			User:       u,
			Stats:      s.data.Stats[u.SpeciesId],
			TurnNumber: turnNumber,
		}
		data.Semver = s.data.Store.Semver
		data.Game.Title = "Raven's Beta"
		data.Game.Turn = s.data.Store.Turn
		if data.Game.Turn > 1 {
			data.Game.LastTurn = s.data.Store.Turn - 1
		}
		data.Game.OrdersDue = fmt.Sprintf("%s by %s. %s", s.data.Turn.Due, s.data.Turn.By, s.data.Turn.TimeZone)
		data.Player.Name = u.Player
		data.Player.Data = u.SpeciesId + "?key?"
		data.Player.IsAuthenticated = u.IsAuthenticated
		data.Player.IsAdmin = u.IsAuthenticated && u.IsAdmin
		data.Player.Species = u.Species
		data.Site.Title = s.data.Site.Title
		data.Site.Slug = s.data.Site.Slug
		data.Site.Copyright.Year = s.data.Site.Copyright.Year
		data.Site.Copyright.Author = s.data.Site.Copyright.Author

		for _, f := range s.data.Files[u.SpeciesId] {
			if f.Turn == turnNumber {
				data.Date = f.Date
				if f.Orders != "" {
					if b, err := ioutil.ReadFile(filepath.Join(files, f.Orders)); err != nil {
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
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Server) handleTurnReport(files string) http.HandlerFunc {
	log.Printf("serving turn reports from %q\n", files)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleTurnReport\n", r.Method, r.URL.Path)
		u := currentUser(r)
		if u.SpeciesId == "" || u.Species == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		turnNumber, err := strconv.Atoi(way.Param(r.Context(), "turn"))
		if err != nil || turnNumber < 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		data := struct {
			Engine *Engine
			Semver string
			Site   struct {
				Title     string
				Slug      string
				Copyright struct {
					Year   int
					Author string
				}
			}
			Game struct {
				Title     string
				Turn      int
				LastTurn  int
				OrdersDue string
			}
			User   UserData
			Player struct {
				Name            string
				Data            string // folder on web server containing this player's data
				IsAdmin         bool
				IsAuthenticated bool
				Species         *cluster.Species
			}
			Stats      *StatsData
			TurnNumber int
			Date       string
			Report     string
		}{
			Engine:     s.data.Engine,
			User:       u,
			Stats:      s.data.Stats[u.SpeciesId],
			TurnNumber: turnNumber,
		}
		data.Semver = s.data.Store.Semver
		data.Game.Title = "Raven's Beta"
		data.Game.Turn = s.data.Store.Turn
		if data.Game.Turn > 1 {
			data.Game.LastTurn = s.data.Store.Turn - 1
		}
		data.Game.OrdersDue = fmt.Sprintf("%s by %s. %s", s.data.Turn.Due, s.data.Turn.By, s.data.Turn.TimeZone)
		data.Player.Name = u.Player
		data.Player.Data = u.SpeciesId + "?key?"
		data.Player.IsAuthenticated = u.IsAuthenticated
		data.Player.IsAdmin = u.IsAuthenticated && u.IsAdmin
		data.Player.Species = u.Species
		data.Site.Title = s.data.Site.Title
		data.Site.Slug = s.data.Site.Slug
		data.Site.Copyright.Year = s.data.Site.Copyright.Year
		data.Site.Copyright.Author = s.data.Site.Copyright.Author

		for _, f := range s.data.Files[u.SpeciesId] {
			if f.Turn == turnNumber {
				data.Date = f.Date
				if f.Report != "" {
					if b, err := ioutil.ReadFile(filepath.Join(files, f.Report)); err != nil {
						log.Printf("server: %s %q: handleTurnReport: %+v\n", r.Method, r.URL.Path, err)
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
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Server) handleTurnUpload(uploads string) http.HandlerFunc {
	log.Printf("uploading turn files to %q\n", uploads)
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleTurnUpload\n", r.Method, r.URL.Path)
		u := currentUser(r)
		if u.SpeciesId == "" || u.Species == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		turnNumber, err := strconv.Atoi(way.Param(r.Context(), "turn"))
		if err != nil || turnNumber < 1 {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		data := struct {
			Engine *Engine
			Semver string
			Site   struct {
				Title     string
				Slug      string
				Copyright struct {
					Year   int
					Author string
				}
			}
			Game struct {
				Title     string
				Turn      int
				LastTurn  int
				OrdersDue string
			}
			User   UserData
			Player struct {
				Name            string
				Data            string // folder on web server containing this player's data
				IsAdmin         bool
				IsAuthenticated bool
				Species         *cluster.Species
			}
			Stats      *StatsData
			TurnNumber int
			OrdersFile string
			Report     string
		}{
			Engine:     s.data.Engine,
			User:       u,
			Stats:      s.data.Stats[u.SpeciesId],
			TurnNumber: turnNumber,
		}
		data.Semver = s.data.Store.Semver
		data.Game.Title = "Raven's Beta"
		data.Game.Turn = s.data.Store.Turn
		if data.Game.Turn > 1 {
			data.Game.LastTurn = s.data.Store.Turn - 1
		}
		data.Game.OrdersDue = fmt.Sprintf("%s by %s. %s", s.data.Turn.Due, s.data.Turn.By, s.data.Turn.TimeZone)
		data.Player.Name = u.Player
		data.Player.Data = u.SpeciesId + "?key?"
		data.Player.IsAuthenticated = u.IsAuthenticated
		data.Player.IsAdmin = u.IsAuthenticated && u.IsAdmin
		data.Player.Species = u.Species
		data.OrdersFile = fmt.Sprintf("sp%02d.t%d.orders.txt", u.Species.No, s.data.Store.Turn)
		data.Site.Title = s.data.Site.Title
		data.Site.Slug = s.data.Site.Slug
		data.Site.Copyright.Year = s.data.Site.Copyright.Year
		data.Site.Copyright.Author = s.data.Site.Copyright.Author
		b, err := s.render("turnUpload", data)
		if err != nil {
			log.Printf("server: %s %q: handleTurnReport: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Server) handleUI() http.HandlerFunc {
	type turnFile struct {
		Turn   int
		Report string
		Orders string
		Combat string
		Date   string
	}
	type otherSpecies struct {
		Name   string
		Status string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("server: %s %q: handleUI\n", r.Method, r.URL.Path)
		u := currentUser(r)
		data := struct {
			Engine *Engine
			Semver string
			Site   struct {
				Title     string
				Slug      string
				Copyright struct {
					Year   int
					Author string
				}
			}
			Game struct {
				Title     string
				Turn      int
				LastTurn  int
				OrdersDue string
			}
			Files  []*turnFile
			Player struct {
				Name            string
				Data            string // folder on web server containing this player's data
				IsAdmin         bool
				IsAuthenticated bool
				Species         *cluster.Species
			}
			User    UserData
			Stats   *StatsData
			Species []*otherSpecies
		}{
			Engine: s.data.Engine,
			User:   u,
			Stats:  s.data.Stats[u.SpeciesId],
		}
		data.Semver = s.data.Store.Semver
		data.Game.Title = "Raven's Beta"
		data.Game.Turn = s.data.Store.Turn
		if data.Game.Turn > 1 {
			data.Game.LastTurn = s.data.Store.Turn - 1
		}
		data.Game.OrdersDue = fmt.Sprintf("%s by %s. %s", s.data.Turn.Due, s.data.Turn.By, s.data.Turn.TimeZone)
		data.Player.Name = u.Player
		data.Player.Data = u.SpeciesId + "?key?"
		data.Player.IsAuthenticated = u.IsAuthenticated
		data.Player.IsAdmin = u.IsAuthenticated && u.IsAdmin
		data.Player.Species = u.Species
		data.Site.Title = s.data.Site.Title
		data.Site.Slug = s.data.Site.Slug
		data.Site.Copyright.Year = s.data.Site.Copyright.Year
		data.Site.Copyright.Author = s.data.Site.Copyright.Author

		var turnFiles []*flist.TurnData
		if u.Species != nil {
			os := make(map[string]string)
			for _, o := range u.Species.Contact {
				os[o.Name] = "Neutral"
			}
			for _, o := range u.Species.Ally {
				os[o.Name] = "Ally"
			}
			for _, o := range u.Species.Enemy {
				os[o.Name] = "Enemy"
			}
			for name, status := range os {
				data.Species = append(data.Species, &otherSpecies{Name: name, Status: status})
			}
			for i := 0; i < len(data.Species); i++ { // bubbly and proud of it
				for j := i + 1; j < len(data.Species); j++ {
					if data.Species[i].Name > data.Species[j].Name {
						data.Species[i], data.Species[j] = data.Species[j], data.Species[i]
					}
				}
			}
			if turnFiles, _ = flist.Fetch(s.data.TurnFiles, u.Species.No); turnFiles != nil {
				for _, f := range turnFiles {
					log.Printf("[turn] %2d %2d %s %q %q\n", u.Species.No, f.Turn, f.Date, f.Orders, f.Report)
				}
			}
		}

		//for _, f := range s.data.Files[u.SpeciesId] {
		//	tf := &turnFile{Turn: f.Turn, Date: f.Date, Report: f.Report, Orders: f.Orders}
		//	data.Files = append(data.Files, tf)
		//}
		for _, f := range turnFiles {
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
		w.Header().Set("Far-Horizons", s.data.Store.Semver)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(b)
	}
}

func (s *Server) postTurnOrders(uploads string) http.HandlerFunc {
	log.Printf("posting turn orders to %q\n", uploads)
	return func(w http.ResponseWriter, r *http.Request) {
		u := currentUser(r)
		if !u.IsAuthenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		turnNumber, err := strconv.Atoi(way.Param(r.Context(), "turn"))
		if err != nil || turnNumber != s.data.Store.Turn {
			log.Printf("postTurnOrders: turnNumber %d s.data.Store.Turn %d\n", turnNumber, s.data.Store.Turn)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("server: %s %q: %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		//log.Printf("server: %s %q: %v\n", r.Method, r.URL.Path, r.PostForm)
		var input struct {
			orders string
		}
		for k, v := range r.Form {
			switch k {
			case "orders":
				if len(v) != 1 || !utf8.ValidString(v[0]) || len(v[0]) < 1 || len(v[0]) > 64*1024 {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.orders = v[0]
			}
		}

		if len(input.orders) < 1 || len(input.orders) > 64*1024 || !utf8.ValidString(input.orders) {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		date := time.Now().UTC().Format(time.RFC3339)
		input.orders = fmt.Sprintf(";; %s T%d %s\n\n", u.SpeciesId, turnNumber, date) + input.orders

		ordersFile := fmt.Sprintf("sp%02d.t%d.orders.txt", u.Species.No, s.data.Store.Turn)
		fullOrdersFile := filepath.Join(uploads, ordersFile)

		log.Printf("server: %s %q: species %s turn %d orders %s\n", r.Method, r.URL.Path, u.SpeciesId, turnNumber, ordersFile)
		if err := ioutil.WriteFile(fullOrdersFile, []byte(input.orders), 0644); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		log.Printf("orders: loading orders file %q\n", ordersFile)

		// parser returns a raw parse tree. we will evaluate it later, maybe.
		var report bytes.Buffer
		var errs []error
		tree, err := parser.Parse(fullOrdersFile, false)
		if err != nil {
			errs = append(errs, err)
		} else if tree == nil {
			errs = append(errs, fmt.Errorf("parser failed to return orders"))
		}
		for _, err := range errs {
			report.WriteString(fmt.Sprintf("parser: %+v\n", err))
		}
		if tree != nil {
			report.WriteString(fmt.Sprintf("%v\n", *tree))
		}

		reportFile := fmt.Sprintf("sp%02d.t%d.report.txt", u.Species.No, s.data.Store.Turn)
		fullreportFile := filepath.Join(uploads, reportFile)

		log.Printf("server: %s %q: species %s turn %d report %s\n", r.Method, r.URL.Path, u.SpeciesId, turnNumber, reportFile)
		if err := ioutil.WriteFile(fullreportFile, report.Bytes(), 0644); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// add a guard here for the race condition
		updated := false
		for _, f := range s.data.Files[u.SpeciesId] {
			if f.Turn == turnNumber {
				f.Report = reportFile
				f.Orders = ordersFile
				f.Date = date
				updated = true
				break
			}
		}
		if !updated {
			s.data.Files[u.SpeciesId] = append(s.data.Files[u.SpeciesId], &FileData{
				SpeciesId: u.SpeciesId,
				SpeciesNo: u.Species.No,
				Turn:      turnNumber,
				Report:    reportFile,
				Orders:    ordersFile,
				Date:      date,
			})
		}

		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func (s *Server) handleHomePage(uploads string) http.HandlerFunc {
	log.Printf("handling home page with %q\n", uploads)
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/home" {
			http.NotFound(w, r)
			return
		}

		// Initialize a slice containing the paths to the two files. Note that the
		// home.page.tmpl file must be the *first* file in the slice.
		files := []string{
			filepath.Join("templates", "home.page.html"),
			filepath.Join("templates", "base.layout.html"),
		}

		// Use the template.ParseFiles() function to read the files and store the
		// templates in a template set. Notice that we can pass the slice of file paths
		// as a variadic parameter?
		ts, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	}
}
