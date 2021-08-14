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

func Parse(b []byte) {
	p := &position{p: b, line: 1}
	for !p.eof() {
		if p.skipSpaces() || p.skipComment() {
			continue
		}
		parseSection(p)
	}
}

type Section struct{}

func parseSection(p *position) *Section {
	k := p.accept([]byte{'s', 't', 'a', 'r', 't'})
	if k == nil {
		return nil
	}
	for !p.eof() {
		if p.skipSpaces() || p.skipComment() {
			continue
		}
	}
	if p.accept([]byte{'c', 'o', 'm', 'b', 'a', 't'}) {
	}
	if p.accept([]byte{'p', 'r', 'e', '-', 'd', 'e', 'p', 'a', 'r', 't', 'u', 'r', 'e'}) {
	}
	if p.accept([]byte{'j', 'u', 'm', 'p', 's'}) {
	}
	if p.accept([]byte{'p', 'r', 'o', 'd', 'u', 'c', 't', 'i', 'o', 'n'}) {
	}
	if p.accept([]byte{'p', 'o', 's', 't', '-', 'a', 'r', 'r', 'i', 'v', 'a', 'l'}) {
	}
	if p.accept([]byte{'s', 't', 'r', 'i', 'k', 'e', 's'}) {
	}
	return nil
}
