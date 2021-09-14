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

// Stars is the stars meta-data
type Stars struct {
	NumStars int    `json:"num_stars"`
	Stars    []Star `json:"stars"`
}

// Star is the star data
type Star struct {
	/* Coordinates. */
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
	/* Dwarf, degenerate, main sequence or giant. */
	Type int `json:"type"`
	/* Star color. Blue, blue-white, etc. */
	Color int `json:"color"`
	/* Star size, from 0 thru 9 inclusive. */
	Size int `json:"size"`
	/* Number of usable planets in star system. */
	NumPlanets int `json:"num_planets"`
	/* TRUE if this is a good potential home system. */
	HomeSystem int `json:"home_system"`
	/* TRUE if wormhole entry/exit. */
	WormHere int `json:"worm_here"`
	/* Coordinates of exit point for wormhole. Valid only if WormHere is TRUE. */
	WormX int `json:"worm_x"`
	WormY int `json:"worm_y"`
	WormZ int `json:"worm_z"`
	/* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	PlanetIndex int `json:"planet_index"`
	/* Message associated with this star system, if any. */
	Message int `json:"message"`
	/* A bit is set if corresponding species has been here. */
	VisitedBy []int `json:"visited_by"`
}

// ReadStars returns either an initialized set of stars or an error.
func ReadStars(name string, bo binary.ByteOrder) (*Stars, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	var sd star_file_t
	if err := binary.Read(r, bo, &sd.NumStars); err != nil {
		return nil, err
	}

	sd.StarBase = make([]star_data, sd.NumStars, sd.NumStars)
	numStars := int(sd.NumStars)
	for i := 0; i < numStars; i++ {
		if err := binary.Read(r, bo, &sd.StarBase[i]); err != nil {
			return nil, err
		}
	}

	var stars Stars
	stars.NumStars = int(sd.NumStars)
	stars.Stars = make([]Star, stars.NumStars, stars.NumStars)
	for i := 0; i < stars.NumStars; i++ {
		stars.Stars[i].X = int(sd.StarBase[i].X)
		stars.Stars[i].Y = int(sd.StarBase[i].Y)
		stars.Stars[i].Z = int(sd.StarBase[i].Z)
		stars.Stars[i].Type = int(sd.StarBase[i].Type)
		stars.Stars[i].Color = int(sd.StarBase[i].Color)
		stars.Stars[i].Size = int(sd.StarBase[i].Size)
		stars.Stars[i].NumPlanets = int(sd.StarBase[i].NumPlanets)
		stars.Stars[i].HomeSystem = int(sd.StarBase[i].HomeSystem)
		stars.Stars[i].WormHere = int(sd.StarBase[i].WormHere)
		stars.Stars[i].WormX = int(sd.StarBase[i].WormX)
		stars.Stars[i].WormY = int(sd.StarBase[i].WormY)
		stars.Stars[i].WormZ = int(sd.StarBase[i].WormZ)
		stars.Stars[i].PlanetIndex = int(sd.StarBase[i].PlanetIndex)
		stars.Stars[i].Message = int(sd.StarBase[i].Message)
		for sp := 1; sp <= MAX_SPECIES; sp++ {
			if speciesBitIsSet(sd.StarBase[i].VisitedBy, sp) {
				stars.Stars[i].VisitedBy = append(stars.Stars[i].VisitedBy, sp)
			}
		}
	}

	return &stars, nil
}
