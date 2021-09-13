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

type star_file_t struct {
	NumStars int32
	StarBase []star_data
}

type star_data struct {
	X           int8 /* Coordinates. */
	Y           int8
	Z           int8
	Type        int8 /* Dwarf, degenerate, main sequence or giant. */
	Color       int8 /* Star color. Blue, blue-white, etc. */
	Size        int8 /* Star size, from 0 thru 9 inclusive. */
	NumPlanets  int8 /* Number of usable planets in star system. */
	HomeSystem  int8 /* TRUE if this is a good potential home system. */
	WormHere    int8 /* TRUE if wormhole entry/exit. */
	WormX       int8
	WormY       int8
	WormZ       int8
	Reserved1   int16     /* Reserved for future use. Zero for now. */
	Reserved2   int16     /* Reserved for future use. Zero for now. */
	PlanetIndex int16     /* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	Message     int32     /* Message associated with this star system, if any. */
	VisitedBy   [2]uint64 /* A bit is set if corresponding species has been here. */
	Reserved3   int32     /* Reserved for future use. Zero for now. */
	Reserved4   int32     /* Reserved for future use. Zero for now. */
	Reserved5   int32     /* Reserved for future use. Zero for now. */
	Padding     [2]uint8  // padding to make Go struct same size as C
}

type Stars struct {
	NumStars int    `json:"num_stars"`
	Stars    []Star `json:"stars"`
}

type Star struct {
	X           int   `json:"x"` /* Coordinates. */
	Y           int   `json:"y"`
	Z           int   `json:"z"`
	Type        int   `json:"type"`        /* Dwarf, degenerate, main sequence or giant. */
	Color       int   `json:"color"`       /* Star color. Blue, blue-white, etc. */
	Size        int   `json:"size"`        /* Star size, from 0 thru 9 inclusive. */
	NumPlanets  int   `json:"num_planets"` /* Number of usable planets in star system. */
	HomeSystem  bool  `json:"home_system"` /* TRUE if this is a good potential home system. */
	WormHere    bool  `json:"worm_here"`   /* TRUE if wormhole entry/exit. */
	WormX       int   `json:"worm_x"`
	WormY       int   `json:"worm_y"`
	WormZ       int   `json:"worm_z"`
	PlanetIndex int   `json:"planet_index"` /* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	Message     int   `json:"message"`      /* Message associated with this star system, if any. */
	VisitedBy   []int `json:"visited_by"`   /* A bit is set if corresponding species has been here. */
}

func ReadStars(name string, bo binary.ByteOrder) (*Stars, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	fmt.Printf("  star: binary file %q contains %d bytes\n", name, len(b))
	r := bytes.NewReader(b)

	var sd star_file_t
	if err := binary.Read(r, bo, &sd.NumStars); err != nil {
		return nil, err
	}
	fmt.Printf("  star: num stars %d\n", sd.NumStars)
	sd.StarBase = make([]star_data, sd.NumStars, sd.NumStars)
	numStars := int(sd.NumStars)
	for i := 0; i < numStars; i++ {
		if err := binary.Read(r, bo, &sd.StarBase[i]); err != nil {
			return nil, err
		}
		//if i < 3 {
		//	ec := json.NewEncoder(os.Stdout)
		//	ec.SetIndent("", "  ")
		//	if err := ec.Encode(sd.StarBase[i]); err != nil {
		//		return nil, err
		//	}
		//}
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
		stars.Stars[i].HomeSystem = sd.StarBase[i].HomeSystem != 0
		stars.Stars[i].WormHere = sd.StarBase[i].WormHere != 0
		stars.Stars[i].WormX = int(sd.StarBase[i].WormX)
		stars.Stars[i].WormY = int(sd.StarBase[i].WormY)
		stars.Stars[i].WormZ = int(sd.StarBase[i].WormZ)
		stars.Stars[i].PlanetIndex = int(sd.StarBase[i].PlanetIndex)
		stars.Stars[i].Message = int(sd.StarBase[i].Message)
		for sp := 1; sp <= MAX_SPECIES; sp++ {
			// skip if bit is not set for this species
			if !speciesBitIsSet(sd.StarBase[i].VisitedBy, sp) {
				continue
			}
			stars.Stars[i].VisitedBy = append(stars.Stars[i].VisitedBy, sp)
		}
	}

	return &stars, nil
}

func (s Stars) Write(name string) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, data, 0644)
}
