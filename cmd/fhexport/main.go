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
	"encoding/json"
	"fmt"
	"github.com/mdhender/fhcms/store/jsondb"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC) // force logs to be UTC

	if err := run(); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(2)
	}

	os.Exit(0)
}

func run() error {
	jdb, err := jsondb.Read("D:\\FarHorizons\\testdata\\t17\\galaxy.json")
	if err != nil {
		return err
	}
	for _, sp := range jdb.Species {
		if sp.Id != 18 {
			continue
		}
		e, err := export(jdb, sp, visited)
		if err != nil {
			return fmt.Errorf("species %d: %w", sp, err)
		}
		if b, err := json.MarshalIndent(e, "", "  "); err != nil {
			return err
		} else if err = ioutil.WriteFile("D:\\FarHorizons\\testdata\\t17\\sp18.export.json", b, 0600); err != nil {
			return err
		}
	}
	return nil
}

func export(ds *jsondb.Store, sp *jsondb.SpeciesData, visited []Coords) (*Export, error) {
	e := Export{TurnNumber: ds.Galaxy.TurnNumber}
	for _, star := range ds.Stars {
		hasVisited := false
		for _, v := range visited {
			if star.X == v.X && star.Y == v.Y && star.Z == v.Z {
				hasVisited = true
				break
			}
		}

		system := &SystemData{Id: fmt.Sprintf("%d.%d.%d", star.X, star.Y, star.Z), Coords: Coords{X: star.X, Y: star.Y, Z: star.Z}}
		system.Name = fmt.Sprintf("%d %d %d", star.X, star.Y, star.Z)
		system.Color = jsondb.IntToStarColor(star.Color)
		system.Size = star.Size
		system.Type = jsondb.IntToStarType(star.Type)
		e.Systems = append(e.Systems, system)

		if hasVisited {
			system.Wormhole = &Coords{X: star.X, Y: star.Y, Z: star.Z}
			for pn := 0; pn < star.NumPlanets; pn++ {
				//p := ds.Planets[pn+star.PlanetIndex]
				planet := &PlanetData{Coords{X: star.X, Y: star.Y, Z: star.Z, Orbit: pn + 1}}
				system.Planets = append(system.Planets, planet)
			}
		}
	}
	sort.Sort(e.Systems)

	species := &SpeciesData{
		Id:       sp.Id,
		Name:     sp.Name,
		GovtName: sp.GovtName,
	}
	e.Species = species

	return &e, nil
}

var visited = []Coords{
	{X: 23, Y: 48, Z: 22},
	{X: 23, Y: 48, Z: 22},
	{X: 23, Y: 48, Z: 22},
	{X: 23, Y: 48, Z: 22},
	{X: 25, Y: 41, Z: 24},
	{X: 25, Y: 41, Z: 24},
	{X: 25, Y: 41, Z: 24},
	{X: 26, Y: 37, Z: 29},
	{X: 27, Y: 43, Z: 18},
	{X: 27, Y: 43, Z: 18},
	{X: 27, Y: 43, Z: 18},
	{X: 27, Y: 43, Z: 18},
	{X: 27, Y: 43, Z: 18},
	{X: 28, Y: 41, Z: 28},
	{X: 30, Y: 33, Z: 28},
	{X: 30, Y: 41, Z: 26},
	{X: 30, Y: 41, Z: 26},
	{X: 30, Y: 41, Z: 26},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 30, Y: 48, Z: 28},
	{X: 33, Y: 38, Z: 27},
}
