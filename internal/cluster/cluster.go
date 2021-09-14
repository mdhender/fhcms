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
	"fmt"
	"github.com/mdhender/fhcms/internal/dat32"
	"sort"
	"strings"
)

// Cluster represents the entire cluster. It contains all the systems,
// planets, species, named planets, colonies, and ships in the game.
type Cluster struct {
	Turn    int                 `json:"turn"`   // current turn number
	Radius  int                 `json:"radius"` // radius of the cluster
	Systems map[string]*System  `json:"systems"`
	Planets map[string]*Planet  `json:"planets"`
	Species map[string]*Species `json:"species"`
}

func ConvertDat32ToCluster(galaxy *dat32.Galaxy, stars *dat32.Stars, planets *dat32.Planets, species []*dat32.Species) (*Cluster, error) {
	c := &Cluster{
		Systems: make(map[string]*System),
		Planets: make(map[string]*Planet),
		Species: make(map[string]*Species),
	}
	c.Turn = galaxy.TurnNumber
	c.Radius = galaxy.Radius

	// add all systems
	for i := 0; i < stars.NumStars; i++ {
		star := &stars.Stars[i]
		cs := &System{
			Coords:  &Coords{X: star.X, Y: star.Y, Z: star.Z},
			Planets: []string{},
		}
		cs.Type.DisplayCode, cs.Type.Description = starTypeTranslate(star.Type)
		cs.Color.DisplayCode, cs.Color.Description = starColorTranslate(star.Type)
		cs.PotentialHomeSystem = star.HomeSystem != 0
		if star.WormHere != 0 {
			cs.Wormhole = &Coords{X: star.WormX, Y: star.WormY, Z: star.WormZ}
		}
		cs.Message = star.Message
		c.Systems[cs.Coords.Id()] = cs
	}

	// add all planets
	for i := 0; i < stars.NumStars; i++ {
		star := &stars.Stars[i]
		coords := Coords{X: star.X, Y: star.Y, Z: star.Z}
		cs := c.Systems[coords.Id()]
		for n := 0; n < star.NumPlanets; n++ {
			planet := planets.Planets[star.PlanetIndex+n]
			cp := &Planet{
				Coords:     &Coords{X: star.X, Y: star.Y, Z: star.Z, Orbit: n + 1},
				Colonies:   make(map[string]*Colony),
				Atmosphere: []GasType{},
			}
			cp.System = cs.Coords.Id()
			for g := 0; g < 4; g++ {
				if planet.Gas[g] != 0 && planet.GasPercent[g] != 0 {
					code, _ := gasToType(planet.Gas[g])
					cp.Atmosphere = append(cp.Atmosphere, GasType{
						Code:       code,
						Percentage: planet.GasPercent[g],
					})
				}
			}
			cp.Diameter = planet.Diameter
			cp.EconEfficiency = planet.EconEfficiency
			cp.Gravity = planet.Gravity
			cp.Message = planet.Message
			cp.MiningDifficulty = planet.MiningDifficulty
			cp.MiningDifficultyIncrease = planet.MDIncrease
			cp.PressureClass = planet.PressureClass
			cp.TemperatureClass = planet.TemperatureClass
			switch planet.Special {
			case 1:
				cp.Special = "Ideal Home Planet"
			case 2:
				cp.Special = "Ideal Colony Planet"
			case 3:
				cp.Special = "Radioactive Hellhole"
			}

			c.Planets[cp.Coords.Id()] = cp
			cs.Planets = append(cs.Planets, cp.Coords.Id())
		}
	}

	// add all species
	for i := 0; i < len(species); i++ {
		sp := species[i]
		id := fmt.Sprintf("SP%02d", sp.Id)
		coords := Coords{X: sp.X, Y: sp.Y, Z: sp.Z, Orbit: sp.PN}
		cs := &Species{
			NamedPlanets: make(map[string]*NamedPlanet),
			Ally:         []string{},
			Colonies:     []string{},
			Contact:      []string{},
			Enemy:        []string{},
			Scanned:      []string{},
			Ships:        make(map[string]*Ship),
			Visited:      []string{},
		}
		cs.Name = sp.Name
		cs.Government.Name = sp.GovtName
		cs.Government.Type = sp.GovtType
		cs.HomeWorld = coords.Id()

		for _, ally := range sp.Ally {
			cs.Ally = append(cs.Ally, fmt.Sprintf("SP%02d", ally))
		}
		cs.AutoOrders = sp.AutoOrders
		// add all colonies for this species
		for n := 0; n < sp.NumNamplas; n++ {
			np := &sp.NamplaBase[n]
			coords := Coords{X: np.X, Y: np.Y, Z: np.Z, Orbit: np.PN}
			cc := &Colony{Inventory: make(map[string]int)}

			// don't add as a colony if not a colony or a home planet!
			cc.Status.HomePlanet = (np.Status & 1) != 0 // status codes for colonies are logically ORed together
			cc.Status.Colony = (np.Status & 2) != 0
			if !(cc.Status.Colony || cc.Status.HomePlanet) {
				continue
			}

			cc.Status.DisbandedColony = (np.Status & 64) != 0
			cc.Status.Hiding = np.Hiding
			cc.Status.Hidden = np.Hidden
			cc.Status.MiningColony = (np.Status & 16) != 0
			cc.Status.Populated = (np.Status & 8) != 0
			cc.Status.ResortColony = (np.Status & 32) != 0

			for item, qty := range np.ItemQuantity {
				if qty > 0 {
					code, _ := itemToCode(item)
					cc.Inventory[code] = qty
				}
			}
			cc.MaBase = np.MaBase
			cc.Message = np.Message
			cc.MiBase = np.MiBase
			cc.PlanetIndex = np.PlanetIndex
			cc.PopUnits = np.PopUnits
			cc.SiegeEff = np.SiegeEff
			cc.Shipyards = np.Shipyards
			cc.Special = np.Special
			cc.Units.Manufacturing.Auto = np.AutoAUs != 0
			cc.Units.Manufacturing.AvailableToInstall = np.AUsToInstall
			cc.Units.Manufacturing.Needed = np.AUsNeeded
			cc.Units.Mining.Auto = np.AutoIUs != 0
			cc.Units.Mining.AvailableToInstall = np.IUsToInstall
			cc.Units.Mining.Needed = np.IUsNeeded
			cc.UseOnAmbush = np.UseOnAmbush

			cs.Colonies = append(cs.Colonies, coords.Id())
			planet := c.Planets[coords.Id()]
			planet.Colonies[id] = cc
		}
		for _, contact := range sp.Contact {
			cs.Contact = append(cs.Contact, fmt.Sprintf("SP%02d", contact))
		}
		cs.EconUnits = sp.EconUnits
		for _, enemy := range sp.Enemy {
			cs.Enemy = append(cs.Enemy, fmt.Sprintf("SP%02d", enemy))
		}
		cs.Fleet.Cost = sp.FleetCost
		cs.Fleet.MaintenancePct = sp.FleetPercentCost
		cs.Gases.Required.Code, _ = gasToType(sp.RequiredGas)
		cs.Gases.Required.MinPct = sp.RequiredGasMin
		cs.Gases.Required.MaxPct = sp.RequiredGasMax
		for _, g := range sp.NeutralGas {
			if g > 0 {
				code, _ := gasToType(g)
				cs.Gases.Neutral = append(cs.Gases.Neutral, code)
			}
		}
		for _, g := range sp.PoisonGas {
			if g > 0 {
				code, _ := gasToType(g)
				cs.Gases.Poison = append(cs.Gases.Poison, code)
			}
		}
		cs.HPOriginalBase = sp.HPOriginalBase
		// add all named planets for this species
		for n := 0; n < sp.NumNamplas; n++ {
			np := &sp.NamplaBase[n]
			coords := Coords{X: np.X, Y: np.Y, Z: np.Z, Orbit: np.PN}
			cs.NamedPlanets[strings.ToUpper(np.Name)] = &NamedPlanet{
				Planet:      coords.Id(),
				DisplayName: np.Name,
			}
		}
		// Scanned will be added later
		// add all ships for this species
		for n := 0; n < sp.NumShips; n++ {
			ship := &sp.ShipBase[n]
			if ship.Name == "Unused" {
				continue
			}
			cc := &Ship{Inventory: make(map[string]int)}
			cc.DisplayName = ship.Name
			cc.Class.Code, cc.Class.Description, cc.Class.Tonnage, cc.Class.Cost = shipToClass(ship.Class)
			cc.Class.FTL = ship.Type == 0
			cc.Age = ship.Age
			cc.Location = &Coords{X: ship.X, Y: ship.Y, Z: ship.Z, Orbit: ship.PN}
			cc.ArrivedViaWormhole = ship.ArrivedViaWormhole
			if ship.DestX != 0 && ship.DestY != 0 && ship.DestZ != 0 {
				cc.Destination = &Coords{X: ship.DestX, Y: ship.DestY, Z: ship.DestZ}
			}
			for item, qty := range ship.ItemQuantity {
				if qty > 0 {
					code, _ := itemToCode(item)
					cc.Inventory[code] = qty
				}
			}
			cc.JustJumped = ship.JustJumped
			cc.LoadingPoint = ship.LoadingPoint
			cc.UnloadingPoint = ship.UnloadingPoint
			cc.RemainingCost = ship.RemainingCost
			cc.Special = ship.Special
			switch ship.Status {
			case 0:
				cc.Status.UnderConstruction = true
			case 1:
				cc.Status.OnSurface = true
			case 2:
				cc.Status.InOrbit = true
			case 3:
				cc.Status.InDeepSpace = true
			case 4:
				cc.Status.JumpedInCombat = true
			case 5:
				cc.Status.ForcedJump = true
			default:
				panic(fmt.Sprintf("assert(ship.status != %d)", i))
			}

			cs.Ships[strings.ToUpper(ship.Name)] = cc
		}
		cs.Tech.MI = Technology{Code: "MI", Level: sp.TechLevel[0], KnowledgeLevel: sp.TechKnowledge[0], ExperiencePoints: sp.TechEps[0]}
		cs.Tech.MA = Technology{Code: "MA", Level: sp.TechLevel[1], KnowledgeLevel: sp.TechKnowledge[1], ExperiencePoints: sp.TechEps[1]}
		cs.Tech.ML = Technology{Code: "ML", Level: sp.TechLevel[2], KnowledgeLevel: sp.TechKnowledge[2], ExperiencePoints: sp.TechEps[2]}
		cs.Tech.GV = Technology{Code: "GV", Level: sp.TechLevel[3], KnowledgeLevel: sp.TechKnowledge[3], ExperiencePoints: sp.TechEps[3]}
		cs.Tech.LS = Technology{Code: "LS", Level: sp.TechLevel[4], KnowledgeLevel: sp.TechKnowledge[4], ExperiencePoints: sp.TechEps[4]}
		cs.Tech.BI = Technology{Code: "BI", Level: sp.TechLevel[5], KnowledgeLevel: sp.TechKnowledge[5], ExperiencePoints: sp.TechEps[5]}
		// Visited will be added later

		c.Species[id] = cs
	}

	// map systems visited
	for i := 0; i < stars.NumStars; i++ {
		star := &stars.Stars[i]
		coords := Coords{X: star.X, Y: star.Y, Z: star.Z}
		cs := c.Systems[coords.Id()]
		for _, spNo := range star.VisitedBy {
			spId := fmt.Sprintf("SP%02d", spNo)
			if sp, ok := c.Species[spId]; ok {
				sp.Visited = append(sp.Visited, coords.Id())
				cs.VisitedBy = append(cs.VisitedBy, spId)
			}
		}
	}

	// map systems scanned.
	// todo: we don't have data for scanned systems, so just pretend
	// that visited === scanned
	for _, sp := range c.Species {
		sp.Scanned = append(sp.Scanned, sp.Visited...)
	}

	// sort some data to keep lists consistent between runs
	for _, s := range c.Systems {
		sort.Strings(s.Planets)
		sort.Strings(s.VisitedBy)
	}
	for _, sp := range c.Species {
		sort.Strings(sp.Colonies)
		sort.Strings(sp.Gases.Neutral)
		sort.Strings(sp.Gases.Poison)
		sort.Strings(sp.Scanned)
		sort.Strings(sp.Visited)
	}

	return c, nil
}
