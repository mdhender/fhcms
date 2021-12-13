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
	"path"
	"strings"
)

// shiftPath splits the given path into the first segment (head) and the rest (tail).
// for example, "/foo/bar/baz" gives head = "foo", tail = "/bar/baz".
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	if i := strings.Index(p[1:], "/") + 1; i <= 0 {
		return p[1:], "/"
	} else {
		return p[1:i], p[i:]
	}
}
