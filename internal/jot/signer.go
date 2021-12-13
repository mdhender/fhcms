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
	"crypto/hmac"
	"crypto/sha256"
)

// Signer interface
type Signer interface {
	Algorithm() string
	Sign(msg []byte) ([]byte, error)
}

// HS256Signer implements a Signer using HMAC256.
type HS256Signer struct {
	key []byte
}

func NewHS256Signer(secret []byte) (*HS256Signer, error) {
	s := HS256Signer{key: make([]byte, len(secret))}
	copy(s.key, secret)
	return &s, nil
}

// Algorithm implements the Signer interface
func (s *HS256Signer) Algorithm() string {
	return "HS256"
}

// Sign implements the Signer interface
func (s *HS256Signer) Sign(msg []byte) ([]byte, error) {
	hm := hmac.New(sha256.New, s.key)
	if _, err := hm.Write(msg); err != nil {
		return nil, err
	}
	return hm.Sum(nil), nil
}
