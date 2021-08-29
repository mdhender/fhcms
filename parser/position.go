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

package parser

import (
	"fmt"
	"unicode"
	"unicode/utf8"
)

type position struct {
	p    []byte
	line int
}

func (p *position) Line() int {
	return p.line
}

func (p *position) clone() *position {
	if p == nil {
		return nil
	}
	return &position{p: p.p, line: p.line}
}

func (p *position) flushLine() {
	for !p.eof() {
		if p.get() == lf {
			return
		}
	}
}

func (p *position) get() rune {
	if p.eof() {
		return utf8.RuneError
	}
	r, w := utf8.DecodeRune(p.p)
	p.p = p.p[w:]
	if r == cr || r == lf {
		p.line++
		if len(p.p) != 0 && ((r == cr && p.p[0] == lf) || (r == lf && p.p[0] == cr)) { // check for CR+LF or LF+CR
			p.p = p.p[1:]
		}
		r = lf
	}
	return r
}

func (p *position) eof() bool {
	return p == nil || len(p.p) == 0
}

func (p *position) hasPrefix(prefix []byte) bool {
	if p.eof() {
		return false
	}
	cp := p.clone()
	for !cp.eof() && len(prefix) != 0 {
		ch := cp.get()
		r, w := utf8.DecodeRune(prefix)
		if r == utf8.RuneError || !(ch == r || unicode.ToUpper(ch) == unicode.ToUpper(r)) {
			return false
		}
		prefix = prefix[w:]
	}
	return len(prefix) == 0
}

func (p *position) hasWord(word []byte) bool {
	if !(p.hasPrefix(word) && unicode.IsSpace(p.peek(len(word)))) {
		return false
	}
	p.skip(len(word))
	return true
}

func (p *position) peek(n int) rune {
	if p.eof() {
		return utf8.RuneError
	}
	cp := p.clone()
	for ; n > 0; n-- {
		cp.get()
	}
	return cp.get()
}

func (p *position) restore(pp *position) {
	if p == nil || pp == nil {
		return
	}
	p.p, p.line = pp.p, pp.line
}

func (p *position) save() *position {
	return p.clone()
}

func (p *position) skip(n int) {
	for ; n > 0; n-- {
		p.get()
	}
}

func (p *position) skipComment() bool {
	if p.peek(0) != ';' {
		return false
	}
	for !(p.eof() || p.peek(0) == lf) {
		p.get()
	}
	return true
}

func (p *position) skipSpaces() bool {
	var skipped int
	for !p.eof() {
		if ch := p.peek(0); ch == lf {
			break
		} else if unicode.IsSpace(ch) {
			skipped++
			p.get()
		} else if ch == ';' {
			p.skipComment()
		} else {
			break
		}
	}
	return skipped != 0
}

func (p *position) slug() string {
	if p == nil {
		return "**nil**"
	} else if p.eof() {
		return "**eof**"
	}
	var s string
	for i := 0; i < 20 && !p.eof(); i++ {
		s += string(p.get())
	}
	return s
}

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
