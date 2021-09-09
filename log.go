/*****************************************************************************
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
****************************************************************************/

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

/* The following routines will post an item to standard output and to
 * an externally defined log file and summary file. */

// log_char can be buggy with long lines or with break before indentation offset.
func log_char(c byte) {
	if logging_disabled {
		return
	}

	/* Check if current line is getting too long. */
	if (c == ' ' || c == '\n') && log_position > 77 {
		// find closest preceding space
		temp_position := bytes.LastIndexByte(log_line[:log_position], ' ')
		if temp_position == -1 {
			panic("assert(temp_position != -1)")
		} else if temp_position == 0 {
			panic("assert(temp_position != 0)")
		}

		// write front of line to files
		s := string(log_line[:temp_position]) + "\n"
		if log_to_file {
			fputs(s, log_file)
		}
		if log_stdout {
			fputs(s, stdout)
		}
		if log_summary {
			fputs(s, summary_file)
		}

		// copy overflow word to beginning of next line
		tmp := log_line[temp_position:log_position]
		log_position = 0
		for log_position = 0; log_position < log_indentation+1; log_position++ {
			log_line[log_position] = ' '
		}
		for i := 0; i < len(tmp); i++ {
			log_line[log_position] = tmp[i]
			log_position++
		}

		if c == ' ' {
			log_line[log_position] = ' '
			log_position++
			return
		}
	}

	// check if line is being manually terminated
	if c == '\n' {
		// write current line to output
		s := string(log_line[:log_position]) + "\n"
		if log_to_file {
			fputs(s, log_file)
		}
		if log_stdout {
			fputs(s, stdout)
		}
		if log_summary {
			fputs(s, summary_file)
		}

		/* Set up for next line. */
		log_position = 0
		log_indentation = 0
		log_start_of_line = true

		return
	}

	// save this character
	log_line[log_position] = c
	log_position++
	log_start_of_line = log_start_of_line && c == ' '
	if log_start_of_line {
		// determine number of indenting spaces for current line.
		log_indentation++
	}
}

func log_int(value int) {
	if logging_disabled {
		return
	}
	log_string(fmt.Sprintf("%d", value))
}

func log_long(value int) {
	if logging_disabled {
		return
	}
	log_string(fmt.Sprintf("%d", value))
}

func log_message(message_filename string) {
	buf, err := ioutil.ReadFile(message_filename)
	if err != nil {
		panic(err)
	}
	fputs(string(buf), log_file)
}

func log_string(s string) {
	if logging_disabled {
		return
	}
	for i := 0; i < len(s); i++ {
		log_char(s[i])
	}
}

func (g *globals) log(format string, a ...interface{}) {
	if g.logging_disabled {
		return
	}
	if g.log_stdout && g.stdout != nil {
		_, _ = fmt.Fprintf(g.stdout, format, a...)
	}
	if g.log_summary && g.summary_file != nil {
		_, _ = fmt.Fprintf(g.summary_file, format, a...)
	}
	if g.log_to_file && g.log_file != nil {
		_, _ = fmt.Fprintf(g.log_file, format, a...)
	}
}
