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

package cmd

import (
	"fmt"
	"github.com/mdhender/fhcms/internal/cluster"
	"log"
	"os"
	"path/filepath"
)

func loader(path string, bigEndian bool) (*cluster.Store, error) {
	if path = filepath.Clean(path); path == "." {
		if cwd, err := os.Getwd(); err != nil {
			return nil, err
		} else if path = filepath.Clean(cwd); path == "." {
			return nil, fmt.Errorf("unable to determine path to data files")
		}
	}
	log.Printf("loader: path      %q\n", path)
	log.Printf("loader: bigEndian %v\n", bigEndian)
	return cluster.FromDat32(path, bigEndian)
}
