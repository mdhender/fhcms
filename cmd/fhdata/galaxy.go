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

package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type GalaxyData struct {
	DNumSpecies int `json:"d_num_species"` /* Design number of species in galaxy. */
	NumSpecies  int `json:"num_species"`   /* Actual number of species allocated. */
	Radius      int `json:"radius"`        /* Galactic radius in parsecs. */
	TurnNumber  int `json:"turn_number"`   /* Current turn number. */
}

type galaxy_data struct {
	DNumSpecies int32 /* Design number of species in galaxy. */
	NumSpecies  int32 /* Actual number of species allocated. */
	Radius      int32 /* Galactic radius in parsecs. */
	TurnNumber  int32 /* Current turn number. */
}

func ReadGalaxy(name string, bo binary.ByteOrder) (*GalaxyData, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	fmt.Printf("galaxy: binary file %q contains %d bytes\n", name, len(b))
	r := bytes.NewReader(b)

	var gd galaxy_data

	if err := binary.Read(r, bo, &gd); err != nil {
		return nil, err
	}

	var g GalaxyData
	g.DNumSpecies = int(gd.DNumSpecies)
	g.NumSpecies = int(gd.NumSpecies)
	g.Radius = int(gd.Radius)
	g.TurnNumber = int(gd.TurnNumber)
	fmt.Printf("galaxy: we think that this is turn %d\n", g.TurnNumber)

	return &g, nil
}

func (g *GalaxyData) Write(name string) error {
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, data, 0644)
}
