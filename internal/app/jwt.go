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
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
)

// JWTAuthenticator is a default authentication middleware to enforce access from the
// Verifier middleware request context values. The Authenticator sends a 401 Unauthorized
// response for any unverified tokens and passes the good ones through. It's just fine
// until you decide to write something similar and customize your client response.
func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, _, err := jwtauth.FromContext(r.Context()); err != nil || token == nil || jwt.Validate(token) != nil {
			if contentType := r.Header.Get("Content-type"); contentType == "application/vnd.api+json" {
				type errorObject struct {
					Status string `json:"status"`
					Code   string `json:"code"`
				}
				var response struct {
					Errors []errorObject `json:"errors"`
				}
				response.Errors = append(response.Errors, errorObject{
					Status: fmt.Sprintf("%d", http.StatusUnauthorized),
					Code:   http.StatusText(http.StatusUnauthorized),
				})
				w.Header().Set("Content-Type", "application/vnd.api+json")
				_ = json.NewEncoder(w).Encode(response)
				return
			} else if contentType := r.Header.Get("Content-type"); contentType == "application/json" {
				response := struct {
					Status string `json:"status"`
					Code   string `json:"code"`
				}{
					Status: fmt.Sprintf("%d", http.StatusUnauthorized),
					Code:   http.StatusText(http.StatusUnauthorized),
				}
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(response)
				return
			}
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}
