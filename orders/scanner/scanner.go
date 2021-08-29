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

package scanner

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	cr = '\r'
	lf = '\n'
	// eof is returned instead of utf8.RuneError at end of input
	eof rune = -1
)

type Token struct {
	Line   int
	Type   TokenType
	Abbr   string
	Value  string
	Coords struct {
		X, Y, Z, PN int
	}
	Number int
}

func (t Token) String() string {
	if t.Type == PlanetCoordinates {
		return fmt.Sprintf("%d: %s %d %d %d %d", t.Line, t.Type.String(), t.Coords.X, t.Coords.Y, t.Coords.Z, t.Coords.PN)
	} else if t.Type == PlanetCoordinates || t.Type == SystemCoordinates {
		return fmt.Sprintf("%d: %s %d %d %d", t.Line, t.Type.String(), t.Coords.X, t.Coords.Y, t.Coords.Z)
	} else if t.Type == Number {
		return fmt.Sprintf("%d: %s %d", t.Line, t.Type.String(), t.Number)
	}
	if t.Abbr != "" && t.Value != "" {
		return fmt.Sprintf("%d: %s %s %q", t.Line, t.Type.String(), t.Abbr, t.Value)
	} else if t.Abbr != "" {
		return fmt.Sprintf("%d: %s %s", t.Line, t.Type.String(), t.Abbr)
	} else if t.Value != "" {
		return fmt.Sprintf("%d: %s %q", t.Line, t.Type.String(), t.Value)
	}
	return fmt.Sprintf("%d: %s", t.Line, t.Type.String())
}

type Scanner struct {
	slug     string
	source   []byte
	buffer   []byte
	tokens   []*Token
	err      error
	line     int
	keywords map[string]TokenType
}

func NewScanner(src []byte) (*Scanner, error) {
	s := &Scanner{
		line:     1,
		keywords: make(map[string]TokenType),
	}
	s.source = append(s.source, src...)
	s.source = append(s.source, lf)
	for _, command := range []string{
		"ALLY", "AMBUSH", "ATTACK", "AUTO",
		"BASE", "BATTLE", "BUILD",
		"CONTINUE",
		"DEEP", "DEBUG", "DESTROY", "DEVELOP", "DISBAND",
		"END", "ENEMY", "ENGAGE", "ESTIMATE",
		"HAVEN", "HIDE", "HIJACK",
		"IBUILD", "ICONTINUE", "INSTALL", "INTERCEPT",
		"JUMP",
		"LAND",
		"MESSAGE", "MOVE", "NAME",
		"NEUTRAL",
		"ORBIT",
		"PJUMP", "PRODUCTION",
		"RECYCLE", "RENAME", "REPAIR", "RESEARCH",
		"SCAN", "SEND", "SHIPYARD", "SUMMARY", "SURRENDER",
		"TARGET", "TEACH", "TECH", "TELESCOPE", "TERRAFORM", "TRANSFER",
		"UNLOAD", "UPGRADE",
		"VISITED",
		"WITHDRAW", "WORMHOLE",
		"ZZZ",
	} {
		s.keywords[command] = Command
	}
	for _, item := range []string{
		"AU", "CU", "DR",
		"FD", "FJ", "FM", "FS",
		"GT", "GW", "GU1", "GU2", "GU3", "GU4", "GU5", "GU6", "GU7", "GU8", "GU9",
		"IU", "JP", "PD", "RM",
		"SG1", "SG2", "SG3", "SG4", "SG5", "SG6", "SG7", "SG8", "SG9", "SU",
		"TP",
	} {
		s.keywords[item] = Item
	}
	for _, ship := range []string{
		"BA", "BAS", "BC", "BCS", "BS", "BSS", "BM", "BMS", "BR", "BRS", "BW", "BWS",
		"CA", "CAS", "CC", "CCS", "CL", "CLS", "CS", "CSS", "CT", "CTS",
		"DD", "DDS", "DN", "DNS", "ES", "ESS", "FF", "FFS", "PB", "PBS", "SD", "SDS",
	} {
		s.keywords[ship] = Ship
	}
	for _, tech := range []string{"BI", "GV", "LS", "MA", "MI", "ML"} {
		s.keywords[tech] = Tech
	}
	s.keywords["TR"] = Transport  // generic transport; not backwards compatible
	s.keywords["TRS"] = Transport // generic sub-light transport; not backwards compatible
	// yes, this means that we're not allowing transports larger than 9,990,000 tons
	for trClass := 1; trClass < 1000; trClass++ {
		s.keywords[fmt.Sprintf("TR%d", trClass)] = Transport
		s.keywords[fmt.Sprintf("TR%dS", trClass)] = Transport
	}
	return s, nil
}

func (s *Scanner) Scan() ([]*Token, error) {
	s.buffer = s.source
	for !s.isAtEnd() {
		t := s.scanToken()
		if t.Type == EOF {
			break
		} else if t.Type == EOL && len(s.tokens) != 0 && s.tokens[len(s.tokens)-1].Type == EOL {
			// suppress blank lines
		} else if t.Type == Unknown && len(s.tokens) != 0 && s.tokens[len(s.tokens)-1].Type == Unknown {
			// compress unknown tokens to reduce log flames
		} else {
			s.tokens = append(s.tokens, t)
		}
	}
	s.tokens = append(s.tokens, &Token{Type: EOF, Line: s.line})
	return s.tokens, s.err
}

func (s *Scanner) isAtEnd() bool {
	return len(s.buffer) == 0
}

func (s *Scanner) advance() (rune, []byte) {
	if len(s.buffer) == 0 {
		s.slug = "***eof***"
		return eof, nil
	}
	r, w := utf8.DecodeRune(s.buffer)
	b := s.buffer[:w]
	s.buffer = s.buffer[w:]

	s.slug = string(s.buffer)
	if len(s.slug) > 20 {
		s.slug = s.slug[:20]
	}

	// CR, CR+LF, LF, and LF+CR are all accepted as end of line markers.
	// All are translated to LF in the token.
	if r == cr {
		s.line = s.line + 1
		if len(s.buffer) != 0 {
			if s.buffer[0] == lf {
				s.buffer = s.buffer[1:]
			}
		}
		return lf, []byte{lf}
	} else if r == lf {
		s.line = s.line + 1
		if len(s.buffer) != 0 {
			if s.buffer[0] == cr {
				s.buffer = s.buffer[1:]
			}
		}
		return lf, []byte{lf}
	}

	return r, b
}

func (s *Scanner) peek() rune {
	if len(s.buffer) == 0 {
		return eof
	}
	r, _ := utf8.DecodeRune(s.buffer)
	// CR, CR+LF, LF, and LF+CR are all accepted as end of line markers.
	// All are translated to LF in the token.
	if r == cr || r == lf {
		return lf
	}
	return r
}

func (s *Scanner) scanToken() *Token {
	// ignore leading whitespace, comments, and (get this) commas
	for r := s.peek(); r != '\n' && (r == ';' || r == ',' || unicode.IsSpace(r)); r = s.peek() {
		s.advance()   // skip past the space, comma, or semi-colon
		if r == ';' { // read to end of line (but not past it) for comments
			for r = s.peek(); r != eof && r != '\n'; r = s.peek() {
				s.advance()
			}
		}
	}
	if s.isAtEnd() {
		return &Token{Type: EOF, Line: s.line}
	}

	if s.peek() == lf {
		s.advance()
		return &Token{Type: EOL, Line: s.line - 1}
	}

	// check for numbers and coordinates
	if '0' <= s.buffer[0] && s.buffer[0] <= '9' {
		// is word a set of planet or system coordinates?
		if pc := isPlanetCoordinates(s.buffer); pc || isSystemCoordinates(s.buffer) {
			tk := &Token{Line: s.line}
			if pc {
				tk.Type = PlanetCoordinates
			} else {
				tk.Type = SystemCoordinates
			}
			var oc []byte
			for x := s.peek(); '0' <= x && x <= '9'; x = s.peek() {
				tk.Coords.X = tk.Coords.X*10 + int(x-'0')
				_, b := s.advance()
				oc = append(oc, b...)
			}
			oc = append(oc, ' ')
			for s.peek() == ' ' {
				s.advance()
			}
			for y := s.peek(); '0' <= y && y <= '9'; y = s.peek() {
				tk.Coords.Y = tk.Coords.Y*10 + int(y-'0')
				_, b := s.advance()
				oc = append(oc, b...)
			}
			oc = append(oc, ' ')
			for s.peek() == ' ' {
				s.advance()
			}
			for z := s.peek(); '0' <= z && z <= '9'; z = s.peek() {
				tk.Coords.Z = tk.Coords.Z*10 + int(z-'0')
				_, b := s.advance()
				oc = append(oc, b...)
			}
			if pc {
				oc = append(oc, ' ')
				for s.peek() == ' ' {
					s.advance()
				}
				for pn := s.peek(); '0' <= pn && pn <= '9'; pn = s.peek() {
					tk.Coords.PN = tk.Coords.PN*10 + int(pn-'0')
					_, b := s.advance()
					oc = append(oc, b...)
				}
			}
			return tk
		}
		// is word a number
		if isNumber(s.buffer) {
			tk := &Token{Type: Number, Line: s.line}
			var oc []byte
			for r := s.peek(); '0' <= r && r <= '9'; r = s.peek() {
				tk.Number = tk.Number*10 + int(r-'0')
				_, b := s.advance()
				oc = append(oc, b...)
			}
			return tk
		}
	}

	if isSectionStart(s.buffer, "COMBAT", "PRE-DEPARTURE", "JUMPS", "PRODUCTION", "POST-ARRIVAL", "STRIKES") {
		var oc []byte
		// start
		for r := s.peek(); r != ' ' && r != '\t' && r != ',' && r != ';' && r != cr && r != lf && r != eof; r = s.peek() {
			_, b := s.advance()
			oc = append(oc, b...)
		}
		// compress run of space or tabs into a single space
		oc = append(oc, ' ')
		for r := s.peek(); r == ' ' || r == '\t'; r = s.peek() {
			s.advance()
		}
		// section name
		for r := s.peek(); r != ' ' && r != '\t' && r != ',' && r != ';' && r != cr && r != lf && r != eof; r = s.peek() {
			_, b := s.advance()
			oc = append(oc, b...)
		}
		return &Token{Type: Section, Line: s.line, Abbr: strings.ToUpper(strings.TrimSpace(string(oc)))}
	}

	// fetch a word from the input.
	// words are terminated by space, tab, comma, semi-colon, cr, lf, or eof.
	var oc []byte
	for r := s.peek(); r != ' ' && r != '\t' && r != ',' && r != ';' && r != cr && r != lf && r != eof; r = s.peek() {
		_, b := s.advance()
		oc = append(oc, b...)
	}
	// we have to eat comma when it is a delimiter
	if r := s.peek(); r == ',' {
		s.advance()
	}
	word := strings.ToUpper(strings.TrimSpace(string(oc)))

	// is word a message?
	if word == "MESSAGE" {
		// expectation is that the buffer contains MESSAGE Species EOL Line+ ZZZ, but we don't enforce that here
		tk := &Token{Type: Message, Line: s.line}
		oc = nil
		// skip spaces and tabs following the word
		for r := s.peek(); r == ' ' || r == '\t'; r = s.peek() {
			s.advance()
		}
		// collect all the characters up to the first line that starts with ZZZ.
		// there may be spaces before the ZZZ, just to make things harder to parse.
		for !s.isAtEnd() {
			// does the line have a prefix of space* [Zz][Zz][Zz]
			i := 0
			for ; i < len(s.buffer) && (s.buffer[i] == ' ' || s.buffer[i] == '\t'); i++ {
				// skip leading spaces for the moment
			}
			if i+3 < len(s.buffer) && (s.buffer[i] == 'Z' || s.buffer[i] == 'z') && (s.buffer[i+1] == 'Z' || s.buffer[i+1] == 'z') && (s.buffer[i+2] == 'Z' || s.buffer[i+2] == 'z') {
				break
			}
			// does not have the prefix, so collect the entire line and add it to the message buffer
			for r := s.peek(); r != lf && r != eof; r = s.peek() {
				_, b := s.advance()
				oc = append(oc, b...)
			}
			if s.peek() == lf {
				s.advance()
			}
			// append a line feed to the message, even if we found end of file
			oc = append(oc, lf)
		}
		// if we got here, the source was missing the ZZZ line
		if len(oc) == 0 {
			// MESSAGE was followed by nothing all, but we have to have at least one line in our token
			oc = append(oc, lf)
		}
		tk.Value = string(oc)
		return tk
	}

	// is word a planet?
	if word == "PL" {
		name := s.scanName()
		return &Token{Type: Planet, Line: s.line, Abbr: "PL", Value: strings.TrimSpace(string(name))}
	}

	// is word a species?
	if word == "SP" {
		name := s.scanName()
		return &Token{Type: Species, Line: s.line, Abbr: "SP", Value: strings.TrimSpace(string(name))}
	}

	// is word a command, item, ship, tech, or transport?
	if kind, ok := s.keywords[word]; ok {
		switch kind {
		case Ship:
			// skip spaces between word and name
			for s.peek() == ' ' {
				s.advance()
			}
			name := s.scanName()
			return &Token{Type: Ship, Line: s.line, Abbr: word, Value: strings.TrimSpace(string(name))}
		case Transport:
			// skip spaces between word and name
			for s.peek() == ' ' {
				s.advance()
			}
			name := s.scanName()
			return &Token{Type: Transport, Line: s.line, Abbr: word, Value: strings.TrimSpace(string(name))}
		default:
			return &Token{Type: kind, Line: s.line, Abbr: word}
		}
	}

	// on unknown tokens, we must flush the rest of the input line
	for r := s.peek(); !(r == lf || r == eof); r = s.peek() {
		s.advance()
	}
	return &Token{Type: Unknown, Line: s.line, Value: string(oc)}
}

// scanName returns the next name in the source.
// names are terminated by tab, comma, semi-colon, cr, lf, or eof.
// the caller must strip leading and trailing spaces as needed.
func (s *Scanner) scanName() []byte {
	var oc []byte
	for r := s.peek(); r != '\t' && r != ',' && r != ';' && r != cr && r != lf && r != eof; r = s.peek() {
		_, b := s.advance()
		oc = append(oc, b...)
	}
	// we have to eat comma when it is a delimiter
	if r := s.peek(); r == ',' {
		s.advance()
	}
	return oc
}

// isNumber returns true if the first word in the source is an integer
// followed by a valid delimiter (space, tab, comma, semi-colon, cr, lf, or eof).
func isNumber(p []byte) bool {
	if len(p) == 0 {
		return false
	}
	for len(p) != 0 && '0' <= p[0] && p[0] <= '9' {
		p = p[1:]
	}
	return len(p) == 0 || (p[0] == ' ' || p[0] == '\t' || p[0] == ',' || p[0] == ';' || p[0] == cr || p[0] == lf)
}

// isPlanetCoordinates returns true if the first four words in the source
// are integers separated by spaces followed by a valid delimiter (space,
// tab, comma, semi-colon, cr, lf, or eof).
func isPlanetCoordinates(p []byte) bool {
	// get x
	if len(p) == 0 {
		return false
	}
	for len(p) != 0 && '0' <= p[0] && p[0] <= '9' {
		p = p[1:]
	}
	if len(p) == 0 || p[0] != ' ' {
		return false
	}
	for len(p) != 0 && p[0] == ' ' {
		p = p[1:]
	}
	// get y
	if len(p) == 0 {
		return false
	}
	for len(p) != 0 && '0' <= p[0] && p[0] <= '9' {
		p = p[1:]
	}
	if len(p) == 0 || p[0] != ' ' {
		return false
	}
	for len(p) != 0 && p[0] == ' ' {
		p = p[1:]
	}
	// get z
	if len(p) == 0 {
		return false
	}
	for len(p) != 0 && '0' <= p[0] && p[0] <= '9' {
		p = p[1:]
	}
	if len(p) == 0 || p[0] != ' ' {
		return false
	}
	for len(p) != 0 && p[0] == ' ' {
		p = p[1:]
	}
	// get pn
	return isNumber(p)
}

// isSectionStart returns true if the first two words in the
// source are START and a valid section name followed by a
// delimiter (space, tab, comma, semi-colon, end of line, end of input).
func isSectionStart(p []byte, names ...string) bool {
	if len(p) < 6 || strings.ToUpper(string(p[:5])) != "START" {
		return false
	}
	p = p[5:]
	if len(p) == 0 || !(p[0] == ' ' || p[0] == '\t') {
		return false
	}
	for len(p) != 0 && (p[0] == ' ' || p[0] == '\t') {
		p = p[1:]
	}
	var oc []byte
	for len(p) != 0 {
		r, w := utf8.DecodeRune(p)
		if r == ' ' || r == '\t' || r == ',' || r == ';' || r == cr || r == lf {
			break
		}
		oc = append(oc, p[:w]...)
		p = p[w:]
	}
	word := strings.ToUpper(strings.TrimSpace(string(oc)))
	for _, name := range names {
		if word == name {
			return true
		}
	}
	return false
}

// isSystemCoordinates returns true if the first four words in the source are integers separated by spaces
func isSystemCoordinates(p []byte) bool {
	// get x
	if len(p) == 0 {
		return false
	}
	for len(p) != 0 && '0' <= p[0] && p[0] <= '9' {
		p = p[1:]
	}
	if len(p) == 0 || p[0] != ' ' {
		return false
	}
	for len(p) != 0 && p[0] == ' ' {
		p = p[1:]
	}
	// get y
	if len(p) == 0 {
		return false
	}
	for len(p) != 0 && '0' <= p[0] && p[0] <= '9' {
		p = p[1:]
	}
	if len(p) == 0 || p[0] != ' ' {
		return false
	}
	for len(p) != 0 && p[0] == ' ' {
		p = p[1:]
	}
	// get z
	return isNumber(p)
}

// runOf returns a slice containing the bytes in the set and another
// slice containing the rest of the input.
func runOf(p []byte, set ...byte) ([]byte, []byte) {
	if len(p) == 0 {
		return nil, nil
	} else if len(set) == 0 {
		return nil, p
	}
	n := 0
	for n < len(p) && bytes.IndexByte(set, p[n]) != -1 {
		n = n + 1
	}
	if n == 0 {
		return nil, p
	}
	return p[:n], p[n:]
}

// runTo returns a slice containing the bytes up to the first delimiter
// and another slice containing the rest of the input.
func runTo(p []byte, delimiters ...byte) ([]byte, []byte) {
	if len(p) == 0 {
		return nil, nil
	} else if len(delimiters) == 0 {
		return p, nil
	}
	n := 0
	for n < len(p) && bytes.IndexByte(delimiters, p[n]) == -1 {
		n = n + 1
	}
	if n == 0 {
		return nil, p
	}
	return p[:n], p[n:]
}
