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
	"unicode"
	"unicode/utf8"
)

const cr, lf = '\r', '\n'

type position struct {
	p    []byte
	line int
}

func (p *position) Line() int {
	return p.line
}

func (p *position) accept(word []byte) []byte {
	if len(word) == 0 {
		return nil
	}
	raw := p.next()
	for tok := raw; !(len(tok) == 0 || len(word) == 0); {
		wr, ww := utf8.DecodeRune(word)
		if wr == utf8.RuneError {
			return nil
		}
		tr, tw := utf8.DecodeRune(tok)
		if !(tr == wr || unicode.ToUpper(tr) == unicode.ToUpper(wr)) {
			return nil
		}
		tok, word = tok[tw:], word[ww:]
	}
	if len(word) != 0 {
		return nil
	}
	p.p = p.p[len(raw):]
	return raw
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

// next returns the next word up to a space, comment, end of line or
// end of input. it returns nil if it starts at end of input. it
// returns an empty slice if it starts on a space, comment, or
// end of line.
func (p *position) next() []byte {
	if p == nil || p.eof() {
		return nil
	}
	n := 0
	for n < len(p.p) {
		r, w := utf8.DecodeRune(p.p[n:])
		if unicode.IsSpace(r) || r == ';' {
			break
		}
		n += w
	}
	return p.p[:n]
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
