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
	"encoding/json"
	"fmt"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/mdhender/fhcms/internal/dat32"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC) // force logs to be UTC

	// byte order is the order in the data file, not the computer we're running on!
	bigendian := false

	root := "."
	err := run(root, "galaxy.dat", "stars.dat", "planets.dat", bigendian)
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}

	os.Exit(0)
}

func run(root, gFile, sFile, pFile string, bigendian bool) error {
	var bo binary.ByteOrder
	if bigendian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}

	galaxy, err := dat32.ReadGalaxy(filepath.Join(root, gFile), bo)
	if err != nil {
		return err
	}

	stars, err := dat32.ReadStars(filepath.Join(root, sFile), bo)
	if err != nil {
		return err
	}

	planets, err := dat32.ReadPlanets(filepath.Join(root, pFile), bo)
	if err != nil {
		return err
	}

	var species []*dat32.Species
	for n := 1; n <= galaxy.NumSpecies; n++ {
		sp, err := dat32.ReadSpecies(filepath.Join(root, fmt.Sprintf("sp%02d.dat", n)), n, bo)
		if err != nil {
			return err
		}
		species = append(species, sp)
	}

	if err := write(filepath.Join(root, "galaxy.json"), galaxy); err != nil {
		return err
	}
	if err := write(filepath.Join(root, "stars.json"), stars); err != nil {
		return err
	}
	if err := write(filepath.Join(root, "planets.json"), planets); err != nil {
		return err
	}
	for _, sp := range species {
		if err := write(filepath.Join(root, fmt.Sprintf("sp%02d.json", sp.Id)), sp); err != nil {
			return err
		}
	}

	c, err := cluster.ConvertDat32ToCluster(galaxy, stars, planets, species)
	if err != nil {
		return err
	} else if err = write(filepath.Join(root, "cluster.json"), c); err != nil {
		return err
	}

	return err
}

func write(name string, data interface{}) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, b, 0644)
}
