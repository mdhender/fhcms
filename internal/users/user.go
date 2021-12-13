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

package users

import (
	"context"
	"log"
	"net/http"
)

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// userKey is the key for User values in Contexts.
// It is unexported; clients must use user.NewContext and user.FromContext
// to create and fetch from the context.
var userKey key

// NewContext returns a new Context that carries value u.
func (u *User) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// FromContext returns the Session value stored in ctx, if any.
func FromContext(ctx context.Context) (*User, bool) {
	u, ok := ctx.Value(userKey).(*User)
	return u, ok
}

type User struct {
	Id              string
	Name            string
	IsAdmin         bool
	IsAuthenticated bool
}

func Current(r *http.Request) *User {
	u, ok := FromContext(r.Context())
	if !ok {
		log.Printf("user: Current: !ok\n")
		return &User{}
	}
	log.Printf("user: Current: ok: id %q\n", u.Id)
	return &User{Id: u.Id, Name: u.Name, IsAdmin: u.IsAdmin, IsAuthenticated: u.IsAuthenticated}
}
