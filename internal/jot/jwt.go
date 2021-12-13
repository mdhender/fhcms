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
	"net/http"
	"time"
)

// NewToken will return an unsigned JWT.
// userId is the user id to assign to the JWT.
// You must use a Factory to sign the token.
func NewToken(ttl time.Duration, userId string) (*JWT, error) {
	var j JWT
	j.h.TokenType = "JWT"
	j.p.IssuedAt = time.Now().Unix()
	j.p.ExpirationTime = time.Now().Add(ttl).Unix()
	j.p.UserID = userId
	return &j, nil
}

// JWT implements my version of the JSON Web Token.
type JWT struct {
	h struct {
		Algorithm   string `json:"alg"` // message authentication code algorithm
		TokenType   string `json:"typ"` // should always be JWT
		KeyID       string `json:"kid"` // identifier used to sign
		ContentType string `json:"cty,omitempty"`
		b64         string // header marshalled to JSON and then base-64 encoded
	}
	p struct {
		// The principal that issued the JWT.
		Issuer string `json:"iss,omitempty"`
		// The subject of the JWT.
		Subject string `json:"sub,omitempty"`
		// The recipients that the JWT is intended for.
		// Each principal intended to process the JWT must identify itself with a value in the audience claim.
		// If the principal processing the claim does not identify itself with a value in the aud claim when this claim is present,
		// then the JWT must be rejected.
		Audience []string `json:"aud,omitempty"`
		// The expiration time on and after which the JWT must not be accepted for processing.
		// The value must be a NumericDate:[9] either an integer or decimal, representing seconds past 1970-01-01 00:00:00Z.
		ExpirationTime int64 `json:"exp,omitempty"`
		// The time on which the JWT will start to be accepted for processing.
		// The value must be a NumericDate.
		NotBefore int64 `json:"nbf,omitempty"`
		// The time at which the JWT was issued.
		// The value must be a NumericDate.
		IssuedAt int64 `json:"iat,omitempty"`
		// Case sensitive unique identifier of the token even among different issuers.
		JWTID string `json:"jti,omitempty"`
		// Private data for use by the application.
		UserID string `json:"user_id,omitempty"` // this is the user id
		b64    string // payload marshalled to JSON and then base-64 encoded
	}
	s        string // signature base-64 encoded
	isSigned bool   // true only if the signature has been verified
}

func (j *JWT) IsValid() bool {
	if j == nil {
		//log.Printf("jwt is nil\n")
		return false
	} else if !j.isSigned {
		//log.Printf("alg %q typ %q signed %v borked\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	} else if j.p.IssuedAt == 0 {
		//log.Printf("alg %q typ %q signed %v no issue timestamp\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	} else if j.p.ExpirationTime == 0 {
		//log.Printf("alg %q typ %q signed %v no expiration timestamp\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	}
	now := time.Now().UTC()
	if !now.After(time.Unix(j.p.IssuedAt, 0)) {
		//log.Printf("alg %q typ %q signed %v !now.After(issuedAt) %s %s\n", j.h.Algorithm, j.h.TokenType, j.isSigned, now.Format("2006-01-02T15:04:05.99999999Z"), time.Unix(j.p.IssuedAt, 0).Format("2006-01-02T15:04:05.99999999Z"))
		return false
	} else if !time.Unix(j.p.ExpirationTime, 0).After(now) {
		//log.Printf("alg %q typ %q signed %v !expiresAt.After(now)\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	} else if j.p.NotBefore != 0 && !now.Before(time.Unix(j.p.NotBefore, 0)) {
		//log.Printf("alg %q typ %q signed %v !now.Before(notBefore)\n", j.h.Algorithm, j.h.TokenType, j.isSigned)
		return false
	}
	return true
}

func (j *JWT) DeleteCookie(w http.ResponseWriter) {
	DeleteCookie(w)
}

func (j *JWT) SetCookie(w http.ResponseWriter) {
	SetCookie(w, j)
}

// String implements the Stringer interface.
// Please don't call this before signing the token.
func (j *JWT) String() string {
	return j.h.b64 + "." + j.p.b64 + "." + j.s
}

// UserID returns the User ID from the JWT.
// We return an error if the JWT is not valid.
func (j *JWT) UserID() string {
	if j == nil {
		return ""
	} else if !j.IsValid() {
		return ""
	} else {
		return j.p.UserID
	}
}
