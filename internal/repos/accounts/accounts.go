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
	"github.com/mdhender/fhcms/internal/models"
)

type Repository struct {
	data []*Account
}

type Account struct {
	Id             string          `json:"id"`
	Email          string          `json:"email"`
	UserName       string          `json:"username"`
	Password       string          `json:"password"`
	Salt           string          `json:"salt"` // salt for hashing the password
	HashedPassword string          `json:"hashed_password"`
	Roles          map[string]bool `json:"roles"` // roles are used for authorization
}

type AccountList struct {
	ById   map[string]*Account
	ByUser map[string]*Account
}

// Authenticate implements models.AccountStore
func (r *Repository) Authenticate(name, password string) (models.Account, bool) {
	for _, acct := range r.data {
		if acct.HashedPassword == hashPassword(acct.Salt, password) && acct.UserName == name {
			return models.Account{
				Id:              acct.Id,
				Email:           acct.Email,
				UserName:        acct.UserName,
				IsActive:        acct.Roles["active"],
				IsAdmin:         acct.Roles["admin"],
				IsAuthenticated: true,
			}, true
		}
	}
	return models.Account{}, false
}

// FetchById implements models.AccountStore
func (r *Repository) FetchById(id string) (models.Account, bool) {
	for _, acct := range r.data {
		if acct.Id == id {
			return models.Account{
				Id:              acct.Id,
				Email:           acct.Email,
				UserName:        acct.UserName,
				IsActive:        acct.Roles["active"],
				IsAdmin:         acct.Roles["admin"],
				IsAuthenticated: false,
			}, true
		}
	}
	return models.Account{}, false
}

// FetchByUserName implements models.AccountStore
func (r *Repository) FetchByUserName(name string) (models.Account, bool) {
	for _, acct := range r.data {
		if acct.UserName == name {
			return models.Account{
				Id:              acct.Id,
				Email:           acct.Email,
				UserName:        acct.UserName,
				IsActive:        acct.Roles["active"],
				IsAdmin:         acct.Roles["admin"],
				IsAuthenticated: false,
			}, true
		}
	}
	return models.Account{}, false
}
