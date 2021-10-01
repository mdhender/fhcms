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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/mdhender/fhcms/config"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/mdhender/fhcms/internal/way"
	"io/ioutil"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cfg := config.DefaultConfig()
	err := cfg.Load()
	if err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(2)
	}

	log.SetFlags(cfg.Log.Flags)

	if cfg.PIDFile {
		pid := os.Getpid()
		if err := ioutil.WriteFile("/tmp/.fhapp.pid", []byte(fmt.Sprintf("%d", pid)), 0600); err != nil {
			log.Printf("unable to create pid file: %+v", err)
			os.Exit(2)
		}
		log.Printf("server: pid %8d: file %q\n", pid, "/tmp/.fhapp.pid")
	}

	if errors := run(cfg); errors != nil {
		for _, err := range errors {
			log.Printf("%+v\n", err)
		}
		os.Exit(2)
	}

	os.Exit(0)
}

func run(cfg *config.Config) (errs []error) {
	sessions, err := NewSessionManager(cfg.Data.Sessions, "fhapp")
	if err != nil {
		return append(errs, err)
	}

	s := &Server{
		router:   way.NewRouter(),
		sessions: sessions,
	}
	s.Addr = net.JoinHostPort(cfg.Server.Host, fmt.Sprintf("%d", cfg.Server.Port))
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20 // 1mb?

	// load data
	galaxyDataFile := filepath.Join(cfg.Data.Path, "galaxy.dat")
	starDataFile := filepath.Join(cfg.Data.Path, "stars.dat")
	planetDataFile := filepath.Join(cfg.Data.Path, "planets.dat")
	speciesDataPath := cfg.Data.Path
	locationDataFile := filepath.Join(cfg.Data.Path, "locations.dat")
	transactionDataFile := filepath.Join(cfg.Data.Path, "transactions.dat")

	if ds, err := cluster.FromDat32(galaxyDataFile, starDataFile, planetDataFile, speciesDataPath, locationDataFile, transactionDataFile, cfg.Data.BigEndian); err != nil {
		return append(errs, err)
	} else {
		s.data.Store = ds
	}
	xlatNo := make(map[int]*cluster.Species)
	for _, sp := range s.data.Store.Species {
		xlatNo[sp.No] = sp
	}

	if err = loader(cfg.Data.Players, &s.data.Players); err != nil {
		return append(errs, err)
	}
	for _, p := range s.data.Players {
		if _, ok := s.data.Store.Species[p.SpeciesId]; !ok {
			log.Printf("players: player %q: no such species %q\n", p.User, p.SpeciesId)
			continue
		}
		log.Printf("players: player %q species_id %q password %q key %q\n", p.User, p.SpeciesId, p.Password, p.Key)
	}
	s.data.Files = make(map[string][]*FileData)
	var files []*struct {
		SpeciesId string `json:"-"`
		SpeciesNo int    `json:"species_no"`
		Turn      int    `json:"turn"`
		Type      string `json:"type"`
		File      string `json:"file"`
		Date      string `json:"date"`
	}
	if err = loader(cfg.Data.Files, &files); err != nil {
		return append(errs, err)
	}
	for _, f := range files {
		if sp, ok := xlatNo[f.SpeciesNo]; ok {
			f.SpeciesId = sp.Id
		}
	}
	for _, f := range files {
		var fd *FileData
		for _, j := range s.data.Files[f.SpeciesId] {
			if j.Turn == f.Turn {
				fd = j
				break
			}
		}
		if fd == nil {
			fd = &FileData{SpeciesId: f.SpeciesId, SpeciesNo: f.SpeciesNo, Turn: f.Turn, Date: f.Date}
			s.data.Files[f.SpeciesId] = append(s.data.Files[f.SpeciesId], fd)
		}
		switch f.Type {
		case "report":
			fd.Report = f.File
		case "orders":
			fd.Orders = f.File
		}
	}
	for _, files := range s.data.Files {
		//log.Printf("files: %2d %v\n", sp, files)
		for i := 0; i < len(files); i++ { // bubbly and proud of it
			for j := i + 1; j < len(files); j++ {
				if files[i].LessThan(files[j]) {
					files[i], files[j] = files[j], files[i]
				}
			}
		}
		//for _, f := range files {
		//	log.Printf("files: %2d %v\n", sp, f)
		//}
	}
	s.data.Site = &Site{}
	if err := loader(cfg.Data.Site, s.data.Site); err != nil {
		return append(errs, err)
	}
	s.data.Stats = make(map[string]*StatsData)
	var stats []*StatsData
	if err := loader(cfg.Data.Stats, &stats); err != nil {
		return append(errs, err)
	}
	for _, stat := range stats {
		if sp, ok := xlatNo[stat.SpeciesNo]; ok {
			stat.SpeciesId = sp.Id
			s.data.Stats[stat.SpeciesId] = stat
		}
	}

	// link in some stuff required for managing sessions
	s.sessions.players = s.data.Players
	s.sessions.species = s.data.Store.Species

	s.routes(cfg.Data.Reports, cfg.Data.Reports)

	s.Handler = s.staticFiles("static", s.sessions.SessionUserHandler(s.router))

	log.Printf("serving %s on address %s\n", "public", s.Addr)
	//return http.ListenAndServe(addr, s.sessions.SessionUserHandler(s.router))
	if err := s.ListenAndServe(); err != nil {
		return append(errs, err)
	}

	log.Printf("oddly: server terminated gracefully\n")
	return nil
}

func loader(name string, a interface{}) error {
	log.Printf("loader: loading %s\n", name)
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, a)
}
