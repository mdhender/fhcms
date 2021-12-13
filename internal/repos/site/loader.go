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

package site

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// loader reads a JSON file containing site data.
func loader(filename string) (*Store, error) {
	var input struct {
		Title string `json:"title"`
	}

	log.Printf("repo: site: site.Load: %q\n", filename)
	if data, err := ioutil.ReadFile(filename); err != nil {
		log.Printf("repo: site: site.Load: %q %+v\n", filename, err)
		return nil, err
	} else if err = json.Unmarshal(data, &input); err != nil {
		log.Printf("repo: site: site.Load: %q %+v\n", filename, err)
		return nil, err
	}
	return &Store{title: input.Title}, nil
}
