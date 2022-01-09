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

package cdb

import (
	"context"
	"github.com/mdhender/fhcms/internal/models"
	"log"
)

// FetchAccount implements the reactor.AuthStore interface
func (db *DB) FetchAccount(id int) (models.Account, bool) {
	log.Printf("[cdb] fetchAccount: id %d\n", id)
	var account struct {
		UserName string
		Admin    bool
		Active   bool
	}
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return models.Account{}, false
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), "select username, admin, active from accounts where id = $1", id).Scan(&account.UserName, &account.Admin, &account.Active)
	if err != nil {
		return models.Account{}, false
	}
	return models.Account{
		Id:       id,
		UserName: account.UserName,
		IsActive: account.Active,
		IsAdmin:  account.Active && account.Admin,
	}, true
}

// FetchProfile implements the reactor.ProfileStore interface
func (db *DB) FetchProfile(id int) (models.Profile, bool) {
	log.Printf("[cdb] fetchProfile: id %d\n", id)
	var profile struct {
		Id       int
		UserName string
		IsAdmin  bool
	}
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return models.Profile{}, false
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), "select id, username, admin from accounts where id = $1", id).Scan(&profile.Id, &profile.UserName, &profile.IsAdmin)
	if err != nil {
		return models.Profile{}, false
	}
	return models.Profile{
		Id:       id,
		UserName: profile.UserName,
		IsAdmin:  profile.IsAdmin,
	}, true
}
