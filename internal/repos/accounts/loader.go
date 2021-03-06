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

package accounts

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
)

func Load(filename string) (*Repository, error) {
	r := &Repository{}
	input := make(map[string]*Account)
	if data, err := ioutil.ReadFile(filename); err != nil {
		return r, err
	} else if err = json.Unmarshal(data, &input); err != nil {
		return r, err
	}
	for id, acct := range input {
		acct.Id = id
		if acct.HashedPassword == "" {
			acct.HashedPassword = hashPassword(acct.Salt, acct.Password)
		}
		acct.Password = ""
		r.data = append(r.data, acct)
	}
	return r, nil
}

// hashPassword returns the SHA-256 hash of the plaintext plus salt.
func hashPassword(salt, plaintext string) string {
	var h []byte
	for _, b := range sha256.Sum256([]byte(salt + plaintext)) {
		h = append(h, b)
	}
	return hex.EncodeToString(h)
}
