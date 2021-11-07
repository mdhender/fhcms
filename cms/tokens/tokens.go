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

import (
	"bytes"
	"strings"
	"unicode/utf8"
)

type Factory struct {
	lines [][]byte
}

type Token struct {
	Kind  Kind
	Line  int
	Value struct {
		Raw         []byte
		Coordinates struct {
			X, Y, Z, PN int
		}
		Abbreviation string
		Number       string
		Text         string
	}
}

func NewFactory(p []byte) *Factory {
	f := &Factory{}
	for line, rest := nextline(p, false); len(rest) != 0; line, rest = nextline(rest, false) {
		words := strings.Fields(strings.TrimSpace(string(line)))
		if len(words) != 0 {
			command := strings.ToLower(words[0])
			switch command {
			case "start":
			default:
			}
		}
		f.lines = append(f.lines, line)
	}
	return &Factory{lines: split(scrub(p))}
}

func (f *Factory) tokenize(i int) []*Token {
	if i < 0 || i > len(f.lines) {
		return nil
	}
	line := f.lines[i]
	words := strings.Fields(strings.TrimSpace(string(line)))
	if len(words) == 0 {
		return nil
	}
	command := strings.ToLower(words[0])
	switch command {
	case "start":
		return []*Token{&Token{Line: i + 1, Kind: Start}}
	default:
	}
}

// nextline returns the next line from the buffer along with the
// remainder of the input. for convenience, it replaces each
// invalid UTF-8 sequences with a single '?'.
//
// if the trim flag is true, it will trim leading spaces as well
// as comments and trailing spaces.
//
// accepts CR, CR+LF, LF, and LF+CR as valid line endings.
//
// the line returned will be nil only if the buffer is empty.
func nextline(p []byte, trim bool) ([]byte, []byte) {
	const cr, lf = '\r', '\n'
	var line []byte
	for len(p) != 0 {
		if ch := p[0]; ch == cr || ch == lf {
			p = p[1:]
			// check for CR+LF or LF+CR
			if len(p) != 0 && ((ch == cr && p[0] == lf) || (ch == lf && p[0] == cr)) {
				p = p[1:]
			}
			break
		}
		r, w := utf8.DecodeRune(p)
		if r == utf8.RuneError {
			line = append(line, '?')
		} else {
			line = append(line, p[:w]...)
		}
		p = p[w:]
	}
	if trim {
		if off := bytes.IndexByte(p, ';'); off != 0 {
			p = p[:off]
		}
		line = bytes.TrimSpace(line)
	}
	return line, p
}

// scrub returns a copy of the input all each invalid utf-8 sequence
// replaced with a single '?' character. We could have used another
// representation like Go's RuneError or '⯑' but the order processor
// could use the representation on a naming command and that could
// be confusing for the players to use later on. The question mark
// is ugly, but much easier for players to understand.
func scrub(p []byte) []byte {
	cp := make([]byte, len(p))
	for len(p) != 0 {
		r, w := utf8.DecodeRune(p)
		if r == utf8.RuneError {
			cp = append(cp, '?')
		} else {
			cp = append(cp, p[:w]...)
		}
		p = p[w:]
	}
	return cp
}

// split returns a copy of the input as a slice of lines accepting
// CR, CR+LF, LF, and LF+CR. None of the lines in the slice returned will
// be nil, but they may be empty slices.
func split(p []byte) (lines [][]byte) {
	const cr, lf = '\r', '\n'
	line := []byte{}
	for len(p) != 0 {
		if p[0] == cr {
			p, lines, line = p[1:], append(lines, bytes.TrimSpace(line)), []byte{}
			if len(p) != 0 && p[0] == lf { // check for CR+LF
				p = p[1:]
			}
		} else if p[0] == lf {
			p, lines, line = p[1:], append(lines, bytes.TrimSpace(line)), []byte{}
			if len(p) != 0 && p[0] == cr { // check for LF+CR
				p = p[1:]
			}
		} else {
			p, line = p[1:], append(line, p[0])
		}
	}
	return append(lines, bytes.TrimSpace(line))
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

// trim returns a slice from the input with comments removed.
func trimComments(p []byte) []byte {
	if off := bytes.IndexByte(p, ';'); off != 0 {
		p = p[:off]
	}
	return p
}
