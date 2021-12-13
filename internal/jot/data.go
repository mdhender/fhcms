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
	"encoding/json"
	"strings"
)

// TokenToJWT expects the token to look like header.payload.signature if it is a valid JWT
func TokenToJWT(token string) (*JWT, error) {
	// token should look like header.payload.signature if it is a valid JWT
	sections := strings.Split(token, ".")
	if len(sections) != 3 || len(sections[0]) == 0 || len(sections[1]) == 0 || len(sections[2]) == 0 {
		return nil, ErrNotJWT
	}

	var j JWT
	j.h.b64 = sections[0]
	j.p.b64 = sections[1]
	j.s = sections[2]

	// the header is base64 encoded JSON
	if rawHeader, err := decode(j.h.b64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(rawHeader, &j.h); err != nil {
		return nil, err
	}

	// the payload is base64 encoded JSON
	if rawPayload, err := decode(j.p.b64); err != nil {
		return nil, err
	} else if err = json.Unmarshal(rawPayload, &j.p); err != nil {
		return nil, err
	}

	return &j, nil
}
