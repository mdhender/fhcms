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

//// authMiddleware checks the request for a cookie or bearer token.
//// if it is missing, invalid, or expired, the request is either denied or redirected to the login handler.
//// there are some routes that treated specially.
////    /login   is passed to the next handler
////    /logout  clears the cookie and redirects to /
////    /api/... confirms the bearer token before passing to the next handler
//func (a *App) authMiddleware(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		route, rest := shiftPath(r.URL.Path)
//		switch route {
//		case "api":
//			// get the bearer token
//			// if missing, invalid, or expired, deny the request
//			r.URL.Path = rest // always strip /api from the request
//			next.ServeHTTP(w, r)
//		case "login":
//			next.ServeHTTP(w, r)
//
//		case "logout":
//			// delete our JWT cookie and redirect back to the home page
//			http.SetCookie(w, &http.Cookie{Name: "jwt", Path: "/", MaxAge: -1, HttpOnly: true})
//			http.Redirect(w, r, "/", http.StatusSeeOther)
//		default:
//			// get the bearer token
//			// if missing, get the cookie
//			if c, err := r.Cookie("jwt"); err == nil {
//
//			}
//			// if missing, invalid, or expired, delete the cookie
//			http.Redirect(w, r, "/login", http.StatusSeeOther)
//			next.ServeHTTP(w, r)
//		}
//	})
//}
//
//func (a *App) loginHandler(w http.ResponseWriter, r *http.Request) {
//	var input struct {
//		Username string `json:"username,omitempty"`
//		Password string `json:"password,omitempty"`
//	}
//
//	if contentType := r.Header.Get("Content-type"); contentType != "application/json" {
//		http.Error(w, http.StatusText(http.StatusUnsupportedMediaType), http.StatusUnsupportedMediaType)
//		return
//	}
//	log.Printf("app: %s %q\n", r.Method, r.URL.Path)
//	r.Body = http.MaxBytesReader(w, r.Body, 1024) // enforce a maximum read of 1kb from the response body
//	// create a json decoder that will accept only our specific fields
//	dec := json.NewDecoder(r.Body)
//	dec.DisallowUnknownFields()
//	if err := dec.Decode(&input); err != nil {
//		log.Printf("app: %s %q: decode %+v\n", r.Method, r.URL.Path, err)
//		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//		return
//	}
//	// call decode again to confirm that the request contained only a single JSON object
//	if err := dec.Decode(&struct{}{}); err != io.EOF {
//		log.Printf("app: %s %q: eof check %+v\n", r.Method, r.URL.Path, err)
//		http.Error(w, "request body must only contain a single JSON object", http.StatusBadRequest)
//		return
//	}
//	log.Printf("app: json: %s %q: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
//
//	for id, acct := range a.Accounts {
//		if acct.Username == input.Username && acct.Password == input.Password {
//			log.Printf("app: json: %s %q: success: username %q password %q\n", r.Method, r.URL.Path, input.Username, input.Password)
//			var response struct {
//				Token string `json:"token"`
//			}
//			iat := time.Now()
//			_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
//				"iat": iat.UTC().Unix(),
//				"exp": iat.Add(tokenTimeToLive).UTC().Unix(),
//				"sub": id,
//			})
//			if err != nil {
//				log.Printf("app: %s %q: gentok %+v\n", r.Method, r.URL.Path, err)
//				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
//				return
//			}
//			response.Token = tokenString
//			w.Header().Set("Content-Type", "application/json")
//			_ = json.NewEncoder(w).Encode(response)
//			return
//		}
//	}
//	http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
//}
