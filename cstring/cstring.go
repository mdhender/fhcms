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

package cstring

import "unicode/utf8"

type CString struct {
	b   []byte // the buffer
	pos *Position
}

type Position struct {
	p    []byte
	line int
}

func (cs *CString) IsEOF() bool {
	return cs == nil || cs.pos.IsEOF()
}

func (cs *CString) Get() rune {
	if cs == nil {
		return utf8.RuneError
	}
	return cs.pos.Get()
}

func (cs *CString) Peek() rune {
	if cs == nil {
		return utf8.RuneError
	}
	return cs.pos.Peek()
}

func (cs *CString) Restore(p *Position) {
	if p == nil {
		return
	}

}

func (p *Position) Get() rune {
	r, w := utf8.DecodeRune(p.p)
	p.p = p.p[w:]
	return r
}

func (p *Position) IsEOF() bool {
	return p == nil || len(p.p) == 0
}

func (p *Position) Peek() rune {
	if p == nil || len(p.p) == 0 {
		return utf8.RuneError
	}
	r, _ := utf8.DecodeRune(p.p)
	return r
}
