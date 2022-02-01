/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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

import "bytes"

// Read returns a command and error.
// A command is the command, parameters, and comments.
// Some commands, like START, have children.
// Some commands, like MESSAGE, will consume multiple lines.
func Read(b []byte) error {
	// split input into lines
	lines := bytes.Split(b, []byte{'\n'})
	// trim leading and trailing whitespace from every line
	for i := 0; i < len(lines); i++ {
		lines[i] = bytes.TrimSpace(lines[i])
	}
	return nil
}

// bdup wastes memory and cycles by needlessly making copies of slices.
func bdup(src []byte) (dst []byte) {
	if src == nil {
		return []byte{}
	}
	dst = make([]byte, len(src))
	copy(dst, src)
	return dst
}

// splitLine splits a line into command and comment.
// The command is all characters to the left of the first ';'.
// The comment is all characters after the first ';'.
// Both the command and comment are trimmed of leading and trailing spaces.
// If either value is empty, an empty slice is returned instead of nil.
func splitLine(line []byte) (command []byte, comment []byte) {
	i := bytes.IndexByte(line, ';')
	if i == -1 {
		return bdup(bytes.TrimSpace(line)), bdup([]byte{})
	}
	return bdup(bytes.TrimSpace(line[:i])), bdup(bytes.TrimSpace(line[i+1:]))
}
