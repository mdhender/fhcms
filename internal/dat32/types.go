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

// Package dat32 exports types that match the data structures in the 32-bit
// version of the C program. Well, match fairly closely.
package dat32

const MAX_ITEMS = 38
const MAX_SPECIES = 100

// galaxy_data is the layout in the binary data file.
type galaxy_data struct {
	/* Design number of species in galaxy. */
	DNumSpecies int32
	/* Actual number of species allocated. */
	NumSpecies int32
	/* Galactic radius in parsecs. */
	Radius int32
	/* Current turn number. */
	TurnNumber int32
}

// star_file_t is a helper that represents the layout in the binary data file.
type star_file_t struct {
	NumStars int32
	StarBase []star_data
}

// star_data is the layout in the binary data file.
type star_data struct {
	/* Coordinates. */
	X int8
	Y int8
	Z int8
	/* Dwarf, degenerate, main sequence or giant. */
	Type int8
	/* Star color. Blue, blue-white, etc. */
	Color int8
	/* Star size, from 0 thru 9 inclusive. */
	Size int8
	/* Number of usable planets in star system. */
	NumPlanets int8
	/* TRUE if this is a good potential home system. */
	HomeSystem int8
	/* TRUE if wormhole entry/exit. */
	WormHere int8
	/* Coordinates of exit point for wormhole. Valid only if WormHere is TRUE. */
	WormX int8
	WormY int8
	WormZ int8
	/* Reserved for future use. Zero for now. */
	Reserved1 int16
	Reserved2 int16
	/* Index (starting at zero) into the file "planets.dat" of the first planet in the star system. */
	PlanetIndex int16
	/* Message associated with this star system, if any. */
	Message int32
	/* A bit is set if corresponding species has been here. */
	VisitedBy [2]uint64
	/* Reserved for future use. Zero for now. */
	Reserved3 int32
	Reserved4 int32
	Reserved5 int32
	// padding to make Go struct same size as C
	Padding [2]uint8
}

// planet_file_t is a helper struct that represents the layout
// of data in the binary file
type planet_file_t struct {
	NumPlanets int32
	PlanetBase []planet_data
}

// planet_data is the layout in the binary data file.
type planet_data struct {
	/* Temperature class, 1-30. */
	TemperatureClass int8
	/* Pressure class, 0-29. */
	PressureClass int8
	/* 0 = not special, 1 = ideal home planet, 2 = ideal colony planet, 3 = radioactive hellhole. */
	Special int8
	/* Reserved for future use. Zero for now. */
	Reserved1 int8
	/* Gas in atmosphere. Zero if none. */
	Gas [4]int8
	/* Percentage of gas in atmosphere. */
	GasPercent [4]int8
	/* Reserved for future use. Zero for now. */
	Reserved2 int16
	/* Diameter in thousands of kilometers. */
	Diameter int16
	/* Surface gravity. Multiple of Earth gravity times 100. */
	Gravity int16
	/* Mining difficulty times 100. */
	MiningDifficulty int16
	/* Economic efficiency. Always 100 for a home planet. */
	EconEfficiency int16
	/* Increase in mining difficulty. */
	MDIncrease int16
	/* Message associated with this planet, if any. */
	Message int32
	/* Reserved for future use. Zero for now. */
	Reserved3 int32
	Reserved4 int32
	Reserved5 int32
}

// species_file_t is a helper that represents the data layout
// in the binary data files.
type species_file_t struct {
	Species  species_data
	NampData []nampla_data
	ShipData []ship_data
}

// species_data is the layout in the binary data file.
type species_data struct {
	/* Name of species. */
	Name [32]uint8
	/* Name of government. */
	GovtName [32]uint8
	/* Type of government. */
	GovtType [32]uint8
	/* Coordinates of home planet. */
	X  uint8
	Y  uint8
	Z  uint8
	PN uint8
	/* Gas required by species. */
	RequiredGas uint8
	/* Minimum needed percentage. */
	RequiredGasMin uint8
	/* Maximum allowed percentage. */
	RequiredGasMax uint8
	/* Reserved for future use. Zero for now. */
	Reserved5 uint8
	/* Gases neutral to species. */
	NeutralGas [6]uint8
	/* Gases poisonous to species. */
	PoisonGas [6]uint8
	/* AUTO command was issued. */
	AutoOrders uint8
	/* Reserved for future use. Zero for now. */
	Reserved3 uint8
	Reserved4 int16
	/* Actual tech levels. */
	TechLevel [6]int16
	/* Tech levels at start of turn. */
	InitTechLevel [6]int16
	/* Unapplied tech level knowledge. */
	TechKnowledge [6]int16
	/* Number of named planets, including home planet and colonies. */
	NumNamplas int32
	/* Number of ships. */
	NumShips int32
	/* Experience points for tech levels. */
	TechEps [6]int32
	/* If non-zero, home planet was bombed either by bombardment or germ warfare and has not yet fully recovered. Value is total economic base before bombing. */
	HPOriginalBase int32
	/* Number of economic units. */
	EconUnits int32
	/* Total fleet maintenance cost. */
	FleetCost int32
	/* Fleet maintenance cost as a percentage times one hundred. */
	FleetPercentCost int32
	/* A bit is set if corresponding species has been met. */
	Contact [2]uint64
	/* A bit is set if corresponding species is considered an ally. */
	Ally [2]uint64
	/* A bit is set if corresponding species is considered an enemy. */
	Enemy [2]uint64
	/* Use for expansion. Initialized to all zeroes. */
	Padding [12]uint8
}

// nampla_data is the layout in the binary data file.
type nampla_data struct {
	/* Name of planet. */
	Name [32]uint8
	/* Coordinates. */
	X  uint8
	Y  uint8
	Z  uint8
	PN uint8
	/* Status of planet. */
	Status uint8
	/* Reserved for future use. Zero for now. */
	Reserved1 uint8
	/* HIDE order given. */
	Hiding uint8
	/* Colony is hidden. */
	Hidden uint8
	/* Reserved for future use. Zero for now. */
	Reserved2 int16
	/* Index (starting at zero) into the file "planets.dat" of this planet. */
	PlanetIndex int16
	/* Siege effectiveness - a percentage between 0 and 99. */
	SiegeEff int16
	/* Number of shipyards on planet. */
	Shipyards int16
	/* Reserved for future use. Zero for now. */
	Reserved4 int32
	/* Incoming ship with only CUs on board. */
	IUsNeeded int32
	/* Incoming ship with only CUs on board. */
	AUsNeeded int32
	/* Number of IUs to be automatically installed. */
	AutoIUs int32
	/* Number of AUs to be automatically installed. */
	AutoAUs int32
	/* Reserved for future use. Zero for now. */
	Reserved5 int32
	/* Colonial mining units to be installed. */
	IUsToInstall int32
	/* Colonial manufacturing units to be installed. */
	AUsToInstall int32
	/* Mining base times 10. */
	MiBase int32
	/* Manufacturing base times 10. */
	MaBase int32
	/* Number of available population units. */
	PopUnits int32
	/* Quantity of each item available. */
	ItemQuantity [MAX_ITEMS]int32
	/* Reserved for future use. Zero for now. */
	Reserved6 int32
	/* Amount to use on ambush. */
	UseOnAmbush int32
	/* Message associated with this planet, if any. */
	Message int32
	/* Different for each application. */
	Special int32
	/* Use for expansion. Initialized to all zeroes. */
	Padding [28]uint8
}

// ship_data is the layout in the binary data file.
type ship_data struct {
	/* Name of ship. */
	Name [32]uint8
	/* Current coordinates. */
	X  uint8
	Y  uint8
	Z  uint8
	PN uint8
	/* Current status of ship. */
	Status uint8
	/* Ship type. */
	Type uint8
	/* Destination if ship was forced to jump from combat. Also used by TELESCOPE command. */
	DestX uint8
	DestY uint8
	DestZ uint8
	/* Set if ship jumped this turn. */
	JustJumped uint8
	/* Ship arrived via wormhole in the PREVIOUS turn. */
	ArrivedViaWormhole uint8
	/* Reserved for future use. Zero for now. */
	Reserved1 uint8
	Reserved2 int16
	Reserved3 int16
	/* Ship class. */
	Class int16
	/* Ship tonnage divided by 10,000. */
	Tonnage int16
	/* Quantity of each item carried. */
	ItemQuantity [MAX_ITEMS]int16
	/* Ship age. */
	Age int16
	/* The cost needed to complete the ship if still under construction. */
	RemainingCost int16
	/* Reserved for future use. Zero for now. */
	Reserved4 int16
	/* Nampla index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet. */
	LoadingPoint int16
	/* Nampla index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet. */
	UnloadingPoint int16
	/* Different for each application. */
	Special int32
	/* Use for expansion. Initialized to all zeroes. */
	Padding [28]uint8
	// padding to make Go struct same size as C
	MorePadding [2]uint8
}

type sp_loc_data struct {
	S       uint8 // species number
	X, Y, Z uint8 // coordinates
}
