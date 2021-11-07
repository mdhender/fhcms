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

type Node struct {
	Parent   *Node
	Children []*Node
	Token    *Token
}

func (n *Node) firstChildOfType(t TokenType) *Node {
	if n != nil {
		for _, c := range n.Children {
			if c.Token.Type == t {
				return c
			}
		}
	}
	return nil
}

func (n *Node) String() string {
	return n.prettyPrint(1)
}

func (n *Node) prettyPrint(l int) string {
	var s string
	for i := 0; i < l; i++ {
		s += " "
	}
	if n.Token != nil {
		s += n.Token.String()
	}
	if n.Children != nil {
		s += "\n"
		for _, c := range n.Children {
			s += c.prettyPrint(l + 3)
		}
	}
	return s + "\n"
}
