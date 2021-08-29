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

package orders

type CommandType int

const (
	UNKNOWN CommandType = iota
	Ally
	Ambush
	Attack
	Auto
	Base
	Battle
	Build
	Continue
	Deep
	Debug
	Destroy
	Develop
	Disband
	End
	Enemy
	Engage
	Estimate
	Haven
	Hide
	Hijack
	Ibuild
	Icontinue
	Install
	Intercept
	Jump
	Land
	Message
	Move
	Name
	Neutral
	Orbit
	Pjump
	Production
	Recycle
	Rename
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
	Zzz
)
