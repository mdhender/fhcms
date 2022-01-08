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
	"fmt"
	"github.com/mdhender/fhcms/internal/jot"
	"net/http"
)

func (s *Server) handleGetLogin(w http.ResponseWriter, r *http.Request) {
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
