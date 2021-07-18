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
	"github.com/mdhender/fhcms/store/jsondb"
)

var __jdb *jsondb.Store
var galaxy galaxy_data
var planet_base []*planet_data // warning: code assumes *planet_data

//*************************************************************************
// get_galaxy_data: get_gal.c
func get_galaxy_data() {
	galaxy.d_num_species = __jdb.Galaxy.DNumSpecies
	galaxy.num_species = __jdb.Galaxy.NumSpecies
	galaxy.radius = __jdb.Galaxy.Radius
	galaxy.turn_number = __jdb.Galaxy.TurnNumber
}

//*************************************************************************
// get_planet_data: get_plan.c
func get_planet_data() {
	/* Allocate enough memory for all planets. */
	num_planets = len(__jdb.Planets)
	planet_base = make([]*planet_data, num_planets, num_planets) // warning: was num_planets + NUM_EXTRA_PLANETS
	for id, planet := range __jdb.Planets {
		p := &planet_data{
			diameter:          planet.Diameter,
			econ_efficiency:   planet.EconEfficiency,
			gas:               planet.Gas,
			gas_percent:       planet.GasPercent,
			gravity:           planet.Gravity,
			md_increase:       planet.MdIncrease,
			message:           planet.Message,
			mining_difficulty: planet.MiningDifficulty,
			pressure_class:    planet.PressureClass,
			special:           planet.Special,
			temperature_class: planet.TemperatureClass,
		}
		planet_base[id] = p
	}
	planet_data_modified = false
}

//*************************************************************************
// get_species_data: utils.c
// get_species_data will read in data files for all species.
// Additional memory must be allocated for routines that build ships or
// name planets.
func get_species_data() {
	/* Allocate enough memory for all species. */ // uhhh, spec_data is a constant size, no need to allocate
	num_species = len(__jdb.Species)
	for id, species := range __jdb.Species {
		s := &species_data{
			id:                 species.Id,
			x:                  species.X,
			y:                  species.Y,
			z:                  species.Z,
			auto_orders:        species.AutoOrders != 0,
			econ_units:         species.EconUnits,
			fleet_cost:         species.FleetCost,
			fleet_percent_cost: species.FleetPercentCost,
			govt_name:          species.GovtName,
			govt_type:          species.GovtType,
			hp_original_base:   species.HpOriginalBase,
			init_tech_level:    species.InitTechLevel,
			name:               species.Name,
			namplas:            make([]*nampla_data, species.NumNamplas, species.NumNamplas), /* Allocate enough memory for all namplas. */
			neutral_gas:        species.NeutralGas,
			num_namplas:        species.NumNamplas,
			num_ships:          species.NumShips,
			poison_gas:         species.PoisonGas,
			required_gas:       species.RequiredGas,
			required_gas_max:   species.RequiredGasMax,
			required_gas_min:   species.RequiredGasMin,
			ships:              make([]*ship_data_, species.NumShips, species.NumShips), /* Allocate enough memory for all ships. */
			tech_eps:           species.TechEps,
			tech_knowledge:     species.TechKnowledge,
			tech_level:         species.TechLevel,
		}
		for _, id := range species.Ally {
			s.ally[id] = true
		}
		for _, id := range species.Contact {
			s.contact[id] = true
		}
		for _, id := range species.Enemy {
			s.enemy[id] = true
		}
		namp_data[id] = s.namplas
		for k := range species.Namplas {
			n := nampla_data{
				AUs_needed:     species.Namplas[k].AUsNeeded,
				AUs_to_install: species.Namplas[k].AUsToInstall,
				auto_AUs:       species.Namplas[k].AutoAUs,
				auto_IUs:       species.Namplas[k].AutoIUs,
				hidden:         species.Namplas[k].Hidden != 0,
				hiding:         species.Namplas[k].Hiding != 0,
				IUs_needed:     species.Namplas[k].IUsNeeded,
				IUs_to_install: species.Namplas[k].IUsToInstall,
				ma_base:        species.Namplas[k].MaBase,
				message:        species.Namplas[k].Message,
				mi_base:        species.Namplas[k].MiBase,
				name:           species.Namplas[k].Name,
				planet_index:   species.Namplas[k].PlanetIndex,
				pn:             species.Namplas[k].Pn,
				pop_units:      species.Namplas[k].PopUnits,
				shipyards:      species.Namplas[k].Shipyards,
				siege_eff:      species.Namplas[k].SiegeEff,
				special:        species.Namplas[k].Special,
				status:         species.Namplas[k].Status,
				use_on_ambush:  species.Namplas[k].UseOnAmbush,
				x:              species.Namplas[k].X,
				y:              species.Namplas[k].Y,
				z:              species.Namplas[k].Z,
			}
			for item, qty := range species.Namplas[k].ItemQuantity {
				n.item_quantity[item] = qty
			}
			namp_data[id][k] = &n
		}
		ship_data[species_index] = s.ships
		for k := range species.Ships {
			s := ship_data_{
				age:                  species.Ships[k].Age,
				arrived_via_wormhole: species.Ships[k].ArrivedViaWormhole != 0,
				class:                species.Ships[k].Class,
				name:                 species.Ships[k].Name,
				dest_x:               species.Ships[k].DestX,
				dest_y:               species.Ships[k].DestY,
				dest_z:               species.Ships[k].DestZ,
				just_jumped:          species.Ships[k].JustJumped != 0,
				loading_point:        species.Ships[k].LoadingPoint,
				pn:                   species.Ships[k].Pn,
				remaining_cost:       species.Ships[k].RemainingCost,
				special:              species.Ships[k].Special,
				status:               species.Ships[k].Status,
				tonnage:              species.Ships[k].Tonnage,
				ttype:                species.Ships[k].Type,
				unloading_point:      species.Ships[k].UnloadingPoint,
				x:                    species.Ships[k].X,
				y:                    species.Ships[k].Y,
				z:                    species.Ships[k].Z,
			}
			for item, qty := range species.Ships[k].ItemQuantity {
				s.item_quantity[item] = qty
			}
			ship_data[id][k] = &s
		}
		spec_data[id] = s
	}
	data_in_memory[species_index] = true
	num_new_namplas[species_index] = 0
	num_new_ships[species_index] = 0
}

//*************************************************************************
// get_star_data: get_star.c
func get_star_data() {
	/* Allocate enough memory for all stars. */
	num_stars = len(__jdb.Stars)
	star_base = make([]*star_data, num_stars, num_stars) // warning: was num_stars + NUM_EXTRA_STARS
	for id, star := range __jdb.Stars {
		s := &star_data{
			color:        star.Color,
			home_system:  star.HomeSystem != 0,
			message:      star.Message,
			num_planets:  star.NumPlanets,
			planet_index: star.PlanetIndex,
			size:         star.Size,
			ttype:        star.Type,
			worm_here:    star.WormHere != 0,
			worm_x:       star.WormX,
			worm_y:       star.WormY,
			worm_z:       star.WormZ,
			x:            star.X,
			y:            star.Y,
			z:            star.Z,
		}
		star_base[id] = s
	}
	star_data_modified = false
}
