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
	"fmt"
	"io/ioutil"
	"log"
)

func (e *Engine) log_char(c byte) {
	if e.logging_disabled != FALSE || c == '\r' {
		return
	}

	/* Check if current line is getting too long. */
	if (c == ' ' || c == '\n') && e.log_position > 77 {
		// find the closest preceding space.
		// BUG: assumes that there is a preceding space.
		temp_position := e.log_position - 1
		for temp_position >= 0 && e.log_line[temp_position] != ' ' {
			temp_position--
		}

		// write front of line to files
		temp_char := e.log_line[temp_position+1]
		e.log_line[temp_position] = '\n'
		e.log_line[temp_position+1] = 0
		if e.log_to_file != FALSE {
			fputs(e.log_line, e.log_file)
		}
		if e.log_stdout != FALSE {
			fputs(e.log_line, e.stdout)
		}
		if e.log_summary != FALSE {
			fputs(e.log_line, e.summary_file)
		}
		e.log_line[temp_position+1] = temp_char

		// copy overflow word to beginning of next line
		e.log_line[e.log_position] = 0
		e.log_position = e.log_indentation + 2
		for i := 0; i < e.log_position; i++ {
			e.log_line[i], e.log_line[i+1] = '.', 0
		}
		for i := 0; i < temp_position+1; i++ {
			e.log_line[e.log_position+i] = e.log_line[temp_position+1+i]
			e.log_line[e.log_position+i+1] = 0
		}
		e.log_position = strlen(e.log_line)

		if c == ' ' {
			e.log_line[e.log_position] = ' '
			e.log_position++
			return
		}
	}

	// check if line is being manually terminated
	if c == '\n' {
		// write current line to output
		e.log_line[e.log_position] = '\n'
		e.log_line[e.log_position+1] = 0
		if e.log_to_file != FALSE {
			fputs(e.log_line, e.log_file)
		}
		if e.log_stdout != FALSE {
			fputs(e.log_line, e.stdout)
		}
		if e.log_summary != FALSE {
			fputs(e.log_line, e.summary_file)
		}

		// set up for next line
		e.log_position, e.log_indentation, e.log_start_of_line = 0, 0, TRUE
		e.log_line[e.log_position] = 0

		return
	}

	// save this character
	e.log_line[e.log_position] = c
	e.log_position++
	e.log_line[e.log_position] = 0

	if e.log_start_of_line != FALSE && c == ' ' {
		// increment number of indenting spaces for current line
		e.log_indentation++
	} else {
		e.log_start_of_line = FALSE
	}
}

func (e *Engine) log_cstring(b []byte) {
	for _, ch := range b {
		e.log_char(ch)
	}
}

func (e *Engine) log_int(i int) {
	e.log_printf("%d", i)
}

func (e *Engine) log_long(i int) {
	e.log_printf("%d", i)
}

// log_message copies the contents of a message file to the current log_file
func (e *Engine) log_message(message_filename string) {
	if b, err := ioutil.ReadFile(message_filename); err != nil {
		fprintf(e.stderr, "\n\tWARNING! utils.c: cannot open message file '%s'!\n\n", message_filename)
		log.Printf("[engine] %q: %+v\n", message_filename, err)
	} else {
		fputs(b, e.log_file)
	}
}

func (e *Engine) log_printf(format string, args ...interface{}) {
	e.log_string(fmt.Sprintf(format, args...))
}

func (e *Engine) log_string(s string) {
	for _, ch := range []byte(s) {
		e.log_char(ch)
	}
}
