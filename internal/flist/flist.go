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

package flist

import (
	"os"
	"strconv"
	"strings"
)

type File struct {
	Turn int
	Type string
	File string
	Date string
}

type TurnData struct {
	Turn   int
	Report string // name of file if it exists
	Orders string // name of file if it exists
	Date   string
}

func Fetch(path string, speciesNo int) ([]*TurnData, error) {
	//root := "/home/mdhender/src/fhcms/cmd/fhapp/testdata/reports"
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var files []*File
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fields := strings.Split(entry.Name(), ".")
		if len(fields) != 4 || !(fields[2] == "orders" || fields[2] == "report") || fields[3] != "txt" {
			continue
		}
		if no, _ := strconv.Atoi(strings.TrimLeft(fields[0], "sp0")); no != speciesNo {
			continue
		}
		turn, _ := strconv.Atoi(strings.TrimLeft(fields[1], "t0"))
		files = append(files, &File{
			Turn: turn,
			Type: fields[2],
			File: entry.Name(),
			Date: info.ModTime().Format("02 Jan 2006 15:04 MST"),
		})
	}

	// sort by turn then file type (with orders before report)
	for i := 0; i < len(files); i++ { // bubbly and proud of it
		for j := i + 1; j < len(files); j++ {
			if files[i].Turn > files[j].Turn {
				files[i], files[j] = files[j], files[i]
			} else if files[i].Turn == files[j].Turn {
				if files[i].Type > files[j].Type {
					files[i], files[j] = files[j], files[i]
				}
			}
		}
	}

	var turns []*TurnData
	var turn *TurnData
	for _, f := range files {
		if turn == nil || turn.Turn != f.Turn {
			turn = &TurnData{
				Turn: f.Turn,
				Date: f.Date,
			}
			turns = append(turns, turn)
		}
		switch f.Type {
		case "orders":
			turn.Orders = f.File
		case "report":
			turn.Report = f.File
		}
	}

	return turns, nil
}
