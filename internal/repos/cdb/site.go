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

func (db *DB) FetchSite() (models.Site, bool) {
	log.Printf("[cdb] fetchSite\n")
	var site struct {
		Title string
	}
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return models.Site{}, false
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), "select title from site").Scan(&site.Title)
	if err != nil {
		log.Printf("[cdb] fetchSite: %+v\n", err)
		return models.Site{Title: "Far Horizons"}, true
	}
	return models.Site{Title: site.Title}, true
}
