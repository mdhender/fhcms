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

// Species is the species from the binary data along with the list
// of named planets and ships from the same file.
type Species struct {
	// species number
	Id int `json:"id"`
	/* Name of species. */
	Name string `json:"name"`
	/* Name of government. */
	GovtName string `json:"govt_name"`
	/* Type of government. */
	GovtType string `json:"govt_type"`
	/* Coordinates of home planet. */
	X  int `json:"x"`
	Y  int `json:"y"`
	Z  int `json:"z"`
	PN int `json:"pn"`
	/* Gas required by species. */
	RequiredGas int `json:"required_gas"`
	/* Minimum needed percentage. */
	RequiredGasMin int `json:"required_gas_min"`
	/* Maximum allowed percentage. */
	RequiredGasMax int `json:"required_gas_max"`
	/* Gases neutral to species. */
	NeutralGas []int `json:"neutral_gas"`
	/* Gases poisonous to species. */
	PoisonGas []int `json:"poison_gas"`
	/* AUTO command was issued. */
	AutoOrders bool `json:"auto_orders"`
	/* Actual tech levels. */
	TechLevel [6]int `json:"tech_level"`
	/* Tech levels at start of turn. */
	InitTechLevel [6]int `json:"init_tech_level"`
	/* Unapplied tech level knowledge. */
	TechKnowledge [6]int `json:"tech_knowledge"`
	/* Number of named planets, including home planet and colonies. */
	NumNamplas int `json:"num_namplas"`
	/* Number of ships. */
	NumShips int `json:"num_ships"`
	/* Experience points for tech levels. */
	TechEps [6]int `json:"tech_eps"`
	/* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	HPOriginalBase int `json:"hp_original_base"`
	/* Number of economic units. */
	EconUnits int `json:"econ_units"`
	/* Total fleet maintenance cost. */
	FleetCost int `json:"fleet_cost"`
	/* Fleet maintenance cost as a percentage times one hundred. */
	FleetPercentCost int `json:"fleet_percent_cost"`
	/* A bit is set if corresponding species has been met. */
	Contact []int `json:"contact"`
	/* A bit is set if corresponding species is considered an ally. */
	Ally []int `json:"ally"`
	/* A bit is set if corresponding species is considered an enemy. */
	Enemy []int `json:"enemy"`
	// All named planets (home planet and colonies)
	NamplaBase []NamedPlanet `json:"nampla_base"`
	// All ships, plus some slots tagged as UNUSED
	ShipBase []Ship `json:"ship_base"`
}

// ReadSpecies returns either an initialized Species or an error.
func ReadSpecies(name string, no int, bo binary.ByteOrder) (*Species, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(b)

	var sd species_file_t
	if err := binary.Read(r, bo, &sd.Species); err != nil {
		return nil, err
	}

	sd.NampData = make([]nampla_data, sd.Species.NumNamplas)
	numColonies := int(sd.Species.NumNamplas)
	for i := 0; i < numColonies; i++ {
		if err := binary.Read(r, bo, &sd.NampData[i]); err != nil {
			return nil, err
		}
	}

	sd.ShipData = make([]ship_data, sd.Species.NumShips)
	numShips := int(sd.Species.NumShips)
	for i := 0; i < numShips; i++ {
		if err := binary.Read(r, bo, &sd.ShipData[i]); err != nil {
			return nil, err
		}
	}

	species := Species{
		Contact: []int{},
		Ally:    []int{},
		Enemy:   []int{},
	}
	species.Id = no
	species.Name = nameToString(sd.Species.Name)
	species.GovtName = nameToString(sd.Species.GovtName)
	species.GovtType = nameToString(sd.Species.GovtType)
	species.X = int(sd.Species.X)
	species.Y = int(sd.Species.Y)
	species.Z = int(sd.Species.Z)
	species.PN = int(sd.Species.PN)
	species.RequiredGas = int(sd.Species.RequiredGas)
	species.RequiredGasMin = int(sd.Species.RequiredGasMin)
	species.RequiredGasMax = int(sd.Species.RequiredGasMax)
	for _, v := range sd.Species.NeutralGas {
		if v != 0 {
			species.NeutralGas = append(species.NeutralGas, int(v))
		}
	}
	for _, v := range sd.Species.PoisonGas {
		if v != 0 {
			species.PoisonGas = append(species.PoisonGas, int(v))
		}
	}
	species.AutoOrders = sd.Species.AutoOrders != 0
	for i := 0; i < 6; i++ {
		species.TechLevel[i] = int(sd.Species.TechLevel[i])
		species.InitTechLevel[i] = int(sd.Species.InitTechLevel[i])
		species.TechKnowledge[i] = int(sd.Species.TechKnowledge[i])
		species.TechEps[i] = int(sd.Species.TechEps[i])
	}
	species.NumNamplas = int(sd.Species.NumNamplas)
	species.NumShips = int(sd.Species.NumShips)
	species.HPOriginalBase = int(sd.Species.HPOriginalBase)
	species.EconUnits = int(sd.Species.EconUnits)
	species.FleetCost = int(sd.Species.FleetCost)
	species.FleetPercentCost = int(sd.Species.FleetPercentCost)
	for spIndex := 0; spIndex < MAX_SPECIES; spIndex++ {
		spNo := spIndex + 1
		if spIndex < 63 && (sd.Species.Contact[0]&(1<<spIndex)) != 0 {
			species.Contact = append(species.Contact, spNo)
		}
		if spIndex < 63 && (sd.Species.Ally[0]&(1<<spIndex)) != 0 {
			species.Ally = append(species.Ally, spNo)
		}
		if spIndex < 63 && (sd.Species.Enemy[0]&(1<<spIndex)) != 0 {
			species.Enemy = append(species.Enemy, spNo)
		}
	}
	species.NamplaBase = make([]NamedPlanet, species.NumNamplas, species.NumNamplas)
	for i := 0; i < species.NumNamplas; i++ {
		species.NamplaBase[i].Name = nameToString(sd.NampData[i].Name)
		species.NamplaBase[i].X = int(sd.NampData[i].X)
		species.NamplaBase[i].Y = int(sd.NampData[i].Y)
		species.NamplaBase[i].Z = int(sd.NampData[i].Z)
		species.NamplaBase[i].PN = int(sd.NampData[i].PN)
		species.NamplaBase[i].Status = int(sd.NampData[i].Status)
		species.NamplaBase[i].Hiding = sd.NampData[i].Hiding != 0
		species.NamplaBase[i].Hidden = sd.NampData[i].Hidden != 0
		species.NamplaBase[i].PlanetIndex = int(sd.NampData[i].PlanetIndex)
		species.NamplaBase[i].SiegeEff = int(sd.NampData[i].SiegeEff)
		species.NamplaBase[i].Shipyards = int(sd.NampData[i].Shipyards)
		species.NamplaBase[i].IUsNeeded = int(sd.NampData[i].IUsNeeded)
		species.NamplaBase[i].AUsNeeded = int(sd.NampData[i].AUsNeeded)
		species.NamplaBase[i].AutoIUs = int(sd.NampData[i].AutoIUs)
		species.NamplaBase[i].AutoAUs = int(sd.NampData[i].AutoAUs)
		species.NamplaBase[i].IUsToInstall = int(sd.NampData[i].IUsToInstall)
		species.NamplaBase[i].AUsToInstall = int(sd.NampData[i].AUsToInstall)
		species.NamplaBase[i].MiBase = int(sd.NampData[i].MiBase)
		species.NamplaBase[i].MaBase = int(sd.NampData[i].MaBase)
		species.NamplaBase[i].PopUnits = int(sd.NampData[i].PopUnits)
		for n := 0; n < len(sd.NampData[i].ItemQuantity); n++ {
			species.NamplaBase[i].ItemQuantity[n] = int(sd.NampData[i].ItemQuantity[n])
		}
		species.NamplaBase[i].UseOnAmbush = int(sd.NampData[i].UseOnAmbush)
		species.NamplaBase[i].Message = int(sd.NampData[i].Message)
		species.NamplaBase[i].Special = int(sd.NampData[i].Special)
	}
	species.ShipBase = make([]Ship, species.NumShips, species.NumShips)
	for i := 0; i < species.NumShips; i++ {
		species.ShipBase[i].Name = nameToString(sd.ShipData[i].Name)
		species.ShipBase[i].X = int(sd.ShipData[i].X)
		species.ShipBase[i].Y = int(sd.ShipData[i].Y)
		species.ShipBase[i].Z = int(sd.ShipData[i].Z)
		species.ShipBase[i].PN = int(sd.ShipData[i].PN)
		species.ShipBase[i].Status = int(sd.ShipData[i].Status)
		species.ShipBase[i].Type = int(sd.ShipData[i].Type)
		species.ShipBase[i].DestX = int(sd.ShipData[i].DestX)
		species.ShipBase[i].DestY = int(sd.ShipData[i].DestY)
		species.ShipBase[i].DestZ = int(sd.ShipData[i].DestZ)
		species.ShipBase[i].JustJumped = sd.ShipData[i].JustJumped != 0
		species.ShipBase[i].ArrivedViaWormhole = sd.ShipData[i].ArrivedViaWormhole != 0
		species.ShipBase[i].Class = int(sd.ShipData[i].Class)
		species.ShipBase[i].Tonnage = int(sd.ShipData[i].Tonnage)
		for n := 0; n < len(sd.ShipData[i].ItemQuantity); n++ {
			species.ShipBase[i].ItemQuantity[n] = int(sd.ShipData[i].ItemQuantity[n])
		}
		species.ShipBase[i].Age = int(sd.ShipData[i].Age)
		species.ShipBase[i].RemainingCost = int(sd.ShipData[i].RemainingCost)
		species.ShipBase[i].LoadingPoint = int(sd.ShipData[i].LoadingPoint)
		species.ShipBase[i].UnloadingPoint = int(sd.ShipData[i].UnloadingPoint)
		species.ShipBase[i].Special = int(sd.ShipData[i].Special)
	}

	return &species, nil
}
