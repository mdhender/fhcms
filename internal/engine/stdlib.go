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

package engine

import (
	"bytes"
)

func isalnum(c byte) bool {
	return isalpha(c) || isdigit(c)
}

func isalpha(c byte) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
}

func isdigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// The strcmp() function compares the two strings s1 and s2. It returns an integer less than, equal to, or greater than zero if s1 is found, respectively, to be less than, to match, or be greater than s2.
//
// -- https://linux.die.net/man/3/strncmp
func strcmp(s1, s2 []byte) int {
	return bytes.Compare(s1[:strlen(s1)], s2[:strlen(s2)])
}

// The strcpy() function copies the string pointed to by src, including the terminating null byte ('\0'), to the buffer pointed to by dest. The strings may not overlap, and the destination string dest must be large enough to receive the copy. Beware of buffer overruns! (See BUGS.)
//
// -- https://linux.die.net/man/3/strcpy
func strcpy(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return src
}

func strlen(s []byte) int {
	if l := bytes.IndexByte(s, 0); l != -1 {
		return l
	}
	return len(s)
}

// The strncmp() function compares the first (at most) n bytes of the two strings s1 and s2. It returns an integer less than, equal to, or greater than zero if s1 is found, respectively, to be less than, to match, or be greater than s2.
//
// -- https://linux.die.net/man/3/strncmp
func strncmp(s1, s2 []byte, n int) int {
	if len(s1) > n {
		s1 = s1[:n]
	}
	if len(s2) > n {
		s2 = s2[:n]
	}
	return bytes.Compare(s1, s2)
}

func toupper(c byte) byte {
	if 'a' <= c && c <= 'z' {
		return ('A' - 1) + c
	}
	return c
}
