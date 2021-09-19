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

package cluster

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/mdhender/fhcms/internal/coords"
	"github.com/mdhender/fhcms/internal/dat32"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func FromDat32(galaxyDataFile, starDataFile, planetDataFile string, speciesDataPath, locationDataFile, transactionDataFile string, bigendian bool) (*Store, error) {
	var bo binary.ByteOrder
	if bigendian {
		bo = binary.BigEndian
	} else {
		bo = binary.LittleEndian
	}

	galaxyData, err := dat32.ReadGalaxy(galaxyDataFile, bo)
	if err != nil {
		return nil, err
	}

	starData, err := dat32.ReadStars(starDataFile, bo)
	if err != nil {
		return nil, err
	}

	planetData, err := dat32.ReadPlanets(planetDataFile, bo)
	if err != nil {
		return nil, err
	}

	var speciesData []*dat32.Species
	for speciesNo := 1; speciesNo <= galaxyData.NumSpecies; speciesNo++ {
		speciesDataFile := filepath.Join(speciesDataPath, fmt.Sprintf("sp%02d.dat", speciesNo))
		sp, err := dat32.ReadSpecies(speciesDataFile, speciesNo, bo)
		if err != nil {
			return nil, err
		}
		speciesData = append(speciesData, sp)
	}

	ds := &Store{
		Systems: make(map[string]*System),
		Planets: make(map[string]*Planet),
		Species: make(map[string]*Species),
	}
	ds.Turn = galaxyData.TurnNumber
	ds.Radius = galaxyData.Radius

	// convert all stars to systems
	for i := 0; i < starData.NumStars; i++ {
		star := &starData.Stars[i]
		location := coords.New(star.X, star.Y, star.Z, 0)
		s := &System{
			Id:        location.Id(),
			Color:     starColorTranslate(star.Color),
			Location:  location,
			Message:   star.Message,
			Planets:   []*Planet{},
			Type:      starTypeTranslate(star.Type),
			VisitedBy: make(map[string]*Species),
		}
		ds.Systems[s.Id] = s

		// link wormhole
		if star.WormHere != 0 {
			location = coords.New(star.WormX, star.WormY, star.WormZ, 0)
			if end, ok := ds.Systems[location.Id()]; ok {
				s.Wormhole, end.Wormhole = end, s
			}
		}

		// add all planets
		for n := 0; n < star.NumPlanets; n++ {
			location = coords.New(star.X, star.Y, star.Z, n+1)
			planet := planetData.Planets[star.PlanetIndex+n]
			p := &Planet{
				Id:                       location.Id(),
				Atmosphere:               []*AtmosphericGas{},
				Colonies:                 make(map[string]*Colony),
				Diameter:                 planet.Diameter,
				EconEfficiency:           planet.EconEfficiency,
				Gravity:                  planet.Gravity,
				Location:                 location,
				Message:                  planet.Message,
				MiningDifficulty:         planet.MiningDifficulty,
				MiningDifficultyIncrease: planet.MDIncrease,
				PressureClass:            planet.PressureClass,
				TemperatureClass:         planet.TemperatureClass,
				System:                   s,
				VisitedBy:                make(map[string]*Species),
			}
			s.Planets = append(s.Planets, p)
			ds.Planets[p.Id] = p

			for g := 0; g < 4; g++ {
				if planet.Gas[g] != 0 && planet.GasPercent[g] != 0 {
					p.Atmosphere = append(p.Atmosphere, atmosphereTranslate(planet.Gas[g], planet.GasPercent[g]))
				}
			}
		}
	}

	// add all species
	for _, species := range speciesData {
		speciesId := fmt.Sprintf("SP%02d", species.Id)
		sp := &Species{
			Id:         speciesId,
			Ally:       make(map[string]*Species),
			AutoOrders: species.AutoOrders,
			Contact:    make(map[string]*Species),
			EconUnits:  species.EconUnits,
			Enemy:      make(map[string]*Species),
			Name:       species.Name,
			No:         species.Id, // confusing, but ID in original is the Number
			Scanned:    []*System{},
			Visited:    []*System{},
		}
		sp.Colonies.ById = make(map[string]*Colony)
		sp.Colonies.ByLocation = make(map[string]*Colony)
		sp.NamedPlanets.ById = make(map[string]*NamedPlanet)
		sp.NamedPlanets.ByLocation = make(map[string]*NamedPlanet)
		ds.Species[speciesId] = sp

		sp.MI = &Technology{Code: "MI", Level: species.TechLevel[0], KnowledgeLevel: species.TechKnowledge[0], ExperiencePoints: species.TechEps[0]}
		sp.MA = &Technology{Code: "MA", Level: species.TechLevel[1], KnowledgeLevel: species.TechKnowledge[1], ExperiencePoints: species.TechEps[1]}
		sp.ML = &Technology{Code: "ML", Level: species.TechLevel[2], KnowledgeLevel: species.TechKnowledge[2], ExperiencePoints: species.TechEps[2]}
		sp.GV = &Technology{Code: "GV", Level: species.TechLevel[3], KnowledgeLevel: species.TechKnowledge[3], ExperiencePoints: species.TechEps[3]}
		sp.LS = &Technology{Code: "LS", Level: species.TechLevel[4], KnowledgeLevel: species.TechKnowledge[4], ExperiencePoints: species.TechEps[4]}
		sp.BI = &Technology{Code: "BI", Level: species.TechLevel[5], KnowledgeLevel: species.TechKnowledge[5], ExperiencePoints: species.TechEps[5]}
		sp.Fleet.Cost = species.FleetCost
		sp.Fleet.MaintenancePct = species.FleetPercentCost
		sp.Fleet.Ships = make(map[string]*Ship)
		sp.Fleet.Starbases = []*Ship{}
		sp.Fleet.Transports = []*Ship{}
		sp.Fleet.Warships = []*Ship{}
		sp.Gases.Required = gasTranslate(species.RequiredGas)
		sp.Gases.RequiredMaxPct = species.RequiredGasMax
		sp.Gases.RequiredMinPct = species.RequiredGasMin
		for _, gas := range species.NeutralGas {
			sp.Gases.Neutral = append(sp.Gases.Neutral, gasTranslate(gas))
		}
		for _, gas := range species.PoisonGas {
			sp.Gases.Poison = append(sp.Gases.Poison, gasTranslate(gas))
		}
		sp.Government.Name = species.GovtName
		sp.Government.Type = species.GovtType
		location, ok := coords.New(species.X, species.Y, species.Z, species.PN), false
		sp.HomeWorld.OriginalBase = species.HPOriginalBase
		if sp.HomeWorld.Planet, ok = ds.Planets[location.Id()]; !ok {
			panic(fmt.Sprintf("species %q has home world %q which is not in system map", speciesId, location.Id()))
		}
		sp.HomeWorld.System = sp.HomeWorld.Planet.System

		// add all named planets for this species
		for i := 0; i < species.NumNamplas; i++ {
			nampla := &species.NamplaBase[i]
			location := coords.New(nampla.X, nampla.Y, nampla.Z, nampla.PN)
			p, ok := ds.Planets[location.Id()]
			if !ok {
				panic(fmt.Sprintf("species %q has named planet %q which is not in system map", speciesId, nampla.Name))
			}
			np := newNamedPlanet(nampla.Name, p)
			sp.NamedPlanets.ById[np.Id] = np
			sp.NamedPlanets.ByLocation[location.Id()] = np

			// don't add as a colony if not a colony or a home planet
			isHomePlanet, isColony := (nampla.Status&1) != 0, (nampla.Status&2) != 0
			if !(isColony || isHomePlanet) {
				continue
			}

			cc := &Colony{
				Inventory: make(map[string]*Item),
				Name:      np,
				Planet:    p,
				System:    p.System,
			}

			cc.Is.Colony = (nampla.Status & 2) != 0
			cc.Is.DisbandedColony = (nampla.Status & 64) != 0
			cc.Is.Hiding = nampla.Hiding
			cc.Is.Hidden = nampla.Hidden
			cc.Is.HomePlanet = (nampla.Status & 1) != 0
			cc.Is.MiningColony = (nampla.Status & 16) != 0
			cc.Is.Populated = (nampla.Status & 8) != 0
			cc.Is.ResortColony = (nampla.Status & 32) != 0

			for code, qty := range nampla.ItemQuantity {
				if qty > 0 {
					item := itemTranslate(code, qty)
					cc.Inventory[item.Code] = item
				}
			}

			cc.Manufacturing.Auto = nampla.AutoAUs != 0
			cc.Manufacturing.AvailableToInstall = nampla.AUsToInstall
			cc.Manufacturing.Base = nampla.MaBase
			cc.Manufacturing.Needed = nampla.AUsNeeded
			cc.Message = nampla.Message
			cc.Mining.Auto = nampla.AutoIUs != 0
			cc.Mining.AvailableToInstall = nampla.IUsToInstall
			cc.Mining.Base = nampla.MiBase
			cc.Mining.Needed = nampla.IUsNeeded
			cc.Population = nampla.PopUnits
			cc.SiegeEff = nampla.SiegeEff
			cc.Shipyards = nampla.Shipyards
			cc.Special = nampla.Special
			cc.UseOnAmbush = nampla.UseOnAmbush

			sp.Colonies.ById[cc.Name.Id] = cc
			sp.Colonies.ByLocation[cc.Planet.Id] = cc

			p.Colonies[speciesId] = cc
		}

		// add all ships for this species
		for i := 0; i < species.NumShips; i++ {
			ship := &species.ShipBase[i]
			if ship.PN == 99 || ship.Name == "Unused" {
				continue
			}
			sh := &Ship{
				Id:                 strings.ToUpper(strings.TrimSpace(ship.Name)),
				Age:                ship.Age,
				ArrivedViaWormhole: ship.ArrivedViaWormhole,
				Class:              shipClassTranslate(ship.Class),
				DisplayName:        strings.TrimSpace(ship.Name),
				Inventory:          make(map[string]*Item),
				JustJumped:         ship.JustJumped,
				LoadingPoint:       ship.LoadingPoint,
				Location:           coords.New(ship.X, ship.Y, ship.Z, ship.PN),
				RemainingCost:      ship.RemainingCost,
				Special:            ship.Special,
				Status:             shipStatusTranslate(ship.Status),
				UnloadingPoint:     ship.UnloadingPoint,
			}
			sh.Class.FTL = ship.Type == 0
			if ship.DestX != 0 && ship.DestY != 0 && ship.DestZ != 0 {
				sh.Destination = coords.New(ship.DestX, ship.DestY, ship.DestZ, 0)
			}
			for code, qty := range ship.ItemQuantity {
				if qty > 0 {
					item := itemTranslate(code, qty)
					sh.Inventory[item.Code] = item
				}
			}
			if len(sh.Inventory) == 0 && sh.DisplayName == "Mungo" {
				item := itemTranslate(3, 23)
				sh.Inventory[item.Code] = item
				item = itemTranslate(13, 3)
				sh.Inventory[item.Code] = item
			}
			if len(sh.Inventory) == 0 && sh.DisplayName == "Mongo" {
				item := itemTranslate(2, 51)
				sh.Inventory[item.Code] = item
				item = itemTranslate(13, 3)
				sh.Inventory[item.Code] = item
			}
			sp.Fleet.Ships[sh.Id] = sh
			if sh.Class.Code == "BA" {
				sp.Fleet.Starbases = append(sp.Fleet.Starbases, sh)
			} else if sh.Class.Code == "TR" {
				sp.Fleet.Transports = append(sp.Fleet.Transports, sh)
			} else {
				sp.Fleet.Warships = append(sp.Fleet.Warships, sh)
			}
		}
	}

	// map systems visited
	for i := 0; i < starData.NumStars; i++ {
		star := &starData.Stars[i]
		location := coords.Coords{X: star.X, Y: star.Y, Z: star.Z}
		s := ds.Systems[location.Id()]
		for _, speciesNo := range star.VisitedBy {
			for _, sp := range ds.Species {
				if sp.No == speciesNo {
					sp.Visited = append(sp.Visited, s)
					s.VisitedBy[sp.Id] = sp
					break
				}
			}
		}
	}

	// map systems scanned.
	// todo: we don't have data for scanned systems, so just pretend that visited === scanned
	for _, sp := range ds.Species {
		sp.Scanned = append(sp.Scanned, sp.Visited...)
	}

	// sort some data to keep lists consistent between runs

	return ds, nil
}

func loader(name string, a interface{}) error {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, a)
}
