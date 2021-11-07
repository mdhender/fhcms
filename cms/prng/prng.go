/*****************************************************************************
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
****************************************************************************/

package prng

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type PRNG struct {
	seed uint64
}

var defaultPRNG = &PRNG{seed: 0xBADC0FFEE}

func New(seed uint64) *PRNG {
	return &PRNG{seed: seed}
}

// Roll returns a random number in the range 1..n.
func Roll(max int) int {
	return defaultPRNG.Roll(max)
}

// Seed the default generator
func Seed(seed string) {
	var hash uint64
	if seed == "0xBADC0FFEE" || seed == "0" || len(seed) == 0 {
		hash = 0xBADC0FFEE
	} else {
		hash = 5381
		for _, ch := range []byte(seed) {
			hash = (hash * 33) ^ uint64(ch)
		}
	}
	defaultPRNG.Seed(hash)
}

// SeedFromTime seeds the generator from the system clock.
func SeedFromTime() {
	defaultPRNG.SeedFromTime()
}

// GetSeed returns the current seed from the default generator.
func (p *PRNG) GetSeed() uint64 {
	return p.seed
}

func (p *PRNG) Roll(max int) int {
	if max < 1 {
		return 1
	}

	// for congruential method, multiply previous value by the prime number 16417.
	congResult := p.seed + (p.seed << 5) + (p.seed << 14) // effectively multiply by 16417

	// for shift-register method, use shift-right 15 and shift-left 17 with no-carry addition (i.e., exclusive-or)
	a := p.seed >> 15
	shiftResult := a ^ p.seed
	a = shiftResult << 17
	shiftResult ^= a

	p.seed = congResult ^ shiftResult

	// avoid returning the low-order bits
	return int(((p.seed&0x0000FFFF)*uint64(max))>>16) + 1
}

// Seed random number generator
func (p *PRNG) Seed(seed uint64) {
	p.seed = seed
	n := p.Roll(100) + p.Roll(200) + p.Roll(300)
	for i := 0; i < n; i++ {
		p.Roll(10)
	}
}

// SeedFromEnv seeds the generator from an environment variable.
func (p *PRNG) SeedFromEnv(name string) error {
	var digits int
	var seed uint64
	for _, ch := range os.Getenv(name) {
		switch ch {
		case '0':
			digits, seed = digits+1, seed*16+0
		case '1':
			digits, seed = digits+1, seed*16+1
		case '2':
			digits, seed = digits+1, seed*16+2
		case '3':
			digits, seed = digits+1, seed*16+3
		case '4':
			digits, seed = digits+1, seed*16+4
		case '5':
			digits, seed = digits+1, seed*16+5
		case '6':
			digits, seed = digits+1, seed*16+6
		case '7':
			digits, seed = digits+1, seed*16+7
		case '8':
			digits, seed = digits+1, seed*16+8
		case '9':
			digits, seed = digits+1, seed*16+9
		case 'a', 'A':
			digits, seed = digits+1, seed*16+10
		case 'b', 'B':
			digits, seed = digits+1, seed*16+11
		case 'c', 'C':
			digits, seed = digits+1, seed*16+12
		case 'd', 'D':
			digits, seed = digits+1, seed*16+13
		case 'e', 'E':
			digits, seed = digits+1, seed*16+14
		case 'f', 'F':
			digits, seed = digits+1, seed*16+15
		}
	}
	if digits == 0 || seed == 0 {
		return fmt.Errorf("invalid seed")
	}
	p.Seed(seed)
	return nil
}

// SeedFromTime seeds the generator from the system clock.
func (p *PRNG) SeedFromTime() {
	rand.Seed(time.Now().UnixNano())
	p.Seed(rand.Uint64())
}
