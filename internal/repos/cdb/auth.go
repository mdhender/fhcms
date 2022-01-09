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

// Authenticate implements the reactor.AuthStore interface
func (db *DB) Authenticate(username, password string) (models.Account, bool) {
	log.Printf("[cdb] authenticate username %q password %q\n", username, password)
	var account struct {
		Id       int
		UserName string
		Admin    bool
		Active   bool
	}
	ctx := context.Background()
	conn, err := db.pool.Acquire(ctx)
	if err != nil {
		return models.Account{}, false
	}
	defer conn.Release()
	err = conn.QueryRow(ctx, "select id, username, admin, active from accounts where username = $1 and password = $2", username, password).Scan(&account.Id, &account.UserName, &account.Admin, &account.Active)
	if err != nil {
		return models.Account{}, false
	}
	return models.Account{
		Id:              account.Id,
		UserName:        username,
		IsActive:        account.Active,
		IsAdmin:         account.Active && account.Admin,
		IsAuthenticated: true,
	}, true
}
