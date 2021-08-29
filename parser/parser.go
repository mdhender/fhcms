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
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
)

var _parser_debug = false

func Parse(name string) (*Node, error) {
	log.Printf("parse:: parsing orders file %q\n", filepath.Base(name))
	_parser_debug = filepath.Base(name) == "sp18.ordx"
	b, err := loadFile(name)
	if err != nil {
		return nil, err
	}
	return parse(b)
}

// loadFile loads a file and returns it with line endings converted to LF only.
func loadFile(name string) ([]byte, error) {
	const cr, lf = '\r', '\n'
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	// replace all CR+LF, CR, and LF+CR endings with LF.
	return bytes.ReplaceAll(bytes.ReplaceAll(bytes.ReplaceAll(b, []byte{cr, lf}, []byte{lf}), []byte{lf, cr}, []byte{lf}), []byte{cr}, []byte{lf}), nil
}

func parse(b []byte) (*Node, error) {
	s := NewScanner(b)
	tokens, err := s.Scan()
	if err != nil {
		return nil, err
	}
	root := &Node{}
	var section, pc, order *Node
	for _, token := range tokens {
		if _parser_debug {
			log.Println(*token)
		}
		if token.Type == End || token.Type == EOF {
			// terminates the current section, production center, and order
			section, pc, order = nil, nil, nil
			continue
		}
		if token.Type == EOL {
			// terminates the current order
			order = nil
			continue
		}
		if token.Type == StartCombat || token.Type == StartJumps || token.Type == StartPostArrival || token.Type == StartPreDeparture || token.Type == StartProduction || token.Type == StartStrikes {
			// start section terminates the current section, production center, and order
			section, pc, order = nil, nil, nil

			if _parser_debug {
				log.Printf("parser: start: %q\n", token.Abbr)
			}
			// we allow the player to specify the same section multiple times.
			// we implement that by appending orders to an existing section.
			if section = root.firstChildOfType(token.Type); section == nil {
				// if there's no such section, we create it.
				section = &Node{Token: token}
				root.Children = append(root.Children, section)
			}
			continue
		}
		if section == nil { // silently ignore all other tokens if we are not in a section
			continue
		}
		if token.Type == Item || token.Type == Number || token.Type == Orbit || token.Type == Planet || token.Type == PlanetCoordinates || token.Type == Ship || token.Type == Species || token.Type == SystemCoordinates || token.Type == Transport {
			// these are arguments to add to the current pc or order
			if order != nil {
				order.Children = append(order.Children, &Node{Parent: order, Token: token})
			}
			continue
		}
		if order != nil {
			// add as an argument to current order
			order.Children = append(order.Children, &Node{Parent: order, Token: token})
			continue
		}
		if pc != nil { // add to current production center
			order = &Node{Parent: pc, Token: token}
			pc.Children = append(pc.Children, order)
			continue
		}
		if _parser_debug {
			log.Printf("parser: order: %q\n", token.String())
		}
		order = &Node{Parent: section, Token: token}
		section.Children = append(section.Children, order)
		if token.Type == Production {
			pc = order
		}
	}
	if _parser_debug && section != nil {
		log.Printf("parser: %s not terminated\n", section.Token.String())
	}

	return root, nil
}
