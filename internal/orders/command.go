/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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

import "bytes"

type COMMAND int

// getCommand returns a command and error.
// A command is the command, parameters, and comments.
// Some commands, like START, have children.
// Some commands, like MESSAGE, will consume multiple lines.
func getCommand(lines [][]byte) (COMMAND, [][]byte, error) {
	if len(lines) == 0 {
		return EOF, nil, nil
	}
	line, _ := splitLine(lines[0])
	command := []byte{' ', ' ', ' '}
	for i := 0; i < 3 && i < len(line); i++ {
		if bytes.IndexByte([]byte{' ', '\t'}, line[i]) != -1 {
			break
		}
		command[i] = toUpper(line[i])
	}
	return EOF, nil, nil
}

// getCommandWord returns a command word and the remainder of the line.
func getCommandWord(line []byte) COMMAND {
	// command is the first three characters or up to the first space or tab
	var command []byte
	for len(command) != 3 && len(line) != 0 {
		var ch byte
		if ch, line = line[0], line[1:]; ch == ' ' || ch == '\t' {
			break
		}
		command = append(command, toUpper(ch))
	}

	switch string(command) {
	case "ALL":
		return ALLY
	case "AMB":
		return AMBUSH
	case "ATT":
		return ATTACK
	case "AUT":
		return AUTO
	case "BAS":
		return BASE
	case "BAT":
		return BATTLE
	case "BUI":
		return BUILD
	case "CON":
		return CONTINUE
	case "DEE":
		return DEEP
	case "DES":
		return DESTROY
	case "DEV":
		return DEVELOP
	case "DIS":
		return DISBAND
	case "END":
		return END
	case "ENE":
		return ENEMY
	case "ENG":
		return ENGAGE
	case "EST":
		return ESTIMATE
	case "HAV":
		return HAVEN
	case "HID":
		return HIDE
	case "HIJ":
		return HIJACK
	case "IBU":
		return IBUILD
	case "ICO":
		return ICONTINUE
	case "INS":
		return INSTALL
	case "INT":
		return INTERCEPT
	case "JUM":
		return JUMP
	case "LAN":
		return LAND
	case "MES":
		return MESSAGE
	case "MOV":
		return MOVE
	case "NAM":
		return NAME
	case "NEU":
		return NEUTRAL
	case "ORB":
		return ORBIT
	case "PJU":
		return PJUMP
	case "PRO":
		return PRODUCTION
	case "REC":
		return RECYCLE
	case "REN":
		return RENAME
	case "REP":
		return REPAIR
	case "RES":
		return RESEARCH
	case "SCA":
		return SCAN
	case "SEN":
		return SEND
	case "SHI":
		return SHIPYARD
	case "STA":
		return START
	case "SUM":
		return SUMMARY
	case "SUR":
		return SURRENDER
	case "TAR":
		return TARGET
	case "TEA":
		return TEACH
	case "TEC":
		return TECH
	case "TEL":
		return TELESCOPE
	case "TER":
		return TERRAFORM
	case "TRA":
		return TRANSFER
	case "UNL":
		return UNLOAD
	case "UPG":
		return UPGRADE
	case "VIS":
		return VISITED
	case "WIT":
		return WITHDRAW
	case "WOR":
		return WORMHOLE
	case "ZZZ":
		return ZZZ
	default:
		return UNDEFINED
	}
}

func toLower(ch byte) byte {
	if 'A' <= ch && ch <= 'Z' {
		return (ch - 'A') + 'a'
	}
	return ch
}

func toUpper(ch byte) byte {
	if 'a' <= ch && ch <= 'z' {
		return (ch - 'a') + 'A'
	}
	return ch
}

// trimSpaces returns the slice after trimming the leading whitespace
func trimSpaces(b []byte) []byte {
	for len(b) != 0 && (b[0] == ' ' || b[0] == '\t' || b[0] == '\r') {
		b = b[1:]
	}
	return b
}

const (
	UNDEFINED COMMAND = iota
	ALLY
	AMBUSH
	ATTACK
	AUTO
	BASE
	BATTLE
	BUILD
	CONTINUE
	DEEP
	DESTROY
	DEVELOP
	DISBAND
	END
	ENEMY
	ENGAGE
	ESTIMATE
	HAVEN
	HIDE
	HIJACK
	IBUILD
	ICONTINUE
	INSTALL
	INTERCEPT
	JUMP
	LAND
	MESSAGE
	MOVE
	NAME
	NEUTRAL
	ORBIT
	PJUMP
	PRODUCTION
	RECYCLE
	RENAME
	REPAIR
	RESEARCH
	SCAN
	SEND
	SHIPYARD
	START
	SUMMARY
	SURRENDER
	TARGET
	TEACH
	TECH
	TELESCOPE
	TERRAFORM
	TRANSFER
	UNLOAD
	UPGRADE
	VISITED
	WITHDRAW
	WORMHOLE
	ZZZ
	EOF
)
