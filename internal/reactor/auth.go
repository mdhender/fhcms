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
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/models"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

type AuthStore interface {
	Authenticate(username, password string) (models.Account, bool)
	FetchAccount(uid int) (models.Account, bool)
}

// we must serve only the login page to unauthenticated visitors
func (s *Server) authOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var a models.Account
		if j := jot.FromRequest(r); s.jf.Validate(j) == nil {
			if id, ok := j.UserID(); ok {
				if x, ok := s.auth.FetchAccount(id); ok {
					a.Id, a.UserName, a.IsActive, a.IsAuthenticated = x.Id, x.UserName, x.IsActive, true
					// the weird triple check for admin below prevents unauthenticated or inactive users from slipping through
					a.IsAdmin = a.IsActive && a.IsAuthenticated && x.IsAdmin
				}
			}
		}
		if !a.IsAuthenticated {
			log.Printf("[reactor] mw: authOnly: %s %s: !authenticated\n", r.Method, r.URL.Path)
			s.handleGetLogin(w, r)
			return
		}
		h(w, r.WithContext(a.NewContext(r.Context())))
	}
}

func (s *Server) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("[reactor] %s %q: handlePostLogin\n", r.Method, r.URL.Path)
	if err := r.ParseForm(); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	log.Printf("[reactor] %s %q: %v\n", r.Method, r.URL.Path, r.PostForm)
	var input struct {
		username string
		password string
	}
	for k, v := range r.Form {
		switch k {
		case "username":
			if len(v) != 1 || !utf8.ValidString(v[0]) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			input.username = v[0]
		case "password":
			if len(v) != 1 || !utf8.ValidString(v[0]) {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			input.password = v[0]
		}
	}

	//// hash the password to prevent simple timing attacks
	//var sh []byte
	//for _, b := range sha256.Sum256([]byte(input.password)) {
	//	sh = append(sh, b)
	//}
	//hashedPassword := hex.EncodeToString(sh)

	log.Printf("[reactor] %s %q: handlePostLogin: username %q password %q\n", r.Method, r.URL.Path, input.username, input.password)
	a, ok := s.auth.Authenticate(input.username, input.password)
	if !ok {
		log.Printf("[reactor] %s %q: handlePostLogin: account %q not found\n", r.Method, r.URL.Path, input.username)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	log.Printf("[reactor] %s %q: handlePostLogin: found username %q\n", r.Method, r.URL.Path, a.UserName)
	// the weird triple check for admin below prevents unauthenticated or inactive users from slipping through
	j, err := s.jf.NewToken(time.Hour*24*7, a.Id, a.UserName, a.IsAuthenticated && a.IsActive && a.IsAdmin)
	if err != nil {
		log.Printf("[reactor] %s %q: handlePostLogin: token %+v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	j.SetCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
