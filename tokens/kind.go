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

package tokens

// In an ORDERS file, we have only a few types of tokens.
//  1. There are keywords like "START," "POST-ARRIVAL," and "END."
//  2. There are commands like "AUTO," "JUMP," and "RESEARCH."
//  3. There are abbreviations like "ML," "MI," "CU," and "AU."
//  3. There are species names like "SP Iowan."
//  4. There are planet names like "PL Rigel" and "PL Taijyn 4."
//  5. There are ship names "TR1 Sam," "TR2S Jo Jo," "CC Jorge,"
//     and "FF Thomas." Note that ship includes transports,
//     warships, and star bases.
//  6. There are planet coordinates like "1 2 3 4."
//  7. There are system coordinates like "1 2 3."
//  8. There are non-negative integers.
//  9. There are messages, which are blocks of text that end
//     with the text "zzz" as the only text on a line.
// When parsing names and messages, invalid UTF-8 characters are
// silently changed to spaces. Each sequence is replaced with a
// single space, not each byte of the sequence.

type Kind int

const (
	Unknown Kind = iota
	// order keywords
	Start
	Combat
	PreDeparture
	Jumps
	Production
	PostArrival
	Strike
	End
	// command keywords
	Ally
	Ambush
	Attack
	Auto
	Base
	Battle
	Build
	Continue
	Deep
	Destroy
	Develop
	Disband
	Enemy
	Engage
	Estimate
	Haven
	Hide
	Hijack
	IBuild
	IContinue
	Install
	Intercept
	Jump
	Land
	Message
	Move
	Name
	Neutral
	Orbit
	PJump
	Recycle
	Repair
	Research
	Scan
	Send
	Shipyard
	Summary
	Surrender
	Target
	Teach
	Tech
	Telescope
	Terraform
	Transfer
	Unload
	Upgrade
	Visited
	Withdraw
	Wormhole
	// item abbreviations
	RM
	PD
	SU
	DR
	CU
	IU
	AU
	FS
	JP
	FM
	FJ
	GT
	FD
	TP
	GW
	SG1
	SG2
	SG3
	SG4
	SG5
	SG6
	SG7
	SG8
	SG9
	GU1
	GU2
	GU3
	GU4
	GU5
	GU6
	GU7
	GU8
	GU9
	SpeciesName
	PlanetName
	ShipName
	PlanetCoordinates
	SystemCoordinates
	Number
	MessageText
)
