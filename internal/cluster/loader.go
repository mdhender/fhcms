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
	"github.com/mdhender/fhcms/internal/dat32"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func FromDat32(path string, bigEndian bool) (*Store, error) {
	galaxyDataFile := filepath.Join(path, "galaxy.dat")
	starDataFile := filepath.Join(path, "stars.dat")
	planetDataFile := filepath.Join(path, "planets.dat")
	speciesDataPath := path
	locationDataFile := filepath.Join(path, "locations.dat")
	//transactionDataFile := filepath.Join(path, "transactions.dat")
	//log.Printf("todo: use transactionDataFile %q\n", transactionDataFile)

	var bo binary.ByteOrder
	if bigEndian {
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
		Systems:     make(map[string]*System),
		Planets:     make(map[string]*Planet),
		Species:     make(map[string]*Species),
		SpeciesBase: make([]*Species, galaxyData.NumSpecies+1, galaxyData.NumSpecies+1),
	}
	ds.Turn = galaxyData.TurnNumber
	ds.Radius = galaxyData.Radius

	// convert all stars to systems
	for i := 0; i < starData.NumStars; i++ {
		star := &starData.Stars[i]
		location := &Coords{X: star.X, Y: star.Y, Z: star.Z}
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
			location = &Coords{X: star.WormX, Y: star.WormY, Z: star.WormZ}
			if end, ok := ds.Systems[location.Id()]; ok {
				s.Wormhole, end.Wormhole = end, s
			}
		}

		// add all planets
		for n := 0; n < star.NumPlanets; n++ {
			location = &Coords{X: star.X, Y: star.Y, Z: star.Z, Orbit: n + 1}
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

		ds.SpeciesBase[species.Id] = sp
		ds.Species[speciesId] = sp

		sp.MI = &Technology{Code: "MI", Level: species.TechLevel[0], KnowledgeLevel: species.TechKnowledge[0], ExperiencePoints: species.TechEps[0]}
		sp.MA = &Technology{Code: "MA", Level: species.TechLevel[1], KnowledgeLevel: species.TechKnowledge[1], ExperiencePoints: species.TechEps[1]}
		sp.ML = &Technology{Code: "ML", Level: species.TechLevel[2], KnowledgeLevel: species.TechKnowledge[2], ExperiencePoints: species.TechEps[2]}
		sp.GV = &Technology{Code: "GV", Level: species.TechLevel[3], KnowledgeLevel: species.TechKnowledge[3], ExperiencePoints: species.TechEps[3]}
		sp.LS = &Technology{Code: "LS", Level: species.TechLevel[4], KnowledgeLevel: species.TechKnowledge[4], ExperiencePoints: species.TechEps[4]}
		sp.BI = &Technology{Code: "BI", Level: species.TechLevel[5], KnowledgeLevel: species.TechKnowledge[5], ExperiencePoints: species.TechEps[5]}

		// We must use the LS tech level at the start of the turn because
		// the distorted species number must be the same throughout the
		// turn, even if the tech level changes during production.
		nibLo, nibHi := sp.No&0x000F, (sp.No>>4)&0x000F // lower and upper four bits
		sp.DistortedNumber = (sp.LS.Level%5+3)*(4*nibLo+nibHi) + (sp.LS.Level%11 + 7)

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
		location, ok := &Coords{X: species.X, Y: species.Y, Z: species.Z, Orbit: species.PN}, false
		sp.HomeWorld.OriginalBase = species.HPOriginalBase
		if sp.HomeWorld.Planet, ok = ds.Planets[location.Id()]; !ok {
			panic(fmt.Sprintf("species %q has home world %q which is not in system map", speciesId, location.Id()))
		}
		sp.HomeWorld.System = sp.HomeWorld.Planet.System

		// add all named planets for this species
		sp.NamedPlanets.Base = make([]*NamedPlanet, species.NumNamplas, species.NumNamplas)
		for i := 0; i < species.NumNamplas; i++ {
			nampla := &species.NamplaBase[i]
			location := &Coords{X: nampla.X, Y: nampla.Y, Z: nampla.Z, Orbit: nampla.PN}
			p, ok := ds.Planets[location.Id()]
			if !ok {
				panic(fmt.Sprintf("species %q has named planet %q which is not in system map", speciesId, nampla.Name))
			}
			np := newNamedPlanet(nampla.Name, p, i)
			sp.NamedPlanets.Base[i] = np
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
			np.Colony = cc

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
					cc.Inventory[item.Abbr] = item
				}
			}

			cc.Manufacturing.AutoAUs = nampla.AutoAUs
			cc.Manufacturing.AvailableToInstall = nampla.AUsToInstall
			cc.Manufacturing.Base = nampla.MaBase
			cc.Manufacturing.Needed = nampla.AUsNeeded
			cc.Message = nampla.Message
			cc.Mining.AutoIUs = nampla.AutoIUs
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
				Index:              i,
				Inventory:          make(map[string]*Item),
				JustJumped:         ship.JustJumped,
				LoadingPoint:       ship.LoadingPoint,
				Location:           &Coords{X: ship.X, Y: ship.Y, Z: ship.Z, Orbit: ship.PN},
				RemainingCost:      ship.RemainingCost,
				Special:            ship.Special,
				Species:            sp,
				Status:             shipStatusTranslate(ship.Status),
				UnloadingPoint:     ship.UnloadingPoint,
			}
			sh.Class.Tonnage = ship.Tonnage
			// Ship types
			switch ship.Type {
			case 0: //   #define FTL             0
				sh.Class.Is.SubLight = false
			case 1: //   #define SUB_LIGHT       1
				sh.Class.Is.SubLight = true
			case 2: //   #define STARBASE        2
				sh.Class.Is.SubLight = true
				sh.Class.Is.Starbase = true
			default:
				panic(fmt.Sprintf("assert(ship.type != %d)", ship.Type))
			}
			if ship.DestX != 0 && ship.DestY != 0 && ship.DestZ != 0 {
				sh.Destination = &Coords{X: ship.DestX, Y: ship.DestY, Z: ship.DestZ}
			}
			sh.Name = strings.TrimSpace(ship.Name)
			var sublightFlag string
			if sh.Class.Is.SubLight {
				sublightFlag = "S"
			}
			if sh.Class.Is.Transport {
				sh.Display.Name = fmt.Sprintf("%s%d%s %s", sh.Class.Code, sh.Class.Tonnage, sublightFlag, sh.Name)
			} else {
				sh.Display.Name = fmt.Sprintf("%s%s %s", sh.Class.Code, sublightFlag, sh.Name)
			}

			sh.Display.Tonnage = fmt.Sprintf("%dk", 10*ship.Tonnage)
			for code, qty := range ship.ItemQuantity {
				if qty > 0 {
					item := itemTranslate(code, qty)
					sh.Inventory[item.Abbr] = item
				}
			}
			sp.Fleet.Base = append(sp.Fleet.Base, sh)
			sp.Fleet.Ships[sh.Id] = sh
			if sh.Class.Is.Starbase {
				sp.Fleet.Starbases = append(sp.Fleet.Starbases, sh)
			} else if sh.Class.Is.Transport {
				sp.Fleet.Transports = append(sp.Fleet.Transports, sh)
			} else if sh.Class.Is.Warship {
				sp.Fleet.Warships = append(sp.Fleet.Warships, sh)
			}
		}
	}

	// populate the maps of contacts, allies, and enemies.
	for _, species := range speciesData {
		speciesId := fmt.Sprintf("SP%02d", species.Id)
		sp := ds.Species[speciesId]

		// add the contact only if it's in the list of species.
		// (it should be a bug if it isn't.)
		for _, spNo := range species.Contact {
			if o, ok := ds.Species[fmt.Sprintf("SP%02d", spNo)]; ok {
				sp.Contact[fmt.Sprintf("SP%02d", spNo)] = o
			}
		}
		// add the ally only if it's in the list of contacts.
		// checking against the contacts guards against the player
		// declaring species to be allies before they've met.
		for _, spNo := range species.Ally {
			if o, ok := sp.Contact[fmt.Sprintf("SP%02d", spNo)]; ok {
				sp.Ally[fmt.Sprintf("SP%02d", spNo)] = o
			}
		}
		// add the enemy only if it's in the list of contacts.
		// checking against the contacts guards against the player
		// declaring species to be enemies before they've met.
		for _, spNo := range species.Enemy {
			if o, ok := sp.Contact[fmt.Sprintf("SP%02d", spNo)]; ok {
				sp.Enemy[fmt.Sprintf("SP%02d", spNo)] = o
			}
		}
	}

	// map systems visited
	for i := 0; i < starData.NumStars; i++ {
		star := &starData.Stars[i]
		location := &Coords{X: star.X, Y: star.Y, Z: star.Z}
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

	ld, err := dat32.ReadLocations(locationDataFile, bo)
	if err != nil {
		return nil, err
	}
	ds.Locations = make([]SpeciesLocationData, len(ld), len(ld))
	for i, loc := range ld {
		ds.Locations[i].Species = ds.SpeciesBase[loc.S]
		ds.Locations[i].Location = NewCoords(loc.X, loc.Y, loc.Z, 0)
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

/* The following routine provides the 'distorted' species number used to
 *      identify a species that uses field distortion units. The input
 *      variable 'species_number' is the same number used in filename
 *      creation for the species. */
func (sp *Species) Distorted() int {
	return sp.DistortedNumber
}
