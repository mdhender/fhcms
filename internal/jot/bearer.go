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
	"strings"
)

// FromBearerToken returns the JWT from the Authorization header.
// If there is no bearer token or if the token is invalid for any reason, it returns nil.
func FromBearerToken(r *http.Request) *JWT {
	log.Printf("jot: bearer: entered\n")
	if headerAuthText := r.Header.Get("Authorization"); headerAuthText != "" {
		log.Printf("jot: bearer: found authorization header\n")
		if authTokens := strings.SplitN(headerAuthText, " ", 2); len(authTokens) == 2 {
			log.Printf("jot: bearer: found authorization token\n")
			if authType, authToken := authTokens[0], strings.TrimSpace(authTokens[1]); authType == "Bearer" {
				log.Printf("jot: bearer: found bearer token\n")
				if j, err := TokenToJWT(authToken); err == nil {
					log.Printf("jot: bearer: returning bearer token\n")
					return j
				} else {
					log.Printf("jot: bearer: token: %+v\n", err)
				}
			}
		}
	}
	return nil
}
