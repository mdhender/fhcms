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

package domain

import (
	"encoding/binary"
	"github.com/mdhender/fhcms/internal/dat32"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

type Specie struct {
	No         string
	Name       string
	Government struct {
		Name string
	}
	Technology [6]struct {
		Name             string
		Code             string
		InitialLevel     int
		CurrentLevel     int
		KnowledgeLevel   int
		ExperiencePoints int
	}
	Stats []*SpecieStat
}

type SpecieStat struct {
	Label string
	Value float64
	Units string
}

func (s *Store) loadSpecie(files, spNo string) (*Specie, error) {
	bigEndian := viper.GetBool("files.big_endian")
	var bo binary.ByteOrder
	if bigEndian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}
	sp, err := dat32.ReadSpecies(filepath.Join(files, "sp"+spNo+".dat"), 0, bo)
	if err != nil {
		log.Printf("[domain] loadSpecie %q %q %+v\n", files, spNo, err)
		return &Specie{}, err
	}
	o := &Specie{
		No:   spNo,
		Name: sp.Name,
	}
	o.Government.Name = sp.GovtName
	for i := 0; i < 6; i++ {
		switch i {
		case 0:
			o.Technology[i].Code = "MI"
			o.Technology[i].Name = "Mining"
		case 1:
			o.Technology[i].Code = "MA"
			o.Technology[i].Name = "Manufacturing"
		case 2:
			o.Technology[i].Code = "ML"
			o.Technology[i].Name = "Military"
		case 3:
			o.Technology[i].Code = "GV"
			o.Technology[i].Name = "Gravitics"
		case 4:
			o.Technology[i].Code = "LS"
			o.Technology[i].Name = "Life Support"
		case 5:
			o.Technology[i].Code = "BI"
			o.Technology[i].Name = "Biology"
		}
		o.Technology[i].CurrentLevel = sp.TechLevel[i]
		o.Technology[i].InitialLevel = sp.InitTechLevel[i]
		o.Technology[i].KnowledgeLevel = sp.TechKnowledge[i]
		o.Technology[i].ExperiencePoints = sp.TechEps[i]
	}

	shipyards := 0
	for _, nampla := range sp.NamplaBase {
		if nampla.PN == 99 {
			continue
		}
		shipyards += nampla.Shipyards
	}
	o.Stats = append(o.Stats, &SpecieStat{Label: "Shipyards", Value: float64(shipyards), Units: "yards"})

	// why the check on fleet maintenance cost?
	var fleetMaintenancePct float64
	if sp.FleetPercentCost < 0 {
		fleetMaintenancePct = 0
	} else if sp.FleetPercentCost < 10000 {
		fleetMaintenancePct = float64(sp.FleetPercentCost) / 100
	} else {
		fleetMaintenancePct = 100
	}
	o.Stats = append(o.Stats, &SpecieStat{Label: "Fleet Maintenance", Value: fleetMaintenancePct, Units: "%"})

	return o, nil
}
