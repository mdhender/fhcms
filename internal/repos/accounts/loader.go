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
	"log"
)

func Load(filename string) (*AccountList, error) {
	var input struct {
		Accounts map[string]*Account `json:"accounts"`
	}
	input.Accounts = make(map[string]*Account)
	if data, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	accts := &AccountList{
		ById:   make(map[string]*Account),
		ByUser: make(map[string]*Account),
	}
	for id, acct := range input.Accounts {
		acct.Id = id
		var h []byte
		for _, b := range sha256.Sum256([]byte(acct.Password)) {
			h = append(h, b)
		}
		// hash the password to prevent simple timing attacks
		acct.HashedPassword = hex.EncodeToString(h)
		log.Printf("accounts: id %q acct %v\n", id, *acct)
		accts.ById[acct.Id] = acct
		accts.ByUser[acct.Username] = acct
	}
	return accts, nil
}
