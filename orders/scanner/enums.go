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

package scanner

import "fmt"

type TokenType int

const (
	Unknown TokenType = iota
	Ally
	Ambush
	Attack
	Auto
	Base
	Battle
	Build
	Command //
	Continue
	Deep
	Debug
	Destroy
	Develop
	Disband
	EOF //
	EOL //
	End
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
	Item //
	Jump
	Land
	Message //
	Move
	Name
	Neutral
	Number //
	Orbit
	PJump
	Planet            //
	PlanetCoordinates //
	Production
	Recycle
	Rename
	Repair
	Research
	Scan
	Section //
	Send
	Ship //
	Shipyard
	Species           //
	SystemCoordinates //
	Summary
	Surrender
	Target
	Teach
	Tech //
	Telescope
	Terraform
	Transfer
	Transport //
	Unload
	Upgrade
	Visited
	Withdraw
	Wormhole
	Zzz
)

func (tk TokenType) String() string {
	switch tk {
	case Ally:
		return "ally"
	case Ambush:
		return "ambush"
	case Attack:
		return "attack"
	case Auto:
		return "auto"
	case Base:
		return "base"
	case Battle:
		return "battle"
	case Build:
		return "build"
	case Command:
		return "command"
	case Continue:
		return "continue"
	case Deep:
		return "deep"
	case Debug:
		return "debug"
	case Destroy:
		return "destroy"
	case Develop:
		return "develop"
	case Disband:
		return "disband"
	case EOF:
		return "eOF"
	case EOL:
		return "eOL"
	case End:
		return "end"
	case Enemy:
		return "enemy"
	case Engage:
		return "engage"
	case Estimate:
		return "estimate"
	case Haven:
		return "haven"
	case Hide:
		return "hide"
	case Hijack:
		return "hijack"
	case IBuild:
		return "iBuild"
	case IContinue:
		return "iContinue"
	case Install:
		return "install"
	case Intercept:
		return "intercept"
	case Item:
		return "item"
	case Jump:
		return "jump"
	case Land:
		return "land"
	case Message:
		return "message"
	case Move:
		return "move"
	case Name:
		return "name"
	case Neutral:
		return "neutral"
	case Number:
		return "number"
	case Orbit:
		return "orbit"
	case PJump:
		return "pJump"
	case Planet:
		return "planet"
	case PlanetCoordinates:
		return "planetCoordinates"
	case Production:
		return "production"
	case Recycle:
		return "recycle"
	case Rename:
		return "rename"
	case Repair:
		return "repair"
	case Research:
		return "research"
	case Scan:
		return "scan"
	case Section:
		return "section"
	case Send:
		return "send"
	case Ship:
		return "ship"
	case Shipyard:
		return "shipyard"
	case Species:
		return "species"
	case SystemCoordinates:
		return "systemCoordinates"
	case Summary:
		return "summary"
	case Surrender:
		return "surrender"
	case Target:
		return "target"
	case Teach:
		return "teach"
	case Tech:
		return "tech"
	case Telescope:
		return "telescope"
	case Terraform:
		return "terraform"
	case Transfer:
		return "transfer"
	case Transport:
		return "transport"
	case Unload:
		return "unload"
	case Upgrade:
		return "upgrade"
	case Visited:
		return "visited"
	case Withdraw:
		return "withdraw"
	case Wormhole:
		return "wormhole"
	case Zzz:
		return "zzz"
	}
	return fmt.Sprintf("token(%d)", tk)
}
