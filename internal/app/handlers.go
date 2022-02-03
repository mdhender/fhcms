/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/jsonapi"
	"io"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

func contact(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "contact\n")
}

func (s *Server) homeHandler(render func(p homePayload) ([]byte, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		u := s.currentUser(r)
		log.Printf("[app] homeHandler: u.id %q\n", u.Id)
		var p homePayload
		p.Account = s.currentUser(r)
		b, err := render(p)
		if err != nil {
			log.Printf("[app] homeHandler: render %+v\n", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(b)
		//_, _ = fmt.Fprintf(w, `<h1>HOME</h1>`)
		//_, _ = fmt.Fprintf(w, `<p>%+v</p>`, s.currentUser(r))
		//_, _ = fmt.Fprintf(w, `<p>todo: implement render</p>`)
	}
}

func (s *Server) loginGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	jot.DeleteCookie(w)
	page := fmt.Sprintf(`<body>
				<h1>Far Horizons Reactor</h1>
				<form action="/login"" method="post">
					<table>
						<tr><td align="right">Username&nbsp;</td><td><input type="text" name="username"></td></tr>
						<tr><td align="right">Password&nbsp;</td><td><input type="password" name="password"></td></tr>
						<tr><td>&nbsp;</td><td align="right"><input type="submit" value="Login"></td></tr>
					</table>
				</form>
			</body>`)
	_, _ = w.Write([]byte(page))
}

func (s *Server) loginPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username,omitempty"`
		Password string `json:"password,omitempty"`
	}

	contentType := r.Header.Get("Content-type")
	switch contentType {
	case "application/json":
		// enforce a maximum read of 1kb from the response body
		r.Body = http.MaxBytesReader(w, r.Body, 1024)
		// create a json decoder that will accept only our specific fields
		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&input); err != nil {
			jsonapi.Error(w, "error", http.StatusText(http.StatusBadRequest))
			return
		}
		// call decode again to confirm that the request contained only a single JSON object
		if err := dec.Decode(&struct{}{}); err != io.EOF {
			jsonapi.Error(w, "error", "request body must only contain a single JSON object")
			return
		}
		log.Printf("[app] json: %s %q: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			log.Printf("[app] form: %s %q: %+v\n", r.Method, r.URL.Path, err)
			http.SetCookie(w, &http.Cookie{Name: "jwt", Path: "/", MaxAge: -1, HttpOnly: true})
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Printf("[app] form: %s %q: %v\n", r.Method, r.URL.Path, r.PostForm)
		for k, v := range r.Form {
			switch k {
			case "username":
				if len(v) != 1 || !utf8.ValidString(v[0]) {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.Username = v[0]
			case "password":
				if len(v) != 1 || !utf8.ValidString(v[0]) {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.Password = v[0]
			}
		}
		log.Printf("[app] form: %s %q: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
	case "text/html":
		if err := r.ParseForm(); err != nil {
			log.Printf("[app] html: %s %q: %+v\n", r.Method, r.URL.Path, err)
			http.SetCookie(w, &http.Cookie{Name: "jwt", Path: "/", MaxAge: -1, HttpOnly: true})
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		log.Printf("[app] html: %s %q: %v\n", r.Method, r.URL.Path, r.PostForm)
		for k, v := range r.Form {
			switch k {
			case "username":
				if len(v) != 1 || !utf8.ValidString(v[0]) {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.Username = v[0]
			case "password":
				if len(v) != 1 || !utf8.ValidString(v[0]) {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
					return
				}
				input.Password = v[0]
			}
		}
		log.Printf("[app] html: %s %q: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
	default:
		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
	}

	acct, ok := s.accountStore.Authenticate(input.Username, input.Password)
	if !ok || !acct.IsAuthenticated {
		switch contentType {
		case "application/json":
			jsonapi.Error(w, "error", http.StatusText(http.StatusBadRequest))
		case "application/x-www-form-urlencoded", "text/html":
			w.Header().Set("Content-Type", "text/html")
			jot.DeleteCookie(w)
			page := fmt.Sprintf(`<body>
				<h1>Far Horizons Reactor</h1>
				<p>Invalid credentials. Please try again.</p>
				<form action="/login"" method="post">
					<table>
						<tr><td align="right">Username&nbsp;</td><td><input type="text" name="username"></td></tr>
						<tr><td align="right">Password&nbsp;</td><td><input type="password" name="password"></td></tr>
						<tr><td>&nbsp;</td><td align="right"><input type="submit" value="Login"></td></tr>
					</table>
				</form>
			</body>`)
			_, _ = w.Write([]byte(page))
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	switch contentType {
	case "application/json":
		log.Printf("[app] json: %s %q: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		var response struct {
			Links struct {
				Self string `json:"self"`
			} `json:"links"`
			Data struct {
				Token string `json:"token"`
			} `json:"data,omitempty"`
		}
		response.Links.Self = r.URL.Path
		response.Data.Token = "value"
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	case "application/x-www-form-urlencoded":
		log.Printf("[app] form: %s %q: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		j, err := s.jf.Token(time.Hour*24*7, acct.Id)
		if err != nil {
			log.Printf("[app] form: %s %q: loginPostHandler: token %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		j.SetCookie(w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	case "text/html":
		log.Printf("[app] html: %s %q: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
		j, err := s.jf.Token(time.Hour*24*7, acct.Id)
		if err != nil {
			log.Printf("[app] html: %s %q: loginPostHandler: token %+v\n", r.Method, r.URL.Path, err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		j.SetCookie(w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	default:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (s *Server) logoutHandler(w http.ResponseWriter, r *http.Request) {
	jot.DeleteCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) manifestJsonV3(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"manifest_version":3,"name":"Far Horizons","version":"7.5.2"}`))
}

func (s *Server) notImplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
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

func GetMessages(w http.ResponseWriter, r *http.Request) {
	if token := r.Header.Get("Token"); token == "" {
		log.Printf("[app] %s %q: token missing from header\n", r.Method, r.URL.Path)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	type Attributes struct {
		From    string `json:"from,omitempty"`
		Dttm    string `json:"dttm,omitempty"`
		Subject string `json:"subject,omitempty"`
		Body    string `json:"body,omitempty"`
	}
	type Message struct {
		Id         int        `json:"id,omitempty"`
		Type       string     `json:"type"`
		Attributes Attributes `json:"attributes"`
	}
	var response struct {
		Links struct {
			Self string `json:"self"`
		} `json:"links"`
		Data []Message `json:"data"`
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60,
		HttpOnly: true,
	})
	response.Links.Self = r.URL.Path
	response.Data = append(response.Data, Message{Id: 1, Type: "messages", Attributes: Attributes{From: "Tomster", Dttm: "2020.07.24, 16:15:03", Subject: "Hey Zoey", Body: "How is it going? Will I see you at EmberConf next year?"}})
	response.Data = append(response.Data, Message{Id: 2, Type: "messages", Attributes: Attributes{From: "EmberConf", Dttm: "2020.07.21, 16:15:03", Subject: "Registration Confirmation for EmberConf 2021", Body: "Thanks so much for registering to join us at EmberConf! You do NOT need to print this confirmation, but photo ID may be required for entry."}})
	w.Header().Set("Content-Type", "application/vnd.api+json")
	_ = json.NewEncoder(w).Encode(response)
	log.Printf("[app] json: %s %q: success\n", r.Method, r.URL.Path)
}
