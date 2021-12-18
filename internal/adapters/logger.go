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

package adapters

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Logger returns an adapter that will log the request.
func Logger(next http.Handler) http.Handler {
	log.Printf("[logger] initializing\n")
	defer log.Printf("[logger] initialized\n")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start, l := time.Now(), &logger{ResponseWriter: w, status: http.StatusOK}
		defer func() {
			elapsed := time.Since(start)
			_, _ = fmt.Fprintf(os.Stderr, "%v %s %s \"%s %s %s\" %d \"%s\" %v\n",
				time.Now().UTC().Format("2006/01/02 15:04:05"),
				r.Host,
				r.RemoteAddr,
				r.Method,
				r.URL.Path,
				r.Proto,
				l.status,
				r.UserAgent(),
				elapsed,
			)
		}()
		next.ServeHTTP(l, r)
	})
}

// create a custom response writer so that we can save the status for logging.
type logger struct {
	http.ResponseWriter
	status int
}

// WriteHeader overrides the default ResponseWriter to save the status
// before calling the original ResponseWriter.WriteHeader.
func (l *logger) WriteHeader(status int) {
	l.status = status
	l.ResponseWriter.WriteHeader(status)
}
