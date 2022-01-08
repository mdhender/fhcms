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
	"time"
)

type Factory struct {
	kid string
	s   Signer
}

// NewFactory returns an initialized factory.
// The signer is used to sign the generated tokens.
// Factories are cheap, so create a new one to rotate keys.
func NewFactory(kid string, s Signer) *Factory {
	return &Factory{kid: kid, s: s}
}

// ID returns the id of the current signer.
func (f *Factory) ID() string {
	return f.kid
}

// Sign will sign a JWT.
// It uses the current values in the header and payload, so it is
// safe to call multiple times.
func (f *Factory) Sign(j *JWT) error {
	if f == nil || f.kid == "" || f.s == nil {
		return ErrBadFactory
	} else if j == nil {
		return ErrInvalid
	}

	j.h.Algorithm = f.s.Algorithm()
	j.h.KeyID = f.kid

	// base64 encode JSON representation of header
	if h, err := json.Marshal(j.h); err != nil {
		return err
	} else {
		j.h.b64 = encode(h)
	}

	// base64 encode JSON representation of payload
	if p, err := json.Marshal(j.p); err != nil {
		return err
	} else {
		j.p.b64 = encode(p)
	}

	// base64 encode JSON representation of signature
	if rawSignature, err := f.s.Sign([]byte(j.h.b64 + "." + j.p.b64)); err != nil {
		return err
	} else {
		j.s = encode(rawSignature)
	}

	return nil
}

// NewToken is a helper to create a new, signed JWT.
// userId is the user id to add to the JWT
func (f *Factory) NewToken(ttl time.Duration, userId int, userName string, isAdmin bool) (*JWT, error) {
	if f == nil || f.kid == "" || f.s == nil {
		return nil, ErrBadFactory
	}

	j, err := NewToken(ttl, userId, userName, isAdmin)
	if err != nil {
		return nil, err
	} else if err = f.Sign(j); err != nil {
		return nil, err
	}

	return j, nil
}

// Validate will return an error if the JWT is not properly signed.
// It tries to update the isSigned to true only if the JWT is properly signed.
func (f *Factory) Validate(j *JWT) error {
	if j == nil {
		return ErrInvalid
	}
	j.isSigned = false
	if f == nil || f.kid == "" || f.s == nil {
		return ErrBadFactory
	} else if expectedSignature, err := f.s.Sign([]byte(j.h.b64 + "." + j.p.b64)); err != nil {
		return err
	} else if j.isSigned = j.s == encode(expectedSignature); !j.isSigned {
		return ErrUnauthorized
	}
	return nil // valid signature
}
