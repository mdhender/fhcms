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

package jot

import (
	"log"
	"net/http"
	"time"
)

func DeleteCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func SetCookie(w http.ResponseWriter, j *JWT) {
	var maxAge int
	if j.p.ExpirationTime != 0 {
		maxAge = int(time.Unix(j.p.ExpirationTime, 0).Sub(time.Now().UTC()).Seconds())
	}
	if maxAge < 15 {
		maxAge = 15
	} else if maxAge > 14*24*60*60 {
		maxAge = 14 * 24 * 60 * 60
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Path:     "/",
		Value:    j.String(),
		MaxAge:   maxAge,
		HttpOnly: true,
	})
}

func FromCookie(r *http.Request) *JWT {
	log.Printf("jot: cookie: entered\n")
	if c, err := r.Cookie("jwt"); err == nil {
		if j, err := TokenToJWT(c.Value); err == nil {
			return j
		} else {
			log.Printf("jot: cookie: token: %+v\n", err)
		}
	} else {
		log.Printf("jot: cookie: %+v\n", err)
	}
	return nil
}

// FromRequest will pull a JWT from a request header.
// It looks for a bearer token first.
// If it can't find one, it looks for a cookie.
func FromRequest(r *http.Request) *JWT {
	log.Printf("jot: request: entered\n")
	if j := FromBearerToken(r); j != nil {
		log.Printf("jot: bearer: found\n")
		return j
	}
	return FromCookie(r)
}
