/*******************************************************************************
fhapp - a player aid for Far Horizons
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

package jdb

import (
	"encoding/json"
	"github.com/mdhender/fhcms/internal/cluster"
	"io/ioutil"
	"path/filepath"
)

type JDB struct {
	Semver  string
	Cluster *cluster.Cluster
}

func Load(data string) (*JDB, error) {
	ds := &JDB{}
	ds.Semver = "7.5.2"
	ds.Cluster = &cluster.Cluster{}
	if err := loader(filepath.Join(data, "cluster.json"), ds.Cluster); err != nil {
		return nil, err
	}
	return ds, nil
}

func loader(name string, a interface{}) error {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, a)
}
