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

package agrep

// Score returns a numeric score indicating how closely two strings match.
// If the score is exactly 10000, then the strings are identical.
// Otherwise, the value returned is the number of character matches,
// allowing for accidental transpositions, insertions, and deletions.
// Excess characters in either string will subtract from the score.
// Thus, it's possible for a score to be negative.
//
// In general, if the strings are at least 7 characters each, then you can assume
// the strings are the same if the highest score equals the length of the correct
// string, length-1, or length-2, AND if the score of the next best match is less
// than the highest score.
// A non-10000 score will never be higher than the length of the correct string.
//
// We assume that the caller has forced both strings to the same case.
func Score(correct, unknown string) int {
	if correct == unknown {
		return 10000
	}

	score := 0
	p1, p2 := []rune(correct), []rune(unknown)

	for len(p1) != 0 && len(p2) != 0 {
		c1, c2 := p1[0], p2[0]
		p1, p2 = p1[1:], p2[1:]
		if c1 == c2 {
			score++
		} else if (len(p1) != 0 && c2 == p1[0]) && (len(p2) != 0 && c1 == p2[0]) {
			// transposed
			score += 2
			p1, p2 = p1[1:], p2[1:]
		} else if len(p2) != 0 && c1 == p2[0] {
			// unneeded character
			score++
			p2 = p2[1:]
		} else if len(p1) != 0 && c2 == p1[0] {
			// missing character
			score++
			p1 = p1[1:]
		}
	}

	// reduce score by excess characters, if any
	score -= len(p1)
	score -= len(p2)

	return score
}
