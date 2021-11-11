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

type planet_file_t struct {
	NumPlanets int32
	PlanetBase []planet_data
}

type planet_data struct {
	TemperatureClass int8    /* Temperature class, 1-30. */
	PressureClass    int8    /* Pressure class, 0-29. */
	Special          int8    /* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	Reserved1        int8    /* Reserved for future use. Zero for now. */
	Gas              [4]int8 /* Gas in atmosphere. Zero if none. */
	GasPercent       [4]int8 /* Percentage of gas in atmosphere. */
	Reserved2        int16   /* Reserved for future use. Zero for now. */
	Diameter         int16   /* Diameter in thousands of kilometers. */
	Gravity          int16   /* Surface gravity. Multiple of Earth gravity times 100. */
	MiningDifficulty int16   /* Mining difficulty times 100. */
	EconEfficiency   int16   /* Economic efficiency. Always 100 for a home planet. */
	MDIncrease       int16   /* Increase in mining difficulty. */
	Message          int32   /* Message associated with this planet, if any. */
	Reserved3        int32   /* Reserved for future use. Zero for now. */
	Reserved4        int32   /* Reserved for future use. Zero for now. */
	Reserved5        int32   /* Reserved for future use. Zero for now. */
}

type Planets struct {
	NumPlanets int      `json:"num_planets"`
	Planets    []Planet `json:"planets"`
}

type Planet struct {
	Id               int    `json:"planet_index"`
	TemperatureClass int    `json:"temperature_class"` /* Temperature class, 1-30. */
	PressureClass    int    `json:"pressure_class"`    /* Pressure class, 0-29. */
	Special          int    `json:"special"`           /* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	Gas              [4]int `json:"gas"`
	GasPercent       [4]int `json:"gas_percent"`
	Diameter         int    `json:"diameter"`          /* Diameter in thousands of kilometers. */
	Gravity          int    `json:"gravity"`           /* Surface gravity. Multiple of Earth gravity times 100. */
	MiningDifficulty int    `json:"mining_difficulty"` /* Mining difficulty times 100. */
	EconEfficiency   int    `json:"econ_efficiency"`   /* Economic efficiency. Always 100 for a home planet. */
	MDIncrease       int    `json:"md_increase"`       /* Increase in mining difficulty. */
	Message          int    `json:"message"`           /* Message associated with this planet, if any. */
}

type GasType struct {
	Code       string `json:"code"`
	Percentage int    `json:"percentage"`
}

func ReadPlanets(name string, bo binary.ByteOrder) (*Planets, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	fmt.Printf("planet: binary file %q contains %d bytes\n", name, len(b))
	r := bytes.NewReader(b)

	var pd planet_file_t
	if err := binary.Read(r, bo, &pd.NumPlanets); err != nil {
		return nil, err
	}
	fmt.Printf("planet: num planets %d\n", pd.NumPlanets)
	pd.PlanetBase = make([]planet_data, pd.NumPlanets, pd.NumPlanets)
	numPlanets := int(pd.NumPlanets)
	for i := 0; i < numPlanets; i++ {
		if err := binary.Read(r, bo, &pd.PlanetBase[i]); err != nil {
			return nil, err
		}
		//if i < 3 {
		//	ec := json.NewEncoder(os.Stdout)
		//	ec.SetIndent("", "  ")
		//	if err := ec.Encode(pd.PlanetBase[i]); err != nil {
		//		return nil, err
		//	}
		//}
	}

	var planets Planets
	planets.NumPlanets = int(pd.NumPlanets)
	planets.Planets = make([]Planet, planets.NumPlanets, planets.NumPlanets)
	for n := 0; n < planets.NumPlanets; n++ {
		planets.Planets[n].Id = n
		planets.Planets[n].TemperatureClass = int(pd.PlanetBase[n].TemperatureClass)
		planets.Planets[n].PressureClass = int(pd.PlanetBase[n].PressureClass)
		planets.Planets[n].Special = int(pd.PlanetBase[n].Special)
		for i := 0; i < len(pd.PlanetBase[n].Gas); i++ {
			planets.Planets[n].Gas[i] = int(pd.PlanetBase[n].Gas[i])
			planets.Planets[n].GasPercent[i] = int(pd.PlanetBase[n].GasPercent[i])
		}
		planets.Planets[n].Diameter = int(pd.PlanetBase[n].Diameter)
		planets.Planets[n].Gravity = int(pd.PlanetBase[n].Gravity)
		planets.Planets[n].MiningDifficulty = int(pd.PlanetBase[n].MiningDifficulty)
		planets.Planets[n].EconEfficiency = int(pd.PlanetBase[n].EconEfficiency)
		planets.Planets[n].MDIncrease = int(pd.PlanetBase[n].MDIncrease)
		planets.Planets[n].Message = int(pd.PlanetBase[n].Message)
	}
	return &planets, nil
}

func (p Planets) Write(name string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, data, 0644)
}
