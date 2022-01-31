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
	"fmt"
	"log"
)

type FILE struct {
	name string
	data *bytes.Buffer
}

// todo: should this do anything?
func fclose(stream *FILE) {
	// f.data = nil
}

// todo: should this do anything?
func fflush(stream *FILE) {
	// do nothing
}

// fgets() reads in at most one less than size characters from stream and stores them into the buffer pointed to by s.
// Reading stops after an EOF or a newline.
// If a newline is read, it is stored into the buffer.
// A terminating null byte ('\0') is stored after the last character in the buffer.
//
// fgets() returns s on success, and NULL on error or when end of file occurs while no characters have been read.
//
// -- https://linux.die.net/man/3/fgets
func fgets(s []byte, size int, stream *FILE) []byte {
	if !(len(s) == cap(s)) {
		panic("assert(len(s) == cap(s))")
	} else if !(size <= cap(s)) {
		panic("assert(size <= cap(s))")
	}
	for i := 0; i < size; i++ {
		s[i] = 0
	}
	if stream == nil || stream.data == nil || stream.data.Len() == 0 {
		return nil
	}
	size-- // read at most one less than the size passed in
	for i := 0; i < size; i++ {
		if ch, err := stream.data.ReadByte(); err == nil {
			if s[i] = ch; ch == '\n' {
				break
			}
		} else {
			break
		}
	}
	return s
}

func fopen(name string, b *bytes.Buffer) *FILE {
	f := &FILE{name: name, data: b}
	return f
}

// Upon successful return, these functions return the number of characters printed (excluding the null byte used to end output to strings).
//
// -- https://linux.die.net/man/3/fprintf
func fprintf(stream *FILE, format string, args ...interface{}) int {
	if stream == nil || stream.data == nil {
		return 0
	}
	s := fmt.Sprintf(format, args...)
	stream.data.WriteString(s)
	return len(s)
}

// fputs() writes the string s to stream, without its terminating null byte ('\0').
//
// fputs() returns a non-negative number on success, or EOF on error
// -- https://linux.die.net/man/3/fputs
func fputs(s []byte, stream *FILE) int {
	if stream == nil || stream.data == nil || s == nil {
		return 0
	}
	l := strlen(s)
	if l > 0 {
		stream.data.Write(s[:l])
	}
	return l
}

func (f *FILE) bytes() []byte {
	log.Printf("[files] %q is bytes()\n", f.name)
	return f.data.Bytes()
}

func (f *FILE) eof() bool {
	return f == nil || f.data == nil || f.data.Len() == 0
}
