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

package dat32

// NamedPlanet is the home planet, a colony, or just a named planet.
// TODO: how to tell the difference?
type NamedPlanet struct {
	/* Name of planet. */
	Name string `json:"name"`
	/* Coordinates. */
	X  int `json:"x"`
	Y  int `json:"y"`
	Z  int `json:"z"`
	PN int `json:"pn"`
	/* Status of planet. */
	Status int `json:"status"`
	/* HIDE order given. */
	Hiding bool `json:"hiding"`
	/* Colony is hidden. */
	Hidden bool `json:"hidden"`
	/* Index (starting at zero) into the file "planets.dat" of this planet. */
	PlanetIndex int `json:"planet_index"`
	/* Siege effectiveness - a percentage between 0 and 99. */
	SiegeEff int `json:"siege_eff"`
	/* Number of shipyards on planet. */
	Shipyards int `json:"shipyards"`
	/* Incoming ship with only CUs on board. */
	IUsNeeded int `json:"ius_needed"`
	/* Incoming ship with only CUs on board. */
	AUsNeeded int `json:"aus_needed"`
	/* Number of IUs to be automatically installed. */
	AutoIUs int `json:"auto_ius"`
	/* Number of AUs to be automatically installed. */
	AutoAUs int `json:"auto_aus"`
	/* Colonial mining units to be installed. */
	IUsToInstall int `json:"ius_to_install"`
	/* Colonial manufacturing units to be installed. */
	AUsToInstall int `json:"aus_to_install"`
	/* Mining base times 10. */
	MiBase int `json:"mi_base"`
	/* Manufacturing base times 10. */
	MaBase int `json:"ma_base"`
	/* Number of available population units. */
	PopUnits int `json:"pop_units"`
	/* Quantity of each item available. */
	ItemQuantity [MAX_ITEMS]int `json:"item_quantity"`
	/* Amount to use on ambush. */
	UseOnAmbush int `json:"use_on_ambush"`
	/* Message associated with this planet, if any. */
	Message int `json:"message"`
	/* Different for each application. */
	Special int `json:"special"`
}
