/*******************************************************************************
fhapp - a player aid for Far Horizons
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

package main

import (
	"context"
	"net/http"
)

type UserData struct {
	Player          string       `json:"player"`
	SpeciesId       int          `json:"species_id"`
	Species         *SpeciesData `json:"-"`
	IsAdmin         bool         `json:"is_admin"`
	IsAuthenticated bool         `json:"is_authenticated"`
}

// userDataContextType is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type userDataContextType int

// userKey is the key for user.UserData values in Contexts.
// It is unexported; clients use user.NewContext and user.FromContext
// instead of using this key directly.
var userDataContextKey userDataContextType

// NewContext returns a new Context that carries value u.
func (u UserData) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userDataContextKey, u)
}

// FromContext returns the User value stored in ctx, if any.
// If there is no value, then the "zero" value for UserData is returned.
func FromContext(ctx context.Context) (UserData, bool) {
	if u, ok := ctx.Value(userDataContextKey).(UserData); ok {
		return u, ok
	}
	return UserData{}, false
}

func currentUser(r *http.Request) UserData {
	u, _ := FromContext(r.Context())
	return u
}
