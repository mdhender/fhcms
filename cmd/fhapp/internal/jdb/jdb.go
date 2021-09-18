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
	"fmt"
	"github.com/mdhender/fhcms/internal/cluster"
	"github.com/mdhender/fhcms/internal/coords"
	"io/ioutil"
	"log"
	"path/filepath"
)

func Load(path string) (*Store, error) {
	ds := &Store{
		Semver:  "7.5.2",
		Planets: make(map[string]*Planet),
		Species: make(map[string]*Species),
		Systems: make(map[string]*System),
	}
	var data cluster.Cluster
	if err := loader(filepath.Join(path, "cluster.json"), &data); err != nil {
		return nil, err
	}

	// add all the systems to the cluster map.
	ds.Systems = make(map[string]*System)
	for id, system := range data.Systems {
		o := &System{
			Id:        id,
			Coords:    &coords.Coords{X: system.Coords.X, Y: system.Coords.Y, Z: system.Coords.Z},
			Message:   system.Message,
			VisitedBy: make(map[string]*Species),
		}
		o.Color.DisplayCode = system.Color.DisplayCode
		o.Color.Description = system.Color.Description
		o.Type.DisplayCode = system.Type.DisplayCode
		o.Type.Description = system.Type.Description
		ds.Systems[id] = o
	}

	// add all the planets to the cluster map.
	ds.Planets = make(map[string]*Planet)
	for id, planet := range data.Planets {
		system, ok := ds.Systems[planet.System]
		assert(ok, "planet %q: system %q does not exist", id, planet.System)
		o := &Planet{
			Id:                       id,
			Coords:                   &coords.Coords{X: planet.Coords.X, Y: planet.Coords.Y, Z: planet.Coords.Z, Orbit: planet.Coords.Orbit},
			System:                   system,
			EconEfficiency:           planet.EconEfficiency,
			Gravity:                  planet.Gravity,
			Message:                  planet.Message,
			MiningDifficulty:         planet.MiningDifficulty,
			MiningDifficultyIncrease: planet.MiningDifficultyIncrease,
			PressureClass:            planet.PressureClass,
			Special:                  planet.Special,
			TemperatureClass:         planet.TemperatureClass,
		}
		ds.Planets[id] = o
	}

	// link systems and planets
	for systemId, system := range data.Systems {
		s, ok := ds.Systems[systemId]
		assert(ok, "system %q does not exist in system map", systemId)
		for _, planetId := range system.Planets {
			p, ok := ds.Planets[planetId]
			assert(ok, "system %q has planet %q which does not exist in system map", systemId, planetId)
			s.Planets = append(s.Planets, p)
		}
	}

	// add all the species to the cluster map.
	for speciesId, species := range data.Species {
		sp := &Species{
			Id:             speciesId,
			Ally:           make(map[string]*Species),
			Contact:        make(map[string]*Species),
			EconUnits:      species.EconUnits,
			Enemy:          make(map[string]*Species),
			HPOriginalBase: species.HPOriginalBase,
			Scanned:        make(map[string]*System),
			Visited:        make(map[string]*System),
		}
		ds.Species[speciesId] = sp

		// link home world to planet
		homeworld, ok := ds.Planets[species.HomeWorld]
		assert(ok, "species %q has home world %q which does not exist in system map", speciesId, species.HomeWorld)
		sp.HomeWorld = homeworld

		sp.BI = &Technology{Code: "BI", Level: species.Tech.BI.Level, KnowledgeLevel: species.Tech.BI.KnowledgeLevel, ExperiencePoints: species.Tech.BI.ExperiencePoints}
		sp.GV = &Technology{Code: "GV", Level: species.Tech.GV.Level, KnowledgeLevel: species.Tech.GV.KnowledgeLevel, ExperiencePoints: species.Tech.GV.ExperiencePoints}
		sp.LS = &Technology{Code: "LS", Level: species.Tech.LS.Level, KnowledgeLevel: species.Tech.LS.KnowledgeLevel, ExperiencePoints: species.Tech.LS.ExperiencePoints}
		sp.MA = &Technology{Code: "MA", Level: species.Tech.MA.Level, KnowledgeLevel: species.Tech.MA.KnowledgeLevel, ExperiencePoints: species.Tech.MA.ExperiencePoints}
		sp.MI = &Technology{Code: "MI", Level: species.Tech.MI.Level, KnowledgeLevel: species.Tech.MI.KnowledgeLevel, ExperiencePoints: species.Tech.MI.ExperiencePoints}
		sp.ML = &Technology{Code: "ML", Level: species.Tech.ML.Level, KnowledgeLevel: species.Tech.ML.KnowledgeLevel, ExperiencePoints: species.Tech.ML.ExperiencePoints}
		sp.Fleet.Cost = species.Fleet.Cost
		sp.Fleet.MaintenancePct = species.Fleet.MaintenancePct
		sp.Fleet.Ships = make(map[string]*Ship)
		sp.Gases.Required.Gas = &GasType{Code: species.Gases.Required.Code}
		sp.Gases.Required.MaxPct = species.Gases.Required.MaxPct
		sp.Gases.Required.MinPct = species.Gases.Required.MinPct
		for _, gas := range species.Gases.Neutral {
			sp.Gases.Neutral = append(sp.Gases.Neutral, &GasType{Code: gas})
		}
		for _, gas := range species.Gases.Poison {
			sp.Gases.Poison = append(sp.Gases.Poison, &GasType{Code: gas})
		}
		sp.Government.Name = species.Government.Name
		sp.Government.Type = species.Government.Type
		sp.Named.Systems = make(map[string]*NamedSystem)
		sp.Named.Planets = make(map[string]*NamedPlanet)

		// link named planets to planets
		for nid, npl := range species.NamedPlanets {
			p, ok := ds.Planets[npl.Planet]
			assert(ok, "species %q has named planet %q which does not exist in systems map\n", speciesId, npl.Planet)
			namedPlanet := &NamedPlanet{Id: nid, DisplayName: npl.DisplayName, Planet: p}
			sp.Named.Planets[namedPlanet.Id] = namedPlanet
		}

		// link colonies to named and host planet
		for _, colonyId := range species.Colonies {
			p, ok := ds.Planets[colonyId]
			if !ok {
				log.Printf("warn: species %q: colony %q does not exist in systems map\n", speciesId, colonyId)
				continue
			}
			planet, ok := data.Planets[p.Id]
			if !ok {
				log.Printf("warn: species %q: colony %q does not exist on planet %q\n", speciesId, colonyId, p.Id)
				continue
			}
			pc, ok := planet.Colonies[species.Id]
			if !ok {
				log.Printf("warn: species %q: colony %q does not exist on planet %q\n", speciesId, colonyId, p.Id)
				continue
			}
			colony := &Colony{System: p.System, Planet: p}
			for _, npl := range sp.Named.Planets {
				if p == npl.Planet {
					colony.Name = npl
				}
			}
			if colony.Name == nil {
				log.Printf("warn: species %q: colony %q does not have a named planet\n", speciesId, colonyId)
				continue
			}

			colony.Colonists.AvailableToInstall = pc.Units.Colonists.AvailableToInstall
			colony.Colonists.Auto = pc.Units.Colonists.Auto
			colony.Colonists.Needed = pc.Units.Colonists.Needed
			colony.Inventory = make(map[string]int)
			for code, qty := range pc.Inventory {
				colony.Inventory[code] = qty
			}
			colony.Is.Colony = pc.Status.Colony
			colony.Is.DisbandedColony = pc.Status.DisbandedColony
			colony.Is.Hiding = pc.Status.Hiding
			colony.Is.Hidden = pc.Status.Hidden
			colony.Is.HomePlanet = pc.Status.HomePlanet
			colony.Is.MiningColony = pc.Status.MiningColony
			colony.Is.Populated = pc.Status.Populated
			colony.Is.ResortColony = pc.Status.ResortColony
			colony.Manufacturing.AvailableToInstall = pc.Units.Manufacturing.AvailableToInstall
			colony.Manufacturing.Auto = pc.Units.Manufacturing.Auto
			colony.Manufacturing.Base = pc.MaBase
			colony.Manufacturing.Needed = pc.Units.Manufacturing.Needed
			colony.Message = pc.Message
			colony.Mining.AvailableToInstall = pc.Units.Mining.AvailableToInstall
			colony.Mining.Auto = pc.Units.Mining.Auto
			colony.Mining.Base = pc.MiBase
			colony.Mining.Needed = pc.Units.Mining.Needed
			colony.Population = pc.PopUnits
			colony.SiegeEff = pc.SiegeEff
			colony.Shipyards = pc.Shipyards
			colony.Special = pc.Special
			colony.UseOnAmbush = pc.UseOnAmbush

			sp.Colonies = append(sp.Colonies, colony)
		}

		// link ships
		for shipId, ship := range species.Ships {
			sh := &Ship{
				Id:            shipId,
				Age:           ship.Age,
				DisplayName:   ship.DisplayName,
				Inventory:     make(map[string]int),
				RemainingCost: ship.RemainingCost,
			}
			if ship.Destination != nil {
				sh.Destination = &coords.Coords{X: ship.Destination.X, Y: ship.Destination.Y, Z: ship.Destination.Z, Orbit: ship.Destination.Orbit}
			}
			for code, qty := range ship.Inventory {
				sh.Inventory[code] = qty
			}
			if ship.Location != nil {
				sh.Location = &coords.Coords{X: ship.Location.X, Y: ship.Location.Y, Z: ship.Location.Z, Orbit: ship.Location.Orbit}
			}
			sh.Status.ForcedJump = ship.Status.ForcedJump
			sh.Status.InDeepSpace = ship.Status.InDeepSpace
			sh.Status.InOrbit = ship.Status.InOrbit
			sh.Status.JumpedInCombat = ship.Status.JumpedInCombat
			sh.Status.OnSurface = ship.Status.OnSurface
			sh.Status.UnderConstruction = ship.Status.UnderConstruction

			sp.Fleet.Ships[sh.Id] = sh
		}

		// link species to allies
		for _, allyId := range species.Ally {
			if ally, ok := ds.Species[allyId]; ok {
				sp.Ally[allyId] = ally
			}
		}

		// link species to contacts
		for _, contactId := range species.Contact {
			if contact, ok := ds.Species[contactId]; ok {
				sp.Contact[contactId] = contact
			}
		}

		// link species to enemies
		for _, enemyId := range species.Enemy {
			if enemy, ok := ds.Species[enemyId]; ok {
				sp.Enemy[enemyId] = enemy
			}
		}

		// link species to systems scanned
		for _, systemId := range species.Scanned {
			s, ok := ds.Systems[systemId]
			if !ok {
				log.Printf("warn: species %q scanned %q which does not exist in system map", speciesId, systemId)
				continue
			}
			sp.Scanned[s.Id] = s
		}

		// link species to systems visited
		for _, systemId := range species.Visited {
			s, ok := ds.Systems[systemId]
			if !ok {
				log.Printf("warn: species %q visited %q which does not exist in system map", speciesId, systemId)
				continue
			}
			s.VisitedBy[sp.Id] = sp
			sp.Visited[s.Id] = s
		}
	}

	return ds, nil
}

func assert(t bool, format string, args ...interface{}) {
	if t {
		return
	}
	panic(fmt.Sprintf(format, args...))
}

func loader(name string, a interface{}) error {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, a)
}

type Store struct {
	Semver  string
	Turn    int // current turn number
	Radius  int // radius of the cluster
	Systems map[string]*System
	Planets map[string]*Planet
	Species map[string]*Species
}
