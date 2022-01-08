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

func (db *DB) FetchGalaxies(uid int) (models.Galaxies, bool) {
	var galaxies models.Galaxies
	conn, err := db.pool.Acquire(context.Background())
	if err != nil {
		return nil, false
	}
	defer conn.Release()
	galaxiesRows, err := conn.Query(context.Background(), "select id, name from galaxies where id in (select galaxy_id from species where account_id = $1)", uid)
	if err != nil {
		log.Printf("[cdb] fetchGalaxies %d: %+v\n", uid, err)
		return galaxies, false
	}
	for galaxiesRows.Next() {
		var galaxiesRowData struct {
			Id   int
			Name string
		}
		if err = galaxiesRows.Scan(&galaxiesRowData.Id, &galaxiesRowData.Name); err != nil {
			log.Printf("[cdb] fetchGalaxies %d: %+v\n", uid, err)
			return nil, false
		}
		c2, err := db.pool.Acquire(context.Background())
		if err != nil {
			return nil, false
		}
		defer c2.Release()
		speciesRows, err := c2.Query(context.Background(), "select id, name from species where account_id = $1 and galaxy_id = $2", uid, galaxiesRowData.Id)
		if err != nil {
			log.Printf("[cdb] fetchGalaxies %d: %+v\n", uid, err)
			return nil, false
		}
		for speciesRows.Next() {
			var speciesRow struct {
				Id   int
				Name string
			}
			if err = speciesRows.Scan(&speciesRow.Id, &speciesRow.Name); err != nil {
				log.Printf("[cdb] fetchGalaxies %d: %+v\n", uid, err)
				return nil, false
			}
			mg := &models.Galaxy{
				Id:   galaxiesRowData.Id,
				Name: galaxiesRowData.Name,
			}
			mg.Specie.Id = speciesRow.Id
			mg.Specie.Name = speciesRow.Name
			mg.Specie.Government.Name = "(missing name)"
			galaxies = append(galaxies, mg)
		}
	}
	return galaxies, true
}
