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
	"io/ioutil"
)

// Galaxy is the galaxy meta-data
type Galaxy struct {
	/* Design number of species in galaxy. */
	DNumSpecies int `json:"d_num_species"`
	/* Actual number of species allocated. */
	NumSpecies int `json:"num_species"`
	/* Galactic radius in parsecs. */
	Radius int `json:"radius"`
	/* Current turn number. */
	TurnNumber int `json:"turn_number"`
}

// ReadGalaxy returns either an initialized Galaxy or an error.
func ReadGalaxy(name string, bo binary.ByteOrder) (*Galaxy, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	var gd galaxy_data

	if err := binary.Read(r, bo, &gd); err != nil {
		return nil, err
	}

	var g Galaxy
	g.DNumSpecies = int(gd.DNumSpecies)
	g.NumSpecies = int(gd.NumSpecies)
	g.Radius = int(gd.Radius)
	g.TurnNumber = int(gd.TurnNumber)

	return &g, nil
}
