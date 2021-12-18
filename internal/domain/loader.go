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

package domain

import (
	"encoding/binary"
	"github.com/mdhender/fhcms/internal/dat32"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

type Specie struct {
	No   string
	Name string
}

func (s *Store) loadSpecie(files, spNo string) (*Specie, error) {
	bigEndian := viper.GetBool("files.big_endian")
	var bo binary.ByteOrder
	if bigEndian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}
	sp, err := dat32.ReadSpecies(filepath.Join(files, "sp"+spNo+".dat"), 0, bo)
	if err != nil {
		log.Printf("[domain] loadSpecie %q %q %+v\n", files, spNo, err)
		return &Specie{}, err
	}
	return &Specie{
		No:   spNo,
		Name: sp.Name,
	}, nil
}
