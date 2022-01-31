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

// regenerate_shields restores by 5 + LS/10 percent of original shield strength per round.
func (e *Engine) regenerate_shields(act *action_data) {
	for unit_index := 0; unit_index < act.num_units_fighting; unit_index++ {
		species_index := act.fighting_species_index[unit_index]

		// percent is the amount regenerated per round
		percent := (e.c_species[species_index].tech_level[LS] / 10) + 5

		// max strength is the original strength, which we can't go above
		max_shield_strength := act.shield_strength[unit_index]

		act.shield_strength_left[unit_index] += (percent * max_shield_strength) / 100
		if act.shield_strength_left[unit_index] > max_shield_strength {
			act.shield_strength_left[unit_index] = max_shield_strength
		}
	}
}
