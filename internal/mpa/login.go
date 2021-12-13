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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/mdhender/fhcms/internal/jot"
	"github.com/mdhender/fhcms/internal/repos/accounts"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

func (s *Server) handleGetLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	jot.DeleteCookie(w)
	page := fmt.Sprintf(`<body>
				<h1>Far Horizons Login</h1>
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

func (s *Server) handlePostLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	log.Printf("server: %s %q: handlePostLogin\n", r.Method, r.URL.Path)
	if err := r.ParseForm(); err != nil {
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

	// hash the password to prevent simple timing attacks
	var sh []byte
	for _, b := range sha256.Sum256([]byte(input.password)) {
		sh = append(sh, b)
	}
	hashedPassword := hex.EncodeToString(sh)
	var a *accounts.Account
	log.Printf("server: %s %q: handlePostLogin: username %q password %q hashed %q\n", r.Method, r.URL.Path, input.username, input.password, hashedPassword)
	acct, ok := s.accts.ByUser[input.username]
	if !ok {
		log.Printf("server: %s %q: handlePostLogin: account %q not found\n", r.Method, r.URL.Path, input.username)
	}
	if acct == nil {
		log.Printf("server: %s %q: handlePostLogin: account %q is nil\n", r.Method, r.URL.Path, input.username)
	} else {
		if acct.Password != input.password {
			log.Printf("server: %s %q: handlePostLogin: account %q pass %q != input %q\n", r.Method, r.URL.Path, acct.Password, input.password)
		}
		if acct.HashedPassword != hashedPassword {
			log.Printf("server: %s %q: handlePostLogin: account %q pass %q != input %q\n", r.Method, r.URL.Path, acct.HashedPassword, hashedPassword)
		}
	}
	if ok && acct != nil && acct.HashedPassword == hashedPassword {
		if acct.Password == input.password {
			log.Printf("matched!\n")
			a = acct
		}
	}
	if a == nil {
		log.Printf("not matched!\n")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
	j, err := s.jf.NewToken(time.Hour*24*7, acct.Id)
	if err != nil {
		log.Printf("server: %s %q: handlePostLogin: token %+v\n", r.Method, r.URL.Path, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	j.SetCookie(w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
