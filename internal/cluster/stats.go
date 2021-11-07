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

type TechLevel struct {
	Name, Code      string
	Total, Min, Max int
}

type GameStats struct {
	TotalSpecies                                                                                                                                                      int
	BankedEconUnits, FleetMaintenancePct, PopulatedPlanets, Production, Shipyards, Starbases, StarbaseTonnage, Transports, TransportTonnage, Warships, WarshipTonnage struct {
		Total, Min, Max, Average float64
	}
	MI, MA, ML, GV, LS, BI TechLevel
	Stats                  []*Stats
}

type Stats struct {
	No                     int    // species number
	Name                   string // species name
	MI, MA, ML, GV, LS, BI int    // tech levels
	FleetMaintenancePct    float64
	Production             float64
	PopulatedPlanets       float64
	Ships                  float64
	Warships               float64
	WarshipTonnage         float64
	Starbases              float64
	StarbaseTonnage        float64
	Transports             float64
	TransportTonnage       float64
	Shipyards              float64
	OffensivePower         float64
	DefensivePower         float64
	BankedEconUnits        float64
}

func (ds *Store) Stats() *GameStats {
	// sort by species number to avoid the random output from a map
	var xlat []*Species
	for _, sp := range ds.Species {
		xlat = append(xlat, sp)
	}
	for i := 0; i < len(xlat); i++ {
		for j := i + 1; j < len(xlat); j++ {
			if xlat[i].No > xlat[j].No {
				xlat[i], xlat[j] = xlat[j], xlat[i]
			}
		}
	}

	gameStats := &GameStats{}
	gameStats.MI = TechLevel{Name: "Mining", Code: "MI"}
	gameStats.MA = TechLevel{Name: "Manufacturing", Code: "MA"}
	gameStats.ML = TechLevel{Name: "Military", Code: "ML"}
	gameStats.GV = TechLevel{Name: "Gravitics", Code: "GV"}
	gameStats.LS = TechLevel{Name: "Life Support", Code: "LS"}
	gameStats.BI = TechLevel{Name: "Biology", Code: "BI"}

	// compute statistics for each species
	for _, species := range xlat {
		gameStats.Stats = append(gameStats.Stats, &Stats{
			No:                  species.No,
			Name:                species.Name,
			MI:                  species.MI.Level,
			MA:                  species.MA.Level,
			ML:                  species.ML.Level,
			GV:                  species.GV.Level,
			LS:                  species.LS.Level,
			BI:                  species.BI.Level,
			BankedEconUnits:     float64(species.EconUnits),
			FleetMaintenancePct: float64(species.Fleet.MaintenancePct),
		})
		stat := gameStats.Stats[len(gameStats.Stats)-1]

		// why the check on fleet maintenance cost?
		if stat.FleetMaintenancePct < 0 {
			stat.FleetMaintenancePct = 0
		} else if stat.FleetMaintenancePct > 10000 {
			stat.FleetMaintenancePct = 10000
		}

		/* Get stats for namplas. */
		for _, nampla := range species.NamedPlanets.ById {
			if nampla.Planet.Orbit == 99 {
				continue
			}
			planet, colony := nampla.Planet, nampla.Colony
			if planet == nil || colony == nil {
				//log.Printf("species %q %v %v\n", species.Name, planet == nil, colony == nil)
				continue
			}

			stat.Shipyards += float64(colony.Shipyards)

			var production_penalty float64
			if ls_needed := lifeSupportNeeded(species, planet); ls_needed == 0 {
				production_penalty = 0
			} else if species.LS.Level > 0 {
				production_penalty = float64((100 * ls_needed) / species.LS.Level)
			} else {
				production_penalty = 100
			}

			raw_material_units := float64((10 * species.MI.Level * colony.Mining.Base) / planet.MiningDifficulty)
			raw_material_units -= (production_penalty * raw_material_units) / 100
			raw_material_units = ((float64(planet.EconEfficiency) * raw_material_units) + 50) / 100

			production_capacity := float64((species.MA.Level * colony.Manufacturing.Base) / 10)
			production_capacity -= (production_penalty * production_capacity) / 100
			production_capacity = ((float64(planet.EconEfficiency) * production_capacity) + 50) / 100

			var n1 float64
			if colony.Is.MiningColony {
				n1 = float64(2*raw_material_units) / 3
			} else if colony.Is.ResortColony {
				n1 = float64(2*production_capacity) / 3
			} else if production_capacity > raw_material_units {
				n1 = raw_material_units
			} else {
				n1 = production_capacity
			}
			n2 := ((stat.FleetMaintenancePct * n1) + 5000) / 10000
			stat.Production += n1 - n2

			if item, ok := colony.Inventory["PD"]; ok && item.Quantity > 0 {
				effectiveTonnage := item.Quantity / 200
				if effectiveTonnage < 1 {
					effectiveTonnage = 1
				}
				stat.DefensivePower += float64(power(effectiveTonnage))
			}

			if colony.Is.Populated {
				stat.PopulatedPlanets++
			}
		}

		/* Get stats for ships. */
		for _, ship := range species.Fleet.Ships {
			if ship.Location == nil || ship.Location.Orbit == 99 || ship.Status.UnderConstruction {
				continue
			}
			stat.Ships++
			trueTonnage := 10_000 * float64(ship.Class.Tonnage)
			if ship.Class.Is.Starbase {
				stat.Starbases++
				stat.StarbaseTonnage += trueTonnage
				stat.DefensivePower += float64(power(ship.Class.Tonnage))
			} else if ship.Class.Is.Transport {
				stat.Transports++
				stat.TransportTonnage += trueTonnage
			} else if ship.Class.Is.Warship {
				stat.Warships++
				stat.WarshipTonnage += trueTonnage
				if ship.Class.Is.SubLight {
					stat.DefensivePower += float64(power(ship.Class.Tonnage))
				} else {
					stat.OffensivePower += float64(power(ship.Class.Tonnage))
				}
			}
		}

		if species.ML.Level == 0 {
			stat.DefensivePower, stat.OffensivePower = 0, 0
		} else {
			stat.OffensivePower += ((float64(species.ML.Level) * stat.OffensivePower) / 50)
			stat.OffensivePower /= 10
			stat.DefensivePower += ((float64(species.ML.Level) * stat.DefensivePower) / 50)
			stat.DefensivePower /= 10
		}
	}

	// compute averages, minimums, and maximums for all species
	for i, sp := range gameStats.Stats {
		gameStats.TotalSpecies++
		if i == 0 {
			gameStats.BankedEconUnits.Min, gameStats.BankedEconUnits.Max = sp.BankedEconUnits, sp.BankedEconUnits
			gameStats.FleetMaintenancePct.Min, gameStats.FleetMaintenancePct.Max = sp.FleetMaintenancePct, sp.FleetMaintenancePct
			gameStats.PopulatedPlanets.Min, gameStats.PopulatedPlanets.Max = sp.PopulatedPlanets, sp.PopulatedPlanets
			gameStats.Production.Min, gameStats.Production.Max = sp.Production, sp.Production
			gameStats.Shipyards.Min, gameStats.Shipyards.Max = sp.Shipyards, sp.Shipyards
			gameStats.Starbases.Min, gameStats.Starbases.Max = sp.Starbases, sp.Starbases
			gameStats.StarbaseTonnage.Min, gameStats.StarbaseTonnage.Max = sp.StarbaseTonnage, sp.StarbaseTonnage
			gameStats.Transports.Min, gameStats.Transports.Max = sp.Transports, sp.Transports
			gameStats.TransportTonnage.Min, gameStats.TransportTonnage.Max = sp.TransportTonnage, sp.TransportTonnage
			gameStats.Warships.Min, gameStats.Warships.Max = sp.Warships, sp.Warships
			gameStats.WarshipTonnage.Min, gameStats.WarshipTonnage.Max = sp.WarshipTonnage, sp.WarshipTonnage

			gameStats.MI.Min, gameStats.MI.Max = sp.MI, sp.MI
			gameStats.MA.Min, gameStats.MA.Max = sp.MA, sp.MA
			gameStats.ML.Min, gameStats.ML.Max = sp.ML, sp.ML
			gameStats.GV.Min, gameStats.GV.Max = sp.GV, sp.GV
			gameStats.LS.Min, gameStats.LS.Max = sp.LS, sp.LS
			gameStats.BI.Min, gameStats.BI.Max = sp.BI, sp.BI
		}
		gameStats.MI.Total += sp.MI
		if sp.MI < gameStats.MI.Min {
			gameStats.MI.Min = sp.MI
		} else if gameStats.MI.Max < sp.MI {
			gameStats.MI.Max = sp.MI
		}
		gameStats.MA.Total += sp.MA
		if sp.MA < gameStats.MA.Min {
			gameStats.MA.Min = sp.MA
		} else if gameStats.MA.Max < sp.MA {
			gameStats.MA.Max = sp.MA
		}
		gameStats.ML.Total += sp.ML
		if sp.ML < gameStats.ML.Min {
			gameStats.ML.Min = sp.ML
		} else if gameStats.ML.Max < sp.ML {
			gameStats.ML.Max = sp.ML
		}
		gameStats.GV.Total += sp.GV
		if sp.GV < gameStats.GV.Min {
			gameStats.GV.Min = sp.GV
		} else if gameStats.GV.Max < sp.GV {
			gameStats.GV.Max = sp.GV
		}
		gameStats.LS.Total += sp.LS
		if sp.LS < gameStats.LS.Min {
			gameStats.LS.Min = sp.LS
		} else if gameStats.LS.Max < sp.LS {
			gameStats.LS.Max = sp.LS
		}
		gameStats.BI.Total += sp.BI
		if sp.BI < gameStats.BI.Min {
			gameStats.BI.Min = sp.BI
		} else if gameStats.BI.Max < sp.BI {
			gameStats.BI.Max = sp.BI
		}

		gameStats.BankedEconUnits.Total += float64(sp.BankedEconUnits)
		if sp.BankedEconUnits < gameStats.BankedEconUnits.Min {
			gameStats.BankedEconUnits.Min = sp.BankedEconUnits
		} else if gameStats.BankedEconUnits.Max < sp.BankedEconUnits {
			gameStats.BankedEconUnits.Max = sp.BankedEconUnits
		}
		gameStats.PopulatedPlanets.Total += float64(sp.PopulatedPlanets)
		if sp.PopulatedPlanets < gameStats.PopulatedPlanets.Min {
			gameStats.PopulatedPlanets.Min = sp.PopulatedPlanets
		} else if gameStats.PopulatedPlanets.Max < sp.PopulatedPlanets {
			gameStats.PopulatedPlanets.Max = sp.PopulatedPlanets
		}
		gameStats.Production.Total += float64(sp.Production)
		if sp.Production < gameStats.Production.Min {
			gameStats.Production.Min = sp.Production
		} else if gameStats.Production.Max < sp.Production {
			gameStats.Production.Max = sp.Production
		}
		gameStats.Shipyards.Total += float64(sp.Shipyards)
		if sp.Shipyards < gameStats.Shipyards.Min {
			gameStats.Shipyards.Min = sp.Shipyards
		} else if gameStats.Shipyards.Max < sp.Shipyards {
			gameStats.Shipyards.Max = sp.Shipyards
		}
		gameStats.Starbases.Total += float64(sp.Starbases)
		if sp.Starbases < gameStats.Starbases.Min {
			gameStats.Starbases.Min = sp.Starbases
		} else if gameStats.Starbases.Max < sp.Starbases {
			gameStats.Starbases.Max = sp.Starbases
		}
		gameStats.StarbaseTonnage.Total += float64(sp.StarbaseTonnage)
		if sp.StarbaseTonnage < gameStats.StarbaseTonnage.Min {
			gameStats.StarbaseTonnage.Min = sp.StarbaseTonnage
		} else if gameStats.StarbaseTonnage.Max < sp.StarbaseTonnage {
			gameStats.StarbaseTonnage.Max = sp.StarbaseTonnage
		}
		gameStats.Transports.Total += float64(sp.Transports)
		if sp.Transports < gameStats.Transports.Min {
			gameStats.Transports.Min = sp.Transports
		} else if gameStats.Transports.Max < sp.Transports {
			gameStats.Transports.Max = sp.Transports
		}
		gameStats.TransportTonnage.Total += sp.TransportTonnage
		if sp.TransportTonnage < gameStats.TransportTonnage.Min {
			gameStats.TransportTonnage.Min = sp.TransportTonnage
		} else if gameStats.TransportTonnage.Max < sp.TransportTonnage {
			gameStats.TransportTonnage.Max = sp.TransportTonnage
		}
		gameStats.Warships.Total += float64(sp.Warships)
		if sp.Warships < gameStats.Warships.Min {
			gameStats.Warships.Min = sp.Warships
		} else if gameStats.Warships.Max < sp.Warships {
			gameStats.Warships.Max = sp.Warships
		}
		gameStats.WarshipTonnage.Total += float64(sp.WarshipTonnage)
		if sp.WarshipTonnage < gameStats.WarshipTonnage.Min {
			gameStats.WarshipTonnage.Min = sp.WarshipTonnage
		} else if gameStats.WarshipTonnage.Max < sp.WarshipTonnage {
			gameStats.WarshipTonnage.Max = sp.WarshipTonnage
		}
	}

	if gameStats.TotalSpecies == 0 {
		// don't calculate values if there are no players?
		return gameStats
	}

	// calculate averages base on number of players
	gameStats.BankedEconUnits.Average = gameStats.BankedEconUnits.Total / float64(gameStats.TotalSpecies)
	gameStats.PopulatedPlanets.Average = gameStats.PopulatedPlanets.Total / float64(gameStats.TotalSpecies)
	gameStats.Production.Average = gameStats.Production.Total / float64(gameStats.TotalSpecies)
	gameStats.Shipyards.Average = gameStats.Shipyards.Total / float64(gameStats.TotalSpecies)
	gameStats.Starbases.Average = gameStats.Starbases.Total / float64(gameStats.TotalSpecies)
	gameStats.StarbaseTonnage.Average = gameStats.StarbaseTonnage.Total / float64(gameStats.TotalSpecies)
	gameStats.Transports.Average = gameStats.Transports.Total / float64(gameStats.TotalSpecies)
	gameStats.TransportTonnage.Average = gameStats.TransportTonnage.Total / float64(gameStats.TotalSpecies)
	gameStats.Warships.Average = gameStats.Warships.Total / float64(gameStats.TotalSpecies)
	gameStats.WarshipTonnage.Average = gameStats.WarshipTonnage.Total / float64(gameStats.TotalSpecies)

	return gameStats
}
