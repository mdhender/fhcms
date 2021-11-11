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
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
)

// sizeof char         is      1
// sizeof long         is      4
// sizeof short        is      2
// sizeof galaxy_data  is     16
// sizeof star_data    is     52
// sizeof planet_data  is     40
// sizeof nampla_data  is    288
// sizeof species_data is    264
// sizeof ship_data    is    172
// sizeof trans_data   is    148

const MAX_ITEMS = 38
const MAX_SPECIES = 100
const NUM_CONTACT_WORDS = ((MAX_SPECIES - 1) / 32) + 1 /* Assume at least 32 bits per long word. */

func main() {
	if NUM_CONTACT_WORDS != 4 {
		panic(fmt.Sprintf("assert(NUM_CONTACT_WORDS == 4)"))
	}

	// byte order is the order in the data file, not the computer we're running on!
	bigendian := false

	err := run(".", "galaxy.dat", "stars.dat", "planets.dat", bigendian)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}
}

func run(root, gFile, sFile, pFile string, bigendian bool) error {
	var bo binary.ByteOrder
	if bigendian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}

	g, err := ReadGalaxy(filepath.Join(root, gFile), bo)
	if err != nil {
		return err
	}

	s, err := ReadStars(filepath.Join(root, sFile), bo)
	if err != nil {
		return err
	}

	p, err := ReadPlanets(filepath.Join(root, pFile), bo)
	if err != nil {
		return err
	}

	var species []*Species
	for n := 1; n <= g.NumSpecies; n++ {
		sp, err := ReadSpecies(filepath.Join(root, fmt.Sprintf("sp%02d.dat", n)), n, bo)
		if err != nil {
			return err
		}
		species = append(species, sp)
	}

	if err := g.Write("galaxy.json"); err != nil {
		return err
	}
	if err := s.Write("stars.json"); err != nil {
		return err
	}
	if err := p.Write("planets.json"); err != nil {
		return err
	}
	for _, sp := range species {
		if err := sp.Write(); err != nil {
			return err
		}
	}

	c := NewCluster(g, s, p, species)
	if err = c.Write("cluster.json"); err != nil {
		return err
	}

	return err
}
