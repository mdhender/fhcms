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

import "context"

// key is an unexported type for keys defined in this package.
// This prevents collisions with keys defined in other packages.
type key int

// jwtKey is the key for jot.JWT values in Cookies and Contexts.
// It is unexported; clients must use jot.NewContext and jot.FromContext
// to create and fetch from the context.
var jwtKey key

// NewContext returns a new Context that carries value j.
func (j *JWT) NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, jwtKey, j)
}

// FromContext returns the JWT value stored in ctx, if any.
func FromContext(ctx context.Context) (*JWT, bool) {
	j, ok := ctx.Value(jwtKey).(*JWT)
	return j, ok
}
