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

// returns TRUE if a ship was found, FALSE if no ship was found
func (e *Engine) get_ship() int {
	var i, n, name_length, best_score, next_best_score, best_ship_index, first_try, minimum_score int
	var best_ship *ship_data

	upper_ship_name := make([]byte, 32, 32)

	// save in case of an error
	temp1_ptr := e.input_line_pointer

	// get ship abbreviation
	if e.get_class_abbr() == PLANET_ID {
		e.input_line_pointer = temp1_ptr
		return FALSE
	}

	temp2_ptr := e.input_line_pointer

	first_try = TRUE

again:

	e.input_line_pointer = temp2_ptr

	if e.abbr_type != SHIP_CLASS && first_try == FALSE {
		// assume abbreviation was accidentally omitted
		e.input_line_pointer = temp1_ptr
	}

	// get ship name
	name_length = e.get_name()

	/* Search all ships for name. */
	e.ship = nil
	for e.ship_index = 0; e.ship_index < e.species.num_ships; e.ship_index++ {
		e.ship = e.ship_base[e.ship_index]
		if e.ship.pn == 99 {
			continue
		}

		// make upper case copy of ship name
		for i = 0; i < 32; i++ {
			upper_ship_name[i] = toupper(e.ship.name[i])
		}

		// compare names
		if strcmp(upper_ship_name, e.upper_name) == 0 {
			e.abbr_type, e.abbr_index = SHIP_CLASS, e.ship.class
			e.correct_spelling_required = FALSE
			return TRUE
		}
	}

	if first_try != FALSE {
		first_try = FALSE
		goto again
	}

	if e.correct_spelling_required != FALSE {
		e.correct_spelling_required = FALSE
		return FALSE
	}

	// possibly a spelling error; find the best match that is approximately the same
	first_try = TRUE

yet_again:

	e.input_line_pointer = temp2_ptr

	if e.abbr_type != SHIP_CLASS && first_try == FALSE {
		// assume abbreviation was accidentally omitted
		e.input_line_pointer = temp1_ptr
	}

	// get ship name
	name_length = e.get_name()

	best_score = -9999
	next_best_score = -9999
	for e.ship_index = 0; e.ship_index < e.species.num_ships; e.ship_index++ {
		e.ship = e.ship_base[e.ship_index]
		if e.ship.pn == 99 {
			continue
		}

		// make upper case copy of ship name
		for i = 0; i < 32; i++ {
			upper_ship_name[i] = toupper(e.ship.name[i])
		}

		n = agrep_score(upper_ship_name, e.upper_name)
		if n > best_score {
			// best match so far
			best_score = n
			best_ship = e.ship
			best_ship_index = e.ship_index
		} else if n > next_best_score {
			next_best_score = n
		}
	}

	if best_ship == nil {
		return FALSE
	}
	name_length = len(best_ship.name)
	minimum_score = name_length - ((name_length / 7) + 1)

	if best_score < minimum_score /* Score too low. */ || name_length < 5 /* No errors allowed. */ || best_score == next_best_score /* Another name with equal score. */ {
		if first_try != FALSE {
			first_try = FALSE
			goto yet_again
		}
		e.correct_spelling_required = FALSE
		return FALSE
	}

	e.ship = best_ship
	e.ship_index = best_ship_index
	e.abbr_type = SHIP_CLASS
	e.abbr_index = e.ship.class
	e.correct_spelling_required = FALSE
	return TRUE
}
