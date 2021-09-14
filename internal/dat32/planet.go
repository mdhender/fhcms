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

// Planets includes the meta-data from the binary data.
type Planets struct {
	NumPlanets int      `json:"num_planets"`
	Planets    []Planet `json:"planets"`
}

// Planet is the planet from the binary data.
type Planet struct {
	Id int `json:"planet_index"`
	/* Temperature class, 1-30. */
	TemperatureClass int `json:"temperature_class"`
	/* Pressure class, 0-29. */
	PressureClass int `json:"pressure_class"`
	/* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	Special    int    `json:"special"`
	Gas        [4]int `json:"gas"`
	GasPercent [4]int `json:"gas_percent"`
	/* Diameter in thousands of kilometers. */
	Diameter int `json:"diameter"`
	/* Surface gravity. Multiple of Earth gravity times 100. */
	Gravity int `json:"gravity"`
	/* Mining difficulty times 100. */
	MiningDifficulty int `json:"mining_difficulty"`
	/* Economic efficiency. Always 100 for a home planet. */
	EconEfficiency int `json:"econ_efficiency"`
	/* Increase in mining difficulty. */
	MDIncrease int `json:"md_increase"`
	/* Message associated with this planet, if any. */
	Message int `json:"message"`
}

// ReadPlanets returns either an initialized set of planets or an error
func ReadPlanets(name string, bo binary.ByteOrder) (*Planets, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	var pd planet_file_t
	if err := binary.Read(r, bo, &pd.NumPlanets); err != nil {
		return nil, err
	}
	pd.PlanetBase = make([]planet_data, pd.NumPlanets, pd.NumPlanets)
	numPlanets := int(pd.NumPlanets)
	for i := 0; i < numPlanets; i++ {
		if err := binary.Read(r, bo, &pd.PlanetBase[i]); err != nil {
			return nil, err
		}
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
