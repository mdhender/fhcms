/*******************************************************************************
Far Horizons Engine
Copyright (C) 2022  Michael D Henderson

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

package engine

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/mdhender/fhcms/internal/dat32"
	"log"
	"path/filepath"
)

// LoadBinary load all data from the original binary files
func (e *Engine) LoadBinary(root string, endian binary.ByteOrder) error {
	galaxy, err := dat32.ReadGalaxy(filepath.Join(root, "galaxy.dat"), endian)
	if err != nil {
		return err
	}
	e.galaxy.d_num_species = galaxy.DNumSpecies
	e.galaxy.num_species = galaxy.NumSpecies
	e.galaxy.radius = galaxy.Radius
	e.galaxy.turn_number = galaxy.TurnNumber
	log.Printf("[engine] loadBinary: loaded galaxy turn %6d\n", e.galaxy.turn_number)

	stars, err := dat32.ReadStars(filepath.Join(root, "stars.dat"), endian)
	if err != nil {
		return err
	}
	e.star_base = make([]*star_data, len(stars.Stars), len(stars.Stars))
	for i := 0; i < len(stars.Stars); i++ {
		sd := &star_data{
			x:            stars.Stars[i].X,
			y:            stars.Stars[i].Y,
			z:            stars.Stars[i].Z,
			_type:        stars.Stars[i].Type,
			color:        stars.Stars[i].Color,
			size:         stars.Stars[i].Size,
			num_planets:  stars.Stars[i].NumPlanets,
			home_system:  stars.Stars[i].HomeSystem,
			worm_here:    stars.Stars[i].WormHere,
			worm_x:       stars.Stars[i].WormX,
			worm_y:       stars.Stars[i].WormY,
			worm_z:       stars.Stars[i].WormZ,
			planet_index: stars.Stars[i].PlanetIndex,
			message:      stars.Stars[i].Message,
			visited_by:   make([]int, e.galaxy.num_species+1, e.galaxy.num_species+1),
		}
		for v := 0; v < len(stars.Stars[i].VisitedBy); v++ {
			if stars.Stars[i].VisitedBy[v] <= e.galaxy.num_species {
				sd.visited_by[stars.Stars[i].VisitedBy[v]-1] = TRUE // visited by is zero-based index?
			}
		}
		e.star_base[i] = sd
	}
	e.num_stars = len(stars.Stars)
	log.Printf("[engine] loadBinary: loaded %6d stars\n", e.num_stars)

	planets, err := dat32.ReadPlanets(filepath.Join(root, "planets.dat"), endian)
	if err != nil {
		return err
	}
	e.planet_base = make([]*planet_data, len(planets.Planets), len(planets.Planets))
	for i := 0; i < len(planets.Planets); i++ {
		e.planet_base[i] = &planet_data{
			temperature_class: planets.Planets[i].TemperatureClass,
			pressure_class:    planets.Planets[i].PressureClass,
			special:           planets.Planets[i].Special,
			gas:               planets.Planets[i].Gas,
			gas_percent:       planets.Planets[i].GasPercent,
			diameter:          planets.Planets[i].Diameter,
			gravity:           planets.Planets[i].Gravity,
			mining_difficulty: planets.Planets[i].MiningDifficulty,
			econ_efficiency:   planets.Planets[i].EconEfficiency,
			md_increase:       planets.Planets[i].MDIncrease,
			message:           planets.Planets[i].Message,
		}
	}
	log.Printf("[engine] loadBinary: loaded %6d planets\n", len(planets.Planets))

	e.spec_data = make([]*species_data, e.galaxy.num_species, e.galaxy.num_species)
	e.spec_logs = make([]*bytes.Buffer, e.galaxy.num_species, e.galaxy.num_species)
	e.spec_orders = make([][]byte, e.galaxy.num_species, e.galaxy.num_species)
	e.namp_data = make([][]*nampla_data, e.galaxy.num_species, e.galaxy.num_species)
	e.ship_data = make([][]*ship_data, e.galaxy.num_species, e.galaxy.num_species)
	for i := 0; i < galaxy.NumSpecies; i++ {
		e.spec_logs[i] = &bytes.Buffer{}
		sp, err := dat32.ReadSpecies(filepath.Join(root, fmt.Sprintf("sp%02d.dat", i+1)), i+1, endian)
		if err != nil {
			return err
		}
		sd := &species_data{
			name:               sp.Name,
			govt_name:          sp.GovtName,
			govt_type:          sp.GovtType,
			x:                  sp.X,
			y:                  sp.Y,
			z:                  sp.Z,
			pn:                 sp.PN,
			required_gas:       sp.RequiredGas,
			required_gas_min:   sp.RequiredGasMin,
			required_gas_max:   sp.RequiredGasMax,
			tech_level:         sp.TechLevel,
			init_tech_level:    sp.InitTechLevel,
			tech_knowledge:     sp.TechKnowledge,
			num_namplas:        sp.NumNamplas,
			num_ships:          sp.NumShips,
			tech_eps:           sp.TechEps,
			hp_original_base:   sp.HPOriginalBase,
			econ_units:         sp.EconUnits,
			fleet_cost:         sp.FleetCost,
			fleet_percent_cost: sp.FleetPercentCost,
			contact:            make([]int, e.galaxy.num_species+1, e.galaxy.num_species+1),
			ally:               make([]int, e.galaxy.num_species+1, e.galaxy.num_species+1),
			enemy:              make([]int, e.galaxy.num_species+1, e.galaxy.num_species+1),
		}
		if sp.AutoOrders {
			sd.auto_orders = TRUE
		}
		for j := 0; j < len(sp.NeutralGas); j++ {
			sd.neutral_gas[j] = sp.NeutralGas[j]
		}
		for j := 0; j < len(sp.PoisonGas); j++ {
			sd.poison_gas[j] = sp.PoisonGas[j]
		}
		for v := 0; v < len(sp.Contact); v++ {
			if sp.Contact[v] <= e.galaxy.num_species {
				sd.contact[sp.Contact[v]] = TRUE
			}
		}
		for v := 0; v < len(sp.Ally); v++ {
			if sp.Ally[v] <= e.galaxy.num_species {
				sd.ally[sp.Ally[v]] = TRUE
			}
		}
		for v := 0; v < len(sp.Enemy); v++ {
			if sp.Enemy[v] <= e.galaxy.num_species {
				sd.enemy[sp.Enemy[v]] = TRUE
			}
		}
		for j := 0; j < sp.NumNamplas; j++ {
			np := &nampla_data{
				name:           sp.NamplaBase[j].Name,
				x:              sp.NamplaBase[j].X,
				y:              sp.NamplaBase[j].Y,
				z:              sp.NamplaBase[j].Z,
				pn:             sp.NamplaBase[j].PN,
				status:         sp.NamplaBase[j].Status,
				planet_index:   sp.NamplaBase[j].PlanetIndex,
				siege_eff:      sp.NamplaBase[j].SiegeEff,
				shipyards:      sp.NamplaBase[j].Shipyards,
				IUs_needed:     sp.NamplaBase[j].IUsNeeded,
				AUs_needed:     sp.NamplaBase[j].AUsNeeded,
				auto_IUs:       sp.NamplaBase[j].AutoIUs,
				auto_AUs:       sp.NamplaBase[j].AutoAUs,
				IUs_to_install: sp.NamplaBase[j].IUsToInstall,
				AUs_to_install: sp.NamplaBase[j].AUsToInstall,
				mi_base:        sp.NamplaBase[j].MiBase,
				ma_base:        sp.NamplaBase[j].MaBase,
				pop_units:      sp.NamplaBase[j].PopUnits,
				item_quantity:  sp.NamplaBase[j].ItemQuantity,
				use_on_ambush:  sp.NamplaBase[j].UseOnAmbush,
				message:        sp.NamplaBase[j].Message,
				special:        sp.NamplaBase[j].Special,
			}
			if sp.NamplaBase[j].Hiding {
				np.hiding = TRUE
			}
			if sp.NamplaBase[j].Hidden {
				np.hidden = TRUE
			}
			e.namp_data[i] = append(e.namp_data[i], np)
		}
		for j := 0; j < sp.NumShips; j++ {
			sh := &ship_data{
				name:            sp.ShipBase[j].Name,
				x:               sp.ShipBase[j].X,
				y:               sp.ShipBase[j].Y,
				z:               sp.ShipBase[j].Z,
				pn:              sp.ShipBase[j].PN,
				status:          sp.ShipBase[j].Status,
				_type:           sp.ShipBase[j].Type,
				dest_x:          sp.ShipBase[j].DestX,
				dest_y:          sp.ShipBase[j].DestY,
				dest_z:          sp.ShipBase[j].DestZ,
				class:           sp.ShipBase[j].Class,
				tonnage:         sp.ShipBase[j].Tonnage,
				item_quantity:   sp.ShipBase[j].ItemQuantity,
				age:             sp.ShipBase[j].Age,
				remaining_cost:  sp.ShipBase[j].RemainingCost,
				loading_point:   sp.ShipBase[j].LoadingPoint,
				unloading_point: sp.ShipBase[j].UnloadingPoint,
				special:         sp.ShipBase[j].Special,
			}
			if sp.ShipBase[j].JustJumped {
				sh.just_jumped = TRUE
			}
			if sp.ShipBase[j].ArrivedViaWormhole {
				sh.arrived_via_wormhole = TRUE
			}
			e.ship_data[i] = append(e.ship_data[i], sh)
		}
		e.spec_data[i] = sd
	}
	log.Printf("[engine] loadBinary: loaded %6d species\n", len(e.spec_data))

	//// create new locations array
	//e.locations = nil
	//// add all colonies
	//for _, sp := range species {
	//	for _, np := range sp.NamplaBase {
	//		if np.Colony != nil && np.Colony.Is.Populated {
	//			e.addLocation(sp, np.Planet.Location)
	//		}
	//	}
	//}
	//// add all ships
	//for _, sp := range species {
	//	for _, ship := range sp.ShipBase {
	//		if ship.Location != nil && ship.Status != nil && !(ship.Status.ForcedJump || ship.Status.JumpedInCombat) {
	//			e.addLocation(sp, ship.Location)
	//		}
	//	}
	//}

	return nil
}

// LoadJson will load all data from JSON files
func LoadJson() error {
	return nil
}

func (e *Engine) addLocation(sp *Species, c *Coords) {
	//for _, l := range e.locations {
	//	if l.Species.Id == sp.Id && l.Location.X == c.X && l.Location.Y == c.Y && l.Location.Z == c.Z {
	//		return // already present
	//	}
	//}
	//e.locations = append(e.locations, &SpeciesLocationData{Species: sp, Location: &Coords{X: c.X, Y: c.Y, Z: c.Z}})
}
