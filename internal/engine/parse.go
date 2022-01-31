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

package engine

import (
	"bytes"
	"strconv"
	"strings"
)

/* The following "get" routines will return 0 if the item found was not
 * of the appropriate type, and 1 or greater if an item of the correct
 * type was found. */

/* Get a class abbreviation and return TECH_ID, ITEM_CLASS, SHIP_CLASS,
 * PLANET_ID, SPECIES_ID or ALLIANCE_ID as appropriate, or UNKNOWN if it
 * cannot be identified. Also, set "abbr_type" to this value. If it is
 * TECH_ID, ITEM_CLASS or SHIP_CLASS, "abbr_index" will contain the
 * abbreviation index. If it is a ship, "tonnage" will contain tonnage/10,000,
 * and "sub_light" will be TRUE or FALSE. (Tonnage value returned is based
 * ONLY on abbreviation.) */

func (e *Engine) get_class_abbr() int {
	e.skip_whitespace()
	e.abbr_type = UNKNOWN

	if !isalnum(e.input_line_pointer[0]) {
		return UNKNOWN
	}
	e.input_abbr[0], e.input_abbr[1] = toupper(e.input_line_pointer[0]), 0
	e.input_line_pointer = e.input_line_pointer[1:]

	if !isalnum(e.input_line_pointer[0]) {
		return UNKNOWN
	}
	e.input_abbr[1], e.input_abbr[2] = toupper(e.input_line_pointer[0]), 0
	e.input_line_pointer = e.input_line_pointer[1:]

	// check for IDs that are followed by one or more digits or letters
	i := 2
	digit_start := e.input_line_pointer
	for isalnum(e.input_line_pointer[0]) {
		e.input_abbr[i], e.input_abbr[i+1] = e.input_line_pointer[0], 0
		i, e.input_line_pointer = i+1, e.input_line_pointer[1:]
	}

	// check tech ID
	for i = 0; i < 6; i++ {
		if strcmp(e.input_abbr, []byte(tech_abbr[i])) == 0 {
			e.abbr_index, e.abbr_type = i, TECH_ID
			return e.abbr_type
		}
	}

	// check item abbreviations
	for i = 0; i < MAX_ITEMS; i++ {
		if strcmp(e.input_abbr, []byte(item_abbr[i])) == 0 {
			e.abbr_index, e.abbr_type = i, ITEM_CLASS
			return e.abbr_type
		}
	}

	// check ship abbreviations
	for i = 0; i < NUM_SHIP_CLASSES; i++ {
		if strncmp(e.input_abbr, []byte(ship_abbr[i]), 2) != 0 {
			continue
		}
		e.input_line_pointer = digit_start
		e.abbr_index = i
		e.tonnage = ship_tonnage[i]
		if i == TR {
			e.tonnage = 0
			for isdigit(e.input_line_pointer[0]) {
				e.tonnage = (10 * e.tonnage) + int(e.input_line_pointer[0]-'0')
				e.input_line_pointer = e.input_line_pointer[1:]
			}
		}
		e.sub_light = FALSE
		if toupper(e.input_line_pointer[0]) == 'S' {
			e.sub_light = TRUE
			e.input_line_pointer = e.input_line_pointer[1:]
		}
		if isalnum(e.input_line_pointer[0]) {
			break // garbage, not a ship class
		}
		e.abbr_type = SHIP_CLASS
		return e.abbr_type
	}

	// check for planet name
	if strcmp(e.input_abbr, []byte("PL")) == 0 {
		e.abbr_type = PLANET_ID
		return (e.abbr_type)
	}

	// check for species name
	if strcmp(e.input_abbr, []byte("SP")) == 0 {
		e.abbr_type = SPECIES_ID
		return e.abbr_type
	}

	e.abbr_type = UNKNOWN
	return e.abbr_type
}

/* Get a command and return its index. */
func (e *Engine) get_command() int {
	e.skip_junk()
	if e.end_of_file != FALSE {
		return -1
	}

	/* Get first three characters of command word. */
	var cmd_s string
	for i := 0; i < 3 && len(e.input_line_pointer) != 0; i++ {
		if !isalpha(e.input_line_pointer[0]) {
			return 0
		}
		cmd_s += string(e.input_line_pointer[:1])
		e.input_line_pointer = e.input_line_pointer[1:]
	}
	cmd_s = strings.ToUpper(cmd_s)

	/* Skip everything after third character of command word. */
	for len(e.input_line_pointer) != 0 && bytes.IndexByte([]byte{'\t', '\r', '\n', ' ', ',', ';'}, e.input_line_pointer[0]) == -1 {
		e.input_line_pointer = e.input_line_pointer[1:]
	}

	/* Find corresponding string in list. */
	cmd_n := UNKNOWN
	for i := 1; i < NUM_COMMANDS; i++ {
		if cmd_s == command_abbr[i] {
			cmd_n = i
			break
		}
	}

	return cmd_n
}

/* Get a name and copy original version to "original_name" and upper
 * case version to "upper_name". Return length of name. */
func (e *Engine) get_name() int {
	e.skip_whitespace()

	var name_length int
	for {
		c := e.input_line_pointer[0]
		if c == ';' {
			break
		}
		e.input_line_pointer = e.input_line_pointer[1:]
		if c == ',' || c == '\t' || c == '\n' {
			break
		}
		if name_length < 31 {
			e.original_name[name_length] = c
			e.upper_name[name_length] = toupper(c)
			name_length++
		}
	}

	// remove any trailing spaces in name
	for name_length > 0 {
		if e.original_name[name_length-1] != ' ' {
			break
		}
		name_length--
	}

	// terminate strings
	e.original_name[name_length] = 0
	e.upper_name[name_length] = 0

	return name_length
}

// get_value reads an integer and places its value in 'value'.
// ignores overflow.
// returns TRUE if it could read an integer, FALSE otherwise
func (e *Engine) get_value() int {
	e.skip_whitespace()
	if len(e.input_line_pointer) == 0 {
		return FALSE
	} else if !(e.input_line_pointer[0] == '+' || e.input_line_pointer[0] == '-' || isdigit(e.input_line_pointer[0])) {
		return FALSE
	}

	input := strings.Builder{}

	// consume sign or first digit
	input.WriteByte(e.input_line_pointer[0])
	e.input_line_pointer = e.input_line_pointer[1:]

	// consume the next run of digits
	for len(e.input_line_pointer) != 0 && isdigit(e.input_line_pointer[0]) {
		input.WriteByte(e.input_line_pointer[0])
		e.input_line_pointer = e.input_line_pointer[1:]
	}

	// convert the input to an int value. ignores overflow.
	e.value, _ = strconv.Atoi(input.String())

	return TRUE
}

/* Skip white space and comments. */
func (e *Engine) skip_junk() {
again:

	/* Read next line. */
	if e.input_line_pointer = fgets(e.input_line, 256, e.input_file); e.input_line_pointer == nil {
		e.end_of_file = TRUE
		return
	}

	if e.just_opened_file != FALSE { /* Skip mail header, if any. */
		if e.input_line_pointer[0] == '\n' {
			goto again
		}
		e.just_opened_file = FALSE
		if strncmp(e.input_line_pointer, []byte("From "), 5) == 0 { /* This is a mail header. */
			for {
				if e.input_line_pointer = fgets(e.input_line, 256, e.input_file); e.input_line_pointer == nil {
					e.end_of_file = TRUE /* Weird. */
					return
				}
				if e.input_line_pointer[0] == '\n' {
					break /* End of header. */
				}
			}
			goto again
		}
	}

	copy(e.original_line, e.input_line) // make a copy

	/* Skip white space and comments. */
	for len(e.input_line_pointer) != 0 {
		switch e.input_line_pointer[0] {
		case ';', '\n': /* Semi-colon. Newline. */
			goto again
		case '\t', ' ', ',', '\r': /* Tab. Space. Comma. CR. */
			e.input_line_pointer = e.input_line_pointer[1:]
		default:
			return
		}
	}
}

func (e *Engine) skip_whitespace() {
	for len(e.input_line_pointer) != 0 {
		switch e.input_line_pointer[0] {
		case '\t', ' ', ',', '\r': /* Tab. Space. Comma. CR. */
			e.input_line_pointer = e.input_line_pointer[1:]
		default:
			return
		}
	}
}
