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

type species_file_t struct {
	Species  species_data
	NampData []nampla_data
	ShipData []ship_data
}

type species_data struct {
	Name             [32]uint8 /* Name of species. */
	GovtName         [32]uint8 /* Name of government. */
	GovtType         [32]uint8 /* Type of government. */
	X                uint8     /* Coordinates of home planet. */
	Y                uint8
	Z                uint8
	PN               uint8
	RequiredGas      uint8     /* Gas required by species. */
	RequiredGasMin   uint8     /* Minimum needed percentage. */
	RequiredGasMax   uint8     /* Maximum allowed percentage. */
	Reserved5        uint8     /* Zero for now. */
	NeutralGas       [6]uint8  /* Gases neutral to species. */
	PoisonGas        [6]uint8  /* Gases poisonous to species. */
	AutoOrders       uint8     /* AUTO command was issued. */
	Reserved3        uint8     /* Zero for now. */
	Reserved4        int16     /* Zero for now. */
	TechLevel        [6]int16  /* Actual tech levels. */
	InitTechLevel    [6]int16  /* Tech levels at start of turn. */
	TechKnowledge    [6]int16  /* Unapplied tech level knowledge. */
	NumNamplas       int32     /* Number of named planets, including home planet and colonies. */
	NumShips         int32     /* Number of ships. */
	TechEps          [6]int32  /* Experience points for tech levels. */
	HPOriginalBase   int32     /* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	EconUnits        int32     /* Number of economic units. */
	FleetCost        int32     /* Total fleet maintenance cost. */
	FleetPercentCost int32     /* Fleet maintenance cost as a percentage times one hundred. */
	Contact          [2]uint64 /* A bit is set if corresponding species has been met. */
	Ally             [2]uint64 /* A bit is set if corresponding species is considered an ally. */
	Enemy            [2]uint64 /* A bit is set if corresponding species is considered an enemy. */
	Padding          [12]uint8 /* Use for expansion. Initialized to all zeroes. */
}

type nampla_data struct {
	Name         [32]uint8        /* Name of planet. */
	X            uint8            /* Coordinates. */
	Y            uint8            /* Coordinates. */
	Z            uint8            /* Coordinates. */
	PN           uint8            /* Coordinates. */
	Status       uint8            /* Status of planet. */
	Reserved1    uint8            /* Zero for now. */
	Hiding       uint8            /* HIDE order given. */
	Hidden       uint8            /* Colony is hidden. */
	Reserved2    int16            /* Zero for now. */
	PlanetIndex  int16            /* Index (starting at zero) into the file "planets.dat" of this planet. */
	SiegeEff     int16            /* Siege effectiveness - a percentage between 0 and 99. */
	Shipyards    int16            /* Number of shipyards on planet. */
	Reserved4    int32            /* Zero for now. */
	IUsNeeded    int32            /* Incoming ship with only CUs on board. */
	AUsNeeded    int32            /* Incoming ship with only CUs on board. */
	AutoIUs      int32            /* Number of IUs to be automatically installed. */
	AutoAUs      int32            /* Number of AUs to be automatically installed. */
	Reserved5    int32            /* Zero for now. */
	IUsToInstall int32            /* Colonial mining units to be installed. */
	AUsToInstall int32            /* Colonial manufacturing units to be installed. */
	MiBase       int32            /* Mining base times 10. */
	MaBase       int32            /* Manufacturing base times 10. */
	PopUnits     int32            /* Number of available population units. */
	ItemQuantity [MAX_ITEMS]int32 /* Quantity of each item available. */
	Reserved6    int32            /* Zero for now. */
	UseOnAmbush  int32            /* Amount to use on ambush. */
	Message      int32            /* Message associated with this planet, if any. */
	Special      int32            /* Different for each application. */
	Padding      [28]uint8        /* Use for expansion. Initialized to all zeroes. */
}

type ship_data struct {
	Name               [32]uint8        /* Name of ship. */
	X                  uint8            /* Current coordinates. */
	Y                  uint8            /* Current coordinates. */
	Z                  uint8            /* Current coordinates. */
	PN                 uint8            /* Current coordinates. */
	Status             uint8            /* Current status of ship. */
	Type               uint8            /* Ship type. */
	DestX              uint8            /* Destination if ship was forced to jump from combat. */
	DestY              uint8            /* Destination if ship was forced to jump from combat. */
	DestZ              uint8            /* Ditto. Also used by TELESCOPE command. */
	JustJumped         uint8            /* Set if ship jumped this turn. */
	ArrivedViaWormhole uint8            /* Ship arrived via wormhole in the PREVIOUS turn. */
	Reserved1          uint8            /* Unused. Zero for now. */
	Reserved2          int16            /* Unused. Zero for now. */
	Reserved3          int16            /* Unused. Zero for now. */
	Class              int16            /* Ship class. */
	Tonnage            int16            /* Ship tonnage divided by 10,000. */
	ItemQuantity       [MAX_ITEMS]int16 /* Quantity of each item carried. */
	Age                int16            /* Ship age. */
	RemainingCost      int16            /* The cost needed to complete the ship if still under construction. */
	Reserved4          int16            /* Unused. Zero for now. */
	LoadingPoint       int16            /* Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet. */
	UnloadingPoint     int16            /* Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet. */
	Special            int32            /* Different for each application. */
	Padding            [28]uint8        /* Use for expansion. Initialized to all zeroes. */
	MorePadding        [2]uint8         // padding to make Go struct same size as C
}

type Species struct {
	Id               int      `json:"id"`        // species number
	Name             string   `json:"name"`      /* Name of species. */
	GovtName         string   `json:"govt_name"` /* Name of government. */
	GovtType         string   `json:"govt_type"` /* Type of government. */
	X                int      `json:"x"`         /* Coordinates of home planet. */
	Y                int      `json:"y"`
	Z                int      `json:"z"`
	PN               int      `json:"pn"`
	RequiredGas      int      `json:"required_gas"`     /* Gas required by species. */
	RequiredGasMin   int      `json:"required_gas_min"` /* Minimum needed percentage. */
	RequiredGasMax   int      `json:"required_gas_max"` /* Maximum allowed percentage. */
	NeutralGas       []int    `json:"neutral_gas"`      /* Gases neutral to species. */
	PoisonGas        []int    `json:"poison_gas"`       /* Gases poisonous to species. */
	AutoOrders       bool     `json:"auto_orders"`      /* AUTO command was issued. */
	TechLevel        [6]int   `json:"tech_level"`
	InitTechLevel    [6]int   `json:"init_tech_level"`
	TechKnowledge    [6]int   `json:"tech_knowledge"`
	NumNamplas       int      `json:"num_namplas"` /* Number of named planets, including home planet and colonies. */
	NumShips         int      `json:"num_ships"`   /* Number of ships. */
	TechEps          [6]int   `json:"tech_eps"`
	HPOriginalBase   int      `json:"hp_original_base"`   /* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	EconUnits        int      `json:"econ_units"`         /* Number of economic units. */
	FleetCost        int      `json:"fleet_cost"`         /* Total fleet maintenance cost. */
	FleetPercentCost int      `json:"fleet_percent_cost"` /* Fleet maintenance cost as a percentage times one hundred. */
	Contact          []int    `json:"contact"`            /* A bit is set if corresponding species has been met. */
	Ally             []int    `json:"ally"`               /* A bit is set if corresponding species is considered an ally. */
	Enemy            []int    `json:"enemy"`              /* A bit is set if corresponding species is considered an enemy. */
	NamplaBase       []Nampla `json:"nampla_base"`
	ShipBase         []Ship   `json:"ship_base"`
}

type Nampla struct {
	Name         string         `json:"name"`
	X            int            `json:"x"`
	Y            int            `json:"y"`
	Z            int            `json:"z"`
	PN           int            `json:"pn"`
	Status       int            `json:"status"`
	Hiding       bool           `json:"hiding"`
	Hidden       bool           `json:"hidden"`
	PlanetIndex  int            `json:"planet_index"`
	SiegeEff     int            `json:"siege_eff"`
	Shipyards    int            `json:"shipyards"`
	IUsNeeded    int            `json:"ius_needed"`
	AUsNeeded    int            `json:"aus_needed"`
	AutoIUs      int            `json:"auto_ius"`
	AutoAUs      int            `json:"auto_aus"`
	IUsToInstall int            `json:"ius_to_install"`
	AUsToInstall int            `json:"aus_to_install"`
	MiBase       int            `json:"mi_base"`
	MaBase       int            `json:"ma_base"`
	PopUnits     int            `json:"pop_units"`
	ItemQuantity [MAX_ITEMS]int `json:"item_quantity"`
	UseOnAmbush  int            `json:"use_on_ambush"`
	Message      int            `json:"message"`
	Special      int            `json:"special"`
}

type Ship struct {
	Name               string         `json:"name"`
	X                  int            `json:"x"`
	Y                  int            `json:"y"`
	Z                  int            `json:"z"`
	PN                 int            `json:"pn"`
	Status             int            `json:"status"`
	Type               int            `json:"type"`
	DestX              int            `json:"dest_x"`
	DestY              int            `json:"dest_y"`
	DestZ              int            `json:"dest_z"`
	JustJumped         bool           `json:"just_jumped"`
	ArrivedViaWormhole bool           `json:"arrived_via_wormhole"`
	Class              int            `json:"class"`
	Tonnage            int            `json:"tonnage"`
	ItemQuantity       [MAX_ITEMS]int `json:"item_quantity"`
	Age                int            `json:"age"`
	RemainingCost      int            `json:"remaining_cost"`
	LoadingPoint       int            `json:"loading_point"`
	UnloadingPoint     int            `json:"unloading_point"`
	Special            int            `json:"special"`
}

type TechLevel struct {
	Code             string `json:"code"`
	Level            int    `json:"level"`             /* Actual tech levels. */
	InitialLevel     int    `json:"initial_level"`     /* Tech levels at start of turn. */
	KnowledgeLevel   int    `json:"knowledge_level"`   /* Unapplied tech level knowledge. */
	ExperiencePoints int    `json:"experience_points"` /* Experience points for tech levels. */
}

func nameToString(name [32]uint8) string {
	var s string
	for _, ch := range name {
		if ch == 0 {
			break
		}
		s += string(ch)
	}
	return s
}

func ReadSpecies(name string, no int, bo binary.ByteOrder) (*Species, error) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	fmt.Printf("species: binary file %q contains %d bytes\n", name, len(b))
	r := bytes.NewReader(b)

	var sd species_file_t
	if err := binary.Read(r, bo, &sd.Species); err != nil {
		return nil, err
	}

	fmt.Printf("species: num colonies %d\n", sd.Species.NumNamplas)
	sd.NampData = make([]nampla_data, sd.Species.NumNamplas)
	numColonies := int(sd.Species.NumNamplas)
	for i := 0; i < numColonies; i++ {
		if err := binary.Read(r, bo, &sd.NampData[i]); err != nil {
			return nil, err
		}
		//if i < 3 {
		//	ec := json.NewEncoder(os.Stdout)
		//	ec.SetIndent("", "  ")
		//	if err := ec.Encode(sd.NampData[i]); err != nil {
		//		return nil, err
		//	}
		//}
	}

	fmt.Printf("species: num colonies %d\n", sd.Species.NumNamplas)
	fmt.Printf("species: num ships    %d\n", sd.Species.NumShips)
	sd.ShipData = make([]ship_data, sd.Species.NumShips)
	numShips := int(sd.Species.NumShips)
	for i := 0; i < numShips; i++ {
		if err := binary.Read(r, bo, &sd.ShipData[i]); err != nil {
			return nil, err
		}
		//if i < 3 {
		//	ec := json.NewEncoder(os.Stdout)
		//	ec.SetIndent("", "  ")
		//	if err := ec.Encode(sd.Ship[i]); err != nil {
		//		return nil, err
		//	}
		//}
	}

	var species Species
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
	//fmt.Printf("sp%02d contact %10d %32b  ally %10d %32b  enemy %10d\n", no, sd.Species.Contact[0], sd.Species.Contact[0], sd.Species.Ally[0], sd.Species.Ally[0], sd.Species.Enemy[0])
	for sp := 1; sp <= MAX_SPECIES; sp++ {
		if sp < 64 && (sd.Species.Contact[0]&(1<<(sp-1))) != 0 {
			species.Contact = append(species.Contact, sp)
		}
		if sp < 64 && (sd.Species.Ally[0]&(1<<(sp-1))) != 0 {
			species.Ally = append(species.Ally, sp)
		}
		if sp < 64 && (sd.Species.Enemy[0]&(1<<(sp-1))) != 0 {
			species.Enemy = append(species.Enemy, sp)
		}
	}
	if species.Contact == nil {
		species.Contact = []int{} // just to make the json look pretty
	}
	if species.Ally == nil {
		species.Ally = []int{} // just to make the json look pretty
	}
	if species.Enemy == nil {
		species.Enemy = []int{} // just to make the json look pretty
	}
	species.NamplaBase = make([]Nampla, species.NumNamplas, species.NumNamplas)
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

func (s Species) Write() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	//fmt.Println(name, string(data))
	return ioutil.WriteFile(fmt.Sprintf("sp%02d.json", s.Id), data, 0644)
}
