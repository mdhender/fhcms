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
	"crypto/sha1"
	"encoding/hex"
	"path"
	"strings"
)

func mkkey(salt, secret string) string {
	h, h2 := sha1.New(), sha1.New()
	_, _ = h.Write([]byte(salt))
	_, _ = h.Write([]byte(secret))
	_, _ = h2.Write(h.Sum(nil))
	_, _ = h2.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

// shiftPath splits off the first component of p, which will be cleaned of relative components before processing.
// head will never contain a slash and tail will always be a rooted path without trailing slash.
// from https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
