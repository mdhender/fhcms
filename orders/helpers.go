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

import (
	"fmt"
	"unicode"
)

func tolower(r rune) rune {
	return unicode.ToLower(r)
}

// split returns a copy of the input split into things that may be of interest.
func split(b []byte) []string {
	//// trim comments
	//if off := bytes.IndexByte(p, ';'); off != 0 {
	//	p = p[:off]
	//}
	//// convert commas to tabs
	//for i := 0; i < len(p); i++ {
	//	if p[i] == ',' {
	//		p[i] = '\t'
	//	}
	//}
	//// compress runs of tabs
	//var b []byte
	//for len(p) != 0 {
	//	if p[0] == '\t' {
	//		if len(b) != 0 && b[len(b)-1] != '\t' {
	//			b = append(b, '\t')
	//		}
	//	} else {
	//		b = append(b, p[0])
	//	}
	//	p = p[1:]
	//}
	//p = b
	//// trim leading and trailing spaces
	//p = bytes.TrimSpace(p)

	var errors []error
	p := position{p: b, line: 1}
	// split into words
	var words []string
	state := "start"
	for !p.eof() {
		if p.skipSpaces() || p.skipComment() {
			continue
		}
		switch state {
		case "start":
			if !p.hasWord([]byte{'s', 't', 'a', 'r', 't'}) {
				p.flushLine() // ignore all lines that are outside a section
				continue
			}
			state = "section"
		case "section":
			for _, v := range []struct {
				section string
				pfx     []byte
			}{
				{"COMBAT", []byte{'c', 'o', 'm', 'b', 'a', 't'}},
				{"PRE-DEPARTURE", []byte{'p', 'r', 'e', '-', 'd', 'e', 'p', 'a', 'r', 't', 'u', 'r', 'e'}},
				{"JUMPS", []byte{'j', 'u', 'm', 'p', 's'}},
				{"PRODUCTION", []byte{'p', 'r', 'o', 'd', 'u', 'c', 't', 'i', 'o', 'n'}},
				{"POST-ARRIVAL", []byte{'p', 'o', 's', 't', '-', 'a', 'r', 'r', 'i', 'v', 'a', 'l'}},
				{"STRIKES", []byte{'s', 't', 'r', 'i', 'k', 'e', 's'}},
			} {
				if !p.hasWord(v.pfx) {
					continue
				}
				state = "eol-before-command-or-end"
				break
			}
		case "eol-before-command-or-end":
			if p.get() != '\n' {
				errors = append(errors, fmt.Errorf("%d: expected eol after start section", p.Line()))
				p.flushLine()
			}
			state = "command-or-end"
		case "command-or-end":
			if p.hasPrefix([]byte{'s', 't', 'a', 'r', 't'}) {
				if unicode.IsSpace(p.peek(5)) {
					p.skip(5)
					state = "section"
					continue
				}
			}
		case "uhm":
			for _, v := range []struct {
				word string
				pfx  []byte
			}{
				{"PL", []byte{'p', 'l'}},
				{"SP", []byte{'s', 'p'}},
			} {
				if p.hasPrefix(v.pfx) {
					if unicode.IsSpace(p.peek(len(v.pfx))) {
						p.skip(len(v.pfx))
						words = append(words, v.word)
						state = "name"
						continue
					}
				}
			}
			for _, v := range []struct {
				word string
				pfx  []byte
			}{
				{"ALLY", []byte{'a', 'l', 'l', 'y'}},
				{"AMBUSH", []byte{'a', 'm', 'b', 'u', 's', 'h'}},
			} {
				if p.hasPrefix(v.pfx) {
					if unicode.IsSpace(p.peek(len(v.pfx))) {
						p.skip(len(v.pfx))
						words = append(words, v.word)
						state = "command"
						continue
					}
				}
			}
		case "command":
		case "name":
		default:
			panic(fmt.Sprintf("unknown state %q", state))
		}
		panic(fmt.Sprintf("unknown %q in state %q", p.slug(), state))
	}

	return words
}

// tabify returns a copy of the input with commas replaced by
// tabs and runs of tabs compressed to a single tab.
func tabify(p []byte) []byte {
	var line []byte
	for len(p) != 0 {
		if p[0] == ',' || p[0] == '\t' {
			if len(line) != 0 && line[len(line)-1] != '\t' {
				line = append(line, '\t')
			}
		} else {
			line = append(line, p[0])
		}
		p = p[1:]
	}
	return line
}
