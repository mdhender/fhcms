/*****************************************************************************
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
****************************************************************************/

package jsondb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Read(filename string) (*Store, error) {
	var ds Store
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	} else if err = json.Unmarshal(b, &ds); err != nil {
		return nil, err
	}
	return &ds, nil
}

func (ds *Store) Write(filename string) error {
	b, err := json.MarshalIndent(ds, "", "  ")
	if err != nil {
		return err
	}
	if filename == "*stdout*" {
		fmt.Println(string(b))
		return nil
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}
	return nil
}
