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

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type ClusterData struct {
	Turn    int                     `json:"turn"`   // current turn number
	Radius  int                     `json:"radius"` // radius of the cluster
	Systems map[string]*SystemData  `json:"systems"`
	Planets map[string]*PlanetData  `json:"planets"`
	Species map[string]*SpeciesData `json:"species"`
}

type SystemData struct {
	Coords *Coords `json:"coords"` // location of the star within the cluster
	Color  struct {
		DisplayCode string `json:"display_code"`
		Description string `json:"description"`
	} `json:"color"`
	Message             int      `json:"message,omitempty"`
	Planets             []string `json:"planets"`                         // list of identifiers for each planet, ordered by orbit
	PotentialHomeSystem bool     `json:"potential_home_system,omitempty"` // true if the system is a good potential home system
	Type                struct {
		DisplayCode string `json:"display_code"`
		Description string `json:"description"`
	} `json:"type"`
	VisitedBy []string `json:"visited_by"`         // list of identifiers for every species that has visited the system
	Wormhole  *Coords  `json:"wormhole,omitempty"` // coordinates of other end of wormhole, nil if not a wormhole
}

type PlanetData struct {
	Coords                   *Coords                `json:"coords"` // location of the planet within the cluster
	System                   string                 `json:"system"` // identifier for system containing the planet
	Atmosphere               []GasType              `json:"atmosphere"`
	Colonies                 map[string]*ColonyData `json:"colonies"` // key is species
	Diameter                 int                    `json:"diameter"`
	EconEfficiency           int                    `json:"econ_efficiency"`
	Gravity                  int                    `json:"gravity"`
	Message                  int                    `json:"message"`
	MiningDifficulty         int                    `json:"mining_difficulty"`
	MiningDifficultyIncrease int                    `json:"mining_difficulty_increase"`
	PressureClass            int                    `json:"pressure_class"`
	TemperatureClass         int                    `json:"temperature_class"`
	Special                  string                 `json:"special,omitempty"`
}

// SpeciesData is used for the exported data.
// TODO: Scan should track the scan results since the planet attributes can change during play.
type SpeciesData struct {
	Name       string `json:"name"` // name of the species
	Government struct {
		Name string `json:"name"`
		Type string `json:"type"`
	} `json:"government"`
	HomeWorld  string   `json:"home_world"` // coordinates of the home world
	Ally       []string `json:"ally"`
	AutoOrders bool     `json:"auto_orders,omitempty"`
	Colonies   []string `json:"colonies"` // coordinates of the planet containing the colony
	Contact    []string `json:"contact"`
	EconUnits  int      `json:"econ_units"`
	Enemy      []string `json:"enemy"`
	Fleet      struct {
		Cost           int `json:"cost"`
		MaintenancePct int `json:"maintenance_pct"` // percentage of production applied to fleet maintenance
	} `json:"fleet"`
	Gases struct {
		Required struct {
			Code   string `json:"code"`
			MinPct int    `json:"min_pct"`
			MaxPct int    `json:"max_pct"`
		} `json:"required"`
		Neutral []string `json:"neutral"`
		Poison  []string `json:"poison"`
	} `json:"gases"`
	HPOriginalBase int                         `json:"hp_original_base"`
	NamedPlanets   map[string]*NamedPlanetData `json:"named_planets"` // key is name of planet, converted to upper case
	Scanned        []string                    `json:"scanned"`       // coordinates of all systems that have been scanned
	Ships          map[string]*ShipData        `json:"ships"`         // key is name of ship, converted to upper case
	Tech           struct {
		MI TechLevelData `json:"mi"`
		MA TechLevelData `json:"ma"`
		ML TechLevelData `json:"ml"`
		GV TechLevelData `json:"gv"`
		LS TechLevelData `json:"ls"`
		BI TechLevelData `json:"bi"`
	} `json:"tech"`
	Visited []string `json:"visited"` // coordinates of all systems that have been visited
}

type NamedPlanetData struct {
	Planet      string `json:"planet"`       // coordinates of the planet being named
	DisplayName string `json:"display_name"` // original name of the planet
}

type ColonyData struct {
	Status struct {
		Colony          bool `json:"colony"`
		DisbandedColony bool `json:"disbanded_colony,omitempty"`
		Hiding          bool `json:"hiding,omitempty"`
		Hidden          bool `json:"hidden,omitempty"`
		HomePlanet      bool `json:"home_planet,omitempty"`
		MiningColony    bool `json:"mining_colony,omitempty"`
		Populated       bool `json:"populated,omitempty"`
		ResortColony    bool `json:"resort_colony,omitempty"`
	} `json:"status"`
	Inventory   map[string]int `json:"inventory"` // key is item code, value is quantity
	MaBase      int            `json:"ma_base"`
	Message     int            `json:"message"`
	MiBase      int            `json:"mi_base"`
	PlanetIndex int            `json:"planet_index"`
	PopUnits    int            `json:"pop_units"`
	SiegeEff    int            `json:"siege_eff"`
	Shipyards   int            `json:"shipyards"`
	Special     int            `json:"special,omitempty"`
	Units       struct {
		Colonists struct {
			Auto               bool `json:"auto,omitempty"`
			AvailableToInstall int  `json:"available_to_install,omitempty"`
			Needed             int  `json:"needed,omitempty"`
		} `json:"colonists"`
		Manufacturing struct {
			Auto               bool `json:"auto,omitempty"`
			AvailableToInstall int  `json:"available_to_install,omitempty"`
			Needed             int  `json:"needed,omitempty"`
		} `json:"manufacturing"`
		Mining struct {
			Auto               bool `json:"auto,omitempty"`
			AvailableToInstall int  `json:"available_to_install,omitempty"`
			Needed             int  `json:"needed,omitempty"`
		} `json:"mining"`
	} `json:"units"`
	UseOnAmbush int `json:"use_on_ambush"`
}

type ShipData struct {
	DisplayName string // original name of the ship
	Class       struct {
		Code        string `json:"code"`
		Description string `json:"description"`
		FTL         bool   `json:"ftl,omitempty"`
		Tonnage     int    `json:"tonnage"`
		Cost        int    `json:"cost"`
	} `json:"class"`
	Age                int            `json:"age"`
	Location           *Coords        `json:"location,omitempty"`
	ArrivedViaWormhole bool           `json:"arrived_via_wormhole,omitempty"`
	Destination        *Coords        `json:"destination,omitempty"`
	Inventory          map[string]int `json:"inventory"` // key is item code, value is quantity
	JustJumped         bool           `json:"just_jumped,omitempty"`
	LoadingPoint       int            `json:"loading_point,omitempty"`
	UnloadingPoint     int            `json:"unloading_point,omitempty"`
	RemainingCost      int            `json:"remaining_cost,omitempty"`
	Special            int            `json:"special,omitempty"`
	Status             struct {
		ForcedJump        bool `json:"forced_jump,omitempty"`
		InDeepSpace       bool `json:"in_deep_space,omitempty"`
		InOrbit           bool `json:"in_orbit,omitempty"`
		JumpedInCombat    bool `json:"jumped_in_combat,omitempty"`
		OnSurface         bool `json:"on_surface,omitempty"`
		UnderConstruction bool `json:"under_construction,omitempty"`
	} `json:"status"`
}

type TechLevelData struct {
	Code             string `json:"code"`
	Level            int    `json:"level"`             // current level
	KnowledgeLevel   int    `json:"knowledge_level"`   // un-applied tech level knowledge
	ExperiencePoints int    `json:"experience_points"` // experience points for tech levels
}

type Coords struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Z     int `json:"z"`
	Orbit int `json:"orbit,omitempty"`
}

func (c *Coords) Id() string {
	if c.Orbit != 0 {
		return fmt.Sprintf("%d.%d.%d.%d", c.X, c.Y, c.Z, c.Orbit)
	}
	return fmt.Sprintf("%d.%d.%d", c.X, c.Y, c.Z)
}

func NewCluster(galaxy *GalaxyData, stars *Stars, planets *Planets, species []*Species) *ClusterData {
	c := &ClusterData{
		Systems: make(map[string]*SystemData),
		Planets: make(map[string]*PlanetData),
		Species: make(map[string]*SpeciesData),
	}
	c.Turn = galaxy.TurnNumber
	c.Radius = galaxy.Radius

	// add all systems
	for i := 0; i < stars.NumStars; i++ {
		star := &stars.Stars[i]
		cs := &SystemData{
			Coords:  &Coords{X: star.X, Y: star.Y, Z: star.Z},
			Planets: []string{},
		}
		cs.Type.DisplayCode, cs.Type.Description = starTypeTranslate(star.Type)
		cs.Color.DisplayCode, cs.Color.Description = starColorTranslate(star.Type)
		cs.PotentialHomeSystem = star.HomeSystem
		if star.WormHere {
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
			cp := &PlanetData{
				Coords:     &Coords{X: star.X, Y: star.Y, Z: star.Z, Orbit: n + 1},
				Colonies:   make(map[string]*ColonyData),
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
		cs := &SpeciesData{
			NamedPlanets: make(map[string]*NamedPlanetData),
			Ally:         []string{},
			Colonies:     []string{},
			Contact:      []string{},
			Enemy:        []string{},
			Scanned:      []string{},
			Ships:        make(map[string]*ShipData),
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
			cc := &ColonyData{Inventory: make(map[string]int)}

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
			cs.NamedPlanets[strings.ToUpper(np.Name)] = &NamedPlanetData{
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
			cc := &ShipData{Inventory: make(map[string]int)}
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
		cs.Tech.MI = TechLevelData{Code: "MI", Level: sp.TechLevel[0], KnowledgeLevel: sp.TechKnowledge[0], ExperiencePoints: sp.TechEps[0]}
		cs.Tech.MA = TechLevelData{Code: "MA", Level: sp.TechLevel[1], KnowledgeLevel: sp.TechKnowledge[1], ExperiencePoints: sp.TechEps[1]}
		cs.Tech.ML = TechLevelData{Code: "ML", Level: sp.TechLevel[2], KnowledgeLevel: sp.TechKnowledge[2], ExperiencePoints: sp.TechEps[2]}
		cs.Tech.GV = TechLevelData{Code: "GV", Level: sp.TechLevel[3], KnowledgeLevel: sp.TechKnowledge[3], ExperiencePoints: sp.TechEps[3]}
		cs.Tech.LS = TechLevelData{Code: "LS", Level: sp.TechLevel[4], KnowledgeLevel: sp.TechKnowledge[4], ExperiencePoints: sp.TechEps[4]}
		cs.Tech.BI = TechLevelData{Code: "BI", Level: sp.TechLevel[5], KnowledgeLevel: sp.TechKnowledge[5], ExperiencePoints: sp.TechEps[5]}
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

	return c
}

func (c *ClusterData) Write(name string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(name, data, 0644)
}
