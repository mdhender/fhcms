/*****************************************************************************
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
****************************************************************************/

package jsondb

type Store struct {
	Version string `json:"semver"`
	Galaxy  struct {
		DNumSpecies int `json:"d_num_species"`
		NumSpecies  int `json:"num_species"`
		Radius      int `json:"radius"`
		TurnNumber  int `json:"turn_number"`
	} `json:"galaxy"`
	Locations []Location `json:"locations"`
	Planets []struct {
		Id               int    `json:"id"`
		Diameter         int    `json:"diameter"`
		EconEfficiency   int    `json:"econ_efficiency"`
		Gas              [4]int `json:"gas"`
		GasPercent       [4]int `json:"gas_percent"`
		Gravity          int    `json:"gravity"`
		MdIncrease       int    `json:"md_increase"`
		Message          int    `json:"message"`
		MiningDifficulty int    `json:"mining_difficulty"`
		PressureClass    int    `json:"pressure_class"`
		Special          int    `json:"special"`
		TemperatureClass int    `json:"temperature_class"`
	} `json:"planets"`
	Species []struct {
		Id               int    `json:"id"`
		Ally             []int  `json:"ally"`
		AutoOrders       int    `json:"auto_orders"`
		Contact          []int  `json:"contact"`
		EconUnits        int    `json:"econ_units"`
		Enemy            []int  `json:"enemy"`
		FleetCost        int    `json:"fleet_cost"`
		FleetPercentCost int    `json:"fleet_percent_cost"`
		GovtName         string `json:"govt_name"`
		GovtType         string `json:"govt_type"`
		HpOriginalBase   int    `json:"hp_original_base"`
		InitTechLevel    [6]int `json:"init_tech_level"`
		Name             string `json:"name"`
		Namplas          []struct {
			Id           int    `json:"id"`
			AUsNeeded    int    `json:"aus_needed"`
			AUsToInstall int    `json:"aus_to_install"`
			AutoAUs      int    `json:"auto_aus"`
			AutoIUs      int    `json:"auto_ius"`
			Hidden       int    `json:"hidden"`
			Hiding       int    `json:"hiding"`
			ItemQuantity []int  `json:"item_quantity"`
			IUsNeeded    int    `json:"ius_needed"`
			IUsToInstall int    `json:"ius_to_install"`
			Name         string `json:"name"`
			PlanetIndex  int    `json:"planet_index"`
			Pn           int    `json:"pn"`
			PopUnits     int    `json:"pop_units"`
			MaBase       int    `json:"ma_base"`
			Message      int    `json:"message"`
			MiBase       int    `json:"mi_base"`
			Shipyards    int    `json:"shipyards"`
			SiegeEff     int    `json:"siege_eff"`
			Status       int    `json:"status"`
			Special      int    `json:"special"`
			UseOnAmbush  int    `json:"use_on_ambush"`
			X            int    `json:"x"`
			Y            int    `json:"y"`
			Z            int    `json:"z"`
		} `json:"namplas"`
		NeutralGas     [6]int `json:"neutral_gas"`
		NumNamplas     int    `json:"num_namplas"`
		NumShips       int    `json:"num_ships"`
		Pn             int    `json:"pn"`
		PoisonGas      [6]int `json:"poison_gas"`
		RequiredGas    int    `json:"required_gas"`
		RequiredGasMax int    `json:"required_gas_max"`
		RequiredGasMin int    `json:"required_gas_min"`
		Ships          []struct {
			Id                 int    `json:"id"`
			Age                int    `json:"age"`
			ArrivedViaWormhole int    `json:"arrived_via_wormhole"`
			Class              int    `json:"class"`
			DestX              int    `json:"dest_x"`
			DestY              int    `json:"dest_y"`
			DestZ              int    `json:"dest_z"`
			JustJumped         int    `json:"just_jumped"`
			ItemQuantity       []int  `json:"item_quantity"`
			LoadingPoint       int    `json:"loading_point"`
			Name               string `json:"name"`
			Pn                 int    `json:"pn"`
			RemainingCost      int    `json:"remaining_cost"`
			Special            int    `json:"special"`
			Status             int    `json:"status"`
			Tonnage            int    `json:"tonnage"`
			Type               int    `json:"type"`
			UnloadingPoint     int    `json:"unloading_point"`
			X                  int    `json:"x"`
			Y                  int    `json:"y"`
			Z                  int    `json:"z"`
		} `json:"ships"`
		TechEps       [6]int `json:"tech_eps"`
		TechKnowledge [6]int `json:"tech_knowledge"`
		TechLevel     [6]int `json:"tech_level"`
		X             int    `json:"x"`
		Y             int    `json:"y"`
		Z             int    `json:"z"`
	} `json:"species"`
	Stars []struct {
		Id          int   `json:"id"`
		Color       int   `json:"color"`
		HomeSystem  int   `json:"home_system"`
		Message     int   `json:"message"`
		NumPlanets  int   `json:"num_planets"`
		PlanetIndex int   `json:"planet_index"`
		Size        int   `json:"size"`
		Type        int   `json:"type"`
		VisitedBy   []int `json:"visited_by"`
		WormHere    int   `json:"worm_here"`
		WormX       int   `json:"worm_x"`
		WormY       int   `json:"worm_y"`
		WormZ       int   `json:"worm_z"`
		X           int   `json:"x"`
		Y           int   `json:"y"`
		Z           int   `json:"z"`
	} `json:"stars"`
}

type Location struct	{
	S int `json:"s"` // species number
	X int `json:"x"`
	Y int `json:"y"`
	Z int `json:"z"`
}

/* Status codes for named planets. These are logically ORed together. */
const HOME_PLANET = 1
const COLONY = 2
const POPULATED = 8
const MINING_COLONY = 16
const RESORT_COLONY = 32
const DISBANDED_COLONY = 64

/* Ship status codes. */
const UNDER_CONSTRUCTION = 0
const ON_SURFACE = 1
const IN_ORBIT = 2
const IN_DEEP_SPACE = 3
const JUMPED_IN_COMBAT = 4
const FORCED_JUMP = 5
