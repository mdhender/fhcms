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

package dat32

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
)

type SpLocData struct {
	S       int // Species number
	X, Y, Z int // Location
}

// ReadLocations returns either a slice of SpLocData or an error.
func ReadLocations(name string, bo binary.ByteOrder) ([]SpLocData, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)
	var ld []SpLocData
	var data sp_loc_data
	for {
		if err = binary.Read(r, bo, &data); err != nil {
			break
		}
		ld = append(ld, SpLocData{S: int(data.S), X: int(data.X), Y: int(data.Y), Z: int(data.Z)})
	}
	if err != io.EOF {
		return nil, err
	}
	return ld, nil
}
