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

// Ship is the ship data from the binary data file
type Ship struct {
	/* Name of ship. In original case. */
	Name string `json:"name"`
	/* Current coordinates. */
	X  int `json:"x"`
	Y  int `json:"y"`
	Z  int `json:"z"`
	PN int `json:"pn"`
	/* Current status of ship. */
	Status int `json:"status"`
	/* Ship type. */
	Type int `json:"type"`
	/* Destination if ship was forced to jump from combat. Also used by TELESCOPE command. */
	DestX int `json:"dest_x"`
	DestY int `json:"dest_y"`
	DestZ int `json:"dest_z"`
	/* Set if ship jumped this turn. */
	JustJumped bool `json:"just_jumped"`
	/* Ship arrived via wormhole in the PREVIOUS turn. */
	ArrivedViaWormhole bool `json:"arrived_via_wormhole"`
	/* Ship class. */
	Class int `json:"class"`
	/* Ship tonnage divided by 10,000. */
	Tonnage int `json:"tonnage"`
	/* Quantity of each item carried. */
	ItemQuantity [MAX_ITEMS]int `json:"item_quantity"`
	/* Ship age. */
	Age int `json:"age"`
	/* The cost needed to complete the ship if still under construction. */
	RemainingCost int `json:"remaining_cost"`
	/* NamedPlanet index for planet where ship was last loaded with CUs. Zero = none. Use 9999 for home planet. */
	LoadingPoint int `json:"loading_point"`
	/* NamedPlanet index for planet that ship should be given orders to jump to where it will unload. Zero = none. Use 9999 for home planet. */
	UnloadingPoint int `json:"unloading_point"`
	/* Different for each application. */
	Special int `json:"special"`
}
