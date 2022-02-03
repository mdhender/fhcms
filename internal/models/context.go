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

package models

import (
	"context"
	"log"
	"net/http"
)

// accountKeyType is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type accountKeyType int

// accountKey is the key for Account values in Contexts.
// It is unexported; clients must use account.NewContext and account.FromContext
// to create and fetch from the context.
var accountKey accountKeyType

// NewContext returns a new Context that carries value a.
func (a Account) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, accountKey, a)
}

// AccountFromContext returns the Session value stored in ctx, if any.
func AccountFromContext(ctx context.Context) (Account, bool) {
	a, ok := ctx.Value(accountKey).(Account)
	return a, ok
}

func CurrentUser(r *http.Request) Account {
	log.Printf("[models] currentUser\n")
	a, ok := AccountFromContext(r.Context())
	if !ok {
		return Account{}
	}
	log.Printf("[models] currentUser: id %d name %q admin %v\n", a.Id, a.UserName, a.IsAdmin)
	return a
}
