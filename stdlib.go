/*
 * Far Horizons Engine
 * Copyright (C) 2021  Michael D Henderson
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var stderr io.Writer
var stdout io.Writer

func init() {
	stderr = os.Stderr
	stdout = os.Stdout
}

func atoi(s string) int {
	s = strings.TrimSpace(s)
	var s2 []byte
	for _, b := range strings.TrimSpace(s) {
		if b == '+' && len(s2) == 0 {
			s2 = append(s2, '+')
		} else if b == '-' && len(s2) == 0 {
			s2 = append(s2, '-')
		} else if b == '0' {
			s2 = append(s2, '0')
		} else if b == '1' {
			s2 = append(s2, '1')
		} else if b == '2' {
			s2 = append(s2, '2')
		} else if b == '3' {
			s2 = append(s2, '3')
		} else if b == '4' {
			s2 = append(s2, '4')
		} else if b == '5' {
			s2 = append(s2, '5')
		} else if b == '6' {
			s2 = append(s2, '6')
		} else if b == '7' {
			s2 = append(s2, '7')
		} else if b == '8' {
			s2 = append(s2, '8')
		} else if b == '9' {
			s2 = append(s2, '9')
		} else {
			break
		}
	}
	i, _ := strconv.Atoi(string(s2))
	return i
}

func exit(i int) {
	panic(fmt.Sprintf("exit(%d)", i))
}

func fflush(w io.Writer) {
	// ignore
}

func fopen(filename string, mode string) *os.File {
	switch mode {
	case "a":
		a, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			log.Printf("fopen: %+v", err)
			return nil
		}
		return a
	case "r":
		r, err := os.OpenFile(filename, os.O_RDONLY, 0644)
		if err == nil {
			log.Printf("fopen: %+v", err)
			return nil
		}
		return r
	case "w":
		w, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			log.Printf("fopen: %+v", err)
			return nil
		}
		return w
	}
	panic(fmt.Sprintf("assert(mode != %q)", mode))
}

func fprintf(w io.Writer, format string, args ...interface{}) {
	if w == nil {
		return
	}
	_, _ = fmt.Fprintf(w, format, args...)
}

func fputs(s string, w io.Writer) {
	_, _ = fmt.Print(s)
}

func isalpha(b byte) bool {
	return ('a' <= b && b <= 'z') || ('A' <= b && b <= 'Z')
}

func isdigit(b byte) bool {
	return '0' <= b && b <= '9'
}

func printf(format string, args ...interface{}) {
	fprintf(os.Stdout, format, args...)
}

// strcmp compares s1 and s2.
// returns -1 if s1 < s2, 0 if s1 == s2, 1 otherwise
func strcmp(s1, s2 string) int {
	if s1 < s2 {
		return -1
	} else if s1 == s2 {
		return 0
	}
	return 1
}

// strstr returns index of substr in s.
// returns -1 if substr is not a substring of s
func strstr(s, substr string) int {
	return strings.Index(s, substr)
}

var __lower = []byte("abcdefghijklmnopqrstuvwxyz")
var __upper = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func toupper(b byte) byte {
	if i := bytes.IndexByte(__lower, b); i != -1 {
		b = __upper[i]
	}
	return b
}
