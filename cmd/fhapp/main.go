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
	"flag"
	"fmt"
	"github.com/mdhender/fhcms/cmd/fhapp/internal/way"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	var addr = flag.String("addr", ":8080", "the address of the application")
	var data = flag.String("data", "D:\\FarHorizons\\testdata\\t17", "the path to the application data")
	var pidFile = flag.String("pid-file", "", "create pid file")
	var players = flag.String("players", "D:\\GoLand\\fhcms\\cmd\\fhapp\\testdata\\players.json", "name of file containing player data")
	flag.Parse()

	log.SetFlags(log.Ldate | log.Ltime | log.LUTC) // force logs to be UTC

	if len(*pidFile) != 0 {
		if err := ioutil.WriteFile(*pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0600); err != nil {
			log.Printf("unable to create pid file: %w", err)
			os.Exit(2)
		}
		log.Printf("created pid file %q\n", *pidFile)
	}

	if err := run(*addr, *data, *players); err != nil {
		log.Fatalf("%+v\n", err)
	}
}

func run(addr, data, players string) error {
	var err error

	s := &Server{
		router:   way.NewRouter(),
		sessions: NewSessionManager("fhapp"),
	}
	s.Addr = addr
	s.ReadTimeout = 5 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxHeaderBytes = 1 << 20 // 1mb?

	// load data
	s.data.Engine = &Engine{Semver: "7.5.2"}
	s.data.DS = &JDB{}
	if err = loader(filepath.Join(data, "galaxy.json"), s.data.DS); err != nil {
		return err
	}
	for _, system := range s.data.DS.Stars {
		for pn := 0; pn < system.NumPlanets; pn++ {
			planet := s.data.DS.Planets[pn+system.PlanetIndex]
			system.Planets = append(system.Planets, planet)
			planet.System = system
		}
	}
	for _, species := range s.data.DS.Species {
		species.Code = fmt.Sprintf("SP%02d", species.Id)
		for pn, nampla := range species.Namplas {
			planet := s.data.DS.Planets[nampla.PlanetIndex]
			if pn == 0 { // by convention, PN 0 is the home colony
				species.HomeSystem = planet.System
				species.HomePlanet = nampla
			}
			nampla.System = planet.System
			nampla.Planet = planet
		}
		for _, ship := range species.Ships {
			if ship.Pn == 99 || ship.Name == "Unused" {
				continue
			}
			if ship.DestX != 0 && ship.DestY != 0 && ship.DestZ != 0 {
				log.Printf("init: ship dest %d %d %d\n", ship.X, ship.Y, ship.Z)
			}
			switch ship.Class {
			case BA: /* Starbase. */
				ship.Code = "BA " + ship.Name
				species.Starbases = append(species.Starbases, ship)
			case BC: /* Battlecruiser. */
				ship.Code = "BC " + ship.Name
				species.Warships = append(species.Warships, ship)
			case BM: /* Battlemoon. */
				ship.Code = "BM " + ship.Name
				species.Warships = append(species.Warships, ship)
			case BR: /* Battlestar. */
				ship.Code = "BR " + ship.Name
				species.Warships = append(species.Warships, ship)
			case BS: /* Battleship. */
				ship.Code = "BS " + ship.Name
				species.Warships = append(species.Warships, ship)
			case BW: /* Battleworld. */
				ship.Code = "BW " + ship.Name
				species.Warships = append(species.Warships, ship)
			case CA: /* Heavy Cruiser. */
				ship.Code = "CA " + ship.Name
				species.Warships = append(species.Warships, ship)
			case CL: /* Light Cruiser. */
				ship.Code = "CL " + ship.Name
				species.Warships = append(species.Warships, ship)
			case CC: /* Command Cruiser. */
				ship.Code = "CC " + ship.Name
				species.Warships = append(species.Warships, ship)
			case CS: /* Strike Cruiser. */
				ship.Code = "CS " + ship.Name
				species.Warships = append(species.Warships, ship)
			case CT: /* Corvette. */
				ship.Code = "CT " + ship.Name
				species.Warships = append(species.Warships, ship)
			case DD: /* Destroyer. */
				ship.Code = "DD " + ship.Name
				species.Warships = append(species.Warships, ship)
			case DN: /* Dreadnought. */
				ship.Code = "DN " + ship.Name
				species.Warships = append(species.Warships, ship)
			case ES: /* Escort. */
				ship.Code = "ES " + ship.Name
				species.Warships = append(species.Warships, ship)
			case FG: /* Frigate. */
				ship.Code = "FG " + ship.Name
				species.Warships = append(species.Warships, ship)
			case PB: /* Picketboat. */
				ship.Code = "PB " + ship.Name
				species.Warships = append(species.Warships, ship)
			case SD: /* Super Dreadnought. */
				ship.Code = "SD " + ship.Name
				species.Warships = append(species.Warships, ship)
			case TR: /* Transport. */
				ship.Code = "TR " + ship.Name
				species.Transports = append(species.Transports, ship)
			default:
				// should panic?
				log.Printf("init: unknown ship class %d\n", ship.Class)
				ship.Code = "?? " + ship.Name
				species.Transports = append(species.Transports, ship)
			}
		}
	}
	if err = loader(players, &s.data.Players); err != nil {
		return err
	}
	for _, p := range s.data.Players {
		log.Printf("players: %v\n", *p)
		if 0 < p.Species && p.Species <= len(s.data.DS.Species) {
			s.data.DS.Species[p.Species-1].Key = p.Key
		}
	}
	s.data.Files = make(map[int][]*FileData)
	var files []*struct {
		SpeciesId int    `json:"species"`
		Turn      int    `json:"turn"`
		Type      string `json:"type"`
		File      string `json:"file"`
		Date      string `json:"date"`
	}
	if err = loader(filepath.Join(data, "files.json"), &files); err != nil {
		return err
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
			fd = &FileData{Species: f.SpeciesId, Turn: f.Turn, Date: f.Date}
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
				if files[i].Turn < files[j].Turn {
					files[i], files[j] = files[j], files[i]
				}
			}
		}
		//for _, f := range files {
		//	log.Printf("files: %2d %v\n", sp, f)
		//}
	}
	s.data.Site = &Site{}
	if err = loader(filepath.Join(data, "site.json"), s.data.Site); err != nil {
		return err
	}
	s.data.Stats = make(map[int]*StatsData)
	var stats []*StatsData
	if err = loader(filepath.Join(data, "stats.json"), &stats); err != nil {
		return err
	}
	for _, stat := range stats {
		s.data.Stats[stat.Species] = stat
	}

	// link in some stuff required for managing sessions
	s.sessions.players = s.data.Players
	s.sessions.species = s.data.DS.Species
	if err = loader(filepath.Join(data, "sessions.json"), s.sessions); err != nil {
		return err
	}
	for id, sess := range s.sessions.sessions {
		log.Printf("sessions: %q %v\n", id, sess)
	}

	s.routes()

	s.Handler = s.sessions.SessionUserHandler(s.router)

	log.Printf("serving %s on address %s\n", "public", s.Addr)
	//return http.ListenAndServe(addr, s.sessions.SessionUserHandler(s.router))
	return s.ListenAndServe()
}

func loader(name string, a interface{}) error {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, a)
}
