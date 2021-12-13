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

import "errors"

var ErrBadFactory = errors.New("bad factory")
var ErrBadRequest = errors.New("bad request")
var ErrExpired = errors.New("expired")
var ErrInvalid = errors.New("invalid jot")
var ErrInvalidSignature = errors.New("invalid signature")
var ErrMissingAuthHeader = errors.New("missing auth header")
var ErrMissingSigner = errors.New("missing signer")
var ErrNotBearer = errors.New("not a bearer token")
var ErrNotJWT = errors.New("not a jwt")
var ErrNotMyKID = errors.New("not my kid")
var ErrUnauthorized = errors.New("unauthorized")
